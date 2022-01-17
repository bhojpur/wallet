package data

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

	"github.com/bhojpur/wallet/pkg/models"

	"github.com/gofrs/uuid"
)

// CustomerContract describe the characteristics of data that should
// be passed along in channels for when a user is created or something.
type CustomerContract struct {
	UserID uuid.UUID
}

type ChanNewCustomers struct {
	Channel chan CustomerContract
	Reader  <-chan CustomerContract
	Writer  chan<- CustomerContract
}

// TransactionContract represents the type of data
// required to record a new transaction in the database.
type TransactionContract struct {
	UserID       uuid.UUID
	AccountID    uuid.UUID
	Amount       float64
	TxnOperation models.TxnOperation // transaction operation type
	Timestamp    time.Time
}

type ChanNewTransactions struct {
	Channel chan TransactionContract
	Reader  <-chan TransactionContract
	Writer  chan<- TransactionContract
}

// type ChanNewTxnEvents struct {
// 	Channel chan models.TxnEvent
// 	Reader  <-chan models.TxnEvent
// 	Writer  chan<- models.TxnEvent
// }
