// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fulfillment

import (
	"testing"
	"time"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestGetWorkorderFulfillments(t *testing.T) {
	fulfillment := model.DummyWorkorderFulfillment()
	fulfillment.IsDeleted = 0
	fulfillment.Save("IsDeleted")

	qs := orm.RequestQuery{}
	m, _, e := getWorkorderFulfillments(&qs)
	assert.NoError(t, e, "Data should be exists.")
	for _, x := range *m {
		assert.Equal(t, int8(0), x.IsDeleted)
	}
}

func TestGetWorkorderFulfillmentByID(t *testing.T) {
	fulfillment := model.DummyWorkorderFulfillment()
	fulfillment.IsDeleted = 0
	fulfillment.Save("IsDeleted")

	item := model.DummyWorkorderFulfillmentItem()
	item.WorkorderFulfillment = fulfillment
	item.Save("WorkorderFulfillment")

	soItem := model.DummySalesOrderItem()
	soItem.SalesOrder = fulfillment.SalesOrder
	soItem.Save("SalesOrder")

	var soItems []*model.SalesOrderItem
	soItems = append(soItems, soItem)

	fulfillment.SalesOrder.SalesOrderItems = soItems
	fulfillment.SalesOrder.Save("SalesOrderItem")

	data, e := getWorkorderFulfillmentByID(fulfillment.ID)
	assert.NoError(t, e)
	assert.Equal(t, int8(0), data.IsDeleted)
	assert.Equal(t, fulfillment.Note, data.Note)
	assert.Equal(t, fulfillment.Code, data.Code)
	assert.Equal(t, fulfillment.SalesOrder.ID, data.SalesOrder.ID)
	assert.Equal(t, fulfillment.DocumentStatus, data.DocumentStatus)
	assert.NotEmpty(t, fulfillment.SalesOrder.SalesOrderItems)
	assert.NotEmpty(t, fulfillment.SalesOrder.SalesOrderItems[0].Note)
	assert.NotEmpty(t, fulfillment.SalesOrder.SalesOrderItems[0].ItemVariant.ID)
	assert.NotEmpty(t, fulfillment.SalesOrder.SalesOrderItems[0].ItemVariant.Note)

	fulfillment.IsDeleted = 1
	fulfillment.Save("IsDeleted")
	data, e = getWorkorderFulfillmentByID(fulfillment.ID)
	assert.Error(t, e)
	assert.Empty(t, data)
}

func TestGetSumQuantityFullfillmentItemBySoitem(t *testing.T) {
	soitem := model.DummySalesOrderItem()
	ful := model.DummyWorkorderFulfillment()
	ful.IsDeleted = 0
	ful.Save("IsDeleted")

	fulfillItem := model.DummyWorkorderFulfillmentItem()
	fulfillItem.SalesOrderItem = soitem
	fulfillItem.Quantity = 2
	fulfillItem.WorkorderFulfillment = ful
	fulfillItem.Save("SalesOrderItem", "Quantity", "WorkorderFulfillment")

	fulfillItem2 := model.DummyWorkorderFulfillmentItem()
	fulfillItem2.SalesOrderItem = soitem
	fulfillItem2.Quantity = 2
	fulfillItem2.WorkorderFulfillment = ful
	fulfillItem2.Save("SalesOrderItem", "Quantity", "WorkorderFulfillment")

	tot, e := getSumQuantityFullfillmentItemBySoitemID(soitem, 0)
	assert.Equal(t, float32(4), tot)
	assert.NoError(t, e)

	tot, e = getSumQuantityFullfillmentItemBySoitemID(soitem, fulfillItem.ID)
	assert.Equal(t, float32(2), tot)
	assert.NoError(t, e)
}

func TestSaveDataFulfillment(t *testing.T) {
	fulItem := model.DummyWorkorderFulfillmentItem()

	res, e := saveDataFulFillment(fulItem.WorkorderFulfillment)
	assert.NoError(t, e)
	assert.NotEmpty(t, res)
}

func TestUpdateDataFulfillment(t *testing.T) {
	fulfillment := model.DummyWorkorderFulfillment()

	var items []*model.WorkorderFulfillmentItem
	item := model.DummyWorkorderFulfillmentItem()
	item.Quantity = 5
	item.WorkorderFulfillment = fulfillment
	item.Save("Quantity", "WorkorderFulfillment")

	item2 := model.DummyWorkorderFulfillmentItem()
	item2.Quantity = 8
	item2.WorkorderFulfillment = fulfillment
	item2.Save("Quantity", "WorkorderFulfillment")
	items = append(items, item, item2)

	fulfillment.WorkorderFulFillmentItems = items
	fulfillment.IsDeleted = 0
	fulfillment.Save("WorkorderFulFillmentItems", "IsDeleted")

	// item yang baru
	var itemsNew []*model.WorkorderFulfillmentItem
	itemNew := &model.WorkorderFulfillmentItem{ID: item.ID, WorkorderFulfillment: fulfillment, SalesOrderItem: model.DummySalesOrderItem(), Quantity: 1, Note: "abc"}
	itemNew2 := &model.WorkorderFulfillmentItem{WorkorderFulfillment: fulfillment, SalesOrderItem: model.DummySalesOrderItem(), Quantity: 3, Note: "xxxx"}
	itemsNew = append(itemsNew, itemNew, itemNew2)

	// update data yang di fulfillment.WorkorderFulFillmentItems
	res, e := updateDataFulFillment(fulfillment, itemsNew)
	assert.NoError(t, e)
	assert.NotEmpty(t, res)
	assert.Equal(t, 3, len(res.WorkorderFulFillmentItems))

	data, _ := getWorkorderFulfillmentByID(fulfillment.ID)
	assert.Equal(t, 3, len(data.WorkorderFulFillmentItems))
}

func TestGetSumQuantityFulfillmentItemByFulfillment(t *testing.T) {
	fulfillment := model.DummyWorkorderFulfillment()
	for a := 0; a < 2; a++ {
		item := model.DummyWorkorderFulfillmentItem()
		item.WorkorderFulfillment = fulfillment
		item.Quantity = 2
		item.Save("WorkorderFulfillment", "Quantity")
	}

	tot, e := getSumQuantityFulfillmentItemByFulfillment(fulfillment)
	assert.NoError(t, e)
	assert.Equal(t, float32(4), tot)
}

func TestGetSumQuantitySalesOrderItemByFulfillment(t *testing.T) {
	fulfillment := model.DummyWorkorderFulfillment()
	for a := 0; a < 2; a++ {
		item := model.DummySalesOrderItem()
		item.SalesOrder = fulfillment.SalesOrder
		item.Quantity = 2
		item.Save("SalesOrder", "Quantity")
	}

	tot, e := getSumQuantitySalesOrderItemByFulfillment(fulfillment)
	assert.NoError(t, e)
	assert.Equal(t, float32(4), tot)
}

func TestGetItemVariantStockByItemVariant(t *testing.T) {
	iv := model.DummyItemVariant()

	for a := 0; a < 2; a++ {
		ivStock := model.DummyItemVariantStock()
		ivStock.ItemVariant = iv
		ivStock.CreatedAt = time.Now()
		ivStock.AvailableStock = 3
		ivStock.Save()
	}

	ivStocks, e := getItemVariantStockByItemVariant(iv)
	assert.NoError(t, e)
	assert.Equal(t, 2, len(ivStocks))
}

func TestUpdateItemVariantAndStock(t *testing.T) {
	fulfillment := model.DummyWorkorderFulfillment()

	iv := model.DummyItemVariant()
	iv.CommitedStock = 20
	iv.Save("CommitedStock")

	soitem := model.DummySalesOrderItem()
	soitem.ItemVariant = iv
	soitem.Save("ItemVariant")

	item1 := model.DummyWorkorderFulfillmentItem()
	item1.Quantity = 10
	item1.SalesOrderItem = soitem
	item1.WorkorderFulfillment = fulfillment
	item1.Save("Quantity", "WorkorderFulfillment", "SalesOrderItem")

	ivStock1 := model.DummyItemVariantStock()
	ivStock1.ItemVariant = iv
	ivStock1.AvailableStock = 12
	ivStock1.CreatedAt = time.Now()
	ivStock1.UnitCost = 1000
	ivStock1.Save("ItemVariant", "UnitCost", "CreatedAt", "AvailableStock")

	// dummy ke 2
	iv2 := model.DummyItemVariant()
	iv2.CommitedStock = 20
	iv2.Save("CommitedStock")

	soItem2 := model.DummySalesOrderItem()
	soItem2.ItemVariant = iv2
	soItem2.Save("ItemVariant")

	item2 := model.DummyWorkorderFulfillmentItem()
	item2.Quantity = 10
	item2.SalesOrderItem = soItem2
	item2.WorkorderFulfillment = fulfillment
	item2.Save("Quantity", "WorkorderFulfillment", "SalesOrderItem")

	for a := 0; a < 3; a++ {
		ivStock2 := model.DummyItemVariantStock()
		ivStock2.ItemVariant = iv2
		ivStock2.AvailableStock = 5
		ivStock2.UnitCost = 1000
		ivStock2.Save("ItemVariant", "AvailableStock", "UnitCost")
	}

	var items []*model.WorkorderFulfillmentItem
	items = append(items, item1, item2)

	fulfillment.WorkorderFulFillmentItems = items
	fulfillment.Save("WorkorderFulFillmentItems")

	totCost, e := updateItemVariantStock(fulfillment)
	assert.NoError(t, e)
	assert.Equal(t, float64(20000), totCost)

}

func TestApproveFulfillment(t *testing.T) {
	fulfillment := model.DummyWorkorderFulfillment()
	so := model.DummySalesOrder()

	iv := model.DummyItemVariant()
	iv.CommitedStock = 20
	iv.Save("CommitedStock")

	soitem := model.DummySalesOrderItem()
	soitem.ItemVariant = iv
	soitem.SalesOrder = so
	soitem.Save("ItemVariant", "SalesOrder")

	item1 := model.DummyWorkorderFulfillmentItem()
	item1.Quantity = 10
	item1.SalesOrderItem = soitem
	item1.WorkorderFulfillment = fulfillment
	item1.Save("Quantity", "WorkorderFulfillment", "SalesOrderItem")

	ivStock1 := model.DummyItemVariantStock()
	ivStock1.ItemVariant = iv
	ivStock1.AvailableStock = 12
	ivStock1.CreatedAt = time.Now()
	ivStock1.UnitCost = 1000
	ivStock1.Save("ItemVariant", "UnitCost", "CreatedAt", "AvailableStock")

	// dummy ke 2
	iv2 := model.DummyItemVariant()
	iv2.CommitedStock = 20
	iv2.Save("CommitedStock")

	soItem2 := model.DummySalesOrderItem()
	soItem2.ItemVariant = iv2
	soItem2.SalesOrder = so
	soItem2.Save("ItemVariant", "SalesOrder")

	item2 := model.DummyWorkorderFulfillmentItem()
	item2.Quantity = 10
	item2.SalesOrderItem = soItem2
	item2.WorkorderFulfillment = fulfillment
	item2.Save("Quantity", "WorkorderFulfillment", "SalesOrderItem")

	for a := 0; a < 3; a++ {
		ivStock2 := model.DummyItemVariantStock()
		ivStock2.ItemVariant = iv2
		ivStock2.AvailableStock = 5
		ivStock2.UnitCost = 1000
		ivStock2.Save("ItemVariant", "AvailableStock", "UnitCost")
	}

	var soitems []*model.SalesOrderItem
	soitems = append(soitems, soitem, soItem2)

	so.SalesOrderItems = soitems
	so.Save("SalesOrderItem")

	var items []*model.WorkorderFulfillmentItem
	items = append(items, item1, item2)

	fulfillment.WorkorderFulFillmentItems = items
	fulfillment.SalesOrder = so
	fulfillment.Save("WorkorderFulFillmentItems", "SalesOrder")

	fulfillment, e := approveFulfillment(fulfillment)
	assert.NoError(t, e)
	assert.NotEmpty(t, fulfillment)
}

func TestApproveFulfillment2(t *testing.T) {
	fulfillment := model.DummyWorkorderFulfillment()
	so := model.DummySalesOrder()

	iv := model.DummyItemVariant()
	iv.CommitedStock = 20
	iv.Save("CommitedStock")

	soitem := model.DummySalesOrderItem()
	soitem.ItemVariant = iv
	soitem.Quantity = 10
	soitem.SalesOrder = so
	soitem.Save("ItemVariant", "SalesOrder", "Quantity")

	item1 := model.DummyWorkorderFulfillmentItem()
	item1.Quantity = 10
	item1.SalesOrderItem = soitem
	item1.WorkorderFulfillment = fulfillment
	item1.Save("Quantity", "WorkorderFulfillment", "SalesOrderItem")

	ivStock1 := model.DummyItemVariantStock()
	ivStock1.ItemVariant = iv
	ivStock1.AvailableStock = 12
	ivStock1.CreatedAt = time.Now()
	ivStock1.UnitCost = 1000
	ivStock1.Save("ItemVariant", "UnitCost", "CreatedAt", "AvailableStock")

	// dummy ke 2
	iv2 := model.DummyItemVariant()
	iv2.CommitedStock = 20
	iv2.Save("CommitedStock")

	soItem2 := model.DummySalesOrderItem()
	soItem2.ItemVariant = iv2
	soItem2.SalesOrder = so
	soItem2.Quantity = 10
	soItem2.Save("ItemVariant", "SalesOrder", "Quantity")

	item2 := model.DummyWorkorderFulfillmentItem()
	item2.Quantity = 10
	item2.SalesOrderItem = soItem2
	item2.WorkorderFulfillment = fulfillment
	item2.Save("Quantity", "WorkorderFulfillment", "SalesOrderItem")

	for a := 0; a < 3; a++ {
		ivStock2 := model.DummyItemVariantStock()
		ivStock2.ItemVariant = iv2
		ivStock2.AvailableStock = 5
		ivStock2.UnitCost = 1000
		ivStock2.Save("ItemVariant", "AvailableStock", "UnitCost")
	}

	var soitems []*model.SalesOrderItem
	soitems = append(soitems, soitem, soItem2)

	so.SalesOrderItems = soitems
	so.Save("SalesOrderItem")

	var items []*model.WorkorderFulfillmentItem
	items = append(items, item1, item2)

	fulfillment.WorkorderFulFillmentItems = items
	fulfillment.SalesOrder = so
	fulfillment.SalesOrder.InvoiceStatus = "finished"
	fulfillment.SalesOrder.Save("InvoiceStatus")
	fulfillment.Save("WorkorderFulFillmentItems", "SalesOrder")

	fulfillment, e := approveFulfillment(fulfillment)
	assert.NoError(t, e)
	assert.NotEmpty(t, fulfillment)
}
