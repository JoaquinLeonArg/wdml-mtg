package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type OwnedCard struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	TournamentID primitive.ObjectID `bson:"tournament_id" json:"tournament_id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	Count        int                `bson:"count" json:"count"`
	CardData     CardData           `bson:"card_data" json:"card_data"`
	CreatedAt    primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt    primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

type CardData struct {
	SetCode         string   `bson:"set_code" json:"set_code"`
	CollectorNumber string   `bson:"collector_number" json:"collector_number"`
	Name            string   `bson:"name" json:"name"`
	Types           []string `bson:"types" json:"types"`
	ManaValue       int      `bson:"mana_value" json:"mana_value"`
	Colors          []string `bson:"colors" json:"colors"`
	ImageURL        string   `bson:"image_url" json:"image_url"`
	BackImageURL    string   `bson:"back_image_url" json:"back_image_url"`
}