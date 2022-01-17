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
	"bytes"
	"errors"
	"fmt"
)

// ERCode defines a machine readable error code type
type ERCode string

// ERMessage defines a string type of an error description
type ERMessage string

const (
	ECONFLICT = ERCode("conflict")  // action cannot be performed e.g. when inserting existing record to db
	EINTERNAL = ERCode("internal")  // internal error
	EINVALID  = ERCode("invalid")   // validation failed
	ENOTFOUND = ERCode("not_found") // entity does not exist
)

// Error is our standard application error
type Error struct {
	// 	Machine readable error code
	Code ERCode

	// 	Human readable message
	Message ERMessage

	// nested error
	Err error
}

func (e Error) Error() string {
	var buf bytes.Buffer

	if e.Message != "" {
		if e.Code != "" {
			fmt.Fprintf(&buf, "<%s> ", e.Code)
		}
		buf.WriteString(string(e.Message))
	} else if e.Err != nil {
		// 	If wrapping an error, print its Error() message.
		// otherwise print the error code & message.
		buf.WriteString(e.Err.Error())
	}

	return buf.String()
}

// ErrorCode returns the code of the root error, if available. Otherwise returns EINTERNAL.
//
// 1. Returns no error code for nil errors.
// 2. Searches the chain of Error.Err until a defined Code is found.
// 3. If no code is defined then return an internal error code (EINTERNAL).
func ErrorCode(err error) ERCode {
	if err == nil {
		return ""
	} else if e, ok := err.(Error); ok && e.Code != "" {
		return e.Code
	} else if ok && e.Err != nil {
		return ErrorCode(e.Err)
	}
	return EINTERNAL
}

// ErrorMessage returns the human-readable message of the error, if available.
// Otherwise returns a generic error message.
//
// 1. Returns no error message for nil errors.
// 2. Searches the chain of Error.Err until a defined Message is found.
// 3. If no message is defined then return a generic error message.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(Error); ok && e.Message != "" {
		return string(e.Message)
	} else if ok && e.Err != nil {
		return ErrorMessage(e.Err)
	}
	return "An internal error has occurred. Please contact technical support."
}

// Is reports whether any error in err's chain matches target.
func Is(err, target error) bool {
	return errors.Is(err, target)
}
