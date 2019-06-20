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

func TestUser_Save(t *testing.T) {
	var m model.User
	faker.Fill(&m, "ID")

	m.Usergroup = model.DummyUsergroup()

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

func TestUser_Delete(t *testing.T) {
	m := model.DummyUser()

	e := m.Delete()
	assert.NoError(t, e)
	assert.Zero(t, m.ID)

	mn := new(model.User)
	mn.ID = 90000
	e = mn.Delete()
	assert.Error(t, e)

	mn = new(model.User)
	e = mn.Delete()
	assert.Error(t, e)
}

func TestUser_Read(t *testing.T) {
	var m model.User

	mn := model.DummyUser()
	m.ID = mn.ID
	m.Read()

	assert.NotZero(t, m.ID)
}

func TestUser_MarshalJSON(t *testing.T) {
	mn := model.DummyUser()

	j, e := mn.MarshalJSON()
	assert.NoError(t, e)
	assert.Contains(t, string(j), common.Encrypt(mn.ID))
}

func TestDummyUserPriviledgeWithUsergroup(t *testing.T) {

	m := model.DummyUserPriviledgeWithUsergroup(1)

	var mm = &model.User{
		ID: m.ID,
	}

	mm.Read("ID")
	assert.Equal(t, m.ID, mm.ID)
	assert.Equal(t, m.Username, mm.Username)
}
