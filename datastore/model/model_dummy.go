// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"

	"git.qasico.com/cuxs/common/faker"
	"git.qasico.com/cuxs/common"
)

// DummyApplicationMenu make a dummy data for model ApplicationMenu
func DummyApplicationMenu() *ApplicationMenu {
	var m ApplicationMenu
	faker.Fill(&m, "ID")

	m.ApplicationModule = DummyApplicationModule()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyApplicationModule make a dummy data for model ApplicationModule
func DummyApplicationModule() *ApplicationModule {
	var m ApplicationModule
	faker.Fill(&m, "ID")

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyApplicationPrivilege make a dummy data for model ApplicationPrivilege
func DummyApplicationPrivilege() *ApplicationPrivilege {
	var m ApplicationPrivilege
	faker.Fill(&m, "ID")

	m.ApplicationModule = DummyApplicationModule()

	m.Usergroup = DummyUsergroup()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyApplicationSetting make a dummy data for model ApplicationSetting
func DummyApplicationSetting() *ApplicationSetting {
	var m ApplicationSetting
	faker.Fill(&m, "ID")

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyBankAccount make a dummy data for model BankAccount
func DummyBankAccount() *BankAccount {
	var m BankAccount
	faker.Fill(&m, "ID")
	m.BankNumber = common.RandomNumeric(10)
	m.BankName = common.RandomStr(10)

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyDirectPlacement make a dummy data for model DirectPlacement
func DummyDirectPlacement() *DirectPlacement {
	var m DirectPlacement
	faker.Fill(&m, "ID")

	m.CreatedBy = DummyUser()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyDirectPlacementItem make a dummy data for model DirectPlacementItem
func DummyDirectPlacementItem() *DirectPlacementItem {
	var m DirectPlacementItem
	faker.Fill(&m, "ID")

	m.ItemVariant = DummyItemVariant()

	m.DirectPlacment = DummyDirectPlacement()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyFinanceExpense make a dummy data for model FinanceExpense
func DummyFinanceExpense() *FinanceExpense {
	var m FinanceExpense
	faker.Fill(&m, "ID")

	m.CreatedBy = DummyUser()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyFinanceRevenue make a dummy data for model FinanceRevenue
func DummyFinanceRevenue() *FinanceRevenue {
	var m FinanceRevenue
	faker.Fill(&m, "ID")

	m.CreatedBy = DummyUser()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyInvoiceReceipt make a dummy data for model InvoiceReceipt
func DummyInvoiceReceipt() *InvoiceReceipt {
	var m InvoiceReceipt
	faker.Fill(&m, "ID")

	m.Partnership = DummyPartnership()

	m.CreatedBy = DummyUser()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyInvoiceReceiptItem make a dummy data for model InvoiceReceiptItem
func DummyInvoiceReceiptItem() *InvoiceReceiptItem {
	var m InvoiceReceiptItem
	faker.Fill(&m, "ID")

	m.InvoiceReceipt = DummyInvoiceReceipt()

	m.SalesInvoice = DummySalesInvoice()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyItem make a dummy data for model Item
func DummyItem() *Item {
	var m Item
	faker.Fill(&m, "ID")

	m.Category = DummyItemCategory()

	m.CreatedBy = DummyUser()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyItemCategory make a dummy data for model ItemCategory
func DummyItemCategory() *ItemCategory {
	var m ItemCategory
	faker.Fill(&m, "ID")

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyItemVariant make a dummy data for model ItemVariant
func DummyItemVariant() *ItemVariant {
	var m ItemVariant
	faker.Fill(&m, "ID")

	m.Item = DummyItem()

	m.Measurement = DummyMeasurement()

	m.CreatedBy = DummyUser()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyItemVariantPrice make a dummy data for model ItemVariantPrice
func DummyItemVariantPrice() *ItemVariantPrice {
	var m ItemVariantPrice
	faker.Fill(&m, "ID")

	m.ItemVariant = DummyItemVariant()

	m.PricingType = DummyPricingType()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyItemVariantStock make a dummy data for model ItemVariantStock
func DummyItemVariantStock() *ItemVariantStock {
	var m ItemVariantStock
	faker.Fill(&m, "ID")

	m.ItemVariant = DummyItemVariant()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyItemVariantStockLog make a dummy data for model ItemVariantStockLog
func DummyItemVariantStockLog() *ItemVariantStockLog {
	var m ItemVariantStockLog
	faker.Fill(&m, "ID")

	m.ItemVariantStock = DummyItemVariantStock()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyMeasurement make a dummy data for model Measurement
func DummyMeasurement() *Measurement {
	var m Measurement
	faker.Fill(&m, "ID")

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyPartnership make a dummy data for model Partnership
func DummyPartnership() *Partnership {
	var m Partnership
	faker.Fill(&m, "ID")

	m.CreatedBy = DummyUser()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyPricingType make a dummy data for model PricingType
func DummyPricingType() *PricingType {
	var m PricingType
	faker.Fill(&m, "ID")
	m.IsDefault = 0

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyPurchaseInvoice make a dummy data for model PurchaseInvoice
func DummyPurchaseInvoice() *PurchaseInvoice {
	var m PurchaseInvoice
	faker.Fill(&m, "ID")

	m.PurchaseOrder = DummyPurchaseOrder()

	m.CreatedBy = DummyUser()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyPurchaseOrder make a dummy data for model PurchaseOrder
func DummyPurchaseOrder() *PurchaseOrder {
	var m PurchaseOrder
	faker.Fill(&m, "ID")

	m.Supplier = DummyPartnership()

	m.CreatedBy = DummyUser()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyPurchaseOrderItem make a dummy data for model PurchaseOrderItem
func DummyPurchaseOrderItem() *PurchaseOrderItem {
	var m PurchaseOrderItem
	faker.Fill(&m, "ID")

	m.PurchaseOrder = DummyPurchaseOrder()

	m.ItemVariant = DummyItemVariant()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyPurchaseReturn make a dummy data for model PurchaseReturn
func DummyPurchaseReturn() *PurchaseReturn {
	var m PurchaseReturn
	faker.Fill(&m, "ID")

	m.PurchaseOrder = DummyPurchaseOrder()

	m.CreatedBy = DummyUser()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyPurchaseReturnItem make a dummy data for model PurchaseReturnItem
func DummyPurchaseReturnItem() *PurchaseReturnItem {
	var m PurchaseReturnItem
	faker.Fill(&m, "ID")

	m.PurchaseReturn = DummyPurchaseReturn()

	m.PurchaseOrderItem = DummyPurchaseOrderItem()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyRecapSales make a dummy data for model RecapSales
func DummyRecapSales() *RecapSales {
	var m RecapSales
	faker.Fill(&m, "ID")

	m.CreatedBy = DummyUser()

	m.Partnership = DummyPartnership()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyRecapSalesItem make a dummy data for model RecapSalesItem
func DummyRecapSalesItem() *RecapSalesItem {
	var m RecapSalesItem
	faker.Fill(&m, "ID")

	m.RecapSales = DummyRecapSales()

	m.SalesOrder = DummySalesOrder()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummySalesInvoice make a dummy data for model SalesInvoice
func DummySalesInvoice() *SalesInvoice {
	var m SalesInvoice
	faker.Fill(&m, "ID")

	m.CreatedBy = DummyUser()

	m.SalesOrder = DummySalesOrder()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummySalesOrder make a dummy data for model SalesOrder
func DummySalesOrder() *SalesOrder {
	var m SalesOrder
	faker.Fill(&m, "ID")

	m.CreatedBy = DummyUser()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummySalesOrderItem make a dummy data for model SalesOrderItem
func DummySalesOrderItem() *SalesOrderItem {
	var m SalesOrderItem
	faker.Fill(&m, "ID")

	m.SalesOrder = DummySalesOrder()

	m.ItemVariant = DummyItemVariant()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummySalesReturn make a dummy data for model SalesReturn
func DummySalesReturn() *SalesReturn {
	var m SalesReturn
	faker.Fill(&m, "ID")

	m.SalesOrder = DummySalesOrder()

	m.CreatedBy = DummyUser()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummySalesReturnItem make a dummy data for model SalesReturnItem
func DummySalesReturnItem() *SalesReturnItem {
	var m SalesReturnItem
	faker.Fill(&m, "ID")

	m.SalesOrderItem = DummySalesOrderItem()

	m.SalesReturn = DummySalesReturn()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyStockopname make a dummy data for model Stockopname
func DummyStockopname() *Stockopname {
	var m Stockopname
	faker.Fill(&m, "ID")

	m.CreatedBy = DummyUser()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyStockopnameItem make a dummy data for model StockopnameItem
func DummyStockopnameItem() *StockopnameItem {
	var m StockopnameItem
	faker.Fill(&m, "ID")

	m.Stockopname = DummyStockopname()

	m.ItemVariantStock = DummyItemVariantStock()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyUser make a dummy data for model User
func DummyUser() *User {
	var m User
	faker.Fill(&m, "ID")

	m.Usergroup = DummyUsergroup()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyUsergroup make a dummy data for model Usergroup
func DummyUsergroup() *Usergroup {
	var m Usergroup
	faker.Fill(&m, "ID")

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyWorkorderFulfillment make a dummy data for model WorkorderFulfillment
func DummyWorkorderFulfillment() *WorkorderFulfillment {
	var m WorkorderFulfillment
	faker.Fill(&m, "ID")

	m.SalesOrder = DummySalesOrder()

	m.CreatedBy = DummyUser()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyWorkorderFulfillmentItem make a dummy data for model WorkorderFulfillmentItem
func DummyWorkorderFulfillmentItem() *WorkorderFulfillmentItem {
	var m WorkorderFulfillmentItem
	faker.Fill(&m, "ID")

	m.WorkorderFulfillment = DummyWorkorderFulfillment()

	m.SalesOrderItem = DummySalesOrderItem()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyWorkorderReceiving make a dummy data for model WorkorderReceiving
func DummyWorkorderReceiving() *WorkorderReceiving {
	var m WorkorderReceiving
	faker.Fill(&m, "ID")

	m.CreatedBy = DummyUser()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyWorkorderReceivingItem make a dummy data for model WorkorderReceivingItem
func DummyWorkorderReceivingItem() *WorkorderReceivingItem {
	var m WorkorderReceivingItem
	faker.Fill(&m, "ID")

	m.WorkorderReceiving = DummyWorkorderReceiving()

	m.PurchaseOrderItem = DummyPurchaseOrderItem()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyWorkorderShipment make a dummy data for model WorkorderShipment
func DummyWorkorderShipment() *WorkorderShipment {
	var m WorkorderShipment
	faker.Fill(&m, "ID")

	m.CreatedBy = DummyUser()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyWorkorderShipmentItem make a dummy data for model WorkorderShipmentItem
func DummyWorkorderShipmentItem() *WorkorderShipmentItem {
	var m WorkorderShipmentItem
	faker.Fill(&m, "ID")

	m.WorkorderFulfillment = DummyWorkorderFulfillment()

	m.WorkorderShipment = DummyWorkorderShipment()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyUserPriviledgeWithUsergroup utk membuat dummy user berdasarkan usergroup yang diberikan
func DummyUserPriviledgeWithUsergroup(usergroupID int64) *User {
	var m User
	faker.Fill(&m, "ID")

	m.Usergroup = &Usergroup{ID: usergroupID}
	m.IsActive = 1

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}

// DummyInvoiceReceiptReturn make a dummy data for model InvoiceReceiptReturn
func DummyInvoiceReceiptReturn() *InvoiceReceiptReturn {
	var m InvoiceReceiptReturn
	faker.Fill(&m, "ID")

	m.InvoiceReceipt = DummyInvoiceReceipt()

	m.SalesReturn = DummySalesReturn()

	if e := m.Save(); e != nil {
		fmt.Printf("error saving %s", e.Error())
	}
	return &m
}
