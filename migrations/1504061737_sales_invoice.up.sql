SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `sales_invoice` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `sales_order_id` BIGINT(20) UNSIGNED NOT NULL,
  `code` VARCHAR(45) NOT NULL,
  `recognition_date` DATE NULL DEFAULT NULL,
  `due_date` DATE NULL DEFAULT NULL,
  `billing_address` TINYTEXT NULL DEFAULT NULL,
  `total_amount` DECIMAL(20,0) NULL DEFAULT 0,
  `total_paid` DECIMAL(20,0) NULL DEFAULT 0,
  `note` TINYTEXT NULL DEFAULT NULL,
  `document_status` ENUM('new', 'active', 'finished') NULL DEFAULT 'new',
  `is_bundled` TINYINT(1) NULL DEFAULT 0 COMMENT 'indicate that invoice has been bundled with other invoice in invoice receipt table',
  `is_deleted` TINYINT(1) NULL DEFAULT 0,
  `created_by` BIGINT(20) UNSIGNED NOT NULL,
  `updated_by` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_purchase_invoice_2_idx` (`created_by` ASC),
  INDEX `fk_purchase_invoice_3_idx` (`updated_by` ASC),
  INDEX `fk_sales_invoice_1_idx` (`sales_order_id` ASC),
  CONSTRAINT `fk_sales_invoice_1`
    FOREIGN KEY (`sales_order_id`)
    REFERENCES `sales_order` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_invoice_2`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_invoice_3`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
