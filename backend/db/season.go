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
	season.Matches = make([]domain.Match, 0)
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

func CreateMatch(seasonID, tournamentID, playerAID, playerBID string) error {

	season, err := GetSeasonByID(seasonID)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInternal, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer session.EndSession(ctx)

	// Find if active match exists
	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		for _, match := range season.Matches {
			if !match.Completed && ((playerAID == match.PlayerAID.Hex() && playerBID == match.PlayerBID.Hex()) ||
				(playerBID == match.PlayerAID.Hex() && playerAID == match.PlayerBID.Hex())) {
				return nil, fmt.Errorf("%w: %s", ErrInternal, "specified players have a match pending")
			}
		}
		dbPlAID, err := primitive.ObjectIDFromHex(playerAID)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInternal, err)
		}
		dbPlBID, err := primitive.ObjectIDFromHex(playerBID)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInternal, err)
		}

		season.Matches = append(season.Matches, domain.Match{
			PlayerAID:   dbPlAID,
			PlayerBID:   dbPlBID,
			PlayerAWins: 0,
			PlayerBWins: 0,
			Completed:   false,
		})

		updateResult, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_SEASONS).
			UpdateByID(ctx, season.ID, bson.M{"$set": season})

		if err != nil || updateResult.MatchedCount == 0 {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return nil, nil
	})

	return err
}

func UpdateMatch(seasonID string, newMatch domain.Match) error {
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
		season, err := GetSeasonByID(seasonID)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		found := false
		season.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
		newMatches := make([]domain.Match, 0, len(season.Matches))
		for _, match := range season.Matches {
			if !match.Completed {
				if match.PlayerAID == newMatch.PlayerAID && match.PlayerBID == newMatch.PlayerBID {
					match = newMatch
					found = true
				}
			}
			newMatches = append(newMatches, match)
		}
		if !found {
			return nil, fmt.Errorf("%w: %s", ErrNotFound, "pending match not found")
		}
		season.Matches = newMatches
		updateResult, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_SEASONS).
			UpdateByID(ctx, season.ID, bson.M{"$set": season})

		if err != nil || updateResult.MatchedCount == 0 {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return nil, nil
	})

	return err
}

func GetMatchesFromSeason(seasonID string, onlyPending bool) ([]domain.Match, error) {
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

	if onlyPending { // Filter active matches
		pendingMatches := make([]domain.Match, 0)
		for _, match := range season.Matches {
			if !match.Completed {
				pendingMatches = append(pendingMatches, match)
			}
		}

		return pendingMatches, nil
	}
	return season.Matches, err
}

func GetMatchesFromPlayer(tournamentID, playerID string, onlyPending bool) ([]domain.Match, error) {
	seasons, err := GetAllSeasons(tournamentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	dbPlayerID, err := primitive.ObjectIDFromHex(playerID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	playerMatches := make([]domain.Match, 0)
	for _, season := range seasons {
		for _, match := range season.Matches {
			if (dbPlayerID == match.PlayerAID) || (dbPlayerID == match.PlayerBID) {
				if onlyPending {
					if !match.Completed {
						playerMatches = append(playerMatches, match)
					}
				} else {
					playerMatches = append(playerMatches, match)
				}
			}
		}
	}
	return playerMatches, nil
}
