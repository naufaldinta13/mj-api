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
	orm.RegisterModel(new(SalesOrderItem))
}

// SalesOrderItem model for sales_order_item table.
type SalesOrderItem struct {
	ID                  int64                `orm:"column(id);auto" json:"-"`
	QuantityFulfillment float32              `orm:"column(quantity_fulfillment)" json:"quantity_fulfillment"`
	SalesOrder          *SalesOrder          `orm:"column(sales_order_id);rel(fk)" json:"sales_order,omitempty"`
	ItemVariant         *ItemVariant         `orm:"column(item_variant_id);rel(fk)" json:"item_variant,omitempty"`
	Quantity            float32              `orm:"column(quantity)" json:"quantity"`
	QuantityPrepare     float32              `orm:"-" json:"quantity_prepare"`
	UnitPrice           float64              `orm:"column(unit_price);null;digits(20);decimals(0)" json:"unit_price"`
	Discount            float32              `orm:"column(discount);null" json:"discount"`
	Subtotal            float64              `orm:"column(subtotal);digits(20);decimals(0)" json:"subtotal"`
	Note                string               `orm:"column(note);null" json:"note"`
	CanBeReturn         float32              `orm:"-" json:"can_be_return, omitempty"`
	ItemVariantStockLog *ItemVariantStockLog `orm:"-" json:"item_variant_stock_log,omitempty"`
	Prefix              float64              `orm:"-" json:"prefix,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *SalesOrderItem) MarshalJSON() ([]byte, error) {
	type Alias SalesOrderItem

	alias := &struct {
		ID            string `json:"id"`
		SalesOrderID  string `json:"sales_order_id"`
		ItemVariantID string `json:"item_variant_id"`
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

	// Encrypt alias.ItemVariantID when m.ItemVariant not nill
	// and the ID is setted
	if m.ItemVariant != nil && m.ItemVariant.ID != int64(0) {
		alias.ItemVariantID = common.Encrypt(m.ItemVariant.ID)
	} else {
		alias.ItemVariant = nil
	}

	if m.ItemVariantStockLog != nil {
		m.Prefix = m.Subtotal - m.ItemVariantStockLog.ItemVariantStock.UnitCost
	}

	return json.Marshal(alias)
}

// Save inserting or updating SalesOrderItem struct into sales_order_item table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to sales_order_item.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *SalesOrderItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting sales_order_item data
// this also will truncated all data from all table
// that have relation with this sales_order_item.
func (m *SalesOrderItem) Delete() (err error) {
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
func (m *SalesOrderItem) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
