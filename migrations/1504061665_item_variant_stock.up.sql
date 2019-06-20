SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `item_variant_stock` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `item_variant_id` BIGINT(20) UNSIGNED NOT NULL,
  `sku_code` VARCHAR(120) NOT NULL,
  `available_stock` FLOAT NULL DEFAULT 0,
  `unit_cost` DECIMAL(20,0) NULL DEFAULT 0 COMMENT 'base price of item variant',
  `created_by` BIGINT(20) UNSIGNED NULL,
  `updated_by` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_item_variant_stock_1_idx` (`item_variant_id` ASC),
  INDEX `fk_item_variant_stock_2_idx` (`created_by` ASC),
  INDEX `fk_item_variant_stock_3_idx` (`updated_by` ASC),
  CONSTRAINT `fk_item_variant_stock_1`
    FOREIGN KEY (`item_variant_id`)
    REFERENCES `item_variant` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_item_variant_stock_2`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_item_variant_stock_3`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB
COMMENT = 'sku table';
