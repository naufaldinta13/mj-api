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
	orm.RegisterModel(new(DirectPlacement))
}

// DirectPlacement model for direct_placement table.
type DirectPlacement struct {
	ID                   int64                  `orm:"column(id);auto" json:"-"`
	CreatedBy            *User                  `orm:"column(created_by);rel(fk)" json:"created_by"`
	CreatedAt            time.Time              `orm:"column(created_at);type(timestamp)" json:"created_at"`
	Note                 string                 `orm:"column(note);null" json:"note"`
	DirectPlacementItems []*DirectPlacementItem `orm:"reverse(many)" json:"direct_placement_items,omitempty"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *DirectPlacement) MarshalJSON() ([]byte, error) {
	type Alias DirectPlacement

	alias := &struct {
		ID          string `json:"id"`
		CreatedByID string `json:"created_by_id"`
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

	return json.Marshal(alias)
}

// Save inserting or updating DirectPlacement struct into direct_placement table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to direct_placement.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *DirectPlacement) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Delete permanently deleting direct_placement data
// this also will truncated all data from all table
// that have relation with this direct_placement.
func (m *DirectPlacement) Delete() (err error) {
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
func (m *DirectPlacement) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
