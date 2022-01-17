package middleware

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
	"strings"

	"github.com/bhojpur/wallet/pkg/auth"
	"github.com/bhojpur/wallet/pkg/errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func AuthByBearerToken(secret string) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		// check that the header is actually set
		header := ctx.Get("Authorization")
		if header == "" {
			return errors.Unauthorized{Message: "authorization header not set"}
		}

		// check that the token value in header is set
		bearer := strings.Split(header, " ")
		if len(bearer) < 2 || bearer[1] == "" {
			return errors.Unauthorized{Message: "authentication token not set"}
		}

		var claims auth.TokenClaims
		token, err := auth.ParseToken(bearer[1], secret, &claims)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return errors.Unauthorized{Message: "invalid signature on token"}
			}

			return errors.Unauthorized{Message: "token has expired or is invalid"}
		}
		if valid := auth.ValidateToken(token); !valid {
			return errors.Unauthorized{Message: "invalid token"}
		}

		ctx.Locals("userDetails", claims.User)

		return ctx.Next()
	}
}
