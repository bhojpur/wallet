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
	"github.com/bhojpur/wallet/pkg/data"
)

type Channels struct {
	ChannelNewUsers        data.ChanNewCustomers
	ChannelNewTransactions data.ChanNewTransactions
	// ChannelTxnEvents       data.ChanNewTxnEvents
}

func NewChannels() *Channels {
	chanNewUsers := make(chan data.CustomerContract, 10)
	chanNewTransactions := make(chan data.TransactionContract, 50)
	// chanNewTxnEvents := make(chan models.TxnEvent, 100)

	return &Channels{
		ChannelNewUsers: data.ChanNewCustomers{
			Channel: chanNewUsers,
			Reader:  chanNewUsers,
			Writer:  chanNewUsers,
		},
		ChannelNewTransactions: data.ChanNewTransactions{
			Channel: chanNewTransactions,
			Reader:  chanNewTransactions,
			Writer:  chanNewTransactions,
		},
		// ChannelTxnEvents: data.ChanNewTxnEvents{
		// 	Channel: chanNewTxnEvents,
		// 	Reader:  chanNewTxnEvents,
		// 	Writer:  chanNewTxnEvents,
		// },
	}
}
