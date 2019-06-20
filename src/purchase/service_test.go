// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchase

import (
	"fmt"
	"testing"
	"time"

	"git.qasico.com/mj/api/datastore/model"

	"git.qasico.com/cuxs/common"
	"git.qasico.com/cuxs/orm"
	"github.com/stretchr/testify/assert"
)

// TestCalculateTotalPaidPO Core Library
func TestCalculateTotalPaidPO(t *testing.T) {
	// Test untuk yang purchase order
	// total_paid di purchase_order adalah sum total_paid dari purchase invoice dari purchase order yang sama
	DummyPO := model.DummyPurchaseOrder()
	DummyPO.TotalPaid = 0
	DummyPO.IsDeleted = 0
	DummyPO.Save()

	DummyPI1 := model.DummyPurchaseInvoice()
	DummyPI1.PurchaseOrder = DummyPO
	DummyPI1.TotalPaid = 2000
	DummyPI1.IsDeleted = 0
	DummyPI1.DocumentStatus = "active"
	DummyPI1.Save()

	DummyPI2 := model.DummyPurchaseInvoice()
	DummyPI2.PurchaseOrder = DummyPO
	DummyPI2.TotalPaid = 3000
	DummyPI2.IsDeleted = 0
	DummyPI2.DocumentStatus = "active"
	DummyPI2.Save()

	DummyPI3 := model.DummyPurchaseInvoice()
	DummyPI3.PurchaseOrder = DummyPO
	DummyPI3.TotalPaid = 5000
	DummyPI3.IsDeleted = 0
	DummyPI3.DocumentStatus = "active"
	DummyPI3.Save()

	// Case 1 Semua document status active
	err := calculateTotalPaidPO(DummyPI1)
	assert.NoError(t, err, "Tidak boleh ada pesan error")
	DummyPO.Read()
	assert.Exactly(t, float64(10000), DummyPO.TotalPaid, "Hasil total paid pada PurchaseOrder harus sama dengan sum yang ada di Purchase Invoice.TotalPaid")

	DummyPI3.DocumentStatus = "finished"
	DummyPI3.Save()

	// Case 2 dua document status active
	err = calculateTotalPaidPO(DummyPI1)
	assert.NoError(t, err, "Tidak boleh ada pesan error")
	DummyPO.Read()
	assert.Exactly(t, float64(10000), DummyPO.TotalPaid, "Hasil total paid pada PurchaseOrder harus sama dengan sum yang ada di Purchase Invoice.TotalPaid")

	DummyPI2.IsDeleted = 1
	DummyPI2.Save()

	// Case 3 dua document status active satu is_delete 1
	err = calculateTotalPaidPO(DummyPI1)
	assert.NoError(t, err, "Tidak boleh ada pesan error")
	DummyPO.Read()
	assert.Exactly(t, float64(10000), DummyPO.TotalPaid, "Hasil total paid pada PurchaseOrder harus sama dengan sum yang ada di Purchase Invoice.TotalPaid")

}

// TestCalculateTotalPaidPI Core Library
func TestCalculateTotalPaidPI(t *testing.T) {
	DummyPO := model.DummyPurchaseOrder()
	DummyPO.TotalPaid = 0
	DummyPO.IsDeleted = 0
	DummyPO.Save()

	DummyPI1 := model.DummyPurchaseInvoice()
	DummyPI1.PurchaseOrder = DummyPO
	DummyPI1.TotalPaid = 0
	DummyPI1.IsDeleted = 0
	DummyPI1.DocumentStatus = "active"
	DummyPI1.Save()

	DummyPI2 := model.DummyPurchaseInvoice()
	DummyPI2.PurchaseOrder = DummyPO
	DummyPI2.TotalPaid = 0
	DummyPI2.IsDeleted = 0
	DummyPI2.DocumentStatus = "active"
	DummyPI2.Save()

	DummyFE11 := model.DummyFinanceExpense()
	DummyFE11.RefID = uint64(DummyPI1.ID)
	DummyFE11.Amount = 2000
	DummyFE11.RefType = "purchase_invoice"
	DummyFE11.DocumentStatus = "cleared"
	DummyFE11.Save()

	DummyFE12 := model.DummyFinanceExpense()
	DummyFE12.RefID = uint64(DummyPI1.ID)
	DummyFE12.Amount = 3000
	DummyFE12.RefType = "purchase_invoice"
	DummyFE12.DocumentStatus = "cleared"
	DummyFE12.Save()

	// Case 1 tidak ada financial_expense.document_status = "uncleared"
	_, err := CalculateTotalPaidPI(DummyPI1)
	assert.NoError(t, err, "Tidak boleh ada pesan error...")
	DummyPI1.Read()
	assert.Exactly(t, float64(5000), DummyPI1.TotalPaid)

	// Case 2 terdapat satu financial_expense.document_status = "uncleared"
	DummyFE12.DocumentStatus = "uncleared"
	DummyFE12.Save()
	_, err = CalculateTotalPaidPI(DummyPI1)
	assert.NoError(t, err, "Tidak boleh ada pesan error...")
	DummyPI1.Read()
	assert.Exactly(t, float64(2000), DummyPI1.TotalPaid)

	DummyFE21 := model.DummyFinanceExpense()
	DummyFE21.RefID = uint64(DummyPI2.ID)
	DummyFE21.Amount = 12000
	DummyFE21.RefType = "purchase_invoice"
	DummyFE21.DocumentStatus = "cleared"
	DummyFE21.Save()

	DummyFE22 := model.DummyFinanceExpense()
	DummyFE22.RefID = uint64(DummyPI2.ID)
	DummyFE22.Amount = 13000
	DummyFE22.RefType = "purchase_invoice"
	DummyFE22.DocumentStatus = "cleared"
	DummyFE22.Save()

	// Case 1 tidak ada financial_expense.document_status = "uncleared"
	_, err = CalculateTotalPaidPI(DummyPI2)
	assert.NoError(t, err, "Tidak boleh ada pesan error...")
	DummyPI2.Read()
	assert.Exactly(t, float64(25000), DummyPI2.TotalPaid)

	// Case 2 terdapat satu financial_expense.document_status = "uncleared"
	DummyFE22.DocumentStatus = "uncleared"
	DummyFE22.Save()
	_, err = CalculateTotalPaidPI(DummyPI2)
	assert.NoError(t, err, "Tidak boleh ada pesan error...")
	DummyPI2.Read()
	assert.Exactly(t, float64(12000), DummyPI2.TotalPaid)

	// Test Total Paid Purchase Order dengan 2 Uncleared
	err = calculateTotalPaidPO(DummyPI1)
	assert.NoError(t, err, "Tidak boleh ada pesan error...")
	DummyPO.Read()
	assert.Equal(t, float64(14000), DummyPO.TotalPaid)

	DummyPI1.DocumentStatus = "new"
	DummyPI1.IsDeleted = 1
	DummyPI1.Save()

	// Test Total Paid Purchase Order dengan 2 Uncleared, satu PurchaseInvoice.is_delete 1 satu PurchaseInvoice.DocumentStatus = "new"
	err = calculateTotalPaidPO(DummyPI1)
	assert.NoError(t, err, "Tidak boleh ada pesan error...")
	DummyPO.Read()
	assert.Equal(t, float64(14000), DummyPO.TotalPaid)
}

// TestChangeDocumentStatus testing fungsi core library ubah status dokumen
func TestChangeDocumentStatusUbahInvoiceErrorMsg(t *testing.T) {
	// Case 1 Ubah Invoice Status di PO

	DummyPO := model.DummyPurchaseOrder()
	DummyPO.DocumentStatus = "active"
	DummyPO.InvoiceStatus = "active"
	DummyPO.ReceivingStatus = "active"
	DummyPO.IsDeleted = 0
	DummyPO.TotalCharge = 100000
	DummyPO.TotalPaid = 0
	DummyPO.Save()

	DummyPI := model.DummyPurchaseInvoice()
	DummyPI.PurchaseOrder = DummyPO
	DummyPI.TotalAmount = float64(100000)
	DummyPI.TotalPaid = 0
	DummyPI.IsDeleted = 0
	DummyPI.DocumentStatus = "active"
	DummyPI.Save()

	DummyFE1 := model.DummyFinanceExpense()
	DummyFE1.RefID = uint64(DummyPI.ID)
	DummyFE1.Amount = 100000
	DummyFE1.RefType = "purchase_invoice"
	DummyFE1.DocumentStatus = "cleared"
	DummyFE1.Save()

	err := ChangeDocumentStatus(DummyPO)
	DummyPI.Read()
	DummyPI.PurchaseOrder.Read()
	assert.NoError(t, err, "Tidak boleh ada pesan error")
}

func TestChangeDocumentStatusUbahInvoiceTotalPaidPI(t *testing.T) {
	// Case 1 Ubah Invoice Status di PO
	DummyPO := model.DummyPurchaseOrder()
	DummyPO.DocumentStatus = "active"
	DummyPO.InvoiceStatus = "active"
	DummyPO.ReceivingStatus = "active"
	DummyPO.IsDeleted = 0
	DummyPO.TotalCharge = 100000
	DummyPO.TotalPaid = 0
	DummyPO.Save()

	DummyPI := model.DummyPurchaseInvoice()
	DummyPI.PurchaseOrder = DummyPO
	DummyPI.TotalAmount = float64(100000)
	DummyPI.TotalPaid = 0
	DummyPI.IsDeleted = 0
	DummyPI.DocumentStatus = "active"
	DummyPI.Save()

	DummyFE1 := model.DummyFinanceExpense()
	DummyFE1.RefID = uint64(DummyPI.ID)
	DummyFE1.Amount = 100000
	DummyFE1.RefType = "purchase_invoice"
	DummyFE1.DocumentStatus = "cleared"
	DummyFE1.Save()

	ChangeDocumentStatus(DummyPO)
	DummyPI.Read("ID")
	DummyPI.PurchaseOrder.Read()
	assert.Exactly(t, float64(100000), DummyPI.TotalPaid)
}

func TestChangeDocumentStatusUbahInvoiceTotalPaidPO(t *testing.T) {
	// Case 1 Ubah Invoice Status di PO

	DummyPO := model.DummyPurchaseOrder()
	DummyPO.DocumentStatus = "active"
	DummyPO.InvoiceStatus = "active"
	DummyPO.ReceivingStatus = "active"
	DummyPO.IsDeleted = 0
	DummyPO.TotalCharge = 100000
	DummyPO.TotalPaid = 0
	DummyPO.Save()

	DummyPI := model.DummyPurchaseInvoice()
	DummyPI.PurchaseOrder = DummyPO
	DummyPI.TotalAmount = float64(100000)
	DummyPI.TotalPaid = 0
	DummyPI.IsDeleted = 0
	DummyPI.DocumentStatus = "active"
	DummyPI.Save()

	DummyFE1 := model.DummyFinanceExpense()
	DummyFE1.RefID = uint64(DummyPI.ID)
	DummyFE1.Amount = 100000
	DummyFE1.RefType = "purchase_invoice"
	DummyFE1.DocumentStatus = "cleared"
	DummyFE1.Save()

	ChangeDocumentStatus(DummyPO)
	DummyPI.Read()
	DummyPI.PurchaseOrder.Read()
	assert.Exactly(t, float64(100000), DummyPI.PurchaseOrder.TotalPaid)
}

func TestChangeDocumentStatusUbahInvoicePOInvoiceStatus(t *testing.T) {
	// Case 1 Ubah Invoice Status di PO

	DummyPO := model.DummyPurchaseOrder()
	DummyPO.DocumentStatus = "active"
	DummyPO.InvoiceStatus = "active"
	DummyPO.ReceivingStatus = "active"
	DummyPO.IsDeleted = 0
	DummyPO.TotalCharge = 100000
	DummyPO.TotalPaid = 0
	DummyPO.Save()

	DummyPI := model.DummyPurchaseInvoice()
	DummyPI.PurchaseOrder = DummyPO
	DummyPI.TotalAmount = float64(100000)
	DummyPI.TotalPaid = 0
	DummyPI.IsDeleted = 0
	DummyPI.DocumentStatus = "active"
	DummyPI.Save()

	DummyFE1 := model.DummyFinanceExpense()
	DummyFE1.RefID = uint64(DummyPI.ID)
	DummyFE1.Amount = 100000
	DummyFE1.RefType = "purchase_invoice"
	DummyFE1.DocumentStatus = "cleared"
	DummyFE1.Save()

	ChangeDocumentStatus(DummyPO)
	DummyPI.Read()
	DummyPI.PurchaseOrder.Read()
	assert.Exactly(t, "finished", DummyPI.PurchaseOrder.InvoiceStatus)
}

func TestChangeDocumentStatusUbahInvoicePIDocumentStatus(t *testing.T) {
	// Case 1 Ubah Invoice Status di PO

	DummyPO := model.DummyPurchaseOrder()
	DummyPO.DocumentStatus = "active"
	DummyPO.InvoiceStatus = "active"
	DummyPO.ReceivingStatus = "active"
	DummyPO.IsDeleted = 0
	DummyPO.TotalCharge = 100000
	DummyPO.TotalPaid = 0
	DummyPO.Save()

	DummyPI := model.DummyPurchaseInvoice()
	DummyPI.PurchaseOrder = DummyPO
	DummyPI.TotalAmount = float64(100000)
	DummyPI.TotalPaid = 0
	DummyPI.IsDeleted = 0
	DummyPI.DocumentStatus = "active"
	DummyPI.Save()

	DummyFE1 := model.DummyFinanceExpense()
	DummyFE1.RefID = uint64(DummyPI.ID)
	DummyFE1.Amount = 100000
	DummyFE1.RefType = "purchase_invoice"
	DummyFE1.DocumentStatus = "cleared"
	DummyFE1.Save()

	ChangeDocumentStatus(DummyPO)
	DummyPI.Read()
	DummyPI.PurchaseOrder.Read()
	assert.Exactly(t, "finished", DummyPI.DocumentStatus)
}

func TestChangeDocumentStatusUbahInvoicePIDocumentStatusWhenIsDelete(t *testing.T) {
	// Case 1 Ubah Invoice Status di PO

	DummyPO := model.DummyPurchaseOrder()
	DummyPO.DocumentStatus = "active"
	DummyPO.InvoiceStatus = "active"
	DummyPO.ReceivingStatus = "active"
	DummyPO.IsDeleted = 0
	DummyPO.TotalCharge = 100000
	DummyPO.TotalPaid = 0
	DummyPO.Save()

	DummyPI := model.DummyPurchaseInvoice()
	DummyPI.PurchaseOrder = DummyPO
	DummyPI.TotalAmount = float64(100000)
	DummyPI.TotalPaid = 0
	DummyPI.IsDeleted = 1
	DummyPI.DocumentStatus = "active"
	DummyPI.Save()

	DummyFE1 := model.DummyFinanceExpense()
	DummyFE1.RefID = uint64(DummyPI.ID)
	DummyFE1.Amount = 100000
	DummyFE1.RefType = "purchase_invoice"
	DummyFE1.DocumentStatus = "cleared"
	DummyFE1.Save()

	ChangeDocumentStatus(DummyPO)
	DummyPI.Read()
	DummyPI.PurchaseOrder.Read()
	assert.NotEqual(t, "finished", DummyPI.PurchaseOrder.DocumentStatus)
}

func TestChangeDocumentStatusWithPOCancelled(t *testing.T) {
	// Case 1 Ubah Invoice Status di PO

	DummyPO := model.DummyPurchaseOrder()
	DummyPO.InvoiceStatus = "active"
	DummyPO.ReceivingStatus = "active"
	DummyPO.IsDeleted = 0
	DummyPO.TotalCharge = 100000
	DummyPO.TotalPaid = 0
	DummyPO.DocumentStatus = "cancelled"
	DummyPO.Save()

	DummyPI := model.DummyPurchaseInvoice()
	DummyPI.PurchaseOrder = DummyPO
	DummyPI.TotalAmount = float64(100000)
	DummyPI.TotalPaid = 0
	DummyPI.IsDeleted = 0
	DummyPI.DocumentStatus = "active"
	DummyPI.Save()

	DummyFE1 := model.DummyFinanceExpense()
	DummyFE1.RefID = uint64(DummyPI.ID)
	DummyFE1.Amount = 100000
	DummyFE1.RefType = "purchase_invoice"
	DummyFE1.DocumentStatus = "cleared"
	DummyFE1.Save()

	ChangeDocumentStatus(DummyPO)
	DummyPI.Read()
	DummyPI.PurchaseOrder.Read()
	assert.NotEqual(t, "finished", DummyPI.PurchaseOrder.DocumentStatus, "actual : "+DummyPI.PurchaseOrder.DocumentStatus)
}

func TestChangeDocumentStatusUbahReceivingStatus1(t *testing.T) {
	// Case 1 Receiving Status semua nya OK

	DummyPO := model.DummyPurchaseOrder()
	DummyPO.DocumentStatus = "active"
	DummyPO.InvoiceStatus = "active"
	DummyPO.ReceivingStatus = "active"
	DummyPO.IsDeleted = 0
	DummyPO.TotalCharge = 100000
	DummyPO.TotalPaid = 0
	DummyPO.Save()

	DummyPI := model.DummyPurchaseInvoice()
	DummyPI.PurchaseOrder = DummyPO
	DummyPI.TotalPaid = 0
	DummyPI.IsDeleted = 0
	DummyPI.DocumentStatus = "active"
	DummyPI.Save()

	DummyFE1 := model.DummyFinanceExpense()
	DummyFE1.RefID = uint64(DummyPI.ID)
	DummyFE1.Amount = 100000
	DummyFE1.RefType = "purchase_invoice"
	DummyFE1.DocumentStatus = "cleared"
	DummyFE1.Save()

	ChangeDocumentStatus(DummyPO)
	DummyPI.Read()
	DummyPI.PurchaseOrder.Read()

	di1 := model.DummyItem()
	div1 := model.DummyItemVariant()
	div1.Item = di1
	div1.Save()

	di2 := model.DummyItem()
	div2 := model.DummyItemVariant()
	div2.Item = di2
	div2.Save()

	DummyPOI1 := model.DummyPurchaseOrderItem()
	DummyPOI1.PurchaseOrder = DummyPO
	DummyPOI1.ItemVariant = div1
	DummyPOI1.Quantity = 10
	DummyPOI1.Save()

	DummyPOI2 := model.DummyPurchaseOrderItem()
	DummyPOI2.PurchaseOrder = DummyPO
	DummyPOI2.ItemVariant = div2
	DummyPOI2.Quantity = 5
	DummyPOI2.Save()

	DummyWR := model.DummyWorkorderReceiving()
	DummyWR.PurchaseOrder = DummyPO
	DummyWR.DocumentStatus = "active"
	DummyWR.IsDeleted = 0
	DummyWR.Save()

	DummyWRI11 := model.DummyWorkorderReceivingItem()
	DummyWRI11.WorkorderReceiving = DummyWR
	DummyWRI11.PurchaseOrderItem = DummyPOI1
	DummyWRI11.Quantity = 6
	DummyWRI11.Save()

	DummyWRI12 := model.DummyWorkorderReceivingItem()
	DummyWRI12.WorkorderReceiving = DummyWR
	DummyWRI12.PurchaseOrderItem = DummyPOI1
	DummyWRI12.Quantity = 4
	DummyWRI12.Save()

	DummyWRI21 := model.DummyWorkorderReceivingItem()
	DummyWRI21.WorkorderReceiving = DummyWR
	DummyWRI21.PurchaseOrderItem = DummyPOI2
	DummyWRI21.Quantity = 5
	DummyWRI21.Save()

	ChangeDocumentStatus(DummyPO)
	DummyPO.Read()
	assert.Exactly(t, "finished", DummyPO.ReceivingStatus)
}

func TestChangeDocumentStatusUbahReceivingStatus2(t *testing.T) {
	DummyPO := model.DummyPurchaseOrder()
	DummyPO.DocumentStatus = "active"
	DummyPO.InvoiceStatus = "active"
	DummyPO.ReceivingStatus = "active"
	DummyPO.IsDeleted = 0
	DummyPO.TotalCharge = 100000
	DummyPO.TotalPaid = 0
	DummyPO.Save()

	DummyPI := model.DummyPurchaseInvoice()
	DummyPI.PurchaseOrder = DummyPO
	DummyPI.TotalPaid = 0
	DummyPI.IsDeleted = 0
	DummyPI.DocumentStatus = "active"
	DummyPI.Save()

	DummyFE1 := model.DummyFinanceExpense()
	DummyFE1.RefID = uint64(DummyPI.ID)
	DummyFE1.Amount = 100000
	DummyFE1.RefType = "purchase_invoice"
	DummyFE1.DocumentStatus = "cleared"
	DummyFE1.Save()

	DummyPI.Read()
	DummyPI.PurchaseOrder.Read()
	di1 := model.DummyItem()
	div1 := model.DummyItemVariant()
	div1.Item = di1
	div1.Save()

	di2 := model.DummyItem()
	div2 := model.DummyItemVariant()
	div2.Item = di2
	div2.Save()

	DummyPOI1 := model.DummyPurchaseOrderItem()
	DummyPOI1.PurchaseOrder = DummyPO
	DummyPOI1.ItemVariant = div1
	DummyPOI1.Quantity = 10
	DummyPOI1.Save()

	DummyPOI2 := model.DummyPurchaseOrderItem()
	DummyPOI2.PurchaseOrder = DummyPO
	DummyPOI2.ItemVariant = div2
	DummyPOI2.Quantity = 5
	DummyPOI2.Save()

	DummyWR := model.DummyWorkorderReceiving()
	DummyWR.PurchaseOrder = DummyPO
	DummyWR.DocumentStatus = "active"
	DummyWR.IsDeleted = 0
	DummyWR.Save()

	DummyWRI11 := model.DummyWorkorderReceivingItem()
	DummyWRI11.WorkorderReceiving = DummyWR
	DummyWRI11.PurchaseOrderItem = DummyPOI1
	DummyWRI11.Quantity = 6
	DummyWRI11.Save()

	DummyWRI12 := model.DummyWorkorderReceivingItem()
	DummyWRI12.WorkorderReceiving = DummyWR
	DummyWRI12.PurchaseOrderItem = DummyPOI1
	DummyWRI12.Quantity = 4
	DummyWRI12.Save()

	DummyWRI21 := model.DummyWorkorderReceivingItem()
	DummyWRI21.WorkorderReceiving = DummyWR
	DummyWRI21.PurchaseOrderItem = DummyPOI2
	DummyWRI21.Quantity = 1
	DummyWRI21.Save()

	ChangeDocumentStatus(DummyPO)
	DummyPO.Read()
	assert.Exactly(t, "active", DummyPO.ReceivingStatus)
}

func TestChangeDocumentStatusUbahReceivingStatusDocumentStatusWithPOCancelled(t *testing.T) {
	// Case 1 Receiving Status semua nya OK

	DummyPO := model.DummyPurchaseOrder()
	DummyPO.DocumentStatus = "cancelled"
	DummyPO.InvoiceStatus = "active"
	DummyPO.ReceivingStatus = "active"
	DummyPO.IsDeleted = 0
	DummyPO.TotalCharge = 100000
	DummyPO.TotalPaid = 0
	DummyPO.Save()

	DummyPI := model.DummyPurchaseInvoice()
	DummyPI.PurchaseOrder = DummyPO
	DummyPI.TotalPaid = 0
	DummyPI.IsDeleted = 0
	DummyPI.DocumentStatus = "active"
	DummyPI.Save()

	DummyFE1 := model.DummyFinanceExpense()
	DummyFE1.RefID = uint64(DummyPI.ID)
	DummyFE1.Amount = 100000
	DummyFE1.RefType = "purchase_invoice"
	DummyFE1.DocumentStatus = "cleared"
	DummyFE1.Save()

	di1 := model.DummyItem()
	div1 := model.DummyItemVariant()
	div1.Item = di1
	div1.Save()

	di2 := model.DummyItem()
	div2 := model.DummyItemVariant()
	div2.Item = di2
	div2.Save()

	DummyPOI1 := model.DummyPurchaseOrderItem()
	DummyPOI1.PurchaseOrder = DummyPO
	DummyPOI1.ItemVariant = div1
	DummyPOI1.Quantity = 10
	DummyPOI1.Save()

	DummyPOI2 := model.DummyPurchaseOrderItem()
	DummyPOI2.PurchaseOrder = DummyPO
	DummyPOI2.ItemVariant = div2
	DummyPOI2.Quantity = 5
	DummyPOI2.Save()

	DummyWR := model.DummyWorkorderReceiving()
	DummyWR.PurchaseOrder = DummyPO
	DummyWR.DocumentStatus = "active"
	DummyWR.IsDeleted = 0
	DummyWR.Save()

	DummyWRI11 := model.DummyWorkorderReceivingItem()
	DummyWRI11.WorkorderReceiving = DummyWR
	DummyWRI11.PurchaseOrderItem = DummyPOI1
	DummyWRI11.Quantity = 6
	DummyWRI11.Save()

	DummyWRI12 := model.DummyWorkorderReceivingItem()
	DummyWRI12.WorkorderReceiving = DummyWR
	DummyWRI12.PurchaseOrderItem = DummyPOI1
	DummyWRI12.Quantity = 4
	DummyWRI12.Save()

	DummyWRI21 := model.DummyWorkorderReceivingItem()
	DummyWRI21.WorkorderReceiving = DummyWR
	DummyWRI21.PurchaseOrderItem = DummyPOI2
	DummyWRI21.Quantity = 5
	DummyWRI21.Save()

	ChangeDocumentStatus(DummyPO)
	DummyPO.Read()
	assert.NotEqual(t, "finished", DummyPO.ReceivingStatus, "actual: "+DummyPO.ReceivingStatus)
}

func TestChangeDocumentStatusUbahReceivingStatusDocumentStatusWithPODeleted(t *testing.T) {
	// Case 1 Receiving Status semua nya OK

	DummyPO := model.DummyPurchaseOrder()
	DummyPO.DocumentStatus = "active"
	DummyPO.InvoiceStatus = "active"
	DummyPO.ReceivingStatus = "active"
	DummyPO.IsDeleted = 1
	DummyPO.TotalCharge = 100000
	DummyPO.TotalPaid = 0
	DummyPO.Save()

	DummyPI := model.DummyPurchaseInvoice()
	DummyPI.PurchaseOrder = DummyPO
	DummyPI.TotalPaid = 0
	DummyPI.IsDeleted = 0
	DummyPI.DocumentStatus = "active"
	DummyPI.Save()

	DummyFE1 := model.DummyFinanceExpense()
	DummyFE1.RefID = uint64(DummyPI.ID)
	DummyFE1.Amount = 100000
	DummyFE1.RefType = "purchase_invoice"
	DummyFE1.DocumentStatus = "cleared"
	DummyFE1.Save()

	ChangeDocumentStatus(DummyPO)
	DummyPI.Read()
	DummyPI.PurchaseOrder.Read()

	di1 := model.DummyItem()
	div1 := model.DummyItemVariant()
	div1.Item = di1
	div1.Save()

	di2 := model.DummyItem()
	div2 := model.DummyItemVariant()
	div2.Item = di2
	div2.Save()

	DummyPOI1 := model.DummyPurchaseOrderItem()
	DummyPOI1.PurchaseOrder = DummyPO
	DummyPOI1.ItemVariant = div1
	DummyPOI1.Quantity = 10
	DummyPOI1.Save()

	DummyPOI2 := model.DummyPurchaseOrderItem()
	DummyPOI2.PurchaseOrder = DummyPO
	DummyPOI2.ItemVariant = div2
	DummyPOI2.Quantity = 5
	DummyPOI2.Save()

	DummyWR := model.DummyWorkorderReceiving()
	DummyWR.PurchaseOrder = DummyPO
	DummyWR.DocumentStatus = "active"
	DummyWR.IsDeleted = 0
	DummyWR.Save()

	DummyWRI11 := model.DummyWorkorderReceivingItem()
	DummyWRI11.WorkorderReceiving = DummyWR
	DummyWRI11.PurchaseOrderItem = DummyPOI1
	DummyWRI11.Quantity = 6
	DummyWRI11.Save()

	DummyWRI12 := model.DummyWorkorderReceivingItem()
	DummyWRI12.WorkorderReceiving = DummyWR
	DummyWRI12.PurchaseOrderItem = DummyPOI1
	DummyWRI12.Quantity = 4
	DummyWRI12.Save()

	DummyWRI21 := model.DummyWorkorderReceivingItem()
	DummyWRI21.WorkorderReceiving = DummyWR
	DummyWRI21.PurchaseOrderItem = DummyPOI2
	DummyWRI21.Quantity = 5
	DummyWRI21.Save()

	ChangeDocumentStatus(DummyPO)
	DummyPO.Read()
	assert.NotEqual(t, "finished", DummyPO.ReceivingStatus, "actual: "+DummyPO.ReceivingStatus)
}

func TestChangeDocumentStatus1(t *testing.T) {
	// Test Ubah Dokumen Status di Purchase
	//Case 1 Semua OK
	DummyPO := model.DummyPurchaseOrder()
	DummyPO.DocumentStatus = "active"
	DummyPO.InvoiceStatus = "active"
	DummyPO.ReceivingStatus = "active"
	DummyPO.IsDeleted = 0
	DummyPO.TotalCharge = 100000
	DummyPO.TotalPaid = 0
	DummyPO.Save()

	DummyPI := model.DummyPurchaseInvoice()
	DummyPI.PurchaseOrder = DummyPO
	DummyPI.TotalPaid = 0
	DummyPI.IsDeleted = 0
	DummyPI.DocumentStatus = "active"
	DummyPI.Save()

	DummyFE1 := model.DummyFinanceExpense()
	DummyFE1.RefID = uint64(DummyPI.ID)
	DummyFE1.Amount = 100000
	DummyFE1.RefType = "purchase_invoice"
	DummyFE1.DocumentStatus = "cleared"
	DummyFE1.Save()

	ChangeDocumentStatus(DummyPO)
	DummyPI.Read()
	DummyPI.PurchaseOrder.Read()

	di1 := model.DummyItem()
	div1 := model.DummyItemVariant()
	div1.Item = di1
	div1.Save()

	di2 := model.DummyItem()
	div2 := model.DummyItemVariant()
	div2.Item = di2
	div2.Save()

	DummyPOI1 := model.DummyPurchaseOrderItem()
	DummyPOI1.PurchaseOrder = DummyPO
	DummyPOI1.ItemVariant = div1
	DummyPOI1.Quantity = 10
	DummyPOI1.Save()

	DummyPOI2 := model.DummyPurchaseOrderItem()
	DummyPOI2.PurchaseOrder = DummyPO
	DummyPOI2.ItemVariant = div2
	DummyPOI2.Quantity = 5
	DummyPOI2.Save()

	DummyWR := model.DummyWorkorderReceiving()
	DummyWR.PurchaseOrder = DummyPO
	DummyWR.DocumentStatus = "active"
	DummyWR.IsDeleted = 0
	DummyWR.Save()

	DummyWRI11 := model.DummyWorkorderReceivingItem()
	DummyWRI11.WorkorderReceiving = DummyWR
	DummyWRI11.PurchaseOrderItem = DummyPOI1
	DummyWRI11.Quantity = 6
	DummyWRI11.Save()

	DummyWRI12 := model.DummyWorkorderReceivingItem()
	DummyWRI12.WorkorderReceiving = DummyWR
	DummyWRI12.PurchaseOrderItem = DummyPOI1
	DummyWRI12.Quantity = 4
	DummyWRI12.Save()

	DummyWRI21 := model.DummyWorkorderReceivingItem()
	DummyWRI21.WorkorderReceiving = DummyWR
	DummyWRI21.PurchaseOrderItem = DummyPOI2
	DummyWRI21.Quantity = 5
	DummyWRI21.Save()

	ChangeDocumentStatus(DummyPO)
	DummyPO.Read()
	assert.Exactly(t, "finished", DummyPO.DocumentStatus)
}

func TestChangeDocumentStatus2(t *testing.T) {
	// Case 2 Receiving Status document status masih active karena Quantity
	// Item yang masuk Workorder Receiving Item tidak sama dengan yang ada di Purchase Order Item

	DummyPO := model.DummyPurchaseOrder()
	DummyPO.DocumentStatus = "active"
	DummyPO.InvoiceStatus = "active"
	DummyPO.ReceivingStatus = "active"
	DummyPO.IsDeleted = 0
	DummyPO.TotalCharge = 100000
	DummyPO.TotalPaid = 0
	DummyPO.Save()

	DummyPI := model.DummyPurchaseInvoice()
	DummyPI.PurchaseOrder = DummyPO
	DummyPI.TotalPaid = 0
	DummyPI.IsDeleted = 0
	DummyPI.DocumentStatus = "active"
	DummyPI.Save()

	DummyFE1 := model.DummyFinanceExpense()
	DummyFE1.RefID = uint64(DummyPI.ID)
	DummyFE1.Amount = 100000
	DummyFE1.RefType = "purchase_invoice"
	DummyFE1.DocumentStatus = "cleared"
	DummyFE1.Save()
	DummyPI.Read()
	DummyPI.PurchaseOrder.Read()

	di1 := model.DummyItem()
	div1 := model.DummyItemVariant()
	div1.Item = di1
	div1.Save()

	di2 := model.DummyItem()
	div2 := model.DummyItemVariant()
	div2.Item = di2
	div2.Save()

	DummyPOI1 := model.DummyPurchaseOrderItem()
	DummyPOI1.PurchaseOrder = DummyPO
	DummyPOI1.ItemVariant = div1
	DummyPOI1.Quantity = 10
	DummyPOI1.Save()

	DummyPOI2 := model.DummyPurchaseOrderItem()
	DummyPOI2.PurchaseOrder = DummyPO
	DummyPOI2.ItemVariant = div2
	DummyPOI2.Quantity = 5
	DummyPOI2.Save()

	DummyWR := model.DummyWorkorderReceiving()
	DummyWR.PurchaseOrder = DummyPO
	DummyWR.DocumentStatus = "active"
	DummyWR.IsDeleted = 0
	DummyWR.Save()

	DummyWRI11 := model.DummyWorkorderReceivingItem()
	DummyWRI11.WorkorderReceiving = DummyWR
	DummyWRI11.PurchaseOrderItem = DummyPOI1
	DummyWRI11.Quantity = 6
	DummyWRI11.Save()

	DummyWRI12 := model.DummyWorkorderReceivingItem()
	DummyWRI12.WorkorderReceiving = DummyWR
	DummyWRI12.PurchaseOrderItem = DummyPOI1
	DummyWRI12.Quantity = 4
	DummyWRI12.Save()

	DummyWRI21 := model.DummyWorkorderReceivingItem()
	DummyWRI21.WorkorderReceiving = DummyWR
	DummyWRI21.PurchaseOrderItem = DummyPOI2
	DummyWRI21.Quantity = 1
	DummyWRI21.Save()

	ChangeDocumentStatus(DummyPO)
	DummyPO.Read()
	assert.Exactly(t, "active", DummyPO.DocumentStatus)
}

func TestChangeDocumentStatus3(t *testing.T) {
	// Case 3 Receiving Status document status masih active karena Quantity
	// Item yang masuk Workorder Receiving Item tidak sama dengan yang ada di Purchase Order Item

	DummyPO := model.DummyPurchaseOrder()
	DummyPO.DocumentStatus = "active"
	DummyPO.InvoiceStatus = "active"
	DummyPO.ReceivingStatus = "active"
	DummyPO.IsDeleted = 0
	DummyPO.TotalCharge = 100000
	DummyPO.TotalPaid = 0
	DummyPO.Save()

	DummyPI := model.DummyPurchaseInvoice()
	DummyPI.PurchaseOrder = DummyPO
	DummyPI.TotalPaid = 0
	DummyPI.IsDeleted = 0
	DummyPI.DocumentStatus = "active"
	DummyPI.Save()

	DummyFE1 := model.DummyFinanceExpense()
	DummyFE1.RefID = uint64(DummyPI.ID)
	DummyFE1.Amount = 100000
	DummyFE1.RefType = "purchase_invoice"
	DummyFE1.DocumentStatus = "cleared"
	DummyFE1.Save()
	DummyPI.Read()
	DummyPI.PurchaseOrder.Read()

	di1 := model.DummyItem()
	div1 := model.DummyItemVariant()
	div1.Item = di1
	div1.Save()

	di2 := model.DummyItem()
	div2 := model.DummyItemVariant()
	div2.Item = di2
	div2.Save()

	DummyPOI1 := model.DummyPurchaseOrderItem()
	DummyPOI1.PurchaseOrder = DummyPO
	DummyPOI1.ItemVariant = div1
	DummyPOI1.Quantity = 10
	DummyPOI1.Save()

	DummyPOI2 := model.DummyPurchaseOrderItem()
	DummyPOI2.PurchaseOrder = DummyPO
	DummyPOI2.ItemVariant = div2
	DummyPOI2.Quantity = 5
	DummyPOI2.Save()

	DummyWR := model.DummyWorkorderReceiving()
	DummyWR.PurchaseOrder = DummyPO
	DummyWR.DocumentStatus = "active"
	DummyWR.IsDeleted = 0
	DummyWR.Save()

	DummyWRI11 := model.DummyWorkorderReceivingItem()
	DummyWRI11.WorkorderReceiving = DummyWR
	DummyWRI11.PurchaseOrderItem = DummyPOI1
	DummyWRI11.Quantity = 6
	DummyWRI11.Save()

	DummyWRI12 := model.DummyWorkorderReceivingItem()
	DummyWRI12.WorkorderReceiving = DummyWR
	DummyWRI12.PurchaseOrderItem = DummyPOI1
	DummyWRI12.Quantity = 4
	DummyWRI12.Save()

	DummyWRI21 := model.DummyWorkorderReceivingItem()
	DummyWRI21.WorkorderReceiving = DummyWR
	DummyWRI21.PurchaseOrderItem = DummyPOI2
	DummyWRI21.Quantity = 5
	DummyWRI21.Save()
	DummyPO.Read()

	DummyWRI11.Delete()
	DummyWRI12.Delete()

	ChangeDocumentStatus(DummyPO)
	DummyPO.Read()
	assert.Exactly(t, "active", DummyPO.DocumentStatus)
}

func TestChangeDocumentStatus4(t *testing.T) {
	// Case 4 Receiving Status document status masih active karena Workorder receiving belum dibuat
	DummyPO := model.DummyPurchaseOrder()
	DummyPO.DocumentStatus = "active"
	DummyPO.InvoiceStatus = "active"
	DummyPO.ReceivingStatus = "active"
	DummyPO.IsDeleted = 0
	DummyPO.TotalCharge = 100000
	DummyPO.TotalPaid = 0
	DummyPO.Save()

	DummyPI := model.DummyPurchaseInvoice()
	DummyPI.PurchaseOrder = DummyPO
	DummyPI.TotalPaid = 0
	DummyPI.IsDeleted = 0
	DummyPI.DocumentStatus = "active"
	DummyPI.Save()

	DummyFE1 := model.DummyFinanceExpense()
	DummyFE1.RefID = uint64(DummyPI.ID)
	DummyFE1.Amount = 100000
	DummyFE1.RefType = "purchase_invoice"
	DummyFE1.DocumentStatus = "cleared"
	DummyFE1.Save()
	DummyPI.Read()
	DummyPI.PurchaseOrder.Read()

	di1 := model.DummyItem()
	div1 := model.DummyItemVariant()
	div1.Item = di1
	div1.Save()

	di2 := model.DummyItem()
	div2 := model.DummyItemVariant()
	div2.Item = di2
	div2.Save()

	DummyPOI1 := model.DummyPurchaseOrderItem()
	DummyPOI1.PurchaseOrder = DummyPO
	DummyPOI1.ItemVariant = div1
	DummyPOI1.Quantity = 10
	DummyPOI1.Save()

	DummyPOI2 := model.DummyPurchaseOrderItem()
	DummyPOI2.PurchaseOrder = DummyPO
	DummyPOI2.ItemVariant = div2
	DummyPOI2.Quantity = 5
	DummyPOI2.Save()

	DummyWR := model.DummyWorkorderReceiving()
	DummyWR.PurchaseOrder = DummyPO
	DummyWR.DocumentStatus = "active"
	DummyWR.IsDeleted = 0
	DummyWR.Save()

	DummyWRI11 := model.DummyWorkorderReceivingItem()
	DummyWRI11.WorkorderReceiving = DummyWR
	DummyWRI11.PurchaseOrderItem = DummyPOI1
	DummyWRI11.Quantity = 6
	DummyWRI11.Save()

	DummyWRI12 := model.DummyWorkorderReceivingItem()
	DummyWRI12.WorkorderReceiving = DummyWR
	DummyWRI12.PurchaseOrderItem = DummyPOI1
	DummyWRI12.Quantity = 4
	DummyWRI12.Save()

	DummyWRI21 := model.DummyWorkorderReceivingItem()
	DummyWRI21.WorkorderReceiving = DummyWR
	DummyWRI21.PurchaseOrderItem = DummyPOI2
	DummyWRI21.Quantity = 5
	DummyWRI21.Save()

	DummyWRI21.Delete()
	DummyWR.Delete()

	ChangeDocumentStatus(DummyPO)
	DummyPO.Read()
	assert.Exactly(t, "active", DummyPO.DocumentStatus)
}

func TestPurchaseReturn1(t *testing.T) {
	// Case 1
	// Belum melakukan return sekalipun

	DummyItem := model.DummyItem()

	DummyItemV1 := model.DummyItemVariant()
	DummyItemV1.Item = DummyItem
	DummyItemV1.BasePrice = 100
	DummyItemV1.Save()

	DummyItemV2 := model.DummyItemVariant()
	DummyItemV2.Item = DummyItem
	DummyItemV2.BasePrice = 100
	DummyItemV2.Save()

	DummyPO := model.DummyPurchaseOrder()

	DummyPOI1 := model.DummyPurchaseOrderItem()
	DummyPOI1.PurchaseOrder = DummyPO
	DummyPOI1.Quantity = 20
	DummyPOI1.ItemVariant = DummyItemV1
	DummyPOI1.Save()

	DummyWR := model.DummyWorkorderReceiving()
	DummyWR.PurchaseOrder = DummyPO
	DummyWR.DocumentStatus = "finished"
	DummyWR.IsDeleted = 0
	DummyWR.Save()

	DummyWRI11 := model.DummyWorkorderReceivingItem()
	DummyWRI11.WorkorderReceiving = DummyWR
	DummyWRI11.PurchaseOrderItem = DummyPOI1
	DummyWRI11.Quantity = 15
	DummyWRI11.Save()

	DummyWRI12 := model.DummyWorkorderReceivingItem()
	DummyWRI12.WorkorderReceiving = DummyWR
	DummyWRI12.PurchaseOrderItem = DummyPOI1
	DummyWRI12.Quantity = 5
	DummyWRI12.Save()

	purchaseitems := ReturningPurchaseOrder(DummyPO, 0)

	for _, v := range *purchaseitems {
		assert.Exactly(t, float32(20), v.CanBeReturn)
	}
}

func TestPurchaseReturn2(t *testing.T) {
	// Case 2
	// Belum melakukan return sekalipun dan terdapat lebih dari satu purchase order item

	DummyItem := model.DummyItem()

	DummyItemV1 := model.DummyItemVariant()
	DummyItemV1.Item = DummyItem
	DummyItemV1.BasePrice = 100
	DummyItemV1.Save()

	DummyItemV2 := model.DummyItemVariant()
	DummyItemV2.Item = DummyItem
	DummyItemV2.BasePrice = 100
	DummyItemV2.Save()

	DummyPO := model.DummyPurchaseOrder()

	DummyPOI1 := model.DummyPurchaseOrderItem()
	DummyPOI1.PurchaseOrder = DummyPO
	DummyPOI1.Quantity = 20
	DummyPOI1.ItemVariant = DummyItemV1
	DummyPOI1.Save()

	DummyPOI2 := model.DummyPurchaseOrderItem()
	DummyPOI2.PurchaseOrder = DummyPO
	DummyPOI2.Quantity = 10
	DummyPOI2.ItemVariant = DummyItemV2
	DummyPOI2.Save()

	DummyWR := model.DummyWorkorderReceiving()
	DummyWR.PurchaseOrder = DummyPO
	DummyWR.DocumentStatus = "finished"
	DummyWR.IsDeleted = 0
	DummyWR.Save()

	DummyWRI11 := model.DummyWorkorderReceivingItem()
	DummyWRI11.WorkorderReceiving = DummyWR
	DummyWRI11.PurchaseOrderItem = DummyPOI1
	DummyWRI11.Quantity = 15
	DummyWRI11.Save()

	DummyWRI12 := model.DummyWorkorderReceivingItem()
	DummyWRI12.WorkorderReceiving = DummyWR
	DummyWRI12.PurchaseOrderItem = DummyPOI1
	DummyWRI12.Quantity = 5
	DummyWRI12.Save()

	DummyWRI21 := model.DummyWorkorderReceivingItem()
	DummyWRI21.WorkorderReceiving = DummyWR
	DummyWRI21.PurchaseOrderItem = DummyPOI2
	DummyWRI21.Quantity = 10
	DummyWRI21.Save()

	purchaseitems := ReturningPurchaseOrder(DummyPO, 0)
	var workreceive float32
	for _, v := range *purchaseitems {
		orm.NewOrm().Raw("SELECT SUM(woi.quantity) AS qty FROM workorder_receiving_item woi "+
			"JOIN workorder_receiving wo ON wo.id = woi.workorder_receiving_id "+
			"WHERE woi.workorder_receiving_id = wo.id AND wo.purchase_order_id = ? AND woi.purchase_order_item_id = ? AND wo.is_deleted = 0 AND wo.document_status = 'finished';", DummyPO.ID, v.ID).QueryRow(&workreceive)
		assert.Exactly(t, workreceive, v.CanBeReturn)
	}
}

func TestPurchaseReturn3(t *testing.T) {
	// Case 3
	// Sudah melakukan return dan terdapat lebih dari satu purchase order item

	DummyItem := model.DummyItem()

	DummyItemV1 := model.DummyItemVariant()
	DummyItemV1.Item = DummyItem
	DummyItemV1.BasePrice = 100
	DummyItemV1.Save()

	DummyItemV2 := model.DummyItemVariant()
	DummyItemV2.Item = DummyItem
	DummyItemV2.BasePrice = 100
	DummyItemV2.Save()

	DummyPO := model.DummyPurchaseOrder()
	DummyPO.IsDeleted = 0
	DummyPO.Save()

	DummyPOI1 := model.DummyPurchaseOrderItem()
	DummyPOI1.PurchaseOrder = DummyPO
	DummyPOI1.Quantity = 20
	DummyPOI1.ItemVariant = DummyItemV1
	DummyPOI1.Save()

	DummyPOI2 := model.DummyPurchaseOrderItem()
	DummyPOI2.PurchaseOrder = DummyPO
	DummyPOI2.Quantity = 10
	DummyPOI2.ItemVariant = DummyItemV2
	DummyPOI2.Save()

	DummyWR := model.DummyWorkorderReceiving()
	DummyWR.PurchaseOrder = DummyPO
	DummyWR.DocumentStatus = "finished"
	DummyWR.IsDeleted = 0
	DummyWR.Save()

	DummyWRI11 := model.DummyWorkorderReceivingItem()
	DummyWRI11.WorkorderReceiving = DummyWR
	DummyWRI11.PurchaseOrderItem = DummyPOI1
	DummyWRI11.Quantity = 15
	DummyWRI11.Save()

	DummyWRI12 := model.DummyWorkorderReceivingItem()
	DummyWRI12.WorkorderReceiving = DummyWR
	DummyWRI12.PurchaseOrderItem = DummyPOI1
	DummyWRI12.Quantity = 5
	DummyWRI12.Save()

	DummyWRI21 := model.DummyWorkorderReceivingItem()
	DummyWRI21.WorkorderReceiving = DummyWR
	DummyWRI21.PurchaseOrderItem = DummyPOI2
	DummyWRI21.Quantity = 10
	DummyWRI21.Save()

	DummyPR := model.DummyPurchaseReturn()
	DummyPR.PurchaseOrder = DummyPO
	DummyPR.DocumentStatus = "active"
	DummyPR.IsDeleted = 0
	DummyPR.Save()

	DummyPRI := model.DummyPurchaseReturnItem()
	DummyPRI.PurchaseOrderItem = DummyPOI1
	DummyPRI.PurchaseReturn = DummyPR
	DummyPRI.Quantity = 1
	DummyPRI.Save()

	purchaseitems := ReturningPurchaseOrder(DummyPO, 0)
	var workreceive, returnitem float32
	for _, v := range *purchaseitems {
		orm.NewOrm().Raw("SELECT SUM(woi.quantity) AS qty FROM workorder_receiving_item woi, workorder_receiving wo WHERE woi.workorder_receiving_id = wo.id AND wo.purchase_order_id = ? AND woi.purchase_order_item_id = ? AND wo.is_deleted = 0;", DummyPO.ID, v.ID).QueryRow(&workreceive)
		orm.NewOrm().Raw("SELECT SUM(pri.quantity) AS qty "+
			"FROM purchase_return_item pri JOIN purchase_return pr ON pr.id = pri.purchase_return_id "+
			"WHERE pri.purchase_return_id = pr.id AND pr.purchase_order_id = ? AND pri.purchase_order_item_id = ? AND pr.is_deleted = 0 AND pr.document_status != 'cancelled';", DummyPO.ID, v.ID).QueryRow(&returnitem)

		assert.Exactly(t, workreceive-returnitem, v.CanBeReturn)
	}
}

func TestPurchaseReturn4(t *testing.T) {
	// Case 4
	// Is deleted pada purchase return yang ada adalah 1

	DummyItem := model.DummyItem()

	DummyItemV1 := model.DummyItemVariant()
	DummyItemV1.Item = DummyItem
	DummyItemV1.BasePrice = 100
	DummyItemV1.Save()

	DummyItemV2 := model.DummyItemVariant()
	DummyItemV2.Item = DummyItem
	DummyItemV2.BasePrice = 100
	DummyItemV2.Save()

	DummyPO := model.DummyPurchaseOrder()

	DummyPOI1 := model.DummyPurchaseOrderItem()
	DummyPOI1.PurchaseOrder = DummyPO
	DummyPOI1.Quantity = 20
	DummyPOI1.ItemVariant = DummyItemV1
	DummyPOI1.Save()

	DummyPOI2 := model.DummyPurchaseOrderItem()
	DummyPOI2.PurchaseOrder = DummyPO
	DummyPOI2.Quantity = 10
	DummyPOI2.ItemVariant = DummyItemV2
	DummyPOI2.Save()

	DummyWR := model.DummyWorkorderReceiving()
	DummyWR.PurchaseOrder = DummyPO
	DummyWR.DocumentStatus = "finished"
	DummyWR.IsDeleted = 0
	DummyWR.Save()

	DummyWRI11 := model.DummyWorkorderReceivingItem()
	DummyWRI11.WorkorderReceiving = DummyWR
	DummyWRI11.PurchaseOrderItem = DummyPOI1
	DummyWRI11.Quantity = 15
	DummyWRI11.Save()

	DummyWRI12 := model.DummyWorkorderReceivingItem()
	DummyWRI12.WorkorderReceiving = DummyWR
	DummyWRI12.PurchaseOrderItem = DummyPOI1
	DummyWRI12.Quantity = 5
	DummyWRI12.Save()

	DummyWRI21 := model.DummyWorkorderReceivingItem()
	DummyWRI21.WorkorderReceiving = DummyWR
	DummyWRI21.PurchaseOrderItem = DummyPOI2
	DummyWRI21.Quantity = 10
	DummyWRI21.Save()

	DummyPR := model.DummyPurchaseReturn()
	DummyPR.PurchaseOrder = DummyPO
	DummyPR.Save()

	DummyPRI := model.DummyPurchaseReturnItem()
	DummyPRI.PurchaseOrderItem = DummyPOI1
	DummyPRI.PurchaseReturn = DummyPR
	DummyPRI.Quantity = 1
	DummyPRI.Save()

	DummyPR.IsDeleted = 1
	DummyPR.Save()

	purchaseitems := ReturningPurchaseOrder(DummyPO, 0)
	var workreceive, returnitem float32
	for _, v := range *purchaseitems {
		orm.NewOrm().Raw("SELECT SUM(woi.quantity) AS qty FROM workorder_receiving_item woi, workorder_receiving wo WHERE woi.workorder_receiving_id = wo.id AND wo.purchase_order_id = ? AND woi.purchase_order_item_id = ? AND wo.is_deleted = 0;", DummyPO.ID, v.ID).QueryRow(&workreceive)
		orm.NewOrm().Raw("SELECT SUM(pri.quantity) AS qty FROM purchase_return_item pri, purchase_return pr WHERE pri.purchase_return_id = pr.id AND pr.purchase_order_id = ? AND pri.purchase_order_item_id = ? AND pr.is_deleted = 0;", DummyPO.ID, v.ID).QueryRow(&returnitem)
		assert.Exactly(t, workreceive-returnitem, v.CanBeReturn)
	}
}

func TestGetAllPurchaseOrderNoData(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM purchase_order").Exec()
	rq := orm.RequestQuery{}
	// test tidak ada data sales_order
	m, total, e := GetAllPurchaseOrder(&rq)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, m)
	assert.NoError(t, e)
}

func TestGetAllPurchaseOrder(t *testing.T) {
	o := orm.NewOrm()
	o.Raw("DELETE FROM purchase_order").Exec()
	// buat dummy
	po := model.DummyPurchaseOrder()
	po.IsDeleted = int8(0)
	po.Save()
	po2 := model.DummyPurchaseOrder()
	po2.IsDeleted = int8(1)
	po2.Save()
	rq := orm.RequestQuery{}
	m, total, e := GetAllPurchaseOrder(&rq)
	assert.NoError(t, e)
	assert.Equal(t, int64(1), total)
	assert.NotEmpty(t, m)
	for _, u := range *m {
		assert.Equal(t, po.ID, u.ID)
	}
}

func TestGetDetailPurchaseOrderNoData(t *testing.T) {
	m, e := GetDetailPurchaseOrder("id", 99999)
	assert.Error(t, e)
	assert.Empty(t, m)
}

func TestGetDetailPurchaseOrderDeletedData(t *testing.T) {
	sls := model.DummyPurchaseOrder()
	sls.IsDeleted = int8(1)
	sls.Save()

	m, e := GetDetailPurchaseOrder("id", sls.ID)
	assert.Error(t, e)
	assert.Empty(t, m)
}

func TestGetDetailPurchaseOrderSuccess(t *testing.T) {
	po := model.DummyPurchaseOrder()
	po.IsDeleted = int8(0)
	po.Save()
	poi := model.DummyPurchaseOrderItem()
	poi.PurchaseOrder = &model.PurchaseOrder{ID: po.ID}
	poi.Save()

	m, e := GetDetailPurchaseOrder("id", po.ID)
	assert.NoError(t, e)
	assert.NotEmpty(t, m)
	assert.Equal(t, po.TotalCharge, m.TotalCharge)
	assert.Equal(t, po.CreatedBy.FullName, m.CreatedBy.FullName)
	for _, u := range m.PurchaseOrderItems {
		assert.Equal(t, poi.ID, u.ID)
		assert.Equal(t, poi.Discount, u.Discount)
		assert.Equal(t, poi.Note, u.Note)
	}

}

func TestCreatePurchaseOrder(t *testing.T) {
	orm.NewOrm().Raw("DELETE FROM purchase_order")
	orm.NewOrm().Raw("DELETE FROM partnership")

	DummyCustomer := model.DummyPartnership()
	DummyCustomer.TotalCredit = 1000
	DummyCustomer.TotalExpenditure = 1000
	DummyCustomer.PartnershipType = "supplier"
	DummyCustomer.Save()

	var ModelPOI []*model.PurchaseOrderItem

	poi := &model.PurchaseOrderItem{
		ItemVariant: model.DummyItemVariant(),
		Quantity:    float32(10),
		UnitPrice:   float64(1000),
		Discount:    float32(0),
		Subtotal:    float64(10000),
		Note:        "what ever...",
	}

	ModelPOI = append(ModelPOI, poi)

	ModelPO := &model.PurchaseOrder{
		Supplier:           DummyCustomer,
		EtaDate:            time.Now(),
		DocumentStatus:     "active",
		InvoiceStatus:      "active",
		ReceivingStatus:    "new",
		RecognitionDate:    time.Now(),
		IsPercentage:       int8(1),
		AutoInvoiced:       int8(1),
		Discount:           float32(10),
		DiscountAmount:     float64(1000),
		Tax:                float32(0),
		TaxAmount:          float64(0),
		ShipmentCost:       float64(1000),
		TotalCharge:        float64(1000),
		CreatedBy:          model.DummyUser(),
		CreatedAt:          time.Now(),
		PurchaseOrderItems: ModelPOI,
	}

	// Test purchase order yang dibuat
	ResultPO, err := CreatePurchaseOrder(ModelPO)
	assert.NoError(t, err, "Tidak boleh error")
	assert.NotNil(t, ResultPO, "Tidak boleh kosong jika berhasil")

	// Test partnership nya
	DummyCustomer.Read()
	assert.Equal(t, float64(2000), DummyCustomer.TotalCredit, "Total credit harus sesuai")
	assert.Equal(t, float64(2000), DummyCustomer.TotalExpenditure, "Total credit harus sesuai")

	// Dapatkan data invoice
	var pinvoice model.PurchaseInvoice
	err = orm.NewOrm().Raw("SELECT pi.* FROM purchase_invoice pi "+
		"WHERE pi.purchase_order_id = ?", ResultPO.ID).QueryRow(&pinvoice)
	assert.Equal(t, "new", pinvoice.DocumentStatus, "Harus sama, actual: "+fmt.Sprint(pinvoice.DocumentStatus, err))
}

func TestUpdatePurchaseOrder(t *testing.T) {
	po := model.DummyPurchaseOrder()

	var items []*model.PurchaseOrderItem
	item := model.DummyPurchaseOrderItem()
	item.Quantity = 5
	item.PurchaseOrder = po
	item.Save("Quantity", "PurchaseOrder")

	item2 := model.DummyPurchaseOrderItem()
	item2.Quantity = 8
	item2.PurchaseOrder = po
	item2.Save("Quantity", "PurchaseOrder")
	items = append(items, item, item2)

	po.PurchaseOrderItems = items
	po.IsDeleted = 0
	po.TotalCharge = 10000
	po.Save("PurchaseOrderItems", "IsDeleted", "TotalCharge")

	// item yang baru
	var itemsNew []*model.PurchaseOrderItem
	itemNew := &model.PurchaseOrderItem{ID: item.ID, PurchaseOrder: po, ItemVariant: model.DummyItemVariant(), Quantity: 1, Note: "abc"}
	itemNew2 := &model.PurchaseOrderItem{PurchaseOrder: po, ItemVariant: model.DummyItemVariant(), Quantity: 3, Note: "xxxx"}
	itemsNew = append(itemsNew, itemNew, itemNew2)

	supplier := model.DummyPartnership()
	supplier.TotalExpenditure = 10000
	supplier.Save("TotalDebt")

	soReq := &model.PurchaseOrder{ID: po.ID, TotalCharge: 8000, PurchaseOrderItems: items, Supplier: supplier, CreatedBy: model.DummyUser(), DocumentStatus: "new", InvoiceStatus: "new", ReceivingStatus: "new"}

	// update data yang di so.PurchaseOrderItems
	res, e := UpdatePurchaseOrder(soReq, itemsNew)
	assert.NoError(t, e)
	assert.NotEmpty(t, res)
	assert.Equal(t, 2, len(res.PurchaseOrderItems))

	data, _ := GetDetailPurchaseOrder("id", po.ID)
	assert.Equal(t, 2, len(data.PurchaseOrderItems))

	supplier = &model.Partnership{ID: soReq.Supplier.ID}
	supplier.Read()
	assert.Equal(t, float64(-2000), supplier.TotalExpenditure)
}

func TestGetItemVariantStockByPurchaseOrder(t *testing.T) {
	//purchase order
	po := model.DummyPurchaseOrder()
	po.IsDeleted = 0
	po.DocumentStatus = "new"
	po.Save("IsDeleted", "DocumentStatus")

	for a := 0; a < 2; a++ {
		// receiving
		receiving := model.DummyWorkorderReceiving()
		receiving.PurchaseOrder = po
		receiving.Save("PurchaseOrder")

		// item variant stock
		ivs := model.DummyItemVariantStock()

		//itemvariantstocklog
		ivslog := model.DummyItemVariantStockLog()
		ivslog.ItemVariantStock = ivs
		ivslog.RefType = "workorder_receiving"
		ivslog.RefID = uint64(receiving.ID)
		ivslog.LogType = "in"
		ivslog.Save("ItemVariantStock", "RefType", "RefID", "LogType")
	}

	res, e := getItemVariantStockByPurchaseOrder(po)
	assert.NoError(t, e)
	assert.NotEmpty(t, res)
	assert.Equal(t, 2, len(res))
}

func TestGetPurchaseInvoiceByPurchaseOrder(t *testing.T) {
	//purchase order
	po := model.DummyPurchaseOrder()
	po.IsDeleted = 0
	po.DocumentStatus = "new"
	po.Save("IsDeleted", "DocumentStatus")

	for a := 0; a < 2; a++ {
		// purchaceinvoice
		invoice := model.DummyPurchaseInvoice()
		invoice.PurchaseOrder = po
		invoice.IsDeleted = 0
		invoice.Save("PurchaseOrder", "IsDeleted")
	}

	res, e := getPurchaseInvoiceByPurchaseOrder(po)
	assert.NoError(t, e)
	assert.NotEmpty(t, res)
	assert.Equal(t, 2, len(res))
}

func TestGetFinanceExpenceByPurchaseInvoice(t *testing.T) {
	//purchase order
	po := model.DummyPurchaseOrder()
	po.IsDeleted = 0
	po.DocumentStatus = "new"
	po.Save("IsDeleted", "DocumentStatus")

	pi := model.DummyPurchaseInvoice()
	pi.Save()

	for a := 0; a < 2; a++ {
		// financeexpense
		finance := model.DummyFinanceExpense()
		finance.RefID = uint64(pi.ID)
		finance.RefType = "purchase_invoice"
		finance.IsDeleted = 0
		finance.Save("RefID", "IsDeleted", "RefType")
	}

	res, e := getFinanceExpenceByPurchaseInvoice(pi)
	assert.NoError(t, e)
	assert.NotEmpty(t, res)
	assert.Equal(t, 2, len(res))
}

func TestGetWorkorderReceivingByPurchaseOrder(t *testing.T) {
	//purchase order
	po := model.DummyPurchaseOrder()
	po.IsDeleted = 0
	po.DocumentStatus = "new"
	po.Save("IsDeleted", "DocumentStatus")

	for a := 0; a < 2; a++ {
		receiving := model.DummyWorkorderReceiving()
		receiving.PurchaseOrder = po
		receiving.IsDeleted = 0
		receiving.Save("PurchaseOrder", "IsDeleted")
	}

	res, e := getWorkorderReceivingByPurchaseOrder(po)
	assert.NoError(t, e)
	assert.NotEmpty(t, res)
	assert.Equal(t, 2, len(res))
}

func TestGetItemVariantStockLogByReceiving(t *testing.T) {
	orm.NewOrm().Raw("DELETE FROM workorder_receiving").Exec()
	orm.NewOrm().Raw("DELETE FROM item_variant_stock_log").Exec()

	// receiving
	receiving := model.DummyWorkorderReceiving()
	receiving.IsDeleted = 0
	receiving.Save("IsDeleted")

	for a := 0; a < 2; a++ {
		//itemvariantstocklog
		ivslog := model.DummyItemVariantStockLog()
		ivslog.RefType = "workorder_receiving"
		ivslog.RefID = uint64(receiving.ID)
		ivslog.LogType = "out"
		ivslog.Quantity = 10
		ivslog.Save("RefType", "RefID", "LogType")
	}

	res, e := getItemVariantStockLogByReceiving(receiving)
	assert.NoError(t, e)
	assert.NotEmpty(t, res)
	assert.Equal(t, 2, len(res))

	for _, i := range res {
		assert.NotEmpty(t, i.Quantity)
		assert.NotEmpty(t, i.ItemVariantStock.AvailableStock)
	}
}

func TestCancelPurchaseOrder(t *testing.T) {
	//purchase order
	po := model.DummyPurchaseOrder()
	po.IsDeleted = 0
	po.DocumentStatus = "new"
	po.InvoiceStatus = "active"
	po.ReceivingStatus = "active"
	po.Save("IsDeleted", "DocumentStatus", "InvoiceStatus", "ReceivingStatus")

	for a := 0; a < 2; a++ {
		// purchaceinvoice
		invoice := model.DummyPurchaseInvoice()
		invoice.PurchaseOrder = po
		invoice.IsDeleted = 0
		invoice.DocumentStatus = "active"
		invoice.Save("PurchaseOrder", "IsDeleted", "DocumentStatus")

		// financeexpense
		finance := model.DummyFinanceExpense()
		finance.RefID = uint64(invoice.ID)
		finance.RefType = "purchase_invoice"
		finance.IsDeleted = 0
		finance.Save("RefID", "IsDeleted", "RefType")

		receiving := model.DummyWorkorderReceiving()
		receiving.PurchaseOrder = po
		receiving.IsDeleted = 0
		receiving.Save("PurchaseOrder", "IsDeleted")

		//itemvariantstocklog
		ivslog := model.DummyItemVariantStockLog()
		ivslog.RefType = "workorder_receiving"
		ivslog.RefID = uint64(receiving.ID)
		ivslog.LogType = "out"
		ivslog.Save("RefType", "RefID", "LogType")
	}

	res, e := cancelPurchaseOrder(po)
	assert.NoError(t, e)
	assert.NotEmpty(t, res)
}

func TestCreatePurchaseOrder2(t *testing.T) {
	DummyCustomer := model.DummyPartnership()
	DummyCustomer.PartnershipType = "supplier"
	DummyCustomer.Save()

	var ModelPOI []*model.PurchaseOrderItem

	poi := &model.PurchaseOrderItem{
		ItemVariant: model.DummyItemVariant(),
		Quantity:    float32(10),
		UnitPrice:   float64(1000),
		Discount:    float32(0),
		Subtotal:    float64(10000),
		Note:        "what ever...",
	}

	ModelPOI = append(ModelPOI, poi)

	ModelPO := &model.PurchaseOrder{
		Supplier:           DummyCustomer,
		EtaDate:            time.Now(),
		DocumentStatus:     "active",
		InvoiceStatus:      "active",
		ReceivingStatus:    "new",
		RecognitionDate:    time.Now(),
		IsPercentage:       int8(1),
		AutoInvoiced:       int8(1),
		Discount:           float32(10),
		DiscountAmount:     float64(1000),
		Tax:                float32(0),
		TaxAmount:          float64(0),
		ShipmentCost:       float64(1000),
		TotalCharge:        float64(0),
		CreatedBy:          model.DummyUser(),
		CreatedAt:          time.Now(),
		PurchaseOrderItems: ModelPOI,
	}

	ResultPO, err := CreatePurchaseOrder(ModelPO)
	assert.NoError(t, err, "Tidak boleh ada error "+fmt.Sprint("Actual:", err))
	assert.NotNil(t, ResultPO, "Tidak boleh kosong")
	assert.Equal(t, 1, len(ResultPO.PurchaseOrderItems), "Harus sesuai Jumlah nya.")
	assert.Equal(t, "active", ResultPO.InvoiceStatus, "Harus sesuai statusnya.")

	var pinvoice model.PurchaseInvoice
	orm.NewOrm().Raw("SELECT * FROM purchase_invoice WHERE purchase_order_id = ?", ResultPO.ID).QueryRow(&pinvoice)
	assert.Equal(t, "new", pinvoice.DocumentStatus, "Harusnya sesuai "+fmt.Sprint("Actual:", pinvoice.DocumentStatus))
}

func TestCreateRequest_Validate_Transform(t *testing.T) {
	DummyPartner := model.DummyPartnership()
	DummyPartner.PartnershipType = "supplier"
	DummyPartner.Save()

	DummyItemVar := model.DummyItemVariant()
	DummyItemVar.IsDeleted = 0
	DummyItemVar.IsArchived = 0
	DummyItemVar.Save()

	request := createRequest{
		SupplierID:      common.Encrypt(DummyPartner.ID),
		Note:            "What Ever",
		IsPercentage:    1,
		AutoInvoiced:    1,
		EtaDate:         time.Now(),
		RecognitionDate: time.Now(),
		ShipmentCost:    5000,
		Discount:        10,
		Tax:             10,
		PurchaseOrderItems: []purchaseOrderItemRequest{
			{
				Note:          "Whaaattttt...",
				Discount:      0,
				UnitPrice:     1000,
				Quantity:      10,
				ItemVariantID: common.Encrypt(DummyItemVar.ID),
			},
		},
	}

	POrderHasil := request.Transform(model.DummyUser())
	// Testing Tax Amount
	assert.Equal(t, float64(900), POrderHasil.TaxAmount, "Harus sama")
	// Testing SubTotal
	assert.Equal(t, float64(10000), POrderHasil.PurchaseOrderItems[0].Subtotal, "Harus sama")
	// Testing Discount Amount
	assert.Equal(t, float64(1000), POrderHasil.DiscountAmount, "Harus sama")
	// Testing Total Charge
	assert.Equal(t, float64((10000-1000)+900+5000), POrderHasil.TotalCharge, "Harus sama")
	// Invoice Status
	assert.Equal(t, "active", POrderHasil.InvoiceStatus, "Harus sama")

	request = createRequest{
		SupplierID:      common.Encrypt(DummyPartner.ID),
		Note:            "What Ever",
		IsPercentage:    0,
		AutoInvoiced:    0,
		EtaDate:         time.Now(),
		RecognitionDate: time.Now(),
		ShipmentCost:    5000,
		DiscountAmount:  5000,
		Tax:             10,
		PurchaseOrderItems: []purchaseOrderItemRequest{
			{
				Note:          "Whaaattttt...",
				Discount:      0,
				UnitPrice:     1000,
				Quantity:      10,
				ItemVariantID: common.Encrypt(DummyItemVar.ID),
			},
		},
	}

	POrderHasil = request.Transform(model.DummyUser())
	// Testing Tax Amount
	assert.Equal(t, float64(1000), POrderHasil.TaxAmount, "Harus sama")
	// Testing SubTotal
	assert.Equal(t, float64(10000), POrderHasil.PurchaseOrderItems[0].Subtotal, "Harus sama")
	// Testing Discount Amount
	assert.Equal(t, float64(5000), POrderHasil.DiscountAmount, "Harus sama")
	// Testing Total Charge
	assert.Equal(t, float64((10000-5000)+1000+5000), POrderHasil.TotalCharge, "Harus sama")
	// Invoice Status
	assert.Equal(t, "new", POrderHasil.InvoiceStatus, "Harus sama")

	request = createRequest{
		SupplierID:      common.Encrypt(DummyPartner.ID),
		Note:            "What Ever",
		IsPercentage:    0,
		AutoInvoiced:    0,
		EtaDate:         time.Now(),
		RecognitionDate: time.Now(),
		ShipmentCost:    5000,
		Tax:             10,
		PurchaseOrderItems: []purchaseOrderItemRequest{
			{
				Note:          "Whaaattttt...",
				Discount:      0,
				UnitPrice:     1000,
				Quantity:      10,
				ItemVariantID: common.Encrypt(DummyItemVar.ID),
			},
		},
	}

	POrderHasil = request.Transform(model.DummyUser())
	// Testing Tax Amount
	assert.Equal(t, float64(1000), POrderHasil.TaxAmount, "Harus sama")
	// Testing SubTotal
	assert.Equal(t, float64(10000), POrderHasil.PurchaseOrderItems[0].Subtotal, "Harus sama")
	// Testing Discount Amount
	assert.Equal(t, float64(0), POrderHasil.DiscountAmount, "Harus sama")
	// Testing Total Charge
	assert.Equal(t, float64(10000+1000+5000), POrderHasil.TotalCharge, "Harus sama")
}
