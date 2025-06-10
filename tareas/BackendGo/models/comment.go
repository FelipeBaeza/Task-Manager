package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TaskID    primitive.ObjectID `json:"taskId" bson:"taskId"`
	AuthorID  primitive.ObjectID `json:"authorId" bson:"authorId"`
	Text      string             `json:"text" bson:"text"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
}
