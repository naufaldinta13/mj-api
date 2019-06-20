SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `finance_expense`
CHANGE `ref_type` `ref_type` ENUM('purchase_invoice', 'sales_return', 'invoice_receipt') NULL DEFAULT NULL;
