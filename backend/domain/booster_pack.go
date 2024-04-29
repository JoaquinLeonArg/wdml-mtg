package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// BoosterPacks collection
type BoosterPack struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	SetCode     string             `bson:"set_code" json:"set_code"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	CardCount   int                `bson:"card_count" json:"card_count"`
	Filter      string             `bson:"filter" json:"filter"`
	Slots       []BoosterPackSlot  `bson:"slots" json:"slots"`
	CreatedAt   primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt   primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

type BoosterPackSlot struct {
	Options []Option `bson:"options" json:"options"`
	Filter  string   `bson:"filter" json:"filter"`
	Count   int      `bson:"count" json:"count"`
}

type Option struct {
	Filter string `bson:"filter" json:"filter"`
	Weight int    `bson:"weight" json:"weight"`
}
