package user_handlers

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
	"github.com/bhojpur/wallet/pkg/models"
	"github.com/bhojpur/wallet/pkg/registry"

	"github.com/gofiber/fiber/v2"
)

func Authenticate(domain *registry.Domain, config config.Config) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		// get the user type authenticating
		userType := ctx.Params("user_type")

		switch models.UserType(userType) {
		case models.UserTypAdmin:
			return AuthenticateAdmin(domain.Admin, config)(ctx)
		case models.UserTypAgent:
			return AuthenticateAgent(domain.Agent, config)(ctx)
		case models.UserTypMerchant:
			return AuthenticateMerchant(domain.Merchant, config)(ctx)
		case models.UserTypSubscriber:
			return AuthenticateSubscriber(domain.Subscriber, config)(ctx)
		default:
			return fiber.ErrNotFound
		}
	}
}

func Register(domain *registry.Domain) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		// get the user type authenticating
		userType := ctx.Params("user_type")

		switch models.UserType(userType) {
		case models.UserTypAdmin:
			return RegisterAdmin(domain.Admin)(ctx)
		case models.UserTypAgent:
			return RegisterAgent(domain.Agent)(ctx)
		case models.UserTypMerchant:
			return RegisterMerchant(domain.Merchant)(ctx)
		case models.UserTypSubscriber:
			return RegisterSubscriber(domain.Subscriber)(ctx)
		default:
			return fiber.ErrNotFound
		}
	}
}
