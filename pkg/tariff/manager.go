package tariff

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

	"github.com/gofrs/uuid"
)

type Manager interface {
	GetCharge(operation models.TxnOperation, src models.UserType, dest models.UserType) (models.Paisas, error)
	GetTariff() ([]Charge, error)
	UpdateCharge(chargeID uuid.UUID, fee models.Paisas) error
}

func NewManager(repository Repository) Manager {
	mgr := &manager{repository}

	go mgr.initTariffSetup()

	return mgr
}

type manager struct {
	repository Repository
}

func (mg manager) GetCharge(operation models.TxnOperation, src models.UserType, dest models.UserType) (models.Paisas, error) {
	tariff, err := mg.repository.Get(operation, src, dest)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return models.Paisas(0), errors.Error{Err: err, Message: errors.ErrTariffNotSet}
	} else if err != nil {
		return models.Paisas(0), err
	}

	return tariff.Fee, nil
}

func (mg manager) GetTariff() ([]Charge, error) {
	charges, err := mg.repository.FetchAll()
	if err != nil {
		return nil, err
	}

	if len(charges) == 0 {
		return nil, errors.Error{Code: errors.ENOTFOUND, Message: errors.ErrTariffNotSet}
	}

	return charges, nil
}

func (mg manager) UpdateCharge(chargeID uuid.UUID, fee models.Paisas) error {
	charge, err := mg.repository.FindByID(chargeID)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return errors.Error{Err: err, Message: errors.ErrChargeNotFound}
	} else if err != nil {
		return err
	}

	charge.Fee = fee
	err = mg.repository.Update(charge)
	if err != nil {
		return err
	}

	return nil
}

// initializes a tariff with zero amount, is used only once during initial setup of charges
func (mg manager) addCharge(txOperation models.TxnOperation, source models.UserType, dest models.UserType) error {

	if ok := models.IsValidTxnOperation(txOperation); !ok {
		return errors.Error{Code: errors.EINVALID, Message: errors.ErrInvalidOperation}
	}

	_, err := mg.repository.Add(Charge{
		Transaction:         txOperation,
		SourceUserType:      source,
		DestinationUserType: dest,
		Fee:                 models.Paisas(0),
	})
	if err != nil {
		return err
	}

	return nil
}

// is a definition of all valid withdrawals, remember withdrawals can only happen
// at an agent's desk
func (mg manager) validWithdrawTx() []ValidTransaction {
	return []ValidTransaction{
		{models.UserTypSubscriber, models.UserTypAgent},
		{models.UserTypMerchant, models.UserTypAgent},
		{models.UserTypAgent, models.UserTypAgent},
	}
}

// is a definition of all valid transfers, remember only agents are allowed
// to transfer to other agents
func (mg manager) validTransferTx() []ValidTransaction {
	return []ValidTransaction{
		{models.UserTypAgent, models.UserTypAgent},           // transfer between an agent to an agent
		{models.UserTypSubscriber, models.UserTypSubscriber}, // transfer between a subscriber to subscriber
		{models.UserTypMerchant, models.UserTypSubscriber},   // transfer between a merchant to subscriber
		{models.UserTypSubscriber, models.UserTypMerchant},   // transfer between a subscriber to merchant -> PAYMENT
		{models.UserTypAgent, models.UserTypMerchant},        // transfer between an agent to merchant -> PAYMENT
	}
}

func (mg manager) initTariffSetup() error {
	// add valid withdraw transactions between customers
	for _, validTx := range mg.validWithdrawTx() {
		err := mg.addCharge(models.TxnOpWithdraw, validTx[0], validTx[1])
		if err != nil {
			return err
		}
	}

	// add valid transfer transactions between customers
	for _, validTx := range mg.validTransferTx() {
		err := mg.addCharge(models.TxnOpTransfer, validTx[0], validTx[1])
		if err != nil {
			return err
		}
	}

	return nil
}
