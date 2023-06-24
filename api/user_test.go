package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	db "github.com/KHarshit1203/simple-bank/db/gen"
	"github.com/KHarshit1203/simple-bank/db/mock"
	"github.com/KHarshit1203/simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg: arg, password: password}
}

func TestCreateUser(t *testing.T) {
	password := util.RandomString(10)

	user := db.User{
		Username: util.RandomString(6),
		FullName: util.RandomString(10),
		Email:    util.RandomEmail(),
	}

	testCases := []struct {
		name          string
		requestBody   gin.H
		buildStub     func(store *mock.MockStore)
		checkResponse func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			requestBody: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStub: func(store *mock.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					Email:    user.Email,
					FullName: user.FullName,
				}

				store.
					EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)

			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, responseRecorder.Code)
				requireBodyMatcherUser(t, responseRecorder, user)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// get mock store
			store := mock.NewMockStore(ctrl)

			// build stub
			tc.buildStub(store)

			// start test server and make request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Marshal request body to JSOn
			requestByte, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			url := "/api/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestByte))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatcherUser(t *testing.T, recorder *httptest.ResponseRecorder, user db.User) {
	require.NotEmpty(t, recorder)

	data, err := io.ReadAll(recorder.Body)
	require.NoError(t, err)

	var gotUser createUserResponse
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.Fullname)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.PasswordChangedAt, gotUser.PasswordChangedAt)
	require.Equal(t, user.CreatedAt, gotUser.CreatedAt)
}
