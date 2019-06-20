// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receiving

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for receiving.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.POST("", h.create, auth.CheckPrivilege("receiving_create"))
	r.GET("", h.get, auth.CheckPrivilege("receiving_read"))
	r.GET("/:id", h.show, auth.CheckPrivilege("receiving_show"))
}

// create endpoint to handle put http method with id.
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest
	var sd *auth.SessionData
	var cwr *model.WorkorderReceiving

	if sd, e = auth.UserSession(ctx); e == nil {
		r.SessionData = sd
		if e = ctx.Bind(&r); e == nil {
			if cwr, e = CreateReceiving(&r); e == nil {
				ctx.Data(cwr)
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
	var data *[]model.WorkorderReceiving

	var tot int64

	if data, tot, e = GetReceiving(rq); e == nil {
		ctx.Data(data, tot)
	}

	return ctx.Serve(e)
}

func (h *Handler) show(c echo.Context) (e error) {
	var id int64
	var receiving *model.WorkorderReceiving
	ctx := c.(*cuxs.Context)

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if receiving, e = GetDetailReceiving("id", id); e == nil {
			ctx.Data(receiving)
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}
