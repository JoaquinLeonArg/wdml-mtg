package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Users collection
type User struct {
	ID                primitive.ObjectID `bson:"_id"`
	Username          string             `bson:"username"`
	Email             string             `bson:"email"`
	Password          string             `bson:"password"`
	Description       string             `bson:"description"`
	ProfilePictureURL string             `bson:"profile_picture_url"`
	CreatedAt         primitive.DateTime `bson:"created_at"`
	UpdatedAt         primitive.DateTime `bson:"updated_at"`
}
