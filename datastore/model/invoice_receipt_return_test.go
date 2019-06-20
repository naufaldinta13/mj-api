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

func TestInvoiceReceiptReturn_Save(t *testing.T) {
	var m model.InvoiceReceiptReturn
	faker.Fill(&m, "ID")

	m.InvoiceReceipt = model.DummyInvoiceReceipt()

	m.SalesReturn = model.DummySalesReturn()

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

func TestInvoiceReceiptReturn_Delete(t *testing.T) {
	m := model.DummyInvoiceReceiptReturn()

	e := m.Delete()
	assert.NoError(t, e)
	assert.Zero(t, m.ID)

	mn := new(model.InvoiceReceiptReturn)
	mn.ID = 90000
	e = mn.Delete()
	assert.Error(t, e)

	mn = new(model.InvoiceReceiptReturn)
	e = mn.Delete()
	assert.Error(t, e)
}

func TestInvoiceReceiptReturn_Read(t *testing.T) {
	var m model.InvoiceReceiptReturn

	mn := model.DummyInvoiceReceiptReturn()
	m.ID = mn.ID
	m.Read()

	assert.NotZero(t, m.ID)
}

func TestInvoiceReceiptReturn_MarshalJSON(t *testing.T) {
	mn := model.DummyInvoiceReceiptReturn()

	j, e := mn.MarshalJSON()
	assert.NoError(t, e)
	assert.Contains(t, string(j), common.Encrypt(mn.ID))
}
