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
	orm.RegisterModel(new(ItemCategory))
}

// ItemCategory model for item_category table.
type ItemCategory struct {
	ID           int64  `orm:"column(id);auto" json:"-"`
	CategoryName string `orm:"column(category_name);size(45)" json:"category_name"`
	Note         string `orm:"column(note);null" json:"note"`
	IsDeleted    int8   `orm:"column(is_deleted);null" json:"is_deleted"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *ItemCategory) MarshalJSON() ([]byte, error) {
	type Alias ItemCategory

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save inserting or updating ItemCategory struct into item_category table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to item_category.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *ItemCategory) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting item_category data
// this also will truncated all data from all table
// that have relation with this item_category.
func (m *ItemCategory) Delete() (err error) {
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
func (m *ItemCategory) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}

// ChangeStatusDelete status Category
func (m *ItemCategory) ChangeStatusDelete() (err error) {
	m.IsDeleted = 1
	o := orm.NewOrm()
	if _, err = o.Update(m, "IsDeleted"); err == nil {
		return err
	}
	return nil
}
