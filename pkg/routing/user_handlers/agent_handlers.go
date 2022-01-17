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
	"net/http"

	"github.com/bhojpur/wallet/pkg/agent"
	"github.com/bhojpur/wallet/pkg/auth"
	"github.com/bhojpur/wallet/pkg/config"
	"github.com/bhojpur/wallet/pkg/models"
	"github.com/bhojpur/wallet/pkg/routing/responses"

	"github.com/gofiber/fiber/v2"
)

func AuthenticateAgent(agentDomain agent.Interactor, config config.Config) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var params agent.LoginParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(err)
		}

		// authenticate by email.
		agt, err := agentDomain.AuthenticateByEmail(params.Email, params.Password)

		// if there is an error authenticating agent.
		if err != nil {
			return err
		}

		// check if agent is a super agent
		agentType := models.UserTypAgent
		if agt.IsSuperAgent() {
			agentType = models.UserTypSuperAgent
		}

		// generate an auth token string
		token, err := auth.GetTokenString(agt.ID, agentType, config.Secret)
		if err != nil {
			return err
		}

		signedUser := models.SignedUser{
			UserID:   agt.ID.String(),
			UserType: agentType,
			Token:    token,
		}
		_ = ctx.Status(http.StatusOK).JSON(signedUser)

		return nil
	}
}

func RegisterAgent(agentDomain agent.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		var params agent.RegistrationParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(err)
		}

		// register agent
		agt, err := agentDomain.Register(params)
		if err != nil {
			return err
		}

		// we use a presenter to reformat the response of agent.
		_ = ctx.Status(http.StatusOK).JSON(responses.RegistrationResponse(agt.ID, models.UserTypAgent))

		return nil
	}
}
