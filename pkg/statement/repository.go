package statement

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

	"github.com/bhojpur/wallet/pkg/errors"
	"github.com/bhojpur/wallet/pkg/storage"

	"github.com/gofrs/uuid"
)

type Repository interface {
	Add(Statement) (Statement, error)
	GetStatements(userID uuid.UUID, from time.Time, limit uint) ([]Statement, error)
}

func NewRepository(database *storage.Database) Repository {
	return &repository{db: database}
}

type repository struct {
	db *storage.Database
}

func (r repository) Add(stmt Statement) (Statement, error) {
	result := r.db.Create(&stmt)
	if err := result.Error; err != nil {
		return Statement{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return stmt, nil
}

func (r repository) GetStatements(userID uuid.UUID, from time.Time, limit uint) ([]Statement, error) {
	var statements []Statement

	result := r.db.Where(
		Statement{UserID: userID},
	).Where(
		"created_at <= ?", from,
	).Order("created_at desc").Limit(int(limit)).Find(&statements)

	if err := result.Error; err != nil {
		return nil, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return statements, nil
}
