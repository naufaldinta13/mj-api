SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `finance_expense` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `ref_id` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
  `ref_type` ENUM('purchase_invoice', 'sales_return') NULL DEFAULT NULL,
  `recognition_date` DATE NOT NULL COMMENT 'paid date',
  `code` VARCHAR(45) NOT NULL,
  `amount` DECIMAL(20,0) UNSIGNED NOT NULL,
  `payment_method` ENUM('cash', 'debit_card', 'credit_card', 'giro') NULL DEFAULT 'cash',
  `bank_name` VARCHAR(45) NULL DEFAULT NULL,
  `bank_number` VARCHAR(45) NULL DEFAULT NULL,
  `bank_holder` VARCHAR(45) NULL DEFAULT NULL,
  `note` TINYTEXT NULL DEFAULT NULL,
  `document_status` ENUM('uncleared', 'cleared') NULL DEFAULT 'uncleared',
  `created_by` BIGINT(20) UNSIGNED NOT NULL,
  `updated_by` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_finance_expense_1_idx` (`created_by` ASC),
  INDEX `fk_finance_expense_2_idx` (`updated_by` ASC),
  CONSTRAINT `fk_finance_expense_1`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_finance_expense_2`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
