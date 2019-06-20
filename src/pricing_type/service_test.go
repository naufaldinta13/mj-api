// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pricingType

import (
	"testing"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestGetPricingTypes(t *testing.T) {
	model.DummyPricingType()
	qs := orm.RequestQuery{}
	_, _, e := GetPricingTypes(&qs)
	assert.NoError(t, e, "Data should be exists.")
}

func TestGetPricingTypeByID(t *testing.T) {
	pType := model.DummyPricingType()
	pType.ParentType = model.DummyPricingType()
	pType.Save()

	m, e := GetPricingTypeByID(pType.ID)
	assert.NoError(t, e)
	assert.Equal(t, pType.ID, m.ID)
	assert.Equal(t, pType.IsDefault, m.IsDefault)
	assert.Equal(t, pType.RuleType, m.RuleType)
	assert.Equal(t, pType.IsPercentage, m.IsPercentage)
	assert.Equal(t, pType.Nominal, m.Nominal)
	assert.Equal(t, pType.Note, m.Note)
	assert.Equal(t, pType.TypeName, m.TypeName)
	assert.Equal(t, pType.ParentType, m.ParentType)
}

func TestUpdateIsDefaultPricingType(t *testing.T) {
	orm.NewOrm().Raw("DELETE FROM pricing_type").Exec()
	// dummy untuk data yang lain
	pt := model.DummyPricingType()
	pt.IsDefault = 1
	pt.Save("IsDefault")

	// update is defaultnya = 1
	e := UpdateIsDefault()
	assert.NoError(t, e)
	// jika di update field is_default = 1 pricing type di atas
	// maka dummy pricing type yang lain => is default nya harus 0
	qs := orm.RequestQuery{}
	m, _, e := GetPricingTypes(&qs)
	assert.NoError(t, e)
	for _, x := range *m {
		assert.Equal(t, int8(0), x.IsDefault, "seharusnya is default = 0")
	}
}

func TestGetParentPricingTypes(t *testing.T) {
	orm.NewOrm().Raw("DELETE FROM pricing_type").Exec()

	pt := model.DummyPricingType()

	pt2 := model.DummyPricingType()
	pt2.ParentType = pt
	pt2.Save()

	qs := orm.RequestQuery{}
	_, total, e := GetParentPricingTypes(&qs)
	assert.NoError(t, e, "Data should be exists.")
	assert.Equal(t, int(total), int(1), "Harusnya ada 1 parent")
}

func TestGetPricingTypeByTypeName(t *testing.T) {
	orm.NewOrm().Raw("DELETE FROM pricing_type").Exec()
	mpt := model.DummyPricingType()

	pricing, err := getPricingByName(mpt.TypeName)
	assert.NotEmpty(t, pricing, "tidak boleh kosong")
	assert.NoError(t, err, "tidak ada error")
}
func TestGetDetailItemVariant(t *testing.T) {
	_, e := GetDetailItemVariant("id", 1000)
	assert.Error(t, e, "Response should be error, beacuse there are no data yet.")

	c := model.DummyItemVariant()
	cd, e := GetDetailItemVariant("id", c.ID)
	assert.NoError(t, e, "Data should be exists.")
	assert.Equal(t, c.ID, cd.ID, "ID Response should be a same.")
}

func TestGetPricingTypeByItemVariant(t *testing.T) {
	orm.NewOrm().Raw("DELETE FROM pricing_type").Exec()
	orm.NewOrm().Raw("DELETE FROM item_variant_prices").Exec()
	orm.NewOrm().Raw("DELETE FROM item_variant").Exec()
	orm.NewOrm().Raw("DELETE FROM item").Exec()

	///buat dummy pricing type yang tidak pake parent
	dmprt := model.DummyPricingType()
	dmprt.Nominal = 10000
	dmprt.ParentType = nil
	dmprt.TypeName = "pusat"
	dmprt.Save()

	///buat dummy pricing type yang pake parent dmprt
	dmprt1 := model.DummyPricingType()
	dmprt1.Nominal = 20000
	dmprt1.TypeName = "parent1"
	dmprt1.ParentType = dmprt
	dmprt1.Save()

	///buat dummy pricing type 2 yang pake parent dmprt
	dmprt2 := model.DummyPricingType()
	dmprt2.Nominal = 30000
	dmprt2.TypeName = "parent2"
	dmprt2.ParentType = dmprt
	dmprt2.Save()

	var ivp []*model.ItemVariantPrice
	dmivp := model.DummyItemVariantPrice()
	dmivp.PricingType = dmprt
	dmivp.Save()
	ivp = append(ivp, dmivp)

	dmiv := model.DummyItemVariant()
	dmiv.IsDeleted = 0
	dmiv.BasePrice = 1000
	dmiv.ItemVariantPrices = ivp
	dmiv.Save()

	pricing, err := GetPricingTypeByItemVariant(dmiv)
	assert.NotEmpty(t, pricing, "tidak boleh kosong")
	assert.NoError(t, err, "tidak ada error")
	assert.Equal(t, 3, len(pricing))
}
