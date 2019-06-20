// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchaseReturn

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/cuxs/validation"
	"github.com/labstack/echo"
)

// Handler collection handler for purchaseReturn.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.POST("", h.create, auth.CheckPrivilege("purchase_return_create"))
	r.GET("", h.get, auth.CheckPrivilege("purchase_return_read"))
	r.GET("/:id", h.show, auth.CheckPrivilege("purchase_return_show"))
	r.PUT("/:id", h.update, auth.CheckPrivilege("purchase_return_update"))
	r.PUT("/:id/cancel", h.cancel, auth.CheckPrivilege("purchase_return_cancel"))
}

// create endpoint to handle post http method.
func (h *Handler) create(c echo.Context) (e error) {
	var r createRequest
	var session *auth.SessionData
	var preturn *model.PurchaseReturn
	ctx := c.(*cuxs.Context)
	if session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			r.Session = session
			if preturn, e = CreatePurchaseReturn(r.Transform()); e == nil {
				ctx.Data(preturn)
			}
		}
	}
	return ctx.Serve(e)
}

// get endpoint to handle get http method
// untuk menampilkan semua data
func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var total int64
	var preturn []*model.PurchaseReturn
	var rq *orm.RequestQuery
	rq = ctx.RequestQuery()

	if preturn, total, e = GetAllPurchaseReturn(rq); e == nil {
		ctx.Data(preturn, total)
	}

	return ctx.Serve(e)
}

// show endpoint to handle get http method
// untuk menampilkan detail data
func (h *Handler) show(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64
	var preturn *model.PurchaseReturn

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if preturn, e = GetDetailPurchaseReturn("id", id); e == nil {
			ctx.Data(preturn)
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// update endpoint to handle put http method
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64
	var r updateRequest
	var preturn *model.PurchaseReturn

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if r.Session, e = auth.UserSession(ctx); e == nil {
			if r.PurchaseReturn, e = GetDetailPurchaseReturn("id", id); e == nil {
				if e = ctx.Bind(&r); e == nil {
					preturn, e = updatePurchaseReturn(r.Transform())
					if e == nil {
						ctx.Data(preturn)
					}
				}
			} else {
				e = echo.ErrNotFound
			}
		}
	}

	return ctx.Serve(e)
}

// cancel endpoint to handle put http method
func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64
	var preturn, purchase *model.PurchaseReturn

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if preturn, e = GetDetailPurchaseReturn("id", id); e == nil {
			if preturn.DocumentStatus != "cancelled" {
				if purchase, e = cancelPurchaseReturn(preturn); e == nil {
					ctx.Data(purchase)
				}
			} else {
				e = validation.SetError("document_status", "Document status is already canceled")
			}
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}
