package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// TournamentPlayers collection
type TournamentPlayer struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`
	TournamentID     primitive.ObjectID `bson:"tournament_id" json:"tournament_id"`
	AccessLevel      AccessLevel        `bson:"access_level" json:"access_level"`
	GameResources    GameResources      `bson:"game_resources" json:"game_resources"`
	TournamentPoints int                `bson:"tournament_points" json:"tournament_points"`
	CreatedAt        primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt        primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

type AccessLevel string

const (
	AccessLevelPlayer        AccessLevel = "al_player"
	AccessLevelModerator     AccessLevel = "al_moderator"
	AccessLevelAdministrator AccessLevel = "al_administrator"
)

type GameResources struct {
	OwnedCards   []Card             `bson:"owned_cards" json:"owned_cards"`
	Decks        []Deck             `bson:"decks" json:"decks"`
	Wildcards    OwnedWildcards     `bson:"wildcards" json:"wildcards"`
	BoosterPacks []OwnedBoosterPack `bson:"booster_packs" json:"booster_packs"`
	Rerolls      int                `bson:"rerolls" json:"rerolls"`
	Coins        int                `bson:"coins" json:"coins"`
}

type Card struct {
	SetCode         string             `bson:"set_code" json:"set_code"`
	CollectorNumber int                `bson:"collector_number" json:"collector_number"`
	Count           int                `bson:"count" json:"count"`
	CardData        CardData           `bson:"card_data" json:"card_data"`
	CreatedAt       primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt       primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

type CardData struct {
	Name       string   `bson:"name" json:"name"`
	SuperTypes []string `bson:"super_types" json:"super_types"`
	Types      []string `bson:"types" json:"types"`
	SubTypes   []string `bson:"sub_types" json:"sub_types"`
	ManaValue  int      `bson:"mana_value" json:"mana_value"`
	Colors     []Color  `bson:"colors" json:"colors"`
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
	CommonCount     int `bson:"common_count" json:"common_count"`
	UncommonCount   int `bson:"uncommon_count" json:"uncommon_count"`
	RareCount       int `bson:"rare_count" json:"rare_count"`
	MythicRareCount int `bson:"mythic_rare_count" json:"mythic_rare_count"`

	// Special cards
	MasterpieceCount int `bson:"masterpiece_count" json:"masterpiece_count"`
}

type OwnedBoosterPack struct {
	SetCode        string      `bson:"set_code" json:"set_code"`
	SetName        string      `bson:"set_name" json:"set_name"`
	BoosterType    BoosterType `bson:"booster_type" json:"booster_type"`
	BoosterGen     BoosterGen  `bson:"booster_gen" json:"booster_gen"`
	BoosterGenData interface{} `bson:"booster_gen_data" json:"booster_gen_data"`
	Available      int         `bson:"available" json:"available"`
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
	CardPool []Card `bson:"card_pool" json:"card_pool"`
	// TODO: Extra customization
}

type Deck struct {
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Cards       []Card             `bson:"cards" json:"cards"`
	CreatedAt   primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt   primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
