package deck

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml_mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml_mtg/backend/api/routes/auth"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"github.com/rs/zerolog/log"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/deck").Subrouter()
	r.HandleFunc("", GetDeckByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/tournament_player", GetDecksForTournamentPlayerHandler).Methods(http.MethodGet)
	r.HandleFunc("", CreateEmptyDeckHandler).Methods(http.MethodPost)
	r.HandleFunc("/card", AddOwnedCardToDeckHandler).Methods(http.MethodPost)
	r.HandleFunc("/card/remove", RemoveCardFromDeckHandler).Methods(http.MethodPost)
}

//
// ENDPOINT: Get deck by ID
//

type GetDeckByIdResponse struct {
	Deck  *domain.Deck       `json:"deck"`
	Cards []domain.OwnedCard `json:"cards"`
}

func GetDeckByIdHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament ID from query
	deckId := r.URL.Query().Get("deck_id")
	if deckId == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	deck, cards, err := GetDeckById(deckId)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get deck data")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetDeckByIdResponse{Deck: deck, Cards: cards}))
}

//
// ENDPOINT: Get deck by tournament player ID
//

type GetDecksByTournamentPlayerIDResponse struct {
	Decks []domain.Deck `json:"decks"`
}

func GetDecksForTournamentPlayerHandler(w http.ResponseWriter, r *http.Request) {
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

	// Get tournament player, then get their decks
	decks, err := GetDecksForTournamentPlayer(tournamentID, userID)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get deck data")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetDecksByTournamentPlayerIDResponse{Decks: decks}))
}

//
// ENDPOINT: Create a new, empty deck
//

type CreateEmptyDeckRequest struct {
	Deck struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"deck"`
	TournamentID string `json:"tournament_id"`
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
	tournamentID := createEmptyDeckRequest.TournamentID
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Create a deck empty of cards, with the name and description provided
	err = CreateEmptyDeck(
		ownerID,
		createEmptyDeckRequest.Deck.Name,
		createEmptyDeckRequest.Deck.Description,
		tournamentID,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(CreateEmptyDeckResponse{}))
}

//
// ENDPOINT: Add cards from a tournament player's collection to one of their decks
//

type AddOwnedCardToDeckRequest struct {
	OwnedCardID string           `json:"owned_card_id"`
	DeckID      string           `json:"deck_id"`
	Amount      int              `json:"amount"`
	Board       domain.DeckBoard `json:"board"`
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

	err = AddOwnedCardToDeck(req.OwnedCardID, req.DeckID, req.Amount, req.Board)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(AddOwnedCardToDeckResponse{}))
}

//
// ENDPOINT: Remove a card from a deck
//

type RemoveCardFromDeckRequest struct {
	OwnedCardID string           `json:"owned_card_id"`
	DeckID      string           `json:"deck_id"`
	Amount      int              `json:"amount"`
	Board       domain.DeckBoard `json:"board"`
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
	var req RemoveCardFromDeckRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	err = RemoveCardFromDeck(
		req.OwnedCardID,
		req.DeckID,
		req.Board,
		req.Amount,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(RemoveCardFromDeckResponse{}))
}
