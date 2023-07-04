package api

import (
	"errors"
	"time"

	db "github.com/KHarshit1203/simple-bank/service/db/gen"
	dbUtil "github.com/KHarshit1203/simple-bank/service/db/utils"
	"github.com/KHarshit1203/simple-bank/service/util"
	"github.com/gofiber/fiber/v2"
)

type createUserRequest struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=10"`
	Fullname string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	Fullname          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// newUserResponse returns the user response
func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		Email:             user.Email,
		Fullname:          user.FullName,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

// createUser implements routing logic for create user api
func (as *ApiServer) createUser(ctx *fiber.Ctx) error {
	ctx.Accepts("application/json")

	var req createUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	if errors := as.validator.validateRequest(req); errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.Fullname,
		Email:          req.Email,
	}

	user, err := as.store.CreateUser(ctx.Context(), arg)
	if err != nil {
		if dbUtil.CheckErrorCode(err, dbUtil.ErrorUniqueKeyViolation.Code) {
			return fiber.NewError(fiber.ErrForbidden.Code, err.Error())
		}
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	resp := newUserResponse(user)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

type loginUserRequest struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=10"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (as *ApiServer) loginUser(ctx *fiber.Ctx) error {
	var req loginUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	if errors := as.validator.validateRequest(req); errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	user, err := as.store.GetUser(ctx.Context(), req.Username)
	if err != nil {
		if errors.Is(err, dbUtil.ErrorRecordNotFound) {
			return fiber.NewError(fiber.ErrNotFound.Code, err.Error())
		}
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	if err := util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		return fiber.NewError(fiber.ErrUnauthorized.Code, err.Error())
	}

	accessToken, _, err := as.tokenMaker.CreateToken(user.Username, as.config.AccessTokenDuration)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(loginUserResponse{AccessToken: accessToken, User: newUserResponse(user)})
}
