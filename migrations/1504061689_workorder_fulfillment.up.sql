SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `workorder_fulfillment` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `sales_order_id` BIGINT(20) UNSIGNED NOT NULL,
  `code` VARCHAR(120) NOT NULL,
  `priority` ENUM('routine', 'rush', 'emergency') NULL DEFAULT 'routine',
  `due_date` TIMESTAMP NULL DEFAULT NULL,
  `shipping_address` TINYTEXT NULL DEFAULT NULL,
  `note` TINYTEXT NULL DEFAULT NULL,
  `document_status` ENUM('new','active', 'finished') NULL DEFAULT 'new',
  `is_deleted` TINYINT(1) NULL DEFAULT 0,
  `is_delivered` TINYINT(1) NULL DEFAULT 0,
  `created_by` BIGINT(20) UNSIGNED NOT NULL,
  `updated_by` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_workorder_fulfillment_1_idx` (`sales_order_id` ASC),
  INDEX `fk_workorder_fulfillment_2_idx` (`created_by` ASC),
  INDEX `fk_workorder_fulfillment_3_idx` (`updated_by` ASC),
  CONSTRAINT `fk_workorder_fulfillment_1`
    FOREIGN KEY (`sales_order_id`)
    REFERENCES `sales_order` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_workorder_fulfillment_2`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_workorder_fulfillment_3`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
