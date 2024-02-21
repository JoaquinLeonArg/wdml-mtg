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

func CreateTournamentPlayer(tournamentPlayer domain.TournamentPlayer) error {
	if tournamentPlayer.ID != primitive.NilObjectID {
		return ErrObjectIDProvided
	}
	tournamentPlayer.ID = primitive.NewObjectID()
	tournamentPlayer.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	tournamentPlayer.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}
	defer session.EndSession(ctx)

	// Find if tournament player exists and if not, create it
	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		resultFind := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENT_PLAYERS).
			FindOne(ctx, bson.M{"tournament_id": tournamentPlayer.TournamentID, "user_id": tournamentPlayer.UserID})
		if err := resultFind.Err(); err != mongo.ErrNoDocuments {
			if err == nil {
				return nil, ErrAlreadyExists
			}
			return nil, fmt.Errorf("%w: %w", ErrInternal, err)
		}

		resultInsert, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENT_PLAYERS).
			InsertOne(ctx,
				tournamentPlayer,
			)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInternal, err)
		}
		return resultInsert, nil
	})

	return err
}

func GetTournamentPlayers(tournamentID string) ([]domain.TournamentPlayer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidID, err)
	}

	// Find players on this tournament
	cursor, err := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_TOURNAMENT_PLAYERS).
		Find(ctx,
			bson.M{"tournamentID": dbTournamentID},
		)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	// Decode tournament players
	var tournamentPlayers []domain.TournamentPlayer
	err = cursor.All(ctx, &tournamentPlayers)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return tournamentPlayers, nil
}

func GetAvailablePacksForTournamentPlayer(tournamentID, userID string) ([]domain.OwnedBoosterPack, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidID, err)
	}
	dbUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidID, err)
	}

	// Find packs for user
	result := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_TOURNAMENT_PLAYERS).
		FindOne(ctx,
			bson.M{"userID": dbUserID, "tournamentID": dbTournamentID},
		)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
		}
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	// Decode user
	var tournamentPlayer *domain.TournamentPlayer
	err = result.Decode(&tournamentPlayer)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return tournamentPlayer.GameResources.BoosterPacks, nil
}

func AddPacksToTournamentPlayer(tournamentID, userID string, packs []domain.OwnedBoosterPack) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidID, err)
	}
	dbUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidID, err)
	}

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}
	defer session.EndSession(ctx)

	// Find if user has packs of the same type and add them, or create new
	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		// Find tournament user
		result := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENT_PLAYERS).
			FindOne(ctx,
				bson.M{"userID": dbUserID, "tournamentID": dbTournamentID},
			)
		if err := result.Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("%w: %w", ErrNotFound, err)
			}
			return nil, fmt.Errorf("%w: %w", ErrInternal, err)
		}
		// Decode user
		var tournamentPlayer *domain.TournamentPlayer
		err = result.Decode(&tournamentPlayer)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInternal, err)
		}
		// Packs the user already has
		seenPacks := make(map[string]int, len(tournamentPlayer.GameResources.BoosterPacks))
		for index, pack := range tournamentPlayer.GameResources.BoosterPacks {
			seenPacks[pack.SetCode] = index
		}
		for _, pack := range packs {
			if index, ok := seenPacks[pack.SetCode]; ok {
				tournamentPlayer.GameResources.BoosterPacks[index].Available += pack.Available
			} else {
				tournamentPlayer.GameResources.BoosterPacks = append(tournamentPlayer.GameResources.BoosterPacks, pack)
			}
		}
		// Update the tournament player
		updateResult, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENT_PLAYERS).
			UpdateByID(ctx, dbUserID, tournamentPlayer)

		if err != nil || updateResult.MatchedCount == 0 {
			return nil, fmt.Errorf("%w: %w", ErrInternal, err)
		}
		return nil, nil
	})

	return err
}
