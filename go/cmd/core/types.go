package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	SessionId string `json:"sessionId"`
}

type Username struct {
	Username string `json:"username"`
}

type UserItem struct {
	ID            primitive.ObjectID `bson:"_id"`
	Username      string             `bson:"username"`
	AccountLinked bool               `bson:"tink_linked"`
}
