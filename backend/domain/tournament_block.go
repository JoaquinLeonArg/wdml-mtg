package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// TournamentBlocks collection
type TournamentBlock struct {
	ID           primitive.ObjectID
	TournamentID primitive.ObjectID
	BlockNumber  int
	Status       BlockStatus
}

type BlockStatus string

const (
	BlockStatusEnded    BlockStatus = "bs_ended"
	BlockStatusOngoing  BlockStatus = "bs_ongoing"
	BlockStatusUpcoming BlockStatus = "bs_upcoming"
)
