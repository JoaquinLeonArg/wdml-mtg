package boosterpacks

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/routes/auth"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	"github.com/rs/zerolog/log"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/boosterpacks").Subrouter()
	r.HandleFunc("/tournament", GetTournamentBoosterPacksHandler).Methods(http.MethodGet)
	r.HandleFunc("/tournament", AddTournamentBoosterPacksHandler).Methods(http.MethodPost)
	r.HandleFunc("/open", OpenBoosterPackHandler).Methods(http.MethodPost)
	r.HandleFunc("/", CreateBoosterPackHandler).Methods(http.MethodPost)
	r.HandleFunc("/", UpdateBoosterPackHandler).Methods(http.MethodPut)
	r.HandleFunc("/buy", BuyStoreBoosterPackHandler).Methods(http.MethodPost)
}

//
// ENDPOINT: Get all available booster pack types for a tournament
//

type GetTournamentBoosterPacksResponse struct {
	BoosterPacks []domain.BoosterPack `json:"booster_packs"`
}

func GetTournamentBoosterPacksHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Get all available booster packs for this tournament
	boosterPacks, err := GetTournamentBoosterPacks()
	if err != nil {
		log.Debug().Err(err).Msg("failed to get vanilla booster pack data")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetTournamentBoosterPacksResponse{BoosterPacks: boosterPacks}))
}

//
// ENDPOINT: Add booster packs to all players on a tournament
//

type AddTournamentBoosterPacksRequest struct {
	Count              int    `json:"count"`
	SetCode            string `json:"set_code"`
	TournamentPlayerId string `json:"tournament_player_id"`
}

type AddTournamentBoosterPacksResponse struct{}

func AddTournamentBoosterPacksHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get user ID from request context
	ownerID, ok := r.Context().Value("user_id").(string)
	if ownerID == "" || !ok {
		log.Debug().
			Msg("failed to read user id from context")
		http.Error(w, "", http.StatusForbidden)
		return
	}

	// Decode body data
	var addTournamentBoosterPacksRequest AddTournamentBoosterPacksRequest
	err := json.NewDecoder(r.Body).Decode(&addTournamentBoosterPacksRequest)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Add the booster packs to each player, checking if the user is allowed to add them and if they are valid packs
	err = AddTournamentBoosterPacks(ownerID, tournamentID, addTournamentBoosterPacksRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(AddTournamentBoosterPacksResponse{}))
}

//
// ENDPOINT: Open a booster pack and add the cards to the player's collection
//

type OpenBoosterPackRequest struct {
	SetCode string `json:"set_code"`
}

type OpenBoosterPackResponse struct {
	CardData []domain.CardData `json:"card_data"`
}

func OpenBoosterPackHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Decode body data
	var openBoosterPackRequest OpenBoosterPackRequest
	err = json.NewDecoder(r.Body).Decode(&openBoosterPackRequest)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Try to open the pack, add the cards to the collection and get them here to send in the response
	cards, err := OpenBoosterPack(userID, tournamentID, openBoosterPackRequest.SetCode)
	if err != nil {
		log.Debug().Err(err).Msg("failed to open booster pack")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(OpenBoosterPackResponse{CardData: cards}))
}

// TODO: Restrict the request body
// Endpoint: Create new booster pack
func CreateBoosterPackHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Decode body data
	var boosterPack domain.BoosterPack
	err := json.NewDecoder(r.Body).Decode(&boosterPack)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Add the booster packs
	err = CreateNewBoosterPack(boosterPack)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(AddTournamentBoosterPacksResponse{}))
}

// TODO: Restrict the request body
// Endpoint: Edit booster pack
func UpdateBoosterPackHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Decode body data
	var boosterPack domain.BoosterPack
	err := json.NewDecoder(r.Body).Decode(&boosterPack)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Add the booster packs
	err = UpdateBoosterPack(boosterPack)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(AddTournamentBoosterPacksResponse{}))
}

type BuyStoreBoosterPackRequest struct {
	BoosterPackID string `json:"booster_pack_id"`
}

type BuyStoreBoosterPackResponse struct {
}

func BuyStoreBoosterPackHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Decode body data
	var buyStoreBoosterPackRequest BuyStoreBoosterPackRequest
	err = json.NewDecoder(r.Body).Decode(&buyStoreBoosterPackRequest)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Check coins, remove them and add the booster to the tournament player
	err = BuyBoosterPack(tournamentID, userID, buyStoreBoosterPackRequest.BoosterPackID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(BuyStoreBoosterPackResponse{}))
}
