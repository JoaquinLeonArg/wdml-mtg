package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// TournamentMatches collection
type TournamentMatch struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	TournamentID   primitive.ObjectID `bson:"tournament_id" json:"tournament_id"`
	FirstPlayerID  primitive.ObjectID `bson:"first_player_id" json:"first_player_id"`
	SecondPlayerID primitive.ObjectID `bson:"second_player_id" json:"second_player_id"`
	BlockID        primitive.ObjectID `bson:"block_id" json:"block_id"`
	Played         bool               `bson:"played" json:"played"`
	Result         MatchResult        `bson:"result" json:"result"`
}

type MatchResult string

const (
	ResultUnplayed        MatchResult = "mr_unplayed"
	ResultFirstPlayerWin  MatchResult = "mr_firstPlayerWin"
	ResultSecondPlayerWin MatchResult = "mr_secondPlayerWin"
	ResultTie             MatchResult = "mr_tie"
	ResultCancelled       MatchResult = "mr_cancelled"
)
