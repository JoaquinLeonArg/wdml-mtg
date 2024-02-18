package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// TournamentPlayers collection
type TournamentPlayer struct {
	ID               primitive.ObjectID
	UserID           primitive.ObjectID
	TournamentID     primitive.ObjectID
	AccessLevel      AccessLevel
	GameResources    GameResources
	TournamentPoints int
	CreatedAt        primitive.DateTime
	UpdatedAt        primitive.DateTime
}

type AccessLevel string

const (
	AccessLevelPlayer        AccessLevel = "al_player"
	AccessLevelModerator     AccessLevel = "al_moderator"
	AccessLevelAdministrator AccessLevel = "al_administrator"
)

type GameResources struct {
	OwnedCards   []OwnedCard
	Decks        []Deck
	Wildcards    OwnedWildcards
	BoosterPacks []OwnedBoosterPack
	Rerolls      int
	Coins        int
}

type OwnedCard struct {
	CardID    primitive.ObjectID
	Count     int
	CreatedAt primitive.DateTime
	UpdatedAt primitive.DateTime
}

type OwnedWildcards struct {
	// By rarity
	CommonCount     int
	UncommonCount   int
	RareCount       int
	MythicRareCount int

	// Special cards
	MasterpieceCount int
}

type OwnedBoosterPack struct {
	SetCode   string
	Available int
}

type Deck struct {
	Name        string
	Description string
	CardIDs     []primitive.ObjectID
	CreatedAt   primitive.DateTime
	UpdatedAt   primitive.DateTime
}
