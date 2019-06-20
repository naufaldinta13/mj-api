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
	orm.RegisterModel(new(ItemVariant))
}

// ItemVariant model for item_variant table.
type ItemVariant struct {
	ID                   int64                  `orm:"column(id);auto" json:"-"`
	Item                 *Item                  `orm:"column(item_id);rel(fk)" json:"item,omitempty"`
	Measurement          *Measurement           `orm:"column(measurement_id);rel(fk)" json:"measurement,omitempty"`
	Barcode              string                 `orm:"column(barcode);size(120);null" json:"barcode"`
	ExternalName         string                 `orm:"column(external_name);size(200);null" json:"external_name"`
	VariantName          string                 `orm:"column(variant_name);size(45);null" json:"variant_name"`
	Image                string                 `orm:"column(image);null" json:"image"`
	BasePrice            float64                `orm:"column(base_price);null;digits(20);decimals(0)" json:"base_price"`
	Note                 string                 `orm:"column(note);null" json:"note"`
	MinimumStock         float32                `orm:"column(minimum_stock);null" json:"minimum_stock"`
	AvailableStock       float32                `orm:"column(available_stock);null" json:"available_stock"`
	CommitedStock        float32                `orm:"column(commited_stock);null" json:"commited_stock"`
	HasExternalName      int8                   `orm:"column(has_external_name);null" json:"has_external_name"`
	IsArchived           int8                   `orm:"column(is_archived);null" json:"is_archived"`
	IsDeleted            int8                   `orm:"column(is_deleted);null" json:"is_deleted"`
	CreatedBy            *User                  `orm:"column(created_by);rel(fk)" json:"created_by"`
	UpdatedBy            *User                  `orm:"column(updated_by);null;rel(fk)" json:"updated_by"`
	CreatedAt            time.Time              `orm:"column(created_at);type(timestamp)" json:"created_at"`
	UpdatedAt            time.Time              `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
	ItemVariantPrices    []*ItemVariantPrice    `orm:"reverse(many)" json:"item_variant_prices,omitempty"`
	ItemVariantStocks    []*ItemVariantStock    `orm:"reverse(many)" json:"item_variant_stocks,omitempty"`
	ItemVariantStockLogs []*ItemVariantStockLog `orm:"-" json:"item_variant_stock_logs,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *ItemVariant) MarshalJSON() ([]byte, error) {
	type Alias ItemVariant

	alias := &struct {
		ID            string `json:"id"`
		MeasurementID string `json:"measurement_id"`
		CreatedByID   string `json:"created_by_id"`
		UpdatedByID   string `json:"updated_by_id"`
		ItemID        string `json:"item_id"`
		AliasName     string `json:"alias_name"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	// Encrypt alias.ItemID when m.Item not nill
	// and the ID is setted
	if m.Item != nil && m.Item.ID != int64(0) {
		alias.ItemID = common.Encrypt(m.Item.ID)
	} else {
		alias.Item = nil
	}

	// Encrypt alias.MeasurementID when m.Measurement not nill
	// and the ID is setted
	if m.Measurement != nil && m.Measurement.ID != int64(0) {
		alias.MeasurementID = common.Encrypt(m.Measurement.ID)
	} else {
		alias.Measurement = nil
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

	if m.VariantName != "" {
		alias.AliasName = m.VariantName
	} else if m.ExternalName != "" {
		alias.AliasName = m.ExternalName
	} else {
		if m.Item != nil {
			alias.AliasName = m.Item.ItemName
		}
	}

	return json.Marshal(alias)
}

// Save inserting or updating ItemVariant struct into item_variant table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to item_variant.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *ItemVariant) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting item_variant data
// this also will truncated all data from all table
// that have relation with this item_variant.
func (m *ItemVariant) Delete() (err error) {
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
func (m *ItemVariant) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}

// Archive item_variant data by id
func (m *ItemVariant) Archive() (err error) {
	m.IsArchived = 1
	if err = m.Save("is_archived"); err == nil {
		// kalau semua item_variant dari item ini di archive maka item nya juga harus di archive
		var t int64
		o := orm.NewOrm()
		if e := o.Raw("select count(*) from item_variant where item_id = ? and is_archived = 0", m.Item.ID).QueryRow(&t); e == nil && t == 0 {
			o.Raw("update item set is_archived = 1 where id = ?", m.Item.ID).Exec()
		}
	}

	return
}

// Unarchive item_variant data by id
func (m *ItemVariant) Unarchive() (err error) {
	m.IsArchived = 0
	if err = m.Save("is_archived"); err == nil {
		// kalau salah satu item_variant dari item ini di unarchive maka item nya juga harus di unarchive
		o := orm.NewOrm()
		o.Raw("update item set is_archived = 0 where id = ?", m.Item.ID).Exec()
	}

	return
}

// DeleteItemVariant for delete item_variant data by id
func (m *ItemVariant) DeleteItemVariant() (err error) {
	m.IsDeleted = 1
	if err = m.Save("is_deleted"); err == nil {
		// kalau semua item_variant dari item ini di delete maka item nya juga harus di delete
		var t int64
		o := orm.NewOrm()
		if e := o.Raw("select count(*) from item_variant where item_id = ? and is_deleted = 0", m.Item.ID).QueryRow(&t); e == nil && t == 0 {
			o.Raw("update item set is_deleted = 1 where id = ?", m.Item.ID).Exec()
		}
	}

	return
}
