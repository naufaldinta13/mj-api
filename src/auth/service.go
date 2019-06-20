// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"errors"
	"strings"
	"time"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/cuxs"
	"git.qasico.com/cuxs/orm"
	"github.com/dgrijalva/jwt-go"
)

// SessionData structur data current user logged in.
type SessionData struct {
	Token             string                    `json:"token"`
	User              *model.User               `json:"user"`
	ApplicationMenu   []model.ApplicationMenu   `json:"app_menu"`
	ApplicationModule []model.ApplicationModule `json:"app_module"`
}

// Login mendapatkan session data dari model user
// user diasumsikan sudah valid untuk login.
// jadi disini tidak ada validasi untuk login, hanya
// untuk mendapatkan session data
func Login(user *model.User) (sd *SessionData, e error) {
	//application menu
	if sd, e = StartSession(user.ID); e == nil {
		// update last login dari user tersebut
		user.LastLogin = time.Now()
		user.Save("LastLogin")
		return sd, nil
	}
	return nil, e
}

// StartSession mendapatkan data user entity dengan token
// untuk menandakan session user yang sedang login.
func StartSession(userID int64, token ...string) (sd *SessionData, e error) {
	sd = new(SessionData)
	var rq orm.RequestQuery
	var menu []model.ApplicationMenu
	var module []model.ApplicationModule

	// buat token baru atau menggunakan yang sebelumnya
	if len(token) == 0 {
		sd.Token = cuxs.JwtToken("id", userID)
	} else {
		sd.Token = token[0]
	}

	// membaca data user terlebih dahulu untuk
	sd.User = &model.User{ID: userID}
	sd.User.Read()

	// get data apllication menu
	if menu, e = GetApplicationMenuByUser(rq, sd.User); e == nil {
		if module, e = GetApplicationModuleByUser(rq, sd.User); e == nil {
			sd.ApplicationMenu = menu
			sd.ApplicationModule = module

			return sd, nil
		}
	}

	return nil, e
}

// UserSession mendapatkan session data dari user yang mengirimkan request.
func UserSession(ctx *cuxs.Context) (*SessionData, error) {
	if u := ctx.Get("user"); u != nil {
		c := u.(*jwt.Token).Claims.(jwt.MapClaims)
		var userID int64

		// id adalah user id
		if c["id"] != nil {
			userID = int64(c["id"].(float64))
		}

		// memakai token sebelumnya
		token := ctx.Get("user").(*jwt.Token).Raw

		return StartSession(userID, token)
	}

	return nil, errors.New("Invalid jwt token")
}

// GetApplicationModule find a single data application_module using field and value condition.
func GetApplicationModule(field string, values ...interface{}) (*model.ApplicationModule, error) {
	m := new(model.ApplicationModule)
	o := orm.NewOrm().QueryTable(m)
	if err := o.Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CheckAuthPrivilege get data from user group module with field application_module_id, and usergroupid
func CheckAuthPrivilege(user *model.User, module *model.ApplicationModule) error {
	m := new(model.ApplicationPrivilege)
	o := orm.NewOrm().QueryTable(m)

	return o.Filter("application_module_id", module.ID).Filter("application_module_id__is_active", 1).Filter("usergroup_id", user.Usergroup.ID).One(m)
}

// GetApplicationModuleIDByPrivilegeUser get data application module by privilege user
func GetApplicationModuleIDByPrivilegeUser(user *model.User) (pids string, e error) {

	o := orm.NewOrm()
	o.Raw("select group_concat(application_module_id) from application_privilege where usergroup_id = ?", user.Usergroup.ID).QueryRow(&pids)

	return pids, e
}

// GetApplicationMenuByUser get data from application menu by user
func GetApplicationMenuByUser(rq orm.RequestQuery, user *model.User) (m []model.ApplicationMenu, err error) {

	mid, _ := GetApplicationModuleIDByPrivilegeUser(user)

	// make new orm query
	q, _ := rq.Query(new(model.ApplicationMenu))

	q = q.Filter("is_active", 1).RelatedSel()
	q = q.Filter("application_module_id__in", strings.Split(mid, ","))

	// get data requested
	var mx []model.ApplicationMenu
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, nil
	}

	// return error some thing went wrong
	return nil, err
}

// GetApplicationModuleByUser get data application module by user
func GetApplicationModuleByUser(rq orm.RequestQuery, user *model.User) (m []model.ApplicationModule, err error) {

	mid, _ := GetApplicationModuleIDByPrivilegeUser(user)

	// make new orm query
	q, _ := rq.Query(new(model.ApplicationModule))

	q = q.Filter("is_active", 1).RelatedSel()

	q = q.Filter("id__in", strings.Split(mid, ","))

	// get data requested
	var mx []model.ApplicationModule
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, nil
	}

	// return error some thing went wrong
	return nil, err
}

// GetUsername get data user by username
func GetUsername(username string) (m *model.User, e error) {

	o := orm.NewOrm()
	if e = o.Raw("select * from user where binary username = ?", username).QueryRow(&m); e == nil {
		return m, nil
	}

	return nil, e
}
