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
	orm.RegisterModel(new(Usergroup))
}

// Usergroup model for usergroup table.
type Usergroup struct {
	ID            int64  `orm:"column(id);auto" json:"-"`
	UsergroupName string `orm:"column(usergroup_name);size(45)" json:"usergroup_name"`
	Note          string `orm:"column(note);null" json:"note"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Usergroup) MarshalJSON() ([]byte, error) {
	type Alias Usergroup

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save inserting or updating Usergroup struct into usergroup table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to usergroup.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Usergroup) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting usergroup data
// this also will truncated all data from all table
// that have relation with this usergroup.
func (m *Usergroup) Delete() (err error) {
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
func (m *Usergroup) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
