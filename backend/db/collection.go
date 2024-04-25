package db

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CardFilter struct {
	Type      CardFilterType
	Operation CardFilterOperation
	Value     string
}

type CardFilterOperation string

const (
	CardFilterOperationEq CardFilterOperation = "="
	CardFilterOperationLt CardFilterOperation = "<"
	CardFilterOperationGt CardFilterOperation = ">"
)

type CardFilterType string

const (
	CardFilterTypeName    CardFilterType = "name"
	CardFilterTypeTags    CardFilterType = "tags"
	CardFilterTypeRarity  CardFilterType = "rarity"
	CardFilterTypeColor   CardFilterType = "color"
	CardFilterTypeTypes   CardFilterType = "types"
	CardFilterTypeOracle  CardFilterType = "oracle"
	CardFilterTypeSetCode CardFilterType = "setcode"
	CardFilterTypeMV      CardFilterType = "mv"
)

func GetCardsFromTournamentPlayer(userID, tournamentID string, filters []CardFilter, count, page int) ([]domain.OwnedCard, int, error) {
	log.Debug().Interface("filters", filters).Int("count", count).Int("page", page).Send()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	dbUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	dbFilters := bson.M{}

	for _, filter := range filters {
		switch filter.Type {
		// By name
		case CardFilterTypeName:
			switch filter.Operation {
			case CardFilterOperationEq:
				dbFilters["card_data.name"] = bson.M{"$regex": filter.Value, "$options": "i"}
			}

		// By tags
		case CardFilterTypeTags:
			yesTags := []string{}
			noTags := []string{}
			for _, v := range strings.Split(filter.Value, " ") {
				if strings.HasPrefix(v, "-") {
					noTags = append(noTags, strings.TrimPrefix(v, "-"))
				} else {
					yesTags = append(yesTags, v)
				}
			}
			switch filter.Operation {
			case CardFilterOperationEq:
				dbFilters["card_data.tags"] = bson.M{"$all": yesTags, "$nin": noTags}
			}

		// By rarity
		case CardFilterTypeRarity:
			switch filter.Operation {
			case CardFilterOperationEq:
				dbFilters["card_data.rarity"] = filter.Value
			}

		// By color
		case CardFilterTypeColor:
			switch filter.Operation {
			case CardFilterOperationLt:
				dbFilters["card_data.colors"] = filter.Value
			case CardFilterOperationEq:
				if filter.Value == "C" {
					dbFilters["card_data.colors"] = bson.M{"$size": 0}
				} else {
					dbFilters["card_data.colors"] = bson.M{"$all": strings.Split(filter.Value, "")}
				}
			}

		// By types
		case CardFilterTypeTypes:
			switch filter.Operation {
			case CardFilterOperationEq:
				dbFilters["card_data.types"] = bson.M{"$all": strings.Split(filter.Value, " ")}
			}

		// By oracle
		case CardFilterTypeOracle:
			switch filter.Operation {
			case CardFilterOperationEq:
				dbFilters["card_data.oracle"] = bson.M{"$regex": filter.Value, "$options": "i"}
			}

		// By set code
		case CardFilterTypeSetCode:
			switch filter.Operation {
			case CardFilterOperationEq:
				dbFilters["card_data.set_code"] = filter.Value
			}

		// By mv
		case CardFilterTypeMV:
			mv, err := strconv.Atoi(filter.Value)
			if err != nil {
				continue
			}
			switch filter.Operation {
			case CardFilterOperationEq:
				dbFilters["card_data.mana_value"] = mv
			}

		}
	}

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer session.EndSession(ctx)

	// Find the user cards that satisfy the filters
	res, err := session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		// Find tournament user
		filter := bson.M{
			"tournament_id": dbTournamentID,
			"user_id":       dbUserID,
		}
		for filterKey, filterValue := range dbFilters {
			filter[filterKey] = filterValue
		}
		log.Debug().Interface("filter", filter).Send()
		cursor, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_CARD_COLLECTION).
			Aggregate(ctx,
				bson.A{
					bson.M{"$match": filter},
					bson.M{"$skip": count * (page - 1)},
					bson.M{"$limit": count},
				},
			)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		count, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_CARD_COLLECTION).
			CountDocuments(ctx, filter)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}

		// Decode cards
		var queryResult []domain.OwnedCard
		err = cursor.All(ctx, &queryResult)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}

		return map[string]interface{}{"count": int(count), "cards": queryResult}, nil
	})

	if err != nil {
		return nil, 0, err
	}

	return res.(map[string]interface{})["cards"].([]domain.OwnedCard), res.(map[string]interface{})["count"].(int), nil
}

func GetOwnedCardById(cardId string) (domain.OwnedCard, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	dbCardId, err := primitive.ObjectIDFromHex(cardId)
	if err != nil {
		return domain.OwnedCard{}, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	// Find card
	result := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_CARD_COLLECTION).
		FindOne(ctx,
			bson.M{"_id": dbCardId},
		)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.OwnedCard{}, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return domain.OwnedCard{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode card
	var card domain.OwnedCard
	err = result.Decode(&card)
	if err != nil {
		return domain.OwnedCard{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return card, nil
}

type CardBySetNum struct {
	Set, Num string
}

func ImportCollection(cards []domain.OwnedCard) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		newValue := make([]interface{}, len(cards))

		for i := range cards {
			newValue[i] = cards[i]
		}

		result, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_CARD_COLLECTION).InsertMany(ctx, newValue)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return domain.OwnedCard{}, fmt.Errorf("%w: %v", ErrNotFound, err)
			}
			return domain.OwnedCard{}, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return result, nil
	})

	return err
}

func UpdateOwnedCard(ownedCard domain.OwnedCard) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	result, err := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_BOOSTER_PACKS).
		UpdateByID(ctx, ownedCard.ID, ownedCard)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("%w", ErrNotFound)
	}
	return nil
}

func TradeUpCards(cardsToRemove map[string]int, cardsToAdd []domain.CardData, tournamentID, ownerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(mongoCtx mongo.SessionContext) (interface{}, error) {
		tournamentPlayer, err := GetTournamentPlayer(tournamentID, ownerID)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
		}
		err = RemoveCardsFromTournamentPlayer(tournamentPlayer.ID.Hex(), cardsToRemove)
		if err != nil {
			return nil, err
		}
		err = AddCardsToTournamentPlayer(tournamentPlayer.ID.Hex(), cardsToAdd)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

func AddCardsToTournamentPlayer(tournamentPlayerID string, cards []domain.CardData) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
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

	_, err = session.WithTransaction(ctx, func(mongoCtx mongo.SessionContext) (interface{}, error) {
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

		// Decode tournament player
		var tournamentPlayer *domain.TournamentPlayer
		err = result.Decode(&tournamentPlayer)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}

		// Add the cards to the tournament player's collection
		// For each card, find if the user already has some of that card, and update or add it accordingly
		cardsToAdd := []domain.OwnedCard{}
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
				// Prepare card to add
				cardsToAdd = append(cardsToAdd, domain.OwnedCard{
					ID:           primitive.NewObjectID(),
					TournamentID: tournamentPlayer.TournamentID,
					UserID:       tournamentPlayer.UserID,
					Tags:         []string{},
					Count:        1,
					CardData:     card,
					CreatedAt:    primitive.NewDateTimeFromTime(time.Now()),
					UpdatedAt:    primitive.NewDateTimeFromTime(time.Now()),
				})

			} else if len(foundCards) == 1 {
				// Update count of existing card
				foundCards[0].Count += 1
				foundCards[0].UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
				result, err := MongoDatabaseClient.
					Database(DB_MAIN).
					Collection(COLLECTION_CARD_COLLECTION).
					UpdateByID(ctx, foundCards[0].ID, bson.M{"$set": foundCards[0]})
				if err != nil || result.MatchedCount == 0 {
					return nil, fmt.Errorf("%w: %v", ErrInternal, err)
				}
			} else {
				// TODO: Consolidate duplicate entries just in case
				return nil, fmt.Errorf("%w: duplicated entries for card found on database", ErrInternal)
			}
		}
		// Add all cards at once
		newValues := make([]interface{}, len(cardsToAdd))
		for i, cardToAdd := range cardsToAdd {
			newValues[i] = cardToAdd
		}

		_, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_CARD_COLLECTION).
			InsertMany(ctx, newValues)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return nil, nil
	})
	return err
}

func RemoveCardsFromTournamentPlayer(tournamentPlayerID string, cardsToRemove map[string]int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
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

	_, err = session.WithTransaction(ctx, func(mongoCtx mongo.SessionContext) (interface{}, error) {
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
		// Remove the cards to the tournament player's collection
		// For each card, find if the user already has some of that card, and update or remove it accordingly
		for cardID, count := range cardsToRemove {
			dbCardID, err := primitive.ObjectIDFromHex(cardID)
			if err != nil {
				return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
			}
			result, err := MongoDatabaseClient.
				Database(DB_MAIN).
				Collection(COLLECTION_CARD_COLLECTION).
				Find(ctx,
					bson.M{
						"tournament_id": tournamentPlayer.TournamentID,
						"user_id":       tournamentPlayer.UserID,
						"_id":           dbCardID,
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
				return nil, fmt.Errorf("%w: %v", ErrInternal, "card to remove not found")
			} else if len(foundCards) == 1 {
				if foundCards[0].Count > count {
					// Update count of existing card
					foundCards[0].Count -= count
					foundCards[0].UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
					result, err := MongoDatabaseClient.
						Database(DB_MAIN).
						Collection(COLLECTION_CARD_COLLECTION).
						UpdateByID(ctx, foundCards[0].ID, bson.M{"$set": foundCards[0]})
					if err != nil || result.MatchedCount == 0 {
						return nil, fmt.Errorf("%w: %v", ErrInternal, err)
					}
				} else {
					// Remove the card entirely
					result, err := MongoDatabaseClient.
						Database(DB_MAIN).
						Collection(COLLECTION_CARD_COLLECTION).
						DeleteOne(ctx, bson.M{"_id": foundCards[0].ID})
					if err != nil || result.DeletedCount == 0 {
						return nil, fmt.Errorf("%w: %v", ErrInternal, err)
					}
				}
			} else {
				// TODO: Consolidate duplicate entries just in case
				return nil, fmt.Errorf("%w: duplicated entries for card found on database", ErrInternal)
			}
		}
		return nil, nil
	})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return nil
}
