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

func TestSalesOrderItem_Save(t *testing.T) {
	var m model.SalesOrderItem
	faker.Fill(&m, "ID")

	m.SalesOrder = model.DummySalesOrder()

	m.ItemVariant = model.DummyItemVariant()

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

func TestSalesOrderItem_Delete(t *testing.T) {
	m := model.DummySalesOrderItem()

	e := m.Delete()
	assert.NoError(t, e)
	assert.Zero(t, m.ID)

	mn := new(model.SalesOrderItem)
	mn.ID = 90000
	e = mn.Delete()
	assert.Error(t, e)

	mn = new(model.SalesOrderItem)
	e = mn.Delete()
	assert.Error(t, e)
}

func TestSalesOrderItem_Read(t *testing.T) {
	var m model.SalesOrderItem

	mn := model.DummySalesOrderItem()
	m.ID = mn.ID
	m.Read()

	assert.NotZero(t, m.ID)
}

func TestSalesOrderItem_MarshalJSON(t *testing.T) {
	mn := model.DummySalesOrderItem()

	j, e := mn.MarshalJSON()
	assert.NoError(t, e)
	assert.Contains(t, string(j), common.Encrypt(mn.ID))
}
