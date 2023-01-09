package api

import (
	"net/http"

	db "github.com/ProstoyVadila/simple_bank/db/sqlc"
	"github.com/ProstoyVadila/simple_bank/e"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type transferRequest struct {
	FromAccountID uuid.UUID `json:"from_account_id" binding:"required,uuid"`
	ToAccountID   uuid.UUID `json:"to_account_id" binding:"required,uuid"`
	Amount        int64     `json:"amount" binding:"required,gt=0"`
	Currency      string    `json:"currency" binding:"required,oneof=PHP KZT USD"`
}

func (s *Server) CreateTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
		Currency:      req.Currency,
	}
	result, err := s.store.TransferTx(ctx, args)
	if err != nil {
		switch err.(type) {
		case e.ErrInvalidCurrencyType:
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
		case e.ErrAccountNotFound:
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, result)
}