package deck

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml_mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"github.com/rs/zerolog/log"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/deck").Subrouter()
	r.HandleFunc("/", GetDeckByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/player/", GetDecksByTournamentPlayerIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/new", CreateEmptyDeckHandler).Methods(http.MethodPost)
	r.HandleFunc("/addCard", AddOwnedCardToDeckHandler).Methods(http.MethodPost)
	r.HandleFunc("/removeCard", RemoveCardFromDeckHandler).Methods(http.MethodPost)
}

//
// ENDPOINT: Get deck by ID
//

type GetDeckByIdResponse struct {
	Deck *domain.Deck `json:"deck"`
}

func GetDeckByIdHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament ID from query
	deckId := r.URL.Query().Get("id")
	if deckId == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Get all available booster packs for this tournament
	deck, err := GetDeckById(deckId)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get deck data")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetDeckByIdResponse{Deck: deck}))
}

//
// ENDPOINT: Get deck by Tournament Player ID
//

type GetDecksByTournamentPlayerIdResponse struct {
	Decks []domain.Deck `json:"decks"`
}

func GetDecksByTournamentPlayerIdHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	tournamentPlayerId := r.URL.Query().Get("tournamentPlayerId")
	if tournamentPlayerId == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	decks, err := GetDecksByTournamentPlayerId(tournamentPlayerId)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get deck data")
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetDecksByTournamentPlayerIdResponse{Decks: decks}))
}

type CreateEmptyDeckRequest struct {
	Deck struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"deck"`
	TournamentId string
}

type CreateEmptyDeckResponse struct{}

func CreateEmptyDeckHandler(w http.ResponseWriter, r *http.Request) {
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
	var createEmptyDeckRequest CreateEmptyDeckRequest
	err := json.NewDecoder(r.Body).Decode(&createEmptyDeckRequest)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Get tournament ID from body
	tournamentID := createEmptyDeckRequest.TournamentId
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Add the booster packs to each player, checking if the user is allowed to add them and if they are valid packs
	err = CreateEmptyDeck(ownerID, createEmptyDeckRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(CreateEmptyDeckResponse{}))
}

type AddOwnedCardToDeckRequest struct {
	Card   domain.OwnedCard `json:"card"`
	DeckId string           `json:"deck_id"`
	Amount int              `json:"amount"`
	Board  domain.DeckBoard `json:"board"`
}

type AddOwnedCardToDeckResponse struct{}

func AddOwnedCardToDeckHandler(w http.ResponseWriter, r *http.Request) {
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
	var addOwnedCardToDeckRequest AddOwnedCardToDeckRequest
	err := json.NewDecoder(r.Body).Decode(&addOwnedCardToDeckRequest)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Add the booster packs to each player, checking if the user is allowed to add them and if they are valid packs
	err = AddOwnedCardToDeck(addOwnedCardToDeckRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(AddOwnedCardToDeckResponse{}))
}

type RemoveCardFromDeckRequest struct {
	Card   domain.DeckCard  `json:"card"`
	Amount int              `json:"amount"`
	Board  domain.DeckBoard `json:"board"`
}

type RemoveCardFromDeckResponse struct{}

func RemoveCardFromDeckHandler(w http.ResponseWriter, r *http.Request) {
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
	var removeCardfromDeckRequest RemoveCardFromDeckRequest
	err := json.NewDecoder(r.Body).Decode(&removeCardfromDeckRequest)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Add the booster packs to each player, checking if the user is allowed to add them and if they are valid packs
	err = RemoveCardFromDeck(removeCardfromDeckRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(RemoveCardFromDeckResponse{}))
}
