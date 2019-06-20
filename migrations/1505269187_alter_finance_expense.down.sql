ALTER TABLE `finance_expense`
 DROP COLUMN `is_deleted`;
 ALTER TABLE `finance_expense`
 ADD `code` VARCHAR(45) NOT NULL;