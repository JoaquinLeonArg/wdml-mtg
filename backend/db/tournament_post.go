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

func GetAllTournamentPosts(tournamentID string) ([]domain.TournamentPost, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	// Find all tournament posts
	cursor, err := MongoDatabaseClient.Database(DB_MAIN).
		Collection(COLLECTION_TOURNAMENT_POSTS).
		Find(ctx, bson.M{
			"tournament_id": dbTournamentID,
		})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	var tournamentPosts []domain.TournamentPost
	err = cursor.All(ctx, &tournamentPosts)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return tournamentPosts, nil
}

func CreateTournamentPost(tournamentPost domain.TournamentPost, tournamentID string) error {
	if tournamentPost.ID != primitive.NilObjectID {
		return ErrObjectIDProvided
	}

	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	tournamentPost.ID = primitive.NewObjectID()
	tournamentPost.TournamentID = dbTournamentID
	tournamentPost.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	tournamentPost.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err = MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_TOURNAMENT_POSTS).
		InsertOne(ctx,
			tournamentPost,
		)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return err
}

func DeleteTournamentPost(tournamentID, tournamentPostID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	dbTournamentPostID, err := primitive.ObjectIDFromHex(tournamentPostID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	// Delete the tournament post
	result, err := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_TOURNAMENT_POSTS).
		DeleteOne(ctx,
			bson.M{
				"_id":           dbTournamentPostID,
				"tournament_id": dbTournamentID,
			})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if result.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}
