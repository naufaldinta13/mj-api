// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fulfillment

import (
	"fmt"
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create fulfillment process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	SalesOrderID              string                     `json:"sales_order_id" valid:"required"`
	Priority                  string                     `json:"priority" valid:"required|in:routine,rush,emergency"`
	DueDate                   time.Time                  `json:"due_date" valid:"required"`
	ShippingAddress           string                     `json:"shipping_address" valid:"required"`
	Note                      string                     `json:"note"`
	WorkorderFulFillmentItems []workorderFulfillmentItem `json:"workorder_fulfillment_items" valid:"required"`
}

type workorderFulfillmentItem struct {
	ID               string  `json:"id"`
	SalesOrderItemID string  `json:"sales_order_item_id" valid:"required"`
	Quantity         float32 `json:"quantity" valid:"required|gt:0"`
	Note             string  `json:"note"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	var soID int64
	var so *model.SalesOrder
	// validasi sales order id
	if soID, e = common.Decrypt(r.SalesOrderID); e != nil {
		o.Failure("sales_order_id", "sales_order_id not valid")
	} else {
		so = &model.SalesOrder{ID: soID}
		if e = so.Read(); e != nil {
			o.Failure("sales_order_id", "sales_order_id doesn't exist")
		} else if e == nil && (so.FulfillmentStatus == "finished" || so.DocumentStatus == "approved_cancel") {
			o.Failure("sales_order_id", "invalid sales order")
		}

		var items []int64
		for _, i := range r.WorkorderFulFillmentItems {
			// validasi sales order item id
			var soitemID int64
			var soitem *model.SalesOrderItem
			if soitemID, e = common.Decrypt(i.SalesOrderItemID); e != nil {
				o.Failure(fmt.Sprintf("workorder_fulfillment_items.%s.sales_order_item_id.invalid", i.SalesOrderItemID), "sales_order_item_id not valid")
			} else {
				soitem = &model.SalesOrderItem{ID: soitemID}
				if e = soitem.Read(); e != nil {
					o.Failure(fmt.Sprintf("workorder_fulfillment_items.%s.sales_order_item_id.invalid", i.SalesOrderItemID), "sales_order_item_id doesn't exist")
				} else if e == nil && soitem.SalesOrder.ID != so.ID {
					o.Failure(fmt.Sprintf("workorder_fulfillment_items.%s.sales_order_item_id.invalid", i.SalesOrderItemID), "sales_order_item_id is not part of sales_order")
				}
				if util.HasElem(items, soitemID) {
					o.Failure(fmt.Sprintf("workorder_fulfillment_items.%s.sales_order_item_id.invalid", i.SalesOrderItemID), "sales_order_item_id cannot same")
				}
				items = append(items, soitemID)

				quantityFul, _ := getSumQuantityFullfillmentItemBySoitemID(soitem, 0)
				if i.Quantity+quantityFul > soitem.Quantity {
					o.Failure(fmt.Sprintf("workorder_fulfillment_items.%s.sales_order_item_id.invalid", i.SalesOrderItemID), "quantity cannot bigger than soitem quantity")
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
func (r *createRequest) Transform(user *model.User) *model.WorkorderFulfillment {
	soID, _ := common.Decrypt(r.SalesOrderID)
	code, _ := util.CodeGen("code_fullfilment", "workorder_fulfillment")

	var items []*model.WorkorderFulfillmentItem
	for _, i := range r.WorkorderFulFillmentItems {
		soitemID, _ := common.Decrypt(i.SalesOrderItemID)
		item := &model.WorkorderFulfillmentItem{
			SalesOrderItem: &model.SalesOrderItem{ID: soitemID},
			Quantity:       i.Quantity,
			Note:           i.Note,
		}
		items = append(items, item)
	}

	fulfillment := &model.WorkorderFulfillment{
		SalesOrder:                &model.SalesOrder{ID: soID},
		Code:                      code,
		Priority:                  r.Priority,
		DueDate:                   r.DueDate,
		ShippingAddress:           r.ShippingAddress,
		Note:                      r.Note,
		DocumentStatus:            "new",
		CreatedBy:                 user,
		CreatedAt:                 time.Now(),
		WorkorderFulFillmentItems: items,
	}
	return fulfillment
}

// updateRequest data struct that stored request data when requesting an update fulfillment process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type updateRequest struct {
	Priority                  string                      `json:"priority" valid:"required|in:routine,rush,emergency"`
	DueDate                   time.Time                   `json:"due_date" valid:"required"`
	ShippingAddress           string                      `json:"shipping_address" valid:"required"`
	Note                      string                      `json:"note"`
	WorkorderFulFillmentItems []workorderFulfillmentItem  `json:"workorder_fulfillment_items" valid:"required"`
	Fulfillment               *model.WorkorderFulfillment `json:"-"`
}

// Validate implement validation.Requests interfaces.
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	if r.Fulfillment.DocumentStatus != "new" {
		o.Failure("document_status", "document_status should be new")
	}

	var e error
	var items []int64
	for a, i := range r.WorkorderFulFillmentItems {
		// validasi id fulfillment item
		var id int64
		if i.ID != "" {
			if id, e = common.Decrypt(i.ID); e != nil {
				o.Failure(fmt.Sprintf("workorder_fulfillment_items.%d.id.invalid", a), "id not valid")
			} else {
				item := &model.WorkorderFulfillmentItem{ID: id}
				if e = item.Read(); e != nil {
					o.Failure(fmt.Sprintf("workorder_fulfillment_items.%d.id.invalid", a), "id doesn't exist")
				}
			}
		}

		// validasi sales order item id
		var soitemID int64
		var soitem *model.SalesOrderItem
		if soitemID, e = common.Decrypt(i.SalesOrderItemID); e != nil {
			o.Failure(fmt.Sprintf("workorder_fulfillment_items.%d.sales_order_item_id.invalid", a), "sales_order_item_id not valid")
		} else {
			soitem = &model.SalesOrderItem{ID: soitemID}
			if e = soitem.Read(); e != nil {
				o.Failure(fmt.Sprintf("workorder_fulfillment_items.%d.sales_order_item_id.invalid", a), "sales_order_item_id doesn't exist")
			} else if e == nil && soitem.SalesOrder.ID != r.Fulfillment.SalesOrder.ID {
				o.Failure(fmt.Sprintf("workorder_fulfillment_items.%d.sales_order_item_id.invalid", a), "sales_order_item_id is not part of sales_order")
			}
			if util.HasElem(items, soitemID) {
				o.Failure(fmt.Sprintf("workorder_fulfillment_items.%d.sales_order_item_id.invalid", a), "sales_order_item_id cannot same")
			}
			items = append(items, soitemID)

			quantityFul, _ := getSumQuantityFullfillmentItemBySoitemID(soitem, id)
			if i.Quantity+quantityFul > soitem.Quantity {
				o.Failure(fmt.Sprintf("workorder_fulfillment_items.%d.sales_order_item_id.invalid", a), "quantity cannot bigger than soitem quantity")
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
func (r *updateRequest) Transform(user *model.User) (*model.WorkorderFulfillment, []*model.WorkorderFulfillmentItem) {
	fulfillment := r.Fulfillment
	fulfillment.Priority = r.Priority
	fulfillment.DueDate = r.DueDate
	fulfillment.ShippingAddress = r.ShippingAddress
	fulfillment.Note = r.Note

	// item dari inputan
	var items []*model.WorkorderFulfillmentItem
	for _, i := range r.WorkorderFulFillmentItems {
		var ID int64
		if i.ID != "" {
			ID, _ = common.Decrypt(i.ID)
		}
		soID, _ := common.Decrypt(i.SalesOrderItemID)

		item := &model.WorkorderFulfillmentItem{
			ID:             ID,
			SalesOrderItem: &model.SalesOrderItem{ID: soID},
			Quantity:       i.Quantity,
			Note:           i.Note,
		}
		items = append(items, item)
	}

	return fulfillment, items
}

// approveRequest data struct that stored request data when requesting an approve fulfillment process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type approveRequest struct {
	Fulfillment *model.WorkorderFulfillment `json:"-"`
}

// Validate implement validation.Requests interfaces.
func (r *approveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var checkVariant = make(map[int64]*model.ItemVariant)

	if r.Fulfillment.DocumentStatus == "finished" {
		o.Failure("document_status", "status document has finished")
	}

	for _, i := range r.Fulfillment.WorkorderFulFillmentItems {
		ivar := &model.ItemVariant{ID: i.SalesOrderItem.ItemVariant.ID}
		ivar.Read("ID")
		////////////////////////////////
		if checkVariant[ivar.ID] == nil {
			ivar.CommitedStock -= i.Quantity
			checkVariant[ivar.ID] = ivar
		} else {
			variant := checkVariant[ivar.ID]
			variant.CommitedStock -= i.Quantity
			checkVariant[ivar.ID] = variant
		}
		////////////////////////////////
	}

	for _, i := range r.Fulfillment.WorkorderFulFillmentItems {

		if i.Quantity > (checkVariant[i.SalesOrderItem.ItemVariant.ID].AvailableStock - checkVariant[i.SalesOrderItem.ItemVariant.ID].CommitedStock) {
			o.Failure("id", "quantity fulfillment item must be lower than stock")
		}
		checkVariant[i.SalesOrderItem.ItemVariant.ID].CommitedStock += i.Quantity
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *approveRequest) Messages() map[string]string {
	return map[string]string{}
}
