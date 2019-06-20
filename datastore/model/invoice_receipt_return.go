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
	orm.RegisterModel(new(InvoiceReceiptReturn))
}

// InvoiceReceiptReturn model for invoice_receipt_return table.
type InvoiceReceiptReturn struct {
	ID             int64           `orm:"column(id);auto" json:"-"`
	InvoiceReceipt *InvoiceReceipt `orm:"column(invoice_receipt_id);rel(fk)" json:"invoice_receipt,omitempty"`
	SalesReturn    *SalesReturn    `orm:"column(sales_return_id);rel(fk)" json:"sales_return,omitempty"`
	Subtotal       float64         `orm:"column(subtotal);digits(20);decimals(0)" json:"subtotal"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *InvoiceReceiptReturn) MarshalJSON() ([]byte, error) {
	type Alias InvoiceReceiptReturn

	alias := &struct {
		ID               string `json:"id"`
		InvoiceReceiptID string `json:"invoice_receipt_id"`
		SalesReturnID    string `json:"sales_return_id"`
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

	// Encrypt alias.SalesReturnID when m.SalesReturn not nill
	// and the ID is setted
	if m.SalesReturn != nil && m.SalesReturn.ID != int64(0) {
		alias.SalesReturnID = common.Encrypt(m.SalesReturn.ID)
	} else {
		alias.SalesReturn = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating InvoiceReceiptReturn struct into invoice_receipt_return table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to invoice_receipt_return.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *InvoiceReceiptReturn) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting invoice_receipt_return data
// this also will truncated all data from all table
// that have relation with this invoice_receipt_return.
func (m *InvoiceReceiptReturn) Delete() (err error) {
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
func (m *InvoiceReceiptReturn) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
