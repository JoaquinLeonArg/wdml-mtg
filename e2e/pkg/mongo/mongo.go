package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/joaquinleonarg/wdml-mtg/backend/db"
	"github.com/joaquinleonarg/wdml-mtg/e2e/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(
		ctx,
		options.Client().
			SetAuth(options.Credential{Username: config.Config.MongoUser, Password: config.Config.MongoPassword}).
			ApplyURI(fmt.Sprintf("mongodb+srv://%s", config.Config.MongoURL)),
	)
	if err != nil {
		panic(err)
	}
	err = client.Database(db.DB_MAIN).Drop(ctx)
	if err != nil {
		panic(err)
	}
}
