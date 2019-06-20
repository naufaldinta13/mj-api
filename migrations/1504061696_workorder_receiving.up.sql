SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `workorder_receiving` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `purchase_order_id` BIGINT(20) UNSIGNED NULL,
  `recognition_date` DATE NOT NULL,
  `code` VARCHAR(45) NOT NULL,
  `pic` VARCHAR(45) NOT NULL,
  `note` TINYTEXT NULL,
  `document_status` ENUM('active', 'finished') NULL,
  `is_deleted` TINYINT(1) NULL DEFAULT 0,
  `created_by` BIGINT(20) UNSIGNED NOT NULL,
  `updated_by` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_workorder_receiving_1_idx` (`purchase_order_id` ASC),
  INDEX `fk_workorder_receiving_2_idx` (`created_by` ASC),
  INDEX `fk_workorder_receiving_3_idx` (`updated_by` ASC),
  CONSTRAINT `fk_workorder_receiving_1`
    FOREIGN KEY (`purchase_order_id`)
    REFERENCES `purchase_order` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_workorder_receiving_2`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_workorder_receiving_3`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
