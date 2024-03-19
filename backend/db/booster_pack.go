package db

import (
	"context"
	"fmt"
	"time"

	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPackBySetCode(setCode string) (*domain.BoosterPack, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Find pack
	result := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_BOOSTER_PACKS).
		FindOne(ctx,
			bson.M{"set_code": setCode},
		)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode pack
	var boosterPack *domain.BoosterPack
	err := result.Decode(&boosterPack)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return boosterPack, nil
}

type BoosterPackNoID struct {
	SetCode     string `bson:"set_code" json:"set_code"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	CardCount   int    `bson:"card_count" json:"card_count"`
	Filter      string `bson:"filter" json:"filter"`
	Slots       []struct {
		Options []domain.Option `bson:"options" json:"options"`
		Filter  string          `bson:"filter" json:"filter"`
		Count   int             `bson:"count" json:"count"`
	} `bson:"slots" json:"slots"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

func CreateBoosterPack(boosterPack domain.BoosterPack) error {
	if boosterPack.ID != primitive.NilObjectID {
		return ErrObjectIDProvided
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer session.EndSession(ctx)
	opts := options.Update().SetUpsert(true)
	// Find if user exists and if not, create it
	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		resultInsert, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_BOOSTER_PACKS).
			UpdateOne(ctx,
				bson.M{
					"set_code": boosterPack.SetCode,
				}, bson.M{
					"$set": BoosterPackNoID{
						SetCode:     boosterPack.SetCode,
						Name:        boosterPack.Name,
						Description: boosterPack.Description,
						CardCount:   boosterPack.CardCount,
						Slots:       boosterPack.Slots,
						CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
						UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
					},
				}, opts)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return resultInsert, nil
	})

	return err
}

func GetAllBoosterPacks() ([]domain.BoosterPack, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cursor, err := MongoDatabaseClient.Database(DB_MAIN).
		Collection(COLLECTION_BOOSTER_PACKS).
		Find(ctx, bson.D{})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	var boosterPacks []domain.BoosterPack
	err = cursor.All(ctx, &boosterPacks)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return boosterPacks, nil
}

// func DeleteAllBoosterPacks() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()
// 	_, err := MongoDatabaseClient.
// 		Database(DB_MAIN).
// 		Collection(COLLECTION_BOOSTER_PACKS).
// 		DeleteMany(ctx, bson.M{})
// 	return err
// }
