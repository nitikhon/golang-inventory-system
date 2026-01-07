package api

import (
	"encoding/json"
	"fmt"
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
	// Create a fresh item to ensure no existing borrowings conflict
	availableItem := entity.Item{
		Name:            "Test Validation Item",
		Description:     "Item for validation tests",
		AvailableAmount: 10,
		TotalAmount:     10,
		Status:          "available",
	}
	err := server.DB.Create(&availableItem).Error
	require.NoError(t, err)

	var user entity.User
	err = server.DB.Where("username = ?", "test_user").First(&user).Error
	require.NoError(t, err)

	t.Run("Authentication", func(t *testing.T) {
		payload := map[string]any{
			"user_id":          user.ID,
			"item_id":          availableItem.ID,
			"borrowing_amount": 1,
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
				name: "DueDate Before Now",
				payload: map[string]any{
					"user_id":          user.ID,
					"item_id":          availableItem.ID,
					"borrowing_amount": 1,
					"due_date":         time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
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
			"approved_by": adminUser.ID,
		}
		body, _ := json.Marshal(payload)

		t.Run("No Token", func(t *testing.T) {
			req := createAuthenticatedRequest("POST", fmt.Sprintf("/api/borrows/approve/%d", pendingBorrowing.ID), body, "")
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("Regular User Token", func(t *testing.T) {
			req := createAuthenticatedRequest("POST", fmt.Sprintf("/api/borrows/approve/%d", pendingBorrowing.ID), body, userToken)
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})
	})

	t.Run("Validation Errors", func(t *testing.T) {
		tests := []struct {
			name           string
			borrowingID    string
			payload        map[string]any
			expectedStatus int
			expectedError  string
		}{
			{
				name:        "Invalid Borrowing ID format",
				borrowingID: "abc",
				payload: map[string]any{
					"approved_by": adminUser.ID,
				},
				expectedStatus: http.StatusBadRequest,
				expectedError:  errormap.ErrInvalidRequestBody,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, _ := json.Marshal(tt.payload)
				req := createAuthenticatedRequest("POST", "/api/borrows/approve/"+tt.borrowingID, body, adminToken)
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
				"approved_by": adminUser.ID,
			}
			body, _ := json.Marshal(payload)
			req := createAuthenticatedRequest("POST", "/api/borrows/approve/999999", body, adminToken)
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
				"approved_by": adminUser.ID,
			}
			body, _ := json.Marshal(payload)
			req := createAuthenticatedRequest("POST", fmt.Sprintf("/api/borrows/approve/%d", activeBorrowing.ID), body, adminToken)
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusConflict, resp.StatusCode)
		})
	})

	t.Run("Successful Approval", func(t *testing.T) {
		payload := map[string]any{
			"approved_by": adminUser.ID,
		}
		body, _ := json.Marshal(payload)

		req := createAuthenticatedRequest("POST", fmt.Sprintf("/api/borrows/approve/%d", pendingBorrowing.ID), body, adminToken)
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

	var user entity.User
	err := server.DB.Where("username = ?", "test_user").First(&user).Error
	require.NoError(t, err)

	var adminUser entity.User
	err = server.DB.Where("username = ?", "test_admin").First(&adminUser).Error
	require.NoError(t, err)

	// Helper to create a pending borrowing
	createPendingBorrowing := func(userID uint) entity.Borrowing {
		var item entity.Item
		err := server.DB.Where("status = ? AND available_amount > 0", "available").First(&item).Error
		require.NoError(t, err)

		b := entity.Borrowing{
			UserID:          userID,
			ItemID:          item.ID,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_PENDING,
			ApprovalStatus:  entity.APPROVAL_PENDING,
			BorrowedAt:      time.Now().Format(time.RFC3339),
			DueDate:         time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
		}
		err = server.DB.Create(&b).Error
		require.NoError(t, err)
		return b
	}

	t.Run("Authentication and Authorization", func(t *testing.T) {
		body, _ := json.Marshal(map[string]any{})

		t.Run("No Token", func(t *testing.T) {
			b := createPendingBorrowing(user.ID) // borrowing belongs to user
			req := createAuthenticatedRequest("POST", fmt.Sprintf("/api/borrows/reject/%d", b.ID), body, "")
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("Unauthorized - Non-Owner Non-Admin tries to reject", func(t *testing.T) {
			b := createPendingBorrowing(adminUser.ID) // borrowing belongs to admin
			// user tries to reject it
			req := createAuthenticatedRequest("POST", fmt.Sprintf("/api/borrows/reject/%d", b.ID), body, userToken)
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})
	})

	t.Run("Validation Errors", func(t *testing.T) {
		tests := []struct {
			name           string
			borrowingID    string
			expectedStatus int
			expectedError  string
		}{
			{
				name:           "Invalid Borrowing ID format",
				borrowingID:    "abc",
				expectedStatus: http.StatusBadRequest,
				expectedError:  errormap.ErrInvalidRequestBody,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, _ := json.Marshal(map[string]any{})
				req := createAuthenticatedRequest("POST", "/api/borrows/reject/"+tt.borrowingID, body, adminToken)
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
			body, _ := json.Marshal(map[string]any{})
			req := createAuthenticatedRequest("POST", "/api/borrows/reject/999999", body, adminToken)
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("Borrowing Not Pending", func(t *testing.T) {
			// Get an already approved borrowing
			var activeBorrowing entity.Borrowing
			err := server.DB.Where("borrowing_status = ?", "active").First(&activeBorrowing).Error
			require.NoError(t, err)

			body, _ := json.Marshal(map[string]any{})
			req := createAuthenticatedRequest("POST", fmt.Sprintf("/api/borrows/reject/%d", activeBorrowing.ID), body, adminToken)
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusConflict, resp.StatusCode)
		})
	})

	t.Run("Successful Owner Cancellation", func(t *testing.T) {
		b := createPendingBorrowing(user.ID)
		body, _ := json.Marshal(map[string]any{})

		// Owner (user) cancels
		req := createAuthenticatedRequest("POST", fmt.Sprintf("/api/borrows/reject/%d", b.ID), body, userToken)
		resp, err := server.App.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rejectedBorrowing entity.Borrowing
		err = json.NewDecoder(resp.Body).Decode(&rejectedBorrowing)
		require.NoError(t, err)

		assert.Equal(t, b.ID, rejectedBorrowing.ID)
		assert.Equal(t, entity.APPROVAL_REJECTED, rejectedBorrowing.ApprovalStatus)
		assert.Equal(t, entity.BORROWING_CANCELLED, rejectedBorrowing.BorrowingStatus)
		assert.Equal(t, user.ID, rejectedBorrowing.RejectedBy)
	})

	t.Run("Successful Admin Rejection", func(t *testing.T) {
		b := createPendingBorrowing(user.ID) // borrowing belongs to user
		body, _ := json.Marshal(map[string]any{})

		// Admin rejects
		req := createAuthenticatedRequest("POST", fmt.Sprintf("/api/borrows/reject/%d", b.ID), body, adminToken)
		resp, err := server.App.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var rejectedBorrowing entity.Borrowing
		err = json.NewDecoder(resp.Body).Decode(&rejectedBorrowing)
		require.NoError(t, err)

		assert.Equal(t, b.ID, rejectedBorrowing.ID)
		assert.Equal(t, entity.APPROVAL_REJECTED, rejectedBorrowing.ApprovalStatus)
		assert.Equal(t, entity.BORROWING_CANCELLED, rejectedBorrowing.BorrowingStatus)
		assert.Equal(t, adminUser.ID, rejectedBorrowing.RejectedBy)
	})
}

func TestGetBorrowingsByBorrowingStatus(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	adminToken := getAuthToken(t, server, "test_admin", "P@ssw0rd")

	t.Run("Validation Errors", func(t *testing.T) {
		req := createAuthenticatedRequest("GET", "/api/borrows/status/?status=ok", nil, adminToken)
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
				req := createAuthenticatedRequest("GET", "/api/borrows/status/", nil, adminToken)
				resp, err := server.App.Test(req)
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var result entity.PaginationResult[entity.Borrowing]
				err = json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err)

				for _, b := range result.Data {
					assert.Equal(t, status, b.BorrowingStatus)
				}
			})
		}

		t.Run("More than one valid status ", func(t *testing.T) {
			req := createAuthenticatedRequest("GET", "/api/borrows/status/", nil, adminToken)
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var result entity.PaginationResult[entity.Borrowing]
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)

			for _, b := range result.Data {
				assert.Contains(t, validStatuses, b.BorrowingStatus)
			}
		})
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

func TestGetBorrowingByUserID(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()
	userToken := getAuthToken(t, server, "test_user", "P@ssw0rd")
	var user entity.User
	err := server.DB.Where("username = ?", "test_user").First(&user).Error
	require.NoError(t, err)
	t.Run("Validation Errors", func(t *testing.T) {
		req := createAuthenticatedRequest("GET", "/api/borrows/user", nil, "")
		resp, err := server.App.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
	t.Run("Successful Retrieval", func(t *testing.T) {
		// Create a specific borrowing for this user
		var item entity.Item
		err := server.DB.Where("status = ?", "available").First(&item).Error
		require.NoError(t, err)
		b := entity.Borrowing{
			UserID:          user.ID,
			ItemID:          item.ID,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_PENDING,
			ApprovalStatus:  entity.APPROVAL_PENDING,
			BorrowedAt:      time.Now().Format(time.RFC3339),
			DueDate:         time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
		}
		err = server.DB.Create(&b).Error
		require.NoError(t, err)
		req := createAuthenticatedRequest("GET", "/api/borrows/user", nil, userToken)
		resp, err := server.App.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result entity.PaginationResult[entity.Borrowing]
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		found := false
		for _, borrowing := range result.Data {
			if borrowing.ID == b.ID {
				found = true
				assert.Equal(t, user.ID, borrowing.UserID)
				assert.Equal(t, item.ID, borrowing.ItemID)
				break
			}
		}
		assert.True(t, found, "Created borrowing should be in the list")
	})
}
func TestUserStats(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()
	userToken := getAuthToken(t, server, "test_user", "P@ssw0rd")
	var user entity.User
	err := server.DB.Where("username = ?", "test_user").First(&user).Error
	require.NoError(t, err)
	var item entity.Item
	err = server.DB.First(&item).Error
	require.NoError(t, err)
	// add one of each type to ensure non-zero check logic works
	server.DB.Create(&entity.Borrowing{
		UserID:          user.ID,
		ItemID:          item.ID,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
	})
	server.DB.Create(&entity.Borrowing{
		UserID:          user.ID,
		ItemID:          item.ID,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_ACTIVE,
	})
	server.DB.Create(&entity.Borrowing{
		UserID:          user.ID,
		ItemID:          item.ID,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_RETURNED,
	})
	req := createAuthenticatedRequest("GET", "/api/borrows/stats", nil, userToken)
	resp, err := server.App.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var stats entity.BorrowingStats
	err = json.NewDecoder(resp.Body).Decode(&stats)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats.OnGoingBorrows, uint(1))
	assert.GreaterOrEqual(t, stats.CurrentlyBorrows, uint(1))
	assert.GreaterOrEqual(t, stats.TotalReturned, uint(1))
}

func TestReturnBorrowing(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	// Helper to get tokens
	adminToken := getAuthToken(t, server, "test_admin", "P@ssw0rd")

	var user entity.User
	err := server.DB.Where("username = ?", "test_user").First(&user).Error
	require.NoError(t, err)

	// Helper to create an active borrowing
	createActiveBorrowing := func() entity.Borrowing {
		var item entity.Item
		// Find item with stock
		err := server.DB.Where("status = ? AND available_amount > 0", "available").First(&item).Error
		require.NoError(t, err)

		b := entity.Borrowing{
			UserID:          user.ID,
			ItemID:          item.ID,
			BorrowingAmount: 1,
			BorrowingStatus: entity.BORROWING_ACTIVE,
			ApprovalStatus:  entity.APPROVAL_APPROVED,
			BorrowedAt:      time.Now().Format(time.RFC3339),
			DueDate:         time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
		}
		err = server.DB.Create(&b).Error
		require.NoError(t, err)
		return b
	}

	t.Run("Validation Errors", func(t *testing.T) {
		req := createAuthenticatedRequest("POST", "/api/borrows/return/abc", nil, adminToken)
		resp, err := server.App.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Business Logic Errors", func(t *testing.T) {
		t.Run("Borrowing Not Found", func(t *testing.T) {
			req := createAuthenticatedRequest("POST", "/api/borrows/return/999999", nil, adminToken)
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("Borrowing Not Active", func(t *testing.T) {
			// Create pending borrowing
			var item entity.Item
			server.DB.First(&item)
			b := entity.Borrowing{
				UserID:          user.ID,
				ItemID:          item.ID,
				BorrowingAmount: 1,
				BorrowingStatus: entity.BORROWING_PENDING,
			}
			server.DB.Create(&b)

			req := createAuthenticatedRequest("POST", fmt.Sprintf("/api/borrows/return/%d", b.ID), nil, adminToken)
			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusConflict, resp.StatusCode)
		})
	})

	t.Run("Successful Return", func(t *testing.T) {
		b := createActiveBorrowing()

		// Capture item amount before return
		var itemBefore entity.Item
		server.DB.First(&itemBefore, b.ItemID)

		req := createAuthenticatedRequest("POST", fmt.Sprintf("/api/borrows/return/%d", b.ID), nil, adminToken)
		resp, err := server.App.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var returnedBorrowing entity.Borrowing
		err = json.NewDecoder(resp.Body).Decode(&returnedBorrowing)
		require.NoError(t, err)

		assert.Equal(t, entity.BORROWING_RETURNED, returnedBorrowing.BorrowingStatus)
		assert.NotEmpty(t, returnedBorrowing.ReturnedAt)

		// Verify Item Stock Increased
		var itemAfter entity.Item
		server.DB.First(&itemAfter, b.ItemID)
		assert.Equal(t, itemBefore.AvailableAmount+b.BorrowingAmount, itemAfter.AvailableAmount)
	})
}
