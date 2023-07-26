package token

import (
	"testing"
	"time"

	"github.com/KHarshit1203/simple-bank/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewPayload(t *testing.T) {
	type args struct {
		username string
		duration time.Duration
	}
	type response struct {
		payload *Payload
		err     error
	}

	tests := []struct {
		name         string
		args         args
		requireCheck func(t *testing.T, arg args, response response)
	}{
		{
			name: "valid_input",
			args: args{username: util.RandomString(10), duration: time.Minute},
			requireCheck: func(t *testing.T, arg args, response response) {
				require.NoError(t, response.err)
				require.NotNil(t, response.payload)
				require.NotEmpty(t, response.payload.ID)

				// valid uuid
				_, err := uuid.Parse(response.payload.ID.String())
				require.NoError(t, err)

				require.Equal(t, arg.username, response.payload.Username)
				require.NotEmpty(t, response.payload.ExpiredAt)
				require.NotEmpty(t, response.payload.IssuedAt)
				require.True(t, response.payload.IssuedAt.Add(arg.duration).Equal(response.payload.ExpiredAt))
			},
		},
		{
			name: "invalid_username",
			args: args{username: "", duration: time.Minute},
			requireCheck: func(t *testing.T, arg args, response response) {
				require.EqualError(t, response.err, "invalid username")
			},
		},
		{
			name: "invalid_duration",
			args: args{username: util.RandomString(10), duration: -1},
			requireCheck: func(t *testing.T, arg args, response response) {
				require.EqualError(t, response.err, "invalid duration, duration must be greater than 0")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := NewPayload(tt.args.username, tt.args.duration)
			tt.requireCheck(t, tt.args, response{payload, err})
		})
	}
}
