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

// CreateURL
//
//	@Summary		Create new Short URL
//	@Description	Create new Short URL
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.CreateURLRequest	true	"Create URL Request"
//	@Success		200		{object}	string
//	@Failure		400		{object}	string
//	@Failure		404		{object}	string
//	@Failure		500		{object}	string
//	@Router			/ [post]
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

// FindURLs
//
//	@Summary		Get all URLs
//	@Description	get all URLs
//	@Produce		json
//	@Success		200	{array}		models.URL
//	@Failure		400	{object}	string
//	@Failure		404	{object}	string
//	@Failure		500	{object}	string
//	@Router			/ [get]
func (uc *URLController) FindURLs(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "data": err.Error()})
		return
	}
	intLimit, err := strconv.Atoi(limit)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "data": err.Error()})
		return
	}

	urls, err := uc.urlService.FindURLs(intPage, intLimit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(urls), "data": urls})

}

// FindURLsByUserID
//
//	@Summary		Get all URLs by User ID
//	@Description	Get all URLs by User ID
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string	true	"User ID"
//	@Success		200		{array}		models.URL
//	@Failure		400		{object}	string
//	@Failure		404		{object}	string
//	@Failure		500		{object}	string
//	@Router			/{userID} [get]
func (uc *URLController) FindURLsByUserID(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")
	var userID = ctx.Param("userID")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "data": err.Error()})
		return
	}
	intLimit, err := strconv.Atoi(limit)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "data": err.Error()})
		return
	}

	urls, err := uc.urlService.FindURLsByUserID(userID, intPage, intLimit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(urls), "data": urls})

}

// FindLongURL
//
//	@Summary		Find Long URL to redirect
//	@Description	Find Long URL to redirect
//	@Accept			json
//	@Produce		json
//	@Param			shortURL	path		string	true	"Short URL"
//	@Success		200			{array}		models.URL
//	@Failure		400			{object}	string
//	@Failure		404			{object}	string
//	@Failure		500			{object}	string
//	@Router			/get/{shortURL} [get]
func (uc *URLController) FindLongURL(ctx *gin.Context) {
	var shortURL = ctx.Param("shortURL")

	longURL := uc.urlService.FindLongURL(shortURL)
	if longURL.LongUrl == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "URL not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": longURL.LongUrl})
}

// RedirectLongURL
//
//	@Summary		Redirect to LongURL
//	@Description	Redirect to LongURL
//	@Accept			json
//	@Produce		json
//	@Param			shortURL	path		string	true	"Short URL"
//	@Success		200			{array}		models.URL
//	@Failure		400			{object}	string
//	@Failure		404			{object}	string
//	@Failure		500			{object}	string
//	@Router			/redirect/{shortURL} [get]
func (uc *URLController) RedirectLongURL(ctx *gin.Context) {
	var shortURL = ctx.Param("shortURL")

	longURL := uc.urlService.FindLongURL(shortURL)

	//Check if URL is nil
	if longURL == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "URL not found"})
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, longURL.LongUrl)
}

// DeleteURL
//
//	@Summary		Delete URL
//	@Description	Delete by ShortURL
//	@Accept			json
//	@Produce		json
//	@Param			shortURL	path		string	true	"Short URL"
//	@Success		200			{object}	string
//	@Failure		400			{object}	string
//	@Failure		404			{object}	string
//	@Failure		500			{object}	string
//	@Router			/{shortURL} [delete]
func (uc *URLController) DeleteURL(ctx *gin.Context) {
	var shortURL = ctx.Param("shortURL")

	err, status := uc.urlService.DeleteURL(shortURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": status})
}

// UpdateURL
//
//	@Summary		Update URL
//	@Description	Update URL
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.CreateURLRequest	true	"Create URL Request"
//	@Success		200		{object}	string
//	@Failure		400		{object}	string
//	@Failure		404		{object}	string
//	@Failure		500		{object}	string
//	@Router			/ [patch]
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
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "data": "URL not found"})
		return
	} else if res.ModifiedCount != 0 {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": "URL updated"})
		return
	}
}
