SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `sales_invoice`
ADD `total_revenued` DECIMAL(20,0) NULL DEFAULT 0 AFTER `total_paid`;
