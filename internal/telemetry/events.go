package telemetry

import (
	"sort"
	"strings"
	"time"
)

type Event struct {
	Name    string
	BatchID string
	Actor   string
	Detail  string
	At      string
	Value   int
}
type Summary struct {
	Count      int
	ByName     map[string]int
	ByActor    map[string]int
	TotalValue int
}

func New(name, batchID, actor, detail string, value int) Event {
	return Event{Name: name, BatchID: batchID, Actor: actor, Detail: detail, At: time.Unix(0, 0).UTC().Format(time.RFC3339), Value: value}
}
func Normalize(event Event) Event {
	event.Name = strings.ToLower(strings.TrimSpace(event.Name))
	event.BatchID = strings.TrimSpace(event.BatchID)
	event.Actor = strings.TrimSpace(event.Actor)
	event.Detail = strings.TrimSpace(event.Detail)
	if event.At == "" {
		event.At = time.Unix(0, 0).UTC().Format(time.RFC3339)
	}
	return event
}
func Valid(event Event) bool {
	return event.Name != "" && event.BatchID != "" && event.Actor != "" && event.At != ""
}
func Sort(events []Event) []Event {
	result := append([]Event(nil), events...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].At == result[j].At {
			return result[i].Name < result[j].Name
		}
		return result[i].At < result[j].At
	})
	return result
}
func FilterBatch(events []Event, batchID string) []Event {
	result := make([]Event, 0)
	for _, event := range events {
		if event.BatchID == batchID {
			result = append(result, event)
		}
	}
	return Sort(result)
}
func FilterActor(events []Event, actor string) []Event {
	result := make([]Event, 0)
	for _, event := range events {
		if event.Actor == actor {
			result = append(result, event)
		}
	}
	return Sort(result)
}
func FilterName(events []Event, name string) []Event {
	result := make([]Event, 0)
	for _, event := range events {
		if event.Name == name {
			result = append(result, event)
		}
	}
	return Sort(result)
}
func Summarize(events []Event) Summary {
	result := Summary{ByName: map[string]int{}, ByActor: map[string]int{}}
	for _, event := range events {
		result.Count++
		result.ByName[event.Name]++
		result.ByActor[event.Actor]++
		result.TotalValue += event.Value
	}
	return result
}
func Names(events []Event) []string {
	seen := map[string]bool{}
	values := make([]string, 0)
	for _, event := range events {
		if seen[event.Name] {
			continue
		}
		seen[event.Name] = true
		values = append(values, event.Name)
	}
	sort.Strings(values)
	return values
}
func Actors(events []Event) []string {
	seen := map[string]bool{}
	values := make([]string, 0)
	for _, event := range events {
		if seen[event.Actor] {
			continue
		}
		seen[event.Actor] = true
		values = append(values, event.Actor)
	}
	sort.Strings(values)
	return values
}
func Total(events []Event) int {
	result := 0
	for _, event := range events {
		result += event.Value
	}
	return result
}
func Latest(events []Event) (Event, bool) {
	sorted := Sort(events)
	if len(sorted) == 0 {
		return Event{}, false
	}
	return sorted[len(sorted)-1], true
}
func Merge(left, right []Event) []Event {
	result := append([]Event(nil), left...)
	result = append(result, right...)
	return Sort(result)
}
func Limit(events []Event, count int) []Event {
	if count < 0 {
		count = 0
	}
	sorted := Sort(events)
	if len(sorted) > count {
		return sorted[:count]
	}
	return sorted
}
func Since(events []Event, at string) []Event {
	result := make([]Event, 0)
	for _, event := range events {
		if event.At >= at {
			result = append(result, event)
		}
	}
	return Sort(result)
}
func NamesForBatch(events []Event, batchID string) []string {
	return Names(FilterBatch(events, batchID))
}
func ActorsForBatch(events []Event, batchID string) []string {
	return Actors(FilterBatch(events, batchID))
}
func Has(events []Event, name string) bool {
	for _, event := range events {
		if event.Name == name {
			return true
		}
	}
	return false
}
func CountName(events []Event, name string) int {
	result := 0
	for _, event := range events {
		if event.Name == name {
			result++
		}
	}
	return result
}
func CountActor(events []Event, actor string) int {
	result := 0
	for _, event := range events {
		if event.Actor == actor {
			result++
		}
	}
	return result
}
