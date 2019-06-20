// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package setting_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/test"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/common/faker"
	"git.qasico.com/cuxs/common/tester"
	"git.qasico.com/cuxs/orm"
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

func das() *model.ApplicationSetting {
	var m model.ApplicationSetting
	faker.Fill(&m, "ID")

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// test get detail application setting when ID not found
func TestGetDetailDataNoToken(t *testing.T) {

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/setting", "GET", http.StatusBadRequest},
	}

	ng := tester.New()
	for _, ep := range routers {
		ng.Method = ep.method
		ng.Path = ep.endpoint
		ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, ep.expected, res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", ep.endpoint, ep.method))
		})
	}
}

// test get all and detail application setting success
func TestGetAllDataAndDetailData(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	aps := das()
	id := common.Encrypt(aps.ID)
	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/setting", "GET", 200},
		{"/v1/setting/" + id, "GET", 200},
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
}

// test get detail application setting when ID not found
func TestGetDetailDataNotFound(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/setting/999999", "GET", http.StatusNotFound},
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
}

func TestUpdateApplicationSetting(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	aps := das()
	id := common.Encrypt(aps.ID)
	fmt.Println("", id)
	var put = []struct {
		req      tester.D
		apsName  string
		value    string
		expected int
	}{
		{tester.D{"application_setting_name": "tes", "value": "tes123"}, "tes", "tes123", http.StatusOK},
		{tester.D{"application_setting_name": "", "value": ""}, "", "", http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range put {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.PUT("/v1/setting/"+id).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))

				if res.Code == http.StatusOK {
					var as *model.ApplicationSetting
					o := orm.NewOrm()
					o.Raw("select * from application_setting where id = ?", aps.ID).QueryRow(&as)
					assert.Equal(t, tes.apsName, as.ApplicationSettingName)
					assert.Equal(t, tes.value, as.Value)
				}

			})
	}
}

// test update application setting, when login usergroup is not sysadmin or owner
func TestUpdateApplicationSetting2(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(4)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	aps := das()
	id := common.Encrypt(aps.ID)
	var put = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"application_setting_name": "tes", "value": "tes123"}, http.StatusUnauthorized},
	}
	ng := tester.New()
	for _, tes := range put {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.PUT("/v1/setting/"+id).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}
