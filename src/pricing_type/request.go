// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pricingType

import (
	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

type createRequest struct {
	ParentTypeID string  `json:"parent_type_id"`
	TypeName     string  `json:"type_name" valid:"required|lte:45"`
	Note         string  `json:"note" valid:"required"`
	RuleType     string  `json:"rule_type" valid:"in:increment,decrement"`
	Nominal      float64 `json:"nominal" valid:"gte:0"`
	IsPercentage int8    `json:"is_percentage" valid:"in:0,1"`
	IsDefault    int8    `json:"is_default" valid:"in:0,1"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if r.ParentTypeID != "" {
		if parentID, e := common.Decrypt(r.ParentTypeID); e != nil {
			o.Failure("parent_type_id", "parent_type_id cannot be decrypt")
		} else {
			var m *model.PricingType
			if m, e = GetPricingTypeByID(parentID); e != nil {
				o.Failure("parent_type_id", "parent_type_id doesn't exist")
			} else {
				if m.ParentType != nil {
					o.Failure("parent_type_id", "not valid")
				}
			}
		}
		// validasi required rule_type
		if r.RuleType == "" {
			o.Failure("rule_type", "rule_type is required")
		}
		// validasi jika is percentage =1 tapi nominal > 100%
		if r.IsPercentage == 1 && r.Nominal > 100 {
			o.Failure("nominal", "should not be more than 100%")
		}

		// validasi nominal tidak boleh lebih kecil dari 0
		if r.Nominal <= 0 {
			o.Failure("nominal", "the field is required.")
		}
	} else {
		// jika parent_type_id gak diisi maka rule_type, nominal, gak bisa diisi
		if r.RuleType != "" {
			o.Failure("rule_type", "rule_type can not filled if parent_type_id is empty")
		}
		// default ruleType
		if r.RuleType == "" {
			r.RuleType = "none"
		}
		if r.Nominal != 0 {
			o.Failure("nominal", "nominal can not filled if parent_type_id is empty")
		}
	}
	// type name must unique
	if _, err := getPricingByName(r.TypeName); err == nil {
		o.Failure("type_name", "Type name is already exists")
	}
	// validasi apabila data pricing type yang akan di create, is defaultnya = 1
	if r.IsDefault == int8(1) {
		if e := UpdateIsDefault(); e != nil {
			o.Failure("is_default", "fail to change other is_default to 0")
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
func (r *createRequest) Transform() (pricing *model.PricingType) {
	ParentID, _ := common.Decrypt(r.ParentTypeID)

	pricing = &model.PricingType{
		TypeName:     r.TypeName,
		RuleType:     r.RuleType,
		Note:         r.Note,
		IsPercentage: r.IsPercentage,
		Nominal:      r.Nominal,
		IsDefault:    r.IsDefault,
	}

	if r.ParentTypeID == "" {
		pricing.ParentType = nil
	} else {
		pricing.ParentType = &model.PricingType{ID: ParentID}
	}

	return
}

// updateRequest data struct that stored request data when requesting an update accountingaccount process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type updateRequest struct {
	PricingType  *model.PricingType `json:"-"`
	ParentTypeID string             `json:"parent_type_id"`
	TypeName     string             `json:"type_name" valid:"required|lte:45"`
	RuleType     string             `json:"rule_type" valid:"in:increment,decrement"`
	Nominal      float64            `json:"nominal"`
	IsPercentage int8               `json:"is_percentage"`
	Note         string             `json:"note" valid:"required"`
	IsDefault    int8               `json:"is_default"`
}

// Validate implement validation.Requests interfaces.
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	// validasi parent_type_id
	if r.ParentTypeID != "" {
		if r.PricingType.ParentType == nil {
			o.Failure("parent_type_id", "cannot be filled")
		} else {
			if parentID, e := common.Decrypt(r.ParentTypeID); e != nil {
				o.Failure("parent_type_id", "parent_type_id cannot be decrypt")
			} else {
				var m *model.PricingType
				if m, e = GetPricingTypeByID(parentID); e != nil {
					o.Failure("parent_type_id", "parent_type_id doesn't exist")
				} else {
					if m.ParentType != nil {
						o.Failure("parent_type_id", "not valid")
					}
				}
			}
			// validasi required rule_type
			if r.RuleType == "" {
				o.Failure("rule_type", "rule_type is required")
			}
			// validasi jika is percentage =1 tapi nominal > 100%
			if r.IsPercentage == 1 && r.Nominal > 100 {
				o.Failure("nominal", "should not be more than 100%")
			}

			// validasi nominal tidak boleh lebih kecil dari 0
			if r.Nominal <= 0 {
				o.Failure("nominal", "the field is required.")
			}
		}
	} else {
		if r.PricingType.ParentType != nil {
			o.Failure("parent_type_id", "cannot empty")
		} else {
			// jika parent_type_id gak diisi maka rule_type, nominal, gak bisa diisi
			if r.RuleType != "" {
				o.Failure("rule_type", "rule_type can not filled if parent_type_id is empty")
			}
			// default ruleType
			if r.RuleType == "" {
				r.RuleType = "none"
			}
			if r.Nominal != 0 {
				o.Failure("nominal", "nominal can not filled if parent_type_id is empty")
			}
		}
	}

	// validasi apabila data pricing type yang akan di update , is defaultnya = 1
	// maka is default gak bisa di ubah menjadi 0
	if r.PricingType.IsDefault == 1 && r.IsDefault == 0 {
		o.Failure("is_default", "is_default cannot be updated")
	}

	// type name must unique
	if opt, err := getPricingByName(r.TypeName); err == nil && opt.ID != r.PricingType.ID {
		o.Failure("type_name", "Type name is already exists")
	}

	// validasi apabila data pricing type yang akan di update, is defaultnya = 0
	// dan yang diinput is_default menjadi 1
	// maka data pricing type lain diubah is_default =0
	if r.PricingType.IsDefault == 0 && r.IsDefault == 1 {
		if e := UpdateIsDefault(); e != nil {
			o.Failure("is_default", "fail to change other is_default to 0")
		}
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *updateRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *updateRequest) Transform(pricingType *model.PricingType) *model.PricingType {
	if r.ParentTypeID != "" {
		parentID, _ := common.Decrypt(r.ParentTypeID)
		pricingType.ParentType = &model.PricingType{ID: parentID}
	} else {
		pricingType.ParentType = nil
	}
	pricingType.Note = r.Note
	pricingType.IsDefault = r.IsDefault
	pricingType.RuleType = r.RuleType
	pricingType.IsPercentage = r.IsPercentage
	pricingType.Nominal = r.Nominal
	pricingType.TypeName = r.TypeName

	return pricingType
}
