// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package report

import (
	"testing"
	"time"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestGetReportSalesItems(t *testing.T) {

	soi := model.DummySalesOrderItem()
	soi.Save()

	ivs := model.DummyItemVariantStock()
	ivs.ItemVariant = soi.ItemVariant
	ivs.Save()

	ivsl := model.DummyItemVariantStockLog()
	ivsl.ItemVariantStock = ivs
	ivsl.Save()

	qs := orm.RequestQuery{}
	_, _, e := GetReportSalesItems(&qs)
	assert.NoError(t, e, "Data should be exists.")
}

func TestGetDetailItemVariantStockLog(t *testing.T) {
	_, e := GetDetailItemVariantStockLog("id", 1000)
	assert.Error(t, e, "Response should be error, beacuse there are no data yet.")

	c := model.DummyItemVariantStockLog()
	cd, e := GetDetailItemVariantStockLog("id", c.ID)
	assert.NoError(t, e, "Data should be exists.")
	assert.Equal(t, c.ID, cd.ID, "ID Response should be a same.")
}

func TestGetReportPurchaseItems(t *testing.T) {

	o := orm.NewOrm()
	o.Raw("delete from purchase_order").Exec()
	o.Raw("delete from partnership").Exec()

	m1 := model.DummyPurchaseOrderItem()
	m1.PurchaseOrder.IsDeleted = int8(0)
	m1.PurchaseOrder.Save()

	m2 := model.DummyPurchaseOrderItem()
	m2.PurchaseOrder.IsDeleted = int8(1)
	m2.PurchaseOrder.Save()

	qs := orm.RequestQuery{}
	_, _, e := GetReportPurchaseItems(&qs)
	assert.NoError(t, e, "Data should be exists.")
}

func TestGetTotalSalesItem(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()
	o.Raw("delete from sales_order_item").Exec()

	ic := model.DummyItemCategory()

	i := model.DummyItem()
	i.Category = ic
	i.Save()

	iv := model.DummyItemVariant()
	iv.Item = i
	iv.Save()

	ic2 := model.DummyItemCategory()

	i2 := model.DummyItem()
	i2.Category = ic2
	i2.Save()

	iv2 := model.DummyItemVariant()
	iv2.Item = i2
	iv2.Save()

	so := model.DummySalesOrder()
	so.DocumentStatus = "active"
	so.IsDeleted = 0
	so.Save()

	soi := model.DummySalesOrderItem()
	soi.Subtotal = 100000
	soi.Quantity = 20
	soi.SalesOrder = so
	soi.ItemVariant = iv
	soi.Save()

	soi2 := model.DummySalesOrderItem()
	soi2.Subtotal = 100000
	soi2.Quantity = 20
	soi2.SalesOrder = so
	soi2.ItemVariant = iv2
	soi2.Save()

	ts, e := GetTotalSalesItem("", "", "", "")
	assert.Equal(t, float64(200000), ts.Revenues)
	assert.Equal(t, float32(40), ts.TotalQuantityTransactions)
	assert.NoError(t, e, "Data should be exists.")

	// test function paramater not null
	ts2, _ := GetTotalSalesItem(common.Encrypt(ic.ID), common.Encrypt(iv.ID), "", "")
	assert.Equal(t, float64(100000), ts2.Revenues)
	assert.Equal(t, float32(20), ts2.TotalQuantityTransactions)
}

func TestGetTotalpurchaseItem(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("delete from purchase_order").Exec()
	o.Raw("delete from purchase_order_item").Exec()

	ic := model.DummyItemCategory()

	i := model.DummyItem()
	i.Category = ic
	i.Save()

	iv := model.DummyItemVariant()
	iv.Item = i
	iv.Save()

	ic2 := model.DummyItemCategory()

	i2 := model.DummyItem()
	i2.Category = ic2
	i2.Save()

	iv2 := model.DummyItemVariant()
	iv2.Item = i2
	iv2.Save()

	po := model.DummyPurchaseOrder()
	po.DocumentStatus = "active"
	po.IsDeleted = 0
	po.Save()

	poi := model.DummyPurchaseOrderItem()
	poi.Subtotal = 100000
	poi.Quantity = 20
	poi.PurchaseOrder = po
	poi.ItemVariant = iv2
	poi.Save()

	poi2 := model.DummyPurchaseOrderItem()
	poi2.Subtotal = 100000
	poi2.Quantity = 20
	poi2.PurchaseOrder = po
	poi2.ItemVariant = iv
	poi2.Save()

	ts2, _ := GetTotalPurchaseItem("", "")
	assert.Equal(t, float64(200000), ts2.Expenses)
	assert.Equal(t, float32(40), ts2.TotalQuantityTransactions)

	ts, e := GetTotalPurchaseItem(common.Encrypt(ic2.ID), common.Encrypt(iv2.ID))
	assert.Equal(t, float64(100000), ts.Expenses)
	assert.Equal(t, float32(20), ts.TotalQuantityTransactions)
	assert.NoError(t, e, "Data should be exists.")
}

func TestGetTotalSalesSummary(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("delete from sales_order").Exec()

	cust := model.DummyPartnership()
	cust.Save()

	so := model.DummySalesOrder()
	so.TotalCost = 100000
	so.IsDeleted = 0
	so.Customer = cust
	so.DocumentStatus = "new"
	so.TotalCharge = 50000
	so.RecognitionDate = time.Now()
	so.Save()

	so2 := model.DummySalesOrder()
	so2.TotalCost = 100000
	so2.IsDeleted = 0
	so2.Customer = cust
	so2.DocumentStatus = "new"
	so2.TotalCharge = 50000
	so2.RecognitionDate = time.Now()
	so2.Save()

	c := common.Encrypt(cust.ID)

	ts, e := GetTotalSales(c, "", "")
	assert.Equal(t, float64(200000), ts.Costs)
	assert.Equal(t, float64(100000), ts.Revenues)
	assert.Equal(t, int64(2), ts.Transactions)
	assert.NoError(t, e, "Data should be exists.")

	ts2, _ := GetTotalSales(c, "2017-10-12", "2117-12-19")
	assert.Equal(t, float64(200000), ts2.Costs)
	assert.Equal(t, float64(100000), ts2.Revenues)
	assert.Equal(t, int64(2), ts2.Transactions)

}

func TestGetTotalPurchaseSummary(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("delete from purchase_order").Exec()

	supp := model.DummyPartnership()
	supp.Save()

	po := model.DummyPurchaseOrder()
	po.IsDeleted = 0
	po.Supplier = supp
	po.DocumentStatus = "new"
	po.TotalCharge = 50000
	po.Save()

	po2 := model.DummyPurchaseOrder()
	po2.IsDeleted = 0
	po2.Supplier = supp
	po2.DocumentStatus = "new"
	po2.TotalCharge = 50000
	po2.Save()

	s := common.Encrypt(supp.ID)

	ts, e := GetTotalPurchase(s, "", "")
	assert.Equal(t, float64(100000), ts.Expenses)
	assert.Equal(t, int64(2), ts.Transactions)
	assert.NoError(t, e, "Data should be exists.")

	ts2, _ := GetTotalPurchase(s, "2017-10-12", "2117-12-19")
	assert.Equal(t, float64(100000), ts2.Expenses)

}
