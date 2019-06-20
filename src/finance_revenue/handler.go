// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package financeRevenue

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/src/sales_invoice"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for financeRevenue.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.CheckPrivilege("revenue_read"))
	r.GET("/:id", h.show, auth.CheckPrivilege("revenue_show"))
	r.POST("", h.create, auth.CheckPrivilege("revenue_create"))
	r.PUT("/:id", h.update, auth.CheckPrivilege("revenue_update"))
	r.PUT("/:id/approve", h.approve, auth.CheckPrivilege("revenue_approve"))
	r.GET("/summary", h.summary, auth.CheckPrivilege("revenue_read"))
}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var total int64
	var data *[]model.FinanceRevenue
	if data, total, e = GetFinanceRevenues(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// show endpoint to get detail finance revenue
func (h *Handler) show(c echo.Context) (e error) {
	var id int64
	var m *model.FinanceRevenue
	ctx := c.(*cuxs.Context)
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if m, e = ShowFinanceRevenue("id", id); e == nil {
			ctx.Data(m)
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// create endpoint to handle post http method.
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var r createRequest
	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			fr := r.Transform()
			if e = fr.Save(); e == nil {
				if fr.RefType == "sales_invoice" {
					salesInvoice.SumTotalRevenuedSalesInvoice(int64(fr.RefID))
					finance, _ := ShowFinanceRevenue("id", fr.ID)
					if e = ApproveRevenue(finance); e == nil {
						ctx.Data(e)
					}
				}
				ctx.Data(e)
			}
		}
	}

	return ctx.Serve(e)
}

// update endpoint to handle put http method.
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest
	var id int64

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if r.Session, e = auth.UserSession(ctx); e == nil {
			if r.Revenue, e = ShowFinanceRevenue("id", id); e == nil {
				if e = ctx.Bind(&r); e == nil {
					fr := r.Transform()
					if e = fr.Save(); e == nil {
						if fr.RefType == "sales_invoice" {
							salesInvoice.SumTotalRevenuedSalesInvoice(int64(fr.RefID))
						}
						ctx.Data(e)
					}
				}
			} else {
				e = echo.ErrNotFound
			}
		}
	}

	return ctx.Serve(e)
}

// approve endpoint to handle put http method.
func (h *Handler) approve(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r approveRequest
	var id int64

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if r.Revenue, e = ShowFinanceRevenue("id", id); e == nil {
			if e = ctx.Bind(&r); e == nil {
				if e = ApproveRevenue(r.Revenue); e == nil {
					ctx.Data(r.Revenue)
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) summary(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	s := ctx.QueryParam("document_status")
	rt := ctx.QueryParam("ref_type")
	pm := ctx.QueryParam("payment_method")
	ds := ctx.QueryParam("date_start")
	de := ctx.QueryParam("date_end")

	ctx.Data(summaryRevenue(s, rt, pm, ds, de))

	return ctx.Serve(e)
}
