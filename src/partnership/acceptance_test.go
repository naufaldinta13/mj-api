// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package partnership_test

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

// TestMain
func TestMain(m *testing.M) {
	test.Setup()

	// run tests
	res := m.Run()

	// cleanup
	test.DataCleanUp()

	os.Exit(res)
}

// TestHandler_URLMappingListPartnershipSuccess untuk mengetest read list partnerships dengan 2 partner
func TestHandler_URLMappingListPartnershipSuccess(t *testing.T) {
	// buat 2 dummy partnership
	model.DummyPartnership()
	model.DummyPartnership()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/partnership"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/partnership", "GET"))
	})
}

// TestHandler_URLMappingListPartnershipTokenFail untuk mengetest read list partnerships tanpa token
func TestHandler_URLMappingListPartnershipTokenFail(t *testing.T) {
	ng := tester.New()
	ng.Method = "GET"
	ng.Path = "/v1/partnership"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/partnership", "GET"))
	})
}

// TestHandler_URLMappingListPartnershipSuccess untuk mengetest show  partnerships menggunakan id benar
func TestHandler_URLMappingDetailPartnershipByIDSuccess(t *testing.T) {
	// buat dummy partnership
	fakePartner := model.DummyPartnership()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/partnership/" + common.Encrypt(fakePartner.ID)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/partnership/"+common.Encrypt(fakePartner.ID), "GET"))
	})
}

// TestHandler_URLMappingListPartnershipTokenFail untuk mengetest show  partnerships dengan id yang salah
func TestHandler_URLMappingDetailPartnershipByIDNoData(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/partnership/9999999"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/partnership/9999999", "GET"))
	})
}

// TestHandler_URLMappingPartnershipUpdateSuccessWithPlafon mengetest update dengan order rule plafon success
func TestHandler_URLMappingPartnershipUpdateSuccessWithPlafon(t *testing.T) {
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(0)
	fakePartner.Save()
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{"order_rule": "plafon", "max_plafon": float64(2), "full_name": "test partner", "email": "@mail.com",
		"phone": "0000", "address": "adds", "city": "city", "province": "prov", "bank_name": "bank", "bank_number": "0bank",
		"bank_holder": "banker", "sales_person": "sale", "visit_day": "monday", "note": "note test"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/partnership/"+common.Encrypt(fakePartner.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	// cek database
	fakePartner.Read("ID")
	assert.Equal(t, "plafon", fakePartner.OrderRule)
	assert.Equal(t, "adds", fakePartner.Address)
	assert.Equal(t, "0000", fakePartner.Phone)
}

// TestHandler_URLMappingPartnershipUpdateSuccessWithoutPlafon mengetest update dengan order rule bukan plafon dan max_plafon 0,success
func TestHandler_URLMappingPartnershipUpdateSuccessWithoutPlafon(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("DELETE FROM partnership").Exec()
	// dummy
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(0)
	fakePartner.Save()
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{"order_rule": "one_bill", "max_plafon": float64(0), "full_name": "test partner", "email": "@mail.com",
		"phone": "0000", "address": "adds", "city": "city", "province": "prov", "bank_name": "bank", "bank_number": "0bank",
		"bank_holder": "banker", "sales_person": "sale", "visit_day": "monday", "note": "note test"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/partnership/"+common.Encrypt(fakePartner.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	// cek database
	fakePartner.Read("ID")
	assert.Equal(t, "one_bill", fakePartner.OrderRule)
	assert.Equal(t, float64(0), fakePartner.MaxPlafon)
	assert.Equal(t, "0000", fakePartner.Phone)
}

// TestHandler_URLMappingPartnershipUpdateFailWithoutPlafon mengetest update dengan order plafon dan max plafon 0,fail
func TestHandler_URLMappingPartnershipUpdateFailWithoutPlafon(t *testing.T) {
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(0)
	fakePartner.Save()
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{"order_rule": "plafon", "max_plafon": float64(0), "full_name": "test partner", "email": "@mail.com",
		"phone": "0000", "address": "adds", "city": "city", "province": "prov", "bank_name": "bank", "bank_number": "0bank",
		"bank_holder": "banker", "sales_person": "sale", "visit_day": "monday", "note": "note test"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/partnership/"+common.Encrypt(fakePartner.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPartnershipUpdateFailBodyEmpty mengetest update dengan body kosong,fail
func TestHandler_URLMappingPartnershipUpdateFailBodyEmpty(t *testing.T) {
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(0)
	fakePartner.Save()
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/partnership/"+common.Encrypt(fakePartner.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPartnershipUpdateFailNoBody mengetest update tanpa body, fail
func TestHandler_URLMappingPartnershipUpdateFailNoBody(t *testing.T) {
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(0)
	fakePartner.Save()
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/partnership/"+common.Encrypt(fakePartner.ID)).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(415), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", "", res.Body.String()))
	})
}

// TestHandler_URLMappingPartnershipUpdateFailPartnerIsDeleted mengetest update dengan partner is deleted, fail
func TestHandler_URLMappingPartnershipUpdateFailPartnerIsDeleted(t *testing.T) {
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(1)
	fakePartner.Save()
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{"order_rule": "plafon", "max_plafon": float64(2), "full_name": "test partner", "email": "@mail.com",
		"phone": "0000", "address": "adds", "city": "city", "province": "prov", "bank_name": "bank", "bank_number": "0bank",
		"bank_holder": "banker", "sales_person": "sale", "visit_day": "monday", "note": "note test"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/partnership/"+common.Encrypt(fakePartner.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPartnershipUpdateFailFullnameIsEmpty mengetest update dengan fullname kosong, fail
func TestHandler_URLMappingPartnershipUpdateFailFullnameIsEmpty(t *testing.T) {
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(0)
	fakePartner.Save()
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{"order_rule": "plafon", "max_plafon": float64(2), "full_name": "", "email": "@mail.com",
		"phone": "0000", "address": "adds", "city": "city", "province": "prov", "bank_name": "bank", "bank_number": "0bank",
		"bank_holder": "banker", "sales_person": "sale", "visit_day": "monday", "note": "note test"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/partnership/"+common.Encrypt(fakePartner.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPartnershipUpdateFailPartnerNotFound mengetest update dengan partnership id tidak ada, fail
func TestHandler_URLMappingPartnershipUpdateFailPartnerNotFound(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{"order_rule": "plafon", "max_plafon": float64(2), "full_name": "name", "email": "@mail.com",
		"phone": "0000", "address": "adds", "city": "city", "province": "prov", "bank_name": "bank", "bank_number": "0bank",
		"bank_holder": "banker", "sales_person": "sale", "visit_day": "monday", "note": "note test"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/partnership/9999999").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPartnershipCreateSuccessWithoutPlafon mengetest create dengan order bukan plafon, success
func TestHandler_URLMappingPartnershipCreateSuccessWithoutPlafon(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("DELETE FROM partnership").Exec()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{"partnership_type": "customer", "order_rule": "one_bill", "max_plafon": float64(0), "full_name": "test partner", "email": "@mail.com",
		"phone": "0000", "address": "adds", "city": "city", "province": "prov", "bank_name": "bank", "bank_number": "0bank",
		"bank_holder": "banker", "sales_person": "sale", "visit_day": "monday", "note": "note test"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/partnership").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	// cek database
	dat := &model.Partnership{PartnershipType: "customer"}
	dat.Read("PartnershipType")
	assert.Equal(t, "one_bill", dat.OrderRule)
	assert.Equal(t, float64(0), dat.MaxPlafon)
	assert.Equal(t, sd.User.ID, dat.CreatedBy.ID)
	assert.Equal(t, int8(0), dat.IsDefault)
	assert.Equal(t, int8(0), dat.IsDeleted)
}

// TestHandler_URLMappingPartnershipCreateSuccessWithPlafon mengetest create dengan order plafon dan max plafon 5, success
func TestHandler_URLMappingPartnershipCreateSuccessWithPlafon(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("DELETE FROM partnership").Exec()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{"partnership_type": "supplier", "order_rule": "plafon", "max_plafon": float64(5), "full_name": "test partner", "email": "@mail.com",
		"phone": "0000", "address": "adds", "city": "city", "province": "prov", "bank_name": "bank", "bank_number": "0bank",
		"bank_holder": "banker", "sales_person": "sale", "visit_day": "monday", "note": "note test"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/partnership").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	// cek database
	dat := &model.Partnership{PartnershipType: "supplier"}
	dat.Read("PartnershipType")
	assert.Equal(t, "plafon", dat.OrderRule)
	assert.Equal(t, float64(5), dat.MaxPlafon)
	assert.Equal(t, sd.User.ID, dat.CreatedBy.ID)
	assert.Equal(t, int8(0), dat.IsDefault)
	assert.Equal(t, int8(0), dat.IsDeleted)
}

// TestHandler_URLMappingPartnershipCreateFailWithPlafon mengetest create dengan order plafon dan max plafon 0, fail
func TestHandler_URLMappingPartnershipCreateFailWithPlafon(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("DELETE FROM partnership").Exec()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{"partnership_type": "supplier", "order_rule": "plafon", "max_plafon": float64(0), "full_name": "test partner", "email": "@mail.com",
		"phone": "11111", "address": "adds", "city": "city", "province": "prov", "bank_name": "bank", "bank_number": "0bank",
		"bank_holder": "banker", "sales_person": "sale", "visit_day": "monday", "note": "note test"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/partnership").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPartnershipCreateFailNoToken mengetest create tanpa token, fail
func TestHandler_URLMappingPartnershipCreateFailNoToken(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("DELETE FROM partnership").Exec()

	// setting body
	scenario := tester.D{"partnership_type": "supplier", "order_rule": "one_bill", "max_plafon": float64(0), "full_name": "test partner", "email": "@mail.com",
		"phone": "11111", "address": "adds", "city": "city", "province": "prov", "bank_name": "bank", "bank_number": "0bank",
		"bank_holder": "banker", "sales_person": "sale", "visit_day": "monday", "note": "note test"}

	ng := tester.New()
	ng.POST("/v1/partnership").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPartnershipCreateFailNoBody mengetest create tanpa body, fail
func TestHandler_URLMappingPartnershipCreateFailNoBody(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("DELETE FROM partnership").Exec()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/partnership").Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(415), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", "", res.Body.String()))
	})
}

// TestHandler_URLMappingPartnershipCreateFailNameEmpty mengetest create dengan full name kosong,fail
func TestHandler_URLMappingPartnershipCreateFailNameEmpty(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("DELETE FROM partnership").Exec()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{"partnership_type": "supplier", "order_rule": "one_bill", "max_plafon": float64(0), "email": "@mail.com",
		"phone": "11111", "address": "adds", "city": "city", "province": "prov", "bank_name": "bank", "bank_number": "0bank",
		"bank_holder": "banker", "sales_person": "sale", "visit_day": "monday", "note": "note test"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/partnership").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPartnershipFailValidationTooMuchCharacter testing ketika karakter yang diinputkan terlalubanyak
func TestHandler_URLMappingPartnershipFailValidationTooMuchCharacter(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("DELETE FROM partnership").Exec()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{"partnership_type": "supplier", "order_rule": "plafon", "max_plafon": float64(99999999999999999999999999999999999999),
		"full_name": common.RandomStr(46), "email": common.RandomStr(46) + "@mail.com",
		"phone": common.RandomStr(46) + "0000", "address": common.RandomStr(255) + "adds", "city": common.RandomStr(46) + "city",
		"province": common.RandomStr(46) + "prov", "bank_name": common.RandomStr(46) + "bank", "bank_number": common.RandomStr(46) + "0bank",
		"bank_holder": common.RandomStr(46) + "banker", "sales_person": common.RandomStr(46) + "sale", "visit_day": common.RandomStr(46) + "monday", "note": common.RandomStr(255) + "note test"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/partnership").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingArchiveSuccess mengetest archive, success
func TestHandler_URLMappingArchiveSuccess(t *testing.T) {
	// buat dummy partnership not archived and not deleted
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(0)
	fakePartner.IsArchived = int8(0)
	fakePartner.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "PUT"
	ng.Path = "/v1/partnership/" + common.Encrypt(fakePartner.ID) + "/archive"
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/partnership/"+common.Encrypt(fakePartner.ID)+"/archive", "PUT"))
	})
	fakePartner.Read("ID")
	assert.Equal(t, int8(1), fakePartner.IsArchived)
	assert.Equal(t, sd.User.ID, fakePartner.UpdatedBy.ID)
}

// TestHandler_URLMappingArchiveFailArchived mengetest archive yang sudah di archive, fail
func TestHandler_URLMappingArchiveFailArchived(t *testing.T) {
	// buat dummy partnership archived and not deleted
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(0)
	fakePartner.IsArchived = int8(1)
	fakePartner.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "PUT"
	ng.Path = "/v1/partnership/" + common.Encrypt(fakePartner.ID) + "/archive"
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/partnership/"+common.Encrypt(fakePartner.ID)+"/archive", "PUT"))
	})
}

// TestHandler_URLMappingArchiveFailArchivedAndDeleted mengetest archive yang sudah di archive dan delete, fail
func TestHandler_URLMappingArchiveFailArchivedAndDeleted(t *testing.T) {
	// buat dummy partnership archived and deleted
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(1)
	fakePartner.IsArchived = int8(1)
	fakePartner.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "PUT"
	ng.Path = "/v1/partnership/" + common.Encrypt(fakePartner.ID) + "/archive"
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/partnership/"+common.Encrypt(fakePartner.ID)+"/archive", "PUT"))
	})
}

// TestHandler_URLMappingArchiveFailPartnershipNotFound mengetest archive dengan id partnership not found, fail
func TestHandler_URLMappingArchiveFailPartnershipNotFound(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "PUT"
	ng.Path = "/v1/partnership/999999/archive"
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/partnership/999999/archive", "PUT"))
	})
}

// TestHandler_URLMappingUnArchiveSuccess mengetest unarchive, success
func TestHandler_URLMappingUnArchiveSuccess(t *testing.T) {
	// buat dummy partnership archived and not deleted
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(0)
	fakePartner.IsArchived = int8(1)
	fakePartner.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "PUT"
	ng.Path = "/v1/partnership/" + common.Encrypt(fakePartner.ID) + "/unarchive"
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/partnership/"+common.Encrypt(fakePartner.ID)+"/unarchive", "PUT"))
	})
	fakePartner.Read("ID")
	assert.Equal(t, int8(0), fakePartner.IsArchived)
	assert.Equal(t, int8(0), fakePartner.IsDeleted)
	assert.Equal(t, sd.User.ID, fakePartner.UpdatedBy.ID)
}

// TestHandler_URLMappingUnArchiveFailIsDeleted mengetest unarchive yang sudah di delete, fail
func TestHandler_URLMappingUnArchiveFailIsDeleted(t *testing.T) {
	// buat dummy partnership archived and deleted
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(1)
	fakePartner.IsArchived = int8(1)
	fakePartner.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "PUT"
	ng.Path = "/v1/partnership/" + common.Encrypt(fakePartner.ID) + "/unarchive"
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/partnership/"+common.Encrypt(fakePartner.ID)+"/unarchive", "PUT"))
	})
}

// TestHandler_URLMappingUnArchiveFailIDNotFound mengetest unarchive dengan id salah, fail
func TestHandler_URLMappingUnArchiveFailIDNotFound(t *testing.T) {
	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "PUT"
	ng.Path = "/v1/partnership/999999/unarchive"
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/partnership/999999/unarchive", "PUT"))
	})
}

// TestHandler_URLMappingUnArchiveFailNotArchived mengetest unarchive yang belum di archive, fail
func TestHandler_URLMappingUnArchiveFailNotArchived(t *testing.T) {
	// buat dummy partnership not archived and not deleted
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(0)
	fakePartner.IsArchived = int8(0)
	fakePartner.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "PUT"
	ng.Path = "/v1/partnership/" + common.Encrypt(fakePartner.ID) + "/unarchive"
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/partnership/"+common.Encrypt(fakePartner.ID)+"/unarchive", "PUT"))
	})
}

// TestHandler_URLMappingUnArchiveFailNotArchived mengetest unarchive tanpa body, fail
func TestHandler_URLMappingUnArchiveFailNoBody(t *testing.T) {
	// buat dummy partnership not archived and not deleted
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(0)
	fakePartner.IsArchived = int8(1)
	fakePartner.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "PUT"
	ng.Path = "/v1/partnership/" + common.Encrypt(fakePartner.ID) + "/unarchive"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(415), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/partnership/"+common.Encrypt(fakePartner.ID)+"/unarchive", "PUT"))
	})
}

// TestHandler_URLMappingDeleteSuccess mengetest delete, success
func TestHandler_URLMappingDeleteSuccess(t *testing.T) {
	// buat dummy partnership archived and not deleted
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(0)
	fakePartner.IsArchived = int8(1)
	fakePartner.Save()
	// setting body
	scenario := tester.D{}
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-Type": "application/json"})
	ng.Method = "DELETE"
	ng.Path = "/v1/partnership/" + common.Encrypt(fakePartner.ID)
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/partnership/"+common.Encrypt(fakePartner.ID), "DELETE"))
	})
	fakePartner.Read("ID")
	assert.Equal(t, int8(1), fakePartner.IsArchived)
	assert.Equal(t, int8(1), fakePartner.IsDeleted)
}

// TestHandler_URLMappingDeleteSuccess mengetest delete dengan is archived 0, fail
func TestHandler_URLMappingDeleteFailNotArchived(t *testing.T) {
	// buat dummy partnership not archived and not deleted
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(0)
	fakePartner.IsArchived = int8(0)
	fakePartner.Save()
	// setting body
	scenario := tester.D{}
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-Type": "application/json"})
	ng.Method = "DELETE"
	ng.Path = "/v1/partnership/" + common.Encrypt(fakePartner.ID)
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/partnership/"+common.Encrypt(fakePartner.ID), "DELETE"))
	})
}

// TestHandler_URLMappingDeleteSuccess mengetest delete dengan is delete 1, fail
func TestHandler_URLMappingDeleteFailAlreadyDelete(t *testing.T) {
	// buat dummy partnership archived and deleted
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(1)
	fakePartner.IsArchived = int8(1)
	fakePartner.Save()
	// setting body
	scenario := tester.D{}
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-Type": "application/json"})
	ng.Method = "DELETE"
	ng.Path = "/v1/partnership/" + common.Encrypt(fakePartner.ID)
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/partnership/"+common.Encrypt(fakePartner.ID), "DELETE"))
	})
}

// TestHandler_URLMappingDeleteSuccess mengetest delete tanpa token, fail
func TestHandler_URLMappingDeleteFailNoToken(t *testing.T) {
	// buat dummy partnership archived and not deleted
	fakePartner := model.DummyPartnership()
	fakePartner.IsDeleted = int8(0)
	fakePartner.IsArchived = int8(1)
	fakePartner.Save()
	// setting body
	scenario := tester.D{}
	ng := tester.New()
	ng.Method = "DELETE"
	ng.Path = "/v1/partnership/" + common.Encrypt(fakePartner.ID)
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/partnership/"+common.Encrypt(fakePartner.ID), "DELETE"))
	})
}

// TestHandler_URLMappingDeleteFailIDNotFound mengetest delete dengan id salah, fail
func TestHandler_URLMappingDeleteFailIDNotFound(t *testing.T) {
	// setting body
	scenario := tester.D{}
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-Type": "application/json"})
	ng.Method = "DELETE"
	ng.Path = "/v1/partnership/99999999"
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/partnership/9999999", "DELETE"))
	})
}
