package server

import (
	"encoding/json"
	"epicwine/pkg/wine"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/matryer/respond.v1"
)

// status endpoint
func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	err := s.Db.VerifyCsv()
	if err != nil {
		data := map[string]string{"status": "error", "msg": err.Error(), "ts": time.Now().Format(time.RFC3339)}
		s.Metric.AddErrors()
		respond.With(w, r, http.StatusFailedDependency, data)
		return
	}

	data := map[string]string{"status": "ok", "ts": time.Now().Format(time.RFC3339)}
	s.Metric.AddSuccess()
	respond.With(w, r, http.StatusOK, data)
}

// list wine endpoint
func (s *Server) listwine(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limit, _ := strconv.Atoi(queryParams.Get("limit"))
	offset, _ := strconv.Atoi(queryParams.Get("offset"))

	winelist, err := s.Db.GetRecords(limit, offset)
	if err != nil {
		s.Metric.AddErrors()
		respond.WithStatus(w, r, http.StatusInternalServerError)
		return
	}
	s.Metric.AddSuccess()
	respond.With(w, r, http.StatusOK, winelist)
}

// put wine endpoint
func (s *Server) putwine(w http.ResponseWriter, r *http.Request) {
	var wreq wine.WineRequestInput
	json.NewDecoder(r.Body).Decode(&wreq)
	// Previous ID
	count, err := s.Db.GetCount()
	if err != nil {
		s.Metric.AddErrors()
		respond.With(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Add One for ID
	wreq.ID = count
	err = s.Db.WriteRecord(wreq)
	if err != nil {
		s.Metric.AddErrors()
		respond.With(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Verify Winelist Count
	wc, _ := s.Db.GetCount()
	s.Metric.SetWineCount(wc)

	s.Metric.AddSuccess()
	respond.WithStatus(w, r, http.StatusCreated)
}

// get wine endpoint
func (s *Server) getwine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	wine, err := s.Db.GetRecord(id)
	if err != nil {
		s.Metric.AddErrors()
		respond.With(w, r, http.StatusNotFound, err.Error())
		return
	}

	s.Metric.AddSuccess()
	respond.With(w, r, http.StatusOK, wine)
}
