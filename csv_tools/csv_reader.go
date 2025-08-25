package csv_tools

import (
	"encoding/csv"
	"fmt"
	"os"
)

func ReadCSV(filePath string, skipHeader bool) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // change if needed
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if skipHeader && len(records) > 0 {
		return records[1:], nil
	}
	return records, nil
}

// AppendCSV appends a new row to the given CSV file.
func AppendCSV(filePath string, row []interface{}) error {
	// Open file in append mode, create if not exists
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush() // ensures the row is written

	stringRow := make([]string, len(row))
	// Convert each element in the row to string
	for i, v := range row {
		switch v := v.(type) {
		case string:
			stringRow[i] = v
		case int:
			stringRow[i] = fmt.Sprintf("%d", v)
		case float64:
			stringRow[i] = fmt.Sprintf("%f", v)
		default:
			stringRow[i] = fmt.Sprintf("%v", v) // fallback for other types
		}
	}
	// Write the row
	if err := writer.Write(stringRow); err != nil {
		return fmt.Errorf("failed to write row: %w", err)
	}

	return nil
}
