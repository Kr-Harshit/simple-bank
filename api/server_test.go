package api

import (
	"fmt"
	"testing"

	db "github.com/KHarshit1203/simple-bank/service/db/gen"
	"github.com/KHarshit1203/simple-bank/service/db/mocks"
	"github.com/KHarshit1203/simple-bank/service/token"
	"github.com/KHarshit1203/simple-bank/service/util"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	type args struct {
		config util.Config
		store  db.Store
	}

	validConfig := util.Config{TokenSymmetricKey: util.RandomString(32)}
	validStore := mocks.NewStore(t)

	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, gotServer Server, gotError error)
	}{
		{
			name: "valid arguments",
			args: args{
				config: validConfig,
				store:  validStore,
			},
			check: func(t *testing.T, gotServer Server, gotError error) {
				req := require.New(t)
				req.NoError(gotError)

				req.NotEmpty(gotServer)

				gotApiServer, ok := gotServer.(*ApiServer)
				req.True(ok)
				req.NotEmpty(gotApiServer.store)
				req.NotEmpty(gotApiServer.config)
				req.Equal(validConfig, gotApiServer.config)
				req.NotEmpty(gotApiServer.router)
				req.Equal("Simple Bank", gotApiServer.router.Config().ServerHeader)
				req.Equal("Simple Bank", gotApiServer.router.Config().AppName)
				req.NotEmpty(gotApiServer.validator)
			},
		},
		{
			name: "empty config input",
			args: args{
				store: validStore,
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
				config: validConfig,
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
				store:  validStore,
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
		t.Run(tt.name, func(t *testing.T) {
			gotServer, gotErr := NewServer(tt.args.config, tt.args.store)
			tt.check(t, gotServer, gotErr)
		})
	}
}
