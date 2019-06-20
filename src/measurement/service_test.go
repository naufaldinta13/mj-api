package measurement

import (
	"fmt"
	"testing"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/common/faker"
	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestGetMeasurementNoData(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM measurement").Exec()
	rq := orm.RequestQuery{}

	m, total, e := GetMeasurement(&rq)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, m)
	assert.NoError(t, e)
}

func TestGetMeasurement(t *testing.T) {
	model.DummyMeasurement()
	qs := orm.RequestQuery{}
	_, _, e := GetMeasurement(&qs)
	assert.NoError(t, e, "Data should be exists.")
}

func msrmnt() *model.Measurement {
	var m model.Measurement
	faker.Fill(&m, "ID")

	m.IsDeleted = 0
	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

func TestShowMeasurement(t *testing.T) {
	_, e := ShowMeasurement("id", 1000)
	assert.Error(t, e, "Response should be error, beacuse there are no data yet.")

	c := msrmnt()
	cd, e := ShowMeasurement("id", c.ID)
	assert.NoError(t, e, "Data should be exists.")
	assert.Equal(t, c.ID, cd.ID, "ID Response should be a same.")
}

func TestGetMeasurementByName(t *testing.T) {
	c := msrmnt()
	cd, e := getMeasurementByName(c.MeasurementName)
	assert.NoError(t, e, "Data should be exists.")
	assert.Equal(t, c.MeasurementName, cd.MeasurementName, "ID Response should be a same.")
}

func TestIsMeasurementUsed(t *testing.T) {
	c := msrmnt()
	dmiv := model.DummyItemVariant()
	dmiv.Measurement = c
	dmiv.Save()
	total := isMeasurementUsed(c.ID)
	assert.Equal(t, true, total)
}
