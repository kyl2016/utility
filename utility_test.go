package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlattenMap(t *testing.T) {
	originData := make(StrMap)
	originData["a"] = 1
	secondLayerData := make(StrMap)
	secondLayerData["b"] = "foo"
	originData["b"] = secondLayerData
	assert.Equal(t, 1, originData["a"])
	result := FlattenMap("", ".", originData)
	assert.Equal(t, result["b.b"], "foo")
}
func TestCheckFloat64IsFloat32(t *testing.T) {
	var f float64
	f = 2.42
	assert.True(t, CanConvertToFloat32Loselessly(f))
	f = 1.401298464324817070923729583289916131280e-41
	assert.True(t, CanConvertToFloat32Loselessly(f))
	f = 1.401298464324817070923729583289916131280e-48
	assert.False(t, CanConvertToFloat32Loselessly(f))
}
func TestCheckFloat64IsInt64(t *testing.T) {
	var f float64
	f = 200
	assert.True(t, CanConvertToInt64Loselessly(f))
}
func TestCheckFloat64IsInt32(t *testing.T) {
	var f float64
	f = 200
	assert.True(t, CanConvertToInt32Loselessly(f))
	f = 20000000000
	assert.False(t, CanConvertToInt32Loselessly(f))
}
func TestStringToChunks(t *testing.T) {
	s := "abcdefg"
	input := []struct {
		chunkSize int
		chunks    []string
	}{
		{1, []string{"a", "b", "c", "d", "e", "f", "g"}},
		{2, []string{"ab", "cd", "ef", "g"}},
		{3, []string{"abc", "def", "g"}},
		{4, []string{"abcd", "efg"}},
		{5, []string{"abcde", "fg"}},
		{6, []string{"abcdef", "g"}},
		{7, []string{"abcdefg"}},
		{8, []string{"abcdefg"}},
		{100, []string{"abcdefg"}},
	}
	for _, item := range input {
		chunks := StringToChunks(s, item.chunkSize)
		assert.Equal(t, item.chunks, chunks)
	}
}
