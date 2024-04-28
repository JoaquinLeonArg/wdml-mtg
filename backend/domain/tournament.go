package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Tournaments collection
type Tournament struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	OwnerID         primitive.ObjectID `bson:"owner_id" json:"owner_id"`
	CurrentSeasonID primitive.ObjectID `bson:"current_season_id" json:"current_season_id"`
	InviteCode      string             `bson:"invite_code" json:"invite_code"`
	Name            string             `bson:"name" json:"name"`
	Description     string             `bson:"description" json:"description"`
	Store           Store              `bson:"store" json:"store"`
	CreatedAt       primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt       primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

type Store struct {
	BoosterPacks []StoreBoosterPack `bson:"booster_packs" json:"booster_packs"`
}

type StoreBoosterPack struct {
	BoosterPackID primitive.ObjectID `bson:"booster_pack_id" json:"booster_pack_id"`
	CoinPrice     int                `bson:"coin_price" json:"coin_price"`
}
