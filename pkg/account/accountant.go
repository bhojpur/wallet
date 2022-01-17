package account

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
	"github.com/bhojpur/wallet/pkg/errors"
	"github.com/bhojpur/wallet/pkg/models"
	"github.com/bhojpur/wallet/pkg/statement"

	"github.com/gofrs/uuid"
)

type Accountant interface {
	DebitAccount(userID uuid.UUID, amount models.Paisas, reason models.TxnOperation) (float64, error)
	CreditAccount(userID uuid.UUID, amount models.Paisas, reason models.TxnOperation) (float64, error)
}

func NewAccountant(accountRepo Repository, ledger statement.Ledger) Accountant {
	return &accountant{repository: accountRepo, ledger: ledger}
}

type accountant struct {
	ledger     statement.Ledger
	repository Repository
}

func (a accountant) isUserAccAccessible(userID uuid.UUID) (*models.Account, error) {
	acc, err := a.repository.GetAccountByUserID(userID)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return nil, errors.Error{Message: errors.AccountNotCreated, Err: err}
	} else if err != nil {
		return nil, err
	}

	if acc.Status == models.StatusFrozen || acc.Status == models.StatusSuspended {
		e := errors.ErrAccountAccess{Reason: string(acc.Status)}
		return nil, errors.Error{Err: e}
	}

	return &acc, nil

}

func (a accountant) CreditAccount(userID uuid.UUID, amount models.Paisas, reason models.TxnOperation) (float64, error) {
	acc, err := a.isUserAccAccessible(userID)
	if err != nil {
		return 0, err
	}

	// update balance with amount: add amount
	amt := acc.Credit(amount)
	*acc, err = a.repository.UpdateBalance(amt, userID)
	if err != nil {
		return 0, err
	}

	err = a.ledger.Record(userID, *acc, reason, amount.ToRupees(), statement.TypeCredit)
	if err != nil {
		return 0, err
	}

	return acc.Balance(), nil
}

func (a accountant) DebitAccount(userID uuid.UUID, amount models.Paisas, reason models.TxnOperation) (float64, error) {
	acc, err := a.isUserAccAccessible(userID)
	if err != nil {
		return 0, err
	}

	// check that balance is more than amount
	if acc.IsBalanceLessThanAmount(amount) {
		e := errors.ErrNotEnoughBalance{
			Message: errors.DebitAmountAboveBalance,
			Amount:  amount.ToRupees(),
			Balance: acc.Balance(),
		}
		return 0, errors.Error{Err: e}
	}

	// update balance with amount: subtract amount
	amt := acc.Debit(amount)
	*acc, err = a.repository.UpdateBalance(amt, userID)
	if err != nil {
		return 0, err
	}

	err = a.ledger.Record(userID, *acc, reason, amount.ToRupees(), statement.TypeDebit)
	if err != nil {
		return 0, err
	}

	return acc.Balance(), nil
}
