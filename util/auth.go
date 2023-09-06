package util

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func ValidateToken(ctx *gin.Context, token string) error {
	claims := &Claims{}
	var jwtSignedKey = []byte("secret_key")
	tokenParse, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtSignedKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return err
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return err
	}

	if !tokenParse.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
		return &gin.Error{}
	}

	ctx.Next()
	return nil

}

func GetTokenInHeaderAndVerify(ctx *gin.Context) error {
	authorizationHeaderKey := ctx.GetHeader("Authorization")
	fields := strings.Fields(authorizationHeaderKey)

	if len(fields) < 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return gin.Error{}
	}

	tokenToValidade := fields[1]
	err := ValidateToken(ctx, tokenToValidade)
	if err != nil {
		return err
	}

	return nil
}
