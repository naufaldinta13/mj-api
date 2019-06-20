// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(InvoiceReceiptItem))
}

// InvoiceReceiptItem model for invoice_receipt_item table.
type InvoiceReceiptItem struct {
	ID             int64           `orm:"column(id);auto" json:"-"`
	InvoiceReceipt *InvoiceReceipt `orm:"column(invoice_receipt_id);rel(fk)" json:"invoice_receipt,omitempty"`
	SalesInvoice   *SalesInvoice   `orm:"column(sales_invoice_id);rel(fk)" json:"sales_invoice,omitempty"`
	Subtotal       float64         `orm:"column(subtotal);digits(20);decimals(0)" json:"subtotal"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *InvoiceReceiptItem) MarshalJSON() ([]byte, error) {
	type Alias InvoiceReceiptItem

	alias := &struct {
		ID               string `json:"id"`
		InvoiceReceiptID string `json:"invoice_receipt_id"`
		SalesInvoiceID   string `json:"sales_invoice_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.InvoiceReceiptID when m.InvoiceReceipt not nill
	// and the ID is setted
	if m.InvoiceReceipt != nil && m.InvoiceReceipt.ID != int64(0) {
		alias.InvoiceReceiptID = common.Encrypt(m.InvoiceReceipt.ID)
	} else {
		alias.InvoiceReceipt = nil
	}

	// Encrypt alias.SalesInvoiceID when m.SalesInvoice not nill
	// and the ID is setted
	if m.SalesInvoice != nil && m.SalesInvoice.ID != int64(0) {
		alias.SalesInvoiceID = common.Encrypt(m.SalesInvoice.ID)
	} else {
		alias.SalesInvoice = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating InvoiceReceiptItem struct into invoice_receipt_item table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to invoice_receipt_item.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *InvoiceReceiptItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting invoice_receipt_item data
// this also will truncated all data from all table
// that have relation with this invoice_receipt_item.
func (m *InvoiceReceiptItem) Delete() (err error) {
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
func (m *InvoiceReceiptItem) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
