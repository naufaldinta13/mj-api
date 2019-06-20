// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_test

import (
	"os"
	"testing"

	"git.qasico.com/mj/api/test"
)

func TestMain(m *testing.M) {
	test.Setup()

	// run tests
	res := m.Run()

	// cleanup
	test.DataCleanUp("bank_account","direct_placement", "direct_placement_item", "finance_expense", "finance_revenue", "invoice_receipt", "invoice_receipt_item", "invoice_receipt_return", "item", "item_category", "item_variant", "item_variant_price", "item_variant_stock", "item_variant_stock_log", "measurement", "partnership", "pricing_type", "purchase_invoice", "purchase_order", "purchase_order_item", "purchase_return", "purchase_return_item", "recap_sales", "recap_sales_item", "sales_invoice", "sales_order", "sales_order_item", "sales_return", "sales_return_item", "stockopname", "stockopname_item", "workorder_fulfillment", "workorder_fulfillment_item", "workorder_receiving", "workorder_receiving_item", "workorder_shipment", "workorder_shipment_item")
	os.Exit(res)
}
