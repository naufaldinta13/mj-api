// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package measurement

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for measurement.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.CheckPrivilege("uom_read"))
	r.GET("/:id", h.show, auth.CheckPrivilege("uom_show"))
	r.POST("", h.create, auth.CheckPrivilege("uom_create"))
	r.PUT("/:id", h.put, auth.CheckPrivilege("uom_update"))
	r.DELETE("/:id", h.delete, auth.CheckPrivilege("uom_delete"))
}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	var data *[]model.Measurement

	var tot int64

	if data, tot, e = GetMeasurement(rq); e == nil {
		ctx.Data(data, tot)
	}

	return ctx.Serve(e)
}

// show endpoint to handle get http method with id.
func (h *Handler) show(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var as *model.Measurement
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if as, e = ShowMeasurement("id", id); e == nil {
			ctx.Data(as)
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// put endpoint to handle put http method with id.
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest
	var sd *auth.SessionData

	if sd, e = auth.UserSession(ctx); e == nil {

		r.SessionData = sd
		if e = ctx.Bind(&r); e == nil {
			u := r.Transform()

			if e = u.Save(); e == nil {
				ctx.Data(u)
			}
		}
	}

	return ctx.Serve(e)
}

// put endpoint to handle put http method with id.
func (h *Handler) put(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest
	var sd *auth.SessionData
	var m *model.Measurement

	if sd, e = auth.UserSession(ctx); e == nil {
		r.SessionData = sd

		if r.ID, e = common.Decrypt(ctx.Param("id")); e == nil {
			if m, e = ShowMeasurement("id", r.ID); e == nil {
				if e = ctx.Bind(&r); e == nil {
					u := r.Transform(m)

					if e = u.Save(); e == nil {
						ctx.Data(u)
					}
				}
			} else {
				e = echo.ErrNotFound
			}
		}
	}

	return ctx.Serve(e)
}

// delete endpoint to handle put http method with id.
func (h *Handler) delete(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r deleteRequest
	var id int64
	var m *model.Measurement

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if m, e = ShowMeasurement("id", id); e == nil {
			r.Measurement = m
			if e = ctx.Bind(&r); e == nil {
				r.Transform()
				if e = m.Save("is_deleted"); e == nil {
					ctx.Data(m)
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}
