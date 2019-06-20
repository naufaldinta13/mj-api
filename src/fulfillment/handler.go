// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fulfillment

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for fullfillment.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.CheckPrivilege("fulfillment_read"))
	r.GET("/:id", h.show, auth.CheckPrivilege("fulfillment_show"))
	r.POST("", h.create, auth.CheckPrivilege("fulfillment_create"))
	r.PUT("/:id", h.update, auth.CheckPrivilege("fulfillment_update"))
	r.PUT("/:id/approve", h.approve, auth.CheckPrivilege("fulfillment_approve"))
}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	var tot int64
	var data *[]model.WorkorderFulfillment
	if data, tot, e = getWorkorderFulfillments(rq); e == nil {
		ctx.Data(data, tot)
	}

	return ctx.Serve(e)
}

// show endpoint to handle get http method with id.
func (h *Handler) show(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var data *model.WorkorderFulfillment
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if data, e = getWorkorderFulfillmentByID(id); e == nil {
			ctx.Data(data)
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// create endpoint to handle get http method.
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest
	var sd *auth.SessionData
	var data *model.WorkorderFulfillment

	if e = ctx.Bind(&r); e == nil {
		if sd, e = auth.UserSession(ctx); e == nil {
			fulfillment := r.Transform(sd.User)
			if data, e = saveDataFulFillment(fulfillment); e == nil {
				ctx.Data(data)
			}
		}
	}

	return ctx.Serve(e)
}

// update endpoint to handle http put method
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest
	var id int64
	var sd *auth.SessionData

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if r.Fulfillment, e = getWorkorderFulfillmentByID(id); e == nil {
			if e = ctx.Bind(&r); e == nil {
				if sd, e = auth.UserSession(ctx); e == nil {
					fulfillment, itemsfulReq := r.Transform(sd.User)
					if fulfillment, e = updateDataFulFillment(fulfillment, itemsfulReq); e == nil {
						ctx.Data(fulfillment)
					}
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// update endpoint to handle http put method
func (h *Handler) approve(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r approveRequest
	var id int64
	var data *model.WorkorderFulfillment

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if r.Fulfillment, e = getWorkorderFulfillmentByID(id); e == nil {
			if e = ctx.Bind(&r); e == nil {
				if data, e = approveFulfillment(r.Fulfillment); e == nil {
					ctx.Data(data)
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}
