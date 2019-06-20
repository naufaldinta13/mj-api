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

func TestItemCategory_Save(t *testing.T) {
	var m model.ItemCategory
	faker.Fill(&m, "ID")

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

func TestItemCategory_Delete(t *testing.T) {
	m := model.DummyItemCategory()

	e := m.Delete()
	assert.NoError(t, e)
	assert.Zero(t, m.ID)

	mn := new(model.ItemCategory)
	mn.ID = 90000
	e = mn.Delete()
	assert.Error(t, e)

	mn = new(model.ItemCategory)
	e = mn.Delete()
	assert.Error(t, e)
}

func TestItemCategory_Read(t *testing.T) {
	var m model.ItemCategory

	mn := model.DummyItemCategory()
	m.ID = mn.ID
	m.Read()

	assert.NotZero(t, m.ID)
}

func TestItemCategory_MarshalJSON(t *testing.T) {
	mn := model.DummyItemCategory()

	j, e := mn.MarshalJSON()
	assert.NoError(t, e)
	assert.Contains(t, string(j), common.Encrypt(mn.ID))
}
