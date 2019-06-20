// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchaseReturn

import (
	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
)

// CreatePurchaseReturn untuk simpan data purchase return ke database
func CreatePurchaseReturn(purchaseReturn *model.PurchaseReturn) (pReturn *model.PurchaseReturn, err error) {

	if err = purchaseReturn.Save(); err == nil {
		for _, row := range purchaseReturn.PurchaseReturnItems {
			row.PurchaseReturn = &model.PurchaseReturn{ID: purchaseReturn.ID}
			if err = row.Save(); err == nil {
				return purchaseReturn, err
			}
		}
	}

	return nil, err
}

// GetAllPurchaseReturn fungsi untuk mengambil semua data purchase return
func GetAllPurchaseReturn(rq *orm.RequestQuery) (preturn []*model.PurchaseReturn, total int64, err error) {
	purchase := new(model.PurchaseReturn)
	req, query := rq.Query(purchase)

	req = req.Filter("is_deleted", 0).RelatedSel()

	if total, err = req.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var preturns []*model.PurchaseReturn
	if _, err = req.All(&preturns, rq.Fields...); err == nil {
		for _, row := range preturns {
			query.LoadRelated(row, "PurchaseReturnItems")
		}
		return preturns, total, err
	}

	return nil, total, err
}

// GetDetailPurchaseReturn digunakan untuk mendapatkan detail data dari purchase return
func GetDetailPurchaseReturn(field string, value ...interface{}) (preturn *model.PurchaseReturn, err error) {
	purchase := new(model.PurchaseReturn)
	query := orm.NewOrm()
	if err = query.QueryTable(purchase).Filter(field, value).Filter("is_deleted", 0).RelatedSel().Limit(1).One(purchase); err == nil {
		query.LoadRelated(purchase, "PurchaseReturnItems", 3)
		return purchase, err
	}
	return nil, err
}

// updatePurchaseReturn untuk simpan data purchase return ke database
func updatePurchaseReturn(purchaseReturn *model.PurchaseReturn) (pReturn *model.PurchaseReturn, err error) {

	if err = purchaseReturn.Save(); err == nil {
		for _, row := range purchaseReturn.PurchaseReturnItems {
			row.PurchaseReturn = &model.PurchaseReturn{ID: purchaseReturn.ID}
			if err = row.Save(); err == nil {
				return purchaseReturn, err
			}
		}
	}

	return nil, err
}

// cancelPurchaseReturn untuk membatalkan purchase return
func cancelPurchaseReturn(purchaseReturn *model.PurchaseReturn) (purchase *model.PurchaseReturn, err error) {
	var finance []*model.FinanceRevenue
	var total int64
	purchaseReturn.DocumentStatus = "cancelled"
	if err = purchaseReturn.Save(); err == nil {
		err = purchaseReturn.Read()
		if purchaseReturn.DocumentStatus == "cancelled" {
			if total, err = orm.NewOrm().Raw("SELECT * FROM finance_revenue WHERE ref_id = ? AND ref_type = 'purchase_return' AND is_deleted = 0;", purchaseReturn.ID).QueryRows(&finance); err == nil && total != 0 {
				for _, row := range finance {
					row.IsDeleted = 1
					if err = row.Save("IsDeleted"); err != nil {
						return nil, err
					}
				}
				return purchaseReturn, err
			}
		}
	}
	return nil, err
}
