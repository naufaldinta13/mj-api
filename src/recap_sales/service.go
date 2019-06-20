// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package recapSales

import (
	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
)

// CreateRecapSales untuk menyimpan data recap sales ke database
func CreateRecapSales(rec *model.RecapSales) (e error) {
	// save recap
	if e = rec.Save(); e == nil {
		var so *model.SalesOrder

		// loop recap sales item
		for _, u := range rec.RecapSalesItems {
			u.RecapSales = &model.RecapSales{ID: rec.ID}

			// save
			if e = u.Save(); e == nil {
				so = &model.SalesOrder{ID: u.SalesOrder.ID}
				so.IsReported = 1
				so.Save("is_reported")
			} else {
				return
			}
		}
	}
	return
}

// GetRecapSales untuk mengambil semua data recap sales dari database
func GetRecapSales(rq *orm.RequestQuery) (m *[]model.RecapSales, total int64, e error) {
	// make new orm query
	q, _ := rq.Query(new(model.RecapSales))
	// get total data
	if total, e = q.Count(); e == nil && total != int64(0) {
		// get data requested
		var mx []model.RecapSales
		if _, e = q.All(&mx, rq.Fields...); e == nil {
			m = &mx
		}
	}
	return
}

// ShowRecapSales untuk mengambil detail data recap sales dari database
func ShowRecapSales(field string, values ...interface{}) (*model.RecapSales, error) {
	rs := new(model.RecapSales)
	o := orm.NewOrm()
	if e := o.QueryTable(rs).Filter(field, values...).RelatedSel().Limit(1).One(rs); e != nil {
		return nil, e
	}
	o.LoadRelated(rs, "RecapSalesItems", 2)
	return rs, nil
}
