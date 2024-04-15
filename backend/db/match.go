package db

import (
	"context"
	"fmt"
	"time"

	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMatchesFromSeason(seasonID string, onlyPending bool) ([]domain.Match, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	dbSeasonID, err := primitive.ObjectIDFromHex(seasonID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	findCriteria := bson.M{
		"season_id": dbSeasonID,
	}
	if onlyPending {
		findCriteria["completed"] = false
	}
	// Find matches from this season
	cursor, err := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_MATCHES).
		Find(ctx,
			findCriteria,
		)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode matches
	var matches []domain.Match
	err = cursor.All(ctx, &matches)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return matches, nil
}
func GetMatchesFromPlayer(playerID string, onlyPending bool, count, page int) ([]domain.Match, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	dbPlayerID, err := primitive.ObjectIDFromHex(playerID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	// Find matches from this player
	findCriteria := bson.M{
		"players_data": bson.M{"$elemMatch": bson.M{"tournament_player_id": dbPlayerID}},
	}
	if onlyPending {
		findCriteria["completed"] = false
	}
	opts := options.Find().SetSkip(int64(count * (page - 1))).SetLimit(int64(count))
	// Find matches from this player
	cursor, err := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_MATCHES).
		Find(ctx,
			findCriteria,
			opts,
		)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode matches
	var matches []domain.Match
	err = cursor.All(ctx, &matches)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return matches, nil
}

func CreateMatch(seasonID string, match domain.Match) error {
	if match.ID != primitive.NilObjectID {
		return ErrObjectIDProvided
	}
	dbSeasonID, err := primitive.ObjectIDFromHex(seasonID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	match.ID = primitive.NewObjectID()
	match.SeasonID = dbSeasonID

	for i := range match.PlayersData {
		match.PlayersData[i].Wins = 0
	}
	match.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	match.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
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
		resultInsert, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_MATCHES).
			InsertOne(ctx, match)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return resultInsert, nil
	})

	log.Debug().Str("match_id", match.ID.String()).Msg("created match")

	return err
}

func UpdateMatch(matchID string, playerWins map[string]int, gamesPlayed int, completed bool) error {
	log.Debug().Interface("wins", playerWins).Str("match_id", matchID).Int("games", gamesPlayed).Send()
	dbMatchID, err := primitive.ObjectIDFromHex(matchID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

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
		// Find match to update
		result := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_MATCHES).
			FindOne(ctx,
				bson.M{"_id": dbMatchID},
			)

		if err := result.Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
			}
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}

		// Decode match
		var match *domain.Match
		err = result.Decode(&match)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}

		for tournament_player_id, wins := range playerWins {
			for index, player_data := range match.PlayersData {
				if player_data.TournamentPlayerID.Hex() == tournament_player_id {
					match.PlayersData[index].Wins = wins
				}
			}
		}
		match.GamesPlayed = gamesPlayed
		match.Completed = completed

		resultInsert, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_MATCHES).
			UpdateByID(ctx, match.ID, bson.M{"$set": match})
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return resultInsert, nil
	})

	return err
}
