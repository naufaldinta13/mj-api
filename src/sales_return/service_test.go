package salesReturn

import (
	"testing"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestGetSalesReturnNoData(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM sales_return").Exec()
	rq := orm.RequestQuery{}

	m, total, e := GetSalesReturn(&rq)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, m)
	assert.NoError(t, e)
}

func TestGetSalesReturn(t *testing.T) {
	model.DummySalesReturn()
	qs := orm.RequestQuery{}
	_, _, e := GetSalesReturn(&qs)
	assert.NoError(t, e, "Data should be exists.")
}

func TestShowSalesReturn(t *testing.T) {
	_, e := ShowSalesReturn("id", 1000)
	assert.Error(t, e, "Response should be error, beacuse there are no data yet.")

	c := model.DummySalesReturn()
	cd, e := ShowSalesReturn("id", c.ID)
	assert.NoError(t, e, "Data should be exists.")
	assert.Equal(t, c.ID, cd.ID, "ID Response should be a same.")
}

func dsr() *model.SalesReturn {
	sr := model.DummySalesReturn()
	sr.DocumentStatus = "cancelled"
	sr.Save()

	fe := model.DummyFinanceExpense()
	fe.IsDeleted = 0
	fe.RefType = "sales_return"
	fe.RefID = uint64(sr.ID)
	fe.Save()

	fe1 := model.DummyFinanceExpense()
	fe1.IsDeleted = 0
	fe1.RefType = "sales_return"
	fe1.RefID = uint64(sr.ID)
	fe1.Save()

	fe2 := model.DummyFinanceExpense()
	fe2.IsDeleted = 0
	fe2.RefType = "sales_return"
	fe2.RefID = uint64(sr.ID)
	fe2.Save()

	return sr
}

func TestUpdateExpense(t *testing.T) {
	sr := dsr()

	e := UpdateExpense(sr)
	assert.NoError(t, e, "Data should be exists.")

}
