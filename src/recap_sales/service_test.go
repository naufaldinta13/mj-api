// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package recapSales

import (
	"testing"
	"time"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

// TestCreateRecapSales test membuat recap sales beserta recap sales item
func TestCreateRecapSales(t *testing.T) {
	var recapItems []*model.RecapSalesItem
	// buat dummy partner
	partner := model.DummyPartnership()
	partner.IsDeleted = int8(0)
	partner.PartnershipType = "customer"
	partner.Save()
	// buat dummy so
	so1 := model.DummySalesOrder()
	so1.IsDeleted = int8(0)
	so1.DocumentStatus = "finished"
	so1.TotalCharge = float64(20000)
	so1.Customer = partner
	so1.Save()

	// init data
	recapItem := &model.RecapSalesItem{SalesOrder: so1}
	recapItems = append(recapItems, recapItem)
	recap := model.RecapSales{
		Partnership:     partner,
		Code:            common.RandomStr(5),
		TotalAmount:     float64(20000),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		RecapSalesItems: recapItems,
	}
	// test
	e := CreateRecapSales(&recap)
	assert.NoError(t, e)
	rec := &model.RecapSales{Partnership: partner, TotalAmount: float64(20000)}
	err := rec.Read("Partnership", "TotalAmount")
	assert.NoError(t, err)
	recItm := &model.RecapSalesItem{SalesOrder: so1, RecapSales: rec}
	err2 := recItm.Read("SalesOrder", "RecapSales")
	assert.NoError(t, err2)
}

// TestCreateRecapSalesFail test membuat recap sales beserta recap sales item gagal
func TestCreateRecapSalesFail(t *testing.T) {
	var recapItems []*model.RecapSalesItem
	// buat dummy partner
	partner := model.DummyPartnership()
	partner.IsDeleted = int8(0)
	partner.PartnershipType = "customer"
	partner.Save()
	// buat dummy so
	so1 := model.DummySalesOrder()
	so1.IsDeleted = int8(0)
	so1.DocumentStatus = "finished"
	so1.TotalCharge = float64(20000)
	so1.Customer = partner
	so1.Save()

	// init data
	recapItem := &model.RecapSalesItem{}
	recapItems = append(recapItems, recapItem)
	recap := model.RecapSales{
		Partnership:     partner,
		Code:            common.RandomStr(5),
		TotalAmount:     float64(20000),
		CreatedBy:       model.DummyUser(),
		CreatedAt:       time.Now(),
		RecapSalesItems: recapItems,
	}
	// test
	e := CreateRecapSales(&recap)
	assert.Error(t, e)

}

// TestGetRecapSales test get all recap sales
func TestGetRecapSales(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM recap_sales").Exec()
	// buat dummy
	model.DummyRecapSales()
	qs := orm.RequestQuery{}
	m, tot, e := GetRecapSales(&qs)
	// test
	assert.NoError(t, e)
	assert.Equal(t, int64(1), tot)
	assert.NotEmpty(t, m)
}

// TestShowRecapSales test get detail recap sales
func TestShowRecapSales(t *testing.T) {
	// buat dummy
	rec := model.DummyRecapSales()
	recItm := model.DummyRecapSalesItem()
	recItm.RecapSales = rec
	recItm.Save()
	// test
	m, e := ShowRecapSales("id", rec.ID)
	assert.NoError(t, e)
	assert.NotEmpty(t, m)
	assert.Equal(t, rec.TotalAmount, m.TotalAmount)
	assert.Equal(t, rec.Partnership.Code, m.Partnership.Code)
	assert.Equal(t, int(1), len(m.RecapSalesItems))
	for _, u := range m.RecapSalesItems {
		assert.Equal(t, recItm.SalesOrder.Code, u.SalesOrder.Code)
	}
}

// TestShowRecapSalesFailID test get detail recap sales id salah
func TestShowRecapSalesFailID(t *testing.T) {
	// test
	m, e := ShowRecapSales("id", 999999999)
	assert.Error(t, e)
	assert.Empty(t, m)
}
