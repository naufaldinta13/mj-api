SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `sales_return` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `sales_order_id` BIGINT(20) UNSIGNED NOT NULL,
  `recognition_date` DATE NOT NULL,
  `code` VARCHAR(45) NOT NULL,
  `total_amount` DECIMAL(20,0) NOT NULL DEFAULT 0 COMMENT 'total amount of sales refund',
  `note` TINYTEXT NULL DEFAULT NULL,
  `document_status` ENUM('new', 'active', 'finished') NULL DEFAULT 'new',
  `is_deleted` TINYINT(1) NULL DEFAULT 0,
  `created_by` BIGINT(20) UNSIGNED NOT NULL,
  `updated_by` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_sales_return_1_idx` (`sales_order_id` ASC),
  INDEX `fk_sales_return_3_idx` (`created_by` ASC),
  INDEX `fk_sales_return_4_idx` (`updated_by` ASC),
  CONSTRAINT `fk_sales_return_1`
    FOREIGN KEY (`sales_order_id`)
    REFERENCES `sales_order` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_return_3`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_return_4`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
