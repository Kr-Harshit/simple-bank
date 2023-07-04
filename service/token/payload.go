package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Payload contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {

	if username == "" {
		return nil, errors.New("invalid username")
	}
	if duration <= 0 {
		return nil, errors.New("invalid duration, duration must be greater than 0")
	}

	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  currentTime.UTC(),
		ExpiredAt: currentTime.Add(duration).UTC(),
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
