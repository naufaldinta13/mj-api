// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package recapSales_test

import (
	"fmt"
	"os"
	"testing"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"
	"git.qasico.com/mj/api/test"

	"net/http"

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
	test.DataCleanUp()

	os.Exit(res)
}

////Post Recap Sales/////////////////////////////////////////////////////////

// TestHandler_URLMappingPostRecapSalesSuccess test create recap sales,success
func TestHandler_URLMappingPostRecapSalesSuccess(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM recap_sales").Exec()
	// buat dummy partner
	partnr := model.DummyPartnership()
	partnr.IsDeleted = int8(0)
	partnr.PartnershipType = "customer"
	partnr.Save()
	// buat dummy so
	so1 := model.DummySalesOrder()
	so1.IsDeleted = int8(0)
	so1.DocumentStatus = "finished"
	so1.TotalCharge = float64(20000)
	so1.Customer = partnr
	so1.Save()
	so2 := model.DummySalesOrder()
	so2.IsDeleted = int8(0)
	so2.DocumentStatus = "finished"
	so2.TotalCharge = float64(30000)
	so2.Customer = partnr
	so2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"partnership": common.Encrypt(partnr.ID),
		"recap_sales_items": []tester.D{
			{
				"sales_order": common.Encrypt(so1.ID),
			},
			{
				"sales_order": common.Encrypt(so2.ID),
			},
		},
	}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/recap-sales").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
	rec := &model.RecapSales{Partnership: &model.Partnership{ID: partnr.ID}}
	e := rec.Read("Partnership")
	assert.NoError(t, e)
	assert.Equal(t, float64(50000), rec.TotalAmount)
	assert.Equal(t, sd.User.ID, rec.CreatedBy.ID)
	// recap sale item
	recItm1 := &model.RecapSalesItem{RecapSales: &model.RecapSales{ID: rec.ID}, SalesOrder: &model.SalesOrder{ID: so1.ID}}
	err1 := recItm1.Read("RecapSales", "SalesOrder")
	assert.NoError(t, err1)
	recItm2 := &model.RecapSalesItem{RecapSales: &model.RecapSales{ID: rec.ID}, SalesOrder: &model.SalesOrder{ID: so2.ID}}
	err2 := recItm2.Read("RecapSales", "SalesOrder")
	assert.NoError(t, err2)
}

// TestHandler_URLMappingPostRecapSalesFailNoToken test create recap sales tanpa token,fail
func TestHandler_URLMappingPostRecapSalesFailNoToken(t *testing.T) {
	// buat dummy partner
	partnr := model.DummyPartnership()
	partnr.IsDeleted = int8(0)
	partnr.PartnershipType = "customer"
	partnr.Save()
	// buat dummy so
	so1 := model.DummySalesOrder()
	so1.IsDeleted = int8(0)
	so1.DocumentStatus = "finished"
	so1.TotalCharge = float64(20000)
	so1.Customer = partnr
	so1.Save()
	so2 := model.DummySalesOrder()
	so2.IsDeleted = int8(0)
	so2.DocumentStatus = "finished"
	so2.TotalCharge = float64(30000)
	so2.Customer = partnr
	so2.Save()

	// setting body
	scenario := tester.D{
		"partnership": common.Encrypt(partnr.ID),
		"recap_sales_items": []tester.D{
			{
				"sales_order": common.Encrypt(so1.ID),
			},
			{
				"sales_order": common.Encrypt(so2.ID),
			},
		},
	}
	// test
	ng := tester.New()
	ng.POST("/v1/recap-sales").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})
}

// TestHandler_URLMappingPostRecapSalesFailPartnerUndecrypt test create recap sales dengan partner id un decrypt,fail
func TestHandler_URLMappingPostRecapSalesFailPartnerUndecrypt(t *testing.T) {
	// buat dummy partner
	partnr := model.DummyPartnership()
	partnr.IsDeleted = int8(0)
	partnr.PartnershipType = "customer"
	partnr.Save()
	// buat dummy so
	so1 := model.DummySalesOrder()
	so1.IsDeleted = int8(0)
	so1.DocumentStatus = "finished"
	so1.TotalCharge = float64(20000)
	so1.Customer = partnr
	so1.Save()
	so2 := model.DummySalesOrder()
	so2.IsDeleted = int8(0)
	so2.DocumentStatus = "finished"
	so2.TotalCharge = float64(30000)
	so2.Customer = partnr
	so2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"partnership": "aaaaaaaa",
		"recap_sales_items": []tester.D{
			{
				"sales_order": common.Encrypt(so1.ID),
			},
			{
				"sales_order": common.Encrypt(so2.ID),
			},
		},
	}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/recap-sales").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

} // TestHandler_URLMappingPostRecapSalesFailSOSame test create recap sales dengan sales order sama,fail
func TestHandler_URLMappingPostRecapSalesFailSOSame(t *testing.T) {
	// buat dummy partner
	partnr := model.DummyPartnership()
	partnr.IsDeleted = int8(0)
	partnr.PartnershipType = "customer"
	partnr.Save()
	// buat dummy so
	so1 := model.DummySalesOrder()
	so1.IsDeleted = int8(0)
	so1.DocumentStatus = "approved_cancel"
	so1.TotalCharge = float64(20000)
	so1.Customer = partnr
	so1.Save()
	so2 := model.DummySalesOrder()
	so2.IsDeleted = int8(0)
	so2.DocumentStatus = "finished"
	so2.TotalCharge = float64(30000)
	so2.Customer = partnr
	so2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"partnership": common.Encrypt(partnr.ID),
		"recap_sales_items": []tester.D{
			{
				"sales_order": common.Encrypt(so1.ID),
			},
			{
				"sales_order": common.Encrypt(so1.ID),
			},
		},
	}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/recap-sales").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMappingPostRecapSalesFailPartnerNotExist test create recap sales dengan partner id not exist,fail
func TestHandler_URLMappingPostRecapSalesFailPartnerNotExist(t *testing.T) {
	// buat dummy partner
	partnr := model.DummyPartnership()
	partnr.IsDeleted = int8(0)
	partnr.PartnershipType = "customer"
	partnr.Save()
	// buat dummy so
	so1 := model.DummySalesOrder()
	so1.IsDeleted = int8(0)
	so1.DocumentStatus = "finished"
	so1.TotalCharge = float64(20000)
	so1.Customer = partnr
	so1.Save()
	so2 := model.DummySalesOrder()
	so2.IsDeleted = int8(0)
	so2.DocumentStatus = "finished"
	so2.TotalCharge = float64(30000)
	so2.Customer = partnr
	so2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"partnership": "99999999",
		"recap_sales_items": []tester.D{
			{
				"sales_order": common.Encrypt(so1.ID),
			},
			{
				"sales_order": common.Encrypt(so2.ID),
			},
		},
	}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/recap-sales").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMappingPostRecapSalesFailSalesUndecrypt test create recap sales dengan sales id un decrypt,fail
func TestHandler_URLMappingPostRecapSalesFailSalesUndecrypt(t *testing.T) {
	// buat dummy partner
	partnr := model.DummyPartnership()
	partnr.IsDeleted = int8(0)
	partnr.PartnershipType = "customer"
	partnr.Save()
	// buat dummy so
	so1 := model.DummySalesOrder()
	so1.IsDeleted = int8(0)
	so1.DocumentStatus = "finished"
	so1.TotalCharge = float64(20000)
	so1.Customer = partnr
	so1.Save()
	so2 := model.DummySalesOrder()
	so2.IsDeleted = int8(0)
	so2.DocumentStatus = "finished"
	so2.TotalCharge = float64(30000)
	so2.Customer = partnr
	so2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"partnership": common.Encrypt(partnr.ID),
		"recap_sales_items": []tester.D{
			{
				"sales_order": "aaaaa",
			},
			{
				"sales_order": common.Encrypt(so2.ID),
			},
		},
	}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/recap-sales").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMappingPostRecapSalesFailSalesNotExist test create recap sales dengan sales id not exist,fail
func TestHandler_URLMappingPostRecapSalesFailSalesNotExist(t *testing.T) {
	// buat dummy partner
	partnr := model.DummyPartnership()
	partnr.IsDeleted = int8(0)
	partnr.PartnershipType = "customer"
	partnr.Save()
	// buat dummy so
	so1 := model.DummySalesOrder()
	so1.IsDeleted = int8(0)
	so1.DocumentStatus = "finished"
	so1.TotalCharge = float64(20000)
	so1.Customer = partnr
	so1.Save()
	so2 := model.DummySalesOrder()
	so2.IsDeleted = int8(0)
	so2.DocumentStatus = "finished"
	so2.TotalCharge = float64(30000)
	so2.Customer = partnr
	so2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"partnership": common.Encrypt(partnr.ID),
		"recap_sales_items": []tester.D{
			{
				"sales_order": "99999999999",
			},
			{
				"sales_order": common.Encrypt(so2.ID),
			},
		},
	}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/recap-sales").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMappingPostRecapSalesFailPartnerTypeSupplier test create recap sales dengan partner type supplier,fail
func TestHandler_URLMappingPostRecapSalesFailPartnerTypeSupplier(t *testing.T) {
	// buat dummy partner
	partnr := model.DummyPartnership()
	partnr.IsDeleted = int8(0)
	partnr.PartnershipType = "supplier"
	partnr.Save()
	// buat dummy so
	so1 := model.DummySalesOrder()
	so1.IsDeleted = int8(0)
	so1.DocumentStatus = "finished"
	so1.TotalCharge = float64(20000)
	so1.Customer = partnr
	so1.Save()
	so2 := model.DummySalesOrder()
	so2.IsDeleted = int8(0)
	so2.DocumentStatus = "finished"
	so2.TotalCharge = float64(30000)
	so2.Customer = partnr
	so2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"partnership": common.Encrypt(partnr.ID),
		"recap_sales_items": []tester.D{
			{
				"sales_order": common.Encrypt(so1.ID),
			},
			{
				"sales_order": common.Encrypt(so2.ID),
			},
		},
	}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/recap-sales").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMappingPostRecapSalesFailDiffPartner test create recap sales dengan partner berbeda,fail
func TestHandler_URLMappingPostRecapSalesFailDiffPartner(t *testing.T) {
	// buat dummy partner
	partnr := model.DummyPartnership()
	partnr.IsDeleted = int8(0)
	partnr.PartnershipType = "customer"
	partnr.Save()
	partner2 := model.DummyPartnership()
	partner2.IsDeleted = int8(0)
	partner2.PartnershipType = "customer"
	partner2.Save()
	// buat dummy so
	so1 := model.DummySalesOrder()
	so1.IsDeleted = int8(0)
	so1.DocumentStatus = "finished"
	so1.TotalCharge = float64(20000)
	so1.Customer = partner2
	so1.Save()
	so2 := model.DummySalesOrder()
	so2.IsDeleted = int8(0)
	so2.DocumentStatus = "finished"
	so2.TotalCharge = float64(30000)
	so2.Customer = partnr
	so2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"partnership": common.Encrypt(partnr.ID),
		"recap_sales_items": []tester.D{
			{
				"sales_order": common.Encrypt(so1.ID),
			},
			{
				"sales_order": common.Encrypt(so2.ID),
			},
		},
	}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/recap-sales").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

// TestHandler_URLMappingPostRecapSalesFailSOPartnerTypeSupplier test create recap sales dengan partner type so supplier,fail
func TestHandler_URLMappingPostRecapSalesFailSOPartnerTypeSupplier(t *testing.T) {
	// buat dummy partner
	partner0 := model.DummyPartnership()
	partner0.IsDeleted = int8(0)
	partner0.PartnershipType = "customer"
	partner0.Save()
	partnr := model.DummyPartnership()
	partnr.IsDeleted = int8(0)
	partnr.PartnershipType = "supplier"
	partnr.Save()
	// buat dummy so
	so1 := model.DummySalesOrder()
	so1.IsDeleted = int8(0)
	so1.DocumentStatus = "finished"
	so1.TotalCharge = float64(20000)
	so1.Customer = partnr
	so1.Save()
	so2 := model.DummySalesOrder()
	so2.IsDeleted = int8(0)
	so2.DocumentStatus = "finished"
	so2.TotalCharge = float64(30000)
	so2.Customer = partnr
	so2.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// setting body
	scenario := tester.D{
		"partnership": common.Encrypt(partner0.ID),
		"recap_sales_items": []tester.D{
			{
				"sales_order": common.Encrypt(so1.ID),
			},
			{
				"sales_order": common.Encrypt(so2.ID),
			},
		},
	}
	// test
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.POST("/v1/recap-sales").SetJSON(scenario).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(422), res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", scenario, res.Body.String()))
	})

}

////GET Recap Sales/////////////////////////////////////////////////////////

// TestHandler_URLMappingGetRecapSalesSuccess mengetest get recap sales,success
func TestHandler_URLMappingGetRecapSalesSuccess(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM recap_sales").Exec()

	// buat dummy recap
	model.DummyRecapSales()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/recap-sales"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/recap-sales", "GET"))
	})
}

////SHOW Recap Sales/////////////////////////////////////////////////////////

// TestHandler_URLMappingGetRecapSalesSuccess mengetest show recap sales,success
func TestHandler_URLMappingShowRecapSalesSuccess(t *testing.T) {
	// buat dummy recap
	rec := model.DummyRecapSales()
	recItm := model.DummyRecapSalesItem()
	recItm.RecapSales = rec
	recItm.Save()

	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/recap-sales/" + common.Encrypt(rec.ID)
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(200), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/recap-sales/"+common.Encrypt(rec.ID), "GET"))
	})
}

// TestHandler_URLMappingShowRecapSalesFailIDWrong mengetest show recap sales dengan id salah,fail
func TestHandler_URLMappingShowRecapSalesFailIDWrong(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/recap-sales/asas"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(400), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/recap-sales/asas", "GET"))
	})
}

// TestHandler_URLMappingShowRecapSalesFailIDNotExist mengetest show recap sales dengan id tidak ada,fail
func TestHandler_URLMappingShowRecapSalesFailIDNotExist(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	ng.Method = "GET"
	ng.Path = "/v1/recap-sales/999999"
	ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, int(404), res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", "/v1/recap-sales/999999", "GET"))
	})
}

// TestHandler_URLMappingCancelRecapSalesSucsess mengetest cancel recap sales sukses
func TestHandler_URLMappingCancelRecapSalesSucsess(t *testing.T) {
	// melakukan proses login
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)
	token := "Bearer " + sd.Token

	// dummy recap sales
	rs := model.DummyRecapSales()
	rs.IsDeleted = 0
	rs.Save()

	ersId := common.Encrypt(rs.ID)

	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token, "Content-type": "application/json"})
	request := tester.D{}
	ng.PUT("/v1/recap-sales/"+ersId+"/cancel").SetJSON(request).Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
		assert.Equal(t, http.StatusOK, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", request, res.Body.String()))
	})
	rs.Read("ID")
	assert.Equal(t, int8(1), rs.IsDeleted)
}
