// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fulfillment_test

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

func TestARouting(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	fulfillment := model.DummyWorkorderFulfillment()
	fulfillment.IsDeleted = 0
	fulfillment.Save("IsDeleted")
	efullID := common.Encrypt(fulfillment.ID)

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/fulfillment", "GET", 200},
		{"/v1/fulfillment/" + efullID, "GET", 200},
		{"/v1/fulfillment/" + efullID + "/approve", "PUT", 415},
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

	fulfillment.IsDeleted = 1
	fulfillment.Save("IsDeleted")
	routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/fulfillment", "GET", 200},
		{"/v1/fulfillment/" + efullID, "GET", 404},
		{"/v1/fulfillment/" + efullID + "/approve", "PUT", 404},
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

	// soItem dan sales order fulfillment status = new
	soItem := model.DummySalesOrderItem()
	soItem.Quantity = 20
	soItem.Save("Quantity")
	soItem.SalesOrder.FulfillmentStatus = "new"
	soItem.SalesOrder.DocumentStatus = "active"
	soItem.SalesOrder.Save("FulfillmentStatus", "DocumentStatus")

	esoItem := common.Encrypt(soItem.ID)
	eso := common.Encrypt(soItem.SalesOrder.ID)
	date := time.Now()

	soItemD := model.DummySalesOrderItem()
	soItemD.SalesOrder = soItem.SalesOrder
	soItemD.Quantity = 4
	soItemD.Save("SalesOrder", "Quantity")

	esoItemD := common.Encrypt(soItemD.ID)

	// soItem2 dengan sales order fulfillment status = finished
	soItem2 := model.DummySalesOrderItem()
	soItem2.Quantity = 3
	soItem2.Save("Quantity")
	soItem2.SalesOrder.FulfillmentStatus = "finished"
	soItem2.SalesOrder.Save("FulfillmentStatus")

	esoItem2 := common.Encrypt(soItem2.ID)
	eso2 := common.Encrypt(soItem2.SalesOrder.ID)

	var data = []struct {
		req      tester.D
		expected int
	}{
		//sukses
		{tester.D{"sales_order_id": eso, "priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so1",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 5, "note": "so item"},
			}}, http.StatusOK},
		//jika note item kosong
		{tester.D{"sales_order_id": eso, "priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so2",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 5, "note": ""},
			}}, http.StatusOK},
		// jika note fulfillment kosong
		{tester.D{"sales_order_id": eso, "priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 5, "note": "so item  a"},
			}}, http.StatusOK},

		// gagal
		//apabila priority kosong
		{tester.D{"sales_order_id": eso, "priority": "", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 10, "note": "so item"},
			}}, http.StatusUnprocessableEntity},
		//apabila priority yang dimasukkan bukan routine, rush, emergency
		{tester.D{"sales_order_id": eso, "priority": "important", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 10, "note": "so item"},
			}}, http.StatusUnprocessableEntity},
		//apabila sales order id kosong
		{tester.D{"sales_order_id": "", "priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 10, "note": "so item"},
			}}, http.StatusUnprocessableEntity},
		//apabila sales order id tidak ada dalam database
		{tester.D{"sales_order_id": "9832434", "priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 10, "note": "so item"},
			}}, http.StatusUnprocessableEntity},
		//apabila sales order id gak valid
		{tester.D{"sales_order_id": "xasdsa", "priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 10, "note": "so item"},
			}}, http.StatusUnprocessableEntity},
		// apabila due date kosong
		{tester.D{"sales_order_id": eso, "priority": "routine", "due_date": nil, "shipping_address": "Jalan Kresek Raya", "note": "so",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 10, "note": "so item"},
			}}, http.StatusUnprocessableEntity},
		//apabila shipping address kosong
		{tester.D{"sales_order_id": eso, "priority": "routine", "due_date": date, "shipping_address": "", "note": "so",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 10, "note": "so item"},
			}}, http.StatusUnprocessableEntity},
		// workorder fulfillments kosong
		{tester.D{"sales_order_id": eso, "priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so",
			"workorder_fulfillment_items": nil}, http.StatusUnprocessableEntity},
		//sales order item id kosong
		{tester.D{"sales_order_id": eso, "priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": "", "quantity": 10, "note": "so item"},
			}}, http.StatusUnprocessableEntity},
		// quantity kosong
		{tester.D{"sales_order_id": eso, "priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 0, "note": "so item"},
			}}, http.StatusUnprocessableEntity},
		// apabila sales order, fulfillment status = finish
		{tester.D{"sales_order_id": eso2, "priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem2, "quantity": 10, "note": "so item"},
			}}, http.StatusUnprocessableEntity},
		// apabila sales order item id gak termasuk sales order id
		{tester.D{"sales_order_id": eso, "priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem2, "quantity": 10, "note": "so item"},
			}}, http.StatusUnprocessableEntity},
		// apabila quantity sales order item id < quantity yang diinput
		{tester.D{"sales_order_id": eso, "priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItemD, "quantity": 10, "note": "so item"},
			}}, http.StatusUnprocessableEntity},
		// apabila quantity nya minus
		{tester.D{"sales_order_id": eso, "priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": -1, "note": "so item"},
			}}, http.StatusUnprocessableEntity},
		// apabila sales order item id nya sama
		{tester.D{"sales_order_id": eso, "priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 10, "note": "so item"},
				{"sales_order_item_id": esoItem, "quantity": 10, "note": "so item"},
			}}, http.StatusUnprocessableEntity},
		//apabila quantity+quantityfulfillment > quanitity order
		{tester.D{"sales_order_id": eso, "priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 6, "note": "so item  a"},
			}}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.POST("/v1/fulfillment").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

func TestHandler_URLMapping_PUT_IfNotID(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// soItem dan sales order fulfillment status = new
	soItem := model.DummySalesOrderItem()
	soItem.Quantity = 20
	soItem.Save("Quantity")
	soItem.SalesOrder.FulfillmentStatus = "new"
	soItem.SalesOrder.Save("FulfillmentStatus")

	esoItem := common.Encrypt(soItem.ID)
	date := time.Now()

	// dummy fullfilment
	fullfillment := model.DummyWorkorderFulfillment()
	fullfillment.DocumentStatus = "new"
	fullfillment.IsDeleted = 0
	fullfillment.SalesOrder = soItem.SalesOrder
	fullfillment.Save("DocumentStatus", "SalesOrder", "IsDeleted")
	efullfilment := common.Encrypt(fullfillment.ID)

	//dummy soItem yang sales ordernya tidak sama dengan dummy soItem yang diatas
	soItemNotSo := model.DummySalesOrderItem()
	soItemNotSo.Quantity = 5
	soItemNotSo.Save("Quantity")
	soItemNotSo.SalesOrder.FulfillmentStatus = "new"
	soItemNotSo.SalesOrder.Save("FulfillmentStatus")

	esoItemNotSo := common.Encrypt(soItemNotSo.ID)

	var data = []struct {
		req         tester.D
		quantityNew float32
		expected    int
	}{
		//sukses
		// lengkap priority = routine
		{tester.D{"priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so1",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 5, "note": "so item"},
			}}, 5, http.StatusOK},
		// lengkap priority = rush
		{tester.D{"priority": "rush", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so1",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 5, "note": "so item"},
			}}, 5, http.StatusOK},
		// lengkap priority = emergency
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so1",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 5, "note": "so item"},
			}}, 5, http.StatusOK},
		// lengkap note fulfillment kosong
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 5, "note": "so item"},
			}}, 5, http.StatusOK},

		//gagal
		//  priority != routine,rush,emergency
		{tester.D{"priority": "very important", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so1",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 1, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// due date kosong
		{tester.D{"priority": "emergency", "due_date": nil, "shipping_address": "Jalan Kresek Raya", "note": "g-so2",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 1, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// jika shipping_address kosong
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "", "note": "g-so3",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 1, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// sales order item id kosong
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so4",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": "", "quantity": 1, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// sales order item nya gak valid
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so0",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": "absd", "quantity": 1, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// sales order item nya gak ada di database
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so6",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": "1232413", "quantity": 1, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// sales order item nya bukan termasuk dari sales order id
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so7",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItemNotSo, "quantity": 1, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// sales order item id nya sama
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so8",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 2, "note": "so item"},
				{"sales_order_item_id": esoItem, "quantity": 2, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// quantity kosong
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so9",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 0, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// quantity minus
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-s010",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": -1, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// quantity > quantity sales order item
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so11",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 21, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// quantity + quantity fulfillment item sebelumnya > quantity sales order item
		// quantity sales order item = 20
		// quantity fulfillment item yang masuk udah ada 15
		// quantity input = 6, maka 15+6 > 20
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so12",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 6, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/fulfillment/"+efullfilment).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
				if res.Code == http.StatusOK {
					fullItem := &model.WorkorderFulfillmentItem{SalesOrderItem: soItem}
					fullItem.Read("SalesOrderItem")
					assert.Equal(t, tes.quantityNew, fullItem.Quantity)
				}
			})
	}
}

func TestHandler_URLMapping_PUT_Fulfillment_Document_Status_Is_Deleted(t *testing.T) {
	// authorization
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// soItem dan sales order fulfillment status = new
	soItem := model.DummySalesOrderItem()
	soItem.Quantity = 20
	soItem.Save("Quantity")
	soItem.SalesOrder.FulfillmentStatus = "new"
	soItem.SalesOrder.Save("FulfillmentStatus")

	esoItem := common.Encrypt(soItem.ID)
	date := time.Now()
	// dummy fullfilment
	fullfillment := model.DummyWorkorderFulfillment()
	fullfillment.DocumentStatus = "finished"
	fullfillment.SalesOrder = soItem.SalesOrder
	fullfillment.IsDeleted = 0
	fullfillment.Save("DocumentStatus", "SalesOrder", "IsDeleted")
	efullfilment := common.Encrypt(fullfillment.ID)

	var data = []struct {
		req      tester.D
		expected int
	}{
		//gagal
		// karena fullfillment.DocumentStatus = "finished"
		{tester.D{"priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so1",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 5, "note": "so item"},
			}}, http.StatusUnprocessableEntity},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/fulfillment/"+efullfilment).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

	fullfillment.DocumentStatus = "active"
	fullfillment.Save("DocumentStatus")
	data = []struct {
		req      tester.D
		expected int
	}{
		//gagal
		// karena fullfillment.DocumentStatus = "active"
		{tester.D{"priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so1",
			"workorder_fulfillment_items": []tester.D{
				{"sales_order_item_id": esoItem, "quantity": 5, "note": "so item"},
			}}, http.StatusUnprocessableEntity},
	}

	for _, tes := range data {
		ng.PUT("/v1/fulfillment/"+efullfilment).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

	fullfillment.DocumentStatus = "new"
	fullfillment.IsDeleted = 1
	fullfillment.Save("DocumentStatus", "IsDeleted")

	request := tester.D{}

	ng.PUT("/v1/fulfillment/"+efullfilment).
		SetJSON(request).
		Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", request, res.Body.String()))
		})
}

func TestHandler_URLMapping_PUT_If_ID_Fulfillment_Item(t *testing.T) {
	// authorization
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// soItem dan sales order fulfillment status = new
	soItem := model.DummySalesOrderItem()
	soItem.Quantity = 20
	soItem.Save("Quantity")
	soItem.SalesOrder.FulfillmentStatus = "new"
	soItem.SalesOrder.Save("FulfillmentStatus")

	esoItem := common.Encrypt(soItem.ID)

	//dummy soItem yang sales ordernya tidak sama dengan dummy soItem yang diatas
	soItemNotSo := model.DummySalesOrderItem()
	soItemNotSo.Quantity = 5
	soItemNotSo.Save("Quantity")
	soItemNotSo.SalesOrder.FulfillmentStatus = "new"
	soItemNotSo.SalesOrder.Save("FulfillmentStatus")

	esoItemNotSo := common.Encrypt(soItemNotSo.ID)

	date := time.Now()
	// dummy fullfilment
	fullfillment := model.DummyWorkorderFulfillment()
	fullfillment.DocumentStatus = "new"
	fullfillment.SalesOrder = soItem.SalesOrder
	fullfillment.IsDeleted = 0
	fullfillment.Save("DocumentStatus", "SalesOrder", "IsDeleted")

	efullfilment := common.Encrypt(fullfillment.ID)

	// dummy fulfillment item
	Itemfull := model.DummyWorkorderFulfillmentItem()
	Itemfull.WorkorderFulfillment = fullfillment
	Itemfull.SalesOrderItem = soItem
	Itemfull.Quantity = 1
	Itemfull.Save("WorkorderFulfillment", "SalesOrderItem", "Quantity")

	eItemfull := common.Encrypt(Itemfull.ID)

	// dummy untuk mencek quantity + quantity yang ada di fulfillment kecuali id fuldillment yang diinput
	Itemfull2 := model.DummyWorkorderFulfillmentItem()
	Itemfull2.WorkorderFulfillment = fullfillment
	Itemfull2.SalesOrderItem = soItem
	Itemfull2.Quantity = 15
	Itemfull2.Save("WorkorderFulfillment", "SalesOrderItem", "Quantity")

	var data = []struct {
		req              tester.D
		quantityItemFull float32
		expected         int
	}{
		//sukses
		// lengkap priority = routine
		{tester.D{"priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so1",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": esoItem, "quantity": 5, "note": "so item"},
			}}, 5, http.StatusOK},
		// lengkap priority = rush
		{tester.D{"priority": "rush", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so1",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": esoItem, "quantity": 5, "note": "so item"},
			}}, 5, http.StatusOK},
		// lengkap priority = emergency
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so1",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": esoItem, "quantity": 5, "note": "so item"},
			}}, 5, http.StatusOK},
		// lengkap note fulfillment kosong
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": esoItem, "quantity": 5, "note": "so item"},
			}}, 5, http.StatusOK},
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so12",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": esoItem, "quantity": 5, "note": "so item"},
			}}, 0, http.StatusOK},

		//gagal
		// apabila id fulfillment item gak valid
		{tester.D{"priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so1",
			"workorder_fulfillment_items": []tester.D{
				{"id": "ancd", "sales_order_item_id": esoItem, "quantity": 5, "note": "so item"},
			}}, 5, http.StatusUnprocessableEntity},
		// apabila id fulfillment item gak ada di database
		{tester.D{"priority": "routine", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "so1",
			"workorder_fulfillment_items": []tester.D{
				{"id": "12313123", "sales_order_item_id": esoItem, "quantity": 5, "note": "so item"},
			}}, 5, http.StatusUnprocessableEntity},
		//  priority != routine,rush,emergency
		{tester.D{"priority": "very important", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so1",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": esoItem, "quantity": 1, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// due date kosong
		{tester.D{"priority": "emergency", "due_date": nil, "shipping_address": "Jalan Kresek Raya", "note": "g-so2",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": esoItem, "quantity": 1, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// jika shipping_address kosong
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "", "note": "g-so3",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": esoItem, "quantity": 1, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// sales order item id kosong
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so4",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": "", "quantity": 1, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// sales order item nya gak valid
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so5",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": "absd", "quantity": 1, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// sales order item nya gak ada di database
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so6",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": "1232413", "quantity": 1, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// sales order item nya bukan termasuk dari sales order id
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so7",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": esoItemNotSo, "quantity": 1, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// sales order item id nya sama
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so8",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": esoItem, "quantity": 2, "note": "so item"},
				{"id": eItemfull, "sales_order_item_id": esoItem, "quantity": 2, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// quantity kosong
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so9",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": esoItem, "quantity": 0, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// quantity minus
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-s010",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": esoItem, "quantity": -1, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// quantity > quantity sales order item
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so11",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": esoItem, "quantity": 21, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
		// quantity + quantity fulfillment item sebelumnya > quantity sales order item
		// quantity sales order item = 20
		// quantity fulfillment item yang masuk udah ada 15 , dilihat dari itemfull2
		// quantity input = 6, maka 15+6 > 20
		{tester.D{"priority": "emergency", "due_date": date, "shipping_address": "Jalan Kresek Raya", "note": "g-so12",
			"workorder_fulfillment_items": []tester.D{
				{"id": eItemfull, "sales_order_item_id": esoItem, "quantity": 6, "note": "so item"},
			}}, 0, http.StatusUnprocessableEntity},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.PUT("/v1/fulfillment/"+efullfilment).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
				//jika di update
				if res.Code == http.StatusOK {
					assert.NotEqual(t, Itemfull.Quantity, tes.quantityItemFull)
				}
			})
	}
}

func TestHandler_URLMapping_PUT_Approve(t *testing.T) {
	// authorization
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	fulfillment := model.DummyWorkorderFulfillment()
	so := model.DummySalesOrder()

	iv := model.DummyItemVariant()
	iv.AvailableStock = 50
	iv.CommitedStock = 20
	iv.Save("AvailableStock", "CommitedStock")

	soitem := model.DummySalesOrderItem()
	soitem.ItemVariant = iv
	soitem.SalesOrder = so
	soitem.Save("ItemVariant", "SalesOrder")

	item1 := model.DummyWorkorderFulfillmentItem()
	item1.Quantity = 10
	item1.SalesOrderItem = soitem
	item1.WorkorderFulfillment = fulfillment
	item1.Save("Quantity", "WorkorderFulfillment", "SalesOrderItem")

	ivStock1 := model.DummyItemVariantStock()
	ivStock1.ItemVariant = iv
	ivStock1.AvailableStock = 12
	ivStock1.CreatedAt = time.Now()
	ivStock1.UnitCost = 1000
	ivStock1.Save("ItemVariant", "UnitCost", "CreatedAt", "AvailableStock")

	// dummy ke 2
	iv2 := model.DummyItemVariant()
	iv2.AvailableStock = 50
	iv2.CommitedStock = 20
	iv2.Save("AvailableStock", "CommitedStock")

	soItem2 := model.DummySalesOrderItem()
	soItem2.ItemVariant = iv2
	soItem2.SalesOrder = so
	soItem2.Save("ItemVariant", "SalesOrder")

	item2 := model.DummyWorkorderFulfillmentItem()
	item2.Quantity = 10
	item2.SalesOrderItem = soItem2
	item2.WorkorderFulfillment = fulfillment
	item2.Save("Quantity", "WorkorderFulfillment", "SalesOrderItem")

	for a := 0; a < 3; a++ {
		ivStock2 := model.DummyItemVariantStock()
		ivStock2.ItemVariant = iv2
		ivStock2.AvailableStock = 5
		ivStock2.UnitCost = 1000
		ivStock2.Save("ItemVariant", "AvailableStock", "UnitCost")
	}

	var soitems []*model.SalesOrderItem
	soitems = append(soitems, soitem, soItem2)

	so.SalesOrderItems = soitems
	so.Save("SalesOrderItem")

	var items []*model.WorkorderFulfillmentItem
	items = append(items, item1, item2)

	fulfillment.WorkorderFulFillmentItems = items
	fulfillment.SalesOrder = so
	fulfillment.DocumentStatus = "active"
	fulfillment.IsDeleted = 0
	fulfillment.Save("WorkorderFulFillmentItems", "SalesOrder", "DocumentStatus", "IsDeleted")

	efulfillment := common.Encrypt(fulfillment.ID)

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	req := tester.D{}
	ng.PUT("/v1/fulfillment/"+efulfillment+"/approve").
		SetJSON(req).
		Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", req, res.Body.String()))
		})

	fulfillment.DocumentStatus = "finished"
	fulfillment.Save("DocumentStatus")
	ng.PUT("/v1/fulfillment/"+efulfillment+"/approve").
		SetJSON(req).
		Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", req, res.Body.String()))
		})

}

func TestHandler_URLMapping_PUT_Approve2Success(t *testing.T) {
	// cleared database
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant_stock").Exec()
	o.Raw("DELETE FROM item_variant").Exec()
	o.Raw("DELETE FROM sales_order").Exec()

	// buat user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy item
	itmVar1 := model.DummyItemVariant()
	itmVar1.IsDeleted = int8(0)
	itmVar1.AvailableStock = float32(200)
	itmVar1.CommitedStock = float32(70)
	itmVar1.Save()

	itmVarStock1 := model.DummyItemVariantStock()
	itmVarStock1.AvailableStock = float32(20)
	itmVarStock1.UnitCost = float64(1000)
	itmVarStock1.ItemVariant = itmVar1
	itmVarStock1.Save()
	itmVarStock2 := model.DummyItemVariantStock()
	itmVarStock2.AvailableStock = float32(80)
	itmVarStock2.UnitCost = float64(2000)
	itmVarStock2.ItemVariant = itmVar1
	itmVarStock2.Save()

	itmVarLog := model.DummyItemVariantStockLog()
	itmVarLog.ItemVariantStock = itmVarStock1
	itmVarLog.Quantity = float32(20)
	itmVarLog.RefType = "direct_placement"
	itmVarLog.FinalStock = float32(20)
	itmVarLog.LogType = "in"
	itmVarLog.Save()
	itmVarLog1 := model.DummyItemVariantStockLog()
	itmVarLog1.ItemVariantStock = itmVarStock2
	itmVarLog1.Quantity = float32(80)
	itmVarLog1.RefType = "direct_placement"
	itmVarLog1.FinalStock = float32(100)
	itmVarLog1.LogType = "in"
	itmVarLog1.Save()

	// buat dummy
	fulfill := model.DummyWorkorderFulfillment()
	fulfill.IsDeleted = int8(0)
	fulfill.DocumentStatus = "active"
	fulfill.Save()

	so := fulfill.SalesOrder
	so.IsDeleted = int8(0)
	so.DocumentStatus = "active"
	so.FulfillmentStatus = "active"
	so.TotalCost = float64(0)
	so.Save()

	fulfillItem := model.DummyWorkorderFulfillmentItem()
	fulfillItem.Quantity = float32(50)
	fulfillItem.WorkorderFulfillment = fulfill
	fulfillItem.Save()
	soi := fulfillItem.SalesOrderItem
	soi.SalesOrder = so
	soi.ItemVariant = itmVar1
	soi.Quantity = float32(50)
	soi.QuantityFulfillment = float32(0)
	soi.Save()
	// set body kosong
	req := tester.D{}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/fulfillment/"+common.Encrypt(fulfill.ID)+"/approve").
		SetJSON(req).
		Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", req, res.Body.String()))
		})
	itmVar1.Read("ID")
	assert.Equal(t, float32(50), itmVar1.AvailableStock)
	itmVarStock1.Read("ID")
	assert.Equal(t, float32(0), itmVarStock1.AvailableStock)
	itmVarStock2.Read("ID")
	assert.Equal(t, float32(50), itmVarStock2.AvailableStock)
	so.Read("ID")
	assert.Equal(t, float64(80000), so.TotalCost)
	assert.Equal(t, "finished", so.FulfillmentStatus)
	soi.Read("ID")
	assert.Equal(t, float32(50), soi.QuantityFulfillment)
}
