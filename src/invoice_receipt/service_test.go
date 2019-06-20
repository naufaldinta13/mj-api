// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package invoiceReceipt

import (
	"testing"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceReceiptsNoData(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM invoice_receipt").Exec()
	rq := orm.RequestQuery{}

	m, total, e := GetInvoiceReceipts(&rq)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, m)
	assert.NoError(t, e)
}

func TestGetInvoiceReceiptsSuccess(t *testing.T) {
	model.DummyInvoiceReceipt()

	qs := orm.RequestQuery{}
	_, _, e := GetInvoiceReceipts(&qs)
	assert.NoError(t, e, "Data should be exists.")
}

func TestGetDetailInvoiceReceipt(t *testing.T) {

	// seharusnya error karena tidak ada data dengan id 999999
	_, err := GetDetailInvoiceReceipt("id", 999999)
	assert.Error(t, err, "SHould be an error cause data with id 999999 doesn't exist")

	// seharusnya tidak error karena ada datanya
	c := model.DummyInvoiceReceipt()
	cd, e := GetDetailInvoiceReceipt("id", c.ID)
	assert.NoError(t, e, "Data should be exists.")
	assert.Equal(t, c.ID, cd.ID, "ID Response should be a same.")
}

func TestGetDetailInvoiceReceipts(t *testing.T) {

	// seharusnya error karena tidak ada data dengan id 999999
	_, err := getDetailInvoiceReceipt("id", 999999)
	assert.Error(t, err, "SHould be an error cause data with id 999999 doesn't exist")

	// seharusnya tidak error karena ada datanya
	c := model.DummyInvoiceReceipt()
	cd, e := getDetailInvoiceReceipt("id", c.ID)
	assert.NoError(t, e, "Data should be exists.")
	assert.Equal(t, c.ID, cd.ID, "ID Response should be a same.")
}

func TestGetSalesOrdersByCustomerIDSuccess(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()

	// buat data customer
	customer := model.DummyPartnership()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.Save()

	// buat data so kedua dengan dummy customer
	so2 := model.DummySalesOrder()
	so2.Customer = customer
	so2.Save()

	// get data so berdasarkan dummy customer
	m, e := GetSalesOrdersByCustomerID(customer.ID)
	// seharusnya tidak error karena ada datanya
	assert.NoError(t, e)
	// seharusnya data yang didapat jumlahnya ada 2
	assert.Equal(t, 2, len(m), "seharusnya data jumlahnya ada 2")
}

func TestGetSalesOrdersByCustomerIDFailNoData(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()

	// get data so berdasarkan dummy customer
	m, _ := GetSalesOrdersByCustomerID(99999999999)
	// seharusnya jumlah data 0 karena tidak ada datanya
	assert.Equal(t, 0, len(m))
}

func TestGetSalesInvoiceBySOIDSuccess(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("delete from sales_invoice").Exec()

	// buat data customer
	customer := model.DummyPartnership()

	// buat data so dengan dummy customer
	so := model.DummySalesOrder()
	so.Customer = customer
	so.Save()

	// buat data si pertama dengan dummy customer
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so
	si1.DocumentStatus = "new"
	si1.Save()

	// buat data si kedua dengan dummy customer
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so
	si1.DocumentStatus = "new"
	si2.Save()

	// get data si berdasarkan dummy customer pada so
	m, e := GetSalesInvoiceBySOID(so.ID)
	// seharusnya tidak error karena ada datanya
	assert.NoError(t, e)
	// seharusnya data yang didapat jumlahnya ada 2
	assert.Equal(t, 2, len(m), "seharusnya data jumlahnya ada 2")
}

func TestGetSalesInvoiceByCustomerIDFailNoData(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("delete from sales_invoice").Exec()

	// get data s1 berdasarkan dummy customer
	m, _ := GetSalesInvoiceBySOID(9999999)
	// seharusnya jumlah data 0 karena tidak ada datanya
	assert.Equal(t, 0, len(m))
}

func TestCheckInvoiceStatusSOSuccess(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()

	// buat data customer
	customer := model.DummyPartnership()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "active"
	so1.Save()

	// buat data so kedua dengan dummy customer
	so2 := model.DummySalesOrder()
	so2.Customer = customer
	so2.InvoiceStatus = "active"
	so2.Save()

	// get data so berdasarkan dummy customer
	m := CheckInvoiceStatusSO(customer.ID)
	// seharusnya tidak error invoice statusnya active
	assert.Equal(t, true, m)
}

func TestCheckInvoiceStatusSOFail(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()

	// buat data customer
	customer := model.DummyPartnership()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.InvoiceStatus = "new"
	so1.Save()

	// buat data so kedua dengan dummy customer
	so2 := model.DummySalesOrder()
	so2.Customer = customer
	so2.InvoiceStatus = "active"
	so2.Save()

	// get data so berdasarkan dummy customer
	m := CheckInvoiceStatusSO(customer.ID)
	// seharusnya error karena invoice statusnya ada yang != active
	assert.Equal(t, false, m)
}

func TestCreateInvoiceReceipt(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_invoice").Exec()

	// buat data customer
	customer := model.DummyPartnership()

	// buat data so pertama dengan dummy customer
	so1 := model.DummySalesOrder()
	so1.Customer = customer
	so1.DocumentStatus = "active"
	so1.InvoiceStatus = "new"
	so1.IsDeleted = int8(0)
	so1.Save()

	// buat data so kedua dengan dummy customer
	so2 := model.DummySalesOrder()
	so2.Customer = customer
	so2.DocumentStatus = "active"
	so2.InvoiceStatus = "new"
	so2.IsDeleted = int8(0)
	so2.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 10000
	si1.IsBundled = 0
	si1.Save()

	// buat data si kedua dengan dummy so kedua
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so2
	si2.DocumentStatus = "active"
	si2.IsBundled = 0
	si2.Save()

	// buat data si ketiga dengan dummy so kedua
	si3 := model.DummySalesInvoice()
	si3.SalesOrder = so2
	si3.DocumentStatus = "active"
	si3.IsBundled = 0
	si3.Save()

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.DocumentStatus = "active"

	fr := model.DummyFinanceRevenue()
	fr.Amount = 10000
	fr.RefID = uint64(si1.ID)
	fr.RefType = "sales_invoice"
	fr.IsDeleted = int8(0)
	fr.Save()

	e := CreateInvoiceReceipt(ir, []*model.SalesOrder{so1, so2})
	assert.NoError(t, e)

	// cek document status sales invoice
	// seharusnya berubah menjadi active dan is_bundled menjadi 1
	si1.Read("ID")
	assert.Equal(t, "active", si1.DocumentStatus)
	assert.Equal(t, int8(1), si1.IsBundled)

	si2.Read("ID")
	assert.Equal(t, "active", si2.DocumentStatus)
	assert.Equal(t, int8(1), si2.IsBundled)

	si3.Read("ID")
	assert.Equal(t, "active", si3.DocumentStatus)
	assert.Equal(t, int8(1), si3.IsBundled)

	// cek jumlah data invoice receipt item yang terbuat
	// seharusnya ada 3 data yg terbuat karena terdapat 3 sales invoice
	var total int64
	o.Raw("select count(*) from invoice_receipt_item where invoice_receipt_id = ?", ir.ID).QueryRow(&total)
	assert.Equal(t, int64(3), total)
}

func TestCreateInvoiceReceiptItemSalesInvoiceStatusNotActive(t *testing.T) {
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
	so1.InvoiceStatus = "active"
	so1.IsDeleted = int8(0)
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "new"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data finance revenue
	fr := model.DummyFinanceRevenue()
	fr.RefID = uint64(si1.ID)
	fr.RefType = "sales_invoice"
	fr.IsDeleted = int8(0)
	fr.Save()

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.DocumentStatus = "new"

	e := CreateInvoiceReceipt(ir, []*model.SalesOrder{so1})
	assert.NoError(t, e)

	// cek subtotal pada invoice receipt item
	var subtotal float64
	o.Raw("select subtotal from invoice_receipt_item where sales_invoice_id = ?", si1.ID).QueryRow(&subtotal)
	assert.Equal(t, float64(90000), subtotal)

	// cek document status sales invoice
	// seharusnya berubah menjadi active dan is_bundled menjadi 1
	si1.Read("ID")
	assert.Equal(t, "active", si1.DocumentStatus)
	assert.Equal(t, int8(1), si1.IsBundled)
}

func TestCreateInvoiceReceiptItemSalesInvoiceStatusActive(t *testing.T) {
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
	so1.InvoiceStatus = "active"
	so1.IsDeleted = int8(0)
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "active"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.Amount = 10000
	fr1.IsDeleted = int8(0)
	fr1.Save()

	// buat data finance revenue pertama
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.Amount = 10000
	fr2.IsDeleted = int8(0)
	fr2.Save()

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.DocumentStatus = "new"

	e := CreateInvoiceReceipt(ir, []*model.SalesOrder{so1})
	assert.NoError(t, e)

	// cek subtotal pada invoice receipt item
	var iri *model.InvoiceReceiptItem
	o.Raw("select * from invoice_receipt_item where sales_invoice_id = ?", si1.ID).QueryRow(&iri)
	assert.Equal(t, float64(70000), iri.Subtotal)

	// cek document status sales invoice
	// seharusnya berubah menjadi active dan is_bundled menjadi 1
	si1.Read("ID")
	assert.Equal(t, "active", si1.DocumentStatus)
	assert.Equal(t, int8(1), si1.IsBundled)
}

func TestCreateInvoiceReceiptSalesInvoiceStatusFinish(t *testing.T) {
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
	so1.InvoiceStatus = "active"
	so1.IsDeleted = int8(0)
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "finished"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.Amount = 10000
	fr1.IsDeleted = int8(0)
	fr1.Save()

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.DocumentStatus = "new"

	e := CreateInvoiceReceipt(ir, []*model.SalesOrder{so1})
	assert.NoError(t, e)

	// cek document status sales invoice
	// seharusnya tetap berstatus active dan is_bundled tetap 0
	si1.Read("ID")
	assert.Equal(t, "finished", si1.DocumentStatus)
	assert.Equal(t, int8(0), si1.IsBundled)
}

func TestCreateInvoiceReceiptItem(t *testing.T) {
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
	so1.InvoiceStatus = "active"
	so1.IsDeleted = int8(0)
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "active"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.Amount = 10000
	fr1.IsDeleted = int8(0)
	fr1.Save()

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.DocumentStatus = "new"
	ir.Save()

	e := CreateInvoiceReceiptItem(si1, ir)
	assert.NoError(t, e)

	// cek subtotal pada invoice receipt item
	var iri *model.InvoiceReceiptItem
	o.Raw("select * from invoice_receipt_item where sales_invoice_id = ?", si1.ID).QueryRow(&iri)
	assert.Equal(t, float64(80000), iri.Subtotal)
}

func TestCreateInvoiceReceiptPaymentSuccess(t *testing.T) {
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
	si1.TotalPaid = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data si kedua dengan dummy so pertama
	si2 := model.DummySalesInvoice()
	si2.SalesOrder = so1
	si2.DocumentStatus = "new"
	si2.TotalAmount = 90000
	si2.TotalPaid = 90000

	si2.IsBundled = 0
	si2.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.DocumentStatus = "cleared"
	fr1.Amount = 45000
	fr1.Save()

	// buat data finance revenue pertama
	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = uint64(si1.ID)
	fr2.RefType = "sales_invoice"
	fr2.DocumentStatus = "cleared"
	fr2.Amount = 45000
	fr2.Save()

	var financeRevenues []*model.FinanceRevenue
	financeRevenues = append(financeRevenues, fr1, fr2)

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.DocumentStatus = "new"
	ir.Save()

	// buat data invoice receipt item
	irm := model.DummyInvoiceReceiptItem()
	irm.SalesInvoice = si1
	irm.InvoiceReceipt = ir
	irm.Save()

	// buat data invoice receipt item kedua
	irm2 := model.DummyInvoiceReceiptItem()
	irm2.SalesInvoice = si2
	irm2.InvoiceReceipt = ir
	irm2.Save()

	var invoiceReceiptItems []*model.InvoiceReceiptItem
	invoiceReceiptItems = append(invoiceReceiptItems, irm, irm2)
	var invoiceReceiptRet []*model.InvoiceReceiptReturn

	e := CreateInvoiceReceiptPayment(ir.ID, invoiceReceiptItems, invoiceReceiptRet, financeRevenues)
	assert.NoError(t, e)

	var fr []*model.FinanceRevenue
	o.Raw("select * from finance_revenue where ref_type = 'sales_invoice'").QueryRows(&fr)
	assert.Equal(t, 2, len(fr))

	// cek document status sales invoice pertama
	si1.Read()
	assert.Equal(t, "finished", si1.DocumentStatus)

	// cek document status sales invoice kedua
	si2.Read()
	assert.Equal(t, "finished", si2.DocumentStatus)

	// cek total paid so
	so1.Read()
	assert.Equal(t, float64(200000), so1.TotalPaid)
	assert.Equal(t, "finished", so1.InvoiceStatus)

	// cek apabila invoice_status dan fulfillment_status pada so finished
	// maka document_statusnya menjadi finished
	so1.Read()
	assert.Equal(t, "finished", so1.DocumentStatus)

	// cek total debt partnership
	customer.Read()
	assert.Equal(t, float64(100000), customer.TotalDebt)
}

func TestCreateInvoiceReceiptItemInvoiceStatusFinish(t *testing.T) {
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
	so1.InvoiceStatus = "active"
	so1.IsDeleted = int8(0)
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	si1 := model.DummySalesInvoice()
	si1.SalesOrder = so1
	si1.DocumentStatus = "finished"
	si1.TotalAmount = 90000
	si1.IsBundled = 0
	si1.Save()

	// buat data finance revenue pertama
	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = uint64(si1.ID)
	fr1.RefType = "sales_invoice"
	fr1.Amount = 10000
	fr1.IsDeleted = int8(0)
	fr1.Save()

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.DocumentStatus = "new"
	ir.Save()

	e := CreateInvoiceReceiptItem(si1, ir)
	assert.NoError(t, e)

	// cek subtotal pada invoice receipt item
	var iri *model.InvoiceReceiptItem
	o.Raw("select * from invoice_receipt_item where sales_invoice_id = ?", si1.ID).QueryRow(&iri)
	assert.Equal(t, float64(90000), iri.Subtotal)
}

func TestCreateInvoiceReceiptReturnInvoiceStatusFinish(t *testing.T) {
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
	so1.InvoiceStatus = "active"
	so1.IsDeleted = int8(0)
	so1.Save()

	// buat data si pertama dengan dummy so pertama
	sr1 := model.DummySalesReturn()
	sr1.SalesOrder = so1
	sr1.DocumentStatus = "finished"
	sr1.TotalAmount = 90000
	sr1.IsBundled = 0
	sr1.Save()

	// buat data finance revenue pertama
	fe1 := model.DummyFinanceExpense()
	fe1.RefID = uint64(sr1.ID)
	fe1.RefType = "sales_return"
	fe1.Amount = 10000
	fe1.IsDeleted = int8(0)
	fe1.Save()

	// buat data invoice receipt
	ir := model.DummyInvoiceReceipt()
	ir.Partnership = customer
	ir.DocumentStatus = "new"
	ir.Save()

	e := CreateInvoiceReceiptReturn(sr1, ir)
	assert.NoError(t, e)

	// cek subtotal pada invoice receipt item
	var iri *model.InvoiceReceiptReturn
	o.Raw("select * from invoice_receipt_return where sales_return_id = ?", sr1.ID).QueryRow(&iri)
	assert.Equal(t, float64(90000), iri.Subtotal)
}
