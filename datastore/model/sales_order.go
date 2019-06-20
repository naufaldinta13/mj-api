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
	orm.RegisterModel(new(SalesOrder))
}

// SalesOrder model for sales_order table.
type SalesOrder struct {
	ID                   int64        `orm:"column(id);auto" json:"-"`
	Reference            *SalesOrder  `orm:"column(reference_id);null;rel(fk)" json:"reference,omitempty"`
	Customer             *Partnership `orm:"column(customer_id);null;rel(fk)" json:"customer,omitempty"`
	Code                 string       `orm:"column(code);size(45)" json:"code"`
	RecognitionDate      time.Time    `orm:"column(recognition_date);type(date);null" json:"recognition_date"`
	EtaDate              time.Time    `orm:"column(eta_date);type(date);null" json:"eta_date"`
	Discount             float32      `orm:"column(discount);null" json:"discount"`
	Tax                  float32      `orm:"column(tax);null" json:"tax"`
	DiscountAmount       float64      `orm:"column(discount_amount);null;digits(20);decimals(0)" json:"discount_amount"`
	TaxAmount            float64      `orm:"column(tax_amount);null;digits(20);decimals(0)" json:"tax_amount"`
	ShipmentAddress      string       `orm:"column(shipment_address);null" json:"shipment_address"`
	ShipmentCost         float64      `orm:"column(shipment_cost);null;digits(20);decimals(0)" json:"shipment_cost"`
	TotalPrice           float64      `orm:"column(total_price);digits(20);decimals(0)" json:"total_price"`
	TotalCharge          float64      `orm:"column(total_charge);digits(20);decimals(0)" json:"total_charge"`
	TotalPaid            float64      `orm:"column(total_paid);null;digits(20);decimals(0)" json:"total_paid"`
	TotalCost            float64      `orm:"column(total_cost);null;digits(20);decimals(0)" json:"total_cost"`
	Note                 string       `orm:"column(note);null" json:"note"`
	DocumentStatus       string       `orm:"column(document_status);null;options(new,active,finished,requested_cancel,approved_cancel)" json:"document_status"`
	InvoiceStatus        string       `orm:"column(invoice_status);null;options(new,active,finished)" json:"invoice_status"`
	FulfillmentStatus    string       `orm:"column(fulfillment_status);null;options(new,active,finished)" json:"fulfillment_status"`
	ShipmentStatus       string       `orm:"column(shipment_status);null;options(new,active,finished)" json:"shipment_status"`
	AutoFulfillment      int8         `orm:"column(auto_fulfillment);null" json:"auto_fulfillment"`
	AutoInvoice          int8         `orm:"column(auto_invoice);null" json:"auto_invoice"`
	AutoPaid             int8         `orm:"column(auto_paid);null" json:"auto_paid"`
	IsPercentageDiscount int8         `orm:"column(is_percentage_discount);null" json:"is_percentage_discount"`
	IsDeleted            int8         `orm:"column(is_deleted);null" json:"is_deleted"`
	CreatedBy            *User        `orm:"column(created_by);rel(fk)" json:"created_by"`
	UpdatedBy            *User        `orm:"column(updated_by);null;rel(fk)" json:"updated_by"`
	RequestCancelBy      *User        `orm:"column(request_cancel_by);null;rel(fk)" json:"request_cancel_by"`
	ApproveCancelBy      *User        `orm:"column(approve_cancel_by);null;rel(fk)" json:"approve_cancel_by"`
	CreatedAt            time.Time    `orm:"column(created_at);type(timestamp)" json:"created_at"`
	UpdatedAt            time.Time    `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
	RequestCancelAt      time.Time    `orm:"column(request_cancel_at);type(timestamp);null" json:"request_cancel_at"`
	ApproveCancelAt      time.Time    `orm:"column(approve_cancel_at);type(timestamp);null" json:"approve_cancel_at"`
	CancelledNote        string       `orm:"column(cancelled_note);null" json:"cancelled_note"`
	IsReported           int8         `orm:"column(is_reported);null" json:"is_reported"`

	TotalRefund           float64                 `orm:"-" json:"total_refund,omitempty"`
	TotalPaidRefund       float64                 `orm:"-" json:"total_paid_refund,omitempty"`
	SalesOrderItems       []*SalesOrderItem       `orm:"-" json:"sales_order_items,omitempty"`
	SalesInvoices         []*SalesInvoice         `orm:"-" json:"sales_invoices,omitempty"`
	WorkorderFulfillments []*WorkorderFulfillment `orm:"-" json:"workorder_fulfillments,omitempty"`
	SalesReturns          []*SalesReturn          `orm:"-" json:"sales_returns,omitempty"`
	InvoiceReceiptItems   []*InvoiceReceiptItem   `orm:"-" json:"invoice_receipt_items,omitempty"`
	InvoiceReceiptReturns []*InvoiceReceiptReturn `orm:"-" json:"invoice_receipt_returns,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *SalesOrder) MarshalJSON() ([]byte, error) {
	type Alias SalesOrder

	alias := &struct {
		ID                string `json:"id"`
		UpdatedByID       string `json:"updated_by_id"`
		RequestCancelByID string `json:"request_cancel_by_id"`
		ApproveCancelByID string `json:"approve_cancel_by_id"`
		ReferenceID       string `json:"reference_id"`
		CustomerID        string `json:"customer_id"`
		BankAccountID     string `json:"bank_account_id"`
		CreatedByID       string `json:"created_by_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.ReferenceID when m.Reference not nill
	// and the ID is setted
	if m.Reference != nil && m.Reference.ID != int64(0) {
		alias.ReferenceID = common.Encrypt(m.Reference.ID)
	} else {
		alias.Reference = nil
	}

	// Encrypt alias.CustomerID when m.Customer not nill
	// and the ID is setted
	if m.Customer != nil && m.Customer.ID != int64(0) {
		alias.CustomerID = common.Encrypt(m.Customer.ID)
	} else {
		alias.Customer = nil
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

	// Encrypt alias.RequestCancelByID when m.RequestCancelBy not nill
	// and the ID is setted
	if m.RequestCancelBy != nil && m.RequestCancelBy.ID != int64(0) {
		alias.RequestCancelByID = common.Encrypt(m.RequestCancelBy.ID)
	} else {
		alias.RequestCancelBy = nil
	}

	// Encrypt alias.ApproveCancelByID when m.ApproveCancelBy not nill
	// and the ID is setted
	if m.ApproveCancelBy != nil && m.ApproveCancelBy.ID != int64(0) {
		alias.ApproveCancelByID = common.Encrypt(m.ApproveCancelBy.ID)
	} else {
		alias.ApproveCancelBy = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating SalesOrder struct into sales_order table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to sales_order.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *SalesOrder) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting sales_order data
// this also will truncated all data from all table
// that have relation with this sales_order.
func (m *SalesOrder) Delete() (err error) {
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
func (m *SalesOrder) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
