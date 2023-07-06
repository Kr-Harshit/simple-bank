package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/KHarshit1203/simple-bank/service/token"
	"github.com/gofiber/fiber/v2"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorizationHeader := c.Get(authorizationHeaderKey, "")
		if len(authorizationHeader) == 0 {
			return fiber.NewError(http.StatusUnauthorized, "missing authorization header")
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return fiber.NewError(http.StatusUnauthorized, "invalid authorization header format")
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			return fiber.NewError(http.StatusUnauthorized, fmt.Sprintf("unsupported authorization type %s", authorizationType))
		}

		acessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(acessToken)
		if err != nil {
			return fiber.NewError(http.StatusUnauthorized, err.Error())
		}

		c.Locals(authorizationPayloadKey, payload)
		return c.Next()
	}

}
