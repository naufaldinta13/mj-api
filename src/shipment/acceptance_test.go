// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package shipment_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/test"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/common/tester"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.Setup()

	// run tests
	res := m.Run()

	// cleanup
	test.DataCleanUp()

	os.Exit(res)
}

func TestARouting(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	shipment := model.DummyWorkorderShipment()
	shipment.IsDeleted = 0
	shipment.Save()

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/shipment", "GET", 200},
		{"/v1/shipment/" + common.Encrypt(shipment.ID), "GET", 200},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, ep := range routers {
		ng.Method = ep.method
		ng.Path = ep.endpoint
		ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, ep.expected, res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", ep.endpoint, ep.method))
		})
	}
}

// TestHandler_URLMappingPostWorkOrderShipmentSuccess mengetest create work order shipment,success
func TestHandler_URLMappingPostWorkOrderShipmentSuccess(t *testing.T) {
	// buat dummy
	woFulfill1 := model.DummyWorkorderFulfillment()
	woFulfill1.IsDeleted = int8(0)
	woFulfill1.IsDelivered = 0
	woFulfill1.DocumentStatus = "finished"
	woFulfill1.Save()

	woFulfill2 := model.DummyWorkorderFulfillment()
	woFulfill2.IsDeleted = int8(0)
	woFulfill2.IsDelivered = 0
	woFulfill2.DocumentStatus = "finished"
	woFulfill2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"priority":     "emergency",
		"truck_number": "123456",
		"note":         "note",
		"workorder_shipment_items": []tester.D{
			{
				"workorder_fulfillment": common.Encrypt(woFulfill1.ID),
			},
			{
				"workorder_fulfillment": common.Encrypt(woFulfill2.ID),
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/shipment").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	m := &model.WorkorderShipment{Priority: "emergency", TruckNumber: "123456", Note: "note"}
	e := m.Read("Priority", "TruckNumber", "Note")
	assert.NoError(t, e)
	n := &model.WorkorderShipmentItem{WorkorderShipment: m, WorkorderFulfillment: woFulfill1}
	err := n.Read("WorkorderShipment", "WorkorderFulfillment")
	assert.NoError(t, err)
}

// TestHandler_URLMappingPostWorkOrderShipmentFailNotFinish mengetest create work order  not finish,fail
func TestHandler_URLMappingPostWorkOrderShipmentFailNotFinish(t *testing.T) {
	// buat dummy
	woFulfill1 := model.DummyWorkorderFulfillment()
	woFulfill1.IsDeleted = int8(0)
	woFulfill1.DocumentStatus = "new"
	woFulfill1.Save()

	woFulfill2 := model.DummyWorkorderFulfillment()
	woFulfill2.IsDeleted = int8(0)
	woFulfill2.DocumentStatus = "active"
	woFulfill2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"priority":     "emergency",
		"truck_number": "123456",
		"note":         "note",
		"workorder_shipment_items": []tester.D{
			{
				"workorder_fulfillment": common.Encrypt(woFulfill1.ID),
			},
			{
				"workorder_fulfillment": common.Encrypt(woFulfill2.ID),
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/shipment").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMappingPostWorkOrderShipmentFailIsDeleted mengetest create work order deleted,fail
func TestHandler_URLMappingPostWorkOrderShipmentFailIsDeleted(t *testing.T) {
	// buat dummy
	woFulfill1 := model.DummyWorkorderFulfillment()
	woFulfill1.IsDeleted = int8(1)
	woFulfill1.DocumentStatus = "finished"
	woFulfill1.Save()

	woFulfill2 := model.DummyWorkorderFulfillment()
	woFulfill2.IsDeleted = int8(0)
	woFulfill2.DocumentStatus = "finished"
	woFulfill2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"priority":     "emergency",
		"truck_number": "123456",
		"note":         "note",
		"workorder_shipment_items": []tester.D{
			{
				"workorder_fulfillment": common.Encrypt(woFulfill1.ID),
			},
			{
				"workorder_fulfillment": common.Encrypt(woFulfill2.ID),
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/shipment").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMappingPostWorkOrderShipmentFailUndecrypt mengetest create work order undecrypt,fail
func TestHandler_URLMappingPostWorkOrderShipmentFailUndecrypt(t *testing.T) {
	// buat dummy
	woFulfill1 := model.DummyWorkorderFulfillment()
	woFulfill1.IsDeleted = int8(0)
	woFulfill1.DocumentStatus = "finished"
	woFulfill1.Save()

	woFulfill2 := model.DummyWorkorderFulfillment()
	woFulfill2.IsDeleted = int8(0)
	woFulfill2.DocumentStatus = "finished"
	woFulfill2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"priority":     "emergency",
		"truck_number": "123456",
		"note":         "note",
		"workorder_shipment_items": []tester.D{
			{
				"workorder_fulfillment": "assas",
			},
			{
				"workorder_fulfillment": common.Encrypt(woFulfill2.ID),
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/shipment").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPUTSukses Berhasil update tanpa ada perubahan dari user
func TestHandler_URLMappingPUTSukses1(t *testing.T) {
	// buat dummy
	woShipment := model.DummyWorkorderShipment()
	woShipment.IsDeleted = 0
	woShipment.Save()

	woFulfill1 := model.DummyWorkorderFulfillment()
	woFulfill1.IsDeleted = 0
	woFulfill1.IsDelivered = 0
	woFulfill1.DocumentStatus = "finished"
	woFulfill1.Save()
	woFulfill2 := model.DummyWorkorderFulfillment()
	woFulfill2.IsDeleted = 0
	woFulfill2.IsDelivered = 0
	woFulfill2.DocumentStatus = "finished"
	woFulfill2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"priority":     woShipment.Priority,
		"truck_number": woShipment.TruckNumber,
		"note":         woShipment.Note,
		"workorder_shipment_items": []tester.D{
			{
				"workorder_fulfillment": common.Encrypt(woFulfill1.ID),
			},
			{
				"workorder_fulfillment": common.Encrypt(woFulfill2.ID),
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/shipment/"+common.Encrypt(woShipment.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPUTSukses Berhasil update terdapat perubahan dari user
func TestHandler_URLMappingPUTSukses2(t *testing.T) {
	// buat dummy
	woShipment := model.DummyWorkorderShipment()
	woShipment.IsDeleted = 0
	woShipment.Save()

	woFulfill1 := model.DummyWorkorderFulfillment()
	woFulfill1.IsDeleted = 0
	woFulfill1.IsDelivered = 0
	woFulfill1.DocumentStatus = "finished"
	woFulfill1.Save()
	woFulfill2 := model.DummyWorkorderFulfillment()
	woFulfill2.IsDeleted = 0
	woFulfill2.IsDelivered = 0
	woFulfill2.DocumentStatus = "finished"
	woFulfill2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"priority":     woShipment.Priority,
		"truck_number": "B1234ASD",
		"note":         "Ubah Note",
		"workorder_shipment_items": []tester.D{
			{
				"workorder_fulfillment": common.Encrypt(woFulfill1.ID),
			},
			{
				"workorder_fulfillment": common.Encrypt(woFulfill2.ID),
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/shipment/"+common.Encrypt(woShipment.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMappingPostWorkOrderShipmentSuccess mengetest create work order shipment,gagal karena isdelivered fulfillmet =1
func TestHandler_URLMappingPostWorkOrderShipmentFailedIfIsDeliveredIs1(t *testing.T) {
	// buat dummy
	woFulfill1 := model.DummyWorkorderFulfillment()
	woFulfill1.IsDeleted = int8(0)
	woFulfill1.IsDelivered = 1
	woFulfill1.DocumentStatus = "finished"
	woFulfill1.Save()

	woFulfill2 := model.DummyWorkorderFulfillment()
	woFulfill2.IsDeleted = int8(0)
	woFulfill1.IsDelivered = 1
	woFulfill2.DocumentStatus = "finished"
	woFulfill2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"priority":     "emergency",
		"truck_number": "123456",
		"note":         "note",
		"workorder_shipment_items": []tester.D{
			{
				"workorder_fulfillment": common.Encrypt(woFulfill1.ID),
			},
			{
				"workorder_fulfillment": common.Encrypt(woFulfill2.ID),
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/shipment").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMappingPUTGagal1 Gagal update terdapat perubahan dari user fulfillment  finish
func TestHandler_URLMappingPUTGagal1(t *testing.T) {
	test.DataCleanUp("workorder_shipment", "workorder_fulfillment")
	// buat dummy
	woShipment := model.DummyWorkorderShipment()
	woShipment.IsDeleted = 0
	woShipment.Save()

	woFulfill1 := model.DummyWorkorderFulfillment()
	woFulfill1.IsDeleted = int8(0)
	woFulfill1.DocumentStatus = "active"
	woFulfill1.Save()

	woFulfill2 := model.DummyWorkorderFulfillment()
	woFulfill2.IsDeleted = int8(0)
	woFulfill2.DocumentStatus = "finished"
	woFulfill2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"priority":     woShipment.Priority,
		"truck_number": "B1234ASD",
		"note":         "Ubah Note",
		"workorder_shipment_items": []tester.D{
			{
				"workorder_fulfillment": common.Encrypt(woFulfill1.ID),
			},
			{
				"workorder_fulfillment": common.Encrypt(woFulfill2.ID),
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/shipment/"+common.Encrypt(woShipment.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMappingPUTGagal1 Gagal update terdapat perubahan dari user fulfillment is delete 1
func TestHandler_URLMappingPUTGagal2(t *testing.T) {
	test.DataCleanUp("workorder_shipment", "workorder_fulfillment")
	// buat dummy
	woShipment := model.DummyWorkorderShipment()
	woShipment.IsDeleted = 0
	woShipment.Save()

	woFulfill1 := model.DummyWorkorderFulfillment()
	woFulfill1.IsDeleted = int8(1)
	woFulfill1.DocumentStatus = "new"
	woFulfill1.Save()

	woFulfill2 := model.DummyWorkorderFulfillment()
	woFulfill2.IsDeleted = int8(1)
	woFulfill2.DocumentStatus = "active"
	woFulfill2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"priority":     woShipment.Priority,
		"truck_number": "B1234ASD",
		"note":         "Ubah Note",
		"workorder_shipment_items": []tester.D{
			{
				"workorder_fulfillment": common.Encrypt(woFulfill1.ID),
			},
			{
				"workorder_fulfillment": common.Encrypt(woFulfill2.ID),
			},
		},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/shipment/"+common.Encrypt(woShipment.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMappingPUTCancelShipment1 Sukses  update SalesOrder ShipmentStatus menjadi finish
func TestHandler_URLMappingPUTCancelShipment1(t *testing.T) {
	test.DataCleanUp("sales_order", "sales_order_item", "workorder_shipment", "workorder_shipment_item", "workorder_fulfillment", "workorder_fulfillment_item")
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

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
	dmWoF1.IsDeleted = int8(0)
	dmWoF1.IsDelivered = int8(0)
	dmWoF1.IsDeleted = int8(0)
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
	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/shipment/" + common.Encrypt(dmWoS.ID) + "/cancel/" + common.Encrypt(dmWoF1.ID), "PUT", 200},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, ep := range routers {
		ng.Method = ep.method
		ng.Path = ep.endpoint
		ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, ep.expected, res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", ep.endpoint, ep.method))
		})
	}
}

// TestHandler_URLMappingPUTCancelShipment2 Sukses  update SalesOrder ShipmentStatus menjadi active
func TestHandler_URLMappingPUTCancelShipment2(t *testing.T) {
	test.DataCleanUp("sales_order", "sales_order_item", "workorder_shipment", "workorder_shipment_item", "workorder_fulfillment", "workorder_fulfillment_item")

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

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
	dmWoF1.IsDeleted = int8(0)
	dmWoF1.IsDelivered = int8(0)
	dmWoF1.IsDeleted = int8(0)
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

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/shipment/" + common.Encrypt(dmWoS.ID) + "/cancel/" + common.Encrypt(dmWoF1.ID), "PUT", 200},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, ep := range routers {
		ng.Method = ep.method
		ng.Path = ep.endpoint
		ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, ep.expected, res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", ep.endpoint, ep.method))
		})
	}
}

// TestHandler_URLMappingPUTCancelShipment3 Sukses update is delete pada fulfillment
func TestHandler_URLMappingPUTCancelShipment3(t *testing.T) {
	test.DataCleanUp("sales_order", "sales_order_item", "workorder_shipment", "workorder_shipment_item", "workorder_fulfillment", "workorder_fulfillment_item")

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	dmSO := model.DummySalesOrder()
	dmSOI := model.DummySalesOrderItem()
	dmWoS := model.DummyWorkorderShipment()
	dmWoSI1 := model.DummyWorkorderShipmentItem()
	dmWoF1 := model.DummyWorkorderFulfillment()
	dmWoFI1 := model.DummyWorkorderFulfillmentItem()

	dmWoS.IsDeleted = int8(0)
	dmWoS.Save()

	dmSOI.SalesOrder = dmSO
	dmSOI.Quantity = float32(20)
	dmSOI.Save()

	dmWoF1.SalesOrder = dmSO
	dmWoF1.IsDeleted = int8(0)
	dmWoF1.IsDelivered = int8(1)
	dmWoF1.IsDeleted = int8(0)
	dmWoF1.Save()

	dmWoFI1.WorkorderFulfillment = dmWoF1
	dmWoFI1.Quantity = float32(5)
	dmWoFI1.Save()

	dmWoSI1.WorkorderShipment = dmWoS
	dmWoSI1.WorkorderFulfillment = dmWoF1
	dmWoSI1.Save()

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/shipment/" + common.Encrypt(dmWoS.ID) + "/cancel/" + common.Encrypt(dmWoF1.ID), "PUT", 200},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, ep := range routers {
		ng.Method = ep.method
		ng.Path = ep.endpoint
		ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, ep.expected, res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", ep.endpoint, ep.method))
		})
	}
}

// TestHandler_URLMappingPUTCancelShipment4 Gagal id pada workordershipment salah
func TestHandler_URLMappingPUTCancelShipment4(t *testing.T) {
	test.DataCleanUp("sales_order", "sales_order_item", "workorder_shipment", "workorder_shipment_item", "workorder_fulfillment", "workorder_fulfillment_item")

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	dmSO := model.DummySalesOrder()
	dmSOI := model.DummySalesOrderItem()
	dmWoS := model.DummyWorkorderShipment()
	dmWoSI1 := model.DummyWorkorderShipmentItem()
	dmWoF1 := model.DummyWorkorderFulfillment()
	dmWoFI1 := model.DummyWorkorderFulfillmentItem()

	dmWoS.IsDeleted = int8(0)
	dmWoS.Save()

	dmSOI.SalesOrder = dmSO
	dmSOI.Quantity = float32(20)
	dmSOI.Save()

	dmWoF1.SalesOrder = dmSO
	dmWoF1.IsDelivered = int8(1)
	dmWoF1.IsDeleted = int8(0)
	dmWoF1.Save()

	dmWoFI1.WorkorderFulfillment = dmWoF1
	dmWoF1.IsDeleted = int8(0)
	dmWoFI1.Quantity = float32(5)
	dmWoFI1.Save()

	dmWoSI1.WorkorderShipment = dmWoS
	dmWoSI1.WorkorderFulfillment = dmWoF1
	dmWoSI1.Save()

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/shipment/" + common.Encrypt(99999) + "/cancel/" + common.Encrypt(dmWoF1.ID), "PUT", http.StatusNotFound},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, ep := range routers {
		ng.Method = ep.method
		ng.Path = ep.endpoint
		ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, ep.expected, res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", ep.endpoint, ep.method))
		})
	}
}

// TestHandler_URLMappingPUTCancelShipment5 Gagal id pada workorder_fulfillment salah
func TestHandler_URLMappingPUTCancelShipment5(t *testing.T) {
	test.DataCleanUp("sales_order", "sales_order_item", "workorder_shipment", "workorder_shipment_item", "workorder_fulfillment", "workorder_fulfillment_item")

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	dmSO := model.DummySalesOrder()
	dmSOI := model.DummySalesOrderItem()
	dmWoS := model.DummyWorkorderShipment()
	dmWoSI1 := model.DummyWorkorderShipmentItem()
	dmWoF1 := model.DummyWorkorderFulfillment()
	dmWoFI1 := model.DummyWorkorderFulfillmentItem()

	dmWoS.IsDeleted = int8(0)
	dmWoS.Save()

	dmSOI.SalesOrder = dmSO
	dmSOI.Quantity = float32(20)
	dmSOI.Save()

	dmWoF1.SalesOrder = dmSO
	dmWoF1.IsDeleted = int8(0)
	dmWoF1.IsDelivered = int8(1)
	dmWoF1.IsDeleted = int8(0)
	dmWoF1.Save()

	dmWoFI1.WorkorderFulfillment = dmWoF1
	dmWoFI1.Quantity = float32(5)
	dmWoFI1.Save()

	dmWoSI1.WorkorderShipment = dmWoS
	dmWoSI1.WorkorderFulfillment = dmWoF1
	dmWoSI1.Save()

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/shipment/" + common.Encrypt(dmWoS.ID) + "/cancel/" + common.Encrypt("9090"), "PUT", http.StatusNotFound},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, ep := range routers {
		ng.Method = ep.method
		ng.Path = ep.endpoint
		ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, ep.expected, res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", ep.endpoint, ep.method))
		})
	}
}

// TestHandler_URLMappingPUTApproveShipment1 Berhasil Shipment Status Sales Order menjadi finished
func TestHandler_URLMappingPUTApproveShipment1(t *testing.T) {
	test.DataCleanUp("sales_order", "sales_order_item", "workorder_shipment", "workorder_shipment_item", "workorder_fulfillment", "workorder_fulfillment_item")

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

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
	dmWoFI1.Quantity = float32(7)
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
	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/shipment/" + common.Encrypt(dmWoS.ID) + "/approve", "PUT", 200},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, ep := range routers {
		ng.Method = ep.method
		ng.Path = ep.endpoint
		ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, ep.expected, res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", ep.endpoint, ep.method))
		})
	}
}

// TestHandler_URLMappingPUTApproveShipment1 Berhasil Shipment Status Sales Order menjadi active
func TestHandler_URLMappingPUTApproveShipment2(t *testing.T) {
	test.DataCleanUp("sales_order", "sales_order_item", "workorder_shipment", "workorder_shipment_item", "workorder_fulfillment", "workorder_fulfillment_item")

	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

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
	dmWoFI1.Quantity = float32(7)
	dmWoFI1.Save()

	dmWoFI2.WorkorderFulfillment = dmWoF2
	dmWoFI2.Quantity = float32(1)
	dmWoFI2.Save()

	dmWoSI1.WorkorderShipment = dmWoS
	dmWoSI1.WorkorderFulfillment = dmWoF1
	dmWoSI1.Save()

	dmWoSI2.WorkorderShipment = dmWoS
	dmWoSI2.WorkorderFulfillment = dmWoF2
	dmWoSI2.Save()
	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/shipment/" + common.Encrypt(dmWoS.ID) + "/approve", "PUT", 200},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, ep := range routers {
		ng.Method = ep.method
		ng.Path = ep.endpoint
		ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, ep.expected, res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", ep.endpoint, ep.method))
		})
	}
}
