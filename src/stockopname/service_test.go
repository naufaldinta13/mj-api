// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stockopname

import (
	"testing"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestGetStockOpnames(t *testing.T) {
	model.DummyStockopname()
	qs := orm.RequestQuery{}
	_, _, e := GetStockOpnames(&qs)
	assert.NoError(t, e, "Data should be exists.")
}

func TestGetStockopnameByID(t *testing.T) {
	stockopname := model.DummyStockopname()

	var items []*model.StockopnameItem
	for i := 0; i <= 2; i++ {
		item := model.DummyStockopnameItem()
		item.Stockopname = stockopname
		item.ItemVariantStock = model.DummyItemVariantStock()
		item.Save("Stockopname", "ItemVariantStock")
		items = append(items, item)
	}

	m, e := GetStockopnameByID(stockopname.ID)
	assert.NoError(t, e)
	assert.Equal(t, stockopname.ID, m.ID)
	assert.Equal(t, stockopname.Code, m.Code)
	assert.Equal(t, stockopname.Note, m.Note)
	assert.NotEmpty(t, m.StockopnameItems)
	assert.Equal(t, 3, len(m.StockopnameItems))
	for i := 0; i <= 2; i++ {
		assert.Equal(t, items[i].ID, m.StockopnameItems[i].ID)
		assert.Equal(t, items[i].Note, m.StockopnameItems[i].Note)
		assert.Equal(t, items[i].Quantity, m.StockopnameItems[i].Quantity)
		assert.Equal(t, items[i].ItemVariantStock.AvailableStock, m.StockopnameItems[i].ItemVariantStock.AvailableStock)
		assert.Equal(t, items[i].ItemVariantStock.SkuCode, m.StockopnameItems[i].ItemVariantStock.SkuCode)
		assert.Equal(t, stockopname.ID, m.StockopnameItems[i].Stockopname.ID)
	}

	stockopname.Delete()
	m, e = GetStockopnameByID(stockopname.ID)
	assert.Error(t, e)
	assert.Empty(t, m)
}

func TestSaveDataStockopname(t *testing.T) {
	// stock available < quantity
	ivStock := model.DummyItemVariantStock()
	ivStock.AvailableStock = 8
	ivStock.Save("AvailableStock")

	stockItem := model.DummyStockopnameItem()
	stockItem.ItemVariantStock = ivStock
	stockItem.Quantity = 12
	stockItem.Save("ItemVariantStock", "Quantity")

	var stockopnameItems []*model.StockopnameItem
	stockopnameItems = append(stockopnameItems, stockItem)

	stockopname := model.DummyStockopname()
	stockopname.StockopnameItems = stockopnameItems

	res, e := saveDataStockopname(stockopname)
	assert.NoError(t, e)
	assert.NotEmpty(t, res)
	assert.Equal(t, 1, len(res.StockopnameItems))
	assert.Equal(t, stockItem.ID, res.StockopnameItems[0].ID)
	assert.Equal(t, stockItem.ItemVariantStock.ID, res.StockopnameItems[0].ItemVariantStock.ID)
	assert.Equal(t, stockItem.Quantity, res.StockopnameItems[0].Quantity)
	assert.Equal(t, stockItem.ItemVariantStock.ID, res.StockopnameItems[0].ItemVariantStock.ID)

	// cek item_variant_stock_log
	stockLog := &model.ItemVariantStockLog{ItemVariantStock: ivStock, LogType: "in", RefType: "stockopname"}
	stockLog.Read("ItemVariantStock", "LogType")
	assert.Equal(t, uint64(stockopname.ID), stockLog.RefID)
	assert.Equal(t, "stockopname", stockLog.RefType)
	assert.Equal(t, float32(4), stockLog.Quantity)
	assert.Equal(t, float32(12), stockLog.FinalStock)

	// stock_available > quantity
	ivStock.AvailableStock = 12
	ivStock.Save("AvailableStock")

	stockItem.Quantity = 8
	stockItem.ItemVariantStock = ivStock
	stockItem.Save("ItemVariantStock", "Quantity")

	stockopnameItems = nil
	stockopnameItems = append(stockopnameItems, stockItem)
	stockopname.StockopnameItems = stockopnameItems
	res, e = saveDataStockopname(stockopname)
	assert.NoError(t, e)
	assert.NotEmpty(t, res)
	assert.Equal(t, 1, len(res.StockopnameItems))
	assert.Equal(t, stockItem.ID, res.StockopnameItems[0].ID)
	assert.Equal(t, stockItem.ItemVariantStock.ID, res.StockopnameItems[0].ItemVariantStock.ID)
	assert.Equal(t, stockItem.Quantity, res.StockopnameItems[0].Quantity)
	assert.Equal(t, stockItem.ItemVariantStock.ID, res.StockopnameItems[0].ItemVariantStock.ID)

	// cek item_variant_stock_log
	stockLog = &model.ItemVariantStockLog{ItemVariantStock: ivStock, LogType: "in", RefType: "stockopname"}
	stockLog.Read("ItemVariantStock", "LogType")
	assert.Equal(t, uint64(stockopname.ID), stockLog.RefID)
	assert.Equal(t, "stockopname", stockLog.RefType)
	assert.Equal(t, float32(4), stockLog.Quantity)
	assert.Equal(t, float32(12), stockLog.FinalStock)
}
