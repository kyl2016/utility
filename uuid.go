package utility

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	uuidTimestampLength = 40
	b32alphabet         = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	b32NumberSet        = "1234567890123456"
	b64characterSet     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	b32WordLength       = 5
	b32Chopper          = 31
)

var serverStartTime = time.Date(2019, time.July, 12, 0, 0, 0, 0, time.UTC)

// GenerateUUID generates base32 uuid
//
// - Returns: timestamp (8 bytes) + randomString
//   - timestamp: consists of 8 alphanumeric characters in set of `b32alphabet`, where 8 = int(uuidTimestampLength/b32WordLength)
//   - randomString: consists of N alphanumeric characters in set of `b32alphabet`, where N = randomLength
func GenerateUUID(randomStringLength uint8) string {
	now := time.Now()
	milliseconds := now.Sub(serverStartTime).Milliseconds()
	encodedTimestamp := uuidEncode(milliseconds, uuidTimestampLength)
	encodedRandomString := GenerateFixedLengthRandomString(randomStringLength)
	return encodedTimestamp + encodedRandomString
}

func uuidEncode(number int64, length int) string {
	encodeLength := int(length / b32WordLength)
	result := make([]byte, encodeLength)
	for i, j := 0, encodeLength-1; i < encodeLength; i, j = i+1, j-1 {
		x := number & b32Chopper
		result[j] = b32alphabet[x]
		number >>= b32WordLength
	}
	return BytesToString(result)
}

func GenerateFixedLengthRandomString(length uint8) string {
	charSetLength := int64(len(b32alphabet))
	b := make([]byte, length)
	for i := range b {
		randInt64 := rand.Int63()
		index := randInt64 % charSetLength
		b[i] = b32alphabet[index]
	}
	return BytesToString(b)
}

func GenerateFixedLengthRandomNumber(length uint8) string {
	rand.Seed(time.Now().UnixNano())
	charSetLength := int32(len(b32NumberSet))
	b := make([]byte, length)
	for i := range b {
		randInt32 := rand.Int31n(10)
		index := randInt32 % charSetLength
		b[i] = b32NumberSet[index]
	}
	res := BytesToString(b)
	fmt.Println(res)
	return res
}
