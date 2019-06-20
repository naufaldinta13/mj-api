// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stockopname_test

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

	stockopname := model.DummyStockopname()
	estockID := common.Encrypt(stockopname.ID)

	//sukses
	var routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/stockopname", "GET", 200},
		{"/v1/stockopname/" + estockID, "GET", 200},
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

	// gagal
	stockopname.Delete()
	routers = []struct {
		endpoint string
		method   string
		expected int
	}{
		{"/v1/stockopname/" + estockID, "GET", 404},
	}

	for _, ep := range routers {
		ng.Method = ep.method
		ng.Path = ep.endpoint
		ng.Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
			assert.Equal(t, ep.expected, res.Code, fmt.Sprintf("Should has 'endpoint %s' with method '%s'", ep.endpoint, ep.method))
		})
	}
}

func TestHandler_URLMapping_POST(t *testing.T) {
	sysadmin := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(sysadmin)
	token := "Bearer " + sd.Token

	date := time.Now()
	ivStock := model.DummyItemVariantStock()
	ivStock.AvailableStock = 8
	ivStock.Save("AvailableStock")

	eivStock := common.Encrypt(ivStock.ID)

	var data = []struct {
		req      tester.D
		expected int
	}{
		//sukses
		// jika note stockopname kosong
		{tester.D{"recognition_date": date, "note": "",
			"stockopname_items": []tester.D{
				{"item_variant_stock_id": eivStock, "quantity": 12, "note": "asd"},
			},
		}, http.StatusOK},
		// lengkap
		{tester.D{"recognition_date": date, "note": "asd",
			"stockopname_items": []tester.D{
				{"item_variant_stock_id": eivStock, "quantity": 12, "note": "asd"},
			},
		}, http.StatusOK},
		// jika note item kosong
		{tester.D{"recognition_date": date, "note": "asd",
			"stockopname_items": []tester.D{
				{"item_variant_stock_id": eivStock, "quantity": 12, "note": ""},
			},
		}, http.StatusOK},
		// lengkap
		{tester.D{"recognition_date": date, "note": "asd",
			"stockopname_items": []tester.D{
				{"item_variant_stock_id": eivStock, "quantity": 12, "note": "asd"},
				{"item_variant_stock_id": common.Encrypt(model.DummyItemVariantStock().ID), "quantity": 12, "note": "asd"},
			},
		}, http.StatusOK},

		// gagal
		// apabila recogniton date tidak diisi
		{tester.D{"recognition_date": nil, "note": "asd",
			"stockopname_items": []tester.D{
				{"item_variant_stock_id": eivStock, "quantity": 12, "note": "asd"},
			},
		}, http.StatusUnprocessableEntity},

		// apabila stockopname items gak diisi
		{tester.D{"recognition_date": date, "note": "asd",
			"stockopname_items": nil,
		}, http.StatusUnprocessableEntity},

		// apabila item_variant_stock_id tidak diisi
		{tester.D{"recognition_date": date, "note": "asd",
			"stockopname_items": []tester.D{
				{"item_variant_stock_id": "", "quantity": 12, "note": "asd"},
			},
		}, http.StatusUnprocessableEntity},

		// apabila quantity barang kosong
		{tester.D{"recognition_date": date, "note": "asd",
			"stockopname_items": []tester.D{
				{"item_variant_stock_id": eivStock, "quantity": 0, "note": "asd"},
			},
		}, http.StatusOK},

		// apabila item_variant_stock_id yg dimasukkan tidak ada dalam database
		{tester.D{"recognition_date": date, "note": "asd",
			"stockopname_items": []tester.D{
				{"item_variant_stock_id": "123414123", "quantity": 12, "note": "asd"},
			},
		}, http.StatusUnprocessableEntity},

		// apabila item_variant_stock_id tidak valid
		{tester.D{"recognition_date": date, "note": "asd",
			"stockopname_items": []tester.D{
				{"item_variant_stock_id": "qwesa", "quantity": 12, "note": "asd"},
			},
		}, http.StatusUnprocessableEntity},

		// apabila item_variant_stock_id yg dimasukkan sama
		{tester.D{"recognition_date": date, "note": "asd",
			"stockopname_items": []tester.D{
				{"item_variant_stock_id": eivStock, "quantity": 12, "note": "asd"},
				{"item_variant_stock_id": eivStock, "quantity": 12, "note": "asd"},
			},
		}, http.StatusUnprocessableEntity},
	}
	ng := tester.New()
	ng.SetHeader(tester.H{"Authorization": token})
	for _, tes := range data {
		ng.POST("/v1/stockopname").
			SetJSON(tes.req).
			Run(test.Router(), func(res tester.HTTPResponse, req tester.HTTPRequest) {
				assert.Equal(t, tes.expected, res.Code, fmt.Sprintf("\nreason: Validation Not Matched,\ndata: %v , \nresponse: %v", tes.req, res.Body.String()))
			})
	}
}
