package tournament

import (
	"errors"

	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml_mtg/backend/errors"
	"github.com/joaquinleonarg/wdml_mtg/backend/pkg/mtgapi"
)

func GetTournamentByID(tournamentID string) (*domain.Tournament, error) {
	return db.GetTournamentByID(tournamentID)
}

func CreateTournament(tournament domain.Tournament) (string, error) {
	tournamentID, err := db.CreateTournament(tournament)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			return "", apiErrors.ErrDuplicatedResource
		}
		if errors.Is(err, db.ErrNotFound) {
			return "", apiErrors.ErrNotFound
		}
		return "", apiErrors.ErrInternal
	}
	_, err = db.CreateTournamentPlayer(
		domain.TournamentPlayer{
			UserID:       tournament.OwnerID,
			TournamentID: tournamentID,
			AccessLevel:  domain.AccessLevelAdministrator,
			GameResources: domain.GameResources{
				OwnedCards: []domain.Card{},
				Decks:      []domain.Deck{},
				Wildcards: domain.OwnedWildcards{
					CommonCount:      0,
					UncommonCount:    0,
					RareCount:        0,
					MythicRareCount:  0,
					MasterpieceCount: 0,
				},
				BoosterPacks: []domain.OwnedBoosterPack{},
				Rerolls:      0,
				Coins:        0,
			},
			TournamentPoints: 0,
		},
	)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			return "", apiErrors.ErrDuplicatedResource
		}
		if errors.Is(err, db.ErrNotFound) {
			return "", apiErrors.ErrNotFound
		}
		return "", apiErrors.ErrInternal
	}
	return tournamentID.Hex(), nil
}

func GetTournamentBoosterPacks() ([]domain.BoosterPackData, error) {
	sets, err := mtgapi.GetAllSets()
	if err != nil {
		return nil, apiErrors.ErrInternal
	}
	var boosterPacks []domain.BoosterPackData

	for _, set := range sets {
		boosterPacks = append(boosterPacks, domain.BoosterPackData{
			SetCode:     string(set.SetCode),
			SetName:     set.Name,
			BoosterType: domain.BoosterTypeDraft,
		})
	}

	return boosterPacks, nil
}
