// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package salesInvoice_test

import (
	"fmt"
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

////GET///////////////////////////////////////

// TestHandler_URLMappingListSalesInvoiceSuccess test read list sales invoice,success
func TestHandler_URLMappingListSalesInvoiceSuccess(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM sales_invoice").Exec()

	// buat dummy sales invoice
	si1 := model.DummySalesInvoice()
	si1.BillingAddress = "test"
	si1.IsDeleted = int8(0)
	si1.Save()
	si2 := model.DummySalesInvoice()
	si2.BillingAddress = "test2"
	si2.IsDeleted = int8(1) // deleted
	si2.Save()
	si3 := model.DummySalesInvoice()
	si3.BillingAddress = "test"
	si3.IsDeleted = int8(0)
	si3.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/sales-invoice"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/sales-invoice", "GET"))
	})
}

// TestHandler_URLMappingListSalesInvoiceNoDataSuccess test read list sales invoice dengan data kosong, success
func TestHandler_URLMappingListSalesInvoiceNoDataSuccess(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM sales_invoice").Exec()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/sales-invoice"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/sales-invoice", "GET"))
	})
}

// TestHandler_URLMappingListSalesInvoiceFailNoToken test read list sales invoice tanpa token, fail
func TestHandler_URLMappingListSalesInvoiceFailNoToken(t *testing.T) {
	// test
	ng := tester.New()
	ng.Method = "GET"
	ng.Path = "/v1/sales-invoice"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/sales-invoice", "GET"))
	})
}

////SHOW///////////////////////////////////////

// TestHandler_URLMappingShowSalesInvoiceSuccess test show sales invoice,success
func TestHandler_URLMappingShowSalesInvoiceSuccess(t *testing.T) {
	// buat dummy sales invoice
	si1 := model.DummySalesInvoice()
	si1.IsDeleted = int8(0)
	si1.Save()
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/sales-invoice/" + common.Encrypt(si1.ID)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/sales-invoice/"+common.Encrypt(si1.ID), "GET"))
	})
}

// TestHandler_URLMappingShowSalesInvoiceNoDataFail test show sales invoice dengan id salah, fail
func TestHandler_URLMappingShowSalesInvoiceNoDataFail(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/sales-invoice/9999999"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/sales-invoice/9999999", "GET"))
	})
}

// TestHandler_URLMappingShowSalesInvoiceDeletedDataFail test show sales invoice dengan id yang di delete, fail
func TestHandler_URLMappingShowSalesInvoiceDeletedDataFail(t *testing.T) {
	// buat dummy sales invoice
	si1 := model.DummySalesInvoice()
	si1.IsDeleted = int8(1)
	si1.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/sales-invoice/" + common.Encrypt(si1.ID)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/sales-invoice/"+common.Encrypt(si1.ID), "GET"))
	})
}

// TestHandler_URLMappingShowSalesInvoiceFailNoToken test show sales invoice tanpa token, fail
func TestHandler_URLMappingShowSalesInvoiceFailNoToken(t *testing.T) {
	// test
	ng := tester.New()
	ng.Method = "GET"
	ng.Path = "/v1/sales-invoice/999999"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/sales-invoice/999999", "GET"))
	})
}

////POST///////////////////////////////////////

// TestHandler_URLMappingPostSalesInvoiceSuccess test create sales invoice,success
func TestHandler_URLMappingPostSalesInvoiceSuccess(t *testing.T) {
	//buat dummy so
	so := model.DummySalesOrder()
	so.IsDeleted = int8(0)
	so.InvoiceStatus = "new"
	so.DocumentStatus = "active"
	so.TotalCharge = float64(20000)
	so.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"recognition_date": "2017-05-08T14:52:36+07:00",
		"sales_order_id":   common.Encrypt(so.ID),
		"due_date":         "2017-05-08T14:52:36+07:00",
		"billing_address":  "address",
		"total_amount":     float64(20000),
		"note":             "note",
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/sales-invoice").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	so.Read("ID")
	assert.Equal(t, "active", so.InvoiceStatus)
	si := &model.SalesInvoice{SalesOrder: &model.SalesOrder{ID: so.ID}}
	e := si.Read("SalesOrder")
	assert.NoError(t, e)
	assert.Equal(t, float64(20000), si.TotalAmount)
	assert.Equal(t, user.ID, si.CreatedBy.ID)
}

// TestHandler_URLMappingPostSalesInvoiceSuccess2 test create sales invoice total amount 0,success
func TestHandler_URLMappingPostSalesInvoiceSuccess2(t *testing.T) {
	//buat dummy so
	so := model.DummySalesOrder()
	so.IsDeleted = int8(0)
	so.DocumentStatus = "active"
	so.InvoiceStatus = "new"
	so.TotalCharge = float64(0)
	so.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"recognition_date": "2017-05-08T14:52:36+07:00",
		"sales_order_id":   common.Encrypt(so.ID),
		"due_date":         "2017-05-08T14:52:36+07:00",
		"billing_address":  "address",
		"total_amount":     float64(0),
		"note":             "note",
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/sales-invoice").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	so.Read("ID")
	assert.Equal(t, "finished", so.InvoiceStatus)
}

// TestHandler_URLMappingPostSalesInvoiceSuccess3 test create sales invoice total amount 0,success
func TestHandler_URLMappingPostSalesInvoiceSuccess3(t *testing.T) {
	//buat dummy so
	so := model.DummySalesOrder()
	so.IsDeleted = int8(0)
	so.InvoiceStatus = "new"
	so.DocumentStatus = "active"
	so.TotalCharge = float64(20000)
	so.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"recognition_date": "2017-05-08T14:52:36+07:00",
		"sales_order_id":   common.Encrypt(so.ID),
		"due_date":         "2017-05-08T14:52:36+07:00",
		"billing_address":  "address",
		"total_amount":     float64(10000),
		"note":             "note",
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/sales-invoice").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	so.Read("ID")
	assert.Equal(t, "active", so.InvoiceStatus)
	si := &model.SalesInvoice{SalesOrder: &model.SalesOrder{ID: so.ID}}
	e := si.Read("SalesOrder")
	assert.NoError(t, e)
	assert.Equal(t, float64(10000), si.TotalAmount)
	assert.Equal(t, user.ID, si.CreatedBy.ID)
}

// TestHandler_URLMappingPostSalesInvoiceFailSOIDWrong test create sales invoice dengan so id salah,fail
func TestHandler_URLMappingPostSalesInvoiceFailSOIDWrong(t *testing.T) {
	//buat dummy so
	so := model.DummySalesOrder()
	so.IsDeleted = int8(0)
	so.DocumentStatus = "active"
	so.InvoiceStatus = "new"
	so.TotalCharge = float64(20000)
	so.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"recognition_date": "2017-05-08T14:52:36+07:00",
		"sales_order_id":   "aaaaaa",
		"due_date":         "2017-05-08T14:52:36+07:00",
		"billing_address":  "address",
		"total_amount":     float64(10000),
		"note":             "note",
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/sales-invoice").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMappingPostSalesInvoiceFailNoToken test create sales invoice tanpa token,fail
func TestHandler_URLMappingPostSalesInvoiceFailNoToken(t *testing.T) {
	//buat dummy so
	so := model.DummySalesOrder()
	so.IsDeleted = int8(0)
	so.DocumentStatus = "active"
	so.InvoiceStatus = "new"
	so.TotalCharge = float64(20000)
	so.Save()

	// setting body
	scenario := tester.D{
		"recognition_date": "2017-05-08T14:52:36+07:00",
		"sales_order_id":   common.Encrypt(so.ID),
		"due_date":         "2017-05-08T14:52:36+07:00",
		"billing_address":  "address",
		"total_amount":     float64(20000),
		"note":             "note",
	}
	ng := tester.New()
	ng.POST("/v1/sales-invoice").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPostSalesInvoiceFailStatusFinish test create sales invoice dengan invoice finish,fail
func TestHandler_URLMappingPostSalesInvoiceFailStatusFinish(t *testing.T) {
	//buat dummy so
	so := model.DummySalesOrder()
	so.IsDeleted = int8(0)
	so.DocumentStatus = "active"
	so.InvoiceStatus = "finished"
	so.TotalCharge = float64(20000)
	so.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"recognition_date": "2017-05-08T14:52:36+07:00",
		"sales_order_id":   common.Encrypt(so.ID),
		"due_date":         "2017-05-08T14:52:36+07:00",
		"billing_address":  "address",
		"total_amount":     float64(20000),
		"note":             "note",
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/sales-invoice").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPostSalesInvoiceFailDeleted test create sales invoice dengan so di delete,fail
func TestHandler_URLMappingPostSalesInvoiceFailDeleted(t *testing.T) {
	//buat dummy so
	so := model.DummySalesOrder()
	so.IsDeleted = int8(1)
	so.DocumentStatus = "active"
	so.InvoiceStatus = "new"
	so.TotalCharge = float64(20000)
	so.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"recognition_date": "2017-05-08T14:52:36+07:00",
		"sales_order_id":   common.Encrypt(so.ID),
		"due_date":         "2017-05-08T14:52:36+07:00",
		"billing_address":  "address",
		"total_amount":     float64(20000),
		"note":             "note",
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/sales-invoice").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPostSalesInvoiceFailTotalAmountMoreThan test create sales invoice dengan total amount lebih,fail
func TestHandler_URLMappingPostSalesInvoiceFailTotalAmountMoreThan(t *testing.T) {
	//buat dummy so
	so := model.DummySalesOrder()
	so.IsDeleted = int8(0)
	so.DocumentStatus = "active"
	so.InvoiceStatus = "new"
	so.TotalCharge = float64(20000)
	so.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"recognition_date": "2017-05-08T14:52:36+07:00",
		"sales_order_id":   common.Encrypt(so.ID),
		"due_date":         "2017-05-08T14:52:36+07:00",
		"billing_address":  "address",
		"total_amount":     float64(21000),
		"note":             "note",
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/sales-invoice").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPostSalesInvoiceFailTotalAmountMoreThan2 test create sales invoice dengan total amount lebih,fail
func TestHandler_URLMappingPostSalesInvoiceFailTotalAmountMoreThan2(t *testing.T) {
	//buat dummy so
	so := model.DummySalesOrder()
	so.IsDeleted = int8(0)
	so.InvoiceStatus = "active"
	so.DocumentStatus = "active"
	so.TotalCharge = float64(20000)
	so.Save()
	si := model.DummySalesInvoice()
	si.SalesOrder = so
	si.IsDeleted = int8(0)
	si.TotalAmount = float64(5000)
	si.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"recognition_date": "2017-05-08T14:52:36+07:00",
		"sales_order_id":   common.Encrypt(so.ID),
		"due_date":         "2017-05-08T14:52:36+07:00",
		"billing_address":  "address",
		"total_amount":     float64(16000),
		"note":             "note",
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/sales-invoice").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

////PUT///////////////////////////////////////

// TestHandler_URLMappingPutSalesInvoiceSuccess test update sales invoice,success
func TestHandler_URLMappingPutSalesInvoiceSuccess(t *testing.T) {
	//buat dummy si
	si := model.DummySalesInvoice()
	si.SalesOrder.TotalCharge = float64(20000)
	si.SalesOrder.DocumentStatus = "active"
	si.SalesOrder.Save()
	si.TotalAmount = float64(1000)
	si.DocumentStatus = "new"
	si.IsDeleted = int8(0)
	si.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"recognition_date": "2017-05-08T14:52:36+07:00",
		"note":             "note",
		"total_amount":     float64(2000),
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/sales-invoice/"+common.Encrypt(si.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	e := si.Read("SalesOrder")
	assert.NoError(t, e)
	assert.Equal(t, "note", si.Note)
	assert.Equal(t, sd.User.ID, si.UpdatedBy.ID)
}

// TestHandler_URLMappingPutSalesInvoiceFailIsDeleted test update sales invoice yang di delete,fail
func TestHandler_URLMappingPutSalesInvoiceFailIsDeleted(t *testing.T) {
	//buat dummy si
	si := model.DummySalesInvoice()
	si.SalesOrder.TotalCharge = float64(20000)
	si.SalesOrder.DocumentStatus = "active"
	si.SalesOrder.Save()
	si.DocumentStatus = "new"
	si.TotalAmount = float64(2000)
	si.IsDeleted = int8(1)
	si.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"recognition_date": "2017-05-08T14:52:36+07:00",
		"note":             "note",
		"total_amount":     float64(2000),
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/sales-invoice/"+common.Encrypt(si.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPutSalesInvoiceFailStatusNotNew test update sales invoice yang status tidak new,fail
func TestHandler_URLMappingPutSalesInvoiceFailStatusNotNew(t *testing.T) {
	//buat dummy si
	si := model.DummySalesInvoice()
	si.SalesOrder.TotalCharge = float64(20000)
	si.SalesOrder.DocumentStatus = "active"
	si.SalesOrder.Save()
	si.DocumentStatus = "active"
	si.TotalAmount = float64(1000)
	si.IsDeleted = int8(0)
	si.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"recognition_date": "2017-05-08T14:52:36+07:00",
		"note":             "note",
		"total_amount":     float64(2000),
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/sales-invoice/"+common.Encrypt(si.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMappingPutSalesInvoiceFailTotalAmountMoreThan2 test update sales invoice dengan total amount lebih,fail
func TestHandler_URLMappingPutSalesInvoiceFailTotalAmountMoreThan2(t *testing.T) {
	//buat dummy so
	so := model.DummySalesOrder()
	so.IsDeleted = int8(0)
	so.InvoiceStatus = "active"
	so.DocumentStatus = "active"
	so.TotalCharge = float64(20000)
	so.Save()
	si := model.DummySalesInvoice()
	si.SalesOrder = so
	si.IsDeleted = int8(0)
	si.DocumentStatus = "new"
	si.TotalAmount = float64(15000)
	si.Save()
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so
	si2.IsDeleted = int8(0)
	si2.DocumentStatus = "finished"
	si2.TotalAmount = float64(5000)
	si2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"recognition_date": "2017-05-08T14:52:36+07:00",
		"note":             "note",
		"total_amount":     float64(16000),
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/sales-invoice/"+common.Encrypt(si.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}
