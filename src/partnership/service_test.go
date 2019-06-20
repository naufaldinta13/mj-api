// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package partnership

import (
	"testing"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

func TestCalculationTotalDebtOneData(t *testing.T) {
	dummypartnership := model.DummyPartnership()
	dummypartnership.PartnershipType = "customer"
	dummypartnership.Save()
	dummysalesorder := model.DummySalesOrder()
	dummysalesorder.Customer = dummypartnership
	dummysalesorder.TotalCharge = float64(1500)
	dummysalesorder.DocumentStatus = "new"
	dummysalesorder.IsDeleted = 0
	dummysalesorder.Save()

	//ini test untuk error atau tidak nya
	//ini update filed total debt di table partnership
	e := CalculationTotalDebt(dummypartnership.ID)
	assert.NoError(t, e, "ini seharusnya tidak errror")

	//ini test untuk data partnershipnya sudah ke update atau belum
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_debt from partnership where id = ?", dummypartnership.ID).QueryRow(&total)
	assert.Equal(t, float64(1500), total)
}

func TestCalculationTotalDebtTwoData(t *testing.T) {
	dummypartnership := model.DummyPartnership()
	dummypartnership.PartnershipType = "customer"
	dummypartnership.Save()
	dummysalesorder := model.DummySalesOrder()
	dummysalesorder.Customer = dummypartnership
	dummysalesorder.TotalCharge = float64(1500)
	dummysalesorder.DocumentStatus = "new"
	dummysalesorder.IsDeleted = 0
	dummysalesorder.Save()
	dummysalesorder2 := model.DummySalesOrder()
	dummysalesorder2.Customer = dummypartnership
	dummysalesorder2.TotalCharge = float64(1500)
	dummysalesorder2.DocumentStatus = "new"
	dummysalesorder2.IsDeleted = 0
	dummysalesorder2.Save()

	//ini test untuk error atau tidak nya
	//ini update filed total debt di table partnership
	e := CalculationTotalDebt(dummypartnership.ID)
	assert.NoError(t, e, "ini seharusnya tidak errror")

	//ini test untuk data partnershipnya sudah ke update atau belum
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_debt from partnership where id = ?", dummypartnership.ID).QueryRow(&total)
	assert.Equal(t, float64(3000), total)
}

func TestCalculationTotalDebtTwoDataWithCanceledOrder(t *testing.T) {
	dummypartnership := model.DummyPartnership()
	dummypartnership.PartnershipType = "customer"
	dummypartnership.Save()
	dummysalesorder := model.DummySalesOrder()
	dummysalesorder.Customer = dummypartnership
	dummysalesorder.TotalCharge = float64(1500)
	dummysalesorder.DocumentStatus = "new"
	dummysalesorder.IsDeleted = 0
	dummysalesorder.Save()
	dummysalesorder2 := model.DummySalesOrder()
	dummysalesorder2.Customer = dummypartnership
	dummysalesorder2.TotalCharge = float64(1500)
	dummysalesorder2.DocumentStatus = "approved_cancel"
	dummysalesorder2.IsDeleted = 0
	dummysalesorder2.Save()

	//ini test untuk error atau tidak nya
	//ini update filed total debt di table partnership
	e := CalculationTotalDebt(dummypartnership.ID)
	assert.NoError(t, e, "ini seharusnya tidak errror")

	//ini test untuk data partnershipnya sudah ke update atau belum
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_debt from partnership where id = ?", dummypartnership.ID).QueryRow(&total)
	assert.Equal(t, float64(1500), total)
}

func TestCalculationTotalDebtTwoDataWithDeletedOrder(t *testing.T) {
	dummypartnership := model.DummyPartnership()
	dummypartnership.PartnershipType = "customer"
	dummypartnership.Save()
	dummysalesorder := model.DummySalesOrder()
	dummysalesorder.Customer = dummypartnership
	dummysalesorder.TotalCharge = float64(1500)
	dummysalesorder.DocumentStatus = "new"
	dummysalesorder.IsDeleted = 1
	dummysalesorder.Save()
	dummysalesorder2 := model.DummySalesOrder()
	dummysalesorder2.Customer = dummypartnership
	dummysalesorder2.TotalCharge = float64(1500)
	dummysalesorder2.DocumentStatus = "active"
	dummysalesorder2.IsDeleted = 0
	dummysalesorder2.Save()

	//ini test untuk error atau tidak nya
	//ini update filed total debt di table partnership
	e := CalculationTotalDebt(dummypartnership.ID)
	assert.NoError(t, e, "ini seharusnya tidak errror")

	//ini test untuk data partnershipnya sudah ke update atau belum
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_debt from partnership where id = ?", dummypartnership.ID).QueryRow(&total)
	assert.Equal(t, float64(1500), total)
}

func TestCalculationTotalDebtWithIsDeleted(t *testing.T) {
	dummypartnership := model.DummyPartnership()
	dummypartnership.PartnershipType = "customer"
	dummypartnership.Save()
	dummysalesorder := model.DummySalesOrder()
	dummysalesorder.Customer = dummypartnership
	dummysalesorder.TotalCharge = float64(1500)
	dummysalesorder.DocumentStatus = "new"
	dummysalesorder.IsDeleted = 1
	dummysalesorder.Save()

	//ini test untuk error atau tidak nya
	//ini update filed total debt di table partnership
	e := CalculationTotalDebt(dummypartnership.ID)
	assert.NoError(t, e, "ini seharusnya tidak errror")

	//ini test untuk data partnershipnya sudah ke update atau belum
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_debt from partnership where id = ?", dummypartnership.ID).QueryRow(&total)
	assert.Equal(t, float64(0), total)
}
func TestCalculationTotalDebtWithCanceled(t *testing.T) {
	dummypartnership := model.DummyPartnership()
	dummypartnership.PartnershipType = "customer"
	dummypartnership.Save()
	dummysalesorder := model.DummySalesOrder()
	dummysalesorder.Customer = dummypartnership
	dummysalesorder.TotalCharge = float64(1500)
	dummysalesorder.IsDeleted = 0
	dummysalesorder.DocumentStatus = "approved_cancel"
	dummysalesorder.Save()

	//ini test untuk error atau tidak nya
	//ini update filed total debt di table partnership
	e := CalculationTotalDebt(dummypartnership.ID)
	assert.NoError(t, e, "ini seharusnya tidak errror")

	//ini test untuk data partnershipnya sudah ke update atau belum
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_debt from partnership where id = ?", dummypartnership.ID).QueryRow(&total)
	assert.Equal(t, float64(0), total)
}

func TestCalculationTotalDebtWithCanceledAndActive(t *testing.T) {
	dummypartnership := model.DummyPartnership()
	dummypartnership.PartnershipType = "customer"
	dummypartnership.Save()
	dummysalesorder := model.DummySalesOrder()
	dummysalesorder.Customer = dummypartnership
	dummysalesorder.TotalCharge = float64(2000)
	dummysalesorder.TotalPaid = float64(1000)
	dummysalesorder.DocumentStatus = "active"
	dummysalesorder.IsDeleted = 0
	dummysalesorder.Save()
	dummysalesorder2 := model.DummySalesOrder()
	dummysalesorder2.Customer = dummypartnership
	dummysalesorder2.TotalCharge = float64(1500)
	dummysalesorder2.DocumentStatus = "approved_cancel"
	dummysalesorder2.IsDeleted = 0
	dummysalesorder2.Save()

	//ini test untuk error atau tidak nya
	//ini update filed total debt di table partnership
	e := CalculationTotalDebt(dummypartnership.ID)
	assert.NoError(t, e, "ini seharusnya tidak errror")

	//ini test untuk data partnershipnya sudah ke update atau belum
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_debt from partnership where id = ?", dummypartnership.ID).QueryRow(&total)
	assert.Equal(t, float64(1000), total)
}

func TestCalculationTotalSpendWithOneData(t *testing.T) {
	dummypartnershipcustomer := model.DummyPartnership()
	dummypartnershipcustomer.PartnershipType = "customer"
	dummypartnershipcustomer.Save()
	dummysalesorder := model.DummySalesOrder()
	dummysalesorder.Customer = dummypartnershipcustomer
	dummysalesorder.TotalCharge = float64(5000)
	dummysalesorder.TotalPaid = float64(2000)
	dummysalesorder.DocumentStatus = "new"
	dummysalesorder.IsDeleted = 0
	dummysalesorder.Save()

	e := CalculationTotalSpend(dummypartnershipcustomer.ID)

	assert.NoError(t, e)
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_spend from partnership where id = ?", dummypartnershipcustomer.ID).QueryRow(&total)
	assert.Equal(t, float64(5000), total)
}
func TestCalculationTotalSpendWithTwoData(t *testing.T) {
	dummypartnershipcustomer := model.DummyPartnership()
	dummypartnershipcustomer.PartnershipType = "customer"
	dummypartnershipcustomer.Save()
	dummysalesorder := model.DummySalesOrder()
	dummysalesorder.Customer = dummypartnershipcustomer
	dummysalesorder.TotalCharge = float64(5000)
	dummysalesorder.TotalPaid = float64(2000)
	dummysalesorder.DocumentStatus = "active"
	dummysalesorder.IsDeleted = 0
	dummysalesorder.Save()
	dummysalesorder2 := model.DummySalesOrder()
	dummysalesorder2.Customer = dummypartnershipcustomer
	dummysalesorder2.TotalCharge = float64(5000)
	dummysalesorder2.TotalPaid = float64(2000)
	dummysalesorder2.DocumentStatus = "new"
	dummysalesorder2.IsDeleted = 0
	dummysalesorder2.Save()

	e := CalculationTotalSpend(dummypartnershipcustomer.ID)

	assert.NoError(t, e)
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_spend from partnership where id = ?", dummypartnershipcustomer.ID).QueryRow(&total)
	assert.Equal(t, float64(10000), total)
}
func TestCalculationTotalSpendWithDocumentStatusActiveAndCancel(t *testing.T) {
	dummypartnershipcustomer := model.DummyPartnership()
	dummypartnershipcustomer.PartnershipType = "customer"
	dummypartnershipcustomer.Save()
	dummysalesorder := model.DummySalesOrder()
	dummysalesorder.Customer = dummypartnershipcustomer
	dummysalesorder.TotalCharge = float64(150000)
	dummysalesorder.TotalPaid = float64(2000)
	dummysalesorder.DocumentStatus = "active"
	dummysalesorder.IsDeleted = 0
	dummysalesorder.Save()
	dummysalesorder2 := model.DummySalesOrder()
	dummysalesorder2.Customer = dummypartnershipcustomer
	dummysalesorder2.TotalCharge = float64(15000)
	dummysalesorder2.TotalPaid = float64(2000)
	dummysalesorder2.DocumentStatus = "approved_cancel"
	dummysalesorder2.IsDeleted = 1
	dummysalesorder2.Save()

	e := CalculationTotalSpend(dummypartnershipcustomer.ID)

	assert.NoError(t, e)
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_spend from partnership where id = ?", dummypartnershipcustomer.ID).QueryRow(&total)
	assert.Equal(t, float64(150000), total)
}
func TestCalculationTotalSpendIsDeleted(t *testing.T) {
	dummypartnershipcustomer := model.DummyPartnership()
	dummypartnershipcustomer.PartnershipType = "customer"
	dummypartnershipcustomer.Save()
	dummysalesorder := model.DummySalesOrder()
	dummysalesorder.Customer = dummypartnershipcustomer
	dummysalesorder.TotalCharge = float64(5000)
	dummysalesorder.TotalPaid = float64(2000)
	dummysalesorder.DocumentStatus = "new"
	dummysalesorder.IsDeleted = 1
	dummysalesorder.Save()

	e := CalculationTotalSpend(dummypartnershipcustomer.ID)

	assert.NoError(t, e)
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_spend from partnership where id = ?", dummypartnershipcustomer.ID).QueryRow(&total)
	assert.Equal(t, float64(0), total)
}

func TestCalculationTotalExpenditureOneData(t *testing.T) {
	dummypartnershipsupplier := model.DummyPartnership()
	dummypartnershipsupplier.PartnershipType = "supplier"
	dummypartnershipsupplier.Save()
	dummypurchaseorder := model.DummyPurchaseOrder()
	dummypurchaseorder.Supplier = dummypartnershipsupplier
	dummypurchaseorder.TotalCharge = float64(6500)
	dummypurchaseorder.DocumentStatus = "new"
	dummypurchaseorder.IsDeleted = 0
	dummypurchaseorder.Save()

	e := CalculationTotalExpenditure(dummypartnershipsupplier.ID)
	assert.NoError(t, e)
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_expenditure from partnership where id = ?", dummypartnershipsupplier.ID).QueryRow(&total)
	assert.Equal(t, float64(6500), total)
}
func TestCalculationTotalExpenditureTwoData(t *testing.T) {
	dummypartnershipsupplier := model.DummyPartnership()
	dummypartnershipsupplier.PartnershipType = "supplier"
	dummypartnershipsupplier.Save()
	dummypurchaseorder := model.DummyPurchaseOrder()
	dummypurchaseorder.Supplier = dummypartnershipsupplier
	dummypurchaseorder.TotalCharge = float64(6500)
	dummypurchaseorder.DocumentStatus = "new"
	dummypurchaseorder.IsDeleted = 0
	dummypurchaseorder.Save()
	dummypurchaseorder2 := model.DummyPurchaseOrder()
	dummypurchaseorder2.Supplier = dummypartnershipsupplier
	dummypurchaseorder2.TotalCharge = float64(6500)
	dummypurchaseorder2.DocumentStatus = "active"
	dummypurchaseorder2.IsDeleted = 0
	dummypurchaseorder2.Save()

	e := CalculationTotalExpenditure(dummypartnershipsupplier.ID)
	assert.NoError(t, e)
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_expenditure from partnership where id = ?", dummypartnershipsupplier.ID).QueryRow(&total)
	assert.Equal(t, float64(13000), total)
}
func TestCalculationTotalExpenditureTwoDataWithCancelled(t *testing.T) {
	dummypartnershipsupplier := model.DummyPartnership()
	dummypartnershipsupplier.PartnershipType = "supplier"
	dummypartnershipsupplier.Save()
	dummypurchaseorder := model.DummyPurchaseOrder()
	dummypurchaseorder.Supplier = dummypartnershipsupplier
	dummypurchaseorder.TotalCharge = float64(6500)
	dummypurchaseorder.DocumentStatus = "cancelled"
	dummypurchaseorder.IsDeleted = 0
	dummypurchaseorder.Save()
	dummypurchaseorder2 := model.DummyPurchaseOrder()
	dummypurchaseorder2.Supplier = dummypartnershipsupplier
	dummypurchaseorder2.TotalCharge = float64(6500)
	dummypurchaseorder2.DocumentStatus = "cancelled"
	dummypurchaseorder2.IsDeleted = 0
	dummypurchaseorder2.Save()

	e := CalculationTotalExpenditure(dummypartnershipsupplier.ID)
	assert.NoError(t, e)
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_expenditure from partnership where id = ?", dummypartnershipsupplier.ID).QueryRow(&total)
	assert.Equal(t, float64(0), total)
}
func TestCalculationTotalExpenditureWithCancelled(t *testing.T) {
	dummypartnershipsupplier := model.DummyPartnership()
	dummypartnershipsupplier.PartnershipType = "supplier"
	dummypartnershipsupplier.Save()
	dummypurchaseorder := model.DummyPurchaseOrder()
	dummypurchaseorder.Supplier = dummypartnershipsupplier
	dummypurchaseorder.TotalCharge = float64(6500)
	dummypurchaseorder.DocumentStatus = "cancelled"
	dummypurchaseorder.IsDeleted = 0
	dummypurchaseorder.Save()

	e := CalculationTotalExpenditure(dummypartnershipsupplier.ID)
	assert.NoError(t, e)
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_expenditure from partnership where id = ?", dummypartnershipsupplier.ID).QueryRow(&total)
	assert.Equal(t, float64(0), total)
}
func TestCalculationTotalExpenditureWithIsDeleted(t *testing.T) {
	dummypartnershipsupplier := model.DummyPartnership()
	dummypartnershipsupplier.PartnershipType = "supplier"
	dummypartnershipsupplier.Save()
	dummypurchaseorder := model.DummyPurchaseOrder()
	dummypurchaseorder.Supplier = dummypartnershipsupplier
	dummypurchaseorder.TotalCharge = float64(6500)
	dummypurchaseorder.DocumentStatus = "new"
	dummypurchaseorder.IsDeleted = 1
	dummypurchaseorder.Save()

	e := CalculationTotalExpenditure(dummypartnershipsupplier.ID)
	assert.NoError(t, e)
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_expenditure from partnership where id = ?", dummypartnershipsupplier.ID).QueryRow(&total)
	assert.Equal(t, float64(0), total)
}
func TestCalculationTotalExpenditureWithNewOrFinishAndCancelled(t *testing.T) {
	dummypartnershipsupplier := model.DummyPartnership()
	dummypartnershipsupplier.PartnershipType = "supplier"
	dummypartnershipsupplier.Save()
	dummypurchaseorder := model.DummyPurchaseOrder()
	dummypurchaseorder.Supplier = dummypartnershipsupplier
	dummypurchaseorder.TotalCharge = float64(200000)
	dummypurchaseorder.DocumentStatus = "new"
	dummypurchaseorder.IsDeleted = 0
	dummypurchaseorder.Save()
	dummypurchaseorder2 := model.DummyPurchaseOrder()
	dummypurchaseorder2.Supplier = dummypartnershipsupplier
	dummypurchaseorder2.TotalCharge = float64(100000)
	dummypurchaseorder2.DocumentStatus = "finished"
	dummypurchaseorder2.IsDeleted = 0
	dummypurchaseorder2.Save()
	dummypurchaseorder3 := model.DummyPurchaseOrder()
	dummypurchaseorder3.Supplier = dummypartnershipsupplier
	dummypurchaseorder3.TotalCharge = float64(6500)
	dummypurchaseorder3.DocumentStatus = "cancelled"
	dummypurchaseorder3.IsDeleted = 1
	dummypurchaseorder3.Save()

	e := CalculationTotalExpenditure(dummypartnershipsupplier.ID)
	assert.NoError(t, e)
	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_expenditure from partnership where id = ?", dummypartnershipsupplier.ID).QueryRow(&total)
	assert.Equal(t, float64(300000), total)
}
func TestCalculationTotalCreditOneData(t *testing.T) {
	dummypartnershipsupplier := model.DummyPartnership()
	dummypartnershipsupplier.PartnershipType = "supplier"
	dummypartnershipsupplier.Save()
	dummypurchaseorder := model.DummyPurchaseOrder()
	dummypurchaseorder.Supplier = dummypartnershipsupplier
	dummypurchaseorder.TotalCharge = float64(5000)
	dummypurchaseorder.DocumentStatus = "new"
	dummypurchaseorder.IsDeleted = 0
	dummypurchaseorder.Save()

	err := CalculationTotalCredit(dummypartnershipsupplier.ID)
	assert.NoError(t, err)

	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_credit from partnership where id = ?", dummypartnershipsupplier.ID).QueryRow(&total)

	assert.Equal(t, float64(5000), total)
}
func TestCalculationTotalCreditTwoData(t *testing.T) {
	dummypartnershipsupplier := model.DummyPartnership()
	dummypartnershipsupplier.PartnershipType = "supplier"
	dummypartnershipsupplier.Save()
	dummypurchaseorder := model.DummyPurchaseOrder()
	dummypurchaseorder.Supplier = dummypartnershipsupplier
	dummypurchaseorder.TotalCharge = float64(5000)
	dummypurchaseorder.DocumentStatus = "new"
	dummypurchaseorder.IsDeleted = 0
	dummypurchaseorder.Save()
	dummypurchaseorder2 := model.DummyPurchaseOrder()
	dummypurchaseorder2.Supplier = dummypartnershipsupplier
	dummypurchaseorder2.TotalCharge = float64(10000)
	dummypurchaseorder2.DocumentStatus = "new"
	dummypurchaseorder2.IsDeleted = 0
	dummypurchaseorder2.Save()

	err := CalculationTotalCredit(dummypartnershipsupplier.ID)
	assert.NoError(t, err)

	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_credit from partnership where id = ?", dummypartnershipsupplier.ID).QueryRow(&total)

	assert.Equal(t, float64(15000), total)
}
func TestCalculationTotalCreditTwoDataCancelled(t *testing.T) {
	dummypartnershipsupplier := model.DummyPartnership()
	dummypartnershipsupplier.PartnershipType = "supplier"
	dummypartnershipsupplier.Save()
	dummypurchaseorder := model.DummyPurchaseOrder()
	dummypurchaseorder.Supplier = dummypartnershipsupplier
	dummypurchaseorder.TotalCharge = float64(5000)
	dummypurchaseorder.DocumentStatus = "cancelled"
	dummypurchaseorder.IsDeleted = 0
	dummypurchaseorder.Save()
	dummypurchaseorder2 := model.DummyPurchaseOrder()
	dummypurchaseorder2.Supplier = dummypartnershipsupplier
	dummypurchaseorder2.TotalCharge = float64(10000)
	dummypurchaseorder2.DocumentStatus = "cancelled"
	dummypurchaseorder2.IsDeleted = 0
	dummypurchaseorder2.Save()

	err := CalculationTotalCredit(dummypartnershipsupplier.ID)
	assert.NoError(t, err)

	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_credit from partnership where id = ?", dummypartnershipsupplier.ID).QueryRow(&total)

	assert.Equal(t, float64(0), total)
}
func TestCalculationTotalCreditWithCancelled(t *testing.T) {
	dummypartnershipsupplier := model.DummyPartnership()
	dummypartnershipsupplier.PartnershipType = "supplier"
	dummypartnershipsupplier.Save()
	dummypurchaseorder := model.DummyPurchaseOrder()
	dummypurchaseorder.Supplier = dummypartnershipsupplier
	dummypurchaseorder.TotalCharge = float64(5000)
	dummypurchaseorder.DocumentStatus = "cancelled"
	dummypurchaseorder.IsDeleted = 0
	dummypurchaseorder.Save()

	err := CalculationTotalCredit(dummypartnershipsupplier.ID)
	assert.NoError(t, err)

	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_credit from partnership where id = ?", dummypartnershipsupplier.ID).QueryRow(&total)

	assert.Equal(t, float64(0), total)
}
func TestCalculationTotalCreditWithIsDeleted(t *testing.T) {
	dummypartnershipsupplier := model.DummyPartnership()
	dummypartnershipsupplier.PartnershipType = "supplier"
	dummypartnershipsupplier.Save()
	dummypurchaseorder := model.DummyPurchaseOrder()
	dummypurchaseorder.Supplier = dummypartnershipsupplier
	dummypurchaseorder.TotalCharge = float64(5000)
	dummypurchaseorder.DocumentStatus = "new"
	dummypurchaseorder.IsDeleted = 1
	dummypurchaseorder.Save()

	err := CalculationTotalCredit(dummypartnershipsupplier.ID)
	assert.NoError(t, err)

	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_credit from partnership where id = ?", dummypartnershipsupplier.ID).QueryRow(&total)

	assert.Equal(t, float64(0), total)
}
func TestCalculationTotalCreditWithFinished(t *testing.T) {
	dummypartnershipsupplier := model.DummyPartnership()
	dummypartnershipsupplier.PartnershipType = "supplier"
	dummypartnershipsupplier.Save()
	dummypurchaseorder := model.DummyPurchaseOrder()
	dummypurchaseorder.Supplier = dummypartnershipsupplier
	dummypurchaseorder.TotalCharge = float64(2000)
	dummypurchaseorder.DocumentStatus = "finished"
	dummypurchaseorder.IsDeleted = 0
	dummypurchaseorder.Save()

	err := CalculationTotalCredit(dummypartnershipsupplier.ID)
	assert.NoError(t, err)

	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_credit from partnership where id = ?", dummypartnershipsupplier.ID).QueryRow(&total)

	assert.Equal(t, float64(2000), total)
}
func TestCalculationTotalCreditWithStatusNewFinisedAndCancelled(t *testing.T) {
	dummypartnershipsupplier := model.DummyPartnership()
	dummypartnershipsupplier.PartnershipType = "supplier"
	dummypartnershipsupplier.Save()
	dummypurchaseorder := model.DummyPurchaseOrder()
	dummypurchaseorder.Supplier = dummypartnershipsupplier
	dummypurchaseorder.TotalCharge = float64(5000)
	dummypurchaseorder.DocumentStatus = "new"
	dummypurchaseorder.IsDeleted = 0
	dummypurchaseorder.Save()
	dummypurchaseorder2 := model.DummyPurchaseOrder()
	dummypurchaseorder2.Supplier = dummypartnershipsupplier
	dummypurchaseorder2.TotalCharge = float64(5000)
	dummypurchaseorder2.DocumentStatus = "finished"
	dummypurchaseorder2.IsDeleted = 0
	dummypurchaseorder2.Save()
	dummypurchaseorder3 := model.DummyPurchaseOrder()
	dummypurchaseorder3.Supplier = dummypartnershipsupplier
	dummypurchaseorder3.TotalCharge = float64(5000)
	dummypurchaseorder3.DocumentStatus = "cancelled"
	dummypurchaseorder3.IsDeleted = 0
	dummypurchaseorder3.Save()

	err := CalculationTotalCredit(dummypartnershipsupplier.ID)
	assert.NoError(t, err)

	o := orm.NewOrm()
	var total float64
	_ = o.Raw("select total_credit from partnership where id = ?", dummypartnershipsupplier.ID).QueryRow(&total)

	assert.Equal(t, float64(10000), total)
}

func TestGetAllPartnershipsNoData(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM partnership").Exec()
	rq := orm.RequestQuery{}
	// test tidak ada data partnership
	m, total, e := GetAllPartnerships(&rq)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, m)
	assert.NoError(t, e)
}

func TestGetAllPartnershipsSuccess(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM partnership").Exec()
	// buat 3 dummy
	pr1 := model.DummyPartnership()
	pr1.IsArchived = int8(0)
	pr1.IsDeleted = int8(0)
	pr1.Save()
	pr2 := model.DummyPartnership()
	pr2.IsArchived = int8(1)
	pr2.IsDeleted = int8(0)
	pr2.Save()
	pr3 := model.DummyPartnership()
	pr3.IsArchived = int8(1)
	pr3.IsDeleted = int8(1)
	pr3.Save()
	rq := orm.RequestQuery{}
	// test ada 2 data partnership
	m, total, e := GetAllPartnerships(&rq)
	assert.NoError(t, e)
	assert.Equal(t, int64(2), total)
	assert.NotEmpty(t, m)
}

func TestGetPartnershipByFieldNoData(t *testing.T) {
	// clear data base partnership
	o := orm.NewOrm()
	o.Raw("DELETE FROM partnership").Exec()
	// ambil partnership dengan id 999999
	m, e := GetPartnershipByField("id", int64(999999))
	assert.Error(t, e)
	assert.Empty(t, m)
}

func TestGetPartnershipByFieldSuccess(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM partnership").Exec()
	fakePartner := model.DummyPartnership()
	// ambil data partnership berdasarkan code
	m, e := GetPartnershipByField("code", fakePartner.Code)
	assert.NoError(t, e)
	assert.NotEmpty(t, m)
	assert.Equal(t, fakePartner.Code, m.Code)
	assert.Equal(t, fakePartner.ID, m.ID)
}
