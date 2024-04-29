package tournament

import (
	"errors"

	"github.com/joaquinleonarg/wdml-mtg/backend/db"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml-mtg/backend/errors"
	"github.com/rs/zerolog/log"
)

func GetTournamentByID(tournamentID string) (*domain.Tournament, error) {
	tournament, err := db.GetTournamentByID(tournamentID)
	if err != nil {
		return nil, apiErrors.ErrInternal
	}
	return tournament, nil
}

func GetTournamentsForUser(userID string) ([]domain.Tournament, error) {
	tournaments, err := db.GetTournamentsForUser(userID)
	if err != nil {
		return nil, apiErrors.ErrInternal
	}
	return tournaments, nil
}

func GetTournamentPlayers(tournamentID string) ([]domain.TournamentPlayer, []domain.User, error) {
	tournament_players, users, err := db.GetTournamentPlayers(tournamentID)

	// Redact sensitive information
	for index := range users {
		users[index].Password = nil
		users[index].Email = ""
	}

	return tournament_players, users, err
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
	// TODO: Consolidate this creation with the other one
	_, err = db.CreateTournamentPlayer(
		domain.TournamentPlayer{
			UserID:       tournament.OwnerID,
			TournamentID: tournamentID,
			AccessLevel:  domain.AccessLevelAdministrator,
			GameResources: domain.GameResources{
				Decks: []domain.Deck{},
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

func UpdateStore(tournamentID, userID string, store domain.Store) error {
	// Check if the user is an admin or moderator
	tournamentPlayer, err := db.GetTournamentPlayer(tournamentID, userID)
	if err != nil {
		return apiErrors.ErrInternal
	}
	if tournamentPlayer.AccessLevel != domain.AccessLevelAdministrator && tournamentPlayer.AccessLevel != domain.AccessLevelModerator {
		return apiErrors.ErrUnauthorized
	}

	// Check that all booster packs exist
	for _, boosterPack := range store.BoosterPacks {
		_, err := db.GetBoosterPackByID(boosterPack.BoosterPackID.Hex())
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return apiErrors.ErrNotFound
			}
			return apiErrors.ErrInternal
		}
	}
	err = db.UpdateTournamentStore(tournamentID, store)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apiErrors.ErrNotFound
		}
		return apiErrors.ErrInternal
	}

	return nil
}

func GetStore(tournamentID string) (*domain.Store, error) {
	tournament, err := db.GetTournamentByID(tournamentID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, apiErrors.ErrNotFound
		}
		return nil, apiErrors.ErrInternal
	}

	// v1.1: create store if it doesn't exist
	if tournament.Store.BoosterPacks == nil {
		log.Info().Str("tournament_id", tournamentID).Msg("initializing tournament store")
		err := db.UpdateTournamentStore(tournamentID, domain.Store{BoosterPacks: []domain.StoreBoosterPack{}})
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return nil, apiErrors.ErrNotFound
			}
			return nil, apiErrors.ErrInternal
		}
	}

	return &tournament.Store, nil
}
