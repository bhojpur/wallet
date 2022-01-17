package tariff

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
	"github.com/bhojpur/wallet/pkg/models"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Charge struct {
	ID uuid.UUID

	Transaction         models.TxnOperation `gorm:"uniqueIndex:idx_unique_tx_identity"`
	SourceUserType      models.UserType     `gorm:"uniqueIndex:idx_unique_tx_identity"`
	DestinationUserType models.UserType     `gorm:"uniqueIndex:idx_unique_tx_identity"`
	Fee                 models.Paisas

	gorm.Model
}

func (t *Charge) BeforeCreate(tx *gorm.DB) error {
	t.ID, _ = uuid.NewV4()
	return nil
}

// ValidTransaction defines a format to identify all allowable transactions between customers
// It is an array of 2 user types since only 2 customers types are allowed in a valid transaction
// Index 0 is the source while index 1 is the destination
//
// Example
// 1. Validation{models.UserTypSubscriber, models.UserTypMerchant}
// Describes a valid transaction (transfer) between a subscriber to a merchant
type ValidTransaction [2]models.UserType
