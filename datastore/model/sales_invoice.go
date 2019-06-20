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
	orm.RegisterModel(new(SalesInvoice))
}

// SalesInvoice model for sales_invoice table.
type SalesInvoice struct {
	ID              int64       `orm:"column(id);auto" json:"-"`
	SalesOrder      *SalesOrder `orm:"column(sales_order_id);rel(fk)" json:"sales_order,omitempty"`
	Code            string      `orm:"column(code);size(45)" json:"code"`
	RecognitionDate time.Time   `orm:"column(recognition_date);type(date);null" json:"recognition_date"`
	DueDate         time.Time   `orm:"column(due_date);type(date);null" json:"due_date"`
	BillingAddress  string      `orm:"column(billing_address);null" json:"billing_address"`
	TotalAmount     float64     `orm:"column(total_amount);null;digits(20);decimals(0)" json:"total_amount"`
	TotalPaid       float64     `orm:"column(total_paid);null;digits(20);decimals(0)" json:"total_paid"`
	TotalRevenued   float64     `orm:"column(total_revenued);null;digits(20);decimals(0)" json:"total_revenued"`
	Note            string      `orm:"column(note);null" json:"note"`
	DocumentStatus  string      `orm:"column(document_status);null;options(new,active,finished)" json:"document_status"`
	IsBundled       int8        `orm:"column(is_bundled);null" json:"is_bundled"`
	IsDeleted       int8        `orm:"column(is_deleted);null" json:"is_deleted"`
	CreatedBy       *User       `orm:"column(created_by);rel(fk)" json:"created_by"`
	UpdatedBy       *User       `orm:"column(updated_by);null;rel(fk)" json:"updated_by"`
	CreatedAt       time.Time   `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	UpdatedAt       time.Time   `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *SalesInvoice) MarshalJSON() ([]byte, error) {
	type Alias SalesInvoice

	alias := &struct {
		ID           string `json:"id"`
		UpdatedByID  string `json:"updated_by_id"`
		SalesOrderID string `json:"sales_order_id"`
		CreatedByID  string `json:"created_by_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.SalesOrderID when m.SalesOrder not nill
	// and the ID is setted
	if m.SalesOrder != nil && m.SalesOrder.ID != int64(0) {
		alias.SalesOrderID = common.Encrypt(m.SalesOrder.ID)
	} else {
		alias.SalesOrder = nil
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

// Save inserting or updating SalesInvoice struct into sales_invoice table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to sales_invoice.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *SalesInvoice) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting sales_invoice data
// this also will truncated all data from all table
// that have relation with this sales_invoice.
func (m *SalesInvoice) Delete() (err error) {
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
func (m *SalesInvoice) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
