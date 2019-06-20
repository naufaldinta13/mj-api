SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `purchase_return`
 ADD `is_bundled` TINYINT(1) NULL AFTER `document_status`;