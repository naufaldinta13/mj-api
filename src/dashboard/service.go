// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dashboard

import (
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/inventory"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/cuxs"
	"git.qasico.com/cuxs/orm"
)

// DSBData menampung semua data dasboard
type DSBData struct {
	TotalSales              int64                `json:"total_sales"`
	TotalPurchase           int64                `json:"total_purchases"`
	TotalReturn             int64                `json:"total_return"`
	TotalItemTerjual        float32              `json:"item_terjual"`
	TotalStockIn            float32              `json:"stock_masuk"`
	TotalStockReturn        float32              `json:"stock_return"`
	GraphicReport           []*Graphic           `json:"graphic_report"`
	GraphicTopBarangRevenue []*TopBarangRevenue  `json:"top_barang_revenue"`
	ToDoListReport          *ToDoList            `json:"todo_list"`
	OverViewReport          *OverView            `json:"overview"`
	NeedToRestock           []*model.ItemVariant `json:"need_to_restock"`
}

// Graphic menampung data penjualan setiap hari nya
type Graphic struct {
	Date         time.Time `json:"date"`
	Transaction  int64     `json:"transaction"`
	TotalRevenue float64   `json:"total_revenue"`
	ItemSold     float32   `json:"item_sold"`
}

// TopBarangRevenue menampung data top 10 revenue item
type TopBarangRevenue struct {
	ItemVariant   *model.ItemVariant `json:"item_variant"`
	TotalSubTotal float64            `json:"total_subtotal"`
}

// ToDoList menampung data ToDoList
type ToDoList struct {
	SalesNeedFulfilled     int64   `json:"sales_need_fulfilled"`
	SalesNeedInvoiced      int64   `json:"sales_need_invoiced"`
	SalesNeedApproveCancel int64   `json:"sales_need_approve_cancel"`
	FulfillmentNeedShipped int64   `json:"fulfillment_need_shipped"`
	ItemNeedReceived       float32 `json:"item_need_received"`
	ShipmentNeedDelivered  int64   `json:"shipment_need_delivered"`
}

// OverView menampung data OverView
type OverView struct {
	TotalPenjualan float64 `json:"total_penjualan"`
	TotalCost      float64 `json:"total_cost"`
	TotalProfit    float64 `json:"total_profit"`
	PiutangUsaha   float64 `json:"piutang_usaha"`
	HutangUsaha    float64 `json:"hutang_usaha"`
}

// GetDashBoard untuk mengambil semua data dashboard
func GetDashBoard(pr *ParamURL) (m *DSBData, e error) {
	var ret int64
	var sales, purchase SOAndPOData
	var stockIn, StockRet, itemJual float32
	var graph []*Graphic
	var barangRev []*TopBarangRevenue
	var needRestock []*model.ItemVariant
	var todo *ToDoList
	month := int(pr.Month)
	year := pr.Year
	if sales, e = GetTotalSales(month, year); e == nil {
		if purchase, e = GetTotalPurchase(month, year); e == nil {
			if ret, e = GetTotalSalesReturn(month, year); e == nil {
				if itemJual, e = GetTotalItemTerjual(month, year); e == nil {
					if stockIn, e = GetTotalStockIn(month, year); e == nil {
						if StockRet, e = GetTotalStockReturn(month, year); e == nil {
							if graph, e = GetReportPerDay(month, year); e == nil {
								if barangRev, e = GetGraphicTopRevenue(month, year); e == nil {
									if todo, e = GetTodoList(); e == nil {
										if needRestock, e = GetNeedToRestock(); e == nil {
											m = &DSBData{
												TotalSales:              sales.Total,
												TotalPurchase:           purchase.Total,
												TotalReturn:             ret,
												TotalItemTerjual:        itemJual,
												TotalStockIn:            stockIn,
												TotalStockReturn:        StockRet,
												GraphicReport:           graph,
												GraphicTopBarangRevenue: barangRev,
												ToDoListReport:          todo,
												OverViewReport:          GetOverview(sales, purchase),
												NeedToRestock:           needRestock,
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

//// Get All Single Data /////////////////////////////////////////////////////////////////

// SOAndPOData menampung data jumlah transaksi, total charge dan total paid
type SOAndPOData struct {
	Total  int64
	Charge float64
	Cost   float64
	Paid   float64
}

// GetTotalSales mengambil jumlah total sales pada bulan tertentu
func GetTotalSales(month, year int) (data SOAndPOData, e error) {
	o := orm.NewOrm()
	e = o.Raw("SELECT COUNT(*) AS total, SUM(so.total_charge) AS charge, SUM(so.total_cost) AS cost, SUM(so.total_paid) AS paid FROM sales_order so "+
		"WHERE so.is_deleted = 0 AND so.document_status != 'approved_cancel' "+
		"AND month(so.recognition_date) = ? AND year(so.recognition_date) = ?", month, year).QueryRow(&data)
	return
}

// GetTotalPurchase mengambil jumlah total purchase pada bulan tertentu
func GetTotalPurchase(month, year int) (data SOAndPOData, e error) {
	o := orm.NewOrm()
	e = o.Raw("SELECT COUNT(*) AS total, SUM(po.total_charge) AS charge, SUM(po.total_paid) AS paid FROM purchase_order po "+
		"WHERE po.is_deleted = 0 AND po.document_status != 'cancelled' "+
		"AND month(po.recognition_date) = ? AND year(po.recognition_date) = ?", month, year).QueryRow(&data)
	return
}

// GetTotalSalesReturn mengambil jumlah total sales return pada bulan tertentu
func GetTotalSalesReturn(month, year int) (total int64, e error) {
	o := orm.NewOrm()
	e = o.Raw("SELECT COUNT(*) FROM sales_return sr WHERE sr.is_deleted = 0 AND sr.document_status != 'cancelled' "+
		"AND month(sr.recognition_date) = ? AND year(sr.recognition_date) = ?", month, year).QueryRow(&total)
	return
}

// GetTotalItemTerjual mengambil jumlah total item terjual pada bulan tertentu
func GetTotalItemTerjual(month, year int) (total float32, e error) {
	o := orm.NewOrm()
	e = o.Raw("SELECT SUM(soi.quantity) As total FROM sales_order_item soi INNER JOIN sales_order so ON so.id = soi.sales_order_id "+
		"WHERE so.is_deleted = 0 AND so.document_status != 'approved_cancel' AND month(so.recognition_date) = ? AND year(so.recognition_date) = ?", month, year).QueryRow(&total)
	return
}

// GetTotalStockIn mengambil jumlah total wo receiving quantity pada bulan tertentu
func GetTotalStockIn(month, year int) (total float32, e error) {
	o := orm.NewOrm()
	e = o.Raw("SELECT SUM(woi.quantity) As total FROM workorder_receiving_item woi INNER JOIN workorder_receiving wo ON wo.id = woi.workorder_receiving_id "+
		"WHERE wo.is_deleted = 0 AND month(wo.recognition_date) = ? AND year(wo.recognition_date) = ?", month, year).QueryRow(&total)
	return
}

// GetTotalStockReturn mengambil jumlah total sales return quantity pada bulan tertentu
func GetTotalStockReturn(month, year int) (total float32, e error) {
	o := orm.NewOrm()
	e = o.Raw("SELECT SUM(sri.quantity) As total FROM sales_return_item sri INNER JOIN sales_return sr ON sr.id = sri.sales_return_id "+
		"WHERE sr.is_deleted = 0 AND month(sr.recognition_date) = ? AND year(sr.recognition_date) = ?", month, year).QueryRow(&total)
	return
}

//// Get All Graphic Report Data //////////////////////////////////////////////////////////////

// GetReportPerDay mengambil report sales per hari pada bulan tertentu Tanpa total revenue
func GetReportPerDay(month, year int) (graph []*Graphic, e error) {
	var trans []transact
	var revenue []rev
	var grap []Graphic
	// get graph without transaction and sum revenue
	if grap, e = getSalesOrderSumQty(month, year); e == nil {
		// get jumlah transaksi
		if trans, e = getSalesOrderTransactionAndDate(month, year); e == nil {
			// get sum revenue
			if revenue, e = getRevenueSumAmount(month, year); e == nil {
				// filter all data
				for _, grp := range grap {
					dataGraph := &Graphic{
						Date:     grp.Date,
						ItemSold: grp.ItemSold,
					}
					// cek loop revenue
					for _, rv := range revenue {
						if rv.Date == grp.Date {
							dataGraph.TotalRevenue = rv.Amount
						}
					}
					// cek transaction
					for _, tr := range trans {
						if tr.Date == grp.Date {
							dataGraph.Transaction = tr.Transaction
						}
					}

					graph = append(graph, dataGraph)
				}
			}
		}
	}
	return
}

// transact untuk menyimpan hasil pengambilan transaction getSalesOrderTransactionAndDate
type transact struct {
	Date        time.Time
	Transaction int64
}

// rev untuk menyimpan hasil pengambilan sum revenue amount getRevenueSumAmount
type rev struct {
	Date   time.Time
	Amount float64
}

// getSalesOrderSumQty mengambil jumlah transaksi quantity sales per hari dalam satu bulan
func getSalesOrderSumQty(month, year int) (gr []Graphic, e error) {
	o := orm.NewOrm()
	_, e = o.Raw("SELECT DATE(so.recognition_date) AS date , sum(soi.quantity) as item_sold "+
		"FROM sales_order_item soi INNER JOIN sales_order so ON so.id = soi.sales_order_id "+
		"WHERE so.is_deleted = 0 AND so.document_status != 'approved_cancel' AND year(so.recognition_date) = ? AND month(so.recognition_date) = ? "+
		"GROUP BY DATE(so.recognition_date) ORDER BY date ASC", year, month).QueryRows(&gr)

	return
}

// getSalesOrderTransactionAndDate mengambil jumlah transaksi sales per hari dalam satu bulan
func getSalesOrderTransactionAndDate(month, year int) (trans []transact, e error) {
	o := orm.NewOrm()
	_, e = o.Raw("SELECT DATE(so.recognition_date) AS date,COUNT(so.id) as transaction "+
		"FROM sales_order so WHERE so.is_deleted = 0 AND so.document_status != 'approved_cancel' AND month(so.recognition_date) = ? AND YEAR(so.recognition_date) = ? "+
		"GROUP BY DATE(so.recognition_date) ORDER BY date ASC", month, year).QueryRows(&trans)
	return
}

// getSalesOrderTransactionAndDate mengambil jumlah transaksi sales per hari dalam satu bulan
func getRevenueSumAmount(month, year int) (revenue []rev, e error) {
	o := orm.NewOrm()
	_, e = o.Raw("SELECT DATE(rev.recognition_date) AS date , sum(rev.amount) as amount "+
		"FROM finance_revenue rev WHERE rev.is_deleted = 0 AND month(rev.recognition_date) = ? AND YEAR(rev.recognition_date) = ? "+
		"GROUP BY DATE(rev.recognition_date) ORDER BY date ASC", month, year).QueryRows(&revenue)
	return
}

//// Get All Graphic Top 10 Revenue Item Data //////////////////////////////////////////////////////////////

// GetGraphicTopRevenue mengambil semua data top 10 revenue dengan model item variant
func GetGraphicTopRevenue(month, year int) (topRev []*TopBarangRevenue, e error) {
	var revItem []revPerItem
	if revItem, e = getTopRevenuePerItem(month, year); e == nil {
		for _, top := range revItem {
			// isi barang revenue
			barang := &TopBarangRevenue{
				TotalSubTotal: top.TotalSubtotal,
			}
			// read item variant
			if barang.ItemVariant, e = inventory.GetDetailItemVariantWithoutRelation("id", top.ItemVar); e != nil {
				break
			}
			topRev = append(topRev, barang)
		}
	}
	return
}

// revPerItem menampung data top 10 revenue data tanpa model item variant
type revPerItem struct {
	ItemVar       int64
	TotalSubtotal float64
}

// getTopRevenuePerItem mengambil top 10 revenue
func getTopRevenuePerItem(month, year int) (revItem []revPerItem, e error) {
	o := orm.NewOrm()
	_, e = o.Raw("SELECT soi.item_variant_id As item_var, SUM(soi.subtotal) AS total_subtotal FROM sales_order_item soi "+
		"INNER JOIN sales_order so ON so.id = soi.sales_order_id "+
		"WHERE so.document_status in ('active','finished','requested_cancel') AND so.is_deleted = 0 "+
		"AND month(so.recognition_date) = ? AND YEAR(so.recognition_date) = ? "+
		"GROUP BY item_variant_id ORDER BY total_subtotal DESC LIMIT 0,5", month, year).QueryRows(&revItem)
	return
}

//// Get Data To Do List ///////////////////////////////////////////////////////////////////////////////////

// GetTodoList mendapatkan data  todolist
func GetTodoList() (todo *ToDoList, e error) {
	var salesNeedFulfill, salesNeedInv, salesNeedCancel, fullfillNeedShip, ShipmentNeedDeliv int64
	var itemNeedReceiv float32
	// get all data ToDoList
	if salesNeedFulfill, e = getCountSalesNeedFulfill(); e == nil {
		if salesNeedInv, e = getCountSalesNeedInvoiced(); e == nil {
			if salesNeedCancel, e = getCountSalesNeedApproveCancel(); e == nil {
				if fullfillNeedShip, e = getCountFullfillNeedShipped(); e == nil {
					if itemNeedReceiv, e = getSumItemNeedReceived(); e == nil {
						if ShipmentNeedDeliv, e = getCountShipmentNeedDeliv(); e == nil {
							todo = &ToDoList{
								SalesNeedFulfilled:     salesNeedFulfill,
								SalesNeedInvoiced:      salesNeedInv,
								SalesNeedApproveCancel: salesNeedCancel,
								FulfillmentNeedShipped: fullfillNeedShip,
								ItemNeedReceived:       itemNeedReceiv,
								ShipmentNeedDelivered:  ShipmentNeedDeliv,
							}
						}
					}
				}
			}
		}
	}
	return
}

// getCountSalesNeedFulfill mengambil jumlah sales yang fulfillment belum finish
func getCountSalesNeedFulfill() (total int64, e error) {
	o := orm.NewOrm()
	e = o.Raw("SELECT COUNT(*) FROM sales_order so WHERE so.is_deleted = 0 " +
		"AND so.document_status != 'approved_cancel' AND so.fulfillment_status != 'finished'").QueryRow(&total)
	return
}

// getCountSalesNeedInvoiced mengambil jumlah sales yang invoice belum finish
func getCountSalesNeedInvoiced() (total int64, e error) {
	o := orm.NewOrm()
	e = o.Raw("SELECT COUNT(*) FROM sales_order so WHERE so.is_deleted = 0 " +
		"AND so.document_status != 'approved_cancel' AND so.invoice_status != 'finished'").QueryRow(&total)
	return
}

// getCountSalesNeedApproveCancel mengambil jumlah sales yang belum approve cancel
func getCountSalesNeedApproveCancel() (total int64, e error) {
	o := orm.NewOrm()
	e = o.Raw("SELECT COUNT(*) FROM sales_order so WHERE so.is_deleted = 0 " +
		"AND so.document_status = 'requested_cancel'").QueryRow(&total)
	return
}

// getCountFullfillNeedShipped mengambil jumlah fullfillment yang perlu shipping
func getCountFullfillNeedShipped() (total int64, e error) {
	o := orm.NewOrm()
	e = o.Raw("SELECT COUNT(*) FROM workorder_fulfillment wf WHERE wf.is_deleted = 0 " +
		"AND wf.is_delivered = 0").QueryRow(&total)
	return
}

// getSumItemNeedReceived mengambil jumlah item yang belum diterima
func getSumItemNeedReceived() (total float32, e error) {
	o := orm.NewOrm()
	e = o.Raw("SELECT ( SUM(poi.quantity)-SUM(wri.quantity)) As total FROM purchase_order_item poi " +
		"INNER JOIN purchase_order po ON po.id = poi.purchase_order_id " +
		"JOIN workorder_receiving_item wri ON wri.purchase_order_item_id = poi.id " +
		"JOIN workorder_receiving wr on wr.id = wri.workorder_receiving_id " +
		"WHERE po.is_deleted = 0 AND po.document_status != 'cancelled' AND wr.document_status = 'finished' " +
		"AND wr.is_deleted = 0").QueryRow(&total)
	return
}

// getCountShipmentNeedDeliv mengambil jumlah shipment yang belum di kirim
func getCountShipmentNeedDeliv() (total int64, e error) {
	o := orm.NewOrm()
	e = o.Raw("SELECT COUNT(*) FROM workorder_shipment ws WHERE ws.is_deleted = 0 " +
		"AND ws.document_status != 'finished' ").QueryRow(&total)
	return
}

//// Get Data OverView /////////////////////////////////////////////////////////////////////////////////////

// GetOverview untuk mendapatkan jumlah overview
func GetOverview(sales, purchase SOAndPOData) (overview *OverView) {
	overview = &OverView{
		TotalPenjualan: sales.Charge,
		TotalCost:      sales.Cost,
		TotalProfit:    sales.Charge - sales.Cost,
		PiutangUsaha:   sales.Charge - sales.Paid,
		HutangUsaha:    purchase.Charge - purchase.Paid,
	}
	return
}

//// Get Data NeedToRestock ////////////////////////////////////////////////////////////////////////////////

// GetNeedToRestock mengambil item variant yang stock available <= stock minimum
func GetNeedToRestock() (m []*model.ItemVariant, e error) {
	// get all id
	var itemVarID []int64
	o := orm.NewOrm()
	if _, e = o.Raw("SELECT iv.id FROM item_variant iv " +
		"WHERE iv.available_stock <= iv.minimum_stock AND iv.is_deleted= 0 limit 20").QueryRows(&itemVarID); e == nil {
		// get detail item variant
		if len(itemVarID) > int(0) {
			m, e = getDetailAllItemVariant("id__in", itemVarID)
		}
	}

	return
}

// getDetailAllItemVariant mengambil detail dari semua item variant sesuai param, menggunakan All
func getDetailAllItemVariant(field string, values ...interface{}) ([]*model.ItemVariant, error) {
	itm := new(model.ItemVariant)
	o := orm.NewOrm().QueryTable(itm)
	var itmVar []*model.ItemVariant
	if _, e := o.Filter(field, values...).RelatedSel().All(&itmVar); e != nil {
		return nil, e
	}
	return itmVar, nil
}

//// Get Param In URL //////////////////////////////////////////////////////////////////////////////////////

// ParamURL untuk menyimpan hasil pengambilan param url
type ParamURL struct {
	Month time.Month
	Year  int
}

// GetParamURL untuk mendapatkan parameter pada url sesuai dengan kriteria:
// filter by month and year ,example: /url?month=6&year=2012
func GetParamURL(ctx *cuxs.Context) (param ParamURL) {
	// take url value
	urlVal := ctx.QueryParams()
	// get certain value
	year := urlVal.Get("year")
	month := urlVal.Get("month")
	// insert the value
	param = filterParam(month, year)
	return
}

// filterParam mengecek nilai month dan year dari param
func filterParam(month, year string) (param ParamURL) {
	// check year
	if year == "" {
		param.Year = time.Now().Year()
	} else {
		if yearVal := common.ToInt(year); yearVal <= int(0) {
			param.Year = time.Now().Year()
		} else {
			param.Year = yearVal
		}
	}

	// check month
	if month == "" {
		param.Month = time.Now().Month()
	} else {
		if monthVal := common.ToInt(month); monthVal < int(1) || monthVal > int(12) {
			param.Month = time.Now().Month()
		} else {
			param.Month = time.Month(monthVal)
		}
	}
	return
}
