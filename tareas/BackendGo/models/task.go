package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
    ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Title            string             `json:"title" bson:"title"`
    Description      string             `json:"description" bson:"description"`
    Completed         bool               `json:"completed" bson:"completed"`
    CreatedAt        primitive.DateTime `json:"created_at" bson:"created_at"`
}