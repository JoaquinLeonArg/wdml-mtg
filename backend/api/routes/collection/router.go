package collection

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml_mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"github.com/rs/zerolog/log"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/collection").Subrouter()
	r.HandleFunc("", GetCollectionHandler).Methods(http.MethodGet)
}

//
// ENDPOINT: Get cards from tournament player's collection
//

type GetCollectionResponse struct {
	Cards   []domain.OwnedCard `json:"cards"`
	Count   int                `json:"count"`
	MaxPage int                `json:"max_page"`
}

func GetCollectionHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournamentID")
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
