package db

import (
	"context"
	"fmt"
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
	Values    []string
}

type CardFilterOperation string

const (
	CardFilterOperationEq CardFilterOperation = "="
	CardFilterOperationLt CardFilterOperation = "<"
	CardFilterOperationGt CardFilterOperation = ">"
)

type CardFilterType string

const (
	CardFilterTypeColor CardFilterType = "color"
)

func GetCardsFromTournamentPlayer(userID, tournamentID string, filters []CardFilter, count, page int) ([]domain.OwnedCard, error) {
	log.Debug().Interface("filters", filters).Int("count", count).Int("page", page)
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

	dbFilters := bson.M{}

	for _, filter := range filters {
		switch filter.Type {
		// By color
		case CardFilterTypeColor:
			switch filter.Operation {
			case CardFilterOperationLt:
				dbFilters["card_data.colors"] = filter.Values
			case CardFilterOperationEq:
				dbFilters["card_data.colors"] = bson.M{"$all": filter.Values}
			}
			// TODO: Other filters

		}
	}

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer session.EndSession(ctx)

	// Find the user cards that satisfy the filters
	cards, err := session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
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

		// Decode cards
		var filteredCards []domain.OwnedCard
		err = cursor.All(ctx, &filteredCards)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return filteredCards, nil
	})

	if err != nil {
		return nil, err
	}

	return cards.([]domain.OwnedCard), nil
}
