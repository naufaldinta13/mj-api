// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package category

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for category.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.CheckPrivilege("category_read"))
	r.POST("", h.create, auth.CheckPrivilege("category_create"))
	r.GET("/:id", h.show, auth.CheckPrivilege("category_show"))
	r.PUT("/:id", h.update, auth.CheckPrivilege("category_update"))
	r.DELETE("/:id", h.delete, auth.CheckPrivilege("category_delete"))
}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	var total int64
	var data *[]model.ItemCategory
	if data, total, e = GetItemCategories(rq); e == nil {
		ctx.Data(data, total)
	}
	return ctx.Serve(e)
}

//// create endpoint to handle post http method.
func (h *Handler) create(c echo.Context) (e error) {
	var r createRequest

	ctx := c.(*cuxs.Context)
	if e = ctx.Bind(&r); e == nil {
		data := r.Transform()
		if e = data.Save(); e == nil {
			ctx.Data(data)
		}
	}
	return ctx.Serve(e)
}

// show endpoint to handle get http method with id.
func (h *Handler) show(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var data *model.ItemCategory
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if data, e = GetItemCategoryByID(id); e == nil {
			ctx.Data(data)
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

//// update endpoint to handle put http method.
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest
	var category *model.ItemCategory

	if r.ID, e = common.Decrypt(ctx.Param("id")); e == nil {
		if category, e = GetItemCategoryByID(r.ID); e == nil {
			if e = ctx.Bind(&r); e == nil {
				data := r.Transform(category)
				if e = data.Save(); e == nil {
					ctx.Data(data)
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// delete endpoint to handle delete http method.
func (h *Handler) delete(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r deleteRequest
	var id int64
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if r.Category, e = GetItemCategoryByID(id); e == nil {
			if e = ctx.Bind(&r); e == nil {
				r.Category.ChangeStatusDelete()
			}
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}
