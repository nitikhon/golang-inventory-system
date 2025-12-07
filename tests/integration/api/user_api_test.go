package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/tests/integration/api/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// same as register
func TestCreateUser(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	t.Run("Success Cases", func(t *testing.T) {
		tests := []struct {
			name    string
			payload map[string]any
		}{
			{
				name: "successful registration",
				payload: map[string]any{
					"username":   "testapi1",
					"email":      "testapi@gmail.com",
					"password":   "P@ssw0rd",
					"phone":      "0981112222",
					"first_name": "Test",
					"last_name":  "Api",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, err := json.Marshal(tt.payload)
				require.NoError(t, err)

				req := httptest.NewRequest("POST", "/api/users/register", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := server.App.Test(req)
				require.NoError(t, err)

				assert.Equal(t, http.StatusCreated, resp.StatusCode)

				var user entity.User
				err = json.NewDecoder(resp.Body).Decode(&user)
				require.NoError(t, err)
				assert.Equal(t, tt.payload["username"], user.Username)
				assert.Equal(t, tt.payload["email"], user.Email)
				assert.Equal(t, tt.payload["phone"], user.Phone)
			})
		}
	})

	t.Run("Validation Errors", func(t *testing.T) {
		tests := []struct {
			name          string
			payload       map[string]any
			expectedError string
		}{
			{
				name: "missing username",
				payload: map[string]any{
					"email":      "test@gmail.com",
					"password":   "P@ssw0rd",
					"phone":      "0981112222",
					"first_name": "Test",
					"last_name":  "Api",
				},
				expectedError: "username is required",
			},
			{
				name: "username too short",
				payload: map[string]any{
					"username":   "ab",
					"email":      "test@gmail.com",
					"password":   "P@ssw0rd",
					"phone":      "0981112222",
					"first_name": "Test",
					"last_name":  "Api",
				},
				expectedError: "username must be 3-20 characters",
			},
			{
				name: "missing email",
				payload: map[string]any{
					"username":   "testuser",
					"password":   "P@ssw0rd",
					"phone":      "0981112222",
					"first_name": "Test",
					"last_name":  "Api",
				},
				expectedError: "email is required",
			},
			{
				name: "invalid email format",
				payload: map[string]any{
					"username":   "testuser",
					"email":      "invalidemail",
					"password":   "P@ssw0rd",
					"phone":      "0981112222",
					"first_name": "Test",
					"last_name":  "Api",
				},
				expectedError: "invalid email format",
			},
			{
				name: "missing password",
				payload: map[string]any{
					"username":   "testuser",
					"email":      "test@gmail.com",
					"phone":      "0981112222",
					"first_name": "Test",
					"last_name":  "Api",
				},
				expectedError: "password is required",
			},
			{
				name: "password too short",
				payload: map[string]any{
					"username":   "testuser",
					"email":      "test@gmail.com",
					"password":   "Pass1",
					"phone":      "0981112222",
					"first_name": "Test",
					"last_name":  "Api",
				},
				expectedError: "password must be at least 8 characters long",
			},
			{
				name: "password missing uppercase",
				payload: map[string]any{
					"username":   "testuser",
					"email":      "test@gmail.com",
					"password":   "p@ssw0rd",
					"phone":      "0981112222",
					"first_name": "Test",
					"last_name":  "Api",
				},
				expectedError: "password must contain at least one uppercase letter",
			},
			{
				name: "password missing lowercase",
				payload: map[string]any{
					"username":   "testuser",
					"email":      "test@gmail.com",
					"password":   "P@SSW0RD",
					"phone":      "0981112222",
					"first_name": "Test",
					"last_name":  "Api",
				},
				expectedError: "password must contain at least one lowercase letter",
			},
			{
				name: "password missing digit",
				payload: map[string]any{
					"username":   "testuser",
					"email":      "test@gmail.com",
					"password":   "P@ssword",
					"phone":      "0981112222",
					"first_name": "Test",
					"last_name":  "Api",
				},
				expectedError: "password must contain at least one digit",
			},
			{
				name: "password missing special character",
				payload: map[string]any{
					"username":   "testuser",
					"email":      "test@gmail.com",
					"password":   "Passw0rd",
					"phone":      "0981112222",
					"first_name": "Test",
					"last_name":  "Api",
				},
				expectedError: "password must contain at least one special character",
			},
			{
				name: "missing phone",
				payload: map[string]any{
					"username":   "testuser",
					"email":      "test@gmail.com",
					"password":   "P@ssw0rd",
					"first_name": "Test",
					"last_name":  "Api",
				},
				expectedError: "phone is required",
			},
			{
				name: "invalid phone format",
				payload: map[string]any{
					"username":   "testuser",
					"email":      "test@gmail.com",
					"password":   "P@ssw0rd",
					"phone":      "12345",
					"first_name": "Test",
					"last_name":  "Api",
				},
				expectedError: "phone number must be exactly 10 digits",
			},
			{
				name: "missing first name",
				payload: map[string]any{
					"username":  "testuser",
					"email":     "test@gmail.com",
					"password":  "P@ssw0rd",
					"phone":     "0981112222",
					"last_name": "Api",
				},
				expectedError: "first name is required",
			},
			{
				name: "first name too long",
				payload: map[string]any{
					"username":   "testuser",
					"email":      "test@gmail.com",
					"password":   "P@ssw0rd",
					"phone":      "0981112223",
					"first_name": strings.Repeat("LONG", 100),
					"last_name":  "Api",
				},
				expectedError: "first name must not exceed 200 characters",
			},
			{
				name: "missing last name",
				payload: map[string]any{
					"username":   "testuser",
					"email":      "test@gmail.com",
					"password":   "P@ssw0rd",
					"phone":      "0981112222",
					"first_name": "Test",
				},
				expectedError: "last name is required",
			},
			{
				name: "last name too long",
				payload: map[string]any{
					"username":   "newuser",
					"email":      "test@gmail.com",
					"password":   "P@ssw0rd",
					"phone":      "0981112222",
					"first_name": "John",
					"last_name":  strings.Repeat("LONG", 100),
				},
				expectedError: "last name must not exceed 200 characters",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, err := json.Marshal(tt.payload)
				require.NoError(t, err)

				req := httptest.NewRequest("POST", "/api/users/register", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := server.App.Test(req)
				require.NoError(t, err)

				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				var response map[string]any
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, response["error"])
			})
		}
	})

	t.Run("Business Logic Errors", func(t *testing.T) {
		tests := []struct {
			name           string
			payload        map[string]any
			expectedStatus int
			expectedError  string
		}{
			{
				name: "duplicate username",
				payload: map[string]any{
					"username":   "test_user",
					"email":      "newemail@gmail.com",
					"password":   "P@ssw0rd",
					"phone":      "0888777666",
					"first_name": "New",
					"last_name":  "User",
				},
				expectedStatus: http.StatusConflict,
				expectedError:  "a user with the provided credentials already exists",
			},
			{
				name: "duplicate email",
				payload: map[string]any{
					"username":   "newusername",
					"email":      "user@test.com",
					"password":   "P@ssw0rd",
					"phone":      "0777666555",
					"first_name": "New",
					"last_name":  "User",
				},
				expectedStatus: http.StatusConflict,
				expectedError:  "a user with the provided credentials already exists",
			},
			{
				name: "duplicate phone",
				payload: map[string]any{
					"username":   "anotherusername",
					"email":      "another@gmail.com",
					"password":   "P@ssw0rd",
					"phone":      "0987654321",
					"first_name": "Another",
					"last_name":  "User",
				},
				expectedStatus: http.StatusConflict,
				expectedError:  "a user with the provided credentials already exists",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, err := json.Marshal(tt.payload)
				require.NoError(t, err)

				req := httptest.NewRequest("POST", "/api/users/register", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := server.App.Test(req)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedStatus, resp.StatusCode)

				var response map[string]string
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, response["error"])
			})
		}
	})
}

func TestUpdateUser(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	// generate admin token for authenticated requests
	adminToken := loginAs(t, server, "test_admin", "P@ssw0rd")

	// create a user to update
	createPayload := map[string]any{
		"username":   "updatetest",
		"email":      "updatetest@gmail.com",
		"password":   "P@ssw0rd",
		"phone":      "0912345678",
		"first_name": "Update",
		"last_name":  "Test",
	}
	body, _ := json.Marshal(createPayload)
	req := httptest.NewRequest("POST", "/api/users/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := server.App.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdUser entity.User
	json.NewDecoder(resp.Body).Decode(&createdUser)

	// only first_name, last_name, phone, email, password allowed
	t.Run("Success Cases", func(t *testing.T) {
		tests := []struct {
			name    string
			payload map[string]any
		}{
			{
				name: "update user successfully",
				payload: map[string]any{
					"id":         createdUser.ID,
					"email":      "updated@gmail.com",
					"password":   "NewP@ssw0rd",
					"phone":      "0998877665",
					"first_name": "Updated",
					"last_name":  "User",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, err := json.Marshal(tt.payload)
				require.NoError(t, err)

				req := httptest.NewRequest("PUT", fmt.Sprintf("/api/users/%d", createdUser.ID), bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+adminToken)

				resp, err := server.App.Test(req)
				require.NoError(t, err)

				if resp.StatusCode != http.StatusOK {
					var errResp map[string]any
					json.NewDecoder(resp.Body).Decode(&errResp)
					t.Logf("Unexpected response: status=%d, body=%v", resp.StatusCode, errResp)
				}
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var user entity.User
				err = json.NewDecoder(resp.Body).Decode(&user)
				require.NoError(t, err)
				assert.Equal(t, createdUser.Username, user.Username)
				assert.Equal(t, tt.payload["email"], user.Email)
			})
		}
	})

	t.Run("Validation Errors", func(t *testing.T) {
		tests := []struct {
			name          string
			payload       map[string]any
			expectedError string
		}{
			{
				name: "missing user ID in body",
				payload: map[string]any{
					"username":   "testuser",
					"email":      "test@gmail.com",
					"password":   "P@ssw0rd",
					"phone":      "0981112222",
					"first_name": "Test",
					"last_name":  "User",
				},
				expectedError: "user ID is required for update",
			},
			{
				name: "invalid email format",
				payload: map[string]any{
					"id":         createdUser.ID,
					"username":   "testuser",
					"email":      "invalidemail",
					"password":   "P@ssw0rd",
					"phone":      "0981112222",
					"first_name": "Test",
					"last_name":  "User",
				},
				expectedError: "invalid email format",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, err := json.Marshal(tt.payload)
				require.NoError(t, err)

				req := httptest.NewRequest("PUT", fmt.Sprintf("/api/users/%d", createdUser.ID), bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+adminToken)

				resp, err := server.App.Test(req)
				require.NoError(t, err)

				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				var response map[string]string
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, response["error"])
			})
		}
	})

	t.Run("Business Logic Failures", func(t *testing.T) {
		tests := []struct {
			name          string
			payload       map[string]any
			expectedError string
		}{
			{
				name: "duplicate email",
				payload: map[string]any{
					"id":         createdUser.ID,
					"email":      "admin@test.com", // existing admin email
					"password":   "P@ssw0rd",
					"phone":      "0998877665",
					"first_name": "Updated",
					"last_name":  "User",
				},
				expectedError: "email already taken",
			},
			{
				name: "duplicate phone",
				payload: map[string]any{
					"id":         createdUser.ID,
					"email":      "updated@gmail.com",
					"password":   "P@ssw0rd",
					"phone":      "1234567890", // existing admin phone
					"first_name": "Updated",
					"last_name":  "User",
				},
				expectedError: "phone already taken",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, err := json.Marshal(tt.payload)
				require.NoError(t, err)

				req := httptest.NewRequest("PUT", fmt.Sprintf("/api/users/%d", createdUser.ID), bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+adminToken)

				resp, err := server.App.Test(req)
				require.NoError(t, err)

				assert.Equal(t, http.StatusConflict, resp.StatusCode)

				var response map[string]string
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, response["error"])
			})
		}
	})

	t.Run("Authorization Errors", func(t *testing.T) {
		t.Run("no token provided", func(t *testing.T) {
			payload := map[string]any{
				"id":         createdUser.ID,
				"username":   "testuser",
				"email":      "test@gmail.com",
				"password":   "P@ssw0rd",
				"phone":      "0981112222",
				"first_name": "Test",
				"last_name":  "User",
			}
			body, _ := json.Marshal(payload)

			req := httptest.NewRequest("PUT", fmt.Sprintf("/api/users/%d", createdUser.ID), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("invalid URL param ID", func(t *testing.T) {
			payload := map[string]any{
				"id":         createdUser.ID,
				"username":   "testuser",
				"email":      "test@gmail.com",
				"password":   "P@ssw0rd",
				"phone":      "0981112222",
				"first_name": "Test",
				"last_name":  "User",
			}
			body, _ := json.Marshal(payload)

			req := httptest.NewRequest("PUT", "/api/users/invalid", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
}

func TestDeleteUser(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	adminToken := loginAs(t, server, "test_admin", "P@ssw0rd")
	userToken := loginAs(t, server, "test_user", "P@ssw0rd")

	// Create a user to delete
	createPayload := map[string]any{
		"username":   "deletetest",
		"email":      "deletetest@gmail.com",
		"password":   "P@ssw0rd",
		"phone":      "0911223344",
		"first_name": "Delete",
		"last_name":  "Test",
	}
	body, _ := json.Marshal(createPayload)
	req := httptest.NewRequest("POST", "/api/users/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := server.App.Test(req)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdUser entity.User
	json.NewDecoder(resp.Body).Decode(&createdUser)

	t.Run("Success Cases", func(t *testing.T) {
		t.Run("admin deletes user successfully", func(t *testing.T) {
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/users/%d", createdUser.ID), nil)
			req.Header.Set("Authorization", "Bearer "+adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		})
	})

	t.Run("Authorization Errors", func(t *testing.T) {
		t.Run("regular user cannot delete", func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/api/users/1", nil)
			req.Header.Set("Authorization", "Bearer "+userToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("no token provided", func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/api/users/1", nil)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})

	t.Run("Validation Errors", func(t *testing.T) {
		t.Run("invalid user ID", func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/api/users/invalid", nil)
			req.Header.Set("Authorization", "Bearer "+adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
}

func TestGetAllUsers(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	adminToken := loginAs(t, server, "test_admin", "P@ssw0rd")
	userToken := loginAs(t, server, "test_user", "P@ssw0rd")

	t.Run("Success Cases", func(t *testing.T) {
		t.Run("admin gets all users", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/", nil)
			req.Header.Set("Authorization", "Bearer "+adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var users []entity.User
			err = json.NewDecoder(resp.Body).Decode(&users)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(users), 2) // At least test_admin and test_user
		})
	})

	t.Run("Authorization Errors", func(t *testing.T) {
		t.Run("regular user cannot get all users", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/", nil)
			req.Header.Set("Authorization", "Bearer "+userToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("no token provided", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/", nil)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})
}

func TestGetUserByID(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	adminToken := loginAs(t, server, "test_admin", "P@ssw0rd")
	userToken := loginAs(t, server, "test_user", "P@ssw0rd")

	t.Run("Success Cases", func(t *testing.T) {
		t.Run("admin gets user by ID", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/1", nil)
			req.Header.Set("Authorization", "Bearer "+adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var user entity.User
			err = json.NewDecoder(resp.Body).Decode(&user)
			require.NoError(t, err)
			assert.Equal(t, uint(1), user.ID)
		})
	})

	t.Run("Authorization Errors", func(t *testing.T) {
		t.Run("regular user cannot get user by ID", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/1", nil)
			req.Header.Set("Authorization", "Bearer "+userToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("no token provided", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/1", nil)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})

	t.Run("Validation Errors", func(t *testing.T) {
		t.Run("invalid user ID", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/invalid", nil)
			req.Header.Set("Authorization", "Bearer "+adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
}

func TestGetUserByUsername(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	adminToken := loginAs(t, server, "test_admin", "P@ssw0rd")
	userToken := loginAs(t, server, "test_user", "P@ssw0rd")

	t.Run("Success Cases", func(t *testing.T) {
		t.Run("admin gets user by username", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/username/test_user", nil)
			req.Header.Set("Authorization", "Bearer "+adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var user entity.User
			err = json.NewDecoder(resp.Body).Decode(&user)
			require.NoError(t, err)
			assert.Equal(t, "test_user", user.Username)
		})
	})

	t.Run("Not Found Errors", func(t *testing.T) {
		t.Run("user not found", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/username/nonexistent", nil)
			req.Header.Set("Authorization", "Bearer "+adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})

	t.Run("Authorization Errors", func(t *testing.T) {
		t.Run("regular user cannot search by username", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/username/test_admin", nil)
			req.Header.Set("Authorization", "Bearer "+userToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("no token provided", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/username/test_user", nil)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})
}

func TestGetUserByPhone(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	adminToken := loginAs(t, server, "test_admin", "P@ssw0rd")
	userToken := loginAs(t, server, "test_user", "P@ssw0rd")

	t.Run("Success Cases", func(t *testing.T) {
		t.Run("admin gets user by phone", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/phone/0987654321", nil)
			req.Header.Set("Authorization", "Bearer "+adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var user entity.User
			err = json.NewDecoder(resp.Body).Decode(&user)
			require.NoError(t, err)
			assert.Equal(t, "0987654321", user.Phone)
		})
	})

	t.Run("Not Found Errors", func(t *testing.T) {
		t.Run("phone not found", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/phone/0000000000", nil)
			req.Header.Set("Authorization", "Bearer "+adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})

	t.Run("Authorization Errors", func(t *testing.T) {
		t.Run("regular user cannot search by phone", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/phone/1234567890", nil)
			req.Header.Set("Authorization", "Bearer "+userToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("no token provided", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/phone/0987654321", nil)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})

	t.Run("Validation Errors", func(t *testing.T) {
		t.Run("invalid phone format", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/phone/12345", nil)
			req.Header.Set("Authorization", "Bearer "+adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
}

func TestGetUserByEmail(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	adminToken := loginAs(t, server, "test_admin", "P@ssw0rd")
	userToken := loginAs(t, server, "test_user", "P@ssw0rd")

	t.Run("Success Cases", func(t *testing.T) {
		t.Run("admin gets user by email", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/email/user@test.com", nil)
			req.Header.Set("Authorization", "Bearer "+adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var user entity.User
			err = json.NewDecoder(resp.Body).Decode(&user)
			require.NoError(t, err)
			assert.Equal(t, "user@test.com", user.Email)
		})
	})

	t.Run("Not Found Errors", func(t *testing.T) {
		t.Run("email not found", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/email/nonexistent@test.com", nil)
			req.Header.Set("Authorization", "Bearer "+adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})

	t.Run("Authorization Errors", func(t *testing.T) {
		t.Run("regular user cannot search by email", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/email/admin@test.com", nil)
			req.Header.Set("Authorization", "Bearer "+userToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("no token provided", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/email/user@test.com", nil)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})

	t.Run("Validation Errors", func(t *testing.T) {
		t.Run("invalid email format", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/email/invalidemail", nil)
			req.Header.Set("Authorization", "Bearer "+adminToken)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
}

func TestLogin(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	t.Run("Success Cases", func(t *testing.T) {
		tests := []struct {
			name     string
			username string
			password string
		}{
			{
				name:     "admin login",
				username: "test_admin",
				password: "P@ssw0rd",
			},
			{
				name:     "regular user login",
				username: "test_user",
				password: "P@ssw0rd",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				payload := map[string]string{
					"username": tt.username,
					"password": tt.password,
				}
				body, _ := json.Marshal(payload)

				req := httptest.NewRequest("POST", "/api/users/login", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := server.App.Test(req)
				require.NoError(t, err)

				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var result map[string]string
				err = json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err)
				assert.NotEmpty(t, result["access_token"])

				// Check refresh token cookie is set
				var hasRefreshCookie bool
				for _, cookie := range resp.Cookies() {
					if cookie.Name == "refresh_token" {
						hasRefreshCookie = true
						assert.NotEmpty(t, cookie.Value)
						assert.True(t, cookie.HttpOnly)
						break
					}
				}
				assert.True(t, hasRefreshCookie, "refresh_token cookie should be set")
			})
		}
	})

	t.Run("Authentication Errors", func(t *testing.T) {
		tests := []struct {
			name          string
			payload       map[string]string
			expectedError string
		}{
			{
				name: "wrong password",
				payload: map[string]string{
					"username": "test_user",
					"password": "WrongPassword1!",
				},
				expectedError: "invalid credentials",
			},
			{
				name: "nonexistent user",
				payload: map[string]string{
					"username": "nonexistent",
					"password": "P@ssw0rd",
				},
				expectedError: "invalid credentials",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, _ := json.Marshal(tt.payload)

				req := httptest.NewRequest("POST", "/api/users/login", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := server.App.Test(req)
				require.NoError(t, err)

				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

				var result map[string]string
				err = json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, result["error"])
			})
		}
	})

	t.Run("Validation Errors", func(t *testing.T) {
		tests := []struct {
			name          string
			payload       map[string]string
			expectedError string
		}{
			{
				name: "empty username",
				payload: map[string]string{
					"username": "",
					"password": "P@ssw0rd",
				},
				expectedError: "invalid credentials",
			},
			{
				name: "username too short",
				payload: map[string]string{
					"username": "ab",
					"password": "P@ssw0rd",
				},
				expectedError: "username must be 3-20 characters",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, _ := json.Marshal(tt.payload)

				req := httptest.NewRequest("POST", "/api/users/login", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := server.App.Test(req)
				require.NoError(t, err)

				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				var result map[string]string
				err = json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, result["error"])
			})
		}
	})

	t.Run("Invalid Request", func(t *testing.T) {
		t.Run("invalid JSON", func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/users/login", bytes.NewReader([]byte(`{invalid`)))
			req.Header.Set("Content-Type", "application/json")

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
}

func TestLogout(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	t.Run("Success Cases", func(t *testing.T) {
		t.Run("successful logout", func(t *testing.T) {
			// First login to get token
			token := loginAs(t, server, "test_user", "P@ssw0rd")

			req := httptest.NewRequest("POST", "/api/users/logout", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := server.App.Test(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var result map[string]string
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			assert.Equal(t, "Logged out successfully", result["message"])

			// Check that refresh token cookie is cleared
			for _, cookie := range resp.Cookies() {
				if cookie.Name == "refresh_token" {
					assert.Empty(t, cookie.Value)
					assert.Equal(t, -1, cookie.MaxAge)
					break
				}
			}
		})
	})

	t.Run("Authorization Errors", func(t *testing.T) {
		t.Run("no token provided", func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/users/logout", nil)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("invalid token", func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/users/logout", nil)
			req.Header.Set("Authorization", "Bearer invalidtoken")

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})
}

func TestMe(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	t.Run("Success Cases", func(t *testing.T) {
		t.Run("get current user info - admin", func(t *testing.T) {
			token := loginAs(t, server, "test_admin", "P@ssw0rd")

			req := httptest.NewRequest("GET", "/api/users/me", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := server.App.Test(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var user entity.User
			err = json.NewDecoder(resp.Body).Decode(&user)
			require.NoError(t, err)
			assert.Equal(t, "test_admin", user.Username)
			assert.Equal(t, "admin@test.com", user.Email)
			assert.True(t, user.IsAdmin)
		})

		t.Run("get current user info - regular user", func(t *testing.T) {
			token := loginAs(t, server, "test_user", "P@ssw0rd")

			req := httptest.NewRequest("GET", "/api/users/me", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := server.App.Test(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var user entity.User
			err = json.NewDecoder(resp.Body).Decode(&user)
			require.NoError(t, err)
			assert.Equal(t, "test_user", user.Username)
			assert.Equal(t, "user@test.com", user.Email)
			assert.False(t, user.IsAdmin)
		})
	})

	t.Run("Authorization Errors", func(t *testing.T) {
		t.Run("no token provided", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/me", nil)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("invalid token", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/me", nil)
			req.Header.Set("Authorization", "Bearer invalidtoken")

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})
}

func TestRefreshToken(t *testing.T) {
	server := setup.NewTestServer(t)
	defer server.Cleanup()

	t.Run("Success Cases", func(t *testing.T) {
		t.Run("refresh token successfully", func(t *testing.T) {
			// Login to get refresh token
			_, refreshToken := loginAsWithCookie(t, server, "test_user", "P@ssw0rd")

			req := httptest.NewRequest("POST", "/api/users/refresh-token", nil)
			req.AddCookie(&http.Cookie{
				Name:  "refresh_token",
				Value: refreshToken,
			})

			resp, err := server.App.Test(req)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var result map[string]string
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			assert.NotEmpty(t, result["access_token"])

			// Check new refresh token cookie is set
			var hasNewRefreshCookie bool
			for _, cookie := range resp.Cookies() {
				if cookie.Name == "refresh_token" {
					hasNewRefreshCookie = true
					assert.NotEmpty(t, cookie.Value)
					break
				}
			}
			assert.True(t, hasNewRefreshCookie, "new refresh_token cookie should be set")
		})
	})

	t.Run("Authorization Errors", func(t *testing.T) {
		t.Run("no refresh token cookie", func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/users/refresh-token", nil)

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

			var result map[string]string
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			assert.Equal(t, "refresh token not found", result["error"])
		})

		t.Run("invalid refresh token", func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/users/refresh-token", nil)
			req.AddCookie(&http.Cookie{
				Name:  "refresh_token",
				Value: "invalidtoken",
			})

			resp, err := server.App.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

			var result map[string]string
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			assert.Equal(t, "invalid refresh token", result["error"])
		})
	})
}
