package error_handlers

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

	"github.com/bhojpur/wallet/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler provides a custom error handling mechanism for fiber framework
func ErrorHandler(ctx *fiber.Ctx, err error) error {

	// if error corresponds to unauthorized
	if e, ok := err.(errors.Unauthorized); ok {
		log.Println(err)
		res := errors.UnauthorizedResponse(e.Error())
		return ctx.Status(res.Status).JSON(res)
	}

	// if error is our custom validation errors slice type
	if e, ok := err.(errors.ValidationErrors); ok {
		log.Println(err)
		res := errors.BadRequestResponse(e.Error())
		return ctx.Status(res.Status).JSON(res)
	}

	if e, ok := err.(errors.Error); ok {
		// we first log the error
		log.Println(e)

		if errors.ErrorCode(e) == errors.EINTERNAL {
			res := errors.InternalServerError(e.Error())
			return ctx.Status(res.Status).JSON(res)
		} else if _, ok := e.Err.(errors.Unauthorized); ok {
			res := errors.UnauthorizedResponse(e.Error())
			return ctx.Status(res.Status).JSON(res)
		} else {
			res := errors.BadRequestResponse(e.Error())
			return ctx.Status(res.Status).JSON(res)
		}
	}

	// if its a fiber error we send back the status code and empty response
	if e, ok := err.(*fiber.Error); ok {
		ctx.Status(e.Code)
		return nil
	}

	// will catch any other error we dont process here and return status 500
	if err != nil {
		log.Println(err)
		msg := "Something has happened. Report Issue."
		res := errors.InternalServerError(msg)
		return ctx.Status(res.Status).JSON(res)
	}

	// Return from handler
	return nil
}
