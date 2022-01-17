package transaction_handlers

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

	"github.com/bhojpur/wallet/pkg/auth"
	"github.com/bhojpur/wallet/pkg/errors"
	"github.com/bhojpur/wallet/pkg/models"
	"github.com/bhojpur/wallet/pkg/ports"
	"github.com/bhojpur/wallet/pkg/routing/responses"
	"github.com/bhojpur/wallet/pkg/transaction"

	"github.com/gofiber/fiber/v2"
)

// Deposit allows user to deposit or credit their account.
func Deposit(txnAdapter ports.TransactorPort) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails auth.UserAuthDetails
		if details, ok := ctx.Locals("userDetails").(auth.UserAuthDetails); !ok {
			return errors.Error{Code: errors.EINTERNAL}
		} else {
			userDetails = details
		}

		// inflate struct with body params
		var p transaction.DepositParams
		_ = ctx.BodyParser(&p)

		// validate params
		err := p.Validate()
		if err != nil {
			return err
		}

		depositor := models.TxnCustomer{
			UserType: userDetails.UserType,
			UserID:   userDetails.UserID,
		}
		err = txnAdapter.Deposit(depositor, p.CustomerNumber, p.CustomerType, p.Amount)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.TransactionResponse())
	}
}

// Withdraw allows user to withdraw or debit their account.
func Withdraw(txnAdapter ports.TransactorPort) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails auth.UserAuthDetails
		if details, ok := ctx.Locals("userDetails").(auth.UserAuthDetails); !ok {
			return errors.Error{Code: errors.EINTERNAL}
		} else {
			userDetails = details
		}

		// inflate struct with body params
		var p transaction.WithdrawParams
		_ = ctx.BodyParser(&p)

		// validate params
		err := p.Validate()
		if err != nil {
			return err
		}

		withdrawer := models.TxnCustomer{
			UserID:   userDetails.UserID,
			UserType: userDetails.UserType,
		}
		err = txnAdapter.Withdraw(withdrawer, p.AgentNumber, p.Amount)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.TransactionResponse())
	}
}

func Transfer(txnAdapter ports.TransactorPort) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails auth.UserAuthDetails
		if details, ok := ctx.Locals("userDetails").(auth.UserAuthDetails); !ok {
			return errors.Error{Code: errors.EINTERNAL}
		} else {
			userDetails = details
		}

		// inflate struct with body params
		var p transaction.TransferParams
		_ = ctx.BodyParser(&p)

		// validate params
		err := p.Validate()
		if err != nil {
			return err
		}

		source := models.TxnCustomer{
			UserID:   userDetails.UserID,
			UserType: userDetails.UserType,
		}
		err = txnAdapter.Transfer(source, p.DestAccountNo, p.DestUserType, p.Amount)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.TransactionResponse())
	}
}
