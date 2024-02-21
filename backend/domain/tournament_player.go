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
	OwnedCards   []Card
	Decks        []Deck
	Wildcards    OwnedWildcards
	BoosterPacks []OwnedBoosterPack
	Rerolls      int
	Coins        int
}

type Card struct {
	SetCode         string
	CollectorNumber int
	Count           int
	CardData        CardData
	CreatedAt       primitive.DateTime
	UpdatedAt       primitive.DateTime
}

type CardData struct {
	Name       string
	SuperTypes []string
	Types      []string
	SubTypes   []string
	ManaValue  int
	Colors     []Color
}

type Color string

const (
	ColorWhite Color = "c_white"
	ColorBlue  Color = "c_blue"
	ColorBlack Color = "c_black"
	ColorRed   Color = "c_red"
	ColorGreen Color = "c_green"
)

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
	SetCode        string
	SetName        string
	BoosterType    BoosterType
	BoosterGen     BoosterGen
	BoosterGenData interface{}
	Available      int
}

type BoosterType string

const (
	BoosterTypeDraft        BoosterType = "bt_draft"
	BoosterTypeBlock        BoosterType = "bt_block"
	BoosterTypeMasterpieces BoosterType = "bt_masterpieces"
	BoosterTypeLands        BoosterType = "bt_lands"
	BoosterTypeOther        BoosterType = "bt_other"
)

type BoosterGen string

const (
	BoosterGenVanilla BoosterGen = "bg_vanilla"
	BoosterGenCustom  BoosterGen = "bg_custom"
)

type BoosterGenDataCustom struct {
	CardPool []Card
}

type Deck struct {
	Name        string
	Description string
	Cards       []primitive.ObjectID
	CreatedAt   primitive.DateTime
	UpdatedAt   primitive.DateTime
}
