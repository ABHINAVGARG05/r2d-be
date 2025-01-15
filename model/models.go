package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}
