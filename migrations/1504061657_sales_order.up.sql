SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `sales_order` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `reference_id` BIGINT(20) UNSIGNED NULL COMMENT 'refer to cancelled SO',
  `customer_id` BIGINT(20) UNSIGNED NULL,
  `code` VARCHAR(45) NOT NULL,
  `recognition_date` DATE NULL DEFAULT NULL COMMENT 'order date',
  `eta_date` DATE NULL DEFAULT NULL COMMENT 'estimated arrival date',
  `discount` FLOAT UNSIGNED NULL DEFAULT 0,
  `tax` FLOAT UNSIGNED NULL DEFAULT 0,
  `discount_amount` DECIMAL(20,0) UNSIGNED NULL DEFAULT 0,
  `tax_amount` DECIMAL(20,0) UNSIGNED NULL DEFAULT 0,
  `shipment_address` TINYTEXT NULL DEFAULT NULL,
  `shipment_cost` DECIMAL(20,0) UNSIGNED NULL DEFAULT 0,
  `total_price` DECIMAL(20,0) UNSIGNED NOT NULL,
  `total_charge` DECIMAL(20,0) UNSIGNED NOT NULL,
  `total_paid` DECIMAL(20,0) UNSIGNED NULL DEFAULT 0,
  `total_cost` DECIMAL(20,0) UNSIGNED NULL DEFAULT 0,
  `note` TINYTEXT NULL DEFAULT NULL,
  `document_status` ENUM('new', 'active', 'finished', 'requested_cancel', 'approved_cancel') NULL DEFAULT 'new' COMMENT '\'requested_cancel\' by cashier usergroup, \'approved_cancel\' by owner / supervisor usergroup',
  `invoice_status` ENUM('new', 'active', 'finished') NULL DEFAULT 'new',
  `fulfillment_status` ENUM('new', 'active', 'finished') NULL DEFAULT 'new',
  `shipment_status` ENUM('new', 'active', 'finished') NULL,
  `auto_fulfillment` TINYINT(1) NULL DEFAULT 0,
  `auto_invoice` TINYINT(1) NULL DEFAULT 0,
  `is_deleted` TINYINT(1) NULL DEFAULT 0,
  `created_by` BIGINT(20) UNSIGNED NOT NULL,
  `updated_by` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
  `request_cancel_by` BIGINT(20) UNSIGNED NULL DEFAULT NULL COMMENT 'refer to user with cashier usergroup ',
  `approve_cancel_by` BIGINT(20) UNSIGNED NULL DEFAULT NULL COMMENT 'refer to user with owner or supervisor usergroup ',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  `request_cancel_at` TIMESTAMP NULL DEFAULT NULL,
  `approve_cancel_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_sales_order_1_idx` (`reference_id` ASC),
  INDEX `fk_sales_order_2_idx` (`customer_id` ASC),
  INDEX `fk_purchase_order_3_idx` (`created_by` ASC),
  INDEX `fk_sales_order_5_idx` (`request_cancel_by` ASC),
  INDEX `fk_sales_order_4_idx` (`updated_by` ASC),
  INDEX `fk_sales_order_6_idx` (`approve_cancel_by` ASC),
  CONSTRAINT `fk_sales_order_1`
    FOREIGN KEY (`reference_id`)
    REFERENCES `sales_order` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_order_2`
    FOREIGN KEY (`customer_id`)
    REFERENCES `partnership` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_order_3`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_order_4`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_order_5`
    FOREIGN KEY (`request_cancel_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_order_6`
    FOREIGN KEY (`approve_cancel_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
