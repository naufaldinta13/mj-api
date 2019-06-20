// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package partnership

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for brand.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.CheckPrivilege("partnership_read"))
	r.GET("/:id", h.show, auth.CheckPrivilege("partnership_show"))
	r.POST("", h.create, auth.CheckPrivilege("partnership_create"))
	r.PUT("/:id", h.update, auth.CheckPrivilege("partnership_update"))
	r.PUT("/:id/archive", h.archive, auth.CheckPrivilege("partnership_archive"))
	r.PUT("/:id/unarchive", h.unarchive, auth.CheckPrivilege("partnership_archive"))
	r.DELETE("/:id", h.delete, auth.CheckPrivilege("partnership_delete"))
}

// get endpoint to handle get http method.
func (h *Handler) get(c echo.Context) (e error) {
	var t int64
	var m *[]model.Partnership
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	if m, t, e = GetAllPartnerships(rq); e == nil {
		ctx.Data(m, t)
	}
	return ctx.Serve(e)
}

// show endpoint to get detail partnership
func (h *Handler) show(c echo.Context) (e error) {
	var id int64
	var m *model.Partnership
	ctx := c.(*cuxs.Context)
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if m, e = GetPartnershipByField("id", id); e == nil {
			ctx.Data(m)
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// create endpoint to create partnership
func (h *Handler) create(c echo.Context) (e error) {
	var session *auth.SessionData
	var r createRequest
	ctx := c.(*cuxs.Context)
	if session, e = auth.UserSession(ctx); e == nil {
		r.Session = session
		if e = ctx.Bind(&r); e == nil {
			m := r.Transform()
			if m.Save(); e == nil {
				ctx.Data(m)
			}
		}
	}
	return ctx.Serve(e)
}

// update endpoint to update partnership
func (h *Handler) update(c echo.Context) (e error) {
	var id int64
	var session *auth.SessionData
	var r updateRequest
	var u *model.Partnership
	ctx := c.(*cuxs.Context)
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if u, e = GetPartnershipByField("id", id); e == nil {
			if session, e = auth.UserSession(ctx); e == nil {
				r.Session = session
				r.PartnerOld = u
				if e = ctx.Bind(&r); e == nil {
					m := r.Transform()
					if m.Save("order_rule", "max_plafon", "full_name", "email", "phone", "address", "city", "province", "bank_name", "bank_holder", "bank_number", "sales_person", "visit_day", "note", "updated_by", "updated_at"); e == nil {
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

// archive endpoint to handle put http method for archive.
func (h *Handler) archive(c echo.Context) (e error) {
	var id int64
	var u *model.Partnership
	var session *auth.SessionData
	var r archiveRequest
	ctx := c.(*cuxs.Context)
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if u, e = GetPartnershipByField("id", id); e == nil {
			if session, e = auth.UserSession(ctx); e == nil {
				r.Session = session
				r.Partner = u
				if e = ctx.Bind(&r); e == nil {
					r.Transform()
					if u.Save("is_archived", "updated_by", "updated_at"); e == nil {
						ctx.Data(u)
					}
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// unarchive endpoint to handle put http method for unarchive.
func (h *Handler) unarchive(c echo.Context) (e error) {
	var id int64
	var u *model.Partnership
	var session *auth.SessionData
	var r unarchiveRequest
	ctx := c.(*cuxs.Context)
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if u, e = GetPartnershipByField("id", id); e == nil {
			if session, e = auth.UserSession(ctx); e == nil {
				r.Session = session
				r.Partner = u
				if e = ctx.Bind(&r); e == nil {
					r.Transform()
					if u.Save("is_archived", "updated_by", "updated_at"); e == nil {
						ctx.Data(u)
					}
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// delete endpoint to handle delete http method for delete.
func (h *Handler) delete(c echo.Context) (e error) {
	var id int64
	var u *model.Partnership
	var r deleteRequest
	ctx := c.(*cuxs.Context)
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if u, e = GetPartnershipByField("id", id); e == nil {
			r.Partner = u
			if e = ctx.Bind(&r); e == nil {
				r.Transform()
				if u.Save("is_deleted"); e == nil {
					ctx.Data(u)
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}
