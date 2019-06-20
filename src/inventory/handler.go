// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package inventory

import (
	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// Handler collection handler for inventory.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/variant", h.getVariant, auth.CheckPrivilege("item_read"))
	r.GET("/variant/:id", h.showVariant, auth.CheckPrivilege("item_show"))
	r.PUT("/variant/:id/archive", h.archiveVariant, auth.CheckPrivilege("item_archive"))
	r.PUT("/variant/:id/unarchive", h.unarchiveVariant, auth.CheckPrivilege("item_archive"))
	r.DELETE("/variant/:id", h.deleteVariant, auth.CheckPrivilege("item_delete"))
	r.GET("/item/:id", h.showItem, auth.CheckPrivilege("item_show"))
	r.POST("/item", h.createItem, auth.CheckPrivilege("item_create"))
	r.PUT("/item/:id/archive", h.archiveItem, auth.CheckPrivilege("item_archive"))
	r.PUT("/item/:id/unarchive", h.unarchiveItem, auth.CheckPrivilege("item_archive"))
	r.DELETE("/item/:id", h.deleteItem, auth.CheckPrivilege("item_delete"))
	r.PUT("/item/:id", h.updateItem, auth.CheckPrivilege("item_update"))
}

// getVariant endpoint to handle get http method.
func (h *Handler) getVariant(c echo.Context) (e error) {
	var t int64
	var m *[]model.ItemVariant
	ctx := c.(*cuxs.Context)
	// get query string from request
	rq := ctx.RequestQuery()
	if m, t, e = GetAllItemVariant(rq); e == nil {
		ctx.Data(m, t)
	}
	return ctx.Serve(e)
}

// showVariant enpoint to handle get http method with id.
func (h *Handler) showVariant(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var m *model.ItemVariant

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if m, e = GetDetailItemVariant("id", id); e == nil {
			ctx.Data(m)
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// archiveVariant endpoint to handle put http method.
func (h *Handler) archiveVariant(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var m *model.ItemVariant
	var r ArchiveItemVariant

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if m, e = GetDetailItemVariant("id", id); e == nil {
			r.isArchived = m.IsArchived
			r.isDeleted = m.IsDeleted
			if e = ctx.Bind(&r); e == nil {
				if e = m.Archive(); e == nil {
					ctx.Data(m)
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// unarchiveVariant endpoint to handle put http method.
func (h *Handler) unarchiveVariant(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var m *model.ItemVariant
	var r UnarchiveItemVariant

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if m, e = GetDetailItemVariant("id", id); e == nil {
			r.isArchived = m.IsArchived
			r.isDeleted = m.IsDeleted
			if e = ctx.Bind(&r); e == nil {
				if e = m.Unarchive(); e == nil {
					ctx.Data(m)
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// deleteVariant endpoint to handle delete http method.
func (h *Handler) deleteVariant(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var m *model.ItemVariant
	var r DeleteItemVariant

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if m, e = GetDetailItemVariant("id", id); e == nil {
			r.ItemVariantID = m.ID
			r.isArchived = m.IsArchived
			r.isDeleted = m.IsDeleted
			if e = ctx.Bind(&r); e == nil {
				if e = m.DeleteItemVariant(); e == nil {
					ctx.Data(m)
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// showItem enpoint to handle get http method with id.
func (h *Handler) showItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var m *model.Item

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if m, e = GetDetailItem("id", id); e == nil {
			ctx.Data(m)
		} else {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// createItem endpoint to handle get http method.
func (h *Handler) createItem(c echo.Context) (e error) {
	var r createItemRequest
	var session *auth.SessionData
	var u *model.Item
	ctx := c.(*cuxs.Context)
	if session, e = auth.UserSession(ctx); e == nil {
		r.Session = session
		if e = ctx.Bind(&r); e == nil {
			if u, e = SavingItem(r.Transform()); e == nil {
				ctx.Data(u)
			}
		}
	}
	return ctx.Serve(e)
}

// archiveItem endpoint to handle put http method.
func (h *Handler) archiveItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var m *model.Item
	var session *auth.SessionData
	var r archiveItem

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if session, e = auth.UserSession(ctx); e == nil {
			if m, e = GetDetailItem("id", id); e == nil {
				r.Item = m
				if e = ctx.Bind(&r); e == nil {
					if e = ArchivedItem(m, session); e == nil {
						ctx.Data(m)
					}
				}
			} else {
				e = echo.ErrNotFound
			}
		}
	}

	return ctx.Serve(e)
}

// unarchiveItem endpoint to handle put http method.
func (h *Handler) unarchiveItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var m *model.Item
	var session *auth.SessionData
	var r unarchiveItem

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if session, e = auth.UserSession(ctx); e == nil {
			if m, e = GetDetailItem("id", id); e == nil {
				r.Item = m
				if e = ctx.Bind(&r); e == nil {
					if e = UnarchivedItem(m, session); e == nil {
						ctx.Data(m)
					}
				}
			} else {
				e = echo.ErrNotFound
			}
		}
	}
	return ctx.Serve(e)
}

// deleteItem endpoint to handle delete http method.
func (h *Handler) deleteItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	var m *model.Item
	var r deleteItem

	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if m, e = GetDetailItem("id", id); e == nil {
			r.Item = m
			if e = ctx.Bind(&r); e == nil {
				if e = DeleterItem(m); e == nil {
					ctx.Data(m)
				}
			}
		} else {
			e = echo.ErrNotFound
		}
	}
	return ctx.Serve(e)
}

// updateItem endpoint to handle get http method.
func (h *Handler) updateItem(c echo.Context) (e error) {
	var r updateItemRequest
	var session *auth.SessionData
	var id int64
	var m *model.Item
	var u *model.Item

	ctx := c.(*cuxs.Context)
	if id, e = common.Decrypt(ctx.Param("id")); e == nil {
		if m, e = GetDetailItem("id", id); e == nil {
			if session, e = auth.UserSession(ctx); e == nil {
				r.Session = session
				r.OldItem = m
				if e = ctx.Bind(&r); e == nil {
					if u, e = SavingUpdateItem(r.Transform()); e == nil {
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
