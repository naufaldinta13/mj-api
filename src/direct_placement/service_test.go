// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package directPlacement

import (
	"testing"
	"time"

	"git.qasico.com/mj/api/datastore/model"

	"github.com/stretchr/testify/assert"
)

//// Test Show direct placement ///////////////////////////

// TestShowDirectPlacement test show direct placement ,success
func TestShowDirectPlacement(t *testing.T) {
	// buat dummy
	direct := model.DummyDirectPlacement()
	directItem := model.DummyDirectPlacementItem()
	directItem.DirectPlacment = direct
	directItem.Save()

	// test
	m, e := ShowDirectPlacement("id", direct.ID)
	assert.NoError(t, e)
	assert.NotEmpty(t, m)
	assert.Equal(t, direct.Note, m.Note)
	assert.Equal(t, direct.CreatedBy.ID, m.CreatedBy.ID)
	assert.NotEmpty(t, m.DirectPlacementItems)
	assert.Equal(t, int(1), len(m.DirectPlacementItems))
	for _, u := range m.DirectPlacementItems {
		assert.Equal(t, directItem.ID, u.ID)
		assert.Equal(t, directItem.Quantity, u.Quantity)
		assert.Equal(t, directItem.ItemVariant.Note, u.ItemVariant.Note)
	}
}

// TestShowDirectPlacementNoItem test show direct placement tanpa directPlacement item ,success
func TestShowDirectPlacementNoItem(t *testing.T) {
	// buat dummy
	direct := model.DummyDirectPlacement()

	// test
	m, e := ShowDirectPlacement("id", direct.ID)
	assert.NoError(t, e)
	assert.NotEmpty(t, m)
	assert.Equal(t, direct.Note, m.Note)
	assert.Equal(t, direct.CreatedBy.ID, m.CreatedBy.ID)
	assert.Empty(t, m.DirectPlacementItems)
	assert.Equal(t, int(0), len(m.DirectPlacementItems))
}

// TestShowDirectPlacementWrongParam test show direct placement dengan param salah ,fail
func TestShowDirectPlacementWrongParam(t *testing.T) {
	// test
	m, e := ShowDirectPlacement("id", int64(99999999))
	assert.Error(t, e)
	assert.Empty(t, m)
}

//// Test create direct placement ///////////////////////////

// TestCreateDirectPlacement test create direct placement ,success
func TestCreateDirectPlacement(t *testing.T) {
	// buat fake data
	var placementItem []*model.DirectPlacementItem

	dirItem := &model.DirectPlacementItem{
		ItemVariant: model.DummyItemVariant(),
		Quantity:    float32(10),
		TotalPrice:  float64(20000),
		UnitPrice:   float64(2000),
	}
	placementItem = append(placementItem, dirItem)

	m := model.DirectPlacement{
		CreatedBy:            model.DummyUser(),
		CreatedAt:            time.Now(),
		Note:                 "note",
		DirectPlacementItems: placementItem,
	}
	// test
	e := CreateDirectPlacement(&m)
	assert.NoError(t, e)
	assert.NotEqual(t, int64(0), m.ID)

	for _, u := range m.DirectPlacementItems {
		assert.NotEmpty(t, u.DirectPlacment)
		assert.Equal(t, m.ID, u.DirectPlacment.ID)
	}
}

// TestCreateDirectPlacementFail test create direct placement ,fail
func TestCreateDirectPlacementFail(t *testing.T) {
	// buat fake data
	var placementItem []*model.DirectPlacementItem

	dirItem := &model.DirectPlacementItem{
		Quantity:   float32(10),
		TotalPrice: float64(20000),
		UnitPrice:  float64(2000),
	}
	placementItem = append(placementItem, dirItem)

	m := model.DirectPlacement{
		CreatedBy:            model.DummyUser(),
		CreatedAt:            time.Now(),
		Note:                 "note",
		DirectPlacementItems: placementItem,
	}
	// test
	e := CreateDirectPlacement(&m)
	assert.Error(t, e)
}
