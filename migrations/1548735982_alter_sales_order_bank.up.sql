SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `sales_order`
ADD COLUMN `bank_account_id` BIGINT(11) UNSIGNED NULL DEFAULT NULL AFTER `customer_id`,
ADD COLUMN `bank_name` VARCHAR(45) NULL DEFAULT NULL AFTER `is_reported`,
ADD COLUMN `bank_number` VARCHAR(45) NULL AFTER `bank_name`;

ALTER TABLE `sales_order` 
ADD INDEX `fk_sales_order_7_idx` (`bank_account_id` ASC);
ALTER TABLE `sales_order` 
ADD CONSTRAINT `fk_sales_order_7`
  FOREIGN KEY (`bank_account_id`)
  REFERENCES `bank_account` (`id`)
  ON DELETE NO ACTION
  ON UPDATE NO ACTION;
