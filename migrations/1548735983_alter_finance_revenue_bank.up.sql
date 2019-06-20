SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `finance_revenue`
ADD COLUMN `bank_account_id` BIGINT(11) UNSIGNED NULL DEFAULT NULL;

ALTER TABLE `finance_revenue`
ADD INDEX `fk_finance_revenue_7_idx` (`bank_account_id` ASC);
ALTER TABLE `finance_revenue`
ADD CONSTRAINT `fk_finance_revenue_7`
  FOREIGN KEY (`bank_account_id`)
  REFERENCES `bank_account` (`id`)
  ON DELETE NO ACTION
  ON UPDATE NO ACTION;

ALTER TABLE `sales_order`
DROP FOREIGN KEY `fk_sales_order_7`;

ALTER TABLE `sales_order`
DROP `bank_account_id`,
DROP `bank_name`,
DROP `bank_number`;