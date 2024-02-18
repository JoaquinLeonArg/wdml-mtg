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
	ErrInternal = fmt.Errorf("internal error")
)

const (
	DB_MAIN          = "wdml_main"
	COLLECTION_USERS = "users"
)

func InitDBConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", config.Config.MongoURL, config.Config.MongoPort)))
	if err != nil {
		return err
	}
	MongoDatabaseClient = client
	return nil
}
