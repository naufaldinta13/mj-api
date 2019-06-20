SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `workorder_shipment` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(120) NOT NULL,
  `priority` ENUM('routine', 'rush', 'emergency') NULL DEFAULT 'routine',
  `truck_number` VARCHAR(45) NULL DEFAULT NULL,
  `document_status` ENUM('active', 'finished') NULL DEFAULT 'active',
  `is_deleted` TINYINT(1) NULL DEFAULT 0,
  `created_by` BIGINT(20) UNSIGNED NOT NULL,
  `updated_by` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_workorder_shipment_1_idx` (`created_by` ASC),
  INDEX `fk_workorder_shipment_2_idx` (`updated_by` ASC),
  CONSTRAINT `fk_workorder_shipment_1`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_workorder_shipment_2`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
