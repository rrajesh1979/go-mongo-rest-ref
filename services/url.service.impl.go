package services

import (
	"errors"
	"go-mongo-rest-ref/models"
	"go-mongo-rest-ref/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"time"
)

type URLServiceImpl struct {
	urlCollection *mongo.Collection
	ctx           context.Context
}

func (u URLServiceImpl) DeleteURL(shortURL string) (error, string) {
	filter := bson.M{"shortURL": shortURL}

	result, err := u.urlCollection.DeleteOne(u.ctx, filter)
	if err != nil {
		return err, "ERROR"
	}

	if result.DeletedCount == 0 {
		return errors.New("URL not found"), "NOT_FOUND"
	}

	return nil, "URL deleted"
}

func (u URLServiceImpl) FindURLsByUserID(userID string, page int, limit int) ([]*models.URL, error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	skip := (page - 1) * limit

	//Find URL
	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))

	query := bson.M{"userID": userID}

	cursor, err := u.urlCollection.Find(u.ctx, query, &opt)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(u.ctx)

	var urls []*models.URL

	for cursor.Next(u.ctx) {
		url := &models.URL{}
		err := cursor.Decode(url)

		if err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(urls) == 0 {
		return []*models.URL{}, nil
	}

	return urls, nil
}

func (u URLServiceImpl) FindURLs(page int, limit int) ([]*models.URL, error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	skip := (page - 1) * limit

	//Find URL
	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))

	query := bson.M{}

	cursor, err := u.urlCollection.Find(u.ctx, query, &opt)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(u.ctx)

	var urls []*models.URL

	for cursor.Next(u.ctx) {
		url := &models.URL{}
		err := cursor.Decode(url)

		if err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(urls) == 0 {
		return []*models.URL{}, nil
	}

	return urls, nil
}

func (u URLServiceImpl) CreateURL(r *models.CreateURLRequest) (*mongo.InsertOneResult, error) {
	//Create URL
	r.CreateAt = time.Now()
	r.UpdatedAt = r.CreateAt
	r.ExpiresAt = r.CreateAt.AddDate(0, 0, int(r.ExpirationDays))
	r.ShortUrl = utils.GenerateShortLink(r.LongURL, r.UserID)
	r.Status = "ACTIVE"
	r.Redirects = 0

	res, err := u.urlCollection.InsertOne(u.ctx, r)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("URL Already Exists")
		}
		return nil, err
	}

	return res, nil
}

func (u URLServiceImpl) UpdateURL(r *models.CreateURLRequest) (*mongo.UpdateResult, error) {
	//Create URL
	var updatedAt = time.Now()
	var expiresAt = updatedAt.AddDate(0, 0, int(r.ExpirationDays))
	var shortURL = utils.GenerateShortLink(r.LongURL, r.UserID)
	var status = "ACTIVE"
	var redirects = 0 //Reset Redirects

	filter := bson.M{"longURL": r.LongURL, "userID": r.UserID}
	update := bson.D{
		{"$set", bson.D{
			{"shortURL", shortURL},
			{"expirationDays", r.ExpirationDays},
			{"updatedAt", updatedAt},
			{"expiresAt", expiresAt},
			{"status", status},
			{"redirects", redirects}}}}

	result, err := u.urlCollection.UpdateOne(u.ctx, filter, update)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("error updating URL")
		}
		return nil, err
	}

	return result, nil
}

func NewURLService(urlsCollection *mongo.Collection, ctx context.Context) URLService {
	return &URLServiceImpl{urlsCollection, ctx}
}
