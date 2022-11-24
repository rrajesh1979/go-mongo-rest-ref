package services

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"go-mongo-rest-ref/config"
	"go-mongo-rest-ref/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)

func setupMongoDBContainer() (func(), *mongo.Client, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0.3",
		ExposedPorts: []string{"27017/tcp", "27018/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "root",
			"MONGO_INITDB_ROOT_PASSWORD": "example",
		},
	}

	mongodbC, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("error starting mongodb container: %s", err)
	}

	host, _ := mongodbC.Host(context.Background())
	p, _ := mongodbC.MappedPort(context.Background(), "27017/tcp")
	port := p.Int()

	closeContainer := func() {
		log.Printf("terminating container")
		if err := mongodbC.Terminate(context.Background()); err != nil {
			log.Fatalf("error terminating mongodb container: %v\n", err)
		}
	}

	uri := fmt.Sprintf("mongodb://%v:%v@%v:%v/",
		"root", "example", host, port)

	log.Printf("uri: %v", uri)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("cannot connect: %v\n", err)
		return closeContainer, nil, err
	}
	return closeContainer, client, nil
}

func TestURLService(t *testing.T) {
	ctx := context.TODO()
	config, err := config.LoadConfig("../")

	terminateC, mongoClient, err := setupMongoDBContainer()
	if err != nil {
		log.Fatalf("cannot setup mongodb: %v", err)
	}
	defer terminateC()

	urlCollection := mongoClient.Database(config.DBName).Collection(config.URLCollection)

	urlService := NewURLService(urlCollection, ctx)
	res, err := urlService.CreateURL(&models.CreateURLRequest{
		UserID:         "rrajesh1979",
		LongURL:        "https://www.google.com",
		ExpirationDays: 10,
	})
	log.Printf("res: %v", res.InsertedID)
	assert.NotEmpty(t, res.InsertedID)

	urls, _ := urlService.FindURLsByUserID("rrajesh1979", 1, 1)
	log.Printf("resU: %v", len(urls))
	assert.Equal(t, 1, len(urls))
}
