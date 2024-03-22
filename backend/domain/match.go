package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Match struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	SeasonID    primitive.ObjectID `bson:"season_id" json:"season_id"`
	PlayersData []MatchPlayerData  `bson:"players_data" json:"players_data"`
	GamesPlayed int                `bson:"games_played" json:"games_played"`
	Gamemode    Gamemode           `bson:"gamemode" json:"gamemode"`
	Completed   bool               `bson:"completed" json:"completed"`
	CreatedAt   primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt   primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

type MatchPlayerData struct {
	TournamentPlayerID primitive.ObjectID `bson:"tournament_player_id" json:"tournament_player_id"`
	Wins               int                `bson:"wins" json:"wins"`
	Tags               []string           `bson:"tags" json:"tags"`
}

type Gamemode string

const (
	Standard       = "gm_1v1"
	Archenemy      = "gm_arc"
	Planechase     = "gm_hop"
	Commander      = "gm_edh"
	TwoHeadedGiant = "gm_2hg"
)
