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

// Agent
type Agent struct {
	ID    uuid.UUID
	Email string `gorm:"not null;unique"` // email is used as account number

	FirstName   string
	LastName    string
	PhoneNumber string `gorm:"not null;unique"`
	PassportNo  string
	Password    string `gorm:"not null"`

	// an agent is usually assigned an agent number that they use for
	// transactions with other customers
	// AgentNumber string `gorm:"column:agent_number;unique"`

	// an extra column/property that tells us if the agent is a super agent
	SuperAgent SuperAgentStatus `gorm:"default:'0'"` // PS: bool values dont work well with gorm during updates

	gorm.Model
}

// BeforeCreate hook will be used to add uuid to entity before adding to db
func (u *Agent) BeforeCreate(tx *gorm.DB) error {
	u.ID, _ = uuid.NewV4()
	return nil
}

func (u Agent) IsSuperAgent() bool {
	return u.SuperAgent == IsSuperAgent
}

// SuperAgentStatus
type SuperAgentStatus string

const (
	IsNotSuperAgent = SuperAgentStatus('0')
	IsSuperAgent    = SuperAgentStatus('1')
)

// Not returns the opposite, if that makes sense
func (status SuperAgentStatus) Not() SuperAgentStatus {
	if status == IsNotSuperAgent {
		return IsSuperAgent
	} else {
		return IsNotSuperAgent
	}
}
