SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `sales_order`
ADD `is_percentage_discount` TINYINT(1) NULL DEFAULT 0 AFTER `auto_invoice`;
