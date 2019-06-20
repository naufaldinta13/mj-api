// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package directPlacement

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/inventory"

	"git.qasico.com/cuxs/orm"
)

// ShowDirectPlacement untuk mengambil data detail berdasarkan parameter beserta direct placement item
func ShowDirectPlacement(field string, values ...interface{}) (*model.DirectPlacement, error) {
	direct := new(model.DirectPlacement)
	o := orm.NewOrm()
	if e := o.QueryTable(direct).Filter(field, values...).RelatedSel().Limit(1).One(direct); e != nil {
		return nil, e
	}
	o.LoadRelated(direct, "DirectPlacementItems", 2)
	return direct, nil
}

// CreateDirectPlacement untuk menyimpan data direct placement ke database beserta direct placement item
func CreateDirectPlacement(direct *model.DirectPlacement) (e error) {
	// save direct placement
	if e = direct.Save(); e == nil {
		for _, u := range direct.DirectPlacementItems {
			u.DirectPlacment = &model.DirectPlacement{ID: direct.ID}
			if e = u.Save(); e == nil {
				inventory.FifoStockIn(u.ItemVariant, u.UnitPrice, u.Quantity, "direct_placement", uint64(u.DirectPlacment.ID))
			}
		}
	}
	return
}
