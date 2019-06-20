// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receiving

import (
	"testing"
	"time"

	"git.qasico.com/mj/api/datastore/model"
	"git.qasico.com/mj/api/src/auth"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestCreateReceivingService(t *testing.T) {
	user := model.DummyUserPriviledgeWithUsergroup(1)
	sdx, _ := auth.Login(user)

	po := model.DummyPurchaseOrder()
	po.Save()
	pox := common.Encrypt(po.ID)

	iv := model.DummyItemVariant()
	iv.AvailableStock = 10
	iv.Save()

	poi := model.DummyPurchaseOrderItem()
	poi.Quantity = 100
	poi.PurchaseOrder = po
	poi.ItemVariant = iv
	poi.Save()
	poix := common.Encrypt(poi.ID)

	wr := model.DummyWorkorderReceiving()
	wr.Save()

	wri := model.DummyWorkorderReceivingItem()
	wri.WorkorderReceiving = wr
	wri.PurchaseOrderItem = poi
	wri.Quantity = 30
	wri.Save()

	var item []receivingItem

	rItem := receivingItem{
		PurchaseOrderItem: poix,
		Quantity:          50,
	}
	item = append(item, rItem)

	r := &createRequest{
		RecognitionDate: time.Now(),
		PurchaseOrder:   pox,
		Note:            "AHAHAHA",
		Pic:             "ACC",
		SessionData:     sdx,
		ReceivingItem:   item,
	}

	//test no error
	wr, e := CreateReceiving(r)
	assert.NoError(t, e, "no errror")
	assert.Equal(t, r.Note, wr.Note)
}

func TestGetReceivingNoData(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM workorder_receiving").Exec()
	rq := orm.RequestQuery{}

	m, total, e := GetReceiving(&rq)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, m)
	assert.NoError(t, e)
}

func TestGetReceiving(t *testing.T) {
	model.DummyWorkorderReceiving()
	qs := orm.RequestQuery{}
	_, _, e := GetReceiving(&qs)
	assert.NoError(t, e, "Data should be exists.")
}
