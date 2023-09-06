package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/methyago/gofinance-backend/db/sqlc"
)

type createCategoryRequest struct {
	Title       string `json:"title" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (server *Server) createCategory(ctx *gin.Context) {
	userClaims := server.GetTokenInHeaderAndVerify(ctx)
	if userClaims == nil {
		return
	}

	var req createCategoryRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateCategoryParams{
		Title:       req.Title,
		Type:        req.Type,
		Description: req.Description,
		UserID:      userClaims.UserID,
	}

	cat, err := server.store.CreateCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cat)
}

type getCategoryRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

func (server *Server) getCategory(ctx *gin.Context) {
	userClaims := server.GetTokenInHeaderAndVerify(ctx)
	if userClaims == nil {
		return
	}

	var req getCategoryRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cat, err := server.store.GetCategory(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cat)
}

type deleteCategoryRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

func (server *Server) deleteCategory(ctx *gin.Context) {
	userClaims := server.GetTokenInHeaderAndVerify(ctx)
	if userClaims == nil {
		return
	}

	var req deleteCategoryRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteCategory(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, true)
}

type updateCategoryIdRequest struct {
	ID int32 `uri:"id" binding:"required"`
}
type updateCategoryRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (server *Server) updateCategory(ctx *gin.Context) {
	userClaims := server.GetTokenInHeaderAndVerify(ctx)
	if userClaims == nil {
		return
	}

	var reqUri updateCategoryIdRequest

	err := ctx.ShouldBindUri(&reqUri)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var reqBody updateCategoryRequest
	err = ctx.ShouldBindJSON(&reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateCategoriesParams{
		ID:          reqUri.ID,
		Title:       reqBody.Title,
		Description: reqBody.Description,
	}

	cat, err := server.store.UpdateCategories(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cat)
}

type listCategoriesRequest struct {
	Type        string `form:"type" json:"type" binding:"required"`
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
}

func (server *Server) getCategories(ctx *gin.Context) {
	userClaims := server.GetTokenInHeaderAndVerify(ctx)
	if userClaims == nil {
		return
	}

	var req listCategoriesRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetCategoriesParams{
		UserID:      userClaims.UserID,
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
	}

	cats, err := server.store.GetCategories(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cats)
}
