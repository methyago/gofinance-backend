package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/methyago/gofinance-backend/db/sqlc"
)

type createAccountRequest struct {
	UserID      int32     `json:"user_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Type        string    `json:"type" binding:"required"`
	Description string    `json:"description" binding:"required"`
	CategoryID  int32     `json:"category_id" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
	Value       int32     `json:"value" binding:"required"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var categoryId = req.CategoryID
	var accountType = req.Type

	cat, err := server.store.GetCategory(ctx, categoryId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	if cat.Type != accountType {
		ctx.JSON(http.StatusBadRequest, gin.H{"error:": "Account type is different of category type"})
		return
	}

	arg := db.CreateAccountParams{
		Title:       req.Title,
		Type:        req.Type,
		Description: req.Description,
		UserID:      req.UserID,
		CategoryID:  req.CategoryID,
		Date:        req.Date,
		Value:       req.Value,
	}

	acc, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, acc)
}

type getAccountRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	acc, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, acc)
}

type deleteAccountRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, true)
}

type updateAccountIdRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

type updateAccountRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Value       int32  `json:"value"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var reqUri updateAccountIdRequest
	err := ctx.ShouldBindUri(&reqUri)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqBody updateAccountRequest
	err = ctx.ShouldBindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateAccountsParams{
		ID:          reqUri.ID,
		Title:       reqBody.Title,
		Description: reqBody.Description,
		Value:       reqBody.Value,
	}

	acc, err := server.store.UpdateAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, acc)
}

type listAccountsRequest struct {
	UserID      int32     `json:"user_id" binding:"required"`
	Type        string    `json:"type" binding:"required"`
	CategoryID  int32     `json:"category_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

func (server *Server) getAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetAccountsParams{
		UserID:      req.UserID,
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
		CategoryID: sql.NullInt32{
			Int32: req.CategoryID,
			Valid: req.CategoryID > 0,
		},
		Date: sql.NullTime{
			Time:  req.Date,
			Valid: !req.Date.IsZero(),
		},
	}

	cats, err := server.store.GetAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cats)
}

type getAccountGraphRequest struct {
	UserID int32  `json:"user_id" binding:"required"`
	Type   string `json:"type" binding:"required"`
}

func (server *Server) getAccountGraph(ctx *gin.Context) {
	var req getAccountGraphRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetAccountGraphParams{
		UserID: req.UserID,
		Type:   req.Type,
	}

	value, err := server.store.GetAccountGraph(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, value)
}

type getAccountReportsRequest struct {
	UserID int32  `json:"user_id" binding:"required"`
	Type   string `json:"type" binding:"required"`
}

func (server *Server) getAccountsReports(ctx *gin.Context) {
	var req getAccountReportsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetAccountsReportsParams{
		UserID: req.UserID,
		Type:   req.Type,
	}

	value, err := server.store.GetAccountsReports(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, value)
}
