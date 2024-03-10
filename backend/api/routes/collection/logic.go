package collection

import (
	"strings"

	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"github.com/rs/zerolog/log"
)

func GetCollectionCards(userID, tournamentID, filters string, count, page int) ([]domain.OwnedCard, error) {
	log.Debug().Str("filters", filters).Send()
	dbFilters := []db.CardFilter{}
	for _, filter := range strings.Split(filters, "+") {
		for _, filterOperation := range []db.CardFilterOperation{
			db.CardFilterOperationEq,
			db.CardFilterOperationLt,
			db.CardFilterOperationGt,
		} {
			for _, filterType := range []db.CardFilterType{
				db.CardFilterTypeColor,
			} {
				if strings.Contains(filter, string(filterOperation)) {
					splitted := strings.Split(filter, string(filterOperation))
					if splitted[0] == string(filterType) {
						dbFilters = append(dbFilters, db.CardFilter{
							Type:      filterType,
							Operation: filterOperation,
							Values:    strings.Split(splitted[1], "|"),
						})
					}
					break
				}
			}
		}
	}
	return db.GetCardsFromTournamentPlayer(userID, tournamentID, dbFilters, count, page)
}
