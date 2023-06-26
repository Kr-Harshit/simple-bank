package token

import (
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

// PaestoMaker is Paesto token maker
type PaestoMaker struct {
	paesto       *paseto.V2
	symmetricKey []byte
}

func NewPaestoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, ErrInavlidPaestoKey
	}

	maker := &PaestoMaker{
		paesto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

// CreateToken creates a new Paesto token for a specific username and duration
func (maker *PaestoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paesto.Encrypt(maker.symmetricKey, payload, nil)
	return token, payload, err
}

// VerifyToken checks if the Paesto token is valid or not
func (maker *PaestoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	if err := maker.paesto.Decrypt(token, maker.symmetricKey, payload, nil); err != nil {
		return nil, ErrInavlidToken
	}

	if err := payload.Valid(); err != nil {
		return nil, ErrExpiredToken
	}
	return payload, nil
}
