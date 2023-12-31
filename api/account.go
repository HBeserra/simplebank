package api

import (
	"database/sql"
	db "github.com/HBeserra/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type createAccountReq struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=BRL USD EUR"`
}

func (s Server) createAccount(ctx *gin.Context) {
	var req createAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) // Returns the error to the client
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := s.Store.CreateAccount(ctx, arg)
	if err != nil {
		log.Print("error can't create account", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"account": account,
	})
}

type getAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s Server) getAccount(ctx *gin.Context) {
	req := getAccountReq{}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := s.Store.GetAccount(ctx, req.ID)
	if err != nil {
		log.Println("Erro no sql:", gin.H{
			"params": req,
			"error":  err,
		})

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"account": account,
	})

}

type listAccReq struct {
	PageNumber int32 `form:"page_number" binding:"min=1"`
	PageSize   int32 `form:"page_size" binding:"min=5,max=50"`
}

func (s Server) listAccounts(ctx *gin.Context) {

	req := &listAccReq{
		PageNumber: 1,
		PageSize:   20,
	}

	if err := ctx.ShouldBindQuery(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: req.PageSize * (req.PageNumber - 1),
	}

	accounts, err := s.Store.ListAccounts(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
