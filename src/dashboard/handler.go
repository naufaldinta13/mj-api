// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dashboard

import (
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for dashboard.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.CheckPrivilege("dashboard_display"))

}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	var data *DSBData
	ctx := c.(*cuxs.Context)
	// get param in url
	if prm := GetParamURL(ctx); e == nil {
		if data, e = GetDashBoard(&prm); e == nil {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}
