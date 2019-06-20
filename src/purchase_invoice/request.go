// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchaseInvoice

import (
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create purchaseInvoice process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	SessionData     *auth.SessionData `json:"-"`
	PurchaseOrder   string            `json:"purchase_order" valid:"required"`
	RecognitionDate time.Time         `json:"recognition_date" valid:"required"`
	DueDate         time.Time         `json:"due_date" valid:"required"`
	TotalAmount     float64           `json:"total_amount" valid:"required|gte:0"`
	Note            string            `json:"note"`
	BillingAddress  string            `json:"billing_address" valid:"required"`
	TotalDiff       float64
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	if id, e := common.Decrypt(r.PurchaseOrder); e == nil {
		po := &model.PurchaseOrder{ID: id, IsDeleted: int8(0)}
		// cek purchase order
		if e = po.Read("ID", "IsDeleted"); e != nil {
			o.Failure("purchase_order", "purchase_order id doesn't exist")
		} else {
			// cek invoice status
			if po.InvoiceStatus == "finished" {
				o.Failure("purchase_order", "invoice for purchase order has been finished")

			}
			// cek total amount dan total charge
			if taPI, e := CalculateTotalAmountPI(id, 0); e == nil {
				if taPI > po.TotalCharge {
					o.Failure("total_amount", "total_amount purchase invoice greater than purchase order")
				}

				r.TotalDiff = po.TotalCharge - taPI

				if r.TotalAmount > r.TotalDiff {
					o.Failure("total_amount", "total_amount purchase invoice cant be greater than purchase order")
				}
			}

		}
	} else {
		o.Failure("purchase_order", "purchase_order id cannot be decrypt")
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *createRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *createRequest) Transform() *model.PurchaseInvoice {

	poID, _ := common.Decrypt(r.PurchaseOrder)
	code, _ := util.CodeGen("code_purchase_invoice", "purchase_invoice")

	pi := &model.PurchaseInvoice{
		Code:            code,
		RecognitionDate: r.RecognitionDate,
		PurchaseOrder:   &model.PurchaseOrder{ID: poID},
		DueDate:         r.DueDate,
		TotalAmount:     r.TotalAmount,
		Note:            r.Note,
		CreatedBy:       r.SessionData.User,
		CreatedAt:       time.Now(),
		DocumentStatus:  "new",
		BillingAddress:  r.BillingAddress,
	}

	return pi
}

// updateRequest data struct that stored request data when requesting an create purchaseInvoice process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type updateRequest struct {
	PI              *model.PurchaseInvoice
	SessionData     *auth.SessionData `json:"-"`
	RecognitionDate time.Time         `json:"recognition_date" valid:"required"`
	TotalAmount     float64           `json:"total_amount" valid:"required|gte:0"`
	Note            string            `json:"note"`
	TotalDiff       float64
}

// Validate implement validation.Requests interfaces.
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	po := &model.PurchaseOrder{ID: r.PI.PurchaseOrder.ID, IsDeleted: int8(0)}
	if e := po.Read("ID", "IsDeleted"); e == nil {
		if po.InvoiceStatus == "finished" {
			o.Failure("purchase_order", "invoice for purchase order has been finished")

		}

		// cek total amount dan total charge
		if taPI, e := CalculateTotalAmountPI(po.ID, r.PI.ID); e == nil {
			if taPI > po.TotalCharge {
				o.Failure("total_amount", "total_amount purchase invoice greater than purchase order")
			}

			r.TotalDiff = po.TotalCharge - taPI

			if r.TotalAmount > r.TotalDiff {
				o.Failure("total_amount", "total_amount purchase invoice cant be greater than purchase order")
			}
		}

	}
	if r.PI.DocumentStatus != "new" {
		o.Failure("purchase_invoice", "can't update document")
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *updateRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *updateRequest) Transform(piID *model.PurchaseInvoice) *model.PurchaseInvoice {

	pi := &model.PurchaseInvoice{
		ID:              piID.ID,
		RecognitionDate: r.RecognitionDate,
		Note:            r.Note,
		TotalAmount:     r.TotalAmount,
		UpdatedAt:       time.Now(),
		UpdatedBy:       r.SessionData.User,
	}

	return pi
}
