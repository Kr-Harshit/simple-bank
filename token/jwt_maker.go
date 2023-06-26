package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

// JWTMaker is a JSON web Token maker
type JWTMaker struct {
	secretKey []byte
}

// NewJwtMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, ErrInavlidJWTKey
	}
	return &JWTMaker{secretKey: []byte(secretKey)}, nil
}

// CreateToken creates a new JSON web token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	claims := jwt.RegisteredClaims{
		Issuer:    "simple-bank",
		IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
		ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
		Subject:   payload.Username,
		ID:        payload.ID.String(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(maker.secretKey)
	return token, payload, err
}

// VerifyToken checks if the JSON web token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(jwtToken *jwt.Token) (interface{}, error) {
		_, ok := jwtToken.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInavlidToken
		}
		return maker.secretKey, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInavlidToken
	}

	claim, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, ErrInavlidToken
	}

	id, err := uuid.Parse(claim.ID)
	if err != nil {
		return nil, ErrParsingToken
	}

	payload := &Payload{
		ID:        id,
		Username:  claim.Subject,
		IssuedAt:  claim.IssuedAt.UTC(),
		ExpiredAt: claim.ExpiresAt.UTC(),
	}

	return payload, nil
}
