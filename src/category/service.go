// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package category

import (
	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
)

// GetItemCategories get all data item_category that matched with query request parameters.
// returning slices of item category, total data without limit and error.
func GetItemCategories(rq *orm.RequestQuery) (m *[]model.ItemCategory, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.ItemCategory))
	q = q.Filter("is_deleted", 0)

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.ItemCategory
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetItemCategoryByID untuk get data item category berdasarkan id item category
// return : data item category dan error
func GetItemCategoryByID(id int64) (m *model.ItemCategory, err error) {
	mx := new(model.ItemCategory)
	o := orm.NewOrm().QueryTable(mx)

	if err = o.Filter("id", id).Filter("is_deleted", 0).RelatedSel().Limit(1).One(mx); err != nil {
		return nil, err
	}
	return mx, nil
}

// getItemCategoryByName untuk get data item category berdasarkan nama item category
// return : data item category dan error
func getItemCategoryByName(CategoryName string) (*model.ItemCategory, error) {
	var category model.ItemCategory
	query := orm.NewOrm()
	e := query.Raw("SELECT * FROM item_category WHERE category_name = ? AND is_deleted = 0", CategoryName).QueryRow(&category)
	return &category, e
}

// GetItemByCategory untuk get data item berdasarkan item category
// return : data item dan error
func GetItemByCategory(category *model.ItemCategory) (m *model.Item, err error) {
	o := orm.NewOrm()
	if err = o.Raw("select i.* from item i inner join item_category c on i.category_id = c.id "+
		"where i.category_id = ? and i.is_deleted = ? and c.is_deleted = ? ", category.ID, 0, 0).QueryRow(&m); err != nil {
		return nil, err
	}
	return m, nil
}
