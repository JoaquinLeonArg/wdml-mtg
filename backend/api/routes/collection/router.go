package collection

import (
	"encoding/csv"
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/routes/auth"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	"github.com/rs/zerolog/log"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/collection").Subrouter()
	r.HandleFunc("", GetCollectionHandler).Methods(http.MethodGet)
	r.HandleFunc("/import", ImportCollectionHandler).Methods(http.MethodPost)
	r.HandleFunc("/tag", SetTagsForCollectionCardHandler).Methods(http.MethodPost)
	r.HandleFunc("/tradeup", TradeUpCardsHandler).Methods(http.MethodPost)
}

//
// ENDPOINT: Get cards from tournament player's collection
//

type EmptyResponse struct{}

type GetCollectionResponse struct {
	Cards   []domain.OwnedCard `json:"cards"`
	Count   int                `json:"count"`
	MaxPage int                `json:"max_page"`
}

func GetCollectionHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Get filters, count and page
	count := 0
	countQuery := r.URL.Query().Get("count")
	if countQuery != "" {
		val, err := strconv.Atoi(countQuery)
		if err != nil {
			log.Debug().
				Msg("failed to read count from query")
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		count = max(val, 75)
	}

	page := 0
	pageQuery := r.URL.Query().Get("page")
	if pageQuery != "" {
		val, err := strconv.Atoi(pageQuery)
		if err != nil {
			log.Debug().
				Msg("failed to read page from query")
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		page = val
	}

	filterQuery := r.URL.Query().Get("filters")

	// Get user ID from request context
	userID, ok := r.Context().Value("user_id").(string)
	if userID == "" || !ok {
		log.Debug().
			Msg("failed to read user id from context")
		http.Error(w, "", http.StatusForbidden)
		return
	}

	cards, total, err := GetCollectionCards(userID, tournamentID, filterQuery, count, page)

	if err != nil {
		log.Debug().Err(err).Msg("failed to get cards from collection")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(
		GetCollectionResponse{
			Cards:   cards,
			Count:   total,
			MaxPage: int(math.Ceil(float64(total) / float64(count))),
		},
	))
}

func ImportCollectionHandler(w http.ResponseWriter, r *http.Request) {
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
	csvReader := csv.NewReader(r.Body)
	allCards, err := csvReader.ReadAll()
	if err != nil {
		log.Debug().Err(err).Msg("failed to import cards")
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = ImportCollection(allCards, userID, tournamentID)
	if err != nil {
		log.Debug().Err(err).Msg("failed to import cards")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(EmptyResponse{}))
}

//
// ENDPOINT: Add tags to collection cards
//

type SetTagsForCollectionCardRequest struct {
	OwnedCardID string   `json:"owned_card_id"`
	Tags        []string `json:"tags"`
}

type SetTagsForCollectionCardResponse struct{}

func SetTagsForCollectionCardHandler(w http.ResponseWriter, r *http.Request) {
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
	var req SetTagsForCollectionCardRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	err = SetTagsToOwnedCard(ownerID, req.OwnedCardID, req.Tags)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(SetTagsForCollectionCardResponse{}))
}

type TradeUpCardsRequest struct {
	Cards map[string]int `json:"cards"`
}

type TradeUpCardsResponse struct {
	Cards []domain.CardData `json:"cards"`
}

func TradeUpCardsHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get user ID from request context
	ownerID, ok := r.Context().Value("user_id").(string)
	if ownerID == "" || !ok {
		log.Debug().
			Msg("failed to read user id from context")
		http.Error(w, "", http.StatusForbidden)
		return
	}

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Decode body data
	var req TradeUpCardsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	cards, err := TradeUpCards(req.Cards, ownerID, tournamentID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(TradeUpCardsResponse{Cards: cards}))
}
