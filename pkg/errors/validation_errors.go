package errors

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
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	jsoniter "github.com/json-iterator/go"
)

// ValidationError
type ValidationError string

func (e ValidationError) Error() string {
	return string(e)
}

// MarshalJSON converts the Errors into a valid JSON.
func (e ValidationError) MarshalJSON() ([]byte, error) {
	errM := ApiErrorResponse{
		Error:   "Validation Error",
		Message: e.Error(),
		Status:  http.StatusBadRequest,
	}
	return jsoniter.Marshal(errM)
}

// ValidationErrors is a slice of validation errors
type ValidationErrors []error

func (e ValidationErrors) Error() string {
	var str []string
	for _, err := range e {
		str = append(str, err.Error())
	}
	return strings.Join(str, ",")
}

const (
	ErrorFirstNameRequired         = ValidationError("first_name is a required field")
	ErrorLastNameRequired          = ValidationError("last_name is a required field")
	ErrorEmailRequired             = ValidationError("email is a required field")
	ErrorPhoneNumberRequired       = ValidationError("phoneNumber is a required field")
	ErrorPasswordRequired          = ValidationError("password is a required field")
	ErrorInvalidUsernameOrPassword = ValidationError("provided wrong username or password")
	ErrorAmountRequired            = ValidationError("amount is a required field")
	ErrorAgentNumberRequired       = ValidationError("agentNumber is a required field")
	ErrorCustomerTypeRequired      = ValidationError("customerType is a required field")
	ErrorAgentIDRequired           = ValidationError("agentID is a required field")
	ErrorAccountNumberRequired     = ValidationError("accountNo is a required field")
	ErrorChargeIDRequired          = ValidationError("chargeId is a required field")
)

// ParseValidationErrorMap takes in the error map that go-ozzo validation
// framework returns and extracts the application error code type as a string
// and returns ErrorCode type of the specific error
func ParseValidationErrorMap(err error) error {

	if err != nil {
		var errs ValidationErrors

		// the validation framework returns a map of errors
		// we check if the error returned matches this map
		if e, ok := err.(validation.Errors); ok {
			// 	we range over the map and convert the map into slice of ValidationError type
			for _, v := range e {
				errs = append(errs, ValidationError(v.Error()))
			}
		}
		return errs
	}

	return nil
}
