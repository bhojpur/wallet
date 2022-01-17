package responses

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
	"time"

	"github.com/bhojpur/wallet/pkg/models"
	"github.com/bhojpur/wallet/pkg/statement"

	"github.com/gofrs/uuid"
)

type transactionStatement struct {
	ID             uuid.UUID           `json:"transactionId"`
	Type           models.TxnOperation `json:"transactionType"`
	CreatedAt      time.Time           `json:"createdAt"`
	CreditedAmount float64             `json:"creditedAmount"`
	DebitedAmount  float64             `json:"debitedAmount"`
	UserID         uuid.UUID           `json:"userId"`
	AccountID      uuid.UUID           `json:"accountId"`
}

type miniStatementResponse struct {
	Message string    `json:"message"`
	UserID  uuid.UUID `json:"userID"`

	Statements []transactionStatement `json:"transactions"`
}

func MiniStatementResponse(userID uuid.UUID, statements []statement.Statement) SuccessResponse {
	var stmts []transactionStatement
	for _, stmt := range statements {
		stmts = append(stmts, transactionStatement{
			ID:             stmt.ID,
			Type:           stmt.Operation,
			CreatedAt:      stmt.CreatedAt,
			CreditedAmount: stmt.CreditAmount,
			DebitedAmount:  stmt.DebitAmount,
			UserID:         stmt.UserID,
			AccountID:      stmt.AccountID,
		})
	}

	msg := "mini statement retrieved for the past 5 transactions"
	data := miniStatementResponse{
		Message:    msg,
		UserID:     userID,
		Statements: stmts,
	}

	return successResponse(msg, data)
}

type transactionResponse struct {
	Message string `json:"message"`
}

func TransactionResponse() SuccessResponse {
	data := transactionResponse{
		Message: "Transaction under processing. You will receive a message shortly.",
	}
	return successResponse("", data)
}

type balanceResponse struct {
	UserID  uuid.UUID `json:"userID"`
	Balance float64   `json:"balance"`
}

func BalanceResponse(userID uuid.UUID, balance float64) SuccessResponse {
	msg := fmt.Sprintf("Your current balance is %v", balance)

	data := balanceResponse{
		UserID:  userID,
		Balance: balance,
	}
	return successResponse(msg, data)
}
