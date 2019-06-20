// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receiving

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/inventory"

	"git.qasico.com/cuxs/orm"
)

// TotalQuantityRI untuk menghitung total quantity semua receiving item
// yang memiliki referensi ke purchase order
func TotalQuantityRI(PoID int64) (total float32, err error) {
	o := orm.NewOrm()

	err = o.Raw("SELECT SUM(ri.quantity) from workorder_receiving_item ri "+
		"WHERE ri.purchase_order_item_id = ?;", PoID).QueryRow(&total)

	return total, err
}

// TotalQuantityPOI untuk menghitung total quantity semua purchase order item
// yang memiliki referensi ke purchase order
func TotalQuantityPOI(PoID int64) (total float32, err error) {
	o := orm.NewOrm()

	err = o.Raw("SELECT SUM(poi.quantity) from purchase_order_item poi "+
		"WHERE poi.purchase_order_id = ?;", PoID).QueryRow(&total)

	return total, err
}

// TotalAllReceivingByPO untuk menghitung total quantity semua receiving
// yang memiliki referensi ke purchase order
func TotalAllReceivingByPO(PoID int64) (total float32, err error) {
	o := orm.NewOrm()

	err = o.Raw("select sum(ri.quantity) from workorder_receiving_item ri "+
		"inner join workorder_receiving r on r.id = ri.workorder_receiving_id "+
		"inner join purchase_order_item poi on poi.id = ri.purchase_order_item_id "+
		"inner join purchase_order po on po.id = poi.purchase_order_id "+
		"where r.is_deleted = 0 AND r.document_status = 'finished' AND r.purchase_order_id = ?", PoID).QueryRow(&total)
	return
}

// CreateReceiving for create
func CreateReceiving(r *createRequest) (wr *model.WorkorderReceiving, e error) {
	wr = r.Transform()

	if e = wr.Save(); e == nil {

		for _, rix := range wr.WorkorderReceivingItems {
			rix.WorkorderReceiving = &model.WorkorderReceiving{ID: wr.ID}
			rix.Save()
			//ambil id variant
			rix.PurchaseOrderItem.Read("ID")

			if _, e = inventory.FifoStockIn(rix.PurchaseOrderItem.ItemVariant, rix.PurchaseOrderItem.UnitPrice, rix.Quantity, "workorder_receiving", uint64(wr.ID)); e != nil {
				return
			}

		}

		qPOI, _ := TotalQuantityPOI(wr.PurchaseOrder.ID)
		qAR, _ := TotalAllReceivingByPO(wr.PurchaseOrder.ID)

		if e = wr.PurchaseOrder.Read("ID"); e == nil {
			if qAR == qPOI {
				wr.PurchaseOrder.ReceivingStatus = "finished"

			} else {
				wr.PurchaseOrder.ReceivingStatus = "active"
			}

			if wr.PurchaseOrder.InvoiceStatus == "finished" && wr.PurchaseOrder.ReceivingStatus == "finished" {
				wr.PurchaseOrder.DocumentStatus = "finished"
			} else {
				wr.PurchaseOrder.DocumentStatus = "active"
			}

			e = wr.PurchaseOrder.Save()
		}

	}

	return
}

// GetReceiving get all data Receiving that matched with query request parameters.
// returning slices of WorkorderReceiving, total data without limit and error.
func GetReceiving(rq *orm.RequestQuery) (m *[]model.WorkorderReceiving, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.WorkorderReceiving))

	// get total data
	if total, err = q.Filter("is_deleted", 0).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.WorkorderReceiving
	if _, err = q.Filter("is_deleted", 0).All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetDetailReceiving for get detail data in table receiving
func GetDetailReceiving(field string, value ...interface{}) (receiving *model.WorkorderReceiving, err error) {
	woreceive := new(model.WorkorderReceiving)
	query := orm.NewOrm()
	row := query.QueryTable(woreceive)
	if err = row.Filter(field, value).Filter("is_deleted", 0).RelatedSel().Limit(1).One(woreceive); err == nil {
		query.LoadRelated(woreceive, "WorkorderReceivingItems", 3)
		return woreceive, err
	}
	return nil, err
}
