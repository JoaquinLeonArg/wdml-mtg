package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateTournament(tournament domain.Tournament) (primitive.ObjectID, error) {
	if tournament.ID != primitive.NilObjectID {
		return primitive.NilObjectID, ErrObjectIDProvided
	}
	tournament.ID = primitive.NewObjectID()
	tournament.InviteCode = uuid.New().String()
	tournament.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	tournament.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer session.EndSession(ctx)

	// Find if tournament exists, if owner exists, and if ok, create it
	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		resultFind := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENTS).
			FindOne(ctx, bson.M{"name": tournament.Name})
		if err := resultFind.Err(); err != mongo.ErrNoDocuments {
			if err == nil {
				return nil, fmt.Errorf("%w", ErrAlreadyExists)
			}
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}

		resultFind = MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_USERS).
			FindOne(ctx, bson.M{"_id": tournament.OwnerID})
		if err := resultFind.Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
			}
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}

		resultInsert, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENTS).
			InsertOne(ctx,
				tournament,
			)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return resultInsert, nil
	})

	return tournament.ID, err
}

func GetTournamentByID(tournamentID string) (*domain.Tournament, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
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
			return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode tournament
	var tournament *domain.Tournament
	err = result.Decode(&tournament)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return tournament, nil
}

func GetTournamentByInviteCode(inviteCode string) (*domain.Tournament, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	// Find tournament
	result := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_TOURNAMENTS).
		FindOne(ctx,
			bson.M{"invite_code": inviteCode},
		)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode tournament
	var tournament *domain.Tournament
	err := result.Decode(&tournament)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return tournament, nil
}

func GetTournamentsForUser(userID string) ([]domain.Tournament, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	tournament_players, err := GetTournamentPlayersForUser(userID)
	if err != nil {
		return nil, err
	}
	tournamentIDs := make([]primitive.ObjectID, 0, len(tournament_players))
	for _, tournament_player := range tournament_players {
		tournamentIDs = append(tournamentIDs, tournament_player.TournamentID)
	}

	// Find tournaments
	cursor, err := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_TOURNAMENTS).
		Find(ctx,
			bson.M{"_id": bson.M{"$in": tournamentIDs}},
		)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode tournaments
	var tournaments []domain.Tournament
	err = cursor.All(ctx, &tournaments)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return tournaments, nil
}
