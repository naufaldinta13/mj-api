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
	orm.RegisterModel(new(ItemVariantPrice))
}

// ItemVariantPrice model for item_variant_price table.
type ItemVariantPrice struct {
	ID          int64        `orm:"column(id);auto" json:"-"`
	ItemVariant *ItemVariant `orm:"column(item_variant_id);rel(fk)" json:"item_variant,omitempty"`
	PricingType *PricingType `orm:"column(pricing_type_id);rel(fk)" json:"pricing_type,omitempty"`
	UnitPrice   float64      `orm:"column(unit_price);digits(20);decimals(0)" json:"unit_price"`
	Note        string       `orm:"column(note);null" json:"note"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *ItemVariantPrice) MarshalJSON() ([]byte, error) {
	type Alias ItemVariantPrice

	alias := &struct {
		ID            string `json:"id"`
		ItemVariantID string `json:"item_variant_id"`
		PricingTypeID string `json:"pricing_type_id"`
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

	// Encrypt alias.PricingTypeID when m.PricingType not nill
	// and the ID is setted
	if m.PricingType != nil && m.PricingType.ID != int64(0) {
		alias.PricingTypeID = common.Encrypt(m.PricingType.ID)
	} else {
		alias.PricingType = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating ItemVariantPrice struct into item_variant_price table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to item_variant_price.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *ItemVariantPrice) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting item_variant_price data
// this also will truncated all data from all table
// that have relation with this item_variant_price.
func (m *ItemVariantPrice) Delete() (err error) {
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
func (m *ItemVariantPrice) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
