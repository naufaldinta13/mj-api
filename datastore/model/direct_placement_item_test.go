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

func TestDirectPlacementItem_Save(t *testing.T) {
	var m model.DirectPlacementItem
	faker.Fill(&m, "ID")

	m.DirectPlacment = model.DummyDirectPlacement()

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

func TestDirectPlacementItem_Delete(t *testing.T) {
	m := model.DummyDirectPlacementItem()

	e := m.Delete()
	assert.NoError(t, e)
	assert.Zero(t, m.ID)

	mn := new(model.DirectPlacementItem)
	mn.ID = 90000
	e = mn.Delete()
	assert.Error(t, e)

	mn = new(model.DirectPlacementItem)
	e = mn.Delete()
	assert.Error(t, e)
}

func TestDirectPlacementItem_Read(t *testing.T) {
	var m model.DirectPlacementItem

	mn := model.DummyDirectPlacementItem()
	m.ID = mn.ID
	m.Read()

	assert.NotZero(t, m.ID)
}

func TestDirectPlacementItem_MarshalJSON(t *testing.T) {
	mn := model.DummyDirectPlacementItem()

	j, e := mn.MarshalJSON()
	assert.NoError(t, e)
	assert.Contains(t, string(j), common.Encrypt(mn.ID))
}
