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
	orm.RegisterModel(new(StockopnameItem))
}

// StockopnameItem model for stockopname_item table.
type StockopnameItem struct {
	ID               int64             `orm:"column(id);auto" json:"-"`
	Stockopname      *Stockopname      `orm:"column(stockopname_id);rel(fk)" json:"stockopname,omitempty"`
	ItemVariantStock *ItemVariantStock `orm:"column(item_variant_stock_id);rel(fk)" json:"item_variant_stock,omitempty"`
	Quantity         float32           `orm:"column(quantity)" json:"quantity"`
	Note             string            `orm:"column(note);null" json:"note"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *StockopnameItem) MarshalJSON() ([]byte, error) {
	type Alias StockopnameItem

	alias := &struct {
		ID                 string `json:"id"`
		StockopnameID      string `json:"stockopname_id"`
		ItemVariantStockID string `json:"item_variant_stock_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.StockopnameID when m.Stockopname not nill
	// and the ID is setted
	if m.Stockopname != nil && m.Stockopname.ID != int64(0) {
		alias.StockopnameID = common.Encrypt(m.Stockopname.ID)
	} else {
		alias.Stockopname = nil
	}

	// Encrypt alias.ItemVariantStockID when m.ItemVariantStock not nill
	// and the ID is setted
	if m.ItemVariantStock != nil && m.ItemVariantStock.ID != int64(0) {
		alias.ItemVariantStockID = common.Encrypt(m.ItemVariantStock.ID)
	} else {
		alias.ItemVariantStock = nil
	}

	return json.Marshal(alias)
}

// Save inserting or updating StockopnameItem struct into stockopname_item table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to stockopname_item.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *StockopnameItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting stockopname_item data
// this also will truncated all data from all table
// that have relation with this stockopname_item.
func (m *StockopnameItem) Delete() (err error) {
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
func (m *StockopnameItem) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
