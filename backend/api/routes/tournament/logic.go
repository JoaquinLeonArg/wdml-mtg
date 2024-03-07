package tournament

import (
	"errors"
	"slices"
	"strings"

	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml_mtg/backend/errors"
	"github.com/joaquinleonarg/wdml_mtg/backend/pkg/scryfall"
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
	sets, err := scryfall.GetAllSets()
	if err != nil {
		return nil, apiErrors.ErrInternal
	}
	var boosterPacks []domain.BoosterPackData

	for _, set := range sets {
		if slices.Contains([]string{"core", "expansion"}, string(set.SetType)) && !set.Digital {
			boosterPacks = append(boosterPacks, domain.BoosterPackData{
				SetCode:     strings.ToUpper(set.Code),
				SetName:     set.Name,
				Expansion:   string(set.SetType),
				BoosterType: domain.BoosterTypeDraft,
			})
		}
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
	if tournament.OwnerID.String() != userID {
		return apiErrors.ErrUnauthorized
	}
	// Check that all set codes exist
	// TODO: Custom boosters
	setNames := make(map[string]string, len(boosterPacks.BoosterPacks))
	setTypes := make(map[string]string, len(boosterPacks.BoosterPacks))
	sets, err := scryfall.GetAllSets()
	if err != nil {
		return apiErrors.ErrInternal
	}
	for _, boosterPacks := range boosterPacks.BoosterPacks {
		for _, set := range sets {
			if set.Code == boosterPacks.Set {
				setNames[set.Code] = set.Name
				setTypes[set.Code] = string(set.SetType)
				break
			}
		}
		return apiErrors.ErrNotFound
	}

	tournament_players, err := db.GetTournamentPlayers(tournamentID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apiErrors.ErrNotFound
		}
		return apiErrors.ErrInternal
	}

	ownedPacks := make([]domain.OwnedBoosterPack, len(boosterPacks.BoosterPacks))
	for _, boosterPack := range boosterPacks.BoosterPacks {
		ownedPacks = append(ownedPacks,
			domain.OwnedBoosterPack{
				BoosterGen:     domain.BoosterGenVanilla,
				BoosterGenData: nil,
				Available:      boosterPack.Count,
				Data: domain.BoosterPackData{
					SetCode:     boosterPack.Set,
					SetName:     setNames[boosterPack.Set],
					BoosterType: domain.BoosterType(boosterPack.Type),
					Expansion:   setTypes[boosterPack.Set],
				},
			},
		)
	}

	for _, tournament_player := range tournament_players {
		// TODO: Handle errors
		_ = db.AddPacksToTournamentPlayer(tournament_player.ID.String(), ownedPacks)
	}
	return nil
}
