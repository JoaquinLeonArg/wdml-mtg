package boosterpacks

import (
	"errors"

	"github.com/joaquinleonarg/wdml-mtg/backend/db"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml-mtg/backend/errors"
)

func GetEventLogs(tournamentID string, count int) ([]domain.EventLog, error) {
	event_logs, err := db.GetEventLogs(tournamentID, count)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, apiErrors.ErrNotFound
		}
		return nil, apiErrors.ErrInternal
	}

	return event_logs, nil

}

func AddOrUpdateEventLog(tournamentID string, eventLogType domain.EventLogType, data interface{}) error {
	if eventLogType == domain.EventLogTypeOpenBoosters {
		if openBoostersData, ok := data.(domain.EventLogDataOpenBoosters); ok {
			return HandleOpenBoostersEventLog(tournamentID, openBoostersData)
		} else {
			return apiErrors.ErrBadRequest
		}
	} else {
		return apiErrors.ErrBadRequest
	}
}

func HandleOpenBoostersEventLog(tournamentID string, data domain.EventLogDataOpenBoosters) error {
	eventLogs, err := db.GetEventLogs(tournamentID, 1)
	if err != nil ||
		len(eventLogs) == 0 ||
		eventLogs[0].Type != domain.EventLogTypeOpenBoosters ||
		eventLogs[0].Data.(domain.EventLogDataOpenBoosters).SetName != data.SetName ||
		eventLogs[0].Data.(domain.EventLogDataOpenBoosters).Username != data.Username {
		err := db.AddEventLog(tournamentID, domain.EventLog{
			Type: domain.EventLogTypeOpenBoosters,
			Data: data,
		})
		if err != nil {
			return apiErrors.ErrInternal
		}
		return nil
	}

	eventLogData, ok := eventLogs[0].Data.(domain.EventLogDataOpenBoosters)
	if !ok {
		return apiErrors.ErrBadRequest
	}

	db.UpdateEventLog(domain.EventLog{
		ID:           eventLogs[0].ID,
		TournamentID: eventLogs[0].TournamentID,
		Type:         eventLogs[0].Type,
		Data: domain.EventLogDataOpenBoosters{
			Username: eventLogData.Username,
			SetName:  eventLogData.SetName,
			Count:    eventLogData.Count + data.Count,
		},
	})

	return nil
}
