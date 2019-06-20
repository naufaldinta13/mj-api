// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bankAccount

import (
	"testing"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestGetBankAccount(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM bank_account").Exec()

	ba := model.DummyBankAccount()
	ba.BankName = "BCA"
	ba.BankNumber = "1872987391723"
	ba.Save("BankName")
	qs := orm.RequestQuery{}
	_, tot, e := GetBankAccounts(&qs)
	assert.NoError(t, e, "Data should be exists.")
	assert.Equal(t, int64(1), tot)
}
