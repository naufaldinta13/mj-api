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
	orm.RegisterModel(new(SalesReturn))
}

// SalesReturn model for sales_return table.
type SalesReturn struct {
	ID               int64              `orm:"column(id);auto" json:"-"`
	SalesOrder       *SalesOrder        `orm:"column(sales_order_id);rel(fk)" json:"sales_order,omitempty"`
	RecognitionDate  time.Time          `orm:"column(recognition_date);type(date)" json:"recognition_date"`
	Code             string             `orm:"column(code);size(45)" json:"code"`
	TotalAmount      float64            `orm:"column(total_amount);digits(20);decimals(0)" json:"total_amount"`
	Note             string             `orm:"column(note);null" json:"note"`
	DocumentStatus   string             `orm:"column(document_status);null;options(new,active,finished,cancelled)" json:"document_status"`
	IsBundled        int8               `orm:"column(is_bundled);null" json:"is_bundled"`
	IsDeleted        int8               `orm:"column(is_deleted);null" json:"is_deleted"`
	CreatedBy        *User              `orm:"column(created_by);rel(fk)" json:"created_by"`
	UpdatedBy        *User              `orm:"column(updated_by);null;rel(fk)" json:"updated_by"`
	CreatedAt        time.Time          `orm:"column(created_at);type(timestamp)" json:"created_at"`
	UpdatedAt        time.Time          `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
	SalesReturnItems []*SalesReturnItem `orm:"reverse(many)" json:"sales_return_items,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *SalesReturn) MarshalJSON() ([]byte, error) {
	type Alias SalesReturn

	alias := &struct {
		ID           string `json:"id"`
		SalesOrderID string `json:"sales_order_id"`
		CreatedByID  string `json:"created_by_id"`
		UpdatedByID  string `json:"updated_by_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.SalesOrderID when m.SalesOrder not nill
	// and the ID is setted
	if m.SalesOrder != nil && m.SalesOrder.ID != int64(0) {
		alias.SalesOrderID = common.Encrypt(m.SalesOrder.ID)
	} else {
		alias.SalesOrder = nil
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

// Save inserting or updating SalesReturn struct into sales_return table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to sales_return.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *SalesReturn) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting sales_return data
// this also will truncated all data from all table
// that have relation with this sales_return.
func (m *SalesReturn) Delete() (err error) {
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
func (m *SalesReturn) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
