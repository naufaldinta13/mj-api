SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `finance_expense`
 ADD `is_deleted` TINYINT(1) NULL DEFAULT 0;
 ALTER TABLE `finance_expense`
 DROP `code`;