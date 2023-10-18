package aes

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAes(t *testing.T) {
	aesProvider, _ := NewAesGcmProvider("2B346A456B444C5361575879776D4867344439495158702F6135546175755275")

	a := map[string]string{
		"UserId": "1",
		"Data":   "2",
	}

	aBytes, _ := json.Marshal(a)
	bytes, _ := aesProvider.Seal(aBytes)
	plain, _ := aesProvider.Open(bytes)
	assert.Equal(t, aBytes, plain)
}
