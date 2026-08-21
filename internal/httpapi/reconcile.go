package httpapi

import (
	"net/http"
	"strings"
)

func (s *Server) reconcile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	batchID := strings.Trim(strings.TrimPrefix(r.URL.Path, "/reports/"), "/")
	if batchID == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "missing batch"})
		return
	}
	report, err := s.service.ReconcileBatch(r.Context(), batchID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, report)
}
