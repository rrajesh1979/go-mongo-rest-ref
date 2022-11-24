package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const userID = "rrajesh1979"

func TestURLShort(t *testing.T) {
	longURL := "https://www.google.com"
	shortURL := GenerateShortLink(longURL, userID)
	assert.Equal(t, shortURL, "9qcffmSX")
}
