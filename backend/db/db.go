package db

import (
	"context"
	"fmt"
	"time"

	"github.com/joaquinleonarg/wdml-mtg/backend/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDatabaseClient *mongo.Client

var (
	ErrInternal         = fmt.Errorf("internal error: %w", mongo.ErrNilValue)
	ErrObjectIDProvided = fmt.Errorf("object id should not be provided: %w", mongo.ErrNilDocument)
	ErrInvalidID        = fmt.Errorf("invalid object id provided: %w", mongo.ErrNilValue)
	ErrNotFound         = fmt.Errorf("not found: %w", mongo.ErrNoDocuments)
	ErrAlreadyExists    = fmt.Errorf("already exists: %w", mongo.ErrEmptySlice)

	ErrUninitialized = fmt.Errorf("uninitialized field: %w", mongo.ErrNilValue)
)

const (
	DB_MAIN                       = "wdml_main"
	COLLECTION_USERS              = "users"
	COLLECTION_TOURNAMENTS        = "tournaments"
	COLLECTION_TOURNAMENT_PLAYERS = "tournament_players"
	COLLECTION_TOURNAMENT_POSTS   = "tournament_posts"
	COLLECTION_CARD_COLLECTION    = "card_collection"
	COLLECTION_BOOSTER_PACKS      = "booster_packs"
	COLLECTION_DECKS              = "decks"
	COLLECTION_SEASONS            = "seasons"
	COLLECTION_MATCHES            = "matches"
	COLLECTION_EVENT_LOGS         = "event_logs"
)

func InitDBConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(
		ctx,
		options.Client().
			SetAuth(options.Credential{Username: config.Config.MongoUser, Password: config.Config.MongoPassword}).
			ApplyURI(fmt.Sprintf("mongodb+srv://%s", config.Config.MongoURL)),
	)
	if err != nil {
		return err
	}
	MongoDatabaseClient = client
	return nil
}
