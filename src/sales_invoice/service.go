// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package salesInvoice

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/orm"
	"git.qasico.com/mj/api/src/sales"
)

// GetSalesInvoice untuk mengambil semua data sales invoice dari database
func GetSalesInvoice(rq *orm.RequestQuery) (m *[]model.SalesInvoice, total int64, e error) {
	// make new orm query
	q, _ := rq.Query(new(model.SalesInvoice))

	nc := orm.NewCondition()
	ncc := nc.And("recognition_date__gte", util.MaxDate()).OrNot("document_status", "finished")
	nccx := nc.And("is_deleted", 0)

	c := nc.AndCond(rq.GetCondition()).AndCond(nccx).AndCond(ncc)
	q = q.SetCond(c)

	// get total data
	if total, e = q.Count(); e == nil && total != int64(0) {
		// get data requested
		var mx []model.SalesInvoice
		if _, e = q.All(&mx, rq.Fields...); e == nil {
			m = &mx
		}
	}
	return
}

// ShowSalesInvoice untuk mengambil data sales invoice berdasarkan param
func ShowSalesInvoice(field string, values ...interface{}) (*model.SalesInvoice, error) {
	m := new(model.SalesInvoice)
	o := orm.NewOrm().QueryTable(m)
	if err := o.Filter(field, values...).Filter("is_deleted", int8(0)).RelatedSel(3).Limit(1).One(m); err != nil {

		return nil, err
	}

	return m, nil
}

// GetSumTotalAmountSalesInvoiceBySalesOrder to get total amount from all sales invoice by so id
func GetSumTotalAmountSalesInvoiceBySalesOrder(SalesOrderID int64) (sumTotAmount float64, e error) {
	o := orm.NewOrm()
	e = o.Raw("SELECT sum(total_amount) from sales_invoice si where si.sales_order_id = ? AND si.is_deleted = 0", SalesOrderID).QueryRow(&sumTotAmount)
	return
}

// CreateSalesInvoice untuk menyimpan data sales invoice dan mengganti invoice status
func CreateSalesInvoice(SI *model.SalesInvoice) (e error) {

	// save sales invoice
	if e = SI.Save(); e == nil {
		SI.SalesOrder.InvoiceStatus = "active"
		SI.SalesOrder.DocumentStatus = "active"
		e = SI.SalesOrder.Save("invoice_status", "DocumentStatus")
	}

	sales.CalculateTotalPaidSI(SI)
	sales.CheckInvoiceStatus(SI.SalesOrder)
	sales.CheckDocumentStatus(SI.SalesOrder)

	return
}

// SumTotalRevenuedSalesInvoice untuk menjumlahkan semua amount pada finance revenue berdasarkan sales invoice
func SumTotalRevenuedSalesInvoice(salesInvoiceID int64) (totalRevenued float64, e error) {
	o := orm.NewOrm()

	// ambil jumlah amount finance revenue
	var amountFinanceRevenue float64
	o.Raw("select sum(amount) from finance_revenue where ref_id = ? AND ref_type = 'sales_invoice' "+
		"AND is_deleted = 0", uint64(salesInvoiceID)).QueryRow(&amountFinanceRevenue)

	// ambil jumlah amount invoice receipt
	var amountInvoiceReceiptItem float64
	o.Raw("select subtotal from invoice_receipt_item where sales_invoice_id = ?", salesInvoiceID).QueryRow(&amountInvoiceReceiptItem)

	totalRevenued = amountFinanceRevenue + amountInvoiceReceiptItem

	// save total revenued sales invoice
	si := &model.SalesInvoice{ID: salesInvoiceID}
	si.TotalRevenued = totalRevenued
	e = si.Save("TotalRevenued")

	return
}
