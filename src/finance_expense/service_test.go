// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package financeExpense

import (
	"testing"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestGetFinanceExpensesNoData(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM finance_expense").Exec()
	rq := orm.RequestQuery{}

	expense := model.DummyFinanceExpense()
	expense.IsDeleted = 1
	expense.Save()

	m, total, e := GetFinanceExpenses(&rq)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, m)
	assert.NoError(t, e)
}

func TestGetFinanceExpensesSuccess(t *testing.T) {
	model.DummyFinanceExpense()

	qs := orm.RequestQuery{}
	_, _, e := GetFinanceExpenses(&qs)
	assert.NoError(t, e, "Data should be exists.")
}

func TestShowFinanceExpense(t *testing.T) {

	// seharusnya error karena tidak ada data dengan id 999999
	_, err := ShowFinanceExpense("id", 999999)
	assert.Error(t, err, "SHould be an error cause data with id 999999 doesn't exist")

	// tes mengambil data yang is deleted nya 1
	d := model.DummyFinanceExpense()
	d.IsDeleted = 1
	d.Save()
	_, err1 := ShowFinanceExpense("id", d.ID)
	assert.Error(t, err1)

	// seharusnya tidak error karena ada datanya
	c := model.DummyFinanceExpense()
	c.IsDeleted = 0
	c.Save()
	cd, e := ShowFinanceExpense("id", c.ID)
	assert.NoError(t, e, "Data should be exists.")
	assert.Equal(t, c.ID, cd.ID, "ID Response should be a same.")
}

// TestGetDetailPurchaseInvoice test get detail dari table purchase invoice
func TestGetDetailPurchaseInvoice(t *testing.T) {
	// buat dummy
	inv := model.DummyPurchaseInvoice()
	inv.IsDeleted = int8(0)
	inv.Save()

	invDel := model.DummyPurchaseInvoice()
	invDel.IsDeleted = int8(1)
	invDel.Save()

	// test
	m, e := GetDetailPurchaseInvoice("id", inv.ID)
	assert.NoError(t, e)
	assert.Equal(t, inv.DocumentStatus, m.DocumentStatus)
	assert.Equal(t, inv.TotalAmount, m.TotalAmount)
	assert.NotEmpty(t, m.PurchaseOrder)
	assert.Equal(t, inv.PurchaseOrder.TotalCharge, m.PurchaseOrder.TotalCharge)
	assert.NotEmpty(t, m.PurchaseOrder.Supplier)
	assert.Equal(t, inv.PurchaseOrder.Supplier.Email, m.PurchaseOrder.Supplier.Email)

	// test2
	m2, err := GetDetailPurchaseInvoice("id", invDel.ID)
	assert.Error(t, err)
	assert.Empty(t, m2)
}

// TestSumExpenseAmount test fungsi penjumlahan semua amount expense
func TestSumExpenseAmount(t *testing.T) {
	// buat dummy
	exp1 := model.DummyFinanceExpense()
	exp1.IsDeleted = int8(0)
	exp1.DocumentStatus = "cleared"
	exp1.Amount = float64(10000)
	exp1.RefType = "purchase_invoice"
	exp1.RefID = uint64(2)
	exp1.Save()

	exp2 := model.DummyFinanceExpense()
	exp2.IsDeleted = int8(1)
	exp2.DocumentStatus = "cleared"
	exp2.Amount = float64(10000)
	exp2.RefType = "purchase_invoice"
	exp2.RefID = uint64(2)
	exp2.Save()

	exp3 := model.DummyFinanceExpense()
	exp3.IsDeleted = int8(0)
	exp3.DocumentStatus = "cleared"
	exp3.Amount = float64(15000)
	exp3.RefType = "purchase_invoice"
	exp3.RefID = uint64(2)
	exp3.Save()

	exp4 := model.DummyFinanceExpense()
	exp4.IsDeleted = int8(0)
	exp4.DocumentStatus = "uncleared"
	exp4.Amount = float64(10000)
	exp4.RefType = "purchase_invoice"
	exp4.RefID = uint64(2)
	exp4.Save()

	exp5 := model.DummyFinanceExpense()
	exp5.IsDeleted = int8(0)
	exp5.DocumentStatus = "cleared"
	exp5.Amount = float64(10000)
	exp5.RefType = "sales_return"
	exp5.RefID = uint64(2)
	exp5.Save()

	// test
	m, e := SumExpenseAmount(uint64(2), "purchase_invoice")
	assert.NoError(t, e)
	assert.Equal(t, float64(25000), m)
}

// TestApproveExpenseError dengan input ref type salah
func TestApproveExpenseError(t *testing.T) {
	// buat dummy
	exp := model.DummyFinanceExpense()
	exp.IsDeleted = int8(0)
	exp.DocumentStatus = "uncleared"
	exp.Amount = float64(10000)
	exp.Save()
	exp.RefType = "test"

	// test
	e := ApproveExpense(exp)
	assert.Error(t, e)
}

// TestApproveExpense1 dengan dummy purchase invoice
func TestApproveExpense1(t *testing.T) {
	// buat dummy
	inv := model.DummyPurchaseInvoice()
	inv.TotalAmount = float64(50000)
	inv.TotalPaid = float64(40000)
	inv.DocumentStatus = "active"
	inv.IsDeleted = int8(0)
	inv.Save()

	po := inv.PurchaseOrder
	po.DocumentStatus = "active"
	po.TotalPaid = float64(40000)
	po.TotalCharge = float64(50000)
	po.InvoiceStatus = "active"
	po.ReceivingStatus = "finished"
	po.IsDeleted = int8(0)
	po.Save()

	supplier := po.Supplier
	supplier.PartnershipType = "supplier"
	supplier.IsDeleted = int8(0)
	supplier.TotalCredit = float64(20000)
	supplier.Save()

	exp1 := model.DummyFinanceExpense()
	exp1.IsDeleted = int8(0)
	exp1.DocumentStatus = "cleared"
	exp1.Amount = float64(20000)
	exp1.RefType = "purchase_invoice"
	exp1.RefID = uint64(inv.ID)
	exp1.Save()

	exp2 := model.DummyFinanceExpense()
	exp2.IsDeleted = int8(0)
	exp2.DocumentStatus = "cleared"
	exp2.Amount = float64(20000)
	exp2.RefType = "purchase_invoice"
	exp2.RefID = uint64(inv.ID)
	exp2.Save()

	exp := model.DummyFinanceExpense()
	exp.IsDeleted = int8(0)
	exp.DocumentStatus = "uncleared"
	exp.Amount = float64(10000)
	exp.RefType = "purchase_invoice"
	exp.RefID = uint64(inv.ID)
	exp.Save()

	// test
	e := ApproveExpense(exp)
	assert.NoError(t, e)

	exp.Read("ID")
	assert.Equal(t, "cleared", exp.DocumentStatus)

	inv.Read("ID")
	assert.Equal(t, "finished", inv.DocumentStatus)
	assert.Equal(t, float64(50000), inv.TotalPaid)

	po.Read("ID")
	assert.Equal(t, "finished", po.DocumentStatus)
	assert.Equal(t, "finished", po.InvoiceStatus)
	assert.Equal(t, float64(50000), po.TotalPaid)

	supplier.Read("ID")
	assert.Equal(t, float64(10000), supplier.TotalCredit)
}

// TestApproveExpense2 dengan dummy sales return
func TestApproveExpense2(t *testing.T) {
	// buat dummy
	ret := model.DummySalesReturn()
	ret.TotalAmount = float64(50000)
	ret.DocumentStatus = "active"
	ret.IsDeleted = int8(0)
	ret.Save()

	exp1 := model.DummyFinanceExpense()
	exp1.IsDeleted = int8(0)
	exp1.DocumentStatus = "cleared"
	exp1.Amount = float64(15000)
	exp1.RefType = "sales_return"
	exp1.RefID = uint64(ret.ID)
	exp1.Save()

	exp2 := model.DummyFinanceExpense()
	exp2.IsDeleted = int8(0)
	exp2.DocumentStatus = "cleared"
	exp2.Amount = float64(30000)
	exp2.RefType = "sales_return"
	exp2.RefID = uint64(ret.ID)
	exp2.Save()

	exp3 := model.DummyFinanceExpense()
	exp3.IsDeleted = int8(1)
	exp3.DocumentStatus = "cleared"
	exp3.Amount = float64(20000)
	exp3.RefType = "sales_return"
	exp3.RefID = uint64(ret.ID)
	exp3.Save()

	exp := model.DummyFinanceExpense()
	exp.IsDeleted = int8(0)
	exp.DocumentStatus = "uncleared"
	exp.Amount = float64(5000)
	exp.RefType = "sales_return"
	exp.RefID = uint64(ret.ID)
	exp.Save()

	// test
	e := ApproveExpense(exp)
	assert.NoError(t, e)

	exp.Read("ID")
	assert.Equal(t, "cleared", exp.DocumentStatus)

	ret.Read("ID")
	assert.Equal(t, "finished", ret.DocumentStatus)

}
