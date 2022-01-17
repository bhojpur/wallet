package customer

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
	"github.com/bhojpur/wallet/pkg/agent"
	"github.com/bhojpur/wallet/pkg/errors"
	"github.com/bhojpur/wallet/pkg/merchant"
	"github.com/bhojpur/wallet/pkg/models"
	"github.com/bhojpur/wallet/pkg/subscriber"

	"github.com/gofrs/uuid"
)

type Finder interface {
	FindAgentByEmail(string) (models.Agent, error)
	FindMerchantByEmail(string) (models.Merchant, error)
	FindSubscriberByEmail(string) (models.Subscriber, error)
	FindIDByEmail(string, models.UserType) (uuid.UUID, error)
}

func NewFinder(agentRepo agent.Repository, merchRepo merchant.Repository, subRepo subscriber.Repository) Finder {
	return &finder{
		agentRepo: agentRepo,
		merchRepo: merchRepo,
		subRepo:   subRepo,
	}
}

type finder struct {
	agentRepo agent.Repository
	merchRepo merchant.Repository
	subRepo   subscriber.Repository
}

func (f finder) FindAgentByEmail(email string) (models.Agent, error) {
	agt, err := f.agentRepo.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return models.Agent{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return models.Agent{}, err
	}

	return agt, nil
}

func (f finder) FindMerchantByEmail(email string) (models.Merchant, error) {
	merch, err := f.merchRepo.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return models.Merchant{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return models.Merchant{}, err
	}

	return merch, nil
}

func (f finder) FindSubscriberByEmail(email string) (models.Subscriber, error) {
	sub, err := f.subRepo.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return models.Subscriber{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return models.Subscriber{}, err
	}

	return sub, nil
}

func (f finder) FindIDByEmail(email string, userType models.UserType) (uuid.UUID, error) {
	switch userType {
	case models.UserTypAgent:
		agt, err := f.FindAgentByEmail(email)
		if err != nil {
			return uuid.Nil, err
		}
		return agt.ID, nil
	case models.UserTypMerchant:
		merch, err := f.FindMerchantByEmail(email)
		if err != nil {
			return uuid.Nil, err
		}
		return merch.ID, nil
	case models.UserTypSubscriber:
		sub, err := f.FindSubscriberByEmail(email)
		if err != nil {
			return uuid.Nil, err
		}
		return sub.ID, nil
	}
	return uuid.Nil, errors.Error{Code: errors.EINVALID, Message: errors.ErrUserNotFound}
}
