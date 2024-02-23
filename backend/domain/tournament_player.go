package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// TournamentPlayers collection
type TournamentPlayer struct {
	ID               primitive.ObjectID `bson:"_id"`
	UserID           primitive.ObjectID `bson:"user_id"`
	TournamentID     primitive.ObjectID `bson:"tournament_id"`
	AccessLevel      AccessLevel        `bson:"access_level"`
	GameResources    GameResources      `bson:"game_resources"`
	TournamentPoints int                `bson:"tournament_points"`
	CreatedAt        primitive.DateTime `bson:"created_at"`
	UpdatedAt        primitive.DateTime `bson:"updated_at"`
}

type AccessLevel string

const (
	AccessLevelPlayer        AccessLevel = "al_player"
	AccessLevelModerator     AccessLevel = "al_moderator"
	AccessLevelAdministrator AccessLevel = "al_administrator"
)

type GameResources struct {
	OwnedCards   []Card             `bson:"owned_cards"`
	Decks        []Deck             `bson:"decks"`
	Wildcards    OwnedWildcards     `bson:"wildcards"`
	BoosterPacks []OwnedBoosterPack `bson:"booster_packs"`
	Rerolls      int                `bson:"rerolls"`
	Coins        int                `bson:"coins"`
}

type Card struct {
	SetCode         string             `bson:"set_code"`
	CollectorNumber int                `bson:"collector_number"`
	Count           int                `bson:"count"`
	CardData        CardData           `bson:"card_data"`
	CreatedAt       primitive.DateTime `bson:"created_at"`
	UpdatedAt       primitive.DateTime `bson:"updated_at"`
}

type CardData struct {
	Name       string   `bson:"name"`
	SuperTypes []string `bson:"super_types"`
	Types      []string `bson:"types"`
	SubTypes   []string `bson:"sub_types"`
	ManaValue  int      `bson:"mana_value"`
	Colors     []Color  `bson:"colors"`
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
	CommonCount     int `bson:"common_count"`
	UncommonCount   int `bson:"uncommon_count"`
	RareCount       int `bson:"rare_count"`
	MythicRareCount int `bson:"mythic_rare_count"`

	// Special cards
	MasterpieceCount int `bson:"masterpiece_count"`
}

type OwnedBoosterPack struct {
	SetCode        string      `bson:"set_code"`
	SetName        string      `bson:"set_name"`
	BoosterType    BoosterType `bson:"booster_type"`
	BoosterGen     BoosterGen  `bson:"booster_gen"`
	BoosterGenData interface{} `bson:"booster_gen_data"`
	Available      int         `bson:"available"`
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
	CardPool []Card `bson:"card_pool"`
	// TODO: Extra customization
}

type Deck struct {
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Cards       []Card             `bson:"cards"`
	CreatedAt   primitive.DateTime `bson:"created_at"`
	UpdatedAt   primitive.DateTime `bson:"updated_at"`
}
