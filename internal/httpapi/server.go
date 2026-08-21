package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"inventoryseal/internal/domain"
	"inventoryseal/internal/service"
)

type Server struct {
	service *service.Service
	mux     *http.ServeMux
}

func New(svc *service.Service) *Server {
	server := &Server{service: svc, mux: http.NewServeMux()}
	server.routes()
	return server
}
func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", s.health)
	s.mux.HandleFunc("/batches", s.batches)
	s.mux.HandleFunc("/batches/", s.batch)
	s.mux.HandleFunc("/actions/", s.action)
	s.mux.HandleFunc("/records/", s.records)
	s.mux.HandleFunc("/reports/", s.reconcile)
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) batches(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		items, err := s.service.FindBatches(r.Context(), r.URL.Query().Get("q"), domain.NormalizeStatus(r.URL.Query().Get("status")))
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, items)
		return
	}
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	var batch domain.Batch
	if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	if err := s.service.CreateBatch(r.Context(), batch); err != nil {
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 201, batch)
}

func (s *Server) batch(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(strings.TrimPrefix(r.URL.Path, "/batches/"), "/"), "/")
	if len(parts) != 1 || parts[0] == "" {
		writeJSON(w, 404, map[string]string{"error": "not found"})
		return
	}
	batch, records, notes, err := s.service.BatchDetails(r.Context(), parts[0])
	if err != nil {
		writeJSON(w, 404, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"batch": batch, "records": records, "notes": notes})
}
