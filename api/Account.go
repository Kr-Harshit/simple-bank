package api

import (
	"net/http"

	db "github.com/KHarshit1203/simple-bank/db/gen"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

type createAccountRequest struct {
	OwnerID  string `json:"owner-id" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errroResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		OwnerID:  req.OwnerID,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errroResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errroResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errroResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errroResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	OwnerID  string `form:"owner-id" binding:"required"`
	PageID   int32  `form:"page-id" binding:"required,min=1"`
	PageSize int32  `form:"page-size" binding:"required,min=5,max=20"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errroResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		OwnerID: req.OwnerID,
		Limit:   req.PageSize,
		Offset:  (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errroResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errroResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errroResponse(err))
		return
	}

	err := server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errroResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errroResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "account deleted")
}

type purgeUserAccountsRequest struct {
	OwnerID string `uri:"owner-id" binding:"required"`
}

func (server *Server) purgeUserAccounts(ctx *gin.Context) {
	var req purgeUserAccountsRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errroResponse(err))
		return
	}

	err := server.store.PurgeUserAccounts(ctx, req.OwnerID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errroResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errroResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "user accounts deleted")
}
