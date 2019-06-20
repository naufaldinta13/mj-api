SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `purchase_order` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `reference_id` BIGINT(20) UNSIGNED NULL DEFAULT NULL COMMENT 'Refer to cancelled PO',
  `supplier_id` BIGINT(20) UNSIGNED NOT NULL,
  `code` VARCHAR(45) NOT NULL COMMENT 'code of PO document',
  `recognition_date` DATE NULL DEFAULT NULL,
  `eta_date` DATE NULL DEFAULT NULL,
  `discount` FLOAT UNSIGNED NULL DEFAULT 0,
  `tax` FLOAT UNSIGNED NULL DEFAULT 0,
  `discount_amount` DECIMAL(20,0) UNSIGNED NULL DEFAULT 0 COMMENT 'nominal of discount',
  `tax_amount` DECIMAL(20,0) UNSIGNED NULL DEFAULT 0 COMMENT 'nominal of tax',
  `shipment_cost` DECIMAL(20,0) UNSIGNED NULL DEFAULT 0,
  `total_charge` DECIMAL(20,0) UNSIGNED NULL DEFAULT 0 COMMENT 'total amount that should be paid',
  `total_paid` DECIMAL(20,0) UNSIGNED NULL DEFAULT 0,
  `note` TINYTEXT NULL DEFAULT NULL,
  `document_status` ENUM('new', 'active', 'finished', 'cancelled') NULL DEFAULT 'new',
  `invoice_status` ENUM('new', 'active', 'finished') NULL DEFAULT 'new',
  `receiving_status` ENUM('new', 'active', 'finished') NULL DEFAULT 'new',
  `auto_invoiced` TINYINT(1) NULL DEFAULT 0,
  `is_deleted` TINYINT(1) NULL DEFAULT 0,
  `created_by` BIGINT(20) UNSIGNED NOT NULL,
  `updated_by` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_purchase_order_1_idx` (`reference_id` ASC),
  INDEX `fk_purchase_order_2_idx` (`supplier_id` ASC),
  INDEX `fk_purchase_order_3_idx` (`created_by` ASC),
  INDEX `fk_purchase_order_4_idx` (`updated_by` ASC),
  CONSTRAINT `fk_purchase_order_1`
    FOREIGN KEY (`reference_id`)
    REFERENCES `purchase_order` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_purchase_order_2`
    FOREIGN KEY (`supplier_id`)
    REFERENCES `partnership` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_purchase_order_3`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_purchase_order_4`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;