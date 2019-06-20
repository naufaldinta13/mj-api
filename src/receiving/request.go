// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receiving

import (
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create receiving process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	SessionData     *auth.SessionData `json:"-"`
	PurchaseOrder   string            `json:"purchase_order" valid:"required"`
	RecognitionDate time.Time         `json:"recognition_date" valid:"required"`
	Pic             string            `json:"pic" valid:"required"`
	Note            string            `json:"note"`
	ReceivingItem   []receivingItem   `json:"work_order_receiving_items" valid:"required"`
}

type receivingItem struct {
	PurchaseOrderItem string  `json:"purchase_order_item" valid:"required"`
	Quantity          float32 `json:"quantity" valid:"required|gt:0"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if id, e := common.Decrypt(r.PurchaseOrder); e != nil {
		o.Failure("purchase_order", "cant be decrypt")
	} else {
		po := &model.PurchaseOrder{ID: id, IsDeleted: int8(0)}
		if e = po.Read("ID", "IsDeleted"); e != nil {
			o.Failure("purchase_order", "doesnt exist")
		} else {
			if po.DocumentStatus == "cancelled" && po.ReceivingStatus == "finished" {
				o.Failure("purchase_order", "invalid purchase order")
			}
		}
	}

	for _, ax := range r.ReceivingItem {

		if idPx, e := common.Decrypt(ax.PurchaseOrderItem); e != nil {
			o.Failure("purchase_order_item", "cant be decrypt")
		} else {
			poi := &model.PurchaseOrderItem{ID: idPx}
			if e = poi.Read("ID"); e != nil {
				o.Failure("purchase_order_item", "doesn't exist")
			} else {
				quantityRI, _ := TotalQuantityRI(idPx)
				quantityRI = quantityRI + ax.Quantity
				if quantityRI > poi.Quantity {
					o.Failure("quantity", "quantity receiving item greater than quantity purchase order item")
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
func (r *createRequest) Transform() *model.WorkorderReceiving {

	poID, _ := common.Decrypt(r.PurchaseOrder)
	code, _ := util.CodeGen("code_receiving", "workorder_receiving")

	var item []*model.WorkorderReceivingItem

	for _, rx := range r.ReceivingItem {
		id, _ := common.Decrypt(rx.PurchaseOrderItem)
		rItem := &model.WorkorderReceivingItem{
			PurchaseOrderItem: &model.PurchaseOrderItem{ID: id},
			Quantity:          rx.Quantity,
		}
		item = append(item, rItem)
	}

	wr := &model.WorkorderReceiving{
		Code:                    code,
		RecognitionDate:         r.RecognitionDate,
		PurchaseOrder:           &model.PurchaseOrder{ID: poID},
		Pic:                     r.Pic,
		Note:                    r.Note,
		CreatedBy:               r.SessionData.User,
		CreatedAt:               time.Now(),
		DocumentStatus:          "finished",
		WorkorderReceivingItems: item,
	}

	return wr
}
