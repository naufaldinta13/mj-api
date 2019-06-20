SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `stockopname` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(120) NOT NULL,
  `recognition_date` DATE NULL DEFAULT NULL,
  `note` TINYTEXT NULL DEFAULT NULL,
  `created_by` BIGINT(20) UNSIGNED NOT NULL,
  `updated_by` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_stockopname_1_idx` (`created_by` ASC),
  INDEX `fk_stockopname_2_idx` (`updated_by` ASC),
  CONSTRAINT `fk_stockopname_1`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_stockopname_2`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;