package token

import (
	"errors"
	"fmt"

	"github.com/aead/chacha20poly1305"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInavlidToken = errors.New("token is invalid")
	ErrParsingToken = errors.New("token payload is invalid")

	// JWT ERRORS
	ErrInavlidJWTKey = fmt.Errorf("invalid key size: must be atleast %d charracter", minSecretKeySize)

	// Paesto ERRORS
	ErrInavlidPaestoKey = fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
)
