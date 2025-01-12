package models

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	Id     bson.ObjectID `bson: "_id" json: "id"`
	Name   string        `bson: "name" json: "name"`
	Gender string        `bson: "gender" json: "gender"`
	Age    int           `bson: "age" json: "age"`
}
