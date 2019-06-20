// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stockopname

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/inventory"

	"git.qasico.com/cuxs/orm"
)

// GetStockOpnames get all data stockopname that matched with query request parameters.
// returning slices of stockopname, total data without limit and error.
func GetStockOpnames(rq *orm.RequestQuery) (m *[]model.Stockopname, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.Stockopname))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.Stockopname
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetStockopnameByID untuk get data stockopname berdasarkan id nya
// return : data stockopname dan error
func GetStockopnameByID(id int64) (m *model.Stockopname, err error) {
	mx := new(model.Stockopname)
	o := orm.NewOrm()

	if err = o.QueryTable(mx).Filter("id", id).RelatedSel().Limit(1).One(mx); err == nil {
		o.LoadRelated(mx, "StockopnameItems", 3)
		return mx, nil
	}
	return nil, err
}

// saveDataStockopname untuk save data stockopname,stockopname_item, item_variant_stock,item_variant_stock_log, item_variant
func saveDataStockopname(stockopname *model.Stockopname) (*model.Stockopname, error) {
	var e error
	// save stockopname
	if e = stockopname.Save(); e == nil {
		for _, i := range stockopname.StockopnameItems {
			i.Stockopname = &model.Stockopname{ID: stockopname.ID}
			if e = i.Save(); e == nil {
				ivStock := &model.ItemVariantStock{ID: i.ItemVariantStock.ID}
				ivStock.Read()
				if i.Quantity > ivStock.AvailableStock {
					inventory.SaveLog(ivStock, uint64(stockopname.ID), "stockopname", "in", i.Quantity-ivStock.AvailableStock)
				} else {
					inventory.SaveLog(ivStock, uint64(stockopname.ID), "stockopname", "out", ivStock.AvailableStock-i.Quantity)
				}
			}
		}
		return stockopname, nil
	}
	return nil, e
}
