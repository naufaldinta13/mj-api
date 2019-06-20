// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"net/url"

	"git.qasico.com/cuxs/cuxs"
	"git.qasico.com/cuxs/validation"
	"github.com/labstack/echo"
)

// Handler collection handler for util.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/generate-code", h.codegen, cuxs.Authorized())
}

// codegen needed some param: prefix,table_name,field_code
// e.g:=?prefix=prefix&table_name=nama_tabel&field_code=kolom_code
func (h *Handler) codegen(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var u url.Values
	var code string
	if u, e = ctx.FormParams(); e == nil {
		if len(u) < int(2) {
			e = validation.SetError("context", "context request query is required")
		} else {
			settingName := u.Get("setting_name")
			if settingName == "" {
				settingName = "code_sales_order"
			}
			table := u.Get("table")
			if table == "" {
				table = "sales_order"
			}
			code, _ = CodeGen(settingName, table)

			if code != "" {
				ctx.Data(code)
			} else {
				e = echo.ErrNotFound
			}
		}
	}

	return ctx.Serve(e)
}
