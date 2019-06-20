// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchaseReturn_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/purchase_return"
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

// TestHandler_URLMappingPOSTCreatePurchaseReturn Sukses Semua aksi berhasil
func TestHandler_URLMappingPOSTCreatePurchaseReturn(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive1 := model.DummyWorkorderReceiving()
	woReceive1.PurchaseOrder = dporder
	woReceive1.IsDeleted = 0
	woReceive1.DocumentStatus = "finished"
	woReceive1.Save()
	woReceiveItem1 := model.DummyWorkorderReceivingItem()
	woReceiveItem1.Quantity = float32(45)
	woReceiveItem1.WorkorderReceiving = woReceive1
	woReceiveItem1.PurchaseOrderItem = dpOrderItem
	woReceiveItem1.Save()

	woReceive2 := model.DummyWorkorderReceiving()
	woReceive2.PurchaseOrder = dporder
	woReceive2.IsDeleted = 0
	woReceive2.DocumentStatus = "active"
	woReceive2.Save()
	woReceiveItem2 := model.DummyWorkorderReceivingItem()
	woReceiveItem2.Quantity = float32(45)
	woReceiveItem2.WorkorderReceiving = woReceive2
	woReceiveItem2.PurchaseOrderItem = dpOrderItem
	woReceiveItem2.Save()

	// setting body
	scenario := tester.D{
		"recognition_date":  time.Now(),
		"purchase_order_id": common.Encrypt(dporder.ID),
		"purchase_return_item": []tester.D{
			{
				"purchase_order_item_id": common.Encrypt(dpOrderItem.ID),
				"quantity":               float32(45),
				"unit_price":             dpOrderItem.UnitPrice,
				"discount":               dpOrderItem.Discount,
				"subtotal":               dpOrderItem.UnitPrice * float64(45),
				"note":                   "Whaaaaaat....",
			},
		},
		"tax_amount":      dporder.TaxAmount,
		"discount_amount": dporder.DiscountAmount,
		"shipment_cost":   dporder.ShipmentCost,
		"total_amount":    dpOrderItem.UnitPrice * float64(45),
		"note":            "Weew...",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/purchase-return").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPOSTCreatePurchaseReturn1 Sukses Semua aksi berhasil
func TestHandler_URLMappingPOSTCreatePurchaseReturn1(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	// setting body
	scenario := tester.D{
		"recognition_date":  time.Now(),
		"purchase_order_id": common.Encrypt(dporder.ID),
		"purchase_return_item": []tester.D{
			{
				"purchase_order_item_id": common.Encrypt(dpOrderItem.ID),
				"quantity":               float32(40),
				"unit_price":             dpOrderItem.UnitPrice,
				"discount":               dpOrderItem.Discount,
				"subtotal":               dpOrderItem.UnitPrice * float64(40),
				"note":                   "Whaaaaaat....",
			},
		},
		"tax_amount":      dporder.TaxAmount,
		"discount_amount": dporder.DiscountAmount,
		"shipment_cost":   dporder.ShipmentCost,
		"total_amount":    dpOrderItem.UnitPrice * float64(40),
		"note":            "Weew...",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/purchase-return").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPOSTCreatePurchaseReturnGagal1 Test validation
func TestHandler_URLMappingPOSTCreatePurchaseReturnGagal1(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	// setting body
	scenario := tester.D{
		"recognition_date":  time.Now(),
		"purchase_order_id": common.Encrypt("abcdefghijklmnopqrstuvwxyz"),
		"purchase_return_item": []tester.D{
			{
				"purchase_order_item_id": common.Encrypt(dpOrderItem.ID),
				"quantity":               float32(40),
				"unit_price":             dpOrderItem.UnitPrice,
				"discount":               dpOrderItem.Discount,
				"subtotal":               dpOrderItem.UnitPrice * float64(40),
				"note":                   "Whaaaaaat....",
			},
		},
		"tax_amount":      dporder.TaxAmount,
		"discount_amount": dporder.DiscountAmount,
		"shipment_cost":   dporder.ShipmentCost,
		"total_amount":    dpOrderItem.UnitPrice * float64(40),
		"note":            "Weew...",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/purchase-return").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPOSTCreatePurchaseReturnGagal2 Test PO Id is not found
func TestHandler_URLMappingPOSTCreatePurchaseReturnGagal2(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.CreatedBy = model.DummyUser()
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()
	orm.NewOrm().Raw("DELETE FROM purchase_oder").Exec()
	// setting body
	scenario := tester.D{
		"recognition_date":  time.Now(),
		"purchase_order_id": common.Encrypt("65536"),
		"purchase_return_item": []tester.D{
			{
				"purchase_order_item_id": common.Encrypt(dpOrderItem.ID),
				"quantity":               float32(40),
				"unit_price":             dpOrderItem.UnitPrice,
				"discount":               dpOrderItem.Discount,
				"subtotal":               dpOrderItem.UnitPrice * float64(40),
				"note":                   "Whaaaaaat....",
			},
		},
		"tax_amount":      dporder.TaxAmount,
		"discount_amount": dporder.DiscountAmount,
		"shipment_cost":   dporder.ShipmentCost,
		"total_amount":    dpOrderItem.UnitPrice * float64(40),
		"note":            "Weew...",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/purchase-return").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPOSTCreatePurchaseReturnGagal3 Test POItems Id is invalid
func TestHandler_URLMappingPOSTCreatePurchaseReturnGagal3(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.CreatedBy = model.DummyUser()
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()
	// setting body
	scenario := tester.D{
		"recognition_date":  time.Now(),
		"purchase_order_id": common.Encrypt("65536"),
		"purchase_return_item": []tester.D{
			{
				"purchase_order_item_id": common.Encrypt("asdfghjkl"),
				"quantity":               float32(40),
				"unit_price":             dpOrderItem.UnitPrice,
				"discount":               dpOrderItem.Discount,
				"subtotal":               dpOrderItem.UnitPrice * float64(40),
				"note":                   "Whaaaaaat....",
			},
		},
		"tax_amount":      dporder.TaxAmount,
		"discount_amount": dporder.DiscountAmount,
		"shipment_cost":   dporder.ShipmentCost,
		"total_amount":    dpOrderItem.UnitPrice * float64(40),
		"note":            "Weew...",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/purchase-return").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPOSTCreatePurchaseReturnGagal4 Test POItems Id is not found
func TestHandler_URLMappingPOSTCreatePurchaseReturnGagal4(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.CreatedBy = model.DummyUser()
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()
	orm.NewOrm().Raw("DELETE FROM purchase_order_item").Exec()
	// setting body
	scenario := tester.D{
		"recognition_date":  time.Now(),
		"purchase_order_id": common.Encrypt("65536"),
		"purchase_return_item": []tester.D{
			{
				"purchase_order_item_id": common.Encrypt("65536"),
				"quantity":               float32(40),
				"unit_price":             dpOrderItem.UnitPrice,
				"discount":               dpOrderItem.Discount,
				"subtotal":               dpOrderItem.UnitPrice * float64(40),
				"note":                   "Whaaaaaat....",
			},
		},
		"tax_amount":      dporder.TaxAmount,
		"discount_amount": dporder.DiscountAmount,
		"shipment_cost":   dporder.ShipmentCost,
		"total_amount":    dpOrderItem.UnitPrice * float64(40),
		"note":            "Weew...",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/purchase-return").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPOSTCreatePurchaseReturnGagal5 Test POItems Qty is over
func TestHandler_URLMappingPOSTCreatePurchaseReturnGagal5(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.CreatedBy = model.DummyUser()
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()
	// setting body
	scenario := tester.D{
		"recognition_date":  time.Now(),
		"purchase_order_id": common.Encrypt(dporder.ID),
		"purchase_return_item": []tester.D{
			{
				"purchase_order_item_id": common.Encrypt(dpOrderItem.ID),
				"quantity":               float32(50),
				"unit_price":             dpOrderItem.UnitPrice,
				"discount":               dpOrderItem.Discount,
				"subtotal":               dpOrderItem.UnitPrice * float64(50),
				"note":                   "Whaaaaaat....",
			},
		},
		"tax_amount":      dporder.TaxAmount,
		"discount_amount": dporder.DiscountAmount,
		"shipment_cost":   dporder.ShipmentCost,
		"total_amount":    dpOrderItem.UnitPrice * float64(50),
		"note":            "Weew...",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/purchase-return").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

func TestHandler_URLMappingGETALL(t *testing.T) {
	// hapus semua data dari tabel purchase order
	o := orm.NewOrm()
	o.Raw("DELETE FROM purchase_order").Exec()

	preturn := model.DummyPurchaseReturn()
	preturnItem := model.DummyPurchaseReturnItem()
	preturnItem.PurchaseReturn = preturn
	preturn.Save()

	// proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/purchase-return"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/purchase-order", "GET"))
	})
}

func TestHandler_URLMappingGETDetail(t *testing.T) {
	// hapus semua data dari tabel purchase order
	o := orm.NewOrm()
	o.Raw("DELETE FROM purchase_return").Exec()

	preturn := model.DummyPurchaseReturn()
	preturn.IsDeleted = 0
	preturn.Save()
	preturnItem := model.DummyPurchaseReturnItem()
	preturnItem.PurchaseReturn = preturn
	preturn.Save()

	// proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/purchase-return/" + common.Encrypt(preturn.ID)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/purchase-return", "GET"))
	})
}

func TestHandler_URLMappingGETDetailErrNotFound(t *testing.T) {
	// hapus semua data dari tabel purchase order
	o := orm.NewOrm()
	o.Raw("DELETE FROM purchase_return").Exec()

	// proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/purchase-return/65536"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusNotFound, res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/purchase-return", "GET"))
	})
}

// TestHandler_URLMappingPOSTCreatePurchaseReturn1 Sukses Semua aksi berhasil
func TestHandler_URLMappingPUTEditPurchaseReturn1(t *testing.T) {
	test.DataCleanUp()
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.Discount = 0
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	preturn := model.PurchaseReturn{
		PurchaseOrder:   &model.PurchaseOrder{ID: dporder.ID},
		RecognitionDate: time.Now(),
		Code:            common.RandomStr(10),
		TotalAmount:     float64(50000),
		Note:            "Wow...",
		IsDeleted:       int8(0),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		PurchaseReturnItems: []*model.PurchaseReturnItem{
			{
				PurchaseOrderItem: &model.PurchaseOrderItem{ID: dpOrderItem.ID},
				Quantity:          float32(10),
				Note:              common.RandomStr(15),
			},
		},
		DocumentStatus: "new",
	}

	pret, _ := purchaseReturn.CreatePurchaseReturn(&preturn)

	// setting body
	scenario := tester.D{
		"recognition_date":  time.Now(),
		"purchase_order_id": common.Encrypt(dporder.ID),
		"purchase_return_item": []tester.D{
			{
				"purchase_order_item_id": common.Encrypt(dpOrderItem.ID),
				"quantity":               float32(40),
				"unit_price":             dpOrderItem.UnitPrice,
				"discount":               dpOrderItem.Discount,
				"subtotal":               dpOrderItem.UnitPrice * float64(40),
				"note":                   "Whaaaaaat....",
			},
		},
		"tax_amount":      dporder.TaxAmount,
		"discount_amount": dporder.DiscountAmount,
		"shipment_cost":   dporder.ShipmentCost,
		"total_amount":    dpOrderItem.UnitPrice * float64(40),
		"note":            "Weew...",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/purchase-return/"+common.Encrypt(pret.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPOSTCreatePurchaseReturnGagal1 Document status sudah active
func TestHandler_URLMappingPUTEditPurchaseReturnGagal1(t *testing.T) {
	test.DataCleanUp()
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.Discount = 0
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	preturn := model.PurchaseReturn{
		PurchaseOrder:   &model.PurchaseOrder{ID: dporder.ID},
		RecognitionDate: time.Now(),
		Code:            common.RandomStr(10),
		TotalAmount:     float64(50000),
		Note:            "Wow...",
		IsDeleted:       int8(0),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		PurchaseReturnItems: []*model.PurchaseReturnItem{
			{
				PurchaseOrderItem: &model.PurchaseOrderItem{ID: dpOrderItem.ID},
				Quantity:          float32(10),
				Note:              common.RandomStr(15),
			},
		},
		DocumentStatus: "new",
	}

	pret, _ := purchaseReturn.CreatePurchaseReturn(&preturn)
	pret.DocumentStatus = "active"
	pret.Save()

	// setting body
	scenario := tester.D{
		"recognition_date":  time.Now(),
		"purchase_order_id": common.Encrypt(dporder.ID),
		"purchase_return_item": []tester.D{
			{
				"purchase_order_item_id": common.Encrypt(dpOrderItem.ID),
				"quantity":               float32(40),
				"unit_price":             dpOrderItem.UnitPrice,
				"discount":               dpOrderItem.Discount,
				"subtotal":               dpOrderItem.UnitPrice * float64(40),
				"note":                   "Whaaaaaat....",
			},
		},
		"tax_amount":      dporder.TaxAmount,
		"discount_amount": dporder.DiscountAmount,
		"shipment_cost":   dporder.ShipmentCost,
		"total_amount":    dpOrderItem.UnitPrice * float64(40),
		"note":            "Weew...",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/purchase-return/"+common.Encrypt(pret.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPOSTCreatePurchaseReturnGagal2 Test validation
func TestHandler_URLMappingPUTEditPurchaseReturnGagal2(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	preturn := model.PurchaseReturn{
		PurchaseOrder:   &model.PurchaseOrder{ID: dporder.ID},
		RecognitionDate: time.Now(),
		Code:            common.RandomStr(10),
		TotalAmount:     float64(50000),
		Note:            "Wow...",
		IsDeleted:       int8(0),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		PurchaseReturnItems: []*model.PurchaseReturnItem{
			{
				PurchaseOrderItem: &model.PurchaseOrderItem{ID: dpOrderItem.ID},
				Quantity:          float32(10),
				Note:              common.RandomStr(15),
			},
		},
		DocumentStatus: "new",
	}

	pret, _ := purchaseReturn.CreatePurchaseReturn(&preturn)

	// setting body
	scenario := tester.D{
		"recognition_date":  time.Time{},
		"purchase_order_id": common.Encrypt("abcdefghijklmnopqrstuvwxyz"),
		"purchase_return_item": []tester.D{
			{
				"purchase_order_item_id": common.Encrypt(dpOrderItem.ID),
				"quantity":               float32(40),
				"unit_price":             dpOrderItem.UnitPrice,
				"discount":               dpOrderItem.Discount,
				"subtotal":               dpOrderItem.UnitPrice * float64(40),
				"note":                   "Whaaaaaat....",
			},
		},
		"tax_amount":      dporder.TaxAmount,
		"discount_amount": dporder.DiscountAmount,
		"shipment_cost":   dporder.ShipmentCost,
		"total_amount":    dpOrderItem.UnitPrice * float64(40),
		"note":            "Weew...",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/purchase-return/"+common.Encrypt(pret.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPOSTCreatePurchaseReturnGagal3 Purchase order tidak ditemukan
func TestHandler_URLMappingPUTEditPurchaseReturnGagal3(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	preturn := model.PurchaseReturn{
		PurchaseOrder:   &model.PurchaseOrder{ID: dporder.ID},
		RecognitionDate: time.Now(),
		Code:            common.RandomStr(10),
		TotalAmount:     float64(50000),
		Note:            "Wow...",
		IsDeleted:       int8(0),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		PurchaseReturnItems: []*model.PurchaseReturnItem{
			{
				PurchaseOrderItem: &model.PurchaseOrderItem{ID: dpOrderItem.ID},
				Quantity:          float32(10),
				Note:              common.RandomStr(15),
			},
		},
		DocumentStatus: "new",
	}
	pret, _ := purchaseReturn.CreatePurchaseReturn(&preturn)
	orm.NewOrm().Raw("DELETE FROM purchase_oder").Exec()

	// setting body
	scenario := tester.D{
		"recognition_date":  time.Time{},
		"purchase_order_id": common.Encrypt("65536"),
		"purchase_return_item": []tester.D{
			{
				"purchase_order_item_id": common.Encrypt(dpOrderItem.ID),
				"quantity":               float32(40),
				"unit_price":             dpOrderItem.UnitPrice,
				"discount":               dpOrderItem.Discount,
				"subtotal":               dpOrderItem.UnitPrice * float64(40),
				"note":                   "Whaaaaaat....",
			},
		},
		"tax_amount":      dporder.TaxAmount,
		"discount_amount": dporder.DiscountAmount,
		"shipment_cost":   dporder.ShipmentCost,
		"total_amount":    dpOrderItem.UnitPrice * float64(40),
		"note":            "Weew...",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/purchase-return/"+common.Encrypt(pret.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPOSTCreatePurchaseReturnGagal4 Purchase order item invalid
func TestHandler_URLMappingPUTEditPurchaseReturnGagal4(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	preturn := model.PurchaseReturn{
		PurchaseOrder:   &model.PurchaseOrder{ID: dporder.ID},
		RecognitionDate: time.Now(),
		Code:            common.RandomStr(10),
		TotalAmount:     float64(50000),
		Note:            "Wow...",
		IsDeleted:       int8(0),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		PurchaseReturnItems: []*model.PurchaseReturnItem{
			{
				PurchaseOrderItem: &model.PurchaseOrderItem{ID: dpOrderItem.ID},
				Quantity:          float32(10),
				Note:              common.RandomStr(15),
			},
		},
		DocumentStatus: "new",
	}
	pret, _ := purchaseReturn.CreatePurchaseReturn(&preturn)

	// setting body
	scenario := tester.D{
		"recognition_date":  time.Time{},
		"purchase_order_id": common.Encrypt(dporder.ID),
		"purchase_return_item": []tester.D{
			{
				"purchase_order_item_id": common.Encrypt("abcdefghijkl"),
				"quantity":               float32(40),
				"unit_price":             dpOrderItem.UnitPrice,
				"discount":               dpOrderItem.Discount,
				"subtotal":               dpOrderItem.UnitPrice * float64(40),
				"note":                   "Whaaaaaat....",
			},
		},
		"tax_amount":      dporder.TaxAmount,
		"discount_amount": dporder.DiscountAmount,
		"shipment_cost":   dporder.ShipmentCost,
		"total_amount":    dpOrderItem.UnitPrice * float64(40),
		"note":            "Weew...",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/purchase-return/"+common.Encrypt(pret.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPOSTCreatePurchaseReturnGagal5 Purchase order item not found
func TestHandler_URLMappingPUTEditPurchaseReturnGagal5(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	preturn := model.PurchaseReturn{
		PurchaseOrder:   &model.PurchaseOrder{ID: dporder.ID},
		RecognitionDate: time.Now(),
		Code:            common.RandomStr(10),
		TotalAmount:     float64(50000),
		Note:            "Wow...",
		IsDeleted:       int8(0),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		PurchaseReturnItems: []*model.PurchaseReturnItem{
			{
				PurchaseOrderItem: &model.PurchaseOrderItem{ID: dpOrderItem.ID},
				Quantity:          float32(10),
				Note:              common.RandomStr(15),
			},
		},
		DocumentStatus: "new",
	}
	pret, _ := purchaseReturn.CreatePurchaseReturn(&preturn)
	orm.NewOrm().Raw("DELETE FROM purchase_order_item").Exec()
	// setting body
	scenario := tester.D{
		"recognition_date":  time.Time{},
		"purchase_order_id": common.Encrypt(dporder.ID),
		"purchase_return_item": []tester.D{
			{
				"purchase_order_item_id": common.Encrypt("65536"),
				"quantity":               float32(40),
				"unit_price":             dpOrderItem.UnitPrice,
				"discount":               dpOrderItem.Discount,
				"subtotal":               dpOrderItem.UnitPrice * float64(40),
				"note":                   "Whaaaaaat....",
			},
		},
		"tax_amount":      dporder.TaxAmount,
		"discount_amount": dporder.DiscountAmount,
		"shipment_cost":   dporder.ShipmentCost,
		"total_amount":    dpOrderItem.UnitPrice * float64(40),
		"note":            "Weew...",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/purchase-return/"+common.Encrypt(pret.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPOSTCreatePurchaseReturnGagal6 Quantity too much input
func TestHandler_URLMappingPUTEditPurchaseReturnGagal6(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	preturn := model.PurchaseReturn{
		PurchaseOrder:   &model.PurchaseOrder{ID: dporder.ID},
		RecognitionDate: time.Now(),
		Code:            common.RandomStr(10),
		TotalAmount:     float64(50000),
		Note:            "Wow...",
		IsDeleted:       int8(0),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		PurchaseReturnItems: []*model.PurchaseReturnItem{
			{
				PurchaseOrderItem: &model.PurchaseOrderItem{ID: dpOrderItem.ID},
				Quantity:          float32(10),
				Note:              common.RandomStr(15),
			},
		},
		DocumentStatus: "new",
	}
	pret, _ := purchaseReturn.CreatePurchaseReturn(&preturn)
	// setting body
	scenario := tester.D{
		"recognition_date":  time.Time{},
		"purchase_order_id": common.Encrypt(dporder.ID),
		"purchase_return_item": []tester.D{
			{
				"purchase_order_item_id": common.Encrypt(dpOrderItem.ID),
				"quantity":               float32(50),
				"unit_price":             dpOrderItem.UnitPrice,
				"discount":               dpOrderItem.Discount,
				"subtotal":               dpOrderItem.UnitPrice * float64(50),
				"note":                   "Whaaaaaat....",
			},
		},
		"tax_amount":      dporder.TaxAmount,
		"discount_amount": dporder.DiscountAmount,
		"shipment_cost":   dporder.ShipmentCost,
		"total_amount":    dpOrderItem.UnitPrice * float64(40),
		"note":            "Weew...",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/purchase-return/"+common.Encrypt(pret.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPOSTCreatePurchaseReturnGagal7 Purchase return not found
func TestHandler_URLMappingPUTEditPurchaseReturnGagal7(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	preturn := model.PurchaseReturn{
		PurchaseOrder:   &model.PurchaseOrder{ID: dporder.ID},
		RecognitionDate: time.Now(),
		Code:            common.RandomStr(10),
		TotalAmount:     float64(50000),
		Note:            "Wow...",
		IsDeleted:       int8(0),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		PurchaseReturnItems: []*model.PurchaseReturnItem{
			{
				PurchaseOrderItem: &model.PurchaseOrderItem{ID: dpOrderItem.ID},
				Quantity:          float32(10),
				Note:              common.RandomStr(15),
			},
		},
		DocumentStatus: "new",
	}
	pret, _ := purchaseReturn.CreatePurchaseReturn(&preturn)
	pret.Delete()
	// setting body
	scenario := tester.D{
		"recognition_date":  time.Time{},
		"purchase_order_id": common.Encrypt(dporder.ID),
		"purchase_return_item": []tester.D{
			{
				"purchase_order_item_id": common.Encrypt(dpOrderItem.ID),
				"quantity":               float32(40),
				"unit_price":             dpOrderItem.UnitPrice,
				"discount":               dpOrderItem.Discount,
				"subtotal":               dpOrderItem.UnitPrice * float64(40),
				"note":                   "Whaaaaaat....",
			},
		},
		"tax_amount":      dporder.TaxAmount,
		"discount_amount": dporder.DiscountAmount,
		"shipment_cost":   dporder.ShipmentCost,
		"total_amount":    dpOrderItem.UnitPrice * float64(40),
		"note":            "Weew...",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/purchase-return/"+common.Encrypt("65536")).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusNotFound, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPUTCancelPurchaseReturn1 Sukses
func TestHandler_URLMappingPUTCancelPurchaseReturn1(t *testing.T) {
	test.DataCleanUp()
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.Discount = 0
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	preturn := model.PurchaseReturn{
		PurchaseOrder:   &model.PurchaseOrder{ID: dporder.ID},
		RecognitionDate: time.Now(),
		Code:            common.RandomStr(10),
		TotalAmount:     float64(50000),
		Note:            "Wow...",
		IsDeleted:       int8(0),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		PurchaseReturnItems: []*model.PurchaseReturnItem{
			{
				PurchaseOrderItem: &model.PurchaseOrderItem{ID: dpOrderItem.ID},
				Quantity:          float32(10),
				Note:              common.RandomStr(15),
			},
		},
		DocumentStatus: "new",
	}

	pret, _ := purchaseReturn.CreatePurchaseReturn(&preturn)

	scenario := tester.D{}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/purchase-return/"+common.Encrypt(pret.ID)+"/cancel").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPUTCancelPurchaseReturnError1 Gagal document status already cancelled
func TestHandler_URLMappingPUTCancelPurchaseReturnError1(t *testing.T) {
	test.DataCleanUp()
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.Discount = 0
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	preturn := model.PurchaseReturn{
		PurchaseOrder:   &model.PurchaseOrder{ID: dporder.ID},
		RecognitionDate: time.Now(),
		Code:            common.RandomStr(10),
		TotalAmount:     float64(50000),
		Note:            "Wow...",
		IsDeleted:       int8(0),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		PurchaseReturnItems: []*model.PurchaseReturnItem{
			{
				PurchaseOrderItem: &model.PurchaseOrderItem{ID: dpOrderItem.ID},
				Quantity:          float32(10),
				Note:              common.RandomStr(15),
			},
		},
		DocumentStatus: "new",
	}

	pret, _ := purchaseReturn.CreatePurchaseReturn(&preturn)
	pret.DocumentStatus = "cancelled"
	pret.Save()
	scenario := tester.D{}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/purchase-return/"+common.Encrypt(pret.ID)+"/cancel").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPUTCancelPurchaseReturnError2 Gagal purchase return id not found
func TestHandler_URLMappingPUTCancelPurchaseReturnError2(t *testing.T) {
	test.DataCleanUp()
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.Discount = 0
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	preturn := model.PurchaseReturn{
		PurchaseOrder:   &model.PurchaseOrder{ID: dporder.ID},
		RecognitionDate: time.Now(),
		Code:            common.RandomStr(10),
		TotalAmount:     float64(50000),
		Note:            "Wow...",
		IsDeleted:       int8(0),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		PurchaseReturnItems: []*model.PurchaseReturnItem{
			{
				PurchaseOrderItem: &model.PurchaseOrderItem{ID: dpOrderItem.ID},
				Quantity:          float32(10),
				Note:              common.RandomStr(15),
			},
		},
		DocumentStatus: "new",
	}

	pret, _ := purchaseReturn.CreatePurchaseReturn(&preturn)
	pret.DocumentStatus = "cancelled"
	pret.IsDeleted = 1
	pret.Save()
	scenario := tester.D{}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/purchase-return/"+common.Encrypt(pret.ID)+"/cancel").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusNotFound, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}
