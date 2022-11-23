package routes

import (
	"github.com/gin-gonic/gin"
	swagFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-mongo-rest-ref/config"
	"go-mongo-rest-ref/controllers"
	"go-mongo-rest-ref/docs"
	"log"
)

type URLRouteController struct {
	urlController controllers.URLController
}

func NewURLControllerRoute(urlController controllers.URLController) URLRouteController {
	return URLRouteController{urlController}
}

func (r *URLRouteController) URLRoute(rg *gin.RouterGroup) {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	router := rg.Group(cfg.APIShort)
	//log.Fatal(router)

	router.GET("/", r.urlController.FindURLs)
	router.POST("/", r.urlController.CreateURL)
	router.GET("/:userID", r.urlController.FindURLsByUserID)
	router.DELETE("/:shortURL", r.urlController.DeleteURL)
	router.PATCH("/", r.urlController.UpdateURL)

	//r := gin.Default()
	docs.SwaggerInfo.BasePath = cfg.APIDefaultPath + cfg.APIVersion
	docs.SwaggerInfo.Title = "GoLang MongoDB REST API Swagger Example API"
	docs.SwaggerInfo.Description = "GoLang MongoDB REST API Swagger Example API"
	docs.SwaggerInfo.Version = cfg.APIVersion
	docs.SwaggerInfo.Host = cfg.Host + ":" + cfg.Port
	docs.SwaggerInfo.BasePath = cfg.APIDefaultPath + cfg.APIVersion
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swagFiles.Handler))
}
