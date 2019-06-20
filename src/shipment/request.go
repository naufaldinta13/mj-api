// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package shipment

import (
	"fmt"
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create shipment process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	Session        *auth.SessionData
	Priority       string          `json:"priority" valid:"required|in:routine,rush,emergency"`
	TruckNumber    string          `json:"truck_number" valid:"required"`
	WOShipmentItem []WOFulfillment `json:"workorder_shipment_items" valid:"required"`
	Note           string          `json:"note"`
}

// WOFulfillment contain workOrder fulfillment id
type WOFulfillment struct {
	WOFulFillmentID string `json:"workorder_fulfillment" valid:"required"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	fIds := make(map[string]bool)
	// loop wo shipment item
	for i, woShipment := range r.WOShipmentItem {
		// check id wo fulfillment
		if id, e := common.Decrypt(woShipment.WOFulFillmentID); e != nil {
			o.Failure(fmt.Sprintf("workorder_shipment_item.%d.workorder_fulfillment.invalid", i), "not valid")
		} else {
			//check duplicate workorder_fulfillment id
			if fIds[woShipment.WOFulFillmentID] == true {
				o.Failure(fmt.Sprintf("workorder_shipment_item.%d.workorder_fulfillment.invalid", i), "workorder_fulfillment id duplicate")
			} else {
				fIds[woShipment.WOFulFillmentID] = true

				woFulfill := &model.WorkorderFulfillment{ID: id, IsDeleted: int8(0)}
				if e = woFulfill.Read("ID", "IsDeleted"); e != nil {
					o.Failure(fmt.Sprintf("workorder_shipment_item.%d.workorder_fulfillment.invalid", i), "not valid")
				} else {
					if woFulfill.DocumentStatus != "finished" {
						o.Failure(fmt.Sprintf("workorder_shipment_item.%d.workorder_fulfillment.invalid", i), "fulfillment still not finished")
					} else if woFulfill.IsDelivered == 1 {
						o.Failure(fmt.Sprintf("workorder_shipment_item.%d.workorder_fulfillment.invalid", i), "fulfillment has been delivered")
					} else {
						// cek apakah fulfillment sudah di pakai di shipment lain
						si := &model.WorkorderShipmentItem{WorkorderFulfillment: woFulfill}
						if e = si.Read("workorder_fulfillment_id"); e == nil {
							o.Failure(fmt.Sprintf("workorder_shipment_item.%d.workorder_fulfillment.invalid", i), "fulfillment already on other shipment")
						}
					}
				}
			}
		}
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *createRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *createRequest) Transform() (model.WorkorderShipment, error) {
	var WoShipItems []*model.WorkorderShipmentItem
	var m model.WorkorderShipment
	var code string
	var e error

	if code, e = util.CodeGen("code_shipment", "workorder_shipment"); e == nil {
		// loop wo shipment item
		for _, u := range r.WOShipmentItem {
			id, _ := common.Decrypt(u.WOFulFillmentID)
			WoShipItem := &model.WorkorderShipmentItem{
				WorkorderFulfillment: &model.WorkorderFulfillment{ID: id},
			}
			WoShipItems = append(WoShipItems, WoShipItem)
		}
		// make model
		m = model.WorkorderShipment{
			Code:           code,
			Priority:       r.Priority,
			TruckNumber:    r.TruckNumber,
			DocumentStatus: "active",
			IsDeleted:      int8(0),
			CreatedAt:      time.Now(),
			CreatedBy:      r.Session.User,
			Note:           r.Note,
			WorkorderShipmentItems: WoShipItems,
		}
	}
	return m, e
}

// updateRequest data struct that stored request data when requesting an update shipment process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type updateRequest struct {
	Session             *auth.SessionData
	Priority            string          `json:"priority" valid:"required|in:routine,rush,emergency"`
	TruckNumber         string          `json:"truck_number" valid:"required"`
	WOShipmentItem      []WOFulfillment `json:"workorder_shipment_items" valid:"required"`
	Note                string          `json:"note"`
	WorkorderShipmentID int64
}

// Validate implement validation.Requests interfaces.
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	fIds := make(map[string]bool)
	// loop wo shipment item
	for i, woShipment := range r.WOShipmentItem {
		// check id wo fulfillment
		if id, e := common.Decrypt(woShipment.WOFulFillmentID); e != nil {
			o.Failure(fmt.Sprintf("workorder_shipment_item.%d.workorder_fulfillment.invalid", i), "not valid")
		} else {
			//check duplicate workorder_fulfillment id
			if fIds[woShipment.WOFulFillmentID] == true {
				o.Failure(fmt.Sprintf("workorder_shipment_item.%d.workorder_fulfillment.invalid", i), "workorder_fulfillment id duplicate")
			} else {
				fIds[woShipment.WOFulFillmentID] = true

				woFulfill := &model.WorkorderFulfillment{ID: id, IsDeleted: int8(0)}
				if e = woFulfill.Read("ID", "IsDeleted"); e != nil {
					o.Failure(fmt.Sprintf("workorder_shipment_item.%d.workorder_fulfillment.invalid", i), "not valid")
				} else {
					if woFulfill.DocumentStatus != "finished" {
						o.Failure(fmt.Sprintf("workorder_shipment_item.%d.workorder_fulfillment.invalid", i), "fulfillment still not finished")
					} else if woFulfill.IsDelivered == 1 {
						o.Failure(fmt.Sprintf("workorder_shipment_item.%d.workorder_fulfillment.invalid", i), "fulfillment has been delivered")
					}
				}
			}
		}
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *updateRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *updateRequest) Transform(shipment *model.WorkorderShipment) (*model.WorkorderShipment, error) {
	var e error

	// loop wo shipment item
	for _, u := range r.WOShipmentItem {
		id, _ := common.Decrypt(u.WOFulFillmentID)
		WoShipItem := &model.WorkorderShipmentItem{
			WorkorderFulfillment: &model.WorkorderFulfillment{ID: id},
		}
		shipment.WorkorderShipmentItems = append(shipment.WorkorderShipmentItems, WoShipItem)
	}
	// make model
	shipment.Priority = r.Priority
	shipment.TruckNumber = r.TruckNumber
	shipment.Note = r.Note

	return shipment, e
}
