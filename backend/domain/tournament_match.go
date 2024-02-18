package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// TournamentMatches collection
type TournamentMatch struct {
	ID             primitive.ObjectID
	TournamentID   primitive.ObjectID
	FirstPlayerID  primitive.ObjectID
	SecondPlayerID primitive.ObjectID
	DueDate        primitive.DateTime
	Played         bool
	Result         MatchResult
}

type MatchResult string

const (
	ResultUnplayed        MatchResult = "mr_unplayed"
	ResultFirstPlayerWin  MatchResult = "mr_firstPlayerWin"
	ResultSecondPlayerWin MatchResult = "mr_secondPlayerWin"
	ResultTie             MatchResult = "mr_tie"
	ResultCancelled       MatchResult = "mr_cancelled"
)
