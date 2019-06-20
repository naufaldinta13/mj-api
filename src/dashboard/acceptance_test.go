// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dashboard_test

import (
	"fmt"
	"os"
	"testing"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/test"

	"git.qasico.com/cuxs/common/tester"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.Setup()

	// run tests
	res := m.Run()

	// cleanup
	test.DataCleanUp()

	os.Exit(res)
}

// TestHandler_URLMappingDashboard untuk mengetest get dasboard
func TestHandler_URLMappingDashboard(t *testing.T) {
	// buat dummy
	model.DummySalesOrderItem()
	model.DummyPurchaseOrderItem()
	model.DummyWorkorderReceivingItem()
	model.DummyWorkorderShipmentItem()
	model.DummyWorkorderFulfillmentItem()
	model.DummyFinanceRevenue()
	model.DummySalesReturnItem()
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/dashboard?month=12&year=2010"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/dashboard", "GET"))
	})
}
