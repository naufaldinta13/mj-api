// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package shipment

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestGetShipmentsDataIsAvailableAndNoError(t *testing.T) {
	orm.NewOrm().Raw("DELETE FROM workorder_shipment").Exec()
	mws1 := model.DummyWorkorderShipment()
	mws2 := model.DummyWorkorderShipment()

	mws1.IsDeleted = 0
	mws1.Save()

	mws2.IsDeleted = 0
	mws2.Save()

	shipments, total, err := GetShipments(&orm.RequestQuery{})

	assert.NotEmpty(t, shipments, "Data shipments tidak Boleh Kosong")
	assert.Equal(t, int64(2), total, "Data total shipments harus sesuai")
	assert.NoError(t, err, "Saat get shipments tidak ada error")
}

func TestGetShipmentsNoData(t *testing.T) {
	orm.NewOrm().Raw("DELETE FROM workorder_shipment").Exec()

	shipments, total, _ := GetShipments(&orm.RequestQuery{})

	assert.Empty(t, shipments, "Data shipments kosong")
	assert.Equal(t, int64(0), total, "Data total shipments harus sesuai")
}

func TestGetDetailShipment(t *testing.T) {
	mws := model.DummyWorkorderShipment()
	mws.IsDeleted = 0
	mws.Save()

	shipment, err := GetDetailShipment("id", mws.ID)

	assert.NotEmpty(t, shipment, "Harus dapatkan 1 data workorder shipment")
	assert.NoError(t, err, "Tidak ada error saat get detail")
}

func TestGetDetailShipmentNoData(t *testing.T) {
	mws := model.DummyWorkorderShipment()
	mws.IsDeleted = 1
	mws.Save()

	_, err := GetDetailShipment("id", mws.ID)

	assert.Error(t, err, "Tidak ada error saat get detail")
}

// TestCreateWorkOrderShipment test create workorder shipment,success
func TestCreateWorkOrderShipment(t *testing.T) {
	var woShipItems []*model.WorkorderShipmentItem
	// buat dummy and init
	woFulfill := model.DummyWorkorderFulfillment()
	woShipItem := &model.WorkorderShipmentItem{WorkorderFulfillment: woFulfill}
	woShipItems = append(woShipItems, woShipItem)
	m := model.WorkorderShipment{
		Code:           "code",
		Priority:       "routine",
		TruckNumber:    "123456",
		DocumentStatus: "finished",
		IsDeleted:      int8(0),
		CreatedAt:      time.Now(),
		CreatedBy:      model.DummyUser(),
		Note:           "note",
		WorkorderShipmentItems: woShipItems,
	}
	e := CreateWorkOrderShipment(&m)
	assert.NoError(t, e)
	assert.NotEqual(t, int64(0), m.ID)
}

// TestCreateWorkOrderShipmentFail test create workorder shipment gagal,fail
func TestCreateWorkOrderShipmentFail(t *testing.T) {
	var woShipItems []*model.WorkorderShipmentItem
	// buat dummy and init
	woShipItem := &model.WorkorderShipmentItem{}
	woShipItems = append(woShipItems, woShipItem)
	m := model.WorkorderShipment{
		Code:           "code",
		Priority:       "routine",
		TruckNumber:    "123456",
		DocumentStatus: "finished",
		IsDeleted:      int8(0),
		CreatedAt:      time.Now(),
		CreatedBy:      model.DummyUser(),
		Note:           "note",
		WorkorderShipmentItems: woShipItems,
	}
	e := CreateWorkOrderShipment(&m)
	assert.Error(t, e)
}

// TestCancelShipmentAndFulfillment1 Sukses
func TestCancelShipmentAndFulfillment1(t *testing.T) {
	dmSO := model.DummySalesOrder()
	dmSOI := model.DummySalesOrderItem()
	dmWoS := model.DummyWorkorderShipment()
	dmWoSI := model.DummyWorkorderShipmentItem()
	dmWoF := model.DummyWorkorderFulfillment()
	dmWoFI := model.DummyWorkorderFulfillmentItem()

	dmSO.IsDeleted = int8(0)
	dmSO.Save()

	dmSOI.SalesOrder = dmSO
	dmSOI.Quantity = float32(10)
	dmSOI.Save()

	dmWoFI.WorkorderFulfillment = dmWoF
	dmWoFI.Save()

	dmWoSI.WorkorderShipment = dmWoS
	dmWoSI.WorkorderFulfillment = dmWoF
	dmWoSI.Save()

	dmWoS.IsDeleted = int8(0)
	dmWoS.Save()

	dmWoF.IsDeleted = int8(0)
	dmWoF.SalesOrder = dmSO
	dmWoF.Save()

	shipmentx, _ := GetDetailShipment("id", dmWoS.ID)
	_, err := cancelShipmentAndFulfillment(shipmentx, dmWoF.ID)
	for _, row := range shipmentx.WorkorderShipmentItems {
		row.Read()
		assert.Empty(t, row.ID, fmt.Sprint("Hapus fulfillment id yang dipilih di workorder shipment item, Actual:", row))
	}
	dmWoF.Read()
	assert.Equal(t, int8(0), dmWoF.IsDelivered, "Is Delivered pada fulfillment menjadi 0, Actual: "+strconv.Itoa(int(dmWoF.IsDelivered)))
	assert.NoError(t, err, "Tidak boleh ada error")
}

// TestCancelShipmentAndFulfillment2 Sukses  update SalesOrder ShipmentStatus menjadi finish
func TestCancelShipmentAndFulfillment2(t *testing.T) {
	dmSO := model.DummySalesOrder()
	dmSOI := model.DummySalesOrderItem()
	dmWoS := model.DummyWorkorderShipment()
	dmWoSI1 := model.DummyWorkorderShipmentItem()
	dmWoSI2 := model.DummyWorkorderShipmentItem()
	dmWoF1 := model.DummyWorkorderFulfillment()
	dmWoF2 := model.DummyWorkorderFulfillment()
	dmWoFI1 := model.DummyWorkorderFulfillmentItem()
	dmWoFI2 := model.DummyWorkorderFulfillmentItem()

	dmWoS.IsDeleted = int8(0)
	dmWoS.Save()

	dmSOI.SalesOrder = dmSO
	dmSOI.Quantity = float32(10)
	dmSOI.Save()

	dmWoF1.SalesOrder = dmSO
	dmWoF1.IsDelivered = int8(0)
	dmWoF1.Save()

	dmWoF2.SalesOrder = dmSO
	dmWoF2.IsDelivered = int8(1)
	dmWoF2.Save()

	dmWoFI1.WorkorderFulfillment = dmWoF1
	dmWoFI1.Quantity = float32(5)
	dmWoFI1.Save()

	dmWoFI2.WorkorderFulfillment = dmWoF2
	dmWoFI2.Quantity = float32(10)
	dmWoFI2.Save()

	dmWoSI1.WorkorderShipment = dmWoS
	dmWoSI1.WorkorderFulfillment = dmWoF1
	dmWoSI1.Save()

	dmWoSI2.WorkorderShipment = dmWoS
	dmWoSI2.WorkorderFulfillment = dmWoF2
	dmWoSI2.Save()

	shipmentx, _ := GetDetailShipment("id", dmWoS.ID)
	shipment, _ := cancelShipmentAndFulfillment(shipmentx, dmWoF1.ID)

	var fulfill model.SalesOrder
	for _, row := range shipment.WorkorderShipmentItems {
		fulfill = model.SalesOrder{ID: row.WorkorderFulfillment.SalesOrder.ID}
		fulfill.Read("ID")
		assert.Equal(t, "active", fulfill.ShipmentStatus, "update shipment status dari sales order dengan referensi fulfillment yang ada menjadi finished, actual: "+fulfill.ShipmentStatus)
		assert.Equal(t, "active", fulfill.DocumentStatus, "update shipment status dari sales order dengan referensi fulfillment yang ada menjadi finished, actual: "+fulfill.ShipmentStatus)
	}
}

// TestCancelShipmentAndFulfillment3 Sukses update SalesOrder ShipmentStatus menjadi active
func TestCancelShipmentAndFulfillment3(t *testing.T) {
	dmSO := model.DummySalesOrder()
	dmSOI := model.DummySalesOrderItem()
	dmWoS := model.DummyWorkorderShipment()
	dmWoSI1 := model.DummyWorkorderShipmentItem()
	dmWoSI2 := model.DummyWorkorderShipmentItem()
	dmWoF1 := model.DummyWorkorderFulfillment()
	dmWoF2 := model.DummyWorkorderFulfillment()
	dmWoFI1 := model.DummyWorkorderFulfillmentItem()
	dmWoFI2 := model.DummyWorkorderFulfillmentItem()

	dmWoS.IsDeleted = int8(0)
	dmWoS.Save()

	dmSOI.SalesOrder = dmSO
	dmSOI.Quantity = float32(20)
	dmSOI.Save()

	dmWoF1.SalesOrder = dmSO
	dmWoF1.IsDelivered = int8(0)
	dmWoF1.Save()

	dmWoF2.SalesOrder = dmSO
	dmWoF2.IsDelivered = int8(1)
	dmWoF2.Save()

	dmWoFI1.WorkorderFulfillment = dmWoF1
	dmWoFI1.Quantity = float32(5)
	dmWoFI1.Save()

	dmWoFI2.WorkorderFulfillment = dmWoF2
	dmWoFI2.Quantity = float32(5)
	dmWoFI2.Save()

	dmWoSI1.WorkorderShipment = dmWoS
	dmWoSI1.WorkorderFulfillment = dmWoF1
	dmWoSI1.Save()

	dmWoSI2.WorkorderShipment = dmWoS
	dmWoSI2.WorkorderFulfillment = dmWoF2
	dmWoSI2.Save()

	shipmentx, _ := GetDetailShipment("id", dmWoS.ID)
	shipment, _ := cancelShipmentAndFulfillment(shipmentx, dmWoF1.ID)

	var fulfill model.WorkorderFulfillment
	for _, row := range shipment.WorkorderShipmentItems {
		orm.DefaultRelsDepth = 3
		orm.NewOrm().QueryTable(new(model.WorkorderFulfillment)).Filter("id", row.WorkorderFulfillment.ID).RelatedSel().Limit(1).One(&fulfill)
		assert.Equal(t, "active", fulfill.SalesOrder.ShipmentStatus, "update shipment status dari sales order dengan referensi fulfillment yang ada menjadi active, actual: "+fulfill.SalesOrder.ShipmentStatus)
	}
}

// TestCancelShipmentAndFulfillment4 Sukses update is delete pada fulfillment
func TestCancelShipmentAndFulfillment4(t *testing.T) {
	dmSO := model.DummySalesOrder()
	dmSOI := model.DummySalesOrderItem()
	dmWoS := model.DummyWorkorderShipment()
	dmWoSI1 := model.DummyWorkorderShipmentItem()
	dmWoF1 := model.DummyWorkorderFulfillment()
	dmWoFI1 := model.DummyWorkorderFulfillmentItem()

	dmSO.DocumentStatus = "new"
	dmSO.ShipmentStatus = "new"
	dmSO.Save("DocumentStatus")

	dmWoS.IsDeleted = int8(0)
	dmWoS.Save()

	dmSOI.SalesOrder = dmSO
	dmSOI.Quantity = float32(20)
	dmSOI.Save()

	dmWoF1.SalesOrder = dmSO
	dmWoF1.IsDelivered = int8(1)
	dmWoF1.Save()

	dmWoFI1.WorkorderFulfillment = dmWoF1
	dmWoFI1.Quantity = float32(5)
	dmWoFI1.Save()

	dmWoSI1.WorkorderShipment = dmWoS
	dmWoSI1.WorkorderFulfillment = dmWoF1
	dmWoSI1.Save()

	shipmentx, _ := GetDetailShipment("id", dmWoS.ID)
	shipment, _ := cancelShipmentAndFulfillment(shipmentx, dmWoF1.ID)

	var fulfill model.WorkorderFulfillment
	for _, row := range shipment.WorkorderShipmentItems {
		orm.DefaultRelsDepth = 3
		orm.NewOrm().QueryTable(new(model.WorkorderFulfillment)).Filter("id", row.WorkorderFulfillment.ID).RelatedSel().Limit(1).One(&fulfill)
		assert.Equal(t, "active", fulfill.SalesOrder.ShipmentStatus, "update shipment status dari sales order dengan referensi fulfillment yang ada menjadi active, actual: "+fulfill.SalesOrder.ShipmentStatus)
		assert.Equal(t, "active", fulfill.SalesOrder.DocumentStatus, "update shipment status dari sales order dengan referensi fulfillment yang ada menjadi active, actual: "+fulfill.SalesOrder.ShipmentStatus)
	}
}

// TestApproveShipment1 Semua skenario di blueprint Berhasil Sales order shipment status finished
func TestApproveShipment1(t *testing.T) {
	dmSO := model.DummySalesOrder()
	dmSOI := model.DummySalesOrderItem()
	dmWoS := model.DummyWorkorderShipment()
	dmWoSI1 := model.DummyWorkorderShipmentItem()
	dmWoSI2 := model.DummyWorkorderShipmentItem()
	dmWoF1 := model.DummyWorkorderFulfillment()
	dmWoF2 := model.DummyWorkorderFulfillment()
	dmWoFI1 := model.DummyWorkorderFulfillmentItem()
	dmWoFI2 := model.DummyWorkorderFulfillmentItem()

	dmWoS.IsDeleted = int8(0)
	dmSO.FulfillmentStatus = "finished"
	dmSO.InvoiceStatus = "active"
	dmWoS.Save()

	dmSOI.SalesOrder = dmSO
	dmSOI.Quantity = float32(10)
	dmSOI.Save()

	dmWoF1.SalesOrder = dmSO
	dmWoF1.IsDelivered = int8(1)
	dmWoF1.Save()

	dmWoF2.SalesOrder = dmSO
	dmWoF2.IsDelivered = int8(1)
	dmWoF2.Save()

	dmWoFI1.WorkorderFulfillment = dmWoF1
	dmWoFI1.Quantity = float32(5)
	dmWoFI1.Save()

	dmWoFI2.WorkorderFulfillment = dmWoF2
	dmWoFI2.Quantity = float32(5)
	dmWoFI2.Save()

	dmWoSI1.WorkorderShipment = dmWoS
	dmWoSI1.WorkorderFulfillment = dmWoF1
	dmWoSI1.Save()

	dmWoSI2.WorkorderShipment = dmWoS
	dmWoSI2.WorkorderFulfillment = dmWoF2
	dmWoSI2.Save()

	shipmentx, _ := GetDetailShipment("id", dmWoS.ID)
	shipment, err := approveShipment(shipmentx)

	assert.NotNil(t, shipment, "Shipment hasil approve tidak boleh kosong")
	assert.Equal(t, "finished", shipment.DocumentStatus, "Update document status menjadi finished di shipment yang dilakukan")
	assert.NoError(t, err, "Tidak boleh ada error")
	for _, row := range shipment.WorkorderShipmentItems {
		assert.Equal(t, int8(1), row.WorkorderFulfillment.IsDelivered, "Update is_delivered pada fulfillment yang dipilih menjadi 1")
		assert.Equal(t, "finished", row.WorkorderFulfillment.SalesOrder.ShipmentStatus, "Update shipment status dari sales order dengan referensi fulfillment yang menjadi finished")
	}
}

// TestApproveShipment2 Semua skenario di blueprint Berhasil  Sales order shipment status active
func TestApproveShipment2(t *testing.T) {
	dmSO := model.DummySalesOrder()
	dmSOI := model.DummySalesOrderItem()
	dmWoS := model.DummyWorkorderShipment()
	dmWoSI1 := model.DummyWorkorderShipmentItem()
	dmWoSI2 := model.DummyWorkorderShipmentItem()
	dmWoF1 := model.DummyWorkorderFulfillment()
	dmWoF2 := model.DummyWorkorderFulfillment()
	dmWoFI1 := model.DummyWorkorderFulfillmentItem()
	dmWoFI2 := model.DummyWorkorderFulfillmentItem()

	dmWoS.IsDeleted = int8(0)
	dmSO.FulfillmentStatus = "finished"
	dmSO.InvoiceStatus = "finished"
	dmWoS.Save()

	dmSOI.SalesOrder = dmSO
	dmSOI.Quantity = float32(10)
	dmSOI.Save()

	dmWoF1.SalesOrder = dmSO
	dmWoF1.IsDelivered = int8(1)
	dmWoF1.Save()

	dmWoF2.SalesOrder = dmSO
	dmWoF2.IsDelivered = int8(1)
	dmWoF2.Save()

	dmWoFI1.WorkorderFulfillment = dmWoF1
	dmWoFI1.Quantity = float32(5)
	dmWoFI1.Save()

	dmWoFI2.WorkorderFulfillment = dmWoF2
	dmWoFI2.Quantity = float32(3)
	dmWoFI2.Save()

	dmWoSI1.WorkorderShipment = dmWoS
	dmWoSI1.WorkorderFulfillment = dmWoF1
	dmWoSI1.Save()

	dmWoSI2.WorkorderShipment = dmWoS
	dmWoSI2.WorkorderFulfillment = dmWoF2
	dmWoSI2.Save()

	shipmentx, _ := GetDetailShipment("id", dmWoS.ID)
	shipment, err := approveShipment(shipmentx)

	assert.NotNil(t, shipment, "Shipment hasil approve tidak boleh kosong")
	assert.Equal(t, "finished", shipment.DocumentStatus, "Update document status menjadi finished di shipment yang dilakukan")
	assert.NoError(t, err, "Tidak boleh ada error")
	for _, row := range shipment.WorkorderShipmentItems {
		assert.Equal(t, int8(1), row.WorkorderFulfillment.IsDelivered, "Update is_delivered pada fulfillment yang dipilih menjadi 1")
		assert.Equal(t, "active", row.WorkorderFulfillment.SalesOrder.ShipmentStatus, "Update shipment status dari sales order dengan referensi fulfillment yang menjadi finished")
	}
}

// TestApproveShipment3 Semua skenario di blueprint Berhasil Sales order document status finished
func TestApproveShipment3(t *testing.T) {
	dmSO := model.DummySalesOrder()
	dmSOI := model.DummySalesOrderItem()
	dmWoS := model.DummyWorkorderShipment()
	dmWoSI1 := model.DummyWorkorderShipmentItem()
	dmWoSI2 := model.DummyWorkorderShipmentItem()
	dmWoF1 := model.DummyWorkorderFulfillment()
	dmWoF2 := model.DummyWorkorderFulfillment()
	dmWoFI1 := model.DummyWorkorderFulfillmentItem()
	dmWoFI2 := model.DummyWorkorderFulfillmentItem()

	dmSO.InvoiceStatus = "finished"
	dmSO.FulfillmentStatus = "finished"
	dmSO.Save()

	dmWoS.IsDeleted = int8(0)
	dmWoS.Save()

	dmSOI.SalesOrder = dmSO
	dmSOI.Quantity = float32(10)
	dmSOI.Save()

	dmWoF1.SalesOrder = dmSO
	dmWoF1.IsDelivered = int8(1)
	dmWoF1.Save()

	dmWoF2.SalesOrder = dmSO
	dmWoF2.IsDelivered = int8(1)
	dmWoF2.Save()

	dmWoFI1.WorkorderFulfillment = dmWoF1
	dmWoFI1.Quantity = float32(5)
	dmWoFI1.Save()

	dmWoFI2.WorkorderFulfillment = dmWoF2
	dmWoFI2.Quantity = float32(5)
	dmWoFI2.Save()

	dmWoSI1.WorkorderShipment = dmWoS
	dmWoSI1.WorkorderFulfillment = dmWoF1
	dmWoSI1.Save()

	dmWoSI2.WorkorderShipment = dmWoS
	dmWoSI2.WorkorderFulfillment = dmWoF2
	dmWoSI2.Save()

	shipmentx, _ := GetDetailShipment("id", dmWoS.ID)
	shipment, err := approveShipment(shipmentx)

	assert.NotNil(t, shipment, "Shipment hasil approve tidak boleh kosong")
	assert.Equal(t, "finished", shipment.DocumentStatus, "Update document status menjadi finished di shipment yang dilakukan")
	assert.NoError(t, err, "Tidak boleh ada error")
	for _, row := range shipment.WorkorderShipmentItems {
		assert.Equal(t, int8(1), row.WorkorderFulfillment.IsDelivered, "Update is_delivered pada fulfillment yang dipilih menjadi 1")
		assert.Equal(t, "finished", row.WorkorderFulfillment.SalesOrder.ShipmentStatus, "Update shipment status dari sales order dengan referensi fulfillment yang menjadi finished")
		assert.Equal(t, "finished", row.WorkorderFulfillment.SalesOrder.DocumentStatus, "Update shipment status dari sales order dengan referensi fulfillment yang menjadi finished")
	}
}

// TestApproveShipment4 Semua skenario di blueprint Berhasil Sales order document status finished
func TestApproveShipment4(t *testing.T) {
	dmSO := model.DummySalesOrder()
	dmSOI := model.DummySalesOrderItem()
	dmWoS := model.DummyWorkorderShipment()
	dmWoSI1 := model.DummyWorkorderShipmentItem()
	dmWoF1 := model.DummyWorkorderFulfillment()
	dmWoFI1 := model.DummyWorkorderFulfillmentItem()

	dmSO.InvoiceStatus = "finished"
	dmSO.FulfillmentStatus = "finished"
	dmSO.InvoiceStatus = "finished"
	dmSO.Save()

	dmWoS.IsDeleted = int8(0)
	dmWoS.Save()

	dmSOI.SalesOrder = dmSO
	dmSOI.Quantity = float32(10)
	dmSOI.Save()

	dmWoF1.SalesOrder = dmSO
	dmWoF1.IsDelivered = int8(1)
	dmWoF1.Save()

	dmWoFI1.WorkorderFulfillment = dmWoF1
	dmWoFI1.Quantity = float32(10)
	dmWoFI1.Save()

	dmWoSI1.WorkorderShipment = dmWoS
	dmWoSI1.WorkorderFulfillment = dmWoF1
	dmWoSI1.Save()

	shipmentx, _ := GetDetailShipment("id", dmWoS.ID)
	shipment, err := approveShipment(shipmentx)

	assert.NotNil(t, shipment, "Shipment hasil approve tidak boleh kosong")
	assert.Equal(t, "finished", shipment.DocumentStatus, "Update document status menjadi finished di shipment yang dilakukan")
	assert.NoError(t, err, "Tidak boleh ada error")
	for _, row := range shipment.WorkorderShipmentItems {
		assert.Equal(t, int8(1), row.WorkorderFulfillment.IsDelivered, "Update is_delivered pada fulfillment yang dipilih menjadi 1")
		assert.Equal(t, "finished", row.WorkorderFulfillment.SalesOrder.ShipmentStatus, "Update shipment status dari sales order dengan referensi fulfillment yang menjadi finished")
		assert.Equal(t, "active", row.WorkorderFulfillment.SalesOrder.DocumentStatus, "Update shipment status dari sales order dengan referensi fulfillment yang menjadi finished")
	}
}

func TestShowShipment(t *testing.T) {
	//bikin SO
	so := model.DummySalesOrder()

	//bikin Item
	iv := model.DummyItemVariant()
	iv2 := model.DummyItemVariant()

	//Bikin SO item
	soi := model.DummySalesOrderItem()
	soi.SalesOrder = so
	soi.ItemVariant = iv
	soi.Save()

	//Bikin SO item
	soi2 := model.DummySalesOrderItem()
	soi2.SalesOrder = so
	soi2.ItemVariant = iv2
	soi2.Save()

	//bikin fulfillment dari SO
	ff := model.DummyWorkorderFulfillment()
	ff.SalesOrder = so
	ff.Save()

	//bikin fulfillment item dari fulfillment dan so item
	ffi := model.DummyWorkorderFulfillmentItem()
	ffi.SalesOrderItem = soi
	ffi.WorkorderFulfillment = ff
	ffi.Save()

	ffi2 := model.DummyWorkorderFulfillmentItem()
	ffi2.SalesOrderItem = soi2
	ffi2.WorkorderFulfillment = ff
	ffi2.Save()

	//bikin shipment
	ws := model.DummyWorkorderShipment()
	ws.IsDeleted = 0
	ws.Save()

	//bikin shipment item dari fulfillment item
	shi := model.DummyWorkorderShipmentItem()
	shi.WorkorderFulfillment = ff
	shi.WorkorderShipment = ws
	shi.Save()

	shipment, err := ShowShipment("id", ws.ID)

	assert.NotEmpty(t, shipment, "Harus dapatkan 1 data workorder shipment")
	assert.NoError(t, err, "Tidak ada error saat get detail")
}
