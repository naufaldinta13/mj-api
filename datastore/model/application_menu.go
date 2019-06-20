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
	orm.RegisterModel(new(ApplicationMenu))
}

// ApplicationMenu model for application_menu table.
type ApplicationMenu struct {
	ID                int64              `orm:"column(id);auto" json:"-"`
	ParentMenu        *ApplicationMenu   `orm:"column(parent_menu_id);null;rel(fk)" json:"parent_menu,omitempty"`
	ApplicationModule *ApplicationModule `orm:"column(application_module_id);rel(fk)" json:"application_module,omitempty"`
	MenuName          string             `orm:"column(menu_name);size(255)" json:"menu_name"`
	Icon              string             `orm:"column(icon);size(255);null" json:"icon"`
	Route             string             `orm:"column(route);size(255)" json:"route"`
	Order             uint               `orm:"column(order)" json:"order"`
	IsActive          int8               `orm:"column(is_active);null" json:"is_active"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *ApplicationMenu) MarshalJSON() ([]byte, error) {
	type Alias ApplicationMenu

	alias := &struct {
		ID                  string `json:"id"`
		ParentMenuID        string `json:"parent_menu_id"`
		ApplicationModuleID string `json:"application_module_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.ParentMenuID when m.ParentMenu not nill
	// and the ID is setted
	if m.ParentMenu != nil && m.ParentMenu.ID != int64(0) {
		alias.ParentMenuID = common.Encrypt(m.ParentMenu.ID)
	} else {
		alias.ParentMenu = nil
	}

	// Encrypt alias.ApplicationModuleID when m.ApplicationModule not nill
	// and the ID is setted
	if m.ApplicationModule != nil && m.ApplicationModule.ID != int64(0) {
		alias.ApplicationModuleID = common.Encrypt(m.ApplicationModule.ID)
	} else {
		alias.ApplicationModule = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating ApplicationMenu struct into application_menu table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to application_menu.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *ApplicationMenu) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting application_menu data
// this also will truncated all data from all table
// that have relation with this application_menu.
func (m *ApplicationMenu) Delete() (err error) {
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
func (m *ApplicationMenu) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
