SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `item_variant` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `item_id` BIGINT(20) UNSIGNED NOT NULL,
  `measurement_id` BIGINT(20) UNSIGNED NOT NULL COMMENT 'FK to measurement table',
  `barcode` VARCHAR(120) NULL DEFAULT NULL COMMENT 'Barcode number of variant',
  `external_name` VARCHAR(200) NULL DEFAULT NULL COMMENT 'Common name of item variant',
  `variant_name` VARCHAR(45) NULL DEFAULT NULL COMMENT 'Variant name',
  `image` TINYTEXT NULL DEFAULT NULL,
  `base_price` DECIMAL(20,0) NULL,
  `note` TINYTEXT NULL DEFAULT NULL,
  `minimum_stock` FLOAT NULL DEFAULT 0,
  `available_stock` FLOAT NULL,
  `commited_stock` FLOAT NULL,
  `has_external_name` TINYINT(1) NULL DEFAULT 0 COMMENT 'to optimize item searching',
  `is_archived` TINYINT(1) NULL DEFAULT 0,
  `is_deleted` TINYINT(1) NULL DEFAULT 0,
  `created_by` BIGINT(20) UNSIGNED NOT NULL,
  `updated_by` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_item_variant_1_idx` (`item_id` ASC),
  INDEX `fk_item_variant_2_idx` (`measurement_id` ASC),
  INDEX `fk_item_variant_3_idx` (`created_by` ASC),
  INDEX `fk_item_variant_4_idx` (`updated_by` ASC),
  CONSTRAINT `fk_item_variant_1`
    FOREIGN KEY (`item_id`)
    REFERENCES `item` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_item_variant_2`
    FOREIGN KEY (`measurement_id`)
    REFERENCES `measurement` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_item_variant_3`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_item_variant_4`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB
COMMENT = 'Variant of item';
