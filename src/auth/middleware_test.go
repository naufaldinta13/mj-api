// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"testing"

	"git.qasico.com/mj/api/datastore/model"

	"github.com/stretchr/testify/assert"
)

func TestSysadmin(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	ctx, h := LoginAs(user)
	//kalo ada data nya
	d := CheckPrivilege("dashboard")(h)
	assert.NoError(t, d(ctx), "should be authorized")

	//jika module is active =0
	modul := model.DummyApplicationModule()
	modul.IsActive = 0
	modul.Save()
	d = CheckPrivilege(modul.Alias)(h)
	assert.Error(t, d(ctx), "should be Application Module Not Exist")
}

func TestOwner(t *testing.T) {
	owner := model.DummyUserPriviledgeWithUsergroup(2)
	ctx, h := LoginAs(owner)

	d := CheckPrivilege("dashboard")(h)
	assert.NoError(t, d(ctx), "should be authorized")

	d = CheckPrivilege("user_create")(h)
	assert.Error(t, d(ctx), "should be unauthorized")
}

func TestSupervisor(t *testing.T) {
	supervisor := model.DummyUserPriviledgeWithUsergroup(3)
	ctx, h := LoginAs(supervisor)

	d := CheckPrivilege("dashboard")(h)
	assert.NoError(t, d(ctx), "should be authorized")

	d = CheckPrivilege("user_create")(h)
	assert.Error(t, d(ctx), "should be unauthorized")
}

func TestCashier(t *testing.T) {
	cashier := model.DummyUserPriviledgeWithUsergroup(4)
	ctx, h := LoginAs(cashier)

	d := CheckPrivilege("dashboard")(h)
	assert.NoError(t, d(ctx), "should be authorized")

	d = CheckPrivilege("user_create")(h)
	assert.Error(t, d(ctx), "should be unauthorized")
}
