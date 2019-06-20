// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package salesReturn

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for salesReturn.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.CheckPrivilege("sales_return_read"))
	r.GET("/:id", h.show, auth.CheckPrivilege("sales_return_show"))
	r.POST("", h.create, auth.CheckPrivilege("sales_return_create"))
	r.PUT("/:id", h.put, auth.CheckPrivilege("sales_return_update"))
	r.PUT("/:id/cancel", h.cancel, auth.CheckPrivilege("sales_return_cancel"))
}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	var data *[]model.SalesReturn

	var tot int64

	if data, tot, e = GetSalesReturn(rq); e == nil {
		ctx.Data(data, tot)
	}

	return ctx.Serve(e)
}

// show endpoint to handle get http method with id.
func (h *Handler) show(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var as *model.SalesReturn
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if as, e = ShowSalesReturn("id", id); e == nil {
			ctx.Data(as)
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// create endpoint to handle put http method with id.
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest
	var sd *auth.SessionData

	if sd, e = auth.UserSession(ctx); e == nil {
		r.SessionData = sd
		if e = ctx.Bind(&r); e == nil {
			sr := r.Transform()

			if e = sr.Save(); e == nil {

				for _, srItem := range sr.SalesReturnItems {
					srItem.SalesReturn = &model.SalesReturn{ID: sr.ID}
					srItem.Save()
				}

				ctx.Data(sr)
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
	var id int64
	var sr *model.SalesReturn

	if sd, e = auth.UserSession(ctx); e == nil {
		r.SessionData = sd

		if id, e = common.Decrypt(ctx.Param("id")); e == nil {
			if sr, e = ShowSalesReturn("id", id); e == nil {
				r.SR = sr
				r.SO = sr.SalesOrder
				if e = ctx.Bind(&r); e == nil {
					u := r.Transform(sr)

					if e = u.Save("recognition_date", "note", "updated_by", "updated_at", "total_amount", "document_status"); e == nil {
						if e = UpdateSalesReturnItem(sr.SalesReturnItems, u.SalesReturnItems); e == nil {
							ctx.Data(u)
						}
					}
				}
			} else {
				e = echo.ErrNotFound
			}
		}
	}

	return ctx.Serve(e)
}

// put endpoint to handle put http method with id.
func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRequest
	var id int64
	var sr *model.SalesReturn

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if sr, e = ShowSalesReturn("id", id); e == nil {
			r.SR = sr
			if e = ctx.Bind(&r); e == nil {
				if e = CancelSalesReturn(sr); e == nil {
					if e = UpdateExpense(sr); e == nil {
						ctx.Data(sr)
					}
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}
