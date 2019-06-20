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
	orm.RegisterModel(new(ApplicationPrivilege))
}

// ApplicationPrivilege model for application_privilege table.
type ApplicationPrivilege struct {
	ID                int64              `orm:"column(id);auto" json:"-"`
	ApplicationModule *ApplicationModule `orm:"column(application_module_id);rel(fk)" json:"application_module,omitempty"`
	Usergroup         *Usergroup         `orm:"column(usergroup_id);rel(fk)" json:"usergroup,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *ApplicationPrivilege) MarshalJSON() ([]byte, error) {
	type Alias ApplicationPrivilege

	alias := &struct {
		ID                  string `json:"id"`
		ApplicationModuleID string `json:"application_module_id"`
		UsergroupID         string `json:"usergroup_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.ApplicationModuleID when m.ApplicationModule not nill
	// and the ID is setted
	if m.ApplicationModule != nil && m.ApplicationModule.ID != int64(0) {
		alias.ApplicationModuleID = common.Encrypt(m.ApplicationModule.ID)
	} else {
		alias.ApplicationModule = nil
	}

	// Encrypt alias.UsergroupID when m.Usergroup not nill
	// and the ID is setted
	if m.Usergroup != nil && m.Usergroup.ID != int64(0) {
		alias.UsergroupID = common.Encrypt(m.Usergroup.ID)
	} else {
		alias.Usergroup = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating ApplicationPrivilege struct into application_privilege table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to application_privilege.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *ApplicationPrivilege) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting application_privilege data
// this also will truncated all data from all table
// that have relation with this application_privilege.
func (m *ApplicationPrivilege) Delete() (err error) {
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
func (m *ApplicationPrivilege) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
