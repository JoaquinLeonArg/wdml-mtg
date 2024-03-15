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
	r.HandleFunc("/", CreateEmptyDeckHandler).Methods(http.MethodPost)
	r.HandleFunc("/card", AddOwnedCardToDeckHandler).Methods(http.MethodPost)
	r.HandleFunc("/card", RemoveCardFromDeckHandler).Methods(http.MethodDelete)
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

	deck, err := GetDeckById(deckId)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get deck data")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
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
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
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

	err = CreateEmptyDeck(ownerID, createEmptyDeckRequest.Deck.Name, createEmptyDeckRequest.Deck.Description, tournamentID)
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
	CardId string           `json:"card_id"`
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
	var req AddOwnedCardToDeckRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	err = AddOwnedCardToDeck(req.CardId, req.DeckId, req.Amount, req.Board)
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
	DeckId string           `json:"deck_id"`
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

	err = RemoveCardFromDeck(removeCardfromDeckRequest.Card, removeCardfromDeckRequest.DeckId, removeCardfromDeckRequest.Amount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(RemoveCardFromDeckResponse{}))
}
