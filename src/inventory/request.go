// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package inventory

import (
	"fmt"
	"math"
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
	"git.qasico.com/mj/api/src/pricing_type"
)

// ArchiveItemVariant data struct that stored request data when requesting an update to archive item_variant process.
type ArchiveItemVariant struct {
	isArchived int8
	isDeleted  int8
}

// Validate implement validation.Requests interfaces.
func (r *ArchiveItemVariant) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if r.isDeleted == int8(1) {
		o.Failure("is_deleted", "cannot archive a deleted item_variant")
	} else {
		if r.isArchived == int8(1) {
			o.Failure("is_archived", "cannot archive an archived item_variant")
		}
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *ArchiveItemVariant) Messages() map[string]string {
	return map[string]string{}
}

// UnarchiveItemVariant data struct that stored request data when requesting an update to archive item_variant process.
type UnarchiveItemVariant struct {
	isArchived int8
	isDeleted  int8
}

// Validate implement validation.Requests interfaces.
func (r *UnarchiveItemVariant) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if r.isDeleted == int8(1) {
		o.Failure("is_deleted", "cannot unarchive a deleted item_variant")
	} else {
		if r.isArchived == int8(0) {
			o.Failure("is_archived", "cannot unarchive an unarchived item_variant")
		}
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *UnarchiveItemVariant) Messages() map[string]string {
	return map[string]string{}
}

// DeleteItemVariant data struct that stored request data when requesting an is_deleted item_variant process.
type DeleteItemVariant struct {
	ItemVariantID int64
	isArchived    int8
	isDeleted     int8
}

// Validate implement validation.Requests interfaces.
func (r *DeleteItemVariant) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if r.isDeleted == int8(1) {
		o.Failure("is_deleted", "cannot delete a deleted item_variant")
	} else {
		if r.isArchived == int8(0) {
			o.Failure("is_archived", "cannot delete an unarchived item_variant")
		}
	}

	if soItems, e := GetDetailSalesOrderItemByItemVariant(r.ItemVariantID); e == nil {
		for _, sx := range soItems {
			md := &model.SalesOrder{ID: sx.SalesOrder.ID}
			if err := md.Read("ID"); err == nil {
				if md.DocumentStatus != "finished" {
					o.Failure("sales_order", "document_status in sales_order must be finish")
				}
			}
		}
	}

	if poItems, e := GetDetailPurchaseOrderByItemVariant(r.ItemVariantID); e == nil {
		for _, sx := range poItems {
			md := &model.PurchaseOrder{ID: sx.PurchaseOrder.ID}
			if err := md.Read("ID"); err == nil {
				if md.DocumentStatus != "finished" {
					o.Failure("purchase_order", "document_status in purchase_order must be finish")
				}
			}
		}
	}

	// cek stock nya terlebih dahulu
	if !validDeleteVariant(r.ItemVariantID) {
		o.Failure("item_variant_id", "item_variant still have stock")
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *DeleteItemVariant) Messages() map[string]string {
	return map[string]string{}
}

// createItemRequest data struct that stored request data when requesting an create item process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createItemRequest struct {
	Session      *auth.SessionData
	ItemType     string           `json:"item_type" valid:"required|in:product,material,service"`
	ItemName     string           `json:"item_name" valid:"required"`
	Category     string           `json:"category_id" valid:"required"`
	Measurement  string           `json:"measurement_id" valid:"required"`
	Note         string           `json:"note"`
	HasVariant   int8             `json:"has_variant"`
	ItemVariants []VariantRequest `json:"item_variants" valid:"required"`
}

// VariantRequest menampung request item variant
type VariantRequest struct {
	ID                string                `json:"id"`
	ExternalName      string                `json:"external_name"`
	VariantName       string                `json:"variant_name"`
	MinimumStock      float32               `json:"minimum_stock" valid:"gte:0"`
	BasePrice         float64               `json:"base_price" valid:"required"`
	Note              string                `json:"note"`
	Image             string                `json:"image"`
	ItemVariantPrices []VariantPriceRequest `json:"item_variant_prices" valid:"required"`
}

// VariantPriceRequest menampung request item variant price
type VariantPriceRequest struct {
	ID          string  `json:"id"`
	PricingType string  `json:"pricing_type_id" valid:"required"`
	UnitPrice   float64 `json:"unit_price" valid:"required|gt:0"`
	Note        string  `json:"note"`
}

// Validate implement validation.Requests interfaces.
func (r *createItemRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if len(r.ItemName) > int(255) {
		o.Failure("item_name", "cannot be more than 255")
	}
	if len(r.Note) > int(255) {
		o.Failure("note", "cannot be more than 255")
	}
	itm := &model.Item{ItemName: r.ItemName, IsDeleted: int8(0)}
	if err := itm.Read("ItemName", "IsDeleted"); err == nil {
		o.Failure("item_name", "already exist")
	}

	if id, e := common.Decrypt(r.Category); e == nil {
		m := &model.ItemCategory{ID: id}
		if e = m.Read("ID"); e == nil {
			if m.IsDeleted == int8(1) {
				o.Failure("category_id", "category_id not exist")
			}
		} else {
			o.Failure("category_id", "category_id not exist")
		}
	} else {
		o.Failure("category_id", "category_id not valid")
	}

	if id, e := common.Decrypt(r.Measurement); e == nil {
		m := &model.Measurement{ID: id}
		if e = m.Read("ID"); e == nil {
			if m.IsDeleted == int8(1) {
				o.Failure("measurement_id", "measurement_id not exist")
			}
		} else {
			o.Failure("measurement_id", "measurement_id not exist")
		}
	} else {
		o.Failure("measurement_id", "measurement_id not valid")
	}

	data, _, _ := pricingType.GetParentPricingTypes(&orm.RequestQuery{})
	varNameCheck := make(map[string]bool)
	for iVar, u := range r.ItemVariants {
		if r.HasVariant == 1 && u.VariantName == "" {
			o.Failure(fmt.Sprintf("item_variants.%d.variant_name.invalid", iVar), "variant_name is required.")
		}

		if len(u.Note) > int(255) {
			o.Failure(fmt.Sprintf("item_variants.%d.note.invalid", iVar), "cannot be more than 255")
		}
		if len(u.VariantName) > int(255) {
			o.Failure(fmt.Sprintf("item_variants.%d.variant_name.invalid", iVar), "cannot be more than 255")
		}
		if u.VariantName != "" {
			if varNameCheck[u.VariantName] == true {
				o.Failure(fmt.Sprintf("item_variants.%d.variant_name.invalid", iVar), "cannot have same variant name")
			} else {
				varNameCheck[u.VariantName] = true
			}
		}

		if math.Mod(u.BasePrice, 1) != 0 {
			o.Failure(fmt.Sprintf("item_variants.%d.base_price.invalid", iVar), "Base price tidak boleh ada koma")
		}

		if len(u.ItemVariantPrices) != len(*data) {
			o.Failure(fmt.Sprintf("item_variants.item_variant_prices.%d.pricing_type_id.invalid", iVar), "Total inputted pricing type is less than existing pricing type")
		} else {
			for iPrice, x := range u.ItemVariantPrices {
				if idPriceType, err := common.Decrypt(x.PricingType); err == nil {
					pr := &model.PricingType{ID: idPriceType}
					if err = pr.Read("ID"); err == nil {
						if pr.ParentType != nil {
							if pr.ParentType.ID != int64(0) {
								o.Failure(fmt.Sprintf("item_variants.%d.item_variant_prices.%d.pricing_type_id.invalid", iVar, iPrice), "pricing_type_id not parent")
							}
						} else {
							if (*data)[iPrice].ID != pr.ID {
								o.Failure(fmt.Sprintf("item_variants.%d.item_variant_prices.%d.pricing_type_id.invalid", iVar, iPrice), "pricing_type_id is not same as in the pricing type list")
							}
						}
					} else {
						o.Failure(fmt.Sprintf("item_variants.%d.item_variant_prices.%d.pricing_type_id.invalid", iVar, iPrice), "pricing_type_id not exist")
					}
				} else {
					o.Failure(fmt.Sprintf("item_variants.%d.item_variant_prices.%d.pricing_type_id.invalid", iVar, iPrice), "pricing_type_id not valid")
				}
				if len(x.Note) > int(255) {
					o.Failure(fmt.Sprintf("item_variants.%d.item_variant_prices.note.%d.invalid", iVar, iPrice), "cannot be more than 255")
				}

				if math.Mod(x.UnitPrice, 1) != 0 {
					o.Failure(fmt.Sprintf("item_variants.%d.item_variant_prices.unit_price.%d.invalid", iVar, iPrice), "Unit price tidak boleh ada koma")
				}
			}
		}
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *createItemRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform untuk mengubah isi model.
func (r *createItemRequest) Transform() model.Item {
	cat, _ := common.Decrypt(r.Category)
	measure, _ := common.Decrypt(r.Measurement)
	var variant []*model.ItemVariant

	// buat banyak model item variant
	for _, u := range r.ItemVariants {
		var Price []*model.ItemVariantPrice

		// buat banyak model item variant price
		for _, x := range u.ItemVariantPrices {
			priceID, _ := common.Decrypt(x.PricingType)

			// buat model item variant price
			itmVariantPrice := &model.ItemVariantPrice{
				PricingType: &model.PricingType{ID: priceID},
				UnitPrice:   x.UnitPrice,
				Note:        x.Note,
			}
			Price = append(Price, itmVariantPrice)
		}

		hasExternal := int8(0)
		if u.ExternalName != "" {
			hasExternal = int8(1)
		}
		// buat model item variant
		itmVariant := &model.ItemVariant{
			Measurement:       &model.Measurement{ID: measure},
			Barcode:           util.GenerateEanCode(),
			ExternalName:      u.ExternalName,
			VariantName:       u.VariantName,
			Image:             u.Image,
			BasePrice:         u.BasePrice,
			Note:              u.Note,
			MinimumStock:      float32(common.FloatPrecision(float64(u.MinimumStock), 2)),
			AvailableStock:    float32(0),
			CommitedStock:     float32(0),
			HasExternalName:   hasExternal,
			IsArchived:        int8(0),
			IsDeleted:         int8(0),
			CreatedBy:         &model.User{ID: r.Session.User.ID},
			CreatedAt:         time.Now(),
			ItemVariantPrices: Price,
		}
		variant = append(variant, itmVariant)
	}

	// buat model item
	m := model.Item{
		Category:     &model.ItemCategory{ID: cat},
		ItemType:     r.ItemType,
		ItemName:     r.ItemName,
		Note:         r.Note,
		HasVariant:   r.HasVariant,
		IsArchived:   int8(0),
		IsDeleted:    int8(0),
		CreatedBy:    &model.User{ID: r.Session.User.ID},
		CreatedAt:    time.Now(),
		ItemVariants: variant,
	}
	return m
}

// archiveItem data struct that stored request data when requesting an update to archive item process.
type archiveItem struct {
	Item *model.Item
}

// Validate implement validation.Requests interfaces.
func (r *archiveItem) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	if r.Item.IsArchived == int8(1) {
		o.Failure("is_archived", "cannot archive an archived item")
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *archiveItem) Messages() map[string]string {
	return map[string]string{}
}

// unarchiveItem data struct that stored request data when requesting an update to unarchive item process.
type unarchiveItem struct {
	Item *model.Item
}

// Validate implement validation.Requests interfaces.
func (r *unarchiveItem) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	if r.Item.IsArchived == int8(0) {
		o.Failure("is_archived", "cannot unarchive an not archived item")
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *unarchiveItem) Messages() map[string]string {
	return map[string]string{}
}

// deleteItem data struct that stored request data when requesting an delete to delete item process.
type deleteItem struct {
	Item *model.Item
}

// Validate implement validation.Requests interfaces.
func (r *deleteItem) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	if r.Item.IsArchived == int8(0) {
		o.Failure("is_archived", "cannot delete not archived item")
	}

	for _, u := range r.Item.ItemVariants {
		if u.IsArchived == int8(0) {
			o.Failure("item_variants.is_archived.invalid", "cannot delete not archived item_variant")
		}

		if m, e := GetDetailPurchaseOrderByItemVariant(u.ID); e == nil {
			if len(m) > int(0) {
				o.Failure("item_variants.item_variant.invalid", "item variant is being used")
			}
		}

		if m, e := GetDetailSalesOrderItemByItemVariant(u.ID); e == nil {
			if len(m) > int(0) {
				o.Failure("item_variants.item_variant.invalid", "item variant is being used")
			}
		}

	}

	// cek stock nya terlebih dahulu
	if !validDeleteItem(r.Item.ID) {
		o.Failure("item_id", "item still have stock")
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *deleteItem) Messages() map[string]string {
	return map[string]string{}
}

// updateItemRequest data struct that stored request data when requesting an create item process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type updateItemRequest struct {
	Session      *auth.SessionData
	OldItem      *model.Item
	ItemName     string           `json:"item_name" valid:"required"`
	Measurement  string           `json:"measurement_id" valid:"required"`
	Category     string           `json:"category_id" valid:"required"`
	Note         string           `json:"note"`
	HasVariant   int8             `json:"has_variant"`
	ItemVariants []VariantRequest `json:"item_variants" valid:"required"`
}

// Validate implement validation.Requests interfaces.
func (r *updateItemRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if len(r.ItemName) > int(200) {
		o.Failure("item_name", "cannot be more than 200")
	}
	if len(r.Note) > int(255) {
		o.Failure("note", "cannot be more than 255")
	}
	itm := &model.Item{ItemName: r.ItemName, IsDeleted: int8(0)}
	if err := itm.Read("ItemName", "IsDeleted"); err == nil {
		fmt.Println("wewe", itm.ID, r.OldItem.ID)

		if itm.ID != r.OldItem.ID {
			o.Failure("item_name", "already exist")
		}

	}

	if id, e := common.Decrypt(r.Measurement); e == nil {
		m := &model.Measurement{ID: id, IsDeleted: int8(0)}
		if e = m.Read("ID", "IsDeleted"); e != nil {
			o.Failure("measurement_id", "measurement_id not exist")
		}
	} else {
		o.Failure("measurement_id", "measurement_id not valid")
	}

	if id, e := common.Decrypt(r.Category); e == nil {
		m := &model.ItemCategory{ID: id, IsDeleted: int8(0)}
		if e = m.Read("ID", "IsDeleted"); e != nil {
			o.Failure("category_id", "category_id not exist")
		}
	} else {
		o.Failure("category_id", "category_id not valid")
	}

	// get all pricing type
	var PrID []int64
	setter := orm.NewOrm()
	setter.Raw("select pt.id from pricing_type pt where pt.parent_type_id is null").QueryRows(&PrID)

	varNameCheck := make(map[string]bool)
	for iVar, u := range r.ItemVariants {
		if r.HasVariant == 1 && u.VariantName == "" {
			o.Failure(fmt.Sprintf("item_variants.%d.variant_name.invalid", iVar), "variant_name is required.")
		}

		if len(u.Note) > int(255) {
			o.Failure(fmt.Sprintf("item_variants.%d.note.invalid", iVar), "cannot be more than 255")
		}
		if len(u.VariantName) > int(255) {
			o.Failure(fmt.Sprintf("item_variants.%d.variant_name.invalid", iVar), "cannot be more than 255")
		}
		if u.VariantName != "" {
			if varNameCheck[u.VariantName] == true {
				o.Failure(fmt.Sprintf("item_variants.%d.variant_name.invalid", iVar), "cannot have same variant name")
			} else {
				varNameCheck[u.VariantName] = true
			}
		}
		if u.ID != "" {
			if varID, err := common.Decrypt(u.ID); err != nil {
				o.Failure(fmt.Sprintf("item_variants.%d.id.invalid", iVar), "id cannot be decrypt")
			} else {
				vr := model.ItemVariant{ID: varID}
				if err = vr.Read("ID"); err != nil {
					o.Failure(fmt.Sprintf("item_variants.%d.id.invalid", iVar), "id does not valid")
				}
			}
		}

		PriceTypeCheck := make(map[int64]bool)
		for iPrice, x := range u.ItemVariantPrices {
			if idPriceType, err := common.Decrypt(x.PricingType); err == nil {
				pr := &model.PricingType{ID: idPriceType}
				if err = pr.Read("ID"); err == nil {
					if pr.ParentType != nil {
						if pr.ParentType.ID != int64(0) {
							o.Failure(fmt.Sprintf("item_variants.%d.item_variant_prices.%d.pricing_type_id.invalid", iVar, iPrice), "pricing_type_id not parent")
						}
					}
					if PriceTypeCheck[idPriceType] == true {
						o.Failure(fmt.Sprintf("item_variants.%d.item_variant_prices.%d.pricing_type_id.invalid", iVar, iPrice), "cannot have same pricing_type_id")
					} else {
						PriceTypeCheck[idPriceType] = true
					}
				} else {
					o.Failure(fmt.Sprintf("item_variants.%d.item_variant_prices.%d.pricing_type_id.invalid", iVar, iPrice), "pricing_type_id not exist")
				}
			} else {
				o.Failure(fmt.Sprintf("item_variants.%d.item_variant_prices.%d.pricing_type_id.invalid", iVar, iPrice), "pricing_type_id not valid")
			}
			if x.ID != "" {
				if varPrID, err := common.Decrypt(x.ID); err != nil {
					o.Failure(fmt.Sprintf("item_variants.%d.item_variant_prices.%d.id.invalid", iVar, iPrice), "id cannot be decrypt")
				} else {
					vrp := &model.ItemVariantPrice{ID: varPrID}
					if err = vrp.Read("ID"); err != nil {
						o.Failure(fmt.Sprintf("item_variants.%d.item_variant_prices.%d.id.invalid", iVar, iPrice), "id does not valid")
					}
				}
			}
			if len(x.Note) > int(255) {
				o.Failure(fmt.Sprintf("item_variants.%d.item_variant_prices.note.%d.invalid", iVar, iPrice), "cannot be more than 255")
			}
		}
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *updateItemRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform untuk mengubah isi model.
func (r *updateItemRequest) Transform() model.Item {
	measure, _ := common.Decrypt(r.Measurement)
	cat, _ := common.Decrypt(r.Category)
	var variant []*model.ItemVariant
	// fill item variant
	for _, u := range r.ItemVariants {
		var Price []*model.ItemVariantPrice

		// buat banyak model item variant price
		for _, x := range u.ItemVariantPrices {
			priceID, _ := common.Decrypt(x.PricingType)

			// buat model item variant price
			itmVariantPrice := &model.ItemVariantPrice{
				PricingType: &model.PricingType{ID: priceID},
				UnitPrice:   x.UnitPrice,
				Note:        x.Note,
			}
			if x.ID != "" {
				idP, _ := common.Decrypt(x.ID)
				itmVariantPrice.ID = idP
			}
			Price = append(Price, itmVariantPrice)
		}

		hasExternal := int8(0)
		if u.ExternalName != "" {
			hasExternal = int8(1)
		}
		// buat model item variant
		itmVariant := &model.ItemVariant{
			Measurement:       &model.Measurement{ID: measure},
			ExternalName:      u.ExternalName,
			VariantName:       u.VariantName,
			Image:             u.Image,
			BasePrice:         u.BasePrice,
			Note:              u.Note,
			MinimumStock:      u.MinimumStock,
			HasExternalName:   hasExternal,
			ItemVariantPrices: Price,
		}
		if u.ID != "" {
			itmVarID, _ := common.Decrypt(u.ID)
			itmVariant.ID = itmVarID
			itmVariant.UpdatedBy = &model.User{ID: r.Session.User.ID}
			itmVariant.UpdatedAt = time.Now()
		} else {
			itmVariant.CreatedBy = &model.User{ID: r.Session.User.ID}
			itmVariant.CreatedAt = time.Now()
			itmVariant.IsArchived = int8(0)
			itmVariant.IsDeleted = int8(0)
			itmVariant.AvailableStock = float32(0)
			itmVariant.CommitedStock = float32(0)
			itmVariant.Barcode = util.GenerateEanCode()
		}
		variant = append(variant, itmVariant)
	}

	// fill item
	m := model.Item{
		ID:           r.OldItem.ID,
		UpdatedBy:    &model.User{ID: r.Session.User.ID},
		UpdatedAt:    time.Now(),
		Category:     &model.ItemCategory{ID: cat},
		Note:         r.Note,
		HasVariant:   r.HasVariant,
		ItemName:     r.ItemName,
		ItemVariants: variant,
	}
	return m
}
