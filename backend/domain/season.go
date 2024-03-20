package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Match struct {
	PlayerAID   primitive.ObjectID `bson:"player_a_id" json:"player_a_id"`
	PlayerBID   primitive.ObjectID `bson:"player_b_id" json:"player_b_id"`
	PlayerAWins int                `bson:"player_a_wins" json:"player_a_wins"`
	PlayerBWins int                `bson:"player_b_wins" json:"player_b_wins"`
	Completed   bool               `bson:"completed" json:"completed"`
}

type Season struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	TournamentID primitive.ObjectID `bson:"tournament_id" json:"tournament_id"`
	Matches      []Match            `bson:"matches" json:"matches"`
	Name         string             `bson:"name" json:"name"`
	Description  string             `bson:"description" json:"description"`
	CreatedAt    primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt    primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
