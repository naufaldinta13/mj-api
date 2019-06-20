// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receiving_test

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

// test TestCreateReceiving success
func TestCreateReceiving(t *testing.T) {

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	po := model.DummyPurchaseOrder()
	po.InvoiceStatus = "finished"
	po.DocumentStatus = "new"
	po.IsDeleted = 0
	po.Save()
	pox := common.Encrypt(po.ID)

	poi := model.DummyPurchaseOrderItem()
	poi.Quantity = 100
	poi.PurchaseOrder = po
	poi.Save()
	poix := common.Encrypt(poi.ID)

	wr := model.DummyWorkorderReceiving()
	wr.PurchaseOrder = po
	wr.DocumentStatus = "finished"
	wr.IsDeleted = 0
	wr.Save()

	wri := model.DummyWorkorderReceivingItem()
	wri.WorkorderReceiving = wr
	wri.PurchaseOrderItem = poi
	wri.Quantity = 10
	wri.Save()

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"recognition_date": time.Now(),
			"purchase_order": pox,
			"note":           "abc",
			"pic":            "ijal funky",
			"work_order_receiving_items": []tester.D{
				{"purchase_order_item": poix, "quantity": 50}, {"purchase_order_item": poix, "quantity": 40},
			}}, http.StatusOK},
		// PO ID cant be decrypt
		{tester.D{"recognition_date": time.Now(),
			"purchase_order": "aaaaa",
			"note":           "abc",
			"pic":            "ijal funky",
			"work_order_receiving_items": []tester.D{
				{"purchase_order_item": poix, "quantity": 50}, {"purchase_order_item": poix, "quantity": 40},
			}}, http.StatusUnprocessableEntity},
		// PO ID doesnt exist
		{tester.D{"recognition_date": time.Now(),
			"purchase_order": "99999999999999",
			"note":           "abc",
			"pic":            "ijal funky",
			"work_order_receiving_items": []tester.D{
				{"purchase_order_item": poix, "quantity": 50}, {"purchase_order_item": poix, "quantity": 40},
			}}, http.StatusUnprocessableEntity},
		// PO Item ID cant be decrypt
		{tester.D{"recognition_date": time.Now(),
			"purchase_order": "aaaaa",
			"note":           "abc",
			"pic":            "ijal funky",
			"work_order_receiving_items": []tester.D{
				{"purchase_order_item": "aaaaa", "quantity": 50}, {"purchase_order_item": poix, "quantity": 40},
			}}, http.StatusUnprocessableEntity},
		// PO Item ID doesnt exist
		{tester.D{"recognition_date": time.Now(),
			"purchase_order": "99999999999999",
			"note":           "abc",
			"pic":            "ijal funky",
			"work_order_receiving_items": []tester.D{
				{"purchase_order_item": "99999999999999", "quantity": 50}, {"purchase_order_item": poix, "quantity": 40},
			}}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.POST("/v1/receiving").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

}

// test TestCreateReceiving success
func TestCreateReceivingIvoiceActive(t *testing.T) {

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	po := model.DummyPurchaseOrder()
	po.IsDeleted = 0
	po.InvoiceStatus = "active"
	po.DocumentStatus = "new"
	po.Save()
	pox := common.Encrypt(po.ID)

	poi := model.DummyPurchaseOrderItem()
	poi.Quantity = 100
	poi.PurchaseOrder = po
	poi.Save()
	poix := common.Encrypt(poi.ID)

	wr := model.DummyWorkorderReceiving()
	wr.PurchaseOrder = po
	wr.DocumentStatus = "finished"
	wr.IsDeleted = 0
	wr.Save()

	wri := model.DummyWorkorderReceivingItem()
	wri.WorkorderReceiving = wr
	wri.PurchaseOrderItem = poi
	wri.Quantity = 10
	wri.Save()

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"recognition_date": time.Now(),
			"purchase_order": pox,
			"note":           "abc",
			"pic":            "ijal funky",
			"work_order_receiving_items": []tester.D{
				{"purchase_order_item": poix, "quantity": 50}, {"purchase_order_item": poix, "quantity": 40},
			}}, http.StatusOK},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.POST("/v1/receiving").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

}

// test get all data work order receiving success
func TestGetAllData(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/receiving", "GET", 200},
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

// test get all  work order receiving with No Token
func TestGetAllDataNoToken(t *testing.T) {

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/receiving", "GET", http.StatusBadRequest},
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

func TestARouting(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	dwr1 := model.DummyWorkorderReceiving()
	dwr1.IsDeleted = 0
	dwr1.Save()
	dwri1 := model.DummyWorkorderReceivingItem()
	dwri1.WorkorderReceiving = dwr1
	dwri1.Save("WorkorderReceiving")

	dwr2 := model.DummyWorkorderReceiving()
	dwr2.IsDeleted = 1
	dwr2.Save()
	dwri2 := model.DummyWorkorderReceivingItem()
	dwri2.WorkorderReceiving = dwr2
	dwri2.Save("WorkorderReceiving")

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/receiving/" + common.Encrypt(dwr1.ID), "GET", 200},
		{"/v1/receiving/" + common.Encrypt(dwr2.ID), "GET", 404},
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
