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

func TestApplicationModule_Save(t *testing.T) {
	var m model.ApplicationModule
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

func TestApplicationModule_Delete(t *testing.T) {
	m := model.DummyApplicationModule()

	e := m.Delete()
	assert.NoError(t, e)
	assert.Zero(t, m.ID)

	mn := new(model.ApplicationModule)
	mn.ID = 90000
	e = mn.Delete()
	assert.Error(t, e)

	mn = new(model.ApplicationModule)
	e = mn.Delete()
	assert.Error(t, e)
}

func TestApplicationModule_Read(t *testing.T) {
	var m model.ApplicationModule

	mn := model.DummyApplicationModule()
	m.ID = mn.ID
	m.Read()

	assert.NotZero(t, m.ID)
}

func TestApplicationModule_MarshalJSON(t *testing.T) {
	mn := model.DummyApplicationModule()

	j, e := mn.MarshalJSON()
	assert.NoError(t, e)
	assert.Contains(t, string(j), common.Encrypt(mn.ID))
}
