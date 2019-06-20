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

func TestItemVariant_Save(t *testing.T) {
	var m model.ItemVariant
	faker.Fill(&m, "ID")

	m.Item = model.DummyItem()

	m.Measurement = model.DummyMeasurement()

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

func TestItemVariant_Delete(t *testing.T) {
	m := model.DummyItemVariant()

	e := m.Delete()
	assert.NoError(t, e)
	assert.Zero(t, m.ID)

	mn := new(model.ItemVariant)
	mn.ID = 90000
	e = mn.Delete()
	assert.Error(t, e)

	mn = new(model.ItemVariant)
	e = mn.Delete()
	assert.Error(t, e)
}

func TestItemVariant_Read(t *testing.T) {
	var m model.ItemVariant

	mn := model.DummyItemVariant()
	m.ID = mn.ID
	m.Read()

	assert.NotZero(t, m.ID)
}

func TestItemVariant_MarshalJSON(t *testing.T) {
	mn := model.DummyItemVariant()

	j, e := mn.MarshalJSON()
	assert.NoError(t, e)
	assert.Contains(t, string(j), common.Encrypt(mn.ID))
}

func TestItemVariant_Archive(t *testing.T) {
	mn := model.DummyItemVariant()
	mn.IsArchived = 0
	mn.Save()

	mn.Archive()
	assert.Equal(t, int8(1), mn.IsArchived)

	// item nya terarchive
	m := model.DummyItem()
	m.IsArchived = 0
	m.Save()

	mv := model.DummyItemVariant()
	mv.Item = m
	mv.IsArchived = 0
	mv.Save()

	mvd := model.DummyItemVariant()
	mvd.Item = m
	mvd.IsArchived = 0
	mvd.Save()

	e := mv.Archive()
	assert.NoError(t, e, "seharusnya success mengarchive variant")

	m.Read()
	assert.Equal(t, int8(0), m.IsArchived, "item seharusnya tidak terarchive")

	e = mvd.Archive()
	assert.NoError(t, e, "seharusnya success mengarchive variant")

	m.Read()
	assert.Equal(t, int8(1), m.IsArchived, "item seharusnya terarchive")
}

func TestItemVariant_Unarchive(t *testing.T) {
	mn := model.DummyItemVariant()
	mn.IsArchived = 1
	mn.Save()

	mn.Unarchive()
	assert.Equal(t, int8(0), mn.IsArchived)

	// item nya terarchive
	m := model.DummyItem()
	m.IsArchived = 1
	m.Save()

	mv := model.DummyItemVariant()
	mv.Item = m
	mv.IsArchived = 1
	mv.Save()

	mvd := model.DummyItemVariant()
	mvd.Item = m
	mvd.IsArchived = 1
	mvd.Save()

	assert.Equal(t, int8(1), m.IsArchived, "item seharusnya archived")

	e := mv.Unarchive()
	assert.NoError(t, e, "seharusnya success mengunarchive variant")

	m.Read()
	assert.Equal(t, int8(0), m.IsArchived, "item seharusnya ter-unarchive")

	e = mvd.Unarchive()
	assert.NoError(t, e, "seharusnya success mengunarchive variant")

	m.Read()
	assert.Equal(t, int8(0), m.IsArchived, "item seharusnya ter unarchive")
}

func TestItemVariant_DeleteItemVariant(t *testing.T) {

	mi := model.DummyItem()
	mi.IsDeleted = 0
	mi.Save()

	mn := model.DummyItemVariant()
	mn.Item = mi
	mn.IsArchived = 1
	mn.Save()

	mn.DeleteItemVariant()
	assert.Equal(t, int8(1), mn.IsDeleted)
	assert.Equal(t, int8(0), mi.IsDeleted)
}
