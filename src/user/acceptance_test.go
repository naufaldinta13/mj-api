// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/test"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/common/tester"
	"git.qasico.com/cuxs/orm"
	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.Setup()

	// run tests
	res := m.Run()

	// cleanup
	test.DataCleanUp()

	os.Exit(res)
}

func TestARouting(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	dummyuser := model.DummyUser()
	dummyuser.IsActive = 1
	dummyuser.Save("IsActive")
	euserID := common.Encrypt(dummyuser.ID)

	// sukses
	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/user", "GET", 200},
		{"/v1/user/" + euserID, "GET", 200},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, ep := range routers {
		ng.Method = ep.method
		ng.Path = ep.endpoint
		ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, ep.expected, res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", ep.endpoint, ep.method))
		})
	}

	// gagal
	dummyuser.Delete()
	routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/user/" + euserID, "GET", 404},
	}

	for _, ep := range routers {
		ng.Method = ep.method
		ng.Path = ep.endpoint
		ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, ep.expected, res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", ep.endpoint, ep.method))
		})
	}
}

func TestHandler_URLMapping_POST(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	uname := random.String(5)

	var data = []struct {
		req      tester.D
		expected int
	}{
		//sukses
		{tester.D{"usergroup_id": "131072", "full_name": "Owner", "username": uname, "password": "qasico123", "confirm_password": "qasico123"}, http.StatusOK},
		// gagal
		// apabila usergroup_id gak diisi
		{tester.D{"usergroup_id": "", "full_name": "Owner", "username": uname, "password": "qasico123", "confirm_password": "qasico123"}, http.StatusUnprocessableEntity},
		// apabila full_name gak diisi
		{tester.D{"usergroup_id": "131072", "full_name": "", "username": uname, "password": "qasico123", "confirm_password": "qasico123"}, http.StatusUnprocessableEntity},
		// apabila username gak diisi
		{tester.D{"usergroup_id": "131072", "full_name": "Owner", "username": "", "password": "qasico123", "confirm_password": "qasico123"}, http.StatusUnprocessableEntity},
		// apabila password gak diisi
		{tester.D{"usergroup_id": "131072", "full_name": "Owner", "username": uname, "password": "", "confirm_password": "qasico123"}, http.StatusUnprocessableEntity},
		// apabila confirm_password gak diisi
		{tester.D{"usergroup_id": "131072", "full_name": "Owner", "username": uname, "password": "qasico123", "confirm_password": ""}, http.StatusUnprocessableEntity},
		// apabila usergroup_id gak ada di database
		{tester.D{"usergroup_id": "131072088629", "full_name": "Owner", "username": uname, "password": "qasico123", "confirm_password": "qasico123"}, http.StatusUnprocessableEntity},
		// apabila usergroup tidak bisa di decrypt
		{tester.D{"usergroup_id": "yiuuaa", "full_name": "Owner", "username": uname, "password": "qasico123", "confirm_password": "qasico123"}, http.StatusUnprocessableEntity},
		// apabila username sama
		{tester.D{"usergroup_id": "131072", "full_name": "Owner", "username": user.Username, "password": "qasico123", "confirm_password": "qasico123"}, http.StatusUnprocessableEntity},
		// apabila confirm_password gak sama dengan password
		{tester.D{"usergroup_id": "131072", "full_name": "Owner", "username": uname, "password": "qasico123", "confirm_password": "qasico"}, http.StatusUnprocessableEntity},
		// apabila password < 5
		{tester.D{"usergroup_id": "131072", "full_name": "Owner", "username": uname, "password": "qasi", "confirm_password": "qasi"}, http.StatusUnprocessableEntity},
		// usergroup id == 1 harusnya gagal
		{tester.D{"usergroup_id": "65536", "full_name": "Andi", "username": common.RandomStr(10), "password": "qasi", "confirm_password": "qasi"}, http.StatusUnprocessableEntity},
		// fullname dengan spasi berhasil
		{tester.D{"usergroup_id": "131072", "full_name": "Su bagja", "username": common.RandomStr(10), "password": "qasico", "confirm_password": "qasico"}, http.StatusOK},
		// fullname dengan special character gagal
		{tester.D{"usergroup_id": "131072", "full_name": "Subagja$^@&", "username": common.RandomStr(10), "password": "qasico", "confirm_password": "qasico"}, http.StatusUnprocessableEntity},
		// username dengan spasi gagal
		{tester.D{"usergroup_id": "131072", "full_name": "Subagja", "username": common.RandomStr(10) + " " + "asda", "password": "qasico", "confirm_password": "qasico"}, http.StatusUnprocessableEntity},
		// username dengan special character selain underscore Gagal
		{tester.D{"usergroup_id": "131072", "full_name": "Subagja", "username": common.RandomStr(10) + "*#)@&(*%" + "asda", "password": "qasico", "confirm_password": "qasico"}, http.StatusUnprocessableEntity},
		// username dengan underscore berhasil
		{tester.D{"usergroup_id": "131072", "full_name": "Subagja", "username": common.RandomStr(10) + "_" + "asda", "password": "qasico", "confirm_password": "qasico"}, http.StatusOK},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.POST("/v1/user").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
	// cek data apakah ada di database
	var d *model.User
	o := orm.NewOrm()
	e := o.Raw("select * from user where username = ?", uname).QueryRow(&d)
	assert.NoError(t, e)
	assert.Equal(t, "131072", common.Encrypt(d.Usergroup.ID))
	assert.Equal(t, "Owner", d.FullName)
	assert.Equal(t, uname, d.Username)
	e = common.PasswordHash(d.Password, "qasico123")
	assert.NoError(t, e)
}

func TestHandler_URLMapping_PUT_Change_Password(t *testing.T) {
	sysadmin := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(sysadmin)
	token := "Bearer " + sd.Token

	dummyuser := model.DummyUser()
	dummyuser.IsActive = 1
	dummyuser.Save("IsActive")
	euserID := common.Encrypt(dummyuser.ID)

	pwd, _ := common.PasswordHasher("qasico")

	var data = []struct {
		req      tester.D
		password string
		expected int
	}{
		//sukses
		{tester.D{"password": "qasico", "confirm_password": "qasico"}, pwd, http.StatusOK},
		// gagal
		// apabila password kosong
		{tester.D{"password": "", "confirm_password": "qasico"}, "", http.StatusUnprocessableEntity},
		//apabila confirm_password kosong
		{tester.D{"password": "qasico", "confirm_password": ""}, "", http.StatusUnprocessableEntity},
		//apabila confirm_password tidak sama dengan password
		{tester.D{"password": "qasico", "confirm_password": "qasic"}, "", http.StatusUnprocessableEntity},
		//apabila panjang password < 5
		{tester.D{"password": "qasi", "confirm_password": "qasi"}, "", http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/user/"+euserID+"/change-password").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))

				//cek data apakah benar di update
				if res.Code == http.StatusOK {
					assert.NotEqual(t, dummyuser.Password, tes.password)
				}
			})
	}
}

func TestHandler_URLMapping_PUT_InActive(t *testing.T) {
	sysadmin := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(sysadmin)
	token := "Bearer " + sd.Token

	//sukses
	dummyuser := model.DummyUser()
	dummyuser.IsActive = 1
	dummyuser.Save("IsActive")
	euserID := common.Encrypt(dummyuser.ID)

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	request := tester.D{}
	ng.PUT("/v1/user/"+euserID+"/inactive").SetJSON(request).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", request, res.Body.String()))
		assert.NotEqual(t, dummyuser.IsActive, int8(0))
	})

	//gagal
	dummyuser.IsActive = 0
	dummyuser.Save("IsActive")
	ng.PUT("/v1/user/"+euserID+"/inactive").SetJSON(request).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusNotFound, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", request, res.Body.String()))
	})
}
