// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package category

import (
	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create category process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	CategoryName string `json:"category_name" valid:"required"`
	Note         string `json:"note"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	_, e := getItemCategoryByName(r.CategoryName)
	if e == nil {
		o.Failure("category_name", "This item category name already exist")
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *createRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *createRequest) Transform() *model.ItemCategory {
	category := &model.ItemCategory{
		CategoryName: r.CategoryName,
		Note:         r.Note,
	}
	return category
}

// updateRequest data struct that stored request data when requesting an update category process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type updateRequest struct {
	ID           int64  `json:"-"`
	CategoryName string `json:"category_name" valid:"required"`
	Note         string `json:"note"`
}

// Validate implement validation.Requests interfaces.
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	category, _ := getItemCategoryByName(r.CategoryName)
	if r.CategoryName == category.CategoryName && r.ID != category.ID {
		o.Failure("category_name", "This item category name already exist")
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *updateRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *updateRequest) Transform(category *model.ItemCategory) *model.ItemCategory {
	category.CategoryName = r.CategoryName
	category.Note = r.Note
	return category
}

// deleteRequest data struct that stored request data when requesting an update category process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type deleteRequest struct {
	Category *model.ItemCategory `json:"category"`
}

// Validate implement validation.Requests interfaces.
func (r *deleteRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	if _, e := GetItemByCategory(r.Category); e == nil {
		o.Failure("category", "category was used by item ")
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *deleteRequest) Messages() map[string]string {
	return map[string]string{}
}
