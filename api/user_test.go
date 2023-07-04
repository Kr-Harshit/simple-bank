package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/KHarshit1203/simple-bank/service/db/gen"
	"github.com/KHarshit1203/simple-bank/service/db/mocks"
	"github.com/KHarshit1203/simple-bank/service/util"
	"github.com/gofiber/fiber/v2"
)

func (at *ApiServerSuite) TestCreateUser() {
	url := "/api/user"

	validUsername := util.RandomString(10)
	validPassword := util.RandomString(10)
	validHashPassword, err := util.HashPassword(validPassword)
	at.NoError(err)
	validEmail := util.RandomEmail()
	validFullName := util.RandomString(20)

	tests := []struct {
		name        string
		requestBody fiber.Map
		setStore    func(t *testing.T, s *mocks.Store)
		check       func(t *testing.T, resp *http.Response, err error)
	}{
		{
			name: "Status OK",
			requestBody: fiber.Map{
				"username":  validUsername,
				"password":  validPassword,
				"email":     validEmail,
				"full_name": validFullName,
			},
			setStore: func(t *testing.T, s *mocks.Store) {
				args := db.CreateUserParams{
					Username:       validUsername,
					HashedPassword: validHashPassword,
					FullName:       validFullName,
					Email:          validEmail,
				}

				s.On("CreateUser", context.Background(), args).Return(db.User{
					Username:       validUsername,
					HashedPassword: validHashPassword,
					FullName:       validFullName,
					Email:          validEmail,
				}, nil)

				// s.AssertCalled(t, "CreateUser")
				// s.AssertExpectations(t)
				// s.AssertNumberOfCalls(t, "CreateUser", 1)
			},
			check: func(t *testing.T, resp *http.Response, err error) {
				// req := require.New(t)

				// req.NoError(err)
				// req.NotEmpty(resp)

				// req.Equal(http.StatusOK, resp.StatusCode)
			},
		},
	}

	for _, tt := range tests {
		at.Run(tt.name, func() {
			tt.setStore(at.T(), at.mockStore)

			requestByte, err := json.Marshal(tt.requestBody)
			at.NoError(err)
			request := httptest.NewRequest("POST", url, bytes.NewReader(requestByte))

			apiServer, ok := at.server.(*ApiServer)
			at.True(ok)

			response, err := apiServer.router.Test(request)

			tt.check(at.T(), response, err)

		})
	}

}
