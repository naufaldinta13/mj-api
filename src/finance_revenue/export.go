// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package financeRevenue

import (
	"fmt"
	"git.qasico.com/cuxs/env"
	"git.qasico.com/cuxs/mailer"
	"git.qasico.com/cuxs/orm"
	"git.qasico.com/mj/api/datastore/model"
	"github.com/labstack/gommon/log"
	"github.com/tealeg/xlsx"
	"strings"
	"time"
)

// Cron mengirimkan email laporan pendapatan bulanan
func Cron() {
	rq := &orm.RequestQuery{
		Offset:     0,
		Limit:      -1,
		Conditions: make([]map[string]string, 0),
	}

	if ds, t, e := getData(rq); t > 0 && e == nil {
		if f, e := xls(ds); e == nil {
			// kirim email
			sendMail(f)
		}
	}
}

func getData(rq *orm.RequestQuery) (m []*model.FinanceRevenue, total int64, err error) {
	// make new orm query
	q, _ := rq.Query(new(model.FinanceRevenue))
	q = q.Filter("is_deleted", 0)

	// last month
	// fungsi ini berjalan pada tanggal 1 bulan baru
	// jadi kita ambil bulan sebelumnya sampai sekarang
	n := time.Now()
	ld := n.AddDate(0, -1, 0)

	q = q.Filter("recognition_date__lt", n.Format("2006-01-02")).Filter("recognition_date__gte", ld.Format("2006-01-02"))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.FinanceRevenue
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func xls(ds []*model.FinanceRevenue) (fileDir string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := env.GetString("EXPORT_DIRECTORY", "")

	now := time.Now()
	filename := fmt.Sprintf("LaporanPendapatan-%s.xlsx", now.Format("200601021504"))
	fileDir = fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Laporan Pendapatan"
		row = sheet.AddRow()
		row.AddCell().Value = ""
		row = sheet.AddRow()
		row.AddCell().Value = ""

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Tanggal"
		row.AddCell().Value = "Transaksi"
		row.AddCell().Value = "Kode Transaksi"
		row.AddCell().Value = "Kode Order"
		row.AddCell().Value = "Customer"
		row.AddCell().Value = "Total Bayar"
		row.AddCell().Value = "Payment Method"
		row.AddCell().Value = "Bank"

		o := orm.NewOrm()
		var total float64
		for i, d := range ds {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = d.RecognitionDate.Format("02/01/2006")
			row.AddCell().Value = strings.Replace(d.RefType, "_", " ", 1)

			var docCode, docCustomer, docType string
			if d.RefType == "sales_invoice" {
				o.Raw("select * from sales_invoice where id = ?", int64(d.RefID)).QueryRow(&d.SalesInvoice)
				if d.SalesInvoice != nil {
					d.SalesInvoice.SalesOrder.Read()
					docType = d.SalesInvoice.Code
					docCode = d.SalesInvoice.SalesOrder.Code
					if d.SalesInvoice.SalesOrder.Customer != nil {
						d.SalesInvoice.SalesOrder.Customer.Read()
						docCustomer = d.SalesInvoice.SalesOrder.Customer.FullName
					} else {
						docCustomer = "Walkin Customer"
					}
				}
			}

			if d.RefType == "purchase_return" {
				o.Raw("select * from purchase_return where id = ?", int64(d.RefID)).QueryRow(&d.PurchaseReturn)
				if d.PurchaseReturn != nil {
					d.PurchaseReturn.PurchaseOrder.Read()
					docType = d.PurchaseReturn.Code
					docCode = d.PurchaseReturn.PurchaseOrder.Code
					d.PurchaseReturn.PurchaseOrder.Supplier.Read()
					docCustomer = d.PurchaseReturn.PurchaseOrder.Supplier.FullName
				}
			}

			row.AddCell().Value = docType
			row.AddCell().Value = docCode
			row.AddCell().Value = docCustomer
			row.AddCell().SetFloat(d.Amount)
			row.AddCell().Value = d.PaymentMethod
			row.AddCell().Value = d.BankName

			total += d.Amount
		}
		row = sheet.AddRow()
		row.AddCell().Value = ""

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "TOTAL"
		row.AddCell().Value = ""
		row.AddCell().Value = ""
		row.AddCell().Value = ""
		row.AddCell().Value = ""
		row.AddCell().Value = ""
		row.AddCell().SetFloat(total)

		err = file.Save(fileDir)
		log.Error(err)
	}

	return fileDir, err
}

func sendMail(file string) {
	m := mailer.NewMessage()
	email := env.GetString("EXPORT_EMAIL", "adam@qasico.com")

	n := time.Now()
	ld := n.AddDate(0, -1, 0)

	m.SetRecipient(email)
	m.SetAddressHeader("Bcc", "james@qasico.com", "James")
	m.SetAddressHeader("Bcc", "adam@qasico.com", "Adam")

	s := fmt.Sprintf("Laporan Pendapatan (%s - %s)", ld.Format("2006/01/02"), n.Format("2006/01/02"))

	m.SetSubject(s)
	m.SetBody("text/html", "Berikut adalah "+s)
	m.Attach(file)

	// initialing smtp dialer
	d := mailer.NewDialer()
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
