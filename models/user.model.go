package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CreateUserRequest struct {
	UserID    string    `json:"userID" bson:"userID" binding:"required"`
	UserName  string    `json:"userName" bson:"userName" binding:"required"`
	Email     string    `json:"email" bson:"email"  binding:"required"`
	UserTier  string    `json:"userTier" bson:"userTier"  binding:"required"`
	CreateAt  time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type User struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"userID" bson:"userID" binding:"required"`
	UserName  string             `json:"userName" bson:"userName" binding:"required"`
	Email     string             `json:"email" bson:"email"  binding:"required"`
	UserTier  string             `json:"userTier" bson:"userTier"  binding:"required"`
	CreateAt  time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
