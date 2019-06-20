// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package setting

import (
	"testing"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestGetApplicationSetting(t *testing.T) {
	model.DummyApplicationSetting()
	qs := orm.RequestQuery{}
	_, _, e := GetApplicationSetting(&qs)
	assert.NoError(t, e, "Data should be exists.")
}

func TestGetAccountingAccount(t *testing.T) {
	_, e := ShowApplicationSetting("id", 1000)
	assert.Error(t, e, "Response should be error, beacuse there are no data yet.")

	c := model.DummyApplicationSetting()
	cd, e := ShowApplicationSetting("id", c.ID)
	assert.NoError(t, e, "Data should be exists.")
	assert.Equal(t, c.ID, cd.ID, "ID Response should be a same.")
}
