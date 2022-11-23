package routes

import (
	"github.com/gin-gonic/gin"
	"go-mongo-rest-ref/config"
	"go-mongo-rest-ref/controllers"
	"log"
)

type URLRouteController struct {
	urlController controllers.URLController
}

func NewURLControllerRoute(urlController controllers.URLController) URLRouteController {
	return URLRouteController{urlController}
}

func (r *URLRouteController) URLRoute(rg *gin.RouterGroup) {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	router := rg.Group(config.APIShort)
	//log.Fatal(router)

	router.GET("/", r.urlController.FindURLs)
	router.POST("/", r.urlController.CreateURL)
	router.GET("/:userID", r.urlController.FindURLsByUserID)
	router.DELETE("/:shortURL", r.urlController.DeleteURL)
	router.PATCH("/", r.urlController.UpdateURL)
}
