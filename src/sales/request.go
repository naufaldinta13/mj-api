// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sales

import (
	"fmt"
	"time"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/inventory"
	"git.qasico.com/mj/api/src/partnership"
	"git.qasico.com/mj/api/src/pricing_type"
)

type createRequest struct {
	Code                 string            `json:"code"`
	RecognitionDate      time.Time         `json:"recognition_date" valid:"required"`
	ReferencesID         string            `json:"references_id"`
	EtaDate              time.Time         `json:"eta_date"`
	ShipmentAddress      string            `json:"shipment_address"`
	CustomerID           string            `json:"customer_id" valid:"required"`
	AutoInvoice          int8              `json:"auto_invoice"`
	AutoFullfilment      int8              `json:"auto_fullfilment"`
	SalesOrderItem       []salesOrderItem  `json:"sales_order_item" valid:"required"`
	Discount             float32           `json:"discount"`
	DiscountAmount       float64           `json:"discount_amount" valid:"gte:0"`
	IsPaid               int8              `json:"is_paid"`
	Tax                  float32           `json:"tax" valid:"gte:0|lte:100"`
	TaxAmount            float64           `json:"-"`
	ShipmentCost         float64           `json:"shipment_cost" valid:"gte:0"`
	TotalPrice           float64           `json:"-"`
	TotalCharge          float64           `json:"-"`
	TotalCost            float64           `json:"-"`
	Note                 string            `json:"note"`
	IsPercentageDiscount int8              `json:"is_percentage_discount" valid:"in:0,1"`
	Session              *auth.SessionData `json:"-"`
}

type salesOrderItem struct {
	ID            string  `json:"id"`
	ItemVariantID string  `json:"item_variant_id" valid:"required"`
	Quantity      float32 `json:"quantity" valid:"required|gt:0"`
	Discount      float32 `json:"discount" valid:"gte:0|lte:100"`
	UnitPrice     float64 `json:"unit_price" valid:"required|gt:0"`
	PricingType   string  `json:"pricing_type" valid:"required"`
	Subtotal      float64 `json:"-"`
	Note          string  `json:"note"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var PartnerID int64
	var partner *model.Partnership
	var sorder *model.SalesOrder
	var code string

	// cek partner id ada atau tidak
	PartnerID, err = common.Decrypt(r.CustomerID)
	if err != nil {
		o.Failure("customer_id.invalid", "Partnership id is invalid")
	}

	// cek partner ecist or not
	partner, err = partnership.GetPartnershipByField("id", PartnerID)

	if err != nil || partner == nil || partner.IsDeleted == 1 || partner.IsArchived == 1 {
		o.Failure("customer_id.invalid", "Partnership id is not found")
	} else {
		code, _ = CodeGen(partner.IsDefault == 1)

		if partner.PartnershipType != "customer" {
			o.Failure("customer_id", "Customer needed to have partner type customer not supplier")
		}

		if partner.IsDefault == int8(1) {
			if r.AutoFullfilment == int8(0) {
				o.Failure("auto_fullfilment", "auto fullfilment need to be filled if Customer is walk in")
			}
			if r.AutoInvoice == int8(0) {
				o.Failure("auto_invoice", "auto invoice need to be filled if Customer is walk in")
			}
			if r.EtaDate.IsZero() {
				o.Failure("eta_date", "ETA Date need to be filled if Customer is walk in")
			}
		} else {
			if r.ShipmentAddress == "" {
				o.Failure("shipment_address", "Shipment address is required")
			}
		}

		if partner.OrderRule == "one_bill" {
			var soContainer *model.SalesOrder
			orm.NewOrm().Raw("SELECT * FROM sales_order WHERE customer_id = ? AND document_status = 'new' OR document_status = 'active' AND invoice_status = 'active';", PartnerID).QueryRow(&soContainer)
			if soContainer != nil {
				o.Failure("customer_id", "Partner still have unfinished invoice")
			}
		} else if partner.OrderRule == "plafon" {
			current := partner.TotalDebt + r.TotalCharge
			if current >= partner.MaxPlafon {
				o.Failure("customer_id", "Partnership has already reached given max plafon")
			}
		}
	}

	if r.IsPaid == 1 {
		if r.AutoInvoice != 1 {
			o.Failure("auto_invoice", "Auto invoice must be checked if Auto Paid is checked")
		}
	}

	so := &model.SalesOrder{Code: code}

	if err := so.Read("Code"); err == nil {
		o.Failure("code", "Code sales order is already being used")
	} else {
		r.Code = code
	}

	tz := time.Time{}
	if r.EtaDate == tz {
		o.Failure("eta_date", "Field is required.")
	}

	// cek status reference id
	if r.ReferencesID != "" {
		refID, err := common.Decrypt(r.ReferencesID)
		if err != nil {
			o.Failure("references_id", "References id is not valid")
		}
		var emptyLoad []string
		sorder, err = GetDetailSalesOrder(refID, emptyLoad)
		if err != nil {
			o.Failure("references_id", "References is not found")
		} else {
			if sorder.DocumentStatus != "approved_cancel" {
				o.Failure("references_id", "References document status is not cancel")
			}
		}
	}

	checkDuplicate := make(map[string]bool)
	var checkVariant = make(map[int64]*model.ItemVariant)

	for _, row := range r.SalesOrderItem {
		IVariantID, _ := common.Decrypt(row.ItemVariantID)
		ivar := &model.ItemVariant{ID: IVariantID}
		ivar.Read("ID")
		////////////////////////////////
		if checkVariant[ivar.ID] == nil {
			ivar.CommitedStock -= row.Quantity
			checkVariant[ivar.ID] = ivar
		} else {
			variant := checkVariant[ivar.ID]
			variant.CommitedStock -= row.Quantity
			checkVariant[ivar.ID] = variant
		}
		////////////////////////////////
	}

	// cek setiap sales order item
	for i, row := range r.SalesOrderItem {
		var UnitA float64
		var PricingID, IVariantID int64
		var ItemVariant *model.ItemVariant
		var IVPrice *model.ItemVariantPrice
		// cek item variant,pricing type dan item variant price
		PricingID, err = common.Decrypt(row.PricingType)

		if err != nil {
			o.Failure(fmt.Sprintf("sales_order_item.%d.pricing_type.invalid", i), "Pricing type id is invalid")
		}

		IVariantID, err = common.Decrypt(row.ItemVariantID)
		if err != nil {
			o.Failure(fmt.Sprintf("sales_order_item.%d.item_variant.invalid", i), "Item Variant id is invalid")
		}

		var pt = &model.PricingType{ID: PricingID}
		pt.Read("ID")
		var iv = &model.ItemVariant{ID: IVariantID}
		iv.Read("ID")

		IVPrice, err = getItemVariantPricing(pt, iv)
		if err == nil {
			if pt.ParentType != nil {
				if pt.RuleType == "increment" {
					if pt.IsPercentage == int8(1) {
						temp := (pt.Nominal * IVPrice.UnitPrice) / float64(100)
						UnitA = IVPrice.UnitPrice + temp
					} else {
						UnitA = IVPrice.UnitPrice + pt.Nominal
					}

					if row.UnitPrice < UnitA {
						o.Failure(fmt.Sprintf("sales_order_item.%d.unit_price.invalid", i), "Unit price is too small")
					}
				} else {
					if pt.IsPercentage == int8(1) {

						temp := (pt.Nominal * IVPrice.UnitPrice) / float64(100)
						UnitA = IVPrice.UnitPrice - temp
					} else {
						UnitA = IVPrice.UnitPrice - pt.Nominal
					}

					if UnitA < 0 {
						o.Failure(fmt.Sprintf("sales_order_item.%d.pricing_type.invalid", i), "Pricing type can make price become zero")
					}
					if row.UnitPrice < UnitA {
						o.Failure(fmt.Sprintf("sales_order_item.%d.unit_price.invalid", i), "Unit price is too small")
					}
				}
			} else {
				if row.UnitPrice < IVPrice.UnitPrice {
					o.Failure(fmt.Sprintf("sales_order_item.%d.unit_price.invalid", i), "Unit price is too small")
				}
			}
		} else {
			o.Failure(fmt.Sprintf("sales_order_item.%d.unit_price.invalid", i), "item variant price doesn't exist")
		}

		ItemVariant, err = inventory.GetDetailItemVariant("id", IVariantID)
		if err != nil || ItemVariant == nil || ItemVariant.IsDeleted == int8(1) || ItemVariant.IsArchived == int8(1) {
			o.Failure(fmt.Sprintf("sales_order_item.%d.item_variant_id.invalid", i), "Item variant id not found")
		} else {

			// cek stock dari item variant sama quantity soi
			if (checkVariant[ItemVariant.ID].AvailableStock - ItemVariant.CommitedStock) < row.Quantity {
				o.Failure(fmt.Sprintf("sales_order_item.%d.quantity.invalid", i), "Stock item is not enough to be sold")
			}

			//check duplicate item variant id
			if checkDuplicate[row.ItemVariantID] == true {
				o.Failure(fmt.Sprintf("sales_order_item.%d.item_variant_id.invalid", i), " item variant id duplicate")
			} else {
				checkDuplicate[row.ItemVariantID] = true
			}
		}

		discamount := (row.UnitPrice * float64(row.Discount)) / float64(100)
		curamount := row.UnitPrice * float64(row.Quantity)
		subtotal := common.FloatPrecision(curamount-discamount, 0)

		r.TotalPrice += subtotal
	}

	if r.IsPercentageDiscount == int8(1) {
		if r.Discount < 0 || r.Discount > float32(100) {
			o.Failure("discount", "discount is less than and equal 0 or greater than 100")
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
func (r *createRequest) Transform() (sorder *model.SalesOrder) {
	idcustomer, _ := common.Decrypt(r.CustomerID)
	var partner *model.Partnership
	partner = &model.Partnership{ID: idcustomer, PartnershipType: "customer"}
	partner.Read("ID", "PartnershipType")

	var disc float32
	var discAmount float64
	// calculate
	if r.IsPercentageDiscount == int8(1) {
		disc = r.Discount
		discAmount = (r.TotalPrice * float64(disc)) / float64(100)
		r.DiscountAmount = discAmount
	} else {
		discAmount = r.DiscountAmount
		if discAmount != float64(0) {
			disc = float32(common.FloatPrecision((discAmount/r.TotalPrice)*float64(100), 2))
		}
	}
	curamount := r.TotalPrice - r.DiscountAmount
	r.TaxAmount = (curamount * float64(r.Tax)) / float64(100)
	r.TotalCharge = common.FloatPrecision(curamount+r.TaxAmount+r.ShipmentCost, 0)
	sorder = &model.SalesOrder{
		Customer:             partner,
		Code:                 r.Code,
		RecognitionDate:      r.RecognitionDate,
		EtaDate:              r.EtaDate,
		Discount:             disc,
		Tax:                  r.Tax,
		DiscountAmount:       discAmount,
		TaxAmount:            r.TaxAmount,
		ShipmentAddress:      r.ShipmentAddress,
		ShipmentCost:         r.ShipmentCost,
		TotalPrice:           r.TotalPrice,
		TotalCharge:          r.TotalCharge,
		TotalCost:            r.TotalCost,
		Note:                 r.Note,
		AutoFulfillment:      r.AutoFullfilment,
		AutoInvoice:          r.AutoInvoice,
		AutoPaid:             r.IsPaid,
		IsPercentageDiscount: r.IsPercentageDiscount,
		InvoiceStatus:        "new",
		FulfillmentStatus:    "new",
		ShipmentStatus:       "new",
		CreatedBy:            r.Session.User,
		CreatedAt:            time.Now(),
	}

	if partner.IsDefault == int8(1) {
		sorder.DocumentStatus = "active"
		sorder.InvoiceStatus = "active"
		sorder.FulfillmentStatus = "active"
		sorder.AutoInvoice = int8(1)
		sorder.AutoFulfillment = int8(1)
	} else {
		sorder.DocumentStatus = "new"
	}

	if sorder.AutoInvoice == int8(1) {
		sorder.DocumentStatus = "active"
		sorder.InvoiceStatus = "active"
	}

	if sorder.AutoFulfillment == int8(1) {
		sorder.DocumentStatus = "active"
		sorder.FulfillmentStatus = "active"
	}

	for _, row := range r.SalesOrderItem {
		idItemVar, _ := common.Decrypt(row.ItemVariantID)
		itemvar, _ := inventory.GetDetailItemVariant("id", idItemVar)

		discamount := ((row.UnitPrice * float64(row.Quantity)) * float64(row.Discount)) / float64(100)
		curamount := row.UnitPrice * float64(row.Quantity)
		subtotal := common.FloatPrecision(curamount-discamount, 0)
		row.Subtotal = subtotal

		soitem := model.SalesOrderItem{
			ItemVariant: itemvar,
			Quantity:    row.Quantity,
			UnitPrice:   row.UnitPrice,
			Discount:    row.Discount,
			Subtotal:    row.Subtotal,
		}

		sorder.SalesOrderItems = append(sorder.SalesOrderItems, &soitem)
	}

	return sorder
}

// updateRequest untuk update sales
type updateRequest struct {
	RecognitionDate      time.Time        `json:"recognition_date" valid:"required"`
	EtaDate              time.Time        `json:"eta_date" valid:"required"`
	SalesOrderItem       []salesOrderItem `json:"sales_order_item" valid:"required"`
	Discount             float32          `json:"discount"`
	DiscountAmount       float64          `json:"discount_amount" valid:"gte:0"`
	Tax                  float32          `json:"tax" valid:"gte:0|lte:100"`
	ShipmentCost         float64          `json:"shipment_cost" valid:"gte:0"`
	Note                 string           `json:"note"`
	IsPercentageDiscount int8             `json:"is_percentage_discount" valid:"in:0,1"`
	TaxAmount            float64          `json:"-"`
	TotalPrice           float64          `json:"-"`
	TotalCharge          float64          `json:"-"`
	TotalCost            float64          `json:"-"`

	SalesOrder *model.SalesOrder `json:"sales_order"`
}

// Validate implement validation.Requests interfaces.
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	// cek document status so
	if r.SalesOrder.DocumentStatus != "new" {
		o.Failure("document_status", "document_status should be new")
	}

	checkDuplicate := make(map[string]bool)
	checkVariant := make(map[int64]*model.ItemVariant)

	tz := time.Time{}
	if r.EtaDate == tz {
		o.Failure("eta_date", "Field is required.")
	}

	for _, row := range r.SalesOrderItem {
		IVariantID, _ := common.Decrypt(row.ItemVariantID)
		ivar := &model.ItemVariant{ID: IVariantID}
		ivar.Read("ID")
		////////////////////////////////
		if checkVariant[ivar.ID] == nil {
			ivar.CommitedStock -= row.Quantity
			checkVariant[ivar.ID] = ivar
		} else {
			variant := checkVariant[ivar.ID]
			variant.CommitedStock -= row.Quantity
			checkVariant[ivar.ID] = variant
		}
		////////////////////////////////
	}

	// cek setiap sales order item
	for i, row := range r.SalesOrderItem {
		// validasi id
		if row.ID != "" {
			if ID, err := common.Decrypt(row.ID); err != nil {
				o.Failure(fmt.Sprintf("sales_order_item.%d.id.invalid", i), "id is not valid")
			} else {
				soItem := &model.SalesOrderItem{ID: ID}
				if err := soItem.Read(); err != nil {
					o.Failure(fmt.Sprintf("sales_order_item.%d.id.invalid", i), "id is not found")
				}

				//check duplicate sales order item id
				if checkDuplicate[row.ID] == true {
					o.Failure(fmt.Sprintf("sales_order_item.%d.id.invalid", i), "id duplicate")
				} else {
					checkDuplicate[row.ID] = true
				}
			}
		}

		var UnitA float64
		var PricingID, IVariantID int64
		var ItemVariant *model.ItemVariant
		var IVPrice *model.ItemVariantPrice
		// cek pricing type
		PricingID, err = common.Decrypt(row.PricingType)
		if err != nil {
			o.Failure(fmt.Sprintf("sales_order_item.%d.pricing_type.invalid", i), "Pricing type id is invalid")
		} else {
			if _, err = pricingType.GetPricingTypeByID(PricingID); err != nil {
				o.Failure(fmt.Sprintf("sales_order_item.%d.pricing_type.invalid", i), "Pricing type id is not found")
			}
		}

		IVariantID, err = common.Decrypt(row.ItemVariantID)
		if err != nil {
			o.Failure(fmt.Sprintf("sales_order_item.%d.item_variant.invalid", i), "Item Variant id is invalid")
		}

		var pt = &model.PricingType{ID: PricingID}
		pt.Read("ID")
		var iv = &model.ItemVariant{ID: IVariantID}
		iv.Read("ID")

		IVPrice, err = getItemVariantPricing(pt, iv)
		if err == nil {
			if pt.ParentType != nil {
				if pt.RuleType == "increment" {
					if pt.IsPercentage == int8(1) {
						temp := (pt.Nominal * IVPrice.UnitPrice) / float64(100)
						UnitA = IVPrice.UnitPrice + temp
					} else {
						UnitA = IVPrice.UnitPrice + pt.Nominal
					}

					if row.UnitPrice < UnitA {
						o.Failure(fmt.Sprintf("sales_order_item.%d.unit_price.invalid", i), "Unit price is too small")
					}
				} else {
					if pt.IsPercentage == int8(1) {
						temp := (pt.Nominal * IVPrice.UnitPrice) / float64(100)
						UnitA = IVPrice.UnitPrice - temp
					} else {
						UnitA = IVPrice.UnitPrice - pt.Nominal
					}

					if UnitA < 0 {
						o.Failure(fmt.Sprintf("sales_order_item.%d.pricing_type.invalid", i), "Pricing type can make price become zero")
					}

					if row.UnitPrice < UnitA {
						o.Failure(fmt.Sprintf("sales_order_item.%d.unit_price.invalid", i), "Unit price is too small")
					}
				}
			} else {
				if row.UnitPrice < IVPrice.UnitPrice {
					o.Failure(fmt.Sprintf("sales_order_item.%d.unit_price.invalid", i), "Unit price is too small")
				}
			}
		} else {
			o.Failure(fmt.Sprintf("sales_order_item.%d.unit_price.invalid", i), "item variant price doesn't exist")
		}

		ItemVariant, err = inventory.GetDetailItemVariant("id", IVariantID)
		if err != nil || ItemVariant == nil || ItemVariant.IsDeleted == int8(1) || ItemVariant.IsArchived == int8(1) {
			o.Failure(fmt.Sprintf("sales_order_item.%d.item_variant_id.invalid", i), "Item variant id not found")
		} else {

			SoItemID, _ := common.Decrypt(row.ID)
			SoItem := &model.SalesOrderItem{ID: SoItemID}
			if e := SoItem.Read("ID"); e != nil {
				SoItem.Quantity = 0
			}

			// cek stock item variant
			if ((checkVariant[ItemVariant.ID].AvailableStock - checkVariant[ItemVariant.ID].CommitedStock) + SoItem.Quantity) < row.Quantity {
				o.Failure(fmt.Sprintf("sales_order_item.%d.quantity.invalid", i), "Stock item is not enough to be sold")
			}

			checkVariant[ItemVariant.ID].CommitedStock += row.Quantity

			//check duplicate item variant id
			if checkDuplicate[row.ItemVariantID] == true {
				o.Failure(fmt.Sprintf("sales_order_item.%d.item_variant_id.invalid", i), " item variant id duplicate")
			} else {
				checkDuplicate[row.ItemVariantID] = true
			}

		}

		// Calculate total price
		discamount := (row.UnitPrice * float64(row.Discount)) / float64(100)
		curamount := row.UnitPrice * float64(row.Quantity)
		subtotal := common.FloatPrecision(curamount-discamount, 0)

		r.TotalPrice += subtotal
	}

	if r.IsPercentageDiscount == int8(1) {
		if r.Discount < 0 || r.Discount > float32(100) {
			o.Failure("discount", "discount is less than and equal 0 or greater than 100")
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
func (r *updateRequest) Transform(user *model.User) (*model.SalesOrder, []*model.SalesOrderItem) {
	var disc float32
	var discAmount float64

	// calculate
	if r.IsPercentageDiscount == int8(1) {
		disc = r.Discount
		discAmount = (r.TotalPrice * float64(disc)) / float64(100)
	} else {
		discAmount = r.DiscountAmount
		if discAmount != float64(0) {
			disc = float32(common.FloatPrecision((discAmount/r.TotalPrice)*float64(100), 2))
		}
	}
	curamount := r.TotalPrice - r.DiscountAmount
	r.TaxAmount = (curamount * float64(r.Tax)) / float64(100)
	r.TotalCharge = common.FloatPrecision(curamount+r.TaxAmount+r.ShipmentCost, 0)

	so := r.SalesOrder
	so.RecognitionDate = r.RecognitionDate
	so.EtaDate = r.EtaDate
	so.Discount = disc
	so.DiscountAmount = discAmount
	so.Tax = r.Tax
	so.TaxAmount = r.TaxAmount
	so.ShipmentCost = r.ShipmentCost
	so.TotalPrice = r.TotalPrice
	so.TotalCharge = r.TotalCharge
	so.TotalCost = r.TotalCost
	so.IsPercentageDiscount = r.IsPercentageDiscount
	so.Note = r.Note
	so.UpdatedAt = time.Now()
	so.UpdatedBy = user

	var items []*model.SalesOrderItem
	for _, row := range r.SalesOrderItem {
		var ID int64
		if row.ID != "" {
			ID, _ = common.Decrypt(row.ID)
		}
		idItemVar, _ := common.Decrypt(row.ItemVariantID)
		itemvar, _ := inventory.GetDetailItemVariant("id", idItemVar)

		discamount := ((row.UnitPrice * float64(row.Quantity)) * float64(row.Discount)) / float64(100)
		curamount := row.UnitPrice * float64(row.Quantity)
		subtotal := common.FloatPrecision(curamount-discamount, 0)
		row.Subtotal = subtotal

		soitem := model.SalesOrderItem{
			ID:          ID,
			ItemVariant: itemvar,
			Quantity:    row.Quantity,
			UnitPrice:   row.UnitPrice,
			Discount:    row.Discount,
			Subtotal:    row.Subtotal,
			Note:        row.Note,
		}
		items = append(items, &soitem)
	}

	return so, items
}

// cancelRequest for note cancel
type cancelRequest struct {
	Sales         *model.SalesOrder
	CancelledNote string `json:"cancelled_note"`
}

// Validate implement validation.Requests interfaces.
func (r *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if r.Sales.DocumentStatus == "approved_cancel" {
		o.Failure("document_status", "sales order has approved cancel")
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *cancelRequest) Messages() map[string]string {
	return map[string]string{}
}

// cancelReq data struct that stored request data when requesting an request cancel sales process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type cancelReq struct {
	Session       *auth.SessionData
	Sales         *model.SalesOrder
	CancelledNote string `json:"cancelled_note"`
}

// Validate implement validation.Requests interfaces.
func (r *cancelReq) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	if r.Sales.DocumentStatus == "requested_cancel" || r.Sales.DocumentStatus == "approved_cancel" {
		o.Failure("document_status", "not valid")
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *cancelReq) Messages() map[string]string {
	return map[string]string{}
}

// Transform untuk mengubah isi model.
func (r *cancelReq) Transform(m *model.SalesOrder) {
	// ubah isi model sales order
	m.CancelledNote = r.CancelledNote
	m.DocumentStatus = "requested_cancel"
	m.RequestCancelAt = time.Now()
	m.RequestCancelBy = r.Session.User
}

// rejectCancelRequest data struct that stored request data when requesting an reject request cancel sales process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type rejectCancelRequest struct {
	SalesOrder *model.SalesOrder
}

// Validate implement validation.Requests interfaces.
func (r *rejectCancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	if r.SalesOrder.DocumentStatus != "requested_cancel" {
		o.Failure("sales_order", "not request for cancel")
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *rejectCancelRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform untuk mengubah isi model.
func (r *rejectCancelRequest) Transform() {
	// ubah isi model sales order
	r.SalesOrder.DocumentStatus = "active"
}
