package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// TournamentBlocks collection
type TournamentBlock struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	TournamentID primitive.ObjectID `bson:"tournament_id" json:"tournament_id"`
	BlockNumber  int                `bson:"block_number" json:"block_number"`
	Status       BlockStatus        `bson:"status" json:"status"`
}

type BlockStatus string

const (
	BlockStatusEnded    BlockStatus = "bs_ended"
	BlockStatusOngoing  BlockStatus = "bs_ongoing"
	BlockStatusUpcoming BlockStatus = "bs_upcoming"
)
