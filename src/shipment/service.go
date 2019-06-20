// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package shipment

import (
	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/mj/api/datastore/model"
)

// GetShipments fungsi untuk get seluruh Work Order Shipment
func GetShipments(rq *orm.RequestQuery) (shipments []*model.WorkorderShipment, total int64, err error) {

	query, _ := rq.Query(new(model.WorkorderShipment))
	query = query.Filter("is_deleted", 0)

	if total, err = query.Count(); err != nil || total == int64(0) {
		return nil, total, err
	}

	var shipmentx []*model.WorkorderShipment
	if _, err = query.All(&shipmentx, rq.Fields...); err == nil {
		return shipmentx, total, nil
	}

	return nil, total, err
}

// GetDetailShipment fungsi untuk get berdasarkan id Work Order Shipment
func GetDetailShipment(field string, value ...interface{}) (shipment *model.WorkorderShipment, err error) {
	models := new(model.WorkorderShipment)

	o := orm.NewOrm()
	if err = o.QueryTable(models).Filter(field, value).Filter("is_deleted", 0).RelatedSel().Limit(1).One(models); err != nil {
		return nil, err
	}

	o.LoadRelated(models, "WorkorderShipmentItems", 3)

	return models, err
}

// ShowShipment fungsi untuk get berdasarkan id Work Order Shipment
func ShowShipment(field string, value ...interface{}) (shipment *model.WorkorderShipment, err error) {
	ws := new(model.WorkorderShipment)

	o := orm.NewOrm()
	if err = o.QueryTable(ws).Filter(field, value).Filter("is_deleted", 0).RelatedSel().Limit(1).One(ws); err != nil {
		return nil, err
	}

	o.LoadRelated(ws, "WorkorderShipmentItems", 3)

	soids := []string{}
	for _, mx := range ws.WorkorderShipmentItems {
		soids = append(soids, common.ToString(mx.WorkorderFulfillment.SalesOrder.ID))
	}

	if len(soids) > 0 {
		qs := o.QueryTable("workorder_fulfillment_item")
		qs.Filter("sales_order_item_id__sales_order_id__id__in", soids).RelatedSel().All(&ws.WorkorderFulfillmentItems)
	}

	return ws, err
}

// CreateWorkOrderShipment untuk menyimpan data work order shipment ke database
func CreateWorkOrderShipment(woShipment *model.WorkorderShipment) (e error) {
	if e = woShipment.Save(); e == nil {
		for _, u := range woShipment.WorkorderShipmentItems {
			u.WorkorderShipment = &model.WorkorderShipment{ID: woShipment.ID}
			if e = u.Save(); e != nil {
				return
			}
		}
	}
	return
}

func cancelShipmentAndFulfillment(shipment *model.WorkorderShipment, fulfillID int64) (ship *model.WorkorderShipment, err error) {
	for _, row := range shipment.WorkorderShipmentItems {
		if row.WorkorderFulfillment.ID == fulfillID {
			// Update is_delivered pada fulfillment yang dipilih menjadi 0
			row.WorkorderFulfillment.IsDelivered = int8(0)
			if err = row.WorkorderFulfillment.Save("IsDelivered"); err != nil {
				return nil, err
			}
		}

		soShipment := model.SalesOrder{ID: row.WorkorderFulfillment.SalesOrder.ID}
		soShipment.Read("ID")

		// update shipment status dari sales order dengan referensi fulfillment yang ada menjadi active
		soShipment.ShipmentStatus = "active"
		soShipment.DocumentStatus = "active"
		if err = soShipment.Save("ShipmentStatus", "DocumentStatus"); err != nil {
			return nil, err
		}

		// hapus fulfillment id yang dipilih di workorder shipment item
		if row.WorkorderFulfillment.ID == fulfillID {
			if err = row.Delete(); err != nil {
				return nil, err
			}
		}
	}
	// update is_delete pada workorder shipment tersebut menjadi 1
	shipmentItem := new(model.WorkorderShipmentItem)
	shipmentItem.WorkorderShipment = shipment
	if err = shipmentItem.Read("WorkorderShipment"); err != nil {
		shipment.IsDeleted = int8(1)
		// simpan update is_delete pada workorder shipment tersebut menjadi 1
		if err = shipment.Save("IsDeleted"); err != nil {
			return nil, err
		}
	}

	return shipment, err
}

func approveShipment(shipment *model.WorkorderShipment) (woShipment *model.WorkorderShipment, err error) {
	var qtyFulfill, qtySales float32
	for _, row := range shipment.WorkorderShipmentItems {
		row.WorkorderFulfillment.IsDelivered = int8(1)
		if err = row.WorkorderFulfillment.Save("IsDelivered"); err != nil {
			return nil, err
		}
		// Jumlahkan semua quantity dari fulfillment yang is_delivered nya adalah 1
		orm.NewOrm().Raw("SELECT SUM(wfi.quantity) AS qtyFulfill FROM workorder_fulfillment wf "+
			"LEFT OUTER JOIN workorder_fulfillment_item wfi ON wfi.workorder_fulfillment_id = wf.id "+
			"WHERE wf.id = ? AND wf.is_delivered = ?;", row.WorkorderFulfillment.ID, 1).QueryRow(&qtyFulfill)
		// ambil semua quantity dari sales order dari fulfillment yang dilakukan
		orm.NewOrm().Raw("SELECT SUM(soi.quantity) AS qtySales FROM workorder_fulfillment wf "+
			"JOIN sales_order so ON so.id = wf.sales_order_id "+
			"JOIN sales_order_item soi ON soi.sales_order_id = so.id "+
			"WHERE wf.id = ?;", row.WorkorderFulfillment.ID).QueryRow(&qtySales)

		soShipment := model.SalesOrder{ID: row.WorkorderFulfillment.SalesOrder.ID}
		if qtySales == qtyFulfill {
			// update shipment status dari sales order dengan referensi fulfillment menjadi finished
			if err = soShipment.Read("ID"); err == nil {
				soShipment.ShipmentStatus = "finished"
			}
			if err = soShipment.Save("ShipmentStatus"); err != nil {
				return nil, err
			}
		} else {
			if err = soShipment.Read("ID"); err == nil {
				soShipment.ShipmentStatus = "active"
			}
			if err = soShipment.Save("ShipmentStatus"); err != nil {
				return nil, err
			}
		}
		soShipment.Read()
		if soShipment.ShipmentStatus == "finished" && soShipment.InvoiceStatus == "finished" && soShipment.FulfillmentStatus == "finished" {
			soShipment.DocumentStatus = "finished"
			if err = soShipment.Save("DocumentStatus"); err != nil {
				return nil, err
			}
		} else {
			soShipment.DocumentStatus = "active"
			if err = soShipment.Save("DocumentStatus"); err != nil {
				return nil, err
			}
		}
	}

	// update document status menjadi finished di shipment yang dilakukan
	shipment.DocumentStatus = "finished"
	if err = shipment.Save("DocumentStatus"); err != nil {
		return nil, err
	}

	shipment.Read()

	return shipment, err
}
