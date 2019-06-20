// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package category_test

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

	category := model.DummyItemCategory()
	category.IsDeleted = 0
	category.Save("IsDeleted")

	ecategoryID := common.Encrypt(category.ID)

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/category", "GET", 200},
		{"/v1/category/" + ecategoryID, "GET", 200},
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

	//gagal
	category.IsDeleted = 1
	category.Save("IsDeleted")
	routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/category/" + ecategoryID, "GET", 404},
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
	name1 := common.RandomStr(5)
	name2 := common.RandomStr(5)
	var data = []struct {
		req          tester.D
		categoryName string
		expected     int
	}{
		//sukses
		{tester.D{"category_name": name1, "note": "segala elektronik"}, name1, http.StatusOK},
		{tester.D{"category_name": name2, "note": ""}, name2, http.StatusOK},
		// gagal
		// apabila category_name kosong
		{tester.D{"category_name": "", "note": ""}, "", http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.POST("/v1/category").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
				if res.Code == http.StatusOK {
					category := &model.ItemCategory{CategoryName: tes.categoryName}
					e := category.Read("CategoryName")
					assert.NoError(t, e)
					assert.Equal(t, tes.categoryName, category.CategoryName)
					assert.Equal(t, int8(0), category.IsDeleted)
				}
			})
	}
}

func TestHandler_URLMapping_POSTWithCategoryNameAlreadyExists(t *testing.T) {
	cleanTable()
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	category1 := model.DummyItemCategory()
	category1.CategoryName = "elektronik"
	category1.IsDeleted = 0
	category1.Save()

	category2 := model.DummyItemCategory()
	category1.CategoryName = "elektroniks"
	category2.IsDeleted = 1
	category2.Save()

	var data = []struct {
		req          tester.D
		categoryName string
		expected     int
	}{
		//gagal category name already exists
		{tester.D{"category_name": "elektronik", "note": "segala elektronik"}, "elektronik", http.StatusUnprocessableEntity},
		//berhasil karena isdeleted 1
		{tester.D{"category_name": "elektroniks", "note": "segala elektronik"}, "elektroniks", http.StatusOK},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.POST("/v1/category").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
				if res.Code == http.StatusOK {
					category := &model.ItemCategory{CategoryName: tes.categoryName}
					e := category.Read("CategoryName")
					assert.NoError(t, e)
					assert.Equal(t, tes.categoryName, category.CategoryName)
					assert.Equal(t, int8(0), category.IsDeleted)
				}
			})
	}
}

func TestHandler_URLMapping_PUT(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	name1 := common.RandomStr(5)
	name2 := common.RandomStr(5)
	category := model.DummyItemCategory()
	category.IsDeleted = 0
	category.Save("IsDeleted")

	ecategoryID := common.Encrypt(category.ID)

	var data = []struct {
		req          tester.D
		categoryName string
		expected     int
	}{
		//sukses
		{tester.D{"category_name": name1, "note": "segala elektronik"}, name1, http.StatusOK},
		{tester.D{"category_name": name2, "note": ""}, name2, http.StatusOK},
		// gagal
		// apabila category_name kosong
		{tester.D{"category_name": "", "note": ""}, "", http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/category/"+ecategoryID).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
				if res.Code == http.StatusOK {
					category := &model.ItemCategory{CategoryName: tes.categoryName}
					e := category.Read("CategoryName")
					assert.NoError(t, e)
					assert.Equal(t, tes.categoryName, category.CategoryName)
					assert.Equal(t, int8(0), category.IsDeleted)
				}
			})
	}

	/*-----------------------------------------------------------------------------------------------------------------------------------------------------*/
	/*-----------------------------------------------------------------------------------------------------------------------------------------------------*/
	// gagal karena categoy not found
	category.IsDeleted = 1
	category.Save("IsDeleted")

	data = []struct {
		req          tester.D
		categoryName string
		expected     int
	}{
		{tester.D{"category_name": "elektronik", "note": "segala elektronik"}, "elektronik", http.StatusNotFound},
		{tester.D{"category_name": "elektronik2", "note": ""}, "elektronik2", http.StatusNotFound},
	}
	for _, tes := range data {
		ng.PUT("/v1/category/"+ecategoryID).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

func TestHandler_URLMapping_PUTWithCategoryNameAlreadyExists(t *testing.T) {
	cleanTable()
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	category1 := model.DummyItemCategory()
	category1.CategoryName = "elektronik"
	category1.IsDeleted = 0
	category1.Save()

	category2 := model.DummyItemCategory()
	category2.CategoryName = "elektronik2"
	category2.IsDeleted = 0
	category2.Save()

	ecategoryID := common.Encrypt(category1.ID)

	var data = []struct {
		req          tester.D
		categoryName string
		expected     int
	}{
		//sukses
		{tester.D{"category_name": "elektronik2", "note": "segala elektronik"}, "elektronik", http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/category/"+ecategoryID).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// cleanTable
func cleanTable() {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_category;").Exec()
}

func TestHandler_URLMapping_DELETE(t *testing.T) {
	sysadmin := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(sysadmin)
	token := "Bearer " + sd.Token
	category := model.DummyItemCategory()
	category.IsDeleted = 0
	category.Save("IsDeleted")
	ecategoryID := common.Encrypt(category.ID)
	//sukses
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-type": "application/json"})
	request := tester.D{}
	ng.DELETE("/v1/category/"+ecategoryID).SetJSON(request).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", request, res.Body.String()))
	})
	//gagal
	// category is_deleted =1
	ng.DELETE("/v1/category/"+ecategoryID).SetJSON(request).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusNotFound, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", request, res.Body.String()))
	})
	// category masih dipakai item
	category.IsDeleted = 0
	category.Save("IsDeleted")
	item := model.DummyItem()
	item.Category = category
	item.IsDeleted = 0
	item.Save("Category", "IsDeleted")

	ng.DELETE("/v1/category/"+ecategoryID).SetJSON(request).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", request, res.Body.String()))
	})
}
