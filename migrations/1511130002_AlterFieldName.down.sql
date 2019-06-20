SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE `item_variant`
CHANGE COLUMN `external_name` `external_name` VARCHAR(200) NULL DEFAULT NULL COMMENT 'Common name of item variant' ;

ALTER TABLE `item_variant`
CHANGE COLUMN `variant_name` `variant_name` VARCHAR(200) NULL DEFAULT NULL COMMENT 'Variant name' ;

ALTER TABLE `item`
CHANGE COLUMN `item_name` `item_name` VARCHAR(200) NOT NULL ;

