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
	Decks        []Deck             `bson:"decks" json:"decks"`
	Wildcards    OwnedWildcards     `bson:"wildcards" json:"wildcards"`
	BoosterPacks []OwnedBoosterPack `bson:"booster_packs" json:"booster_packs"`
	Rerolls      int                `bson:"rerolls" json:"rerolls"`
	Coins        int                `bson:"coins" json:"coins"`
}

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
	BoosterGen     BoosterGen      `bson:"booster_gen" json:"booster_gen"`
	BoosterGenData interface{}     `bson:"booster_gen_data" json:"booster_gen_data"`
	Available      int             `bson:"available" json:"available"`
	Data           BoosterPackData `bson:"data" json:"data"`
}

type BoosterPackData struct {
	SetCode     string      `bson:"set_code" json:"set_code"`
	SetName     string      `bson:"set_name" json:"set_name"`
	BoosterType BoosterType `bson:"booster_type" json:"booster_type"`
	Expansion   string      `bson:"expansion" json:"expansion"`
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
	CardPool []OwnedCard `bson:"card_pool" json:"card_pool"`
	// TODO: Extra customization
}

type Deck struct {
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Cards       []OwnedCard        `bson:"cards" json:"cards"`
	CreatedAt   primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt   primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
