// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchase

import (
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
	r.GET("", h.get, auth.CheckPrivilege("purchase_order_read"))
	r.GET("/:id", h.show, auth.CheckPrivilege("purchase_order_show"))
	r.POST("", h.create, auth.CheckPrivilege("purchase_order_create"))
	r.PUT("/:id", h.update, auth.CheckPrivilege("purchase_order_update"))
	r.PUT("/:id/cancel", h.cancel, auth.CheckPrivilege("purchase_order_cancel"))
}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	var t int64
	var m *[]model.PurchaseOrder
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	if m, t, e = GetAllPurchaseOrder(rq); e == nil {
		ctx.Data(m, t)
	}
	return ctx.Serve(e)
}

// show endpoint to get detail sales_order
func (h *Handler) show(c echo.Context) (e error) {
	var id int64
	var m *model.PurchaseOrder
	ctx := c.(*cuxs.Context)
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if m, e = GetDetailPurchaseOrder("id", id); e == nil {
			ctx.Data(m)
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// get endpoint to handle get http method.
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest
	var sd *auth.SessionData
	var data *model.PurchaseOrder

	if e = ctx.Bind(&r); e == nil {
		if sd, e = auth.UserSession(ctx); e == nil {
			po := r.Transform(sd.User)
			if data, e = CreatePurchaseOrder(po); e == nil {
				ctx.Data(data)
			}
		}
	}
	return ctx.Serve(e)
}

// get endpoint to handle get http method.
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest
	var sd *auth.SessionData
	var id int64
	var data *model.PurchaseOrder

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if r.PurchaseOrder, e = GetDetailPurchaseOrder("id", id); e == nil {
			if e = ctx.Bind(&r); e == nil {
				if sd, e = auth.UserSession(ctx); e == nil {
					po, items := r.Transform(sd.User)
					if data, e = UpdatePurchaseOrder(po, items); e == nil {
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

// cancel endpoint to handle get http method.
func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRequest
	var id int64
	var data *model.PurchaseOrder

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if r.PurchaseOrder, e = GetDetailPurchaseOrder("id", id); e == nil {
			if e = ctx.Bind(&r); e == nil {
				po := r.Transform()
				if data, e = cancelPurchaseOrder(po); e == nil {
					ctx.Data(data)
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}
