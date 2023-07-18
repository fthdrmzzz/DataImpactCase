package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents the user model
type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Password   string             `bson:"password" json:"password"`
	IsActive   bool               `bson:"isActive" json:"isActive"`
	Balance    string             `bson:"balance" json:"balance"`
	Age        string             `bson:"age" json:"age"`
	Name       string             `bson:"name" json:"name"`
	Gender     string             `bson:"gender" json:"gender"`
	Company    string             `bson:"company" json:"company"`
	Email      string             `bson:"email" json:"email"`
	Phone      string             `bson:"phone" json:"phone"`
	Address    string             `bson:"address" json:"address"`
	About      string             `bson:"about" json:"about"`
	Registered string             `bson:"registered" json:"registered"`
	Latitude   float64            `bson:"latitude" json:"latitude"`
	Longitude  float64            `bson:"longitude" json:"longitude"`
	Tags       []string           `bson:"tags" json:"tags"`
	Friends    []Friend           `bson:"friends" json:"friends"`
	Data       string             `bson:"data" json:"data"`
}

// Friend represents the friend model
type Friend struct {
	ID   int    `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}
