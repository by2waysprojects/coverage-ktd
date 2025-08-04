package reporting

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/by2waysprojects/coverage-ktd/model"
)

func GenerateJSONReport(data model.ReportData) error {
	filename := filepath.Join(reportsDir, fmt.Sprintf("report-%s.json",
		data.EndTime.Format("20060102-150405")))

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	fmt.Printf("JSON report generated: %s\n", filename)
	return nil
}
