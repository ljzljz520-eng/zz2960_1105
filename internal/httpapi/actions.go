package httpapi

import (
	"encoding/json"
	"inventoryseal/internal/domain"
	"net/http"
	"strings"
)

func (s *Server) action(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(strings.TrimPrefix(r.URL.Path, "/actions/"), "/"), "/")
	if len(parts) != 2 {
		writeJSON(w, 404, map[string]string{"error": "not found"})
		return
	}
	batchID, action := parts[0], parts[1]
	actor := r.URL.Query().Get("actor")
	if actor == "" {
		actor = "http"
	}
	var err error
	switch action {
	case "review":
		err = s.service.ReviewBatch(r.Context(), batchID, actor)
	case "archive":
		err = s.service.ArchiveBatch(r.Context(), batchID, actor)
	case "publish":
		_, err = s.service.ConfirmAndPublish(r.Context(), batchID, actor)
	default:
		writeJSON(w, 404, map[string]string{"error": "unknown action"})
		return
	}
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]string{"status": "ok"})
}

func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	batchID := strings.TrimPrefix(r.URL.Path, "/records/")
	if batchID == "" {
		writeJSON(w, 404, map[string]string{"error": "missing batch"})
		return
	}
	if r.Method == http.MethodGet {
		_, records, _, err := s.service.BatchDetails(r.Context(), batchID)
		if err != nil {
			writeJSON(w, 404, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, records)
		return
	}
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	var record domain.Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	record.BatchID = batchID
	if err := s.service.AddRecord(r.Context(), record); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 201, record)
}
