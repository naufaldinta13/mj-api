// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_test

import (
	"testing"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/common/faker"
	"github.com/stretchr/testify/assert"
)

func TestPurchaseInvoice_Save(t *testing.T) {
	var m model.PurchaseInvoice
	faker.Fill(&m, "ID")

	m.PurchaseOrder = model.DummyPurchaseOrder()

	m.CreatedBy = model.DummyUser()

	e := m.Save()
	assert.NoError(t, e)
	assert.NotZero(t, m.ID)

	mn := m
	faker.Fill(&mn, "ID")
	e = mn.Save()
	assert.NoError(t, e)

	mn.ID = 999999
	e = mn.Save()
	assert.NoError(t, e)
}

func TestPurchaseInvoice_Delete(t *testing.T) {
	m := model.DummyPurchaseInvoice()

	e := m.Delete()
	assert.NoError(t, e)
	assert.Zero(t, m.ID)

	mn := new(model.PurchaseInvoice)
	mn.ID = 90000
	e = mn.Delete()
	assert.Error(t, e)

	mn = new(model.PurchaseInvoice)
	e = mn.Delete()
	assert.Error(t, e)
}

func TestPurchaseInvoice_Read(t *testing.T) {
	var m model.PurchaseInvoice

	mn := model.DummyPurchaseInvoice()
	m.ID = mn.ID
	m.Read()

	assert.NotZero(t, m.ID)
}

func TestPurchaseInvoice_MarshalJSON(t *testing.T) {
	mn := model.DummyPurchaseInvoice()

	j, e := mn.MarshalJSON()
	assert.NoError(t, e)
	assert.Contains(t, string(j), common.Encrypt(mn.ID))
}