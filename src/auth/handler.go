// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for auth.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.POST("", h.signin)
	r.GET("/me", h.me, cuxs.Authorized())
}

// signin endpoint to handle post http method.
func (h *Handler) signin(c echo.Context) (e error) {
	var r SignInRequest
	var sd *SessionData

	ctx := c.(*cuxs.Context)
	if e = ctx.Bind(&r); e == nil {
		if sd, e = Login(r.User); e == nil {
			ctx.Data(sd)
		}
	}
	return ctx.Serve(e)
}

// me endpoint untuk get sesion data yang lagi login.
func (h *Handler) me(c echo.Context) (e error) {
	var sd *SessionData

	ctx := c.(*cuxs.Context)
	// get current user dan data application menu
	if sd, e = UserSession(ctx); e == nil {
		ctx.Data(sd)
	}
	return ctx.Serve(e)
}
