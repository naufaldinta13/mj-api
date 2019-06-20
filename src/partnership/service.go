// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package partnership

import (
	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
)

// CalculationTotalDebt ini untuk kalkulasi total hutang
//digunakan untuk sales order
func CalculationTotalDebt(partnershipID int64) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("update partnership p set p.total_debt = (SELECT sum(so.total_charge-so.total_paid) as total_debts FROM sales_order so "+
		"WHERE so.customer_id = ? and so.document_status != 'approved_cancel' and so.is_deleted = 0) where p.id = ?;", partnershipID, partnershipID).Exec()
	return err
}

// CalculationTotalSpend ini untuk kalkulasi menghabiskan total
//digunakan untuk sales order
func CalculationTotalSpend(partnershipID int64) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("update partnership p set p.total_spend = (SELECT sum(so.total_charge) as total_spends FROM sales_order so "+
		"WHERE so.customer_id = ? and so.document_status != 'approved_cancel' and so.is_deleted = 0) where p.id = ?;", partnershipID, partnershipID).Exec()
	return err
}

// CalculationTotalExpenditure ini untuk kalkulasi total pengeluaran
//digunakan untuk purchase order
func CalculationTotalExpenditure(partnershipID int64) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("update partnership p set p.total_expenditure = (SELECT sum(po.total_charge) as total_expenditures FROM purchase_order po "+
		"WHERE po.supplier_id = ? and po.document_status != 'cancelled' and po.is_deleted = 0) where p.id = ?;", partnershipID, partnershipID).Exec()
	return err
}

// CalculationTotalCredit ini untuk kalkulasi total piutang
//digunakan untuk purchase order
func CalculationTotalCredit(partnershipID int64) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("update partnership p set p.total_credit = (SELECT sum(po.total_charge) as total_credits FROM purchase_order po "+
		"WHERE po.supplier_id = ? and po.document_status != 'cancelled' and po.is_deleted = 0) where p.id = ?;", partnershipID, partnershipID).Exec()
	return err
}

// GetAllPartnerships untuk mengambil semua data partnership dari database
func GetAllPartnerships(rq *orm.RequestQuery) (m *[]model.Partnership, total int64, e error) {
	// make new orm query
	q, _ := rq.Query(new(model.Partnership))
	// get total data
	if total, e = q.Filter("is_deleted", int8(0)).Count(); e == nil && total != int64(0) {
		// get data requested
		var mx []model.Partnership
		if _, e = q.Filter("is_deleted", int8(0)).All(&mx, rq.Fields...); e == nil {
			m = &mx
		}
	}
	return
}

// GetPartnershipByField untuk mengambil semua data partnership dari database
func GetPartnershipByField(field string, values ...interface{}) (*model.Partnership, error) {
	ps := new(model.Partnership)
	o := orm.NewOrm().QueryTable(ps)
	if e := o.Filter(field, values...).RelatedSel().Limit(1).One(ps); e != nil {
		return nil, e
	}
	return ps, nil
}
