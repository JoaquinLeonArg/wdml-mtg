package db

import (
	"context"
	"fmt"
	"time"

	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateTournament(tournament domain.Tournament) error {
	if tournament.ID != primitive.NilObjectID {
		return ErrObjectIDProvided
	}
	tournament.ID = primitive.NewObjectID()
	tournament.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	tournament.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}
	defer session.EndSession(ctx)

	// Find if tournament exists, if owner exists, and if ok, create it
	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		resultFind := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENTS).
			FindOne(ctx, bson.M{"title": tournament.Title})
		if err := resultFind.Err(); err != mongo.ErrNoDocuments {
			if err == nil {
				return nil, fmt.Errorf("%w", ErrAlreadyExists)
			}
			return nil, fmt.Errorf("%w: %w", ErrInternal, err)
		}

		resultFind = MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_USERS).
			FindOne(ctx, bson.M{"_id": tournament.OwnerID})
		if err := resultFind.Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return nil, fmt.Errorf("%w: %w", ErrInternal, err)
		}

		resultInsert, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENTS).
			InsertOne(ctx,
				tournament,
			)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInternal, err)
		}
		return resultInsert, nil
	})

	return err
}

func GetTournamentByID(tournamentID string) (*domain.Tournament, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidID, err)
	}

	// Find tournament
	result := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_TOURNAMENTS).
		FindOne(ctx,
			bson.M{"_id": dbTournamentID},
		)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	// Decode tournament
	var tournament *domain.Tournament
	err = result.Decode(&tournament)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return tournament, nil
}
