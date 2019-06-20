// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package inventory

import (
	"fmt"
	"testing"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common/faker"
	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

// ClearStockAndStockLog membersihkan table item variant stock dan item variant stock log
func ClearStockAndStockLog() {
	o := orm.NewOrm()
	o.Raw("DELETE from item_variant_stock").Exec()
	o.Raw("DELETE from item_variant_stock_log").Exec()
}

// FakeItemVariantNoLogForItem membuat dummy item variant dan langsung dimasukkan ke item
func FakeItemVariantNoLogForItem(itm *model.Item, itmVarArchive int8, itmVarDelete int8) {
	var m model.ItemVariant
	faker.Fill(&m, "ID")
	m.Item = &model.Item{ID: itm.ID}
	m.Measurement = model.DummyMeasurement()
	m.CreatedBy = model.DummyUser()
	m.IsArchived = itmVarArchive
	m.IsDeleted = itmVarDelete
	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	itm.ItemVariants = append(itm.ItemVariants, &m)
}

// FakeItemVariantStock membuat dummy item variant stock berikut log
func FakeItemVariantStock(itemVarID int64, quantity float32, finQuantity float32, logType string, refType string, refID uint64) *model.ItemVariantStock {
	var m model.ItemVariantStock
	faker.Fill(&m, "ID")
	m.AvailableStock = quantity
	m.ItemVariant = &model.ItemVariant{ID: itemVarID}
	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	FakeLog(m.ID, quantity, finQuantity, logType, refType, refID)
	return &m
}

// FakeLog membuat dummy custom log
func FakeLog(itemVarStockID int64, quantity float32, finQuantity float32, logType string, refType string, refID uint64) *model.ItemVariantStockLog {
	var stockLog model.ItemVariantStockLog
	faker.Fill(&stockLog, "ID")
	stockLog.RefID = refID
	stockLog.LogType = logType
	stockLog.Quantity = quantity
	stockLog.RefType = refType
	stockLog.FinalStock = finQuantity
	stockLog.ItemVariantStock = &model.ItemVariantStock{ID: itemVarStockID}
	if e := stockLog.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &stockLog
}

// TestFifoStockOut melakukan test pada fungsi GetFifoStock
func TestFifoStockOut(t *testing.T) {
	//buat dummy item variant
	itemVar1 := model.DummyItemVariant()
	itemVar2 := model.DummyItemVariant()
	// buat stock item variant
	varStock1 := FakeItemVariantStock(itemVar1.ID, float32(10), float32(10), "in", "workorder_receiving", uint64(1))
	FakeItemVariantStock(itemVar2.ID, float32(5), float32(5), "in", "workorder_receiving", uint64(1))
	varStock3 := FakeItemVariantStock(itemVar1.ID, float32(20), float32(20), "in", "workorder_receiving", uint64(1))
	FakeItemVariantStock(itemVar1.ID, float32(30), float32(30), "in", "workorder_receiving", uint64(1))
	// buat dummy log pengurangan
	varStock1.AvailableStock = float32(5)
	varStock1.Save()
	FakeLog(varStock1.ID, float32(5), float32(5), "out", "stockopname", uint64(1))
	//ambil fifo stock 10,--->total stock 25 (5 dan 20)
	//1. mengambil stock dengan jumlah cukup
	res1, err1 := FifoStockOut(itemVar1.ID, float32(10), "workorder_fulfillment", uint64(1))
	assert.NoError(t, err1)
	assert.Equal(t, int(2), len(res1))

	varStock1.Read("ID")
	assert.Equal(t, float32(0), varStock1.AvailableStock)
	varStock3.Read("ID")
	assert.Equal(t, float32(15), varStock3.AvailableStock)
	stockLog := &model.ItemVariantStockLog{ItemVariantStock: varStock3, LogType: "out"}
	stockLog.Read("ItemVariantStock", "LogType")
	assert.Equal(t, float32(5), stockLog.Quantity)
	assert.Equal(t, float32(15), stockLog.FinalStock)

	//2. mengambil stock dengan jumlah stock tidak mencukupi
	_, err3 := FifoStockOut(itemVar2.ID, float32(10), "workorder_fulfillment", uint64(1))
	assert.Error(t, err3)
}

func TestFifoStockIn(t *testing.T) {
	// buat dummy item variant
	itemVar1 := model.DummyItemVariant()
	varStock1 := FakeItemVariantStock(itemVar1.ID, float32(10), float32(10), "in", "workorder_receiving", uint64(1))

	// buat dummy log pengurangan
	// saat ini belum dibuat shared lib untuk pengambilan stock
	varStock1.AvailableStock = float32(5)
	varStock1.Save()
	FakeLog(varStock1.ID, float32(5), float32(5), "out", "stockopname", uint64(1))
	//ambil fifo stock 10,--->total stock 25 (5 dan 20)

	res, e := FifoStockIn(itemVar1, float64(10000), float32(25), "direct_placement", uint64(2))

	assert.NoError(t, e)
	assert.Equal(t, itemVar1.ID, res.ItemVariant.ID)
	assert.Equal(t, float32(25), res.AvailableStock)
	assert.Equal(t, float64(10000), res.UnitCost)
	assert.Equal(t, float32(30), res.ItemVariant.AvailableStock)

	// cek stock log
	sLog := &model.ItemVariantStockLog{ItemVariantStock: res, RefID: uint64(2), RefType: "direct_placement"}
	sLog.Read("ItemVariantStock", "RefID", "RefType")
	assert.Equal(t, float32(25), sLog.Quantity)
	assert.Equal(t, "in", sLog.LogType)
	assert.Equal(t, float32(25), sLog.FinalStock)

}

func TestCheckCancelStockTrue(t *testing.T) {
	ClearStockAndStockLog()
	// buat dummy item variant
	itemVar1 := model.DummyItemVariant()
	itemVar2 := model.DummyItemVariant()
	varStock1 := FakeItemVariantStock(itemVar1.ID, float32(20), float32(20), "in", "workorder_receiving", uint64(1))
	FakeItemVariantStock(itemVar2.ID, float32(25), float32(25), "in", "workorder_receiving", uint64(1))
	// tambah in log untuk varStock1 dengan refID berbeda
	varStock1.AvailableStock = float32(30)
	varStock1.Save("available_stock")
	FakeLog(varStock1.ID, float32(10), float32(30), "in", "workorder_receiving", uint64(2))

	// cek cancel dengan return true
	_, res, e := CheckCancelStock(uint64(1), "workorder_receiving")
	assert.NoError(t, e)
	assert.True(t, res)
}

func TestCheckCancelStockFalse(t *testing.T) {
	// buat dummy item variant
	itemVar1 := model.DummyItemVariant()
	itemVar2 := model.DummyItemVariant()
	varStock1 := FakeItemVariantStock(itemVar1.ID, float32(20), float32(20), "in", "stockopname", uint64(1))
	FakeItemVariantStock(itemVar2.ID, float32(25), float32(25), "in", "stockopname", uint64(1))
	// tambah out log untuk varStock1 dengan refID sama dan refType beda
	varStock1.AvailableStock = float32(5)
	varStock1.Save("available_stock")
	FakeLog(varStock1.ID, float32(15), float32(5), "out", "workorder_fulfillment", uint64(2))

	// cek cancel dengan return false
	_, res, e := CheckCancelStock(uint64(1), "stockopname")

	assert.Error(t, e)
	assert.False(t, res)
}

func TestCancelStockSuccess(t *testing.T) {
	// buat dummy item variant
	itemVar1 := model.DummyItemVariant()
	itemVar2 := model.DummyItemVariant()

	// variant stock pertama
	varStock1 := FakeItemVariantStock(itemVar1.ID, float32(20), float32(20), "in", "direct_placement", uint64(1))
	// buat log stockopname dengan log type out
	varStock1.AvailableStock = float32(10)
	varStock1.Save("available_stock")
	FakeLog(varStock1.ID, float32(10), float32(10), "out", "stockopname", uint64(2))

	// variant stock pertama
	varStock2 := FakeItemVariantStock(itemVar2.ID, float32(25), float32(25), "in", "stockopname", uint64(2))
	// tambah in log untuk varStock2 dengan refID berbeda
	varStock2.AvailableStock = float32(30)
	varStock2.Save("available_stock")
	FakeLog(varStock2.ID, float32(5), float32(30), "in", "workorder_receiving", uint64(3))

	// cancel dengan hasil success
	e := CancelStock(uint64(2), "stockopname")
	assert.NoError(t, e)
	varStock1.Read("ID")
	assert.Equal(t, float32(20), varStock1.AvailableStock)
	varStock2.Read("ID")
	assert.Equal(t, float32(5), varStock2.AvailableStock)

	sLOG1 := &model.ItemVariantStockLog{ItemVariantStock: &model.ItemVariantStock{ID: varStock1.ID}, RefID: uint64(2), RefType: "stockopname", LogType: "in"}
	sLOG1.Read("ItemVariantStock", "RefID", "RefType", "LogType")
	assert.Equal(t, float32(10), sLOG1.Quantity)
	assert.Equal(t, float32(20), sLOG1.FinalStock)
	sLOG2 := &model.ItemVariantStockLog{ItemVariantStock: &model.ItemVariantStock{ID: varStock2.ID}, RefID: uint64(2), RefType: "stockopname", LogType: "out"}
	sLOG2.Read("ItemVariantStock", "RefID", "RefType", "LogType")
	assert.Equal(t, float32(25), sLOG2.Quantity)
	assert.Equal(t, float32(5), sLOG2.FinalStock)
}

func TestCancelStockFail(t *testing.T) {
	// buat dummy item variant
	itemVar1 := model.DummyItemVariant()
	itemVar2 := model.DummyItemVariant()
	varStock1 := FakeItemVariantStock(itemVar1.ID, float32(10), float32(10), "in", "workorder_receiving", uint64(1))
	varStock1.AvailableStock = float32(5)
	varStock1.Save("available_stock")
	FakeLog(varStock1.ID, float32(5), float32(5), "out", "stockopname", uint64(2))
	varStock2 := FakeItemVariantStock(itemVar2.ID, float32(25), float32(25), "in", "stockopname", uint64(2))
	// tambah skenario out--->in dan in--->out
	varStock1.AvailableStock = float32(15)
	varStock1.Save("available_stock")
	FakeLog(varStock1.ID, float32(10), float32(15), "in", "direct_placement", uint64(3))
	varStock2.AvailableStock = float32(10)
	varStock2.Save("available_stock")
	FakeLog(varStock2.ID, float32(15), float32(10), "out", "workorder_fulfillment", uint64(4))

	// cancel dengan hasil fail
	e := CancelStock(uint64(2), "stockopname")
	assert.Error(t, e)
}

func TestGetAllItemVariantNoItemVariant(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant").Exec()
	rq := orm.RequestQuery{}
	// test tidak ada data itemVariant
	m, total, e := GetAllItemVariant(&rq)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, m)
	assert.NoError(t, e)
}

func TestGetAllItemVariant(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant").Exec()
	// buat dummy item variant
	f1 := model.DummyItemVariant()
	f1.IsArchived = int8(0)
	f1.IsDeleted = int8(0)
	f1.Save()
	f2 := model.DummyItemVariant()
	f2.IsArchived = int8(1)
	f2.IsDeleted = int8(0)
	f2.Save()
	// dummy yang sudah dihapus
	f3 := model.DummyItemVariant()
	f3.IsArchived = int8(1)
	f3.IsDeleted = int8(1)
	f3.Save()
	rq := orm.RequestQuery{}
	// test tidak ada data itemVariant
	m, total, e := GetAllItemVariant(&rq)
	assert.Equal(t, int64(2), total)
	assert.NotEmpty(t, m)
	assert.NoError(t, e)
}

func TestGetDetailItemVariant(t *testing.T) {
	_, e := GetDetailItemVariant("id", 1000)
	assert.Error(t, e, "Response should be error, beacuse there are no data yet.")

	c := model.DummyItemVariant()
	cd, e := GetDetailItemVariant("id", c.ID)
	assert.NoError(t, e, "Data should be exists.")
	assert.Equal(t, c.ID, cd.ID, "ID Response should be a same.")
}

func TestGetDetailPurchaseOrderByItemVariant(t *testing.T) {

	c := model.DummyItemVariant()
	cd, e := GetDetailPurchaseOrderByItemVariant(c.ID)
	for _, x := range cd {
		assert.NoError(t, e, "Data should be exists.")
		assert.Equal(t, c.ID, x.ItemVariant.ID, "ID Response should be a same.")
	}
}

func TestGetDetailSalesOrderItemByItemVariant(t *testing.T) {

	c := model.DummyItemVariant()
	cd, e := GetDetailSalesOrderItemByItemVariant(c.ID)
	for _, x := range cd {
		assert.NoError(t, e, "Data should be exists.")
		assert.Equal(t, c.ID, x.ItemVariant.ID, "ID Response should be a same.")
	}
}

func TestGetDetailItemSuccess(t *testing.T) {
	// tambah dummy item
	itm := model.DummyItem()
	itm.IsDeleted = int8(0)
	itm.Save()
	// buat dummy item variant
	ivr := model.DummyItemVariant()
	ivr.Item = &model.Item{ID: itm.ID}
	ivr.IsDeleted = int8(0)
	ivr.Save()
	// buat dummy item variant deleted
	ivr2 := model.DummyItemVariant()
	ivr2.Item = &model.Item{ID: itm.ID}
	ivr2.IsDeleted = int8(1)
	ivr2.Save()
	// buat dummy item variant price
	ivPrice := model.DummyItemVariantPrice()
	ivPrice.ItemVariant = &model.ItemVariant{ID: ivr.ID}
	ivPrice.Save()

	m, e := GetDetailItem("id", itm.ID)
	assert.NoError(t, e)
	assert.NotEmpty(t, m)
	assert.Equal(t, itm.ItemType, m.ItemType)
	assert.Equal(t, itm.ItemName, m.ItemName)
	assert.Equal(t, int(1), len(m.ItemVariants))
	for _, u := range m.ItemVariants {
		assert.Equal(t, ivr.VariantName, u.VariantName)
		assert.Equal(t, ivr.ExternalName, u.ExternalName)
		assert.Equal(t, int(1), len(u.ItemVariantPrices))
		for _, x := range u.ItemVariantPrices {
			assert.Equal(t, ivPrice.UnitPrice, x.UnitPrice)
		}
	}
}

func TestGetDetailItemFailDeleted(t *testing.T) {
	// tambah dummy item
	itm := model.DummyItem()
	itm.IsDeleted = int8(1)
	itm.Save()

	m, e := GetDetailItem("id", itm.ID)
	assert.Error(t, e)
	assert.Empty(t, m)
}

func TestGetDetailItemFailNoID(t *testing.T) {
	m, e := GetDetailItem("id", int64(9999999))
	assert.Error(t, e)
	assert.Empty(t, m)
}

func TestSavingItemSuccess(t *testing.T) {
	// dummy item variant price
	var itmVP model.ItemVariantPrice
	faker.Fill(&itmVP, "ID")
	itmVP.ItemVariant = model.DummyItemVariant()
	itmVP.PricingType = model.DummyPricingType()
	// dummy item variant
	var itmV model.ItemVariant
	faker.Fill(&itmV, "ID")
	itmV.Item = model.DummyItem()
	itmV.Measurement = model.DummyMeasurement()
	itmV.CreatedBy = model.DummyUser()
	itmV.ItemVariantPrices = append(itmV.ItemVariantPrices, &itmVP)
	// dummy item
	var itm model.Item
	faker.Fill(&itm, "ID")
	itm.Category = model.DummyItemCategory()
	itm.CreatedBy = model.DummyUser()
	itm.ItemVariants = append(itm.ItemVariants, &itmV)
	//save item
	m, e := SavingItem(itm)
	assert.NoError(t, e)
	assert.NotEmpty(t, m)
	assert.Equal(t, itm.ItemName, m.ItemName)
	assert.Equal(t, itm.ItemType, m.ItemType)
	for _, u1 := range m.ItemVariants {
		assert.Equal(t, itmV.Barcode, u1.Barcode)
		for _, u2 := range u1.ItemVariantPrices {
			assert.Equal(t, itmVP.UnitPrice, u2.UnitPrice)
		}
	}
}

func TestSavingItemFailItmVarPriceError(t *testing.T) {
	// dummy item variant price no pricing type
	var itmVP model.ItemVariantPrice
	faker.Fill(&itmVP, "ID")
	// dummy item variant
	var itmV model.ItemVariant
	faker.Fill(&itmV, "ID")
	itmV.Item = model.DummyItem()
	itmV.Measurement = model.DummyMeasurement()
	itmV.CreatedBy = model.DummyUser()
	itmV.ItemVariantPrices = append(itmV.ItemVariantPrices, &itmVP)
	// dummy item
	var itm model.Item
	faker.Fill(&itm, "ID")
	itm.Category = model.DummyItemCategory()
	itm.CreatedBy = model.DummyUser()
	itm.ItemVariants = append(itm.ItemVariants, &itmV)
	//save item
	m, e := SavingItem(itm)
	assert.Error(t, e)
	assert.Empty(t, m)
}

func TestSavingItemFailItemVariantError(t *testing.T) {
	// dummy item variant price
	var itmVP model.ItemVariantPrice
	faker.Fill(&itmVP, "ID")
	itmVP.ItemVariant = model.DummyItemVariant()
	itmVP.PricingType = model.DummyPricingType()
	// dummy item variant no created by and measurement
	var itmV model.ItemVariant
	faker.Fill(&itmV, "ID")
	itmV.ItemVariantPrices = append(itmV.ItemVariantPrices, &itmVP)
	// dummy item
	var itm model.Item
	faker.Fill(&itm, "ID")
	itm.Category = model.DummyItemCategory()
	itm.CreatedBy = model.DummyUser()
	itm.ItemVariants = append(itm.ItemVariants, &itmV)
	//save item
	m, e := SavingItem(itm)
	assert.Error(t, e)
	assert.Empty(t, m)
}

func TestArchivedItemSuccess(t *testing.T) {
	// buat dummy
	usr := model.DummyUserPriviledgeWithUsergroup(int64(1))
	session := &auth.SessionData{User: usr}

	itm := model.DummyItem()
	itm.IsArchived = int8(0)
	itm.IsDeleted = int8(0)
	itm.Save()
	// memasukkan dummy item var ke item
	FakeItemVariantNoLogForItem(itm, int8(0), int8(0))

	e := ArchivedItem(itm, session)
	assert.NoError(t, e)
	assert.Equal(t, int8(1), itm.IsArchived)
	assert.Equal(t, int8(0), itm.IsDeleted)
	assert.Equal(t, usr.ID, itm.UpdatedBy.ID)

	for _, u := range itm.ItemVariants {
		assert.Equal(t, int8(1), u.IsArchived)
		assert.Equal(t, int8(0), u.IsDeleted)
		assert.Equal(t, usr.ID, u.UpdatedBy.ID)
	}
}

func TestUnarchivedItemSuccess(t *testing.T) {
	// buat dummy
	usr := model.DummyUserPriviledgeWithUsergroup(int64(1))
	session := &auth.SessionData{User: usr}

	itm := model.DummyItem()
	itm.IsArchived = int8(1)
	itm.IsDeleted = int8(0)
	itm.Save()
	// memasukkan dummy item var ke item
	FakeItemVariantNoLogForItem(itm, int8(1), int8(0))

	e := UnarchivedItem(itm, session)
	assert.NoError(t, e)
	assert.Equal(t, int8(0), itm.IsArchived)
	assert.Equal(t, int8(0), itm.IsDeleted)
	assert.Equal(t, usr.ID, itm.UpdatedBy.ID)

	for _, u := range itm.ItemVariants {
		assert.Equal(t, int8(0), u.IsArchived)
		assert.Equal(t, int8(0), u.IsDeleted)
		assert.Equal(t, usr.ID, u.UpdatedBy.ID)
	}
}

func TestDeleterItem(t *testing.T) {
	// buat dummy
	itm := model.DummyItem()
	itm.IsArchived = int8(1)
	itm.IsDeleted = int8(0)
	itm.Save()
	// memasukkan dummy item var ke item
	FakeItemVariantNoLogForItem(itm, int8(1), int8(0))

	e := DeleterItem(itm)
	assert.NoError(t, e)
	assert.Equal(t, int8(1), itm.IsArchived)
	assert.Equal(t, int8(1), itm.IsDeleted)

	for _, u := range itm.ItemVariants {
		assert.Equal(t, int8(1), u.IsArchived)
		assert.Equal(t, int8(1), u.IsDeleted)
	}
}
