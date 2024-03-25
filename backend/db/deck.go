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

func GetDeckByID(deckID string) (*domain.Deck, []domain.OwnedCard, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	dbDeckID, err := primitive.ObjectIDFromHex(deckID)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	// Find deck
	result := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_DECKS).
		FindOne(ctx,
			bson.M{"_id": dbDeckID},
		)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode deck
	var deck *domain.Deck
	err = result.Decode(&deck)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Get card ids
	cardIDs := []primitive.ObjectID{}
	for _, card := range deck.Cards {
		cardIDs = append(cardIDs, card.OwnedCardID)
	}

	// Find cards
	cursor, err := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_CARD_COLLECTION).
		Find(ctx, bson.M{
			"_id": bson.M{"$in": cardIDs},
		})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode cards
	var cards []domain.OwnedCard
	err = cursor.All(ctx, &cards)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return deck, cards, nil
}

func GetDecksForTournamentPlayer(tournamentPlayerID string) ([]domain.Deck, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	dbTournamentPlayerID, err := primitive.ObjectIDFromHex(tournamentPlayerID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	filter := bson.M{"tournament_player_id": dbTournamentPlayerID}
	cursor, err := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_DECKS).
		Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode decks
	var decks []domain.Deck
	err = cursor.All(ctx, &decks)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return decks, nil
}

func CreateEmptyDeck(deck domain.Deck) error {
	if deck.ID != primitive.NilObjectID {
		return ErrObjectIDProvided
	}
	deck.ID = primitive.NewObjectID()
	deck.Cards = make([]domain.DeckCard, 0)
	deck.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	deck.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
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
		// TODO: Check that deck with the same name doesn't exist
		resultInsert, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_DECKS).
			InsertOne(ctx, deck)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return resultInsert, nil
	})

	log.Debug().Str("deck_id", deck.ID.String()).Msg("created deck")

	return err
}

func AddOwnedCardToDeck(cardID string, deckID string, amount int, board domain.DeckBoard) error {
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
		card, err := GetOwnedCardById(cardID)
		if err != nil {
			return nil, err
		}
		isBasic := false
		for _, cardType := range card.CardData.Types {
			if cardType == "Basic" {
				isBasic = true
				break
			}
		}

		deck, _, err := GetDeckByID(deckID)
		if err != nil {
			return nil, err
		}

		var foundCard domain.DeckCard
		foundAmount := 0
		foundIndex := -1
		for index, deckCard := range deck.Cards {
			if deckCard.OwnedCardID == card.ID {
				foundAmount += deckCard.Count
				// This will find the amount of the given card in the deck
				if deckCard.Board == board {
					foundCard = deckCard
					foundIndex = index
				}
			}
		}
		if foundIndex != -1 {
			if foundAmount+amount <= 4 || isBasic {
				if foundAmount+amount <= card.Count {
					foundCard.Count += amount
					deck.Cards[foundIndex] = foundCard
				} else {
					return nil, fmt.Errorf("%w: %s", ErrInternal, "not enough cards in collection")
				}
			} else {
				return nil, fmt.Errorf("%w: %s", ErrInternal, "too many copies of card in deck")
			}
		} else {
			deck.Cards = append(deck.Cards, domain.DeckCard{
				OwnedCardID: card.ID,
				Count:       amount,
				Board:       board,
			})
		}

		updateResult, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_DECKS).
			UpdateByID(ctx, deck.ID, bson.M{"$set": deck})

		if err != nil || updateResult.MatchedCount == 0 {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}

		return nil, nil
	})
	return err
}

func RemoveDeckCardFromDeck(ownedCardID, deckID string, board domain.DeckBoard, amount int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	dbOwnedCardID, err := primitive.ObjectIDFromHex(ownedCardID)
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

	// Check if card already exists in deck

	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		deck, _, err := GetDeckByID(deckID)
		if err != nil {
			return nil, err
		}

		newDeckCards := make([]domain.DeckCard, 0)
		for _, deckCard := range deck.Cards {
			if dbOwnedCardID == deckCard.OwnedCardID &&
				board == deckCard.Board {
				if deckCard.Count-amount <= 0 {
					continue
				} else {
					deckCard.Count -= amount
					newDeckCards = append(newDeckCards, deckCard)
				}
			} else {
				newDeckCards = append(newDeckCards, deckCard)
			}
		}

		deck.Cards = newDeckCards
		updateResult, err := MongoDatabaseClient.
			Database(DB_MAIN).
			Collection(COLLECTION_DECKS).
			UpdateByID(ctx, deck.ID, bson.M{"$set": deck})

		if err != nil || updateResult.MatchedCount == 0 {
			return nil, fmt.Errorf("%w: %v", ErrInternal, err)
		}
		return nil, nil
	})
	return err
}
