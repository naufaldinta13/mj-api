// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package salesInvoice

import (
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create salesInvoice process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	Session         *auth.SessionData
	RecognitionDate time.Time `json:"recognition_date" valid:"required"`
	SalesOrderID    string    `json:"sales_order_id" valid:"required"`
	DueDate         time.Time `json:"due_date" valid:"required"`
	BillingAddress  string    `json:"billing_address" valid:"required"`
	TotalAmount     float64   `json:"total_amount" valid:"required|gte:0"`
	Note            string    `json:"note"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	// check sales order
	if id, e := common.Decrypt(r.SalesOrderID); e != nil {
		o.Failure("sales_order_id", "not valid")
	} else {
		sales := &model.SalesOrder{ID: id, IsDeleted: int8(0)}
		if e = sales.Read("ID", "IsDeleted"); e != nil {
			o.Failure("sales_order_id", "not valid")
		} else {
			// check status
			if sales.InvoiceStatus == "finished" || sales.DocumentStatus == "approved_cancel" || sales.DocumentStatus == "requested_cancel" {
				o.Failure("sales_order_id", "cannot create invoice")
			} else {
				// check total amount from all SI in sales order
				totalAmountInvoice, _ := GetSumTotalAmountSalesInvoiceBySalesOrder(id)
				if r.TotalAmount+totalAmountInvoice > (sales.TotalCharge) {
					o.Failure("total_amount", "exceed the limit")
				}

				if r.TotalAmount == 0 && sales.TotalCharge != 0 {
					o.Failure("total_amount", "total amount cannot be 0")
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
func (r *createRequest) Transform() (*model.SalesInvoice, error) {
	// prepare the data
	var code string
	var e error
	si := &model.SalesInvoice{}
	if code, e = util.CodeGen("code_sales_invoice", "sales_invoice"); e == nil {
		soID, _ := common.Decrypt(r.SalesOrderID)
		SO := &model.SalesOrder{ID: soID}
		SO.Read("ID")

		// insert
		si.SalesOrder = SO
		si.Code = code
		si.RecognitionDate = r.RecognitionDate
		si.DueDate = r.DueDate
		si.BillingAddress = r.BillingAddress
		si.TotalAmount = r.TotalAmount
		si.Note = r.Note
		si.DocumentStatus = "new"
		si.IsBundled = int8(0)
		si.IsDeleted = int8(0)
		si.CreatedBy = r.Session.User
		si.CreatedAt = time.Now()
	}
	return si, e
}

// updateRequest data struct that stored request data when requesting an update salesInvoice process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type updateRequest struct {
	Session         *auth.SessionData
	SalesInvoiceOld *model.SalesInvoice
	RecognitionDate time.Time `json:"recognition_date" valid:"required"`
	TotalAmount     float64   `json:"total_amount" valid:"required|gte:0"`
	Note            string    `json:"note"`
}

// Validate implement validation.Requests interfaces.
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	// check status sales invoice
	if r.SalesInvoiceOld.DocumentStatus != "new" || r.SalesInvoiceOld.SalesOrder.DocumentStatus == "approved_cancel" || r.SalesInvoiceOld.SalesOrder.DocumentStatus == "requested_cancel" {
		o.Failure("sales_invoice", "cannot be update")
	} else {
		// check total amount from all SI in sales order
		totalAmount, _ := GetSumTotalAmountSalesInvoiceBySalesOrder(r.SalesInvoiceOld.SalesOrder.ID)
		if (r.TotalAmount + (totalAmount - r.SalesInvoiceOld.TotalAmount)) > (r.SalesInvoiceOld.SalesOrder.TotalCharge) {
			o.Failure("total_amount", "exceed the limit")
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
func (r *updateRequest) Transform() {
	r.SalesInvoiceOld.RecognitionDate = r.RecognitionDate
	r.SalesInvoiceOld.Note = r.Note
	r.SalesInvoiceOld.TotalAmount = r.TotalAmount
	r.SalesInvoiceOld.UpdatedBy = r.Session.User
	r.SalesInvoiceOld.UpdatedAt = time.Now()
}
