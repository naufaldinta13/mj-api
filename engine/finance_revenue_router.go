// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import "git.qasico.com/mj/api/src/finance_revenue"

func init() {
	handlers["finance-revenue"] = &financeRevenue.Handler{}
}
