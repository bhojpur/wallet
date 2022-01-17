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

	"github.com/bhojpur/wallet/pkg/admin"
	"github.com/bhojpur/wallet/pkg/agent"
	"github.com/bhojpur/wallet/pkg/auth"
	"github.com/bhojpur/wallet/pkg/config"
	"github.com/bhojpur/wallet/pkg/models"
	"github.com/bhojpur/wallet/pkg/routing/responses"
	"github.com/bhojpur/wallet/pkg/tariff"

	"github.com/gofiber/fiber/v2"
)

func AuthenticateAdmin(adminDomain admin.Interactor, config config.Config) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var params admin.LoginParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return err
		}

		// authenticate by email.
		adm, err := adminDomain.AuthenticateByEmail(params.Email, params.Password)

		// if there is an error authenticating admin.
		if err != nil {
			return err
		}

		// generate an auth token string
		token, err := auth.GetTokenString(adm.ID, models.UserTypAdmin, config.Secret)
		if err != nil {
			return err
		}

		signedUser := models.SignedUser{
			UserID:   adm.ID.String(),
			UserType: models.UserTypAdmin,
			Token:    token,
		}
		_ = ctx.Status(http.StatusOK).JSON(signedUser)

		return nil
	}
}

func RegisterAdmin(adminDomain admin.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		var params admin.RegistrationParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return err
		}

		// register admin
		adm, err := adminDomain.Register(params)
		if err != nil {
			return err
		}

		// we use a presenter to reformat the response of admin.
		_ = ctx.Status(http.StatusOK).JSON(responses.RegistrationResponse(adm.ID, models.UserTypAdmin))

		return nil
	}
}

func AssignFloat(adminDomain admin.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		var params admin.AssignFloatParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return err
		}

		balance, err := adminDomain.AssignFloat(params)
		if err != nil {
			return err
		}

		// we use a presenter to reformat the response of admin.
		_ = ctx.Status(http.StatusOK).JSON(responses.SuccessResponse{
			Status:  "success",
			Message: "Float has been assigned.",
			Data: map[string]interface{}{
				"balance": balance,
			},
		})

		return nil
	}
}

func UpdateCharge(manager tariff.Manager) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		var params admin.UpdateChargeParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return err
		}

		err = manager.UpdateCharge(params.ChargeID, params.Amount)
		if err != nil {
			return err
		}

		_ = ctx.Status(http.StatusOK).JSON(responses.SuccessResponse{
			Status:  "success",
			Message: "charge updated",
		})

		return nil
	}
}

func GetTariff(manager tariff.Manager) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		charges, err := manager.GetTariff()
		if err != nil {
			return err
		}

		_ = ctx.Status(http.StatusOK).JSON(responses.TariffResponse(charges))

		return nil
	}
}

func UpdateSuperAgentStatus(agentDomain agent.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		var params agent.MakeSuperAgentParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return err
		}

		err = agentDomain.UpdateSuperAgentStatus(params.Email)
		if err != nil {
			return err
		}

		_ = ctx.Status(http.StatusOK).JSON(responses.SuccessResponse{
			Status:  "success",
			Message: "Super Agent Status updated",
		})

		return nil
	}
}
