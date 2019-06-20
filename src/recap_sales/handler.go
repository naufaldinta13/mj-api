// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package recapSales

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for recapSales.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.CheckPrivilege("report_customer_sales"))
	r.POST("", h.create, auth.CheckPrivilege("report_customer_sales"))
	r.GET("/:id", h.show, auth.CheckPrivilege("report_customer_sales"))
	r.PUT("/:id/cancel", h.cancel, auth.CheckPrivilege("report_customer_sales"))
}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	var t int64
	var m *[]model.RecapSales
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	if m, t, e = GetRecapSales(rq); e == nil {
		ctx.Data(m, t)
	}
	return ctx.Serve(e)
}

// create endpoint to handle post http method.
func (h *Handler) create(c echo.Context) (e error) {
	var r createRequest
	var session *auth.SessionData
	var m model.RecapSales
	ctx := c.(*cuxs.Context)
	if session, e = auth.UserSession(ctx); e == nil {
		r.Session = session
		if e = ctx.Bind(&r); e == nil {
			if m, e = r.Transform(); e == nil {
				if e = CreateRecapSales(&m); e == nil {
					ctx.Data(m)
				}
			}
		}
	}

	return ctx.Serve(e)
}

// show endpoint to handle get http method with id.
func (h *Handler) show(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var m *model.RecapSales
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if m, e = ShowRecapSales("id", id); e == nil {
			ctx.Data(m)
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// cancel endpoint to handle get http method.
func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRequest
	var id int64

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if r.RecapSales, e = ShowRecapSales("id", id); e == nil {
			if e = ctx.Bind(&r); e == nil {
				rs := r.Transform()
				if e = rs.Save("IsDeleted"); e == nil {
					ctx.Data(rs)
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}
