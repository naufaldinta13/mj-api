SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `purchase_order` ADD `is_percentage` tinyint(1) NULL DEFAULT '0' COMMENT 'jika is_percentage 1 maka discount amount wajib diisi ' AFTER `receiving_status`;