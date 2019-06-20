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
	orm.RegisterModel(new(Partnership))
}

// Partnership model for partnership table.
type Partnership struct {
	ID               int64     `orm:"column(id);auto" json:"-"`
	Code             string    `orm:"column(code);size(120)" json:"code"`
	PartnershipType  string    `orm:"column(partnership_type);null;options(customer,supplier)" json:"partnership_type"`
	OrderRule        string    `orm:"column(order_rule);null;options(none,one_bill,plafon)" json:"order_rule"`
	FullName         string    `orm:"column(full_name);size(45)" json:"full_name"`
	Email            string    `orm:"column(email);size(45);null" json:"email"`
	Phone            string    `orm:"column(phone);size(45);null" json:"phone"`
	Address          string    `orm:"column(address);null" json:"address"`
	City             string    `orm:"column(city);size(45);null" json:"city"`
	Province         string    `orm:"column(province);size(45);null" json:"province"`
	BankName         string    `orm:"column(bank_name);size(45);null" json:"bank_name"`
	BankNumber       string    `orm:"column(bank_number);size(45);null" json:"bank_number"`
	BankHolder       string    `orm:"column(bank_holder);size(45);null" json:"bank_holder"`
	MaxPlafon        float64   `orm:"column(max_plafon);null;digits(20);decimals(0)" json:"max_plafon"`
	TotalDebt        float64   `orm:"column(total_debt);null;digits(20);decimals(0)" json:"total_debt"`
	TotalCredit      float64   `orm:"column(total_credit);null;digits(20);decimals(0)" json:"total_credit"`
	TotalSpend       float64   `orm:"column(total_spend);null;digits(20);decimals(0)" json:"total_spend"`
	TotalExpenditure float64   `orm:"column(total_expenditure);null;digits(20);decimals(0)" json:"total_expenditure"`
	SalesPerson      string    `orm:"column(sales_person);size(45);null" json:"sales_person"`
	VisitDay         string    `orm:"column(visit_day);size(45);null" json:"visit_day"`
	Note             string    `orm:"column(note);null" json:"note"`
	IsArchived       int8      `orm:"column(is_archived);null" json:"is_archived"`
	IsDeleted        int8      `orm:"column(is_deleted);null" json:"is_deleted"`
	IsDefault        int8      `orm:"column(is_default);null" json:"is_default"`
	CreatedBy        *User     `orm:"column(created_by);rel(fk)" json:"created_by"`
	UpdatedBy        *User     `orm:"column(updated_by);null;rel(fk)" json:"updated_by"`
	CreatedAt        time.Time `orm:"column(created_at);type(timestamp)" json:"created_at"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *Partnership) MarshalJSON() ([]byte, error) {
	type Alias Partnership

	alias := &struct {
		ID          string `json:"id"`
		CreatedByID string `json:"created_by_id"`
		UpdatedByID string `json:"updated_by_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.UpdatedByID when m.UpdatedBy not nill
	// and the ID is setted
	if m.UpdatedBy != nil && m.UpdatedBy.ID != int64(0) {
		alias.UpdatedByID = common.Encrypt(m.UpdatedBy.ID)
	} else {
		alias.UpdatedBy = nil
	}

	// Encrypt alias.CreatedByID when m.CreatedBy not nill
	// and the ID is setted
	if m.CreatedBy != nil && m.CreatedBy.ID != int64(0) {
		alias.CreatedByID = common.Encrypt(m.CreatedBy.ID)
	} else {
		alias.CreatedBy = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating Partnership struct into partnership table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to partnership.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *Partnership) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting partnership data
// this also will truncated all data from all table
// that have relation with this partnership.
func (m *Partnership) Delete() (err error) {
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
func (m *Partnership) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
