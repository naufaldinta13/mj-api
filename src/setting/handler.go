// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package setting

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for setting.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, cuxs.Authorized())
	r.GET("/:id", h.show, cuxs.Authorized())
	r.PUT("/:id", h.put, auth.CheckPrivilege("application_update"))
}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	var data *[]model.ApplicationSetting

	var tot int64

	if data, tot, e = GetApplicationSetting(rq); e == nil {
		ctx.Data(data, tot)
	}

	return ctx.Serve(e)
}

// show endpoint to handle get http method with id.
func (h *Handler) show(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var as *model.ApplicationSetting
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if as, e = ShowApplicationSetting("id", id); e == nil {
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
	var as *model.ApplicationSetting

	if sd, e = auth.UserSession(ctx); e == nil {

		r.SessionData = sd

		if id, e = common.Decrypt(ctx.Param("id")); e == nil {
			if as, e = ShowApplicationSetting("id", id); e == nil {
				r.ApplicationSetting = as
				if e = ctx.Bind(&r); e == nil {
					if as, e = UpdateApplicationSetting(&r); e == nil {
						ctx.Data(as)
					}
				}
			}
		}
	}

	return ctx.Serve(e)
}
