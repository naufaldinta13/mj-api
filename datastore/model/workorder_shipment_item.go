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
	orm.RegisterModel(new(WorkorderShipmentItem))
}

// WorkorderShipmentItem model for workorder_shipment_item table.
type WorkorderShipmentItem struct {
	ID                   int64                 `orm:"column(id);auto" json:"-"`
	WorkorderShipment    *WorkorderShipment    `orm:"column(workorder_shipment_id);rel(fk)" json:"workorder_shipment,omitempty"`
	WorkorderFulfillment *WorkorderFulfillment `orm:"column(workorder_fulfillment_id);rel(fk)" json:"workorder_fulfillment,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *WorkorderShipmentItem) MarshalJSON() ([]byte, error) {
	type Alias WorkorderShipmentItem

	alias := &struct {
		ID                     string `json:"id"`
		WorkorderShipmentID    string `json:"workorder_shipment_id"`
		WorkorderFulfillmentID string `json:"workorder_fulfillment_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.WorkorderShipmentID when m.WorkorderShipment not nill
	// and the ID is setted
	if m.WorkorderShipment != nil && m.WorkorderShipment.ID != int64(0) {
		alias.WorkorderShipmentID = common.Encrypt(m.WorkorderShipment.ID)
	} else {
		alias.WorkorderShipment = nil
	}

	// Encrypt alias.WorkorderFulfillmentID when m.WorkorderFulfillment not nill
	// and the ID is setted
	if m.WorkorderFulfillment != nil && m.WorkorderFulfillment.ID != int64(0) {
		alias.WorkorderFulfillmentID = common.Encrypt(m.WorkorderFulfillment.ID)
	} else {
		alias.WorkorderFulfillment = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating WorkorderShipmentItem struct into workorder_shipment_item table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to workorder_shipment_item.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *WorkorderShipmentItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting workorder_shipment_item data
// this also will truncated all data from all table
// that have relation with this workorder_shipment_item.
func (m *WorkorderShipmentItem) Delete() (err error) {
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
func (m *WorkorderShipmentItem) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
