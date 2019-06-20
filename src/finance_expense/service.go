// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package financeExpense

import (
	"errors"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/sales_return"

	"git.qasico.com/cuxs/orm"
)

// GetFinanceExpenses get all data finance_expense that matched with query request parameters.
// returning slices of workorder fullfillment, total data without limit and error.
func GetFinanceExpenses(rq *orm.RequestQuery) (m *[]model.FinanceExpense, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.FinanceExpense))
	q = q.Filter("is_deleted", 0)

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.FinanceExpense
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// ShowFinanceExpense untuk get detail finance expense by id
func ShowFinanceExpense(field string, values ...interface{}) (*model.FinanceExpense, error) {
	m := new(model.FinanceExpense)
	o := orm.NewOrm().QueryTable(m)
	if err := o.Filter(field, values...).Filter("is_deleted", 0).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetDetailPurchaseInvoice untuk mengambil data detail purchase invoice berdasarkan param
func GetDetailPurchaseInvoice(field string, values ...interface{}) (*model.PurchaseInvoice, error) {
	m := new(model.PurchaseInvoice)
	o := orm.NewOrm()
	if err := o.QueryTable(m).Filter(field, values...).Filter("is_deleted", int8(0)).RelatedSel(3).Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SumExpenseAmount untuk menghitung total amount dari semua expense sesui param
func SumExpenseAmount(refID uint64, refType string) (total float64, err error) {
	o := orm.NewOrm()
	err = o.Raw("SELECT SUM(fe.amount) AS amount FROM finance_expense fe WHERE fe.ref_id = ? AND fe.ref_type = ? AND fe.document_status = 'cleared' AND fe.is_deleted = 0;", refID, refType).QueryRow(&total)
	return
}

// ApproveExpense untuk melakukan approve pada finance expense
func ApproveExpense(exp *model.FinanceExpense) (e error) {
	// ubah status expense menjadi cleared
	exp.DocumentStatus = "cleared"
	if e = exp.Save("document_status"); e == nil {
		// cek tipe refType
		if exp.RefType == "purchase_invoice" {
			e = purchaseInvoiceExpense(exp)
		} else if exp.RefType == "sales_return" {
			e = salesReturnExpense(exp)
		} else {
			e = errors.New("refType is wrong")
		}
	}
	return
}

// purchaseInvoiceExpense proses update total paid po invoice
func purchaseInvoiceExpense(exp *model.FinanceExpense) (e error) {
	var poInvoice *model.PurchaseInvoice
	if poInvoice, e = GetDetailPurchaseInvoice("id", exp.RefID); e == nil {
		// update total paid po invoice
		poInvoice.TotalPaid = poInvoice.TotalPaid + exp.Amount
		// update total paid po
		poInvoice.PurchaseOrder.TotalPaid = poInvoice.PurchaseOrder.TotalPaid + exp.Amount
		// update total credit partnership po
		poInvoice.PurchaseOrder.Supplier.TotalCredit = poInvoice.PurchaseOrder.Supplier.TotalCredit - exp.Amount
		// cek status dan save yang diupdate
		e = purchaseInvoiceExpenseStatusCheck(exp, poInvoice)
	}
	return
}

// purchaseInvoiceExpenseStatusCheck untuk mengecek status pembayaran po
func purchaseInvoiceExpenseStatusCheck(exp *model.FinanceExpense, poInvoice *model.PurchaseInvoice) (e error) {
	var sum float64
	if sum, e = SumExpenseAmount(exp.RefID, exp.RefType); e == nil {
		// cek total amount dan total charge po invoice dan ubah invoice status po
		if sum == poInvoice.TotalAmount {
			poInvoice.DocumentStatus = "finished"
		} else {
			poInvoice.DocumentStatus = "active"
		}
		// cek total paid po dan total charge po
		if poInvoice.PurchaseOrder.TotalPaid == poInvoice.PurchaseOrder.TotalCharge {
			poInvoice.PurchaseOrder.InvoiceStatus = "finished"
		} else {
			poInvoice.PurchaseOrder.InvoiceStatus = "active"
		}
		// cek  semua status po
		if poInvoice.PurchaseOrder.InvoiceStatus == "finished" && poInvoice.PurchaseOrder.ReceivingStatus == "finished" {
			poInvoice.PurchaseOrder.DocumentStatus = "finished"
		} else {
			poInvoice.PurchaseOrder.DocumentStatus = "active"
		}
		// save all update
		if e = poInvoice.PurchaseOrder.Supplier.Save("total_credit"); e == nil {
			if e = poInvoice.PurchaseOrder.Save("total_paid", "document_status", "invoice_status"); e == nil {
				e = poInvoice.Save("total_paid", "document_status")
			}
		}
	}
	return
}

// salesReturnExpense proses update amount pada sales return
func salesReturnExpense(exp *model.FinanceExpense) (e error) {
	var sum float64
	var salesRet *model.SalesReturn
	if sum, e = SumExpenseAmount(exp.RefID, exp.RefType); e == nil {
		if salesRet, e = salesReturn.ShowSalesReturn("id", exp.RefID); e == nil {
			// cek amount sales return dan amount dari semua expense
			if salesRet.TotalAmount == sum {
				salesRet.DocumentStatus = "finished"
			} else {
				salesRet.DocumentStatus = "active"
			}
			e = salesRet.Save("document_status")
		}
	}
	return
}
