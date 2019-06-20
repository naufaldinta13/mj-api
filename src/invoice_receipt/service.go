// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package invoiceReceipt

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/sales_invoice"

	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
)

// GetInvoiceReceipts get all data Invoice Receipt that matched with query request parameters.
func GetInvoiceReceipts(rq *orm.RequestQuery) (m *[]model.InvoiceReceipt, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.InvoiceReceipt))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.InvoiceReceipt
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetDetailInvoiceReceipt untuk get detail invoice receipt by id
func GetDetailInvoiceReceipt(field string, values ...interface{}) (*model.InvoiceReceipt, error) {
	m := new(model.InvoiceReceipt)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	// get all sales order based on invoice_receipt_item and invoice_receipt_return
	// harus di distinct agar sales_order yang di dapat tidak duplikat
	if _, err := o.Raw(`select distinct so.* from invoice_receipt ir
			LEFT join invoice_receipt_item iri on iri.invoice_receipt_id = ir.id
			LEFT join sales_invoice si on si.id = iri.sales_invoice_id
			LEFT join invoice_receipt_return irr on irr.invoice_receipt_id = ir.id
			LEFT join sales_return sr on sr.id = irr.sales_return_id
			LEFT join sales_order so on so.id = si.sales_order_id or so.id = sr.sales_order_id
			where ir.id = ?`, m.ID).QueryRows(&m.SalesOrders); err == nil {

		for _, so := range m.SalesOrders {
			// invoice_receipt_item
			qs := o.QueryTable("invoice_receipt_item")
			qs.Filter("sales_invoice_id__sales_order_id__id", so.ID).Filter("invoice_receipt_id__id", m.ID).RelatedSel("sales_invoice_id").All(&so.InvoiceReceiptItems)

			// invoice_receipt_return
			qsr := o.QueryTable("invoice_receipt_return")
			qsr.Filter("sales_return_id__sales_order_id__id", so.ID).Filter("invoice_receipt_id__id", m.ID).RelatedSel("sales_return_id").All(&so.InvoiceReceiptReturns)
		}
	}

	return m, nil
}

// getDetailInvoiceReceipt untuk get detail invoice receipt by id
func getDetailInvoiceReceipt(field string, values ...interface{}) (*model.InvoiceReceipt, error) {
	m := new(model.InvoiceReceipt)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "InvoiceReceiptItems", 1)
	o.LoadRelated(m, "InvoiceReceiptReturns", 1)

	return m, nil
}

// GetSalesOrdersByCustomerID untuk mengambil banyak data sales order by customer id
func GetSalesOrdersByCustomerID(customerID int64) (m []*model.SalesOrder, err error) {

	o := orm.NewOrm()
	if _, err := o.Raw("select * from sales_order where customer_id = ?", customerID).QueryRows(&m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetSalesInvoiceBySOID untuk mengambil banyak data sales invoice by sales order id
func GetSalesInvoiceBySOID(soID int64) (m []*model.SalesInvoice, err error) {

	o := orm.NewOrm()
	if _, err := o.Raw("select * from sales_invoice where sales_order_id = ?", soID).QueryRows(&m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetNotBundledSalesInvoiceBySOID untuk mengambil banyak data sales invoice by sales order id yang belum dibuat invoice receipt
func GetNotBundledSalesInvoiceBySOID(soID int64) (m []*model.SalesInvoice, err error) {

	o := orm.NewOrm()
	if _, err := o.Raw("select * from sales_invoice where is_bundled = 0 AND sales_order_id = ?", soID).QueryRows(&m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetNotBundledSalesReturnBySOID untuk mengambil banyak data sales return by sales order id yang belum dibuat invoice receipt
func GetNotBundledSalesReturnBySOID(soID int64) (m []*model.SalesReturn, err error) {

	o := orm.NewOrm()
	if _, err := o.Raw("select * from sales_return where is_bundled = 0 AND sales_order_id = ?", soID).QueryRows(&m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetSalesReturnBySOID untuk mengambil banyak data sales return by sales order id
func GetSalesReturnBySOID(soID int64) (m []*model.SalesReturn, err error) {

	o := orm.NewOrm()
	if _, err := o.Raw("select * from sales_return where sales_order_id = ?", soID).QueryRows(&m); err != nil {
		return nil, err
	}

	return m, nil
}

// CheckInvoiceStatusSO untuk cek invoice status active pada SO by customer id
func CheckInvoiceStatusSO(customerID int64) bool {

	// ambil sales order by customer id
	m, _ := GetSalesOrdersByCustomerID(customerID)

	for _, x := range m {
		// jika terdapat sales order dengan invoice status != active maka akan berhenti
		if x.InvoiceStatus != "active" {
			return false
		}
	}

	return true
}

// CreateInvoiceReceiptItem untuk membuat invoice receipt item
func CreateInvoiceReceiptItem(salesInvoice *model.SalesInvoice, invoiceReceipt *model.InvoiceReceipt) error {
	var e error
	var amountFinanceRevenue float64

	// sum amount pada finance revenue untuk perhitungan subtotal pada invoice receipt item
	o := orm.NewOrm()
	o.Raw("select sum(amount) from finance_revenue where ref_type = 'sales_invoice' and ref_id = ? and is_deleted = ?", salesInvoice.ID, 0).QueryRow(&amountFinanceRevenue)

	m := &model.InvoiceReceiptItem{
		SalesInvoice:   salesInvoice,
		InvoiceReceipt: invoiceReceipt,
	}

	// cek document status di sales invoice
	if salesInvoice.DocumentStatus == "active" {
		// jika statusnya active maka perhitungannya
		// total amount (pada sales invoice) - hasil sum amount (pada finance revenue)
		m.Subtotal = salesInvoice.TotalAmount - amountFinanceRevenue
	} else {
		// jika statusnya selain active maka
		// hanya mengambil total amount dari sales invoice
		m.Subtotal = salesInvoice.TotalAmount
	}

	if e = m.Save(); e != nil {
		return e
	}

	return nil
}

// CreateInvoiceReceiptReturn untuk membuat invoice receipt return
func CreateInvoiceReceiptReturn(salesReturn *model.SalesReturn, invoiceReceipt *model.InvoiceReceipt) error {
	var e error
	var amountFinanceExpense float64

	// sum amount pada finance expense untuk perhitungan subtotal pada invoice receipt return
	o := orm.NewOrm()
	o.Raw("select sum(amount) from finance_expense ref_type = 'sales_return' and where ref_id = ? and is_deleted = ?", salesReturn.ID, 0).QueryRow(&amountFinanceExpense)

	m := &model.InvoiceReceiptReturn{
		SalesReturn:    salesReturn,
		InvoiceReceipt: invoiceReceipt,
	}

	// cek document status di sales return
	if salesReturn.DocumentStatus == "active" {
		// jika statusnya active maka perhitungannya
		// total amount (pada sales return) - hasil sum amount (pada finance expense)
		m.Subtotal = salesReturn.TotalAmount - amountFinanceExpense
	} else {
		// jika statusnya selain active maka
		// hanya mengambil total amount dari sales return
		m.Subtotal = salesReturn.TotalAmount
	}

	if e = m.Save(); e != nil {
		return e
	}

	return nil
}

// CreateInvoiceReceipt untuk membuat data invoice receipt beserta invoice receipt item
func CreateInvoiceReceipt(invoiceReceipt *model.InvoiceReceipt, salesOrders []*model.SalesOrder) error {

	var e error
	var salesInv []*model.SalesInvoice
	var salesRetrun []*model.SalesReturn

	// dapatkan data-data sales porder berdasarkan customer
	for _, x := range salesOrders {
		// dapatkan sales invoice berdasarkan sales order
		si, _ := GetSalesInvoiceBySOID(x.ID)
		for _, usi := range si {
			salesInv = append(salesInv, usi)
		}

		// dapatkan sales return berdasarkan sales order
		sr, _ := GetSalesReturnBySOID(x.ID)
		for _, usr := range sr {
			salesRetrun = append(salesRetrun, usr)
		}
	}

	for _, mx := range salesInv {
		// cek pada sales invoice jika is bundled == 0
		if mx.IsBundled == 0 {
			// cek pada sales invoice jika document status != finished
			if mx.DocumentStatus != "finished" {
				invoiceReceipt.Save()

				// buat data invoice receipt item
				CreateInvoiceReceiptItem(mx, invoiceReceipt)

				// ubah document status menjadi active
				mx.DocumentStatus = "active"
				// ubah is bundled menjadi 1
				mx.IsBundled = int8(1)
				if e = mx.Save("document_status", "is_bundled"); e != nil {
					return e
				}

				// sum subtotal pada invoice receipt item
				// untuk field total invoice pada invoice receipt
				var totalInvoice float64
				o := orm.NewOrm()
				o.Raw("select sum(subtotal) from invoice_receipt_item where invoice_receipt_id = ?", invoiceReceipt.ID).QueryRow(&totalInvoice)

				// masukan hasil sum subtotal ke field invoice receipt
				invoiceReceipt.TotalInvoice = totalInvoice
				invoiceReceipt.Save("total_invoice")

			}
		}
	}

	for _, ms := range salesRetrun {
		// cek pada sales invoice jika is bundled == 0
		if ms.IsBundled == 0 {
			// cek pada sales invoice jika document status != finished
			if ms.DocumentStatus != "finished" {

				// buat data invoice receipt return
				CreateInvoiceReceiptReturn(ms, invoiceReceipt)

				// ubah document status menjadi active
				ms.DocumentStatus = "active"

				// ubah is bundled menjadi 1
				ms.IsBundled = int8(1)
				if e = ms.Save("document_status", "is_bundled"); e != nil {
					return e
				}

				// sum subtotal pada invoice receipt return
				// untuk field total return pada invoice receipt
				var totalReturn float64
				o := orm.NewOrm()
				o.Raw("select sum(subtotal) from invoice_receipt_return where invoice_receipt_id = ?", invoiceReceipt.ID).QueryRow(&totalReturn)

				// masukan hasil sum subtotal ke field invoice receipt
				invoiceReceipt.TotalReturn = totalReturn
				invoiceReceipt.Save("total_return")

			}
		}
	}

	invoiceReceipt.TotalAmount = invoiceReceipt.TotalInvoice - invoiceReceipt.TotalReturn
	// cek total amount pada invoice receipt
	if invoiceReceipt.TotalAmount >= 0 {
		invoiceReceipt.Save("TotalAmount")
	} else {
		e = validation.SetError("total_amount", "total_amount can't be less then 0")
	}

	for _, sInv := range salesInv {
		salesInvoice.SumTotalRevenuedSalesInvoice(sInv.ID)
	}

	return e
}

// FetchInvoiceReceiptItemPayment fetch data invoice receipt item untuk membuat invoice receipt payment
func FetchInvoiceReceiptItemPayment(InvoiceReceiptID int64, invoiceReceiptItems []*model.InvoiceReceiptItem) error {
	var err error
	var ri *model.InvoiceReceipt

	type ActionInvoiceReceiptItem struct {
		InvoiceReceiptItem *model.InvoiceReceiptItem
		Action             string
	}

	var InvoiceReceiptItemAction = make([]*ActionInvoiceReceiptItem, 0)

	if ri, err = getDetailInvoiceReceipt("id", InvoiceReceiptID); err == nil {
		for _, v := range ri.InvoiceReceiptItems {
			action := &ActionInvoiceReceiptItem{InvoiceReceiptItem: v, Action: "delete"}
			InvoiceReceiptItemAction = append(InvoiceReceiptItemAction, action)
		}
	}

	for _, x := range invoiceReceiptItems {
		for _, c := range InvoiceReceiptItemAction {
			if c.InvoiceReceiptItem.ID == x.ID {
				c.Action = "save"
			}

		}
	}

	var totalInvoice float64
	for _, vir := range InvoiceReceiptItemAction {
		if vir.Action == "delete" {
			// update is bundle sales invoice
			vir.InvoiceReceiptItem.SalesInvoice.IsBundled = 0
			vir.InvoiceReceiptItem.SalesInvoice.Save("IsBundled")
			totalInvoice += vir.InvoiceReceiptItem.Subtotal
			vir.InvoiceReceiptItem.Delete()

			// hitung ulang total_revenued nya
			salesInvoice.SumTotalRevenuedSalesInvoice(vir.InvoiceReceiptItem.SalesInvoice.ID)
		}
	}

	invoiceReceipt := new(model.InvoiceReceipt)
	invoiceReceipt.ID = InvoiceReceiptID
	if err = invoiceReceipt.Read("ID"); err == nil {
		invoiceReceipt.TotalInvoice = invoiceReceipt.TotalInvoice - totalInvoice
		invoiceReceipt.Save("TotalInvoice")
	}

	return err

}

// FetchInvoiceReceiptReturnPayment fetch data invoice receipt return untuk membuat invoice receipt payment
func FetchInvoiceReceiptReturnPayment(InvoiceReceiptID int64, invoiceReceiptReturns []*model.InvoiceReceiptReturn) error {
	var err error
	var ri *model.InvoiceReceipt

	type ActionInvoiceReceiptReturn struct {
		InvoiceReceiptReturn *model.InvoiceReceiptReturn
		Action               string
	}

	var InvoiceReceiptReturnAction = make([]*ActionInvoiceReceiptReturn, 0)

	if ri, err = getDetailInvoiceReceipt("id", InvoiceReceiptID); err == nil {
		for _, n := range ri.InvoiceReceiptReturns {
			action := &ActionInvoiceReceiptReturn{InvoiceReceiptReturn: n, Action: "delete"}
			InvoiceReceiptReturnAction = append(InvoiceReceiptReturnAction, action)
		}
	}

	for _, x := range invoiceReceiptReturns {
		for _, c := range InvoiceReceiptReturnAction {
			if c.InvoiceReceiptReturn.ID == x.ID {
				c.Action = "save"
			}

		}
	}

	var totalReturn float64
	for _, vir := range InvoiceReceiptReturnAction {
		if vir.Action == "delete" {
			// update is bundle sales invoice
			vir.InvoiceReceiptReturn.SalesReturn.IsBundled = 0
			vir.InvoiceReceiptReturn.SalesReturn.Save("IsBundled")
			totalReturn += vir.InvoiceReceiptReturn.Subtotal

			vir.InvoiceReceiptReturn.Delete()
		}
	}

	invoiceReceipt := new(model.InvoiceReceipt)
	invoiceReceipt.ID = InvoiceReceiptID
	if err = invoiceReceipt.Read("ID"); err == nil {
		invoiceReceipt.TotalReturn = invoiceReceipt.TotalReturn - totalReturn
		invoiceReceipt.Save("TotalReturn")
	}

	return err

}

// CreateInvoiceReceiptPayment untuk membuat data invoice receipt payment
func CreateInvoiceReceiptPayment(InvoiceReceiptID int64, invoiceReceiptItems []*model.InvoiceReceiptItem, invoiceReceiptReturns []*model.InvoiceReceiptReturn, financeRevenues []*model.FinanceRevenue) error {
	var err error

	FetchInvoiceReceiptItemPayment(InvoiceReceiptID, invoiceReceiptItems)
	FetchInvoiceReceiptReturnPayment(InvoiceReceiptID, invoiceReceiptReturns)

	invoiceReceipt := new(model.InvoiceReceipt)
	invoiceReceipt.ID = InvoiceReceiptID
	if err = invoiceReceipt.Read("ID"); err == nil {
		invoiceReceipt.TotalAmount = invoiceReceipt.TotalInvoice - invoiceReceipt.TotalReturn
		invoiceReceipt.DocumentStatus = "finished"
		invoiceReceipt.Save("TotalAmount", "DocumentStatus")
	}

	for _, mx := range financeRevenues {
		if e := mx.Save(); e != nil {
			return e
		}
	}

	o := orm.NewOrm()
	var m []*model.InvoiceReceiptItem
	o.Raw("select * from invoice_receipt_item where invoice_receipt_id = ?", InvoiceReceiptID).QueryRows(&m)

	var msi *model.SalesInvoice
	var amountSI float64
	for _, mir := range m {
		// update document_status pada sales_invoice menjadi finished
		msi = mir.SalesInvoice
		msi.Read("ID")
		msi.TotalPaid = msi.TotalPaid + mir.Subtotal

		// ambil total_amount dari slaes_invoice
		o.Raw("select total_amount from sales_invoice where id = ? ", mir.SalesInvoice.ID).QueryRow(&amountSI)
		// cek total paid si
		if msi.TotalAmount == msi.TotalPaid {
			msi.DocumentStatus = "finished"
		}

		msi.Save("DocumentStatus", "TotalPaid")
		salesInvoice.SumTotalRevenuedSalesInvoice(msi.ID)

		/////////////////////////////
		var so *model.SalesOrder
		o.Raw("select * from sales_order where id = ?", msi.SalesOrder.ID).QueryRow(&so)

		// update total_paid sales_order
		so.TotalPaid += amountSI
		so.Save("TotalPaid")

		var partner *model.Partnership
		o.Raw("select * from partnership where id = ?", so.Customer.ID).QueryRow(&partner)

		// update total_debt partnership
		partner.TotalDebt -= so.TotalPaid
		partner.Save("TotalDebt")

		// update invoice_status menjadi finished atau active pada sales_order
		so.Read()
		if so.TotalPaid == so.TotalCharge {
			// jika jumlah total_paid sama dengan total_charge
			// update invoice_status menjadi finished
			so.InvoiceStatus = "finished"
			so.Save("InvoiceStatus")
		} else {
			// jika jumlah total_amount tidak sama dengan total_charge
			// update invoice_status menjadi active
			so.InvoiceStatus = "active"
			so.Save("InvoiceStatus")
		}

		// update document_status menjadi finished atau partial pada sales_order
		so.Read()
		if so.InvoiceStatus == "finished" && so.FulfillmentStatus == "finished" {
			// jika invoice_status dan fulfillment_status sama dengan finished
			// update document_status menjadi finished
			so.DocumentStatus = "finished"
			so.Save("DocumentStatus")
		} else {
			// jika invoice_status dan fulfillment_status tidak sama dengan finished
			// update document_status menjadi active
			so.DocumentStatus = "active"
			so.Save("DocumentStatus")
		}
	}

	// cek total return invoice receipt
	if invoiceReceipt.TotalReturn != 0 {
		// buat data finance revenue
		financeRevenue := new(model.FinanceRevenue)
		financeRevenue.RefType = "invoice_receipt"
		financeRevenue.RefID = uint64(invoiceReceipt.ID)
		financeRevenue.PaymentMethod = "cash"
		financeRevenue.DocumentStatus = "cleared"
		financeRevenue.Amount = invoiceReceipt.TotalReturn
		financeRevenue.Save()

		// buat data finance expense
		financeExpense := new(model.FinanceExpense)
		financeExpense.RefType = "invoice_receipt"
		financeExpense.RefID = uint64(invoiceReceipt.ID)
		financeExpense.PaymentMethod = "cash"
		financeExpense.DocumentStatus = "cleared"
		financeExpense.Amount = invoiceReceipt.TotalReturn
		financeExpense.Save()
	}

	var salesReturn = &model.SalesReturn{}
	for _, uu := range invoiceReceiptReturns {
		salesReturn.ID = uu.SalesReturn.ID
		if err = salesReturn.Read("ID"); err == nil {
			var sumexp float64
			o.Raw("select sum(amount) from finance_expense where ref_type = ? and ref_id = ? and is_delete = 0", "sales_return", uu.SalesReturn.ID).QueryRow(&sumexp)
			sumexp += uu.Subtotal

			// harus sesuai dengan finance expenses dengan total amount
			if sumexp == uu.SalesReturn.TotalAmount {
				salesReturn.DocumentStatus = "finished"
				salesReturn.Save("DocumentStatus")
			}
		}
	}

	return nil
}
