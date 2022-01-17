package registry

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
	"github.com/bhojpur/wallet/pkg/account"
	"github.com/bhojpur/wallet/pkg/admin"
	"github.com/bhojpur/wallet/pkg/agent"
	"github.com/bhojpur/wallet/pkg/config"
	"github.com/bhojpur/wallet/pkg/customer"
	"github.com/bhojpur/wallet/pkg/merchant"
	"github.com/bhojpur/wallet/pkg/ports"
	"github.com/bhojpur/wallet/pkg/statement"
	"github.com/bhojpur/wallet/pkg/storage"
	"github.com/bhojpur/wallet/pkg/subscriber"
	"github.com/bhojpur/wallet/pkg/tariff"
	"github.com/bhojpur/wallet/pkg/transaction"
)

type Domain struct {
	Admin      admin.Interactor
	Agent      agent.Interactor
	Merchant   merchant.Interactor
	Subscriber subscriber.Interactor

	Account     account.Interactor
	Transaction transaction.Interactor
	Statement   statement.Interactor
	Tariff      tariff.Manager

	Transactor ports.TransactorPort
}

func NewDomain(config config.Config, database *storage.Database, channels *Channels) *Domain {
	adminRepo := admin.NewRepository(database)
	agentRepo := agent.NewRepository(database)
	merchantRepo := merchant.NewRepository(database)
	subscriberRepo := subscriber.NewRepository(database)

	accRepo := account.NewRepository(database)
	txnRepo := transaction.NewRepository(database)
	statementRepo := statement.NewRepository(database)
	tariffRepo := tariff.NewRepository(database)

	// initialize ports and adapters
	ledger := statement.NewLedger(statementRepo)
	tariffManager := tariff.NewManager(tariffRepo)
	accountant := account.NewAccountant(accRepo, ledger)
	customerFinder := customer.NewFinder(agentRepo, merchantRepo, subscriberRepo)
	transactor := transaction.NewTransactor(accountant, tariffManager)

	return &Domain{
		Admin:       admin.NewInteractor(config, adminRepo, accountant, customerFinder),
		Agent:       agent.NewInteractor(config, agentRepo, channels.ChannelNewUsers),
		Merchant:    merchant.NewInteractor(config, merchantRepo, channels.ChannelNewUsers),
		Subscriber:  subscriber.NewInteractor(config, subscriberRepo, channels.ChannelNewUsers),
		Account:     account.NewInteractor(accRepo, channels.ChannelNewUsers, channels.ChannelNewTransactions),
		Transaction: transaction.NewInteractor(txnRepo, channels.ChannelNewTransactions),
		Statement:   statement.NewInteractor(statementRepo),
		Transactor:  ports.NewTransactor(customerFinder, transactor),
		Tariff:      tariffManager,
	}
}
