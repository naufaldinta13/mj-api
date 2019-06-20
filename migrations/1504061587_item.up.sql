SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `item` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `category_id` BIGINT(20) UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Default category 1 (uncategorized)',
  `item_type` ENUM('product', 'material', 'service') NULL DEFAULT 'product',
  `item_name` VARCHAR(200) NOT NULL,
  `note` TINYTEXT NULL DEFAULT NULL,
  `has_variant` TINYINT(1) NULL DEFAULT 0,
  `is_archived` TINYINT(1) NULL DEFAULT 0,
  `is_deleted` TINYINT(1) NULL DEFAULT 0,
  `created_by` BIGINT(20) UNSIGNED NOT NULL,
  `updated_by` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_item_1_idx` (`category_id` ASC),
  INDEX `fk_item_2_idx` (`created_by` ASC),
  INDEX `fk_item_3_idx` (`updated_by` ASC),
  CONSTRAINT `fk_item_1`
    FOREIGN KEY (`category_id`)
    REFERENCES `item_category` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_item_2`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_item_3`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
