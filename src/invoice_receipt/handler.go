// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package invoiceReceipt

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for salesInvoice.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.CheckPrivilege("invoice_receipt_read"))
	r.GET("/:id", h.show, auth.CheckPrivilege("invoice_receipt_show"))
	r.POST("", h.create, auth.CheckPrivilege("invoice_receipt_create"))
	r.PUT("/:id/payment", h.payment, auth.CheckPrivilege("invoice_receipt_pay"))
}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	var data *[]model.InvoiceReceipt
	var tot int64

	if data, tot, e = GetInvoiceReceipts(rq); e == nil {
		ctx.Data(data, tot)
	}

	return ctx.Serve(e)
}

// show endpoint to handle get http method with id.
func (h *Handler) show(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64
	var m *model.InvoiceReceipt

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if m, e = GetDetailInvoiceReceipt("id", id); e == nil {
			ctx.Data(m)
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// create endpoint to handle post http method.
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest
	var sd *auth.SessionData

	if sd, e = auth.UserSession(ctx); e == nil {
		r.SessionData = sd
		if e = ctx.Bind(&r); e == nil {
			m, mso := r.Transform()
			if e = CreateInvoiceReceipt(m, mso); e == nil {
				ctx.Data(m)
			}
		}
	}

	return ctx.Serve(e)
}

// payment endpoint to handle put http method.
func (h *Handler) payment(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var r PaymentInvoiceReceipt
	var sd *auth.SessionData
	var mir *model.InvoiceReceipt
	var mfr []*model.FinanceRevenue
	var miri []*model.InvoiceReceiptItem
	var mirr []*model.InvoiceReceiptReturn

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if sd, e = auth.UserSession(ctx); e == nil {
			if mir, e = getDetailInvoiceReceipt("id", id); e == nil {
				r.SessionData = sd
				r.InvoiceReceiptID = mir.ID
				if e = ctx.Bind(&r); e == nil {
					miri, mirr, mfr = r.Transform(mir.ID)
					if e = CreateInvoiceReceiptPayment(id, miri, mirr, mfr); e == nil {
						ctx.Data(mfr)
					}
				}
			} else {
				e = echo.ErrNotFound
			}
		}
	}

	return ctx.Serve(e)
}
