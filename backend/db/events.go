package db

import (
	"context"
	"fmt"
	"time"

	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetEventLogs(tournamentID string, count int) ([]domain.EventLog, error) {
	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	opts := options.Find().SetLimit(int64(count))
	// Get the event logs
	cursor, err := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_EVENT_LOGS).
		Find(ctx, bson.M{"tournament_id": dbTournamentID}, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode event logs
	var eventLogs []domain.EventLog
	err = cursor.All(ctx, &eventLogs)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return eventLogs, nil
}

func AddEventLog(tournamentID string, eventLog domain.EventLog) error {
	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	if eventLog.ID != primitive.NilObjectID {
		return ErrObjectIDProvided
	}
	eventLog.ID = primitive.NewObjectID()
	eventLog.TournamentID = dbTournamentID

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		// Insert the event
		resultInsert, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_EVENT_LOGS).
			InsertOne(ctx, eventLog)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return resultInsert, nil
	})

	return err
}

func UpdateEventLog(eventLog domain.EventLog) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		// Update the event
		resultInsert, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_EVENT_LOGS).
			UpdateByID(ctx, eventLog.ID, bson.M{"$set": eventLog})
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return resultInsert, nil
	})

	return err
}
