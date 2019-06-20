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
	orm.RegisterModel(new(WorkorderReceiving))
}

// WorkorderReceiving model for workorder_receiving table.
type WorkorderReceiving struct {
	ID                      int64                     `orm:"column(id);auto" json:"-"`
	PurchaseOrder           *PurchaseOrder            `orm:"column(purchase_order_id);null;rel(fk)" json:"purchase_order,omitempty"`
	RecognitionDate         time.Time                 `orm:"column(recognition_date);type(date)" json:"recognition_date"`
	Code                    string                    `orm:"column(code);size(45)" json:"code"`
	Pic                     string                    `orm:"column(pic);size(45)" json:"pic"`
	Note                    string                    `orm:"column(note);null" json:"note"`
	DocumentStatus          string                    `orm:"column(document_status);null;options(active,finished)" json:"document_status"`
	IsDeleted               int8                      `orm:"column(is_deleted);null" json:"is_deleted"`
	CreatedBy               *User                     `orm:"column(created_by);rel(fk)" json:"created_by"`
	UpdatedBy               *User                     `orm:"column(updated_by);null;rel(fk)" json:"updated_by"`
	CreatedAt               time.Time                 `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	UpdatedAt               time.Time                 `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
	WorkorderReceivingItems []*WorkorderReceivingItem `orm:"reverse(many)" json:"work_order_receiving_items,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *WorkorderReceiving) MarshalJSON() ([]byte, error) {
	type Alias WorkorderReceiving

	alias := &struct {
		ID              string `json:"id"`
		PurchaseOrderID string `json:"purchase_order_id"`
		CreatedByID     string `json:"created_by_id"`
		UpdatedByID     string `json:"updated_by_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
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

	// Encrypt alias.UpdatedByID when m.UpdatedBy not nill
	// and the ID is setted
	if m.UpdatedBy != nil && m.UpdatedBy.ID != int64(0) {
		alias.UpdatedByID = common.Encrypt(m.UpdatedBy.ID)
	} else {
		alias.UpdatedBy = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating WorkorderReceiving struct into workorder_receiving table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to workorder_receiving.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *WorkorderReceiving) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting workorder_receiving data
// this also will truncated all data from all table
// that have relation with this workorder_receiving.
func (m *WorkorderReceiving) Delete() (err error) {
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
func (m *WorkorderReceiving) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
