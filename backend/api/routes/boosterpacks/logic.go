package boosterpacks

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	scryfallapi "github.com/BlueMonday/go-scryfall"
	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml_mtg/backend/errors"
	boostergen "github.com/joaquinleonarg/wdml_mtg/backend/internal/booster_gen"
	"github.com/joaquinleonarg/wdml_mtg/backend/pkg/scryfall"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTournamentBoosterPacks() ([]domain.BoosterPackData, error) {
	packs, err := db.GetAllBoosterPacks()
	if err != nil {
		return nil, apiErrors.ErrInternal
	}
	var boosterPacks []domain.BoosterPackData

	for _, pack := range packs {
		boosterPacks = append(boosterPacks, domain.BoosterPackData{
			SetCode:     strings.ToUpper(pack.SetCode),
			SetName:     pack.Name,
			Description: pack.Description,
			BoosterType: domain.BoosterTypeDraft,
		})
	}

	return boosterPacks, nil
}

func AddTournamentBoosterPacks(userID, tournamentID string, boosterPacks AddTournamentBoosterPacksRequest) error {
	// Check if the tournament exists
	tournament, err := db.GetTournamentByID(tournamentID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apiErrors.ErrNotFound
		}
		return apiErrors.ErrInternal
	}
	// Check if user is the owner
	if tournament.OwnerID.Hex() != userID {
		return apiErrors.ErrUnauthorized
	}
	// Check that all set codes exist
	// TODO: Custom boosters
	setNames := make(map[string]string, len(boosterPacks.BoosterPacks))
	setDescriptions := make(map[string]string, len(boosterPacks.BoosterPacks))
	packs, err := db.GetAllBoosterPacks()
	if err != nil {
		return apiErrors.ErrInternal
	}
	for _, boosterPacks := range boosterPacks.BoosterPacks {
		found := false
		for _, pack := range packs {
			if pack.SetCode == strings.ToLower(boosterPacks.Set) {
				setNames[pack.SetCode] = pack.Name
				setDescriptions[pack.SetCode] = string(pack.Description)
				found = true
				break
			}
		}
		if !found {
			return apiErrors.ErrNotFound
		}
	}

	tournament_players, err := db.GetTournamentPlayers(tournamentID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apiErrors.ErrNotFound
		}
		return apiErrors.ErrInternal
	}

	ownedPacks := make([]domain.OwnedBoosterPack, 0, len(boosterPacks.BoosterPacks))
	for _, boosterPack := range boosterPacks.BoosterPacks {
		if boosterPack.Count == 0 || boosterPack.Set == "" || boosterPack.Type == "" {
			continue
		}
		ownedPacks = append(ownedPacks,
			domain.OwnedBoosterPack{
				BoosterGen:     domain.BoosterGenVanilla,
				BoosterGenData: nil,
				Available:      boosterPack.Count,
				Data: domain.BoosterPackData{
					SetCode:     boosterPack.Set,
					SetName:     setNames[boosterPack.Set],
					BoosterType: domain.BoosterType(boosterPack.Type),
					Description: setDescriptions[boosterPack.Set],
				},
			},
		)
	}

	if len(ownedPacks) == 0 {
		return apiErrors.ErrNoData
	}

	for _, tournament_player := range tournament_players {
		err = db.AddPacksToTournamentPlayer(tournament_player.ID.Hex(), ownedPacks)
		// TODO: Handle error better
		if err != nil {
			log.Warn().Err(err).Msg(fmt.Sprintf("failed to add booster packs to tournament player %s", tournament_player.ID.Hex()))
		}
	}
	return nil
}

func OpenBoosterPack(userID, tournamentID string, boosterPackData domain.BoosterPackData) ([]domain.CardData, error) {
	var cards []domain.CardData
	var err error
	if boosterPackData.BoosterType == domain.BoosterTypeDraft {
		// Vanilla booster
		_, err := boostergen.CheckIfBoosterExists(boosterPackData.SetCode)
		if err != nil {
			log.Debug().Str("set", boosterPackData.SetCode).Err(err).Msg("booster pack does not exist")
			return nil, apiErrors.ErrNotFound
		}
		cards, err = boostergen.GenerateBooster(strings.ToLower(boosterPackData.SetCode), boostergen.GetBoosterDataFromDb)

		if err != nil {
			log.Debug().Err(err).Msg("failed to generate booster pack")
			return nil, apiErrors.ErrInternal
		}
	} else {
		return nil, apiErrors.ErrNotFound
	}

	err = db.ConsumeBoosterPackForTournamentPlayer(userID, tournamentID, boosterPackData, cards)
	if err != nil {
		log.Debug().Err(err).Msg("failed to open booster pack")
		if errors.Is(err, db.ErrNotFound) {
			return nil, apiErrors.ErrNotFound
		}
		return nil, apiErrors.ErrInternal
	}
	return cards, nil
}

func GenerateVanillaBoosterPack(boosterPackData domain.BoosterPackData) ([]domain.CardData, error) {
	// Vanilla gen
	cards, err := scryfall.GetSetCards(boosterPackData.SetCode)

	if err != nil || len(cards) == 0 {
		log.Debug().Str("set", boosterPackData.SetCode).Err(err).Msg("failed to generate booster pack")
		return nil, apiErrors.ErrInternal
	}

	commons := make([]scryfallapi.Card, 0)
	uncommons := make([]scryfallapi.Card, 0)
	rares := make([]scryfallapi.Card, 0)
	for _, card := range cards {
		if card.Rarity == "common" && !strings.HasPrefix(card.TypeLine, "Basic") {
			commons = append(commons, card)
		}
		if card.Rarity == "uncommon" {
			uncommons = append(uncommons, card)
		}
		if card.Rarity == "rare" {
			rares = append(rares, card)
		}
		if card.Rarity == "mythic" {
			rares = append(rares, card)
		}
	}

	boosterPack := make([]domain.CardData, 0, 15)
	addCardsOfRarity := func(n int, rarity []scryfallapi.Card) []domain.CardData {
		newCards := []domain.CardData{}
		for len(newCards) < n {
			card := rarity[rand.Int()%len(rarity)]
			colors := []string{}
			for _, col := range card.Colors {
				colors = append(colors, string(col))
			}
			types := scryfall.ParseScryfallTypeline(card.TypeLine)

			newCard := domain.CardData{
				SetCode:         strings.ToUpper(card.Set),
				CollectorNumber: card.CollectorNumber,
				Name:            card.Name,
				Types:           types,
				ManaValue:       int(card.CMC),
				Colors:          colors,
			}
			cardFront, cardBack := scryfall.GetImageFromFaces(card)
			newCard.ImageURL = cardFront
			newCard.BackImageURL = cardBack

			newCards = append(newCards, newCard)
		}
		return newCards
	}

	boosterPack = append(boosterPack, addCardsOfRarity(11, commons)...)
	boosterPack = append(boosterPack, addCardsOfRarity(3, uncommons)...)
	boosterPack = append(boosterPack, addCardsOfRarity(1, rares)...)

	log.Debug().Interface("booster", boosterPack).Send()

	return boosterPack, nil
}

func CreateNewBoosterPack(boosterPack domain.BoosterPack) error {
	err := db.CreateBoosterPack(domain.BoosterPack{
		SetCode:     boosterPack.SetCode,
		Name:        boosterPack.Name,
		Description: boosterPack.Description,
		CardCount:   boosterPack.CardCount,
		Slots:       boosterPack.Slots,
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	})
	if errors.Is(err, db.ErrAlreadyExists) {
		return apiErrors.ErrDuplicatedResource
	}
	return nil
}
