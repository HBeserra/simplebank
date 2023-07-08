package api

import (
	db "github.com/HBeserra/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type createAccountReq struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required" oneof:"BRL USD EUR"`
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
