// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package inventory

import (
	"errors"
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/stock"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/orm"
)

// VariantStock untuk menampung data item variant stock
// yang diambil berdasarkan fifo
type VariantStock struct {
	ItemVariantStock model.ItemVariantStock
	QuantityNeed     float32
}

// FifoStockOut mengambil item variant stock dengan urutan FIFO berdasarkan parameter itemVariantID dan quantity
func FifoStockOut(itemVarID int64, quantity float32, refType string, refID uint64) (stockLog []*model.ItemVariantStockLog, e error) {
	var VSstock []model.ItemVariantStock
	itemVarStock := new(model.ItemVariantStock)
	o := orm.NewOrm().QueryTable(itemVarStock)
	if _, e = o.Filter("item_variant_id", itemVarID).Exclude("available_stock", 0).OrderBy("id").RelatedSel().All(&VSstock); e == nil {
		var varStock []*VariantStock
		if varStock, e = filterTotalItemVariantStock(VSstock, quantity); e == nil {
			stockLog, e = variantStockOut(varStock, refType, refID)
		}
	}
	return
}

// FifoStockIn untuk membuat item variant stock baru dan log nya
func FifoStockIn(itemVariant *model.ItemVariant, unitCost float64, quantity float32, refType string, refID uint64) (varStock *model.ItemVariantStock, e error) {
	sc, _ := util.GenerateCodeSKU(itemVariant.ID)
	varStock = &model.ItemVariantStock{ItemVariant: itemVariant, SkuCode: sc, AvailableStock: quantity, UnitCost: unitCost, CreatedAt: time.Now()}
	if e = varStock.Save(); e == nil {
		sLog := &model.ItemVariantStockLog{ItemVariantStock: varStock, RefID: refID, RefType: refType, LogType: "in", Quantity: quantity, FinalStock: quantity}
		if e = sLog.Save(); e == nil {
			stock.CalculateAvailableStockItemVariant(itemVariant)
			varStock.ItemVariant.Read()
		}
	}

	return
}

// CheckCancelStock untuk melakukan pengecekan sebelum melakukan cancel item variant stock
func CheckCancelStock(refID uint64, refType string) (stockLog []model.ItemVariantStockLog, res bool, e error) {
	if stockLog, e = getItemVariantLogByRef(refID, refType); e == nil {
		for _, u := range stockLog {
			if res = filterCancelStock(u, refType); res == false {
				e = errors.New("stock log tidak ditemukan")

				return
			}
		}
	}

	return
}

// CancelStock untuk melakukan cancel dan mengembalikan stock pada item variant stock
func CancelStock(refID uint64, refType string) (e error) {
	var res bool
	var stockLog []model.ItemVariantStockLog
	if stockLog, res, e = CheckCancelStock(refID, refType); e == nil {
		if res == true {
			for _, u := range stockLog {
				if u.LogType == "out" {
					_, e = SaveLog(u.ItemVariantStock, refID, refType, "in", u.Quantity)
				} else {
					_, e = SaveLog(u.ItemVariantStock, refID, refType, "out", u.Quantity)
				}
				if e == nil {
					stock.CalculateAvailableStockItemVariant(u.ItemVariantStock.ItemVariant)
				}
			}
		} else {
			e = errors.New("stock already out")
		}
	}

	return
}

// GetAllItemVariant untuk mengambil semua data item variant dari database
func GetAllItemVariant(rq *orm.RequestQuery) (m *[]model.ItemVariant, total int64, e error) {
	// make new orm query
	q, o := rq.Query(new(model.ItemVariant))

	// get total data
	if total, e = q.Filter("is_deleted", int8(0)).Count(); e == nil && total != int64(0) {
		// get data requested
		var mx []model.ItemVariant
		if _, e = q.Filter("is_deleted", int8(0)).All(&mx, rq.Fields...); e == nil {
			// if embeds
			var mxs []model.ItemVariant
			for _, u := range mx {
				o.LoadRelated(&u, "ItemVariantPrices")
				mxs = append(mxs, u)
			}

			m = &mxs
		}
	}
	return
}

// getItemVariantLogByRef untuk mengambil data item variant stock log
func getItemVariantLogByRef(refID uint64, refType string) (ItemVarLog []model.ItemVariantStockLog, e error) {
	var t int64
	o := orm.NewOrm().QueryTable(new(model.ItemVariantStockLog))
	if t, e = o.Filter("ref_id", refID).Filter("ref_type", refType).OrderBy("id").RelatedSel().All(&ItemVarLog); e != nil || t == int64(0) {
		e = errors.New("data empty")
	}
	return
}

// filterCancelStock proses mengecek satu persatu data log untuk cancel item variant stock
func filterCancelStock(IVLog model.ItemVariantStockLog, refType string) (res bool) {
	o := orm.NewOrm()
	res = false
	var sLog *model.ItemVariantStockLog
	if refType == "workorder_receiving" || refType == "direct_placement" {
		if e := o.Raw("SELECT * FROM item_variant_stock_log ivl WHERE ivl.item_variant_stock_id = ? AND ivl.log_type='out' LIMIT 1", IVLog.ItemVariantStock.ID).QueryRow(&sLog); e != nil && sLog == nil {
			res = true
		}
	} else if refType == "stockopname" {
		if IVLog.LogType == "in" {
			if e := o.Raw("SELECT * FROM item_variant_stock_log ivl WHERE ivl.item_variant_stock_id = ? AND ivl.log_type='out' LIMIT 1", IVLog.ItemVariantStock.ID).QueryRow(&sLog); e != nil && sLog == nil {
				res = true
			}
		} else {
			res = true
		}
	} else {
		res = true
	}
	return
}

// filterTotalItemVariantStock untuk memilah item variant stock yang mempunyai jumlah yang stock mencukupi
func filterTotalItemVariantStock(varStock []model.ItemVariantStock, quantity float32) (variantST []*VariantStock, e error) {
	var qty float32
	tempQuantity := quantity
	for _, u := range varStock {
		// kalau quantity sudah cukup atau lebih
		if qty >= quantity {
			break
		}
		var qn float32
		if u.AvailableStock < tempQuantity {
			qn = u.AvailableStock
		} else {
			qn = tempQuantity
		}
		tempQuantity = tempQuantity - u.AvailableStock
		res := &VariantStock{ItemVariantStock: u, QuantityNeed: qn}
		variantST = append(variantST, res)
		qty = qty + u.AvailableStock
	}
	// setelah di cek semua tapi quantity masih kurang
	if qty < quantity {
		e = errors.New("out of range")
	}
	return
}

// variantStockOut mengurangi stock pada item variant dan membuat lognya
func variantStockOut(varStock []*VariantStock, refType string, refID uint64) (itemStockLog []*model.ItemVariantStockLog, e error) {
	// loop range for every variant stock
	for _, u := range varStock {
		// create log
		var sLog *model.ItemVariantStockLog
		if sLog, e = SaveLog(&u.ItemVariantStock, refID, refType, "out", u.QuantityNeed); e != nil {
			break
		}
		itemStockLog = append(itemStockLog, sLog)
	}
	return
}

// SaveLog untuk menyimpan data item variant stock log dan item variant stock stock
func SaveLog(ivs *model.ItemVariantStock, refID uint64, refType string, logType string, quantity float32) (sLog *model.ItemVariantStockLog, e error) {
	o := orm.NewOrm()
	sLog = &model.ItemVariantStockLog{ItemVariantStock: ivs, RefID: refID, RefType: refType, LogType: logType, Quantity: quantity}
	if logType == "out" {
		sLog.FinalStock = ivs.AvailableStock - quantity
	} else {
		sLog.FinalStock = ivs.AvailableStock + quantity
	}
	if e = sLog.Save(); e == nil {
		// calculate all stock log for stock available in its sku (item variant stock)
		var in, out float32
		o.Raw("SELECT sum(quantity) FROM item_variant_stock_log WHERE item_variant_stock_id = ? AND log_type='out'", ivs.ID).QueryRow(&out)
		o.Raw("SELECT sum(quantity) FROM item_variant_stock_log WHERE item_variant_stock_id = ? AND log_type='in'", ivs.ID).QueryRow(&in)
		ivs.AvailableStock = in - out
		ivs.UpdatedAt = time.Now()
		if e = ivs.Save("available_stock", "updated_at"); e == nil {
			//perbarui available stock di variant
			stock.CalculateAvailableStockItemVariant(ivs.ItemVariant)
		}
	}
	return
}

// GetDetailItemVariantWithoutRelation untuk mengambil data item variant berdasarkan id tanpa relasi apapun
func GetDetailItemVariantWithoutRelation(field string, values ...interface{}) (*model.ItemVariant, error) {
	m := new(model.ItemVariant)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetDetailItemVariant untuk mengambil data item variant berdasarkan id
func GetDetailItemVariant(field string, values ...interface{}) (*model.ItemVariant, error) {
	m := new(model.ItemVariant)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	if _, err := o.Raw("select * from item_variant iv "+
		"inner join item_variant_stock ivs on ivs.item_variant_id = iv.id "+
		"inner join item_variant_stock_log ivsl on ivsl.item_variant_stock_id = ivs.id "+
		"where iv.id = ? ORDER BY ivsl.id DESC limit 0,25 ", m.ID).QueryRows(&m.ItemVariantStockLogs); err != nil {

		return nil, err
	}

	o.LoadRelated(m, "ItemVariantStocks", 2)
	o.LoadRelated(m, "ItemVariantPrices", 2)

	return m, nil
}

// GetDetailSalesOrderItemByItemVariant for get detail sales order item by item variant
func GetDetailSalesOrderItemByItemVariant(id int64) ([]*model.SalesOrderItem, error) {
	var m []*model.SalesOrderItem

	o := orm.NewOrm()
	_, e := o.Raw("select * from sales_order_item where item_variant_id = ?", id).QueryRows(&m)

	return m, e
}

// GetDetailPurchaseOrderByItemVariant for get detail purchase order item by item variant
func GetDetailPurchaseOrderByItemVariant(id int64) ([]*model.PurchaseOrderItem, error) {
	var m []*model.PurchaseOrderItem

	o := orm.NewOrm()
	_, e := o.Raw("select * from purchase_order_item where item_variant_id = ?", id).QueryRows(&m)

	return m, e
}

// GetDetailItem untuk mengambil data item berdasarkan param
func GetDetailItem(field string, values ...interface{}) (*model.Item, error) {
	m := new(model.Item)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter(field, values...).Filter("is_deleted", int8(0)).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	var mx []model.ItemVariant
	if _, err := o.QueryTable(new(model.ItemVariant)).Filter("item_id", m.ID).Filter("is_deleted", int8(0)).RelatedSel().All(&mx); err != nil {
		return nil, err
	}
	for _, u := range mx {
		o.LoadRelated(&u, "ItemVariantPrices", 2)
		mm := u
		m.ItemVariants = append(m.ItemVariants, &mm)
	}

	return m, nil
}

// SavingItem untuk menyimpan data item, item variant dan item variant price
func SavingItem(itm model.Item) (*model.Item, error) {
	var e error
	if e = itm.Save(); e == nil {
		for _, variant := range itm.ItemVariants {
			variant.Item = &model.Item{ID: itm.ID}
			if e = variant.Save(); e != nil {
				return nil, e
			}
			for _, price := range variant.ItemVariantPrices {
				price.ItemVariant = &model.ItemVariant{ID: variant.ID}
				if e = price.Save(); e != nil {
					return nil, e
				}
			}
		}
	}
	return &itm, e
}

// ArchivedItem untuk meng-archive item dan item variant
func ArchivedItem(itm *model.Item, session *auth.SessionData) (e error) {
	updateTime := time.Now()
	itm.IsArchived = int8(1)
	itm.UpdatedBy = &model.User{ID: session.User.ID}
	itm.UpdatedAt = updateTime
	if e = itm.Save("is_archived", "updated_by", "updated_at"); e == nil {
		for _, u := range itm.ItemVariants {
			if u.IsArchived == int8(0) {
				u.IsArchived = int8(1)
				u.UpdatedBy = &model.User{ID: session.User.ID}
				u.UpdatedAt = updateTime
				if e = u.Save("is_archived", "updated_by", "updated_at"); e != nil {
					return
				}
			}
		}
	}
	return
}

// UnarchivedItem untuk meng-unarchive item dan item variant
func UnarchivedItem(itm *model.Item, session *auth.SessionData) (e error) {
	updateTime := time.Now()
	itm.IsArchived = int8(0)
	itm.UpdatedBy = &model.User{ID: session.User.ID}
	itm.UpdatedAt = updateTime
	if e = itm.Save("is_archived", "updated_by", "updated_at"); e == nil {
		for _, u := range itm.ItemVariants {
			u.IsArchived = int8(0)
			u.UpdatedBy = &model.User{ID: session.User.ID}
			u.UpdatedAt = updateTime
			if e = u.Save("is_archived", "updated_by", "updated_at"); e != nil {
				return
			}
		}
	}
	return
}

// DeleterItem untuk men-delete item dan item variant
func DeleterItem(itm *model.Item) (e error) {
	itm.IsDeleted = int8(1)
	if e = itm.Save("is_deleted"); e == nil {
		for _, u := range itm.ItemVariants {
			u.IsDeleted = int8(1)
			if e = u.Save("is_deleted"); e != nil {
				return
			}
		}
	}
	return
}

// SavingUpdateItem untuk menyimpan data item, item variant dan item variant price update
func SavingUpdateItem(itm model.Item) (*model.Item, error) {
	var e error
	if e = itm.Save("category_id", "updated_at", "updated_by", "note", "has_variant", "item_name"); e == nil {
		for _, variant := range itm.ItemVariants {
			if variant.ID != int64(0) {
				if e = variant.Save("updated_at", "updated_by", "measurement_id", "external_name", "variant_name", "image", "base_price", "note", "minimum_stock", "has_external_name"); e != nil {
					return nil, e
				}
			} else {
				variant.Item = &model.Item{ID: itm.ID}
				if e = variant.Save(); e != nil {
					return nil, e
				}
			}

			fetchItemVariantPrice(variant.ItemVariantPrices, variant.ID)

			//for _, price := range variant.ItemVariantPrices {
			//	price.ItemVariant = &model.ItemVariant{ID: variant.ID}
			//
			//	if e = price.Save(); e != nil {
			//		return nil, e
			//	}
			//}
		}
	}
	m, err := GetDetailItem("id", itm.ID)
	return m, err
}

func fetchItemVariantPrice(newItemVariantPrice []*model.ItemVariantPrice, variantID int64) (e error) {
	var oldItemVariantPrice []*model.ItemVariantPrice

	o := orm.NewOrm()
	o.Raw("select * from item_variant_price where item_variant_id = ?", variantID).QueryRows(&oldItemVariantPrice)

	var newVariantPriceID []int64
	for _, newVariantPrice := range newItemVariantPrice {
		newVariantPrice.ItemVariant = &model.ItemVariant{ID: variantID}
		if newVariantPrice.ID != int64(0) {
			newVariantPriceID = append(newVariantPriceID, newVariantPrice.ID)
		}
		e = newVariantPrice.Save()
	}
	for _, oldVariantPrice := range oldItemVariantPrice {
		if !util.HasElem(newVariantPriceID, oldVariantPrice.ID) {
			oldVariantPrice.Delete()
		}
	}

	return e

}

func validDeleteVariant(variantID int64) bool {
	o := orm.NewOrm()

	var s int64
	o.Raw("select sum(available_stock) from item_variant_stock where item_variant_id = ?", variantID).QueryRow(&s)
	if s > 0 {
		return false
	}

	return true
}

func validDeleteItem(itemID int64) bool {
	o := orm.NewOrm()

	var s int64
	o.Raw("select sum(ivs.available_stock) from item_variant_stock ivs inner join item_variant iv on iv.id = ivs.item_variant_id where iv.item_id = ?", itemID).QueryRow(&s)
	if s > 0 {
		return false
	}

	return true
}
