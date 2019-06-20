// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package financeRevenue

import (
	"errors"
	"fmt"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/partnership"
	"git.qasico.com/mj/api/src/sales"
	"git.qasico.com/mj/api/src/sales_invoice"

	"git.qasico.com/cuxs/orm"
)

// GetFinanceRevenues get all data finance_revenue that matched with query request parameters.
// returning slices of workorder fullfillment, total data without limit and error.
func GetFinanceRevenues(rq *orm.RequestQuery) (m *[]model.FinanceRevenue, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.FinanceRevenue))
	q = q.Filter("is_deleted", 0)

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.FinanceRevenue
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// SumAmountFinanceRevenue untuk sum amount pada finance revenue berdasarkan ref_if dan ref_type
// parameter id untuk digunakan untuk update
func SumAmountFinanceRevenue(refID uint64, refType string, id int64) (totalAmount float64) {

	orm.NewOrm().Raw("select sum(amount) as totalAmount from finance_revenue where ref_id = ? and ref_type = ? and is_deleted = 0 and id != ?", refID, refType, id).QueryRow(&totalAmount)

	return totalAmount
}

// ShowFinanceRevenue untuk mengambil data detail finance order berdasarkan param
func ShowFinanceRevenue(field string, values ...interface{}) (*model.FinanceRevenue, error) {
	m := new(model.FinanceRevenue)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter(field, values...).Filter("is_deleted", int8(0)).RelatedSel(1).Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SumRevenueAmountCleared untuk menghitung total amount dari semua revenue sesuai param dan status cleared
func SumRevenueAmountCleared(refID uint64, refType string) (total float64, err error) {
	o := orm.NewOrm()
	err = o.Raw("SELECT SUM(re.amount) AS total FROM finance_revenue re WHERE re.ref_id = ? AND re.ref_type = ? AND re.document_status = ? AND re.is_deleted = 0;", refID, refType, "cleared").QueryRow(&total)
	return
}

// ApproveRevenue untuk melakukan approve pada finance revenue
func ApproveRevenue(rev *model.FinanceRevenue) (e error) {
	// ubah status revenue menjadi cleared
	rev.DocumentStatus = "cleared"
	if e = rev.Save("document_status"); e == nil {
		// cek tipe refType
		if rev.RefType == "sales_invoice" {
			e = salesInvoiceRevenue(rev)
		} else if rev.RefType == "purchase_return" {
			e = purchaseReturnRevenue(rev)
		} else {
			e = errors.New("refType is wrong")
		}
	}
	return
}

// salesInvoiceRevenue proses update total paid so invoice
func salesInvoiceRevenue(rev *model.FinanceRevenue) (e error) {
	var soInvoice *model.SalesInvoice
	if soInvoice, e = salesInvoice.ShowSalesInvoice("id", rev.RefID); e == nil {
		sales.CalculateTotalPaidSI(soInvoice)
		sales.CheckInvoiceStatus(soInvoice.SalesOrder)
		sales.CheckDocumentStatus(soInvoice.SalesOrder)
		partnership.CalculationTotalDebt(soInvoice.SalesOrder.Customer.ID)
		partnership.CalculationTotalSpend(soInvoice.SalesOrder.Customer.ID)
	}
	return
}

// purchaseReturnRevenue proses update amount pada purchase return revenue
func purchaseReturnRevenue(rev *model.FinanceRevenue) (e error) {
	var sum float64
	var purchaseRet = &model.PurchaseReturn{ID: int64(rev.RefID), IsDeleted: int8(0)}
	if sum, e = SumRevenueAmountCleared(rev.RefID, rev.RefType); e == nil {
		if e = purchaseRet.Read("ID", "IsDeleted"); e == nil {
			// cek amount purchase return dan amount dari semua revenue
			if purchaseRet.TotalAmount == sum {
				purchaseRet.DocumentStatus = "finished"
			} else {
				purchaseRet.DocumentStatus = "active"
			}
			e = purchaseRet.Save("document_status")
		}
	}
	return
}

func summaryRevenue(s string, rt string, pm string, ds string, de string) (total float64) {
	qb, _ := orm.NewQueryBuilder("mysql")

	qb = qb.Select("sum(amount) as total").From("finance_revenue")
	qb.Where(fmt.Sprintf("recognition_date >= '%s'", util.MaxDate()))

	if s != "" {
		qb.And(fmt.Sprintf("document_status = '%s'", s))
	}

	if rt != "" {
		qb.And(fmt.Sprintf("ref_type = '%s'", rt))
	}

	if pm != "" {
		qb.And(fmt.Sprintf("payment_method = '%s'", pm))
	}

	if ds != "" {
		qb.And(fmt.Sprintf("recognition_date >= '%s'", ds))
	}

	if de != "" {
		qb.And(fmt.Sprintf("recognition_date <= '%s'", de))
	}

	sql := qb.String()
	o := orm.NewOrm()

	o.Raw(sql).QueryRow(&total)

	return
}
