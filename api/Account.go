package api

import (
	"errors"
	"net/http"

	db "github.com/KHarshit1203/simple-bank/service/db/gen"
	dbUtil "github.com/KHarshit1203/simple-bank/service/db/utils"
	"github.com/gofiber/fiber/v2"
)

type createAccountRequest struct {
	Owner    string `json:"owner" validate:"required"`
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

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
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
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	if errors := as.validator.validateRequest(req); errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	account, err := as.store.GetAccount(ctx.Context(), req.ID)
	if err != nil {
		if errors.Is(err, dbUtil.ErrorRecordNotFound) {
			return fiber.NewError(fiber.ErrNotFound.Code, err.Error())
		}
		fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}
	return ctx.Status(fiber.StatusOK).JSON(account)
}

type listAccountRequest struct {
	Owner    string `form:"owner" validate:"required"`
	PageID   int32  `form:"page-id" validate:"required,min=1"`
	PageSize int32  `form:"page-size" validate:"required,min=5,max=20"`
}

func (as *ApiServer) listAccount(ctx *fiber.Ctx) error {
	var req listAccountRequest

	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	if errors := as.validator.validateRequest(req); errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	arg := db.ListAccountsParams{
		Owner:  req.Owner,
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

	err := as.store.DeleteAccount(ctx.Context(), req.ID)
	if err != nil {
		if errors.Is(err, dbUtil.ErrorRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON("account deleted")
}

type purgeUserAccountsRequest struct {
	OwnerID string `uri:"owner-id" validate:"required"`
}

func (as *ApiServer) purgeUserAccounts(ctx *fiber.Ctx) error {
	var req purgeUserAccountsRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	if errors := as.validator.validateRequest(req); errors != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(errors)
	}

	err := as.store.PurgeUserAccounts(ctx.Context(), req.OwnerID)
	if err != nil {
		if errors.Is(err, dbUtil.ErrorRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, err.Error())
		}
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON("user accounts deleted")
}
