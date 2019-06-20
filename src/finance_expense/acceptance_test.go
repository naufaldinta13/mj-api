// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package financeExpense_test

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
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	exp := model.DummyFinanceExpense()
	exp.IsDeleted = 0
	exp.Save()

	expID := common.Encrypt(exp.ID)

	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/finance-expense", "GET", 200},
		{"/v1/finance-expense/" + expID, "GET", 200},
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

// TestHandler_URLMappingCreateFinanceExpenseSukses1 Purchase invoice credit card
func TestHandler_URLMappingCreateFinanceExpenseSukses1(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	porder := model.DummyPurchaseOrder()
	porder.TotalCharge = float64(90000000)
	porder.Save()

	pinvoice := model.DummyPurchaseInvoice()
	pinvoice.DocumentStatus = "new"
	pinvoice.TotalAmount = float64(90000)
	pinvoice.PurchaseOrder = porder
	pinvoice.IsDeleted = 0
	pinvoice.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.DocumentStatus = "uncleared"
	fexpense.Amount = float64(50000)
	fexpense.RefID = uint64(pinvoice.ID)
	fexpense.RefType = "purchase_invoice"
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "purchase_invoice",
		"ref_id":           common.Encrypt(pinvoice.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(40000),
		"bank_number":      common.RandomNumeric(9),
		"bank_name":        "BRI",
		"bank_holder":      "Andi Sumarno",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

	pinvoice.Read()
	assert.Equal(t, "active", pinvoice.DocumentStatus, "document status pada purchase invoice harus berubah menjadi active")
}

// TestHandler_URLMappingCreateFinanceExpenseSukses2 Purchase invoice debit card
func TestHandler_URLMappingCreateFinanceExpenseSukses2(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	porder := model.DummyPurchaseOrder()
	porder.TotalCharge = float64(90000000)
	porder.Save()

	pinvoice := model.DummyPurchaseInvoice()
	pinvoice.DocumentStatus = "new"
	pinvoice.TotalAmount = float64(90000)
	pinvoice.PurchaseOrder = porder
	pinvoice.IsDeleted = 0
	pinvoice.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "purchase_invoice",
		"ref_id":           common.Encrypt(pinvoice.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "debit_card",
		"amount":           float64(30000),
		"bank_number":      common.RandomNumeric(9),
		"bank_name":        "BRI",
		"bank_holder":      "Andi Sumarno",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateFinanceExpenseSukses3 Purchase invoice giro
func TestHandler_URLMappingCreateFinanceExpenseSukses3(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	porder := model.DummyPurchaseOrder()
	porder.TotalCharge = float64(90000000)
	porder.Save()

	pinvoice := model.DummyPurchaseInvoice()
	pinvoice.DocumentStatus = "new"
	pinvoice.TotalAmount = float64(90000)
	pinvoice.PurchaseOrder = porder
	pinvoice.IsDeleted = 0
	pinvoice.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "purchase_invoice",
		"ref_id":           common.Encrypt(pinvoice.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "giro",
		"amount":           float64(40000),
		"giro_number":      common.RandomNumeric(9),
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateFinanceExpenseSukses4 Sales return credit card
func TestHandler_URLMappingCreateFinanceExpenseSukses4(t *testing.T) {
	test.DataCleanUp("finance_expense", "sales_return", "purchase_invoice")
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sorder := model.DummySalesOrder()
	sorder.TotalCharge = float64(9000000)
	sorder.Save()

	sreturn := model.DummySalesReturn()
	sreturn.DocumentStatus = "active"
	sreturn.TotalAmount = float64(90000)
	sreturn.SalesOrder = sorder
	sreturn.IsDeleted = 0
	sreturn.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "sales_return",
		"ref_id":           common.Encrypt(sreturn.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(40000),
		"bank_number":      common.RandomNumeric(9),
		"bank_name":        "BRI",
		"bank_holder":      "Andi Sumarno",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

	sreturn.Read()
	assert.Equal(t, "active", sreturn.DocumentStatus, "document status pada sales return harus berubah menjadi active")
}

// TestHandler_URLMappingCreateFinanceExpenseSukses5 Sales return debit card
func TestHandler_URLMappingCreateFinanceExpenseSukses5(t *testing.T) {
	test.DataCleanUp("finance_expense", "sales_return", "purchase_invoice")
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sorder := model.DummySalesOrder()
	sorder.TotalCharge = float64(9000000)
	sorder.Save()

	sreturn := model.DummySalesReturn()
	sreturn.DocumentStatus = "active"
	sreturn.TotalAmount = float64(90000)
	sreturn.SalesOrder = sorder
	sreturn.IsDeleted = 0
	sreturn.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "sales_return",
		"ref_id":           common.Encrypt(sreturn.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "debit_card",
		"amount":           float64(40000),
		"bank_number":      common.RandomNumeric(9),
		"bank_name":        "BRI",
		"bank_holder":      "Andi Sumarno",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateFinanceExpenseSukses5 Sales return giro
func TestHandler_URLMappingCreateFinanceExpenseSukses6(t *testing.T) {
	test.DataCleanUp("finance_expense", "sales_return", "purchase_invoice")
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sorder := model.DummySalesOrder()
	sorder.TotalCharge = float64(9000000)
	sorder.Save()

	sreturn := model.DummySalesReturn()
	sreturn.DocumentStatus = "active"
	sreturn.TotalAmount = float64(9000)
	sreturn.SalesOrder = sorder
	sreturn.IsDeleted = 0
	sreturn.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "sales_return",
		"ref_id":           common.Encrypt(sreturn.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "giro",
		"amount":           float64(400),
		"giro_number":      common.RandomNumeric(9),
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateFinanceExpenseGagal ref_id sales_return not found
func TestHandler_URLMappingCreateFinanceExpenseGagalRefIDSRNotFound(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sreturn := model.DummySalesReturn()
	sreturn.DocumentStatus = "active"
	sreturn.TotalAmount = float64(100000)
	sreturn.IsDeleted = 0
	sreturn.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.DocumentStatus = "cleared"
	fexpense.Amount = float64(100000)
	fexpense.RefID = uint64(sreturn.ID)
	fexpense.RefType = "sales_return"
	fexpense.Save()

	sreturn.Delete()

	// setting body
	scenario := tester.D{
		"ref_type":         "sales_return",
		"ref_id":           common.Encrypt(sreturn.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(40000),
		"bank_number":      common.RandomNumeric(9),
		"bank_name":        "BRI",
		"bank_holder":      "Andi Sumarno",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateFinanceExpenseGagal ref_id sales_return not found
func TestHandler_URLMappingCreateFinanceExpenseGagalRefIDPINotFound(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	pinvoice := model.DummyPurchaseInvoice()
	pinvoice.DocumentStatus = "active"
	pinvoice.TotalAmount = float64(100000)
	pinvoice.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.DocumentStatus = "cleared"
	fexpense.Amount = float64(100000)
	fexpense.RefID = uint64(pinvoice.ID)
	fexpense.RefType = "sales_return"
	fexpense.Save()

	pinvoice.Delete()

	// setting body
	scenario := tester.D{
		"ref_type":         "purchase_invoice",
		"ref_id":           common.Encrypt(pinvoice.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(40000),
		"bank_number":      common.RandomNumeric(9),
		"bank_name":        "BRI",
		"bank_holder":      "Andi Sumarno",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateFinanceExpenseGagal ref_id invalid
func TestHandler_URLMappingCreateFinanceExpenseGagal(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sreturn := model.DummySalesReturn()
	sreturn.DocumentStatus = "active"
	sreturn.TotalAmount = float64(100000)
	sreturn.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.DocumentStatus = "cleared"
	fexpense.Amount = float64(100000)
	fexpense.RefID = uint64(sreturn.ID)
	fexpense.RefType = "sales_return"
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "sales_return",
		"ref_id":           "abcdefghijklmn",
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(40000),
		"bank_number":      common.RandomNumeric(9),
		"bank_name":        "BRI",
		"bank_holder":      "Andi Sumarno",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateFinanceExpenseGagal1 yang required tidak diisi
func TestHandler_URLMappingCreateFinanceExpenseGagal1(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sreturn := model.DummySalesReturn()
	sreturn.DocumentStatus = "new"
	sreturn.TotalAmount = float64(100000)
	sreturn.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.DocumentStatus = "finished"
	fexpense.Amount = float64(50000)
	fexpense.RefID = uint64(sreturn.ID)
	fexpense.RefType = "purchase_invoice"
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "",
		"ref_id":           common.Encrypt(2120),
		"recognition_date": time.Time{},
		"payment_method":   "",
		"amount":           float64(40000),
		"bank_number":      common.RandomNumeric(9),
		"bank_name":        "BRI",
		"bank_holder":      "Andi Sumarno",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateFinanceExpenseGagal2 payment method credit card bank number dll tidak diisi
func TestHandler_URLMappingCreateFinanceExpenseGagal2(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sorder := model.DummySalesOrder()
	sorder.TotalCharge = float64(9000000)
	sorder.Save()

	sreturn := model.DummySalesReturn()
	sreturn.DocumentStatus = "active"
	sreturn.TotalAmount = float64(9000)
	sreturn.SalesOrder = sorder
	sreturn.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.DocumentStatus = "finished"
	fexpense.Amount = float64(50000)
	fexpense.RefID = uint64(sreturn.ID)
	fexpense.RefType = "purchase_invoice"
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "purchase_invoice",
		"ref_id":           common.Encrypt(sreturn.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(40000),
		"bank_number":      "",
		"bank_name":        "",
		"bank_holder":      "",
		"giro_number":      "01234567890",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateFinanceExpenseGagal3 payment method debit card bank number dll tidak diisi
func TestHandler_URLMappingCreateFinanceExpenseGagal3(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sorder := model.DummySalesOrder()
	sorder.TotalCharge = float64(9000000)
	sorder.Save()

	sreturn := model.DummySalesReturn()
	sreturn.DocumentStatus = "active"
	sreturn.TotalAmount = float64(9000)
	sreturn.SalesOrder = sorder
	sreturn.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.DocumentStatus = "finished"
	fexpense.Amount = float64(50000)
	fexpense.RefID = uint64(sreturn.ID)
	fexpense.RefType = "purchase_invoice"
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "purchase_invoice",
		"ref_id":           common.Encrypt(sreturn.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "debit_card",
		"amount":           float64(40000),
		"bank_number":      "",
		"bank_name":        "",
		"bank_holder":      "",
		"giro_number":      "1234567890",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateFinanceExpenseGagal4 payment method giro giro number dll tidak diisi
func TestHandler_URLMappingCreateFinanceExpenseGagal4(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sorder := model.DummySalesOrder()
	sorder.TotalCharge = float64(9000000)
	sorder.Save()

	sreturn := model.DummySalesReturn()
	sreturn.DocumentStatus = "active"
	sreturn.TotalAmount = float64(9000)
	sreturn.SalesOrder = sorder
	sreturn.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.DocumentStatus = "finished"
	fexpense.Amount = float64(50000)
	fexpense.RefID = uint64(sreturn.ID)
	fexpense.RefType = "purchase_invoice"
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "purchase_invoice",
		"ref_id":           common.Encrypt(sreturn.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "giro",
		"amount":           float64(0),
		"giro_number":      "",
		"bank_number":      "123",
		"bank_name":        "123",
		"bank_holder":      "123",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateFinanceExpenseGagal5 dokument status purchase invoice sudah finish
func TestHandler_URLMappingCreateFinanceExpenseGagal5(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	porder := model.DummyPurchaseOrder()
	porder.TotalCharge = float64(90000000)
	porder.Save()

	pinvoice := model.DummyPurchaseInvoice()
	pinvoice.DocumentStatus = "finished"
	pinvoice.TotalAmount = float64(90000)
	pinvoice.PurchaseOrder = porder
	pinvoice.IsDeleted = 0
	pinvoice.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.DocumentStatus = "finished"
	fexpense.Amount = float64(50000)
	fexpense.RefID = uint64(pinvoice.ID)
	fexpense.RefType = "purchase_invoice"
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "purchase_invoice",
		"ref_id":           common.Encrypt(pinvoice.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(40000),
		"bank_number":      common.RandomNumeric(9),
		"bank_name":        "BRI",
		"bank_holder":      "Andi Sumarno",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateFinanceExpenseGagal6 dokument status sales_return sudah finish
func TestHandler_URLMappingCreateFinanceExpenseGagal6(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sorder := model.DummySalesOrder()
	sorder.TotalCharge = float64(9000000)
	sorder.Save()

	sreturn := model.DummySalesReturn()
	sreturn.DocumentStatus = "finished"
	sreturn.TotalAmount = float64(9000)
	sreturn.SalesOrder = sorder
	sreturn.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.DocumentStatus = "finished"
	fexpense.Amount = float64(50000)
	fexpense.RefID = uint64(sreturn.ID)
	fexpense.RefType = "sales_return"
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "sales_return",
		"ref_id":           common.Encrypt(sreturn.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(40000),
		"bank_number":      common.RandomNumeric(9),
		"bank_name":        "BRI",
		"bank_holder":      "Andi Sumarno",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateFinanceExpenseGagal7 pembayaran lebih purchase invoice
func TestHandler_URLMappingCreateFinanceExpenseGagal7(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	porder := model.DummyPurchaseOrder()
	porder.TotalCharge = float64(90000000)
	porder.Save()

	pinvoice := model.DummyPurchaseInvoice()
	pinvoice.DocumentStatus = "new"
	pinvoice.TotalAmount = float64(900000)
	pinvoice.PurchaseOrder = porder
	pinvoice.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.DocumentStatus = "cleared"
	fexpense.Amount = float64(90000000)
	fexpense.RefID = uint64(pinvoice.ID)
	fexpense.RefType = "purchase_invoice"
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "purchase_invoice",
		"ref_id":           common.Encrypt(pinvoice.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(4000000),
		"bank_number":      common.RandomNumeric(9),
		"bank_name":        "BRI",
		"bank_holder":      "Andi Sumarno",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateFinanceExpenseGagal8 pembayaran lebih sales_return
func TestHandler_URLMappingCreateFinanceExpenseGagal8(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sorder := model.DummySalesOrder()
	sorder.TotalCharge = float64(9000000)
	sorder.IsDeleted = 0
	sorder.Save()

	sreturn := model.DummySalesReturn()
	sreturn.DocumentStatus = "active"
	sreturn.TotalAmount = float64(9000)
	sreturn.SalesOrder = sorder
	sreturn.IsDeleted = 0
	sreturn.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.DocumentStatus = "cleared"
	fexpense.Amount = float64(9000000)
	fexpense.RefID = uint64(sreturn.ID)
	fexpense.RefType = "sales_return"
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "sales_return",
		"ref_id":           common.Encrypt(sreturn.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(40000),
		"bank_number":      common.RandomNumeric(9),
		"bank_name":        "BRI",
		"bank_holder":      "Andi Sumarno",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingCreateFinanceExpenseGagal9 payment method giro giro number dll tidak diisi
func TestHandler_URLMappingCreateFinanceExpenseGagal9(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sorder := model.DummySalesOrder()
	sorder.TotalCharge = float64(9000000)
	sorder.Save()

	sreturn := model.DummySalesReturn()
	sreturn.DocumentStatus = "active"
	sreturn.TotalAmount = float64(9000)
	sreturn.SalesOrder = sorder
	sreturn.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.DocumentStatus = "finished"
	fexpense.Amount = float64(50000)
	fexpense.RefID = uint64(sreturn.ID)
	fexpense.RefType = "purchase_invoice"
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "purchase_invoice",
		"ref_id":           common.Encrypt(sreturn.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "cash",
		"amount":           float64(0),
		"giro_number":      "123",
		"bank_number":      "123",
		"bank_name":        "123",
		"bank_holder":      "123",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/finance-expense").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPutFinanceExpenseSukses1 Berhasil update untuk purchase invoice from giro to credit card
func TestHandler_URLMappingPutFinanceExpenseSukses1(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	porder := model.DummyPurchaseOrder()
	porder.IsDeleted = int8(0)
	porder.TotalCharge = float64(100000)
	porder.Save()

	pinvoice := model.DummyPurchaseInvoice()
	pinvoice.PurchaseOrder = porder
	pinvoice.DocumentStatus = "active"
	pinvoice.TotalAmount = float64(100000)
	pinvoice.IsDeleted = 0
	pinvoice.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.RefID = uint64(pinvoice.ID)
	fexpense.DocumentStatus = "uncleared"
	fexpense.RefType = "purchase_invoice"
	fexpense.PaymentMethod = "giro"
	fexpense.BankNumber = "GIRO123456789CEK"
	fexpense.Amount = float64(35000)
	fexpense.IsDeleted = int8(0)
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"ref_type":         "purchase_invoice",
		"ref_id":           common.Encrypt(pinvoice.ID),
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(15000),
		"bank_number":      "123CreditCard",
		"bank_name":        "123",
		"bank_holder":      "123",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(fexpense.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	fexpense.Read()
	assert.Equal(t, float64(15000), fexpense.Amount, "Seharusnya berubah menjadi 15000")
	assert.Equal(t, "credit_card", fexpense.PaymentMethod, "Payment method harusnya berubah menjadi credit card")
	assert.Equal(t, "123CreditCard", fexpense.BankNumber, "Bank number berubah menjadi 123CreditCard")
}

// TestHandler_URLMappingPutFinanceExpenseSukses2 Berhasil update untuk sales return from giro to credit card
func TestHandler_URLMappingPutFinanceExpenseSukses2(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sorder := model.DummySalesOrder()
	sorder.IsDeleted = int8(0)
	sorder.TotalCharge = float64(100000)
	sorder.Save()

	sreturn := model.DummySalesReturn()
	sreturn.SalesOrder = sorder
	sreturn.DocumentStatus = "active"
	sreturn.TotalAmount = float64(100000)
	sreturn.IsDeleted = 0
	sreturn.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.RefID = uint64(sreturn.ID)
	fexpense.DocumentStatus = "uncleared"
	fexpense.RefType = "sales_return"
	fexpense.PaymentMethod = "giro"
	fexpense.BankNumber = "GIRO123456789CEK"
	fexpense.Amount = float64(35000)
	fexpense.IsDeleted = int8(0)
	fexpense.Save()

	fmt.Println("", fexpense, "FEREFOD", fexpense.RefID)

	// setting body
	scenario := tester.D{
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(15000),
		"bank_number":      "123CreditCard",
		"bank_name":        "123",
		"bank_holder":      "123",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(fexpense.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	fexpense.Read()
	assert.Equal(t, float64(15000), fexpense.Amount, "Seharusnya berubah menjadi 15000")
	assert.Equal(t, "credit_card", fexpense.PaymentMethod, "Payment method harusnya berubah menjadi credit card")
	assert.Equal(t, "123CreditCard", fexpense.BankNumber, "Bank number berubah menjadi 123CreditCard")
}

// TestHandler_URLMappingPutFinanceExpenseSukses3 Berhasil update untuk purchase invoice from credit card to giro
func TestHandler_URLMappingPutFinanceExpenseSukses3(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	porder := model.DummyPurchaseOrder()
	porder.IsDeleted = int8(0)
	porder.TotalCharge = float64(100000)
	porder.Save()

	pinvoice := model.DummyPurchaseInvoice()
	pinvoice.PurchaseOrder = porder
	pinvoice.DocumentStatus = "active"
	pinvoice.TotalAmount = float64(100000)
	pinvoice.IsDeleted = 0
	pinvoice.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.RefID = uint64(pinvoice.ID)
	fexpense.DocumentStatus = "uncleared"
	fexpense.RefType = "purchase_invoice"
	fexpense.PaymentMethod = "credit_card"
	fexpense.BankNumber = "CC123456789CEK"
	fexpense.Amount = float64(35000)
	fexpense.IsDeleted = int8(0)
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "giro",
		"amount":           float64(15000),
		"giro_number":      "123Giro",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(fexpense.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	fexpense.Read()
	assert.Equal(t, float64(15000), fexpense.Amount, "Seharusnya berubah menjadi 15000")
	assert.Equal(t, "giro", fexpense.PaymentMethod, "Payment method harusnya berubah menjadi giro")
	assert.Equal(t, "123Giro", fexpense.BankNumber, "Bank number berubah menjadi 123Giro")
}

// TestHandler_URLMappingPutFinanceExpenseSukses4 Berhasil update untuk purchase invoice from credit card to cash
func TestHandler_URLMappingPutFinanceExpenseSukses4(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	porder := model.DummyPurchaseOrder()
	porder.IsDeleted = int8(0)
	porder.TotalCharge = float64(100000)
	porder.Save()

	pinvoice := model.DummyPurchaseInvoice()
	pinvoice.PurchaseOrder = porder
	pinvoice.DocumentStatus = "active"
	pinvoice.TotalAmount = float64(100000)
	pinvoice.IsDeleted = 0
	pinvoice.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.RefID = uint64(pinvoice.ID)
	fexpense.DocumentStatus = "uncleared"
	fexpense.RefType = "purchase_invoice"
	fexpense.PaymentMethod = "credit_card"
	fexpense.BankNumber = "CC123456789CEK"
	fexpense.Amount = float64(35000)
	fexpense.IsDeleted = int8(0)
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "cash",
		"amount":           float64(15000),
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(fexpense.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	fexpense.Read()
	assert.Equal(t, float64(15000), fexpense.Amount, "Seharusnya berubah menjadi 15000")
	assert.Equal(t, "cash", fexpense.PaymentMethod, "Payment method harusnya berubah menjadi cash")
	assert.Equal(t, "", fexpense.BankNumber, "Bank number berubah menjadi Kosong")
}

// TestHandler_URLMappingPutFinanceExpenseGagal1 Gagal update untuk purchase invoice credit card for inputing giro number
// and bank number, name, holder is blank
func TestHandler_URLMappingPutFinanceExpenseGagal1(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	porder := model.DummyPurchaseOrder()
	porder.IsDeleted = int8(0)
	porder.TotalCharge = float64(100000)
	porder.Save()

	pinvoice := model.DummyPurchaseInvoice()
	pinvoice.PurchaseOrder = porder
	pinvoice.DocumentStatus = "active"
	pinvoice.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.RefID = uint64(pinvoice.ID)
	fexpense.DocumentStatus = "uncleared"
	fexpense.RefType = "purchase_invoice"
	fexpense.PaymentMethod = "credit_card"
	fexpense.BankNumber = "CC123456789CEK"
	fexpense.Amount = float64(35000)
	fexpense.IsDeleted = int8(0)
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(0),
		"bank_number":      "",
		"bank_name":        "",
		"bank_holder":      "",
		"giro_number":      "Giro123",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(fexpense.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPutFinanceExpenseGagal2 Gagal update untuk purchase invoice giro for not inputing giro number
// and bank number, name, holder is filled
func TestHandler_URLMappingPutFinanceExpenseGagal2(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	porder := model.DummyPurchaseOrder()
	porder.IsDeleted = int8(0)
	porder.TotalCharge = float64(100000)
	porder.Save()

	pinvoice := model.DummyPurchaseInvoice()
	pinvoice.PurchaseOrder = porder
	pinvoice.DocumentStatus = "active"
	pinvoice.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.RefID = uint64(pinvoice.ID)
	fexpense.DocumentStatus = "uncleared"
	fexpense.RefType = "purchase_invoice"
	fexpense.PaymentMethod = "credit_card"
	fexpense.BankNumber = "CC123456789CEK"
	fexpense.Amount = float64(35000)
	fexpense.IsDeleted = int8(0)
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "giro",
		"amount":           float64(0),
		"bank_number":      "CC123456",
		"bank_name":        "BCA",
		"bank_holder":      "Andi Soemarno",
		"giro_number":      "",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(fexpense.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPutFinanceExpenseGagal3 Gagal update untuk purchase invoice cash for inputing giro number,
// bank number, name, holder is filled
func TestHandler_URLMappingPutFinanceExpenseGagal3(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	porder := model.DummyPurchaseOrder()
	porder.IsDeleted = int8(0)
	porder.TotalCharge = float64(100000)
	porder.Save()

	pinvoice := model.DummyPurchaseInvoice()
	pinvoice.PurchaseOrder = porder
	pinvoice.DocumentStatus = "active"
	pinvoice.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.RefID = uint64(pinvoice.ID)
	fexpense.DocumentStatus = "uncleared"
	fexpense.RefType = "purchase_invoice"
	fexpense.PaymentMethod = "credit_card"
	fexpense.BankNumber = "CC123456789CEK"
	fexpense.Amount = float64(35000)
	fexpense.IsDeleted = int8(0)
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "cash",
		"amount":           float64(0),
		"bank_number":      "CC123456",
		"bank_name":        "BCA",
		"bank_holder":      "Andi Soemarno",
		"giro_number":      "Giro123",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(fexpense.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPutFinanceExpenseGagal4 Gagal update ref id is invalid
func TestHandler_URLMappingPutFinanceExpenseGagal4(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	porder := model.DummyPurchaseOrder()
	porder.IsDeleted = int8(0)
	porder.TotalCharge = float64(100000)
	porder.Save()

	pinvoice := model.DummyPurchaseInvoice()
	pinvoice.PurchaseOrder = porder
	pinvoice.DocumentStatus = "finished"
	pinvoice.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.RefID = 111111
	fexpense.DocumentStatus = "cleared"
	fexpense.RefType = "purchase_invoice"
	fexpense.PaymentMethod = "giro"
	fexpense.BankNumber = "GIRO123456789CEK"
	fexpense.Amount = float64(35000)
	fexpense.IsDeleted = int8(0)
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(15000),
		"bank_number":      "123CreditCard",
		"bank_name":        "123",
		"bank_holder":      "123",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(fexpense.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPutFinanceExpenseGagal5 Gagal update document status purchase invoice is finished
func TestHandler_URLMappingPutFinanceExpenseGagal5(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	porder := model.DummyPurchaseOrder()
	porder.IsDeleted = int8(0)
	porder.TotalCharge = float64(100000)
	porder.Save()

	pinvoice := model.DummyPurchaseInvoice()
	pinvoice.PurchaseOrder = porder
	pinvoice.DocumentStatus = "finished"
	pinvoice.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.RefID = uint64(pinvoice.ID)
	fexpense.DocumentStatus = "cleared"
	fexpense.RefType = "purchase_invoice"
	fexpense.PaymentMethod = "giro"
	fexpense.BankNumber = "GIRO123456789CEK"
	fexpense.Amount = float64(35000)
	fexpense.IsDeleted = int8(0)
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(15000),
		"bank_number":      "123CreditCard",
		"bank_name":        "123",
		"bank_holder":      "123",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(fexpense.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPutFinanceExpenseGagal6 Gagal update financial expense purchase invoice paid too much
func TestHandler_URLMappingPutFinanceExpenseGagal6(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	porder := model.DummyPurchaseOrder()
	porder.IsDeleted = int8(0)
	porder.TotalCharge = float64(100000)
	porder.Save()

	pinvoice := model.DummyPurchaseInvoice()
	pinvoice.PurchaseOrder = porder
	pinvoice.DocumentStatus = "active"
	pinvoice.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.RefID = uint64(pinvoice.ID)
	fexpense.DocumentStatus = "uncleared"
	fexpense.RefType = "purchase_invoice"
	fexpense.PaymentMethod = "giro"
	fexpense.BankNumber = "GIRO123456789CEK"
	fexpense.Amount = float64(100000)
	fexpense.IsDeleted = int8(0)
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(15000),
		"bank_number":      "123CreditCard",
		"bank_name":        "123",
		"bank_holder":      "123",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(fexpense.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPutFinanceExpenseGagal7 Gagal update ref id is invalid
func TestHandler_URLMappingPutFinanceExpenseGagal7(t *testing.T) {
	test.DataCleanUp("sales_order", "sales_return", "finance_expense")
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sorder := model.DummySalesOrder()
	sorder.IsDeleted = int8(0)
	sorder.TotalCharge = float64(100000)
	sorder.Save()

	sreturn := model.DummySalesReturn()
	sreturn.SalesOrder = sorder
	sreturn.DocumentStatus = "finished"
	sreturn.IsDeleted = int8(0)
	sreturn.Save()

	fexpense := model.DummyFinanceExpense()
	feID := common.Encrypt(fexpense.ID)
	fexpense.RefID = 11111
	fexpense.DocumentStatus = "uncleared"
	fexpense.RefType = "sales_return"
	fexpense.PaymentMethod = "giro"
	fexpense.BankNumber = "GIRO123456789CEK"
	fexpense.Amount = float64(35000)
	fexpense.IsDeleted = int8(0)
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(15000),
		"bank_number":      "123CreditCard",
		"bank_name":        "123",
		"bank_holder":      "123",
		"note":             "sales_return",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+feID).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	test.DataCleanUp("sales_order", "sales_return", "finance_expense")
}

// TestHandler_URLMappingPutFinanceExpenseGagal8 Gagal update document status sales return is finished
func TestHandler_URLMappingPutFinanceExpenseGagal8(t *testing.T) {
	test.DataCleanUp("sales_order", "sales_return", "finance_expense")
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sorder := model.DummySalesOrder()
	sorder.IsDeleted = int8(0)
	sorder.TotalCharge = float64(100000)
	sorder.Save()

	sreturn := model.DummySalesReturn()
	sreturn.SalesOrder = sorder
	sreturn.DocumentStatus = "finished"
	sreturn.IsDeleted = int8(0)
	sreturn.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.RefID = uint64(sreturn.ID)
	fexpense.DocumentStatus = "uncleared"
	fexpense.RefType = "sales_return"
	fexpense.PaymentMethod = "giro"
	fexpense.BankNumber = "GIRO123456789CEK"
	fexpense.Amount = float64(35000)
	fexpense.IsDeleted = int8(0)
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(15000),
		"bank_number":      "123CreditCard",
		"bank_name":        "123",
		"bank_holder":      "123",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(fexpense.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	test.DataCleanUp("sales_order", "sales_return", "finance_expense")
}

// TestHandler_URLMappingPutFinanceExpenseGagal9 Gagal update financial expense sales return paid too much
func TestHandler_URLMappingPutFinanceExpenseGagal9(t *testing.T) {
	test.DataCleanUp("sales_order", "sales_return", "finance_expense")
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	sorder := model.DummySalesOrder()
	sorder.IsDeleted = int8(0)
	sorder.TotalCharge = float64(100000)
	sorder.Save()

	sreturn := model.DummySalesReturn()
	sreturn.SalesOrder = sorder
	sreturn.DocumentStatus = "active"
	sreturn.IsDeleted = int8(0)
	sreturn.Save()

	fexpense := model.DummyFinanceExpense()
	fexpense.RefID = uint64(sreturn.ID)
	fexpense.DocumentStatus = "uncleared"
	fexpense.RefType = "sales_return"
	fexpense.PaymentMethod = "giro"
	fexpense.BankNumber = "GIRO123456789CEK"
	fexpense.Amount = float64(100000)
	fexpense.IsDeleted = int8(0)
	fexpense.Save()

	// setting body
	scenario := tester.D{
		"recognition_date": time.Date(2017, 10, 15, 8, 30, 0, 0, time.Local),
		"payment_method":   "credit_card",
		"amount":           float64(15000),
		"bank_number":      "123CreditCard",
		"bank_name":        "123",
		"bank_holder":      "123",
		"note":             "purchase_invoice",
	}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(fexpense.ID)).SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	test.DataCleanUp("sales_order", "sales_return", "finance_expense")
}

// TestHandler_URLMappingPutFinanceExpenseIDNotFound Gagal update ID not found
func TestHandler_URLMappingPutFinanceExpenseIDNotFound(t *testing.T) {
	test.DataCleanUp("sales_order", "sales_return", "finance_expense")
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/999999999").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusNotFound, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	test.DataCleanUp("sales_order", "sales_return", "finance_expense")
}

//// approve finance expense /////////////////////////////////////////////////////////////

// TestHandler_URLMappingApproveFinanceExpenseSuccess Berhasil approve untuk finance expense po invoice finish
func TestHandler_URLMappingApproveFinanceExpenseSuccess(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy
	inv := model.DummyPurchaseInvoice()
	inv.TotalAmount = float64(50000)
	inv.TotalPaid = float64(40000)
	inv.DocumentStatus = "active"
	inv.IsDeleted = int8(0)
	inv.Save()

	po := inv.PurchaseOrder
	po.DocumentStatus = "active"
	po.TotalPaid = float64(40000)
	po.TotalCharge = float64(50000)
	po.InvoiceStatus = "active"
	po.ReceivingStatus = "finished"
	po.IsDeleted = int8(0)
	po.Save()

	supplier := po.Supplier
	supplier.PartnershipType = "supplier"
	supplier.IsDeleted = int8(0)
	supplier.TotalCredit = float64(20000)
	supplier.Save()

	exp1 := model.DummyFinanceExpense()
	exp1.IsDeleted = int8(0)
	exp1.DocumentStatus = "cleared"
	exp1.Amount = float64(20000)
	exp1.RefType = "purchase_invoice"
	exp1.RefID = uint64(inv.ID)
	exp1.Save()

	exp2 := model.DummyFinanceExpense()
	exp2.IsDeleted = int8(0)
	exp2.DocumentStatus = "cleared"
	exp2.Amount = float64(20000)
	exp2.RefType = "purchase_invoice"
	exp2.RefID = uint64(inv.ID)
	exp2.Save()

	exp := model.DummyFinanceExpense()
	exp.IsDeleted = int8(0)
	exp.DocumentStatus = "uncleared"
	exp.Amount = float64(10000)
	exp.RefType = "purchase_invoice"
	exp.RefID = uint64(inv.ID)
	exp.Save()

	// setting body
	scenario := tester.D{}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(exp.ID)+"/approve").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

	exp.Read("ID")
	assert.Equal(t, "cleared", exp.DocumentStatus)

	inv.Read("ID")
	assert.Equal(t, "finished", inv.DocumentStatus)
	assert.Equal(t, float64(50000), inv.TotalPaid)

	po.Read("ID")
	assert.Equal(t, "finished", po.DocumentStatus)
	assert.Equal(t, "finished", po.InvoiceStatus)
	assert.Equal(t, float64(50000), po.TotalPaid)

	supplier.Read("ID")
	assert.Equal(t, float64(10000), supplier.TotalCredit)
}

// TestHandler_URLMappingApproveFinanceExpenseSuccess2 Berhasil approve untuk finance expense po invoice active
func TestHandler_URLMappingApproveFinanceExpenseSuccess2(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy
	inv := model.DummyPurchaseInvoice()
	inv.TotalAmount = float64(50000)
	inv.TotalPaid = float64(40000)
	inv.DocumentStatus = "active"
	inv.IsDeleted = int8(0)
	inv.Save()

	po := inv.PurchaseOrder
	po.DocumentStatus = "active"
	po.TotalPaid = float64(40000)
	po.TotalCharge = float64(50000)
	po.InvoiceStatus = "active"
	po.ReceivingStatus = "finished"
	po.IsDeleted = int8(0)
	po.Save()

	supplier := po.Supplier
	supplier.PartnershipType = "supplier"
	supplier.IsDeleted = int8(0)
	supplier.TotalCredit = float64(20000)
	supplier.Save()

	exp1 := model.DummyFinanceExpense()
	exp1.IsDeleted = int8(0)
	exp1.DocumentStatus = "cleared"
	exp1.Amount = float64(20000)
	exp1.RefType = "purchase_invoice"
	exp1.RefID = uint64(inv.ID)
	exp1.Save()

	exp2 := model.DummyFinanceExpense()
	exp2.IsDeleted = int8(0)
	exp2.DocumentStatus = "cleared"
	exp2.Amount = float64(20000)
	exp2.RefType = "purchase_invoice"
	exp2.RefID = uint64(inv.ID)
	exp2.Save()

	exp := model.DummyFinanceExpense()
	exp.IsDeleted = int8(0)
	exp.DocumentStatus = "uncleared"
	exp.Amount = float64(5000)
	exp.RefType = "purchase_invoice"
	exp.RefID = uint64(inv.ID)
	exp.Save()

	// setting body
	scenario := tester.D{}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(exp.ID)+"/approve").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

	exp.Read("ID")
	assert.Equal(t, "cleared", exp.DocumentStatus)

	inv.Read("ID")
	assert.Equal(t, "active", inv.DocumentStatus)
	assert.Equal(t, float64(45000), inv.TotalPaid)

	po.Read("ID")
	assert.Equal(t, "active", po.DocumentStatus)
	assert.Equal(t, "active", po.InvoiceStatus)
	assert.Equal(t, float64(45000), po.TotalPaid)

	supplier.Read("ID")
	assert.Equal(t, float64(15000), supplier.TotalCredit)
}

// TestHandler_URLMappingApproveFinanceExpenseSuccess3 Berhasil approve untuk finance expense sales return finish
func TestHandler_URLMappingApproveFinanceExpenseSuccess3(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy
	ret := model.DummySalesReturn()
	ret.TotalAmount = float64(50000)
	ret.DocumentStatus = "active"
	ret.IsDeleted = int8(0)
	ret.Save()

	exp1 := model.DummyFinanceExpense()
	exp1.IsDeleted = int8(0)
	exp1.DocumentStatus = "cleared"
	exp1.Amount = float64(15000)
	exp1.RefType = "sales_return"
	exp1.RefID = uint64(ret.ID)
	exp1.Save()

	exp2 := model.DummyFinanceExpense()
	exp2.IsDeleted = int8(0)
	exp2.DocumentStatus = "cleared"
	exp2.Amount = float64(30000)
	exp2.RefType = "sales_return"
	exp2.RefID = uint64(ret.ID)
	exp2.Save()

	exp3 := model.DummyFinanceExpense()
	exp3.IsDeleted = int8(1)
	exp3.DocumentStatus = "cleared"
	exp3.Amount = float64(20000)
	exp3.RefType = "sales_return"
	exp3.RefID = uint64(ret.ID)
	exp3.Save()

	exp := model.DummyFinanceExpense()
	exp.IsDeleted = int8(0)
	exp.DocumentStatus = "uncleared"
	exp.Amount = float64(5000)
	exp.RefType = "sales_return"
	exp.RefID = uint64(ret.ID)
	exp.Save()

	// setting body
	scenario := tester.D{}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(exp.ID)+"/approve").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

	exp.Read("ID")
	assert.Equal(t, "cleared", exp.DocumentStatus)

	exp.Read("ID")
	assert.Equal(t, "cleared", exp.DocumentStatus)

	ret.Read("ID")
	assert.Equal(t, "finished", ret.DocumentStatus)
}

// TestHandler_URLMappingApproveFinanceExpenseSuccess4 Berhasil approve untuk finance expense sales return active
func TestHandler_URLMappingApproveFinanceExpenseSuccess4(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy
	ret := model.DummySalesReturn()
	ret.TotalAmount = float64(50000)
	ret.DocumentStatus = "active"
	ret.IsDeleted = int8(0)
	ret.Save()

	exp1 := model.DummyFinanceExpense()
	exp1.IsDeleted = int8(0)
	exp1.DocumentStatus = "cleared"
	exp1.Amount = float64(15000)
	exp1.RefType = "sales_return"
	exp1.RefID = uint64(ret.ID)
	exp1.Save()

	exp2 := model.DummyFinanceExpense()
	exp2.IsDeleted = int8(0)
	exp2.DocumentStatus = "cleared"
	exp2.Amount = float64(30000)
	exp2.RefType = "sales_return"
	exp2.RefID = uint64(ret.ID)
	exp2.Save()

	exp3 := model.DummyFinanceExpense()
	exp3.IsDeleted = int8(1)
	exp3.DocumentStatus = "cleared"
	exp3.Amount = float64(20000)
	exp3.RefType = "sales_return"
	exp3.RefID = uint64(ret.ID)
	exp3.Save()

	exp := model.DummyFinanceExpense()
	exp.IsDeleted = int8(0)
	exp.DocumentStatus = "uncleared"
	exp.Amount = float64(2000)
	exp.RefType = "sales_return"
	exp.RefID = uint64(ret.ID)
	exp.Save()

	// setting body
	scenario := tester.D{}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(exp.ID)+"/approve").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

	exp.Read("ID")
	assert.Equal(t, "cleared", exp.DocumentStatus)

	ret.Read("ID")
	assert.Equal(t, "active", ret.DocumentStatus)
}

// TestHandler_URLMappingApproveFinanceExpenseFail gagal approve untuk finance expense sales return cleared
func TestHandler_URLMappingApproveFinanceExpenseFail(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy
	ret := model.DummySalesReturn()
	ret.TotalAmount = float64(50000)
	ret.DocumentStatus = "active"
	ret.IsDeleted = int8(0)
	ret.Save()

	exp1 := model.DummyFinanceExpense()
	exp1.IsDeleted = int8(0)
	exp1.DocumentStatus = "cleared"
	exp1.Amount = float64(15000)
	exp1.RefType = "sales_return"
	exp1.RefID = uint64(ret.ID)
	exp1.Save()

	exp2 := model.DummyFinanceExpense()
	exp2.IsDeleted = int8(0)
	exp2.DocumentStatus = "cleared"
	exp2.Amount = float64(30000)
	exp2.RefType = "sales_return"
	exp2.RefID = uint64(ret.ID)
	exp2.Save()

	exp3 := model.DummyFinanceExpense()
	exp3.IsDeleted = int8(1)
	exp3.DocumentStatus = "cleared"
	exp3.Amount = float64(20000)
	exp3.RefType = "sales_return"
	exp3.RefID = uint64(ret.ID)
	exp3.Save()

	exp := model.DummyFinanceExpense()
	exp.IsDeleted = int8(0)
	exp.DocumentStatus = "cleared"
	exp.Amount = float64(2000)
	exp.RefType = "sales_return"
	exp.RefID = uint64(ret.ID)
	exp.Save()

	// setting body
	scenario := tester.D{}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(exp.ID)+"/approve").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingApproveFinanceExpenseFail2 gagal approve untuk finance expense sales return deleted
func TestHandler_URLMappingApproveFinanceExpenseFail2(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// buat dummy
	ret := model.DummySalesReturn()
	ret.TotalAmount = float64(50000)
	ret.DocumentStatus = "active"
	ret.IsDeleted = int8(0)
	ret.Save()

	exp1 := model.DummyFinanceExpense()
	exp1.IsDeleted = int8(0)
	exp1.DocumentStatus = "cleared"
	exp1.Amount = float64(15000)
	exp1.RefType = "sales_return"
	exp1.RefID = uint64(ret.ID)
	exp1.Save()

	exp2 := model.DummyFinanceExpense()
	exp2.IsDeleted = int8(0)
	exp2.DocumentStatus = "cleared"
	exp2.Amount = float64(30000)
	exp2.RefType = "sales_return"
	exp2.RefID = uint64(ret.ID)
	exp2.Save()

	exp3 := model.DummyFinanceExpense()
	exp3.IsDeleted = int8(1)
	exp3.DocumentStatus = "cleared"
	exp3.Amount = float64(20000)
	exp3.RefType = "sales_return"
	exp3.RefID = uint64(ret.ID)
	exp3.Save()

	exp := model.DummyFinanceExpense()
	exp.IsDeleted = int8(1)
	exp.DocumentStatus = "uncleared"
	exp.Amount = float64(2000)
	exp.RefType = "sales_return"
	exp.RefID = uint64(ret.ID)
	exp.Save()

	// setting body
	scenario := tester.D{}

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.PUT("/v1/finance-expense/"+common.Encrypt(exp.ID)+"/approve").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusNotFound, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}
