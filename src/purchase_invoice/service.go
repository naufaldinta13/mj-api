// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchaseInvoice

import (
	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

// CreatePurchaseInvoice to create purchase invoice
func CreatePurchaseInvoice(r *createRequest) (pi *model.PurchaseInvoice, err error) {

	poID, _ := common.Decrypt(r.PurchaseOrder)
	po := &model.PurchaseOrder{ID: poID}

	if err = po.Read("ID"); err == nil {
		po.DocumentStatus = "active"
		po.InvoiceStatus = "active"
		po.Save("InvoiceStatus", "DocumentStatus")

		pi = r.Transform()
		err = pi.Save()
	}

	return
}

// CalculateTotalAmountPI Core library untuk menghitung total amount semua purchase invoice
// yang memiliki referensi ke purchase order
func CalculateTotalAmountPI(PoID int64, piID int64) (total float64, err error) {
	o := orm.NewOrm()

	err = o.Raw("SELECT SUM(pi.total_amount) from purchase_invoice pi "+
		"WHERE pi.purchase_order_id = ? AND pi.id != ? AND pi.is_deleted = 0;", PoID, piID).QueryRow(&total)

	return total, err
}

// GetPurchaseInvoice get all data Purchase Invoice that matched with query request parameters.
// returning slices of Purchase Invoice, total data without limit and error.
func GetPurchaseInvoice(rq *orm.RequestQuery) (m *[]model.PurchaseInvoice, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.PurchaseInvoice))

	// get total data
	if total, err = q.Filter("is_deleted", 0).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.PurchaseInvoice
	if _, err = q.Filter("is_deleted", 0).All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// ShowPurchaseInvoice find a single data Purchase Invoice using field and value condition.
func ShowPurchaseInvoice(field string, values ...interface{}) (*model.PurchaseInvoice, error) {
	m := new(model.PurchaseInvoice)
	o := orm.NewOrm().QueryTable(m)
	if err := o.Filter(field, values...).Filter("is_deleted", 0).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UpdatePurchaseInvoice to update purchase invoice
func UpdatePurchaseInvoice(r *updateRequest) (pi *model.PurchaseInvoice, err error) {

	pi = r.Transform(r.PI)
	if err = pi.Save("recognition_date", "note", "updated_by", "updated_at", "total_amount"); err == nil {

	}
	return pi, nil
}
