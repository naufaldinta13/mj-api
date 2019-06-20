// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stockopname

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for stockopname.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.CheckPrivilege("stock_opname_read"))
	r.GET("/:id", h.show, auth.CheckPrivilege("stock_opname_show"))
	r.POST("", h.create, auth.CheckPrivilege("stock_opname_create"))
}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	var tot int64
	var data *[]model.Stockopname
	if data, tot, e = GetStockOpnames(rq); e == nil {
		ctx.Data(data, tot)
	}

	return ctx.Serve(e)
}

// show endpoint to handle get http method with id.
func (h *Handler) show(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var data *model.Stockopname
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if data, e = GetStockopnameByID(id); e == nil {
			ctx.Data(data)
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// create endpoint to handle post http method with id.
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest
	var sd *auth.SessionData

	if e = ctx.Bind(&r); e == nil {
		if sd, e = auth.UserSession(ctx); e == nil {
			stockopname := r.Transform(sd.User)
			if data, e := saveDataStockopname(stockopname); e == nil {
				ctx.Data(data)
			}
		}
	}
	return ctx.Serve(e)
}
