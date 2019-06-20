SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE `sales_order`
DROP FOREIGN KEY `fk_sales_order_7`;
ALTER TABLE `sales_order`
DROP INDEX `fk_sales_order_7_idx` ;

ALTER TABLE `sales_order`
DROP COLUMN `bank_account_id`,
DROP COLUMN `bank_name`,
DROP COLUMN `bank_number`;

