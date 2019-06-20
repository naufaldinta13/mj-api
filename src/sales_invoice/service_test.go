// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package salesInvoice

import (
	"testing"
	"time"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

// TestGetSalesInvoiceDataEmpty test sales invoice dengan database empty
func TestGetSalesInvoiceDataEmpty(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM sales_invoice").Exec()

	rq := orm.RequestQuery{}
	// test tidak ada data sales invoice
	m, total, e := GetSalesInvoice(&rq)
	assert.Equal(t, int64(0), total)
	assert.NoError(t, e)
	assert.Empty(t, m)
}

// TestGetSalesInvoiceSuccess test sales invoice dengan 3 data di database (1 di delete),success
func TestGetSalesInvoiceSuccess(t *testing.T) {
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

	rq := orm.RequestQuery{}
	m, total, e := GetSalesInvoice(&rq)
	assert.NoError(t, e)
	assert.NotEmpty(t, m)
	assert.Equal(t, int64(2), total)
	for _, u := range *m {
		assert.Equal(t, "test", u.BillingAddress)
	}

}

// TestShowSalesInvoiceSuccess test show sales invoice, success
func TestShowSalesInvoiceSuccess(t *testing.T) {
	// buat dummy
	si := model.DummySalesInvoice()
	si.IsDeleted = int8(0)
	si.Save()
	// test
	m, e := ShowSalesInvoice("id", si.ID)
	assert.NoError(t, e)
	assert.Equal(t, si.BillingAddress, m.BillingAddress)
	assert.Equal(t, si.ID, m.ID)
	assert.Equal(t, si.TotalPaid, m.TotalPaid)
	assert.Equal(t, si.SalesOrder.TotalCharge, m.SalesOrder.TotalCharge)
}

// TestShowSalesInvoiceFail test show sales invoice dengan id tidak ada, fail
func TestShowSalesInvoiceFail(t *testing.T) {
	// test
	m, e := ShowSalesInvoice("id", 99999999)
	assert.Error(t, e)
	assert.Empty(t, m)
}

// TestGetSumTotalAmountSalesInvoiceBySalesOrder test get sum total amount dari sales invoice
func TestGetSumTotalAmountSalesInvoiceBySalesOrder(t *testing.T) {
	so := model.DummySalesOrder()
	so.TotalCharge = float64(50000)
	so.InvoiceStatus = "new"
	so.IsDeleted = int8(0)
	so.Save()
	si := model.DummySalesInvoice()
	si.IsDeleted = int8(0)
	si.TotalAmount = float64(30000)
	si.SalesOrder = so
	si.Save()
	s2 := model.DummySalesInvoice()
	s2.IsDeleted = int8(1)
	s2.TotalAmount = float64(20000)
	s2.SalesOrder = so
	s2.Save()
	// test
	total, e := GetSumTotalAmountSalesInvoiceBySalesOrder(so.ID)
	assert.NoError(t, e)
	assert.Equal(t, float64(30000), total)
}

// TestGetSumTotalAmountSalesInvoiceBySalesOrder2 test get sum total amount dari sales invoice
func TestGetSumTotalAmountSalesInvoiceBySalesOrder2(t *testing.T) {
	so := model.DummySalesOrder()
	so.TotalCharge = float64(50000)
	so.InvoiceStatus = "new"
	so.IsDeleted = int8(0)
	so.Save()
	si := model.DummySalesInvoice()
	si.IsDeleted = int8(0)
	si.TotalAmount = float64(30000)
	si.SalesOrder = so
	si.Save()
	s2 := model.DummySalesInvoice()
	s2.IsDeleted = int8(1)
	s2.TotalAmount = float64(20000)
	s2.SalesOrder = so
	s2.Save()
	s3 := model.DummySalesInvoice()
	s3.IsDeleted = int8(0)
	s3.TotalAmount = float64(10000)
	s3.SalesOrder = so
	s3.Save()
	// test
	total, e := GetSumTotalAmountSalesInvoiceBySalesOrder(so.ID)
	assert.NoError(t, e)
	assert.Equal(t, float64(40000), total)
}

// TestCreateSalesInvoice1 test create SI menjadi status active
func TestCreateSalesInvoice1(t *testing.T) {
	so := model.DummySalesOrder()
	so.TotalCharge = float64(50000)
	so.InvoiceStatus = "new"
	so.IsDeleted = int8(0)
	so.Save()
	si := model.DummySalesInvoice()
	si.IsDeleted = int8(0)
	si.TotalAmount = float64(20000)
	si.SalesOrder = so
	si.Save()
	invoice := &model.SalesInvoice{
		SalesOrder:      so,
		Code:            "codeee",
		RecognitionDate: time.Now(),
		DueDate:         time.Now(),
		BillingAddress:  "address",
		TotalAmount:     float64(10000),
		Note:            "note",
		DocumentStatus:  "new",
		IsBundled:       int8(0),
		IsDeleted:       int8(0),
		CreatedBy:       so.CreatedBy,
		CreatedAt:       time.Now(),
	}
	// test
	e := CreateSalesInvoice(invoice)
	assert.NoError(t, e)
	so.Read("ID")
	assert.Equal(t, "active", so.InvoiceStatus)
}

// TestCreateSalesInvoice2 test create SI menjadi status finished
func TestCreateSalesInvoice2(t *testing.T) {
	so := model.DummySalesOrder()
	so.TotalCharge = float64(50000)
	so.InvoiceStatus = "new"
	so.IsDeleted = int8(0)
	so.Save()
	si := model.DummySalesInvoice()
	si.IsDeleted = int8(0)
	si.TotalAmount = float64(40000)
	si.SalesOrder = &model.SalesOrder{ID: so.ID}
	si.Save()
	invoice := &model.SalesInvoice{
		SalesOrder:      so,
		Code:            "codeee",
		RecognitionDate: time.Now(),
		DueDate:         time.Now(),
		BillingAddress:  "address",
		TotalAmount:     float64(10000),
		Note:            "note",
		DocumentStatus:  "new",
		IsBundled:       int8(0),
		IsDeleted:       int8(0),
		CreatedBy:       so.CreatedBy,
		CreatedAt:       time.Now(),
	}
	// test
	e := CreateSalesInvoice(invoice)
	assert.NoError(t, e)
	so.Read("ID")
	assert.Equal(t, "active", so.InvoiceStatus)
}

// TestCreateSalesInvoice3 test create SI menjadi status finished tanpa membuat SI
func TestCreateSalesInvoice3(t *testing.T) {
	so := model.DummySalesOrder()
	so.TotalCharge = float64(0)
	so.InvoiceStatus = "new"
	so.IsDeleted = int8(0)
	so.Save()
	si := &model.SalesInvoice{
		SalesOrder:      so,
		Code:            "codeee",
		RecognitionDate: time.Now(),
		DueDate:         time.Now(),
		BillingAddress:  "address",
		TotalAmount:     float64(0),
		Note:            "note",
		DocumentStatus:  "new",
		IsBundled:       int8(0),
		IsDeleted:       int8(0),
		CreatedBy:       so.CreatedBy,
		CreatedAt:       time.Now(),
	}
	// test
	e := CreateSalesInvoice(si)
	assert.NoError(t, e)
	so.Read("ID")
	assert.Equal(t, "finished", so.InvoiceStatus)
}

// TestSumTotalRevenuedSalesInvoice untuk mengetest sum total revenued sales invoice
func TestSumTotalRevenuedSalesInvoice(t *testing.T) {

	orm.NewOrm().Raw("delete from finance_revenue").Exec()
	orm.NewOrm().Raw("delete from invcoice_receipt").Exec()
	orm.NewOrm().Raw("delete from sales_invoice").Exec()

	// buat data sales invoice
	si := model.DummySalesInvoice()
	si.Save()

	// buat data finance revenue pertama
	fr := model.DummyFinanceRevenue()
	fr.RefID = uint64(si.ID)
	fr.RefType = "sales_invoice"
	fr.Amount = 10000
	fr.IsDeleted = 0
	fr.Save()

	// buat data finance revenue kedua
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si.ID)
	fr2.RefType = "sales_invoice"
	fr2.Amount = 10000
	fr2.IsDeleted = 0
	fr2.Save()

	// buat data invoice reveipt
	ir := model.DummyInvoiceReceipt()
	ir.TotalAmount = 10000
	ir.IsDeleted = 0
	ir.Save()

	// buat data invoice receipt item
	iri := model.DummyInvoiceReceiptItem()
	iri.SalesInvoice = si
	iri.InvoiceReceipt = ir
	iri.Save()

	amount, e := SumTotalRevenuedSalesInvoice(si.ID)
	assert.NoError(t, e)

	// cek total revenued sales invoice
	si.Read("ID")
	assert.Equal(t, amount, si.TotalRevenued)
}
