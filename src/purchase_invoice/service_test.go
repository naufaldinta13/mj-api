// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchaseInvoice

import (
	"testing"
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestCalculateTotalAmountPI(t *testing.T) {
	po := model.DummyPurchaseOrder()
	po.Save()

	pi := model.DummyPurchaseInvoice()
	pi.TotalAmount = 1000
	pi.IsDeleted = 0
	pi.PurchaseOrder = po
	pi.Save()

	pi2 := model.DummyPurchaseInvoice()
	pi2.TotalAmount = 1000
	pi2.PurchaseOrder = po
	pi2.IsDeleted = 0
	pi2.Save()

	//test no error
	total, e := CalculateTotalAmountPI(po.ID, 0)
	assert.NoError(t, e, "no errror")
	assert.Equal(t, float64(2000), total)

	//test calculate for update Pi
	totalu, e := CalculateTotalAmountPI(po.ID, pi2.ID)
	assert.NoError(t, e, "no errror")
	assert.Equal(t, float64(1000), totalu)

}

func TestCreatePurchaseInvoice(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)

	po := model.DummyPurchaseOrder()
	po.TotalCharge = 4000
	po.Save()
	pox := common.Encrypt(po.ID)

	pi := model.DummyPurchaseInvoice()
	pi.TotalAmount = 1000
	pi.PurchaseOrder = po
	pi.Save()

	pi2 := model.DummyPurchaseInvoice()
	pi2.TotalAmount = 1000
	pi2.PurchaseOrder = po
	pi2.Save()

	r := &createRequest{
		RecognitionDate: time.Now(),
		PurchaseOrder:   pox,
		DueDate:         time.Now(),
		TotalAmount:     3000,
		Note:            "AHAHAHA",
		SessionData:     sd,
		BillingAddress:  "JEKARDAH",
	}

	//test no error
	pi, e := CreatePurchaseInvoice(r)
	assert.NoError(t, e, "no errror")
	assert.Equal(t, r.Note, pi.Note)

	po.Read("ID")
	assert.Equal(t, "active", po.InvoiceStatus)
}

func TestGetPurchaseInvoiceNoData(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM purchase_invoice").Exec()
	rq := orm.RequestQuery{}

	m, total, e := GetPurchaseInvoice(&rq)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, m)
	assert.NoError(t, e)
}

func TestGetPurchaseInvoice(t *testing.T) {
	model.DummyPurchaseInvoice()
	qs := orm.RequestQuery{}
	_, _, e := GetPurchaseInvoice(&qs)
	assert.NoError(t, e, "Data should be exists.")
}

func TestShowPurchaseInvoice(t *testing.T) {
	_, e := ShowPurchaseInvoice("id", 1000)
	assert.Error(t, e, "Response should be error, beacuse there are no data yet.")

	c := model.DummyPurchaseInvoice()
	c.IsDeleted = 0
	c.Save()

	cd, e := ShowPurchaseInvoice("id", c.ID)
	assert.NoError(t, e, "Data should be exists.")
	assert.Equal(t, c.ID, cd.ID, "ID Response should be a same.")
}

func TestUpdatePurchaseInvoice(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sd, _ := auth.Login(user)

	po := model.DummyPurchaseOrder()
	po.TotalCharge = 4000
	po.Save()
	//pox := common.Encrypt(po.ID)

	pi := model.DummyPurchaseInvoice()
	pi.TotalAmount = 1500
	pi.PurchaseOrder = po
	pi.Save()

	pi2 := model.DummyPurchaseInvoice()
	pi2.TotalAmount = 1000
	pi2.PurchaseOrder = po
	pi2.Save()

	r := &updateRequest{
		PI:              pi,
		SessionData:     sd,
		RecognitionDate: time.Now(),
		TotalAmount:     3000,
		Note:            "AHAHAHA",
	}

	//test no error
	pix, e := UpdatePurchaseInvoice(r)
	assert.NoError(t, e, "no errror")
	assert.NotEqual(t, r.TotalAmount, pi.TotalAmount)
	assert.Equal(t, r.Note, pix.Note)
}
