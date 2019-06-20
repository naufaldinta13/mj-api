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
	orm.RegisterModel(new(User))
}

// User model for user table.
type User struct {
	ID            int64      `orm:"column(id);auto" json:"-"`
	Usergroup     *Usergroup `orm:"column(usergroup_id);rel(fk)" json:"usergroup,omitempty"`
	FullName      string     `orm:"column(full_name);size(100)" json:"full_name"`
	Username      string     `orm:"column(username);size(100)" json:"username"`
	Password      string     `orm:"column(password);size(145)" json:"password"`
	RememberToken string     `orm:"column(remember_token);size(128);null" json:"remember_token"`
	IsActive      int8       `orm:"column(is_active);null" json:"is_active"`
	LastLogin     time.Time  `orm:"column(last_login);type(timestamp);null" json:"last_login"`
	CreatedBy     *User      `orm:"column(created_by);null;rel(fk)" json:"created_by"`
	UpdatedBy     *User      `orm:"column(updated_by);null;rel(fk)" json:"updated_by"`
	CreatedAt     time.Time  `orm:"column(created_at);type(timestamp)" json:"created_at"`
	UpdatedAt     time.Time  `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *User) MarshalJSON() ([]byte, error) {
	type Alias User

	alias := &struct {
		ID          string `json:"id"`
		UpdatedByID string `json:"updated_by_id"`
		UsergroupID string `json:"usergroup_id"`
		CreatedByID string `json:"created_by_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.UsergroupID when m.Usergroup not nill
	// and the ID is setted
	if m.Usergroup != nil && m.Usergroup.ID != int64(0) {
		alias.UsergroupID = common.Encrypt(m.Usergroup.ID)
	} else {
		alias.Usergroup = nil
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

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *User) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting user data
// this also will truncated all data from all table
// that have relation with this user.
func (m *User) Delete() (err error) {
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
func (m *User) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}

//GetUser is used by jwtUser to get model user
func (m *User) GetUser(id int64) (interface{}, error) {
	m.ID = id
	e := m.Read()
	return m, e
}

// Inactive status User
func (m *User) Inactive() (err error) {
	m.IsActive = 0
	m.UpdatedAt = time.Now()

	o := orm.NewOrm()
	if _, err = o.Update(m, "is_active", "updated_at"); err == nil {
		return err
	}

	return nil
}
