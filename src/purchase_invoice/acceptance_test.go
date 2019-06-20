// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchaseInvoice_test

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

// test TestCreatePurchaseInvoice success
func TestCreatePurchaseInvoice(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM purchase_invoice").Exec()

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	po := model.DummyPurchaseOrder()
	po.InvoiceStatus = "active"
	po.TotalCharge = 5000
	po.IsDeleted = 0
	po.Save()
	pox := common.Encrypt(po.ID)

	pi := model.DummyPurchaseInvoice()
	pi.TotalAmount = 1000
	pi.PurchaseOrder = po
	pi.DocumentStatus = "new"
	pi.IsDeleted = 0
	pi.Save()

	pi2 := model.DummyPurchaseInvoice()
	pi2.TotalAmount = 1000
	pi2.PurchaseOrder = po
	pi.DocumentStatus = "new"
	pi.IsDeleted = 0
	pi2.Save()

	// clear database

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"recognition_date": time.Now(),
			"purchase_order":  pox,
			"due_date":        time.Now(),
			"total_amount":    3000,
			"note":            "AHAHAHA",
			"billing_address": "THE POK"}, http.StatusOK},
		{tester.D{"recognition_date": time.Now(),
			"purchase_order":  "999999999",
			"due_date":        time.Now(),
			"total_amount":    30000,
			"note":            "AHAHAHA",
			"billing_address": "THE POK"}, http.StatusUnprocessableEntity},
		{tester.D{"recognition_date": time.Now(),
			"purchase_order":  "aaaaa",
			"due_date":        time.Now(),
			"total_amount":    30000,
			"note":            "AHAHAHA",
			"billing_address": "THE POK"}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.POST("/v1/purchase-invoice").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

	//cek count purchase invoice
	var total int64
	_ = o.Raw("select count(*) from purchase_invoice").QueryRow(&total)
	assert.Equal(t, int64(3), total)
}

// test TestCreatePurchaseInvoice where invoice status = finished in purchase order
func TestCreatePurchaseInvoiceStatusFinished(t *testing.T) {

	o := orm.NewOrm()
	o.Raw("DELETE FROM purchase_invoice").Exec()

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	po := model.DummyPurchaseOrder()
	po.InvoiceStatus = "finished"
	po.TotalCharge = 5000
	po.IsDeleted = 0
	po.Save()
	pox := common.Encrypt(po.ID)

	pi := model.DummyPurchaseInvoice()
	pi.TotalAmount = 1000
	pi.PurchaseOrder = po
	pi.Save()

	pi2 := model.DummyPurchaseInvoice()
	pi2.TotalAmount = 1000
	pi2.PurchaseOrder = po
	pi2.Save()

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"recognition_date": time.Now(),
			"purchase_order":  pox,
			"due_date":        time.Now(),
			"total_amount":    3000,
			"note":            "AHAHAHA",
			"billing_address": "THE POK"}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.POST("/v1/purchase-invoice").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// test TestCreatePurchaseInvoice where invoice sum total amount pi > total charge po
func TestCreatePurchaseInvoiceTotalAmountGreaterThanPO(t *testing.T) {

	o := orm.NewOrm()
	o.Raw("DELETE FROM purchase_invoice").Exec()

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	po := model.DummyPurchaseOrder()
	po.InvoiceStatus = "active"
	po.TotalCharge = 5000
	po.IsDeleted = 0
	po.Save()
	pox := common.Encrypt(po.ID)

	pi := model.DummyPurchaseInvoice()
	pi.TotalAmount = 8000
	pi.IsDeleted = 0
	pi.PurchaseOrder = po
	pi.Save()

	pi2 := model.DummyPurchaseInvoice()
	pi2.TotalAmount = 8000
	pi2.IsDeleted = 0
	pi2.PurchaseOrder = po
	pi2.Save()

	var create = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"recognition_date": time.Now(),
			"purchase_order":  pox,
			"due_date":        time.Now(),
			"total_amount":    200000,
			"note":            "AHAHAHA",
			"billing_address": "THE POK"}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range create {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.POST("/v1/purchase-invoice").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// test get all data purchase invoice success
func TestGetAllData(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/purchase-invoice", "GET", 200},
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

// test get all data purchase invoice with No Token
func TestGetAllDataNoToken(t *testing.T) {

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/purchase-invoice", "GET", http.StatusBadRequest},
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

// test get detail data purchase invoice success
func TestShowPurchaseInvoice(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	pi := model.DummyPurchaseInvoice()
	pi.IsDeleted = 0
	pi.Save()
	id := common.Encrypt(pi.ID)

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/purchase-invoice/" + id, "GET", 200},
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

// test get detail data purchase invoice with No Token
func TestShowPurchaseInvoiceNoToken(t *testing.T) {
	pi := model.DummyPurchaseInvoice()
	pi.IsDeleted = 0
	pi.Save()
	id := common.Encrypt(pi.ID)

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/purchase-invoice/" + id, "GET", http.StatusBadRequest},
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

// test get detail purchase invoice when ID not found
func TestShowPurchaseInvoiceNotFound(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/purchase-invoice/999999", "GET", http.StatusNotFound},
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

// test TestUpdatePurchaseInvoice success
func TestUpdatePurchaseInvoice(t *testing.T) {

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	po := model.DummyPurchaseOrder()
	po.InvoiceStatus = "active"
	po.IsDeleted = 0
	po.TotalCharge = 5000
	po.Save()

	pi := model.DummyPurchaseInvoice()
	pi.TotalAmount = 1000
	pi.PurchaseOrder = po
	pi.DocumentStatus = "new"
	pi.IsDeleted = 0
	pi.Save()

	pi2 := model.DummyPurchaseInvoice()
	pi2.TotalAmount = 1000
	pi2.PurchaseOrder = po
	pi2.DocumentStatus = "new"
	pi2.IsDeleted = 0
	pi2.Save()

	id := common.Encrypt(pi.ID)

	var update = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"recognition_date": time.Now(),
			"total_amount": float64(2000),
			"note":         "UPDATE"}, http.StatusOK},
	}
	ng := tester.New()
	for _, tes := range update {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.PUT("/v1/purchase-invoice/"+id).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
	pix := &model.PurchaseInvoice{ID: pi.ID}
	pix.Read("ID")

	assert.Equal(t, float64(2000), pix.TotalAmount)
}

// test TestUpdatePurchaseInvoice where invoice status = finished in purchase order
func TestTestUpdatePurchaseInvoiceStatusFinished(t *testing.T) {

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	po := model.DummyPurchaseOrder()
	po.InvoiceStatus = "finished"
	po.IsDeleted = 0
	po.TotalCharge = 5000
	po.Save()

	pi := model.DummyPurchaseInvoice()
	pi.TotalAmount = 1000
	pi.PurchaseOrder = po
	pi.IsDeleted = 0
	pi.Save()

	pi2 := model.DummyPurchaseInvoice()
	pi2.TotalAmount = 1000
	pi2.PurchaseOrder = po
	pi2.IsDeleted = 0
	pi2.Save()

	id := common.Encrypt(pi.ID)

	var update = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"recognition_date": time.Now(),
			"total_amount": 3000,
			"note":         "AHAHAHA"}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range update {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.PUT("/v1/purchase-invoice/"+id).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// test TestUpdatePurchaseInvoice where invoice sum total amount pi > total charge po
func TestTestUpdatePurchaseInvoiceTotalAmountGreaterThanPO(t *testing.T) {

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	po := model.DummyPurchaseOrder()
	po.InvoiceStatus = "active"
	po.IsDeleted = 0
	po.TotalCharge = 5000
	po.Save()

	pi := model.DummyPurchaseInvoice()
	pi.TotalAmount = 8000
	pi.PurchaseOrder = po
	pi.IsDeleted = 0
	pi.Save()

	pi2 := model.DummyPurchaseInvoice()
	pi2.TotalAmount = 8000
	pi2.PurchaseOrder = po
	pi2.IsDeleted = 0
	pi2.Save()

	id := common.Encrypt(pi2.ID)

	var update = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"recognition_date": time.Now(),
			"total_amount": 20000,
			"note":         "AHAHAHA"}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	for _, tes := range update {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.PUT("/v1/purchase-invoice/"+id).
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}

// test TestUpdatePurchaseInvoice id not found
func TestUpdatePurchaseInvoiceIdNotFound(t *testing.T) {

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	var update = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{}, http.StatusNotFound},
	}
	ng := tester.New()
	for _, tes := range update {
		ng.SetHeader(tester.H{"Authorization": token})
		ng.PUT("/v1/purchase-invoice/99999999999").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}

}
