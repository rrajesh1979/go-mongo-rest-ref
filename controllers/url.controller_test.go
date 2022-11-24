package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"go-mongo-rest-ref/config"
	"go-mongo-rest-ref/models"
	"go-mongo-rest-ref/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go/types"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
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

func TestNewURLController(t *testing.T) {
	ctx := context.TODO()
	cfg, err := config.LoadConfig("../")

	terminateC, mongoClient, err := setupMongoDBContainer()
	if err != nil {
		log.Fatalf("cannot setup mongodb: %v", err)
	}
	defer terminateC()

	router := gin.Default()

	urlCollection := mongoClient.Database(cfg.DBName).Collection(cfg.URLCollection)
	urlService := services.NewURLService(urlCollection, ctx)
	urlController := NewURLController(urlService)

	// Setup Data in MongoDB
	setupTestData(urlService)

	findURLs := urlController.FindURLs
	findURLsByUserID := urlController.FindURLsByUserID
	deleteByShortURL := urlController.DeleteURL
	createURL := urlController.CreateURL
	updateURL := urlController.UpdateURL

	router.GET("/", findURLs)
	router.GET("/:userID", findURLsByUserID)
	router.DELETE("/:shortURL", deleteByShortURL)
	router.POST("/", createURL)
	router.PATCH("/", updateURL)

	var crA = []byte(`{
			"expirationDays": 15,
			"longURL": "https://www.google.com",
			"userID": "rrajesh1979"
	}`)

	var upA = []byte(`{
			"expirationDays": 55,
			"longURL": "https://www.google.com",
			"userID": "rrajesh1979"
	}`)

	allURLsRequest, _ := http.NewRequest("GET", "/", nil)
	urlsByUserIDRequestA, _ := http.NewRequest("GET", "/rrajesh1979?page=1&limit=5", nil)
	urlsByUserIDRequestB, _ := http.NewRequest("GET", "/johndoe?page=1&limit=1", nil)
	deleteURLRequestA, _ := http.NewRequest("DELETE", "/9qcffmSX", nil)
	createURLRequestA, _ := http.NewRequest("POST", "/", bytes.NewBuffer(crA))
	updateURLRequestA, _ := http.NewRequest("PATCH", "/", bytes.NewBuffer(upA))
	w := httptest.NewRecorder()

	// Test FindURLs
	router.ServeHTTP(w, allURLsRequest)
	allURLsResponse, _ := io.ReadAll(w.Body)
	var allURLsResponseJSON URLResponse
	err1 := json.Unmarshal(allURLsResponse, &allURLsResponseJSON)
	if err1 != nil {
		return
	}
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 4, allURLsResponseJSON.Results)

	// Test FindURLsByUserID
	router.ServeHTTP(w, urlsByUserIDRequestA)
	userURLsResponseA, _ := io.ReadAll(w.Body)
	var userURLsResponseJSONA URLResponse
	err2 := json.Unmarshal(userURLsResponseA, &userURLsResponseJSONA)
	if err2 != nil {
		return
	}
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 2, userURLsResponseJSONA.Results)

	// Test FindURLsByUserID
	router.ServeHTTP(w, urlsByUserIDRequestB)
	userURLsResponseB, _ := io.ReadAll(w.Body)
	var userURLsResponseJSONB URLResponse
	err3 := json.Unmarshal(userURLsResponseB, &userURLsResponseJSONB)
	if err3 != nil {
		return
	}
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1, userURLsResponseJSONB.Results)

	// Test DeleteURL
	router.ServeHTTP(w, deleteURLRequestA)
	deleteURLResponseA, _ := io.ReadAll(w.Body)
	var deleteURLResponseJSONA URLDeleteResponse
	err4 := json.Unmarshal(deleteURLResponseA, &deleteURLResponseJSONA)
	if err4 != nil {
		return
	}
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", deleteURLResponseJSONA.Status)
	assert.Equal(t, "URL deleted", deleteURLResponseJSONA.Data)

	//Test CreateURL
	router.ServeHTTP(w, createURLRequestA)
	createURLResponseA, _ := io.ReadAll(w.Body)
	var createURLResponseJSONA URLCreateResponse
	err5 := json.Unmarshal(createURLResponseA, &createURLResponseJSONA)
	if err5 != nil {
		return
	}
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", createURLResponseJSONA.Status)
	assert.NotEmpty(t, createURLResponseJSONA.Data)

	//Test UpdateURL
	router.ServeHTTP(w, updateURLRequestA)
	updateURLResponseA, _ := io.ReadAll(w.Body)
	var updateURLResponseJSONA URLCreateResponse
	err6 := json.Unmarshal(updateURLResponseA, &updateURLResponseJSONA)
	if err6 != nil {
		return
	}
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", updateURLResponseJSONA.Status)
	assert.Equal(t, "URL updated", updateURLResponseJSONA.Data)

}

func setupTestData(urlService services.URLService) {
	_, errA := urlService.CreateURL(&models.CreateURLRequest{
		UserID:         "rrajesh1979",
		LongURL:        "https://www.google.com",
		ExpirationDays: 10,
	})
	if errA != nil {
		return
	}

	_, errB := urlService.CreateURL(&models.CreateURLRequest{
		UserID:         "rrajesh1979",
		LongURL:        "https://www.mongodb.com",
		ExpirationDays: 20,
	})
	if errB != nil {
		return
	}

	_, errC := urlService.CreateURL(&models.CreateURLRequest{
		UserID:         "johndoe",
		LongURL:        "https://www.mongodb.com",
		ExpirationDays: 30,
	})
	if errC != nil {
		return
	}

	_, errD := urlService.CreateURL(&models.CreateURLRequest{
		UserID:         "johndoe",
		LongURL:        "https://www.google.com",
		ExpirationDays: 40,
	})
	if errD != nil {
		return
	}

}

type URLResponse struct {
	Data    types.Array `json:"data" binding:"required"`
	Results int         `json:"results"`
	Status  string      `json:"status" binding:"required"`
}
type URLDeleteResponse struct {
	Data   string `json:"data" binding:"required"`
	Status string `json:"status" binding:"required"`
}
type URLCreateResponse struct {
	Data   string `json:"data" binding:"required"`
	Status string `json:"status" binding:"required"`
}
