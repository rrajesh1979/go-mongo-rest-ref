package test

import (
	"github.com/stretchr/testify/assert"
	"go-mongo-rest-ref/utils"
	"testing"
)

const userID = "rrajesh1979"

func TestURLShort(t *testing.T) {
	//t.Run("Create URL", testCreateURL)
	//t.Run("Find URL", testFindURL)
	//t.Run("Find URL by User ID", testFindURLByUserID)
	//t.Run("Update URL", testUpdateURL)
	//t.Run("Delete URL", testDeleteURL)
	longURL := "https://www.google.com"
	shortURL := utils.GenerateShortLink(longURL, userID)
	assert.Equal(t, shortURL, "9qcffmSX")
}
