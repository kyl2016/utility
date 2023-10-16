package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUUID(t *testing.T) {
	timeStampLenght := 8
	for i := 1; i <= 100; i++ {
		randomString := GenerateUUID(uint8(i))
		assert.Equal(t, len(randomString), timeStampLenght+i)
	}
}

func TestRandomNumber(t *testing.T) {
	timeStampLength := 10
	for i := 1; i <= 50; i++ {
		randomString := GenerateFixedLengthRandomNumber(uint8(timeStampLength))
		assert.Equal(t, len(randomString), timeStampLength)
	}
}
