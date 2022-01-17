package agent

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
	Add(models.Agent) (models.Agent, error)
	Delete(models.Agent) error
	FetchAll() ([]models.Agent, error)
	FindByID(uuid.UUID) (models.Agent, error)
	FindByEmail(string) (models.Agent, error)
	Update(models.Agent) error
}

func NewRepository(database *storage.Database) Repository {
	return &repository{db: database}
}

type repository struct {
	db *storage.Database
}

func (r repository) searchBy(row models.Agent) (models.Agent, error) {
	var agent models.Agent
	result := r.db.Where(row).First(&agent)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Agent{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return models.Agent{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return agent, nil
}

// Add an agent if not in db.
func (r repository) Add(agent models.Agent) (models.Agent, error) {
	// add new agent to agents table, if query return violation of unique key column,
	// we know that the agent with given record exists and return that agent instead
	result := r.db.Model(models.Agent{}).Create(&agent)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return agent, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrUserExists}
		}
		return models.Agent{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return agent, nil
}

// Delete a agent
func (r repository) Delete(agent models.Agent) error {
	result := r.db.Delete(&agent)
	if result.Error != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

// FetchAll gets all agents in db
func (r repository) FetchAll() ([]models.Agent, error) {
	var agents []models.Agent
	result := r.db.Find(&agents)
	if err := result.Error; err != nil {
		return nil, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	// we might not need to return this error
	if result.RowsAffected == 0 {
		return nil, errors.Error{Code: errors.ENOTFOUND}
	}

	return agents, nil
}

// FindByID searches agent by primary id
func (r repository) FindByID(id uuid.UUID) (models.Agent, error) {
	agent, err := r.searchBy(models.Agent{ID: id})
	return agent, err
}

// FindByEmail searches agent by email
func (r repository) FindByEmail(email string) (models.Agent, error) {
	agent, err := r.searchBy(models.Agent{Email: email})
	return agent, err
}

// Update
func (r repository) Update(agent models.Agent) error {
	var u models.Agent
	result := r.db.Debug().Model(&u).Where(models.Agent{ID: agent.ID}).Omit("id").Updates(agent)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}
