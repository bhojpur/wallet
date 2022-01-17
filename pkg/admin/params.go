package admin

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
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofrs/uuid"
)

// LoginParams are properties required during login of an admin
type LoginParams struct {
	Email    string `json:"email" schema:"email" form:"email"`
	Password string `json:"password" schema:"password" form:"password"`
}

func (req LoginParams) Validate() error {

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Password, validation.Required.Error(string(errors.ErrorPasswordRequired))),
		validation.Field(&req.Email, validation.Required.Error(string(errors.ErrorEmailRequired)), is.EmailFormat),
	)

	return errors.ParseValidationErrorMap(err)
}

// RegistrationParams are properties required during registration of a new admin
type RegistrationParams struct {
	FirstName string `json:"firstName" schema:"firstName" form:"firstName"`
	LastName  string `json:"lastName" schema:"lastName" form:"lastName"`
	Email     string `json:"email" schema:"email" form:"email"`
	Password  string `json:"password" schema:"password" form:"password"`
}

func (req RegistrationParams) Validate() error {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Email, validation.Required.Error(string(errors.ErrorEmailRequired)), is.EmailFormat),
		validation.Field(&req.Password, validation.Required.Error(string(errors.ErrorPasswordRequired))),
		validation.Field(&req.FirstName, validation.Required.Error(string(errors.ErrorFirstNameRequired))),
		validation.Field(&req.LastName, validation.Required.Error(string(errors.ErrorLastNameRequired))),
	)

	return errors.ParseValidationErrorMap(err)
}

type AssignFloatParams struct {
	AgentAccountNumber string        `json:"accountNo" schema:"accountNo" form:"accountNo"`
	Amount             models.Rupees `json:"amount" schema:"amount" form:"amount"`
}

func (req AssignFloatParams) Validate() error {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.AgentAccountNumber, validation.Required.Error(string(errors.ErrorAccountNumberRequired))),
		validation.Field(&req.Amount, validation.Required.Error(string(errors.ErrorAmountRequired))),
	)

	return errors.ParseValidationErrorMap(err)
}

type UpdateChargeParams struct {
	ChargeID uuid.UUID     `json:"chargeId" schema:"chargeId" form:"chargeId"`
	Amount   models.Paisas `json:"amount" schema:"amount" form:"amount"`
}

func (req UpdateChargeParams) Validate() error {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.ChargeID, validation.Required.Error(string(errors.ErrorChargeIDRequired))),
		validation.Field(&req.Amount, validation.Required.Error(string(errors.ErrorAmountRequired))),
	)

	return errors.ParseValidationErrorMap(err)
}
