// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"git.qasico.com/mj/api/src/auth"
)

// GetUserByID untuk get data user berdasarkan id user
// return : data user dan error
func GetUserByID(id int64) (m *model.User, err error) {
	mx := new(model.User)
	o := orm.NewOrm().QueryTable(mx)

	if err = o.Filter("id", id).RelatedSel().Limit(1).One(mx); err != nil {
		return nil, err
	}
	return mx, nil
}

// GetUsers get all data user that matched with query request parameters.
// returning slices of users, total data without limit and error.
func GetUsers(rq *orm.RequestQuery, session *auth.SessionData) (m *[]model.User, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.User))

	if session.User.Usergroup.ID != int64(1) {
		q = q.Filter("id", session.User.ID)
	}

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.User
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
