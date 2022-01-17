package auth

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"log"
	"time"

	"github.com/bhojpur/wallet/pkg/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
)

// TokenParsingError ...
type TokenParsingError struct {
	message string
}

func (err TokenParsingError) Error() string {
	return err.message
}

type UserAuthDetails struct {
	UserID   uuid.UUID       `json:"userId"`
	UserType models.UserType `json:"userType"`
}

type TokenClaims struct {
	User UserAuthDetails `json:"user"`

	jwt.StandardClaims
}

func generateToken(userID uuid.UUID, userType models.UserType) *jwt.Token {

	issuedAt := time.Now().Unix()
	expirationTime := time.Now().Add(6 * time.Hour).Unix()

	claims := TokenClaims{
		User: UserAuthDetails{
			UserID:   userID,
			UserType: userType,
		},

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			IssuedAt:  issuedAt,
		},
	}

	// We build a token, we give it and expiry of 6 hours.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token
}

// GetTokenString generates a jwt access token for a user
func GetTokenString(userID uuid.UUID, userType models.UserType, secret string) (string, error) {
	token := generateToken(userID, userType)

	str, err := token.SignedString([]byte(secret))
	if err != nil { // we have an error generating the token i.e. "500"
		log.Println(err)
		return "", TokenParsingError{message: err.Error()}
	}
	return str, nil
}

func ParseToken(token, secret string, claims *TokenClaims) (*jwt.Token, error) {
	tok, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	return tok, err
}

func ValidateToken(tok *jwt.Token) bool {
	if !tok.Valid {
		return false
	}
	return true
}
