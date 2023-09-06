package api

import (
	"bytes"
	"crypto/sha512"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type loginResponseStruct struct {
	UserID   int32  `json:"user_id"`
	UserName string `json:"username"`
	Token    string `json:"token"`
}

func (server *Server) login(ctx *gin.Context) {
	var req loginRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedInput := sha512.Sum512_256([]byte(req.Password))
	trimmedHash := bytes.Trim(hashedInput[:], "\x00")
	preparedPassword := string(trimmedHash)
	plainTextInBytes := []byte(preparedPassword)
	hashTextInBytes := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(hashTextInBytes, plainTextInBytes)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &MyCustomClaims{
		Username: req.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var jwtSignedKey = []byte("secret_key")
	generatedTokenString, err := generatedToken.SignedString(jwtSignedKey)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	result := loginResponseStruct{
		UserID:   user.ID,
		UserName: user.Username,
		Token:    generatedTokenString,
	}
	ctx.JSON(http.StatusOK, result)
}

func validateToken(ctx *gin.Context, token string) (error, string) {
	var jwtSignedKey = []byte("secret_key")
	tokenParse, err := jwt.ParseWithClaims(token, &MyCustomClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return jwtSignedKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return err, ""
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return err, ""
	}

	claims := tokenParse.Claims.(*MyCustomClaims)

	if !tokenParse.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
		return &gin.Error{}, ""
	}

	ctx.Next()
	return nil, claims.Username

}

type UserClaims struct {
	UserID   int32
	UserName string
}

func (server *Server) GetTokenInHeaderAndVerify(ctx *gin.Context) *UserClaims {
	authorizationHeaderKey := ctx.GetHeader("Authorization")
	fields := strings.Fields(authorizationHeaderKey)

	if len(fields) < 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return nil
	}

	tokenToValidade := fields[1]
	err, username := validateToken(ctx, tokenToValidade)
	if err != nil || username == "" {
		return nil
	}

	user, err := server.store.GetUser(ctx, username)
	if err != nil {
		return nil
	}

	return &UserClaims{
		UserID:   user.ID,
		UserName: user.Username,
	}

}
