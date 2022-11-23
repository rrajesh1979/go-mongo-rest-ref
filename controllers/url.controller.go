package controllers

import (
	"go-mongo-rest-ref/models"
	"go-mongo-rest-ref/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type URLController struct {
	urlService services.URLService
}

func NewURLController(urlService services.URLService) URLController {
	return URLController{urlService}
}

func (uc *URLController) CreateURL(ctx *gin.Context) {
	var cure *models.CreateURLRequest
	if err := ctx.ShouldBindJSON(&cure); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := uc.urlService.CreateURL(cure)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": res.InsertedID})
}

func (uc *URLController) FindURLs(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	intLimit, err := strconv.Atoi(limit)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	urls, err := uc.urlService.FindURLs(intPage, intLimit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(urls), "data": urls})

}

func (uc *URLController) FindURLsByUserID(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")
	var userID = ctx.Param("userID")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	intLimit, err := strconv.Atoi(limit)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	urls, err := uc.urlService.FindURLsByUserID(userID, intPage, intLimit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(urls), "data": urls})

}

func (uc *URLController) DeleteURL(ctx *gin.Context) {
	var shortURL = ctx.Param("shortURL")

	err, status := uc.urlService.DeleteURL(shortURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": status})
}

func (uc *URLController) UpdateURL(ctx *gin.Context) {
	var cure *models.CreateURLRequest
	if err := ctx.ShouldBindJSON(&cure); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := uc.urlService.UpdateURL(cure)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if res.MatchedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "URL not found"})
		return
	} else if res.ModifiedCount != 0 {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "URL updated"})
		return
	}
}
