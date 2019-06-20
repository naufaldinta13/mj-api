SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `invoice_receipt`
 ADD `total_invoice` DECIMAL(20,0) NULL AFTER `recognition_date`;
ALTER TABLE `invoice_receipt`
 ADD `total_return` DECIMAL(20,0) NULL AFTER `total_invoice`;