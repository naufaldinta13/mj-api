// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchaseReturn

import (
	"fmt"
	"testing"
	"time"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestCreatePurchaseReturnSukses(t *testing.T) {
	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	preturn := model.PurchaseReturn{
		PurchaseOrder:   &model.PurchaseOrder{ID: dporder.ID},
		RecognitionDate: time.Now(),
		Code:            common.RandomStr(10),
		TotalAmount:     float64(50000),
		Note:            "Wow...",
		IsDeleted:       int8(0),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		PurchaseReturnItems: []*model.PurchaseReturnItem{
			{
				PurchaseOrderItem: &model.PurchaseOrderItem{ID: dpOrderItem.ID},
				Quantity:          float32(10),
				Note:              common.RandomStr(15),
			},
		},
		DocumentStatus: "new",
	}

	pret, err := CreatePurchaseReturn(&preturn)
	assert.NoError(t, err, "Tidak boleh ada error")
	assert.NotNil(t, pret, "Tidak boleh kosong")
}

func TestCreatePurchaseReturnGagal(t *testing.T) {
	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	preturn := model.PurchaseReturn{
		PurchaseOrder:   nil,
		RecognitionDate: time.Now(),
		Code:            common.RandomStr(10),
		TotalAmount:     float64(50000),
		Note:            "Wow...",
		IsDeleted:       int8(0),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		PurchaseReturnItems: []*model.PurchaseReturnItem{
			{
				PurchaseOrderItem: &model.PurchaseOrderItem{ID: dpOrderItem.ID},
				Quantity:          float32(10),
				Note:              common.RandomStr(15),
			},
		},
		DocumentStatus: "new",
	}

	pret, err := CreatePurchaseReturn(&preturn)
	assert.Error(t, err, "terdapat error")
	assert.Nil(t, pret, "data kosong")
}

func TestGetAllPurchaseReturnErrorNoData(t *testing.T) {
	orm.NewOrm().Raw("DELETE FROM purchase_return;").Exec()
	data, total, err := GetAllPurchaseReturn(&orm.RequestQuery{})
	assert.Nil(t, data, "Data kosong tidak ada berisi")
	assert.NoError(t, err, "Tidak ada error")
	assert.Equal(t, int64(0), total, "Total data 0")
}

func TestGetAllPurchaseReturn(t *testing.T) {
	orm.NewOrm().Raw("DELETE FROM purchase_return_item;").Exec()
	orm.NewOrm().Raw("DELETE FROM purchase_return;").Exec()

	DummyPRI := model.DummyPurchaseReturnItem()
	DummyPRI.PurchaseReturn.IsDeleted = 0
	DummyPRI.PurchaseReturn.Save()

	data, total, err := GetAllPurchaseReturn(&orm.RequestQuery{})
	assert.NotNil(t, data, "Data tidak berisi")
	assert.NoError(t, err, "Tidak ada error")
	assert.Equal(t, int64(1), total, "Total data 1"+fmt.Sprint(data))
	for _, row := range data {
		assert.Equal(t, int8(0), row.IsDeleted)
		assert.NotNil(t, row.PurchaseReturnItems)
	}
}

func TestGetDetailPurchaseReturnNoData(t *testing.T) {
	orm.NewOrm().Raw("DELETE FROM purchase_return;").Exec()
	data, err := GetDetailPurchaseReturn("id", 1)
	assert.Nil(t, data, "Data kosong tidak ada berisi")
	assert.Error(t, err, "terdapat error")
}

func TestGetDetailPurchaseReturn(t *testing.T) {
	orm.NewOrm().Raw("DELETE FROM purchase_return_item;").Exec()

	dpreturn := model.DummyPurchaseReturn()
	dpreturn.IsDeleted = 0
	dpreturn.Save()
	dpreturnItem := model.DummyPurchaseReturnItem()
	dpreturnItem.PurchaseReturn = dpreturn
	dpreturnItem.Save()

	data, err := GetDetailPurchaseReturn("id", dpreturn.ID)
	assert.NotNil(t, data, "Data tidak berisi")
	assert.NoError(t, err, "Tidak ada error")
}

func TestUpdatePurchaseReturnSukses(t *testing.T) {
	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	preturn := model.PurchaseReturn{
		PurchaseOrder:   &model.PurchaseOrder{ID: dporder.ID},
		RecognitionDate: time.Now(),
		Code:            common.RandomStr(10),
		TotalAmount:     float64(50000),
		Note:            "Wow...",
		IsDeleted:       int8(0),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		PurchaseReturnItems: []*model.PurchaseReturnItem{
			{
				PurchaseOrderItem: &model.PurchaseOrderItem{ID: dpOrderItem.ID},
				Quantity:          float32(10),
				Note:              common.RandomStr(15),
			},
		},
		DocumentStatus: "new",
	}

	pret, _ := CreatePurchaseReturn(&preturn)
	// Update data yang berhasil dibuat
	pret.PurchaseReturnItems[0].Quantity = float32(10)
	pret.Note = "My Note"

	purchase, err := updatePurchaseReturn(pret)
	assert.NoError(t, err, "Tidak boleh ada error")
	assert.NotNil(t, purchase, "Tidak boleh kosong")
	assert.Equal(t, float32(10), purchase.PurchaseReturnItems[0].Quantity, "Quantity harus nya sama")
	assert.Equal(t, "My Note", purchase.Note, "Note harus nya sama")
}

func TestUpdatePurchaseReturnGagal(t *testing.T) {
	iVariant := model.DummyItemVariant()
	iVariant.Item = model.DummyItem()
	iVariant.IsDeleted = 0
	iVariant.Save()

	dporder := model.DummyPurchaseOrder()
	dporder.Supplier = model.DummyPartnership()
	dporder.IsDeleted = int8(0)
	dporder.DocumentStatus = "active"
	dporder.Tax = float32(10)
	dporder.TaxAmount = float64(4500)
	dporder.Discount = float32(10)
	dporder.DiscountAmount = float64(5000)
	dporder.ShipmentCost = float64(0)
	dporder.TotalCharge = float64(45000 + 4500)
	dporder.Save()

	dpOrderItem := model.DummyPurchaseOrderItem()
	dpOrderItem.UnitPrice = float64(1000)
	dpOrderItem.Quantity = float32(50)
	dpOrderItem.PurchaseOrder = dporder
	dpOrderItem.Save()

	woReceive := model.DummyWorkorderReceiving()
	woReceive.PurchaseOrder = dporder
	woReceive.IsDeleted = 0
	woReceive.DocumentStatus = "finished"
	woReceive.Save()

	woReceiveItem := model.DummyWorkorderReceivingItem()
	woReceiveItem.Quantity = float32(45)
	woReceiveItem.WorkorderReceiving = woReceive
	woReceiveItem.PurchaseOrderItem = dpOrderItem
	woReceiveItem.Save()

	preturn := model.PurchaseReturn{
		PurchaseOrder:   &model.PurchaseOrder{ID: dporder.ID},
		RecognitionDate: time.Now(),
		Code:            common.RandomStr(10),
		TotalAmount:     float64(50000),
		Note:            "Wow...",
		IsDeleted:       int8(0),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		PurchaseReturnItems: []*model.PurchaseReturnItem{
			{
				PurchaseOrderItem: &model.PurchaseOrderItem{ID: dpOrderItem.ID},
				Quantity:          float32(10),
				Note:              common.RandomStr(15),
			},
		},
		DocumentStatus: "new",
	}

	pret, _ := CreatePurchaseReturn(&preturn)
	// Update data yang berhasil dibuat
	pret.PurchaseOrder = nil
	pret.PurchaseReturnItems[0].Quantity = float32(10)
	pret.Note = "My Note"

	purchase, err := updatePurchaseReturn(pret)
	assert.Error(t, err, "Tidak boleh ada error")
	assert.Nil(t, purchase, "Tidak boleh kosong")
}

func TestCancelPurchaseReturnSukses(t *testing.T) {

	DummyPRI := model.DummyPurchaseReturnItem()
	DummyPRI.PurchaseReturn.IsDeleted = 0
	DummyPRI.PurchaseReturn.DocumentStatus = "new"
	DummyPRI.PurchaseReturn.Save()

	DummyFR := model.DummyFinanceRevenue()
	DummyFR.RefID = uint64(DummyPRI.PurchaseReturn.ID)
	DummyFR.RefType = "purchase_return"
	DummyFR.IsDeleted = 0
	DummyFR.Save()

	ResPR, err := cancelPurchaseReturn(DummyPRI.PurchaseReturn)
	assert.NoError(t, err, "tidak boleh ada error jika sukses")
	assert.NotNil(t, ResPR, "Tidak boleh kosong datanya jika berhasil")
	var finance []*model.FinanceRevenue
	orm.NewOrm().Raw("SELECT * FROM finance_revenue WHERE ref_id = ? AND ref_type = 'purchase_return';", ResPR.ID).QueryRows(&finance)
	for _, row := range finance {
		assert.Equal(t, int8(1), row.IsDeleted, "Jika sukses harusnya ikut terhapus")
	}
}

func TestCancelPurchaseReturnGagal(t *testing.T) {

	DummyPRI := model.DummyPurchaseReturnItem()
	DummyPRI.PurchaseReturn.IsDeleted = 0
	DummyPRI.PurchaseReturn.DocumentStatus = "new"
	DummyPRI.PurchaseReturn.Save()

	DummyFR := model.DummyFinanceRevenue()
	DummyFR.RefID = uint64(DummyPRI.PurchaseReturn.ID)
	DummyFR.RefType = "sales_invoice"
	DummyFR.IsDeleted = 0
	DummyFR.Save()

	ResPR, err := cancelPurchaseReturn(DummyPRI.PurchaseReturn)
	assert.NoError(t, err, "tidak boleh ada error")
	assert.Nil(t, ResPR, "Tidak boleh kosong datanya jika berhasil")
}
