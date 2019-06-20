// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchase

import (
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/inventory"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/orm"
)

// CalculateTotalPaidPO Core library untuk menghitung total paid purchase order
func calculateTotalPaidPO(PurchaseInvoice *model.PurchaseInvoice) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("UPDATE purchase_order po SET po.total_paid = (SELECT SUM(pi.total_paid) AS total_paid from purchase_invoice pi "+
		"WHERE pi.purchase_order_id = ?) WHERE po.id = ?;", PurchaseInvoice.PurchaseOrder.ID, PurchaseInvoice.PurchaseOrder.ID).Exec()
	return err
}

// CalculateTotalPaidPI Core library untuk menghitung total paid purchase invoice
func CalculateTotalPaidPI(PurchaseInvoice *model.PurchaseInvoice) (pi *model.PurchaseInvoice, err error) {
	o := orm.NewOrm()
	_, err = o.Raw("UPDATE purchase_invoice pi SET pi.total_paid = (SELECT SUM(fe.amount) AS total_paid from finance_expense fe "+
		"WHERE fe.ref_id = ? AND fe.ref_type = 'purchase_invoice' AND fe.document_status = 'cleared') WHERE pi.id = ?;", uint64(PurchaseInvoice.ID), PurchaseInvoice.ID).Exec()
	calculateTotalPaidPO(PurchaseInvoice)
	PurchaseInvoice.Read()
	pi = PurchaseInvoice
	return
}

// ChangeDocumentStatus Core Library untuk merubah DocumentStatus InvoiceStatus dan ReceivingStatus
func ChangeDocumentStatus(order *model.PurchaseOrder) (err error) {
	// Ubah Invoice Status
	changeInvoiceStatus(order)
	// Ubah Receiving Status
	changeReceivingStatus(order)

	// Ubah Document Status
	if order.InvoiceStatus == "finished" && order.ReceivingStatus == "finished" {
		order.DocumentStatus = "finished"
		order.Save()
	}

	return
}

// changeInvoiceStatus fungsi untuk merubah Status Invoice
func changeInvoiceStatus(order *model.PurchaseOrder) (err error) {
	o := orm.NewOrm()
	var Invoice model.PurchaseInvoice
	err = o.Raw("SELECT pi.* FROM purchase_invoice pi JOIN purchase_order po ON pi.purchase_order_id = po.id WHERE pi.purchase_order_id = ? AND pi.is_deleted = ? AND po.document_status != ? AND po.is_deleted = ?;", order.ID, 0, "cancelled", 0).QueryRow(&Invoice)

	var pi *model.PurchaseInvoice
	if err == nil {
		pi, err = CalculateTotalPaidPI(&Invoice)
		if err == nil {
			order.Read()
			if order.TotalPaid == order.TotalCharge {
				order.InvoiceStatus = "finished"
			}
			err = order.Save()

			if pi.TotalPaid == pi.TotalAmount {
				pi.DocumentStatus = "finished"
				pi.Save()
			}
		}
	}
	return err
}

// changeReceivingStatus fungsi untuk merubah Receiving Status
func changeReceivingStatus(order *model.PurchaseOrder) (err error) {
	var receivings []model.WorkorderReceivingItem
	o := orm.NewOrm().QueryTable(new(model.WorkorderReceivingItem))
	o.Filter("workorder_receiving_id__purchase_order_id__id", order.ID).Filter("workorder_receiving_id__is_deleted", 0).RelatedSel().All(&receivings)
	var totalA float32
	for _, v := range receivings {
		totalA += v.Quantity
	}

	var orderitems []model.PurchaseOrderItem
	o = orm.NewOrm().QueryTable(new(model.PurchaseOrderItem))
	_, err = o.Filter("purchase_order_id__id", order.ID).Exclude("purchase_order_id__document_status__in", "cancelled").Filter("purchase_order_id__is_deleted", 0).RelatedSel().All(&orderitems)
	var totalB float32
	if len(orderitems) > 0 {
		for _, v := range orderitems {
			totalB += v.Quantity
		}

		if totalA == totalB {
			order.ReceivingStatus = "finished"
			order.Save()
		}
	}

	return
}

// ReturningPurchaseOrder Core Library PO untuk kembalikan barang
func ReturningPurchaseOrder(order *model.PurchaseOrder, purchaseRetID int64) (orderitems *[]model.PurchaseOrderItem) {
	var orderitemx []model.PurchaseOrderItem
	var wor, pri float32
	o := orm.NewOrm()
	o.QueryTable(new(model.PurchaseOrderItem)).Filter("purchase_order_id__id", order.ID).RelatedSel().All(&orderitemx)

	for i := 0; i < len(orderitemx); i++ {

		o.Raw("SELECT SUM(woi.quantity) AS qty FROM workorder_receiving_item woi "+
			"JOIN workorder_receiving wo ON wo.id = woi.workorder_receiving_id "+
			"WHERE woi.workorder_receiving_id = wo.id AND wo.purchase_order_id = ? AND woi.purchase_order_item_id = ? AND wo.is_deleted = 0 AND wo.document_status = 'finished';", order.ID, orderitemx[i].ID).QueryRow(&wor)

		o.Raw("SELECT SUM(pri.quantity) AS qty "+
			"FROM purchase_return_item pri JOIN purchase_return pr ON pr.id = pri.purchase_return_id "+
			"WHERE pri.purchase_return_id = pr.id AND pr.purchase_order_id = ? AND pri.purchase_order_item_id = ? AND pr.is_deleted = 0 AND pr.document_status != 'cancelled' AND NOT pr.id = ?;", order.ID, orderitemx[i].ID, purchaseRetID).QueryRow(&pri)

		if pri <= 0 {
			orderitemx[i].CanBeReturn += wor
		} else {
			orderitemx[i].CanBeReturn += (wor - pri)
		}
	}

	orderitems = &orderitemx

	return
}

// GetAllPurchaseOrder untuk mengambil semua data purchase order dari database
func GetAllPurchaseOrder(rq *orm.RequestQuery) (m *[]model.PurchaseOrder, total int64, e error) {
	// make new orm query
	q, _ := rq.Query(new(model.PurchaseOrder))
	// get total data
	if total, e = q.Filter("is_deleted", int8(0)).Count(); e == nil && total != int64(0) {
		// get data requested
		var mx []model.PurchaseOrder
		if _, e = q.Filter("is_deleted", int8(0)).All(&mx, rq.Fields...); e == nil {
			m = &mx
		}
	}
	return
}

// GetDetailPurchaseOrder untuk mengambil data detail purchase order berdasarkan param
func GetDetailPurchaseOrder(field string, values ...interface{}) (*model.PurchaseOrder, error) {
	m := new(model.PurchaseOrder)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter(field, values...).Filter("is_deleted", int8(0)).RelatedSel(1).Limit(1).One(m); err != nil {
		return nil, err
	}
	o.LoadRelated(m, "PurchaseOrderItems", 3)
	o.LoadRelated(m, "PurchaseInvoices", 3)
	o.LoadRelated(m, "PurchaseReturns", 3)
	o.LoadRelated(m, "WorkorderReceivings", 3)

	return m, nil
}

// CreatePurchaseOrder berguna untuk menyimpan data Purchase Order
func CreatePurchaseOrder(po *model.PurchaseOrder) (*model.PurchaseOrder, error) {
	var e error

	if e = po.Save(); e == nil {
		for _, item := range po.PurchaseOrderItems {
			item.PurchaseOrder = &model.PurchaseOrder{ID: po.ID}
			if e = item.Save(); e == nil {
				// update total credit
				supplier := &model.Partnership{ID: po.Supplier.ID}
				if e = supplier.Read(); e == nil {
					supplier.TotalCredit += po.TotalCharge
					supplier.TotalExpenditure += po.TotalCharge
					supplier.Save("TotalCredit", "TotalExpenditure")
				}
			}
		}
		//check is auto invoice
		if po.AutoInvoiced == 1 {
			// create purchase invoice
			code, _ := util.CodeGen("code_purchase_invoice", "purchase_invoice")
			pi := &model.PurchaseInvoice{
				PurchaseOrder:  po,
				Code:           code,
				TotalAmount:    po.TotalCharge,
				DocumentStatus: "new",
				CreatedBy:      po.CreatedBy,
				CreatedAt:      time.Now(),
			}
			if e = pi.Save(); e == nil {
				return po, e
			}
		} else {
			return po, e
		}
	}
	return nil, e
}

// UpdatePurchaseOrder untuk update data purchase order dan purchase order item
// return purchase order dan error
// untuk saat ini  itemsReq adalah item2 dari inputan / request
func UpdatePurchaseOrder(po *model.PurchaseOrder, itemsReq []*model.PurchaseOrderItem) (*model.PurchaseOrder, error) {
	var e error
	oldPo := &model.PurchaseOrder{ID: po.ID}
	oldPo.Read()
	//save dulu perubahan di po
	if e = po.Save(); e == nil {
		var itemsReqID []int64
		// looping yang dari req
		for _, req := range itemsReq {
			if req.ID != 0 {
				// append itemsReqID
				itemsReqID = append(itemsReqID, req.ID)
				// update soitem di database dari data yang ada request ada di database
				item := &model.PurchaseOrderItem{ID: req.ID}
				item.Read()
				item.ItemVariant = req.ItemVariant
				item.Quantity = req.Quantity
				item.UnitPrice = req.UnitPrice
				item.Discount = req.Discount
				item.Subtotal = req.Subtotal
				item.Note = req.Note
				item.Save("ItemVariant", "Quantity", "UnitPrice", "Discount", "Subtotal", "Note")
			} else {
				req.PurchaseOrder = &model.PurchaseOrder{ID: po.ID}
				req.Save()
			}
		}

		// looping id dari database
		for _, i := range po.PurchaseOrderItems {
			if !util.HasElem(itemsReqID, i.ID) {
				i.Delete()
			}
		}

		// update partnership
		po.Supplier.Read()
		po.Supplier.TotalExpenditure -= oldPo.TotalCharge
		po.Supplier.TotalExpenditure += po.TotalCharge
		po.Supplier.Save("TotalExpenditure")

		po, _ = GetDetailPurchaseOrder("id", po.ID)
		return po, nil
	}

	return nil, e
}

// getItemVariantStockByPurchaseOrder untuk get data item variant stock
// return []*item variant stock dan error
func getItemVariantStockByPurchaseOrder(po *model.PurchaseOrder) ([]*model.ItemVariantStock, error) {
	o := orm.NewOrm()

	var ivs []*model.ItemVariantStock
	_, e := o.Raw("select ivs.* from item_variant_stock ivs inner join item_variant_stock_log ivslog on ivslog.item_variant_stock_id = ivs.id "+
		"inner join workorder_receiving wr on wr.id = ivslog.ref_id inner join purchase_order po on po.id = wr.purchase_order_id "+
		"where wr.purchase_order_id = ? and ivslog.ref_type = ? and ivslog.log_type = ?", po.ID, "workorder_receiving", "in").QueryRows(&ivs)
	if e == nil {
		return ivs, nil
	}
	return nil, e
}

// cancelPurchaseOrder -> cancel po
func cancelPurchaseOrder(po *model.PurchaseOrder) (*model.PurchaseOrder, error) {
	var e error

	// update
	if e = po.Save("DocumentStatus", "CancelledNote"); e == nil {
		po.Supplier.TotalExpenditure = po.Supplier.TotalExpenditure - po.TotalCharge
		po.Supplier.Save("TotalExpenditure")

		if po.ReceivingStatus != "new" {
			// workorder receiving
			var receiving []*model.WorkorderReceiving
			if receiving, e = getWorkorderReceivingByPurchaseOrder(po); e == nil {
				for _, i := range receiving {
					if e = inventory.CancelStock(uint64(i.ID), "workorder_receiving"); e != nil {
						break
					}
					i.IsDeleted = 1
					i.Save("IsDeleted")
				}
			}
		}

		if po.InvoiceStatus != "new" {
			// update total credit supplier
			po.Supplier.TotalCredit = po.Supplier.TotalCredit - (po.TotalCharge - po.TotalPaid)
			po.Supplier.Save("TotalCredit")

			// update purchase invoice
			var purchaseInvoices []*model.PurchaseInvoice
			if purchaseInvoices, e = getPurchaseInvoiceByPurchaseOrder(po); e == nil {
				for _, i := range purchaseInvoices {
					i.IsDeleted = 1
					i.Save("IsDeleted")
					if i.DocumentStatus != "new" {
						// update finance
						var finances []*model.FinanceExpense
						if finances, e = getFinanceExpenceByPurchaseInvoice(i); e == nil {
							for _, fx := range finances {
								fx.IsDeleted = 1
								fx.Save("IsDeleted")
							}
						}
					}
				}
			}
		}
		return po, nil
	}
	return nil, e
}

// getPurchaseInvoiceByPurchaseOrder -> get data purchase invoice berdasarkan purchase order
func getPurchaseInvoiceByPurchaseOrder(po *model.PurchaseOrder) (purchaseInvoices []*model.PurchaseInvoice, e error) {
	o := orm.NewOrm()

	if _, e = o.Raw("select * from purchase_invoice where purchase_order_id = ? and is_deleted = ?", po.ID, 0).QueryRows(&purchaseInvoices); e == nil {
		return purchaseInvoices, nil
	}
	return nil, e
}

// getFinanceExpenceByPurchaseInvoice -> get data finance expense berdasarkan purchase order
func getFinanceExpenceByPurchaseInvoice(pi *model.PurchaseInvoice) (finances []*model.FinanceExpense, e error) {
	o := orm.NewOrm()

	if _, e = o.Raw("select * from finance_expense where ref_id = ? and "+
		"ref_type = 'purchase_invoice' and is_deleted = ?", pi.ID, 0).QueryRows(&finances); e == nil {
		return finances, nil
	}
	return nil, e
}

// getWorkorderReceivingByPurchaseOrder -> get data workorder receiving berdasarkan purchase order
func getWorkorderReceivingByPurchaseOrder(po *model.PurchaseOrder) (receiving []*model.WorkorderReceiving, e error) {
	o := orm.NewOrm()

	if _, e = o.Raw("select * from workorder_receiving where purchase_order_id = ? and is_deleted = ?", po.ID, 0).QueryRows(&receiving); e == nil {
		return receiving, nil
	}
	return nil, e
}

// getItemVariantStockLogByReceiving -> get data workorder receiving berdasarkan purchase order
func getItemVariantStockLogByReceiving(receiving *model.WorkorderReceiving) (ivsLog []*model.ItemVariantStockLog, e error) {
	m := new(model.ItemVariantStockLog)
	o := orm.NewOrm()
	if _, e = o.QueryTable(m).Filter("ref_id", uint64(receiving.ID)).RelatedSel().All(&ivsLog); e == nil {
		return ivsLog, nil
	}
	return nil, e
}
