package test

import (
	"github.com/stretchr/testify/assert"
	"go-mongo-rest-ref/utils"
	"testing"
)

const userID = "rrajesh1979"

func TestURLShort(t *testing.T) {
	longURL := "https://www.google.com"
	shortURL := utils.GenerateShortLink(longURL, userID)
	assert.Equal(t, shortURL, "9qcffmSX")
}
