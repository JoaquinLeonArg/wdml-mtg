package db

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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
