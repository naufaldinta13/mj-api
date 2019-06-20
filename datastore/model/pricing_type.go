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
	orm.RegisterModel(new(PricingType))
}

// PricingType model for pricing_type table.
type PricingType struct {
	ID           int64        `orm:"column(id);auto" json:"-"`
	ParentType   *PricingType `orm:"column(parent_type_id);null;rel(fk)" json:"parent_type"`
	TypeName     string       `orm:"column(type_name);size(45)" json:"type_name"`
	Note         string       `orm:"column(note);null" json:"note"`
	RuleType     string       `orm:"column(rule_type);options(increment,decrement,none);null" json:"rule_type"`
	Nominal      float64      `orm:"column(nominal);digits(20);decimals(0)" json:"nominal"`
	IsPercentage int8         `orm:"column(is_percentage);null" json:"is_percentage"`
	IsDefault    int8         `orm:"column(is_default);null" json:"is_default"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *PricingType) MarshalJSON() ([]byte, error) {
	type Alias PricingType

	alias := &struct {
		ID           string `json:"id"`
		ParentTypeID string `json:"parent_type_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.ParentTypeID when m.ParentType not nill
	// and the ID is setted
	if m.ParentType != nil && m.ParentType.ID != int64(0) {
		alias.ParentTypeID = common.Encrypt(m.ParentType.ID)
	} else {
		alias.ParentType = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating PricingType struct into pricing_type table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to pricing_type.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *PricingType) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting pricing_type data
// this also will truncated all data from all table
// that have relation with this pricing_type.
func (m *PricingType) Delete() (err error) {
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
func (m *PricingType) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
