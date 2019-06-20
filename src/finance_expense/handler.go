// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package financeExpense

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for financeExpense.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.CheckPrivilege("expense_read"))
	r.GET("/:id", h.show, auth.CheckPrivilege("expense_show"))
	r.POST("", h.create, auth.CheckPrivilege("expense_create"))
	r.PUT("/:id", h.update, auth.CheckPrivilege("expense_update"))
	r.PUT("/:id/approve", h.approve, auth.CheckPrivilege("expense_approve"))

}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	var tot int64
	var data *[]model.FinanceExpense
	if data, tot, e = GetFinanceExpenses(rq); e == nil {
		ctx.Data(data, tot)
	}

	return ctx.Serve(e)
}

// show endpoint to handle get http method with id.
func (h *Handler) show(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64
	var m *model.FinanceExpense

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if m, e = ShowFinanceExpense("id", id); e == nil {
			ctx.Data(m)
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// create endpoint to handle post http method.
func (h *Handler) create(c echo.Context) (e error) {
	var r createRequest
	var fexpense *model.FinanceExpense
	ctx := c.(*cuxs.Context)
	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			fexpense = r.Transform()
			if e = fexpense.Save(); e == nil {
				ctx.Data(fexpense)
			}
		}
	}

	return ctx.Serve(e)
}

// update endpoint to handle put http method by id.
func (h *Handler) update(c echo.Context) (e error) {
	var r updateRequest
	var fexpense *model.FinanceExpense
	ctx := c.(*cuxs.Context)
	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = common.Decrypt(ctx.Param("id")); e == nil {
			if r.FinanceExpense, e = ShowFinanceExpense("id", r.ID); e == nil {
				if e = ctx.Bind(&r); e == nil {
					fexpense = r.Transform()
					if e = fexpense.Save("RecognitionDate", "Amount", "PaymentMethod", "BankNumber", "BankName", "BankHolder", "Note", "UpdatedAt", "UpdatedBy"); e == nil {
						ctx.Data(fexpense)
					}
				}
			} else {
				e = echo.ErrNotFound
			}
		}
	}

	return ctx.Serve(e)
}

// approve endpoint to handle put http method by id.
func (h *Handler) approve(c echo.Context) (e error) {
	var r approveRequest
	var id int64
	var fexpense *model.FinanceExpense
	ctx := c.(*cuxs.Context)
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if fexpense, e = ShowFinanceExpense("id", id); e == nil {
			r.FinanceExpense = fexpense
			if e = ctx.Bind(&r); e == nil {
				if e = ApproveExpense(fexpense); e == nil {
					ctx.Data(fexpense)
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}
