// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package salesReturn_test

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

// test get all and detail sales return success
func TestGetAllData(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sr := model.DummySalesReturn()
	sr.IsDeleted = 0
	sr.Save()
	id := common.Encrypt(sr.ID)

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/sales-return", "GET", 200},
		{"/v1/sales-return/" + id, "GET", 200},
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

// test get all and detail data sales return with No Token
func TestGetAllDataNoToken(t *testing.T) {
	sr := model.DummySalesReturn()
	sr.IsDeleted = 0
	sr.Save()
	id := common.Encrypt(sr.ID)

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/sales-return", "GET", http.StatusBadRequest},
		{"/v1/sales-return/" + id, "GET", http.StatusBadRequest},
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

// test get all and detail data sales return where usergroup is cashier
func TestGetAllDataIsCashier(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(4)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sr := model.DummySalesReturn()
	sr.IsDeleted = 0
	sr.Save()
	id := common.Encrypt(sr.ID)

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/sales-return", "GET", http.StatusUnauthorized},
		{"/v1/sales-return/" + id, "GET", http.StatusUnauthorized},
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

// test get detail sales return when ID not found
func TestGetDetailDataNotFound(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/sales-return/999999", "GET", http.StatusNotFound},
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

// test TestCreateSalesReturn success
func TestCreateSalesReturn(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("DELETE FROM sales_return").Exec()

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	so := model.DummySalesOrder()
	so.DocumentStatus = "new"
	so.IsDeleted = 0
	so.Save()
	soID := common.Encrypt(so.ID)
	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.Quantity = 10
	soi.UnitPrice = 20000
	soi.Subtotal = float64(soi.Quantity) * soi.UnitPrice
	soi.Save()
	soiID := common.Encrypt(soi.ID)

	soi2 := model.DummySalesOrderItem()
	soi2.SalesOrder = so
	soi2.Quantity = 10
	soi2.UnitPrice = 10000
	soi2.Subtotal = float64(soi2.Quantity) * soi2.UnitPrice
	soi2.Save()
	soi2ID := common.Encrypt(soi2.ID)

	f := model.DummyWorkorderFulfillment()
	f.SalesOrder = so
	f.IsDelivered = 1
	f.IsDeleted = 0
	f.DocumentStatus = "finished"
	f.CreatedBy = model.DummyUser()
	f.Save()

	fi := model.DummyWorkorderFulfillmentItem()
	fi.WorkorderFulfillment = f
	fi.Quantity = 10
	fi.SalesOrderItem = soi
	fi.Save()

	fi2 := model.DummyWorkorderFulfillmentItem()
	fi2.WorkorderFulfillment = f
	fi2.Quantity = 10
	fi2.SalesOrderItem = soi2
	fi2.Save()

	so1 := model.DummySalesOrder()
	so1.DocumentStatus = "new"
	so1.IsDeleted = 0
	so1.Save()

	soi3 := model.DummySalesOrderItem()
	soi3.SalesOrder = so1
	soi3.Quantity = 10
	soi3.UnitPrice = 10000
	soi3.Subtotal = float64(soi3.Quantity) * soi3.UnitPrice
	soi3.Save()
	soi3ID := common.Encrypt(soi3.ID)

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"sales_order": soID, "recognition_date": "2017-08-21T00:00:00Z", "note": "jkt", "sales_return_items": []tester.D{
			{"sales_order_item": soiID, "quantity": 10, "note": "tes"},
			{"sales_order_item": soi2ID, "quantity": 10, "note": "tes"},
		}}, http.StatusOK},
		{tester.D{"sales_order": soID, "note": "", "sales_return_items": []tester.D{
			{"sales_order_item": "abc", "quantity": -3, "note": "tes"},
		}}, http.StatusUnprocessableEntity},
		{tester.D{"sales_order": soID, "note": "", "sales_return_items": []tester.D{
			{"sales_order_item": "999999999999999", "quantity": 10, "note": "tes"},
		}}, http.StatusUnprocessableEntity},
		{tester.D{"sales_order": soID, "note": "", "sales_return_items": []tester.D{
			{"sales_order_item": soi3ID, "quantity": -3, "note": "tes"},
		}}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.POST("/v1/sales-return").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

	//cek count sales return
	var total int64
	_ = o.Raw("select count(*) from sales_return").QueryRow(&total)
	assert.Equal(t, int64(1), total)

	//cek count sales return item
	var totalSri int64
	_ = o.Raw("select count(*) from sales_return_item").QueryRow(&totalSri)
	assert.Equal(t, int64(2), totalSri)
}

// test TestCreateSalesReturn quantity more than can be return
func TestCreateSalesReturnFail(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("DELETE FROM sales_return").Exec()

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	so := model.DummySalesOrder()
	so.DocumentStatus = "new"
	so.IsDeleted = 0
	so.Save()
	soID := common.Encrypt(so.ID)
	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.Quantity = 10
	soi.UnitPrice = 20000
	soi.Subtotal = float64(soi.Quantity) * soi.UnitPrice
	soi.Save()
	soiID := common.Encrypt(soi.ID)

	f := model.DummyWorkorderFulfillment()
	f.SalesOrder = so
	f.IsDelivered = 1
	f.IsDeleted = 0
	f.DocumentStatus = "finished"
	f.CreatedBy = model.DummyUser()
	f.Save()

	fi := model.DummyWorkorderFulfillmentItem()
	fi.WorkorderFulfillment = f
	fi.Quantity = 10
	fi.SalesOrderItem = soi
	fi.Save()

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"sales_order": soID, "recognition_date": "2017-08-21T00:00:00Z", "note": "jkt", "sales_return_items": []tester.D{
			{"sales_order_item": soiID, "quantity": 20, "note": "tes"},
		}}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.POST("/v1/sales-return").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

	//cek count sales return
	var total int64
	_ = o.Raw("select count(*) from sales_return").QueryRow(&total)
	assert.Equal(t, int64(0), total)

	//cek count sales return item
	var totalSri int64
	_ = o.Raw("select count(*) from sales_return_item").QueryRow(&totalSri)
	assert.Equal(t, int64(0), totalSri)
}

// test TestCreateSalesReturn No Token
func TestCreateSalesReturnNoToken(t *testing.T) {

	so := model.DummySalesOrder()
	soID := common.Encrypt(so.ID)
	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.Save()
	soiID := common.Encrypt(soi.ID)

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"sales_order": soID, "recognition_date": time.Now(), "total_amount": 20000, "note": "jkt", "sales_return_item": []tester.D{
			{"sales_order_item": soiID, "quantity": 5, "note": "tes"},
		}}, http.StatusBadRequest},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.POST("/v1/sales-return").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// test TestUpdateSalesReturn success
func TestUpdateSalesReturn(t *testing.T) {

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	so := model.DummySalesOrder()
	so.DocumentStatus = "new"
	so.IsDeleted = 0
	so.Save()

	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.Quantity = 20
	soi.UnitPrice = 20000
	soi.Subtotal = float64(soi.Quantity) * soi.UnitPrice
	soi.Save()
	soiID := common.Encrypt(soi.ID)

	soi2 := model.DummySalesOrderItem()
	soi2.SalesOrder = so
	soi2.Quantity = 20
	soi2.UnitPrice = 10000
	soi2.Subtotal = float64(soi.Quantity) * soi.UnitPrice
	soi2.Save()
	soi2ID := common.Encrypt(soi2.ID)

	f := model.DummyWorkorderFulfillment()
	f.SalesOrder = so
	f.IsDelivered = 1
	f.IsDeleted = 0
	f.DocumentStatus = "finished"
	f.CreatedBy = model.DummyUser()
	f.Save()

	fi := model.DummyWorkorderFulfillmentItem()
	fi.WorkorderFulfillment = f
	fi.Quantity = 20
	fi.SalesOrderItem = soi
	fi.Save()

	fi2 := model.DummyWorkorderFulfillmentItem()
	fi2.WorkorderFulfillment = f
	fi2.Quantity = 20
	fi2.SalesOrderItem = soi2
	fi2.Save()

	srx := model.DummySalesReturn()
	srx.IsDeleted = 0
	srx.SalesOrder = so
	srx.Save()

	srix := model.DummySalesReturnItem()
	srix.SalesReturn = srx
	srix.SalesOrderItem = soi
	srix.Quantity = 4
	srix.Save()

	sri2x := model.DummySalesReturnItem()
	sri2x.SalesReturn = srx
	sri2x.SalesOrderItem = soi2
	sri2x.Quantity = 4
	sri2x.Save()

	sr := model.DummySalesReturn()
	sr.IsDeleted = 0
	sr.DocumentStatus = "new"
	sr.SalesOrder = so
	sr.Save()

	sri := model.DummySalesReturnItem()
	sri.SalesReturn = sr
	sri.SalesOrderItem = soi
	sri.Quantity = 5
	sri.Save()

	sri2 := model.DummySalesReturnItem()
	sri2.SalesReturn = sr
	sri2.SalesOrderItem = soi2
	sri2.Quantity = 5
	sri2.Save()

	so1 := model.DummySalesOrder()
	so1.DocumentStatus = "new"
	so1.IsDeleted = 0
	so1.Save()

	soi3 := model.DummySalesOrderItem()
	soi3.SalesOrder = so1
	soi3.Quantity = 10
	soi3.UnitPrice = 10000
	soi3.Subtotal = float64(soi3.Quantity) * soi3.UnitPrice
	soi3.Save()
	soi3ID := common.Encrypt(soi3.ID)

	id := common.Encrypt(sr.ID)

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"recognition_date": "2017-08-21T00:00:00Z", "note": "UPDATE", "sales_return_items": []tester.D{
			{"id": common.Encrypt(sri.ID), "sales_order_item": soiID, "quantity": 12, "note": "update"},
			{"id": common.Encrypt(sri2.ID), "sales_order_item": soi2ID, "quantity": 12, "note": "update"},
		}}, http.StatusOK},
		{tester.D{"recognition_date": "2017-08-21T00:00:00Z", "note": "UPDATE", "sales_return_items": []tester.D{
			{"id": common.Encrypt(sri.ID), "sales_order_item": soiID, "quantity": 111, "note": "update"},
			{"id": common.Encrypt(sri2.ID), "sales_order_item": "aaa", "quantity": 13, "note": "update"},
		}}, http.StatusUnprocessableEntity},
		{tester.D{"note": "", "sales_return_items": []tester.D{
			{"sales_order_item": "999999999999999", "quantity": 10, "note": "tes"},
		}}, http.StatusUnprocessableEntity},
		{tester.D{"note": "", "sales_return_items": []tester.D{
			{"sales_order_item": soi3ID, "quantity": -3, "note": "tes"},
		}}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.PUT("/v1/sales-return/"+id).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

	sr.Read("ID")
	assert.Equal(t, "UPDATE", sr.Note)

}

// test TestUpdateSalesReturn success
func TestUpdateSalesReturnCancelled(t *testing.T) {

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	so := model.DummySalesOrder()
	so.DocumentStatus = "new"
	so.IsDeleted = 0
	so.Save()

	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.Quantity = 20
	soi.UnitPrice = 20000
	soi.Subtotal = float64(soi.Quantity) * soi.UnitPrice
	soi.Save()
	soiID := common.Encrypt(soi.ID)

	f := model.DummyWorkorderFulfillment()
	f.SalesOrder = so
	f.IsDelivered = 1
	f.IsDeleted = 0
	f.DocumentStatus = "finished"
	f.CreatedBy = model.DummyUser()
	f.Save()

	fi := model.DummyWorkorderFulfillmentItem()
	fi.WorkorderFulfillment = f
	fi.Quantity = 20
	fi.SalesOrderItem = soi
	fi.Save()

	srx := model.DummySalesReturn()
	srx.IsDeleted = 0
	srx.SalesOrder = so
	srx.Save()

	srix := model.DummySalesReturnItem()
	srix.SalesReturn = srx
	srix.SalesOrderItem = soi
	srix.Quantity = 4
	srix.Save()

	sr := model.DummySalesReturn()
	sr.IsDeleted = 0
	sr.SalesOrder = so
	sr.DocumentStatus = "cancelled"
	sr.Save()

	sri := model.DummySalesReturnItem()
	sri.SalesReturn = sr
	sri.SalesOrderItem = soi
	sri.Quantity = 5
	sri.Save()

	id := common.Encrypt(sr.ID)

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"recognition_date": "2017-08-21T00:00:00Z", "note": "UPDATE", "sales_return_items": []tester.D{
			{"id": common.Encrypt(sri.ID), "sales_order_item": soiID, "quantity": 12, "note": "update"},
		}}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.PUT("/v1/sales-return/"+id).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// test update sales return is not found id
func TestUpdateSalesReturnIDNotFound(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	var put = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"recognition_date": "2017-08-21T00:00:00Z", "note": "UPDATE", "sales_return_items": []tester.D{
			{"id": "65536", "quantity": 11, "note": "update"},
			{"id": "65536", "quantity": 13, "note": "update"},
		}}, http.StatusNotFound},
	}
	ng := tester.New()
	for _, tes := range put {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.PUT("/v1/sales-return/9999999999999").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// test TestUpdateSalesReturnNoToken No Token
func TestUpdateSalesReturnNoToken(t *testing.T) {

	so := model.DummySalesOrder()
	soID := common.Encrypt(so.ID)
	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.Save()
	soiID := common.Encrypt(soi.ID)

	sr := model.DummySalesReturn()
	sr.Save()
	sri := model.DummySalesReturnItem()
	sri.SalesOrderItem = soi
	sri.SalesReturn = sr
	sri.Save()
	sri2 := model.DummySalesReturnItem()
	sri2.SalesOrderItem = soi
	sri2.SalesReturn = sr
	sri2.Save()

	id := common.Encrypt(sr.ID)
	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"sales_order": soID, "recognition_date": time.Now(), "total_amount": 20000, "note": "jkt", "sales_return_items": []tester.D{
			{"sales_order_item": soiID, "quantity": 5, "note": "tes"},
		}}, http.StatusBadRequest},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.PUT("/v1/sales-return/"+id).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// test TestUpdateSalesReturn success where
// sales return ID doesn't exist will be delete in database
func TestUpdateSalesReturnDeleteOldDataSRItem(t *testing.T) {

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	so := model.DummySalesOrder()
	so.DocumentStatus = "new"
	so.IsDeleted = 0
	so.Save()

	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.Quantity = 20
	soi.UnitPrice = 20000
	soi.Subtotal = float64(soi.Quantity) * soi.UnitPrice
	soi.Save()
	soiID := common.Encrypt(soi.ID)

	soi2 := model.DummySalesOrderItem()
	soi2.SalesOrder = so
	soi2.Quantity = 20
	soi2.UnitPrice = 10000
	soi2.Subtotal = float64(soi.Quantity) * soi.UnitPrice
	soi2.Save()

	f := model.DummyWorkorderFulfillment()
	f.SalesOrder = so
	f.IsDelivered = 1
	f.IsDeleted = 0
	f.DocumentStatus = "finished"
	f.CreatedBy = model.DummyUser()
	f.Save()

	fi := model.DummyWorkorderFulfillmentItem()
	fi.WorkorderFulfillment = f
	fi.Quantity = 20
	fi.SalesOrderItem = soi
	fi.Save()

	fi2 := model.DummyWorkorderFulfillmentItem()
	fi2.WorkorderFulfillment = f
	fi2.Quantity = 20
	fi2.SalesOrderItem = soi2
	fi2.Save()

	sr := model.DummySalesReturn()
	sr.IsDeleted = 0
	sr.DocumentStatus = "new"
	sr.SalesOrder = so
	sr.Save()

	sri := model.DummySalesReturnItem()
	sri.SalesReturn = sr
	sri.SalesOrderItem = soi
	sri.Quantity = 5
	sri.Save()

	sri2 := model.DummySalesReturnItem()
	sri2.SalesReturn = sr
	sri2.SalesOrderItem = soi2
	sri2.Quantity = 5
	sri2.Save()

	id := common.Encrypt(sr.ID)

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"recognition_date": "2017-08-21T00:00:00Z", "note": "UPDATE", "sales_return_items": []tester.D{
			{"id": common.Encrypt(sri.ID), "sales_order_item": soiID, "quantity": 12, "note": "update"},
		}}, http.StatusOK},
		{tester.D{"recognition_date": "2017-08-21T00:00:00Z", "note": "UPDATE", "sales_return_items": []tester.D{
			{"id": common.Encrypt(sri.ID), "sales_order_item": soiID, "quantity": 111, "note": "update"},
			{"id": common.Encrypt(sri2.ID), "sales_order_item": "aaa", "quantity": 13, "note": "update"},
		}}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.PUT("/v1/sales-return/"+id).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

	sr.Read("ID")
	assert.Equal(t, "UPDATE", sr.Note)

}

func TestCancelSalesReturn(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sr := model.DummySalesReturn()
	sr.IsDeleted = 0
	sr.DocumentStatus = "new"
	sr.Note = "Cancel"
	sr.Save()

	fe := model.DummyFinanceExpense()
	fe.IsDeleted = 0
	fe.RefType = "sales_return"
	fe.Note = "update SR coy"
	fe.RefID = uint64(sr.ID)
	fe.Save()

	id := common.Encrypt(sr.ID)

	var cancel = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{}, http.StatusOK},
	}
	ng := tester.New()
	for _, tes := range cancel {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.PUT("/v1/sales-return/"+id+"/cancel").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

func TestCancelSalesIDNotFound(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	id := common.Encrypt(999999999)

	var cancel = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{}, http.StatusNotFound},
	}
	ng := tester.New()
	for _, tes := range cancel {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.PUT("/v1/sales-return/"+id+"/cancel").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

func TestCancelSalesCancelled(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sr := model.DummySalesReturn()
	sr.IsDeleted = 0
	sr.Note = "Cancel"
	sr.DocumentStatus = "cancelled"
	sr.Save()
	id := common.Encrypt(sr.ID)

	var cancel = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range cancel {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.PUT("/v1/sales-return/"+id+"/cancel").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// test TestCreateSalesReturnFailNoFulfillment doesnt exist
func TestCreateSalesReturnFailNoFulfillment(t *testing.T) {

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	so := model.DummySalesOrder()
	so.DocumentStatus = "new"
	so.IsDeleted = 0
	so.Save()
	soID := common.Encrypt(so.ID)
	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.Quantity = 10
	soi.UnitPrice = 20000
	soi.Subtotal = float64(soi.Quantity) * soi.UnitPrice
	soi.Save()
	soiID := common.Encrypt(soi.ID)

	so1 := model.DummySalesOrder()
	so1.Save()

	f := model.DummyWorkorderFulfillment()
	f.SalesOrder = so1
	f.IsDelivered = 1
	f.IsDeleted = 0
	f.DocumentStatus = "finished"
	f.CreatedBy = model.DummyUser()
	f.Save()

	fi := model.DummyWorkorderFulfillmentItem()
	fi.WorkorderFulfillment = f
	fi.Quantity = 10
	fi.SalesOrderItem = soi
	fi.Save()

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"sales_order": soID, "recognition_date": "2017-08-21T00:00:00Z", "note": "jkt", "sales_return_items": []tester.D{
			{"sales_order_item": soiID, "quantity": 10, "note": "tes"},
		}}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.POST("/v1/sales-return").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// test TestCreateSalesReturnFailNoSalesOrder quantity more than can be return
func TestCreateSalesReturnFailNoSalesOrder(t *testing.T) {

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"sales_order": "65536", "recognition_date": "2017-08-21T00:00:00Z", "note": "jkt", "sales_return_items": []tester.D{
			{"sales_order_item": "65536", "quantity": 10, "note": "tes"},
		}}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.POST("/v1/sales-return").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}
