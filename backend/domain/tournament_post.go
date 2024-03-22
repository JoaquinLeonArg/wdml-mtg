package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type TournamentPost struct {
	ID                 primitive.ObjectID       `bson:"_id" json:"id"`
	TournamentID       primitive.ObjectID       `bson:"tournament_id" json:"tournament_id"`
	TournamentPlayerID primitive.ObjectID       `bson:"tournament_player_id" json:"tournament_player_id"`
	Title              string                   `bson:"title" json:"title"`
	Blocks             []TournamentPostBlock    `bson:"blocks" json:"blocks"`
	Reactions          []TournamentPostReaction `bson:"reactions" json:"reactions"`
	CreatedAt          primitive.DateTime       `bson:"created_at" json:"created_at"`
	UpdatedAt          primitive.DateTime       `bson:"updated_at" json:"updated_at"`
}

type TournamentPostBlock struct {
	Collapsable bool                    `bson:"collapsable" json:"collapsable"`
	Title       string                  `bson:"title" json:"title"`
	Content     []TournamentPostContent `bson:"content" json:"content"`
}

type TournamentPostContent struct {
	Type      TournamentPostType `bson:"type" json:"type"`
	Content   string             `bson:"content" json:"content"`
	ExtraData map[string]string  `bson:"extra_data" json:"extra_data"`
}

type TournamentPostType string

const (
	TournamentPostTypeText  TournamentPostType = "tpt_text"
	TournamentPostTypeImage TournamentPostType = "tpt_image"
)

type TournamentPostReaction struct {
	UserID       primitive.ObjectID         `bson:"user_id" json:"user_id"`
	ReactionType TournamentPostReactionType `bson:"reaction_type" json:"reaction_type"`
}

type TournamentPostReactionType string

const (
	TournamentPostReactionTypeHeart TournamentPostReactionType = "tprt_heart"
	TournamentPostReactionTypeSkull TournamentPostReactionType = "tprt_skull"
)
