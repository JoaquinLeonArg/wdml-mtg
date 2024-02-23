package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Tournaments collection
type Tournament struct {
	ID          primitive.ObjectID `bson:"_id"`
	OwnerID     primitive.ObjectID `bson:"owner_id"`
	InviteCode  string             `bson:"invite_code"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	CreatedAt   primitive.DateTime `bson:"created_at"`
	UpdatedAt   primitive.DateTime `bson:"updated_at"`
}
