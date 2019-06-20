// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package measurement_test

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

// test get all measurement success
func TestGetAllData(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	m := model.DummyMeasurement()
	m.IsDeleted = 0
	m.Save()
	id := common.Encrypt(m.ID)

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/measurement", "GET", 200},
		{"/v1/measurement/" + id, "GET", 200},
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

// test get all measurement No Token
func TestGetAllDataNoToken(t *testing.T) {
	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/measurement", "GET", http.StatusBadRequest},
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

// test get detail measurement when is deleted = 1
func TestGetDetailDataIsDeleted(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	m := model.DummyMeasurement()
	m.IsDeleted = 1
	m.Save()
	id := common.Encrypt(m.ID)

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/measurement/" + id, "GET", http.StatusNotFound},
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

// test get detail measurement when ID not found
func TestGetDetailDataNotFound(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/measurement/999999", "GET", http.StatusNotFound},
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

// test TestCreateMeasurement success
func TestCreateMeasurement(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("DELETE FROM measurement").Exec()

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"measurement_name": "tes", "note": "tes123"}, http.StatusOK},
		{tester.D{"measurement_name": ""}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.POST("/v1/measurement").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

	var total int64
	_ = o.Raw("select count(*) from measurement").QueryRow(&total)
	assert.Equal(t, int64(1), total)
}

// test TestCreateMeasurementFailMeasurementNameAlreadyExists
func TestCreateMeasurementNameAlreadyExists(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("DELETE FROM measurement").Exec()

	md1 := model.DummyMeasurement()
	md1.IsDeleted = 0
	md1.Save()

	md2 := model.DummyMeasurement()
	md2.IsDeleted = 1
	md2.Save()
	random := common.RandomStr(3)

	md3 := model.DummyMeasurement()
	md3.MeasurementName = "kodi" + random
	md3.IsDeleted = 0
	md3.Save()

	urand := "Kodi" + random

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"measurement_name": md1.MeasurementName, "note": "tes123"}, http.StatusUnprocessableEntity},
		{tester.D{"measurement_name": md2.MeasurementName, "note": "tes123"}, http.StatusOK},
		{tester.D{"measurement_name": urand, "note": "tes123"}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.POST("/v1/measurement").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// test TestCreateMeasurement No Token
func TestCreateMeasurementNoToken(t *testing.T) {

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"measurement_name": "tes", "note": "tes123"}, http.StatusBadRequest},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.POST("/v1/measurement").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// test update measurement success
func TestUpdateMeasurement(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	m := model.DummyMeasurement()
	m.IsDeleted = 0
	m.Save()
	id := common.Encrypt(m.ID)
	var put = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"measurement_name": "update name", "note": "update123"}, http.StatusOK},
		{tester.D{"measurement_name": ""}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range put {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.PUT("/v1/measurement/"+id).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

	m.Read("ID")
	assert.Equal(t, "update name", m.MeasurementName)
}

// test TestUpdateMeasurementNameAlreadyExists
func TestUpdateMeasurementNameAlreadyExists(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("DELETE FROM measurement").Exec()

	md1 := model.DummyMeasurement()
	md1.IsDeleted = 0
	md1.Save()

	md2 := model.DummyMeasurement()
	md2.IsDeleted = 1
	md2.Save()

	md3 := model.DummyMeasurement()
	md3.IsDeleted = 0
	md3.Save()

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"measurement_name": md3.MeasurementName, "note": "tes123"}, http.StatusUnprocessableEntity},
		{tester.D{"measurement_name": md2.MeasurementName, "note": "tes123"}, http.StatusOK},
	}

	ide := common.Encrypt(md1.ID)
	ng := tester.New()
	for _, tes := range create {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.PUT("/v1/measurement/"+ide).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// test update measurement is delete = 1
func TestUpdateMeasurementIsDeleted(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	m := model.DummyMeasurement()
	m.IsDeleted = 1
	m.Save()
	id := common.Encrypt(m.ID)
	var put = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"measurement_name": "update", "note": "update123"}, http.StatusNotFound},
	}
	ng := tester.New()
	for _, tes := range put {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.PUT("/v1/measurement/"+id).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// test update measurement, no token
func TestUpdateMeasurementNoToken(t *testing.T) {

	m := model.DummyMeasurement()
	m.IsDeleted = 0
	m.Save()
	id := common.Encrypt(m.ID)
	var put = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"measurement_name": "update", "note": "update123"}, http.StatusBadRequest},
	}
	ng := tester.New()
	for _, tes := range put {
		ng.PUT("/v1/measurement/"+id).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// Test Delete Measurement success
func TestDeleteMeasurementSuccess(t *testing.T) {

	m := model.DummyMeasurement()
	m.IsDeleted = int8(0)
	m.Save()

	scenario := tester.D{}

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-Type": "application/json"})
	ng.Method = "DELETE"
	ng.Path = "/v1/measurement/" + common.Encrypt(m.ID)
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/measurement/"+common.Encrypt(m.ID), "DELETE"))
	})
	m.Read("ID")
	assert.Equal(t, int8(1), m.IsDeleted)
}

// Test TestDeleteMeasurementFailBecauseBeingUsed
func TestDeleteMeasurementFailBecauseBeingUsed(t *testing.T) {

	m := model.DummyMeasurement()
	m.IsDeleted = int8(0)
	m.Save()

	miv := model.DummyItemVariant()
	miv.Measurement = m
	miv.Save()

	scenario := tester.D{}

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-Type": "application/json"})
	ng.Method = "DELETE"
	ng.Path = "/v1/measurement/" + common.Encrypt(m.ID)
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/measurement/"+common.Encrypt(m.ID), "DELETE"))
	})
	m.Read("ID")
	assert.Equal(t, int8(0), m.IsDeleted)
}

// Test Delete Measurement where ID not Found
func TestDeleteMeasurementNotFound(t *testing.T) {

	scenario := tester.D{}

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-Type": "application/json"})
	ng.Method = "DELETE"
	ng.Path = "/v1/measurement/7777777"
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/partnership/9999999", "DELETE"))
	})
}
