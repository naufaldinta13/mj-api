// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package directPlacement_test

import (
	"fmt"
	"os"
	"testing"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/test"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/common/tester"
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

//// Show direct-placement /////////////////////////////////////////////////////////////////////////////

// TestHandler_URLMappingDetailDirectPlacementSuccess untuk mengetest show  direct placement menggunakan id, success
func TestHandler_URLMappingDetailDirectPlacementSuccess(t *testing.T) {
	// buat dummy partnership
	direct := model.DummyDirectPlacement()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/direct-placement/" + common.Encrypt(direct.ID)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/direct-placement/"+common.Encrypt(direct.ID), "GET"))
	})
}

// TestHandler_URLMappingDetailDirectPlacementFailNoToken untuk mengetest show  direct placement tanpa token, fail
func TestHandler_URLMappingDetailDirectPlacementFailNoToken(t *testing.T) {
	// buat dummy partnership
	direct := model.DummyDirectPlacement()

	// test
	ng := tester.New()
	ng.Method = "GET"
	ng.Path = "/v1/direct-placement/" + common.Encrypt(direct.ID)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/direct-placement/"+common.Encrypt(direct.ID), "GET"))
	})
}

// TestHandler_URLMappingDetailDirectPlacementFailWrongID untuk mengetest show  direct placement menggunakan id salah, fail
func TestHandler_URLMappingDetailDirectPlacementFailWrongID(t *testing.T) {

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/direct-placement/9999999"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/direct-placement/9999999", "GET"))
	})
}

// TestHandler_URLMappingDetailDirectPlacementFailUndecryptID untuk mengetest show  direct placement menggunakan undecrypt id, fail
func TestHandler_URLMappingDetailDirectPlacementFailUndecryptID(t *testing.T) {

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/direct-placement/aaaa"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/direct-placement/aaa", "GET"))
	})
}

//// Post direct-placement /////////////////////////////////////////////////////////////////////////////

// TestHandler_URLMappingDirectPlacementCreateSuccess mengetest create direct placement,success
func TestHandler_URLMappingDirectPlacementCreateSuccess(t *testing.T) {
	// dummy
	itemVar := model.DummyItemVariant()
	itemVar.IsDeleted = int8(0)
	itemVar.IsArchived = int8(0)
	itemVar.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"note": "note",
		"direct_placement_items": []tester.D{
			{
				"item_variant_id": common.Encrypt(itemVar.ID),
				"quantity":        float32(10),
				"unit_price":      float64(2000),
				"total_price":     float64(20000),
			},
		},
	}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/direct-placement").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingDirectPlacementCreateFailEmptyBody mengetest create direct empty body,fail
func TestHandler_URLMappingDirectPlacementCreateFailEmptyBody(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/direct-placement").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingDirectPlacementCreateFailItemVarIsDeleted mengetest create direct placement dengan item variant deleted,fail
func TestHandler_URLMappingDirectPlacementCreateFailItemVarIsDeleted(t *testing.T) {
	// dummy
	itemVar := model.DummyItemVariant()
	itemVar.IsDeleted = int8(1)
	itemVar.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"note": "note",
		"direct_placement_items": []tester.D{
			{
				"item_variant_id": common.Encrypt(itemVar.ID),
				"quantity":        float32(10),
				"unit_price":      float64(2000),
				"total_price":     float64(20000),
			},
		},
	}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/direct-placement").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingDirectPlacementCreateFailItemVarIsDeleted mengetest create direct placement dengan item variant archived,fail
func TestHandler_URLMappingDirectPlacementCreateFailItemVarIsArchived(t *testing.T) {
	// dummy
	itemVar := model.DummyItemVariant()
	itemVar.IsArchived = int8(1)
	itemVar.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"note": "note",
		"direct_placement_items": []tester.D{
			{
				"item_variant_id": common.Encrypt(itemVar.ID),
				"quantity":        float32(10),
				"unit_price":      float64(2000),
				"total_price":     float64(20000),
			},
		},
	}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/direct-placement").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingDirectPlacementCreateFailItemVarUndecrypt mengetest create direct placement dengan item variant undecrypt,fail
func TestHandler_URLMappingDirectPlacementCreateFailItemVarUndecrypt(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"note": "note",
		"direct_placement_items": []tester.D{
			{
				"item_variant_id": "aaa",
				"quantity":        float32(10),
				"unit_price":      float64(2000),
				"total_price":     float64(20000),
			},
		},
	}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/direct-placement").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingDirectPlacementCreateFailQuantityZero mengetest create direct placement dengan quantity 0,fail
func TestHandler_URLMappingDirectPlacementCreateFailQuantityZero(t *testing.T) {
	// dummy
	itemVar := model.DummyItemVariant()
	itemVar.IsDeleted = int8(0)
	itemVar.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"note": "note",
		"direct_placement_items": []tester.D{
			{
				"item_variant_id": common.Encrypt(itemVar.ID),
				"quantity":        float32(0),
				"unit_price":      float64(2000),
				"total_price":     float64(0),
			},
		},
	}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/direct-placement").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}
