package qase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"qa-int-automation/tester"
	"strings"
	"time"
)

type QaseClient struct {
	apiToken    string
	projectCode string
	client      *http.Client
}

func NewClient(apiToken, projectCode string) *QaseClient {
	return &QaseClient{
		apiToken:    apiToken,
		projectCode: projectCode,
		client:      &http.Client{Timeout: 15 * time.Second},
	}
}

type QaseTestCase struct {
	Title        string            `json:"title"`
	Automation   int               `json:"automation"`
	Steps        []QaseTestStep    `json:"steps"`
	Tags         []string          `json:"tags"`
	CustomFields map[string]string `json:"custom_fields"`
}

type QaseTestStep struct {
	Action   string `json:"action"`
	Expected string `json:"expected_result"`
}

func (c *QaseClient) SyncTestCases(ctx context.Context, testCases []tester.TestCase) error {
	existingCases, err := c.getExistingCases(ctx)
	if err != nil {
		return fmt.Errorf("failed to get existing cases: %w", err)
	}

	for _, tc := range testCases {
		qaseCase := QaseTestCase{
			Title:      tc.Name,
			Automation: 1, // Mark as automated
			Steps: []QaseTestStep{
				{
					Action:   fmt.Sprintf("Send %s payload", strings.TrimPrefix(tc.Name, "TC-")),
					Expected: fmt.Sprintf("Should return HTTP %d", tc.ExpectedStatus),
				},
			},
			Tags: []string{"auto-sync", "line-webhook"},
			CustomFields: map[string]string{
				"payload_file": tc.PayloadFile,
			},
		}

		if caseID, exists := existingCases[tc.Name]; exists {
			// Update existing case
			if err := c.updateCase(ctx, caseID, qaseCase); err != nil {
				return fmt.Errorf("failed to update case %s: %w", tc.Name, err)
			}
			log.Printf("Updated Qase case %d: %s", caseID, tc.Name)
		} else {
			// Create new case
			caseID, err := c.createCase(ctx, qaseCase)
			if err != nil {
				return fmt.Errorf("failed to create case %s: %w", tc.Name, err)
			}
			log.Printf("Created Qase case %d: %s", caseID, tc.Name)
		}
	}
	return nil
}

func (c *QaseClient) getExistingCases(ctx context.Context) (map[string]int, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("https://api.qase.io/v1/case/%s?limit=100", c.projectCode),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Token", c.apiToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Entities []struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
		} `json:"entities"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	cases := make(map[string]int)
	for _, entity := range result.Entities {
		cases[entity.Title] = entity.ID
	}

	return cases, nil
}

func (c *QaseClient) createCase(ctx context.Context, tc QaseTestCase) (int, error) {
	body, err := json.Marshal(struct {
		Title      string         `json:"title"`
		Automation int            `json:"automation"`
		Steps      []QaseTestStep `json:"steps"`
	}{
		Title:      tc.Title,
		Automation: tc.Automation,
		Steps:      tc.Steps,
	})
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("https://api.qase.io/v1/case/%s", c.projectCode),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", c.apiToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.ID, nil
}

func (c *QaseClient) updateCase(ctx context.Context, caseID int, tc QaseTestCase) error {
	body, err := json.Marshal(tc)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"PATCH",
		fmt.Sprintf("https://api.qase.io/v1/case/%s/%d", c.projectCode, caseID),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", c.apiToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
