// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fulfillment

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/inventory"
	"git.qasico.com/mj/api/src/sales"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/orm"
)

// getWorkorderFulfillments get all data stockopname that matched with query request parameters.
// returning slices of workorder fullfillment, total data without limit and error.
func getWorkorderFulfillments(rq *orm.RequestQuery) (m *[]model.WorkorderFulfillment, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.WorkorderFulfillment))
	q = q.Filter("is_deleted", 0)

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.WorkorderFulfillment
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
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

func getSumQuantityFullfillmentItemBySoitemID(item *model.SalesOrderItem, id int64) (total float32, e error) {
	o := orm.NewOrm()
	if e = o.Raw("select sum(quantity) as total from workorder_fulfillment_item i "+
		"inner join workorder_fulfillment f on f.id = i.workorder_fulfillment_id "+
		"where i.sales_order_item_id = ? and f.is_deleted = ? and i.id != ? ", item.ID, 0, id).QueryRow(&total); e == nil {
		return total, nil
	}
	return 0, e
}

// saveDataFulFillment untuk save data fulfillment, fulfillment_item, update status fulfillment_status sales order = "active"
func saveDataFulFillment(fulfillment *model.WorkorderFulfillment) (*model.WorkorderFulfillment, error) {
	var e error
	// save fulfillment
	if e = fulfillment.Save(); e == nil {
		for _, i := range fulfillment.WorkorderFulFillmentItems {
			i.WorkorderFulfillment = &model.WorkorderFulfillment{ID: fulfillment.ID}
			i.Save()
		}
		so := &model.SalesOrder{ID: fulfillment.SalesOrder.ID}
		so.Read()
		so.FulfillmentStatus = "active"
		so.Save("FulfillmentStatus")

		return fulfillment, nil
	}
	return nil, e
}

// updateDataFulFillment untuk update data fulfillment dan fulfillment_item
// return data fulfillment dan error
// untuk saat ini itemsfulReq adalah item2 dari inputan / request
func updateDataFulFillment(fulfillment *model.WorkorderFulfillment, itemsfulReq []*model.WorkorderFulfillmentItem) (*model.WorkorderFulfillment, error) {
	var e error
	//save dulu perubahan di fulfillment
	if e = fulfillment.Save(); e == nil {

		// looping id dari database
		var itemsID []int64
		for _, i := range fulfillment.WorkorderFulFillmentItems {
			id := i.ID
			itemsID = append(itemsID, id)
		}

		// looping yang dari req
		for _, req := range itemsfulReq {
			// update item fulfillmentItems di database
			if util.HasElem(itemsID, req.ID) {
				item := &model.WorkorderFulfillmentItem{ID: req.ID}
				e = item.Read()
				item.SalesOrderItem = req.SalesOrderItem
				item.Quantity = req.Quantity
				item.Note = req.Note
				item.Save("SalesOrderItem", "Quantity", "Note")

			} else {
				req.WorkorderFulfillment = &model.WorkorderFulfillment{ID: fulfillment.ID}
				req.Save()
			}
		}

		fulfillment, _ = getWorkorderFulfillmentByID(fulfillment.ID)
		return fulfillment, nil
	}

	return nil, e
}

func approveFulfillment(fulfillment *model.WorkorderFulfillment) (*model.WorkorderFulfillment, error) {
	var e error

	fulfillment.DocumentStatus = "finished"
	if e = fulfillment.Save(); e == nil {
		if e = sales.CheckFulfillmentStatus(fulfillment.SalesOrder); e == nil {
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

// getSumQuantityFulfillmentItemByFulfillment berguna untuk mengambil jumlah quantity dari fullfillmentItem berdasarkan fulfillment
func getSumQuantityFulfillmentItemByFulfillment(fulfillment *model.WorkorderFulfillment) (total float32, e error) {
	o := orm.NewOrm()
	if e = o.Raw("select sum(quantity) as total from workorder_fulfillment_item where workorder_fulfillment_id = ?", fulfillment.ID).QueryRow(&total); e == nil {
		return total, nil
	}
	return 0, e
}

// getSumQuantitySalesOrderItemByFulfillment berguna untuk mengambil jumlah quantity dari sales order item berdasarkan fulfillment
func getSumQuantitySalesOrderItemByFulfillment(fulfillment *model.WorkorderFulfillment) (total float32, e error) {
	o := orm.NewOrm()
	if e = o.Raw("select sum(quantity) as total from sales_order_item where sales_order_id = ?", fulfillment.SalesOrder.ID).QueryRow(&total); e == nil {
		return total, nil
	}
	return 0, e
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

func getItemVariantStockByItemVariant(itemVariant *model.ItemVariant) (ivStocks []*model.ItemVariantStock, e error) {
	o := orm.NewOrm()
	if _, e = o.Raw("select * from item_variant_stock where item_variant_id = ? and available_stock > ? order by created_at asc", itemVariant.ID, 0).QueryRows(&ivStocks); e == nil {
		return ivStocks, nil
	}
	return nil, e
}
