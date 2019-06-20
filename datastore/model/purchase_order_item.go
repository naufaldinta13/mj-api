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
	orm.RegisterModel(new(PurchaseOrderItem))
}

// PurchaseOrderItem model for purchase_order_item table.
type PurchaseOrderItem struct {
	ID            int64          `orm:"column(id);auto" json:"-"`
	PurchaseOrder *PurchaseOrder `orm:"column(purchase_order_id);rel(fk)" json:"purchase_order,omitempty"`
	ItemVariant   *ItemVariant   `orm:"column(item_variant_id);rel(fk)" json:"item_variant,omitempty"`
	Quantity      float32        `orm:"column(quantity)" json:"quantity"`
	UnitPrice     float64        `orm:"column(unit_price);null;digits(20);decimals(0)" json:"unit_price"`
	Discount      float32        `orm:"column(discount);null" json:"discount"`
	Subtotal      float64        `orm:"column(subtotal);digits(20);decimals(0)" json:"subtotal"`
	Note          string         `orm:"column(note);null" json:"note"`
	CanBeReturn   float32        `orm:"-" json:"can_be_return,omitempty"` // Ketika Receiving Item telah dibuat maka update CanBeReturn Sesuai dengan Quantity Item yang di Receive
	Partnership   *Partnership   `orm:"-" json:"partnership,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PurchaseOrderItem) MarshalJSON() ([]byte, error) {
	type Alias PurchaseOrderItem

	alias := &struct {
		ID              string `json:"id"`
		PurchaseOrderID string `json:"purchase_order_id"`
		ItemVariantID   string `json:"item_variant_id"`
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

	// Encrypt alias.ItemVariantID when m.ItemVariant not nill
	// and the ID is setted
	if m.ItemVariant != nil && m.ItemVariant.ID != int64(0) {
		alias.ItemVariantID = common.Encrypt(m.ItemVariant.ID)
	} else {
		alias.ItemVariant = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating PurchaseOrderItem struct into purchase_order_item table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to purchase_order_item.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *PurchaseOrderItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting purchase_order_item data
// this also will truncated all data from all table
// that have relation with this purchase_order_item.
func (m *PurchaseOrderItem) Delete() (err error) {
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
func (m *PurchaseOrderItem) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
