// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package salesInvoice

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for salesInvoice.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.CheckPrivilege("sales_invoice_read"))
	r.POST("", h.create, auth.CheckPrivilege("sales_invoice_create"))
	r.GET("/:id", h.show, auth.CheckPrivilege("sales_invoice_show"))
	r.PUT("/:id", h.update, auth.CheckPrivilege("sales_invoice_update"))
}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var t int64
	var m *[]model.SalesInvoice
	// get query string from request
	rq := ctx.RequestQuery()
	if m, t, e = GetSalesInvoice(rq); e == nil {
		ctx.Data(m, t)
	}

	return ctx.Serve(e)
}

// create endpoint to handle post http method.
func (h *Handler) create(c echo.Context) (e error) {
	var r createRequest
	var session *auth.SessionData
	var u *model.SalesInvoice
	ctx := c.(*cuxs.Context)
	if session, e = auth.UserSession(ctx); e == nil {
		r.Session = session
		if e = ctx.Bind(&r); e == nil {
			if u, e = r.Transform(); e == nil {
				if e = CreateSalesInvoice(u); e == nil {
					ctx.Data(u)
				}
			}
		}
	}
	return ctx.Serve(e)
}

// show endpoint to handle get http method with id.
func (h *Handler) show(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var m *model.SalesInvoice
	var id int64
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if m, e = ShowSalesInvoice("id", id); e == nil {
			ctx.Data(m)
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// update endpoint to handle put http method.
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest
	var id int64
	var m *model.SalesInvoice
	var session *auth.SessionData
	if session, e = auth.UserSession(ctx); e == nil {
		if id, e = common.Decrypt(ctx.Param("id")); e == nil {
			if m, e = ShowSalesInvoice("id", id); e == nil {
				r.SalesInvoiceOld = m
				r.Session = session
				if e = ctx.Bind(&r); e == nil {
					r.Transform()
					if e = m.Save("recognition_date", "note", "total_amount", "updated_by", "updated_at"); e == nil {
						ctx.Data(m)
					}
				}
			} else {
				e = echo.ErrNotFound
			}
		}

	}

	return ctx.Serve(e)
}
