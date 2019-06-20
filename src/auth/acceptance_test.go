// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth_test

import (
	"fmt"
	"net/http"
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

func TestARouting(t *testing.T) {
	// mentest endpoint yang memiliki middleware session
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/auth", "POST", 415},
		{"/v1/auth/me", "GET", 200},
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, ep := range routers {
		ng.Method = ep.method
		ng.Path = ep.endpoint
		ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, ep.expected, res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", ep.endpoint, ep.method))
		})
	}
}

func TestAuthenticationHandler_URLMapping(t *testing.T) {
	var auth = []struct {
		req      tester.D
		expected int
	}{
		{tester.D{"username": "", "password": ""}, http.StatusUnprocessableEntity},
		{tester.D{"username": "sysadmin", "password": "xxxxx"}, http.StatusUnprocessableEntity},
		{tester.D{"username": "qasico", "password": "qasico123"}, http.StatusUnprocessableEntity},
		{tester.D{"username": "sysadmin", "password": "qasico123"}, http.StatusOK},
	}
	ng := tester.New()
	for _, tes := range auth {
		ng.POST("/v1/auth").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}
