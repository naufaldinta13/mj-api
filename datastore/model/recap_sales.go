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
	orm.RegisterModel(new(RecapSales))
}

// RecapSales model for recap_sales table.
type RecapSales struct {
	ID              int64             `orm:"column(id);auto" json:"-"`
	Partnership     *Partnership      `orm:"column(partnership_id);rel(fk)" json:"partnership,omitempty"`
	Code            string            `orm:"column(code);size(45)" json:"code"`
	TotalAmount     float64           `orm:"column(total_amount);digits(20);decimals(0)" json:"total_amount"`
	CreatedBy       *User             `orm:"column(created_by);rel(fk)" json:"created_by"`
	CreatedAt       time.Time         `orm:"column(created_at);type(timestamp);auto_now" json:"created_at"`
	IsDeleted       int8              `orm:"column(is_deleted);null" json:"is_deleted"`
	RecapSalesItems []*RecapSalesItem `orm:"reverse(many)" json:"recap_sales_item,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *RecapSales) MarshalJSON() ([]byte, error) {
	type Alias RecapSales

	alias := &struct {
		ID            string `json:"id"`
		CreatedByID   string `json:"created_by_id"`
		PartnershipID string `json:"partnership_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.CreatedByID when m.CreatedBy not nill
	// and the ID is setted
	if m.CreatedBy != nil && m.CreatedBy.ID != int64(0) {
		alias.CreatedByID = common.Encrypt(m.CreatedBy.ID)
	} else {
		alias.CreatedBy = nil
	}

	// Encrypt alias.PartnershipID when m.Partnership not nill
	// and the ID is setted
	if m.Partnership != nil && m.Partnership.ID != int64(0) {
		alias.PartnershipID = common.Encrypt(m.Partnership.ID)
	} else {
		alias.Partnership = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating RecapSales struct into recap_sales table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to recap_sales.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *RecapSales) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting recap_sales data
// this also will truncated all data from all table
// that have relation with this recap_sales.
func (m *RecapSales) Delete() (err error) {
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
func (m *RecapSales) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
