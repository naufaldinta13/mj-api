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

func TestInvoiceReceiptItem_Save(t *testing.T) {
	var m model.InvoiceReceiptItem
	faker.Fill(&m, "ID")

	m.InvoiceReceipt = model.DummyInvoiceReceipt()

	m.SalesInvoice = model.DummySalesInvoice()

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

func TestInvoiceReceiptItem_Delete(t *testing.T) {
	m := model.DummyInvoiceReceiptItem()

	e := m.Delete()
	assert.NoError(t, e)
	assert.Zero(t, m.ID)

	mn := new(model.InvoiceReceiptItem)
	mn.ID = 90000
	e = mn.Delete()
	assert.Error(t, e)

	mn = new(model.InvoiceReceiptItem)
	e = mn.Delete()
	assert.Error(t, e)
}

func TestInvoiceReceiptItem_Read(t *testing.T) {
	var m model.InvoiceReceiptItem

	mn := model.DummyInvoiceReceiptItem()
	m.ID = mn.ID
	m.Read()

	assert.NotZero(t, m.ID)
}

func TestInvoiceReceiptItem_MarshalJSON(t *testing.T) {
	mn := model.DummyInvoiceReceiptItem()

	j, e := mn.MarshalJSON()
	assert.NoError(t, e)
	assert.Contains(t, string(j), common.Encrypt(mn.ID))
}
