// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package directPlacement

import (
	"fmt"
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create directPlacement process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	Session              *auth.SessionData
	Note                 string          `json:"note"`
	DirectPlacementItems []PlacementItem `json:"direct_placement_items" valid:"required"`
}

// PlacementItem untuk menampung direct placement item
type PlacementItem struct {
	ItemVariant string  `json:"item_variant_id" valid:"required"`
	Quantity    float32 `json:"quantity" valid:"required|gt:0"`
	UnitPrice   float64 `json:"unit_price" valid:"required"`
	TotalPrice  float64 `json:"total_price" valid:"required|gt:0"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	// check item variant
	for i, directItm := range r.DirectPlacementItems {
		if id, e := common.Decrypt(directItm.ItemVariant); e != nil {
			o.Failure(fmt.Sprintf("direct_placement_items.%d.item_variant_id.invalid", i), "not valid")
		} else {
			variant := &model.ItemVariant{ID: id, IsDeleted: int8(0)}
			if e = variant.Read("ID", "IsDeleted"); e != nil {
				o.Failure(fmt.Sprintf("direct_placement_items.%d.item_variant_id.invalid", i), "was deleted")
			}
			variant = &model.ItemVariant{ID: id, IsArchived: int8(0)}
			if e = variant.Read("ID", "IsArchived"); e != nil {
				o.Failure(fmt.Sprintf("direct_placement_items.%d.item_variant_id.invalid", i), "was archived")
			}
		}
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *createRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *createRequest) Transform() model.DirectPlacement {

	// transform direct placement item
	var placementItem []*model.DirectPlacementItem
	for _, u := range r.DirectPlacementItems {
		id, _ := common.Decrypt(u.ItemVariant)
		quantity := common.FloatPrecision(float64(u.Quantity), 2)
		unitPrice := common.FloatPrecision(u.TotalPrice/float64(quantity), 2)

		dirItem := &model.DirectPlacementItem{
			ItemVariant: &model.ItemVariant{ID: id},
			Quantity:    float32(quantity),
			TotalPrice:  u.TotalPrice,
			UnitPrice:   unitPrice,
		}
		placementItem = append(placementItem, dirItem)
	}

	// transform direct placement
	m := model.DirectPlacement{
		CreatedBy:            r.Session.User,
		CreatedAt:            time.Now(),
		Note:                 r.Note,
		DirectPlacementItems: placementItem,
	}
	return m
}
