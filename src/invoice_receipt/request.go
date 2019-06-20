// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package invoiceReceipt

import (
	"fmt"
	"strconv"
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create invoiceReceipt process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	SalesOrder      []SalesOrderRequest `json:"sales_order" valid:"required"`
	PartnershipID   string              `json:"partnership_id" valid:"required"`
	RecognitionDate time.Time           `json:"recognition_date" valid:"required"`
	Note            string              `json:"note"`
	SessionData     *auth.SessionData
}

// SalesOrderRequest data struct stored request data sales order id
type SalesOrderRequest struct {
	SalesOrderID string `json:"id"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if m, e := common.Decrypt(r.PartnershipID); e != nil {
		o.Failure("partnership_id", "partnership_id is not valid")
	} else {
		// cek partnersip
		partner := &model.Partnership{ID: m}
		if err := partner.Read("ID"); err != nil {
			o.Failure("partnership_id", "partnership_id not found")
		} else {
			// cek partnership is_deleted dan is_archived
			if partner.IsDeleted == int8(1) || partner.IsArchived == int8(1) {
				o.Failure("partnership_id", "partnership_id not found")
			} else {
				for k, i := range r.SalesOrder {
					if soID, e := common.Decrypt(i.SalesOrderID); e != nil {
						o.Failure(fmt.Sprintf("sales_orders.%d.id.invalid", k), "sales_order_id is not valid")
					} else {
						// cek sales order
						so := &model.SalesOrder{ID: soID}
						if err := so.Read("ID"); err != nil {
							o.Failure(fmt.Sprintf("sales_orders.%d.id.invalid", k), "sales_order_id is not found")
						} else {
							// cek customer sales order harus sama
							if partner.ID != so.Customer.ID {
								o.Failure(fmt.Sprintf("sales_orders.%d.id.invalid", k), "customer_id must be same")
							}
							// cek document status sales order
							if so.DocumentStatus != "active" && so.DocumentStatus != "requested_cancel" {
								o.Failure(fmt.Sprintf("sales_orders.%d.document_status.invalid", k), "document_status must be active or requested_cancel")
							}
							// cek invoice status sales order
							if so.InvoiceStatus == "new" {
								o.Failure(fmt.Sprintf("sales_orders.%d.invoice_status.invalid", k), "sales order has no invoice yet")
							}

							sInvoices, _ := GetNotBundledSalesInvoiceBySOID(so.ID)

							for _, si := range sInvoices {
								if err := si.Read("ID"); err == nil {
									// cek is_bundled sales invoice
									if si.IsBundled == int8(1) {
										o.Failure(fmt.Sprintf("sales_invoice.%d.is_bundled.invalid", k), "is_bundled in sales invoice must be 0")
									}

									// cek is_deleted sales invoice
									if si.IsDeleted == int8(1) {
										o.Failure(fmt.Sprintf("sales_invoice.%d.is_deleted.invalid", k), "sales invoice was deleted")
									}
								}
							}

							sReturns, _ := GetNotBundledSalesReturnBySOID(so.ID)
							for _, sr := range sReturns {
								if err := sr.Read("ID"); err == nil {
									// cek is_bundled sales return
									if sr.IsBundled == int8(1) {
										o.Failure(fmt.Sprintf("sales_return.%d.is_bundled.invalid", k), "is_bundled in sales return must be 0")
									}

									// cek is_deleted sales return
									if sr.IsDeleted == int8(1) {
										o.Failure(fmt.Sprintf("sales_return.%d.is_deleted.invalid", k), "sales return was deleted")
									}
								}
							}
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
func (r *createRequest) Transform() (*model.InvoiceReceipt, []*model.SalesOrder) {
	partnershipID, _ := common.Decrypt(r.PartnershipID)
	code, _ := util.CodeGen("code_invoice_receipt", "invoice_receipt")
	m := &model.InvoiceReceipt{
		Partnership:     &model.Partnership{ID: partnershipID},
		Code:            code,
		RecognitionDate: r.RecognitionDate,
		Note:            r.Note,
		DocumentStatus:  "new",
		CreatedBy:       &model.User{ID: r.SessionData.User.ID},
		CreatedAt:       time.Now(),
	}

	var mso []*model.SalesOrder
	for _, i := range r.SalesOrder {
		soid, _ := common.Decrypt(i.SalesOrderID)
		nso := &model.SalesOrder{ID: soid}

		mso = append(mso, nso)
	}

	return m, mso
}

// PaymentInvoiceReceipt data struct that stored request data when requesting an update payment invoiceReceipt process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type PaymentInvoiceReceipt struct {
	InvoiceReceiptItem   []PaymentInvoiceReceiptItem   `json:"invoice_receipt_item" valid:"required"`
	InvoiceReceiptReturn []PaymentInvoiceReceiptReturn `json:"invoice_receipt_return"`
	FinanceRevenue       []FinanceRevenueRequest       `json:"finance_revenue" valid:"required"`
	InvoiceReceiptID     int64
	SessionData          *auth.SessionData
}

// PaymentInvoiceReceiptItem data struct for create data payment invoiceReceiptItem
type PaymentInvoiceReceiptItem struct {
	InvoiceReceiptItemID string `json:"id"`
}

// PaymentInvoiceReceiptReturn data struct for create data payment invoiceReceiptReturn
type PaymentInvoiceReceiptReturn struct {
	InvoiceReceiptReturnID string `json:"id"`
}

// FinanceRevenueRequest data struct for create data payment invoiceReceipt
type FinanceRevenueRequest struct {
	RecognitionDate time.Time `json:"recognition_date"`
	PaymentMethod   string    `json:"payment_method" valid:"in:cash,debit_card,credit_card,giro"`
	BankNumber      string    `json:"bank_number"`
	BankAccountId   string    `json:"bank_account_id"`
	BankName        string    `json:"bank_name"`
	BankHolder      string    `json:"bank_holder"`
	Amount          float64   `json:"amount"`
	Note            string    `json:"note"`
}

// Validate implement validation.Requests interfaces.
func (r *PaymentInvoiceReceipt) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	var subtotalInvoice float64
	for x, k := range r.InvoiceReceiptItem {
		ln := strconv.Itoa(x)

		if paymentInvoiceReceiptItemID, e := common.Decrypt(k.InvoiceReceiptItemID); e != nil {
			o.Failure("invoice_receipt_item."+ln+".invoice_receipt_item_id.invalid", "invoice_receipt_item_id can't be decrypt")
		} else {
			// cek invoice reciept item
			invoiceReceiptItem := new(model.InvoiceReceiptItem)
			invoiceReceiptItem.ID = paymentInvoiceReceiptItemID
			if e := invoiceReceiptItem.Read("ID"); e != nil {
				o.Failure("invoice_receipt_item."+ln+".invoice_receipt_item_id.invalid", "invoice_receipt_item_id doesn't exist")
			} else {
				salesInvoice := new(model.SalesInvoice)
				salesInvoice.ID = invoiceReceiptItem.SalesInvoice.ID
				if err := salesInvoice.Read("ID"); err != nil {
					o.Failure("invoice_receipt_item."+ln+".sales_invoice.id.invalid", "sales_invoice not found")
				} else {
					// cek document status sales invoice
					if salesInvoice.DocumentStatus == "finished" {
						o.Failure("invoice_receipt_item."+ln+".sales_invoice.document_status.invalid", "document status has been finished")
					} else {
						subtotalInvoice += invoiceReceiptItem.Subtotal
					}

					// cek is_deleted sales invoice
					if salesInvoice.IsDeleted == int8(1) {
						o.Failure("sales_invoice."+ln+".is_deleted.document_status.invalid", "sales invoice was deleted")
					}
				}
			}
		}
	}

	// cek payment method
	for i, m := range r.FinanceRevenue {
		kn := strconv.Itoa(i)

		if m.BankNumber != "" {
			if m.PaymentMethod != "debit_card" && m.PaymentMethod != "credit_card" {
				o.Failure("finance_revenue."+kn+".bank_number.invalid", "payment_method must be debit_card or credit_card")
			}
		}

		if m.BankName != "" {
			if m.PaymentMethod != "debit_card" && m.PaymentMethod != "credit_card" {
				o.Failure("finance_revenue."+kn+".bank_name.invalid", "payment_method must be debit_card or credit_card")
			}
		}

		if m.BankHolder != "" {
			if m.PaymentMethod != "debit_card" && m.PaymentMethod != "credit_card" {
				o.Failure("finance_revenue."+kn+".bank_holder.invalid", "payment_method must be debit_card or credit_card")
			}
		}
	}

	if m, e := getDetailInvoiceReceipt("id", r.InvoiceReceiptID); e == nil {
		// cek document status invoice receipt
		if m.DocumentStatus != "new" {
			o.Failure("document_status", "document_status in invoice_receipt must be new")
		} else {
			var amountFinanceRevenue float64
			for _, mx := range r.FinanceRevenue {
				amountFinanceRevenue += mx.Amount
			}

			// update total amount invoice receipt menjadi subtotalInvoice
			m.TotalAmount = subtotalInvoice
			if m.TotalAmount != amountFinanceRevenue {
				o.Failure("amount", "amount and total_amount in invoice_receipt must be equal")
			}
			//
			//if o.Valid {
			//	m.Save("TotalAmount")
			//}
		}
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *PaymentInvoiceReceipt) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *PaymentInvoiceReceipt) Transform(InvoiceReceiptID int64) ([]*model.InvoiceReceiptItem, []*model.InvoiceReceiptReturn, []*model.FinanceRevenue) {
	var financeRevenues []*model.FinanceRevenue
	var paymentInvoiceReceiptItems []*model.InvoiceReceiptItem
	var paymentInvoiceReceiptReturns []*model.InvoiceReceiptReturn

	for _, receiptItem := range r.InvoiceReceiptItem {
		miri, _ := common.Decrypt(receiptItem.InvoiceReceiptItemID)
		ri := &model.InvoiceReceiptItem{
			ID: miri,
		}
		ri.Read("ID")
		paymentInvoiceReceiptItems = append(paymentInvoiceReceiptItems, ri)
	}

	for _, receiptReturn := range r.InvoiceReceiptReturn {
		mirr, _ := common.Decrypt(receiptReturn.InvoiceReceiptReturnID)
		rr := &model.InvoiceReceiptReturn{
			ID: mirr,
		}
		rr.Read("ID")
		paymentInvoiceReceiptReturns = append(paymentInvoiceReceiptReturns, rr)
	}

	for _, m := range r.FinanceRevenue {
		bai, _ := common.Decrypt(m.BankAccountId)
		var ba = new(model.BankAccount)

		if bai != 0 {
			ba = &model.BankAccount{ID: bai}
			ba.Read("ID")
		}else{
			ba = nil
		}

		fr := &model.FinanceRevenue{
			RefID:           uint64(InvoiceReceiptID),
			RefType:         "invoice_receipt",
			RecognitionDate: m.RecognitionDate,
			DocumentStatus:  "cleared",
			PaymentMethod:   m.PaymentMethod,
			BankAccount:     ba,
			BankNumber:      m.BankNumber,
			BankName:        m.BankName,
			BankHolder:      m.BankHolder,
			Amount:          m.Amount,
			Note:            m.Note,
			CreatedBy:       &model.User{ID: r.SessionData.User.ID},
			CreatedAt:       time.Now(),
		}
		financeRevenues = append(financeRevenues, fr)
	}

	return paymentInvoiceReceiptItems, paymentInvoiceReceiptReturns, financeRevenues
}
