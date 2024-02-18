package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Posts collection
type Post struct {
	ID primitive.ObjectID
}
