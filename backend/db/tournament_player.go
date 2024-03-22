package db

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateTournamentPlayer(tournamentPlayer domain.TournamentPlayer) (primitive.ObjectID, error) {
	if tournamentPlayer.ID != primitive.NilObjectID {
		return primitive.NilObjectID, ErrObjectIDProvided
	}
	tournamentPlayer.ID = primitive.NewObjectID()
	tournamentPlayer.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	tournamentPlayer.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer session.EndSession(ctx)

	// Find if tournament player exists and if not, create it
	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		resultFind := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENT_PLAYERS).
			FindOne(ctx, bson.M{"tournament_id": tournamentPlayer.TournamentID, "user_id": tournamentPlayer.UserID})
		if err := resultFind.Err(); err != mongo.ErrNoDocuments {
			if err == nil {
				return nil, ErrAlreadyExists
			}
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}

		resultInsert, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENT_PLAYERS).
			InsertOne(ctx,
				tournamentPlayer,
			)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return resultInsert, nil
	})

	return tournamentPlayer.TournamentID, err
}

func GetTournamentPlayerByID(tournamentPlayerID string) (*domain.TournamentPlayer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbTournamentPlayerID, err := primitive.ObjectIDFromHex(tournamentPlayerID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	// Find players on this tournament
	result := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_TOURNAMENT_PLAYERS).
		FindOne(ctx,
			bson.M{"_id": dbTournamentPlayerID},
		)
	if err := result.Err(); err == mongo.ErrNoDocuments {
		if err == nil {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode tournament players
	var tournamentPlayer *domain.TournamentPlayer
	if err := result.Decode(&tournamentPlayer); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return tournamentPlayer, nil
}

func GetTournamentPlayers(tournamentID string) ([]domain.TournamentPlayer, []domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	// Find players on this tournament
	cursor, err := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_TOURNAMENT_PLAYERS).
		Find(ctx,
			bson.M{"tournament_id": dbTournamentID},
		)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode tournament players
	var tournamentPlayers []domain.TournamentPlayer
	err = cursor.All(ctx, &tournamentPlayers)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	userIDs := make([]primitive.ObjectID, 0, len(tournamentPlayers))
	for _, tournamentPlayer := range tournamentPlayers {
		userIDs = append(userIDs, tournamentPlayer.UserID)
	}

	// Find users
	cursor, err = MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_USERS).
		Find(ctx,
			bson.M{"_id": bson.M{"$in": userIDs}},
		)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode tournament players
	var users []domain.User
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return tournamentPlayers, users, nil
}

func GetTournamentPlayersForUser(userID string) ([]domain.TournamentPlayer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	// Find tournament players for this user
	cursor, err := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_TOURNAMENT_PLAYERS).
		Find(ctx,
			bson.M{"user_id": dbUserID},
		)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode tournament players
	var tournamentPlayers []domain.TournamentPlayer
	err = cursor.All(ctx, &tournamentPlayers)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return tournamentPlayers, nil
}

func GetTournamentPlayer(tournamentID, userID string) (*domain.TournamentPlayer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	dbUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	// Find packs for user
	result := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_TOURNAMENT_PLAYERS).
		FindOne(ctx,
			bson.M{"user_id": dbUserID, "tournament_id": dbTournamentID},
		)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode user
	var tournamentPlayer *domain.TournamentPlayer
	err = result.Decode(&tournamentPlayer)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return tournamentPlayer, nil
}

func GetAvailablePacksForTournamentPlayer(tournamentID, userID string) ([]domain.OwnedBoosterPack, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	dbUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	// Find packs for user
	result := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_TOURNAMENT_PLAYERS).
		FindOne(ctx,
			bson.M{"user_id": dbUserID, "tournament_id": dbTournamentID},
		)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode user
	var tournamentPlayer *domain.TournamentPlayer
	err = result.Decode(&tournamentPlayer)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return tournamentPlayer.GameResources.BoosterPacks, nil
}

func ConsumeBoosterPackForTournamentPlayer(userID, tournamentID string, setCode string, cards []domain.CardData) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	dbUserID, err := primitive.ObjectIDFromHex(userID)
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

	// Find if user has packs of the same type and add them, or create new
	_, err = session.WithTransaction(ctx, func(mongoCtx mongo.SessionContext) (interface{}, error) {
		// Find tournament user
		result := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENT_PLAYERS).
			FindOne(ctx,
				bson.M{"user_id": dbUserID, "tournament_id": dbTournamentID},
			)
		if err := result.Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
			}
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		// Decode user
		var tournamentPlayer *domain.TournamentPlayer
		err = result.Decode(&tournamentPlayer)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}

		removed := false
		newPacks := make([]domain.OwnedBoosterPack, 0, len(tournamentPlayer.GameResources.BoosterPacks))
		// Find and remove the booster pack
		for _, boosterPack := range tournamentPlayer.GameResources.BoosterPacks {
			if boosterPack.SetCode == setCode {
				if boosterPack.Available == 1 && !removed {
					removed = true
					continue
				}
				if boosterPack.Available > 1 && !removed {
					boosterPack.Available -= 1
					newPacks = append(newPacks, boosterPack)
					removed = true
					continue
				}
			}
			newPacks = append(newPacks, boosterPack)
		}
		if !removed {
			return nil, fmt.Errorf("%w: %s", ErrNotFound, "booster pack not available for tournament player")
		}
		tournamentPlayer.GameResources.BoosterPacks = newPacks

		// Update the tournament player
		updateResult, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENT_PLAYERS).
			UpdateByID(ctx, tournamentPlayer.ID, bson.M{"$set": tournamentPlayer})

		if err != nil || updateResult.MatchedCount == 0 {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}

		// Add the cards to the tournament player's collection
		// For each card, find if the user already has some of that card, and update or add it accordingly
		for _, card := range cards {
			result, err := MongoDatabaseClient.
				Database(DB_MAIN).
				Collection(COLLECTION_CARD_COLLECTION).
				Find(ctx,
					bson.M{
						"tournament_id":              tournamentPlayer.TournamentID,
						"user_id":                    tournamentPlayer.UserID,
						"card_data.set_code":         card.SetCode,
						"card_data.collector_number": card.CollectorNumber,
					},
				)
			if err != nil {
				return nil, fmt.Errorf("%w: %v", ErrInternal, err)
			}
			var foundCards []domain.OwnedCard
			if err := result.All(ctx, &foundCards); err != nil {
				return nil, fmt.Errorf("%w: %v", ErrInternal, err)
			}
			if len(foundCards) == 0 {
				// Add a new card
				_, err := MongoDatabaseClient.
					Database(DB_MAIN).
					Collection(COLLECTION_CARD_COLLECTION).
					InsertOne(ctx, domain.OwnedCard{
						ID:           primitive.NewObjectID(),
						TournamentID: tournamentPlayer.TournamentID,
						UserID:       tournamentPlayer.UserID,
						Tags:         []string{},
						Count:        1,
						CardData:     card,
						CreatedAt:    primitive.NewDateTimeFromTime(time.Now()),
						UpdatedAt:    primitive.NewDateTimeFromTime(time.Now()),
					})

				if err != nil {
					return nil, fmt.Errorf("%w: %v", ErrInternal, err)
				}
			} else if len(foundCards) == 1 {
				// Update count of existing card
				if foundCards[0].Count >= 4 && !slices.Contains(foundCards[0].CardData.Types, "Basic") {
					switch foundCards[0].CardData.Rarity {
					case "mythic":
						AddCoinsToTournamentPlayer(domain.MYTHIC_TO_COIN, userID, tournamentID)
					case "rare":
						AddCoinsToTournamentPlayer(domain.RARE_TO_COIN, userID, tournamentID)
					case "uncommon":
						AddCoinsToTournamentPlayer(domain.UNCOMMON_TO_COIN, userID, tournamentID)
					case "common":
						AddCoinsToTournamentPlayer(domain.COMMON_TO_COIN, userID, tournamentID)
					}
				} else {
					foundCards[0].Count += 1
				}
				foundCards[0].UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
				_, err := MongoDatabaseClient.
					Database(DB_MAIN).
					Collection(COLLECTION_CARD_COLLECTION).
					UpdateByID(ctx, foundCards[0].ID, bson.M{"$set": foundCards[0]})
				if err != nil {
					return nil, fmt.Errorf("%w: %v", ErrInternal, err)
				}

			} else {
				// TODO: Consolidate duplicate entries just in case
				return nil, fmt.Errorf("%w: duplicated entries for card found on database", ErrInternal)
			}
		}

		return nil, nil
	})

	return err
}

func AddPacksToTournamentPlayer(tournamentPlayerID string, pack domain.OwnedBoosterPack) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbTournamentPlayerID, err := primitive.ObjectIDFromHex(tournamentPlayerID)
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

	// Find if user has packs of the same type and add them, or create new
	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		// Find tournament user
		result := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENT_PLAYERS).
			FindOne(ctx,
				bson.M{"_id": dbTournamentPlayerID},
			)
		if err := result.Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
			}
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		// Decode user
		var tournamentPlayer *domain.TournamentPlayer
		err = result.Decode(&tournamentPlayer)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		// Packs the user already has
		seenPacks := make(map[string]int, len(tournamentPlayer.GameResources.BoosterPacks))
		for index, pack := range tournamentPlayer.GameResources.BoosterPacks {
			seenPacks[pack.SetCode] = index
		}
		if index, ok := seenPacks[pack.SetCode]; ok {
			tournamentPlayer.GameResources.BoosterPacks[index].Available += pack.Available
		} else {
			tournamentPlayer.GameResources.BoosterPacks = append(tournamentPlayer.GameResources.BoosterPacks, pack)
		}
		// Update the tournament player
		updateResult, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENT_PLAYERS).
			UpdateByID(ctx, dbTournamentPlayerID, bson.M{"$set": tournamentPlayer})

		if err != nil || updateResult.MatchedCount == 0 {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return nil, nil
	})

	return err
}

func AddCoinsToTournamentPlayer(coins int, userID, tournamentID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	dbUserID, err := primitive.ObjectIDFromHex(userID)
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

	_, err = session.WithTransaction(ctx, func(mongoCtx mongo.SessionContext) (interface{}, error) {
		// Find tournament user
		result := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENT_PLAYERS).
			FindOne(ctx,
				bson.M{"user_id": dbUserID, "tournament_id": dbTournamentID},
			)
		if err := result.Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
			}
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		// Decode user
		var tournamentPlayer *domain.TournamentPlayer
		err = result.Decode(&tournamentPlayer)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		tournamentPlayer.GameResources.Coins += coins
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

func AddPointsToTournamentPlayer(points int, userID, tournamentID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	dbUserID, err := primitive.ObjectIDFromHex(userID)
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

	_, err = session.WithTransaction(ctx, func(mongoCtx mongo.SessionContext) (interface{}, error) {
		// Find tournament player
		result := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_TOURNAMENT_PLAYERS).
			FindOne(ctx,
				bson.M{"user_id": dbUserID, "tournament_id": dbTournamentID},
			)
		if err := result.Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
			}
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		// Decode user
		var tournamentPlayer *domain.TournamentPlayer
		err = result.Decode(&tournamentPlayer)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		tournamentPlayer.TournamentPoints += points
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
