package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cdevoogd/dashboard/backend/dash"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) handleGetAllApplications(w http.ResponseWriter, r *http.Request) {
	records, err := s.db.GetAllApplications()
	if err != nil {
		msg := fmt.Sprint("error querying all applications: ", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	apps := make([]*dash.ApplicationResponse, len(records))
	for i, record := range records {
		apps[i] = record.ToApplicationResponse()
	}

	s.writeJSON(w, apps)
}

func (s *Server) handleGetApplication(w http.ResponseWriter, r *http.Request) {
	appID := mux.Vars(r)["id"]
	if appID == "" {
		s.handleEmptyPathParameter(w, r)
		return
	}

	record, err := s.db.GetApplication(appID)
	if err != nil {
		if s.db.IsNotFoundError(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		msg := fmt.Sprintf("error querying for application %s: %s", appID, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, record.ToApplicationResponse())
}

func (s *Server) handleCreateApplication(w http.ResponseWriter, r *http.Request) {
	app := &dash.ApplicationRequest{}
	err := json.NewDecoder(r.Body).Decode(app)
	if err != nil {
		msg := fmt.Sprint("error decoding json: ", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	appID, err := uuid.NewRandom()
	if err != nil {
		msg := fmt.Sprint("error generating uuid: ", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	record := app.ToApplicationRecord(appID.String())
	err = s.db.AddApplication(record)
	if err != nil {
		msg := fmt.Sprint("error inserting application: ", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, record.ToApplicationResponse())
}

func (s *Server) handleUpdateApplication(w http.ResponseWriter, r *http.Request) {
	appID := mux.Vars(r)["id"]
	if appID == "" {
		s.handleEmptyPathParameter(w, r)
		return
	}

	app := &dash.ApplicationRequest{}
	err := json.NewDecoder(r.Body).Decode(app)
	if err != nil {
		msg := fmt.Sprint("error decoding json: ", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	record := app.ToApplicationRecord(appID)
	err = s.db.UpdateApplication(record)
	if err != nil {
		msg := fmt.Sprint("error updating application: ", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, record.ToApplicationResponse())
}

func (s *Server) handleDeleteApplication(w http.ResponseWriter, r *http.Request) {
	appID := mux.Vars(r)["id"]
	if appID == "" {
		s.handleEmptyPathParameter(w, r)
		return
	}

	err := s.db.DeleteApplication(appID)
	if err != nil {
		if s.db.IsNotFoundError(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		msg := fmt.Sprintf("error deleting application %s: %s", appID, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}
