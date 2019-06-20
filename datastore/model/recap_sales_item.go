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
	orm.RegisterModel(new(RecapSalesItem))
}

// RecapSalesItem model for recap_sales_item table.
type RecapSalesItem struct {
	ID         int64       `orm:"column(id);auto" json:"-"`
	RecapSales *RecapSales `orm:"column(recap_sales_id);rel(fk)" json:"recap_sales,omitempty"`
	SalesOrder *SalesOrder `orm:"column(sales_order_id);rel(fk)" json:"sales_order,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *RecapSalesItem) MarshalJSON() ([]byte, error) {
	type Alias RecapSalesItem

	alias := &struct {
		ID           string `json:"id"`
		RecapSalesID string `json:"recap_sales_id"`
		SalesOrderID string `json:"sales_order_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.RecapSalesID when m.RecapSales not nill
	// and the ID is setted
	if m.RecapSales != nil && m.RecapSales.ID != int64(0) {
		alias.RecapSalesID = common.Encrypt(m.RecapSales.ID)
	} else {
		alias.RecapSales = nil
	}

	// Encrypt alias.SalesOrderID when m.SalesOrder not nill
	// and the ID is setted
	if m.SalesOrder != nil && m.SalesOrder.ID != int64(0) {
		alias.SalesOrderID = common.Encrypt(m.SalesOrder.ID)
	} else {
		alias.SalesOrder = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating RecapSalesItem struct into recap_sales_item table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to recap_sales_item.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *RecapSalesItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting recap_sales_item data
// this also will truncated all data from all table
// that have relation with this recap_sales_item.
func (m *RecapSalesItem) Delete() (err error) {
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
func (m *RecapSalesItem) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
