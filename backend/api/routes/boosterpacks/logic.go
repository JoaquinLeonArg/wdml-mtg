package boosterpacks

import (
	"errors"
	"strings"
	"time"

	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml_mtg/backend/errors"
	boostergen "github.com/joaquinleonarg/wdml_mtg/backend/internal/booster_gen"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTournamentBoosterPacks() ([]domain.BoosterPack, error) {
	boosterPacks, err := db.GetAllBoosterPacks()
	if err != nil {
		return nil, apiErrors.ErrInternal
	}

	return boosterPacks, nil
}

func AddTournamentBoosterPacks(userID, tournamentID string, boosterPack AddTournamentBoosterPacksRequest) error {
	// Check if tournament player has permissions
	tournament_player, err := db.GetTournamentPlayer(tournamentID, userID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apiErrors.ErrNotFound
		}
		return apiErrors.ErrInternal
	}

	if tournament_player.AccessLevel != domain.AccessLevelAdministrator {
		return apiErrors.ErrUnauthorized
	}

	// Get the set's data
	packs, err := db.GetAllBoosterPacks()
	if err != nil {
		return apiErrors.ErrInternal
	}
	var setName string
	var setDescription string
	found := false
	for _, pack := range packs {
		if pack.SetCode == strings.ToLower(boosterPack.SetCode) {
			setName = pack.Name
			setDescription = string(pack.Description)
			found = true
			break
		}
	}
	if !found {
		return apiErrors.ErrNotFound
	}

	err = db.AddPacksToTournamentPlayer(tournament_player.ID.Hex(), domain.OwnedBoosterPack{
		Available:   boosterPack.Count,
		SetCode:     boosterPack.SetCode,
		Name:        setName,
		Description: setDescription,
	})
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apiErrors.ErrNotFound
		}
		return apiErrors.ErrInternal
	}

	return nil
}

func OpenBoosterPack(userID, tournamentID string, setCode string) ([]domain.CardData, error) {
	cards, err := boostergen.GenerateBooster(strings.ToLower(setCode), boostergen.GetBoosterDataFromDb)

	if err != nil {
		log.Debug().Err(err).Msg("failed to generate booster pack")
		return nil, apiErrors.ErrInternal
	}

	err = db.ConsumeBoosterPackForTournamentPlayer(userID, tournamentID, setCode, cards)
	if err != nil {
		log.Debug().Err(err).Msg("failed to open booster pack")
		if errors.Is(err, db.ErrNotFound) {
			return nil, apiErrors.ErrNotFound
		}
		if errors.Is(err, db.ErrInvalidID) {
			return nil, apiErrors.ErrBadRequest
		}
		return nil, apiErrors.ErrInternal
	}
	return cards, nil
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
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			return apiErrors.ErrDuplicatedResource
		}
		return apiErrors.ErrInternal
	}
	return err
}

func UpdateBoosterPack(boosterPack domain.BoosterPack) error {
	err := db.UpdateBoosterPack(domain.BoosterPack{
		SetCode:     boosterPack.SetCode,
		Name:        boosterPack.Name,
		Description: boosterPack.Description,
		CardCount:   boosterPack.CardCount,
		Slots:       boosterPack.Slots,
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	})
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			return apiErrors.ErrDuplicatedResource
		}
		return apiErrors.ErrInternal
	}
	return err
}
