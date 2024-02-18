package db

import (
	"context"
	"fmt"
	"time"

	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
)

func GetUserByUsername(username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Find user
	result := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_USERS).
		FindOne(ctx,
			domain.User{
				Username: username,
			},
		)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
		}
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	// Decode user
	var user *domain.User
	err := result.Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	return user, nil
}

func CreateUser(user domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}
	defer session.EndSession(context.TODO())
	err = session.StartTransaction()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}

	// Find if user exists and if not, create it
	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		resultFind := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_USERS).
			FindOne(ctx, bson.M{"$or": []domain.User{
				{Username: user.Username},
				{Email: user.Email},
			}})
		if resultFind.Err() != mongo.ErrNoDocuments {
			if err == nil {
				return nil, fmt.Errorf("%w", ErrUserAlreadyExists)
			}
			return nil, fmt.Errorf("%w: %w", ErrInternal, err)
		}

		resultInsert, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_USERS).
			InsertOne(ctx,
				user,
			)
		return resultInsert, fmt.Errorf("%w: %w", ErrInternal, err)
	})

	return err
}
