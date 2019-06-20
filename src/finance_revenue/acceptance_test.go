// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package financeRevenue_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/test"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/common/tester"
	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.Setup()

	// run tests
	res := m.Run()

	// cleanup
	//test.DataCleanUp()

	os.Exit(res)
}

func TestARouting(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	fr := model.DummyFinanceRevenue()
	fr.IsDeleted = 0
	fr.Save()

	efr := common.Encrypt(fr.ID)

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/finance-revenue", "GET", 200},
		{"/v1/finance-revenue/" + efr, "GET", 200},
		{"/v1/finance-revenue/" + efr, "PUT", 415},
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

	fr.IsDeleted = 1
	fr.Save()
	routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/finance-revenue", "GET", 200},
		{"/v1/finance-revenue/" + efr, "GET", 404},
		{"/v1/finance-revenue/" + efr, "PUT", 404},
	}

	for _, ep := range routers {
		ng.Method = ep.method
		ng.Path = ep.endpoint
		ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, ep.expected, res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", ep.endpoint, ep.method))
		})
	}
}

// tes membuat finance revenue dengan ref_type nya sales invoice
func TestHandler_URLMappingCreateFinanceRevenueSalesInvoiceSuccess(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	orm.NewOrm().Raw("DELETE from invoice_receipt").Exec()

	// buat dummy sales invoice

	si := model.DummySalesInvoice()
	si.DocumentStatus = "new"
	si.TotalAmount = 1000000
	si.IsDeleted = 0
	si.Save()

	customer := model.DummyPartnership()

	so := si.SalesOrder
	so.DocumentStatus = "active"
	so.Customer = customer
	so.TotalPaid = float64(40000)
	so.TotalCharge = float64(50000)
	so.InvoiceStatus = "active"
	so.ShipmentStatus = "finished"
	so.FulfillmentStatus = "finished"
	so.IsDeleted = int8(0)
	so.Save()

	customer.PartnershipType = "customer"
	customer.IsDeleted = int8(0)
	customer.TotalDebt = float64(20000)
	customer.Save()

	data := tester.D{"ref_id": common.Encrypt(si.ID), "ref_type": "sales_invoice", "recognition_date": time.Now(), "payment_method": "credit_card", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

	si.Read("ID")
	// cek document status pada sales invoice, harusnya menjadi active
	assert.Equal(t, "active", si.DocumentStatus)

	// cek total revenued pada sales invoice
	assert.Equal(t, float64(50000), si.TotalRevenued)

}

// tes fail membuat finance revenue dengan ref_type nya sales invoice not found
func TestHandler_URLMappingCreateFinanceRevenueSalesInvoiceNotFound(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy sales invoice
	si := model.DummySalesInvoice()
	si.DocumentStatus = "new"
	si.TotalAmount = 1000000
	si.Save()

	si.Delete()

	data := tester.D{"ref_id": common.Encrypt(si.ID), "ref_type": "sales_invoice", "recognition_date": time.Now(), "payment_method": "credit_card", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

}

// tes membuat finance revenue dengan ref_type nya purchase return
func TestHandler_URLMappingCreateFinanceRevenuePurchaseReturnSuccess(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy purchase return
	pr := model.DummyPurchaseReturn()
	pr.DocumentStatus = "new"
	pr.TotalAmount = 1000000
	pr.Save()

	data := tester.D{"ref_id": common.Encrypt(pr.ID), "ref_type": "purchase_return", "recognition_date": time.Now(), "payment_method": "credit_card", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "purchase return"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

	pr.Read("ID")
	// cek document status pada purchase return, harusnya menjadi active
	assert.Equal(t, "active", pr.DocumentStatus)

}

// tes fail membuat finance revenue dengan ref_type nya purchase return
func TestHandler_URLMappingCreateFinanceRevenuePurchaseReturnFailNotFound(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy purchase return
	pr := model.DummyPurchaseReturn()
	pr.DocumentStatus = "new"
	pr.TotalAmount = 1000000
	pr.Save()

	pr.Delete()

	data := tester.D{"ref_id": common.Encrypt(pr.ID), "ref_type": "purchase_return", "recognition_date": time.Now(), "payment_method": "credit_card", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "purchase return"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

}

// tes gagal membuat finance revenue karena document status pada sales invoice sudah finish
// gagal karena document status pada sales invoice sudah finish
func TestHandler_URLMappingCreateFinanceRevenueSalesInvoiceFailDocumentStatus(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy sales invoice
	si := model.DummySalesInvoice()
	si.DocumentStatus = "finished"
	si.TotalAmount = 1000000
	si.Save()

	data := tester.D{"ref_id": common.Encrypt(si.ID), "ref_type": "sales_invoice", "recognition_date": time.Now(), "payment_method": "credit_card", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

}

// tes gagal membuat finance revenue dengan ref_type nya purchase return
// gagal karena document status pada purchase return sudah finish
func TestHandler_URLMappingCreateFinanceRevenuePurchaseReturnFailDocumentStatus(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy purchase return
	pr := model.DummyPurchaseReturn()
	pr.DocumentStatus = "finished"
	pr.TotalAmount = 1000000
	pr.Save()

	data := tester.D{"ref_id": common.Encrypt(pr.ID), "ref_type": "purchase_return", "recognition_date": time.Now(), "payment_method": "credit_card", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "purchase return"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

}

// tes membuat finance revenue dengan payment debit card
func TestHandler_URLMappingCreateFinanceRevenueSuccessPaymentDebitCard(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy sales invoice
	si := model.DummySalesInvoice()
	si.DocumentStatus = "new"
	si.TotalAmount = 1000000
	si.IsDeleted = 0
	si.Save()

	customer := model.DummyPartnership()

	so := si.SalesOrder
	so.DocumentStatus = "active"
	so.Customer = customer
	so.TotalPaid = float64(40000)
	so.TotalCharge = float64(50000)
	so.InvoiceStatus = "active"
	so.ShipmentStatus = "finished"
	so.FulfillmentStatus = "finished"
	so.IsDeleted = int8(0)
	so.Save()

	customer.PartnershipType = "customer"
	customer.IsDeleted = int8(0)
	customer.TotalDebt = float64(20000)
	customer.Save()

	data := tester.D{"ref_id": common.Encrypt(si.ID), "ref_type": "sales_invoice", "recognition_date": time.Now(), "payment_method": "debit_card", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

	si.Read("ID")
	// cek document status pada sales invoice, harusnya menjadi active
	assert.Equal(t, "active", si.DocumentStatus)

}

// tes membuat finance revenue dengan payment giro
func TestHandler_URLMappingCreateFinanceRevenueSuccessPaymentGiro(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy sales invoice
	si := model.DummySalesInvoice()
	si.DocumentStatus = "new"
	si.TotalAmount = 1000000
	si.IsDeleted = 0
	si.Save()

	customer := model.DummyPartnership()

	so := si.SalesOrder
	so.DocumentStatus = "active"
	so.Customer = customer
	so.TotalPaid = float64(40000)
	so.TotalCharge = float64(50000)
	so.InvoiceStatus = "active"
	so.ShipmentStatus = "finished"
	so.FulfillmentStatus = "finished"
	so.IsDeleted = int8(0)
	so.Save()

	customer.PartnershipType = "customer"
	customer.IsDeleted = int8(0)
	customer.TotalDebt = float64(20000)
	customer.Save()

	data := tester.D{"ref_id": common.Encrypt(si.ID), "ref_type": "sales_invoice", "recognition_date": time.Now(), "payment_method": "giro", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

	si.Read("ID")
	// cek document status pada sales invoice, harusnya menjadi active
	assert.Equal(t, "active", si.DocumentStatus)

}

// tes membuat finance revenue dengan payment cash
func TestHandler_URLMappingCreateFinanceRevenueSuccessPaymentCash(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy sales invoice
	si := model.DummySalesInvoice()
	si.DocumentStatus = "new"
	si.TotalAmount = 1000000
	si.IsDeleted = 0
	si.Save()

	customer := model.DummyPartnership()

	so := si.SalesOrder
	so.DocumentStatus = "active"
	so.Customer = customer
	so.TotalPaid = float64(40000)
	so.TotalCharge = float64(50000)
	so.InvoiceStatus = "active"
	so.ShipmentStatus = "finished"
	so.FulfillmentStatus = "finished"
	so.IsDeleted = int8(0)
	so.Save()

	customer.PartnershipType = "customer"
	customer.IsDeleted = int8(0)
	customer.TotalDebt = float64(20000)
	customer.Save()

	data := tester.D{"ref_id": common.Encrypt(si.ID), "ref_type": "sales_invoice", "recognition_date": time.Now(), "payment_method": "cash", "amount": 50000, "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

	si.Read("ID")
	// cek document status pada sales invoice, harusnya menjadi active
	assert.Equal(t, "active", si.DocumentStatus)

}

// tes fail membuat finance revenue dengan payment credit card
func TestHandler_URLMappingCreateFinanceRevenueFailPaymentCreditCard(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy sales invoice
	si := model.DummySalesInvoice()
	si.DocumentStatus = "new"
	si.TotalAmount = 1000000
	si.IsDeleted = 0
	si.Save()

	// sukses
	data := tester.D{"ref_id": common.Encrypt(si.ID), "ref_type": "sales_invoice", "recognition_date": time.Now(), "payment_method": "credit_card", "amount": 50000, "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

}

// tes fail membuat finance revenue dengan payment debit card
func TestHandler_URLMappingCreateFinanceFailPaymentDebitCard(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy sales invoice
	si := model.DummySalesInvoice()
	si.DocumentStatus = "new"
	si.TotalAmount = 1000000
	si.IsDeleted = 0
	si.Save()

	data := tester.D{"ref_id": common.Encrypt(si.ID), "ref_type": "sales_invoice", "recognition_date": time.Now(), "payment_method": "debit_card", "amount": 50000, "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

}

// tes fail membuat finance revenue dengan payment giro
func TestHandler_URLMappingCreateFinanceRevenueFailPaymentGiro(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy sales invoice
	si := model.DummySalesInvoice()
	si.DocumentStatus = "new"
	si.TotalAmount = 1000000
	si.IsDeleted = 0
	si.Save()

	data := tester.D{"ref_id": common.Encrypt(si.ID), "ref_type": "sales_invoice", "recognition_date": time.Now(), "payment_method": "giro", "amount": 50000,
		"bank_name": "Tes", "bank_holder": "Tes", "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

}

// tes fail membuat finance revenue dengan payment cash
func TestHandler_URLMappingCreateFinanceRevenueFailPaymentCash(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy sales invoice
	si := model.DummySalesInvoice()
	si.DocumentStatus = "new"
	si.TotalAmount = 1000000
	si.IsDeleted = 0
	si.Save()

	data := tester.D{"ref_id": common.Encrypt(si.ID), "ref_type": "sales_invoice", "recognition_date": time.Now(), "payment_method": "cash", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

}

// tes gagal membuat finance revenue dengan ref_type nya sales invoice
// karena jumlah bayar lebih dari total amount sales invoice
func TestHandler_URLMappingCreateFinanceRevenueSalesInvoiceFailAmount(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy sales invoice
	si := model.DummySalesInvoice()
	si.DocumentStatus = "new"
	si.TotalAmount = 40000
	si.IsDeleted = 0
	si.Save()

	data := tester.D{"ref_id": common.Encrypt(si.ID), "ref_type": "sales_invoice", "recognition_date": time.Now(), "payment_method": "credit_card", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

}

// tes membuat finance revenue dengan ref_type nya purchase return
// karena jumlah bayar lebih dari total amount purchase return
func TestHandler_URLMappingCreateFinanceRevenuePurchaseReturnFailAmount(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy purchase return
	pr := model.DummyPurchaseReturn()
	pr.DocumentStatus = "new"
	pr.TotalAmount = 40000
	pr.Save()

	data := tester.D{"ref_id": common.Encrypt(pr.ID), "ref_type": "purchase_return", "recognition_date": time.Now(), "payment_method": "credit_card", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "purchase return"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

}

// tes update finance revenue dengan ref_type sales invoice
func TestHandler_URLMappingUpdateFinanceRevenueSalesInvoiceSuccess(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	orm.NewOrm().Raw("delete from invoice_receipt").Exec()
	orm.NewOrm().Raw("delete from finance_revenue").Exec()

	si := model.DummySalesInvoice()
	si.TotalAmount = 100000
	si.IsDeleted = 0
	si.Save()

	// buat dummy finance revenue
	fr := model.DummyFinanceRevenue()
	fr.RefID = uint64(si.ID)
	fr.RefType = "sales_invoice"
	fr.DocumentStatus = "uncleared"
	fr.IsDeleted = 0
	fr.Save()

	efr := common.Encrypt(fr.ID)

	data := tester.D{"recognition_date": time.Now(), "payment_method": "debit_card", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-revenue/"+efr).SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

	si.Read("ID")
	// cek total revenued pada sales invoice
	assert.Equal(t, float64(50000), si.TotalRevenued)
}

// tes update finance revenue dengan payment debit card
func TestHandler_URLMappingUpdateFinanceRevenueSuccessPaymentDebitCard(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	test.DbClean("finance_revenue")

	pr := model.DummyPurchaseReturn()
	pr.TotalAmount = 100000
	pr.Save()

	// buat dummy finance revenue
	fr := model.DummyFinanceRevenue()
	fr.RefID = uint64(pr.ID)
	fr.RefType = "purchase_return"
	fr.DocumentStatus = "uncleared"
	fr.IsDeleted = 0
	fr.Save()

	efr := common.Encrypt(fr.ID)

	data := tester.D{"recognition_date": time.Now(), "payment_method": "debit_card", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-revenue/"+efr).SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})
}

// tes update finance revenue jika total amount finance revenue > purchase return totalAmount
func TestHandler_URLMappingUpdateFinanceRevenueFailAmountPurchaseReturn(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	test.DbClean("finance_revenue")

	pr := model.DummyPurchaseReturn()
	pr.TotalAmount = 10000
	pr.Save()

	// buat dummy finance revenue
	fr := model.DummyFinanceRevenue()
	fr.RefID = uint64(pr.ID)
	fr.RefType = "purchase_return"
	fr.DocumentStatus = "uncleared"
	fr.IsDeleted = 0
	fr.Save()

	efr := common.Encrypt(fr.ID)

	data := tester.D{"recognition_date": time.Now(), "payment_method": "debit_card", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-revenue/"+efr).SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})
}

// tes update finance revenue jika total amount finance revenue > sales invoice totalAmount
func TestHandler_URLMappingUpdateFinanceRevenueFailAmountSalesInvoice(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	test.DbClean("finance_revenue")

	si := model.DummySalesInvoice()
	si.TotalAmount = 10000
	si.IsDeleted = 0
	si.Save()

	// buat dummy finance revenue
	fr := model.DummyFinanceRevenue()
	fr.RefID = uint64(si.ID)
	fr.RefType = "sales_invoice"
	fr.DocumentStatus = "uncleared"
	fr.IsDeleted = 0
	fr.Save()

	efr := common.Encrypt(fr.ID)

	data := tester.D{"recognition_date": time.Now(), "payment_method": "debit_card", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-revenue/"+efr).SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})
}

func TestHandler_URLMappingUpdateFinanceRevenueFailPaymentCreditCard(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	test.DbClean("finance_revenue")

	si := model.DummySalesInvoice()
	si.TotalAmount = 10000
	si.IsDeleted = 0
	si.Save()

	// buat dummy finance revenue
	fr := model.DummyFinanceRevenue()
	fr.RefID = uint64(si.ID)
	fr.RefType = "sales_invoice"
	fr.DocumentStatus = "uncleared"
	fr.IsDeleted = 0
	fr.Save()

	efr := common.Encrypt(fr.ID)

	data := tester.D{"ref_id": common.Encrypt(si.ID), "ref_type": "sales_invoice", "recognition_date": time.Now(), "payment_method": "credit_card", "amount": 50000, "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-revenue/"+efr).SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

}

// tes fail update finance revenue dengan payment cash
func TestHandler_URLMappingUpdateFinanceRevenueFailPaymentCash(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	test.DbClean("finance_revenue")

	si := model.DummySalesInvoice()
	si.TotalAmount = 10000
	si.IsDeleted = 0
	si.Save()

	// buat dummy finance revenue
	fr := model.DummyFinanceRevenue()
	fr.RefID = uint64(si.ID)
	fr.RefType = "sales_invoice"
	fr.DocumentStatus = "uncleared"
	fr.IsDeleted = 0
	fr.Save()

	efr := common.Encrypt(fr.ID)

	data := tester.D{"ref_id": common.Encrypt(si.ID), "ref_type": "sales_invoice", "recognition_date": time.Now(), "payment_method": "cash", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-revenue/"+efr).SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

}

// tes fail update finance revenue dengan payment giro
func TestHandler_URLMappingUpdateFinanceRevenueFailPaymentGiro(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	si := model.DummySalesInvoice()
	si.TotalAmount = 10000
	si.IsDeleted = 0
	si.Save()

	// buat dummy finance revenue
	fr := model.DummyFinanceRevenue()
	fr.RefID = uint64(si.ID)
	fr.RefType = "sales_invoice"
	fr.DocumentStatus = "uncleared"
	fr.IsDeleted = 0
	fr.Save()

	efr := common.Encrypt(fr.ID)

	data := tester.D{"recognition_date": time.Now(), "payment_method": "giro", "amount": 50000, "bank_name": "Tes", "bank_holder": "Tes", "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-revenue/"+efr).SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})
}

// tes fail update finance revenue dengan document status finance revenue adalah cleared
func TestHandler_URLMappingUpdateFinanceRevenueFailDocumentStatus(t *testing.T) {

	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	si := model.DummySalesInvoice()
	si.TotalAmount = 10000
	si.Save()

	// buat dummy finance revenue
	fr := model.DummyFinanceRevenue()
	fr.RefID = uint64(si.ID)
	fr.RefType = "sales_invoice"
	fr.DocumentStatus = "cleared"
	fr.IsDeleted = 0
	fr.Save()

	efr := common.Encrypt(fr.ID)

	data := tester.D{"recognition_date": time.Now(), "payment_method": "giro", "amount": 50000, "bank_name": "Tes", "bank_holder": "Tes", "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-revenue/"+efr).SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})
}

//// approve finance revenue ///////////////////////////////////////////////////////////////////

// TestHandler_URLMappingApproveFinanceRevenueSuccess1 berhasil dengan sales invoice,finish
func TestHandler_URLMappingApproveFinanceRevenueSuccess1(t *testing.T) {
	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy
	customer := model.DummyPartnership()

	inv := model.DummySalesInvoice()
	inv.TotalAmount = float64(50000)
	inv.TotalPaid = float64(40000)
	inv.DocumentStatus = "active"
	inv.IsDeleted = int8(0)
	inv.Save()

	so := inv.SalesOrder
	so.DocumentStatus = "active"
	so.Customer = customer
	so.TotalPaid = float64(40000)
	so.TotalCharge = float64(50000)
	so.InvoiceStatus = "active"
	so.ShipmentStatus = "finished"
	so.FulfillmentStatus = "finished"
	so.IsDeleted = int8(0)
	so.Save()

	customer.PartnershipType = "customer"
	customer.IsDeleted = int8(0)
	customer.TotalDebt = float64(20000)
	customer.Save()

	rev1 := model.DummyFinanceRevenue()
	rev1.IsDeleted = int8(0)
	rev1.DocumentStatus = "cleared"
	rev1.Amount = float64(20000)
	rev1.RefType = "sales_invoice"
	rev1.RefID = uint64(inv.ID)
	rev1.Save()

	rev2 := model.DummyFinanceRevenue()
	rev2.IsDeleted = int8(0)
	rev2.DocumentStatus = "cleared"
	rev2.Amount = float64(20000)
	rev2.RefType = "sales_invoice"
	rev2.RefID = uint64(inv.ID)
	rev2.Save()

	rev := model.DummyFinanceRevenue()
	rev.IsDeleted = int8(0)
	rev.DocumentStatus = "uncleared"
	rev.Amount = float64(10000)
	rev.RefType = "sales_invoice"
	rev.RefID = uint64(inv.ID)
	rev.Save()

	data := tester.D{}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-revenue/"+common.Encrypt(rev.ID)+"/approve").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})
	rev.Read("ID")
	assert.Equal(t, "cleared", rev.DocumentStatus)

	inv.Read("ID")
	assert.Equal(t, "finished", inv.DocumentStatus)
	assert.Equal(t, float64(50000), inv.TotalPaid)

	so.Read("ID")
	assert.Equal(t, "finished", so.InvoiceStatus)
	assert.Equal(t, float64(50000), so.TotalPaid)

	customer.Read("ID")
	assert.Equal(t, float64(0), customer.TotalDebt)
}

// TestHandler_URLMappingApproveFinanceRevenueSuccess2 berhasil dengan sales invoice,active
func TestHandler_URLMappingApproveFinanceRevenueSuccess2(t *testing.T) {
	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy
	customer := model.DummyPartnership()

	inv := model.DummySalesInvoice()
	inv.TotalAmount = float64(50000)
	inv.TotalPaid = float64(40000)
	inv.DocumentStatus = "active"
	inv.IsDeleted = int8(0)
	inv.Save()

	so := inv.SalesOrder
	so.DocumentStatus = "active"
	so.Customer = customer
	so.TotalPaid = float64(40000)
	so.TotalCharge = float64(50000)
	so.InvoiceStatus = "active"
	so.ShipmentStatus = "finished"
	so.FulfillmentStatus = "finished"
	so.IsDeleted = int8(0)
	so.Save()

	customer.PartnershipType = "customer"
	customer.IsDeleted = int8(0)
	customer.TotalDebt = float64(20000)
	customer.Save()

	rev1 := model.DummyFinanceRevenue()
	rev1.IsDeleted = int8(0)
	rev1.DocumentStatus = "cleared"
	rev1.Amount = float64(20000)
	rev1.RefType = "sales_invoice"
	rev1.RefID = uint64(inv.ID)
	rev1.Save()

	rev2 := model.DummyFinanceRevenue()
	rev2.IsDeleted = int8(0)
	rev2.DocumentStatus = "cleared"
	rev2.Amount = float64(20000)
	rev2.RefType = "sales_invoice"
	rev2.RefID = uint64(inv.ID)
	rev2.Save()

	rev := model.DummyFinanceRevenue()
	rev.IsDeleted = int8(0)
	rev.DocumentStatus = "uncleared"
	rev.Amount = float64(5000)
	rev.RefType = "sales_invoice"
	rev.RefID = uint64(inv.ID)
	rev.Save()

	data := tester.D{}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-revenue/"+common.Encrypt(rev.ID)+"/approve").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})
	rev.Read("ID")
	assert.Equal(t, "cleared", rev.DocumentStatus)

	inv.Read("ID")
	assert.Equal(t, "active", inv.DocumentStatus)
	assert.Equal(t, float64(45000), inv.TotalPaid)

	so.Read("ID")
	assert.Equal(t, "active", so.DocumentStatus)
	assert.Equal(t, "active", so.InvoiceStatus)
	assert.Equal(t, float64(45000), so.TotalPaid)

	customer.Read("ID")
	assert.Equal(t, float64(5000), customer.TotalDebt)
}

// TestHandler_URLMappingApproveFinanceRevenueSuccess3 berhasil dengan purchase return,finish
func TestHandler_URLMappingApproveFinanceRevenueSuccess3(t *testing.T) {
	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy
	ret := model.DummyPurchaseReturn()
	ret.TotalAmount = float64(50000)
	ret.DocumentStatus = "active"
	ret.IsDeleted = int8(0)
	ret.Save()

	rev1 := model.DummyFinanceRevenue()
	rev1.IsDeleted = int8(0)
	rev1.DocumentStatus = "cleared"
	rev1.Amount = float64(15000)
	rev1.RefType = "purchase_return"
	rev1.RefID = uint64(ret.ID)
	rev1.Save()

	rev2 := model.DummyFinanceRevenue()
	rev2.IsDeleted = int8(0)
	rev2.DocumentStatus = "cleared"
	rev2.Amount = float64(30000)
	rev2.RefType = "purchase_return"
	rev2.RefID = uint64(ret.ID)
	rev2.Save()

	rev3 := model.DummyFinanceRevenue()
	rev3.IsDeleted = int8(1)
	rev3.DocumentStatus = "cleared"
	rev3.Amount = float64(20000)
	rev3.RefType = "purchase_return"
	rev3.RefID = uint64(ret.ID)
	rev3.Save()

	rev := model.DummyFinanceRevenue()
	rev.IsDeleted = int8(0)
	rev.DocumentStatus = "uncleared"
	rev.Amount = float64(5000)
	rev.RefType = "purchase_return"
	rev.RefID = uint64(ret.ID)
	rev.Save()

	data := tester.D{}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-revenue/"+common.Encrypt(rev.ID)+"/approve").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})
	rev.Read("ID")
	assert.Equal(t, "cleared", rev.DocumentStatus)

	ret.Read("ID")
	assert.Equal(t, "finished", ret.DocumentStatus)
}

// TestHandler_URLMappingApproveFinanceRevenueSuccess4 berhasil dengan purchase return,active
func TestHandler_URLMappingApproveFinanceRevenueSuccess4(t *testing.T) {
	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy
	ret := model.DummyPurchaseReturn()
	ret.TotalAmount = float64(50000)
	ret.DocumentStatus = "active"
	ret.IsDeleted = int8(0)
	ret.Save()

	rev1 := model.DummyFinanceRevenue()
	rev1.IsDeleted = int8(0)
	rev1.DocumentStatus = "cleared"
	rev1.Amount = float64(15000)
	rev1.RefType = "purchase_return"
	rev1.RefID = uint64(ret.ID)
	rev1.Save()

	rev2 := model.DummyFinanceRevenue()
	rev2.IsDeleted = int8(0)
	rev2.DocumentStatus = "cleared"
	rev2.Amount = float64(30000)
	rev2.RefType = "purchase_return"
	rev2.RefID = uint64(ret.ID)
	rev2.Save()

	rev3 := model.DummyFinanceRevenue()
	rev3.IsDeleted = int8(1)
	rev3.DocumentStatus = "cleared"
	rev3.Amount = float64(20000)
	rev3.RefType = "purchase_return"
	rev3.RefID = uint64(ret.ID)
	rev3.Save()

	rev := model.DummyFinanceRevenue()
	rev.IsDeleted = int8(0)
	rev.DocumentStatus = "uncleared"
	rev.Amount = float64(2000)
	rev.RefType = "purchase_return"
	rev.RefID = uint64(ret.ID)
	rev.Save()

	data := tester.D{}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-revenue/"+common.Encrypt(rev.ID)+"/approve").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})
	rev.Read("ID")
	assert.Equal(t, "cleared", rev.DocumentStatus)

	ret.Read("ID")
	assert.Equal(t, "active", ret.DocumentStatus)
}

// TestHandler_URLMappingApproveFinanceRevenueFailStatus gagal dengan purchase return,cleared
func TestHandler_URLMappingApproveFinanceRevenueFailStatus(t *testing.T) {
	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy
	ret := model.DummyPurchaseReturn()
	ret.TotalAmount = float64(50000)
	ret.DocumentStatus = "active"
	ret.IsDeleted = int8(0)
	ret.Save()

	rev1 := model.DummyFinanceRevenue()
	rev1.IsDeleted = int8(0)
	rev1.DocumentStatus = "cleared"
	rev1.Amount = float64(15000)
	rev1.RefType = "purchase_return"
	rev1.RefID = uint64(ret.ID)
	rev1.Save()

	rev2 := model.DummyFinanceRevenue()
	rev2.IsDeleted = int8(0)
	rev2.DocumentStatus = "cleared"
	rev2.Amount = float64(30000)
	rev2.RefType = "purchase_return"
	rev2.RefID = uint64(ret.ID)
	rev2.Save()

	rev3 := model.DummyFinanceRevenue()
	rev3.IsDeleted = int8(1)
	rev3.DocumentStatus = "cleared"
	rev3.Amount = float64(20000)
	rev3.RefType = "purchase_return"
	rev3.RefID = uint64(ret.ID)
	rev3.Save()

	rev := model.DummyFinanceRevenue()
	rev.IsDeleted = int8(0)
	rev.DocumentStatus = "cleared"
	rev.Amount = float64(1000)
	rev.RefType = "purchase_return"
	rev.RefID = uint64(ret.ID)
	rev.Save()

	data := tester.D{}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-revenue/"+common.Encrypt(rev.ID)+"/approve").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})
}

// TestHandler_URLMappingApproveFinanceRevenueFailDeleted gagal dengan purchase return,deleted
func TestHandler_URLMappingApproveFinanceRevenueFailDeleted(t *testing.T) {
	// login user
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy
	ret := model.DummyPurchaseReturn()
	ret.TotalAmount = float64(50000)
	ret.DocumentStatus = "active"
	ret.IsDeleted = int8(0)
	ret.Save()

	rev1 := model.DummyFinanceRevenue()
	rev1.IsDeleted = int8(0)
	rev1.DocumentStatus = "cleared"
	rev1.Amount = float64(15000)
	rev1.RefType = "purchase_return"
	rev1.RefID = uint64(ret.ID)
	rev1.Save()

	rev2 := model.DummyFinanceRevenue()
	rev2.IsDeleted = int8(0)
	rev2.DocumentStatus = "cleared"
	rev2.Amount = float64(30000)
	rev2.RefType = "purchase_return"
	rev2.RefID = uint64(ret.ID)
	rev2.Save()

	rev3 := model.DummyFinanceRevenue()
	rev3.IsDeleted = int8(1)
	rev3.DocumentStatus = "cleared"
	rev3.Amount = float64(20000)
	rev3.RefType = "purchase_return"
	rev3.RefID = uint64(ret.ID)
	rev3.Save()

	rev := model.DummyFinanceRevenue()
	rev.IsDeleted = int8(1)
	rev.DocumentStatus = "uncleared"
	rev.Amount = float64(2000)
	rev.RefType = "purchase_return"
	rev.RefID = uint64(ret.ID)
	rev.Save()

	data := tester.D{}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-revenue/"+common.Encrypt(rev.ID)+"/approve").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusNotFound, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})
}

// tes kasir membuat finance revenue sukses
func TestHandler_URLMappingCreateFinanceRevenueCashierSuccess(t *testing.T) {

	// login user kasir
	user := model.DummyUserPriviledgeWithUsergroup(4)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy sales invoice
	si := model.DummySalesInvoice()
	si.DocumentStatus = "new"
	si.TotalAmount = 1000000
	si.TotalRevenued = 0
	si.IsDeleted = 0
	si.Save()

	customer := model.DummyPartnership()

	so := si.SalesOrder
	so.DocumentStatus = "active"
	so.Customer = customer
	so.TotalPaid = float64(40000)
	so.TotalCharge = float64(50000)
	so.InvoiceStatus = "active"
	so.ShipmentStatus = "finished"
	so.FulfillmentStatus = "finished"
	so.IsDeleted = int8(0)
	so.Save()

	customer.PartnershipType = "customer"
	customer.IsDeleted = int8(0)
	customer.TotalDebt = float64(20000)
	customer.Save()

	data := tester.D{"ref_id": common.Encrypt(si.ID), "ref_type": "sales_invoice", "recognition_date": time.Now(), "payment_method": "debit_card", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

	si.Read("ID")
	// cek document status pada sales invoice, harusnya menjadi active
	assert.Equal(t, "active", si.DocumentStatus)
}

// tes kasir membuat finance revenue gagal
func TestHandler_URLMappingCreateFinanceRevenueCashierFail(t *testing.T) {

	// login user kasir
	user := model.DummyUserPriviledgeWithUsergroup(4)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy sales invoice
	si := model.DummySalesInvoice()
	si.DocumentStatus = "new"
	si.TotalAmount = 1000000
	si.TotalRevenued = 1000
	si.IsDeleted = 0
	si.Save()

	data := tester.D{"ref_id": common.Encrypt(si.ID), "ref_type": "sales_invoice", "recognition_date": time.Now(), "payment_method": "debit_card", "amount": 50000,
		"bank_number": common.RandomNumeric(9), "bank_name": "Tes", "bank_holder": "Tes", "note": "sales invoice"}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-revenue").SetJSON(data).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", data, res.Body.String()))
	})

	si.Read("ID")
	// cek document status pada sales invoice, harusnya tetap new
	assert.Equal(t, "new", si.DocumentStatus)
}
