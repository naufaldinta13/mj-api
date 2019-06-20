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
	orm.RegisterModel(new(ApplicationModule))
}

// ApplicationModule model for application_module table.
type ApplicationModule struct {
	ID           int64              `orm:"column(id);auto" json:"-"`
	ParentModule *ApplicationModule `orm:"column(parent_module_id);null;rel(fk)" json:"parent_module,omitempty"`
	ModuleName   string             `orm:"column(module_name);size(255);null" json:"module_name"`
	Alias        string             `orm:"column(alias);size(255)" json:"alias"`
	Note         string             `orm:"column(note);null" json:"note"`
	IsActive     int8               `orm:"column(is_active);null" json:"is_active"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *ApplicationModule) MarshalJSON() ([]byte, error) {
	type Alias ApplicationModule

	alias := &struct {
		ID             string `json:"id"`
		ParentModuleID string `json:"parent_module_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.ParentModuleID when m.ParentModule not nill
	// and the ID is setted
	if m.ParentModule != nil && m.ParentModule.ID != int64(0) {
		alias.ParentModuleID = common.Encrypt(m.ParentModule.ID)
	} else {
		alias.ParentModule = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating ApplicationModule struct into application_module table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to application_module.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *ApplicationModule) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting application_module data
// this also will truncated all data from all table
// that have relation with this application_module.
func (m *ApplicationModule) Delete() (err error) {
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
func (m *ApplicationModule) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
