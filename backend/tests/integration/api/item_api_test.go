package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/util/errormap"
	"github.com/nitikhon/golang-inventory-system/tests/integration/api/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetItems(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		validateItems  bool
	}{
		{
			name:           "Get all items successfully",
			method:         "GET",
			url:            "/api/items/",
			expectedStatus: http.StatusOK,
			validateItems:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.validateItems {
				var result entity.PaginationResult[entity.Item]
				err = json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err)

				// Should have the seeded items (5 items from setup)
				assert.GreaterOrEqual(t, len(result.Data), 5)

				// Validate structure of first item
				if len(result.Data) > 0 {
					item := result.Data[0]
					assert.NotEmpty(t, item.Name)
					assert.NotEmpty(t, item.Description)
					assert.GreaterOrEqual(t, item.TotalAmount, 0)
					assert.GreaterOrEqual(t, item.AvailableAmount, 0)
					assert.Contains(t, []string{"available", "borrowed", "maintenance", "lost"}, item.Status)
				}
			}
		})
	}
}

func TestGetItemById(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	// Get the first item ID from the database
	var firstItem entity.Item
	err := server.DB.First(&firstItem).Error
	require.NoError(t, err)

	t.Run("Successful Retrieval", func(t *testing.T) {
		t.Run("Get existing item by valid ID", func(t *testing.T) {
			url := fmt.Sprintf("/api/items/%d", firstItem.ID)
			req := httptest.NewRequest("GET", url, nil)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var item entity.Item
			err = json.NewDecoder(resp.Body).Decode(&item)
			require.NoError(t, err)
			assert.Equal(t, firstItem.ID, item.ID)
			assert.Equal(t, firstItem.Name, item.Name)
			assert.Equal(t, firstItem.Description, item.Description)
		})
	})

	t.Run("Error Cases", func(t *testing.T) {
		tests := []struct {
			name           string
			itemID         string
			expectedStatus int
		}{
			{
				name:           "Get item with non-existent ID",
				itemID:         "99999",
				expectedStatus: http.StatusNotFound,
			},
			{
				name:           "Get item with invalid ID format",
				itemID:         "invalid",
				expectedStatus: http.StatusBadRequest,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				url := fmt.Sprintf("/api/items/%s", tt.itemID)
				req := httptest.NewRequest("GET", url, nil)

				resp, err := server.App.Test(req)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			})
		}
	})
}

func TestCreateItem(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	adminToken := getAuthToken(t, server, "test_admin", "P@ssw0rd")
	userToken := getAuthToken(t, server, "test_user", "P@ssw0rd")

	t.Run("Authentication and Authorization", func(t *testing.T) {
		tests := []struct {
			name           string
			payload        map[string]any
			token          string
			expectedStatus int
		}{
			{
				name: "Create item without authentication",
				payload: map[string]any{
					"name":             "Unauthorized Item",
					"description":      "Should fail",
					"available_amount": 1,
					"total_amount":     1,
				},
				token:          "",
				expectedStatus: http.StatusUnauthorized,
			},
			{
				name: "Create item as regular user (should fail)",
				payload: map[string]any{
					"name":             "User Item",
					"description":      "Regular user trying to create",
					"available_amount": 1,
					"total_amount":     1,
				},
				token:          userToken,
				expectedStatus: http.StatusForbidden,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, err := json.Marshal(tt.payload)
				require.NoError(t, err)

				req := createAuthenticatedRequest("POST", "/api/items/", body, tt.token)

				resp, err := server.App.Test(req)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			})
		}
	})

	t.Run("Successful Creation", func(t *testing.T) {
		t.Run("Create item successfully as admin", func(t *testing.T) {
			payload := map[string]any{
				"name":             "Test Laptop",
				"description":      "Test laptop for development",
				"available_amount": 5,
				"total_amount":     10,
				"status":           "available",
			}

			body, err := json.Marshal(payload)
			require.NoError(t, err)

			req := createAuthenticatedRequest("POST", "/api/items/", body, adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			var createdItem entity.Item
			err = json.NewDecoder(resp.Body).Decode(&createdItem)
			require.NoError(t, err)
			assert.Equal(t, payload["name"], createdItem.Name)
			assert.Equal(t, payload["description"], createdItem.Description)
			assert.Greater(t, createdItem.ID, uint(0))
		})
	})

	t.Run("Validation Errors", func(t *testing.T) {
		t.Run("Create item with missing required fields", func(t *testing.T) {
			payload := map[string]any{
				"description": "Missing name field",
			}

			body, err := json.Marshal(payload)
			require.NoError(t, err)

			req := createAuthenticatedRequest("POST", "/api/items/", body, adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})

	t.Run("Business Logic Errors", func(t *testing.T) {
		t.Run("Create item with duplicate name", func(t *testing.T) {
			payload := map[string]any{
				"name":             "Portable Projector",
				"description":      "Duplicate name test",
				"available_amount": 1,
				"total_amount":     1,
			}

			body, err := json.Marshal(payload)
			require.NoError(t, err)

			req := createAuthenticatedRequest("POST", "/api/items/", body, adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusConflict, resp.StatusCode)
		})
	})
}

func TestPutUpdateItem(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	adminToken := getAuthToken(t, server, "test_admin", "P@ssw0rd")
	userToken := getAuthToken(t, server, "test_user", "P@ssw0rd")

	// Get an existing item to update
	var existingItem entity.Item
	err := server.DB.First(&existingItem).Error
	require.NoError(t, err)

	t.Run("Authentication and Authorization", func(t *testing.T) {
		tests := []struct {
			name           string
			payload        map[string]any
			token          string
			expectedStatus int
		}{
			{
				name: "Update item without authentication",
				payload: map[string]any{
					"ID":          existingItem.ID,
					"name":        "Unauthorized Update",
					"description": "Should fail",
				},
				token:          "",
				expectedStatus: http.StatusUnauthorized,
			},
			{
				name: "Update item as regular user (should fail)",
				payload: map[string]any{
					"ID":          existingItem.ID,
					"name":        "User Update",
					"description": "Regular user trying to update",
				},
				token:          userToken,
				expectedStatus: http.StatusForbidden,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, err := json.Marshal(tt.payload)
				require.NoError(t, err)

				req := createAuthenticatedRequest("PUT", "/api/items/", body, tt.token)

				resp, err := server.App.Test(req)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			})
		}
	})

	t.Run("Successful Update", func(t *testing.T) {
		t.Run("Update item successfully as admin", func(t *testing.T) {
			payload := map[string]any{
				"ID":               existingItem.ID,
				"name":             existingItem.Name + " updated",
				"description":      "Updated description",
				"available_amount": existingItem.AvailableAmount + 1,
				"total_amount":     existingItem.TotalAmount + 1,
				"status":           "maintenance",
			}

			body, err := json.Marshal(payload)
			require.NoError(t, err)

			req := createAuthenticatedRequest("PUT", "/api/items/", body, adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var updatedItem entity.Item
			err = json.NewDecoder(resp.Body).Decode(&updatedItem)
			require.NoError(t, err)
			assert.Equal(t, payload["name"], updatedItem.Name)
			assert.Equal(t, payload["description"], updatedItem.Description)
			assert.Equal(t, payload["status"], updatedItem.Status)
		})
	})

	t.Run("Error Cases", func(t *testing.T) {
		t.Run("Update non-existent item", func(t *testing.T) {
			payload := map[string]any{
				"ID":               99999,
				"name":             "Non-existent",
				"description":      "This item doesn't exist",
				"available_amount": 1,
				"total_amount":     1,
				"status":           "maintenance",
			}

			body, err := json.Marshal(payload)
			require.NoError(t, err)

			req := createAuthenticatedRequest("PUT", "/api/items/", body, adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})
}

func TestDeleteItem(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	adminToken := getAuthToken(t, server, "test_admin", "P@ssw0rd")
	userToken := getAuthToken(t, server, "test_user", "P@ssw0rd")

	t.Run("Authentication and Authorization", func(t *testing.T) {
		tests := []struct {
			name           string
			itemID         string
			token          string
			expectedStatus int
		}{
			{
				name:           "Delete item without authentication",
				itemID:         "1",
				token:          "",
				expectedStatus: http.StatusUnauthorized,
			},
			{
				name:           "Delete item as regular user (should fail)",
				itemID:         "1",
				token:          userToken,
				expectedStatus: http.StatusForbidden,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				url := fmt.Sprintf("/api/items/%s", tt.itemID)
				req := createAuthenticatedRequest("DELETE", url, nil, tt.token)

				resp, err := server.App.Test(req)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			})
		}
	})

	t.Run("Successful Deletion", func(t *testing.T) {
		t.Run("Delete item successfully as admin", func(t *testing.T) {
			// Create a new item specifically for deletion test
			newItem := entity.Item{
				Name:            "Item To Delete",
				Description:     "This item will be deleted in test",
				AvailableAmount: 1,
				TotalAmount:     1,
				Status:          "available",
			}
			err := server.DB.Create(&newItem).Error
			require.NoError(t, err)

			url := fmt.Sprintf("/api/items/%d", newItem.ID)
			req := createAuthenticatedRequest("DELETE", url, nil, adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusNoContent, resp.StatusCode)

			// Verify item is actually deleted
			var deletedItem entity.Item
			result := server.DB.First(&deletedItem, newItem.ID)
			assert.Error(t, result.Error) // Should not be found
		})
	})

	t.Run("Error Cases", func(t *testing.T) {
		tests := []struct {
			name           string
			itemID         string
			expectedStatus int
		}{
			{
				name:           "Delete non-existent item",
				itemID:         "99999",
				expectedStatus: http.StatusInternalServerError,
			},
			{
				name:           "Delete item with invalid ID format",
				itemID:         "invalid",
				expectedStatus: http.StatusBadRequest,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				url := fmt.Sprintf("/api/items/%s", tt.itemID)
				req := createAuthenticatedRequest("DELETE", url, nil, adminToken)

				resp, err := server.App.Test(req)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			})
		}
	})
}

func TestPatchUpdateItem(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	adminToken := getAuthToken(t, server, "test_admin", "P@ssw0rd")
	userToken := getAuthToken(t, server, "test_user", "P@ssw0rd")

	var existingItem entity.Item
	err := server.DB.First(&existingItem).Error
	require.NoError(t, err)

	t.Run("Authentication and Authorization", func(t *testing.T) {
		tests := []struct {
			name           string
			itemID         string
			payload        map[string]any
			token          string
			expectedStatus int
		}{
			{
				name:           "Patch update without authentication",
				itemID:         strconv.FormatUint(uint64(existingItem.ID), 10),
				payload:        map[string]any{"name": "Unauthorized Update"},
				token:          "",
				expectedStatus: http.StatusUnauthorized,
			},
			{
				name:           "Patch update as regular user (should fail)",
				itemID:         strconv.FormatUint(uint64(existingItem.ID), 10),
				payload:        map[string]any{"name": "User Update"},
				token:          userToken,
				expectedStatus: http.StatusForbidden,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, err := json.Marshal(tt.payload)
				require.NoError(t, err)

				url := fmt.Sprintf("/api/items/%s", tt.itemID)
				req := createAuthenticatedRequest("PATCH", url, body, tt.token)

				resp, err := server.App.Test(req)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			})
		}
	})

	t.Run("Successful Partial Updates", func(t *testing.T) {
		tests := []struct {
			name           string
			payload        map[string]any
			expectedStatus int
			validate       func(t *testing.T, updated entity.Item, original entity.Item)
		}{
			{
				name:           "Update only name",
				payload:        map[string]any{"name": "Updated Name Only"},
				expectedStatus: http.StatusOK,
				validate: func(t *testing.T, updated entity.Item, original entity.Item) {
					assert.Equal(t, "Updated Name Only", updated.Name)
					assert.Equal(t, original.Description, updated.Description)
					assert.Equal(t, original.AvailableAmount, updated.AvailableAmount)
					assert.Equal(t, original.TotalAmount, updated.TotalAmount)
					assert.Equal(t, original.Status, updated.Status)
				},
			},
			{
				name:           "Update only description",
				payload:        map[string]any{"description": "New Description"},
				expectedStatus: http.StatusOK,
				validate: func(t *testing.T, updated entity.Item, original entity.Item) {
					assert.Equal(t, "New Description", updated.Description)
					assert.Equal(t, original.Name, updated.Name)
				},
			},
			{
				name:           "Update only available_amount",
				payload:        map[string]any{"available_amount": 2},
				expectedStatus: http.StatusOK,
				validate: func(t *testing.T, updated entity.Item, original entity.Item) {
					assert.Equal(t, 2, updated.AvailableAmount)
					assert.Equal(t, original.TotalAmount, updated.TotalAmount)
				},
			},
			{
				name:           "Update only total_amount",
				payload:        map[string]any{"total_amount": 20},
				expectedStatus: http.StatusOK,
				validate: func(t *testing.T, updated entity.Item, original entity.Item) {
					assert.Equal(t, 20, updated.TotalAmount)
				},
			},
			{
				name:           "Update only status",
				payload:        map[string]any{"status": "maintenance"},
				expectedStatus: http.StatusOK,
				validate: func(t *testing.T, updated entity.Item, original entity.Item) {
					assert.Equal(t, "maintenance", updated.Status)
				},
			},
			{
				name: "Update multiple fields at once",
				payload: map[string]any{
					"name":             "Multi Field Update",
					"description":      "Updated multiple fields",
					"available_amount": 5,
					"total_amount":     15,
					"status":           "borrowed",
				},
				expectedStatus: http.StatusOK,
				validate: func(t *testing.T, updated entity.Item, original entity.Item) {
					assert.Equal(t, "Multi Field Update", updated.Name)
					assert.Equal(t, "Updated multiple fields", updated.Description)
					assert.Equal(t, 5, updated.AvailableAmount)
					assert.Equal(t, 15, updated.TotalAmount)
					assert.Equal(t, "borrowed", updated.Status)
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Create a fresh item for each test to avoid side effects
				testItem := entity.Item{
					Name:            "patch test item " + tt.name,
					Description:     "Original description",
					AvailableAmount: 5,
					TotalAmount:     10,
					Status:          "available",
				}
				err := server.DB.Create(&testItem).Error
				require.NoError(t, err)

				body, err := json.Marshal(tt.payload)
				require.NoError(t, err)

				url := fmt.Sprintf("/api/items/%d", testItem.ID)
				req := createAuthenticatedRequest("PATCH", url, body, adminToken)

				resp, err := server.App.Test(req)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				if tt.expectedStatus == http.StatusOK {
					var updatedItem entity.Item
					err = json.NewDecoder(resp.Body).Decode(&updatedItem)
					require.NoError(t, err)
					tt.validate(t, updatedItem, testItem)
				}
			})
		}
	})

	t.Run("Validation Errors", func(t *testing.T) {
		// Create a test item for validation tests
		testItem := entity.Item{
			Name:            "validation test item",
			Description:     "For validation testing",
			AvailableAmount: 5,
			TotalAmount:     10,
			Status:          "available",
		}
		err := server.DB.Create(&testItem).Error
		require.NoError(t, err)

		tests := []struct {
			name           string
			payload        map[string]any
			expectedStatus int
			expectedError  string
		}{
			{
				name:           "Empty name",
				payload:        map[string]any{"name": ""},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "name should not be empty",
			},
			{
				name:           "Negative available_amount",
				payload:        map[string]any{"available_amount": -1},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "available_amount cannot be negative",
			},
			{
				name:           "Zero total_amount",
				payload:        map[string]any{"total_amount": 0},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "total_amount must be greater than 0",
			},
			{
				name:           "Negative total_amount",
				payload:        map[string]any{"total_amount": -5},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "total_amount must be greater than 0",
			},
			{
				name:           "available_amount exceeds total_amount (new available > existing total)",
				payload:        map[string]any{"available_amount": 15}, // testItem.TotalAmount is 10
				expectedStatus: http.StatusBadRequest,
				expectedError:  "available_amount cannot exceed total_amount",
			},
			{
				name:           "available_amount exceeds total_amount (both provided)",
				payload:        map[string]any{"available_amount": 10, "total_amount": 5},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "available_amount cannot exceed total_amount",
			},
			{
				name:           "total_amount less than existing available_amount",
				payload:        map[string]any{"total_amount": 3}, // testItem.AvailableAmount is 5
				expectedStatus: http.StatusBadRequest,
				expectedError:  "updated total_amount is less than item's available amount",
			},
			{
				name:           "Empty status",
				payload:        map[string]any{"status": ""},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "status must be specified",
			},
			{
				name:           "Invalid status value",
				payload:        map[string]any{"status": "invalid_status"},
				expectedStatus: http.StatusBadRequest,
				expectedError:  "status must be one of: available, borrowed, maintenance, lost",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, err := json.Marshal(tt.payload)
				require.NoError(t, err)

				url := fmt.Sprintf("/api/items/%d", testItem.ID)
				req := createAuthenticatedRequest("PATCH", url, body, adminToken)

				resp, err := server.App.Test(req)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				var errorResponse map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errorResponse)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResponse["error"])
			})
		}
	})

	t.Run("Error Cases", func(t *testing.T) {
		tests := []struct {
			name           string
			itemID         string
			payload        map[string]any
			expectedStatus int
			expectedError  string
		}{
			{
				name:           "Non-existent item ID",
				itemID:         "99999",
				payload:        map[string]any{"name": "Update Non-existent"},
				expectedStatus: http.StatusNotFound,
				expectedError:  "item not found",
			},
			{
				name:           "Invalid item ID format",
				itemID:         "invalid",
				payload:        map[string]any{"name": "Update Invalid ID"},
				expectedStatus: http.StatusBadRequest,
			},
			{
				name:           "Invalid JSON payload",
				itemID:         strconv.FormatUint(uint64(existingItem.ID), 10),
				payload:        nil, // Will send invalid JSON
				expectedStatus: http.StatusBadRequest,
				expectedError:  "invalid item payload",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var body []byte
				var err error

				if tt.payload != nil {
					body, err = json.Marshal(tt.payload)
					require.NoError(t, err)
				} else {
					body = []byte("{invalid json}")
				}

				url := fmt.Sprintf("/api/items/%s", tt.itemID)
				req := createAuthenticatedRequest("PATCH", url, body, adminToken)

				resp, err := server.App.Test(req)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				if tt.expectedError != "" {
					var errorResponse map[string]string
					err = json.NewDecoder(resp.Body).Decode(&errorResponse)
					require.NoError(t, err)
					assert.Equal(t, tt.expectedError, errorResponse["error"])
				}
			})
		}
	})

	t.Run("Duplicate Name Conflict", func(t *testing.T) {
		// Create two items
		item1 := entity.Item{
			Name:            "Unique Item One",
			Description:     "First item",
			AvailableAmount: 5,
			TotalAmount:     10,
			Status:          "available",
		}
		item2 := entity.Item{
			Name:            "Unique Item Two",
			Description:     "Second item",
			AvailableAmount: 3,
			TotalAmount:     8,
			Status:          "available",
		}
		err := server.DB.Create(&item1).Error
		require.NoError(t, err)
		err = server.DB.Create(&item2).Error
		require.NoError(t, err)

		// Try to update item2's name to item1's name
		body, err := json.Marshal(map[string]any{"name": "Unique Item One"})
		require.NoError(t, err)

		url := fmt.Sprintf("/api/items/%d", item2.ID)
		req := createAuthenticatedRequest("PATCH", url, body, adminToken)

		resp, err := server.App.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		var errorResponse map[string]string
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, errormap.ErrItemNameAlreadyExists, errorResponse["error"])
	})

	t.Run("Valid Status Values", func(t *testing.T) {
		validStatuses := []string{"available", "borrowed", "maintenance", "lost"}

		for _, status := range validStatuses {
			t.Run("Status: "+status, func(t *testing.T) {
				// Create a test item for each status test
				testItem := entity.Item{
					Name:            "status test " + status,
					Description:     "Testing status update",
					AvailableAmount: 5,
					TotalAmount:     10,
					Status:          "available",
				}
				err := server.DB.Create(&testItem).Error
				require.NoError(t, err)

				body, err := json.Marshal(map[string]any{"status": status})
				require.NoError(t, err)

				url := fmt.Sprintf("/api/items/%d", testItem.ID)
				req := createAuthenticatedRequest("PATCH", url, body, adminToken)

				resp, err := server.App.Test(req)
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var updatedItem entity.Item
				err = json.NewDecoder(resp.Body).Decode(&updatedItem)
				require.NoError(t, err)
				assert.Equal(t, status, updatedItem.Status)
			})
		}
	})
}
