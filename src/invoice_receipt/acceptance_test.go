// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package invoiceReceipt_test

import (
	"fmt"
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

// TestHandler_URLMappingGetAllInvoiceReceiptFailNoToken untuk mengetest show invoice receipt tanpa token
func TestHandler_URLMappingGetAllInvoiceReceiptFailNoToken(t *testing.T) {
	// dummy sales invoice
	model.DummyInvoiceReceipt()

	ng := tester.New()
	ng.Method = "GET"
	ng.Path = "/v1/invoice-receipt"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/invoice-receipt", "GET"))
	})
}

// TestHandler_URLMappingGetAllInvoiceReceiptSuccess untuk mengetest show invoice receipt menggunakan id benar
func TestHandler_URLMappingGetAllInvoiceReceiptSuccess(t *testing.T) {
	// dummy sales invoice
	model.DummyInvoiceReceipt()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/invoice-receipt"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/invoice-receipt", "GET"))
	})
}

// TestHandler_URLMappingDetailInvoiceReceiptFailNoToken untuk mengetest show invoice receipt tanpa token
func TestHandler_URLMappingDetailInvoiceReceiptFailNoToken(t *testing.T) {
	// dummy sales invoice
	m := model.DummyInvoiceReceipt()

	ng := tester.New()
	ng.Method = "GET"
	ng.Path = "/v1/invoice-receipt/" + common.Encrypt(m.ID)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/invoice-receipt/"+common.Encrypt(m.ID), "GET"))
	})
}

// TestHandler_URLMappingDetailInvoiceReceiptSuccess untuk mengetest show invoice receipt menggunakan id benar
func TestHandler_URLMappingDetailInvoiceReceiptSuccess(t *testing.T) {
	// dummy sales invoice
	m := model.DummyInvoiceReceipt()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/invoice-receipt/" + common.Encrypt(m.ID)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/invoice-receipt/"+common.Encrypt(m.ID), "GET"))
	})
}

// TestHandler_URLMappingDetailInvoiceReceiptFailIDNotExist untuk mengetest show invoice receipt menggunakan id yg tidak ada
func TestHandler_URLMappingDetailInvoiceReceiptFailIDNotExist(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/invoice-receipt/9999999999"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/invoice-receipt/9999999999", "GET"))
	})
}

// TestHandler_URLMappingCreateInvoiceReceiptSuccess test membuat invoice receipt success
func TestHandler_URLMappingCreateInvoiceReceiptSuccess(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()

	// buat data customer
	customer := model.DummyPartnership()
	customer.IsArchived = int8(0)
	customer.IsDeleted = int8(0)
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.DocumentStatus = "requested_cancel"
	so1.InvoiceStatus = "active"
	so1.IsDeleted = int8(0)
	so1.Save()

	// buat data so kedua dengan dummy customer
	so2 := model.DummySalesOrder()
	so2.Customer = customer
	so2.DocumentStatus = "requested_cancel"
	so2.InvoiceStatus = "active"
	so2.IsDeleted = int8(0)
	so2.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "active"
	si1.TotalAmount = 90000
	si1.IsBundled = int8(0)
	si1.IsDeleted = int8(0)
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so2
	si2.DocumentStatus = "active"
	si2.TotalAmount = 90000
	si2.IsBundled = int8(0)
	si2.IsDeleted = int8(0)
	si2.Save()

	// buat data sales return
	sr1 := model.DummySalesReturn()
	sr1.SalesOrder = so1
	sr1.DocumentStatus = "active"
	sr1.TotalAmount = 10000
	sr1.IsBundled = int8(0)
	sr1.IsDeleted = int8(0)
	sr1.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.Amount = 10000
	fr1.IsDeleted = int8(0)
	fr1.Save()

	// buat data finance revenue kedua
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.Amount = 10000
	fr2.IsDeleted = int8(0)
	fr2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	tgl := time.Now()

	// setting body
	scenario := tester.D{
		"sales_order": []tester.D{
			{"id": common.Encrypt(so1.ID)},
			{"id": common.Encrypt(so2.ID)},
		},
		"partnership_id":   common.Encrypt(customer.ID),
		"recognition_date": tgl,
		"note":             "notes",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/invoice-receipt").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

	// cek data invoice receipt
	data := &model.InvoiceReceipt{Partnership: customer}
	data.Read("Partnership")
	assert.Equal(t, "notes", data.Note)
	assert.Equal(t, customer.ID, data.Partnership.ID)
	assert.Equal(t, sd.User.ID, data.CreatedBy.ID)

	// cek data invoice receipt item
	receiptItem := &model.InvoiceReceiptItem{InvoiceReceipt: data}
	receiptItem.Read("InvoiceReceipt")
	assert.Equal(t, data.ID, receiptItem.InvoiceReceipt.ID)
	assert.Equal(t, si1.ID, receiptItem.SalesInvoice.ID)
	assert.Equal(t, si1.TotalAmount-(fr1.Amount+fr2.Amount), receiptItem.Subtotal)
	//
	//// cek data sales invoice
	//si := &model.SalesInvoice{SalesOrder: so1}
	//si.Read("SalesOrder")
	//assert.Equal(t, "active", si.DocumentStatus)
	//assert.Equal(t, int8(1), si.IsBundled)
	//assert.Equal(t, data.TotalAmount+float64(20000), si.TotalRevenued)

	// cek data sales return
	sr := &model.SalesReturn{SalesOrder: so1}
	sr.Read("SalesOrder")
	assert.Equal(t, "active", sr.DocumentStatus)
	assert.Equal(t, int8(1), sr.IsBundled)
}

// TestHandler_URLMappingCreateInvoiceReceiptSuccessWithNotActiveDocumentStatusSalesInvoice test membuat invoice receipt dengan document status selain active pada sales invoice
func TestHandler_URLMappingCreateInvoiceReceiptSuccessWithNotActiveDocumentStatusSalesInvoice(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()

	// buat data customer
	customer := model.DummyPartnership()
	customer.IsArchived = int8(0)
	customer.IsDeleted = int8(0)
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "new"
	so1.IsDeleted = int8(0)
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.Amount = 10000
	fr1.Save()

	// buat data finance revenue kedua
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.Amount = 10000
	fr2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	tgl := time.Now()

	// setting body
	scenario := tester.D{
		"sales_order": []tester.D{
			{"id": common.Encrypt(so1.ID)},
		},
		"partnership_id":   common.Encrypt(customer.ID),
		"recognition_date": tgl,
		"note":             "notes",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/invoice-receipt").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateInvoiceReceiptFailInvoiceStatusInSONotActive test membuat invoice receipt dengan invoice status selain active pada sales order
func TestHandler_URLMappingCreateInvoiceReceiptFailInvoiceStatusInSONotActive(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()

	// buat data customer
	customer := model.DummyPartnership()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "new"
	so1.IsDeleted = int8(0)
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.Amount = 10000
	fr1.Save()

	// buat data finance revenue kedua
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.Amount = 10000
	fr2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	tgl := time.Now()

	// setting body
	scenario := tester.D{"partnership_id": common.Encrypt(customer.ID), "recognition_date": tgl, "note": "notes"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/invoice-receipt").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateInvoiceReceiptFailErrorDecryptPartnershipID test membuat invoice receipt error decrypt partnership
func TestHandler_URLMappingCreateInvoiceReceiptFailErrorDecryptPartnershipID(t *testing.T) {

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	tgl := time.Now()

	// setting body
	scenario := tester.D{"partnership_id": "qwewqee", "recognition_date": tgl, "note": "notes"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/invoice-receipt").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateInvoiceReceiptFailCustomerNotSame test membuat invoice receipt dengan beda customer
func TestHandler_URLMappingCreateInvoiceReceiptFailCustomerNotSame(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()

	// buat data partnerhsip
	partner := model.DummyPartnership()
	partner.IsArchived = int8(0)
	partner.IsDeleted = int8(0)
	partner.Save()

	// buat data customer
	customer := model.DummyPartnership()
	customer.IsArchived = int8(0)
	customer.IsDeleted = int8(0)
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "new"
	so1.IsDeleted = int8(0)
	so1.Save()

	// buat data so kedua dengan dummy customer
	so2 := model.DummySalesOrder()
	so2.Customer = partner
	so2.InvoiceStatus = "new"
	so2.IsDeleted = int8(0)
	so2.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so2
	si2.DocumentStatus = "new"
	si2.TotalAmount = 90000
	si2.IsBundled = 0
	si2.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.Amount = 10000
	fr1.Save()

	// buat data finance revenue kedua
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.Amount = 10000
	fr2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	tgl := time.Now()

	// setting body
	scenario := tester.D{
		"sales_order": []tester.D{
			{"id": common.Encrypt(so1.ID)},
			{"id": common.Encrypt(so2.ID)},
		},
		"partnership_id":   common.Encrypt(customer.ID),
		"recognition_date": tgl,
		"note":             "notes",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/invoice-receipt").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateInvoiceReceiptFailPartnershipNotFound test membuat invoice receipt partnership not found
func TestHandler_URLMappingCreateInvoiceReceiptFailPartnershipNotFound(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()

	// buat data customer
	customer := model.DummyPartnership()
	customer.IsArchived = int8(0)
	customer.IsDeleted = int8(0)
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "new"
	so1.IsDeleted = int8(0)
	so1.Save()

	// buat data so kedua dengan dummy customer
	so2 := model.DummySalesOrder()
	so2.Customer = customer
	so2.InvoiceStatus = "new"
	so2.IsDeleted = int8(0)
	so2.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so2
	si2.DocumentStatus = "new"
	si2.TotalAmount = 90000
	si2.IsBundled = 0
	si2.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.Amount = 10000
	fr1.Save()

	// buat data finance revenue kedua
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.Amount = 10000
	fr2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	tgl := time.Now()

	// setting body
	scenario := tester.D{
		"sales_order": []tester.D{
			{"id": common.Encrypt(so1.ID)},
			{"id": common.Encrypt(so2.ID)},
		},
		"partnership_id":   "9999",
		"recognition_date": tgl,
		"note":             "notes",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/invoice-receipt").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateInvoiceReceiptFailPartnership1 test membuat invoice receipt partnership ter-delete
func TestHandler_URLMappingCreateInvoiceReceiptFailPartnership1(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()

	// buat data customer
	customer := model.DummyPartnership()
	customer.IsArchived = int8(0)
	customer.IsDeleted = int8(1)
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "new"
	so1.IsDeleted = int8(0)
	so1.Save()

	// buat data so kedua dengan dummy customer
	so2 := model.DummySalesOrder()
	so2.Customer = customer
	so2.InvoiceStatus = "new"
	so2.IsDeleted = int8(0)
	so2.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so2
	si2.DocumentStatus = "new"
	si2.TotalAmount = 90000
	si2.IsBundled = 0
	si2.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.Amount = 10000
	fr1.Save()

	// buat data finance revenue kedua
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.Amount = 10000
	fr2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	tgl := time.Now()

	// setting body
	scenario := tester.D{
		"sales_order": []tester.D{
			{"id": common.Encrypt(so1.ID)},
			{"id": common.Encrypt(so2.ID)},
		},
		"partnership_id":   common.Encrypt(customer.ID),
		"recognition_date": tgl,
		"note":             "notes",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/invoice-receipt").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateInvoiceReceiptFailPartnership2 test membuat invoice receipt partnership ter-arhcive
func TestHandler_URLMappingCreateInvoiceReceiptFailPartnership2(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()

	// buat data customer
	customer := model.DummyPartnership()
	customer.IsArchived = int8(1)
	customer.IsDeleted = int8(0)
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "new"
	so1.IsDeleted = int8(0)
	so1.Save()

	// buat data so kedua dengan dummy customer
	so2 := model.DummySalesOrder()
	so2.Customer = customer
	so2.InvoiceStatus = "new"
	so2.IsDeleted = int8(0)
	so2.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so2
	si2.DocumentStatus = "new"
	si2.TotalAmount = 90000
	si2.IsBundled = 0
	si2.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.Amount = 10000
	fr1.Save()

	// buat data finance revenue kedua
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.Amount = 10000
	fr2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	tgl := time.Now()

	// setting body
	scenario := tester.D{
		"sales_order": []tester.D{
			{"id": common.Encrypt(so1.ID)},
			{"id": common.Encrypt(so2.ID)},
		},
		"partnership_id":   common.Encrypt(customer.ID),
		"recognition_date": tgl,
		"note":             "notes",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/invoice-receipt").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateInvoiceReceiptFailSalesOrderDecrypt test membuat invoice receipt gagal decrypt sales order
func TestHandler_URLMappingCreateInvoiceReceiptFailSalesOrderDecrypt(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()

	// buat data customer
	customer := model.DummyPartnership()
	customer.IsArchived = int8(0)
	customer.IsDeleted = int8(0)
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "new"
	so1.IsDeleted = int8(0)
	so1.Save()

	// buat data so kedua dengan dummy customer
	so2 := model.DummySalesOrder()
	so2.Customer = customer
	so2.InvoiceStatus = "new"
	so2.IsDeleted = int8(0)
	so2.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so2
	si2.DocumentStatus = "new"
	si2.TotalAmount = 90000
	si2.IsBundled = 1
	si2.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.Amount = 10000
	fr1.Save()

	// buat data finance revenue kedua
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.Amount = 10000
	fr2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	tgl := time.Now()

	// setting body
	scenario := tester.D{
		"sales_order": []tester.D{
			{"id": "0"},
			{"id": "0"},
		},
		"partnership_id":   common.Encrypt(customer.ID),
		"recognition_date": tgl,
		"note":             "notes",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/invoice-receipt").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateInvoiceReceiptFailSalesOrderNotFound test membuat invoice receipt sales order not found
func TestHandler_URLMappingCreateInvoiceReceiptFailSalesOrderNotFound(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()

	// buat data customer
	customer := model.DummyPartnership()
	customer.IsArchived = int8(0)
	customer.IsDeleted = int8(0)
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "new"
	so1.IsDeleted = int8(0)
	so1.Save()

	// buat data so kedua dengan dummy customer
	so2 := model.DummySalesOrder()
	so2.Customer = customer
	so2.InvoiceStatus = "new"
	so2.IsDeleted = int8(0)
	so2.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so2
	si2.DocumentStatus = "new"
	si2.TotalAmount = 90000
	si2.IsBundled = 0
	si2.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.Amount = 10000
	fr1.Save()

	// buat data finance revenue kedua
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.Amount = 10000
	fr2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	tgl := time.Now()

	// setting body
	scenario := tester.D{
		"sales_order": []tester.D{
			{"id": "9999"},
			{"id": "9999"},
		},
		"partnership_id":   common.Encrypt(customer.ID),
		"recognition_date": tgl,
		"note":             "notes",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/invoice-receipt").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateInvoiceReceiptPaymentSuccess test membuat invoice receipt payment dengan document status selain active pada sales invoice
func TestHandler_URLMappingCreateInvoiceReceiptPaymentSuccess(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()
	o.Raw("delete from finance_revenue").Exec()
	o.Raw("delete from invoice_receipt_return").Exec()

	// buat data customer
	customer := model.DummyPartnership()
	customer.TotalDebt = 300000
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "active"
	so1.TotalCharge = 200000
	so1.TotalPaid = 20000
	so1.IsDeleted = int8(0)
	so1.FulfillmentStatus = "finished"
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 90000
	si1.TotalPaid = 90000
	si1.IsBundled = 0
	si1.IsDeleted = 0
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so1
	si2.DocumentStatus = "new"
	si2.TotalAmount = 90000
	si2.TotalPaid = 90000

	si2.IsBundled = 0
	si2.IsDeleted = 0
	si2.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.PaymentMethod = "credit_card"
	fr1.DocumentStatus = "cleared"
	fr1.BankName = "asd"
	fr1.BankHolder = "asd"
	fr1.BankNumber = "213"
	fr1.Amount = 45000
	fr1.Save()

	// buat data finance revenue pertama
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.PaymentMethod = "credit_card"
	fr2.DocumentStatus = "cleared"
	fr2.BankName = "asd"
	fr2.BankHolder = "asd"
	fr2.BankNumber = "213"
	fr2.Amount = 45000
	fr2.Save()

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.TotalAmount = 90000
	ir.DocumentStatus = "new"
	ir.TotalInvoice = 160000
	ir.TotalReturn = 120000
	ir.Save()

	// buat data invoice receipt item
	irm := model.DummyInvoiceReceiptItem()
	irm.SalesInvoice = si1
	irm.InvoiceReceipt = ir
	irm.Subtotal = 50000
	irm.Save()

	// buat data invoice receipt item kedua
	irm2 := model.DummyInvoiceReceiptItem()
	irm2.SalesInvoice = si2
	irm2.InvoiceReceipt = ir
	irm2.Subtotal = 40000
	irm2.Save()

	// buat data invoice receipt item ketiga
	irm3 := model.DummyInvoiceReceiptItem()
	irm3.SalesInvoice = si2
	irm3.InvoiceReceipt = ir
	irm3.Subtotal = 80000
	irm3.Save()

	// buat data sales return
	sr := model.DummySalesReturn()
	sr.SalesOrder = so1
	sr.DocumentStatus = "active"
	sr.IsDeleted = 0
	sr.IsBundled = 0
	sr.Save()

	// buat data invoice receipt return
	irr := model.DummyInvoiceReceiptReturn()
	irr.InvoiceReceipt = ir
	irr.SalesReturn = sr
	irr.Subtotal = 40000
	irr.Save()

	// buat data invoice receipt return kedua
	irr2 := model.DummyInvoiceReceiptReturn()
	irr2.InvoiceReceipt = ir
	irr2.SalesReturn = sr
	irr2.Subtotal = 40000
	irr2.Save()

	// buat data invoice receipt return ketiga
	irr3 := model.DummyInvoiceReceiptReturn()
	irr3.InvoiceReceipt = ir
	irr3.SalesReturn = sr
	irr3.Subtotal = 60000
	irr3.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"invoice_receipt_item": []tester.D{
			{"id": common.Encrypt(irm.ID)},
			{"id": common.Encrypt(irm2.ID)},
		},
		"invoice_receipt_return": []tester.D{
			{"id": common.Encrypt(irr.ID)},
			{"id": common.Encrypt(irr2.ID)},
		},
		"finance_revenue": []tester.D{
			{"ref_id": uint64(si1.ID), "ref_type": "sales_invoice", "payment_method": "credit_card", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
			{"ref_id": uint64(si2.ID), "ref_type": "sales_invoice", "payment_method": "credit_card", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/invoice-receipt/"+common.Encrypt(ir.ID)+"/payment").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

	// cek document status sales invoice pertama
	si1.Read()
	assert.Equal(t, "finished", si1.DocumentStatus)

	// cek document status sales invoice kedua
	si2.Read()
	assert.Equal(t, "finished", si2.DocumentStatus)

	// cek total paid so
	so1.Read()
	assert.Equal(t, float64(200000), so1.TotalPaid)

	// cek invoice status pada so
	assert.Equal(t, "finished", so1.InvoiceStatus)

	// cek apabila invoice_status dan fulfillment_status pada so finished
	// maka document_statusnya menjadi finished
	so1.Read()
	assert.Equal(t, "finished", so1.DocumentStatus)

	// cek total debt partnership
	customer.Read()
	assert.Equal(t, float64(100000), customer.TotalDebt)

	//cek total invoice pada invoice receipt
	ir.Read()
	assert.Equal(t, float64(80000), ir.TotalInvoice)

	// cek total return pada invoice receipt
	assert.Equal(t, float64(60000), ir.TotalReturn)

	// cek total amount pada invoice receipt
	assert.Equal(t, ir.TotalInvoice-ir.TotalReturn, ir.TotalAmount)
}

// TestHandler_URLMappingCreateInvoiceReceiptPaymentSuccess2 test membuat invoice receipt payment dengan document status active pada sales invoice
func TestHandler_URLMappingCreateInvoiceReceiptPaymentSuccess2(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()
	o.Raw("delete from finance_revenue").Exec()

	// buat data customer
	customer := model.DummyPartnership()
	customer.TotalDebt = 300000
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "new"
	so1.TotalCharge = 9900000
	so1.TotalPaid = 20000
	so1.IsDeleted = int8(0)
	so1.FulfillmentStatus = "finished"
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "active"
	si1.TotalAmount = 90000
	si1.TotalPaid = 90000

	si1.IsBundled = 0
	si1.IsDeleted = 0
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so1
	si2.DocumentStatus = "active"
	si2.TotalAmount = 90000
	si2.TotalPaid = 90000
	si2.IsBundled = 0
	si2.IsDeleted = 0
	si2.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.PaymentMethod = "credit_card"
	fr1.BankName = "asd"
	fr1.BankHolder = "asd"
	fr1.BankNumber = "213"
	fr1.Amount = 45000
	fr1.Save()

	// buat data finance revenue pertama
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.PaymentMethod = "credit_card"
	fr2.BankName = "asd"
	fr2.BankHolder = "asd"
	fr2.BankNumber = "213"
	fr2.Amount = 45000
	fr2.Save()

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.TotalAmount = 90000
	ir.DocumentStatus = "new"
	ir.TotalReturn = 100000000
	ir.IsDeleted = 0
	ir.TotalInvoice = 160000
	ir.Save()

	// buat data invoice receipt item
	irm := model.DummyInvoiceReceiptItem()
	irm.SalesInvoice = si1
	irm.InvoiceReceipt = ir
	irm.Subtotal = 50000
	irm.Save()

	// buat data invoice receipt item kedua
	irm2 := model.DummyInvoiceReceiptItem()
	irm2.SalesInvoice = si2
	irm2.InvoiceReceipt = ir
	irm2.Subtotal = 40000
	irm2.Save()

	// buat data invoice receipt item ketiga
	irm3 := model.DummyInvoiceReceiptItem()
	irm3.SalesInvoice = si2
	irm3.InvoiceReceipt = ir
	irm3.Subtotal = 80000
	irm3.Save()

	// buat data sales return
	sr := model.DummySalesReturn()
	sr.SalesOrder = so1
	sr.DocumentStatus = "active"
	sr.IsDeleted = 0
	sr.IsBundled = 0
	sr.Save()

	// buat data invoice receipt return
	irr := model.DummyInvoiceReceiptReturn()
	irr.InvoiceReceipt = ir
	irr.SalesReturn = sr
	irr.Subtotal = 40000
	irr.Save()

	// buat data invoice receipt return kedua
	irr2 := model.DummyInvoiceReceiptReturn()
	irr2.InvoiceReceipt = ir
	irr2.SalesReturn = sr
	irr2.Subtotal = 40000
	irr2.Save()

	// buat data invoice receipt return ketiga
	irr3 := model.DummyInvoiceReceiptReturn()
	irr3.InvoiceReceipt = ir
	irr3.SalesReturn = sr
	irr3.Subtotal = 60000
	irr3.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"invoice_receipt_item": []tester.D{
			{"id": common.Encrypt(irm.ID)},
			{"id": common.Encrypt(irm2.ID)},
		},
		"invoice_receipt_return": []tester.D{
			{"id": common.Encrypt(irr.ID)},
			{"id": common.Encrypt(irr2.ID)},
		},
		"finance_revenue": []tester.D{
			{"ref_id": uint64(si1.ID), "ref_type": "sales_invoice", "payment_method": "credit_card", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
			{"ref_id": uint64(si2.ID), "ref_type": "sales_invoice", "payment_method": "credit_card", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/invoice-receipt/"+common.Encrypt(ir.ID)+"/payment").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

	// cek document status sales invoice pertama
	si1.Read()
	assert.Equal(t, "finished", si1.DocumentStatus)

	// cek document status sales invoice kedua
	si2.Read()
	assert.Equal(t, "finished", si2.DocumentStatus)

	// cek total paid so
	so1.Read()
	assert.Equal(t, float64(200000), so1.TotalPaid)

	// cek invoice status pada so
	assert.Equal(t, "active", so1.InvoiceStatus)

	// cek apabila invoice_status dan fulfillment_status pada so finished
	// maka document_statusnya menjadi finished
	so1.Read()
	assert.Equal(t, "active", so1.DocumentStatus)

	// cek total debt partnership
	customer.Read()
	assert.Equal(t, float64(100000), customer.TotalDebt)

	//cek total invoice pada invoice receipt
	ir.Read()
	assert.Equal(t, float64(80000), ir.TotalInvoice)

	// cek total amount pada invoice receipt
	assert.Equal(t, ir.TotalInvoice-ir.TotalReturn, ir.TotalAmount)
}

// TestHandler_URLMappingCreateInvoiceReceiptPaymentFailDocumentStatusInvoiceReceipt test membuat invoice receipt payment dengan document status invoice receipt selain new
func TestHandler_URLMappingCreateInvoiceReceiptPaymentFailDocumentStatusInvoiceReceipt(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()
	o.Raw("delete from finance_revenue").Exec()

	// buat data customer
	customer := model.DummyPartnership()
	customer.TotalDebt = 300000
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "active"
	so1.TotalCharge = 200000
	so1.TotalPaid = 20000
	so1.IsDeleted = int8(0)
	so1.InvoiceStatus = "new"
	so1.FulfillmentStatus = "finished"
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so1
	si2.DocumentStatus = "new"
	si2.TotalAmount = 90000
	si2.IsBundled = 0
	si2.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.PaymentMethod = "credit_card"
	fr1.BankName = "asd"
	fr1.BankHolder = "asd"
	fr1.BankNumber = "213"
	fr1.Amount = 45000
	fr1.Save()

	// buat data finance revenue pertama
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.PaymentMethod = "credit_card"
	fr2.BankName = "asd"
	fr2.BankHolder = "asd"
	fr2.BankNumber = "213"
	fr2.Amount = 45000
	fr2.Save()

	// buat data sales return
	sr := model.DummySalesReturn()
	sr.SalesOrder = so1
	sr.DocumentStatus = "active"
	sr.IsDeleted = 0
	sr.Save()

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.TotalAmount = 90000
	ir.DocumentStatus = "active"
	ir.Save()

	// buat data invoice receipt return
	irr := model.DummyInvoiceReceiptReturn()
	irr.InvoiceReceipt = ir
	irr.SalesReturn = sr
	irr.Subtotal = 40000
	irr.Save()

	// buat data invoice receipt return kedua
	irr2 := model.DummyInvoiceReceiptReturn()
	irr2.InvoiceReceipt = ir
	irr2.SalesReturn = sr
	irr2.Subtotal = 40000
	irr2.Save()

	// buat data invoice receipt item
	irm := model.DummyInvoiceReceiptItem()
	irm.SalesInvoice = si1
	irm.InvoiceReceipt = ir
	irm.Subtotal = 50000
	irm.Save()

	// buat data invoice receipt item kedua
	irm2 := model.DummyInvoiceReceiptItem()
	irm2.SalesInvoice = si2
	irm2.InvoiceReceipt = ir
	irm2.Subtotal = 40000
	irm2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"invoice_receipt_item": []tester.D{
			{"id": common.Encrypt(irm.ID)},
			{"id": common.Encrypt(irm2.ID)},
		},
		"invoice_receipt_return": []tester.D{
			{"id": common.Encrypt(irr.ID)},
			{"id": common.Encrypt(irr2.ID)},
		},
		"finance_revenue": []tester.D{
			{"ref_id": uint64(si1.ID), "ref_type": "sales_invoice", "payment_method": "credit_card", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
			{"ref_id": uint64(si2.ID), "ref_type": "sales_invoice", "payment_method": "credit_card", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/invoice-receipt/"+common.Encrypt(ir.ID)+"/payment").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateInvoiceReceiptPaymentFailInvoiceReceiptNotFound test membuat invoice receipt payment dengan invoice receipt yang tidak ada datanya
func TestHandler_URLMappingCreateInvoiceReceiptPaymentFailInvoiceReceiptNotFound(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()
	o.Raw("delete from finance_revenue").Exec()

	// buat data customer
	customer := model.DummyPartnership()
	customer.TotalDebt = 300000
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "active"
	so1.TotalCharge = 200000
	so1.TotalPaid = 20000
	so1.IsDeleted = int8(0)
	so1.InvoiceStatus = "new"
	so1.FulfillmentStatus = "finished"
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so1
	si2.DocumentStatus = "new"
	si2.TotalAmount = 90000
	si2.IsBundled = 0
	si2.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.PaymentMethod = "credit_card"
	fr1.BankName = "asd"
	fr1.BankHolder = "asd"
	fr1.BankNumber = "213"
	fr1.Amount = 45000
	fr1.Save()

	// buat data finance revenue pertama
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.PaymentMethod = "credit_card"
	fr2.BankName = "asd"
	fr2.BankHolder = "asd"
	fr2.BankNumber = "213"
	fr2.Amount = 45000
	fr2.Save()

	// buat data sales return
	sr := model.DummySalesReturn()
	sr.SalesOrder = so1
	sr.DocumentStatus = "active"
	sr.IsDeleted = 0
	sr.Save()

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.TotalAmount = 90000
	ir.DocumentStatus = "new"
	ir.Save()

	// buat data invoice receipt return
	irr := model.DummyInvoiceReceiptReturn()
	irr.InvoiceReceipt = ir
	irr.SalesReturn = sr
	irr.Subtotal = 40000
	irr.Save()

	// buat data invoice receipt return kedua
	irr2 := model.DummyInvoiceReceiptReturn()
	irr2.InvoiceReceipt = ir
	irr2.SalesReturn = sr
	irr2.Subtotal = 40000
	irr2.Save()

	// buat data invoice receipt item
	irm := model.DummyInvoiceReceiptItem()
	irm.SalesInvoice = si1
	irm.InvoiceReceipt = ir
	irm.Subtotal = 50000
	irm.Save()

	// buat data invoice receipt item kedua
	irm2 := model.DummyInvoiceReceiptItem()
	irm2.SalesInvoice = si2
	irm2.InvoiceReceipt = ir
	irm.Subtotal = 40000
	irm2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"invoice_receipt_item": []tester.D{
			{"id": common.Encrypt(irm.ID)},
			{"id": common.Encrypt(irm2.ID)},
		},
		"invoice_receipt_return": []tester.D{
			{"id": common.Encrypt(irr.ID)},
			{"id": common.Encrypt(irr2.ID)},
		},
		"finance_revenue": []tester.D{
			{"ref_id": uint64(si1.ID), "ref_type": "sales_invoice", "payment_method": "credit_card", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
			{"ref_id": uint64(si2.ID), "ref_type": "sales_invoice", "payment_method": "credit_card", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/invoice-receipt/99999999999999/payment").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateInvoiceReceiptPaymentFailInvoiceReceiptItemNotFound test membuat invoice receipt payment dengan invoice receipt item yang tidak ada datanya
func TestHandler_URLMappingCreateInvoiceReceiptPaymentFailInvoiceReceiptItemNotFound(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()
	o.Raw("delete from finance_revenue").Exec()

	// buat data customer
	customer := model.DummyPartnership()
	customer.TotalDebt = 300000
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "active"
	so1.TotalCharge = 200000
	so1.TotalPaid = 20000
	so1.IsDeleted = int8(0)
	so1.InvoiceStatus = "new"
	so1.FulfillmentStatus = "finished"
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so1
	si2.DocumentStatus = "new"
	si2.TotalAmount = 90000
	si2.IsBundled = 0
	si2.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.PaymentMethod = "credit_card"
	fr1.BankName = "asd"
	fr1.BankHolder = "asd"
	fr1.BankNumber = "213"
	fr1.Amount = 45000
	fr1.Save()

	// buat data finance revenue pertama
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.PaymentMethod = "credit_card"
	fr2.BankName = "asd"
	fr2.BankHolder = "asd"
	fr2.BankNumber = "213"
	fr2.Amount = 45000
	fr2.Save()

	// buat data sales return
	sr := model.DummySalesReturn()
	sr.SalesOrder = so1
	sr.DocumentStatus = "active"
	sr.IsDeleted = 0
	sr.Save()

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.TotalAmount = 90000
	ir.DocumentStatus = "new"
	ir.Save()

	// buat data invoice receipt return
	irr := model.DummyInvoiceReceiptReturn()
	irr.InvoiceReceipt = ir
	irr.SalesReturn = sr
	irr.Subtotal = 40000
	irr.Save()

	// buat data invoice receipt return kedua
	irr2 := model.DummyInvoiceReceiptReturn()
	irr2.InvoiceReceipt = ir
	irr2.SalesReturn = sr
	irr2.Subtotal = 40000
	irr2.Save()

	// buat data invoice receipt item
	irm := model.DummyInvoiceReceiptItem()
	irm.SalesInvoice = si1
	irm.InvoiceReceipt = ir
	irm.Subtotal = 50000
	irm.Save()

	// buat data invoice receipt item kedua
	irm2 := model.DummyInvoiceReceiptItem()
	irm2.SalesInvoice = si2
	irm2.InvoiceReceipt = ir
	irm2.Subtotal = 40000
	irm2.Save()
	irm.ID = 5
	irm.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"invoice_receipt_item": []tester.D{
			{"id": common.Encrypt(irm.ID)},
			{"id": common.Encrypt(irm2.ID)},
		},
		"invoice_receipt_return": []tester.D{
			{"id": common.Encrypt(irr.ID)},
			{"id": common.Encrypt(irr2.ID)},
		},
		"finance_revenue": []tester.D{
			{"ref_id": uint64(si1.ID), "ref_type": "sales_invoice", "payment_method": "credit_card", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
			{"ref_id": uint64(si2.ID), "ref_type": "sales_invoice", "payment_method": "credit_card", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/invoice-receipt/"+common.Encrypt(ir.ID)+"/payment").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateInvoiceReceiptPaymentFailInvoiceReceiptItemErrorDecrypt test membuat invoice receipt payment dengan invoice receipt item gagal decrypt
func TestHandler_URLMappingCreateInvoiceReceiptPaymentFailInvoiceReceiptItemErrorDecrypt(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()
	o.Raw("delete from finance_revenue").Exec()

	// buat data customer
	customer := model.DummyPartnership()
	customer.TotalDebt = 300000
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "active"
	so1.TotalCharge = 200000
	so1.TotalPaid = 20000
	so1.IsDeleted = int8(0)
	so1.InvoiceStatus = "new"
	so1.FulfillmentStatus = "finished"
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so1
	si2.DocumentStatus = "new"
	si2.TotalAmount = 90000
	si2.IsBundled = 0
	si2.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.PaymentMethod = "credit_card"
	fr1.BankName = "asd"
	fr1.BankHolder = "asd"
	fr1.BankNumber = "213"
	fr1.Amount = 45000
	fr1.Save()

	// buat data finance revenue pertama
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.PaymentMethod = "credit_card"
	fr2.BankName = "asd"
	fr2.BankHolder = "asd"
	fr2.BankNumber = "213"
	fr2.Amount = 45000
	fr2.Save()

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.TotalAmount = 90000
	ir.DocumentStatus = "new"
	ir.Save()

	// buat data sales return
	sr := model.DummySalesReturn()
	sr.SalesOrder = so1
	sr.DocumentStatus = "active"
	sr.IsDeleted = 0
	sr.Save()

	// buat data invoice receipt return
	irr := model.DummyInvoiceReceiptReturn()
	irr.InvoiceReceipt = ir
	irr.SalesReturn = sr
	irr.Subtotal = 40000
	irr.Save()

	// buat data invoice receipt return kedua
	irr2 := model.DummyInvoiceReceiptReturn()
	irr2.InvoiceReceipt = ir
	irr2.SalesReturn = sr
	irr2.Subtotal = 40000
	irr2.Save()

	// buat data invoice receipt item
	irm := model.DummyInvoiceReceiptItem()
	irm.SalesInvoice = si1
	irm.InvoiceReceipt = ir
	irm.Subtotal = 50000
	irm.Save()

	// buat data invoice receipt item kedua
	irm2 := model.DummyInvoiceReceiptItem()
	irm2.SalesInvoice = si2
	irm2.InvoiceReceipt = ir
	irm2.Subtotal = 40000
	irm2.Save()

	irm.Delete()

	si1.Delete()
	si2.Delete()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"invoice_receipt_item": []tester.D{
			{"id": common.Encrypt(irm.ID)},
			{"id": common.Encrypt(irm2.ID)},
		},
		"invoice_receipt_return": []tester.D{
			{"id": common.Encrypt(irr.ID)},
			{"id": common.Encrypt(irr2.ID)},
		},
		"finance_revenue": []tester.D{
			{"ref_id": uint64(si1.ID), "ref_type": "sales_invoice", "payment_method": "credit_card", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
			{"ref_id": uint64(si2.ID), "ref_type": "sales_invoice", "payment_method": "credit_card", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/invoice-receipt/"+common.Encrypt(ir.ID)+"/payment").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateInvoiceReceiptPaymentFailTotalAmount test membuat invoice receipt payment dengan jumlah amount dan total amount invoice receipt tidak sama
func TestHandler_URLMappingCreateInvoiceReceiptPaymentFailTotalAmount(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()
	o.Raw("delete from finance_revenue").Exec()

	// buat data customer
	customer := model.DummyPartnership()
	customer.TotalDebt = 300000
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "active"
	so1.TotalCharge = 200000
	so1.TotalPaid = 20000
	so1.IsDeleted = int8(0)
	so1.InvoiceStatus = "new"
	so1.FulfillmentStatus = "finished"
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so1
	si2.DocumentStatus = "new"
	si2.TotalAmount = 90000
	si2.IsBundled = 0
	si2.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.PaymentMethod = "credit_card"
	fr1.BankName = "asd"
	fr1.BankHolder = "asd"
	fr1.BankNumber = "213"
	fr1.Amount = 45000
	fr1.Save()

	// buat data finance revenue pertama
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.PaymentMethod = "credit_card"
	fr2.BankName = "asd"
	fr2.BankHolder = "asd"
	fr2.BankNumber = "213"
	fr2.Amount = 45000
	fr2.Save()

	// buat data sales return
	sr := model.DummySalesReturn()
	sr.SalesOrder = so1
	sr.DocumentStatus = "active"
	sr.IsDeleted = 0
	sr.Save()

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.TotalAmount = 900000
	ir.DocumentStatus = "new"
	ir.Save()

	// buat data invoice receipt return
	irr := model.DummyInvoiceReceiptReturn()
	irr.InvoiceReceipt = ir
	irr.SalesReturn = sr
	irr.Subtotal = 40000
	irr.Save()

	// buat data invoice receipt return kedua
	irr2 := model.DummyInvoiceReceiptReturn()
	irr2.InvoiceReceipt = ir
	irr2.SalesReturn = sr
	irr2.Subtotal = 40000
	irr2.Save()

	// buat data invoice receipt item
	irm := model.DummyInvoiceReceiptItem()
	irm.SalesInvoice = si1
	irm.InvoiceReceipt = ir
	irm.Subtotal = 50000
	irm.Save()

	// buat data invoice receipt item kedua
	irm2 := model.DummyInvoiceReceiptItem()
	irm2.SalesInvoice = si2
	irm2.InvoiceReceipt = ir
	irm2.Subtotal = 50000
	irm2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"invoice_receipt_item": []tester.D{
			{"id": common.Encrypt(irm.ID)},
			{"id": common.Encrypt(irm2.ID)},
		},
		"invoice_receipt_return": []tester.D{
			{"id": common.Encrypt(irr.ID)},
			{"id": common.Encrypt(irr2.ID)},
		},
		"finance_revenue": []tester.D{
			{"ref_id": uint64(si1.ID), "ref_type": "sales_invoice", "payment_method": "credit_card", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
			{"ref_id": uint64(si2.ID), "ref_type": "sales_invoice", "payment_method": "credit_card", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/invoice-receipt/"+common.Encrypt(ir.ID)+"/payment").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateInvoiceReceiptPaymentFailPaymentMethod test membuat invoice receipt payment fail payment method
// bank_name, bank_holder, bank_number hanya bisa diisi jika payment methodnya debit_card atau credit_card
func TestHandler_URLMappingCreateInvoiceReceiptPaymentFailPaymentMethod(t *testing.T) {
	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()
	o.Raw("delete from finance_revenue").Exec()

	// buat data customer
	customer := model.DummyPartnership()
	customer.TotalDebt = 300000
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "active"
	so1.TotalCharge = 200000
	so1.TotalPaid = 20000
	so1.IsDeleted = int8(0)
	so1.InvoiceStatus = "new"
	so1.FulfillmentStatus = "finished"
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 90000
	si1.IsBundled = int8(1)
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so1
	si2.DocumentStatus = "new"
	si2.TotalAmount = 90000
	si2.IsBundled = int8(1)
	si2.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.PaymentMethod = "giro"
	fr1.BankName = "asd"
	fr1.BankHolder = "asd"
	fr1.BankNumber = "213"
	fr1.Amount = 45000
	fr1.Save()

	// buat data finance revenue pertama
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.PaymentMethod = "giro"
	fr2.BankName = "asd"
	fr2.BankHolder = "asd"
	fr2.BankNumber = "213"
	fr2.Amount = 45000
	fr2.Save()

	// buat data sales return
	sr := model.DummySalesReturn()
	sr.SalesOrder = so1
	sr.DocumentStatus = "active"
	sr.IsDeleted = 0
	sr.IsBundled = 1
	sr.Save()

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.TotalAmount = 90000
	ir.DocumentStatus = "new"
	ir.Save()

	// buat data invoice receipt return
	irr := model.DummyInvoiceReceiptReturn()
	irr.InvoiceReceipt = ir
	irr.SalesReturn = sr
	irr.Subtotal = 40000
	irr.Save()

	// buat data invoice receipt return kedua
	irr2 := model.DummyInvoiceReceiptReturn()
	irr2.InvoiceReceipt = ir
	irr2.SalesReturn = sr
	irr2.Subtotal = 40000
	irr2.Save()

	// buat data invoice receipt item
	irm := model.DummyInvoiceReceiptItem()
	irm.SalesInvoice = si1
	irm.InvoiceReceipt = ir
	irm.Subtotal = 50000
	irm.Save()

	// buat data invoice receipt item kedua
	irm2 := model.DummyInvoiceReceiptItem()
	irm2.SalesInvoice = si2
	irm2.InvoiceReceipt = ir
	irm2.Subtotal = 40000
	irm2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"invoice_receipt_item": []tester.D{
			{"id": common.Encrypt(irm.ID)},
			{"id": common.Encrypt(irm2.ID)},
		},
		"invoice_receipt_return": []tester.D{
			{"id": common.Encrypt(irr.ID)},
			{"id": common.Encrypt(irr2.ID)},
		},
		"finance_revenue": []tester.D{
			{"ref_id": uint64(si1.ID), "ref_type": "sales_invoice", "payment_method": "cash", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
			{"ref_id": uint64(si2.ID), "ref_type": "sales_invoice", "payment_method": "cash", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/invoice-receipt/"+common.Encrypt(ir.ID)+"/payment").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

//  test membuat invoice receipt payment dengan document status sales invoice telah finished
func TestHandler_URLMappingCreateInvoiceReceiptPaymentFailedIfDocumentStatusSalesInvoiceHasBeenFinished(t *testing.T) {

	// clear database
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()
	o.Raw("delete from invoice_receipt").Exec()
	o.Raw("delete from invoice_receipt_item").Exec()
	o.Raw("delete from finance_revenue").Exec()

	// buat data customer
	customer := model.DummyPartnership()
	customer.TotalDebt = 300000
	customer.Save()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "active"
	so1.TotalCharge = 200000
	so1.TotalPaid = 20000
	so1.IsDeleted = int8(0)
	so1.InvoiceStatus = "new"
	so1.FulfillmentStatus = "finished"
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.DocumentStatus = "finished"
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so1
	si2.TotalAmount = 90000
	si2.IsBundled = 0
	si2.DocumentStatus = "finished"
	si2.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.PaymentMethod = "credit_card"
	fr1.BankName = "asd"
	fr1.BankHolder = "asd"
	fr1.BankNumber = "213"
	fr1.Amount = 45000
	fr1.Save()

	// buat data finance revenue pertama
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.PaymentMethod = "credit_card"
	fr2.BankName = "asd"
	fr2.BankHolder = "asd"
	fr2.BankNumber = "213"
	fr2.Amount = 45000
	fr2.Save()

	// buat data sales return
	sr := model.DummySalesReturn()
	sr.SalesOrder = so1
	sr.DocumentStatus = "active"
	sr.IsDeleted = 0
	sr.Save()

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.TotalAmount = 90000
	ir.DocumentStatus = "new"
	ir.Save()

	// buat data invoice receipt return
	irr := model.DummyInvoiceReceiptReturn()
	irr.InvoiceReceipt = ir
	irr.SalesReturn = sr
	irr.Subtotal = 40000
	irr.Save()

	// buat data invoice receipt return kedua
	irr2 := model.DummyInvoiceReceiptReturn()
	irr2.InvoiceReceipt = ir
	irr2.SalesReturn = sr
	irr2.Subtotal = 40000
	irr2.Save()

	// buat data invoice receipt item
	irm := model.DummyInvoiceReceiptItem()
	irm.SalesInvoice = si1
	irm.InvoiceReceipt = ir
	irm.Subtotal = 50000
	irm.Save()

	// buat data invoice receipt item kedua
	irm2 := model.DummyInvoiceReceiptItem()
	irm2.SalesInvoice = si2
	irm2.InvoiceReceipt = ir
	irm2.Subtotal = 40000
	irm2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"invoice_receipt_item": []tester.D{
			{"id": common.Encrypt(irm.ID)},
			{"id": common.Encrypt(irm2.ID)},
		},
		"invoice_receipt_return": []tester.D{
			{"id": common.Encrypt(irr.ID)},
			{"id": common.Encrypt(irr2.ID)},
		},
		"finance_revenue": []tester.D{
			{"ref_id": uint64(si1.ID), "ref_type": "sales_invoice", "payment_method": "credit_card", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
			{"ref_id": uint64(si2.ID), "ref_type": "sales_invoice", "payment_method": "credit_card", "bank_name": "asd", "bank_holder": "asd", "bank_number": "123", "amount": float64(45000)},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/invoice-receipt/"+common.Encrypt(ir.ID)+"/payment").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}
