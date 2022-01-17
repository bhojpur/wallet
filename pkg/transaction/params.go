package transaction

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

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type DepositParams struct {
	Amount models.Rupees `json:"amount" schema:"amount" form:"amount"`
	// In a production system, the customer number is usually a generated number for the case
	// of a merchant of agent and a mobile number for a subscriber , but we are going to use the
	// customer's email as a replacement
	CustomerNumber string          `json:"accountNo" schema:"accountNo" form:"accountNo"`
	CustomerType   models.UserType `json:"customerType" schema:"customerType" form:"customerType"`
}

func (req DepositParams) Validate() error {

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Amount, validation.Required.Error(string(errors.ErrorAmountRequired))),
		validation.Field(&req.CustomerNumber, validation.Required.Error(string(errors.ErrorAccountNumberRequired))),
		validation.Field(&req.CustomerType, validation.Required.Error(string(errors.ErrorCustomerTypeRequired))),
	)

	return errors.ParseValidationErrorMap(err)
}

type TransferParams struct {
	Amount models.Rupees `json:"amount" schema:"amount" form:"amount"`

	// In a production system, the account number is usually
	// a generated number, but we are going to use the customer's
	// email as a replacement
	DestAccountNo string          `json:"accountNo" schema:"accountNo" form:"accountNo"`
	DestUserType  models.UserType `json:"customerType" schema:"customerType" form:"customerType"`
}

func (req TransferParams) Validate() error {

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Amount, validation.Required.Error(string(errors.ErrorAmountRequired))),
		validation.Field(&req.DestAccountNo, validation.Required.Error(string(errors.ErrorAccountNumberRequired))),
		validation.Field(&req.DestUserType, validation.Required.Error(string(errors.ErrorCustomerTypeRequired))),
	)

	return errors.ParseValidationErrorMap(err)
}

type WithdrawParams struct {
	Amount models.Rupees `json:"amount" schema:"amount" form:"amount"`
	// In a production system, the agent number is usually
	// a generated number, but we are going to use the agent's
	// email as a replacement
	AgentNumber string `json:"agentNumber" schema:"agentNumber" form:"agentNumber"`
}

func (req WithdrawParams) Validate() error {

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Amount, validation.Required.Error(string(errors.ErrorAmountRequired))),
		validation.Field(&req.AgentNumber, validation.Required.Error(string(errors.ErrorAgentNumberRequired))),
	)

	return errors.ParseValidationErrorMap(err)
}
