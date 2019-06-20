// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package inventory_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/inventory"
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

// ClearPriceType membersihkan table pricing_type
func ClearPriceType() {
	o := orm.NewOrm()
	o.Raw("DELETE FROM pricing_type").Exec()
}

// TestHandler_URLMappingListItemVariantSuccess untuk mengetest read list item variant
func TestHandler_URLMappingListItemVariantSuccess(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant").Exec()
	model.DummyItemVariant()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/inventory/variant"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/inventory/variant", "GET"))
	})
}

// TestHandler_URLMappingListItemVariantFailNoToken untuk mengetest read list item variant tanpa token, fail
func TestHandler_URLMappingListItemVariantFailNoToken(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant").Exec()
	model.DummyItemVariant()

	ng := tester.New()
	ng.Method = "GET"
	ng.Path = "/v1/inventory/variant"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/inventory/variant", "GET"))
	})
}

// TestHandler_URLMappingItemVariantDetailSuccess untuk mengetest show item variant by id
func TestHandler_URLMappingItemVariantDetailSuccess(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant").Exec()
	m := model.DummyItemVariant()

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/inventory/variant/" + common.Encrypt(m.ID)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/inventory/variant/"+common.Encrypt(m.ID), "GET"))
	})
}

// TestHandler_URLMappingItemVariantDetailFail untuk mengetest fail show item variant by id
func TestHandler_URLMappingItemVariantDetailFail(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant").Exec()

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/inventory/variant/999999"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Should return with error '%s'", "404"))
	})
}

// TestHandler_URLMappingItemVariantDetailFailNoToken untuk mengetest show item variant tanpa token, fail
func TestHandler_URLMappingItemVariantDetailFailNoToken(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant").Exec()
	m := model.DummyItemVariant()

	ng := tester.New()
	ng.Method = "GET"
	ng.Path = "/v1/inventory/variant/" + common.Encrypt(m.ID)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/inventory/variant/"+common.Encrypt(m.ID), "GET"))
	})
}

// TestHandler_URLMappingItemVariantArchiveSuccess untuk mengetest archive item variant
func TestHandler_URLMappingItemVariantArchiveSuccess(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant").Exec()

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// data item yang mempunyai dua item variant
	item := model.DummyItem()
	item.IsArchived = 0
	item.IsDeleted = 0
	item.Save()

	// data item variant pertama yang belum terarchive
	itemVariant := model.DummyItemVariant()
	itemVariant.Item = item
	itemVariant.IsArchived = 0
	itemVariant.IsDeleted = 0
	itemVariant.Save()

	// data item variant kedua yang sudah terarchive
	itemVariant2 := model.DummyItemVariant()
	itemVariant2.Item = item
	itemVariant2.IsArchived = 1
	itemVariant2.IsDeleted = 0
	itemVariant2.Save()

	var data = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{}, http.StatusOK},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/inventory/variant/"+common.Encrypt(itemVariant.ID)+"/archive").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))

				if res.Code == http.StatusOK {

					// tes untuk mengecek data item variant pertama sudah terarchive atau belum
					var iv *model.ItemVariant
					o.Raw("select * from item_variant where id = ?", itemVariant.ID).QueryRow(&iv)
					assert.Equal(t, int8(1), iv.IsArchived)

					// tes untuk mengecek data item sudah ikut terarchive atau belum
					// seharusnya data item ikut terarchive karena kedua data item variant sudah terarchive
					var i *model.Item
					o.Raw("select * from item where id = ?", item.ID).QueryRow(&i)
					assert.Equal(t, int8(1), i.IsArchived)

				}

			})
	}
}

// TestHandler_URLMappingItemVariantArchiveFailWasArchived untuk mengetest archive item variant yang sudah di archive
func TestHandler_URLMappingItemVariantArchiveFailWasArchived(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	itemVariant := model.DummyItemVariant()
	itemVariant.IsArchived = 1
	itemVariant.IsDeleted = 0
	itemVariant.Save()

	var data = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/inventory/variant/"+common.Encrypt(itemVariant.ID)+"/archive").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// TestHandler_URLMappingItemVariantArchiveFailWasDeleted untuk mengetest archive item variant yang sudah di delete
func TestHandler_URLMappingItemVariantArchiveFailWasDeleted(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	itemVariant := model.DummyItemVariant()
	itemVariant.IsArchived = 1
	itemVariant.IsDeleted = 1
	itemVariant.Save()

	var data = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/inventory/variant/"+common.Encrypt(itemVariant.ID)+"/archive").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// TestHandler_URLMappingItemVariantUnarchiveSuccess untuk mengetest unarchive item variant
func TestHandler_URLMappingItemVariantUnarchiveSuccess(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant").Exec()

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// data item yang mempunyai dua item variant
	item := model.DummyItem()
	item.IsArchived = 1
	item.IsDeleted = 0
	item.Save()

	// data item variant pertama yang sudah terarchive
	itemVariant := model.DummyItemVariant()
	itemVariant.Item = item
	itemVariant.IsArchived = 1
	itemVariant.IsDeleted = 0
	itemVariant.Save()

	// data item variant kedua yang belum terarchive
	itemVariant2 := model.DummyItemVariant()
	itemVariant2.Item = item
	itemVariant2.IsArchived = 0
	itemVariant2.IsDeleted = 0
	itemVariant2.Save()

	var data = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{}, http.StatusOK},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/inventory/variant/"+common.Encrypt(itemVariant.ID)+"/unarchive").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))

				if res.Code == http.StatusOK {

					// tes untuk mengecek data item variant pertama sudah ter-unarchive atau belum
					var iv *model.ItemVariant
					o.Raw("select * from item_variant where id = ?", itemVariant.ID).QueryRow(&iv)
					assert.Equal(t, int8(0), iv.IsArchived)

					// tes untuk mengecek data item sudah ikut ter-unarchive atau belum
					// seharusnya data item ikut ter-unarchive karena kedua data item variant sudah ter-unarchive
					var i *model.Item
					o.Raw("select * from item where id = ?", item.ID).QueryRow(&i)
					assert.Equal(t, int8(0), i.IsArchived)

				}

			})
	}
}

// TestHandler_URLMappingItemVariantUnarchiveFailWasUnarchived untuk mengetest unarchive item variant yang sudah di unarchive
func TestHandler_URLMappingItemVariantUnarchiveFailWasUnarchived(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	itemVariant := model.DummyItemVariant()
	itemVariant.IsArchived = 0
	itemVariant.IsDeleted = 0
	itemVariant.Save()

	var data = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/inventory/variant/"+common.Encrypt(itemVariant.ID)+"/unarchive").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// TestHandler_URLMappingItemVariantUnarchiveFailWasDeleted untuk mengetest unarchive item variant yang sudah di delete
func TestHandler_URLMappingItemVariantUnarchiveFailWasDeleted(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	itemVariant := model.DummyItemVariant()
	itemVariant.IsArchived = 1
	itemVariant.IsDeleted = 1
	itemVariant.Save()

	var data = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/inventory/variant/"+common.Encrypt(itemVariant.ID)+"/unarchive").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// TestHandler_URLMappingItemVariantDeleteSuccess untuk mengetest delete item variant
func TestHandler_URLMappingItemVariantDeleteSuccess(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant").Exec()

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// data item yang mempunyai dua item variant
	item := model.DummyItem()
	item.IsArchived = 1
	item.IsDeleted = 0
	item.Save()

	// data item variant pertama yang belum terdelete
	itemVariant := model.DummyItemVariant()
	itemVariant.Item = item
	itemVariant.IsArchived = 1
	itemVariant.IsDeleted = 0
	itemVariant.Save()

	// data item variant kedua yang sudah terdelete
	itemVariant2 := model.DummyItemVariant()
	itemVariant2.Item = item
	itemVariant2.IsArchived = 1
	itemVariant2.IsDeleted = 1
	itemVariant2.Save()

	// data sales order yang finish
	so := model.DummySalesOrder()
	so.DocumentStatus = "finished"
	so.IsDeleted = 0
	so.Save()

	// data sales order item
	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.ItemVariant = itemVariant
	soi.Save()

	// data purchase order yang finish
	po := model.DummyPurchaseOrder()
	po.DocumentStatus = "finished"
	po.IsDeleted = 0
	po.Save()

	// data purchase order item
	poi := model.DummyPurchaseOrderItem()
	poi.PurchaseOrder = po
	poi.ItemVariant = itemVariant
	poi.Save()

	var data = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{}, http.StatusOK},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-type": "application/json"})
	for _, tes := range data {
		ng.DELETE("/v1/inventory/variant/"+common.Encrypt(itemVariant.ID)).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))

				if res.Code == http.StatusOK {

					// tes untuk mengecek data item variant pertama sudah terdelete atau belum
					var iv *model.ItemVariant
					o.Raw("select * from item_variant where id = ?", itemVariant.ID).QueryRow(&iv)
					assert.Equal(t, int8(1), iv.IsDeleted)

					// tes untuk mengecek data item sudah ikut terdelete atau belum
					// seharusnya data item ikut terdelete karena kedua data item variant sudah terdelete
					var i *model.Item
					o.Raw("select * from item where id = ?", item.ID).QueryRow(&i)
					assert.Equal(t, int8(1), i.IsDeleted)

				}

			})
	}
}

// TestHandler_URLMappingItemVariantDeleteFailWasDeleted untuk mengetest delete item variant yang sudah di delete
func TestHandler_URLMappingItemVariantDeleteFailWasDeleted(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant").Exec()

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// data item yang mempunyai dua item variant
	item := model.DummyItem()
	item.IsArchived = 1
	item.IsDeleted = 0
	item.Save()

	// data item variant pertama yang sudah terdelete
	itemVariant := model.DummyItemVariant()
	itemVariant.Item = item
	itemVariant.IsArchived = 1
	itemVariant.IsDeleted = 1
	itemVariant.Save()

	// data item variant kedua yang sudah terdelete
	itemVariant2 := model.DummyItemVariant()
	itemVariant2.Item = item
	itemVariant2.IsArchived = 1
	itemVariant2.IsDeleted = 1
	itemVariant2.Save()

	// data sales order yang finish
	so := model.DummySalesOrder()
	so.DocumentStatus = "finished"
	so.IsDeleted = 0
	so.Save()

	// data sales order item
	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.ItemVariant = itemVariant
	soi.Save()

	// data purchase order yang finish
	po := model.DummyPurchaseOrder()
	po.DocumentStatus = "finished"
	po.IsDeleted = 0
	po.Save()

	// data purchase order item
	poi := model.DummyPurchaseOrderItem()
	poi.PurchaseOrder = po
	poi.ItemVariant = itemVariant
	poi.Save()

	var data = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-type": "application/json"})
	for _, tes := range data {
		ng.DELETE("/v1/inventory/variant/"+common.Encrypt(itemVariant.ID)).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// TestHandler_URLMappingItemVariantDeleteFailUnarchive untuk mengetest delete item variant yang belum di archive
func TestHandler_URLMappingItemVariantDeleteFailUnarchive(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant").Exec()

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// data item yang mempunyai dua item variant
	item := model.DummyItem()
	item.IsArchived = 1
	item.IsDeleted = 0
	item.Save()

	// data item variant pertama yang belum di archive
	itemVariant := model.DummyItemVariant()
	itemVariant.Item = item
	itemVariant.IsArchived = 0
	itemVariant.IsDeleted = 0
	itemVariant.Save()

	// data item variant kedua yang sudah terdelete
	itemVariant2 := model.DummyItemVariant()
	itemVariant2.Item = item
	itemVariant2.IsArchived = 1
	itemVariant2.IsDeleted = 1
	itemVariant2.Save()

	// data sales order yang finish
	so := model.DummySalesOrder()
	so.DocumentStatus = "finished"
	so.IsDeleted = 0
	so.Save()

	// data sales order item
	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.ItemVariant = itemVariant
	soi.Save()

	// data purchase order yang finish
	po := model.DummyPurchaseOrder()
	po.DocumentStatus = "finished"
	po.IsDeleted = 0
	po.Save()

	// data purchase order item
	poi := model.DummyPurchaseOrderItem()
	poi.PurchaseOrder = po
	poi.ItemVariant = itemVariant
	poi.Save()

	var data = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-type": "application/json"})
	for _, tes := range data {
		ng.DELETE("/v1/inventory/variant/"+common.Encrypt(itemVariant.ID)).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// TestHandler_URLMappingItemVariantDeleteFailSO untuk mengetest delete item variant pada SO yang masih aktif
func TestHandler_URLMappingItemVariantDeleteFailSO(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant").Exec()

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// data item yang mempunyai dua item variant
	item := model.DummyItem()
	item.IsArchived = 1
	item.IsDeleted = 0
	item.Save()

	// data item variant pertama yang belum terdelete
	itemVariant := model.DummyItemVariant()
	itemVariant.Item = item
	itemVariant.IsArchived = 1
	itemVariant.IsDeleted = 0
	itemVariant.Save()

	// data item variant kedua yang sudah terdelete
	itemVariant2 := model.DummyItemVariant()
	itemVariant2.Item = item
	itemVariant2.IsArchived = 1
	itemVariant2.IsDeleted = 1
	itemVariant2.Save()

	// data sales order yang belum finish
	so := model.DummySalesOrder()
	so.DocumentStatus = "new"
	so.IsDeleted = 0
	so.Save()

	// data sales order item
	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.ItemVariant = itemVariant
	soi.Save()

	// data purchase order yang finish
	po := model.DummyPurchaseOrder()
	po.DocumentStatus = "finished"
	po.IsDeleted = 0
	po.Save()

	// data purchase order item
	poi := model.DummyPurchaseOrderItem()
	poi.PurchaseOrder = po
	poi.ItemVariant = itemVariant
	poi.Save()

	var data = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-type": "application/json"})
	for _, tes := range data {
		ng.DELETE("/v1/inventory/variant/"+common.Encrypt(itemVariant.ID)).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// TestHandler_URLMappingItemVariantDeleteFailPO untuk mengetest delete item variant pada PO yang masih aktif
func TestHandler_URLMappingItemVariantDeleteFailPO(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant").Exec()

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// data item yang mempunyai dua item variant
	item := model.DummyItem()
	item.IsArchived = 1
	item.IsDeleted = 0
	item.Save()

	// data item variant pertama yang belum terdelete
	itemVariant := model.DummyItemVariant()
	itemVariant.Item = item
	itemVariant.IsArchived = 1
	itemVariant.IsDeleted = 0
	itemVariant.Save()

	// data item variant kedua yang sudah terdelete
	itemVariant2 := model.DummyItemVariant()
	itemVariant2.Item = item
	itemVariant2.IsArchived = 1
	itemVariant2.IsDeleted = 1
	itemVariant2.Save()

	// data sales order yang finish
	so := model.DummySalesOrder()
	so.DocumentStatus = "finished"
	so.IsDeleted = 0
	so.Save()

	// data sales order item
	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.ItemVariant = itemVariant
	soi.Save()

	// data purchase order yang belum finish
	po := model.DummyPurchaseOrder()
	po.DocumentStatus = "new"
	po.IsDeleted = 0
	po.Save()

	// data purchase order item
	poi := model.DummyPurchaseOrderItem()
	poi.PurchaseOrder = po
	poi.ItemVariant = itemVariant
	poi.Save()

	var data = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-type": "application/json"})
	for _, tes := range data {
		ng.DELETE("/v1/inventory/variant/"+common.Encrypt(itemVariant.ID)).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// TestHandler_URLMappingGetDetailItemSuccess mengetest item ,success
func TestHandler_URLMappingGetDetailItemSuccess(t *testing.T) {
	// buat dummy item not deleted
	itm := model.DummyItem()
	itm.IsDeleted = int8(0)
	itm.Save()
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/inventory/item/" + common.Encrypt(itm.ID)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/item/"+common.Encrypt(itm.ID), "GET"))
	})
}

// TestHandler_URLMappingGetDetailItemFailNoID mengetest item dengan id kosong,fail
func TestHandler_URLMappingGetDetailItemFailNoID(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/inventory/item/999999"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/item/999999", "GET"))
	})
}

// TestHandler_URLMappingGetDetailItemFailIsDeleted mengetest item yang di delete,fail
func TestHandler_URLMappingGetDetailItemFailIsDeleted(t *testing.T) {
	// buat dummy item deleted
	itm := model.DummyItem()
	itm.IsDeleted = int8(1)
	itm.Save()
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/inventory/item/" + common.Encrypt(itm.ID)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/item/"+common.Encrypt(itm.ID), "GET"))
	})
}

// TestHandler_URLMappingGetDetailItemFailNoToken mengetest item tanpa token,fail
func TestHandler_URLMappingGetDetailItemFailNoToken(t *testing.T) {
	// buat dummy item not deleted
	itm := model.DummyItem()
	itm.IsDeleted = int8(0)
	itm.Save()

	ng := tester.New()
	ng.Method = "GET"
	ng.Path = "/v1/inventory/item/" + common.Encrypt(itm.ID)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/item/"+common.Encrypt(itm.ID), "GET"))
	})
}

// TestHandler_URLMappingItemCreateSuccess mengetest create item,success
func TestHandler_URLMappingItemCreateSuccess(t *testing.T) {
	test.DataCleanUp("pricing_type")
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemCreateSuccess mengetest create item,success minimum angka dibulatkan dua angka dibelakang koma
func TestHandler_URLMappingItemCreateSuccessMinimumItemDibulatkan(t *testing.T) {
	test.DataCleanUp("pricing_type")

	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10.44444),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
		itemV := model.ItemVariant{ExternalName: name2}
		itemV.Read("ExternalName")
		assert.Equal(t, float32(10.44), itemV.MinimumStock, "Minimum stock pada item variant harus dibulatkan. Actual:"+fmt.Sprint(itemV.MinimumStock))
	})
}

// TestHandler_URLMappingItemCreateSuccessSomeEmptyField mengetest create item dengan beberapa field kosong,success
func TestHandler_URLMappingItemCreateSuccessSomeEmptyField(t *testing.T) {
	test.DataCleanUp("pricing_type")
	name := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(0),
		"item_variants": []tester.D{
			{
				"external_name": "",
				"variant_name":  "",
				"minimum_stock": float32(0),
				"base_price":    float64(1000),
				"note":          "",
				"image":         "",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": "",
				"variant_name":  "",
				"minimum_stock": float32(0),
				"base_price":    float64(1000),
				"note":          "",
				"image":         "",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemCreateFailItemVariantUnitPriceCannotHaveComma1 mengetest create item,gagal unit price tidak boleh ada koma
func TestHandler_URLMappingItemCreateFailItemVariantUnitPriceCannotHaveComma1(t *testing.T) {
	test.DataCleanUp("pricing_type")
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000.00099),
						"note":            "",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000.00099),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000.00099),
						"note":            "",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000.00099),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000.00099),
						"note":            "",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000.00099),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000.00099),
						"note":            "",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000.00099),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemCreateFailItemVariantUnitPriceCannotHaveComma2 mengetest create item,gagal unit price tidak boleh ada koma
func TestHandler_URLMappingItemCreateFailItemVariantUnitPriceCannotHaveComma2(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(0.55),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(0.55),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemCreateFailItemVariantUnitPriceCannotHaveComma1 mengetest create item,gagal unit price tidak boleh ada koma
func TestHandler_URLMappingItemCreateFailItemVariantBasePriceCannotHaveComma1(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000.99),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000.99),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemCreateFailItemVariantUnitPriceCannotHaveComma2 mengetest create item,gagal unit price tidak boleh ada koma
func TestHandler_URLMappingItemCreateFailItemVariantBasePriceCannotHaveComma2(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(0.5),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(0.5),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemCreateFailItemVariantMinimumItemCannotBeMinus mengetest create item,gagal min item tidak boleh minus
func TestHandler_URLMappingItemCreateFailItemVariantMinimumItemCannotBeMinus(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(-10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(-10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemCreateFailItemNameUnique mengetest create item dengan item name yang sudah ada,fail
func TestHandler_URLMappingItemCreateFailItemNameUnique(t *testing.T) {
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// dummy item
	itm := model.DummyItem()
	itm.IsDeleted = int8(0)
	itm.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      itm.ItemName,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemCreateFailItemNameUniqueDeleted mengetest create item dengan item name yang sudah ada tapi di delete,success
func TestHandler_URLMappingItemCreateFailItemNameUniqueDeleted(t *testing.T) {
	test.DataCleanUp("pricing_type")
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// dummy item
	itm := model.DummyItem()
	itm.IsDeleted = int8(1)
	itm.IsArchived = int8(1)
	itm.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      itm.ItemName,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemCreateFailItemTypeWrong mengetest create item dengan item type salah,fail
func TestHandler_URLMappingItemCreateFailItemTypeWrong(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "xxxxx",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemFailCreateFieldEmpty mengetest create item dengan beberapa field empty,fail
func TestHandler_URLMappingItemFailCreateFieldEmpty(t *testing.T) {
	name := common.RandomStr(7)
	name1 := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	name4 := common.RandomStr(7)
	// dummy
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name1,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": "",
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name4,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(0),
						"note":            "note var price",
					},
					{
						"pricing_type_id": "",
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemFailCreateSameName mengetest create item dengan variant name sama dengan yang lain,fail
func TestHandler_URLMappingItemFailCreateSameName(t *testing.T) {
	name := common.RandomStr(7)
	name1 := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name1,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemCreateSuccessSameNameWithDatabase mengetest create item dengan variant name sama dengan di database,success
func TestHandler_URLMappingItemCreateSuccessSameNameWithDatabase(t *testing.T) {
	test.DataCleanUp("pricing_type")
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	varx := model.DummyItemVariant()
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  varx.VariantName,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemFailCreateItemNameMore200 mengetest create item dengan item name >200,fail
func TestHandler_URLMappingItemFailCreateItemNameMore200(t *testing.T) {
	name := common.RandomStr(201)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemFailCreateVarNameMore45 mengetest create item dengan variant name >45,fail
func TestHandler_URLMappingItemFailCreateVarNameMore45(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(46)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemFailCreatePricingTypeIDWrong mengetest create item dengan pricing type id salah,fail
func TestHandler_URLMappingItemFailCreatePricingTypeIDWrong(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": "999999",
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemFailCreateCategoryIDWrong mengetest create item dengan category id salah,fail
func TestHandler_URLMappingItemFailCreateCategoryIDWrong(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()

	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    "xxxxx",
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemFailCreateCategoryIDNotExist mengetest create item dengan category id tidak ada,fail
func TestHandler_URLMappingItemFailCreateCategoryIDNotExist(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()

	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    "9999999",
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemFailCreateMeasurementIDWrong mengetest create item dengan measurement id salah,fail
func TestHandler_URLMappingItemFailCreateMeasurementIDWrong(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": "xxxxx",
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemFailCreateMeasurementIDNotExist mengetest create item dengan measurement id tidak ada,fail
func TestHandler_URLMappingItemFailCreateMeasurementIDNotExist(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": "9999999",
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemFailCreateUnitPriceZero mengetest create item dengan unit price 0,fail
func TestHandler_URLMappingItemFailCreateUnitPriceZero(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(0),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(-1),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemFailCreateCategoryAndMeasureIsDelete mengetest create item dengan measurement dan category di delete,fail
func TestHandler_URLMappingItemFailCreateCategoryAndMeasureIsDelete(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(1)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(1)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemCreateFailCategoryDeleted mengetest create item dengan category di delete,fail
func TestHandler_URLMappingItemCreateFailCategoryDeleted(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(1)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemFailCreateNoteMore255 mengetest create item dengan note lebih 255 char,fail
func TestHandler_URLMappingItemFailCreateNoteMore255(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	note := common.RandomStr(256)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           note,
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            note,
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          note,
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemFailCreatePricingTypeHasParent mengetest create item dengan pricing type bukan parent,fail
func TestHandler_URLMappingItemFailCreatePricingTypeHasParent(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakeParentPricing := model.DummyPricingType()
	fakePriceType1 := model.DummyPricingType()
	fakePriceType1.ParentType = fakeParentPricing
	fakePriceType1.Save()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemFailCreateNoToken mengetest create item tanpa token,fail
func TestHandler_URLMappingItemFailCreateNoToken(t *testing.T) {
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemArchiveSuccess mengetest item archive, success
func TestHandler_URLMappingItemArchiveSuccess(t *testing.T) {
	// buat dummy item not archived and not deleted
	fakeItem := model.DummyItem()
	fakeItem.IsDeleted = int8(0)
	fakeItem.IsArchived = int8(0)
	fakeItem.Save()
	// dummy item variant 1
	itmVar := model.DummyItemVariant()
	itmVar.Item = &model.Item{ID: fakeItem.ID}
	itmVar.IsArchived = int8(0)
	itmVar.IsDeleted = int8(0)
	itmVar.Save()
	// dummy item variant 2
	itmVar2 := model.DummyItemVariant()
	itmVar2.Item = &model.Item{ID: fakeItem.ID}
	itmVar2.IsArchived = int8(1)
	itmVar2.IsDeleted = int8(0)
	itmVar2.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "PUT"
	ng.Path = "/v1/inventory/item/" + common.Encrypt(fakeItem.ID) + "/archive"
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/inventory/item/"+common.Encrypt(fakeItem.ID)+"/archive", "PUT"))
	})
	fakeItem.Read("ID")
	assert.Equal(t, int8(1), fakeItem.IsArchived)
	assert.Equal(t, int8(0), fakeItem.IsDeleted)
	assert.Equal(t, sd.User.ID, fakeItem.UpdatedBy.ID)

	itmVar.Read("ID")
	assert.Equal(t, int8(1), itmVar.IsArchived)
	assert.Equal(t, int8(0), itmVar.IsDeleted)
	assert.Equal(t, sd.User.ID, itmVar.UpdatedBy.ID)

	itmVar2.Read("ID")
	assert.Equal(t, int8(1), itmVar2.IsArchived)
	assert.Equal(t, int8(0), itmVar2.IsDeleted)
}

// TestHandler_URLMappingItemArchiveFailAlreadyArchived mengetest item archive yang sudah di archive, fail
func TestHandler_URLMappingItemArchiveFailAlreadyArchived(t *testing.T) {
	// buat dummy item  archived and not deleted
	fakeItem := model.DummyItem()
	fakeItem.IsDeleted = int8(0)
	fakeItem.IsArchived = int8(1)
	fakeItem.Save()
	// dummy item variant 1
	itmVar := model.DummyItemVariant()
	itmVar.Item = &model.Item{ID: fakeItem.ID}
	itmVar.IsArchived = int8(1)
	itmVar.IsDeleted = int8(0)
	itmVar.Save()
	// dummy item variant 2
	itmVar2 := model.DummyItemVariant()
	itmVar2.Item = &model.Item{ID: fakeItem.ID}
	itmVar2.IsArchived = int8(1)
	itmVar2.IsDeleted = int8(0)
	itmVar2.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "PUT"
	ng.Path = "/v1/inventory/item/" + common.Encrypt(fakeItem.ID) + "/archive"
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/inventory/item/"+common.Encrypt(fakeItem.ID)+"/archive", "PUT"))
	})
}

// TestHandler_URLMappingItemArchiveFailAlreadyDeleted mengetest item archive yang sudah di delete, fail
func TestHandler_URLMappingItemArchiveFailAlreadyDeleted(t *testing.T) {
	// buat dummy item  archived and  deleted
	fakeItem := model.DummyItem()
	fakeItem.IsDeleted = int8(1)
	fakeItem.IsArchived = int8(1)
	fakeItem.Save()
	// dummy item variant 1
	itmVar := model.DummyItemVariant()
	itmVar.Item = &model.Item{ID: fakeItem.ID}
	itmVar.IsArchived = int8(1)
	itmVar.IsDeleted = int8(1)
	itmVar.Save()
	// dummy item variant 2
	itmVar2 := model.DummyItemVariant()
	itmVar2.Item = &model.Item{ID: fakeItem.ID}
	itmVar2.IsArchived = int8(1)
	itmVar2.IsDeleted = int8(1)
	itmVar2.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "PUT"
	ng.Path = "/v1/inventory/item/" + common.Encrypt(fakeItem.ID) + "/archive"
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/inventory/item/"+common.Encrypt(fakeItem.ID)+"/archive", "PUT"))
	})
}

// TestHandler_URLMappingItemUnarchiveSuccess mengetest item unarchive, success
func TestHandler_URLMappingItemUnarchiveSuccess(t *testing.T) {
	// buat dummy item  archived and not deleted
	fakeItem := model.DummyItem()
	fakeItem.IsDeleted = int8(0)
	fakeItem.IsArchived = int8(1)
	fakeItem.Save()
	// dummy item variant 1
	itmVar := model.DummyItemVariant()
	itmVar.Item = &model.Item{ID: fakeItem.ID}
	itmVar.IsArchived = int8(1)
	itmVar.IsDeleted = int8(0)
	itmVar.Save()
	// dummy item variant 2
	itmVar2 := model.DummyItemVariant()
	itmVar2.Item = &model.Item{ID: fakeItem.ID}
	itmVar2.IsArchived = int8(1)
	itmVar2.IsDeleted = int8(0)
	itmVar2.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "PUT"
	ng.Path = "/v1/inventory/item/" + common.Encrypt(fakeItem.ID) + "/unarchive"
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/inventory/item/"+common.Encrypt(fakeItem.ID)+"/unarchive", "PUT"))
	})
	fakeItem.Read("ID")
	assert.Equal(t, int8(0), fakeItem.IsArchived)
	assert.Equal(t, int8(0), fakeItem.IsDeleted)
	assert.Equal(t, sd.User.ID, fakeItem.UpdatedBy.ID)

	itmVar.Read("ID")
	assert.Equal(t, int8(0), itmVar.IsArchived)
	assert.Equal(t, int8(0), itmVar.IsDeleted)
	assert.Equal(t, sd.User.ID, itmVar.UpdatedBy.ID)

	itmVar2.Read("ID")
	assert.Equal(t, int8(0), itmVar2.IsArchived)
	assert.Equal(t, int8(0), itmVar2.IsDeleted)
}

// TestHandler_URLMappingItemUnarchiveFailNotArchived mengetest item unarchive yang sudah belum archive, fail
func TestHandler_URLMappingItemUnarchiveFailNotArchived(t *testing.T) {
	// buat dummy item not archived and not deleted
	fakeItem := model.DummyItem()
	fakeItem.IsDeleted = int8(0)
	fakeItem.IsArchived = int8(0)
	fakeItem.Save()
	// dummy item variant 1
	itmVar := model.DummyItemVariant()
	itmVar.Item = &model.Item{ID: fakeItem.ID}
	itmVar.IsArchived = int8(0)
	itmVar.IsDeleted = int8(0)
	itmVar.Save()
	// dummy item variant 2
	itmVar2 := model.DummyItemVariant()
	itmVar2.Item = &model.Item{ID: fakeItem.ID}
	itmVar2.IsArchived = int8(0)
	itmVar2.IsDeleted = int8(0)
	itmVar2.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "PUT"
	ng.Path = "/v1/inventory/item/" + common.Encrypt(fakeItem.ID) + "/unarchive"
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/inventory/item/"+common.Encrypt(fakeItem.ID)+"/unarchive", "PUT"))
	})
}

// TestHandler_URLMappingItemUnarchiveFailAlreadyDeleted mengetest item unarchive yang sudah di delete, fail
func TestHandler_URLMappingItemUnarchiveFailAlreadyDeleted(t *testing.T) {
	// buat dummy item  archived and  deleted
	fakeItem := model.DummyItem()
	fakeItem.IsDeleted = int8(1)
	fakeItem.IsArchived = int8(1)
	fakeItem.Save()
	// dummy item variant 1
	itmVar := model.DummyItemVariant()
	itmVar.Item = &model.Item{ID: fakeItem.ID}
	itmVar.IsArchived = int8(1)
	itmVar.IsDeleted = int8(1)
	itmVar.Save()
	// dummy item variant 2
	itmVar2 := model.DummyItemVariant()
	itmVar2.Item = &model.Item{ID: fakeItem.ID}
	itmVar2.IsArchived = int8(1)
	itmVar2.IsDeleted = int8(1)
	itmVar2.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "PUT"
	ng.Path = "/v1/inventory/item/" + common.Encrypt(fakeItem.ID) + "/unarchive"
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/inventory/item/"+common.Encrypt(fakeItem.ID)+"/unarchive", "PUT"))
	})
}

// TestHandler_URLMappingItemDeleteSuccess mengetest item delete, success
func TestHandler_URLMappingItemDeleteSuccess(t *testing.T) {
	// buat dummy item archived and not deleted
	fakeItem := model.DummyItem()
	fakeItem.IsDeleted = int8(0)
	fakeItem.IsArchived = int8(1)
	fakeItem.Save()
	// dummy item variant 1
	itmVar := model.DummyItemVariant()
	itmVar.Item = &model.Item{ID: fakeItem.ID}
	itmVar.IsArchived = int8(1)
	itmVar.IsDeleted = int8(0)
	itmVar.Save()
	// dummy item variant 2
	itmVar2 := model.DummyItemVariant()
	itmVar2.Item = &model.Item{ID: fakeItem.ID}
	itmVar2.IsArchived = int8(1)
	itmVar2.IsDeleted = int8(1)
	itmVar2.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-Type": "application/json"})
	ng.Method = "DELETE"
	ng.Path = "/v1/inventory/item/" + common.Encrypt(fakeItem.ID)
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/inventory/item/"+common.Encrypt(fakeItem.ID), "DELETE"))
	})
	fakeItem.Read("ID")
	assert.Equal(t, int8(1), fakeItem.IsArchived)
	assert.Equal(t, int8(1), fakeItem.IsDeleted)

	itmVar.Read("ID")
	assert.Equal(t, int8(1), itmVar.IsArchived)
	assert.Equal(t, int8(1), itmVar.IsDeleted)

	itmVar2.Read("ID")
	assert.Equal(t, int8(1), itmVar2.IsArchived)
	assert.Equal(t, int8(1), itmVar2.IsDeleted)
}

// TestHandler_URLMappingItemDeleteFailNotArchived mengetest item delete dengan is archived 0, fail
func TestHandler_URLMappingItemDeleteFailNotArchived(t *testing.T) {
	// buat dummy item not archived and not deleted
	fakeItem := model.DummyItem()
	fakeItem.IsDeleted = int8(0)
	fakeItem.IsArchived = int8(0)
	fakeItem.Save()
	// dummy item variant 1
	itmVar := model.DummyItemVariant()
	itmVar.Item = &model.Item{ID: fakeItem.ID}
	itmVar.IsArchived = int8(1)
	itmVar.IsDeleted = int8(0)
	itmVar.Save()
	// dummy item variant 2
	itmVar2 := model.DummyItemVariant()
	itmVar2.Item = &model.Item{ID: fakeItem.ID}
	itmVar2.IsArchived = int8(1)
	itmVar2.IsDeleted = int8(1)
	itmVar2.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-Type": "application/json"})
	ng.Method = "DELETE"
	ng.Path = "/v1/inventory/item/" + common.Encrypt(fakeItem.ID)
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/inventory/item/"+common.Encrypt(fakeItem.ID), "DELETE"))
	})
}

// TestHandler_URLMappingItemDeleteFailNotArchivedVariant mengetest item delete dengan is archived 0 pada variant, fail
func TestHandler_URLMappingItemDeleteFailNotArchivedVariant(t *testing.T) {
	// buat dummy item  archived and not deleted
	fakeItem := model.DummyItem()
	fakeItem.IsDeleted = int8(0)
	fakeItem.IsArchived = int8(1)
	fakeItem.Save()
	// dummy item variant 1
	itmVar := model.DummyItemVariant()
	itmVar.Item = &model.Item{ID: fakeItem.ID}
	itmVar.IsArchived = int8(0)
	itmVar.IsDeleted = int8(0)
	itmVar.Save()
	// dummy item variant 2
	itmVar2 := model.DummyItemVariant()
	itmVar2.Item = &model.Item{ID: fakeItem.ID}
	itmVar2.IsArchived = int8(1)
	itmVar2.IsDeleted = int8(1)
	itmVar2.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-Type": "application/json"})
	ng.Method = "DELETE"
	ng.Path = "/v1/inventory/item/" + common.Encrypt(fakeItem.ID)
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/inventory/item/"+common.Encrypt(fakeItem.ID), "DELETE"))
	})
}

// TestHandler_URLMappingItemDeleteFailAlreadyDeleted mengetest item delete dengan is delete 1, fail
func TestHandler_URLMappingItemDeleteFailAlreadyDeleted(t *testing.T) {
	// buat dummy item  archived and  deleted
	fakeItem := model.DummyItem()
	fakeItem.IsDeleted = int8(1)
	fakeItem.IsArchived = int8(1)
	fakeItem.Save()
	// dummy item variant 1
	itmVar := model.DummyItemVariant()
	itmVar.Item = &model.Item{ID: fakeItem.ID}
	itmVar.IsArchived = int8(1)
	itmVar.IsDeleted = int8(1)
	itmVar.Save()
	// dummy item variant 2
	itmVar2 := model.DummyItemVariant()
	itmVar2.Item = &model.Item{ID: fakeItem.ID}
	itmVar2.IsArchived = int8(1)
	itmVar2.IsDeleted = int8(1)
	itmVar2.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-Type": "application/json"})
	ng.Method = "DELETE"
	ng.Path = "/v1/inventory/item/" + common.Encrypt(fakeItem.ID)
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/inventory/item/"+common.Encrypt(fakeItem.ID), "DELETE"))
	})
}

// TestHandler_URLMappingItemDeleteFailUsedByPO mengetest item delete yang sudah digunakan PO, fail
func TestHandler_URLMappingItemDeleteFailUsedByPO(t *testing.T) {
	// buat dummy item  archived and not deleted
	fakeItem := model.DummyItem()
	fakeItem.IsDeleted = int8(0)
	fakeItem.IsArchived = int8(1)
	fakeItem.Save()
	// dummy item variant 1
	itmVar := model.DummyItemVariant()
	itmVar.Item = &model.Item{ID: fakeItem.ID}
	itmVar.IsArchived = int8(1)
	itmVar.IsDeleted = int8(0)
	itmVar.Save()
	// dummy item variant 2
	itmVar2 := model.DummyItemVariant()
	itmVar2.Item = &model.Item{ID: fakeItem.ID}
	itmVar2.IsArchived = int8(1)
	itmVar2.IsDeleted = int8(1)
	itmVar2.Save()
	// dummy PO item
	POI := model.DummyPurchaseOrderItem()
	POI.ItemVariant = &model.ItemVariant{ID: itmVar.ID}
	POI.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-Type": "application/json"})
	ng.Method = "DELETE"
	ng.Path = "/v1/inventory/item/" + common.Encrypt(fakeItem.ID)
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/inventory/item/"+common.Encrypt(fakeItem.ID), "DELETE"))
	})
}

// TestHandler_URLMappingItemDeleteFailUsedBySO mengetest item delete yang sudah digunakan SO, fail
func TestHandler_URLMappingItemDeleteFailUsedBySO(t *testing.T) {
	// buat dummy item  archived and not deleted
	fakeItem := model.DummyItem()
	fakeItem.IsDeleted = int8(0)
	fakeItem.IsArchived = int8(1)
	fakeItem.Save()
	// dummy item variant 1
	itmVar := model.DummyItemVariant()
	itmVar.Item = &model.Item{ID: fakeItem.ID}
	itmVar.IsArchived = int8(1)
	itmVar.IsDeleted = int8(0)
	itmVar.Save()
	// dummy item variant 2
	itmVar2 := model.DummyItemVariant()
	itmVar2.Item = &model.Item{ID: fakeItem.ID}
	itmVar2.IsArchived = int8(1)
	itmVar2.IsDeleted = int8(1)
	itmVar2.Save()
	// dummy PO item
	SOI := model.DummySalesOrderItem()
	SOI.ItemVariant = &model.ItemVariant{ID: itmVar.ID}
	SOI.Save()

	// setting body
	scenario := tester.D{}

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-Type": "application/json"})
	ng.Method = "DELETE"
	ng.Path = "/v1/inventory/item/" + common.Encrypt(fakeItem.ID)
	ng.SetJSON(scenario)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("Validation Not Matched Or Should has 'endpoint %s' with method '%s'", "/v1/inventory/item/"+common.Encrypt(fakeItem.ID), "DELETE"))
	})
}

// TestHandler_URLMappingItemUpdateSuccess mengetest update item,success
func TestHandler_URLMappingItemUpdateSuccess(t *testing.T) {
	// bersihkan table pricing_type
	ClearPriceType()

	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy component
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()
	// dummy item
	itm := model.DummyItem()
	itm.IsDeleted = int8(0)
	itm.IsArchived = int8(0)
	itm.Category = category
	itm.Save()
	variant1 := model.DummyItemVariant()
	variant1.Item = &model.Item{ID: itm.ID}
	variant1.IsDeleted = int8(0)
	variant1.IsArchived = int8(1)
	variant1.Save()
	Price := model.DummyItemVariantPrice()
	delPriceTypeID := Price.PricingType.ID
	Price.ItemVariant = &model.ItemVariant{ID: variant1.ID}
	Price.PricingType = &model.PricingType{ID: fakePriceType1.ID}
	Price.Save()
	delPriceType := &model.PricingType{ID: delPriceTypeID}
	delPriceType.Delete()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"id":            common.Encrypt(variant1.ID),
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"id":              common.Encrypt(Price.ID),
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/inventory/item/"+common.Encrypt(itm.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	resItm, e := inventory.GetDetailItem("id", itm.ID)
	assert.NoError(t, e)
	assert.Equal(t, name, resItm.ItemName)
	assert.Equal(t, itm.ItemType, resItm.ItemType)
	assert.NotEqual(t, itm.Note, resItm.Note)
	assert.Equal(t, int(2), len(resItm.ItemVariants))
	for _, u := range resItm.ItemVariants {
		assert.Equal(t, "image", u.Image)
		if u.ID == variant1.ID {
			assert.Equal(t, name2, u.VariantName)
			assert.Equal(t, int(2), len(u.ItemVariantPrices))
		}
	}
}

// TestHandler_URLMappingItemUpdateFailPricingTypeNotParent mengetest update item dengan pricing type bukan parent,fail
func TestHandler_URLMappingItemUpdateFailPricingTypeNotParent(t *testing.T) {
	// bersihkan table pricing_type
	ClearPriceType()

	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	// dummy component
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType3.ParentType = fakePriceType2
	fakePriceType3.Save()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()
	// dummy item
	itm := model.DummyItem()
	itm.IsDeleted = int8(0)
	itm.IsArchived = int8(0)
	itm.Category = category
	itm.Save()
	variant1 := model.DummyItemVariant()
	variant1.Item = &model.Item{ID: itm.ID}
	variant1.IsDeleted = int8(0)
	variant1.IsArchived = int8(1)
	variant1.Save()
	Price := model.DummyItemVariantPrice()
	Price.ItemVariant = &model.ItemVariant{ID: variant1.ID}
	Price.PricingType = &model.PricingType{ID: fakePriceType1.ID}
	Price.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_name":      name,
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"category_id":    common.Encrypt(category.ID),
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"id":            common.Encrypt(variant1.ID),
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"id":              common.Encrypt(Price.ID),
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/inventory/item/"+common.Encrypt(itm.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemUpdateFailWrongID mengetest update item dengan id yang salah,fail
func TestHandler_URLMappingItemUpdateFailWrongID(t *testing.T) {
	// bersihkan table pricing_type
	ClearPriceType()

	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy component
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()
	// dummy item
	itm := model.DummyItem()
	itm.IsDeleted = int8(0)
	itm.IsArchived = int8(0)
	itm.Category = category
	itm.Save()
	variant1 := model.DummyItemVariant()
	variant1.Item = &model.Item{ID: itm.ID}
	variant1.IsDeleted = int8(0)
	variant1.IsArchived = int8(1)
	variant1.Save()
	Price := model.DummyItemVariantPrice()
	Price.ItemVariant = &model.ItemVariant{ID: variant1.ID}
	Price.PricingType = &model.PricingType{ID: fakePriceType1.ID}
	Price.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_name":      name,
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"category_id":    common.Encrypt(category.ID),
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"id":            common.Encrypt(variant1.ID),
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"id":              common.Encrypt(Price.ID),
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}
	itm.ID = int64(999999)
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/inventory/item/"+common.Encrypt(itm.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemUpdateFailItemNameSame mengetest update item dengan item name yang sama,fail
func TestHandler_URLMappingItemUpdateFailItemNameSame(t *testing.T) {
	// bersihkan table pricing_type
	ClearPriceType()

	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy component
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()
	// dummy item
	itm := model.DummyItem()
	itm.IsDeleted = int8(0)
	itm.IsArchived = int8(0)
	itm.Category = category
	itm.Save()
	itm2 := model.DummyItem()
	itm2.ItemName = "wew"
	itm2.IsDeleted = int8(0)
	itm2.IsArchived = int8(0)
	itm2.Category = category
	itm2.Save()
	variant1 := model.DummyItemVariant()
	variant1.Item = &model.Item{ID: itm.ID}
	variant1.IsDeleted = int8(0)
	variant1.IsArchived = int8(1)
	variant1.Save()
	Price := model.DummyItemVariantPrice()
	Price.ItemVariant = &model.ItemVariant{ID: variant1.ID}
	Price.PricingType = &model.PricingType{ID: fakePriceType1.ID}
	Price.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	fake := model.DummyItem()
	fake.ItemName = "wew"
	fake.Save()
	// setting body
	scenario := tester.D{
		"item_name":      fake.ItemName,
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"category_id":    common.Encrypt(category.ID),
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"id":            common.Encrypt(variant1.ID),
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"id":              common.Encrypt(Price.ID),
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/inventory/item/"+common.Encrypt(itm.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMappingItemUpdateFailSameVariantName mengetest update item dengan variant name sama,fail
func TestHandler_URLMappingItemUpdateFailSameVariantName(t *testing.T) {
	// bersihkan table pricing_type
	ClearPriceType()

	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy component
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()
	// dummy item
	itm := model.DummyItem()
	itm.IsDeleted = int8(0)
	itm.IsArchived = int8(0)
	itm.Category = category
	itm.Save()
	variant1 := model.DummyItemVariant()
	variant1.Item = &model.Item{ID: itm.ID}
	variant1.IsDeleted = int8(0)
	variant1.IsArchived = int8(1)
	variant1.Save()
	Price := model.DummyItemVariantPrice()
	Price.ItemVariant = &model.ItemVariant{ID: variant1.ID}
	Price.PricingType = &model.PricingType{ID: fakePriceType1.ID}
	Price.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_name":      name,
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"category_id":    common.Encrypt(category.ID),
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"id":            common.Encrypt(variant1.ID),
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"id":              common.Encrypt(Price.ID),
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/inventory/item/"+common.Encrypt(itm.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemUpdateFailDuplicatePricingType mengetest update item dengan pricing type duplicate,fail
func TestHandler_URLMappingItemUpdateFailDuplicatePricingType(t *testing.T) {
	// bersihkan table pricing_type
	ClearPriceType()

	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	// dummy component
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()
	// dummy item
	itm := model.DummyItem()
	itm.IsDeleted = int8(0)
	itm.IsArchived = int8(0)
	itm.Category = category
	itm.Save()
	variant1 := model.DummyItemVariant()
	variant1.Item = &model.Item{ID: itm.ID}
	variant1.IsDeleted = int8(0)
	variant1.IsArchived = int8(1)
	variant1.Save()
	Price := model.DummyItemVariantPrice()
	Price.ItemVariant = &model.ItemVariant{ID: variant1.ID}
	Price.PricingType = &model.PricingType{ID: fakePriceType1.ID}
	Price.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_name":      name,
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"category_id":    common.Encrypt(category.ID),
		"item_variants": []tester.D{
			{
				"id":            common.Encrypt(variant1.ID),
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"id":              common.Encrypt(Price.ID),
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/inventory/item/"+common.Encrypt(itm.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemUpdateFailItemNameMore200 mengetest update item dengan item name >200,fail
func TestHandler_URLMappingItemUpdateFailItemNameMore200(t *testing.T) {
	// bersihkan table pricing_type
	ClearPriceType()

	name := common.RandomStr(201)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy component
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()
	// dummy item
	itm := model.DummyItem()
	itm.IsDeleted = int8(0)
	itm.IsArchived = int8(0)
	itm.Category = category
	itm.Save()
	variant1 := model.DummyItemVariant()
	variant1.Item = &model.Item{ID: itm.ID}
	variant1.IsDeleted = int8(0)
	variant1.IsArchived = int8(1)
	variant1.Save()
	Price := model.DummyItemVariantPrice()
	Price.ItemVariant = &model.ItemVariant{ID: variant1.ID}
	Price.PricingType = &model.PricingType{ID: fakePriceType1.ID}
	Price.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_name":      name,
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"category_id":    common.Encrypt(category.ID),
		"item_variants": []tester.D{
			{
				"id":            common.Encrypt(variant1.ID),
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"id":              common.Encrypt(Price.ID),
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/inventory/item/"+common.Encrypt(itm.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMappingItemUpdateFailNoteAndVarNameMore255 mengetest update item dengan input variant name >45 dan note >255,fail
func TestHandler_URLMappingItemUpdateFailNoteAndVarNameMore255(t *testing.T) {
	// bersihkan table pricing_type
	ClearPriceType()

	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(56)
	note := common.RandomStr(256)
	// dummy component
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()
	// dummy item
	itm := model.DummyItem()
	itm.IsDeleted = int8(0)
	itm.IsArchived = int8(0)
	itm.Category = category
	itm.Save()
	variant1 := model.DummyItemVariant()
	variant1.Item = &model.Item{ID: itm.ID}
	variant1.IsDeleted = int8(0)
	variant1.IsArchived = int8(1)
	variant1.Save()
	Price := model.DummyItemVariantPrice()
	Price.ItemVariant = &model.ItemVariant{ID: variant1.ID}
	Price.PricingType = &model.PricingType{ID: fakePriceType1.ID}
	Price.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_name":      name,
		"measurement_id": common.Encrypt(measure.ID),
		"note":           note,
		"category_id":    common.Encrypt(category.ID),
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"id":            common.Encrypt(variant1.ID),
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"id":              common.Encrypt(Price.ID),
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            note,
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          note,
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/inventory/item/"+common.Encrypt(itm.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemUpdateFailAllIDWrong mengetest update item dengan semua id salah,fail
func TestHandler_URLMappingItemUpdateFailAllIDWrong(t *testing.T) {
	// bersihkan table pricing_type
	ClearPriceType()

	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy component
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()
	// dummy item
	itm := model.DummyItem()
	itm.IsDeleted = int8(0)
	itm.IsArchived = int8(0)
	itm.Category = category
	itm.Save()
	variant1 := model.DummyItemVariant()
	variant1.Item = &model.Item{ID: itm.ID}
	variant1.IsDeleted = int8(0)
	variant1.IsArchived = int8(1)
	variant1.Save()
	Price := model.DummyItemVariantPrice()
	Price.ItemVariant = &model.ItemVariant{ID: variant1.ID}
	Price.PricingType = &model.PricingType{ID: fakePriceType1.ID}
	Price.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_name":      name,
		"measurement_id": "9999999",
		"category_id":    "9999999",
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"id":            "9999999",
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"id":              "9999999",
						"pricing_type_id": "9999999",
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/inventory/item/"+common.Encrypt(itm.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemUpdateFailAllIDUndecrypt mengetest update item dengan semua id un decrypt,fail
func TestHandler_URLMappingItemUpdateFailAllIDUndecrypt(t *testing.T) {
	// bersihkan table pricing_type
	ClearPriceType()

	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy component
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()
	// dummy item
	itm := model.DummyItem()
	itm.IsDeleted = int8(0)
	itm.IsArchived = int8(0)
	itm.Category = category
	itm.Save()
	variant1 := model.DummyItemVariant()
	variant1.Item = &model.Item{ID: itm.ID}
	variant1.IsDeleted = int8(0)
	variant1.IsArchived = int8(1)
	variant1.Save()
	Price := model.DummyItemVariantPrice()
	Price.ItemVariant = &model.ItemVariant{ID: variant1.ID}
	Price.PricingType = &model.PricingType{ID: fakePriceType1.ID}
	Price.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_name":      name,
		"measurement_id": "aaaaa",
		"category_id":    "aaaa",
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"id":            "aaaa",
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"id":              "aaaa",
						"pricing_type_id": "aaaa",
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/inventory/item/"+common.Encrypt(itm.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemCreateSuccessSameNameWithDatabase mengetest create item dengan variant name sama dengan di database,success
func TestHandler_URLMappingItemCreateErrorPricingTypeIDNotSameInList(t *testing.T) {
	test.DataCleanUp("pricing_type")
	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	varx := model.DummyItemVariant()
	// dummy
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	fakePriceType3 := model.DummyPricingType()
	fakePriceType4 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_type":      "product",
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"external_name": name2,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  varx.VariantName,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType3.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType4.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/inventory/item").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingItemUpdateSuccessDeleteItemVariantPrice mengetest update item,success delete item variant price
func TestHandler_URLMappingItemUpdateSuccessDeleteItemVariantPrice(t *testing.T) {
	// bersihkan table pricing_type
	ClearPriceType()

	name := common.RandomStr(7)
	name2 := common.RandomStr(7)
	name3 := common.RandomStr(7)
	// dummy component
	fakePriceType1 := model.DummyPricingType()
	fakePriceType2 := model.DummyPricingType()
	category := model.DummyItemCategory()
	category.IsDeleted = int8(0)
	category.Save()
	measure := model.DummyMeasurement()
	measure.IsDeleted = int8(0)
	measure.Save()
	// dummy item
	itm := model.DummyItem()
	itm.IsDeleted = int8(0)
	itm.IsArchived = int8(0)
	itm.Category = category
	itm.Save()
	variant1 := model.DummyItemVariant()
	variant1.Item = &model.Item{ID: itm.ID}
	variant1.IsDeleted = int8(0)
	variant1.IsArchived = int8(1)
	variant1.Save()
	Price := model.DummyItemVariantPrice()
	delPriceTypeID := Price.PricingType.ID
	Price.ItemVariant = &model.ItemVariant{ID: variant1.ID}
	Price.PricingType = &model.PricingType{ID: fakePriceType1.ID}
	Price.Save()
	Price2 := model.DummyItemVariantPrice()
	Price2.ItemVariant = &model.ItemVariant{ID: variant1.ID}
	Price2.PricingType = &model.PricingType{ID: fakePriceType1.ID}
	Price2.Save()
	Price3 := model.DummyItemVariantPrice()
	Price3.ItemVariant = &model.ItemVariant{ID: variant1.ID}
	Price3.PricingType = &model.PricingType{ID: fakePriceType1.ID}
	Price3.Save()
	delPriceType := &model.PricingType{ID: delPriceTypeID}
	delPriceType.Delete()

	var ivp []*model.ItemVariantPrice
	orm.NewOrm().Raw("select * from item_variant_price where item_variant_id = ?", variant1.ID).QueryRows(&ivp)

	assert.Equal(t, 3, len(ivp), "dummy data item variant price harus 3, karenan akan dihapus salah satu item variant pricenya")

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"item_name":      name,
		"category_id":    common.Encrypt(category.ID),
		"measurement_id": common.Encrypt(measure.ID),
		"note":           "note",
		"has_variant":    int8(1),
		"item_variants": []tester.D{
			{
				"id":            common.Encrypt(variant1.ID),
				"external_name": name2,
				"variant_name":  name2,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"id":              common.Encrypt(Price.ID),
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
			{
				"external_name": name3,
				"variant_name":  name3,
				"minimum_stock": float32(10),
				"base_price":    float64(1000),
				"note":          "note var",
				"image":         "image",
				"item_variant_prices": []tester.D{
					{
						"pricing_type_id": common.Encrypt(fakePriceType1.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
					{
						"pricing_type_id": common.Encrypt(fakePriceType2.ID),
						"unit_price":      float64(2000),
						"note":            "note var price",
					},
				},
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/inventory/item/"+common.Encrypt(itm.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	resItm, e := inventory.GetDetailItem("id", itm.ID)
	assert.NoError(t, e)
	assert.Equal(t, name, resItm.ItemName)
	assert.Equal(t, itm.ItemType, resItm.ItemType)
	assert.NotEqual(t, itm.Note, resItm.Note)
	assert.Equal(t, int(2), len(resItm.ItemVariants))
	for _, u := range resItm.ItemVariants {
		assert.Equal(t, "image", u.Image)
		if u.ID == variant1.ID {
			assert.Equal(t, name2, u.VariantName)
			assert.Equal(t, int(2), len(u.ItemVariantPrices))
		}
	}
}
