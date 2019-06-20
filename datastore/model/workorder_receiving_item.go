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
	orm.RegisterModel(new(WorkorderReceivingItem))
}

// WorkorderReceivingItem model for workorder_receiving_item table.
type WorkorderReceivingItem struct {
	ID                 int64               `orm:"column(id);auto" json:"-"`
	WorkorderReceiving *WorkorderReceiving `orm:"column(workorder_receiving_id);rel(fk)" json:"workorder_receiving,omitempty"`
	PurchaseOrderItem  *PurchaseOrderItem  `orm:"column(purchase_order_item_id);rel(fk)" json:"purchase_order_item,omitempty"`
	Quantity           float32             `orm:"column(quantity)" json:"quantity"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *WorkorderReceivingItem) MarshalJSON() ([]byte, error) {
	type Alias WorkorderReceivingItem

	alias := &struct {
		ID                   string `json:"id"`
		PurchaseOrderItemID  string `json:"purchase_order_item_id"`
		WorkorderReceivingID string `json:"workorder_receiving_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.WorkorderReceivingID when m.WorkorderReceiving not nill
	// and the ID is setted
	if m.WorkorderReceiving != nil && m.WorkorderReceiving.ID != int64(0) {
		alias.WorkorderReceivingID = common.Encrypt(m.WorkorderReceiving.ID)
	} else {
		alias.WorkorderReceiving = nil
	}

	// Encrypt alias.PurchaseOrderItemID when m.PurchaseOrderItem not nill
	// and the ID is setted
	if m.PurchaseOrderItem != nil && m.PurchaseOrderItem.ID != int64(0) {
		alias.PurchaseOrderItemID = common.Encrypt(m.PurchaseOrderItem.ID)
	} else {
		alias.PurchaseOrderItem = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating WorkorderReceivingItem struct into workorder_receiving_item table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to workorder_receiving_item.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *WorkorderReceivingItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting workorder_receiving_item data
// this also will truncated all data from all table
// that have relation with this workorder_receiving_item.
func (m *WorkorderReceivingItem) Delete() (err error) {
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
func (m *WorkorderReceivingItem) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
