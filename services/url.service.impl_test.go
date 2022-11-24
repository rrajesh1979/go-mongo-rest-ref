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
	cfg, err := config.LoadConfig("../")

	terminateC, mongoClient, err := setupMongoDBContainer()
	if err != nil {
		log.Fatalf("cannot setup mongodb: %v", err)
	}
	defer terminateC()

	urlCollection := mongoClient.Database(cfg.DBName).Collection(cfg.URLCollection)
	urlService := NewURLService(urlCollection, ctx)

	// Test CreateURL
	resIA, err := urlService.CreateURL(&models.CreateURLRequest{
		UserID:         "rrajesh1979",
		LongURL:        "https://www.google.com",
		ExpirationDays: 10,
	})
	log.Printf("res: %v", resIA.InsertedID)
	assert.NotEmpty(t, resIA.InsertedID)

	resIB, err := urlService.CreateURL(&models.CreateURLRequest{
		UserID:         "johndoe",
		LongURL:        "https://www.mongodb.com",
		ExpirationDays: 20,
	})
	log.Printf("res: %v", resIB.InsertedID)
	assert.NotEmpty(t, resIB.InsertedID)

	// Test FindURLsByUserID
	urlsA, _ := urlService.FindURLsByUserID("rrajesh1979", 1, 1)
	log.Printf("resU: %v", len(urlsA))
	assert.Equal(t, 1, len(urlsA))

	// Test FindURLByShortURL
	urlsB, _ := urlService.FindURLsByUserID("johndoe", 0, 0)
	log.Printf("resU: %v", len(urlsB))
	assert.Equal(t, 1, len(urlsB))

	// Test FindURLs
	urlsC, _ := urlService.FindURLs(1, 2)
	log.Printf("resU: %v", len(urlsC))
	assert.Equal(t, 2, len(urlsC))

	// Test DeleteURL
	_, resD := urlService.DeleteURL("9qcffmSX")
	log.Printf("resD: %v", resD)
	assert.Equal(t, resD, "URL deleted")

	// Test UpdateURL
	resE, _ := urlService.UpdateURL(&models.CreateURLRequest{
		UserID:         "johndoe",
		LongURL:        "https://www.mongodb.com",
		ExpirationDays: 30,
	})
	log.Printf("resE: %v", resE)
	assert.Equal(t, resE.ModifiedCount, int64(1))

}
