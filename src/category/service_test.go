// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package category

import (
	"testing"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestGetItemCategories(t *testing.T) {
	category := model.DummyItemCategory()
	category.IsDeleted = 0
	category.Save("IsDeleted")

	qs := orm.RequestQuery{}
	m, _, e := GetItemCategories(&qs)
	assert.NoError(t, e, "Data should be exists.")
	for _, x := range *m {
		assert.Equal(t, int8(0), x.IsDeleted)
	}
}

func TestGetItemCategoryByName(t *testing.T) {
	category := model.DummyItemCategory()
	category.IsDeleted = 0
	category.Save()
	resCategory, err := getItemCategoryByName(category.CategoryName)
	assert.NoError(t, err, "Tidak boleh ada error")
	assert.NotNil(t, resCategory)
}

func TestGetItemCategoryByID(t *testing.T) {
	category := model.DummyItemCategory()
	category.IsDeleted = 0
	category.Save("IsDeleted")

	m, e := GetItemCategoryByID(category.ID)
	assert.NoError(t, e)
	assert.Equal(t, category.ID, m.ID)
	assert.Equal(t, category.CategoryName, m.CategoryName)
	assert.Equal(t, category.Note, m.Note)
	assert.Equal(t, category.IsDeleted, m.IsDeleted)

	category.IsDeleted = 1
	category.Save("IsDeleted")
	m, e = GetItemCategoryByID(category.ID)
	assert.Error(t, e)
	assert.Empty(t, m)
}

func TestGetItemByCategory(t *testing.T) {
	category := model.DummyItemCategory()
	category.IsDeleted = 0
	category.Save("IsDeleted")
	item := model.DummyItem()
	item.Category = category
	item.IsDeleted = 0
	item.Save("Category", "IsDeleted")
	m, e := GetItemByCategory(category)
	assert.NoError(t, e)
	assert.Equal(t, category.ID, m.Category.ID)
	assert.Equal(t, item.ID, m.ID)
	assert.Equal(t, item.IsDeleted, m.IsDeleted)
	assert.Equal(t, item.Note, m.Note)
	category.IsDeleted = 1
	category.Save("IsDeleted")
	m, e = GetItemByCategory(category)
	assert.Error(t, e)
	assert.Empty(t, m)
}
