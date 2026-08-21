package archive

import (
	"encoding/json"
	"fmt"
	"sort"

	"inventoryseal/internal/crypto"
	"inventoryseal/internal/domain"
)

type Manifest struct {
	BatchID string             `json:"batch_id"`
	Status  domain.BatchStatus `json:"status"`
	Records []Entry            `json:"records"`
	Digest  string             `json:"digest"`
	Notes   int                `json:"notes"`
	Audits  int                `json:"audits"`
}
type Entry struct {
	ID     string `json:"id"`
	Result string `json:"result"`
	Count  int    `json:"count"`
}

func BuildManifest(batch domain.Batch, records []domain.Record, notes []domain.CollaborationNote, audits []domain.AuditEvent) Manifest {
	entries := make([]Entry, 0, len(records))
	for _, record := range records {
		entries = append(entries, Entry{ID: record.ID, Result: record.Result, Count: record.Observed})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })
	value := Manifest{BatchID: batch.ID, Status: batch.Status, Records: entries, Notes: len(notes), Audits: len(audits)}
	data, _ := json.Marshal(value)
	value.Digest = crypto.Digest(data)
	return value
}

func EncodeManifest(value Manifest) ([]byte, error) { return json.MarshalIndent(value, "", "  ") }
func DecodeManifest(data []byte) (Manifest, error) {
	var value Manifest
	err := json.Unmarshal(data, &value)
	return value, err
}
func Filename(batchID string, sequence int) string {
	return fmt.Sprintf("%s-export-%03d.json", batchID, sequence)
}
func VerifyManifest(value Manifest) bool {
	copy := value
	copy.Digest = ""
	data, _ := json.Marshal(copy)
	return value.Digest == crypto.Digest(data)
}
