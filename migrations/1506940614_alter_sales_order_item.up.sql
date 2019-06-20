SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `sales_order_item`
ADD `quantity_fulfillment` float unsigned NOT NULL AFTER `quantity`;