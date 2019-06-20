// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchaseReturn

import (
	"fmt"
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/purchase"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create purchaseReturn process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	Code                string               `json:"code"`
	RecognitionDate     time.Time            `json:"recognition_date" valid:"required"`
	PurchaseOrderID     string               `json:"purchase_order_id" valid:"required"`
	PurchaseReturnItems []ItemPurchaseReturn `json:"purchase_return_item" valid:"required"`
	TotalAmount         float64
	Note                string `json:"note"`
	Session             *auth.SessionData
}

// ItemPurchaseReturn Digunakan untuk create request
type ItemPurchaseReturn struct {
	PurchaseOrderItemID string  `json:"purchase_order_item_id" valid:"required"`
	Quantity            float32 `json:"quantity" valid:"gt:0|lte:100"`
	Note                string  `json:"note"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var currentAmount float64
	var porderID, porderitemID int64
	var porder *model.PurchaseOrder

	// insert generated code
	r.Code, _ = util.CodeGen("code_purchase_return", "purchase_return")

	// check purchase order id
	if porderID, err = common.Decrypt(r.PurchaseOrderID); err != nil {
		o.Failure("purchase_order_id", "Purchase order id is invalid")
	} else {
		if porder, err = purchase.GetDetailPurchaseOrder("id", porderID); err != nil {
			o.Failure("purchase_order_id", "Purchase order id is not found")
		} else {
			// ambil can be return setiap purchase order item
			porderitems := purchase.ReturningPurchaseOrder(porder, 0)
			// untuk mengecek duplicate poi id
			poiID := make(map[string]bool)
			// check validasi purchase return item
			for i, row := range r.PurchaseReturnItems {
				// check if puchase order item id decrypt able
				if porderitemID, err = common.Decrypt(row.PurchaseOrderItemID); err != nil {
					o.Failure(fmt.Sprintf("purchase_return_item.%d.purchase_order_item_id.invalid", i), "Purchase order item id is invalid")
				} else {
					orderItem := model.PurchaseOrderItem{ID: porderitemID}
					if err = orderItem.Read("ID"); err != nil {
						o.Failure(fmt.Sprintf("purchase_return_item.%d.purchase_order_item_id.invalid", i), "Purchase order item id is not found")
					} else {
						// cek duplicate poi id
						if poiID[row.PurchaseOrderItemID] == true {
							o.Failure(fmt.Sprintf("purchase_return_item.%d.purchase_order_item_id.invalid", i), "Purchase order item id is duplicate")
						} else {
							poiID[row.PurchaseOrderItemID] = true
						}

						// cek apakah po id pada poi sama
						if orderItem.PurchaseOrder.ID != porder.ID {
							o.Failure(fmt.Sprintf("purchase_return_item.%d.purchase_order_item_id.invalid", i), "Purchase order item id is invalid")
						}

						// cek apakah qty yg mau di return tidak melebihi yang bisa direturn
						for _, itemPOI := range *porderitems {
							if itemPOI.ID == orderItem.ID {
								if row.Quantity > itemPOI.CanBeReturn {
									o.Failure(fmt.Sprintf("purchase_return_item.%d.quantity", i), "Quantity inputted is too large")
								}
								// menghitung sum harga
								current := common.FloatPrecision(float64(row.Quantity)*itemPOI.UnitPrice, 0)
								discAmount := common.FloatPrecision((current*float64(itemPOI.Discount))/float64(100), 0)
								current = current - discAmount
								currentAmount += current
							}
						}
					}
				}
			}
		}
	}
	// insert total amount
	// Warning!! Perhitungan Diskon Tidak ada
	r.TotalAmount = currentAmount

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *createRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *createRequest) Transform() *model.PurchaseReturn {
	var preturn model.PurchaseReturn
	var pReturnItem []*model.PurchaseReturnItem
	var decryptID int64

	for _, row := range r.PurchaseReturnItems {
		decryptID, _ = common.Decrypt(row.PurchaseOrderItemID)
		pritem := model.PurchaseReturnItem{
			PurchaseOrderItem: &model.PurchaseOrderItem{ID: decryptID},
			Quantity:          row.Quantity,
			Note:              row.Note,
		}
		pReturnItem = append(pReturnItem, &pritem)
	}

	decryptID, _ = common.Decrypt(r.PurchaseOrderID)
	preturn = model.PurchaseReturn{
		PurchaseOrder:       &model.PurchaseOrder{ID: decryptID},
		RecognitionDate:     r.RecognitionDate,
		Code:                r.Code,
		TotalAmount:         r.TotalAmount,
		Note:                r.Note,
		IsDeleted:           int8(0),
		CreatedBy:           r.Session.User,
		CreatedAt:           time.Now(),
		PurchaseReturnItems: pReturnItem,
		DocumentStatus:      "new",
	}

	return &preturn
}

// updateRequest data struct that stored request data when requesting an update purchaseReturn process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type updateRequest struct {
	RecognitionDate     time.Time                  `json:"recognition_date" valid:"required"`
	PurchaseOrderID     string                     `json:"purchase_order_id" valid:"required"`
	PurchaseReturnItems []UpdateItemPurchaseReturn `json:"purchase_return_item" valid:"required"`
	TotalAmount         float64
	Note                string `json:"note"`
	Session             *auth.SessionData
	PurchaseReturn      *model.PurchaseReturn
}

// UpdateItemPurchaseReturn Digunakan untuk update request
type UpdateItemPurchaseReturn struct {
	PurchaseOrderItemID string  `json:"purchase_order_item_id" valid:"required"`
	Quantity            float32 `json:"quantity" valid:"gt:0|lte:100"`
	Note                string  `json:"note"`
}

// Validate implement validation.Requests interfaces.
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var currentAmount float64
	var porderID, porderitemID int64
	var porder *model.PurchaseOrder

	if r.PurchaseReturn.DocumentStatus != "new" {
		o.Failure("document_status", "Document already executed by the other document")
	}

	// check purchase order id
	if porderID, err = common.Decrypt(r.PurchaseOrderID); err != nil {
		o.Failure("purchase_order_id", "Purchase order id is invalid")
	} else {
		if porder, err = purchase.GetDetailPurchaseOrder("id", porderID); err != nil {
			o.Failure("purchase_order_id", "Purchase order id is not found")
		} else {
			// ambil can be return setiap purchase order item
			porderitems := purchase.ReturningPurchaseOrder(porder, r.PurchaseReturn.ID)

			// untuk mengecek duplicate poi id
			poiID := make(map[string]bool)
			for i, row := range r.PurchaseReturnItems {
				// check if puchase order item id decrypt able
				if porderitemID, err = common.Decrypt(row.PurchaseOrderItemID); err != nil {
					o.Failure(fmt.Sprintf("purchase_return_item.%d.purchase_order_item_id.invalid", i), "Purchase order item id is invalid")
				} else {
					orderItem := model.PurchaseOrderItem{ID: porderitemID}
					if err = orderItem.Read("ID"); err != nil {
						o.Failure(fmt.Sprintf("purchase_return_item.%d.purchase_order_item_id.invalid", i), "Purchase order item id is not found")
					} else {

						// cek duplicate poi id
						if poiID[row.PurchaseOrderItemID] == true {
							o.Failure(fmt.Sprintf("purchase_return_item.%d.purchase_order_item_id.invalid", i), "Purchase order item id is duplicate")
						} else {
							poiID[row.PurchaseOrderItemID] = true
						}

						// cek apakah po id pada poi sama
						if orderItem.PurchaseOrder.ID != porder.ID {
							o.Failure(fmt.Sprintf("purchase_return_item.%d.purchase_order_item_id.invalid", i), "Purchase order item id is invalid")
						}

						// cek apakah qty yg mau di return tidak melebihi yang bisa direturn
						for _, itemPOI := range *porderitems {
							if itemPOI.ID == orderItem.ID {
								if row.Quantity > itemPOI.CanBeReturn {
									o.Failure(fmt.Sprintf("purchase_return_item.%d.quantity", i), "Quantity inputted is too large")
								}
								// menghitung sum harga
								current := common.FloatPrecision(float64(row.Quantity)*itemPOI.UnitPrice, 0)
								discAmount := common.FloatPrecision((current*float64(itemPOI.Discount))/float64(100), 0)
								current = current - discAmount
								currentAmount += current
							}
						}
					}
				}
			}
		}
	}

	// insert total amount
	// Warning!! Perhitungan Diskon Tidak ada
	r.TotalAmount = currentAmount

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *updateRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *updateRequest) Transform() *model.PurchaseReturn {

	for i, row := range r.PurchaseReturnItems {
		r.PurchaseReturn.PurchaseReturnItems[i].Quantity = row.Quantity
		r.PurchaseReturn.PurchaseReturnItems[i].Note = row.Note
	}

	r.PurchaseReturn.RecognitionDate = r.RecognitionDate
	r.PurchaseReturn.TotalAmount = r.TotalAmount
	r.PurchaseReturn.Note = r.Note
	r.PurchaseReturn.IsDeleted = int8(0)
	r.PurchaseReturn.UpdatedAt = time.Now()
	r.PurchaseReturn.UpdatedBy = r.Session.User

	return r.PurchaseReturn
}
