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
	orm.RegisterModel(new(DirectPlacementItem))
}

// DirectPlacementItem model for direct_placement_item table.
type DirectPlacementItem struct {
	ID             int64            `orm:"column(id);auto" json:"-"`
	DirectPlacment *DirectPlacement `orm:"column(direct_placment_id);rel(fk)" json:"direct_placment,omitempty"`
	ItemVariant    *ItemVariant     `orm:"column(item_variant_id);rel(fk)" json:"item_variant,omitempty"`
	Quantity       float32          `orm:"column(quantity);null" json:"quantity"`
	UnitPrice      float64          `orm:"column(unit_price);digits(20);decimals(0)" json:"unit_price"`
	TotalPrice     float64          `orm:"column(total_price);digits(20);decimals(0)" json:"total_price"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *DirectPlacementItem) MarshalJSON() ([]byte, error) {
	type Alias DirectPlacementItem

	alias := &struct {
		ID               string `json:"id"`
		DirectPlacmentID string `json:"direct_placment_id"`
		ItemVariantID    string `json:"item_variant_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.DirectPlacmentID when m.DirectPlacment not nill
	// and the ID is setted
	if m.DirectPlacment != nil && m.DirectPlacment.ID != int64(0) {
		alias.DirectPlacmentID = common.Encrypt(m.DirectPlacment.ID)
	} else {
		alias.DirectPlacment = nil
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

// Save inserting or updating DirectPlacementItem struct into direct_placement_item table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to direct_placement_item.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *DirectPlacementItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting direct_placement_item data
// this also will truncated all data from all table
// that have relation with this direct_placement_item.
func (m *DirectPlacementItem) Delete() (err error) {
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
func (m *DirectPlacementItem) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
