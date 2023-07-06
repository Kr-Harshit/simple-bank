package api

import (
	"errors"
	"net/http"

	db "github.com/KHarshit1203/simple-bank/service/db/gen"
	dbUtil "github.com/KHarshit1203/simple-bank/service/db/utils"
	"github.com/KHarshit1203/simple-bank/service/token"
	"github.com/gofiber/fiber/v2"
)

type createAccountRequest struct {
	Currency string `json:"currency" validate:"required,currency"`
}

func (as *ApiServer) createAccount(ctx *fiber.Ctx) error {
	var req createAccountRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, err.Error())
	}

	if errors := as.validator.validateRequest(req); errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	authPayload := ctx.Locals(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := as.store.CreateAccount(ctx.Context(), arg)
	if err != nil {
		if dbUtil.CheckErrorCode(err, dbUtil.ErrorForeignKeyViolation.Code) || dbUtil.CheckErrorCode(err, dbUtil.ErrorUniqueKeyViolation.Code) {
			return fiber.NewError(fiber.ErrForbidden.Code, err.Error())
		}
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())

	}

	return ctx.Status(fiber.StatusOK).JSON(account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" validate:"required,min=1"`
}

func (as *ApiServer) getAccount(ctx *fiber.Ctx) error {
	var req getAccountRequest
	if err := ctx.ParamsParser(&req); err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	if errors := as.validator.validateRequest(req); errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	authPayload := ctx.Locals(authorizationPayloadKey).(*token.Payload)
	account, err := as.store.GetAccount(ctx.Context(), req.ID)
	if err != nil {
		if errors.Is(err, dbUtil.ErrorRecordNotFound) {
			return fiber.NewError(fiber.ErrNotFound.Code, err.Error())
		}
		fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}
	if account.Owner != authPayload.Username {
		return fiber.NewError(http.StatusNotFound, "no account found")
	}
	return ctx.Status(fiber.StatusOK).JSON(account)
}

type listAccountRequest struct {
	PageID   int32 `query:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" validate:"required,min=5,max=20"`
}

func (as *ApiServer) listAccount(ctx *fiber.Ctx) error {
	var req listAccountRequest

	if err := ctx.QueryParser(&req); err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	if errors := as.validator.validateRequest(req); errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	authPayload := ctx.Locals(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := as.store.ListAccounts(ctx.Context(), arg)
	if err != nil {
		if errors.Is(err, dbUtil.ErrorRecordNotFound) {
			return fiber.NewError(fiber.ErrNotFound.Code, err.Error())
		}
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(accounts)
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" validate:"required,min=1"`
}

func (as *ApiServer) deleteAccount(ctx *fiber.Ctx) error {
	var req deleteAccountRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	if errors := as.validator.validateRequest(req); errors != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(errors)
	}

	authPayload := ctx.Locals(authorizationPayloadKey).(*token.Payload)
	account, err := as.store.GetAccount(ctx.Context(), req.ID)
	if err != nil {
		if errors.Is(err, dbUtil.ErrorRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	if account.Owner != authPayload.Username {
		return fiber.NewError(http.StatusUnauthorized, "user doesn't have access to account")
	}

	err = as.store.DeleteAccount(ctx.Context(), db.DeleteAccountParams{
		Owner: authPayload.Username,
		ID:    account.ID,
	})
	if err != nil {
		if errors.Is(err, dbUtil.ErrorRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON("account deleted")
}

func (as *ApiServer) purgeUserAccounts(ctx *fiber.Ctx) error {
	authPayload := ctx.Locals(authorizationPayloadKey).(*token.Payload)

	err := as.store.PurgeUserAccounts(ctx.Context(), authPayload.Username)
	if err != nil {
		if errors.Is(err, dbUtil.ErrorRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON("user accounts deleted")
}
