package ledger

import (
	"inventoryseal/internal/domain"
	"sort"
)

type Totals struct {
	Total     int
	Match     int
	Overage   int
	Shortage  int
	Confirmed int
	Published int
}
type ActorTotals struct {
	Actor   string
	Actions int
}

func Count(records []domain.Record) Totals {
	result := Totals{Total: len(records)}
	for _, record := range records {
		if record.Result == "match" {
			result.Match++
		}
		if record.Result == "overage" {
			result.Overage++
		}
		if record.Result == "shortage" {
			result.Shortage++
		}
		if record.Confirmed {
			result.Confirmed++
		}
		if record.Published {
			result.Published++
		}
	}
	return result
}

func Actors(events []domain.AuditEvent) []ActorTotals {
	counts := map[string]int{}
	for _, event := range events {
		counts[event.Actor]++
	}
	result := make([]ActorTotals, 0, len(counts))
	for actor, actions := range counts {
		result = append(result, ActorTotals{actor, actions})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Actor < result[j].Actor })
	return result
}

func Timeline(events []domain.AuditEvent) []domain.AuditEvent {
	result := append([]domain.AuditEvent(nil), events...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Sequence == result[j].Sequence {
			return result[i].ID < result[j].ID
		}
		return result[i].Sequence < result[j].Sequence
	})
	return result
}

func Completion(records []domain.Record) float64 {
	if len(records) == 0 {
		return 0
	}
	confirmed := 0
	for _, record := range records {
		if record.Confirmed {
			confirmed++
		}
	}
	return float64(confirmed) / float64(len(records))
}
