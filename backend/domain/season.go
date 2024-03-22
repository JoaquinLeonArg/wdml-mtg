package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Season struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	TournamentID primitive.ObjectID `bson:"tournament_id" json:"tournament_id"`
	Name         string             `bson:"name" json:"name"`
	Description  string             `bson:"description" json:"description"`
	CreatedAt    primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt    primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
