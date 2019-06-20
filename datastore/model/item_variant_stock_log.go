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
	orm.RegisterModel(new(ItemVariantStockLog))
}

// ItemVariantStockLog model for item_variant_stock_log table.
type ItemVariantStockLog struct {
	ID                   int64                 `orm:"column(id);auto" json:"-"`
	ItemVariantStock     *ItemVariantStock     `orm:"column(item_variant_stock_id);rel(fk)" json:"item_variant_stock,omitempty"`
	RefID                uint64                `orm:"column(ref_id)" json:"ref_id"`
	RefType              string                `orm:"column(ref_type);null;options(workorder_fulfillment,workorder_receiving,stockopname,direct_placement)" json:"ref_type"`
	LogType              string                `orm:"column(log_type);null;options(in,out)" json:"log_type"`
	Quantity             float32               `orm:"column(quantity)" json:"quantity"`
	FinalStock           float32               `orm:"column(final_stock)" json:"final_stock"`
	Stockopname          *Stockopname          `orm:"-" json:"stockopname,omitempty"`
	WorkorderFulfillment *WorkorderFulfillment `orm:"-" json:"workorder_fulfillment,omitempty"`
	WorkorderReceiving   *WorkorderReceiving   `orm:"-" json:"workorder_receiving,omitempty"`
	DirectPlacement      *DirectPlacement      `orm:"-" json:"direct_placement,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *ItemVariantStockLog) MarshalJSON() ([]byte, error) {
	type Alias ItemVariantStockLog

	alias := &struct {
		ID                 string `json:"id"`
		ItemVariantStockID string `json:"item_variant_stock_id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.ItemVariantStockID when m.ItemVariantStock not nill
	// and the ID is setted
	if m.ItemVariantStock != nil && m.ItemVariantStock.ID != int64(0) {
		alias.ItemVariantStockID = common.Encrypt(m.ItemVariantStock.ID)
	} else {
		alias.ItemVariantStock = nil
	}

	o := orm.NewOrm()
	if m.RefType == "stockopname" {
		o.Raw("select * from stockopname where id = ?", int64(m.RefID)).QueryRow(&m.Stockopname)
	}

	if m.RefType == "workorder_fulfillment" {
		o.Raw("select * from workorder_fulfillment where id = ?", int64(m.RefID)).QueryRow(&m.WorkorderFulfillment)
	}

	if m.RefType == "workorder_receiving" {
		o.Raw("select * from workorder_receiving where id = ?", int64(m.RefID)).QueryRow(&m.WorkorderReceiving)
	}

	if m.RefType == "direct_placement" {
		o.Raw("select * from direct_placement where id = ?", int64(m.RefID)).QueryRow(&m.DirectPlacement)
	}

	m.ItemVariantStock.Read()
	m.ItemVariantStock.ItemVariant.Read()
	m.ItemVariantStock.ItemVariant.Measurement.Read()
	m.ItemVariantStock.ItemVariant.Item.Read()

	return json.Marshal(alias)
}

// Save inserting or updating ItemVariantStockLog struct into item_variant_stock_log table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to item_variant_stock_log.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *ItemVariantStockLog) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting item_variant_stock_log data
// this also will truncated all data from all table
// that have relation with this item_variant_stock_log.
func (m *ItemVariantStockLog) Delete() (err error) {
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
func (m *ItemVariantStockLog) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
