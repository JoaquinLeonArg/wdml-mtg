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
)

const (
	DB_MAIN          = "wdml_main"
	COLLECTION_USERS = "users"
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
