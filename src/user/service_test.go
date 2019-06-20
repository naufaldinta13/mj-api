// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"strconv"
	"testing"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"git.qasico.com/mj/api/src/auth"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	user := model.DummyUser()
	user.IsActive = 1
	user.Save("is_active")

	m, e := GetUserByID(user.ID)

	assert.NoError(t, e)
	assert.Equal(t, user.ID, m.ID)
	assert.Equal(t, user.Password, m.Password)
	assert.Equal(t, user.Username, m.Username)
	assert.Equal(t, user.Usergroup, m.Usergroup)
	assert.Equal(t, user.IsActive, m.IsActive)
	assert.Equal(t, user.CreatedBy, m.CreatedBy)
	assert.Equal(t, user.FullName, m.FullName)
	assert.Equal(t, user.RememberToken, m.RememberToken)
}

func TestGetUsers(t *testing.T) {
	orm.NewOrm().Raw("DELETE FROM user").Exec()
	model.DummyUser()
	model.DummyUser()
	model.DummyUser()
	model.DummyUser()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)

	qs := orm.RequestQuery{}
	m, _, e := GetUsers(&qs, sd)
	assert.NoError(t, e, "Data should be exists.")

	assert.Equal(t, int(5), len(*m), "Actual"+strconv.Itoa(len(*m)))
}
func TestGetUsersNotUsergroup1(t *testing.T) {
	orm.NewOrm().Raw("DELETE FROM user").Exec()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(2)
	sd, _ := auth.Login(user)

	dug := model.DummyUsergroup()
	dug.UsergroupName = "admin"
	dug.Save()

	dug2 := model.DummyUsergroup()
	dug2.UsergroupName = "owner"
	dug2.Save()

	du := model.DummyUser()
	du.Usergroup = dug2
	du.Save()

	du1 := model.DummyUser()
	du1.Usergroup = dug
	du1.Save()

	qs := orm.RequestQuery{}
	m, _, e := GetUsers(&qs, sd)
	assert.NoError(t, e, "Data should be exists.")

	assert.Equal(t, int(1), len(*m), "Actual"+strconv.Itoa(len(*m)))
}
