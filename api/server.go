package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/methyago/gofinance-backend/db/sqlc"
)

type Server struct {
	store  *db.SQLStore
	router *gin.Engine
}

func NewServer(store *db.SQLStore) Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/user", server.createUser)
	router.GET("/get/:username", server.getUser)
	router.GET("/get/:id", server.getUserById)

	server.router = router
	return *server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"api has error:": err.Error()}

}
