// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/cuxs"
	"git.qasico.com/cuxs/orm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestGetApplicationModule(t *testing.T) {
	_, e := GetApplicationModule("id", 1000)
	assert.Error(t, e, "Response should be error, because there are no data yet.")

	c := model.DummyApplicationModule()
	cd, e := GetApplicationModule("id", c.ID)
	assert.NoError(t, e, "Data should be exists.")
	assert.Equal(t, c.ID, cd.ID, "ID Response should be a same.")
}

func TestCheckAuthPrivilege(t *testing.T) {
	user := model.DummyUser()
	module := model.DummyApplicationModule()
	module.IsActive = 1
	module.Save()
	//dummy priviledge
	priviledge := model.DummyApplicationPrivilege()
	priviledge.Usergroup = user.Usergroup
	priviledge.ApplicationModule = module
	priviledge.Save()

	// jika data ada
	e := CheckAuthPrivilege(user, module)
	assert.NoError(t, e)

	// jika modul gak aktif
	module1 := model.DummyApplicationModule()
	module1.IsActive = 0
	module1.Save()
	priviledge.Usergroup = user.Usergroup
	priviledge.ApplicationModule = module1
	priviledge.Save()
	e = CheckAuthPrivilege(user, module1)
	assert.Error(t, e)

	// jika data gak ada
	modulRand := model.DummyApplicationModule()
	e = CheckAuthPrivilege(user, modulRand)
	assert.Error(t, e)
}

func LoginTokenTest(token string) (*SessionData, error) {
	validAuth := "Bearer " + token

	req, _ := http.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, validAuth)
	res := httptest.NewRecorder()

	e := cuxs.New()
	c := e.NewContext(req, res)
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	h := cuxs.Authorized()(handler)
	if h(c) == nil {
		ctx := cuxs.NewContext(c)

		sd, e := UserSession(ctx)

		return sd, e
	}

	return nil, errors.New("Tidak dapat melakukan login")
}

func DbClean(table ...string) {
	orm := orm.NewOrm()
	for _, t := range table {
		_, e := orm.Raw(fmt.Sprintf("Delete From %s where id > ?", t), 0).Exec()
		if e != nil {
			panic(e)
		}
		orm.Raw(fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = 1;", t)).Exec()
	}
}

func TestSessionAsSysAdmin(t *testing.T) {
	DbClean("application_menu")
	// user dengan usergroup 1
	user := model.DummyUserPriviledgeWithUsergroup(1)
	for a := 0; a < 3; a++ {
		modul := model.DummyApplicationModule()

		privilege := model.DummyApplicationPrivilege()
		privilege.ApplicationModule = modul
		privilege.Usergroup = user.Usergroup
		privilege.Save("ApplicationModule", "Usergroup")

		menu := model.DummyApplicationMenu()
		menu.ApplicationModule = modul
		menu.IsActive = 1
		menu.Save("ApplicationModule", "IsActive")
	}

	// apabila dia login
	// mengecek data application_menu
	sd, e := StartSession(user.ID)
	assert.NoError(t, e, "Seharusnya start session berhasil")
	assert.NotEmpty(t, sd.Token, "Seharusnya token ada")
	assert.Equal(t, 3, len(sd.ApplicationMenu), "Seharusnya data application menu ada 3")
	assert.Equal(t, int64(1), sd.User.Usergroup.ID, "Seharusnya usergroup nya adalah 1")

	// mencoba untuk login menggunakan token
	// dan mendapatkan user session dari context
	r, e := LoginTokenTest(sd.Token)
	assert.NoError(t, e)
	assert.NotEmpty(t, r.Token, "Seharusnya token ada")
	assert.Equal(t, 3, len(r.ApplicationMenu), "Seharusnya data application menu ada 3")
	assert.Equal(t, int64(1), r.User.Usergroup.ID, "Seharusnya usergroup nya adalah 1")
}

func TestSessionAsOwner(t *testing.T) {
	DbClean("application_menu")
	// user dengan usergroup 1
	user := model.DummyUserPriviledgeWithUsergroup(2)
	for a := 0; a < 3; a++ {
		modul := model.DummyApplicationModule()

		privilege := model.DummyApplicationPrivilege()
		privilege.ApplicationModule = modul
		privilege.Usergroup = user.Usergroup
		privilege.Save("ApplicationModule", "Usergroup")

		menu := model.DummyApplicationMenu()
		menu.ApplicationModule = modul
		menu.IsActive = 1
		menu.Save("ApplicationModule", "IsActive")
	}

	// apabila dia login
	// mengecek data application_menu
	sd, e := StartSession(user.ID)
	assert.NoError(t, e, "Seharusnya start session berhasil")
	assert.NotEmpty(t, sd.Token, "Seharusnya token ada")
	assert.Equal(t, 3, len(sd.ApplicationMenu), "Seharusnya data application menu ada 3")
	assert.Equal(t, int64(2), sd.User.Usergroup.ID, "Seharusnya usergroup nya adalah 2")

	// mencoba untuk login menggunakan token
	// dan mendapatkan user session dari context
	r, e := LoginTokenTest(sd.Token)
	assert.NoError(t, e)
	assert.NotEmpty(t, r.Token, "Seharusnya token ada")
	assert.Equal(t, 3, len(r.ApplicationMenu), "Seharusnya data application menu ada 3")
	assert.Equal(t, int64(2), r.User.Usergroup.ID, "Seharusnya usergroup nya adalah 2")
}

func TestSessionAsSupervisor(t *testing.T) {
	DbClean("application_menu")
	// user dengan usergroup 1
	user := model.DummyUserPriviledgeWithUsergroup(3)
	for a := 0; a < 3; a++ {
		modul := model.DummyApplicationModule()

		privilege := model.DummyApplicationPrivilege()
		privilege.ApplicationModule = modul
		privilege.Usergroup = user.Usergroup
		privilege.Save("ApplicationModule", "Usergroup")

		menu := model.DummyApplicationMenu()
		menu.ApplicationModule = modul
		menu.IsActive = 1
		menu.Save("ApplicationModule", "IsActive")
	}

	// apabila dia login
	// mengecek data application_menu
	sd, e := StartSession(user.ID)
	assert.NoError(t, e, "Seharusnya start session berhasil")
	assert.NotEmpty(t, sd.Token, "Seharusnya token ada")
	assert.Equal(t, 3, len(sd.ApplicationMenu), "Seharusnya data application menu ada 3")
	assert.Equal(t, int64(3), sd.User.Usergroup.ID, "Seharusnya usergroup nya adalah 3")

	// mencoba untuk login menggunakan token
	// dan mendapatkan user session dari context
	r, e := LoginTokenTest(sd.Token)
	assert.NoError(t, e)
	assert.NotEmpty(t, r.Token, "Seharusnya token ada")
	assert.Equal(t, 3, len(r.ApplicationMenu), "Seharusnya data application menu ada 3")
	assert.Equal(t, int64(3), r.User.Usergroup.ID, "Seharusnya usergroup nya adalah 3")
}

func TestSessionAsCashier(t *testing.T) {
	DbClean("application_menu")
	// user dengan usergroup 1
	user := model.DummyUserPriviledgeWithUsergroup(4)
	for a := 0; a < 3; a++ {
		modul := model.DummyApplicationModule()

		privilege := model.DummyApplicationPrivilege()
		privilege.ApplicationModule = modul
		privilege.Usergroup = user.Usergroup
		privilege.Save("ApplicationModule", "Usergroup")

		menu := model.DummyApplicationMenu()
		menu.ApplicationModule = modul
		menu.IsActive = 1
		menu.Save("ApplicationModule", "IsActive")
	}

	// apabila dia login
	// mengecek data application_menu
	sd, e := StartSession(user.ID)
	assert.NoError(t, e, "Seharusnya start session berhasil")
	assert.NotEmpty(t, sd.Token, "Seharusnya token ada")
	assert.Equal(t, 3, len(sd.ApplicationMenu), "Seharusnya data application menu ada 3")
	assert.Equal(t, int64(4), sd.User.Usergroup.ID, "Seharusnya usergroup nya adalah 4")

	// mencoba untuk login menggunakan token
	// dan mendapatkan user session dari context
	r, e := LoginTokenTest(sd.Token)
	assert.NoError(t, e)
	assert.NotEmpty(t, r.Token, "Seharusnya token ada")
	assert.Equal(t, 3, len(r.ApplicationMenu), "Seharusnya data application menu ada 3")
	assert.Equal(t, int64(4), r.User.Usergroup.ID, "Seharusnya usergroup nya adalah 4")
}

func TestLogin(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	for a := 0; a < 3; a++ {
		modul := model.DummyApplicationModule()

		privilege := model.DummyApplicationPrivilege()
		privilege.ApplicationModule = modul
		privilege.Usergroup = user.Usergroup
		privilege.Save("ApplicationModule", "Usergroup")

		menu := model.DummyApplicationMenu()
		menu.ApplicationModule = modul
		menu.IsActive = 1
		menu.Save("ApplicationModule", "IsActive")
	}

	// mengecek lastlogin
	sd, e := Login(user)
	assert.NoError(t, e)
	assert.NotEqual(t, user.LastLogin, sd.User.LastLogin)
}

func TestGetApplicationModuleIDByPrivilegeUser(t *testing.T) {

	user := model.DummyUser()
	user.Save()

	pri := model.DummyApplicationPrivilege()
	pri.Usergroup = user.Usergroup
	pri.Save()

	pri2 := model.DummyApplicationPrivilege()
	pri2.Usergroup = user.Usergroup
	pri2.Save()

	_, e := GetApplicationModuleIDByPrivilegeUser(user)
	assert.NoError(t, e, "seharusnya tidak error")
}

func TestGetApplicationMenuByUser(t *testing.T) {

	o := orm.NewOrm()
	o.Raw("delete from application_menu")
	var rq orm.RequestQuery

	user := model.DummyUser()
	user.IsActive = 1
	user.Save()

	module := model.DummyApplicationModule()
	module.IsActive = 1
	module.Save()

	module2 := model.DummyApplicationModule()
	module2.IsActive = 1
	module2.Save()

	pri := model.DummyApplicationPrivilege()
	pri.ApplicationModule = module
	pri.Usergroup = user.Usergroup
	pri.Save()

	pri2 := model.DummyApplicationPrivilege()
	pri2.ApplicationModule = module2
	pri2.Usergroup = user.Usergroup
	pri2.Save()

	menu := model.DummyApplicationMenu()
	menu.ApplicationModule = module
	menu.IsActive = 1
	menu.Save()

	menu2 := model.DummyApplicationMenu()
	menu2.ApplicationModule = module2
	menu2.IsActive = 1
	menu2.Save()

	m, e := GetApplicationMenuByUser(rq, user)
	assert.NoError(t, e, "seharusnya tidak error")
	assert.Equal(t, 2, len(m))
}

func TestGetUsername(t *testing.T) {

	user := model.DummyUser()
	user.Save()

	m, e := GetUsername(user.Username)
	assert.NoError(t, e, "seharusnya tidak error")
	assert.Equal(t, user.Username, m.Username)
}
