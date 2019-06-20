// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package report

import (
	"fmt"
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/partnership"
	"git.qasico.com/mj/api/src/util"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
)

// GetReportSalesItems digunakan untuk mengambil data-data sales order item
func GetReportSalesItems(rq *orm.RequestQuery) (m []*model.SalesOrderItem, total int64, err error) {
	q, _ := rq.Query(new(model.SalesOrderItem))
	q = q.RelatedSel().Filter("sales_order_id__is_deleted", 0)

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.SalesOrderItem
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		for _, u := range mx {
			var stock *model.ItemVariantStockLog
			if stock, err = GetDetailItemVariantStockLog("item_variant_stock_id__item_variant_id__id", u.ItemVariant.ID); err == nil {
				u.ItemVariantStockLog = stock
			}
		}
		return mx, total, nil
	}
	//if something wrong
	return nil, total, err
}

// GetDetailItemVariantStockLog digunakan untuk mengambil data item variant stock log
func GetDetailItemVariantStockLog(field string, values ...interface{}) (*model.ItemVariantStockLog, error) {
	m := new(model.ItemVariantStockLog)
	o := orm.NewOrm().QueryTable(m)
	if err := o.Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetReportPurchaseItems digunakan untuk mengambil data-data purchase order item
func GetReportPurchaseItems(rq *orm.RequestQuery) (m []*model.PurchaseOrderItem, total int64, err error) {
	q, _ := rq.Query(new(model.PurchaseOrderItem))
	q = q.RelatedSel().Filter("purchase_order_id__is_deleted", 0)

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.PurchaseOrderItem
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		for _, u := range mx {
			var partner *model.Partnership
			if partner, err = partnership.GetPartnershipByField("id", u.PurchaseOrder.Supplier.ID); err == nil {
				u.Partnership = partner
			}
		}

		return mx, total, nil
	}
	//if something wrong
	return nil, total, err
}

// SalesData menampung data jumlah transaksi, total cost, total revenue
type SalesData struct {
	Costs        float64 `json:"costs"`
	Revenues     float64 `json:"revenues"`
	Transactions int64   `json:"transactions"`
}

// GetTotalSales mengambil jumlah total sales
func GetTotalSales(custID string, gteDate string, lteDate string) (*SalesData, error) {
	var e error
	var cost float64
	var revenue float64
	var transactions int64

	qb, _ := orm.NewQueryBuilder("mysql")
	qb = qb.Select("Count(*) AS total", "SUM(sales_order.total_cost) AS costs", "SUM(sales_order.total_charge) AS charge ").
		From("sales_order").
		Where("sales_order.is_deleted = 0 AND sales_order.document_status in( 'new', 'active', 'finished') ")

	if custID != "" {
		cID, _ := common.Decrypt(custID)
		qb.And(fmt.Sprintf("sales_order.customer_id = %d", cID))
	}
	if gteDate != "" {
		gte, _ := util.FormatDateToTimestamp(time.RFC3339, gteDate)
		qb.And(fmt.Sprintf("sales_order.recognition_date >= '%s'", gte.String()))
	}

	if lteDate != "" {
		lte, _ := util.FormatDateToTimestamp(time.RFC3339, lteDate)
		qb.And(fmt.Sprintf("sales_order.recognition_date <= '%s'", lte.String()))
	}

	qb.And(fmt.Sprintf("(sales_order.recognition_date >= '%s' OR NOT sales_order.invoice_status = 'finished')", util.MaxDate()))

	sql := qb.String()
	o := orm.NewOrm()

	if e = o.Raw(sql).QueryRow(&transactions, &cost, &revenue); e == nil {
		sd := &SalesData{
			Costs:        cost,
			Revenues:     revenue,
			Transactions: transactions,
		}

		return sd, nil
	}
	return nil, e
}

// PurchaseData menampung data jumlah transaksi, total expense
type PurchaseData struct {
	Expenses     float64 `json:"expenses"`
	Transactions int64   `json:"transactions"`
}

// GetTotalPurchase mengambil jumlah total purchase
func GetTotalPurchase(suppID string, gteDate string, lteDate string) (*PurchaseData, error) {
	var e error
	var expenses float64
	var transactions int64

	qb, _ := orm.NewQueryBuilder("mysql")
	qb = qb.Select("Count(*) AS total", "SUM(purchase_order.total_charge) AS charge ").
		From("purchase_order").
		Where("purchase_order.is_deleted = 0 AND purchase_order.document_status != 'cancelled' ")

	if gteDate != "" {
		gte, _ := util.FormatDateToTimestamp(time.RFC3339, gteDate)
		qb.And(fmt.Sprintf("purchase_order.recognition_date >= '%s'", gte.String()))
	}

	if lteDate != "" {
		lte, _ := util.FormatDateToTimestamp(time.RFC3339, lteDate)
		qb.And(fmt.Sprintf("purchase_order.recognition_date <= '%s'", lte.String()))
	}

	qb.And(fmt.Sprintf("purchase_order.recognition_date >= '%s'", util.MaxDate()))

	if suppID != "" {
		cID, _ := common.Decrypt(suppID)
		qb.And(fmt.Sprintf("purchase_order.supplier_id = %d", cID))
	}

	sql := qb.String()
	o := orm.NewOrm()

	if e = o.Raw(sql).QueryRow(&transactions, &expenses); e == nil {
		sd := &PurchaseData{
			Expenses:     expenses,
			Transactions: transactions,
		}

		return sd, nil
	}
	return nil, e
}

// SalesItemData menampung data jumlah transaksi, total revenue
type SalesItemData struct {
	Revenues                  float64 `json:"revenues"`
	TotalQuantityTransactions float32 `json:"total_quantity_transactions"`
}

// GetTotalSalesItem mengambil jumlah total sales
func GetTotalSalesItem(category string, iVar string, gteDate string, lteDate string) (*SalesItemData, error) {
	var e error
	var quantity float32
	var revenue float64

	qb, _ := orm.NewQueryBuilder("mysql")
	qb = qb.Select("SUM(sales_order_item.subtotal) AS subtotal", "SUM(sales_order_item.quantity) AS totalquantity ").
		From("sales_order_item").
		InnerJoin("sales_order").On("sales_order_item.sales_order_id = sales_order.id")

	if category != "" {
		cID, _ := common.Decrypt(category)
		qb.InnerJoin("item_variant").On("sales_order_item.item_variant_id = item_variant.id").
			InnerJoin("item").On("item_variant.item_id = item.id").
			Where(fmt.Sprintf("item.category_id  = %d", cID))
	}

	if iVar != "" {
		ivID, _ := common.Decrypt(iVar)
		qb.And(fmt.Sprintf("sales_order_item.item_variant_id = %d", ivID))
	}

	if gteDate != "" {
		gte, _ := util.FormatDateToTimestamp(time.RFC3339, gteDate)
		qb.And(fmt.Sprintf("sales_order.recognition_date >= '%s'", gte.String()))
	}

	if lteDate != "" {
		lte, _ := util.FormatDateToTimestamp(time.RFC3339, lteDate)
		qb.And(fmt.Sprintf("sales_order.recognition_date <= '%s'", lte.String()))
	}

	qb.And("document_status").
		In("'new','active','finished'")

	qb.And(fmt.Sprintf("(sales_order.recognition_date >= '%s' OR NOT sales_order.invoice_status = 'finished')", util.MaxDate()))

	sql := qb.String()
	o := orm.NewOrm()

	if e = o.Raw(sql).QueryRow(&revenue, &quantity); e == nil {

		sd := &SalesItemData{
			Revenues:                  revenue,
			TotalQuantityTransactions: quantity,
		}

		return sd, nil
	}

	return nil, e
}

// PurchaseItemData menampung data jumlah transaksi, total expense
type PurchaseItemData struct {
	TotalQuantityTransactions float32 `json:"total_quantity_transactions"`
	Expenses                  float64 `json:"expenses"`
}

// GetTotalPurchaseItem mengambil jumlah total purchase
func GetTotalPurchaseItem(category string, iVar string) (*PurchaseItemData, error) {
	var e error
	var quantity float32
	var expense float64

	qb, _ := orm.NewQueryBuilder("mysql")
	qb = qb.Select("SUM(purchase_order_item.subtotal) AS subtotal", "SUM(purchase_order_item.quantity) AS totalquantity ").
		From("purchase_order_item").
		InnerJoin("purchase_order").On("purchase_order_item.purchase_order_id = purchase_order.id")

	if category != "" {
		cID, _ := common.Decrypt(category)
		qb.InnerJoin("item_variant").On("purchase_order_item.item_variant_id = item_variant.id").
			InnerJoin("item").On("item_variant.item_id = item.id").
			Where(fmt.Sprintf("item.category_id  = %d", cID))
	}

	if iVar != "" {
		ivID, _ := common.Decrypt(iVar)
		qb.And(fmt.Sprintf("purchase_order_item.item_variant_id = %d", ivID))
	}

	qb.And("document_status").
		In("'new','active','finished'")

	qb.And(fmt.Sprintf("purchase_order.recognition_date >= '%s'", util.MaxDate()))

	sql := qb.String()
	o := orm.NewOrm()

	if e = o.Raw(sql).QueryRow(&expense, &quantity); e == nil {
		sd := &PurchaseItemData{
			Expenses:                  expense,
			TotalQuantityTransactions: quantity,
		}

		return sd, nil
	}

	return nil, e
}

// RevenueData menampung data jumlah transaksi, total balance
type RevenueData struct {
	TotalBalance float64 `json:"total_balance"`
	Transactions int64   `json:"transactions"`
}

// GetTotalRevenueBank mengambil jumlah total revenue dengan transaction_type
func GetTotalRevenueBank(custID string, bankID string) (*RevenueData, error) {
	var e error
	var total float64
	var transactions int64

	qb, _ := orm.NewQueryBuilder("mysql")
	qb = qb.Select("SUM(finance_revenue.amount) AS subtotal, Count(finance_revenue.id) AS total").
		From("finance_revenue")

	if bankID != "" {
		bID := common.ToInt(bankID)
		qb.InnerJoin("sales_invoice").On("sales_invoice.id = finance_revenue.ref_id").
			InnerJoin("sales_order").On("sales_order.id = sales_invoice.sales_order_id").
			Where(fmt.Sprintf("finance_revenue.bank_account_id  = %d", bID))
	}

	//if lteDate != "" {
	//	lte, _ := util.FormatDateToTimestamp(time.RFC3339, lteDate)
	//	qb.And(fmt.Sprintf("sales_order.recognition_date <= '%s'", lte.String()))
	//}
	qb.And("finance_revenue.is_deleted = 0 AND finance_revenue.payment_method = 'debit_card' and finance_revenue.ref_type = 'sales_invoice' ")

	qb.And(fmt.Sprintf("finance_revenue.recognition_date >= '%s'", util.MaxDate()))

	sql := qb.String()
	o := orm.NewOrm()

	if e = o.Raw(sql).QueryRow(&total, &transactions); e == nil {
		rd := &RevenueData{
			TotalBalance: total,
			Transactions: transactions,
		}

		return rd, nil
	}
	return nil, e
}
