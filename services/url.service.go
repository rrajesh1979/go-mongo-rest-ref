package services

import (
	"go-mongo-rest-ref/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type URLService interface {
	FindURLs(page int, limit int) ([]*models.URL, error)
	FindURLsByUserID(userID string, page int, limit int) ([]*models.URL, error)
	CreateURL(url *models.CreateURLRequest) (*mongo.InsertOneResult, error)
	DeleteURL(shortURL string) (error, string)
	UpdateURL(*models.CreateURLRequest) (*mongo.UpdateResult, error)
}
