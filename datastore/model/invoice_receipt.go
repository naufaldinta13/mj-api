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
	orm.RegisterModel(new(InvoiceReceipt))
}

// InvoiceReceipt model for invoice_receipt table.
type InvoiceReceipt struct {
	ID                    int64                   `orm:"column(id);auto" json:"-"`
	Partnership           *Partnership            `orm:"column(partnership_id);rel(fk)" json:"partnership,omitempty"`
	Code                  string                  `orm:"column(code);size(120)" json:"code"`
	RecognitionDate       time.Time               `orm:"column(recognition_date);type(date);null" json:"recognition_date"`
	TotalInvoice          float64                 `orm:"column(total_invoice);null;digits(20);decimals(0)" json:"total_invoice"`
	TotalReturn           float64                 `orm:"column(total_return);null;digits(20);decimals(0)" json:"total_return"`
	TotalAmount           float64                 `orm:"column(total_amount);null;digits(20);decimals(0)" json:"total_amount"`
	Note                  string                  `orm:"column(note);null" json:"note"`
	DocumentStatus        string                  `orm:"column(document_status);null;options(new,active,finished)" json:"document_status"`
	CreatedBy             *User                   `orm:"column(created_by);rel(fk)" json:"created_by"`
	UpdatedBy             *User                   `orm:"column(updated_by);null;rel(fk)" json:"updated_by"`
	CreatedAt             time.Time               `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	UpdatedAt             time.Time               `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
	IsDeleted             int8                    `orm:"column(is_deleted);null" json:"is_deleted"`
	InvoiceReceiptItems   []*InvoiceReceiptItem   `orm:"reverse(many)" json:"invoice_receipt_items,omitempty"`
	InvoiceReceiptReturns []*InvoiceReceiptReturn `orm:"reverse(many)" json:"invoice_receipt_returns,omitempty"`
	SalesOrders           []*SalesOrder           `orm:"-" json:"sales_orders,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *InvoiceReceipt) MarshalJSON() ([]byte, error) {
	type Alias InvoiceReceipt

	alias := &struct {
		ID            string `json:"id"`
		CreatedByID   string `json:"created_by_id"`
		UpdatedByID   string `json:"updated_by_id"`
		PartnershipID string `json:"partnership_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.PartnershipID when m.Partnership not nill
	// and the ID is setted
	if m.Partnership != nil && m.Partnership.ID != int64(0) {
		alias.PartnershipID = common.Encrypt(m.Partnership.ID)
	} else {
		alias.Partnership = nil
	}

	// Encrypt alias.CreatedByID when m.CreatedBy not nill
	// and the ID is setted
	if m.CreatedBy != nil && m.CreatedBy.ID != int64(0) {
		alias.CreatedByID = common.Encrypt(m.CreatedBy.ID)
	} else {
		alias.CreatedBy = nil
	}

	// Encrypt alias.UpdatedByID when m.UpdatedBy not nill
	// and the ID is setted
	if m.UpdatedBy != nil && m.UpdatedBy.ID != int64(0) {
		alias.UpdatedByID = common.Encrypt(m.UpdatedBy.ID)
	} else {
		alias.UpdatedBy = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating InvoiceReceipt struct into invoice_receipt table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to invoice_receipt.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *InvoiceReceipt) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting invoice_receipt data
// this also will truncated all data from all table
// that have relation with this invoice_receipt.
func (m *InvoiceReceipt) Delete() (err error) {
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
func (m *InvoiceReceipt) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
