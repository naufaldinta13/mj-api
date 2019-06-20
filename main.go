// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"git.qasico.com/mj/api/engine"
	"git.qasico.com/mj/api/src/finance_revenue"

	"git.qasico.com/cuxs/common/log"
	"git.qasico.com/cuxs/cuxs"
	"github.com/robfig/cron"
)

// init preparing application instances.
func init() {
	log.DebugMode = cuxs.IsDebug()
	log.Log = log.New()

	if e := cuxs.DbSetup(); e != nil {
		panic(e)
	}

	cronTask()
}

// main creating new instances application
// and serving application server.
func main() {
	// starting server
	cuxs.StartServer(engine.Router())
}

// cronTask run function like a cronjob
func cronTask() {
	c := cron.New()

	// run auto invoice monthly
	c.AddFunc("0 0 * * sat", func() {
		financeRevenue.Cron()
	})

	c.Start()
}
