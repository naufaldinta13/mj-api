// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"testing"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestSumStockCommited(t *testing.T) {

	// data item variant
	mi := model.DummyItemVariant()
	mi.CommitedStock = 5
	mi.Save()

	// data sales order
	so := model.DummySalesOrder()
	so.DocumentStatus = "new"
	so.IsDeleted = 0
	so.Save()

	// data sales order item pertama
	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.ItemVariant = mi
	soi.Quantity = 20
	soi.Save()

	// data sales order item kedua
	soi2 := model.DummySalesOrderItem()
	soi2.SalesOrder = so
	soi2.ItemVariant = mi
	soi2.Quantity = 20
	soi2.Save()

	// data workorder fulfillment
	wf := model.DummyWorkorderFulfillment()
	wf.SalesOrder = so
	wf.IsDeleted = 0
	wf.DocumentStatus = "finished"
	wf.CreatedBy = model.DummyUser()
	wf.Save()

	// data workorder fulfillment pertama
	wfi := model.DummyWorkorderFulfillmentItem()
	wfi.WorkorderFulfillment = wf
	wfi.SalesOrderItem = soi
	wfi.Quantity = 10
	wfi.Save()

	// data workorder fulfillment kedua
	wfi2 := model.DummyWorkorderFulfillmentItem()
	wfi2.WorkorderFulfillment = wf
	wfi2.SalesOrderItem = soi
	wfi2.Quantity = 20
	wfi2.Save()

	// total quantity sales order item (40) - total quantity workorder fulfillment item (30)
	qty, e := SumStockCommited(mi)
	assert.NoError(t, e)
	assert.Equal(t, float32(10), qty, "seharusnya total quantity berjumlah 10")

	// cek perubahan stock commit pada item variant
	var newStockCommit float32

	o := orm.NewOrm()
	o.Raw("select commited_stock from item_variant where id = ?", mi.ID).QueryRow(&newStockCommit)
	assert.Equal(t, float32(10), newStockCommit, "seharusnya stock commit berubah menjadi 10")
}

func TestSumStockCommitedwithDeletedSO(t *testing.T) {

	// jika pada so is_deleted = 1

	// data item variant
	mi := model.DummyItemVariant()
	mi.CommitedStock = 5
	mi.Save()

	// data sales order
	so := model.DummySalesOrder()
	so.DocumentStatus = "approved_cancel"
	so.IsDeleted = 1
	so.Save()

	// data sales order item pertama
	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.ItemVariant = mi
	soi.Quantity = 10
	soi.Save()

	// data sales order item kedua
	soi2 := model.DummySalesOrderItem()
	soi2.SalesOrder = so
	soi2.ItemVariant = mi
	soi2.Quantity = 20
	soi2.Save()

	// data workorder fulfillment
	wf := model.DummyWorkorderFulfillment()
	wf.SalesOrder = so
	wf.IsDeleted = 0
	wf.DocumentStatus = "finished"
	wf.CreatedBy = model.DummyUser()
	wf.Save()

	// data workorder fulfillment pertama
	wfi := model.DummyWorkorderFulfillmentItem()
	wfi.WorkorderFulfillment = wf
	wfi.SalesOrderItem = soi
	wfi.Quantity = 10
	wfi.Save()

	// data workorder fulfillment kedua
	wfi2 := model.DummyWorkorderFulfillmentItem()
	wfi2.WorkorderFulfillment = wf
	wfi2.SalesOrderItem = soi
	wfi2.Quantity = 20
	wfi2.Save()

	qty, e := SumStockCommited(mi)
	assert.NoError(t, e)
	assert.Equal(t, float32(0), qty, "seharusnya total quantity berjumlah 0")

	var newStockCommit3 float32

	// cek perubahan stock commit pada item variant
	o := orm.NewOrm()
	o.Raw("select commited_stock from item_variant where id = ?", mi.ID).QueryRow(&newStockCommit3)
	assert.Equal(t, float32(0), newStockCommit3, "seharusnya stock commit berjumlah 0")
}

func TestSumStockCommitedwithApprovedCancelSO(t *testing.T) {

	// jika document status pada so == approved_cancel

	// data item variant
	mi := model.DummyItemVariant()
	mi.CommitedStock = 5
	mi.Save()

	// data sales order
	so := model.DummySalesOrder()
	so.DocumentStatus = "approved_cancel"
	so.IsDeleted = 0
	so.Save()

	// data sales order item pertama
	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.ItemVariant = mi
	soi.Quantity = 10
	soi.Save()

	// data sales order item kedua
	soi2 := model.DummySalesOrderItem()
	soi2.SalesOrder = so
	soi2.ItemVariant = mi
	soi2.Quantity = 20
	soi2.Save()

	// data workorder fulfillment
	wf := model.DummyWorkorderFulfillment()
	wf.SalesOrder = so
	wf.IsDeleted = 0
	wf.DocumentStatus = "finished"
	wf.CreatedBy = model.DummyUser()
	wf.Save()

	// data workorder fulfillment pertama
	wfi := model.DummyWorkorderFulfillmentItem()
	wfi.WorkorderFulfillment = wf
	wfi.SalesOrderItem = soi
	wfi.Quantity = 10
	wfi.Save()

	// data workorder fulfillment kedua
	wfi2 := model.DummyWorkorderFulfillmentItem()
	wfi2.WorkorderFulfillment = wf
	wfi2.SalesOrderItem = soi
	wfi2.Quantity = 20
	wfi2.Save()

	// total quantity sales order item (40) - total quantity workorder fulfillment item (0)
	qty, e := SumStockCommited(mi)
	assert.NoError(t, e)
	assert.Equal(t, float32(0), qty, "seharusnya total quantity berjumlah 0")

	var newStockCommit2 float32

	// cek perubahan stock commit pada item variant
	o := orm.NewOrm()
	o.Raw("select commited_stock from item_variant where id = ?", mi.ID).QueryRow(&newStockCommit2)
	assert.Equal(t, float32(0), newStockCommit2, "seharusnya stock commit berjumlah 0")
}

func TestSumStockCommitedwithActivedFulfillment(t *testing.T) {

	// jika document status pada workorder fulfillment != finished
	mi2 := model.DummyItemVariant()
	mi2.CommitedStock = 9
	mi2.Save()

	so2 := model.DummySalesOrder()
	so2.DocumentStatus = "new"
	so2.IsDeleted = 0
	so2.Save()

	soi3 := model.DummySalesOrderItem()
	soi3.SalesOrder = so2
	soi3.ItemVariant = mi2
	soi3.Quantity = 20
	soi3.Save()

	soi4 := model.DummySalesOrderItem()
	soi4.SalesOrder = so2
	soi4.ItemVariant = mi2
	soi4.Quantity = 20
	soi4.Save()

	wf2 := model.DummyWorkorderFulfillment()
	wf2.SalesOrder = so2
	wf2.IsDeleted = 0
	wf2.DocumentStatus = "active"
	wf2.CreatedBy = model.DummyUser()
	wf2.Save()

	wfi3 := model.DummyWorkorderFulfillmentItem()
	wfi3.WorkorderFulfillment = wf2
	wfi3.SalesOrderItem = soi3
	wfi3.Quantity = 10
	wfi3.Save()

	wfi4 := model.DummyWorkorderFulfillmentItem()
	wfi4.WorkorderFulfillment = wf2
	wfi4.SalesOrderItem = soi3
	wfi4.Quantity = 20
	wfi4.Save()

	qty, e := SumStockCommited(mi2)
	assert.NoError(t, e)
	assert.Equal(t, float32(40), qty, "seharusnya total quantity berjumlah 40")

	var newStockCommit4 float32

	o := orm.NewOrm()
	o.Raw("select commited_stock from item_variant where id = ?", mi2.ID).QueryRow(&newStockCommit4)
	assert.Equal(t, float32(40), newStockCommit4, "seharusnya stock commit berjumlah 40")
}

func TestSumStockCommitedwithDeletedFulfillment(t *testing.T) {

	// jika pada workorder fulfillment is_deleted = 1
	mi3 := model.DummyItemVariant()
	mi3.CommitedStock = 9
	mi3.Save()

	so3 := model.DummySalesOrder()
	so3.DocumentStatus = "new"
	so3.IsDeleted = 0
	so3.Save()

	soi5 := model.DummySalesOrderItem()
	soi5.SalesOrder = so3
	soi5.ItemVariant = mi3
	soi5.Quantity = 20
	soi5.Save()

	soi6 := model.DummySalesOrderItem()
	soi6.SalesOrder = so3
	soi6.ItemVariant = mi3
	soi6.Quantity = 20
	soi6.Save()

	wf3 := model.DummyWorkorderFulfillment()
	wf3.SalesOrder = so3
	wf3.IsDeleted = 1
	wf3.DocumentStatus = "finished"
	wf3.CreatedBy = model.DummyUser()
	wf3.Save()

	wfi5 := model.DummyWorkorderFulfillmentItem()
	wfi5.WorkorderFulfillment = wf3
	wfi5.SalesOrderItem = soi5
	wfi5.Quantity = 10
	wfi5.Save()

	wfi6 := model.DummyWorkorderFulfillmentItem()
	wfi6.WorkorderFulfillment = wf3
	wfi6.SalesOrderItem = soi5
	wfi6.Quantity = 20
	wfi6.Save()

	qty, e := SumStockCommited(mi3)
	assert.NoError(t, e)
	assert.Equal(t, float32(40), qty, "seharusnya total quantity berjumlah 40")

	var newStockCommit5 float32

	o := orm.NewOrm()
	o.Raw("select commited_stock from item_variant where id = ?", mi3.ID).QueryRow(&newStockCommit5)
	assert.Equal(t, float32(40), newStockCommit5, "seharusnya stock commit berjumlah 40")
}

func TestSumAvailableStockItemVariantStock(t *testing.T) {

	mi := model.DummyItemVariant()
	mi.Save()

	variantStock := model.DummyItemVariantStock()
	variantStock.ItemVariant = mi
	variantStock.AvailableStock = 10
	variantStock.Save()

	variantStock2 := model.DummyItemVariantStock()
	variantStock2.ItemVariant = mi
	variantStock2.AvailableStock = 15
	variantStock2.Save()

	qty, e := sumAvailableStockItemVariantStock(mi)
	assert.NoError(t, e)
	assert.Equal(t, float32(25), qty, "seharusnya total quantity berjumlah 25")
}

func TestCalculateAvailableStockVariant(t *testing.T) {

	// data item variant
	mi := model.DummyItemVariant()
	mi.AvailableStock = 10
	mi.Save()

	// data sales order
	so := model.DummySalesOrder()
	so.DocumentStatus = "new"
	so.IsDeleted = 0
	so.Save()

	// data sales order item pertama
	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.ItemVariant = mi
	soi.Quantity = 30
	soi.Save()

	// data sales order item kedua
	soi2 := model.DummySalesOrderItem()
	soi2.SalesOrder = so
	soi2.ItemVariant = mi
	soi2.Quantity = 20
	soi2.Save()

	// data workorder fulfillment
	wf := model.DummyWorkorderFulfillment()
	wf.SalesOrder = so
	wf.IsDeleted = 0
	wf.DocumentStatus = "finished"
	wf.CreatedBy = model.DummyUser()
	wf.Save()

	// data workorder fulfillment pertama
	wfi := model.DummyWorkorderFulfillmentItem()
	wfi.WorkorderFulfillment = wf
	wfi.SalesOrderItem = soi
	wfi.Quantity = 10
	wfi.Save()

	// data workorder fulfillment kedua
	wfi2 := model.DummyWorkorderFulfillmentItem()
	wfi2.WorkorderFulfillment = wf
	wfi2.SalesOrderItem = soi
	wfi2.Quantity = 10
	wfi2.Save()

	// data variant stock pertama
	variantStock := model.DummyItemVariantStock()
	variantStock.ItemVariant = mi
	variantStock.AvailableStock = 10
	variantStock.Save()

	// data variant stock kedua
	variantStock2 := model.DummyItemVariantStock()
	variantStock2.ItemVariant = mi
	variantStock2.AvailableStock = 40
	variantStock2.Save()

	total, e := CalculateAvailableStockItemVariant(mi)
	assert.NoError(t, e)
	assert.Equal(t, float32(0), total, "seharusnya total berjumlah 20")

	var availableStock float32

	// harusnya available stock pada item variant berubah menjadi 20
	o := orm.NewOrm()
	o.Raw("select available_stock from item_variant where id = ?", mi.ID).QueryRow(&availableStock)
	assert.Equal(t, float32(50), availableStock, "seharusnya available stock pada item variant berjumlah 20")
}

func TestCalculateAvailableStockVariantwithDeletedSO(t *testing.T) {

	// jika pada so is_deleted = 1
	mi3 := model.DummyItemVariant()
	mi3.AvailableStock = 5
	mi3.Save()

	so3 := model.DummySalesOrder()
	so3.DocumentStatus = "new"
	so3.IsDeleted = 1
	so3.Save()

	soi5 := model.DummySalesOrderItem()
	soi5.SalesOrder = so3
	soi5.ItemVariant = mi3
	soi5.Quantity = 20
	soi5.Save()

	soi6 := model.DummySalesOrderItem()
	soi6.SalesOrder = so3
	soi6.ItemVariant = mi3
	soi6.Quantity = 10
	soi6.Save()

	wf3 := model.DummyWorkorderFulfillment()
	wf3.SalesOrder = so3
	wf3.IsDeleted = 0
	wf3.DocumentStatus = "finished"
	wf3.CreatedBy = model.DummyUser()
	wf3.Save()

	wfi5 := model.DummyWorkorderFulfillmentItem()
	wfi5.WorkorderFulfillment = wf3
	wfi5.SalesOrderItem = soi5
	wfi5.Quantity = 10
	wfi5.Save()

	wfi6 := model.DummyWorkorderFulfillmentItem()
	wfi6.WorkorderFulfillment = wf3
	wfi6.SalesOrderItem = soi5
	wfi6.Quantity = 20
	wfi6.Save()

	// data variant stock pertama
	variantStock5 := model.DummyItemVariantStock()
	variantStock5.ItemVariant = mi3
	variantStock5.AvailableStock = 10
	variantStock5.Save()

	// data variant stock kedua
	variantStock6 := model.DummyItemVariantStock()
	variantStock6.ItemVariant = mi3
	variantStock6.AvailableStock = 40
	variantStock6.Save()

	total, e := CalculateAvailableStockItemVariant(mi3)
	assert.NoError(t, e)
	assert.Equal(t, float32(0), total, "seharusnya total berjumlah 50")

	var availableStock3 float32

	// harusnya available stock pada item variant berubah menjadi 50
	o := orm.NewOrm()
	o.Raw("select available_stock from item_variant where id = ?", mi3.ID).QueryRow(&availableStock3)
	assert.Equal(t, float32(50), availableStock3, "seharusnya available stock pada item variant berjumlah 50")
}

func TestCalculateAvailableStockVariantwithApprovedCancelSO(t *testing.T) {

	// jika pada so is_deleted = 1
	mi3 := model.DummyItemVariant()
	mi3.AvailableStock = 5
	mi3.Save()

	so3 := model.DummySalesOrder()
	so3.DocumentStatus = "approved_cancel"
	so3.IsDeleted = 0
	so3.Save()

	soi5 := model.DummySalesOrderItem()
	soi5.SalesOrder = so3
	soi5.ItemVariant = mi3
	soi5.Quantity = 20
	soi5.Save()

	soi6 := model.DummySalesOrderItem()
	soi6.SalesOrder = so3
	soi6.ItemVariant = mi3
	soi6.Quantity = 10
	soi6.Save()

	wf3 := model.DummyWorkorderFulfillment()
	wf3.SalesOrder = so3
	wf3.IsDeleted = 0
	wf3.DocumentStatus = "finished"
	wf3.CreatedBy = model.DummyUser()
	wf3.Save()

	wfi5 := model.DummyWorkorderFulfillmentItem()
	wfi5.WorkorderFulfillment = wf3
	wfi5.SalesOrderItem = soi5
	wfi5.Quantity = 10
	wfi5.Save()

	wfi6 := model.DummyWorkorderFulfillmentItem()
	wfi6.WorkorderFulfillment = wf3
	wfi6.SalesOrderItem = soi5
	wfi6.Quantity = 20
	wfi6.Save()

	// data variant stock pertama
	variantStock5 := model.DummyItemVariantStock()
	variantStock5.ItemVariant = mi3
	variantStock5.AvailableStock = 10
	variantStock5.Save()

	// data variant stock kedua
	variantStock6 := model.DummyItemVariantStock()
	variantStock6.ItemVariant = mi3
	variantStock6.AvailableStock = 40
	variantStock6.Save()

	total, e := CalculateAvailableStockItemVariant(mi3)
	assert.NoError(t, e)
	assert.Equal(t, float32(0), total, "seharusnya total berjumlah 50")

	var availableStock3 float32

	// harusnya available stock pada item variant berubah menjadi 50
	o := orm.NewOrm()
	o.Raw("select available_stock from item_variant where id = ?", mi3.ID).QueryRow(&availableStock3)
	assert.Equal(t, float32(50), availableStock3, "seharusnya available stock pada item variant berjumlah 50")
}

func TestCalculateAvailableStockVariantwithDeletedFulfillment(t *testing.T) {

	// jika pada workorder fulfillment is_deleted = 1
	mi3 := model.DummyItemVariant()
	mi3.AvailableStock = 5
	mi3.Save()

	so3 := model.DummySalesOrder()
	so3.DocumentStatus = "new"
	so3.IsDeleted = 0
	so3.Save()

	soi5 := model.DummySalesOrderItem()
	soi5.SalesOrder = so3
	soi5.ItemVariant = mi3
	soi5.Quantity = 20
	soi5.Save()

	soi6 := model.DummySalesOrderItem()
	soi6.SalesOrder = so3
	soi6.ItemVariant = mi3
	soi6.Quantity = 10
	soi6.Save()

	wf3 := model.DummyWorkorderFulfillment()
	wf3.SalesOrder = so3
	wf3.IsDeleted = 1
	wf3.DocumentStatus = "finished"
	wf3.CreatedBy = model.DummyUser()
	wf3.Save()

	wfi5 := model.DummyWorkorderFulfillmentItem()
	wfi5.WorkorderFulfillment = wf3
	wfi5.SalesOrderItem = soi5
	wfi5.Quantity = 10
	wfi5.Save()

	wfi6 := model.DummyWorkorderFulfillmentItem()
	wfi6.WorkorderFulfillment = wf3
	wfi6.SalesOrderItem = soi5
	wfi6.Quantity = 20
	wfi6.Save()

	// data variant stock pertama
	variantStock5 := model.DummyItemVariantStock()
	variantStock5.ItemVariant = mi3
	variantStock5.AvailableStock = 10
	variantStock5.Save()

	// data variant stock kedua
	variantStock6 := model.DummyItemVariantStock()
	variantStock6.ItemVariant = mi3
	variantStock6.AvailableStock = 40
	variantStock6.Save()

	total, e := CalculateAvailableStockItemVariant(mi3)
	assert.NoError(t, e)
	assert.Equal(t, float32(0), total, "seharusnya total berjumlah 20")

	var availableStock3 float32

	// harusnya available stock pada item variant berubah menjadi 20
	o := orm.NewOrm()
	o.Raw("select available_stock from item_variant where id = ?", mi3.ID).QueryRow(&availableStock3)
	assert.Equal(t, float32(50), availableStock3, "seharusnya available stock pada item variant berjumlah 20")
}

func TestCalculateAvailableStockVariantwithActivedFulfillment(t *testing.T) {

	// jika document status pada workorder fulfillment == active (fulfillment != finished)
	mi2 := model.DummyItemVariant()
	mi2.AvailableStock = 5
	mi2.Save()

	so2 := model.DummySalesOrder()
	so2.DocumentStatus = "new"
	so2.IsDeleted = 0
	so2.Save()

	soi3 := model.DummySalesOrderItem()
	soi3.SalesOrder = so2
	soi3.ItemVariant = mi2
	soi3.Quantity = 20
	soi3.Save()

	soi4 := model.DummySalesOrderItem()
	soi4.SalesOrder = so2
	soi4.ItemVariant = mi2
	soi4.Quantity = 20
	soi4.Save()

	wf2 := model.DummyWorkorderFulfillment()
	wf2.SalesOrder = so2
	wf2.IsDeleted = 0
	wf2.DocumentStatus = "active"
	wf2.CreatedBy = model.DummyUser()
	wf2.Save()

	wfi3 := model.DummyWorkorderFulfillmentItem()
	wfi3.WorkorderFulfillment = wf2
	wfi3.SalesOrderItem = soi3
	wfi3.Quantity = 10
	wfi3.Save()

	wfi4 := model.DummyWorkorderFulfillmentItem()
	wfi4.WorkorderFulfillment = wf2
	wfi4.SalesOrderItem = soi3
	wfi4.Quantity = 20
	wfi4.Save()

	// data variant stock pertama
	variantStock3 := model.DummyItemVariantStock()
	variantStock3.ItemVariant = mi2
	variantStock3.AvailableStock = 10
	variantStock3.Save()

	// data variant stock kedua
	variantStock4 := model.DummyItemVariantStock()
	variantStock4.ItemVariant = mi2
	variantStock4.AvailableStock = 40
	variantStock4.Save()

	total, e := CalculateAvailableStockItemVariant(mi2)
	assert.NoError(t, e)
	assert.Equal(t, float32(0), total, "seharusnya total berjumlah 10")

	var availableStock2 float32

	// harusnya available stock pada item variant berubah menjadi 10
	o := orm.NewOrm()
	o.Raw("select available_stock from item_variant where id = ?", mi2.ID).QueryRow(&availableStock2)
	assert.Equal(t, float32(50), availableStock2, "seharusnya available stock pada item variant berjumlah 10")
}

func TestGetStockLog(t *testing.T) {
	sl := model.DummyItemVariantStockLog()
	sl.Save()

	qs := orm.RequestQuery{}
	_, _, e := GetStockLog(&qs)
	assert.NoError(t, e, "Data should be exists.")
}

func TestGetStockLogWithFulfillment(t *testing.T) {

	i := model.DummyItem()
	i.Save()

	iv := model.DummyItemVariant()
	iv.Item = i
	iv.Save()

	so := model.DummySalesOrder()
	so.Save()

	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.ItemVariant = iv
	soi.Save()

	wf := model.DummyWorkorderFulfillment()
	wf.SalesOrder = so
	wf.Save()

	wfi := model.DummyWorkorderFulfillmentItem()
	wfi.WorkorderFulfillment = wf
	wfi.SalesOrderItem = soi
	wfi.Save()

	ivs := model.DummyItemVariantStock()
	ivs.ItemVariant = iv
	ivs.Save()
	//(workorder_fulfillment,workorder_receiving,stockopname,direct_placement)"
	sl := model.DummyItemVariantStockLog()
	sl.ItemVariantStock = ivs
	sl.RefType = "workorder_fulfillment"
	sl.RefID = uint64(so.ID)
	sl.Save()

	qs := orm.RequestQuery{}
	_, _, e := GetStockLog(&qs)
	assert.NoError(t, e, "Data should be exists.")
}

func TestGetStockLogWithStockOpname(t *testing.T) {

	i := model.DummyItem()
	i.Save()

	iv := model.DummyItemVariant()
	iv.Item = i
	iv.Save()

	sop := model.DummyStockopname()
	sop.Save()

	ivs := model.DummyItemVariantStock()
	ivs.ItemVariant = iv
	ivs.Save()

	sopi := model.DummyStockopnameItem()
	sopi.Stockopname = sop
	sopi.ItemVariantStock = ivs
	sopi.Save()

	sl := model.DummyItemVariantStockLog()
	sl.ItemVariantStock = ivs
	sl.RefType = "stockopname"
	sl.RefID = uint64(sop.ID)
	sl.Save()

	qs := orm.RequestQuery{}
	_, _, e := GetStockLog(&qs)
	assert.NoError(t, e, "Data should be exists.")
}

func TestGetStockLogWithReceiving(t *testing.T) {

	i := model.DummyItem()
	i.Save()

	iv := model.DummyItemVariant()
	iv.Item = i
	iv.Save()

	po := model.DummyPurchaseOrder()
	po.Save()

	poi := model.DummyPurchaseOrderItem()
	poi.PurchaseOrder = po
	poi.ItemVariant = iv
	poi.Save()

	wr := model.DummyWorkorderReceiving()
	wr.PurchaseOrder = po
	wr.Save()

	wri := model.DummyWorkorderReceivingItem()
	wri.WorkorderReceiving = wr
	wri.PurchaseOrderItem = poi
	wri.Save()

	ivs := model.DummyItemVariantStock()
	ivs.ItemVariant = iv
	ivs.Save()

	sl := model.DummyItemVariantStockLog()
	sl.ItemVariantStock = ivs
	sl.RefType = "workorder_receiving"
	sl.RefID = uint64(po.ID)
	sl.Save()

	qs := orm.RequestQuery{}
	_, _, e := GetStockLog(&qs)
	assert.NoError(t, e, "Data should be exists.")
}
