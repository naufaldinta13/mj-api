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
	orm.RegisterModel(new(BankAccount))
}

// BankAccount model for bank_account table.
type BankAccount struct {
	ID         int64  `orm:"column(id);auto" json:"-"`
	BankName   string `orm:"column(bank_name);size(45);null" json:"bank_name"`
	BankNumber string `orm:"column(bank_number);size(45);null" json:"bank_number"`
	IsDefault  int8   `orm:"column(is_default);null" json:"is_default"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *BankAccount) MarshalJSON() ([]byte, error) {
	type Alias BankAccount

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save inserting or updating BankAccount struct into bank_account table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to bank_accounts.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *BankAccount) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting bank_accounts data
// this also will truncated all data from all table
// that have relation with this bank_accounts.
func (m *BankAccount) Delete() (err error) {
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
func (m *BankAccount) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
