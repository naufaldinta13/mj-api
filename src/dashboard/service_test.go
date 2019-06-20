// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dashboard

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/cuxs"
	"git.qasico.com/cuxs/orm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

//// GetParam ///////////////////////////////////////////////////////////////////////////////

// TestGetParamURL1 test get param url dengan rq param di url benar
func TestGetParamURL1(t *testing.T) {
	// buat context
	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/?year=2017&month=12", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	c := cuxs.NewContext(ctx)
	// test
	m := GetParamURL(c)
	assert.Equal(t, int(2017), m.Year)
	assert.Equal(t, time.December, m.Month)
}

// TestGetParamURL2 test get param url dengan rq param di url salah
func TestGetParamURL2(t *testing.T) {
	// buat context
	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/?year=20aa&month=as", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	c := cuxs.NewContext(ctx)
	// test
	m := GetParamURL(c)
	assert.Equal(t, int(time.Now().Year()), m.Year)
	assert.Equal(t, time.Now().Month(), m.Month)
}

// TestGetParamURL3 test get param url dengan rq param di url salah
func TestGetParamURL3(t *testing.T) {
	// buat context
	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/?year=-20171&month=13", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	c := cuxs.NewContext(ctx)
	// test
	m := GetParamURL(c)
	assert.Equal(t, int(time.Now().Year()), m.Year)
	assert.Equal(t, time.Now().Month(), m.Month)
}

// TestGetParamURL4 test get param url dengan rq param di url kosong
func TestGetParamURL4(t *testing.T) {
	// buat context
	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	c := cuxs.NewContext(ctx)
	// test
	m := GetParamURL(c)
	assert.Equal(t, int(time.Now().Year()), m.Year)
	assert.Equal(t, time.Now().Month(), m.Month)
}

//// GetTotalSales ///////////////////////////////////////////////////////////////////////////////

// TestGetTotalSales test func get jumlah sales order yang tidak di delete
func TestGetTotalSales(t *testing.T) {
	tm := time.Date(2012, time.November, 12, 15, 00, 00, 00, time.Local)
	// hapus database
	o := orm.NewOrm()
	o.Raw("DELETE FROM sales_order").Exec()
	// buat dummy
	so1 := model.DummySalesOrder()
	so1.DocumentStatus = "new"
	so1.TotalCharge = float64(21000)
	so1.TotalPaid = float64(10000)
	so1.IsDeleted = int8(0)
	so1.Save()

	so2 := model.DummySalesOrder()
	so2.DocumentStatus = "new"
	so2.TotalCharge = float64(5000)
	so2.TotalPaid = float64(5000)
	so2.IsDeleted = int8(1)
	so2.Save()

	so3 := model.DummySalesOrder()
	so3.DocumentStatus = "new"
	so3.TotalCharge = float64(2000)
	so3.TotalPaid = float64(2000)
	so3.IsDeleted = int8(0)
	so3.Save()

	so4 := model.DummySalesOrder()
	so4.DocumentStatus = "new"
	so4.TotalCharge = float64(500)
	so4.TotalPaid = float64(100)
	so4.IsDeleted = int8(0)
	so4.RecognitionDate = tm
	so4.Save()
	// test
	total, e := GetTotalSales(int(time.Now().Month()), time.Now().Year())
	assert.NoError(t, e)
	assert.Equal(t, int64(2), total.Total)
	assert.Equal(t, float64(23000), total.Charge)
	assert.Equal(t, float64(12000), total.Paid)
}

// TestGetTotalSalesEmpty test func get jumlah sales order yang tidak di delete,empty
func TestGetTotalSalesEmpty(t *testing.T) {
	tm := time.Date(2012, time.November, 12, 15, 00, 00, 00, time.Local)
	// hapus database
	o := orm.NewOrm()
	o.Raw("DELETE FROM sales_order").Exec()
	// buat dummy
	so1 := model.DummySalesOrder()
	so1.IsDeleted = int8(1)
	so1.Save()

	so2 := model.DummySalesOrder()
	so2.IsDeleted = int8(1)
	so2.Save()

	so3 := model.DummySalesOrder()
	so3.IsDeleted = int8(1)
	so3.Save()

	so4 := model.DummySalesOrder()
	so4.IsDeleted = int8(0)
	so4.RecognitionDate = tm
	so4.Save()
	// test
	total, e := GetTotalSales(int(time.Now().Month()), time.Now().Year())
	assert.NoError(t, e)
	assert.Equal(t, int64(0), total.Total)
}

//// GetTotalPurchase ///////////////////////////////////////////////////////////////////////////////

// TestGetTotalPurchase test func get jumlah purchase order yang tidak di delete
func TestGetTotalPurchase(t *testing.T) {
	tm := time.Date(2012, time.November, 12, 15, 00, 00, 00, time.Local)
	tm2 := time.Date(time.Now().Year(), time.Now().Month(), 12, 15, 00, 00, 00, time.Local)
	// hapus database
	o := orm.NewOrm()
	o.Raw("DELETE FROM purchase_order").Exec()
	// buat dummy
	po1 := model.DummyPurchaseOrder()
	po1.DocumentStatus = "new"
	po1.TotalCharge = float64(21000)
	po1.TotalPaid = float64(10000)
	po1.IsDeleted = int8(0)
	po1.RecognitionDate = tm2
	po1.Save()

	po2 := model.DummyPurchaseOrder()
	po2.TotalCharge = float64(5000)
	po2.TotalPaid = float64(1000)
	po2.DocumentStatus = "new"
	po2.IsDeleted = int8(1)
	po2.Save()

	po3 := model.DummyPurchaseOrder()
	po3.DocumentStatus = "new"
	po3.TotalCharge = float64(2000)
	po3.TotalPaid = float64(5000)
	po3.IsDeleted = int8(0)
	po3.Save()

	po4 := model.DummyPurchaseOrder()
	po4.DocumentStatus = "new"
	po4.TotalCharge = float64(100)
	po4.TotalPaid = float64(500)
	po4.IsDeleted = int8(0)
	po4.RecognitionDate = tm
	po4.Save()
	// test
	total, e := GetTotalPurchase(int(time.Now().Month()), time.Now().Year())
	assert.NoError(t, e)
	assert.Equal(t, int64(2), total.Total)
	assert.Equal(t, float64(23000), total.Charge)
	assert.Equal(t, float64(15000), total.Paid)
}

// TestGetTotalPurchaseEmpty test func get jumlah purchase order yang tidak di delete,empty
func TestGetTotalPurchaseEmpty(t *testing.T) {
	tm := time.Date(2012, time.November, 12, 15, 00, 00, 00, time.Local)
	// hapus database
	o := orm.NewOrm()
	o.Raw("DELETE FROM purchase_order").Exec()
	// buat dummy
	po1 := model.DummyPurchaseOrder()
	po1.DocumentStatus = "new"
	po1.IsDeleted = int8(1)
	po1.Save()

	po2 := model.DummyPurchaseOrder()
	po2.DocumentStatus = "new"
	po2.IsDeleted = int8(1)
	po2.Save()

	po3 := model.DummyPurchaseOrder()
	po3.DocumentStatus = "new"
	po3.IsDeleted = int8(1)
	po3.Save()

	po4 := model.DummyPurchaseOrder()
	po4.DocumentStatus = "new"
	po4.IsDeleted = int8(0)
	po4.RecognitionDate = tm
	po4.Save()
	// test
	total, e := GetTotalPurchase(int(time.Now().Month()), time.Now().Year())
	assert.NoError(t, e)
	assert.Equal(t, int64(0), total.Total)
}

//// GetTotalSalesReturn ///////////////////////////////////////////////////////////////////////////////

// TestGetTotalSalesReturn test func get jumlah sales return yang tidak di delete
func TestGetTotalSalesReturn(t *testing.T) {
	tm := time.Date(2012, time.November, 12, 15, 00, 00, 00, time.Local)
	tm2 := time.Date(time.Now().Year(), time.Now().Month(), 12, 15, 00, 00, 00, time.Local)
	// hapus database
	o := orm.NewOrm()
	o.Raw("DELETE FROM sales_return").Exec()
	// buat dummy
	sr1 := model.DummySalesReturn()
	sr1.IsDeleted = int8(0)
	sr1.DocumentStatus = "new"
	sr1.RecognitionDate = tm2
	sr1.Save()

	sr2 := model.DummySalesReturn()
	sr2.IsDeleted = int8(1)
	sr2.DocumentStatus = "new"
	sr2.Save()

	sr3 := model.DummySalesReturn()
	sr3.IsDeleted = int8(0)
	sr3.DocumentStatus = "new"
	sr3.Save()

	sr4 := model.DummySalesReturn()
	sr4.IsDeleted = int8(0)
	sr4.DocumentStatus = "new"
	sr4.RecognitionDate = tm
	sr4.Save()
	// test
	total, e := GetTotalSalesReturn(int(time.Now().Month()), time.Now().Year())
	assert.NoError(t, e)
	assert.Equal(t, int64(2), total)
}

// TestGetTotalSalesReturnEmpty test func get jumlah sales return yang tidak di delete,empty
func TestGetTotalSalesReturnEmpty(t *testing.T) {
	tm := time.Date(2012, time.November, 12, 15, 00, 00, 00, time.Local)
	// hapus database
	o := orm.NewOrm()
	o.Raw("DELETE FROM sales_return").Exec()
	// buat dummy
	sr1 := model.DummySalesReturn()
	sr1.IsDeleted = int8(1)
	sr1.Save()

	sr2 := model.DummySalesReturn()
	sr2.IsDeleted = int8(1)
	sr2.Save()

	sr3 := model.DummySalesReturn()
	sr3.IsDeleted = int8(1)
	sr3.Save()

	sr4 := model.DummySalesReturn()
	sr4.IsDeleted = int8(0)
	sr4.RecognitionDate = tm
	sr4.Save()
	// test
	total, e := GetTotalSalesReturn(int(time.Now().Month()), time.Now().Year())
	assert.NoError(t, e)
	assert.Equal(t, int64(0), total)
}

//// GetTotalItemTerjual ///////////////////////////////////////////////////////////////////////////////

// TestGetTotalItemTerjual test func get jumlah quantity soi terjual yang tidak di delete
func TestGetTotalItemTerjual(t *testing.T) {
	tm := time.Date(2012, time.November, 12, 15, 00, 00, 00, time.Local)
	tm2 := time.Date(time.Now().Year(), time.Now().Month(), 12, 15, 00, 00, 00, time.Local)
	// hapus database
	o := orm.NewOrm()
	o.Raw("DELETE FROM sales_order").Exec()
	// buat dummy
	so1 := model.DummySalesOrder()
	so1.IsDeleted = int8(0)
	so1.DocumentStatus = "new"
	so1.RecognitionDate = tm2
	so1.Save()
	soi1 := model.DummySalesOrderItem()
	soi1.Quantity = float32(100)
	soi1.SalesOrder = so1
	soi1.Save()
	soi12 := model.DummySalesOrderItem()
	soi12.Quantity = float32(50)
	soi12.SalesOrder = so1
	soi12.Save()

	so2 := model.DummySalesOrder()
	so2.IsDeleted = int8(1)
	so2.DocumentStatus = "new"
	so2.RecognitionDate = time.Now()
	so2.Save()
	soi2 := model.DummySalesOrderItem()
	soi2.Quantity = float32(150)
	soi2.SalesOrder = so2
	soi2.Save()

	so3 := model.DummySalesOrder()
	so3.IsDeleted = int8(0)
	so3.DocumentStatus = "new"
	so3.RecognitionDate = time.Now()
	so3.Save()
	soi3 := model.DummySalesOrderItem()
	soi3.Quantity = float32(200)
	soi3.SalesOrder = so3
	soi3.Save()

	so4 := model.DummySalesOrder()
	so4.IsDeleted = int8(0)
	so4.RecognitionDate = tm
	so4.DocumentStatus = "new"
	so4.Save()
	soi4 := model.DummySalesOrderItem()
	soi4.Quantity = float32(200)
	soi4.SalesOrder = so4
	soi4.Save()
	// test
	total, e := GetTotalItemTerjual(int(time.Now().Month()), time.Now().Year())
	assert.NoError(t, e)
	assert.Equal(t, float32(350), total)
}

//// GetTotalStockIn ///////////////////////////////////////////////////////////////////////////////

// TestGetTotalStockIn test func get jumlah quantity woi receiving yang tidak di delete
func TestGetTotalStockIn(t *testing.T) {
	tm := time.Date(2012, time.November, 12, 15, 00, 00, 00, time.Local)
	tm2 := time.Date(time.Now().Year(), time.Now().Month(), 12, 15, 00, 00, 00, time.Local)
	// hapus database
	o := orm.NewOrm()
	o.Raw("DELETE FROM workorder_receiving").Exec()
	// buat dummy
	wo1 := model.DummyWorkorderReceiving()
	wo1.IsDeleted = int8(0)
	wo1.RecognitionDate = time.Now()
	wo1.Save()
	woi1 := model.DummyWorkorderReceivingItem()
	woi1.Quantity = float32(50)
	woi1.WorkorderReceiving = wo1
	woi1.Save()

	wo2 := model.DummyWorkorderReceiving()
	wo2.IsDeleted = int8(1)
	wo2.RecognitionDate = time.Now()
	wo2.Save()
	woi2 := model.DummyWorkorderReceivingItem()
	woi2.Quantity = float32(150)
	woi2.WorkorderReceiving = wo2
	woi2.Save()

	wo3 := model.DummyWorkorderReceiving()
	wo3.IsDeleted = int8(0)
	wo3.RecognitionDate = tm2
	wo3.Save()
	woi3 := model.DummyWorkorderReceivingItem()
	woi3.Quantity = float32(70)
	woi3.WorkorderReceiving = wo3
	woi3.Save()

	wo4 := model.DummyWorkorderReceiving()
	wo4.IsDeleted = int8(0)
	wo4.RecognitionDate = tm
	wo4.Save()
	woi4 := model.DummyWorkorderReceivingItem()
	woi4.Quantity = float32(50)
	woi4.WorkorderReceiving = wo4
	woi4.Save()

	// test
	total, e := GetTotalStockIn(int(time.Now().Month()), time.Now().Year())
	assert.NoError(t, e)
	assert.Equal(t, float32(120), total)
}

//// GetTotalStockReturn ///////////////////////////////////////////////////////////////////////////////

// TestGetTotalStockReturn test func get jumlah quantity sales return item yang tidak di delete
func TestGetTotalStockReturn(t *testing.T) {
	tm := time.Date(2012, time.November, 12, 15, 00, 00, 00, time.Local)
	tm2 := time.Date(time.Now().Year(), time.Now().Month(), 12, 15, 00, 00, 00, time.Local)
	// hapus database
	o := orm.NewOrm()
	o.Raw("DELETE FROM sales_return").Exec()
	// buat dummy
	sr1 := model.DummySalesReturn()
	sr1.IsDeleted = int8(0)
	sr1.RecognitionDate = time.Now()
	sr1.Save()
	sri1 := model.DummySalesReturnItem()
	sri1.Quantity = float32(50)
	sri1.SalesReturn = sr1
	sri1.Save()

	sr2 := model.DummySalesReturn()
	sr2.IsDeleted = int8(1)
	sr2.RecognitionDate = time.Now()
	sr2.Save()
	sri2 := model.DummySalesReturnItem()
	sri2.Quantity = float32(150)
	sri2.SalesReturn = sr2
	sri2.Save()

	sr3 := model.DummySalesReturn()
	sr3.IsDeleted = int8(0)
	sr3.RecognitionDate = tm2
	sr3.Save()
	sri3 := model.DummySalesReturnItem()
	sri3.Quantity = float32(100)
	sri3.SalesReturn = sr3
	sri3.Save()

	sr4 := model.DummySalesReturn()
	sr4.IsDeleted = int8(0)
	sr4.RecognitionDate = tm
	sr4.Save()
	sri4 := model.DummySalesReturnItem()
	sri4.Quantity = float32(50)
	sri4.SalesReturn = sr4
	sri4.Save()

	// test
	total, e := GetTotalStockReturn(int(time.Now().Month()), time.Now().Year())
	assert.NoError(t, e)
	assert.Equal(t, float32(150), total)
}

//// GetReportPerDay ///////////////////////////////////////////////////////////////////////////////

// TestGetReportPerDay test func get graphic sales order per hari
func TestGetReportPerDay(t *testing.T) {
	tm := time.Date(2012, time.November, 12, 15, 00, 00, 00, time.Local)
	tm2 := time.Date(2010, time.January, 12, 15, 00, 00, 00, time.Local)
	tm3 := time.Date(2010, time.January, 20, 07, 00, 00, 00, time.Local)
	tm4 := time.Date(2010, time.January, 25, 17, 00, 00, 00, time.Local)
	// hapus database
	o := orm.NewOrm()
	o.Raw("DELETE FROM sales_order").Exec()
	o.Raw("DELETE FROM finance_revenue").Exec()

	// buat dummy
	so1 := model.DummySalesOrder() //ok
	so1.IsDeleted = int8(0)
	so1.DocumentStatus = "active"
	so1.RecognitionDate = tm2
	so1.Save()
	soi1 := model.DummySalesOrderItem()
	soi1.Quantity = float32(100)
	soi1.SalesOrder = so1
	soi1.Save()
	soi12 := model.DummySalesOrderItem()
	soi12.Quantity = float32(50)
	soi12.SalesOrder = so1
	soi12.Save()
	sinv1 := model.DummySalesInvoice()
	sinv1.IsDeleted = int8(0)
	sinv1.SalesOrder = so1
	sinv1.Save()
	rev1 := model.DummyFinanceRevenue()
	rev1.RecognitionDate = tm2
	rev1.IsDeleted = int8(0)
	rev1.Amount = float64(20000)
	rev1.RefType = "sales_invoice"
	rev1.Save()
	rev12 := model.DummyFinanceRevenue()
	rev12.RecognitionDate = tm2
	rev12.IsDeleted = int8(0)
	rev12.Amount = float64(3000)
	rev12.RefType = "sales_invoice"
	rev12.Save()

	so2 := model.DummySalesOrder() // deleted
	so2.IsDeleted = int8(1)
	so2.DocumentStatus = "active"
	so2.RecognitionDate = tm3
	so2.Save()
	soi2 := model.DummySalesOrderItem()
	soi2.Quantity = float32(150)
	soi2.SalesOrder = so2
	soi2.Save()
	sinv2 := model.DummySalesInvoice()
	sinv2.IsDeleted = int8(0)
	sinv2.SalesOrder = so2
	sinv2.Save()
	rev2 := model.DummyFinanceRevenue()
	rev2.RecognitionDate = tm3
	rev2.IsDeleted = int8(1)
	rev2.Amount = float64(10000)
	rev2.RefType = "sales_invoice"
	rev2.Save()

	so3 := model.DummySalesOrder() // ok
	so3.IsDeleted = int8(0)
	so3.DocumentStatus = "active"
	so3.RecognitionDate = tm3
	so3.Save()
	soi3 := model.DummySalesOrderItem()
	soi3.Quantity = float32(200)
	soi3.SalesOrder = so3
	soi3.Save()
	soi32 := model.DummySalesOrderItem()
	soi32.Quantity = float32(200)
	soi32.SalesOrder = so3
	soi32.Save()
	sinv3 := model.DummySalesInvoice()
	sinv3.IsDeleted = int8(0)
	sinv3.SalesOrder = so3
	sinv3.Save()
	rev3 := model.DummyFinanceRevenue()
	rev3.RecognitionDate = tm3
	rev3.IsDeleted = int8(0)
	rev3.Amount = float64(5000)
	rev3.RefType = "sales_invoice"
	rev3.Save()

	so4 := model.DummySalesOrder() // ok
	so4.IsDeleted = int8(0)
	so4.DocumentStatus = "active"
	so4.RecognitionDate = tm4
	so4.Save()
	soi4 := model.DummySalesOrderItem()
	soi4.Quantity = float32(10)
	soi4.SalesOrder = so4
	soi4.Save()
	sinv4 := model.DummySalesInvoice()
	sinv4.IsDeleted = int8(0)
	sinv4.SalesOrder = so4
	sinv4.Save()
	rev4 := model.DummyFinanceRevenue()
	rev4.RecognitionDate = tm4
	rev4.IsDeleted = int8(0)
	rev4.Amount = float64(12000)
	rev4.RefType = "sales_invoice"
	rev4.Save()
	sinv42 := model.DummySalesInvoice()
	sinv42.IsDeleted = int8(1)
	sinv42.SalesOrder = so4
	sinv42.Save()
	rev42 := model.DummyFinanceRevenue()
	rev42.RecognitionDate = tm4
	rev42.IsDeleted = int8(1)
	rev42.Amount = float64(2000)
	rev42.RefType = "sales_invoice"
	rev42.Save()

	so5 := model.DummySalesOrder() // passed the date
	so5.IsDeleted = int8(0)
	so5.DocumentStatus = "active"
	so5.RecognitionDate = tm
	so5.Save()
	soi5 := model.DummySalesOrderItem()
	soi5.Quantity = float32(200)
	soi5.SalesOrder = so5
	soi5.Save()

	// test
	grp, e := GetReportPerDay(int(01), int(2010))
	assert.NoError(t, e)
	for i, r := range grp {
		if i == int(0) {
			assert.Equal(t, float64(23000), r.TotalRevenue)
			assert.Equal(t, int64(1), r.Transaction)
			assert.Equal(t, float32(150), r.ItemSold)
		}
		if i == int(1) {
			assert.Equal(t, float64(5000), r.TotalRevenue)
			assert.Equal(t, int64(1), r.Transaction)
			assert.Equal(t, float32(400), r.ItemSold)
		}
		if i == int(2) {
			assert.Equal(t, float64(12000), r.TotalRevenue)
			assert.Equal(t, int64(1), r.Transaction)
			assert.Equal(t, float32(10), r.ItemSold)
		}
	}
}

//// GetGraphicTopRevenue ///////////////////////////////////////////////////////////////////////////////

// TestGetGraphicTopRevenue test func get graphic top 10 revenue
func TestGetGraphicTopRevenue(t *testing.T) {
	// hapus database
	o := orm.NewOrm()
	o.Raw("DELETE FROM item").Exec()
	o.Raw("DELETE FROM sales_order").Exec()
	// buat dummy
	soi := model.DummySalesOrderItem() // itm beda tapi 1 so
	soi.Quantity = float32(10)
	soi.UnitPrice = float64(1500)
	soi.Discount = float32(0)
	soi.Subtotal = float64(15000)
	soi.Save()
	soi2 := model.DummySalesOrderItem()
	soi2.SalesOrder = soi.SalesOrder
	soi2.Quantity = float32(10)
	soi2.UnitPrice = float64(2000)
	soi2.Discount = float32(0)
	soi2.Subtotal = float64(20000)
	soi2.Save()
	soi22 := model.DummySalesOrderItem()
	soi22.SalesOrder = soi.SalesOrder
	soi22.ItemVariant = soi.ItemVariant
	soi22.Quantity = float32(10)
	soi22.UnitPrice = float64(2000)
	soi22.Discount = float32(0)
	soi22.Subtotal = float64(20000)
	soi22.Save()

	so := soi.SalesOrder
	so.IsDeleted = int8(0)
	so.DocumentStatus = "active"
	so.RecognitionDate = time.Now()
	so.Save()

	itm := soi.ItemVariant
	itm.IsDeleted = int8(0)
	itm.Save()
	itm2 := soi2.ItemVariant
	itm2.IsDeleted = int8(0)
	itm2.Save()

	so1 := model.DummySalesOrder() // 1 itm 1 soi
	so1.IsDeleted = int8(0)
	so1.DocumentStatus = "active"
	so.RecognitionDate = time.Now()
	so1.Save()
	soi3 := model.DummySalesOrderItem()
	soi3.SalesOrder = so1
	soi3.Quantity = float32(10)
	soi3.UnitPrice = float64(500)
	soi3.Discount = float32(0)
	soi3.Subtotal = float64(5000)
	soi3.Save()

	itm3 := soi3.ItemVariant
	itm3.IsDeleted = int8(0)
	itm3.Save()

	// test
	top, e := GetGraphicTopRevenue(int(time.Now().Month()), time.Now().Year())
	assert.NoError(t, e)
	for i, u := range top {
		if i == int(0) {
			assert.Equal(t, float64(35000), u.TotalSubTotal)
			assert.Equal(t, itm.ID, u.ItemVariant.ID)
			assert.NotEmpty(t, u.ItemVariant.Item)
		}
		if i == int(1) {
			assert.Equal(t, float64(20000), u.TotalSubTotal)
			assert.Equal(t, itm2.ID, u.ItemVariant.ID)
			assert.NotEmpty(t, u.ItemVariant.Item)
		}
		if i == int(2) {
			assert.Equal(t, float64(5000), u.TotalSubTotal)
			assert.Equal(t, itm3.ID, u.ItemVariant.ID)
			assert.NotEmpty(t, u.ItemVariant.Item)
		}
	}
}

//// GetTodoList ///////////////////////////////////////////////////////////////////////////////

// TestGetTodoList test func get to do list
func TestGetTodoList(t *testing.T) {
	// hapus database
	o := orm.NewOrm()
	o.Raw("DELETE FROM item").Exec()
	o.Raw("DELETE FROM sales_order").Exec()
	o.Raw("DELETE FROM purchase_order").Exec()
	o.Raw("DELETE FROM workorder_shipment").Exec()
	// buat dummy
	// sales fulfill-->2
	// sales inv-->2
	// sales cancel-->1
	// fulfill ship-->0
	// shipment-->1
	// need recaive-->30
	sales1 := model.DummySalesOrder() // so
	sales1.IsDeleted = int8(0)
	sales1.DocumentStatus = "active"
	sales1.FulfillmentStatus = "active"
	sales1.InvoiceStatus = "active"
	sales1.Save()

	sales2 := model.DummySalesOrder()
	sales2.IsDeleted = int8(0)
	sales2.DocumentStatus = "requested_cancel"
	sales2.FulfillmentStatus = "finished"
	sales2.InvoiceStatus = "active"
	sales2.Save()

	sales3 := model.DummySalesOrder()
	sales3.IsDeleted = int8(0)
	sales3.DocumentStatus = "finished"
	sales3.FulfillmentStatus = "active"
	sales3.InvoiceStatus = "finished"
	sales3.Save()

	fulfill := model.DummyWorkorderFulfillment() // fullfill gagal
	delSo := fulfill.SalesOrder
	fulfill.IsDeleted = int8(0)
	fulfill.IsDelivered = int8(1)
	fulfill.SalesOrder = sales1
	fulfill.Save()

	delSo.Delete()

	ship := model.DummyWorkorderShipment() // shipment sukses
	ship.IsDeleted = int8(0)
	ship.DocumentStatus = "active"
	ship.Save()
	ship2 := model.DummyWorkorderShipment() // shipment gagal
	ship2.IsDeleted = int8(0)
	ship2.DocumentStatus = "finished"
	ship2.Save()

	po := model.DummyPurchaseOrder() // po sukses
	po.Note = "po-1"
	po.DocumentStatus = "active"
	po.IsDeleted = int8(0)
	po.ReceivingStatus = "new"
	po.Save()
	poi := model.DummyPurchaseOrderItem()
	delPo := poi.PurchaseOrder
	poi.PurchaseOrder = po
	poi.Quantity = float32(25)
	poi.Note = "note-1"
	poi.Save()
	delPo.Delete()
	poi2 := model.DummyPurchaseOrderItem()
	delPo2 := poi2.PurchaseOrder
	poi2.PurchaseOrder = po
	poi2.Quantity = float32(20)
	poi2.Note = "note-2"
	poi2.Save()
	delPo2.Delete()

	po2 := model.DummyPurchaseOrder() //po gagal
	po2.DocumentStatus = "cancelled"
	po2.IsDeleted = int8(0)
	po2.ReceivingStatus = "finished"
	po2.Note = "note-del-1"
	po2.Save()
	poi3 := model.DummyPurchaseOrderItem()
	delPo3 := poi3.PurchaseOrder
	poi3.PurchaseOrder = po2
	poi3.Quantity = float32(20)
	poi3.Note = "note-del-2"
	poi3.Save()
	delPo3.Delete()

	wr := model.DummyWorkorderReceiving() //ok
	wr.PurchaseOrder = po
	wr.IsDeleted = int8(0)
	wr.DocumentStatus = "finished"
	wr.Note = "note-1"
	wr.Save()
	wrItem := model.DummyWorkorderReceivingItem()
	delwrItem := wrItem.PurchaseOrderItem
	delwrItem2 := wrItem.WorkorderReceiving
	wrItem.PurchaseOrderItem = poi
	wrItem.Quantity = float32(10)
	wrItem.WorkorderReceiving = wr
	wrItem.Save()
	delwrItem.PurchaseOrder.Delete()
	delwrItem2.Delete()
	wrItem2 := model.DummyWorkorderReceivingItem()
	delwrItem3 := wrItem2.PurchaseOrderItem
	delwrItem4 := wrItem2.WorkorderReceiving
	wrItem2.PurchaseOrderItem = poi2
	wrItem2.Quantity = float32(5)
	wrItem2.WorkorderReceiving = wr
	wrItem2.Save()
	delwrItem3.PurchaseOrder.Delete()
	delwrItem4.Delete()

	// test
	m, e := GetTodoList()
	assert.NoError(t, e)
	assert.Equal(t, m.SalesNeedFulfilled, int64(2))
	assert.Equal(t, m.SalesNeedInvoiced, int64(2))
	assert.Equal(t, m.SalesNeedApproveCancel, int64(1))
	assert.Equal(t, m.FulfillmentNeedShipped, int64(0))
	assert.Equal(t, m.ShipmentNeedDelivered, int64(1))
	assert.Equal(t, m.ItemNeedReceived, float32(30))

}

//// GetOverview ///////////////////////////////////////////////////////////////////////////////

func TestGetOverview(t *testing.T) {
	sales := SOAndPOData{
		Charge: float64(20000),
		Cost:   float64(10000),
		Paid:   float64(15000),
	}
	purchase := SOAndPOData{
		Charge: float64(10000),
		Paid:   float64(2000),
	}
	m := GetOverview(sales, purchase)
	assert.Equal(t, float64(20000), m.TotalPenjualan)
	assert.Equal(t, float64(10000), m.TotalCost)
	assert.Equal(t, float64(10000), m.TotalProfit)
	assert.Equal(t, float64(5000), m.PiutangUsaha)
	assert.Equal(t, float64(8000), m.HutangUsaha)
}

//// GetNeedToRestock /////////////////////////////////////////////////////////////////////////////

func TestGetNeedToRestock(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant").Exec()
	// buat dummy item
	itm1 := model.DummyItemVariant()
	itm1.AvailableStock = float32(11)
	itm1.MinimumStock = float32(10)
	itm1.IsDeleted = int8(0)
	itm1.Save()

	itm2 := model.DummyItemVariant()
	itm2.AvailableStock = float32(6)
	itm2.MinimumStock = float32(6)
	itm2.IsDeleted = int8(0)
	itm2.Save()

	itm3 := model.DummyItemVariant()
	itm3.AvailableStock = float32(4)
	itm3.MinimumStock = float32(5)
	itm3.IsDeleted = int8(0)
	itm3.Save()

	itm4 := model.DummyItemVariant()
	itm4.AvailableStock = float32(0)
	itm4.MinimumStock = float32(0)
	itm4.IsDeleted = int8(0)
	itm4.Save()

	// test
	m, e := GetNeedToRestock()
	assert.NoError(t, e)
	for _, j := range m {
		assert.NotEmpty(t, j.Item)
	}
}

func TestGetNeedToRestockEmpty(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM item_variant").Exec()

	// test
	m, e := GetNeedToRestock()
	assert.NoError(t, e)
	assert.Empty(t, m)
}
