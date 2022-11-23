package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CreateURLRequest struct {
	LongURL        string    `json:"longURL" bson:"longURL" binding:"required"`
	ShortUrl       string    `json:"shortURL,omitempty" bson:"shortURL,omitempty"`
	ExpirationDays int64     `json:"expirationDays,omitempty" bson:"expirationDays,omitempty"`
	UserID         string    `json:"userID" bson:"userID"  binding:"required"`
	Status         string    `json:"status,omitempty" bson:"status,omitempty"`
	Redirects      int64     `json:"redirects,omitempty" bson:"redirects,omitempty"`
	ExpiresAt      time.Time `json:"expiresAt,omitempty" bson:"expiresAt,omitempty"`
	CreateAt       time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type URL struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	LongUrl   string             `json:"longURL" bson:"longURL" binding:"required"`
	ShortUrl  string             `json:"shortURL,omitempty" bson:"shortURL,omitempty"`
	UserId    string             `json:"userID" bson:"userID"  binding:"required"`
	Status    string             `json:"status,omitempty" bson:"status,omitempty"`
	Redirects int64              `json:"redirects,omitempty" bson:"redirects,omitempty"`
	ExpiresAt time.Time          `json:"expiresAt,omitempty" bson:"expiresAt,omitempty"`
	CreateAt  time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
