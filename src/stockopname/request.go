// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stockopname

import (
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create user process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	RecognitionDate  time.Time         `json:"recognition_date" valid:"required"`
	Note             string            `json:"note"`
	StockopnameItems []stockopnameItem `json:"stockopname_items" valid:"required"`
}

type stockopnameItem struct {
	ItemVariantStockID string  `json:"item_variant_stock_id" valid:"required"`
	Quantity           float32 `json:"quantity" valid:"required|gte:0"`
	Note               string  `json:"note"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	var stockIDS []int64

	for _, i := range r.StockopnameItems {
		// validasi item_variant_stock_id
		var ivStockID int64
		if ivStockID, e = common.Decrypt(i.ItemVariantStockID); e != nil {
			o.Failure("item_variant_stock_id", "item_variant_stock_id is not valid")
		} else {
			// item variant stock id gak boleh sama
			if util.HasElem(stockIDS, ivStockID) {
				o.Failure("item_variant_stock_id", "item_variant_stock_id cannot same")
			}
			ivStock := &model.ItemVariantStock{ID: ivStockID}
			if e = ivStock.Read(); e != nil {
				o.Failure("item_variant_stock_id", "item_variant_stock_id doesn't exist")
			}
		}
		stockIDS = append(stockIDS, ivStockID)
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *createRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *createRequest) Transform(user *model.User) *model.Stockopname {
	code, _ := util.CodeGen("code_stockopname", "stockopname")
	stockopname := &model.Stockopname{
		Code:            code,
		RecognitionDate: r.RecognitionDate,
		Note:            r.Note,
		CreatedBy:       user,
		CreatedAt:       time.Now(),
	}
	var stockopnameItems []*model.StockopnameItem
	for _, i := range r.StockopnameItems {
		ivStockID, _ := common.Decrypt(i.ItemVariantStockID)
		stockItem := &model.StockopnameItem{
			ItemVariantStock: &model.ItemVariantStock{ID: ivStockID},
			Quantity:         i.Quantity,
			Note:             i.Note,
		}
		stockopnameItems = append(stockopnameItems, stockItem)
	}
	stockopname.StockopnameItems = stockopnameItems

	return stockopname
}
