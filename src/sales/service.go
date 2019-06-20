// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sales

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/inventory"
	"git.qasico.com/mj/api/src/partnership"
	"git.qasico.com/mj/api/src/stock"
	"git.qasico.com/mj/api/src/util"
)

// CalculateTotalPaidSI by Sales Invoice
func CalculateTotalPaidSI(si *model.SalesInvoice) error {

	o := orm.NewOrm()

	// sum amount
	_, err := o.Raw("update sales_invoice si set si.total_paid = (SELECT sum(amount) from finance_revenue where ref_id = ? AND ref_type = 'sales_invoice' AND document_status = 'cleared') where si.id = ?; ", uint64(si.ID), si.ID).Exec()

	si.Read()

	if si.TotalPaid >= si.TotalAmount {
		o.Raw("update sales_invoice si set si.document_status = 'finished' where si.id = ?; ", si.ID).Exec()
	}

	CalculateTotalPaidSO(si.SalesOrder)
	return err
}

// CalculateTotalPaidSO by sales order
func CalculateTotalPaidSO(so *model.SalesOrder) error {

	o := orm.NewOrm()

	// sum total paid
	_, err := o.Raw("update sales_order so set so.total_paid = (SELECT sum(total_paid) from sales_invoice where sales_order_id = ?) where so.id = ?", so.ID, so.ID).Exec()
	return err

}

// CheckFulfillmentStatus update fulfillment status to finished
func CheckFulfillmentStatus(so *model.SalesOrder) (e error) {
	o := orm.NewOrm()

	var quantitySOItem float64

	if e = o.Raw("SELECT sum(quantity) from sales_order_item soi inner join sales_order so on soi.sales_order_id = so.id"+
		" where so.document_status != 'approved_cancel' AND so.is_deleted = 0"+
		" AND sales_order_id = ?", so.ID).QueryRow(&quantitySOItem); e != nil {
		return e
	}

	var quantityFItem float64
	if e = o.Raw("select sum(wfi.quantity) from workorder_fulfillment_item wfi "+
		"inner join workorder_fulfillment wf on wf.id = wfi.workorder_fulfillment_id "+
		"inner join sales_order_item soi on soi.id = wfi.sales_order_item_id "+
		"inner join sales_order so on so.id = soi.sales_order_id "+
		"where wf.is_deleted = 0 and wf.document_status = 'finished' and so.document_status != 'approved_cancel' "+
		"and so.is_deleted = 0 and wf.sales_order_id = ?", so.ID).QueryRow(&quantityFItem); e != nil {

		return e
	}

	if quantitySOItem > 0 && quantityFItem == quantitySOItem {
		_, e = o.Raw("update sales_order so set so.fulfillment_status = 'finished' where id = ?", so.ID).Exec()
		return e
	}

	CheckDocumentStatus(so)

	return nil
}

// CheckInvoiceStatus update invoice status to finished
func CheckInvoiceStatus(so *model.SalesOrder) (e error) {
	CalculateTotalPaidSO(so)
	so.Read()

	o := orm.NewOrm()
	if so.TotalPaid == so.TotalCharge {
		_, e = o.Raw("update sales_order so set so.invoice_status = 'finished' where id = ?", so.ID).Exec()
		return e
	}
	return nil
}

// CheckDocumentStatus update document status to finished
func CheckDocumentStatus(so *model.SalesOrder) (e error) {
	o := orm.NewOrm()

	if so.InvoiceStatus == "finished" && so.FulfillmentStatus == "finished" {
		_, e = o.Raw("update sales_order so set so.invoice_status = 'finished' where id = ?", so.ID).Exec()
		return e
	}
	return nil
}

// CanBeReturnSales fungsion check quantity return
func CanBeReturnSales(so *model.SalesOrder, sr *model.SalesReturn) (soi []*model.SalesOrderItem, e error) {
	o := orm.NewOrm()
	var soITEM []*model.SalesOrderItem
	var quantityFItem float32
	var quantitySRItem float32
	o.QueryTable(new(model.SalesOrderItem)).Filter("sales_order_id__id", so.ID).RelatedSel().All(&soITEM)
	for i := 0; i < len(soITEM); i++ {

		if e = o.Raw("select sum(wfi.quantity) from workorder_fulfillment_item wfi "+
			"inner join workorder_fulfillment wf on wf.id = wfi.workorder_fulfillment_id "+
			"inner join sales_order_item soi on soi.id = wfi.sales_order_item_id "+
			"inner join sales_order so on so.id = soi.sales_order_id "+
			"where wf.document_status = 'finished' and wf.is_deleted = 0 and wf.is_delivered = 1 and so.document_status != 'approved_cancel' "+
			"and so.is_deleted = 0 and wf.sales_order_id = ? and wfi.sales_order_item_id = ?", so.ID, soITEM[i].ID).QueryRow(&quantityFItem); e != nil {

			return nil, e
		}

		//ini untuk create
		if sr == nil {
			if e = o.Raw("SELECT sum(quantity) from sales_return_item sri inner join sales_return sr on sri.sales_return_id = sr.id"+
				" where sr.is_deleted = 0"+
				" AND sr.sales_order_id = ? AND sr.document_status != 'cancelled' AND sri.sales_order_item_id = ?", so.ID, soITEM[i].ID).QueryRow(&quantitySRItem); e != nil {
				return nil, e
			}
			//ini untuk update
		} else {
			if e = o.Raw("SELECT sum(quantity) from sales_return_item sri inner join sales_return sr on sri.sales_return_id = sr.id"+
				" where sr.is_deleted = 0"+
				" AND sr.sales_order_id = ? AND sr.document_status != 'cancelled' AND sri.sales_order_item_id = ?"+
				" AND sr.ID != ?", so.ID, soITEM[i].ID, sr.ID).QueryRow(&quantitySRItem); e != nil {
				return nil, e
			}
		}
		if quantitySRItem <= 0 {
			soITEM[i].CanBeReturn += quantityFItem
		} else {
			soITEM[i].CanBeReturn += (quantityFItem - quantitySRItem)
		}

	}

	soi = soITEM

	return soi, nil
}

// GetAllSalesOrder untuk mengambil semua data sales order dari database
func GetAllSalesOrder(rq *orm.RequestQuery) (m *[]model.SalesOrder, total int64, e error) {
	// make new orm query
	q, _ := rq.Query(new(model.SalesOrder))

	nc := orm.NewCondition()
	ncc := nc.And("recognition_date__gte", util.MaxDate()).OrNot("invoice_status", "finished")
	nccx := nc.And("is_deleted", 0)

	c := nc.AndCond(rq.GetCondition()).AndCond(nccx).AndCond(ncc)
	q = q.SetCond(c)

	// get total data
	if total, e = q.Count(); e == nil && total != int64(0) {
		// get data requested
		var mx []model.SalesOrder
		if _, e = q.All(&mx, rq.Fields...); e == nil {
			m = &mx
		}
	}
	return
}

// GetAllSalesOrderOld untuk mengambil semua data sales order dari database
func GetAllSalesOrderOld(rq *orm.RequestQuery) (m *[]model.SalesOrder, total int64, e error) {
	// make new orm query
	q, _ := rq.Query(new(model.SalesOrder))

	nc := orm.NewCondition()
	nccx := nc.And("is_deleted", 0)

	c := nc.AndCond(rq.GetCondition()).AndCond(nccx)
	q = q.SetCond(c)

	// get total data
	if total, e = q.Count(); e == nil && total != int64(0) {
		// get data requested
		var mx []model.SalesOrder
		if _, e = q.All(&mx, rq.Fields...); e == nil {
			m = &mx
		}
	}
	return
}

// GetDetailSalesOrder untuk mengambil data detail sales order berdasarkan param
func GetDetailSalesOrder(id int64, loadRelated []string) (*model.SalesOrder, error) {
	m := new(model.SalesOrder)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter("id", id).Filter("is_deleted", int8(0)).RelatedSel(1).Limit(1).One(m); err != nil {
		return nil, err
	}

	if util.HasElem(loadRelated, "sales_order_items") {
		// sales order items
		m.SalesOrderItems, _ = CanBeReturnSales(m, nil)
	}

	if util.HasElem(loadRelated, "sales_invoices") {
		// sales invoices
		o.Raw("select * from sales_invoice where sales_order_id = ? and is_deleted = ?", m.ID, 0).QueryRows(&m.SalesInvoices)
	}

	if util.HasElem(loadRelated, "workorder_fulfillments") {
		// workorder fulfillment
		o.Raw("select * from workorder_fulfillment where sales_order_id = ? and is_deleted = ?", m.ID, 0).QueryRows(&m.WorkorderFulfillments)
	}

	if util.HasElem(loadRelated, "sales_returns") {
		var srids string //sales return id id
		var tr float64   // total sales refund
		var tpsr float64 // total paid sales refund

		// workorder fulfillment
		o.Raw("select * from sales_return where sales_order_id = ? and is_deleted = ?", m.ID, 0).QueryRows(&m.SalesReturns)
		o.Raw("select group_concat(sr.id) from sales_return sr where sales_order_id = ? and is_deleted = ?", m.ID, 0).QueryRow(&srids)

		o.Raw("select sum(total_amount) from sales_return where sales_order_id = ? and is_deleted = ?", m.ID, 0).QueryRow(&tr)
		o.Raw("select sum(amount) from finance_expense where ref_type = ? and ref_id in ("+srids+" ) and is_deleted = ?", "sales_return", 0).QueryRow(&tpsr)

		m.TotalRefund = tr
		m.TotalPaidRefund = tpsr
	}

	return m, nil
}

// CreateSalesOrder untuk simpan sales order
func CreateSalesOrder(order *createRequest) (sales *model.SalesOrder, err error) {
	sales = order.Transform()
	
	// Masukkan data ke dalam database sesuai dengan list inputan
	if err = sales.Save(); err != nil {
		return nil, err
	}
	// wofulfillmentitem
	var wofulfillmentitems []model.WorkorderFulfillmentItem
	// Simpan Sales Order Items
	for _, row := range sales.SalesOrderItems {
		row.SalesOrder = &model.SalesOrder{ID: sales.ID}
		if err = row.Save(); err != nil {
			return nil, err
		}
		row.ItemVariant.Read()
		// Update avalibale stock dengan available stock lama dikurang dengan quantity input.
		// dan update commited stock dengan commited stock lama ditambah dengan quantity input
		// pada item variant yang dipilih
		_, err = stock.CalculateAvailableStockItemVariant(row.ItemVariant)
		if err != nil {
			return nil, err
		}
		// Update total debt dengan total debt lama ditambah dengan total charge SO pada partner ship SO tersebut
		if err = partnership.CalculationTotalDebt(sales.Customer.ID); err == nil {
			if err = partnership.CalculationTotalSpend(sales.Customer.ID); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
		wofulfillmentitem := model.WorkorderFulfillmentItem{
			SalesOrderItem: row,
			Quantity:       row.Quantity,
		}
		wofulfillmentitems = append(wofulfillmentitems, wofulfillmentitem)
	}

	var partner *model.Partnership
	partner, err = partnership.GetPartnershipByField("id", sales.Customer.ID)
	if err != nil {
		return nil, err
	}

	if sales.AutoInvoice == 1 || partner.IsDefault == 1 || sales.AutoPaid == 1 {
		// Masukkan data sales invoice dengan referensi sales order id dan semua sales order item id dengan quantity yang sama
		// (untuk data yang diinput dapat melihat list inputan)
		code, err := util.CodeGen("code_sales_invoice", "sales_invoice")
		if err != nil {
			return nil, err
		}
		sales.Customer.Read()
		sales.Read()
		sinvoice := model.SalesInvoice{
			SalesOrder:      sales,
			Code:            code,
			BillingAddress:  sales.Customer.Address,
			RecognitionDate: sales.RecognitionDate,
			TotalAmount:     sales.TotalCharge,
			DocumentStatus:  "new",
			CreatedAt:       sales.CreatedAt,
			CreatedBy:       order.Session.User,
		}
		if err = sinvoice.Save(); err != nil {
			return nil, err
		}

		// cek apakah yang membuat sales order walk-in customer
		if partner.IsDefault == int8(1) || sales.AutoPaid == int8(1) {
			// buat kan finance revenue untuk sales invoice yang tealah dibuat pada auto-invoice
			if err = createSOWalkInCustomerAutoInvoice(sales, &sinvoice); err != nil {
				return nil, err
			}
		}

	}
	if sales.AutoFulfillment == 1 || partner.IsDefault == 1 {
		// Masukkan data fullfillment dengan referensi sales order id dan semua sales order item id dengan quantity yang sama
		// (untuk data yang diinput dapat melihat list inputan)
		code, err := util.CodeGen("code_fullfilment", "workorder_fulfillment")
		if err != nil {
			return nil, err
		}
		sales.Customer.Read()
		// Update fullfilment  status menjadi finished
		wofulfillment := model.WorkorderFulfillment{
			SalesOrder:      sales,
			Code:            code,
			ShippingAddress: sales.ShipmentAddress,
			DocumentStatus:  "new",
			CreatedAt:       sales.CreatedAt,
			CreatedBy:       order.Session.User,
			Priority:        "routine",
		}

		if err = wofulfillment.Save(); err != nil {
			return nil, err
		}

		for _, row := range wofulfillmentitems {
			row.WorkorderFulfillment = &wofulfillment
			if err = row.Save(); err != nil {
				return nil, err
			}
		}
		wofulfillmentID := wofulfillment.ID
		datawofulfillment, _ := getWorkorderFulfillmentByID(wofulfillmentID)
		approveFulfillment(datawofulfillment)
	}
	// update document status pada sales order menjadi active
	if sales.AutoInvoice == int8(0) && sales.AutoFulfillment == int8(0) && partner.IsDefault == int8(0) && sales.AutoPaid == int8(0) {
		sales.DocumentStatus = "new"
	}

	if err = sales.Save("DocumentStatus"); err != nil {
		return nil, err
	}

	return
}

// createSOWalkInCustomerAutoInvoice membuat finance revenue untuk sales order yang walk-in customer dan auto-invoice
func createSOWalkInCustomerAutoInvoice(so *model.SalesOrder, si *model.SalesInvoice) (e error) {
	var fr *model.FinanceRevenue
	// ubah status sales order invoice menjadi 'finished'
	so.InvoiceStatus = "finished"
	// ubah status sales invoice manjadi 'finished'
	si.DocumentStatus = "finished"

	// simpan sales order dan sales invoice
	if e = so.Save(); e == nil {
		if e = si.Save(); e == nil {
			// buat finance revenue dari sales invoice
			if fr, e = createFinanceRevenue(si); e == nil {
				// jumlahkan semua total revenue dari sales invoice
				if _, e = sumTotalRevenuedSalesInvoice(si.ID); e == nil {
					// approve-kan finance revenue yang telah dibuat
					e = approveRevenue(fr)
				}
			}
		}
	}
	return
}

// createFinanceRevenue untuk membuat finance revenue dari sales invoice
func createFinanceRevenue(si *model.SalesInvoice) (fr *model.FinanceRevenue, e error) {
	fr = &model.FinanceRevenue{
		RefID:           uint64(si.ID),
		RefType:         "sales_invoice",
		RecognitionDate: si.RecognitionDate,
		PaymentMethod:   "cash",
		Amount:          si.TotalAmount,
		DocumentStatus:  "uncleared",
		IsDeleted:       int8(0),
		CreatedAt:       time.Now(),
		CreatedBy:       si.CreatedBy,
	}
	e = fr.Save()
	return
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// sumTotalRevenuedSalesInvoice untuk menjumlahkan semua amount pada finance revenue berdasarkan sales invoice
func sumTotalRevenuedSalesInvoice(salesInvoiceID int64) (totalRevenued float64, e error) {
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

// approveRevenue untuk melakukan approve pada finance revenue
func approveRevenue(rev *model.FinanceRevenue) (e error) {
	// ubah status revenue menjadi cleared
	rev.DocumentStatus = "cleared"
	if e = rev.Save("document_status"); e == nil {
		// cek tipe refType
		if rev.RefType == "sales_invoice" {
			e = salesInvoiceRevenue(rev)
		} else {
			e = errors.New("refType is wrong")
		}
	}
	return
}

// salesInvoiceRevenue proses update total paid so invoice
func salesInvoiceRevenue(rev *model.FinanceRevenue) (e error) {
	var soInvoice *model.SalesInvoice
	if soInvoice, e = showSalesInvoice("id", rev.RefID); e == nil {
		CalculateTotalPaidSI(soInvoice)
		CheckInvoiceStatus(soInvoice.SalesOrder)
		CheckDocumentStatus(soInvoice.SalesOrder)
		partnership.CalculationTotalDebt(soInvoice.SalesOrder.Customer.ID)
		partnership.CalculationTotalSpend(soInvoice.SalesOrder.Customer.ID)
	}
	return
}

// showSalesInvoice untuk mengambil data sales invoice berdasarkan param
func showSalesInvoice(field string, values ...interface{}) (*model.SalesInvoice, error) {
	m := new(model.SalesInvoice)
	o := orm.NewOrm().QueryTable(m)
	if err := o.Filter(field, values...).Filter("is_deleted", int8(0)).RelatedSel(3).Limit(1).One(m); err != nil {

		return nil, err
	}

	return m, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////

// CancelSalesOrder untuk cancel sales order
func CancelSalesOrder(so *model.SalesOrder, user *model.User) (e error) {
	o := orm.NewOrm()
	so.IsDeleted = 0
	so.DocumentStatus = "approved_cancel"
	so.Customer.TotalSpend = so.Customer.TotalSpend - so.TotalCharge
	if e = calculateStockCommited(so); e != nil {
		return e
	}
	so.ApproveCancelAt = time.Now()
	so.ApproveCancelBy = user

	var si []*model.SalesInvoice

	if so.InvoiceStatus != "new" {
		if o.QueryTable(new(model.SalesInvoice)).Filter("sales_order_id", so.ID).All(&si); e == nil {
			for _, six := range si {
				//update is delete pada seluruh sales invoice yg memiliki referensi  SO, menjadi is_delete = 1
				six.IsDeleted = 1
				six.Save("IsDeleted")
				//check document status
				if six.DocumentStatus != "new" {

					// update total debt dengan total debt lama dikurang dengan total charge SO, pd partnership tersebut
					so.Customer.TotalDebt = so.Customer.TotalDebt - (so.TotalCharge - so.TotalPaid)
					if e = so.Customer.Save("total_debt"); e == nil {
					}

					// cek is_bundle
					if six.IsBundled != 1 {

						//update is delete pada seluruh finance revenue yg memiliki referensi sales invoice, menjadi is_delete = 1
						var fr []*model.FinanceRevenue
						if o.QueryTable(new(model.FinanceRevenue)).Filter("ref_id", six.ID).Filter("ref_type", "sales_invoice").All(&fr); e == nil {
							for _, frx := range fr {
								frx.IsDeleted = 1
								frx.Save("IsDeleted")
							}
						}
						// jika is bundled = 1
					} else {
						if e = checkInvoiceReceiptItem(six); e != nil {
							return
						}

					}

				}
			}
		}
	}

	// cek fulfillment status
	var fulfillment []*model.WorkorderFulfillment
	if so.FulfillmentStatus != "new" {
		// ambil workorder_fulfillment berdasarkan parameter sales order diatas
		if fulfillment, e = getWorkorderFulfillments("sales_order_id", so.ID); e == nil {
			for _, ffx := range fulfillment {
				// cek dokumen status work order fulfillment finished
				if ffx.DocumentStatus == "finished" {
					if e = inventory.CancelStock(uint64(ffx.ID), "workorder_fulfillment"); e != nil {
						break
					}
				}

				ffx.IsDeleted = 1
				if e = ffx.Save("IsDeleted"); e != nil {
					return
				}
			}
		}
	} else {
		//if e = calculateStockAvailable(so); e != nil {
		//	return e
		//}
	}

	if e = so.Customer.Save("total_spend"); e == nil {
		if e = so.Save("cancelled_note", "is_deleted", "document_status", "approve_cancel_at", "approve_cancel_by"); e != nil {
			return e
		}
	}

	return
}

// getSalesInvoiceByPartnerID digunakan untuk validasi create sales order
func getSalesInvoiceByPartnerID(partnership *model.Partnership) (*model.SalesInvoice, error) {
	var sinvoice model.SalesInvoice

	query := orm.NewOrm()
	err := query.Raw("SELECT si.* FROM sales_invoice si "+
		"JOIN sales_order so ON so.id = si.sales_order_id "+
		"WHERE so.customer_id = ? AND si.document_status = 'new' OR si.document_status = 'active'", partnership.ID).QueryRow(&sinvoice)

	return &sinvoice, err
}

// getItemVariantPricing mengambil item variant price
func getItemVariantPricing(PricingType *model.PricingType, ItemVariant *model.ItemVariant) (IVPrice *model.ItemVariantPrice, err error) {
	var Item model.ItemVariantPrice
	query := orm.NewOrm().QueryTable(Item)

	if PricingType.ParentType == nil {
		if err = query.Filter("item_variant_id__id", ItemVariant.ID).Filter("pricing_type_id__id", PricingType.ID).RelatedSel().Limit(1).One(&Item); err == nil {
			return &Item, err
		}
	} else {
		if err = query.Filter("item_variant_id__id", ItemVariant.ID).Filter("pricing_type_id__id", PricingType.ParentType.ID).RelatedSel().Limit(1).One(&Item); err == nil {
			return &Item, err
		}
	}
	return nil, err
}

func calculateStockAvailable(so *model.SalesOrder) (e error) {
	o := orm.NewOrm()
	o.Raw("select * from sales_order_item where sales_order_id = ?", so.ID).QueryRows(&so.SalesOrderItems)

	for _, sox := range so.SalesOrderItems {
		variant := sox.ItemVariant
		variant.Read("ID")

		//available stock + (quantity SO - quantity FFI)
		variant.AvailableStock = variant.AvailableStock + (sox.Quantity - sox.QuantityFulfillment)

		if e = variant.Save("commited_stock", "available_stock"); e != nil {
			return e
		}

	}
	return
}

func calculateStockCommited(so *model.SalesOrder) (e error) {
	o := orm.NewOrm()
	o.Raw("select * from sales_order_item where sales_order_id = ?", so.ID).QueryRows(&so.SalesOrderItems)

	for _, sox := range so.SalesOrderItems {
		variant := sox.ItemVariant
		variant.Read("ID")
		//committed stock - (quantity SO - quantity FFI)
		variant.CommitedStock = variant.CommitedStock - (sox.Quantity - sox.QuantityFulfillment)

		if e = variant.Save("commited_stock", "available_stock"); e != nil {
			return e
		}

	}
	return
}

func getWorkorderFulfillments(field string, values ...interface{}) (m []*model.WorkorderFulfillment, err error) {
	mx := new(model.WorkorderFulfillment)
	o := orm.NewOrm()
	var ff []*model.WorkorderFulfillment
	if o.QueryTable(mx).Filter(field, values...).All(&ff); err == nil {
		for _, x := range ff {
			o.LoadRelated(x, "WorkorderFulFillmentItems", 3)
		}
		return ff, nil
	}
	return nil, err
}

// UpdateSalesOrder untuk update data sales order dan sales order item
// return sales order dan error
// untuk saat ini  itemsReq adalah item2 dari inputan / request
func UpdateSalesOrder(so *model.SalesOrder, itemsReq []*model.SalesOrderItem) (*model.SalesOrder, error) {
	var e error
	oldSo := &model.SalesOrder{ID: so.ID}
	oldSo.Read()
	diff := so.TotalCharge - oldSo.TotalCharge
	//save dulu perubahan di so
	if e = so.Save(); e == nil {
		var itemsReqID []int64
		// looping yang dari req
		for _, req := range itemsReq {
			if req.ID != 0 {
				// append itemsReqID
				itemsReqID = append(itemsReqID, req.ID)
				// update soitem di database dari data yang ada request ada di database
				item := &model.SalesOrderItem{ID: req.ID}
				item.Read()
				item.ItemVariant = req.ItemVariant
				item.Quantity = req.Quantity
				item.UnitPrice = req.UnitPrice
				item.Discount = req.Discount
				item.Subtotal = req.Subtotal
				item.Note = req.Note
				item.Save("ItemVariant", "Quantity", "UnitPrice", "Discount", "Subtotal", "Note")
			} else {
				req.SalesOrder = &model.SalesOrder{ID: so.ID}
				req.Save()
			}

			// update available stock dan commited stock
			stock.CalculateAvailableStockItemVariant(req.ItemVariant)

		}

		// looping id dari database
		for _, i := range so.SalesOrderItems {
			if !util.HasElem(itemsReqID, i.ID) {
				i.Delete()
			}
		}

		// update partnership
		customer := &model.Partnership{ID: so.Customer.ID}
		customer.Read()
		customer.TotalDebt += diff
		customer.TotalSpend += diff
		customer.Save("TotalDebt", "TotalSpend")
		var emptyLoad []string
		emptyLoad = append(emptyLoad, "sales_order_items")
		so, _ = GetDetailSalesOrder(so.ID, emptyLoad)
		return so, nil
	}

	return nil, e
}

func checkInvoiceReceiptItem(si *model.SalesInvoice) (err error) {
	iri := new(model.InvoiceReceiptItem)
	o := orm.NewOrm()

	// get InvoiceReceiptItem dengan parameter SalesInvoiceID
	var irix []*model.InvoiceReceiptItem
	if o.QueryTable(iri).Filter("sales_invoice_id", si.ID).All(&irix); err == nil {
		for _, iReceiptItem := range irix {
			StIReceipt := iReceiptItem.Subtotal
			ir := iReceiptItem.InvoiceReceipt
			ir.Read("ID")

			// update total_amount di invoice receipt
			// total_amount = total_amount - subtotal(invoice receipt)
			ir.TotalAmount = ir.TotalAmount - iReceiptItem.Subtotal
			if err = ir.Save("total_amount"); err == nil {
				// hapus data invoice receipt yg memiliki sales_invoice diatas
				if err = iReceiptItem.Delete(); err == nil {
					// ambil amount pada invoice receipt yg diubah
					irx := iReceiptItem.InvoiceReceipt
					irx.Read("ID")

					// cek dokumen status invoice receipt
					if irx.DocumentStatus == "finished" {
						var fr []*model.FinanceRevenue
						if o.QueryTable(new(model.FinanceRevenue)).Filter("ref_id", si.ID).Filter("ref_type", "sales_invoice").All(&fr); err == nil {
							for _, frx := range fr {

								// update amount pada finance revenue
								// amount = amount - subtotal(invoice receipt yg diubah)
								frx.Amount = frx.Amount - StIReceipt
								if err = frx.Save("Amount"); err == nil {

									// cek amount
									if frx.Amount == 0 {
										// update is_deleted = 1, pada finance revenue
										frx.IsDeleted = 1
										if err = frx.Save(); err != nil {
											return
										}
									}
								}
							}
						}
					}
					// cek amount invoice receipt
					if irx.TotalAmount == 0 {
						//update is_deleted = 1, pada invoice receipt yg diubah
						irx.IsDeleted = int8(1)
						if err = irx.Save("IsDeleted"); err == nil {
						}
					}
				} else {
					break
				}

			}
		}
	}
	return
}

func checkFulfillmentInWorkorderShipmentItem(field string, values ...interface{}) (err error) {
	mx := new(model.WorkorderShipmentItem)
	o := orm.NewOrm()

	// get WorkorderShipmentItem dengan parameter WorkorderFulFillmentID
	var si []*model.WorkorderShipmentItem
	if o.QueryTable(mx).Filter(field, values...).All(&si); err == nil {
		for _, six := range si {
			shipmentID := six.WorkorderShipment.ID
			//hapus data WorkorderShipmentItem dengan parameter WorkorderFulFillmentID
			if err = six.Delete(); err == nil {
				// ambil keseluruhan item pada workorder shipment item dengan parameter workorder shipment id
				m := &model.WorkorderShipmentItem{WorkorderShipment: &model.WorkorderShipment{ID: shipmentID}}
				if err = m.Read("WorkorderShipment"); err != nil {
					//update is_delete = 1, pada workorder shipment
					sh := &model.WorkorderShipment{ID: shipmentID, IsDeleted: int8(1)}
					if err = sh.Save("is_deleted"); err != nil {
						return err
					}
				}

			} else {
				break
			}
		}
		return
	}
	return
}

// CodeGen will generated a code base on code prefix and last code
// e.g appSettingName: "code_sales_order"
// NOTE untuk sekarang code hanya bisa sampai 99999, mis : WO#SO-00001 s/d WO#SO-99999
func CodeGen(isWalkinCustomer bool) (code string, e error) {
	// codeIndexs for number part in code
	var codeIndex, max int
	var min, ruleDigit, appSettingName string

	if isWalkinCustomer {
		appSettingName = "code_sales_order_wic"
	} else {
		appSettingName = "code_sales_order"
	}

	// get an initial data
	if code, max, min, ruleDigit, e = util.InitCode(appSettingName); e == nil {
		// check whether last data exist or not
		if lastCode, err := lastData(isWalkinCustomer); err == nil {
			// get a number part of code in lastCode
			number := regexp.MustCompile(`[\d]+$`).FindString(lastCode)
			codeIndex = common.ToInt(number)

			// check whether a codeIndex already maximum number or not
			if codeIndex == max {
				// then change a number to initial number e.g:000001
				number = min
				code = code + number
			} else {
				// if code index is not maximum then increment a code number
				codeIndex = codeIndex + 1
				number = fmt.Sprintf("%0"+ruleDigit, codeIndex)
				code = code + number
			}
		} else {
			// generated a new code if last data not exist
			code = fmt.Sprintf(code+"%s", min)
		}
	}

	return code, e
}

func lastData(isWalkinCustomer bool) (m string, e error) {
	qb, _ := orm.NewQueryBuilder("mysql")

	qb = qb.Select("code").From("sales_order")
	if isWalkinCustomer {
		qb.Where("code like 'S#OW-%'")
	} else {
		qb.Where("code like 'S#O-%'")
	}

	o := orm.NewOrm()
	e = o.Raw(qb.OrderBy("id").Desc().String()).QueryRow(&m)

	return
}

// ShowSalesOrderFulfillment digunakan untuk mendapatkan detail sales order yang quantity sales order item nya
// belum terpenuhi semua berdasarkan quantity yang ada di workorder fulfillment item
func ShowSalesOrderFulfillment(field string, value ...interface{}) (*model.SalesOrder, error) {
	var err error
	var quantity float32
	var container []*model.SalesOrderItem
	var sales *model.SalesOrder
	sales = new(model.SalesOrder)

	o := orm.NewOrm()
	// ambil sales order
	if err = o.QueryTable(sales).Filter(field, value...).Limit(1).One(sales); err == nil {
		orm.DefaultRelsDepth = 3
		// Ambil sales order item
		_, err = o.QueryTable(new(model.SalesOrderItem)).Filter("sales_order_id__id", sales.ID).RelatedSel().All(&container)
		for _, row := range container {
			// sum quantity workorder fulfillment
			quantity, err = sumQtyFulfillment(row.ID)
			if err == nil {
				if row.Quantity > quantity {
					// QuantityPrepare digunakan untuk memberitahu berapa quantity sisa yang belum difulfill
					row.QuantityPrepare = row.Quantity - quantity
					sales.SalesOrderItems = append(sales.SalesOrderItems, row)
				}
			} else {
				return nil, err
			}
		}

		return sales, err
	}
	return nil, err
}

// sumQtyFulfillment digunakan untuk mendapatkan total quantity yang didapat dari quantity workorder fulfillment item
func sumQtyFulfillment(SOItemID int64) (float32, error) {
	var quantity float32
	var err error

	o := orm.NewOrm()

	err = o.Raw(`SELECT SUM(woi.quantity) AS quantity FROM workorder_fulfillment_item woi
	LEFT JOIN workorder_fulfillment wf ON wf.id = woi.workorder_fulfillment_id
	WHERE woi.sales_order_item_id = ? AND wf.is_deleted = ?`, SOItemID, 0).QueryRow(&quantity)

	return quantity, err
}

func approveFulfillment(fulfillment *model.WorkorderFulfillment) (*model.WorkorderFulfillment, error) {
	var e error

	fulfillment.DocumentStatus = "finished"
	if e = fulfillment.Save(); e == nil {
		if e = CheckFulfillmentStatus(fulfillment.SalesOrder); e == nil {
			var totalCost float64
			if totalCost, e = updateItemVariantStock(fulfillment); e == nil {
				// update sales order.total_cost
				fulfillment.SalesOrder.TotalCost += totalCost
				if e = fulfillment.SalesOrder.Save("total_cost"); e == nil {
					return fulfillment, nil
				}
			}
		}
	}
	return nil, e
}
func updateItemVariantStock(fulfillment *model.WorkorderFulfillment) (totalCost float64, e error) {

	for _, item := range fulfillment.WorkorderFulFillmentItems {
		var logs []*model.ItemVariantStockLog
		// ambil item sejumlah quantity dari item var stock berdasarkan FIFO
		if logs, e = inventory.FifoStockOut(item.SalesOrderItem.ItemVariant.ID, item.Quantity, "workorder_fulfillment", uint64(fulfillment.ID)); e == nil {
			// hitung cost dari logs
			for _, log := range logs {
				totalCost = totalCost + (log.ItemVariantStock.UnitCost * float64(log.Quantity))
			}
		} else {
			break
		}
		// update sales order item quantity berdasarkan workorder fulfillment item
		item.SalesOrderItem.QuantityFulfillment = item.SalesOrderItem.QuantityFulfillment + item.Quantity
		item.SalesOrderItem.Save("quantity_fulfillment")
	}

	return
}

// getWorkorderFulfillmentByID untuk get data workorder fulfillment berdasarkan id nya
// return : data workorder fulfillment dan error
func getWorkorderFulfillmentByID(id int64) (m *model.WorkorderFulfillment, err error) {
	mx := new(model.WorkorderFulfillment)
	o := orm.NewOrm()

	if err = o.QueryTable(mx).Filter("id", id).Filter("is_deleted", 0).RelatedSel().Limit(1).One(mx); err == nil {
		o.LoadRelated(mx, "WorkorderFulFillmentItems", 3)
		return mx, nil
	}
	return nil, err
}
