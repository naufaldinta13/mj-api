// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package shipment

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for shipment.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.CheckPrivilege("shipment_show"))
	r.GET("/:id", h.show, auth.CheckPrivilege("shipment_read"))
	r.POST("", h.create, auth.CheckPrivilege("shipment_create"))
	r.PUT("/:id", h.update, auth.CheckPrivilege("shipment_update"))
	r.PUT("/:wsid/cancel/:fulid", h.shipmentCancel, auth.CheckPrivilege("shipment_cancel"))
	r.PUT("/:id/approve", h.approve, auth.CheckPrivilege("shipment_approve"))
}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()

	if data, total, e := GetShipments(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

func (h *Handler) create(c echo.Context) (e error) {
	var session *auth.SessionData
	var r createRequest
	var m model.WorkorderShipment
	ctx := c.(*cuxs.Context)
	if session, e = auth.UserSession(ctx); e == nil {
		r.Session = session
		if e = ctx.Bind(&r); e == nil {
			if m, e = r.Transform(); e == nil {
				if e = CreateWorkOrderShipment(&m); e == nil {
					ctx.Data(m)
				}
			}
		}
	}
	return ctx.Serve(e)
}

// show endpoint to handle get http method with id.
func (h *Handler) show(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var shipment *model.WorkorderShipment
	var id int64
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {

		if shipment, e = ShowShipment("id", id); e == nil {
			ctx.Data(shipment)
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// update endpoint to handle put http method.
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var m, shipment *model.WorkorderShipment
	var r updateRequest
	if r.WorkorderShipmentID, e = common.Decrypt(ctx.Param("id")); e == nil {
		if shipment, e = GetDetailShipment("id", r.WorkorderShipmentID); e == nil {
			if e = ctx.Bind(&r); e == nil {
				if m, e = r.Transform(shipment); e == nil {
					if e = CreateWorkOrderShipment(m); e == nil {
						ctx.Data(m)
					}
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

func (h *Handler) shipmentCancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var wsid, fulid int64
	var shipx, shipment *model.WorkorderShipment
	var fulfill model.WorkorderFulfillment
	if wsid, e = common.Decrypt(ctx.Param("wsid")); e == nil {
		if fulid, e = common.Decrypt(ctx.Param("fulid")); e == nil {
			if shipment, e = GetDetailShipment("id", wsid); e == nil {
				fulfill = model.WorkorderFulfillment{ID: fulid, IsDeleted: int8(0)}
				if e = fulfill.Read("ID", "IsDeleted"); e == nil {
					if shipx, e = cancelShipmentAndFulfillment(shipment, fulfill.ID); e == nil {
						ctx.Data(&shipx)
					}
				} else {
					e = echo.ErrNotFound
				}
			} else {
				e = echo.ErrNotFound
			}
		}
	}
	return ctx.Serve(e)
}

func (h *Handler) approve(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64
	var shipment, shipmentx *model.WorkorderShipment
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if shipment, e = GetDetailShipment("id", id); e == nil {
			if shipmentx, e = approveShipment(shipment); e == nil {
				ctx.Data(shipmentx)
			}
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}
