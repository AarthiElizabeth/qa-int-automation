package tester

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	defaultHTTPTimeout = 30 * time.Second
)

func RunTestCase(baseURL string, tc TestCase, tenantID, channelID, secret string) Result {
	// Validate inputs
	if baseURL == "" || tc.PayloadFile == "" {
		return Result{
			Name:     tc.Name,
			Passed:   false,
			Message:  "Invalid test configuration",
			Response: "Missing required parameters",
		}
	}

	// Prepare request URL
	query := url.Values{}
	query.Set("tenant_id", tenantID)
	query.Set("channel_id", channelID)
	fullURL := fmt.Sprintf("%s?%s", baseURL, query.Encode())

	// Read payload file
	payload, err := os.ReadFile(tc.PayloadFile)
	if err != nil {
		return Result{
			Name:     tc.Name,
			Passed:   false,
			Message:  "Failed to read payload file",
			Response: err.Error(),
		}
	}

	// Generate signature
	signature := GenerateSignature(secret, payload)

	// Create HTTP request with timeout
	ctx, cancel := context.WithTimeout(context.Background(), defaultHTTPTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", fullURL, bytes.NewBuffer(payload))
	if err != nil {
		return Result{
			Name:     tc.Name,
			Passed:   false,
			Message:  "Failed to create HTTP request",
			Response: err.Error(),
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Line-Signature", signature)

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Result{
			Name:     tc.Name,
			Passed:   false,
			Message:  "HTTP request failed",
			Response: err.Error(),
		}
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{
			Name:     tc.Name,
			Passed:   false,
			Message:  "Failed to read response body",
			Response: err.Error(),
		}
	}

	// Verify status code
	passed := resp.StatusCode == tc.ExpectedStatus
	msg := "Test passed"
	if !passed {
		msg = fmt.Sprintf("Expected status %d but got %d", tc.ExpectedStatus, resp.StatusCode)
	}

	return Result{
		Name:     tc.Name,
		Passed:   passed,
		Status:   resp.StatusCode,
		Message:  msg,
		Response: string(body),
	}
}
