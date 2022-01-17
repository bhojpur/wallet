package errors

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
	"fmt"
	"strconv"

	"github.com/bhojpur/wallet/pkg/models"

	"github.com/gofrs/uuid"
)

const (
	AccountNotCreated          = ERMessage("user's account has not been created, report issue")
	DepositAmountBelowMinimum  = ERMessage("cannot deposit amounts less than")
	WithdrawAmountBelowMinimum = ERMessage("cannot withdraw amounts less than")
	TransferAmountBelowMinimum = ERMessage("cannot transfer amounts less than")
	DebitAmountAboveBalance    = ERMessage("cannot debit amount, account balance not enough")

	UserCantHaveAccount = ERMessage("user is not allowed to hold an account")
)

// ErrUserHasAccount
func ErrUserHasAccount(userID, accountID uuid.UUID) ERMessage {
	return ERMessage(fmt.Sprintf("user %v has account with id %v", userID, accountID))
}

// ErrAccountAccess ...
type ErrAccountAccess struct {
	Reason  string
	message string
}

func (err ErrAccountAccess) Error() string {
	msg := fmt.Sprintf("couldn't access account. %v", err.Reason)
	return msg
}

// errAmountBelowMinimum
type errAmountBelowMinimum struct {
	MinAmount models.Rupees // minimum amount allowable for deposit or withdraw
	Message   ERMessage
}

func (err errAmountBelowMinimum) Error() string {
	return string(err.Message) + " " + strconv.Itoa(int(err.MinAmount))
}

func ErrAmountBelowMinimum(min models.Rupees, message ERMessage) error {
	return errAmountBelowMinimum{MinAmount: min, Message: message}
}

// ErrNotEnoughBalance
type ErrNotEnoughBalance struct {
	Message ERMessage
	Amount  models.Rupees
	Balance float64
}

func (err ErrNotEnoughBalance) Error() string {
	return fmt.Sprintf("%s. Amount: %v Balance: %v", string(err.Message), err.Amount, err.Balance)
}
