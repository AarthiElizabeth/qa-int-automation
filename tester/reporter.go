package tester

import (
	"fmt"
	"os"
	"time"
)

func GenerateHTMLReport(results []Result, filePath string) error {
	// Calculate statistics
	totalTests := len(results)
	passedTests := 0
	for _, r := range results {
		if r.Passed {
			passedTests++
		}
	}
	passRate := float64(passedTests) / float64(totalTests) * 100

	html := `<!DOCTYPE html>
<html>
<head>
	<title>Webhook Test Report</title>
	<style>
		body { font-family: Arial, sans-serif; margin: 20px; }
		h1 { color: #333; }
		.stats { background: #f8f8f8; padding: 15px; border-radius: 5px; margin-bottom: 20px; }
		table { width: 100%; border-collapse: collapse; margin-top: 20px; }
		th { background-color: #f2f2f2; text-align: left; }
		td, th { padding: 12px; border: 1px solid #ddd; }
		pre { white-space: pre-wrap; word-wrap: break-word; max-height: 200px; overflow-y: auto; }
		.pass { background-color: #e8f5e9; }
		.fail { background-color: #ffebee; }
		.status-pass { color: #2e7d32; }
		.status-fail { color: #c62828; }
	</style>
</head>
<body>
	<h1>LINE Webhook Test Report</h1>
	<div class="stats">
		<p><strong>Execution Time:</strong> ` + time.Now().Format(time.RFC1123) + `</p>
		<p><strong>Total Tests:</strong> ` + fmt.Sprintf("%d", totalTests) + `</p>
		<p><strong>Passed:</strong> ` + fmt.Sprintf("%d (%.1f%%)", passedTests, passRate) + `</p>
		<p><strong>Failed:</strong> ` + fmt.Sprintf("%d", totalTests-passedTests) + `</p>
	</div>
	<table>
		<tr>
			<th>Test Name</th>
			<th>Status</th>
			<th>Code</th>
			<th>Message</th>
			<th>Response</th>
		</tr>`

	for _, r := range results {
		statusClass := "pass"
		statusText := "<span class='status-pass'>PASS</span>"
		if !r.Passed {
			statusClass = "fail"
			statusText = "<span class='status-fail'>FAIL</span>"
		}

		html += fmt.Sprintf(
			`<tr class="%s">
				<td><strong>%s</strong></td>
				<td>%s</td>
				<td>%d</td>
				<td>%s</td>
				<td><pre>%s</pre></td>
			</tr>`,
			statusClass, r.Name, statusText, r.Status, r.Message, r.Response,
		)
	}

	html += `</table>
</body>
</html>`

	if err := os.MkdirAll("results", 0755); err != nil {
		return fmt.Errorf("failed to create results directory: %w", err)
	}

	return os.WriteFile(filePath, []byte(html), 0644)
}
