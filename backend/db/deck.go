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

func GetDeckByID(deckID string) (*domain.Deck, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	dbDeckID, err := primitive.ObjectIDFromHex(deckID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	result := MongoDatabaseClient.
		Database(DB_MAIN).
		Collection(COLLECTION_DECKS).
		FindOne(ctx,
			bson.M{"_id": dbDeckID},
		)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Decode deck
	var deck *domain.Deck
	err = result.Decode(&deck)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return deck, nil
}

func GetDecksForTournamentPlayer(tournamentPlayerID string) ([]domain.Deck, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	card, err := GetOwnedCardById(cardID)
	if err != nil {
		return err
	}

	deck, err := GetDeckByID(deckID)
	if err != nil {
		return err
	}

	log.Debug().Interface("card", card.CardData.Name).Send()

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
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
			if foundAmount+amount <= 4 {
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

func RemoveDeckCardFromDeck(card domain.DeckCard, deckID string, amount int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ownedCard, err := GetOwnedCardById(card.OwnedCardID.Hex())
	if err != nil {
		return err
	}
	log.Debug().Interface("card", ownedCard.CardData.Name).Send()

	// Begin transaction
	session, err := MongoDatabaseClient.
		StartSession()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer session.EndSession(ctx)

	deck, err := GetDeckByID(deckID)
	if err != nil {
		return err
	}

	// Check if card already exists in deck

	_, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		newDeckCards := make([]domain.DeckCard, 0)
		for _, deckCard := range deck.Cards {
			if card.OwnedCardID == ownedCard.ID &&
				deckCard.Board == card.Board {
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
