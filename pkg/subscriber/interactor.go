package subscriber

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
	"github.com/bhojpur/wallet/pkg/config"
	"github.com/bhojpur/wallet/pkg/data"
	"github.com/bhojpur/wallet/pkg/errors"
	"github.com/bhojpur/wallet/pkg/helpers"
	"github.com/bhojpur/wallet/pkg/models"
)

type Interactor interface {
	AuthenticateByEmail(email, password string) (models.Subscriber, error)
	Register(RegistrationParams) (models.Subscriber, error)
}

func NewInteractor(config config.Config, subsRepo Repository, custChan data.ChanNewCustomers) Interactor {
	return &interactor{
		config:           config,
		repository:       subsRepo,
		customersChannel: custChan,
	}
}

type interactor struct {
	customersChannel data.ChanNewCustomers
	config           config.Config
	repository       Repository
}

// AuthenticateByEmail verifies a subscriber by the provided unique email address
func (ui interactor) AuthenticateByEmail(email, password string) (models.Subscriber, error) {
	// search for subscriber by email.
	subscriber, err := ui.repository.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return models.Subscriber{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return models.Subscriber{}, err
	}

	// validate password
	if err := helpers.ComparePasswordToHash(subscriber.Password, password); err != nil {
		return models.Subscriber{}, errors.Unauthorized{Message: errors.ErrorMessage(err)}
	}

	return subscriber, nil
}

// Register takes in a subscriber registration parameters and creates a new subscriber
// then adds the subscriber to db.
func (ui interactor) Register(params RegistrationParams) (models.Subscriber, error) {
	subscriber := models.Subscriber{
		FirstName:   params.FirstName,
		LastName:    params.LastName,
		Email:       params.Email,
		PhoneNumber: params.PhoneNumber,
		Password:    params.Password,
		PassportNo:  params.PassportNo,
	}

	// hash subscriber password before adding to db.
	passwordHash, err := helpers.HashPassword(subscriber.Password)
	if err != nil { // if we get an error, it means our hashing func dint work
		return models.Subscriber{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	// change password to hashed string
	subscriber.Password = passwordHash
	sub, err := ui.repository.Add(subscriber)
	if err != nil {
		return models.Subscriber{}, err
	}

	// tell channel listeners that a new subscriber has been created.
	ui.postNewSubscriberToChannel(&sub)
	return sub, nil
}

// take the newly created subscriber and post them to channel
// that listens for newly created customers and acts upon them
// like creating an account for them automatically.
func (ui interactor) postNewSubscriberToChannel(subscriber *models.Subscriber) {
	newSubscriber := parseToNewSubscriber(*subscriber)
	go func() { ui.customersChannel.Writer <- newSubscriber }()
}
