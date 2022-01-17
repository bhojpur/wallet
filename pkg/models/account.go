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
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// AccountStatus (active,dormant,frozen,suspended)
type AccountStatus string

// AccountType (savings,current,utility)
type AccountType string

const (
	StatusActive    = AccountStatus("active")
	StatusDormant   = AccountStatus("dormant")
	StatusFrozen    = AccountStatus("frozen")
	StatusSuspended = AccountStatus("suspended")
)

const (
	// different types of accounts a user could hold
	// we will use current account only.
	AccTypeSavings = AccountType("savings")
	AccTypeCurrent = AccountType("current")
	AccTypeUtility = AccountType("utility")
)

// Account entity definition
type Account struct {
	ID uuid.UUID

	// balance will be stored in Paisas
	AvailableBalance Paisas `gorm:"column:available_balance"`

	Status      AccountStatus `gorm:"column:status"`
	AccountType AccountType   `gorm:"column:account_type"`
	UserID      uuid.UUID     `gorm:"column:user_id;not null;unique"` // a user can only have one account

	gorm.Model
}

// Balance converts balance from Paisas
func (acc Account) Balance() float64 {
	return float64(acc.AvailableBalance / 100)
}

// Credit add an amount to account balance and return it
func (acc Account) Credit(amount Paisas) Paisas {
	// convert incoming amount into Paisas and add to account balance
	return amount + acc.AvailableBalance
}

// Debit subtract an amount from account balance and return it
func (acc Account) Debit(amount Paisas) Paisas {
	// convert incoming amount into Paisas and subtract to account balance
	return acc.AvailableBalance - amount
}

// IsBalanceLessThanAmount converts amount into Paisas and returns true if balance is less than amount
func (acc Account) IsBalanceLessThanAmount(amount Paisas) bool {
	return acc.AvailableBalance < amount
}
