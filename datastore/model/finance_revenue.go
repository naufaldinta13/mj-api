// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(FinanceRevenue))
}

// FinanceRevenue model for finance_revenue table.
type FinanceRevenue struct {
	ID              int64           `orm:"column(id);auto" json:"-"`
	BankAccount     *BankAccount    `orm:"column(bank_account_id);null;rel(fk)" json:"bank_account,omitempty"`
	RefID           uint64          `orm:"column(ref_id);null" json:"ref_id"`
	RecognitionDate time.Time       `orm:"column(recognition_date);type(date);null" json:"recognition_date"`
	RefType         string          `orm:"column(ref_type);null;options(sales_invoice,purchase_return,invoice_receipt)" json:"ref_type"`
	Amount          float64         `orm:"column(amount);digits(20);decimals(0)" json:"amount"`
	PaymentMethod   string          `orm:"column(payment_method);null;options(cash,debit_card,credit_card,giro)" json:"payment_method"`
	BankName        string          `orm:"column(bank_name);size(45);null" json:"bank_name"`
	BankNumber      string          `orm:"column(bank_number);size(45);null" json:"bank_number"`
	BankHolder      string          `orm:"column(bank_holder);size(45);null" json:"bank_holder"`
	Note            string          `orm:"column(note);null" json:"note"`
	DocumentStatus  string          `orm:"column(document_status);null;options(uncleared,cleared)" json:"document_status"`
	CreatedBy       *User           `orm:"column(created_by);rel(fk)" json:"created_by"`
	UpdatedBy       *User           `orm:"column(updated_by);null;rel(fk)" json:"updated_by"`
	CreatedAt       time.Time       `orm:"column(created_at);type(timestamp)" json:"created_at"`
	UpdatedAt       time.Time       `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
	IsDeleted       int8            `orm:"column(is_deleted);null" json:"is_deleted"`
	InvoiceReceipt  *InvoiceReceipt `orm:"-" json:"invoice_receipt,omitempty"`
	SalesInvoice    *SalesInvoice   `orm:"-" json:"sales_invoice,omitempty"`
	PurchaseReturn  *PurchaseReturn `orm:"-" json:"purchase_return,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *FinanceRevenue) MarshalJSON() ([]byte, error) {
	type Alias FinanceRevenue

	alias := &struct {
		ID            string `json:"id"`
		BankAccountID string `json:"bank_account_id"`
		CreatedByID   string `json:"created_by_id"`
		UpdatedByID   string `json:"updated_by_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.BankAccountID when m.BankAccount not nill
	// and the ID is setted
	if m.BankAccount != nil && m.BankAccount.ID != int64(0) {
		alias.BankAccountID = common.Encrypt(m.BankAccount.ID)
	} else {
		alias.BankAccount = nil
	}

	// Encrypt alias.UpdatedByID when m.UpdatedBy not nill
	// and the ID is setted
	if m.UpdatedBy != nil && m.UpdatedBy.ID != int64(0) {
		alias.UpdatedByID = common.Encrypt(m.UpdatedBy.ID)
	} else {
		alias.UpdatedBy = nil
	}

	// Encrypt alias.CreatedByID when m.CreatedBy not nill
	// and the ID is setted
	if m.CreatedBy != nil && m.CreatedBy.ID != int64(0) {
		alias.CreatedByID = common.Encrypt(m.CreatedBy.ID)
	} else {
		alias.CreatedBy = nil
	}

	o := orm.NewOrm()
	if m.RefType == "invoice_receipt" {
		o.Raw("select * from invoice_receipt where id = ?", int64(m.RefID)).QueryRow(&m.InvoiceReceipt)
		if m.InvoiceReceipt != nil {
			m.InvoiceReceipt.Partnership.Read()
		}
	}

	if m.RefType == "sales_invoice" {
		o.Raw("select * from sales_invoice where id = ?", int64(m.RefID)).QueryRow(&m.SalesInvoice)
		if m.SalesInvoice != nil {
			m.SalesInvoice.SalesOrder.Read()

			if m.SalesInvoice.SalesOrder.Customer != nil {
				m.SalesInvoice.SalesOrder.Customer.Read()
			}
		}
	}

	if m.RefType == "purchase_return" {
		o.Raw("select * from purchase_return where id = ?", int64(m.RefID)).QueryRow(&m.PurchaseReturn)
		if m.PurchaseReturn != nil {
			m.PurchaseReturn.PurchaseOrder.Read()
			m.PurchaseReturn.PurchaseOrder.Supplier.Read()
		}
	}

	return json.Marshal(alias)
}

// Save inserting or updating FinanceRevenue struct into finance_revenue table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to finance_revenue.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *FinanceRevenue) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting finance_revenue data
// this also will truncated all data from all table
// that have relation with this finance_revenue.
func (m *FinanceRevenue) Delete() (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		var i int64
		if i, err = o.Delete(m); i == 0 && err == nil {
			err = orm.ErrNoAffected
		}
		return
	}
	return orm.ErrNoRows
}

// Read execute select based on data struct that already
// assigned.
func (m *FinanceRevenue) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
