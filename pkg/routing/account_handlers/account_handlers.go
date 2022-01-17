package account_handlers

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
	"net/http"

	"github.com/bhojpur/wallet/pkg/account"
	"github.com/bhojpur/wallet/pkg/auth"
	"github.com/bhojpur/wallet/pkg/errors"
	"github.com/bhojpur/wallet/pkg/models"
	"github.com/bhojpur/wallet/pkg/routing/responses"
	"github.com/bhojpur/wallet/pkg/statement"

	"github.com/gofiber/fiber/v2"
)

// BalanceEnquiry ...
func BalanceEnquiry(interactor account.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails auth.UserAuthDetails
		if details, ok := ctx.Locals("userDetails").(auth.UserAuthDetails); !ok {
			return errors.Error{Code: errors.EINTERNAL}
		} else {
			userDetails = details
		}

		// we check if user is admin, we return error
		if userDetails.UserType == models.UserTypAdmin {
			return errors.Error{Code: errors.EINVALID, Message: errors.UserCantHaveAccount}
		}

		balance, err := interactor.GetBalance(userDetails.UserID)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.BalanceResponse(userDetails.UserID, balance))
	}
}

// MiniStatement returns a small short summary of the
// most recent transactions on an account.
func MiniStatement(interactor statement.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails auth.UserAuthDetails
		if details, ok := ctx.Locals("userDetails").(auth.UserAuthDetails); !ok {
			return errors.Error{Code: errors.EINTERNAL}
		} else {
			userDetails = details
		}

		statements, err := interactor.GetStatement(userDetails.UserID)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.MiniStatementResponse(userDetails.UserID, statements))
	}
}
