// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sales

import (
	"fmt"
	"testing"
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/common/faker"
	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func dummyFinanceRevenue() *model.SalesInvoice {
	var financeRevenue model.FinanceRevenue
	si := dummySalesInvoice()
	faker.Fill(&financeRevenue, "ID")

	financeRevenue.CreatedBy = model.DummyUser()
	financeRevenue.Amount = 10000
	financeRevenue.RefID = uint64(si.ID)
	financeRevenue.DocumentStatus = "cleared"
	financeRevenue.RefType = "sales_invoice"
	financeRevenue.RecognitionDate = time.Now()
	if e := financeRevenue.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return si
}

// dummy finance revenue document status = "uncleared"
func dummyFinanceRevenue2() *model.SalesInvoice {
	var financeRevenue model.FinanceRevenue
	si := dummySalesInvoice()
	faker.Fill(&financeRevenue, "ID")

	financeRevenue.CreatedBy = model.DummyUser()
	financeRevenue.Amount = 10000
	financeRevenue.RefID = uint64(si.ID)
	financeRevenue.DocumentStatus = "uncleared"
	financeRevenue.RefType = "sales_invoice"
	financeRevenue.RecognitionDate = time.Now()
	if e := financeRevenue.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return si
}

// dummy finance revenue RefType = "purchase_return"
func dummyFinanceRevenue3() *model.SalesInvoice {
	var financeRevenue model.FinanceRevenue
	si := dummySalesInvoice()
	faker.Fill(&financeRevenue, "ID")

	financeRevenue.CreatedBy = model.DummyUser()
	financeRevenue.Amount = 10000
	financeRevenue.RefID = uint64(si.ID)
	financeRevenue.DocumentStatus = "cleared"
	financeRevenue.RefType = "purchase_return"
	financeRevenue.RecognitionDate = time.Now()
	if e := financeRevenue.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return si
}

func dummySalesInvoice() *model.SalesInvoice {
	var salesOrder model.SalesOrder
	faker.Fill(&salesOrder, "ID")
	//
	salesOrder.CreatedBy = model.DummyUser()
	salesOrder.TotalCharge = 40000

	if e := salesOrder.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}

	var salesInvoice model.SalesInvoice
	faker.Fill(&salesInvoice, "ID")

	salesInvoice.SalesOrder = &salesOrder
	salesInvoice.TotalPaid = 30000
	salesInvoice.CreatedBy = model.DummyUser()
	salesInvoice.IsDeleted = 0
	if e := salesInvoice.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}

	var salesInvoice2 model.SalesInvoice
	faker.Fill(&salesInvoice2, "ID")

	salesInvoice2.SalesOrder = &salesOrder
	salesInvoice2.TotalPaid = 10000
	salesInvoice2.CreatedBy = model.DummyUser()
	salesInvoice2.IsDeleted = 0
	if e := salesInvoice2.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &salesInvoice
}

func TestCalculateTotalPaidSI(t *testing.T) {
	fr := dummyFinanceRevenue()

	//test no error
	// update total paid sales invoice
	e := CalculateTotalPaidSI(fr)
	assert.NoError(t, e, "no errror")

	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_paid from sales_invoice where id = ?", fr.ID).QueryRow(&total)
	assert.Equal(t, float64(10000), total)

}

// TestCalculateTotalPaidSI2 Test Calculate Total Paid Sales Invoice
// revenue document status = "uncleared"
func TestCalculateTotalPaidSI2(t *testing.T) {
	fr := dummyFinanceRevenue2()

	//test no error
	// update total paid sales invoice
	e := CalculateTotalPaidSI(fr)
	assert.NoError(t, e, "no errror")

	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_paid from sales_invoice where id = ?", fr.ID).QueryRow(&total)
	assert.Equal(t, float64(0), total)

}

// TestCalculateTotalPaidSI3 Test Calculate Total Paid Sales Invoice
// where finance revenue RefType = "purchase_return"
func TestCalculateTotalPaidSI3(t *testing.T) {
	fr := dummyFinanceRevenue3()

	//test no error
	// update total paid sales invoice
	e := CalculateTotalPaidSI(fr)
	assert.NoError(t, e, "no errror")

	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_paid from sales_invoice where id = ?", fr.ID).QueryRow(&total)
	assert.Equal(t, float64(0), total)

}

func TestCalculateTotalPaidSO(t *testing.T) {
	si := dummySalesInvoice()

	//test no error
	// update total paid sales order
	e := CalculateTotalPaidSO(si.SalesOrder)
	assert.NoError(t, e, "no errror")

	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_paid from sales_order where id = ?", si.SalesOrder.ID).QueryRow(&total)
	assert.Equal(t, float64(40000), total)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////

func dummyWorkorderFulfillmentItem() *model.WorkorderFulfillmentItem {
	var so model.SalesOrder
	faker.Fill(&so, "ID")

	so.DocumentStatus = "active"
	so.IsDeleted = 0
	so.FulfillmentStatus = "active"
	so.CreatedBy = model.DummyUser()

	if e := so.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}

	var soi model.SalesOrderItem
	faker.Fill(&soi, "ID")

	soi.SalesOrder = &so
	soi.Quantity = 100
	soi.ItemVariant = model.DummyItemVariant()

	if e := soi.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}

	var f model.WorkorderFulfillment
	faker.Fill(&f, "ID")

	f.SalesOrder = &so
	f.IsDelivered = 1
	f.IsDeleted = 0
	f.CreatedBy = model.DummyUser()

	if e := f.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}

	var m model.WorkorderFulfillmentItem
	faker.Fill(&m, "ID")

	m.WorkorderFulfillment = &f
	m.Quantity = 100
	m.SalesOrderItem = &soi

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

func TestCheckFulfillmentStatus(t *testing.T) {
	wfi := dummyWorkorderFulfillmentItem()
	wfi.WorkorderFulfillment.DocumentStatus = "finished"
	wfi.WorkorderFulfillment.Save()
	e := CheckFulfillmentStatus(wfi.SalesOrderItem.SalesOrder)
	wfi.SalesOrderItem.SalesOrder.Read()
	assert.NoError(t, e, "no errror")
	assert.Equal(t, "finished", wfi.SalesOrderItem.SalesOrder.FulfillmentStatus, "Seharusnya fulfillment status finished")
	fmt.Println("qty wfi", wfi.Quantity)
	fmt.Println("qty item", wfi.SalesOrderItem.Quantity)
	fmt.Println("fs", wfi.SalesOrderItem.SalesOrder.FulfillmentStatus)
}

func TestCheckFulfillmentStatusNotFinished(t *testing.T) {
	wfi2 := dummyWorkorderFulfillmentItem()
	wfi2.Quantity = 80
	wfi2.Save()
	e := CheckFulfillmentStatus(wfi2.SalesOrderItem.SalesOrder)
	assert.NoError(t, e, "no errror")
	wfi2.SalesOrderItem.SalesOrder.Read()
	assert.NotEqual(t, "finished", wfi2.SalesOrderItem.SalesOrder.FulfillmentStatus, "Seharusnya fulfillment status != finished")
	fmt.Println("fs2", wfi2.SalesOrderItem.SalesOrder.FulfillmentStatus)
}

func TestCheckFulfillmentStatusMultiItemNotFinished(t *testing.T) {
	wfi2 := dummyWorkorderFulfillmentItem()
	wfi2.Quantity = 10
	wfi2.Save()

	//new fulfillment item
	wfi3 := dummyWorkorderFulfillmentItem()

	//one sales order item and sales order with multi fulfillment item
	wfi3.Quantity = 80
	wfi3.WorkorderFulfillment = wfi2.WorkorderFulfillment
	wfi3.SalesOrderItem = wfi2.SalesOrderItem
	wfi3.Save()
	wfi3.WorkorderFulfillment.SalesOrder = wfi2.WorkorderFulfillment.SalesOrder
	wfi3.WorkorderFulfillment.Save()

	e := CheckFulfillmentStatus(wfi2.SalesOrderItem.SalesOrder)
	assert.NoError(t, e, "no errror")

	wfi2.SalesOrderItem.SalesOrder.Read()
	assert.NotEqual(t, "finished", wfi2.SalesOrderItem.SalesOrder.FulfillmentStatus, "Seharusnya fulfillment status != finished")
	fmt.Println("fs2", wfi2.SalesOrderItem.SalesOrder.FulfillmentStatus)
}

func TestCheckFulfillmentStatusMultiItemFinished(t *testing.T) {
	wfi2 := dummyWorkorderFulfillmentItem()
	wfi2.WorkorderFulfillment.DocumentStatus = "finished"
	wfi2.WorkorderFulfillment.Save()
	wfi2.Quantity = 20
	wfi2.Save()

	//new fulfillment item
	wfi3 := dummyWorkorderFulfillmentItem()

	//one sales order item and sales order with multi fulfillment item
	wfi3.Quantity = 80
	wfi3.WorkorderFulfillment = wfi2.WorkorderFulfillment
	wfi3.SalesOrderItem = wfi2.SalesOrderItem
	wfi3.Save()
	wfi3.WorkorderFulfillment.SalesOrder = wfi2.WorkorderFulfillment.SalesOrder
	wfi3.WorkorderFulfillment.Save()

	e := CheckFulfillmentStatus(wfi2.SalesOrderItem.SalesOrder)
	assert.NoError(t, e, "no errror")

	wfi2.SalesOrderItem.SalesOrder.Read()
	assert.Equal(t, "finished", wfi2.SalesOrderItem.SalesOrder.FulfillmentStatus, "Seharusnya fulfillment status = finished")
	fmt.Println("fs2", wfi2.SalesOrderItem.SalesOrder.FulfillmentStatus)
}

func TestCheckFulfillmentStatusMultiItemNotDelivered(t *testing.T) {
	wfi2 := dummyWorkorderFulfillmentItem()
	wfi2.Quantity = 20
	wfi2.Save()

	//new fulfillment item
	wfi3 := dummyWorkorderFulfillmentItem()

	//one sales order item and sales order with multi fulfillment item
	wfi3.Quantity = 70
	wfi3.WorkorderFulfillment = wfi2.WorkorderFulfillment
	wfi3.WorkorderFulfillment.IsDelivered = 0
	wfi3.SalesOrderItem = wfi2.SalesOrderItem
	wfi3.Save()
	wfi3.WorkorderFulfillment.SalesOrder = wfi2.WorkorderFulfillment.SalesOrder
	wfi3.WorkorderFulfillment.Save()

	e := CheckFulfillmentStatus(wfi2.SalesOrderItem.SalesOrder)
	assert.NoError(t, e, "no errror")

	wfi2.SalesOrderItem.SalesOrder.Read()
	assert.NotEqual(t, "finished", wfi2.SalesOrderItem.SalesOrder.FulfillmentStatus, "Seharusnya fulfillment status != finished")
	fmt.Println("fs2", wfi2.SalesOrderItem.SalesOrder.FulfillmentStatus)
}

func TestCheckFulfillmentStatusMultiItemSoCancelled(t *testing.T) {
	wfi2 := dummyWorkorderFulfillmentItem()
	wfi2.Quantity = 20
	wfi2.Save()

	//new fulfillment item
	wfi3 := dummyWorkorderFulfillmentItem()

	//one sales order item and sales order with multi fulfillment item
	wfi3.Quantity = 80
	wfi3.WorkorderFulfillment = wfi2.WorkorderFulfillment
	wfi3.SalesOrderItem = wfi2.SalesOrderItem
	wfi3.Save()
	wfi3.WorkorderFulfillment.SalesOrder = wfi2.WorkorderFulfillment.SalesOrder
	wfi3.WorkorderFulfillment.Save()
	wfi3.WorkorderFulfillment.SalesOrder.DocumentStatus = "approved_cancel"
	wfi3.WorkorderFulfillment.SalesOrder.Save()

	e := CheckFulfillmentStatus(wfi2.SalesOrderItem.SalesOrder)
	assert.NoError(t, e, "no errror")

	wfi2.SalesOrderItem.SalesOrder.Read()
	assert.NotEqual(t, "finished", wfi2.SalesOrderItem.SalesOrder.FulfillmentStatus, "Seharusnya fulfillment status != finished")
	fmt.Println("fs2", wfi2.SalesOrderItem.SalesOrder.FulfillmentStatus)
}

func TestCheckFulfillmentDeleted(t *testing.T) {
	wfi2 := dummyWorkorderFulfillmentItem()
	wfi2.Quantity = 20
	wfi2.Save()

	//new fulfillment item
	wfi3 := dummyWorkorderFulfillmentItem()

	//one sales order item and sales order with multi fulfillment item
	wfi3.Quantity = 80
	wfi3.WorkorderFulfillment = wfi2.WorkorderFulfillment
	wfi3.SalesOrderItem = wfi2.SalesOrderItem
	wfi3.Save()
	wfi3.WorkorderFulfillment.SalesOrder = wfi2.WorkorderFulfillment.SalesOrder
	wfi3.WorkorderFulfillment.IsDeleted = 1
	wfi3.WorkorderFulfillment.Save()

	e := CheckFulfillmentStatus(wfi2.SalesOrderItem.SalesOrder)
	assert.NoError(t, e, "no errror")

	wfi2.SalesOrderItem.SalesOrder.Read()
	assert.NotEqual(t, "finished", wfi2.SalesOrderItem.SalesOrder.FulfillmentStatus, "Seharusnya fulfillment status != finished")
	fmt.Println("fs2", wfi2.SalesOrderItem.SalesOrder.FulfillmentStatus)
}

func TestCheckFulfillmentStatusMultiItemSoDeleted(t *testing.T) {
	wfi2 := dummyWorkorderFulfillmentItem()
	wfi2.Quantity = 20
	wfi2.Save()

	//new fulfillment item
	wfi3 := dummyWorkorderFulfillmentItem()

	//one sales order item and sales order with multi fulfillment item
	wfi3.Quantity = 80
	wfi3.WorkorderFulfillment = wfi2.WorkorderFulfillment
	wfi3.SalesOrderItem = wfi2.SalesOrderItem
	wfi3.Save()
	wfi3.WorkorderFulfillment.SalesOrder = wfi2.WorkorderFulfillment.SalesOrder
	wfi3.WorkorderFulfillment.Save()
	wfi3.WorkorderFulfillment.SalesOrder.IsDeleted = 1
	wfi3.WorkorderFulfillment.SalesOrder.Save()

	e := CheckFulfillmentStatus(wfi2.SalesOrderItem.SalesOrder)
	assert.NoError(t, e, "no errror")

	wfi2.SalesOrderItem.SalesOrder.Read()
	assert.NotEqual(t, "finished", wfi2.SalesOrderItem.SalesOrder.FulfillmentStatus, "Seharusnya fulfillment status != finished")
	fmt.Println("fs2", wfi2.SalesOrderItem.SalesOrder.FulfillmentStatus)
}

func dummySalesInvoice1() *model.SalesInvoice {
	var salesOrder model.SalesOrder
	faker.Fill(&salesOrder, "ID")
	//
	salesOrder.CreatedBy = model.DummyUser()
	salesOrder.InvoiceStatus = "active"
	salesOrder.TotalCharge = 40000
	salesOrder.TotalPaid = 40000

	if e := salesOrder.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}

	var salesInvoice model.SalesInvoice
	faker.Fill(&salesInvoice, "ID")

	salesInvoice.SalesOrder = &salesOrder
	salesInvoice.TotalPaid = 40000
	salesInvoice.CreatedBy = model.DummyUser()
	salesInvoice.IsDeleted = 0
	if e := salesInvoice.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}

	return &salesInvoice
}

func TestCheckInvoiceStatus(t *testing.T) {
	si := dummySalesInvoice1()
	e := CheckInvoiceStatus(si.SalesOrder)
	si.SalesOrder.Read()
	assert.NoError(t, e, "no errror")
	assert.Equal(t, "finished", si.SalesOrder.InvoiceStatus, "Seharusnya invoice status finished")
	fmt.Println("tc1", si.SalesOrder.TotalCharge)
	fmt.Println("tp1", si.SalesOrder.TotalPaid)

	si2 := dummySalesInvoice1()
	si2.SalesOrder.TotalPaid = 30000
	si2.SalesOrder.Save()
	si2.TotalPaid = 30000
	si2.Save()
	e = CheckInvoiceStatus(si2.SalesOrder)
	si2.SalesOrder.Read()
	assert.NoError(t, e, "no errror")
	assert.NotEqual(t, "finished", si2.SalesOrder.InvoiceStatus, "Seharusnya invoice status tidak finished")
	fmt.Println("invoice status 2", si2.SalesOrder.InvoiceStatus)
	fmt.Println("tc2", si2.SalesOrder.TotalCharge)
	fmt.Println("tp2", si2.SalesOrder.TotalPaid)
}

func dummySO() *model.SalesOrder {
	var salesOrder model.SalesOrder
	faker.Fill(&salesOrder, "ID")
	//
	salesOrder.CreatedBy = model.DummyUser()
	salesOrder.TotalCharge = 40000
	salesOrder.TotalPaid = 40000
	salesOrder.FulfillmentStatus = "finished"
	salesOrder.InvoiceStatus = "finished"

	if e := salesOrder.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}

	return &salesOrder
}

func TestCheckDocumentStatus(t *testing.T) {
	so := dummySO()

	ex := CheckDocumentStatus(so)
	assert.NoError(t, ex, "no errror")
}

/////////////////////////////////////////////////////////////////////////////////////////////////

func dummySalesReturnItem() (*model.SalesOrder, *model.SalesReturn) {
	so := model.DummySalesOrder()
	so.DocumentStatus = "new"
	so.IsDeleted = 0
	so.Save()

	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.Quantity = 100
	soi.UnitPrice = 10000
	soi.Subtotal = float64(soi.Quantity) * soi.UnitPrice
	soi.Save()

	soi2 := model.DummySalesOrderItem()
	soi2.SalesOrder = so
	soi2.Quantity = 100
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
	fi.Quantity = 100
	fi.SalesOrderItem = soi
	fi.Save()

	fi2 := model.DummyWorkorderFulfillmentItem()
	fi2.WorkorderFulfillment = f
	fi2.Quantity = 100
	fi2.SalesOrderItem = soi2
	fi2.Save()

	sr := model.DummySalesReturn()
	sr.IsDeleted = 0
	sr.DocumentStatus = "new"
	sr.SalesOrder = so
	sr.Save()

	sri := model.DummySalesReturnItem()
	sri.SalesOrderItem = soi
	sri.Quantity = 80
	sri.SalesReturn = sr
	sri.Save()

	sri2 := model.DummySalesReturnItem()
	sri2.SalesOrderItem = soi2
	sri2.Quantity = 90
	sri2.SalesReturn = sr
	sri2.Save()

	return so, sr
}

func TestSalesCanBeReturn(t *testing.T) {
	so, _ := dummySalesReturnItem()

	_, e := CanBeReturnSales(so, nil)
	assert.NoError(t, e, "no error")

	soxItem, _ := CanBeReturnSales(so, nil)
	var total, st float32

	for _, a := range soxItem {
		total = a.CanBeReturn
		st += total
	}
	assert.Equal(t, float32(30), st)
}

func TestSalesCanBeReturnForUpdate(t *testing.T) {
	so, sr := dummySalesReturnItem()

	_, e := CanBeReturnSales(so, sr)
	assert.NoError(t, e, "no error")

	soxItem, _ := CanBeReturnSales(so, sr)
	var total, st float32

	for _, a := range soxItem {
		total = a.CanBeReturn
		st += total
	}
	assert.Equal(t, float32(200), st)
}
func TestGetAllSalesOrderNoData(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM sales_order").Exec()
	rq := orm.RequestQuery{}
	// test tidak ada data sales_order
	m, total, e := GetAllSalesOrder(&rq)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, m)
	assert.NoError(t, e)
}

func TestGetAllSalesOrder(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM sales_order").Exec()
	// buat dummy
	sls := model.DummySalesOrder()
	sls.IsDeleted = int8(0)
	sls.Save()
	sls2 := model.DummySalesOrder()
	sls2.IsDeleted = int8(1)
	sls2.Save()
	rq := orm.RequestQuery{}
	m, total, e := GetAllSalesOrder(&rq)
	assert.NoError(t, e)
	assert.Equal(t, int64(1), total)
	assert.NotEmpty(t, m)
	for _, u := range *m {
		assert.Equal(t, sls.ID, u.ID)
	}
}

func TestGetDetailSalesOrderNoData(t *testing.T) {
	m, e := GetDetailSalesOrder(99999, nil)
	assert.Error(t, e)
	assert.Empty(t, m)
}

func TestGetDetailSalesOrderDeletedData(t *testing.T) {
	sls := model.DummySalesOrder()
	sls.IsDeleted = int8(1)
	sls.Save()

	m, e := GetDetailSalesOrder(sls.ID, nil)
	assert.Error(t, e)
	assert.Empty(t, m)
}

func TestGetDetailSalesOrderSuccess(t *testing.T) {
	sls := model.DummySalesOrder()
	sls.TotalCharge = 100000
	sls.IsDeleted = int8(0)
	sls.Save()

	soi := model.DummySalesOrderItem()
	soi.SalesOrder = &model.SalesOrder{ID: sls.ID}
	soi.Save()

	si := model.DummySalesInvoice()
	si.SalesOrder = sls
	si.IsDeleted = 0
	si.Save()

	wf := model.DummyWorkorderFulfillment()
	wf.SalesOrder = sls
	wf.IsDeleted = 0
	wf.Save()

	wfi := model.DummyWorkorderFulfillmentItem()
	wfi.WorkorderFulfillment = wf
	wfi.Quantity = 5
	wfi.Save()

	sr := model.DummySalesReturn()
	sr.SalesOrder = sls
	sr.TotalAmount = 20000
	sr.IsDeleted = 0
	sr.DocumentStatus = "active"
	sr.Save()

	sr2 := model.DummySalesReturn()
	sr2.SalesOrder = sls
	sr2.TotalAmount = 30000
	sr2.IsDeleted = 0
	sr2.DocumentStatus = "active"
	sr2.Save()

	f := model.DummyFinanceExpense()
	f.RefID = uint64(sr2.ID)
	f.RefType = "sales_return"
	f.Amount = 5000
	f.IsDeleted = 0
	f.Save()

	m, e := GetDetailSalesOrder(sls.ID, []string{"sales_order_items", "sales_invoices", "workorder_fulfillments", "sales_returns"})
	assert.NoError(t, e)
	assert.NotEmpty(t, m)
	assert.Equal(t, sls.TotalCharge, m.TotalCharge)
	assert.Equal(t, sls.CreatedBy.FullName, m.CreatedBy.FullName)

	for _, u := range m.SalesOrderItems {
		assert.Equal(t, soi.ID, u.ID)
		assert.Equal(t, soi.Discount, u.Discount)
		assert.Equal(t, soi.Note, u.Note)
	}

	for _, u := range m.SalesInvoices {
		assert.Equal(t, si.ID, u.ID)
		assert.Equal(t, si.Note, u.Note)
	}

	for _, u := range m.WorkorderFulfillments {
		assert.Equal(t, wf.ID, u.ID)
		assert.Equal(t, wf.DocumentStatus, u.DocumentStatus)
		assert.Equal(t, wf.Note, u.Note)
	}

	assert.Equal(t, m.TotalRefund, float64(50000))
	assert.Equal(t, m.TotalPaidRefund, float64(5000))
}

func TestGetSalesInvoiceByPartnerID(t *testing.T) {
	dummyPartner := model.DummyPartnership()
	sinvoice := model.DummySalesInvoice()
	sorder := model.DummySalesOrder()

	sorder.Customer = dummyPartner
	sorder.Save("Customer")

	sinvoice.SalesOrder = sorder
	sinvoice.DocumentStatus = "active"
	sinvoice.Save("SalesOrder", "DocumentStatus")

	data, err := getSalesInvoiceByPartnerID(dummyPartner)
	assert.NotNil(t, data, fmt.Sprint("Data tidak boleh kosong, Actual:", data))
	assert.NoError(t, err, fmt.Sprint("Tidak boleh ada pesan error, Actual:", err))
}

func TestCreateSalesOrderNoErrorAutoIsONE(t *testing.T) {
	DummyItemVariant := model.DummyItemVariant()
	DummyItemVariant.BasePrice = 1000
	DummyItemVariant.Save("BasePrice")
	DummyPartner := model.DummyPartnership()

	code, _ := util.CodeGen("code_sales_order", "sales_order")
	session := new(auth.SessionData)
	session.User = model.DummyUser()
	sorder := createRequest{
		Code:            code,
		RecognitionDate: time.Now(),
		EtaDate:         time.Now(),
		CustomerID:      common.Encrypt(DummyPartner.ID),
		ShipmentAddress: "Every corner in the world",
		AutoInvoice:     1,
		AutoFullfilment: 1,
		SalesOrderItem: []salesOrderItem{
			{
				ItemVariantID: common.Encrypt(DummyItemVariant.ID),
				Quantity:      float32(10),
				Discount:      float32(0),
				UnitPrice:     DummyItemVariant.BasePrice,
				PricingType:   common.Encrypt(model.DummyPricingType().ID),
				Subtotal:      float64(10) * DummyItemVariant.BasePrice,
				Note:          "None",
			},
		},
		Discount:     float32(10),
		Tax:          float32(10),
		ShipmentCost: float64(0),
		TotalCharge:  float64(10) * DummyItemVariant.BasePrice,
		Note:         "Walk in customer",
		Session:      session,
	}

	sales, err := CreateSalesOrder(&sorder)
	// Cek error
	assert.NoError(t, err, "Tidak boleh ada error")
	// Cek sales tidak boleh kosong
	assert.NotNil(t, sales, fmt.Sprint("Tidak boleh nil atau kosong, actual: ", sales))
	// Cek sales order item harus terbuat
	for _, row := range sales.SalesOrderItems {
		assert.NotNil(t, row, "Tidak boleh kosong")
	}
	// saat auto invoice di checklist sales invoice otomatis dibuat
	sinvoices := []model.SalesInvoice{}
	orm.NewOrm().Raw("SELECT * FROM sales_invoice WHERE sales_order_id = ?", sales.ID).QueryRows(&sinvoices)
	for _, row := range sinvoices {
		assert.NotNil(t, row, "Sales invoice tidak boleh kosong saat auto invoice")
	}
	// saat auto fulfillment di checklist workorder fulfillment otomatis dibuat
	wofulfillment := model.WorkorderFulfillment{}
	orm.NewOrm().Raw("SELECT * FROM workorder_fulfillment WHERE sales_order_id = ?", sales.ID).QueryRow(&wofulfillment)
	assert.NotNil(t, wofulfillment, "Workorder fulfillment tidak boleh kosong saat auto fulfillment")
	// saat auto fulfillment di checklist workorder fulfillment item otomatis dibuat
	for _, row := range wofulfillment.WorkorderFulFillmentItems {
		assert.NotNil(t, row, "Workorder fulfillment item tidak boleh kosong saat auto fulfillment")
	}
}

func TestCreateSalesOrderNoErrorAutoIsZERO(t *testing.T) {
	orm.NewOrm().Raw("DELETE FROM partnership").Exec()
	orm.NewOrm().Raw("DELETE FROM sales_order").Exec()
	orm.NewOrm().Raw("DELETE FROM sales_invoice").Exec()

	DummyItemVariant := model.DummyItemVariant()
	DummyItemVariant.BasePrice = 1000
	DummyItemVariant.Save("BasePrice")
	DummyPartner := model.DummyPartnership()
	DummyPartner.IsDefault = int8(0)
	DummyPartner.Save("IsDefault")

	code, _ := util.CodeGen("code_sales_order", "sales_order")
	session := new(auth.SessionData)
	session.User = model.DummyUser()
	sorder := createRequest{
		Code:            code,
		RecognitionDate: time.Now(),
		EtaDate:         time.Now(),
		CustomerID:      common.Encrypt(DummyPartner.ID),
		ShipmentAddress: "Every corner in the world",
		AutoInvoice:     int8(0),
		AutoFullfilment: int8(0),
		SalesOrderItem: []salesOrderItem{
			{
				ItemVariantID: common.Encrypt(DummyItemVariant.ID),
				Quantity:      float32(10),
				Discount:      float32(0),
				UnitPrice:     DummyItemVariant.BasePrice,
				PricingType:   common.Encrypt(model.DummyPricingType().ID),
				Subtotal:      float64(10) * DummyItemVariant.BasePrice,
				Note:          "None",
			},
		},
		Discount:     float32(10),
		Tax:          float32(10),
		ShipmentCost: float64(0),
		TotalCharge:  float64(10) * DummyItemVariant.BasePrice,
		Note:         "none",
		Session:      session,
	}

	sales, err := CreateSalesOrder(&sorder)
	// Cek error
	assert.NoError(t, err, "Tidak boleh ada error")
	// Cek sales tidak boleh kosong
	assert.NotNil(t, sales, fmt.Sprint("Tidak boleh nil atau kosong, actual: ", sales))
	// Cek sales order item harus terbuat
	for _, row := range sales.SalesOrderItems {
		assert.NotNil(t, row, "Tidak boleh kosong")
	}
	// saat auto invoice di checklist sales invoice otomatis dibuat
	sinvoices := []model.SalesInvoice{}
	orm.NewOrm().Raw("SELECT * FROM sales_invoice WHERE sales_order_id = ?", sales.ID).QueryRows(&sinvoices)
	for _, row := range sinvoices {
		assert.Nil(t, row, "Sales invoice tidak boleh kosong saat auto invoice")
	}
	// saat auto fulfillment di checklist workorder fulfillment otomatis dibuat
	wofulfillment := model.WorkorderFulfillment{}
	orm.NewOrm().Raw("SELECT * FROM workorder_fulfillment WHERE sales_order_id = ?", sales.ID).QueryRow(&wofulfillment)
	assert.Nil(t, wofulfillment.SalesOrder, "Workorder fulfillment kosong saat auto fulfillment 0")
	// saat auto fulfillment di checklist workorder fulfillment item otomatis dibuat
	for _, row := range wofulfillment.WorkorderFulFillmentItems {
		assert.Nil(t, row, "Workorder fulfillment item kosong saat auto fulfillment 0")
	}
}

func TestUpdateSalesOrder(t *testing.T) {
	so := model.DummySalesOrder()

	var items []*model.SalesOrderItem
	item := model.DummySalesOrderItem()
	item.Quantity = 5
	item.SalesOrder = so
	item.Save()

	item2 := model.DummySalesOrderItem()
	item2.Quantity = 8
	item2.SalesOrder = so
	item2.Save()
	items = append(items, item, item2)

	so.SalesOrderItems = items
	so.IsDeleted = 0
	so.TotalCharge = 10000
	so.Save()

	// item yang baru
	var itemsNew []*model.SalesOrderItem
	itemNew := &model.SalesOrderItem{ID: item.ID, SalesOrder: so, ItemVariant: model.DummyItemVariant(), Quantity: 1, Note: "abc"}
	itemNew2 := &model.SalesOrderItem{SalesOrder: so, ItemVariant: model.DummyItemVariant(), Quantity: 3, Note: "xxxx"}
	itemsNew = append(itemsNew, itemNew, itemNew2)

	customer := model.DummyPartnership()
	customer.TotalDebt = 10000
	customer.Save("TotalDebt")

	soReq := &model.SalesOrder{ID: so.ID, TotalCharge: 8000, SalesOrderItems: items, Customer: customer, CreatedBy: model.DummyUser(), DocumentStatus: "new", InvoiceStatus: "new", FulfillmentStatus: "new", ShipmentStatus: "new"}

	// update data yang di so.SalesOrderItem
	res, e := UpdateSalesOrder(soReq, itemsNew)
	assert.NoError(t, e)
	assert.NotEmpty(t, res)
	assert.Equal(t, 2, len(res.SalesOrderItems))
	var emptyLoad []string
	emptyLoad = append(emptyLoad, "sales_order_items")
	data, _ := GetDetailSalesOrder(so.ID, emptyLoad)
	assert.Equal(t, 2, len(data.SalesOrderItems))

	customer = &model.Partnership{ID: soReq.Customer.ID}
	customer.Read()
	assert.Equal(t, float64(8000), customer.TotalDebt)
}

func TestGetWorkorderFulfillmentID(t *testing.T) {

	so := model.DummySalesOrder()
	so.IsDeleted = 0
	so.Save()

	wf := model.DummyWorkorderFulfillment()
	wf.IsDeleted = 0
	wf.CreatedBy = model.DummyUser()
	wf.SalesOrder = so
	wf.Save()

	wfi := model.DummyWorkorderFulfillmentItem()
	wfi.WorkorderFulfillment = wf
	wfi.Save()

	m, e := getWorkorderFulfillments("sales_order_id", so.ID)
	assert.NoError(t, e)
	assert.NotEmpty(t, m)
	for _, u := range m {
		for _, x := range u.WorkorderFulFillmentItems {
			assert.Equal(t, wfi.ID, x.ID)
			assert.Equal(t, wfi.Quantity, x.Quantity)
		}
	}
}

func TestCalculateStock(t *testing.T) {
	so := model.DummySalesOrder()
	so.DocumentStatus = "new"
	so.Note = "Cancel SO"
	so.InvoiceStatus = "active"
	so.IsDeleted = 0
	so.TotalCharge = 10000
	so.Save()

	iv := model.DummyItemVariant()
	iv.CommitedStock = 100
	iv.AvailableStock = 100
	iv.Save()

	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.ItemVariant = iv
	soi.Quantity = 100
	soi.QuantityFulfillment = 50
	soi.UnitPrice = 10000
	soi.Save()

	var sois []*model.SalesOrderItem
	sois = append(sois, soi)
	so.SalesOrderItems = sois

	ivcs := iv.CommitedStock
	ivcs = ivcs - (soi.Quantity - soi.QuantityFulfillment)

	e := calculateStockCommited(so)
	assert.NoError(t, e)

	er := calculateStockAvailable(so)
	assert.NoError(t, er)

	data := &model.ItemVariant{ID: iv.ID}
	data.Read("ID")
	assert.Equal(t, iv.ID, data.ID)
	assert.Equal(t, ivcs, data.CommitedStock)
}

func TestCheckFulfillmentInWorkorderShipmentItem(t *testing.T) {
	p := model.DummyPartnership()
	p.TotalDebt = 10000
	p.TotalCredit = 5000
	p.Save()

	so := model.DummySalesOrder()
	so.DocumentStatus = "new"
	so.Note = "Cancel SO"
	so.InvoiceStatus = "active"
	so.IsDeleted = 0
	so.TotalCharge = 10000
	so.Customer = p
	so.Save()

	iv := model.DummyItemVariant()
	iv.CommitedStock = 100
	iv.Save()

	ivs := model.DummyItemVariantStock()
	ivs.ItemVariant = iv
	ivs.Save()

	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.ItemVariant = iv
	soi.Quantity = 100
	soi.QuantityFulfillment = 50
	soi.UnitPrice = 10000
	soi.Save()

	soi2 := model.DummySalesOrderItem()
	soi2.SalesOrder = so
	soi2.Quantity = 100
	soi2.UnitPrice = 10000
	soi2.Save()

	si := model.DummySalesInvoice()
	si.SalesOrder = so
	si.DocumentStatus = "active"
	si.IsDeleted = 0
	si.IsBundled = 1
	si.Save()

	fr := model.DummyFinanceRevenue()
	fr.IsDeleted = 0
	fr.RefType = "sales_invoice"
	fr.Note = "CANCEL SO coy"
	fr.RefID = uint64(si.ID)
	fr.Save()

	wf := model.DummyWorkorderFulfillment()
	wf.SalesOrder = so
	wf.IsDelivered = 1
	wf.IsDeleted = 0
	wf.DocumentStatus = "finished"
	wf.CreatedBy = model.DummyUser()
	wf.Save()

	wfi := model.DummyWorkorderFulfillmentItem()
	wfi.WorkorderFulfillment = wf
	wfi.Quantity = 100
	wfi.SalesOrderItem = soi
	wfi.Save()

	ivsl := model.DummyItemVariantStockLog()
	ivsl.ItemVariantStock = ivs
	ivsl.RefType = "workorder_fulfillment"
	ivsl.RefID = uint64(wf.ID)
	ivsl.Save()

	ws := model.DummyWorkorderShipment()
	ws.IsDeleted = 0
	ws.Save()

	wsi := model.DummyWorkorderShipmentItem()
	wsi.WorkorderShipment = ws
	wsi.WorkorderFulfillment = wf
	wsi.Save()

	ir := model.DummyInvoiceReceipt()
	ir.DocumentStatus = "finished"
	ir.Save()

	iri := model.DummyInvoiceReceiptItem()
	iri.SalesInvoice = si
	iri.InvoiceReceipt = ir
	iri.Save()

	var wfx []*model.WorkorderFulfillmentItem
	wfx = append(wfx, wfi)
	wf.WorkorderFulFillmentItems = wfx

	m, e := getWorkorderFulfillments("sales_order_id", so.ID)
	assert.NoError(t, e)
	assert.NotEmpty(t, m)

	for _, u := range m {
		e = checkFulfillmentInWorkorderShipmentItem("workorder_fulfillment_id", u.ID)
		assert.NoError(t, e)

		data := &model.WorkorderShipment{ID: ws.ID}
		data.Read("ID")

		assert.NotEqual(t, ws.IsDeleted, data.IsDeleted)
	}

}

func TestCheckInvoiceReceiptItem(t *testing.T) {
	p := model.DummyPartnership()
	p.TotalDebt = 10000
	p.TotalCredit = 5000
	p.Save()

	so := model.DummySalesOrder()
	so.DocumentStatus = "new"
	so.Note = "Cancel SO"
	so.InvoiceStatus = "active"
	so.IsDeleted = 0
	so.TotalCharge = 10000
	so.Customer = p
	so.Save()

	iv := model.DummyItemVariant()
	iv.CommitedStock = 100
	iv.Save()

	ivs := model.DummyItemVariantStock()
	ivs.ItemVariant = iv
	ivs.Save()

	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.ItemVariant = iv
	soi.Quantity = 100
	soi.QuantityFulfillment = 50
	soi.UnitPrice = 10000
	soi.Save()

	soi2 := model.DummySalesOrderItem()
	soi2.SalesOrder = so
	soi2.Quantity = 100
	soi2.UnitPrice = 10000
	soi2.Save()

	si := model.DummySalesInvoice()
	si.SalesOrder = so
	si.DocumentStatus = "active"
	si.IsDeleted = 0
	si.IsBundled = 1
	si.Save()

	fr := model.DummyFinanceRevenue()
	fr.IsDeleted = 0
	fr.RefType = "sales_invoice"
	fr.Note = "CANCEL SO coy"
	fr.RefID = uint64(si.ID)
	fr.Save()

	ir := model.DummyInvoiceReceipt()
	ir.DocumentStatus = "finished"
	ir.Save()

	iri := model.DummyInvoiceReceiptItem()
	iri.SalesInvoice = si
	iri.InvoiceReceipt = ir
	iri.Save()

	e := checkInvoiceReceiptItem(si)
	assert.NoError(t, e)
}

func TestShowSalesOrderFulfillment(t *testing.T) {
	fulfill := model.DummyWorkorderFulfillmentItem()

	// workorder fulfillment item
	fulfill.Quantity = float32(10)
	fulfill.Save()

	// workorder fulfillment
	fulfill.WorkorderFulfillment.IsDeleted = 0
	fulfill.WorkorderFulfillment.Save("IsDeleted")

	// sales order item
	fulfill.SalesOrderItem.Quantity = float32(20)
	fulfill.SalesOrderItem.Save("Quantity")

	sales, err := ShowSalesOrderFulfillment("id", fulfill.SalesOrderItem.SalesOrder.ID)
	assert.NotNil(t, sales, "Tidak boleh kosong jika berhasil mengambil data")
	assert.NotNil(t, sales.SalesOrderItems[0], "Tidak boleh kosong jika berhasil mengambil data")
	assert.NotNil(t, sales.SalesOrderItems[0].ItemVariant, "Tidak boleh kosong jika berhasil mengambil data")
	assert.NotNil(t, sales.SalesOrderItems[0].ItemVariant.Item, "Tidak boleh kosong jika berhasil mengambil data")
	assert.NotNil(t, sales.SalesOrderItems[0].ItemVariant.Item.Category, "Tidak boleh kosong jika berhasil mengambil data")
	assert.NoError(t, err, "Tidak boleh ada error")
}

func TestShowSalesOrderFulfillmentErrorWoFulfillItemIsDelete0(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM sales_order_item")
	o.Raw("DELETE FROM sales_order")
	o.Raw("DELETE FROM workoorder_fulfillment_item")
	o.Raw("DELETE FROM workorder_fulfillment")

	fulfill := model.DummyWorkorderFulfillmentItem()

	// workorder fulfillment item
	fulfill.Quantity = float32(10)
	fulfill.Save()

	// workorder fulfillment
	fulfill.WorkorderFulfillment.IsDeleted = 1
	fulfill.WorkorderFulfillment.Save("IsDeleted")

	// sales order item
	fulfill.SalesOrderItem.Quantity = float32(20)
	fulfill.SalesOrderItem.Save("Quantity")

	sales, err := ShowSalesOrderFulfillment("id", fulfill.SalesOrderItem.SalesOrder.ID)
	assert.NotNil(t, sales, "Tidak boleh kosong jika berhasil mengambil data")
	assert.NoError(t, err, "Tidak boleh ada error")
}
func TestApproveFulfillment(t *testing.T) {
	fulfillment := model.DummyWorkorderFulfillment()
	so := model.DummySalesOrder()

	iv := model.DummyItemVariant()
	iv.CommitedStock = 20
	iv.Save("CommitedStock")

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
	iv2.CommitedStock = 20
	iv2.Save("CommitedStock")

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
	fulfillment.Save("WorkorderFulFillmentItems", "SalesOrder")

	fulfillment, e := approveFulfillment(fulfillment)
	assert.NoError(t, e)
	assert.NotEmpty(t, fulfillment)
}

func TestApproveFulfillment2(t *testing.T) {
	fulfillment := model.DummyWorkorderFulfillment()
	so := model.DummySalesOrder()

	iv := model.DummyItemVariant()
	iv.CommitedStock = 20
	iv.Save("CommitedStock")

	soitem := model.DummySalesOrderItem()
	soitem.ItemVariant = iv
	soitem.Quantity = 10
	soitem.SalesOrder = so
	soitem.Save("ItemVariant", "SalesOrder", "Quantity")

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
	iv2.CommitedStock = 20
	iv2.Save("CommitedStock")

	soItem2 := model.DummySalesOrderItem()
	soItem2.ItemVariant = iv2
	soItem2.SalesOrder = so
	soItem2.Quantity = 10
	soItem2.Save("ItemVariant", "SalesOrder", "Quantity")

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
	fulfillment.SalesOrder.InvoiceStatus = "finished"
	fulfillment.SalesOrder.Save("InvoiceStatus")
	fulfillment.Save("WorkorderFulFillmentItems", "SalesOrder")

	fulfillment, e := approveFulfillment(fulfillment)
	assert.NoError(t, e)
	assert.NotEmpty(t, fulfillment)
}
func TestUpdateItemVariantAndStock(t *testing.T) {
	fulfillment := model.DummyWorkorderFulfillment()

	iv := model.DummyItemVariant()
	iv.CommitedStock = 20
	iv.Save("CommitedStock")

	soitem := model.DummySalesOrderItem()
	soitem.ItemVariant = iv
	soitem.Save("ItemVariant")

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
	iv2.CommitedStock = 20
	iv2.Save("CommitedStock")

	soItem2 := model.DummySalesOrderItem()
	soItem2.ItemVariant = iv2
	soItem2.Save("ItemVariant")

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

	var items []*model.WorkorderFulfillmentItem
	items = append(items, item1, item2)

	fulfillment.WorkorderFulFillmentItems = items
	fulfillment.Save("WorkorderFulFillmentItems")

	totCost, e := updateItemVariantStock(fulfillment)
	assert.NoError(t, e)
	assert.Equal(t, float64(20000), totCost)

}
func TestGetWorkorderFulfillmentByID(t *testing.T) {
	fulfillment := model.DummyWorkorderFulfillment()
	fulfillment.IsDeleted = 0
	fulfillment.Save("IsDeleted")

	item := model.DummyWorkorderFulfillmentItem()
	item.WorkorderFulfillment = fulfillment
	item.Save("WorkorderFulfillment")

	soItem := model.DummySalesOrderItem()
	soItem.SalesOrder = fulfillment.SalesOrder
	soItem.Save("SalesOrder")

	var soItems []*model.SalesOrderItem
	soItems = append(soItems, soItem)

	fulfillment.SalesOrder.SalesOrderItems = soItems
	fulfillment.SalesOrder.Save("SalesOrderItem")

	data, e := getWorkorderFulfillmentByID(fulfillment.ID)
	assert.NoError(t, e)
	assert.Equal(t, int8(0), data.IsDeleted)
	assert.Equal(t, fulfillment.Note, data.Note)
	assert.Equal(t, fulfillment.Code, data.Code)
	assert.Equal(t, fulfillment.SalesOrder.ID, data.SalesOrder.ID)
	assert.Equal(t, fulfillment.DocumentStatus, data.DocumentStatus)
	assert.NotEmpty(t, fulfillment.SalesOrder.SalesOrderItems)
	assert.NotEmpty(t, fulfillment.SalesOrder.SalesOrderItems[0].Note)
	assert.NotEmpty(t, fulfillment.SalesOrder.SalesOrderItems[0].ItemVariant.ID)
	assert.NotEmpty(t, fulfillment.SalesOrder.SalesOrderItems[0].ItemVariant.Note)

	fulfillment.IsDeleted = 1
	fulfillment.Save("IsDeleted")
	data, e = getWorkorderFulfillmentByID(fulfillment.ID)
	assert.Error(t, e)
	assert.Empty(t, data)
}
