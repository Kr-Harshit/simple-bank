package api

import (
	"fmt"
	"testing"

	db "github.com/KHarshit1203/simple-bank/service/db/gen"
	"github.com/KHarshit1203/simple-bank/service/db/mocks"
	"github.com/KHarshit1203/simple-bank/service/token"
	"github.com/KHarshit1203/simple-bank/service/util"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ApiServerSuite struct {
	suite.Suite
	mockConfig util.Config
	mockStore  *mocks.Store
	server     Server
}

func TestApiServerSuite(t *testing.T) {
	suite.Run(t, &ApiServerSuite{})
}

func (at *ApiServerSuite) SetupSubTest() {
	mockConfig := util.Config{TokenSymmetricKey: util.RandomString(32)}
	mockStore := mocks.NewStore(at.T())
	testServer, err := NewServer(mockConfig, mockStore)
	at.NoError(err)

	at.mockConfig = mockConfig
	at.mockStore = mockStore
	at.server = testServer
}

func (at *ApiServerSuite) TestNewServer() {
	type args struct {
		config util.Config
		store  db.Store
	}

	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, gotServer Server, gotError error)
	}{
		{
			name: "valid arguments",
			args: args{
				config: at.mockConfig,
				store:  at.mockStore,
			},
			check: func(t *testing.T, gotServer Server, gotError error) {
				req := require.New(t)
				req.NoError(gotError)

				req.NotEmpty(gotServer)

				gotApiServer, ok := gotServer.(*ApiServer)
				req.True(ok)
				req.NotEmpty(gotApiServer.store)
				req.NotEmpty(gotApiServer.config)
				req.Equal(at.mockConfig, gotApiServer.config)
				req.NotEmpty(gotApiServer.router)
				req.Equal("Simple Bank", gotApiServer.router.Config().ServerHeader)
				req.Equal("Simple Bank", gotApiServer.router.Config().AppName)
				req.NotEmpty(gotApiServer.validator)
			},
		},
		{
			name: "empty config input",
			args: args{
				store: at.mockStore,
			},
			check: func(t *testing.T, gotServer Server, gotError error) {
				req := require.New(t)
				req.Error(gotError)
				req.EqualError(gotError, "invalid config")
				req.Empty(gotServer)
			},
		},
		{
			name: "empty store input",
			args: args{
				config: at.mockConfig,
			},
			check: func(t *testing.T, gotServer Server, gotError error) {
				req := require.New(t)
				req.Error(gotError)
				req.EqualError(gotError, "invalid store")
				req.Empty(gotServer)
			},
		},
		{
			name: "invalid config input",
			args: args{
				config: util.Config{TokenSymmetricKey: util.RandomString(10)},
				store:  at.mockStore,
			},
			check: func(t *testing.T, gotServer Server, gotError error) {
				req := require.New(t)
				req.Error(gotError)
				req.EqualError(gotError, fmt.Sprintf("cannot create token maker: %v", token.ErrInavlidPaestoKey.Error()))
				req.Empty(gotServer)
			},
		},
	}

	for _, tt := range tests {
		at.Run(tt.name, func() {
			gotServer, gotErr := NewServer(tt.args.config, tt.args.store)
			tt.check(at.T(), gotServer, gotErr)
		})
	}
}
