// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchaseInvoice

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for purchaseInvoice.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.POST("", h.create, auth.CheckPrivilege("purchase_invoice_create"))
	r.GET("", h.get, auth.CheckPrivilege("purchase_invoice_read"))
	r.GET("/:id", h.show, auth.CheckPrivilege("purchase_invoice_show"))
	r.PUT("/:id", h.put, auth.CheckPrivilege("purchase_invoice_update"))
}

// create endpoint to handle put http method with id.
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest
	var sd *auth.SessionData
	var cpi *model.PurchaseInvoice

	if sd, e = auth.UserSession(ctx); e == nil {
		r.SessionData = sd
		if e = ctx.Bind(&r); e == nil {
			if cpi, e = CreatePurchaseInvoice(&r); e == nil {
				ctx.Data(cpi)
			}
		}
	}
	return ctx.Serve(e)
}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	var data *[]model.PurchaseInvoice

	var tot int64

	if data, tot, e = GetPurchaseInvoice(rq); e == nil {
		ctx.Data(data, tot)
	}

	return ctx.Serve(e)
}

// show endpoint to handle get http method with id.
func (h *Handler) show(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var as *model.PurchaseInvoice
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if as, e = ShowPurchaseInvoice("id", id); e == nil {
			ctx.Data(as)
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// put endpoint to handle put http method with id.
func (h *Handler) put(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest
	var sd *auth.SessionData
	var id int64
	var pi *model.PurchaseInvoice
	var pix *model.PurchaseInvoice
	if sd, e = auth.UserSession(ctx); e == nil {
		r.SessionData = sd
		if id, e = common.Decrypt(ctx.Param("id")); e == nil {
			if pi, e = ShowPurchaseInvoice("id", id); e == nil {
				r.PI = pi
				if e = ctx.Bind(&r); e == nil {
					if pix, e = UpdatePurchaseInvoice(&r); e == nil {
						ctx.Data(pix)
					}
				}
			} else {
				e = echo.ErrNotFound
			}
		}
	}

	return ctx.Serve(e)
}
