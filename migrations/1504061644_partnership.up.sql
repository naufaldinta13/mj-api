SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `partnership` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(120) NOT NULL,
  `partnership_type` ENUM('customer', 'supplier') NULL DEFAULT 'customer' COMMENT 'Type of business relationship',
  `order_rule` ENUM('none', 'one_bill', 'plafon') NULL DEFAULT 'none' COMMENT 'Rule of ordering : onebill (cannot order if any invoice still unpaid), plafon (cannot order if exceed max plafon)',
  `full_name` VARCHAR(45) NOT NULL,
  `email` VARCHAR(45) NULL DEFAULT NULL,
  `phone` VARCHAR(45) NULL DEFAULT NULL,
  `address` TINYTEXT NULL DEFAULT NULL,
  `city` VARCHAR(45) NULL DEFAULT NULL,
  `province` VARCHAR(45) NULL DEFAULT NULL,
  `bank_name` VARCHAR(45) NULL DEFAULT NULL COMMENT 'name of bank ',
  `bank_number` VARCHAR(45) NULL DEFAULT NULL COMMENT 'number of bank account',
  `bank_holder` VARCHAR(45) NULL DEFAULT NULL COMMENT 'name of bank account holder',
  `max_plafon` DECIMAL(20,0) NULL DEFAULT 0 COMMENT 'maximum amount that can be lent',
  `total_debt` DECIMAL(20,0) NULL DEFAULT 0 COMMENT 'total of payable',
  `total_credit` DECIMAL(20,0) NULL DEFAULT 0 COMMENT 'total of receivable',
  `total_spend` DECIMAL(20,0) NULL DEFAULT 0 COMMENT 'total of sales transaction',
  `total_expenditure` DECIMAL(20,0) NULL,
  `sales_person` VARCHAR(45) NULL DEFAULT NULL,
  `visit_day` VARCHAR(45) NULL DEFAULT NULL,
  `note` TINYTEXT NULL DEFAULT NULL,
  `is_archived` TINYINT(1) NULL DEFAULT 0,
  `is_deleted` TINYINT(1) NULL,
  `is_default` TINYINT(1) NULL DEFAULT 0 COMMENT 'Default customer for SO',
  `created_by` BIGINT(20) UNSIGNED NOT NULL,
  `updated_by` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_partnership_1_idx` (`created_by` ASC),
  INDEX `fk_partnership_2_idx` (`updated_by` ASC),
  CONSTRAINT `fk_partnership_1`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_partnership_2`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
