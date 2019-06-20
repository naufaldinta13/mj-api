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
	orm.RegisterModel(new(PurchaseReturn))
}

// PurchaseReturn model for purchase_return table.
type PurchaseReturn struct {
	ID                  int64                 `orm:"column(id);auto" json:"-"`
	PurchaseOrder       *PurchaseOrder        `orm:"column(purchase_order_id);rel(fk)" json:"purchase_order,omitempty"`
	RecognitionDate     time.Time             `orm:"column(recognition_date);type(date)" json:"recognition_date"`
	Code                string                `orm:"column(code);size(45)" json:"code"`
	TotalAmount         float64               `orm:"column(total_amount);digits(20);decimals(0)" json:"total_amount"`
	Note                string                `orm:"column(note);null" json:"note"`
	DocumentStatus      string                `orm:"column(document_status);null;options(new,active,finished,cancelled)" json:"document_status"`
	IsDeleted           int8                  `orm:"column(is_deleted);null" json:"is_deleted"`
	CreatedBy           *User                 `orm:"column(created_by);rel(fk)" json:"created_by"`
	UpdatedBy           *User                 `orm:"column(updated_by);null;rel(fk)" json:"updated_by"`
	CreatedAt           time.Time             `orm:"column(created_at);type(timestamp)" json:"created_at"`
	UpdatedAt           time.Time             `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
	PurchaseReturnItems []*PurchaseReturnItem `orm:"reverse(many)" json:"purchase_return_items,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PurchaseReturn) MarshalJSON() ([]byte, error) {
	type Alias PurchaseReturn

	alias := &struct {
		ID              string `json:"id"`
		CreatedByID     string `json:"created_by_id"`
		UpdatedByID     string `json:"updated_by_id"`
		PurchaseOrderID string `json:"purchase_order_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.UpdatedByID when m.UpdatedBy not nill
	// and the ID is setted
	if m.UpdatedBy != nil && m.UpdatedBy.ID != int64(0) {
		alias.UpdatedByID = common.Encrypt(m.UpdatedBy.ID)
	} else {
		alias.UpdatedBy = nil
	}

	// Encrypt alias.PurchaseOrderID when m.PurchaseOrder not nill
	// and the ID is setted
	if m.PurchaseOrder != nil && m.PurchaseOrder.ID != int64(0) {
		alias.PurchaseOrderID = common.Encrypt(m.PurchaseOrder.ID)
	} else {
		alias.PurchaseOrder = nil
	}

	// Encrypt alias.CreatedByID when m.CreatedBy not nill
	// and the ID is setted
	if m.CreatedBy != nil && m.CreatedBy.ID != int64(0) {
		alias.CreatedByID = common.Encrypt(m.CreatedBy.ID)
	} else {
		alias.CreatedBy = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating PurchaseReturn struct into purchase_return table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to purchase_return.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *PurchaseReturn) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting purchase_return data
// this also will truncated all data from all table
// that have relation with this purchase_return.
func (m *PurchaseReturn) Delete() (err error) {
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
func (m *PurchaseReturn) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
