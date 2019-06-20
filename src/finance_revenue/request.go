// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package financeRevenue

import (
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create financeRevenue process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	RefID           string    `json:"ref_id" valid:"required"`
	RefType         string    `json:"ref_type" valid:"required|in:sales_invoice,purchase_return"`
	RecognitionDate time.Time `json:"recognition_date" valid:"required"`
	PaymentMethod   string    `json:"payment_method" valid:"required|in:cash,debit_card,credit_card,giro"`
	Amount          float64   `json:"amount" valid:"required"`
	BankName        string    `json:"bank_name"`
	BankNumber      string    `json:"bank_number"`
	BankHolder      string    `json:"bank_holder"`
	BankAccountID   string    `json:"bank_account_id"`
	Note            string    `json:"note"`
	Session         *auth.SessionData
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	refID, err := common.Decrypt(r.RefID)
	if err != nil {
		o.Failure("ref_id", "can't decrypt ref_id")
	}

	if r.RefType == "sales_invoice" {
		// cek sales invoice
		si := &model.SalesInvoice{ID: refID}
		if e := si.Read("ID"); e != nil {
			o.Failure("ref_id", "can't found sales_invoice")
		} else {
			// cek document_status sales invoice
			if si.DocumentStatus == "finished" {
				o.Failure("sales_invoice", "sales_invoice already paid")
			} else {
				// cek total revenued pada sales invoice
				// kasir tidak bisa membuat invoice jika total revenued lebih dari 0
				if r.Session.User.Usergroup.ID == int64(4) && si.TotalRevenued > float64(0) {
					o.Failure("sales_invoice", "sales invoice already paid")
				} else {
					// jumlahkan semua amount dari finance_revenue yang memiliki referensi dari sales invoice yang diinput
					amountFr := SumAmountFinanceRevenue(uint64(si.ID), r.RefType, 0)
					amountFr += r.Amount
					// cek jumlah (sum amount finance revenue + amount inputan) > total amount di sales invoice
					if amountFr > si.TotalAmount {
						o.Failure("amount", "the amount of the pay is greater than the total charge of the sales_invoice")
					} else {
						si.DocumentStatus = "active"
						si.Save("DocumentStatus")
					}
				}
			}
		}
	}

	if r.RefType == "purchase_return" {
		// cek purchase return
		pr := &model.PurchaseReturn{ID: refID}
		if e := pr.Read("ID"); e != nil {
			o.Failure("ref_id", "can't found purchase_return")
		} else {
			// cek document_status purchase return
			if pr.DocumentStatus == "finished" {
				o.Failure("purchase_return", "purchase_return already paid")
			} else {
				// jumlahkan semua amount dari finance_revenue yang memiliki referensi dari purchase return yang diinput
				amountFr := SumAmountFinanceRevenue(uint64(pr.ID), r.RefType, 0)
				amountFr += r.Amount
				// cek jumlah (sum amount finance revenue + amount inputan) > total amount di purchase return
				if amountFr > pr.TotalAmount {
					o.Failure("amount", "the amount of the pay is greater than the total charge of the purchase_return")
				} else {
					pr.DocumentStatus = "active"
					pr.Save("DocumentStatus")
				}
			}
		}
	}

	if r.PaymentMethod == "debit_card" || r.PaymentMethod == "credit_card" {
		if r.BankNumber == "" && r.PaymentMethod == "debit_card" {
			o.Failure("bank_number", "Bank number is required")
		}
		if r.BankAccountID == "" {
			o.Failure("bank_account_id", "Please select your Bank")
		}
		if r.BankName == "" {
			o.Failure("bank_name", "Bank name is required")
		}
		if r.BankHolder == "" {
			o.Failure("bank_holder", "Bank holder is required")
		}
	} else if r.PaymentMethod == "giro" {
		if r.BankNumber == "" {
			o.Failure("bank_number", "Bank number is required")
		}
		if r.BankName != "" {
			o.Failure("bank_name", "Bank name can't be inputted if payment method is giro")
		}
	} else {
		if r.BankNumber != "" {
			o.Failure("bank_number", "Bank number can't be inputted if payment method is cash")
		}
		if r.BankName != "" {
			o.Failure("bank_name", "Bank name can't be inputted if payment method is cash")
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
func (r *createRequest) Transform() *model.FinanceRevenue {
	ba, _ := common.Decrypt(r.BankAccountID)
	bankAccount := &model.BankAccount{ID: ba}

	if ba == 0 {
		bankAccount = nil
	}

	refID, _ := common.Decrypt(r.RefID)
	fr := &model.FinanceRevenue{
		RefID:           uint64(refID),
		RefType:         r.RefType,
		BankAccount:     bankAccount,
		RecognitionDate: r.RecognitionDate,
		PaymentMethod:   r.PaymentMethod,
		Amount:          r.Amount,
		BankName:        r.BankName,
		BankNumber:      r.BankNumber,
		BankHolder:      r.BankHolder,
		Note:            r.Note,
		DocumentStatus:  "uncleared",
		IsDeleted:       int8(0),
		CreatedAt:       time.Now(),
		CreatedBy:       r.Session.User,
	}

	return fr
}

// updateRequest data struct that stored request data when requesting an update financeRevenue process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type updateRequest struct {
	RecognitionDate time.Time             `json:"recognition_date" valid:"required"`
	PaymentMethod   string                `json:"payment_method" valid:"required|in:cash,debit_card,credit_card,giro"`
	Amount          float64               `json:"amount" valid:"required"`
	BankName        string                `json:"bank_name"`
	BankNumber      string                `json:"bank_number"`
	BankHolder      string                `json:"bank_holder"`
	Note            string                `json:"note"`
	Session         *auth.SessionData     `json:"-"`
	Revenue         *model.FinanceRevenue `json:"-"`
}

// Validate implement validation.Requests interfaces.
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if r.PaymentMethod == "debit_card" || r.PaymentMethod == "credit_card" {
		if r.BankNumber == "" {
			o.Failure("bank_number", "Bank number is required")
		}
		if r.BankName == "" {
			o.Failure("bank_name", "Bank name is required")
		}
		if r.BankHolder == "" {
			o.Failure("bank_holder", "Bank holder is required")
		}
	} else if r.PaymentMethod == "giro" {
		if r.BankNumber == "" {
			o.Failure("bank_number", "Bank number is required")
		}
		if r.BankName != "" {
			o.Failure("bank_name", "Bank name can't be inputted if payment method is giro")
		}
	} else {
		if r.BankNumber != "" {
			o.Failure("bank_number", "Bank number can't be inputted if payment method is cash")
		}
		if r.BankName != "" {
			o.Failure("bank_name", "Bank name can't be inputted if payment method is cash")
		}
	}

	if r.Revenue.DocumentStatus == "cleared" {
		o.Failure("document_status", "Document status has been cleared")
	}

	if r.Revenue.RefType == "sales_invoice" {
		// total Amount finance revenue berdasarkan reftype dan ref id
		totAmount := SumAmountFinanceRevenue(r.Revenue.RefID, "sales_invoice", r.Revenue.ID)
		totAmount += r.Amount
		// sales invoice dari data
		si := &model.SalesInvoice{ID: int64(r.Revenue.RefID)}
		si.Read()

		if totAmount > si.TotalAmount {
			o.Failure("sales_invoice", "the amount of the pay is greater than the total charge of the sales_invoice")
		}
	} else {
		// total Amount finance revenue berdasarkan reftype dan ref id
		totAmount := SumAmountFinanceRevenue(r.Revenue.RefID, "purchase_return", r.Revenue.ID)
		totAmount += r.Amount
		// purchase return dari data
		pr := &model.PurchaseReturn{ID: int64(r.Revenue.RefID)}
		pr.Read()
		if totAmount > pr.TotalAmount {
			o.Failure("purchase_return", "the amount of the pay is greater than the total charge of the purchase_return")
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
func (r *updateRequest) Transform() *model.FinanceRevenue {
	revenue := r.Revenue
	revenue.RecognitionDate = r.RecognitionDate
	revenue.PaymentMethod = r.PaymentMethod
	revenue.Amount = r.Amount
	revenue.BankName = r.BankName
	revenue.BankHolder = r.BankHolder
	revenue.BankNumber = r.BankNumber
	revenue.Note = r.Note
	revenue.UpdatedBy = r.Session.User
	revenue.UpdatedAt = time.Now()

	return revenue
}

// approveRequest untuk menampung model revenue untuk validasi
type approveRequest struct {
	Revenue *model.FinanceRevenue `json:"-"`
}

// Validate implement validation.Requests interfaces.
func (r *approveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	if r.Revenue.DocumentStatus == "cleared" {
		o.Failure("document_status", "Document status has been cleared")
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *approveRequest) Messages() map[string]string {
	return map[string]string{}
}
