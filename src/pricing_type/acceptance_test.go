// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pricingType_test

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

func TestARouting(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	pricingType := model.DummyPricingType()
	ePTypeID := common.Encrypt(pricingType.ID)

	// success
	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/pricing-type", "GET", 200},
		{"/v1/pricing-type/" + ePTypeID, "GET", 200},
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
	//-------------------------------------------------------------------------------------------------------------------------//
	//-------------------------------------------------------------------------------------------------------------------------//
	//-------------------------------------------------------------------------------------------------------------------------//
	// gagal
	test.DbClean("pricing_type")
	routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/pricing-type", "GET", 200},
		{"/v1/pricing-type/" + ePTypeID, "GET", 404},
	}

	for _, ep := range routers {
		ng.Method = ep.method
		ng.Path = ep.endpoint
		ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, ep.expected, res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", ep.endpoint, ep.method))
		})
	}
}

func TestHandler_URLMappingPOST(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	name := common.RandomStr(5)

	mpt := model.DummyPricingType()
	mpt.IsDefault = 1
	mpt.Save()

	mpt2 := model.DummyPricingType()
	mpt2.ParentType = mpt
	mpt2.IsDefault = 1
	mpt2.Save()

	var data = []struct {
		req      tester.D
		expected int
	}{
		//sukses
		{tester.D{"parent_type_id": "", "type_name": name, "note": name + "asd", "rule_type": "", "nominal": float64(0), "is_percentage": int8(0), "is_default": int8(1)}, http.StatusOK},
		{tester.D{"parent_type_id": "", "type_name": common.RandomStr(45), "note": name + "asd", "rule_type": "", "nominal": float64(0), "is_percentage": int8(0), "is_default": int8(1)}, http.StatusOK},
		{tester.D{"parent_type_id": common.Encrypt(mpt.ID), "type_name": common.RandomStr(10), "note": name + "asd", "rule_type": "increment", "nominal": float64(500), "is_percentage": int8(0), "is_default": int8(0)}, http.StatusOK},
		// gagal
		// undecrypt
		{tester.D{"parent_type_id": "aaaa", "type_name": common.RandomStr(10), "note": name + "asd", "rule_type": "increment", "nominal": float64(500), "is_percentage": int8(0), "is_default": int8(0)}, http.StatusUnprocessableEntity},
		// id tidak terdaftar
		{tester.D{"parent_type_id": "999999999", "type_name": common.RandomStr(10), "note": name + "asd", "rule_type": "increment", "nominal": float64(500), "is_percentage": int8(0), "is_default": int8(0)}, http.StatusUnprocessableEntity},
		// type name harus unique
		{tester.D{"parent_type_id": "", "type_name": mpt2.TypeName, "note": name + "asd", "rule_type": "", "nominal": float64(0), "is_percentage": int8(0), "is_default": int8(1)}, http.StatusUnprocessableEntity},
		// harus parent
		{tester.D{"parent_type_id": common.Encrypt(mpt2.ID), "type_name": common.RandomStr(10), "note": name + "asd", "rule_type": "increment", "nominal": float64(500), "is_percentage": int8(0), "is_default": int8(0)}, http.StatusUnprocessableEntity},
		// semua field kosong
		{tester.D{"parent_type_id": "", "type_name": "", "note": "", "rule_type": "", "nominal": float64(0), "is_percentage": int8(0), "is_default": int8(0)}, http.StatusUnprocessableEntity},
		// karakter type name melebihi 45 karakter
		{tester.D{"parent_type_id": "", "type_name": common.RandomStr(46), "note": name + "asd", "rule_type": "", "nominal": float64(0), "is_percentage": int8(0), "is_default": int8(1)}, http.StatusUnprocessableEntity},
		// Tidak memiliki parent tetapi rule type nya diisi jadi gagal
		{tester.D{"parent_type_id": "", "type_name": common.RandomStr(10), "note": name + "asd", "rule_type": "decrement", "nominal": float64(0), "is_percentage": int8(0), "is_default": int8(0)}, http.StatusUnprocessableEntity},
		// Tidak memiliki parent tetapi nominal nya tidak kosong
		{tester.D{"parent_type_id": "", "type_name": common.RandomStr(10), "note": name + "asd", "rule_type": "decrement", "nominal": float64(1000), "is_percentage": int8(0), "is_default": int8(0)}, http.StatusUnprocessableEntity},
		// Memiliki parent tetapi is default == 1 jadi gagal
		{tester.D{"parent_type_id": common.Encrypt(mpt.ID), "type_name": common.RandomStr(10), "note": name + "asd", "rule_type": "", "nominal": float64(0), "is_percentage": int8(0), "is_default": int8(1)}, http.StatusUnprocessableEntity},
		//  memiliki parent dengan rule kosong tetapi is percentage nya 1 dan nominal lebih dari 100
		{tester.D{"parent_type_id": common.Encrypt(mpt.ID), "type_name": common.RandomStr(10), "note": name + "asd", "rule_type": "", "nominal": float64(500), "is_percentage": int8(1), "is_default": int8(0)}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.POST("/v1/pricing-type").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// TestHandler_URLMapping_PutSuccess test update pricing type success
func TestHandler_URLMapping_PutSuccess(t *testing.T) {
	// hapus all pricing type
	orm.NewOrm().Raw("DELETE FROM pricing_type").Exec()

	// buat dummy parent pricing type
	pr := model.DummyPricingType()
	pr.IsDefault = int8(1)
	pr.Save()

	// buat dummy user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"parent_type_id": "",
		"rule_type":      "",
		"nominal":        0,
		"is_percentage":  1,
		"note":           "untuk pelanggan",
		"type_name":      pr.TypeName,
		"is_default":     1,
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/pricing-type/"+common.Encrypt(pr.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	// cek database
	pr.Read("ID")
	assert.Equal(t, "none", pr.RuleType)
	assert.Equal(t, float64(0), pr.Nominal)
	assert.Equal(t, "untuk pelanggan", pr.Note)
}

// TestHandler_URLMapping_PutSuccess2 test update pricing type success
func TestHandler_URLMapping_PutSuccess2(t *testing.T) {
	// hapus all pricing type
	orm.NewOrm().Raw("DELETE FROM pricing_type").Exec()

	// buat dummy parent pricing type
	pr := model.DummyPricingType()
	pr.IsDefault = int8(1)
	pr.Save()

	pr1 := model.DummyPricingType()
	pr1.ParentType = pr
	pr1.IsDefault = int8(0)
	pr1.Save()

	// buat dummy user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"parent_type_id": common.Encrypt(pr.ID),
		"rule_type":      "increment",
		"nominal":        100,
		"is_percentage":  1,
		"note":           "untuk pelanggan",
		"is_default":     0,
		"type_name":      pr1.TypeName,
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/pricing-type/"+common.Encrypt(pr1.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	// cek database
	pr1.Read("ID")
	assert.Equal(t, "increment", pr1.RuleType)
	assert.Equal(t, float64(100), pr1.Nominal)
	assert.Equal(t, "untuk pelanggan", pr1.Note)
}

// TestHandler_URLMapping_PutSuccess3 test update pricing type success
func TestHandler_URLMapping_PutSuccess3(t *testing.T) {
	// hapus all pricing type
	orm.NewOrm().Raw("DELETE FROM pricing_type").Exec()

	// buat dummy parent pricing type
	pr := model.DummyPricingType()
	pr.IsDefault = int8(1)
	pr.Save()

	pr1 := model.DummyPricingType()
	pr1.ParentType = pr
	pr1.IsDefault = int8(0)
	pr1.Save()

	// buat dummy user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"parent_type_id": common.Encrypt(pr.ID),
		"rule_type":      "increment",
		"nominal":        100,
		"is_percentage":  1,
		"note":           "untuk pelanggan",
		"is_default":     1,
		"type_name":      pr1.TypeName,
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/pricing-type/"+common.Encrypt(pr1.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	// cek database
	pr1.Read("ID")
	assert.Equal(t, "increment", pr1.RuleType)
	assert.Equal(t, float64(100), pr1.Nominal)
	assert.Equal(t, "untuk pelanggan", pr1.Note)
	pr.Read("ID")
	assert.Equal(t, int8(0), pr.IsDefault)
}

// TestHandler_URLMapping_PutFail1 test update pricing type fail
func TestHandler_URLMapping_PutFail1(t *testing.T) {
	// hapus all pricing type
	orm.NewOrm().Raw("DELETE FROM pricing_type").Exec()

	// buat dummy parent pricing type
	pr := model.DummyPricingType()
	pr.IsDefault = int8(0)
	pr.Save()

	pr1 := model.DummyPricingType()
	pr1.ParentType = pr
	pr1.IsDefault = int8(1)
	pr1.Save()

	// buat dummy user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"parent_type_id": "9999999",
		"rule_type":      "",
		"nominal":        500,
		"is_percentage":  1,
		"note":           "untuk pelanggan",
		"is_default":     0,
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/pricing-type/"+common.Encrypt(pr1.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMapping_PutFail2 test update pricing type fail
func TestHandler_URLMapping_PutFail2(t *testing.T) {
	// hapus all pricing type
	orm.NewOrm().Raw("DELETE FROM pricing_type").Exec()

	// buat dummy parent pricing type
	pr := model.DummyPricingType()
	pr.IsDefault = int8(0)
	pr.Save()

	pr1 := model.DummyPricingType()
	pr1.ParentType = pr
	pr1.IsDefault = int8(1)
	pr1.Save()

	// buat dummy user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"parent_type_id": "aaaa",
		"rule_type":      "increment",
		"nominal":        0,
		"is_percentage":  0,
		"note":           "untuk pelanggan",
		"is_default":     1,
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/pricing-type/"+common.Encrypt(pr1.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMapping_PutFail3 test update pricing type fail
func TestHandler_URLMapping_PutFail3(t *testing.T) {
	// hapus all pricing type
	orm.NewOrm().Raw("DELETE FROM pricing_type").Exec()

	// buat dummy parent pricing type
	pr := model.DummyPricingType()
	pr.IsDefault = int8(0)
	pr.Save()

	pr1 := model.DummyPricingType()
	pr1.ParentType = pr
	pr1.IsDefault = int8(1)
	pr1.Save()

	// buat dummy user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"parent_type_id": common.Encrypt(pr1.ID),
		"rule_type":      "increment",
		"nominal":        0,
		"is_percentage":  1,
		"note":           "untuk pelanggan",
		"is_default":     1,
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/pricing-type/"+common.Encrypt(pr1.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMapping_PutFail4 test update pricing type fail
func TestHandler_URLMapping_PutFail4(t *testing.T) {
	// hapus all pricing type
	orm.NewOrm().Raw("DELETE FROM pricing_type").Exec()

	// buat dummy parent pricing type
	pr := model.DummyPricingType()
	pr.IsDefault = int8(1)
	pr.Save()

	// buat dummy user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"parent_type_id": "",
		"rule_type":      "increment",
		"nominal":        10,
		"is_percentage":  1,
		"note":           "untuk pelanggan",
		"is_default":     1,
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/pricing-type/"+common.Encrypt(pr.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMapping_PutFail5 test update pricing type fail
func TestHandler_URLMapping_PutFail5(t *testing.T) {
	// hapus all pricing type
	orm.NewOrm().Raw("DELETE FROM pricing_type").Exec()

	// buat dummy parent pricing type
	pr := model.DummyPricingType()
	pr.IsDefault = int8(1)
	pr.Save()

	pr1 := model.DummyPricingType()
	pr1.ParentType = pr
	pr1.IsDefault = int8(0)
	pr1.Save()

	// buat dummy user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"parent_type_id": "",
		"rule_type":      "increment",
		"nominal":        0,
		"is_percentage":  1,
		"note":           "untuk pelanggan",
		"is_default":     1,
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/pricing-type/"+common.Encrypt(pr1.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMapping_PutFail6 test update pricing type fail
func TestHandler_URLMapping_PutFail6(t *testing.T) {
	// hapus all pricing type
	orm.NewOrm().Raw("DELETE FROM pricing_type").Exec()

	// buat dummy parent pricing type
	pr := model.DummyPricingType()
	pr.IsDefault = int8(0)
	pr.Save()

	pr1 := model.DummyPricingType()
	pr1.ParentType = nil
	pr1.IsDefault = int8(0)
	pr1.Save()

	// buat dummy user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"parent_type_id": common.Encrypt(pr.ID),
		"rule_type":      "increment",
		"nominal":        0,
		"is_percentage":  1,
		"note":           "untuk pelanggan",
		"is_default":     1,
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/pricing-type/"+common.Encrypt(pr1.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}
