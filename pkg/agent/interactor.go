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
	"log"

	"github.com/bhojpur/wallet/pkg/config"
	"github.com/bhojpur/wallet/pkg/data"
	"github.com/bhojpur/wallet/pkg/errors"
	"github.com/bhojpur/wallet/pkg/helpers"
	"github.com/bhojpur/wallet/pkg/models"
)

type Interactor interface {
	AuthenticateByEmail(email, password string) (models.Agent, error)
	Register(RegistrationParams) (models.Agent, error)
	UpdateSuperAgentStatus(email string) error
}

func NewInteractor(config config.Config, agentRepo Repository, custChan data.ChanNewCustomers) Interactor {
	return &interactor{
		config:           config,
		repository:       agentRepo,
		customersChannel: custChan,
	}
}

type interactor struct {
	customersChannel data.ChanNewCustomers
	config           config.Config
	repository       Repository
}

// AuthenticateByEmail verifies an agent by the provided unique email address
func (ui interactor) AuthenticateByEmail(email, password string) (models.Agent, error) {
	// search for agent by email.
	agent, err := ui.repository.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return models.Agent{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return models.Agent{}, err
	}

	// validate password
	if err := helpers.ComparePasswordToHash(agent.Password, password); err != nil {
		return models.Agent{}, errors.Unauthorized{Message: errors.ErrorMessage(err)}
	}

	return agent, nil
}

// Register takes in a agent object and adds the agent to db.
func (ui interactor) Register(params RegistrationParams) (models.Agent, error) {
	agent := models.Agent{
		FirstName:   params.FirstName,
		LastName:    params.LastName,
		Email:       params.Email,
		PhoneNumber: params.PhoneNumber,
		Password:    params.Password,
		PassportNo:  params.PassportNo,
	}

	// hash agent password before adding to db.
	passwordHash, err := helpers.HashPassword(agent.Password)
	if err != nil { // if we get an error, it means our hashing func dint work
		return models.Agent{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	// change password to hashed string
	agent.Password = passwordHash
	agt, err := ui.repository.Add(agent)
	if err != nil {
		return models.Agent{}, err
	}

	// tell channel listeners that a new agent has been created.
	ui.postNewAgentToChannel(&agt)
	return agt, nil
}

func (ui interactor) UpdateSuperAgentStatus(email string) error {
	agent, err := ui.repository.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return err
	}

	log.Printf("status before: %v", agent.SuperAgent)
	agent.SuperAgent = agent.SuperAgent.Not()
	log.Printf("status after: %v", agent.SuperAgent)
	err = ui.repository.Update(agent)
	if err != nil {
		return err
	}

	return nil
}

// take the newly created agent and post them to channel
// that listens for newly created customers and acts upon them
// like creating an account for them automatically.
func (ui interactor) postNewAgentToChannel(agent *models.Agent) {
	newAgent := parseToNewAgent(*agent)
	go func() { ui.customersChannel.Writer <- newAgent }()
}
