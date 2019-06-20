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
	orm.RegisterModel(new(SalesReturnItem))
}

// SalesReturnItem model for sales_return_item table.
type SalesReturnItem struct {
	ID             int64           `orm:"column(id);auto" json:"-"`
	SalesReturn    *SalesReturn    `orm:"column(sales_return_id);rel(fk)" json:"sales_return,omitempty"`
	SalesOrderItem *SalesOrderItem `orm:"column(sales_order_item_id);rel(fk)" json:"sales_order_item,omitempty"`
	Quantity       float32         `orm:"column(quantity)" json:"quantity"`
	Note           string          `orm:"column(note);null" json:"note"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *SalesReturnItem) MarshalJSON() ([]byte, error) {
	type Alias SalesReturnItem

	alias := &struct {
		ID               string `json:"id"`
		SalesReturnID    string `json:"sales_return_id"`
		SalesOrderItemID string `json:"sales_order_item_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.SalesReturnID when m.SalesReturn not nill
	// and the ID is setted
	if m.SalesReturn != nil && m.SalesReturn.ID != int64(0) {
		alias.SalesReturnID = common.Encrypt(m.SalesReturn.ID)
	} else {
		alias.SalesReturn = nil
	}

	// Encrypt alias.SalesOrderItemID when m.SalesOrderItem not nill
	// and the ID is setted
	if m.SalesOrderItem != nil && m.SalesOrderItem.ID != int64(0) {
		alias.SalesOrderItemID = common.Encrypt(m.SalesOrderItem.ID)
	} else {
		alias.SalesOrderItem = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating SalesReturnItem struct into sales_return_item table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to sales_return_item.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *SalesReturnItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting sales_return_item data
// this also will truncated all data from all table
// that have relation with this sales_return_item.
func (m *SalesReturnItem) Delete() (err error) {
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
func (m *SalesReturnItem) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
