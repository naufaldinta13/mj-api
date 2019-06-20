// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package directPlacement

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for directPlacement.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/:id", h.show, auth.CheckPrivilege("direct_placement_show"))
	r.POST("", h.create, auth.CheckPrivilege("direct_placement_create"))
}

// create endpoint to handle post http method.
func (h *Handler) create(c echo.Context) (e error) {
	var r createRequest
	var session *auth.SessionData
	ctx := c.(*cuxs.Context)
	if session, e = auth.UserSession(ctx); e == nil {
		r.Session = session
		if e = ctx.Bind(&r); e == nil {
			m := r.Transform()
			if e = CreateDirectPlacement(&m); e == nil {
				ctx.Data(m)
			}
		}
	}
	return ctx.Serve(e)
}

// show endpoint to handle get http method with id.
func (h *Handler) show(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var m *model.DirectPlacement
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if m, e = ShowDirectPlacement("id", id); e == nil {
			ctx.Data(m)
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}
