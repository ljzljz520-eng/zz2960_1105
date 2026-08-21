package importer

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"inventoryseal/internal/domain"
)

type Row struct {
	ID       string
	Label    string
	Expected int
	Observed int
}

func ParseRows(input io.Reader) ([]Row, []string, error) {
	scanner := bufio.NewScanner(input)
	rows := make([]Row, 0)
	warnings := make([]string, 0)
	line := 0
	for scanner.Scan() {
		line++
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		parts := strings.Split(text, ",")
		if len(parts) != 4 {
			warnings = append(warnings, fmt.Sprintf("line %d: expected four columns", line))
			continue
		}
		expected, err1 := strconv.Atoi(strings.TrimSpace(parts[2]))
		observed, err2 := strconv.Atoi(strings.TrimSpace(parts[3]))
		if err1 != nil || err2 != nil {
			warnings = append(warnings, fmt.Sprintf("line %d: counts must be numbers", line))
			continue
		}
		rows = append(rows, Row{ID: strings.TrimSpace(parts[0]), Label: strings.TrimSpace(parts[1]), Expected: expected, Observed: observed})
	}
	if err := scanner.Err(); err != nil {
		return nil, warnings, err
	}
	return rows, warnings, nil
}

func ToRecord(row Row, batchID, actor string) domain.Record {
	return domain.Record{ID: row.ID, BatchID: batchID, Label: row.Label, Expected: row.Expected, Observed: row.Observed, Result: domain.EvaluateRecord(domain.Record{Expected: row.Expected, Observed: row.Observed}), UpdatedBy: actor, Version: 1}
}
