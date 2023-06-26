package token

import (
	"testing"
	"time"

	"github.com/KHarshit1203/simple-bank/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	type args struct {
		symmetricKey string
		username     string
		duration     time.Duration
	}
	type response struct {
		maker   Maker
		token   string
		payload *Payload
		err     error
	}

	tests := []struct {
		name         string
		args         args
		requireCheck func(t *testing.T, arg args, resp response)
	}{
		{
			name: "invalid_symmetric_key",
			args: args{
				symmetricKey: "",
			},
			requireCheck: func(t *testing.T, arg args, resp response) {
				require.EqualError(t, resp.err, ErrInavlidJWTKey.Error())
				require.Empty(t, resp.maker)
			},
		},
		{
			name: "valid_input",
			args: args{
				symmetricKey: util.RandomString(32),
				username:     util.RandomString(32),
				duration:     time.Minute,
			},
			requireCheck: func(t *testing.T, arg args, resp response) {
				// check Maker
				require.NotEmpty(t, resp.maker)

				// check create Token
				require.NoError(t, resp.err)
				require.NotEmpty(t, resp.payload)
				require.NotEmpty(t, resp.token)

				// check verify token
				payload, err := resp.maker.VerifyToken(resp.token)
				require.NoError(t, err)
				require.Equal(t, resp.payload.Username, payload.Username)
				require.Equal(t, resp.payload.ID, payload.ID)
				require.Equal(t, resp.payload.ExpiredAt.Truncate(jwt.TimePrecision), payload.ExpiredAt)
				require.Equal(t, resp.payload.IssuedAt.Truncate(jwt.TimePrecision), payload.IssuedAt)
			},
		},
		{
			name: "invalid_create_token_input",
			args: args{
				symmetricKey: util.RandomString(32),
				username:     "",
				duration:     time.Minute,
			},
			requireCheck: func(t *testing.T, arg args, resp response) {
				// check Maker
				require.NotEmpty(t, resp.maker)

				require.Error(t, resp.err)
				require.Empty(t, resp.payload)
				require.Empty(t, resp.token)
			},
		},
		{
			name: "incorrect_token",
			args: args{
				symmetricKey: util.RandomString(32),
				username:     util.RandomString(10),
				duration:     time.Minute,
			},
			requireCheck: func(t *testing.T, arg args, resp response) {
				// check Maker
				require.NotEmpty(t, resp.maker)

				// check create Token
				require.NoError(t, resp.err)
				require.NotEmpty(t, resp.payload)
				require.NotEmpty(t, resp.token)

				incorrectToken := util.RandomString(len(resp.token))

				// check verify token
				payload, err := resp.maker.VerifyToken(incorrectToken)
				require.EqualError(t, err, ErrInavlidToken.Error())
				require.Empty(t, payload)
			},
		},
		{
			name: "expired_token",
			args: args{
				symmetricKey: util.RandomString(32),
				username:     util.RandomString(32),
				duration:     time.Microsecond,
			},
			requireCheck: func(t *testing.T, arg args, resp response) {
				// check verify token
				time.Sleep(time.Millisecond * 10)

				payload, err := resp.maker.VerifyToken(resp.token)
				require.EqualError(t, err, ErrExpiredToken.Error())
				require.Empty(t, payload)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maker, err := NewJWTMaker(tt.args.symmetricKey)
			if err != nil {
				tt.requireCheck(t, tt.args, response{maker: maker, err: err})
			} else {
				token, payload, err := maker.CreateToken(tt.args.username, tt.args.duration)
				tt.requireCheck(t, tt.args, response{maker, token, payload, err})
			}
		})
	}
}
