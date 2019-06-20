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
	orm.RegisterModel(new(PurchaseOrder))
}

// PurchaseOrder model for purchase_order table.
type PurchaseOrder struct {
	ID                  int64                 `orm:"column(id);auto" json:"-"`
	Reference           *PurchaseOrder        `orm:"column(reference_id);null;rel(fk)" json:"reference,omitempty"`
	Supplier            *Partnership          `orm:"column(supplier_id);rel(fk)" json:"supplier,omitempty"`
	Code                string                `orm:"column(code);size(45)" json:"code"`
	RecognitionDate     time.Time             `orm:"column(recognition_date);type(date);null" json:"recognition_date"`
	EtaDate             time.Time             `orm:"column(eta_date);type(date);null" json:"eta_date"`
	Discount            float32               `orm:"column(discount);null" json:"discount"`
	Tax                 float32               `orm:"column(tax);null" json:"tax"`
	DiscountAmount      float64               `orm:"column(discount_amount);null;digits(20);decimals(0)" json:"discount_amount"`
	TaxAmount           float64               `orm:"column(tax_amount);null;digits(20);decimals(0)" json:"tax_amount"`
	ShipmentCost        float64               `orm:"column(shipment_cost);null;digits(20);decimals(0)" json:"shipment_cost"`
	TotalCharge         float64               `orm:"column(total_charge);null;digits(20);decimals(0)" json:"total_charge"`
	TotalPaid           float64               `orm:"column(total_paid);null;digits(20);decimals(0)" json:"total_paid"`
	Note                string                `orm:"column(note);null" json:"note"`
	DocumentStatus      string                `orm:"column(document_status);null;options(new,active,finished,cancelled)" json:"document_status"`
	InvoiceStatus       string                `orm:"column(invoice_status);null;options(new,active,finished)" json:"invoice_status"`
	ReceivingStatus     string                `orm:"column(receiving_status);null;options(new,active,finished)" json:"receiving_status"`
	AutoInvoiced        int8                  `orm:"column(auto_invoiced);null" json:"auto_invoiced"`
	IsPercentage        int8                  `orm:"column(is_percentage);null" json:"is_percentage"`
	IsDeleted           int8                  `orm:"column(is_deleted);null" json:"is_deleted"`
	CreatedBy           *User                 `orm:"column(created_by);rel(fk)" json:"created_by"`
	UpdatedBy           *User                 `orm:"column(updated_by);null;rel(fk)" json:"updated_by"`
	CreatedAt           time.Time             `orm:"column(created_at);type(timestamp)" json:"created_at"`
	UpdatedAt           time.Time             `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
	PurchaseOrderItems  []*PurchaseOrderItem  `orm:"reverse(many)" json:"purchase_order_items,omitempty"`
	PurchaseInvoices    []*PurchaseInvoice    `orm:"reverse(many)" json:"purchase_invoices,omitempty"`
	PurchaseReturns     []*PurchaseReturn     `orm:"reverse(many)" json:"purchase_returns,omitempty"`
	WorkorderReceivings []*WorkorderReceiving `orm:"reverse(many)" json:"workorder_receivings,omitempty"`
	CancelledNote       string                `orm:"column(cancelled_note);null" json:"cancelled_note"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PurchaseOrder) MarshalJSON() ([]byte, error) {
	type Alias PurchaseOrder

	alias := &struct {
		ID          string `json:"id"`
		ReferenceID string `json:"reference_id"`
		SupplierID  string `json:"supplier_id"`
		CreatedByID string `json:"created_by_id"`
		UpdatedByID string `json:"updated_by_id"`
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

	// Encrypt alias.SupplierID when m.Supplier not nill
	// and the ID is setted
	if m.Supplier != nil && m.Supplier.ID != int64(0) {
		alias.SupplierID = common.Encrypt(m.Supplier.ID)
	} else {
		alias.Supplier = nil
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

// Save inserting or updating PurchaseOrder struct into purchase_order table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to purchase_order.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *PurchaseOrder) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting purchase_order data
// this also will truncated all data from all table
// that have relation with this purchase_order.
func (m *PurchaseOrder) Delete() (err error) {
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
func (m *PurchaseOrder) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
