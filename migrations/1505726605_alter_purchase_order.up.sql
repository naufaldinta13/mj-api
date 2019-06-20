SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `purchase_order`
ADD `cancelled_note` tinytext NULL;