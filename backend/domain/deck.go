package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Deck struct {
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
	TournamentPlayerID primitive.ObjectID `bson:"tournament_player_id" json:"tournament_player_id"`
	Name               string             `bson:"name" json:"name"`
	Description        string             `bson:"description" json:"description"`
	Cards              []DeckCard         `bson:"cards" json:"cards"`
	CreatedAt          primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt          primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

type DeckCard struct {
	OwnedCardID primitive.ObjectID `bson:"owned_card_id" json:"owned_card_id"`
	Count       int                `bson:"count" json:"count"`
	Board       DeckBoard          `bson:"board" josn:"board"`
}

type DeckBoard string

const (
	SideBoard  MatchResult = "b_sideboard"
	MainBoard  MatchResult = "b_mainboard"
	MaybeBoard MatchResult = "b_maybeboard"
)
