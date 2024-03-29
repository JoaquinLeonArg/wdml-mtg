package boosterpacks

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	"github.com/rs/zerolog/log"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/event_log").Subrouter()
	r.HandleFunc("", GetEventLogsHandler).Methods(http.MethodGet)
	r.HandleFunc("", AddEventLogHandler).Methods(http.MethodPost)
}

//
// ENDPOINT: Get the last N event logs from a tournament
//

type GetEventLogsResponse struct {
	EventLogs []domain.EventLog `json:"event_logs"`
}

func GetEventLogsHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Get last 30 event logs for this tournament
	eventLogs, err := GetEventLogs(tournamentID, 30)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get event logs")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetEventLogsResponse{EventLogs: eventLogs}))
}

//
// ENDPOINT: Add or update an event log
//

type AddEventLogRequest struct {
	EventLogType domain.EventLogType `json:"event_log_type"`
	EventLogData interface{}         `json:"event_log_data"`
}

type AddEventLogResponse struct{}

func AddEventLogHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get user ID from request context
	ownerID, ok := r.Context().Value("user_id").(string)
	if ownerID == "" || !ok {
		log.Debug().
			Msg("failed to read user id from context")
		http.Error(w, "", http.StatusForbidden)
		return
	}

	// Decode body data
	var req AddEventLogRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Add or update an event log, depending on the type and content
	err = AddOrUpdateEventLog(tournamentID, req.EventLogType, req.EventLogData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(AddEventLogResponse{}))
}
