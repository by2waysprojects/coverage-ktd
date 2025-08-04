package reporting

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/by2waysprojects/coverage-ktd/model"
)

const (
	reportsDir   = "/reports"
	htmlTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>Security Report - {{.StartTime.Format "2006-01-02"}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 2rem; }
        .header { border-bottom: 2px solid #333; padding-bottom: 1rem; }
        .summary { background: #f8f9fa; padding: 1.5rem; margin: 2rem 0; border-radius: 8px; }
        table { width: 100%; border-collapse: collapse; margin-top: 1rem; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #e9ecef; }
        .detected { background-color: #ffe3e3; }
        .confidence { width: 120px; }
        .percentage { font-weight: bold; color: {{if ge .DetectionRate 50}}#dc3545{{else}}#28a745{{end}}; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Security Test Report</h1>
        <p>Generated: {{.EndTime.Format "2006-01-02 15:04:05"}}</p>
    </div>

    <div class="summary">
        <h2>Executive Summary</h2>
        <p>Total Tests Executed: {{.TotalTests}}</p>
        <p>Threats Detected: <span class="percentage">{{.Detected}} ({{printf "%.1f" .DetectionRate}}%)</span></p>
        <p>Duration: {{.EndTime.Sub .StartTime | printf "%v"}}</p>
    </div>

    <h2>Detailed Results</h2>
    <table>
        <thead>
            <tr>
                <th>Test ID</th>
                <th>Type</th>
                <th>Status</th>
                <th class="confidence">Confidence</th>
                <th>Payload</th>
                <th>Description</th>
            </tr>
        </thead>
        <tbody>
            {{range .Tests}}
            <tr class="{{if .Detected}}detected{{end}}">
                <td>{{.ID}}</td>
                <td>{{.Type}}</td>
                <td>{{if .Detected}}⚠️ Detected{{else}}✅ Clean{{end}}</td>
                <td>{{printf "%.1f" .Confidence}}%</td>
                <td><code>{{.Payload}}</code></td>
                <td>{{.Description}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>`
)

func GenerateHTMLReport(data model.ReportData) error {
	filename := filepath.Join(reportsDir, fmt.Sprintf("report-%s.html",
		data.EndTime.Format("20060102-150405")))

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating HTML file: %w", err)
	}
	defer file.Close()

	tmpl := template.Must(template.New("report").Parse(htmlTemplate))

	data.DetectionRate = (float64(data.Detected) / float64(data.TotalTests)) * 100
	if data.TotalTests == 0 {
		data.DetectionRate = 0
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	fmt.Printf("HTML report generated: %s\n", filename)
	return nil
}
