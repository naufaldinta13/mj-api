// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pricingType

import (
	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
)

// GetPricingTypes get all data pricing_type that matched with query request parameters.
// returning slices of PricingTypes, total data without limit and error.
func GetPricingTypes(rq *orm.RequestQuery) (m *[]model.PricingType, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.PricingType))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.PricingType
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetParentPricingTypes get all data pricing_type that matched with query request parameters.
// returning slices Parent of PricingTypes, total data without limit and error.
func GetParentPricingTypes(rq *orm.RequestQuery) (m *[]model.PricingType, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.PricingType))

	// get total data
	if total, err = q.Filter("parent_type_id__isnull", true).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.PricingType
	if _, err = q.Filter("parent_type_id__isnull", true).All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetPricingTypeByID untuk get data pricing type berdasarkan id pricing type
// return : data pricing type dan error
func GetPricingTypeByID(id int64) (m *model.PricingType, err error) {
	mx := new(model.PricingType)
	o := orm.NewOrm().QueryTable(mx)

	if err = o.Filter("id", id).RelatedSel().Limit(1).One(mx); err != nil {
		return nil, err
	}
	return mx, nil
}

// UpdateIsDefault mengganti data pricing type yang is default = 1 menjadi 0
func UpdateIsDefault() (e error) {
	m := &model.PricingType{IsDefault: int8(1)}
	if err := m.Read("IsDefault"); err == nil {
		m.IsDefault = int8(0)
		e = m.Save("is_default")
	}
	return e
}

// getPricingByName untuk mendapatkan pricing type berdasarkan type_name
func getPricingByName(name string) (*model.PricingType, error) {
	var pricing model.PricingType
	o := orm.NewOrm()
	err := o.Raw("SELECT * FROM pricing_type WHERE type_name = ?", name).QueryRow(&pricing)
	return &pricing, err
}

// GetDetailItemVariant untuk mengambil data item variant berdasarkan id
func GetDetailItemVariant(field string, values ...interface{}) (*model.ItemVariant, error) {
	m := new(model.ItemVariant)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	if _, err := o.Raw("select * from item_variant iv "+
		"inner join item_variant_stock ivs on ivs.item_variant_id = iv.id "+
		"inner join item_variant_stock_log ivsl on ivsl.item_variant_stock_id = ivs.id "+
		"where iv.id = ?", m.ID).QueryRows(&m.ItemVariantStockLogs); err != nil {

		return nil, err
	}

	o.LoadRelated(m, "ItemVariantStocks", 2)
	o.LoadRelated(m, "ItemVariantPrices", 2)

	return m, nil
}

// GetPricingTypeByItemVariant untuk mendapatkan pricing type berdasarkan item variant
func GetPricingTypeByItemVariant(iv *model.ItemVariant) (m []*model.PricingType, err error) {
	///cek dulu ada item variant pricenya ga
	///looping terlebih dahulu item variant pricenya
	if totivp := len(iv.ItemVariantPrices); totivp > 0 {
		for _, ivp := range iv.ItemVariantPrices {
			///append ke dalam array return pricing dari item variant price
			m = append(m, ivp.PricingType)

			///buat penampung parentnya pricing type
			var mpt []*model.PricingType
			///get parentnya pricing type dan masukan array yang sudah dibuat penampungnya
			totMpt, _ := orm.NewOrm().Raw("select * from pricing_type where parent_type_id = ?", ivp.PricingType.ID).QueryRows(&mpt)
			if totMpt > 0 {
				for _, pt := range mpt {
					///append ke dalam array return pricing dari parent pricing type
					m = append(m, pt)
				}
			}
		}
		return m, nil
	}

	return nil, err
}
