// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for user.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.POST("", h.create, auth.CheckPrivilege("user_create"))
	r.GET("", h.get, auth.CheckPrivilege("user_read"))
	r.GET("/:id", h.show, auth.CheckPrivilege("user_show"))
	r.PUT("/:id/change-password", h.changePassword, auth.CheckPrivilege("user_change_password"))
	r.PUT("/:id/inactive", h.inactive, auth.CheckPrivilege("user_inactive"))
}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	var total int64
	var data *[]model.User
	var session *auth.SessionData
	if session, e = auth.UserSession(ctx); e == nil {
		if data, total, e = GetUsers(rq, session); e == nil {
			ctx.Data(data, total)
		}
	}
	return ctx.Serve(e)
}

// show endpoint to handle get http method with id.
func (h *Handler) show(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var data *model.User
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if data, e = GetUserByID(id); e == nil {
			ctx.Data(data)
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// create endpoint to handle post http method
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if e = ctx.Bind(&r); e == nil {
		data := r.Transform()
		if e = data.Save(); e == nil {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// changePassword endpoint to handle put http method.
func (h *Handler) changePassword(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r changePwdRequest
	var id int64
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if r.User, e = GetUserByID(id); e == nil {
			if e = ctx.Bind(&r); e == nil {
				data := r.Transform(r.User)
				if e = data.Save("Password"); e == nil {
					ctx.Data(data)
				}
			}
		}
	}
	return ctx.Serve(e)
}

// inactive endpoint to handle put http method.
func (h *Handler) inactive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64
	var user *model.User

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if user, e = GetUserByID(id); e == nil && user.IsActive == 1 {
			if e = user.Inactive(); e == nil {
				ctx.Data(user)
			}
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}
