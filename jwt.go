package utility

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrJWTAlgInvalid = errors.New("jwt_alg_invalid")
)

func MakeJWTToken(alg string, key []byte, payload StrMap) (string, error) {
	claims := jwt.MapClaims{}
	for k, v := range payload {
		claims[k] = v
	}
	claims["iat"] = time.Now().Unix() * 1000
	method := jwt.GetSigningMethod(normalizedJWTAlg(alg))
	token := jwt.NewWithClaims(method, claims)
	return token.SignedString(key)
}

func ParseJWTPayload(tokenStr string, alg string, key []byte) (payload interface{}, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		tokenAlg := token.Method.Alg()
		if tokenAlg != normalizedJWTAlg(alg) {
			return nil, fmt.Errorf("signing method invalid %w: %v", ErrJWTAlgInvalid, tokenAlg)
		}
		return key, nil
	})
	if err == nil {
		switch v := token.Claims.(type) {
		case jwt.MapClaims:
			payload = StrMap(v)
		default:
			payload = v
		}
	}
	return
}

func normalizedJWTAlg(alg string) string {
	if alg == "" {
		return "HS256"
	}
	return strings.ToUpper(alg)
}
