package db

import (
	"context"
	"fmt"
	"time"

	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetPackBySetCode(setCode string) (*domain.BoosterPack, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
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

func CreateBoosterPack(boosterPack domain.BoosterPack) error {
	if boosterPack.ID != primitive.NilObjectID {
		return ErrObjectIDProvided
	}
	boosterPack.ID = primitive.NewObjectID()
	boosterPack.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	boosterPack.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer session.EndSession(ctx)

	// Find if user exists and if not, create it
	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		resultFind := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_BOOSTER_PACKS).
			FindOne(ctx, bson.M{"set_code": boosterPack.SetCode})
		if err := resultFind.Err(); err != mongo.ErrNoDocuments {
			if err == nil {
				return nil, fmt.Errorf("%w", ErrAlreadyExists)
			}
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}

		resultInsert, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_BOOSTER_PACKS).
			InsertOne(ctx,
				boosterPack,
			)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return resultInsert, nil
	})

	return err
}

func UpdateBoosterPack(boosterPack domain.BoosterPack) error {
	if boosterPack.ID != primitive.NilObjectID {
		return ErrObjectIDProvided
	}
	boosterPack.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		resultInsert, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_BOOSTER_PACKS).
			UpdateOne(ctx,
				bson.M{
					"set_code": boosterPack.SetCode,
				}, bson.M{
					"$set": bson.M{
						"set_code":    boosterPack.SetCode,
						"name":        boosterPack.Name,
						"description": boosterPack.Description,
						"card_count":  boosterPack.CardCount,
						"slots":       boosterPack.Slots,
						"filter":      boosterPack.Filter,
						"updated_at":  primitive.NewDateTimeFromTime(time.Now()),
					},
				})
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return resultInsert, nil
	})

	return err
}

func GetAllBoosterPacks() ([]domain.BoosterPack, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
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

func GetBoosterPackByID(boosterPackID string) (*domain.BoosterPack, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	dbBoosterPackID, err := primitive.ObjectIDFromHex(boosterPackID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	result := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_BOOSTER_PACKS).
		FindOne(ctx, bson.M{"_id": dbBoosterPackID})

	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode booster pack
	var boosterPack *domain.BoosterPack
	err = result.Decode(&boosterPack)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return boosterPack, nil
}

func BuyBoosterPack(tournamentID, userID, boosterPackID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	dbBoosterPackID, err := primitive.ObjectIDFromHex(boosterPackID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	dbUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		result := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_BOOSTER_PACKS).
			FindOne(ctx, bson.M{"_id": dbBoosterPackID})

		if err := result.Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
			}
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		// Decode booster pack
		var boosterPack *domain.BoosterPack
		err = result.Decode(&boosterPack)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}

		result = MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENTS).
			FindOne(ctx, bson.M{"_id": dbTournamentID})

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

		result = MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENT_PLAYERS).
			FindOne(ctx, bson.M{"user_id": dbUserID, "tournament_id": dbTournamentID})

		if err := result.Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
			}
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		// Decode tournament player
		var tournamentPlayer *domain.TournamentPlayer
		err = result.Decode(&tournamentPlayer)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}

		// Check and substract coins
		var foundStoreBoosterPack domain.StoreBoosterPack
		found := false
		for _, storeBoosterPack := range tournament.Store.BoosterPacks {
			if storeBoosterPack.BoosterPackID == dbBoosterPackID {
				foundStoreBoosterPack = storeBoosterPack
				found = true
			}
		}
		if !found {
			return nil, ErrNotFound
		}
		if tournamentPlayer.GameResources.Coins < foundStoreBoosterPack.CoinPrice {
			return nil, ErrBadRequest
		}
		tournamentPlayer.GameResources.Coins -= foundStoreBoosterPack.CoinPrice

		// Packs the user already has
		seenPacks := make(map[string]int, len(tournamentPlayer.GameResources.BoosterPacks))
		for index, pack := range tournamentPlayer.GameResources.BoosterPacks {
			seenPacks[pack.SetCode] = index
		}
		if index, ok := seenPacks[boosterPack.SetCode]; ok {
			tournamentPlayer.GameResources.BoosterPacks[index].Available += 1
		} else {
			tournamentPlayer.GameResources.BoosterPacks = append(tournamentPlayer.GameResources.BoosterPacks, domain.OwnedBoosterPack{
				Available:   1,
				SetCode:     boosterPack.SetCode,
				Name:        boosterPack.Name,
				Description: boosterPack.Description,
			})
		}
		// Update the tournament player
		updateResult, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENT_PLAYERS).
			UpdateByID(ctx, tournamentPlayer.ID, bson.M{"$set": tournamentPlayer})

		if err != nil || updateResult.MatchedCount == 0 {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}

		return nil, nil
	})

	return err
}
