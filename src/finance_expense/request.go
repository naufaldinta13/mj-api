package financeExpense

import (
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/sales_return"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create category process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	RefType         string    `json:"ref_type" valid:"required|in:purchase_invoice,sales_return"`
	RefID           string    `json:"ref_id" valid:"required"`
	RecognitionDate time.Time `json:"recognition_date" valid:"required"`
	PaymentMethod   string    `json:"payment_method" valid:"required|in:cash,debit_card,credit_card,giro"`
	BankNumber      string    `json:"bank_number"`
	BankName        string    `json:"bank_name"`
	BankHolder      string    `json:"bank_holder"`
	GiroNumber      string    `json:"giro_number"`
	Amount          float64   `json:"amount" valid:"required"`
	Note            string    `json:"note"`
	Session         *auth.SessionData
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	refID, er := common.Decrypt(r.RefID)
	if er != nil {
		o.Failure("ref_id", "Ref id is invalid")
	}
	// Ambil ref_type yang dipilih
	// ambil document_status dari purchase_invoice yang dipilih
	if r.RefType == "purchase_invoice" {
		var pinvoice model.PurchaseInvoice

		query := orm.NewOrm().QueryTable(new(model.PurchaseInvoice)).RelatedSel()
		err := query.Filter("id", refID).Filter("is_deleted", 0).Limit(1).One(&pinvoice)
		if err != nil {
			o.Failure("ref_id", "Ref id for purchase invoice is not found")
		} else {
			// cek document status purchase invoice == finished?
			if pinvoice.DocumentStatus == "finished" {
				o.Failure("ref_id", "Purchase invoice already paid")
			} else {
				// jumlahkan semua amount dari finance_expense yang memiliki referensi dari purchase invoice yang diinput
				var sumFE float64
				orm.NewOrm().Raw("SELECT SUM(fe.amount) AS sumFE FROM finance_expense fe WHERE fe.ref_id = ? AND fe.ref_type = ?;", pinvoice.ID, "purchase_invoice").QueryRow(&sumFE)
				// cek jumlah amount diatas  > total amount di purchase invoice
				if (sumFE + r.Amount) > pinvoice.TotalAmount {
					o.Failure("amount", "You pay too much from total charge in purchase order")
				} else {
					// update document status pada purchase invoice yang dipilih menjadi active
					pinvoice.DocumentStatus = "active"
					pinvoice.Save("DocumentStatus")
				}
			}
		}
	} else {
		var sreturn *model.SalesReturn
		sreturn, err := salesReturn.ShowSalesReturn("id", refID)
		if err != nil && sreturn == nil || sreturn.IsDeleted == int8(1) || sreturn.DocumentStatus == "cancelled" {
			o.Failure("ref_id", "Ref id for sales return is not found")
		} else {
			// ambil document status sales return yang dipilih
			// check document status == finished
			if sreturn.DocumentStatus == "finished" {
				o.Failure("ref_id", "Sales return already paid")
			} else {
				// jumlahkan semua amount dari expense yang memiliki referensi dari sales return yang diinputkan
				var sumSR float64
				orm.NewOrm().Raw("SELECT SUM(fe.amount) AS amount FROM finance_expense fe WHERE fe.ref_id = ? AND fe.ref_type = ?;", sreturn.ID, "sales_return").QueryRow(&sumSR)
				// check jumlah amount diatas > total amount sales return
				if (sumSR + r.Amount) > sreturn.TotalAmount {
					o.Failure("amount", "You pay too much from total amount in sales return")
				} else {
					// update document status pada sales return yang dipilih menjadi active
					sreturn.DocumentStatus = "active"
					sreturn.Save("DocumentStatus")
				}
			}
		}
	}

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
		if r.GiroNumber != "" {
			o.Failure("giro_number", "Giro number can't be inputted if payment method is debit card")
		}
	} else if r.PaymentMethod == "giro" {
		if r.GiroNumber == "" {
			o.Failure("giro_number", "Giro number is required")
		}
		if r.BankNumber != "" {
			o.Failure("bank_number", "Bank number can't be inputted if payment method is giro")
		}
		if r.BankName != "" {
			o.Failure("bank_name", "Bank name can't be inputted if payment method is giro")
		}
		if r.BankHolder != "" {
			o.Failure("bank_holder", "Bank holder can't be inputted if payment method is giro")
		}
	} else {
		if r.BankNumber != "" {
			o.Failure("bank_number", "Bank number can't be inputted if payment method is cash")
		}
		if r.BankName != "" {
			o.Failure("bank_name", "Bank name can't be inputted if payment method is cash")
		}
		if r.BankHolder != "" {
			o.Failure("bank_holder", "Bank holder can't be inputted if payment method is cash")
		}
		if r.GiroNumber != "" {
			o.Failure("giro_number", "Giro number can't be inputted if payment method is cash")
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
func (r *createRequest) Transform() *model.FinanceExpense {
	var fexpense model.FinanceExpense
	refID, _ := common.Decrypt(r.RefID)

	fexpense = model.FinanceExpense{
		RefType:         r.RefType,
		RefID:           uint64(refID),
		RecognitionDate: r.RecognitionDate,
		Amount:          r.Amount,
		PaymentMethod:   r.PaymentMethod,
		BankName:        r.BankName,
		BankHolder:      r.BankHolder,
		DocumentStatus:  "uncleared",
		Note:            r.Note,
		CreatedAt:       time.Now(),
		CreatedBy:       r.Session.User,
		IsDeleted:       int8(0),
	}

	if r.PaymentMethod == "giro" {
		fexpense.BankNumber = r.GiroNumber
	} else {
		fexpense.BankNumber = r.BankNumber
	}

	return &fexpense
}

type updateRequest struct {
	ID              int64
	RecognitionDate time.Time `json:"recognition_date" valid:"required"`
	PaymentMethod   string    `json:"payment_method" valid:"required|in:cash,debit_card,credit_card,giro"`
	BankNumber      string    `json:"bank_number"`
	BankName        string    `json:"bank_name"`
	BankHolder      string    `json:"bank_holder"`
	GiroNumber      string    `json:"giro_number"`
	Amount          float64   `json:"amount" valid:"required"`
	Note            string    `json:"note"`
	FinanceExpense  *model.FinanceExpense
	Session         *auth.SessionData
}

// Validate implement validation.Requests interfaces.
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if r.FinanceExpense.DocumentStatus == "cleared" {
		o.Failure("document_status", "Finance expense document status already cleared")
	}

	// Ambil ref_type yang dipilih
	// ambil document_status dari purchase_invoice yang dipilih
	if r.FinanceExpense.RefType == "purchase_invoice" {
		var pinvoice model.PurchaseInvoice

		query := orm.NewOrm().QueryTable(new(model.PurchaseInvoice)).RelatedSel()
		err := query.Filter("id", r.FinanceExpense.RefID).Filter("is_deleted", 0).Limit(1).One(&pinvoice)
		if err != nil {
			o.Failure("ref_id", "Ref id for purchase invoice is not found")
		} else {
			// cek document status purchase invoice == finished?
			if pinvoice.DocumentStatus == "finished" {
				o.Failure("ref_id", "Purchase invoice already paid")
			} else {
				// jumlahkan semua amount dari finance_expense yang memiliki referensi dari purchase invoice yang diinput
				var sumFE float64
				orm.NewOrm().Raw("SELECT SUM(fe.amount) AS sumFE FROM finance_expense fe WHERE fe.ID !=? fe.ref_id = ? AND fe.ref_type = ?;", r.FinanceExpense.ID, pinvoice.ID, "purchase_invoice").QueryRow(&sumFE)
				// cek jumlah amount diatas  > total amount di purchase invoice
				if (sumFE + r.Amount) > pinvoice.TotalAmount {
					o.Failure("ref_id", "You pay too much from total charge in purchase order")
				} else {
					// update document status pada purchase invoice yang dipilih menjadi active
					pinvoice.DocumentStatus = "active"
					pinvoice.Save("DocumentStatus")
				}
			}
		}
	} else {
		var sreturn *model.SalesReturn
		sreturn, err := salesReturn.ShowSalesReturn("id", r.FinanceExpense.RefID)
		if err != nil || sreturn.IsDeleted == int8(1) || sreturn.DocumentStatus == "cancelled" {
			o.Failure("ref_id", "Ref id for sales return is not found")
		} else {
			// ambil document status sales return yang dipilih
			// check document status == finished
			if sreturn.DocumentStatus == "finished" {
				o.Failure("ref_id", "Sales return already paid")
			} else {
				// jumlahkan semua amount dari expense yang memiliki referensi dari sales return yang diinputkan
				var sumSR float64
				orm.NewOrm().Raw("SELECT SUM(fe.amount) AS amount FROM finance_expense fe WHERE fe.ID != ? fe.ref_id = ? AND fe.ref_type = ?;", r.FinanceExpense.ID, sreturn.ID, "sales_return").QueryRow(&sumSR)
				// check jumlah amount diatas > total amount sales return
				if (sumSR + r.Amount) > sreturn.TotalAmount {
					o.Failure("ref_id", "You pay too much from total amount in sales return")
				} else {
					// update document status pada sales return yang dipilih menjadi active
					sreturn.DocumentStatus = "active"
					sreturn.Save("DocumentStatus")
				}
			}
		}
	}

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
		if r.GiroNumber != "" {
			o.Failure("giro_number", "Giro number can't be inputted if payment method is debit card")
		}
	} else if r.PaymentMethod == "giro" {
		if r.GiroNumber == "" {
			o.Failure("giro_number", "Giro number is required")
		}
		if r.BankNumber != "" {
			o.Failure("bank_number", "Bank number can't be inputted if payment method is giro")
		}
		if r.BankName != "" {
			o.Failure("bank_name", "Bank name can't be inputted if payment method is giro")
		}
		if r.BankHolder != "" {
			o.Failure("bank_holder", "Bank holder can't be inputted if payment method is giro")
		}
	} else {
		if r.BankNumber != "" {
			o.Failure("bank_number", "Bank number can't be inputted if payment method is cash")
		}
		if r.BankName != "" {
			o.Failure("bank_name", "Bank name can't be inputted if payment method is cash")
		}
		if r.BankHolder != "" {
			o.Failure("bank_holder", "Bank holder can't be inputted if payment method is cash")
		}
		if r.GiroNumber != "" {
			o.Failure("giro_number", "Giro number can't be inputted if payment method is cash")
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
func (r *updateRequest) Transform() *model.FinanceExpense {
	var fexpense model.FinanceExpense

	fexpense = model.FinanceExpense{
		ID:              r.FinanceExpense.ID,
		RecognitionDate: r.RecognitionDate,
		Amount:          r.Amount,
		PaymentMethod:   r.PaymentMethod,
		BankNumber:      r.BankNumber,
		BankName:        r.BankName,
		BankHolder:      r.BankHolder,
		Note:            r.Note,
		UpdatedAt:       time.Now(),
		UpdatedBy:       r.Session.User,
	}

	if r.PaymentMethod == "giro" {
		fexpense.BankNumber = r.GiroNumber
	} else {
		fexpense.BankNumber = r.BankNumber
	}

	return &fexpense
}

// approveRequest struct untuk menampung model expense
type approveRequest struct {
	FinanceExpense *model.FinanceExpense
}

// Validate implement validation.Requests interfaces.
func (r *approveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	if r.FinanceExpense.DocumentStatus == "cleared" {
		o.Failure("document_status", "Finance expense document status already cleared")
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *approveRequest) Messages() map[string]string {
	return map[string]string{}
}
