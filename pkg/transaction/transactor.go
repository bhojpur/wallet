package transaction

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

	"github.com/bhojpur/wallet/pkg/account"
	"github.com/bhojpur/wallet/pkg/errors"
	"github.com/bhojpur/wallet/pkg/models"
	"github.com/bhojpur/wallet/pkg/tariff"
)

type Transactor interface {
	Transact(Transaction) error
}

func NewTransactor(accountant account.Accountant, manager tariff.Manager) Transactor {
	return &transactor{accountant: accountant, tariff: manager}
}

type transactor struct {
	accountant account.Accountant
	tariff     tariff.Manager
}

// in mobile money a deposit will happen from the account of an agent to the other customer. The source is the agent's
// account and destination is the account of the other customer.
func (tr transactor) deposit(source, destination models.TxnCustomer, amount models.Rupees) error {
	if amount < minimumDepositAmount {
		e := errors.ErrAmountBelowMinimum(minimumDepositAmount, errors.DepositAmountBelowMinimum)
		return errors.Error{Err: e}
	}

	// the source should always be an agent
	// a super agent too is allowed to do deposits to other agents
	if !source.UserType.IsAgent() {
		return errors.Error{Code: errors.EINVALID, Message: errors.DepositOnlyAtAgent}
	}

	// a super agent is only allowed to deposit to another agent's account
	if source.UserType == models.UserTypSuperAgent && destination.UserType != models.UserTypAgent {
		return errors.Error{Code: errors.EINVALID, Message: errors.SuperAgentCantDeposit}
	}

	// a merchant is not allowed to deposit
	if destination.UserType == models.UserTypMerchant {
		return errors.Error{Code: errors.EINVALID, Message: errors.CustomerCantDeposit}
	}

	// get the charge applicable to this transaction
	// usually depositing has no transaction cost

	srcNewBal, err := tr.accountant.DebitAccount(source.UserID, amount.ToPaisas(), models.TxnOpDeposit)
	if err != nil {
		return err
	}

	destNewBal, err := tr.accountant.CreditAccount(destination.UserID, amount.ToPaisas(), models.TxnOpDeposit)
	if err != nil {
		return err
	}

	log.Printf("Source balance: %v || Dest balance: %v", srcNewBal, destNewBal)

	return nil
}

// in mobile money a withdrawal will happen from the account of the customer withdrawing to the agent. The source is the
// customer's account and the destination is the account of the agent
func (tr transactor) withdraw(source, destination models.TxnCustomer, amount models.Rupees) error {
	if amount < minimumWithdrawalAmount {
		e := errors.ErrAmountBelowMinimum(minimumWithdrawalAmount, errors.WithdrawAmountBelowMinimum)
		return errors.Error{Err: e}
	}

	// a super agent cannot perform withdrawals for customers or withdraw
	if destination.UserType == models.UserTypSuperAgent || source.UserType == models.UserTypSuperAgent {
		return errors.Error{Code: errors.EINVALID, Message: errors.SuperAgentCantWithdraw}
	}

	// the destination should always be an agent
	if destination.UserType != models.UserTypAgent {
		return errors.Error{Code: errors.EINVALID, Message: errors.WithdrawalOnlyAtAgent}
	}

	// we can implement a double withdrawal check here. That will prevent a user from
	// withdrawing same amount twice within a stipulated time interval because of system lag.

	// get the charge applicable to this transaction
	charge, err := tr.tariff.GetCharge(models.TxnOpWithdraw, source.UserType, destination.UserType)
	if err != nil {
		return err
	}

	// we apply a transaction fee to the transaction
	// when withdrawing the source is charged the fee (customer)
	amt := amount.ToPaisas() + charge

	srcNewBal, err := tr.accountant.DebitAccount(source.UserID, amt, models.TxnOpWithdraw)
	if err != nil {
		return err
	}

	destNewBal, err := tr.accountant.CreditAccount(destination.UserID, amount.ToPaisas(), models.TxnOpWithdraw)
	if err != nil {
		return err
	}

	log.Printf("Source balance: %v || Dest balance: %v", srcNewBal, destNewBal)

	return nil
}

func (tr transactor) transfer(source, destination models.TxnCustomer, amount models.Rupees) error {
	if amount < minimumTransferAmount {
		e := errors.ErrAmountBelowMinimum(minimumTransferAmount, errors.TransferAmountBelowMinimum)
		return errors.Error{Err: e}
	}

	// a super agent is not allowed to make a transfer
	// can only do a deposit
	if source.UserType == models.UserTypSuperAgent || destination.UserType == models.UserTypSuperAgent {
		return errors.Error{Code: errors.EINVALID, Message: errors.SuperAgentCantTransfer}
	}

	// get the charge applicable to this transaction
	charge, err := tr.tariff.GetCharge(models.TxnOpTransfer, source.UserType, destination.UserType)
	if err != nil {
		return err
	}

	// we apply the transaction fee to the transaction
	// when making a transfer, it is the source that is charged
	amt := amount.ToPaisas() + charge

	srcNewBal, err := tr.accountant.DebitAccount(source.UserID, amt, models.TxnOpTransfer)
	if err != nil {
		return err
	}

	destNewBal, err := tr.accountant.CreditAccount(destination.UserID, amount.ToPaisas(), models.TxnOpTransfer)
	if err != nil {
		return err
	}

	log.Printf("Source balance: %v || Dest balance: %v", srcNewBal, destNewBal)

	return nil
}

func (tr transactor) Transact(transaction Transaction) error {
	// cannot transact within the same account
	if transaction.Source.UserID == transaction.Destination.UserID {
		return errors.Error{Code: errors.EINVALID, Message: errors.TransactionWithSameAccount}
	}

	switch transaction.TxnOperation {
	case models.TxnOpDeposit:
		return tr.deposit(transaction.Source, transaction.Destination, transaction.Amount)
	case models.TxnOpWithdraw:
		return tr.withdraw(transaction.Source, transaction.Destination, transaction.Amount)
	case models.TxnOpTransfer:
		return tr.transfer(transaction.Source, transaction.Destination, transaction.Amount)
	}
	return errors.Error{Code: errors.EINVALID}
}
