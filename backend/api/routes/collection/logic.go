package collection

import (
	"strings"

	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"github.com/rs/zerolog/log"
)

func GetCollectionCards(userID, tournamentID, filters string, count, page int) ([]domain.OwnedCard, int, error) {
	log.Debug().Str("filters", filters).Send()
	dbFilters := []db.CardFilter{}
	for _, filter := range strings.Split(filters, "+") {
		for _, filterOperation := range []db.CardFilterOperation{
			db.CardFilterOperationEq,
			db.CardFilterOperationLt,
			db.CardFilterOperationGt,
		} {
			for _, filterType := range []db.CardFilterType{
				db.CardFilterTypeName,
				db.CardFilterTypeTags,
				db.CardFilterTypeRarity,
				db.CardFilterTypeColor,
				db.CardFilterTypeTypes,
				db.CardFilterTypeOracle,
				db.CardFilterTypeSetCode,
				db.CardFilterTypeMV,
			} {
				if strings.Contains(filter, string(filterOperation)) {
					splitted := strings.Split(filter, string(filterOperation))
					if splitted[1] == "" {
						continue
					}
					if splitted[0] == string(filterType) {
						dbFilters = append(dbFilters, db.CardFilter{
							Type:      filterType,
							Operation: filterOperation,
							Value:     splitted[1],
						})
					}
					continue
				}
			}
		}
	}
	return db.GetCardsFromTournamentPlayer(userID, tournamentID, dbFilters, count, page)
}

func GetOwnedCardById(cardId string) (domain.OwnedCard, error) {
	return db.GetOwnedCardById(cardId)
}
