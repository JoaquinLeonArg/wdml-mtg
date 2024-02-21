package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Tournaments collection
type Tournament struct {
	ID          primitive.ObjectID
	OwnerID     primitive.ObjectID
	Title       string
	Description string
	CreatedAt   primitive.DateTime
	UpdatedAt   primitive.DateTime
}
