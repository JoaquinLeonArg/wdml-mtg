package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type EventLog struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	TournamentID primitive.ObjectID `bson:"tournament_id" json:"tournament_id"`
	Type         EventLogType       `bson:"type" json:"type"`
	Data         interface{}        `bson:"data" json:"data"`
}

type EventLogType string

const (
	EventLogTypeOpenBoosters       EventLogType = "elt_open_boosters"
	EventLogTypeWinMatch           EventLogType = "elt_win_match"
	EventLogTypeAddMythic          EventLogType = "elt_add_mythic"
	EventLogTypeDistributeBoosters EventLogType = "elt_distribute_boosters"
)

type EventLogDataOpenBoosters struct {
	Username string `bson:"username" json:"username"`
	SetName  string `bson:"set_name" json:"set_name"`
	Count    int    `bson:"count" json:"count"`
}
