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

const (
	SPECIAL_TO_COIN  = 10
	MYTHIC_TO_COIN   = 20
	RARE_TO_COIN     = 10
	UNCOMMON_TO_COIN = 3
	COMMON_TO_COIN   = 1
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
	Available   int    `bson:"available" json:"available"`
	SetCode     string `bson:"set_code" json:"set_code"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}
