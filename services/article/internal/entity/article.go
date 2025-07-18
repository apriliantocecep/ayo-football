package entity

import "go.mongodb.org/mongo-driver/v2/bson"

type Article struct {
	ID      bson.ObjectID `bson:"_id"`
	Content string        `json:"content"`
	UserId  string        `json:"user_id" bson:"user_id"`
	Status  string        `json:"status"`
}
