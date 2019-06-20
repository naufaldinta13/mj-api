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

func TestWorkorderReceiving_Save(t *testing.T) {
	var m model.WorkorderReceiving
	faker.Fill(&m, "ID")

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

func TestWorkorderReceiving_Delete(t *testing.T) {
	m := model.DummyWorkorderReceiving()

	e := m.Delete()
	assert.NoError(t, e)
	assert.Zero(t, m.ID)

	mn := new(model.WorkorderReceiving)
	mn.ID = 90000
	e = mn.Delete()
	assert.Error(t, e)

	mn = new(model.WorkorderReceiving)
	e = mn.Delete()
	assert.Error(t, e)
}

func TestWorkorderReceiving_Read(t *testing.T) {
	var m model.WorkorderReceiving

	mn := model.DummyWorkorderReceiving()
	m.ID = mn.ID
	m.Read()

	assert.NotZero(t, m.ID)
}

func TestWorkorderReceiving_MarshalJSON(t *testing.T) {
	mn := model.DummyWorkorderReceiving()

	j, e := mn.MarshalJSON()
	assert.NoError(t, e)
	assert.Contains(t, string(j), common.Encrypt(mn.ID))
}
