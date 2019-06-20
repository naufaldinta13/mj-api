// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sales

import (
	"strings"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for brand.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.CheckPrivilege("sales_order_read"))
	r.GET("/old", h.getOld, auth.CheckPrivilege("sales_order_read"))
	r.GET("/:id", h.show, auth.CheckPrivilege("sales_order_show"))
	r.POST("", h.create, auth.CheckPrivilege("sales_order_create"))
	r.PUT("/:id", h.update, auth.CheckPrivilege("sales_order_update"))
	r.PUT("/:id/cancel", h.cancel, auth.CheckPrivilege("sales_order_cancel"))
	r.PUT("/:id/cancel/request", h.cancelReq, auth.CheckPrivilege("sales_order_request_cancel"))
	r.PUT("/:id/cancel/reject", h.rejectCancelReq, auth.CheckPrivilege("sales_order_reject_cancel"))
	r.GET("/fulfillment/:id", h.fulfillment, auth.CheckPrivilege("sales_order_show"))
}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	var t int64
	var m *[]model.SalesOrder
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	if m, t, e = GetAllSalesOrder(rq); e == nil {
		ctx.Data(m, t)
	}
	return ctx.Serve(e)
}

// getOld endpoint to handle get http method.
func (h *Handler) getOld(c echo.Context) (e error) {
	var t int64
	var m *[]model.SalesOrder
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	if m, t, e = GetAllSalesOrderOld(rq); e == nil {
		ctx.Data(m, t)
	}
	return ctx.Serve(e)
}

// show endpoint to get detail sales_order
func (h *Handler) show(c echo.Context) (e error) {
	var id int64
	var m *model.SalesOrder
	ctx := c.(*cuxs.Context)
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		loadRelated := []string{}

		if p, e := ctx.FormParams(); e == nil {
			if len(p) > int(0) {
				loadParam := p.Get("load")
				if loadParam != "" {
					loadRelated = strings.Split(loadParam, ",")
				}
			}
		}

		if m, e = GetDetailSalesOrder(id, loadRelated); e == nil {
			ctx.Data(m)
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// create endpoint to create sales_order
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest
	var session *auth.SessionData
	var sorder *model.SalesOrder

	if session, e = auth.UserSession(ctx); e == nil {
		r.Session = session
		e = ctx.Bind(&r)
		if e == nil {
			// Transform dilakukan didalam CreateSalesOrder yang akan menghasilkan
			// SalesOrder yang sudah disimpan dan pesan error jika terjadi error
			sorder, e = CreateSalesOrder(&r)
			if e == nil {
				ctx.Data(sorder)
			}
		}
	}

	return ctx.Serve(e)
}

// update endpoint to update sales_order
func (h *Handler) update(c echo.Context) (e error) {
	var id int64
	ctx := c.(*cuxs.Context)
	var r updateRequest
	var sd *auth.SessionData
	var data *model.SalesOrder
	var emptyLoad []string
	emptyLoad = append(emptyLoad, "sales_order_items")

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if r.SalesOrder, e = GetDetailSalesOrder(id, emptyLoad); e == nil {
			if e = ctx.Bind(&r); e == nil {
				if sd, e = auth.UserSession(ctx); e == nil {
					so, items := r.Transform(sd.User)
					if data, e = UpdateSalesOrder(so, items); e == nil {
						ctx.Data(data)
					}
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// cancel endpoint to handle put http method
func (h *Handler) cancel(c echo.Context) (e error) {

	ctx := c.(*cuxs.Context)
	var so *model.SalesOrder
	var r cancelRequest
	var id int64
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e == nil {
		if id, e = common.Decrypt(ctx.Param("id")); e == nil {
			if so, e = GetDetailSalesOrder(id, nil); e == nil {
				r.Sales = so
				if e = ctx.Bind(&r); e == nil {
					so.CancelledNote = r.CancelledNote
					if e = CancelSalesOrder(so, session.User); e == nil {
						ctx.Data(so)
					}
				}

			} else {
				e = echo.ErrNotFound
			}
		}
	}
	return ctx.Serve(e)
}

// cancelReq endpoint to cancel request sales
func (h *Handler) cancelReq(c echo.Context) (e error) {
	var id int64
	var session *auth.SessionData
	var r cancelReq
	var u *model.SalesOrder
	ctx := c.(*cuxs.Context)
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if u, e = GetDetailSalesOrder(id, nil); e == nil {
			if session, e = auth.UserSession(ctx); e == nil {
				r.Session = session
				r.Sales = u
				if e = ctx.Bind(&r); e == nil {
					r.Transform(u)
					if e = u.Save("cancelled_note", "document_status", "request_cancel_by", "request_cancel_at"); e == nil {
						ctx.Data(u)
					}
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// update endpoint to update partnership
func (h *Handler) rejectCancelReq(c echo.Context) (e error) {
	var id int64
	var r rejectCancelRequest
	var u *model.SalesOrder
	ctx := c.(*cuxs.Context)
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if u, e = GetDetailSalesOrder(id, nil); e == nil {
			r.SalesOrder = u
			if e = ctx.Bind(&r); e == nil {
				r.Transform()
				if e = u.Save("document_status"); e == nil {
					ctx.Data(u)
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// show endpoint to handle get http method with id.
func (h *Handler) fulfillment(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64
	var sales *model.SalesOrder

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if sales, e = ShowSalesOrderFulfillment("id", id); e == nil {
			ctx.Data(sales)
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}
