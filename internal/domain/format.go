package domain

import (
	"encoding/json"
	"sort"
)

func EncodeBatch(batch Batch) ([]byte, error) { return json.Marshal(batch) }
func DecodeBatch(data []byte) (Batch, error) {
	var value Batch
	err := json.Unmarshal(data, &value)
	return value, err
}
func EncodeRecord(record Record) ([]byte, error) { return json.Marshal(record) }
func DecodeRecord(data []byte) (Record, error) {
	var value Record
	err := json.Unmarshal(data, &value)
	return value, err
}
func EncodeAudit(event AuditEvent) ([]byte, error) { return json.Marshal(event) }
func DecodeAudit(data []byte) (AuditEvent, error) {
	var value AuditEvent
	err := json.Unmarshal(data, &value)
	return value, err
}
func EncodeNote(note CollaborationNote) ([]byte, error) { return json.Marshal(note) }
func DecodeNote(data []byte) (CollaborationNote, error) {
	var value CollaborationNote
	err := json.Unmarshal(data, &value)
	return value, err
}
func EncodeSnapshot(snapshot ExportSnapshot) ([]byte, error) { return json.Marshal(snapshot) }
func DecodeSnapshot(data []byte) (ExportSnapshot, error) {
	var value ExportSnapshot
	err := json.Unmarshal(data, &value)
	return value, err
}

func SortRecords(records []Record) []Record {
	copyOf := append([]Record(nil), records...)
	sort.Slice(copyOf, func(i, j int) bool { return copyOf[i].ID < copyOf[j].ID })
	return copyOf
}
