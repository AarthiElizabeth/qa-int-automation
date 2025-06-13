package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"qa-int-automation/junitexport"
	"qa-int-automation/tester"

	"gopkg.in/yaml.v2"
)

type QaseResult struct {
	CaseID     int    `json:"case_id"`
	Status     string `json:"status"`
	Comment    string `json:"comment"`
	TimeMS     int64  `json:"time_ms"`
	Stacktrace string `json:"stacktrace,omitempty"`
}

type QasePayload struct {
	Results []QaseResult `json:"results"`
}

func main() {
	log.Println("Starting test execution")
	startTime := time.Now()

	// Load configuration
	config, err := loadConfig("config/testcases.yaml")
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Run test cases
	qaseResults, testResults := runTestSuite(config)

	// Generate reports
	if err := generateReports(testResults); err != nil {
		log.Printf("Report generation error: %v", err)
	}

	// Send results to QASE
	if qaseRunID := os.Getenv("QASE_RUN_ID"); qaseRunID != "" {
		if err := sendToQase(qaseRunID, qaseResults); err != nil {
			log.Printf("QASE reporting error: %v", err)
		}
	} else {
		log.Println("Skipping QASE reporting - no run ID specified")
	}

	log.Printf("Test execution completed in %v", time.Since(startTime))
}

func loadConfig(path string) (*tester.TestSuite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var suite tester.TestSuite
	if err := yaml.Unmarshal(data, &suite); err != nil {
		return nil, fmt.Errorf("YAML unmarshal error: %w", err)
	}

	// Validate configuration
	if suite.BaseURL == "" || suite.TenantID == "" || suite.ChannelID == "" || suite.Secret == "" {
		return nil, fmt.Errorf("missing required configuration fields")
	}

	return &suite, nil
}

func runTestSuite(suite *tester.TestSuite) ([]QaseResult, []tester.Result) {
	caseMap := map[string]int{
		"TC-001 Text Event":      101,
		"TC-002 Follow Event":    102,
		"TC-003 Unfollow Event":  103,
		"TC-004 Postback Event":  104,
		"TC-005 Rich Menu Event": 105,
		"TC-006 Multiple Choice": 106,
	}

	var qaseResults []QaseResult
	var testResults []tester.Result

	for _, tc := range suite.TestCases {
		log.Printf("Running test case: %s", tc.Name)
		start := time.Now()
		result := tester.RunTestCase(suite.BaseURL, tc, suite.TenantID, suite.ChannelID, suite.Secret)
		duration := time.Since(start)

		status := "passed"
		if !result.Passed {
			status = "failed"
			log.Printf("Test failed: %s - %s", tc.Name, result.Message)
		}

		if caseID, exists := caseMap[tc.Name]; exists {
			qaseResults = append(qaseResults, QaseResult{
				CaseID:     caseID,
				Status:     status,
				Comment:    result.Message,
				TimeMS:     duration.Milliseconds(),
				Stacktrace: result.Response,
			})
		}

		testResults = append(testResults, result)
	}

	return qaseResults, testResults
}

func generateReports(results []tester.Result) error {
	if err := os.MkdirAll("results", 0755); err != nil {
		return fmt.Errorf("failed to create results directory: %w", err)
	}

	// HTML Report
	if err := tester.GenerateHTMLReport(results, "results/report.html"); err != nil {
		return fmt.Errorf("HTML report generation failed: %w", err)
	}

	// JUnit Report
	if err := junitexport.ExportToJUnit(results, "results/test-results.xml"); err != nil {
		return fmt.Errorf("JUnit report generation failed: %w", err)
	}

	return nil
}

func sendToQase(runID string, results []QaseResult) error {
	apiToken := os.Getenv("QASE_API_TOKEN")
	projectCode := os.Getenv("QASE_PROJECT_CODE")
	if apiToken == "" || projectCode == "" {
		return fmt.Errorf("QASE credentials not set")
	}

	payload := QasePayload{
		Results: results,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("https://api.qase.io/v1/result/bulk/%s?run_id=%s", projectCode, runID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", apiToken)
	req.Header.Set("X-Qase-Run", runID)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("QASE API error: %d - %s", resp.StatusCode, string(body))
	}

	log.Printf("Successfully reported %d results to QASE run %s", len(results), runID)
	return nil
}
