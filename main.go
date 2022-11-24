package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-mongo-rest-ref/config"
	"go-mongo-rest-ref/controllers"
	"go-mongo-rest-ref/routes"
	"go-mongo-rest-ref/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client

	urlShortService    services.URLService
	urlShortController controllers.URLController
	urlCollection      *mongo.Collection
	URLRouteController routes.URLRouteController
)

func init() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()

	// Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	urlCollection = mongoclient.Database(config.DBName).Collection(config.URLCollection)
	urlShortService = services.NewURLService(urlCollection, ctx)
	urlShortController = controllers.NewURLController(urlShortService)
	URLRouteController = routes.NewURLControllerRoute(urlShortController)

	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoclient.Disconnect(ctx)

	startGinServer(config)
}

func startGinServer(config config.Config) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.Origin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group(config.APIDefaultPath + config.APIVersion)
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	URLRouteController.URLRoute(router)
	log.Fatal(server.Run(":" + config.Port))
}
