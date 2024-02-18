package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Cards collection
type Card struct {
	ID        primitive.ObjectID
	SetCode   string
	SetNumber int
	Name      string
	ImageURL  string
}
