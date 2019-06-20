package salesReturn

import (
	"fmt"
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/sales"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create setting process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	SalesOrder      string            `json:"sales_order"`
	RecognitionDate time.Time         `json:"recognition_date" valid:"required"`
	Note            string            `json:"note"`
	ReturnItem      []returnItem      `json:"sales_return_items" valid:"required"`
	SessionData     *auth.SessionData `json:"-"`
}

type returnItem struct {
	SalesOrderItem string  `json:"sales_order_item" valid:"required"`
	Quantity       float32 `json:"quantity" valid:"required|gt:0"`
	Note           string  `json:"note"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if idx, e := common.Decrypt(r.SalesOrder); e == nil {
		so := &model.SalesOrder{ID: idx}
		if e = so.Read("ID"); e != nil {
			o.Failure("sales_order", "sales order doesn't exist")
		} else {
			wf := &model.WorkorderFulfillment{SalesOrder: so}
			if erx := wf.Read("SalesOrder"); erx != nil {
				o.Failure("sales_order", "fulfillment doesn't exist")
			} else {
				CanBeReturn, _ := sales.CanBeReturnSales(so, nil)

				// check item return
				for i, item := range r.ReturnItem {
					if id, e := common.Decrypt(item.SalesOrderItem); e == nil {
						soi := &model.SalesOrderItem{ID: id}
						if e = soi.Read("ID"); e != nil {
							o.Failure(fmt.Sprintf("sales_return_items.%d.sales_order_item.invalid", i), "sales order item id doesn't exist")
						} else {
							if soi.SalesOrder.ID != so.ID {
								o.Failure(fmt.Sprintf("sales_return_items.%d.sales_order_item.invalid", i), "sales order item must from the sales order")
							}

							for _, itemCanBeReturn := range CanBeReturn {
								if itemCanBeReturn.ID == id {
									if item.Quantity > itemCanBeReturn.CanBeReturn {
										o.Failure(fmt.Sprintf("sales_return_items.%d.quantity.invalid", i), "quantity cant be greater than quantity item can be return")

									}
								}
							}
						}
					} else {
						o.Failure(fmt.Sprintf("sales_return_items.%d.sales_order_item.invalid", i), "sales order item id cannot be decrypt")
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
func (r *createRequest) Transform() *model.SalesReturn {

	var item []*model.SalesReturnItem
	sid, _ := common.Decrypt(r.SalesOrder)

	var up float64
	var subTotal float64
	var ta float64
	for _, srItem := range r.ReturnItem {
		id, _ := common.Decrypt(srItem.SalesOrderItem)
		soi := &model.SalesOrderItem{ID: id}
		soi.Read("ID")
		qtyPres := common.FloatPrecision(float64(srItem.Quantity), 2)
		up = soi.UnitPrice
		subTotal = up * float64(srItem.Quantity)
		ta = ta + subTotal
		rItem := &model.SalesReturnItem{
			SalesOrderItem: soi,
			Note:           srItem.Note,
			Quantity:       float32(qtyPres),
		}

		item = append(item, rItem)
	}

	code, _ := util.CodeGen("code_sales_return", "sales_return")
	srx := &model.SalesReturn{

		SalesOrder:       &model.SalesOrder{ID: sid},
		RecognitionDate:  r.RecognitionDate,
		Code:             code,
		TotalAmount:      ta,
		Note:             r.Note,
		DocumentStatus:   "new",
		IsDeleted:        int8(0),
		CreatedBy:        &model.User{ID: r.SessionData.User.ID},
		CreatedAt:        time.Now(),
		SalesReturnItems: item,
	}

	return srx
}

// updateRequest data struct that stored request data when requesting an create setting process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type updateRequest struct {
	SR              *model.SalesReturn
	SO              *model.SalesOrder
	RecognitionDate time.Time          `json:"recognition_date" valid:"required"`
	Note            string             `json:"note"`
	ReturnItem      []returnItemUpdate `json:"sales_return_items" valid:"required"`
	SessionData     *auth.SessionData  `json:"-"`
}

type returnItemUpdate struct {
	ID             string  `json:"id"`
	SalesOrderItem string  `json:"sales_order_item" valid:"required"`
	Quantity       float32 `json:"quantity" valid:"required|gt:0"`
	Note           string  `json:"note"`
}

// Validate implement validation.Requests interfaces.
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if r.SR.DocumentStatus == "cancelled" || r.SR.DocumentStatus == "finished" {
		o.Failure("document_status", "can't update document")
	}

	CanBeReturn, _ := sales.CanBeReturnSales(r.SO, r.SR)
	// check item return
	for i, item := range r.ReturnItem {
		if slsID, e := common.Decrypt(item.SalesOrderItem); e == nil {
			soi := &model.SalesOrderItem{ID: slsID}
			if e = soi.Read("ID"); e != nil {
				o.Failure(fmt.Sprintf("sales_return_items.%d.sales_order_item.invalid", i), "sales order item id doesn't exist")
			} else {
				if soi.SalesOrder.ID != r.SO.ID {
					o.Failure(fmt.Sprintf("sales_return_items.%d.sales_order_item.invalid", i), "sales order item must from the sales order")
				}

				for _, itemCanBeReturn := range CanBeReturn {
					if itemCanBeReturn.ID == slsID {
						if item.Quantity > itemCanBeReturn.CanBeReturn {
							o.Failure(fmt.Sprintf("sales_return_items.%d.quantity.invalid", i), "quantity cant be greater than quantity item can be return")

						}
					}
				}
			}
		} else {
			o.Failure(fmt.Sprintf("sales_return_items.%d.sales_order_item.invalid", i), "sales order item id cannot be decrypt")
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
func (r *updateRequest) Transform(sr *model.SalesReturn) *model.SalesReturn {

	var item []*model.SalesReturnItem

	var up float64
	var subTotal float64
	var ta float64

	for _, srItem := range r.ReturnItem {

		id, _ := common.Decrypt(srItem.SalesOrderItem)
		soi := &model.SalesOrderItem{ID: id}
		soi.Read("ID")
		qtyPres := common.FloatPrecision(float64(srItem.Quantity), 2)
		up = soi.UnitPrice
		subTotal = up * float64(srItem.Quantity)
		ta = ta + subTotal

		rItem := &model.SalesReturnItem{
			SalesOrderItem: soi,
			Note:           srItem.Note,
			Quantity:       float32(qtyPres),
			SalesReturn:    &model.SalesReturn{ID: r.SR.ID},
		}

		// if id of sales return item exist
		if srItem.ID != "" {
			itmID, _ := common.Decrypt(srItem.ID)
			rItem.ID = itmID
		}

		item = append(item, rItem)
	}

	srx := &model.SalesReturn{
		ID:               sr.ID,
		RecognitionDate:  r.RecognitionDate,
		TotalAmount:      ta,
		Note:             r.Note,
		DocumentStatus:   "new",
		UpdatedBy:        &model.User{ID: r.SessionData.User.ID},
		UpdatedAt:        time.Now(),
		SalesReturnItems: item,
	}
	return srx
}

// cancelRequest data struct that stored request data when requesting an create setting process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type cancelRequest struct {
	SR *model.SalesReturn
}

// Validate implement validation.Requests interfaces.
func (r *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	if r.SR.DocumentStatus == "cancelled" || r.SR.DocumentStatus == "finished" {
		o.Valid = false
		o.Failure("document_status", "can't cancel document")
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *cancelRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *cancelRequest) Transform() {
}
