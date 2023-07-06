package api

import (
	"fmt"
	"reflect"

	db "github.com/KHarshit1203/simple-bank/service/db/gen"
	"github.com/KHarshit1203/simple-bank/service/token"
	"github.com/KHarshit1203/simple-bank/service/util"
	"github.com/gofiber/fiber/v2"
)

type Server interface {
	// setupRoutes adds router routes to server method and path
	setupRoutes()

	// Start starts the server on a given address
	Start(address string) error
}

type ApiServer struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *fiber.App
	validator  ApiValidator
}

func NewServer(config util.Config, store db.Store) (Server, error) {
	if reflect.ValueOf(config).IsZero() {
		return nil, fmt.Errorf("invalid config")
	}
	if store == nil {
		return nil, fmt.Errorf("invalid store")
	}

	tokenMaker, err := token.NewPaestoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	router := fiber.New(fiber.Config{
		AppName:      "Simple Bank",
		ServerHeader: "Simple Bank",
	})

	validator := NewValidator()
	validator.RegisterValidation("currency", validCurrency)

	apiserver := &ApiServer{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		router:     router,
		validator:  validator,
	}

	apiserver.setupRoutes()

	return apiserver, nil
}

func (as *ApiServer) setupRoutes() {
	as.router.Route("/api", func(api fiber.Router) {
		// user api's
		api.Post("/user", as.createUser)
		api.Post("/login", as.loginUser)

		authRoutes := api.Group("/").Use(authMiddleware(as.tokenMaker))

		// account api's
		authRoutes.Route("/accounts", func(router fiber.Router) {
			router.Post("/", as.createAccount)
			router.Get("/:id", as.getAccount)
			router.Get("/", as.listAccount)
			router.Delete("/:id", as.deleteAccount)
		})
		authRoutes.Delete("/accounts/purge", as.purgeUserAccounts)
		authRoutes.Post("/transfer", as.createTransfer)
	})
}

func (as *ApiServer) Start(address string) error {
	return as.router.Listen(address)
}
