// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bankAccount

import (
	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
)

// GetBankAccounts get all data item_category that matched with query request parameters.
// returning slices of bank account, total data without limit and error.
func GetBankAccounts(rq *orm.RequestQuery) (m *[]model.BankAccount, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.BankAccount))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []model.BankAccount
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return &mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}