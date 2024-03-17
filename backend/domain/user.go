package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Users collection
type User struct {
	ID                primitive.ObjectID `bson:"_id" json:"id"`
	Username          string             `bson:"username" json:"username"`
	Email             string             `bson:"email" json:"email"`
	Password          []byte             `bson:"password" json:"password"`
	Description       string             `bson:"description" json:"description"`
	ProfilePictureURL string             `bson:"profile_picture_url" json:"profile_picture_url"`
	CreatedAt         primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt         primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
