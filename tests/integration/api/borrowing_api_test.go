package api

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/util/errormap"
	"github.com/nitikhon/golang-inventory-system/tests/integration/api/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBorrowItem(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	// Helper to get tokens
	userToken := getAuthToken(t, server, "test_user", "P@ssw0rd")

	// Get existings items/users for valid IDs
	var availableItem entity.Item
	err := server.DB.Where("status = ?", "available").First(&availableItem).Error
	require.NoError(t, err)

	var user entity.User
	err = server.DB.Where("username = ?", "test_user").First(&user).Error
	require.NoError(t, err)

	t.Run("Authentication", func(t *testing.T) {
		payload := map[string]any{
			"user_id":          user.ID,
			"item_id":          availableItem.ID,
			"borrowing_amount": 1,
			"borrowed_at":      time.Now().Format(time.RFC3339),
			"due_date":         time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)
		req := createAuthenticatedRequest("POST", "/api/borrows/", body, "")
		resp, err := server.App.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Validation Errors", func(t *testing.T) {
		tests := []struct {
			name           string
			payload        map[string]any
			expectedStatus int
			expectedError  string
		}{
			{
				name: "Missing User ID",
				payload: map[string]any{
					"item_id":          availableItem.ID,
					"borrowing_amount": 1,
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "user ID is required",
			},
			{
				name: "Missing Item ID",
				payload: map[string]any{
					"user_id":          user.ID,
					"borrowing_amount": 1,
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "item ID is required",
			},
			{
				name: "Invalid Borrowing Amount (Zero)",
				payload: map[string]any{
					"user_id":          user.ID,
					"item_id":          availableItem.ID,
					"borrowing_amount": 0,
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "borrowing amount must be greater than zero",
			},
			{
				name: "ReturnedAt Should Be Empty",
				payload: map[string]any{
					"user_id":          user.ID,
					"item_id":          availableItem.ID,
					"borrowing_amount": 1,
					"returned_at":      time.Now().Format(time.RFC3339),
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "returned at date must be empty when borrowing an item",
			},
			{
				name: "DueDate Before BorrowedAt",
				payload: map[string]any{
					"user_id":          user.ID,
					"item_id":          availableItem.ID,
					"borrowing_amount": 1,
					"borrowed_at":      time.Now().Add(24 * time.Hour).Format(time.RFC3339),
					"due_date":         time.Now().Format(time.RFC3339),
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "due date cannot be before the borrowed at date",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, _ := json.Marshal(tt.payload)
				req := createAuthenticatedRequest("POST", "/api/borrows/", body, userToken)
				resp, err := server.App.Test(req)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				var respBody map[string]string
				json.NewDecoder(resp.Body).Decode(&respBody)
				if tt.expectedError != "" {
					assert.Contains(t, respBody["error"], tt.expectedError)
				}
			})
		}
	})

	t.Run("Business Logic Errors", func(t *testing.T) {
		t.Run("Item Not Available or Not Found", func(t *testing.T) {
			payload := map[string]any{
				"user_id":          user.ID,
				"item_id":          999999,
				"borrowing_amount": 1,
				"borrowed_at":      time.Now().Format(time.RFC3339),
			}
			body, _ := json.Marshal(payload)
			req := createAuthenticatedRequest("POST", "/api/borrows/", body, userToken)
			resp, err := server.App.Test(req)
			require.NoError(t, err)

			var respBody map[string]string
			json.NewDecoder(resp.Body).Decode(&respBody)
			t.Log(respBody)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("Not Enough Items", func(t *testing.T) {
			var limitedItem entity.Item
			err := server.DB.Where("available_amount > 0 AND available_amount < 100").First(&limitedItem).Error
			require.NoError(t, err)

			payload := map[string]any{
				"user_id":          user.ID,
				"item_id":          limitedItem.ID,
				"borrowing_amount": limitedItem.AvailableAmount + 1,
				"borrowed_at":      time.Now().Format(time.RFC3339),
			}
			body, _ := json.Marshal(payload)
			req := createAuthenticatedRequest("POST", "/api/borrows/", body, userToken)
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusConflict, resp.StatusCode)

			var respBody map[string]string
			json.NewDecoder(resp.Body).Decode(&respBody)
			assert.Equal(t, errormap.ErrItemNotEnough, respBody["error"])
		})
	})

	t.Run("Successful Borrowing", func(t *testing.T) {
		// Use a seeded item that has availability and no existing borrowings
		var testItem entity.Item
		err := server.DB.Where("name = ?", "external hard drive").First(&testItem).Error
		require.NoError(t, err)

		borrowedAt := time.Now().Truncate(time.Second).UTC()
		dueDate := borrowedAt.Add(7 * 24 * time.Hour)

		payload := map[string]any{
			"user_id":          user.ID,
			"item_id":          testItem.ID,
			"borrowing_amount": 1,
			"borrowed_at":      borrowedAt.Format(time.RFC3339),
			"due_date":         dueDate.Format(time.RFC3339),
			"description":      "Test borrowing",
		}
		body, _ := json.Marshal(payload)

		req := createAuthenticatedRequest("POST", "/api/borrows/", body, userToken)
		resp, err := server.App.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var borrowing entity.Borrowing
		err = json.NewDecoder(resp.Body).Decode(&borrowing)
		require.NoError(t, err)

		assert.Equal(t, user.ID, borrowing.UserID)
		assert.Equal(t, testItem.ID, borrowing.ItemID)
		assert.Equal(t, 1, borrowing.BorrowingAmount)
		assert.Equal(t, "Test borrowing", borrowing.Description)
		assert.Equal(t, entity.BORROWING_PENDING, borrowing.BorrowingStatus)
	})
}

func TestApproveBorrowing(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	// Helper to get tokens
	adminToken := getAuthToken(t, server, "test_admin", "P@ssw0rd")
	userToken := getAuthToken(t, server, "test_user", "P@ssw0rd")

	// Get a pending borrowing from seed data ("wireless keyboard" borrowing is pending)
	var pendingBorrowing entity.Borrowing
	err := server.DB.Where("borrowing_status = ? AND approval_status = ?", "pending", "pending").First(&pendingBorrowing).Error
	require.NoError(t, err)

	var adminUser entity.User
	err = server.DB.Where("username = ?", "test_admin").First(&adminUser).Error
	require.NoError(t, err)

	t.Run("Authentication and Authorization", func(t *testing.T) {
		payload := map[string]any{
			"id":          pendingBorrowing.ID,
			"approved_by": adminUser.ID,
		}
		body, _ := json.Marshal(payload)

		t.Run("No Token", func(t *testing.T) {
			req := createAuthenticatedRequest("POST", "/api/borrows/approve", body, "")
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("Regular User Token", func(t *testing.T) {
			req := createAuthenticatedRequest("POST", "/api/borrows/approve", body, userToken)
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})
	})

	t.Run("Validation Errors", func(t *testing.T) {
		tests := []struct {
			name           string
			payload        map[string]any
			expectedStatus int
			expectedError  string
		}{
			{
				name: "Missing Borrowing ID",
				payload: map[string]any{
					"approved_by": adminUser.ID,
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "borrowing ID is required",
			},
			{
				name: "Missing ApprovedBy",
				payload: map[string]any{
					"id": pendingBorrowing.ID,
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "approvedBy is required",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, _ := json.Marshal(tt.payload)
				req := createAuthenticatedRequest("POST", "/api/borrows/approve", body, adminToken)
				resp, err := server.App.Test(req)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				var respBody map[string]string
				json.NewDecoder(resp.Body).Decode(&respBody)
				if tt.expectedError != "" {
					assert.Contains(t, respBody["error"], tt.expectedError)
				}
			})
		}
	})

	t.Run("Business Logic Errors", func(t *testing.T) {
		t.Run("Borrowing Not Found", func(t *testing.T) {
			payload := map[string]any{
				"id":          999999,
				"approved_by": adminUser.ID,
			}
			body, _ := json.Marshal(payload)
			req := createAuthenticatedRequest("POST", "/api/borrows/approve", body, adminToken)
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("Borrowing Not Pending", func(t *testing.T) {
			// Get an already approved borrowing (First one in seed data is approved)
			var activeBorrowing entity.Borrowing
			err := server.DB.Where("borrowing_status = ?", "active").First(&activeBorrowing).Error
			require.NoError(t, err)

			payload := map[string]any{
				"id":          activeBorrowing.ID,
				"approved_by": adminUser.ID,
			}
			body, _ := json.Marshal(payload)
			req := createAuthenticatedRequest("POST", "/api/borrows/approve", body, adminToken)
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusConflict, resp.StatusCode)
		})
	})

	t.Run("Successful Approval", func(t *testing.T) {
		payload := map[string]any{
			"id":          pendingBorrowing.ID,
			"approved_by": adminUser.ID,
		}
		body, _ := json.Marshal(payload)

		req := createAuthenticatedRequest("POST", "/api/borrows/approve", body, adminToken)
		resp, err := server.App.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var approvedBorrowing entity.Borrowing
		err = json.NewDecoder(resp.Body).Decode(&approvedBorrowing)
		require.NoError(t, err)

		// Check basic response fields
		assert.Equal(t, pendingBorrowing.ID, approvedBorrowing.ID)
		assert.Equal(t, entity.APPROVAL_APPROVED, approvedBorrowing.ApprovalStatus)
		assert.Equal(t, adminUser.ID, approvedBorrowing.ApprovedBy)

		// Verify DB State
		var dbBorrowing entity.Borrowing
		err = server.DB.First(&dbBorrowing, pendingBorrowing.ID).Error
		t.Log(dbBorrowing)
		t.Log(dbBorrowing.ApprovedAt)
		require.NoError(t, err)
		assert.Equal(t, entity.APPROVAL_APPROVED, dbBorrowing.ApprovalStatus)
		assert.Equal(t, adminUser.ID, dbBorrowing.ApprovedBy)
		assert.NotEmpty(t, dbBorrowing.ApprovedAt)
	})
}

func TestRejectBorrowing(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	// Helper to get tokens
	adminToken := getAuthToken(t, server, "test_admin", "P@ssw0rd")
	userToken := getAuthToken(t, server, "test_user", "P@ssw0rd")

	// Get a pending borrowing from seed data
	var pendingBorrowing entity.Borrowing
	err := server.DB.Where("borrowing_status = ? AND approval_status = ?", "pending", "pending").First(&pendingBorrowing).Error
	require.NoError(t, err)

	var adminUser entity.User
	err = server.DB.Where("username = ?", "test_admin").First(&adminUser).Error
	require.NoError(t, err)

	t.Run("Authentication and Authorization", func(t *testing.T) {
		payload := map[string]any{
			"id":          pendingBorrowing.ID,
			"rejected_by": adminUser.ID,
		}
		body, _ := json.Marshal(payload)

		t.Run("No Token", func(t *testing.T) {
			req := createAuthenticatedRequest("POST", "/api/borrows/reject", body, "")
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("Regular User Token", func(t *testing.T) {
			req := createAuthenticatedRequest("POST", "/api/borrows/reject", body, userToken)
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})
	})

	t.Run("Validation Errors", func(t *testing.T) {
		tests := []struct {
			name           string
			payload        map[string]any
			expectedStatus int
			expectedError  string
		}{
			{
				name: "Missing Borrowing ID",
				payload: map[string]any{
					"rejected_by": adminUser.ID,
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "borrowing ID is required",
			},
			{
				name: "Missing RejectedBy",
				payload: map[string]any{
					"id": pendingBorrowing.ID,
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "rejectedBy is required",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, _ := json.Marshal(tt.payload)
				req := createAuthenticatedRequest("POST", "/api/borrows/reject", body, adminToken)
				resp, err := server.App.Test(req)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				var respBody map[string]string
				json.NewDecoder(resp.Body).Decode(&respBody)
				if tt.expectedError != "" {
					assert.Contains(t, respBody["error"], tt.expectedError)
				}
			})
		}
	})

	t.Run("Business Logic Errors", func(t *testing.T) {
		t.Run("Borrowing Not Found", func(t *testing.T) {
			payload := map[string]any{
				"id":          999999,
				"rejected_by": adminUser.ID,
			}
			body, _ := json.Marshal(payload)
			req := createAuthenticatedRequest("POST", "/api/borrows/reject", body, adminToken)
			resp, err := server.App.Test(req)

			var respBody map[string]string
			json.NewDecoder(resp.Body).Decode(&respBody)
			t.Log(respBody)

			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("Borrowing Not Pending", func(t *testing.T) {
			// Get an already approved borrowing
			var activeBorrowing entity.Borrowing
			err := server.DB.Where("borrowing_status = ?", "active").First(&activeBorrowing).Error
			require.NoError(t, err)

			payload := map[string]any{
				"id":          activeBorrowing.ID,
				"rejected_by": adminUser.ID,
			}
			body, _ := json.Marshal(payload)
			req := createAuthenticatedRequest("POST", "/api/borrows/reject", body, adminToken)
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusConflict, resp.StatusCode)
		})
	})

	t.Run("Successful Rejection", func(t *testing.T) {
		payload := map[string]any{
			"id":          pendingBorrowing.ID,
			"rejected_by": adminUser.ID,
		}
		body, _ := json.Marshal(payload)

		req := createAuthenticatedRequest("POST", "/api/borrows/reject", body, adminToken)
		resp, err := server.App.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rejectedBorrowing entity.Borrowing
		err = json.NewDecoder(resp.Body).Decode(&rejectedBorrowing)
		require.NoError(t, err)

		// Check basic response fields
		assert.Equal(t, pendingBorrowing.ID, rejectedBorrowing.ID)
		assert.Equal(t, entity.APPROVAL_REJECTED, rejectedBorrowing.ApprovalStatus)
		assert.Equal(t, entity.BORROWING_CANCELLED, rejectedBorrowing.BorrowingStatus)
		assert.Equal(t, adminUser.ID, rejectedBorrowing.RejectedBy)

		// Verify DB State
		var dbBorrowing entity.Borrowing
		err = server.DB.First(&dbBorrowing, pendingBorrowing.ID).Error
		require.NoError(t, err)
		assert.Equal(t, entity.APPROVAL_REJECTED, dbBorrowing.ApprovalStatus)
		assert.Equal(t, entity.BORROWING_CANCELLED, dbBorrowing.BorrowingStatus)
		assert.Equal(t, adminUser.ID, dbBorrowing.RejectedBy)
		assert.NotEmpty(t, dbBorrowing.RejectedAt)
	})
}

func TestGetBorrowingsByBorrowingStatus(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	adminToken := getAuthToken(t, server, "test_admin", "P@ssw0rd")

	t.Run("Validation Errors", func(t *testing.T) {
		req := createAuthenticatedRequest("GET", "/api/borrows/status/invalid_status", nil, adminToken)
		resp, err := server.App.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody map[string]string
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, "invalid borrowing status", respBody["error"])
	})

	t.Run("Successful Retrieval", func(t *testing.T) {
		validStatuses := []string{
			entity.BORROWING_PENDING,
			entity.BORROWING_ACTIVE,
			entity.BORROWING_CANCELLED,
		}

		for _, status := range validStatuses {
			t.Run("Status: "+status, func(t *testing.T) {
				req := createAuthenticatedRequest("GET", "/api/borrows/status/"+status, nil, adminToken)
				resp, err := server.App.Test(req)
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var borrowings []entity.Borrowing
				err = json.NewDecoder(resp.Body).Decode(&borrowings)
				require.NoError(t, err)

				for _, b := range borrowings {
					assert.Equal(t, status, b.BorrowingStatus)
				}
			})
		}
	})
}

func TestGetBorrowingsByApprovalStatus(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	adminToken := getAuthToken(t, server, "test_admin", "P@ssw0rd")

	t.Run("Validation Errors", func(t *testing.T) {
		req := createAuthenticatedRequest("GET", "/api/borrows/approval-status/invalid_status", nil, adminToken)
		resp, err := server.App.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody map[string]string
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, "invalid approval status", respBody["error"])
	})

	t.Run("Successful Retrieval", func(t *testing.T) {
		validStatuses := []string{
			entity.APPROVAL_PENDING,
			entity.APPROVAL_APPROVED,
		}

		for _, status := range validStatuses {
			t.Run("Status: "+status, func(t *testing.T) {
				req := createAuthenticatedRequest("GET", "/api/borrows/approval-status/"+status, nil, adminToken)
				resp, err := server.App.Test(req)
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var borrowings []entity.Borrowing
				err = json.NewDecoder(resp.Body).Decode(&borrowings)
				require.NoError(t, err)

				for _, b := range borrowings {
					assert.Equal(t, status, b.ApprovalStatus)
				}
			})
		}
	})
}
