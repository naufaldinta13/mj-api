// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package report

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/purchase"
	"git.qasico.com/mj/api/src/sales"

	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
	"git.qasico.com/mj/api/src/finance_revenue"
	"fmt"
)

// Handler collection handler for report.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/sales-item", h.salesItem, auth.CheckPrivilege("report_sales_item"))
	r.GET("/sales", h.sales, auth.CheckPrivilege("report_sales"))
	r.GET("/purchase", h.purchase, auth.CheckPrivilege("report_purchase"))
	r.GET("/purchase-item", h.purchaseItem, auth.CheckPrivilege("report_purchase_item"))
	r.GET("/bank", h.bank, auth.CheckPrivilege("report_bank"))
	r.GET("/sales-item/summary", h.salesItemSummary, auth.CheckPrivilege("report_sales_item"))
	r.GET("/sales/summary", h.salesSummary, auth.CheckPrivilege("report_sales"))
	r.GET("/purchase/summary", h.purchaseSummary, auth.CheckPrivilege("report_purchase"))
	r.GET("/purchase-item/summary", h.purchaseItemSummary, auth.CheckPrivilege("report_purchase_item"))
	r.GET("/bank/summary", h.bankSummary, auth.CheckPrivilege("report_bank"))
}

// salesItem endpoint to handle get http method.
func (h *Handler) salesItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var data []*model.SalesOrderItem
	var t int64

	// get query string from request
	rq := ctx.RequestQuery()
	if data, t, e = GetReportSalesItems(rq); e == nil {
		ctx.Data(data, t)
	}

	return ctx.Serve(e)
}

// sales endpoint to handle get http method.
func (h *Handler) sales(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var data *[]model.SalesOrder
	var t int64

	// get query string from request
	rq := ctx.RequestQuery()
	if data, t, e = sales.GetAllSalesOrder(rq); e == nil {
		ctx.Data(data, t)
	}
	fmt.Println("errror", e)

	return ctx.Serve(e)
}

// purchase endpoint to handle get http method.
func (h *Handler) purchase(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var data *[]model.PurchaseOrder
	var t int64
	// get query string from request
	rq := ctx.RequestQuery()
	if data, t, e = purchase.GetAllPurchaseOrder(rq); e == nil {
		ctx.Data(data, t)
	}

	return ctx.Serve(e)
}

// purchaseItem endpoint to handle get http method.
func (h *Handler) purchaseItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var data []*model.PurchaseOrderItem
	var t int64
	// get query string from request
	rq := ctx.RequestQuery()
	if data, t, e = GetReportPurchaseItems(rq); e == nil {
		ctx.Data(data, t)
	}

	return ctx.Serve(e)
}

// bank endpoint to handle get http method.
func (h *Handler) bank(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var data *[]model.FinanceRevenue
	var t int64
	// get query string from request
	rq := ctx.RequestQuery()
	if data, t, e = financeRevenue.GetFinanceRevenues(rq); e == nil {
		ctx.Data(data, t)
	}

	return ctx.Serve(e)
}

// salesItem endpoint to handle get http method.
func (h *Handler) salesItemSummary(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var data *SalesItemData
	var category, ivar, startDate, endDate string

	param, _ := ctx.FormParams()

	category = param.Get("category")
	ivar = param.Get("item_variant_id")
	startDate = param.Get("start_date")
	endDate = param.Get("end_date")

	if data, e = GetTotalSalesItem(category, ivar, startDate, endDate); e == nil {
		ctx.Data(data)
	}

	return ctx.Serve(e)
}

// sales endpoint to handle get http method.
func (h *Handler) salesSummary(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var data *SalesData
	var custID, startDate, endDate string

	param, _ := ctx.FormParams()

	startDate = param.Get("start_date")
	endDate = param.Get("end_date")
	custID = param.Get("customer_id")

	if data, e = GetTotalSales(custID, startDate, endDate); e == nil {
		ctx.Data(data)
	}

	return ctx.Serve(e)
}

// purchase endpoint to handle get http method.
func (h *Handler) purchaseSummary(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var data *PurchaseData
	var suppID, startDate, endDate string

	param, _ := ctx.FormParams()

	startDate = param.Get("start_date")
	endDate = param.Get("end_date")
	suppID = param.Get("supplier_id")

	if data, e = GetTotalPurchase(suppID, startDate, endDate); e == nil {
		ctx.Data(data)
	}

	return ctx.Serve(e)
}

// purchaseItem endpoint to handle get http method.
func (h *Handler) purchaseItemSummary(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var data *PurchaseItemData
	var category, ivar string

	param, _ := ctx.FormParams()

	category = param.Get("category")
	ivar = param.Get("item_variant_id")

	if data, e = GetTotalPurchaseItem(category, ivar); e == nil {
		ctx.Data(data)
	}

	return ctx.Serve(e)
}

// bank endpoint to handle get http method.
func (h *Handler) bankSummary(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var data *RevenueData
	var custID, bankID string

	param, _ := ctx.FormParams()

	custID = param.Get("customer_id")
	bankID = param.Get("bank_id")

	if data, e = GetTotalRevenueBank(custID, bankID); e == nil {
		ctx.Data(data)
	}

	return ctx.Serve(e)
}
