package api

import (
	"fmt"
	"net/http"

	db "github.com/KHarshit1203/simple-bank/service/db/gen"
	"github.com/KHarshit1203/simple-bank/service/token"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx"
)

type transferRequest struct {
	FromAccountID int64   `json:"from_account_id" validate:"required,min=1"`
	ToAccountID   int64   `json:"to_account_id" validate:"required,min=1"`
	Amount        float32 `json:"amount" validate:"required,gt=0"`
	Currency      string  `json:"currency" validate:"required,currency"`
}

func (as *ApiServer) createTransfer(ctx *fiber.Ctx) error {
	var req transferRequest

	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	if errors := as.validator.validateRequest(req); errors != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(errors)
	}

	if req.FromAccountID == req.ToAccountID {
		return fiber.NewError(http.StatusBadRequest, "from_account_id cannot be same as to_account_id")
	}

	FromAccount, err := as.validAccount(ctx, req.FromAccountID, req.Currency)
	if err != nil {
		return err
	}
	authPaylod := ctx.Locals(authorizationPayloadKey).(*token.Payload)
	if FromAccount.Owner != authPaylod.Username {
		return fiber.NewError(http.StatusUnauthorized, "user doesn't access to account")
	}

	if _, err := as.validAccount(ctx, req.ToAccountID, req.Currency); err != nil {
		return err
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := as.store.TransferTx(ctx.Context(), arg)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(result)
}

func (as *ApiServer) validAccount(ctx *fiber.Ctx, accountID int64, currency string) (db.Account, error) {
	account, err := as.store.GetAccount(ctx.Context(), accountID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return account, fiber.NewError(http.StatusNotFound, err.Error())
		}
		return account, fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	if account.Currency != currency {
		return account, fiber.NewError(http.StatusBadRequest, fmt.Sprintf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency))
	}
	return account, nil
}
