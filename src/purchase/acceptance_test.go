// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchase_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

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

// TestHandler_URLMappingListPurchaseOrderSuccess untuk mengetest read list Purchase order, success
func TestHandler_URLMappingListPurchaseOrderSuccess(t *testing.T) {
	// hapus semua data dari tabel purchase order
	o := orm.NewOrm()
	o.Raw("DELETE FROM purchase_order").Exec()
	// buat dummy purchase order 2x
	model.DummyPurchaseOrder()
	model.DummyPurchaseOrder()

	// proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/purchase-order"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/purchase-order", "GET"))
	})
}

// TestHandler_URLMappingListPurchaseOrderNoToken untuk mengetest read list Purchase order tanpa token, fail
func TestHandler_URLMappingListPurchaseOrderNoToken(t *testing.T) {
	ng := tester.New()
	ng.Method = "GET"
	ng.Path = "/v1/purchase-order"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/purchase-order", "GET"))
	})
}

// TestHandler_URLMappingGetDetailPurchaseOrderSuccess untuk mengetest show Purchase order, success
func TestHandler_URLMappingGetDetailPurchaseOrderSuccess(t *testing.T) {
	// buat dummy purchase order
	slsItm := model.DummyPurchaseOrderItem()
	slsItm.PurchaseOrder.IsDeleted = int8(0)
	slsItm.PurchaseOrder.Save()
	// proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/purchase-order/" + common.Encrypt(slsItm.PurchaseOrder.ID)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/purchase-order/"+common.Encrypt(slsItm.PurchaseOrder.ID), "GET"))
	})
}

// TestHandler_URLMappingGetDetailPurchaseOrderNoToken untuk mengetest show Purchase order tanpa token, fail
func TestHandler_URLMappingGetDetailPurchaseOrderNoToken(t *testing.T) {
	// buat dummy purchase order
	slsItm := model.DummyPurchaseOrderItem()
	ng := tester.New()
	ng.Method = "GET"
	ng.Path = "/v1/purchase-order/" + common.Encrypt(slsItm.PurchaseOrder.ID)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/purchase-order/"+common.Encrypt(slsItm.PurchaseOrder.ID), "GET"))
	})
}

// TestHandler_URLMappingGetDetailPurchaseOrderFailNoData untuk mengetest show Purchase order dengan id salah, fail
func TestHandler_URLMappingGetDetailPurchaseOrderFailNoData(t *testing.T) {
	// proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/purchase-order/9999999"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/purchase-order/99999", "GET"))
	})
}

func TestHandler_URLMappingPOST(t *testing.T) {
	test.DataCleanUp("purchase_order")
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	ref := model.DummyPurchaseOrder()
	ref.DocumentStatus = "cancelled"
	ref.Save()
	erefId := common.Encrypt(ref.ID)

	supplier := model.DummyPartnership()
	supplier.PartnershipType = "supplier"
	supplier.Save("PartnershipType")
	esupplierID := common.Encrypt(supplier.ID)

	iv := model.DummyItemVariant()
	iv.IsDeleted = 0
	iv.IsArchived = 0
	iv.Save("IsDeleted", "IsArchived")
	eiv := common.Encrypt(iv.ID)

	var data = []struct {
		req      tester.D
		expected int
	}{
		//sukses
		{tester.D{"recognition_date": time.Now(), "reference_id": erefId, "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": esupplierID, "tax": 10, "discount": 10, "is_percentage": int8(1), "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 12, "unit_price": 1000, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusOK},
		// apabila recognition_date kosong
		{tester.D{"recognition_date": nil, "reference_id": erefId, "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": esupplierID, "tax": 10, "is_percentage": int8(1), "discount": 10, "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 12, "unit_price": 1000, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusUnprocessableEntity},
		// reference_id kosong
		{tester.D{"recognition_date": time.Now(), "reference_id": "", "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": esupplierID, "tax": 10, "discount": 10, "is_percentage": int8(1), "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 12, "unit_price": 1000, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusOK},
		//reference id gak valid
		{tester.D{"recognition_date": time.Now(), "reference_id": "jjkdhajd", "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": esupplierID, "tax": 10, "is_percentage": int8(1), "discount": 10, "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 12, "unit_price": 1000, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusUnprocessableEntity},
		//reference id gak ada di database
		{tester.D{"recognition_date": time.Now(), "reference_id": "99999999", "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": esupplierID, "tax": 10, "is_percentage": int8(1), "discount": 10, "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 12, "unit_price": 1000, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusUnprocessableEntity},
		//eta date kosong
		{tester.D{"recognition_date": time.Now(), "reference_id": erefId, "auto_invoiced": 1, "eta_date": nil, "supplier_id": esupplierID, "tax": 10, "is_percentage": int8(1), "discount": 10, "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 12, "unit_price": 1000, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusUnprocessableEntity},
		//supplier_id kosong
		{tester.D{"recognition_date": time.Now(), "reference_id": erefId, "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": "", "tax": 10, "is_percentage": int8(1), "discount": 10, "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 12, "unit_price": 1000, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusUnprocessableEntity},
		//supplier_id gak valid
		{tester.D{"recognition_date": time.Now(), "reference_id": erefId, "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": "xxxx", "tax": 10, "is_percentage": int8(1), "discount": 10, "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 12, "unit_price": 1000, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusUnprocessableEntity},
		//supplier_id gak ada di dalam database
		{tester.D{"recognition_date": time.Now(), "reference_id": erefId, "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": "9999999", "tax": 10, "is_percentage": int8(1), "discount": 10, "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 12, "unit_price": 1000, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusUnprocessableEntity},
		//tax kosong
		{tester.D{"recognition_date": time.Now(), "reference_id": erefId, "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": esupplierID, "tax": 0, "discount": 10, "is_percentage": int8(1), "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 12, "unit_price": 1000, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusOK},
		//discount kosong
		{tester.D{"recognition_date": time.Now(), "reference_id": erefId, "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": esupplierID, "tax": 10, "discount_amount": 1000, "is_percentage": int8(0), "shipment_cost": 5000, "note": "diskon kosong yang ada discount amount",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 12, "unit_price": 1000, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusOK},
		//shipment_cost kosong
		{tester.D{"recognition_date": time.Now(), "reference_id": erefId, "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": esupplierID, "tax": 10, "discount": 10, "is_percentage": int8(1), "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 12, "unit_price": 1000, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusOK},
		//note po kosong
		{tester.D{"recognition_date": time.Now(), "reference_id": erefId, "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": esupplierID, "tax": 10, "discount": 10, "is_percentage": int8(1), "shipment_cost": 5000, "note": "",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 12, "unit_price": 1000, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusOK},
		//item variant kosong
		{tester.D{"recognition_date": time.Now(), "reference_id": erefId, "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": esupplierID, "tax": 10, "is_percentage": int8(1), "discount": 10, "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": "", "quantity": 12, "unit_price": 1000, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusUnprocessableEntity},
		//item variant gak valid
		{tester.D{"recognition_date": time.Now(), "reference_id": erefId, "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": esupplierID, "tax": 10, "is_percentage": int8(1), "discount": 10, "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": "asada", "quantity": 12, "unit_price": 1000, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusUnprocessableEntity},
		//item variant gak ada di database
		{tester.D{"recognition_date": time.Now(), "reference_id": erefId, "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": esupplierID, "tax": 10, "is_percentage": int8(1), "discount": 10, "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": "999999999", "quantity": 12, "unit_price": 1000, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusUnprocessableEntity},
		//quantity kosong
		{tester.D{"recognition_date": time.Now(), "reference_id": erefId, "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": esupplierID, "tax": 10, "is_percentage": int8(1), "discount": 10, "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 0, "unit_price": 1000, "discount": 0, "note": "catatan PO Items"}},
		}, http.StatusUnprocessableEntity},
		//unit price kosong
		{tester.D{"recognition_date": time.Now(), "reference_id": erefId, "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": esupplierID, "tax": 10, "is_percentage": int8(1), "discount": 10, "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 12, "unit_price": 0, "discount": 1, "note": "catatan PO Items"}},
		}, http.StatusUnprocessableEntity},
		//poitems diskon kosong
		{tester.D{"recognition_date": time.Now(), "reference_id": erefId, "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": esupplierID, "tax": 10, "discount": 10, "is_percentage": int8(1), "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 12, "unit_price": 1000, "discount": 0, "note": "catatan PO Items"}},
		}, http.StatusOK},
		//note kosong
		{tester.D{"recognition_date": time.Now(), "reference_id": erefId, "auto_invoiced": 1, "eta_date": time.Now(), "supplier_id": esupplierID, "tax": 10, "discount": 10, "is_percentage": int8(1), "shipment_cost": 5000, "note": "catatan",
			"purchase_order_items": []tester.D{{"item_variant_id": eiv, "quantity": 12, "unit_price": 1000, "discount": 1, "note": ""}},
		}, http.StatusOK},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.POST("/v1/purchase-order").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

func TestHandler_URLMappingPUT(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	DummyPO := model.DummyPurchaseOrderItem()
	DummyPO.PurchaseOrder.IsDeleted = 0
	DummyPO.PurchaseOrder.DocumentStatus = "new"
	DummyPO.PurchaseOrder.TotalCharge = 10000
	DummyPO.PurchaseOrder.Save()

	epoid := common.Encrypt(DummyPO.PurchaseOrder.ID)

	DummyPO2 := model.DummyPurchaseOrderItem()
	DummyPO2.PurchaseOrder.IsDeleted = 0
	DummyPO2.PurchaseOrder.DocumentStatus = "finished"
	DummyPO2.PurchaseOrder.TotalCharge = 10000
	DummyPO2.PurchaseOrder.Save()

	epoid2 := common.Encrypt(DummyPO2.PurchaseOrder.ID)

	DummyIV := model.DummyItemVariantPrice()
	DummyIV.ItemVariant.IsArchived = 0
	DummyIV.ItemVariant.IsDeleted = 0
	DummyIV.ItemVariant.Save()
	eitemvar := common.Encrypt(DummyIV.ItemVariant.ID)

	DummyPartner := model.DummyPartnership()
	DummyPartner.PartnershipType = "supplier"
	DummyPartner.IsDeleted = 0
	DummyPartner.IsArchived = 0
	DummyPartner.Save()
	esupplierID := common.Encrypt(DummyPartner.ID)

	DummyPartner1 := model.DummyPartnership()
	DummyPartner1.PartnershipType = "customer"
	DummyPartner1.IsDeleted = 0
	DummyPartner1.IsArchived = 0
	DummyPartner1.Save()
	esupplierID1 := common.Encrypt(DummyPartner1.ID)

	DummyIV2 := model.DummyItemVariantPrice()
	DummyIV2.ItemVariant.IsArchived = 1
	DummyIV2.ItemVariant.IsDeleted = 1
	DummyIV2.ItemVariant.Save()
	eitemvar2 := common.Encrypt(DummyIV2.ItemVariant.ID)

	var data = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{
			"recognition_date": time.Now(),
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID,
			"tax":              10,
			"is_percentage":    1,
			"discount":         0,
			"shipment_cost":    5000,
			"note":             "is percentage 1 tetapi discount 0",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": "abcdefg",
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusUnprocessableEntity},
		{tester.D{
			"recognition_date": time.Now(),
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID,
			"tax":              10,
			"is_percentage":    1,
			"discount":         0,
			"shipment_cost":    5000,
			"note":             "is percentage 1 tetapi discount 0",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": "65536000",
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusUnprocessableEntity},
		{tester.D{
			"recognition_date": time.Now(),
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID,
			"tax":              10,
			"is_percentage":    1,
			"discount":         0,
			"shipment_cost":    5000,
			"note":             "is percentage 1 tetapi discount 0",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eitemvar2,
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusUnprocessableEntity},
		{tester.D{
			"recognition_date": time.Now(),
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID,
			"tax":              10,
			"is_percentage":    1,
			"discount":         0,
			"shipment_cost":    5000,
			"note":             "is percentage 1 tetapi discount 0",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eitemvar,
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusOK},
		{tester.D{
			"recognition_date": time.Now(),
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID,
			"tax":              10,
			"is_percentage":    0,
			"discount":         10,
			"shipment_cost":    5000,
			"note":             "is percentage 0 tetapi discount amount kosong",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eitemvar,
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusOK},
		{tester.D{
			"recognition_date": time.Now(),
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID1,
			"tax":              10,
			"is_percentage":    1,
			"discount":         10,
			"shipment_cost":    5000,
			"note":             "is percentage 0 tetapi discount amount kosong",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eitemvar,
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusUnprocessableEntity},
		{tester.D{
			"recognition_date": time.Now(),
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      "abcdefg",
			"tax":              10,
			"is_percentage":    1,
			"discount":         10,
			"shipment_cost":    5000,
			"note":             "is percentage 0 tetapi discount amount kosong",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eitemvar,
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusUnprocessableEntity},
		{tester.D{
			"recognition_date": time.Now(),
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      "65536000",
			"tax":              10,
			"is_percentage":    1,
			"discount":         10,
			"shipment_cost":    5000,
			"note":             "is percentage 0 tetapi discount amount kosong",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eitemvar,
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusUnprocessableEntity},
		{tester.D{
			"recognition_date": time.Now(),
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID,
			"tax":              10,
			"is_percentage":    0,
			"discount":         0,
			"discount_amount":  2000,
			"shipment_cost":    50000,
			"note":             "is percentage 0 tetapi discount amount kosong",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eitemvar,
					"quantity":        10,
					"unit_price":      1000,
					"discount":        0,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusOK},
		{tester.D{
			"recognition_date": time.Now(),
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID,
			"tax":              10,
			"is_percentage":    0,
			"discount":         0,
			"discount_amount":  0,
			"shipment_cost":    50000,
			"note":             "is percentage 0 tetapi discount amount kosong",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eitemvar,
					"quantity":        10,
					"unit_price":      1000,
					"discount":        0,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusOK},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/purchase-order/"+epoid).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

	var data2 = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{}, http.StatusNotFound},
	}

	for _, tes := range data2 {
		ng.PUT("/v1/purchase-order/65536000").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

	var data3 = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{
			"recognition_date": time.Now(),
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID,
			"tax":              10,
			"is_percentage":    1,
			"discount":         0,
			"shipment_cost":    5000,
			"note":             "is percentage 1 tetapi discount 0",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eitemvar,
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusUnprocessableEntity},
	}

	for _, tes := range data3 {
		ng.PUT("/v1/purchase-order/"+epoid2).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// untuk validasi baru dan semuanya error
func TestHandler_URLMappingPOST1(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	ref1 := model.DummyPurchaseOrder()
	ref1.DocumentStatus = "new"
	ref1.Save("DocumentStatus")
	erefId1 := common.Encrypt(ref1.ID)

	ref2 := model.DummyPurchaseOrder()
	ref2.DocumentStatus = "cancelled"
	ref2.Save("DocumentStatus")
	erefId2 := common.Encrypt(ref2.ID)

	supplier1 := model.DummyPartnership()
	supplier1.PartnershipType = "supplier"
	supplier1.Save("PartnershipType")
	esupplierID1 := common.Encrypt(supplier1.ID)

	supplier2 := model.DummyPartnership()
	supplier2.PartnershipType = "customer"
	supplier2.Save("PartnershipType")
	esupplierID2 := common.Encrypt(supplier2.ID)

	supplier3 := model.DummyPartnership()
	supplier3.FullName = "Subagja"
	supplier3.PartnershipType = "supplier"
	supplier3.Save("PartnershipType", "FullName")
	esupplierID3 := common.Encrypt(supplier3.ID)

	iv := model.DummyItemVariant()
	iv.IsArchived = 1
	iv.Save()
	eiv := common.Encrypt(iv.ID)

	iv2 := model.DummyItemVariant()
	iv2.IsArchived = 1
	iv2.Save()
	eiv2 := common.Encrypt(iv2.ID)

	iv3 := model.DummyItemVariant()
	iv3.IsArchived = 0
	iv3.IsDeleted = 0
	iv3.Save()
	eiv3 := common.Encrypt(iv3.ID)

	var data = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{
			"recognition_date": time.Now(),
			"reference_id":     erefId1,
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID1,
			"tax":              10,
			"discount":         10,
			"shipment_cost":    5000,
			"note":             "catatan",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eiv,
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusUnprocessableEntity},
		{tester.D{
			"recognition_date": time.Now(),
			"reference_id":     erefId2,
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID1,
			"tax":              10,
			"discount":         10,
			"shipment_cost":    5000,
			"note":             "catatan",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eiv2,
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusUnprocessableEntity},
		{tester.D{
			"recognition_date": time.Now(),
			"reference_id":     erefId2,
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID2,
			"tax":              10,
			"discount":         10,
			"shipment_cost":    5000,
			"note":             "catatan",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eiv,
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusUnprocessableEntity},
		{tester.D{
			"recognition_date": time.Now(),
			"reference_id":     erefId2,
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID1,
			"tax":              10,
			"discount":         10,
			"shipment_cost":    5000,
			"note":             "catatan",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": "abcdefghij",
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusUnprocessableEntity},
		{tester.D{
			"recognition_date": time.Now(),
			"reference_id":     erefId2,
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID1,
			"tax":              10,
			"discount":         10,
			"shipment_cost":    5000,
			"note":             "catatan",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": "65536000",
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusUnprocessableEntity},
		{tester.D{
			"recognition_date": time.Now(),
			"reference_id":     erefId2,
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID1,
			"tax":              10,
			"is_percentage":    1,
			"discount":         0,
			"shipment_cost":    5000,
			"note":             "catatan",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eiv3,
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusOK},
		{tester.D{
			"recognition_date": time.Now(),
			"reference_id":     erefId2,
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID1,
			"tax":              10,
			"is_percentage":    1,
			"discount":         0,
			"discount_amount":  1000,
			"shipment_cost":    5000,
			"note":             "catatan",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eiv3,
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusOK},
		{tester.D{
			"recognition_date": time.Now(),
			"reference_id":     erefId2,
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID1,
			"tax":              10,
			"is_percentage":    0,
			"discount":         90,
			"discount_amount":  0,
			"shipment_cost":    5000,
			"note":             "catatan",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eiv,
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusUnprocessableEntity},
		{tester.D{
			"recognition_date": time.Now(),
			"reference_id":     erefId2,
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID1,
			"tax":              10,
			"is_percentage":    0,
			"discount":         0,
			"discount_amount":  0,
			"shipment_cost":    5000,
			"note":             "catatan",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eiv3,
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusOK},
		{tester.D{
			"recognition_date": time.Now(),
			"reference_id":     erefId2,
			"auto_invoiced":    0,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID1,
			"tax":              10,
			"is_percentage":    0,
			"discount":         0,
			"discount_amount":  0,
			"shipment_cost":    5000,
			"note":             "catatan",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eiv3,
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusOK},
		{tester.D{
			"recognition_date": time.Now(),
			"reference_id":     erefId2,
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID3,
			"tax":              10,
			"is_percentage":    0,
			"discount":         0,
			"discount_amount":  0,
			"shipment_cost":    5000,
			"note":             "catatan",
			"purchase_order_items": []tester.D{
				{
					"item_variant_id": eiv3,
					"quantity":        12,
					"unit_price":      1000,
					"discount":        1,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusOK},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.POST("/v1/purchase-order").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
	var porder model.PurchaseOrder
	err := orm.NewOrm().Raw("SELECT po.* FROM purchase_order po "+
		"WHERE po.supplier_id = ?", supplier3.ID).QueryRow(&porder)
	assert.Equal(t, "active", porder.DocumentStatus, "Harus sama, actual: "+fmt.Sprint(porder.DocumentStatus, err))
	assert.Equal(t, "active", porder.InvoiceStatus, "Harus sama, actual: "+fmt.Sprint(porder.InvoiceStatus, err))

	var pinvoice model.PurchaseInvoice
	err = orm.NewOrm().Raw("SELECT pi.* FROM purchase_invoice pi "+
		"WHERE pi.purchase_order_id = ?", porder.ID).QueryRow(&pinvoice)
	assert.Equal(t, "new", pinvoice.DocumentStatus, "Harus sama, actual: "+fmt.Sprint(pinvoice.DocumentStatus, err))
}

func TestHandler_URLMappingPUT_Cancel(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	//purchase order
	po := model.DummyPurchaseOrder()
	po.IsDeleted = 0
	po.DocumentStatus = "new"
	po.Save("IsDeleted", "DocumentStatus")
	epo := common.Encrypt(po.ID)

	for a := 0; a < 2; a++ {
		// purchaceinvoice
		invoice := model.DummyPurchaseInvoice()
		invoice.PurchaseOrder = po
		invoice.IsDeleted = 0
		invoice.DocumentStatus = "active"
		invoice.Save("PurchaseOrder", "IsDeleted", "DocumentStatus")

		// financeexpense
		finance := model.DummyFinanceExpense()
		finance.RefID = uint64(po.ID)
		finance.IsDeleted = 0
		finance.Save("RefID", "IsDeleted")

		receiving := model.DummyWorkorderReceiving()
		receiving.PurchaseOrder = po
		receiving.IsDeleted = 0
		receiving.Save("PurchaseOrder", "IsDeleted")

		iv := model.DummyItemVariantStock()

		//itemvariantstocklog
		ivslog := model.DummyItemVariantStockLog()
		ivslog.RefType = "workorder_receiving"
		ivslog.RefID = uint64(receiving.ID)
		ivslog.ItemVariantStock = iv
		ivslog.LogType = "in"
		ivslog.Save("RefType", "RefID", "LogType", "ItemVariantStock")
	}

	var data = []struct {
		req      tester.D
		expected int
	}{
		//sukses
		{tester.D{"cancelled_note": "note1"}, http.StatusOK},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/purchase-order/"+epo+"/cancel").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

	// jika document_status ="cancelled"
	po.DocumentStatus = "cancelled"
	po.IsDeleted = 0
	po.Save("DocumentStatus", "IsDeleted")
	data = []struct {
		req      tester.D
		expected int
	}{
		//sukses
		{tester.D{"cancelled_note": "note1"}, http.StatusUnprocessableEntity},
	}
	for _, tes := range data {
		ng.PUT("/v1/purchase-order/"+epo+"/cancel").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

	// not found
	po.IsDeleted = 1
	po.Save("IsDeleted")
	data = []struct {
		req      tester.D
		expected int
	}{
		//sukses
		{tester.D{"cancelled_note": "note1"}, http.StatusNotFound},
	}

	for _, tes := range data {
		ng.PUT("/v1/purchase-order/"+epo+"/cancel").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

func TestHandler_URLMappingPUT2(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	DummyPO := model.DummyPurchaseOrderItem()
	DummyPO.PurchaseOrder.IsDeleted = 0
	DummyPO.PurchaseOrder.DocumentStatus = "new"
	DummyPO.PurchaseOrder.TotalCharge = 10000
	DummyPO.PurchaseOrder.Discount = 10
	DummyPO.PurchaseOrder.DiscountAmount = 900
	DummyPO.PurchaseOrder.Tax = 10
	DummyPO.PurchaseOrder.TaxAmount = 810
	DummyPO.PurchaseOrder.ShipmentCost = 1090
	DummyPO.PurchaseOrder.Save()
	DummyPO.Discount = 10
	DummyPO.Quantity = 10
	DummyPO.UnitPrice = 1000
	DummyPO.Subtotal = 9000
	DummyPO.Save()
	epoid := common.Encrypt(DummyPO.PurchaseOrder.ID)

	DummyIV := model.DummyItemVariantPrice()
	DummyIV.ItemVariant.IsArchived = 0
	DummyIV.ItemVariant.IsDeleted = 0
	DummyIV.ItemVariant.Save()
	eitemvar := common.Encrypt(DummyIV.ItemVariant.ID)

	DummyPartner := model.DummyPartnership()
	DummyPartner.PartnershipType = "supplier"
	DummyPartner.IsDeleted = 0
	DummyPartner.IsArchived = 0
	DummyPartner.TotalExpenditure = 100000
	DummyPartner.Save()
	esupplierID := common.Encrypt(DummyPartner.ID)

	var data = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{
			"recognition_date": time.Now(),
			"auto_invoiced":    1,
			"eta_date":         time.Now(),
			"supplier_id":      esupplierID,
			"tax":              10,
			"is_percentage":    1,
			"discount":         0,
			"shipment_cost":    5000,
			"note":             "is percentage 1 tetapi discount 0",
			"purchase_order_items": []tester.D{
				{
					"id":              common.Encrypt(DummyPO.ID),
					"item_variant_id": eitemvar,
					"quantity":        5,
					"unit_price":      2000,
					"discount":        20,
					"note":            "catatan PO Items",
				},
			},
		}, http.StatusOK},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/purchase-order/"+epoid).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// test untuk mengecek data
func TestCancelPurchaseCheckData(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	//purchase order
	po := model.DummyPurchaseOrder()
	po.IsDeleted = 0
	po.DocumentStatus = "new"
	po.ReceivingStatus = "active"
	po.InvoiceStatus = "active"
	po.Save("IsDeleted", "DocumentStatus", "ReceivingStatus", "InvoiceStatus")
	epo := common.Encrypt(po.ID)

	for a := 0; a < 2; a++ {
		// purchaceinvoice
		invoice := model.DummyPurchaseInvoice()
		invoice.PurchaseOrder = po
		invoice.IsDeleted = 0
		invoice.DocumentStatus = "active"
		invoice.Save("PurchaseOrder", "IsDeleted", "DocumentStatus")

		// financeexpense
		finance := model.DummyFinanceExpense()
		finance.RefID = uint64(po.ID)
		finance.IsDeleted = 0
		finance.Save("RefID", "IsDeleted")

		receiving := model.DummyWorkorderReceiving()
		receiving.PurchaseOrder = po
		receiving.IsDeleted = 0
		receiving.Save("PurchaseOrder", "IsDeleted")

		iv := model.DummyItemVariantStock()

		//itemvariantstocklog
		ivslog := model.DummyItemVariantStockLog()
		ivslog.RefType = "workorder_receiving"
		ivslog.RefID = uint64(receiving.ID)
		ivslog.ItemVariantStock = iv
		ivslog.LogType = "in"
		ivslog.Save("RefType", "RefID", "LogType", "ItemVariantStock")
	}

	var data = []struct {
		req      tester.D
		expected int
	}{
		//sukses
		{tester.D{"cancelled_note": "note1"}, http.StatusOK},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/purchase-order/"+epo+"/cancel").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

	o := orm.NewOrm()
	// get data purchase invoice by Po
	var purchaseInvoices []model.PurchaseInvoice
	o.Raw("select * from purchase_invoice where purchase_order_id = ? and is_deleted = ?", po.ID, 0).QueryRows(&purchaseInvoices)
	for _, pi := range purchaseInvoices {
		// get data finance expense berdasarkan purchase invoice
		var finances []model.FinanceExpense
		o.Raw("select * from finance_expense where ref_id = ? and ref_type = 'purchase_invoice' ", pi.ID).QueryRows(&finances)
		for _, fx := range finances {
			assert.Equal(t, int8(1), fx.IsDeleted)
		}
	}

}

// TestCancelPurchaseCheckDataIfInvoiceStatusIsNew untuk mengecek data apabila inv
func TestCancelPurchaseCheckDataIfInvoiceStatusIsNew(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	//purchase order
	po := model.DummyPurchaseOrder()
	po.IsDeleted = 0
	po.DocumentStatus = "new"
	po.ReceivingStatus = "new"
	po.InvoiceStatus = "active"
	po.Save("IsDeleted", "DocumentStatus", "ReceivingStatus", "InvoiceStatus")
	epo := common.Encrypt(po.ID)

	for a := 0; a < 2; a++ {
		// purchaceinvoice
		invoice := model.DummyPurchaseInvoice()
		invoice.PurchaseOrder = po
		invoice.IsDeleted = 0
		invoice.DocumentStatus = "active"
		invoice.Save("PurchaseOrder", "IsDeleted", "DocumentStatus")

		// financeexpense
		finance := model.DummyFinanceExpense()
		finance.RefID = uint64(po.ID)
		finance.IsDeleted = 0
		finance.Save("RefID", "IsDeleted")

		receiving := model.DummyWorkorderReceiving()
		receiving.PurchaseOrder = po
		receiving.IsDeleted = 0
		receiving.Save("PurchaseOrder", "IsDeleted")

		iv := model.DummyItemVariantStock()

		//itemvariantstocklog
		ivslog := model.DummyItemVariantStockLog()
		ivslog.RefType = "workorder_receiving"
		ivslog.RefID = uint64(receiving.ID)
		ivslog.ItemVariantStock = iv
		ivslog.LogType = "in"
		ivslog.Save("RefType", "RefID", "LogType", "ItemVariantStock")
	}

	var data = []struct {
		req      tester.D
		expected int
	}{
		//sukses
		{tester.D{"cancelled_note": "note1"}, http.StatusOK},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/purchase-order/"+epo+"/cancel").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

	o := orm.NewOrm()
	// get datareceiving by Po
	var receiving []model.WorkorderFulfillment
	o.Raw("select * from workorder_receiving where purchase_order_id = ? and is_deleted = ?", po.ID, 0).QueryRows(&receiving)
	for _, r := range receiving {
		assert.Equal(t, int8(0), r.IsDeleted)
	}
}
