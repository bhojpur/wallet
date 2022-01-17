package models

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
	"time"

	"github.com/gofrs/uuid"
)

type TxnOperation string

const (
	TxnOpDeposit  = TxnOperation("DEPOSIT")
	TxnOpWithdraw = TxnOperation("WITHDRAW")
	TxnOpTransfer = TxnOperation("TRANSFER")

	// only used when an admin is assigning float to a super agent
	TxnFloatAssignment = TxnOperation("FLOAT_ASSIGNMENT")
)

type TxnState string

const (
	TxStateCreated = TxnState("CREATED")
	TxStateFailed  = TxnState("FAILED")
)

type Transaction struct {
	ID        uuid.UUID
	Operation TxnOperation
	Timestamp time.Time
	Amount    float64
	UserID    uuid.UUID
	AccountID uuid.UUID
}

// TxnEvent is a description of a transaction operation event. We have defined operations
// as deposit, withdrawal and transfer. In the end all 3 operations can be modelled as one;
// "transfer" operations.
//
// Example:
// 1. During a deposit, money is moved from an agent's account to the depositor's account
// 2. During a withdrawal, money is moved from the customer withdrawing to the agent's account.
//
// A transfer operation/transaction, needs to have a source and destination and the amount being
// transferred.
// type TxnEvent struct {
// 	Source      TxnCustomer // where money is coming from
// 	Destination TxnCustomer // where money is going
// 	// we can further use this field to describe the specific type of transaction/transfer
// 	TxnOperation TxnOperation
// 	// transaction state to track the transaction
// 	TxnState TxnState
// 	// amount of money if Rupees being transacted
// 	Amount Rupees
// }

// TxnCustomer is a description of a customer involved in a transaction. We can describe them
// by their user id and user type; We have defined a customer being an agent, merchant or subscriber.
type TxnCustomer struct {
	UserID   uuid.UUID
	UserType UserType
}

// IsValidTxnOperation returns true if the given operation is among the defined
func IsValidTxnOperation(operation TxnOperation) bool {
	validOps := [3]TxnOperation{TxnOpDeposit, TxnOpWithdraw, TxnOpTransfer}
	for _, op := range validOps {
		if op == operation {
			return true
		}
	}
	return false
}
