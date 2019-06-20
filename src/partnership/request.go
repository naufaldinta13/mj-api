// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package partnership

import (
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create partnership process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	Session         *auth.SessionData
	PartnershipType string  `json:"partnership_type" valid:"required|in:customer,supplier"`
	OrderRule       string  `json:"order_rule" valid:"required|in:none,one_bill,plafon"`
	MaxPlafon       float64 `json:"max_plafon" valid:"range:0,9999999999999999999"`
	FullName        string  `json:"full_name" valid:"required|lte:45"`
	Email           string  `json:"email" valid:"lte:45"`
	Phone           string  `json:"phone" valid:"lte:45"`
	Address         string  `json:"address" valid:"lte:255"`
	City            string  `json:"city" valid:"lte:45"`
	Province        string  `json:"province" valid:"lte:45"`
	BankName        string  `json:"bank_name" valid:"lte:45"`
	BankNumber      string  `json:"bank_number" valid:"lte:45"`
	BankHolder      string  `json:"bank_holder" valid:"lte:45"`
	SalesPerson     string  `json:"sales_person" valid:"lte:45"`
	VisitDay        string  `json:"visit_day" valid:"lte:45"`
	Note            string  `json:"note" valid:"lte:255"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	// jika order rule plafon maka max plafon tidak boleh 0
	if r.OrderRule == "plafon" {
		if r.MaxPlafon <= float64(0) {
			o.Failure("max_plafon", "must be greater than 0")
		}
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *createRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform untuk mengubah isi model.
func (r *createRequest) Transform() *model.Partnership {
	m := &model.Partnership{
		Code:            common.RandomNumeric(7),
		PartnershipType: r.PartnershipType,
		OrderRule:       r.OrderRule,
		FullName:        r.FullName,
		Email:           r.Email,
		Phone:           r.Phone,
		Address:         r.Address,
		City:            r.City,
		Province:        r.Province,
		BankName:        r.BankName,
		BankNumber:      r.BankNumber,
		BankHolder:      r.BankHolder,
		MaxPlafon:       r.MaxPlafon,
		SalesPerson:     r.SalesPerson,
		VisitDay:        r.VisitDay,
		Note:            r.Note,
		IsArchived:      int8(0),
		IsDeleted:       int8(0),
		IsDefault:       int8(0),
		CreatedBy:       &model.User{ID: r.Session.User.ID},
		CreatedAt:       time.Now(),
	}
	return m
}

// updateRequest data struct that stored request data when requesting an update partnership process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type updateRequest struct {
	PartnerOld  *model.Partnership
	Session     *auth.SessionData
	OrderRule   string  `json:"order_rule" valid:"required|in:none,one_bill,plafon"`
	MaxPlafon   float64 `json:"max_plafon" valid:"range:0,9999999999999999999"`
	FullName    string  `json:"full_name" valid:"required|lte:45"`
	Email       string  `json:"email" valid:"lte:45"`
	Phone       string  `json:"phone" valid:"lte:45"`
	Address     string  `json:"address" valid:"lte:255"`
	City        string  `json:"city" valid:"lte:45"`
	Province    string  `json:"province" valid:"lte:45"`
	BankName    string  `json:"bank_name" valid:"lte:45"`
	BankNumber  string  `json:"bank_number" valid:"lte:45"`
	BankHolder  string  `json:"bank_holder" valid:"lte:45"`
	SalesPerson string  `json:"sales_person" valid:"lte:45"`
	VisitDay    string  `json:"visit_day" valid:"lte:45"`
	Note        string  `json:"note" valid:"lte:255"`
}

// Validate implement validation.Requests interfaces.
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	// jika order rule plafon maka max plafon tidak boleh 0
	if r.OrderRule == "plafon" {
		if r.MaxPlafon <= float64(0) {
			o.Failure("max_plafon", "must be greater than 0")
		}
	}
	// tidak bisa diupdate kalau sudah di delete
	if r.PartnerOld.IsDeleted == int8(1) {
		o.Failure("is_deleted", "cannot be update")
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *updateRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform untuk mengubah isi model.
func (r *updateRequest) Transform() *model.Partnership {
	m := r.PartnerOld
	m.OrderRule = r.OrderRule
	m.MaxPlafon = r.MaxPlafon
	m.FullName = r.FullName
	m.Email = r.Email
	m.Phone = r.Phone
	m.Address = r.Address
	m.City = r.City
	m.Province = r.Province
	m.BankName = r.BankName
	m.BankNumber = r.BankNumber
	m.BankHolder = r.BankHolder
	m.SalesPerson = r.SalesPerson
	m.VisitDay = r.VisitDay
	m.Note = r.Note
	m.UpdatedBy = &model.User{ID: r.Session.User.ID}
	m.UpdatedAt = time.Now()
	return m
}

// archiveRequest data struct that stored request data when requesting an put archive partnership process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type archiveRequest struct {
	Session *auth.SessionData
	Partner *model.Partnership
}

// Validate implement validation.Requests interfaces.
func (r *archiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	if r.Partner.IsArchived == int8(1) {
		o.Failure("is_archived", "already archived")
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *archiveRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform untuk mengubah isi model.
func (r *archiveRequest) Transform() {
	r.Partner.IsArchived = int8(1)
	r.Partner.UpdatedBy = &model.User{ID: r.Session.User.ID}
	r.Partner.UpdatedAt = time.Now()
}

// unarchiveRequest data struct that stored request data when requesting an put unarchive partnership process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type unarchiveRequest struct {
	Session *auth.SessionData
	Partner *model.Partnership
}

// Validate implement validation.Requests interfaces.
func (r *unarchiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	if r.Partner.IsArchived == int8(0) {
		o.Failure("is_archived", "not archived")
	}
	if r.Partner.IsDeleted == int8(1) {
		o.Failure("is_deleted", "already deleted")
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *unarchiveRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform untuk mengubah isi model.
func (r *unarchiveRequest) Transform() {
	r.Partner.IsArchived = int8(0)
	r.Partner.UpdatedBy = &model.User{ID: r.Session.User.ID}
	r.Partner.UpdatedAt = time.Now()
}

// deleteRequest data struct that stored request data when requesting an delete partnership process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type deleteRequest struct {
	Partner *model.Partnership
}

// Validate implement validation.Requests interfaces.
func (r *deleteRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	if r.Partner.IsArchived == int8(0) {
		o.Failure("is_archived", "not archived")
	}
	if r.Partner.IsDeleted == int8(1) {
		o.Failure("is_deleted", "already deleted")
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *deleteRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform untuk mengubah isi model.
func (r *deleteRequest) Transform() {
	r.Partner.IsDeleted = int8(1)
}
