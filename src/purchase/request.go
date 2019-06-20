// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchase

import (
	"fmt"
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/mj/api/src/inventory"
)

// createRequest data struct that stored request data when requesting an create purchase order process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	RecognitionDate    time.Time                  `json:"recognition_date" valid:"required"`
	ReferenceID        string                     `json:"reference_id"`
	AutoInvoiced       int8                       `json:"auto_invoiced"`
	EtaDate            time.Time                  `json:"eta_date" valid:"required"`
	SupplierID         string                     `json:"supplier_id" valid:"required"`
	Tax                float32                    `json:"tax" valid:"lte:100"`
	Discount           float32                    `json:"discount" valid:"gte:0|lte:100"`
	DiscountAmount     float64                    `json:"discount_amount" valid:"gte:0"`
	IsPercentage       int8                       `json:"is_percentage"`
	ShipmentCost       float64                    `json:"shipment_cost" valid:"gte:0"`
	Note               string                     `json:"note"`
	PurchaseOrderItems []purchaseOrderItemRequest `json:"purchase_order_items" valid:"required"`
}

type purchaseOrderItemRequest struct {
	ID            string  `json:"id"`
	ItemVariantID string  `json:"item_variant_id" valid:"required"`
	Quantity      float32 `json:"quantity" valid:"required|gt:0"`
	UnitPrice     float64 `json:"unit_price" valid:"required|gt:0"`
	Discount      float32 `json:"discount" valid:"lte:100"`
	Note          string  `json:"note"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	// validasi reference id
	if r.ReferenceID != "" {
		if refID, e := common.Decrypt(r.ReferenceID); e != nil {
			o.Failure("reference_id", "reference id is not valid")
		} else {
			po := &model.PurchaseOrder{ID: refID, IsDeleted: int8(0)}
			if e = po.Read(); e != nil {
				o.Failure("reference_id", "reference_id doesn't exist")
			} else {
				if po.DocumentStatus != "cancelled" {
					o.Failure("reference_id", "Document status in reference is not cancelled")
				}
			}
		}
	}
	// validasi supplier
	if supID, e := common.Decrypt(r.SupplierID); e != nil {
		o.Failure("supplier_id", "supplier_id is not valid")
	} else {
		supplier := &model.Partnership{ID: supID}
		if e = supplier.Read(); e != nil {
			o.Failure("supplier_id", "supplier_id doesn't exist")
		} else {
			if supplier.PartnershipType != "supplier" {
				o.Failure("supplier_id", "Partnership inputted supplier is not valid")
			}
		}
	}

	checkDuplicate := make(map[string]bool)
	for i, item := range r.PurchaseOrderItems {
		// validasi itemm variant
		if ivID, e := common.Decrypt(item.ItemVariantID); e != nil {
			o.Failure(fmt.Sprintf("purchase_order_items.%d.item_variant_id.invalid", i), "item_variant_id not valid")
		} else {
			iv := &model.ItemVariant{ID: ivID}
			if e = iv.Read(); e != nil {
				o.Failure(fmt.Sprintf("purchase_order_items.%d.item_variant_id.invalid", i), "item_variant_id doesn't exist")
			} else {
				if iv.IsArchived == int8(1) || iv.IsDeleted == int8(1) {
					o.Failure(fmt.Sprintf("purchase_order_items.%d.item_variant_id.invalid", i), "item_variant_id is already archived or does not exists")
				}
				//check duplicate item variant id
				if checkDuplicate[item.ItemVariantID] == true {
					o.Failure(fmt.Sprintf("sales_order_item.%d.item_variant_id.invalid", i), " item variant id duplicate")
				} else {
					checkDuplicate[item.ItemVariantID] = true
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

func (r *createRequest) Transform(user *model.User) *model.PurchaseOrder {

	var items []*model.PurchaseOrderItem
	var total float64
	for _, i := range r.PurchaseOrderItems {
		ivID, _ := common.Decrypt(i.ItemVariantID)
		// subtotal = quantity * unit_price * 1 – discount
		discountPoItem := ((i.UnitPrice * float64(i.Quantity)) * float64(i.Discount)) / float64(100)

		subtotal := (float64(i.Quantity) * i.UnitPrice) - float64(discountPoItem)

		poItem := &model.PurchaseOrderItem{
			ItemVariant: &model.ItemVariant{ID: ivID},
			Quantity:    float32(common.FloatPrecision(float64(i.Quantity), 2)),
			UnitPrice:   i.UnitPrice,
			Discount:    i.Discount,
			Subtotal:    common.FloatPrecision(subtotal, 0),
			Note:        i.Note,
		}
		items = append(items, poItem)
		total += subtotal
	}
	var po = new(model.PurchaseOrder)

	code, _ := util.CodeGen("code_purchase_order", "purchase_order")
	supplierID, _ := common.Decrypt(r.SupplierID)

	if r.ReferenceID != "" {
		refID, _ := common.Decrypt(r.ReferenceID)
		po.Reference = &model.PurchaseOrder{ID: refID}
	}

	var currentTotal, currentDiscount float64
	currentDiscount = float64(float64(r.Discount)*total) / float64(100)
	currentTotal = total - currentDiscount

	po.Code = code
	po.Supplier = &model.Partnership{ID: supplierID}
	po.RecognitionDate = r.RecognitionDate
	po.EtaDate = r.EtaDate
	po.ShipmentCost = common.FloatPrecision(r.ShipmentCost, 0)
	po.AutoInvoiced = r.AutoInvoiced
	po.Note = r.Note
	po.DocumentStatus = "new"
	po.InvoiceStatus = "new"
	po.ReceivingStatus = "new"
	po.CreatedBy = user
	po.CreatedAt = time.Now()
	po.Discount = float32(common.FloatPrecision(float64(r.Discount), 2))
	po.Tax = float32(common.FloatPrecision(float64(r.Tax), 2))
	po.TaxAmount = common.FloatPrecision(float64(currentTotal*float64(r.Tax))/float64(100), 0)
	po.IsPercentage = r.IsPercentage
	po.PurchaseOrderItems = items

	if r.IsPercentage == int8(1) {
		po.DiscountAmount = common.FloatPrecision(currentDiscount, 0)
		po.TotalCharge = common.FloatPrecision(currentTotal+po.TaxAmount+r.ShipmentCost, 0)
	} else {
		if r.DiscountAmount > float64(0) {
			po.DiscountAmount = common.FloatPrecision(r.DiscountAmount, 0)
			po.Discount = float32(common.FloatPrecision(float64(r.DiscountAmount/total)*float64(100), 0))
			currentTotal = total - r.DiscountAmount
			po.TotalCharge = common.FloatPrecision(currentTotal+po.TaxAmount+r.ShipmentCost, 0)
		} else {
			po.TotalCharge = common.FloatPrecision(total+po.TaxAmount+r.ShipmentCost, 0)
		}
	}

	if r.AutoInvoiced == int8(1) {
		po.DocumentStatus = "active"
		po.InvoiceStatus = "active"
	}

	return po
}

// updateRequest data struct that stored request data when requesting an update purchase order process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type updateRequest struct {
	RecognitionDate    time.Time                  `json:"recognition_date" valid:"required"`
	EtaDate            time.Time                  `json:"eta_date" valid:"required"`
	SupplierID         string                     `json:"supplier_id" valid:"required"`
	Tax                float32                    `json:"tax" valid:"lte:100"`
	Discount           float32                    `json:"discount" valid:"lte:100"`
	DiscountAmount     float64                    `json:"discount_amount" valid:"gte:0"`
	IsPercentage       int8                       `json:"is_percentage"`
	ShipmentCost       float64                    `json:"shipment_cost" valid:"gte:0"`
	Note               string                     `json:"note"`
	PurchaseOrderItems []purchaseOrderItemRequest `json:"purchase_order_items" valid:"required"`
	PurchaseOrder      *model.PurchaseOrder       `json:"-"`
}

// Validate implement validation.Requests interfaces.
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if r.PurchaseOrder.DocumentStatus != "new" {
		o.Failure("document_status", "document_status should be new")
	}

	// validasi supplier
	if supID, e := common.Decrypt(r.SupplierID); e != nil {
		o.Failure("supplier_id", "supplier_id is not valid")
	} else {
		supplier := &model.Partnership{ID: supID}
		if e = supplier.Read(); e != nil {
			o.Failure("supplier_id", "supplier_id doesn't exist")
		} else {
			if supplier.PartnershipType != "supplier" {
				o.Failure("supplier_id", "Partnership inputted supplier is not valid")
			}
		}
	}

	checkDuplicate := make(map[string]bool)
	for i, item := range r.PurchaseOrderItems {
		// validasi itemm variant
		if ivID, e := common.Decrypt(item.ItemVariantID); e != nil {
			o.Failure(fmt.Sprintf("purchase_order_items.%d.item_variant_id.invalid", i), "item_variant_id not valid")
		} else {
			iv := &model.ItemVariant{ID: ivID}
			if e = iv.Read(); e != nil {
				o.Failure(fmt.Sprintf("purchase_order_items.%d.item_variant_id.invalid", i), "item_variant_id doesn't exist")
			} else {
				if iv.IsArchived == int8(1) || iv.IsDeleted == int8(1) {
					o.Failure(fmt.Sprintf("purchase_order_items.%d.item_variant_id.invalid", i), "item_variant_id is already archived or does not exists")
				}
				//check duplicate item variant id
				if checkDuplicate[item.ItemVariantID] == true {
					o.Failure(fmt.Sprintf("sales_order_item.%d.item_variant_id.invalid", i), " item variant id duplicate")
				} else {
					checkDuplicate[item.ItemVariantID] = true
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

func (r *updateRequest) Transform(user *model.User) (*model.PurchaseOrder, []*model.PurchaseOrderItem) {
	var items []*model.PurchaseOrderItem
	var total float64
	var currentTotal, currentDiscount float64

	for _, i := range r.PurchaseOrderItems {
		ivID, _ := common.Decrypt(i.ItemVariantID)
		tempDiscount := ((float64(i.Quantity) * i.UnitPrice) * float64(i.Discount)) / float64(100)
		tempSubTotal := float64(i.Quantity) * i.UnitPrice
		DecryptPOI, _ := common.Decrypt(i.ID)
		poItem := &model.PurchaseOrderItem{
			ID:          DecryptPOI,
			ItemVariant: &model.ItemVariant{ID: ivID},
			Quantity:    float32(common.FloatPrecision(float64(i.Quantity), 2)),
			UnitPrice:   i.UnitPrice,
			Discount:    float32(common.FloatPrecision(float64(i.Discount), 2)),
			Subtotal:    common.FloatPrecision(tempSubTotal-tempDiscount, 0),
			Note:        i.Note,
		}
		items = append(items, poItem)
		total += common.FloatPrecision(tempSubTotal-tempDiscount, 0)
	}

	currentDiscount = (float64(r.Discount) * total) / float64(100)
	currentTotal = total - currentDiscount

	po := r.PurchaseOrder
	supplierID, _ := common.Decrypt(r.SupplierID)

	po.Supplier = &model.Partnership{ID: supplierID}
	po.RecognitionDate = r.RecognitionDate
	po.EtaDate = r.EtaDate
	po.ShipmentCost = r.ShipmentCost
	po.Note = r.Note
	po.UpdatedBy = user
	po.UpdatedAt = time.Now()
	po.DiscountAmount = common.FloatPrecision(r.DiscountAmount, 0)
	po.Discount = float32(common.FloatPrecision(float64(r.Discount), 2))
	po.Tax = float32(common.FloatPrecision(float64(r.Tax), 2))
	po.TaxAmount = common.FloatPrecision(float64(currentTotal*float64(r.Tax))/float64(100), 0)
	po.IsPercentage = r.IsPercentage
	po.PurchaseOrderItems = items

	if r.IsPercentage == int8(1) {
		po.DiscountAmount = common.FloatPrecision(currentDiscount, 0)
		po.TotalCharge = common.FloatPrecision(currentTotal+po.TaxAmount+r.ShipmentCost, 0)
	} else {
		if r.DiscountAmount > float64(0) {
			po.DiscountAmount = common.FloatPrecision(r.DiscountAmount, 0)
			po.Discount = float32(common.FloatPrecision(float64(r.DiscountAmount/total)*float64(100), 0))
			currentTotal = total - r.DiscountAmount
			po.TotalCharge = common.FloatPrecision(currentTotal+po.TaxAmount+r.ShipmentCost, 0)
		} else {
			po.TotalCharge = common.FloatPrecision(total+po.TaxAmount+r.ShipmentCost, 0)
		}
	}

	return po, items
}

// cancelRequest data struct that stored request data when requesting an cancel purchase order process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type cancelRequest struct {
	CancelledNote string               `json:"cancelled_note"`
	PurchaseOrder *model.PurchaseOrder `json:"-"`
}

// Validate implement validation.Requests interfaces.
func (r *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if r.PurchaseOrder.DocumentStatus == "cancelled" {
		o.Failure("document_status", "document status is cancelled")
	}

	if _, t, e := inventory.CheckCancelStock(uint64(r.PurchaseOrder.ID), "purchase_order"); e == nil && t == true {
		o.Failure("item_variant_stock", "item have been sold")
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *cancelRequest) Messages() map[string]string {
	return map[string]string{}
}

func (r *cancelRequest) Transform() *model.PurchaseOrder {
	r.PurchaseOrder.DocumentStatus = "cancelled"
	r.PurchaseOrder.CancelledNote = r.CancelledNote

	return r.PurchaseOrder
}
