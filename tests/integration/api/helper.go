package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nitikhon/golang-inventory-system/tests/integration/api/setup"
	"github.com/stretchr/testify/require"
)

// Helper function to get authentication token
func getAuthToken(t *testing.T, server *setup.TestServer, username, password string) string {
	loginPayload := map[string]string{
		"username": username,
		"password": password,
	}

	body, err := json.Marshal(loginPayload)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := server.App.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResponse struct {
		AccessToken string `json:"access_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	require.NoError(t, err)

	return loginResponse.AccessToken
}

// Helper function to create authenticated request
func createAuthenticatedRequest(method, url string, body []byte, token string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

// Helper function to login and get access token
func loginAs(t *testing.T, server *setup.TestServer, username, password string) string {
	payload := map[string]string{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/users/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := server.App.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Login failed for user: %s", username)

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	return result["access_token"]
}

// Helper function to login and get both access token and refresh token cookie
func loginAsWithCookie(t *testing.T, server *setup.TestServer, username, password string) (string, string) {
	payload := map[string]string{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/users/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := server.App.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Login failed for user: %s", username)

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Extract refresh token from cookie
	var refreshToken string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "refresh_token" {
			refreshToken = cookie.Value
			break
		}
	}

	return result["access_token"], refreshToken
}
