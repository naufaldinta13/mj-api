// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"net/http"
	"net/http/httptest"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/cuxs"
	"github.com/labstack/echo"
)

// LoginAs is function to check authorization user_is_login
func LoginAs(user *model.User) (echo.Context, echo.HandlerFunc) {
	sd, _ := Login(user)
	token := sd.Token
	token = "Bearer " + token

	e := cuxs.New()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, token)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	var x = func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	return c, x
}
