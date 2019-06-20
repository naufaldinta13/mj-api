// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"time"

	"git.qasico.com/mj/api/datastore/model"

	"regexp"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/validation"
)

// createRequest data struct that stored request data when requesting an create user process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type createRequest struct {
	UsergroupID     string `json:"usergroup_id" valid:"required"`
	FullName        string `json:"full_name" valid:"required"`
	Username        string `json:"username" valid:"required"`
	Password        string `json:"password" valid:"required|gte:5"`
	ConfirmPassword string `json:"confirm_password" valid:"required"`
}

// Validate implement validation.Requests interfaces.
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	// validasi usergroup
	if ugID, e := common.Decrypt(r.UsergroupID); e != nil {
		o.Failure("usergroup_id", "usergroup_id cannot be decrypt")
	} else {
		usergroup := &model.Usergroup{ID: ugID}
		if e = usergroup.Read(); e != nil {
			o.Failure("usergroup_id", "usergroup_id doesn't exist")
		}
	}

	if ugID, _ := common.Decrypt(r.UsergroupID); ugID == 1 {
		o.Failure("user_group", "User group not valid")
	}
	// validasi username
	user := &model.User{Username: r.Username}
	if e := user.Read("Username"); e == nil {
		o.Failure("username", "username has been registered")
	}

	res, err := regexp.MatchString("^[a-zA-Z0-9_]*$", r.Username)
	if res == false && err == nil {
		o.Failure("username", "Can't use space or special character in this input field")
	}

	// validasi FullName
	res, err = regexp.MatchString("^[a-zA-Z0-9 ]*$", r.FullName)
	if res == false && err == nil {
		o.Failure("fullname", "Can't use space or special character in this input field")
	}

	// validasi confirm password
	if r.ConfirmPassword != r.Password {
		o.Failure("confirm_password", "confirm_password not same with password")
	}
	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *createRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *createRequest) Transform() *model.User {
	ugID, _ := common.Decrypt(r.UsergroupID)
	pwd, _ := common.PasswordHasher(r.Password)
	user := &model.User{
		Usergroup: &model.Usergroup{ID: ugID},
		FullName:  r.FullName,
		Username:  r.Username,
		Password:  pwd,
		IsActive:  1,
		CreatedAt: time.Now(),
	}
	return user
}

// changePwdRequest data struct that stored request data when requesting an update user process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type changePwdRequest struct {
	Password        string      `json:"password" valid:"required|gte:5"`
	ConfirmPassword string      `json:"confirm_password" valid:"required"`
	User            *model.User `json:"-"`
}

// Validate implement validation.Requests interfaces.
func (r *changePwdRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if r.ConfirmPassword != r.Password {
		o.Failure("confirm_password", "confirm_password not same with password")
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *changePwdRequest) Messages() map[string]string {
	return map[string]string{}
}

// Transform transforming request into model.
func (r *changePwdRequest) Transform(user *model.User) *model.User {
	pwd, _ := common.PasswordHasher(r.Password)
	user.Password = pwd
	user.UpdatedAt = time.Now()

	return user
}
