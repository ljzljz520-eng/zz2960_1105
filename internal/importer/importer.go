package importer

import (
	"context"
	"io"

	"inventoryseal/internal/service"
)

type Report struct {
	Imported int
	Rejected int
	Warnings []string
}

func Import(ctx context.Context, svc *service.Service, batchID, actor string, input io.Reader) (Report, error) {
	rows, warnings, err := ParseRows(input)
	if err != nil {
		return Report{Warnings: warnings}, err
	}
	report := Report{Warnings: warnings}
	for _, row := range rows {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		if err := svc.AddRecord(ctx, ToRecord(row, batchID, actor)); err != nil {
			report.Rejected++
			report.Warnings = append(report.Warnings, row.ID+": "+err.Error())
			continue
		}
		report.Imported++
	}
	return report, nil
}

func ValidateRows(rows []Row) []string {
	issues := make([]string, 0)
	for _, row := range rows {
		if row.ID == "" {
			issues = append(issues, "missing id")
		}
		if row.Expected < 0 || row.Observed < 0 {
			issues = append(issues, row.ID+": negative count")
		}
	}
	return issues
}
