package db

import (
	"context"
	"fmt"
	"time"

	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllSeasons(tournamentID string) ([]domain.Season, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	// Find seasons on this tournament
	cursor, err := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_SEASONS).
		Find(ctx,
			bson.M{"tournament_id": dbTournamentID},
		)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode seasons
	var seasons []domain.Season
	err = cursor.All(ctx, &seasons)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return seasons, nil
}

func GetSeasonByID(seasonID string) (*domain.Season, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	dbSeasonID, err := primitive.ObjectIDFromHex(seasonID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	// Find season
	result := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_SEASONS).
		FindOne(ctx,
			bson.M{"_id": dbSeasonID},
		)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode season
	var season *domain.Season
	err = result.Decode(&season)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return season, nil
}

func CreateEmptySeason(season domain.Season) error {
	if season.ID != primitive.NilObjectID {
		return ErrObjectIDProvided
	}
	season.ID = primitive.NewObjectID()
	season.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	season.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		// TODO: Check that season with the same name doesn't exist
		resultInsert, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_SEASONS).
			InsertOne(ctx, season)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return resultInsert, nil
	})

	log.Debug().Str("season_id", season.ID.String()).Msg("created season")

	return err
}
