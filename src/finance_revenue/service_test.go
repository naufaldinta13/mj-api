// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package financeRevenue

import (
	"testing"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestGetFinanceRevenuesNoData(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM finance_revenue").Exec()
	rq := orm.RequestQuery{}

	revenue := model.DummyFinanceRevenue()
	revenue.IsDeleted = 1
	revenue.Save()

	m, total, e := GetFinanceRevenues(&rq)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, m)
	assert.NoError(t, e)
}

func TestGetFinanceRevenuesSuccess(t *testing.T) {
	model.DummyFinanceRevenue()

	qs := orm.RequestQuery{}
	_, _, e := GetFinanceRevenues(&qs)
	assert.NoError(t, e, "Data should be exists.")
}

func TestSumAmountFinanceRevenue(t *testing.T) {

	fr1 := model.DummyFinanceRevenue()
	fr1.RefID = 1
	fr1.RefType = "sales_invoice"
	fr1.Amount = 10000
	fr1.IsDeleted = 0
	fr1.Save()

	fr2 := model.DummyFinanceRevenue()
	fr2.RefID = 1
	fr2.RefType = "sales_invoice"
	fr2.Amount = 10000
	fr2.IsDeleted = 0
	fr2.Save()

	// harusnya hasilnya 20000
	total := SumAmountFinanceRevenue(1, "sales_invoice", 0)
	assert.Equal(t, total, float64(20000))

	// jika parameter id diisi maka nilai amount fr2 tidak ikut ditambahkan
	total = SumAmountFinanceRevenue(1, "sales_invoice", fr2.ID)
	assert.Equal(t, total, float64(10000))

	// tes dengan data finance revenue yg terhapus
	fr3 := model.DummyFinanceRevenue()
	fr3.RefID = 1
	fr3.RefType = "sales_invoice"
	fr3.Amount = 10000
	fr3.IsDeleted = 1
	fr3.Save()

	// harusnya hasilnya tetap 20000
	total = SumAmountFinanceRevenue(1, "sales_invoice", 0)
	assert.Equal(t, total, float64(20000))
}

func TestShowFinanceRevenue(t *testing.T) {
	finance := model.DummyFinanceRevenue()
	finance.IsDeleted = 0
	finance.Save("IsDeleted")

	data, e := ShowFinanceRevenue("id", finance.ID)
	assert.NoError(t, e)
	assert.NotEmpty(t, data)
	assert.Equal(t, data.ID, finance.ID)
}

// TestApproveRevenue dengan dummy ref type salah
func TestApproveRevenue(t *testing.T) {
	// buat dummy

	rev := model.DummyFinanceRevenue()
	rev.IsDeleted = int8(0)
	rev.DocumentStatus = "uncleared"
	rev.Amount = float64(10000)
	rev.Save()
	rev.RefType = "test"

	// test
	e := ApproveRevenue(rev)
	assert.Error(t, e)
}

// TestApproveRevenue1 dengan dummy sales invoice
func TestApproveRevenue1(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM finance_revenue").Exec()
	// buat dummy
	customer := model.DummyPartnership()

	inv := model.DummySalesInvoice()
	inv.TotalAmount = float64(20000)
	inv.TotalPaid = float64(15000)
	inv.DocumentStatus = "active"
	inv.IsDeleted = int8(0)
	inv.Save()

	so := inv.SalesOrder
	so.DocumentStatus = "active"
	so.Customer = customer
	so.TotalPaid = float64(15000)
	so.TotalCharge = float64(20000)
	so.InvoiceStatus = "active"
	so.ShipmentStatus = "finished"
	so.FulfillmentStatus = "finished"
	so.IsDeleted = int8(0)
	so.Save()

	customer.PartnershipType = "customer"
	customer.IsDeleted = int8(0)
	customer.TotalDebt = float64(15000)
	customer.Save()

	rev1 := model.DummyFinanceRevenue()
	rev1.IsDeleted = int8(0)
	rev1.DocumentStatus = "cleared"
	rev1.Amount = float64(10000)
	rev1.RefType = "sales_invoice"
	rev1.RefID = uint64(inv.ID)
	rev1.Save()

	rev2 := model.DummyFinanceRevenue()
	rev2.IsDeleted = int8(0)
	rev2.DocumentStatus = "cleared"
	rev2.Amount = float64(5000)
	rev2.RefType = "sales_invoice"
	rev2.RefID = uint64(inv.ID)
	rev2.Save()

	rev := model.DummyFinanceRevenue()
	rev.IsDeleted = int8(0)
	rev.DocumentStatus = "uncleared"
	rev.Amount = float64(5000)
	rev.RefType = "sales_invoice"
	rev.RefID = uint64(inv.ID)
	rev.Save()
	// test
	e := ApproveRevenue(rev)
	assert.NoError(t, e)

	rev.Read("ID")
	assert.Equal(t, "cleared", rev.DocumentStatus)

	inv.Read("ID")
	assert.Equal(t, "finished", inv.DocumentStatus)
	assert.Equal(t, float64(20000), inv.TotalPaid)

	so.Read("ID")
	assert.Equal(t, "finished", so.InvoiceStatus)
	assert.Equal(t, float64(20000), so.TotalPaid)

	customer.Read("ID")
	assert.Equal(t, float64(0), customer.TotalDebt)
}

// TestApproveRevenue2 dengan dummy purchase return
func TestApproveRevenue2(t *testing.T) {
	// buat dummy
	ret := model.DummyPurchaseReturn()
	ret.TotalAmount = float64(50000)
	ret.DocumentStatus = "active"
	ret.IsDeleted = int8(0)
	ret.Save()

	rev1 := model.DummyFinanceRevenue()
	rev1.IsDeleted = int8(0)
	rev1.DocumentStatus = "cleared"
	rev1.Amount = float64(15000)
	rev1.RefType = "purchase_return"
	rev1.RefID = uint64(ret.ID)
	rev1.Save()

	rev2 := model.DummyFinanceRevenue()
	rev2.IsDeleted = int8(0)
	rev2.DocumentStatus = "cleared"
	rev2.Amount = float64(30000)
	rev2.RefType = "purchase_return"
	rev2.RefID = uint64(ret.ID)
	rev2.Save()

	rev3 := model.DummyFinanceRevenue()
	rev3.IsDeleted = int8(1)
	rev3.DocumentStatus = "cleared"
	rev3.Amount = float64(20000)
	rev3.RefType = "purchase_return"
	rev3.RefID = uint64(ret.ID)
	rev3.Save()

	rev := model.DummyFinanceRevenue()
	rev.IsDeleted = int8(0)
	rev.DocumentStatus = "uncleared"
	rev.Amount = float64(5000)
	rev.RefType = "purchase_return"
	rev.RefID = uint64(ret.ID)
	rev.Save()

	// test
	e := ApproveRevenue(rev)
	assert.NoError(t, e)

	rev.Read("ID")
	assert.Equal(t, "cleared", rev.DocumentStatus)

	ret.Read("ID")
	assert.Equal(t, "finished", ret.DocumentStatus)

}
