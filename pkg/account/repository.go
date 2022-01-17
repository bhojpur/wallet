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
	"github.com/bhojpur/wallet/pkg/storage"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	GetAccountByUserID(uuid.UUID) (models.Account, error)
	UpdateBalance(amount models.Paisas, userID uuid.UUID) (models.Account, error)

	Create(userId uuid.UUID) (models.Account, error)
}

type repository struct {
	db *storage.Database
}

func NewRepository(database *storage.Database) Repository {
	return &repository{db: database}
}

// GetAccountByUserID fetches an account tied to a user's id
func (r repository) GetAccountByUserID(userID uuid.UUID) (models.Account, error) {
	var acc models.Account
	result := r.db.Where(models.Account{UserID: userID}).First(&acc)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Account{}, errors.Error{Code: errors.ENOTFOUND}
	} else if result.Error != nil {
		return models.Account{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return acc, nil
}

// UpdateBalance
func (r repository) UpdateBalance(amount models.Paisas, userID uuid.UUID) (models.Account, error) {
	var acc models.Account
	result := r.db.Model(models.Account{}).Where(models.Account{UserID: userID}).Updates(models.Account{AvailableBalance: amount}).Scan(&acc)
	if err := result.Error; err != nil {
		return models.Account{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return acc, nil
}

// Create a now account for userId
func (r repository) Create(userId uuid.UUID) (models.Account, error) {
	// check if user has an account and return it, otherwise create an account for user
	acc := zeroAccount(userId)
	result := r.db.Where(models.Account{UserID: userId}).FirstOrCreate(&acc)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := result.Error.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return models.Account{}, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrUserHasAccount(userId, acc.ID)}
		}
		return models.Account{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return acc, nil
}

func zeroAccount(userId uuid.UUID) models.Account {
	id, _ := uuid.NewV4()

	return models.Account{
		ID: id,
		// balance:     0, // no need to initialize with zero value, Go will do that for us
		Status:      models.StatusActive,
		AccountType: models.AccTypeCurrent,
		UserID:      userId,
	}
}
