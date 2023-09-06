package api

import (
	"bytes"
	"crypto/sha512"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
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

	expirationTime := time.Now().Add(100 * time.Minute)
	claims := &Claims{
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
	ctx.JSON(http.StatusOK, generatedTokenString)
}
