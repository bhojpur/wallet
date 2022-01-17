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

	"github.com/bhojpur/wallet/pkg/data"
	"github.com/bhojpur/wallet/pkg/errors"
	"github.com/bhojpur/wallet/pkg/models"
)

const (
	minimumDepositAmount    = models.Rupees(10) // least possible amount that can be deposited into an account
	minimumWithdrawalAmount = models.Rupees(1)  // least possible amount that can be withdrawn from an account
	minimumTransferAmount   = models.Rupees(10) // least possible amount that can be transferred to another account
)

type Interactor interface {
	AddTransaction(models.Transaction) error
}

type interactor struct {
	repository   Repository
	transChannel data.ChanNewTransactions
	// txnEventsChannel data.ChanNewTxnEvents
}

func NewInteractor(repository Repository, transChan data.ChanNewTransactions) Interactor {
	intr := &interactor{
		repository:   repository,
		transChannel: transChan,
	}

	go intr.listenOnCreatedTransactions()

	return intr
}

func (i interactor) AddTransaction(tx models.Transaction) error {
	_, err := i.repository.Add(tx)
	if err != nil {
		// if we get an error we are going to add the
		// transaction into a buffer object that will
		// retry adding the transaction at a later time

		return err
	}
	return nil
}

// func (i interactor) listenOnTxnEvents() {
// 	for {
// 		select {
// 		case event := <-i.txnEventsChannel.Reader:
// 			err := i.transact(event)
// 			if err != nil { // if we get an error processing the transaction event, we can log, or send to customer
// 				log.Println(err)
// 				continue
// 			}
// 		}
// 	}
// }

func (i interactor) listenOnCreatedTransactions() {
	for {
		select {
		case tx := <-i.transChannel.Reader:
			transaction := parseToTransaction(tx)

			err := i.AddTransaction(*transaction)
			if err != nil { // if we get an error, it is unexpected, we log it
				log.Printf("error happened when adding transaction to db %v", err.(errors.Error).Err)
				return
			}
			log.Printf("Transaction %v has been successfully added.", transaction.ID)
		}
	}
}
