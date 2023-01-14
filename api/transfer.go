package api

import (
	"net/http"

	"github.com/ProstoyVadila/simple_bank/api/middleware"
	db "github.com/ProstoyVadila/simple_bank/db/sqlc"
	"github.com/ProstoyVadila/simple_bank/e"
	"github.com/ProstoyVadila/simple_bank/token"
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
		ctx.Error(err)
		return
	}

	authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)
	fromAccount, toAccount, err := s.validateAccounts(ctx, authPayload, req.FromAccountID, req.ToAccountID, req.Currency)
	if err != nil {
		switch err.(type) {
		case e.ErrUnauthorized, e.ErrInvalidCurrencyType, e.ErrInvalidUUID:
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		default:
			ctx.Error(err)
		}
		return
	}

	args := db.TransferTxParams{
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		Amount:      req.Amount,
		Currency:    req.Currency,
	}
	result, err := s.store.TransferTx(ctx, args)
	if err != nil {
		switch err.(type) {
		case e.ErrUnauthorized:
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *Server) validateAccounts(
	ctx *gin.Context,
	authPayload *token.Payload,
	fromAccountID utils.UUIDString,
	toAccountID utils.UUIDString,
	currency string,
) (db.Account, db.Account, error) {
	fromAccID, err := fromAccountID.UUID()
	if err != nil {
		return db.Account{}, db.Account{}, err
	}
	toAccID, err := toAccountID.UUID()
	if err != nil {
		return db.Account{}, db.Account{}, err
	}
	fromAccount, err := s.store.GetAccount(ctx, fromAccID)
	if err != nil {
		return db.Account{}, db.Account{}, err
	}
	toAccount, err := s.store.GetAccount(ctx, toAccID)
	if err != nil {
		return db.Account{}, db.Account{}, err
	}
	if fromAccount.Currency != currency || toAccount.Currency != currency {
		return db.Account{}, db.Account{}, e.ErrInvalidCurrencyType{Curr: currency}
	}
	if fromAccount.OwnerName != authPayload.Username {
		return db.Account{}, db.Account{}, e.ErrUnauthorized{
			Msg: "from account does not belong to the user",
			Obj: fromAccount.OwnerName,
		}
	}

	return fromAccount, toAccount, nil
}
