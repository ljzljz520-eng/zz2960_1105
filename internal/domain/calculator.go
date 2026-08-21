package domain

import "sort"

type Range struct {
	Minimum int
	Maximum int
}
type Measurement struct {
	RecordID  string
	Delta     int
	Absolute  int
	Direction string
}

func Difference(record Record) int { return record.Observed - record.Expected }

func MeasurementOf(record Record) Measurement {
	delta := Difference(record)
	direction := "equal"
	if delta > 0 {
		direction = "over"
	}
	if delta < 0 {
		direction = "under"
	}
	absolute := delta
	if absolute < 0 {
		absolute = -absolute
	}
	return Measurement{RecordID: record.ID, Delta: delta, Absolute: absolute, Direction: direction}
}

func Measurements(records []Record) []Measurement {
	result := make([]Measurement, 0, len(records))
	for _, record := range records {
		result = append(result, MeasurementOf(record))
	}
	sort.Slice(result, func(i, j int) bool { return result[i].RecordID < result[j].RecordID })
	return result
}

func Within(record Record, allowed Range) bool {
	delta := Difference(record)
	return delta >= allowed.Minimum && delta <= allowed.Maximum
}

func Outliers(records []Record, allowed Range) []Record {
	result := make([]Record, 0)
	for _, record := range records {
		if !Within(record, allowed) {
			result = append(result, record)
		}
	}
	return SortRecords(result)
}

func AverageDelta(records []Record) float64 {
	if len(records) == 0 {
		return 0
	}
	total := 0
	for _, record := range records {
		total += Difference(record)
	}
	return float64(total) / float64(len(records))
}
