package db

import (
	"context"
	"fmt"
	"time"

	"github.com/joaquinleonarg/wdml_mtg/backend/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDatabaseClient *mongo.Client

var (
	ErrInternal         = fmt.Errorf("internal error")
	ErrObjectIDProvided = fmt.Errorf("object id should not be provided")
	ErrInvalidID        = fmt.Errorf("invalid object id provided")
	ErrNotFound         = fmt.Errorf("not found")
	ErrAlreadyExists    = fmt.Errorf("already exists")
)

const (
	DB_MAIN                       = "wdml_main"
	COLLECTION_USERS              = "users"
	COLLECTION_TOURNAMENTS        = "tournaments"
	COLLECTION_TOURNAMENT_PLAYERS = "tournament_players"
	COLLECTION_CARD_COLLECTION    = "card_collection"
	COLLECTION_BOOSTER_PACKS      = "booster_packs"
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
