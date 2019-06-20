// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package test

import (
	"bytes"
	"fmt"

	"git.qasico.com/mj/api/engine"

	"git.qasico.com/cuxs/common/log"
	"git.qasico.com/cuxs/cuxs"
	"git.qasico.com/cuxs/env"
	"git.qasico.com/cuxs/orm"
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

// Setup testing bootstrap setup.
func Setup() {
	var output bytes.Buffer
	log.Log.Out = &output
	log.Log.Level = logrus.ErrorLevel

	env.Load("../../.env")

	cuxs.Config.DbHost = env.GetString("TESTDB_HOST", "0.0.0.0:3306")
	cuxs.Config.DbName = env.GetString("TESTDB_NAME", "")
	cuxs.Config.DbUser = env.GetString("TESTDB_USERNAME", "")
	cuxs.Config.DbPassword = env.GetString("TESTDB_PASSWORD", "")

	if e := cuxs.DbSetup(); e != nil {
		panic(e)
	}
}

// Router get engine routers.
func Router() *echo.Echo {
	return engine.Router()
}

// DbClean cleaning all data from databases.
func DbClean(table ...string) {
	orm := orm.NewOrm()
	for _, t := range table {
		_, e := orm.Raw(fmt.Sprintf("Delete From %s where id > ?", t), 0).Exec()
		if e != nil {
			panic(e)
		}
		orm.Raw(fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = 1;", t)).Exec()
	}
}

// DataCleanUp clean all data without resetting initial data.
func DataCleanUp(tables ...string) {
	DbClean(tables...)

	var table = []struct {
		Table string
		ID    int
	}{
		{"application_menu", 34},
		{"application_privilege", 485},
		{"application_module", 161},
		{"usergroup", 4},
		{"user", 1},
		{"application_setting", 16},
	}

	orm := orm.NewOrm()
	for _, d := range table {
		orm.Raw(fmt.Sprintf("Delete From %s where id > ?", d.Table), d.ID).Exec()
		orm.Raw(fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = %d;", d.Table, d.ID)).Exec()
	}
}
