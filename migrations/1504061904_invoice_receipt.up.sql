SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `invoice_receipt` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `partnership_id` BIGINT(20) UNSIGNED NOT NULL,
  `code` VARCHAR(120) NOT NULL,
  `recognition_date` DATE NULL DEFAULT NULL,
  `total_amount` DECIMAL(20,0) NULL DEFAULT 0 COMMENT 'total amount of all invoices in the invoice receipt',
  `note` TINYTEXT NULL DEFAULT NULL,
  `document_status` ENUM('new', 'active', 'finished') NULL DEFAULT 'new',
  `created_by` BIGINT(20) UNSIGNED NOT NULL,
  `updated_by` BIGINT(20) UNSIGNED NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_invoice_receipt_2_idx` (`created_by` ASC),
  INDEX `fk_invoice_receipt_3_idx` (`updated_by` ASC),
  INDEX `fk_invoice_receipt_1_idx` (`partnership_id` ASC),
  CONSTRAINT `fk_invoice_receipt_1`
    FOREIGN KEY (`partnership_id`)
    REFERENCES `partnership` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_invoice_receipt_2`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_invoice_receipt_3`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB
COMMENT = 'sales invoice receipt (tanda terima invoice)';
