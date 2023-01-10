package api

import (
	"net/http"

	db "github.com/ProstoyVadila/simple_bank/db/sqlc"
	"github.com/ProstoyVadila/simple_bank/e"
	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID utils.UUIDString `json:"from_account_id" binding:"required,uuid"`
	ToAccountID   utils.UUIDString `json:"to_account_id" binding:"required,uuid"`
	Currency      string           `json:"currency" binding:"required,currency"`
	Amount        int64            `json:"amount" binding:"required,gt=0"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondWithValidationError(ctx, err)
		return
	}

	fromAccountID, _ := req.FromAccountID.UUID()
	toAccountID, _ := req.ToAccountID.UUID()

	args := db.TransferTxParams{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
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
