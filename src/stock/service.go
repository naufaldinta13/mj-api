// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
)

// SumStockCommited for sum quantity in sales_order_item
func SumStockCommited(m *model.ItemVariant) (total float32, e error) {

	o := orm.NewOrm()

	var totalQuantitySOItem, totalQuantityFulfillmentItem float32

	// total quantity sales order item
	if e = o.Raw("select sum(quantity) from sales_order_item soi "+
		"inner join sales_order so on so.id = soi.sales_order_id where so.document_status != 'approved_cancel' "+
		"and so.is_deleted = 0 and soi.item_variant_id = ?", m.ID).QueryRow(&totalQuantitySOItem); e == nil {

		// total quantity workorder fulfillment item
		o.Raw("select sum(wfi.quantity) from workorder_fulfillment_item wfi "+
			"inner join workorder_fulfillment wf on wf.id = wfi.workorder_fulfillment_id "+
			"inner join sales_order_item soi on soi.id = wfi.sales_order_item_id "+
			"inner join sales_order so on so.id = soi.sales_order_id "+
			"where wf.is_deleted = 0 and wf.document_status = 'finished' and so.document_status != 'approved_cancel' "+
			"and so.is_deleted = 0 and soi.item_variant_id = ? ", m.ID).QueryRow(&totalQuantityFulfillmentItem)

		// total quantity sales order item - total quantity workorder fulfillment item
		total = totalQuantitySOItem - totalQuantityFulfillmentItem

		// update stock commited item variant
		mi := &model.ItemVariant{
			ID:            m.ID,
			CommitedStock: total,
		}
		e = mi.Save("commited_stock")
	}

	return
}

// sumAvailableStockItemVariantStock for sum available stock (item variant stock table)
func sumAvailableStockItemVariantStock(m *model.ItemVariant) (total float32, e error) {

	o := orm.NewOrm()

	e = o.Raw("select sum(available_stock) from item_variant_stock "+
		"where item_variant_id = ?", m.ID).QueryRow(&total)

	return
}

// CalculateAvailableStockItemVariant for calculate available stock (item variant table)
func CalculateAvailableStockItemVariant(m *model.ItemVariant) (total float32, e error) {

	var totalAvailableStock float32

	// total stock commit
	if _, e = SumStockCommited(m); e == nil {

		// total available stock (item variant stock table)
		if totalAvailableStock, e = sumAvailableStockItemVariantStock(m); e == nil {

			// update available stock item variant
			mi := &model.ItemVariant{
				ID:             m.ID,
				AvailableStock: totalAvailableStock,
			}

			e = mi.Save("available_stock")
		}

	}

	return
}

// GetStockLog for get item_stock_log by entity
func GetStockLog(rq *orm.RequestQuery) (m *[]model.ItemVariantStockLog, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.ItemVariantStockLog))

	// get total data
	if total, err = q.RelatedSel().Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.ItemVariantStockLog
	if _, err = q.RelatedSel().All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
