package policy

import (
	"fmt"
	"strings"

	"inventoryseal/internal/domain"
)

type Rule struct {
	Name        string
	Description string
	Check       func(domain.Record) bool
}
type Result struct {
	Passed     bool
	Violations []string
}

func DefaultRules() []Rule {
	return []Rule{
		{Name: "identity", Description: "record ids and labels are present", Check: func(r domain.Record) bool { return r.ID != "" && r.Label != "" }},
		{Name: "counts", Description: "counts are non-negative", Check: func(r domain.Record) bool { return r.Expected >= 0 && r.Observed >= 0 }},
		{Name: "evaluation", Description: "result matches observed count", Check: func(r domain.Record) bool { return r.Result == domain.EvaluateRecord(r) }},
		{Name: "version", Description: "record version is positive", Check: func(r domain.Record) bool { return r.Version > 0 }},
	}
}

func Evaluate(record domain.Record, rules []Rule) Result {
	violations := make([]string, 0)
	for _, rule := range rules {
		if !rule.Check(record) {
			violations = append(violations, rule.Name+": "+rule.Description)
		}
	}
	return Result{Passed: len(violations) == 0, Violations: violations}
}

func Explain(result Result) string {
	if result.Passed {
		return "passed"
	}
	return strings.Join(result.Violations, "; ")
}
func Require(result Result) error {
	if result.Passed {
		return nil
	}
	return fmt.Errorf("policy violation: %s", Explain(result))
}
