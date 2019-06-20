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
	orm.RegisterModel(new(ItemVariantStock))
}

// ItemVariantStock model for item_variant_stock table.
type ItemVariantStock struct {
	ID             int64        `orm:"column(id);auto" json:"-"`
	ItemVariant    *ItemVariant `orm:"column(item_variant_id);rel(fk)" json:"item_variant,omitempty"`
	SkuCode        string       `orm:"column(sku_code);size(120)" json:"sku_code"`
	AvailableStock float32      `orm:"column(available_stock);null" json:"available_stock"`
	UnitCost       float64      `orm:"column(unit_cost);null;digits(20);decimals(0)" json:"unit_cost"`
	CreatedBy      *User        `orm:"column(created_by);null;rel(fk)" json:"created_by"`
	UpdatedBy      *User        `orm:"column(updated_by);null;rel(fk)" json:"updated_by"`
	CreatedAt      time.Time    `orm:"column(created_at);type(timestamp)" json:"created_at"`
	UpdatedAt      time.Time    `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *ItemVariantStock) MarshalJSON() ([]byte, error) {
	type Alias ItemVariantStock

	alias := &struct {
		ID            string `json:"id"`
		ItemVariantID string `json:"item_variant_id"`
		CreatedByID   string `json:"created_by_id"`
		UpdatedByID   string `json:"updated_by_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.ItemVariantID when m.ItemVariant not nill
	// and the ID is setted
	if m.ItemVariant != nil && m.ItemVariant.ID != int64(0) {
		alias.ItemVariantID = common.Encrypt(m.ItemVariant.ID)
	} else {
		alias.ItemVariant = nil
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

// Save inserting or updating ItemVariantStock struct into item_variant_stock table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to item_variant_stock.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *ItemVariantStock) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting item_variant_stock data
// this also will truncated all data from all table
// that have relation with this item_variant_stock.
func (m *ItemVariantStock) Delete() (err error) {
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
func (m *ItemVariantStock) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
