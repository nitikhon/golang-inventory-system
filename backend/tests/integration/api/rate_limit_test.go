package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/nitikhon/golang-inventory-system/internal/util"
)

var (
	// this only implement rate limiting for DDoS protection
	ddosReqs = util.ParseInt(os.Getenv("RATE_LIMIT_BOT_MAX")) + 42
	ddosUrl  = "http://localhost:8080/non-existent-path-for-dl-test"

	// this implement rate limiting for normal requests
	normalReqs = util.ParseInt(os.Getenv("RATE_LIMIT_USER_MAX")) + 42
	normalUrl  = "http://localhost:8080/api/users/me"
)

func TestNormalRateLimiting(t *testing.T) {
	fmt.Println("Starting Rate Limiter Verification Script...")

	token := login("john_doe", "user123")

	// Spam requests
	fmt.Printf("Sending %d requests to a protected endpoint...\n", normalReqs)

	var wg sync.WaitGroup
	results := make(chan int, normalReqs)

	start := time.Now()

	for range normalReqs {
		wg.Go(
			func() {
				status := sendRequest(token, normalUrl)
				results <- status
				// Add a tiny delay to not overwhelm network stack too drastically,
				// though for rate limiting we want speed.
				// time.Sleep(10 * time.Millisecond)
			},
		)
	}

	wg.Wait()
	close(results)
	duration := time.Since(start)

	statusCounts := make(map[int]int)
	for status := range results {
		statusCounts[status]++
	}

	fmt.Printf("\n--- Results ---\n")
	fmt.Printf("Time taken: %v\n", duration)
	fmt.Printf("Total Requests: %d\n", normalReqs)
	for code, count := range statusCounts {
		fmt.Printf("Status %d: %d\n", code, count)
	}

	if statusCounts[429] > 0 {
		fmt.Println("\n✅ SUCCESS: Rate limiting triggered (429 Too Many Requests received).")
	} else {
		fmt.Println("\n❌ FAILURE: No 429 errors received. Rate limiter might not be active or configured correctly.")
	}
}

func TestDDOSRateLimiting(t *testing.T) {
	fmt.Println("Starting Rate Limiter Verification Script...")

	// Spam requests
	fmt.Printf("Sending %d requests to a protected endpoint...\n", ddosReqs)

	var wg sync.WaitGroup
	results := make(chan int, ddosReqs)

	start := time.Now()

	for range ddosReqs {
		wg.Go(
			func() {
				status := sendRequest("", ddosUrl)
				results <- status
				// Add a tiny delay to not overwhelm network stack too drastically,
				// though for rate limiting we want speed.
				time.Sleep(100 * time.Millisecond)
			},
		)
	}

	wg.Wait()
	close(results)
	duration := time.Since(start)

	statusCounts := make(map[int]int)
	for status := range results {
		statusCounts[status]++
	}

	fmt.Printf("\n--- Results ---\n")
	fmt.Printf("Time taken: %v\n", duration)
	fmt.Printf("Total Requests: %d\n", ddosReqs)
	for code, count := range statusCounts {
		fmt.Printf("Status %d: %d\n", code, count)
	}

	if statusCounts[429] > 0 {
		fmt.Println("\n✅ SUCCESS: Rate limiting triggered (429 Too Many Requests received).")
	} else {
		fmt.Println("\n❌ FAILURE: No 429 errors received. Rate limiter might not be active or configured correctly.")
	}
}

func sendRequest(token, url string) int {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request error: %v\n", err)
		return 0
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func login(username string, password string) string {
	loginPayload := map[string]string{
		"username": username,
		"password": password,
	}

	body, _ := json.Marshal(loginPayload)

	req, _ := http.NewRequest("POST", "http://localhost:8080/api/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request error: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	var loginResponse struct {
		AccessToken string `json:"access_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	if err != nil {
		fmt.Printf("Response error: %v\n", err)
		return ""
	}
	return loginResponse.AccessToken
}
