// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pricingType

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for pricingType.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, cuxs.Authorized())
	r.POST("", h.create, auth.CheckPrivilege("pricing_type_create"))
	r.GET("/:id", h.show, auth.CheckPrivilege("pricing_type_show"))
	r.PUT("/:id", h.update, auth.CheckPrivilege("pricing_type_update"))
	r.GET("/:idItemVariant/item-variant-price", h.showItemVariantPrice, cuxs.Authorized())

}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	var total int64
	var data *[]model.PricingType
	if data, total, e = GetPricingTypes(rq); e == nil {
		ctx.Data(data, total)
	}
	return ctx.Serve(e)
}

// create endpoint to handle create http method.
func (h *Handler) create(c echo.Context) (e error) {
	var r createRequest
	ctx := c.(*cuxs.Context)

	if e = ctx.Bind(&r); e == nil {
		pricing := r.Transform()
		if e = pricing.Save(); e == nil {
			ctx.Data(pricing)
		}
	}

	return ctx.Serve(e)
}

// get endpoint to handle get http method.
func (h *Handler) show(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64
	var data *model.PricingType

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if data, e = GetPricingTypeByID(id); e == nil {
			ctx.Data(data)
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// update endpoint handler untuk get id http method
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest
	var id int64

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if r.PricingType, e = GetPricingTypeByID(id); e == nil {
			if e = ctx.Bind(&r); e == nil {
				data := r.Transform(r.PricingType)
				if e = data.Save(); e == nil {
					ctx.Data(data)
				}
			}
		}
	}

	return ctx.Serve(e)
}

// update endpoint handler untuk get id http method
func (h *Handler) showItemVariantPrice(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64
	var data []*model.PricingType
	var itemVariant *model.ItemVariant

	if id, e = common.Decrypt(ctx.Param("idItemVariant")); e == nil {
		if itemVariant, e = GetDetailItemVariant("id", id); e == nil {
			if data, e = GetPricingTypeByItemVariant(itemVariant); e == nil {
				ctx.Data(data)
			} else {
				e = echo.ErrNotFound
			}
		}
	}

	return ctx.Serve(e)
}
