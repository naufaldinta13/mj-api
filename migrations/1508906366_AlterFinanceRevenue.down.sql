SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `finance_revenue`
CHANGE `ref_type` `ref_type` ENUM('sales_invoice', 'purchase_return') NULL DEFAULT NULL;
