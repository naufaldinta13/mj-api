SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `direct_placement_item`
CHANGE `unti_price` `unit_price` DECIMAL(20,0) UNSIGNED NOT NULL DEFAULT 0;
