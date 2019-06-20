SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `recap_sales` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `recognition_date` DATE NOT NULL,
  `partnership_id` BIGINT(20) UNSIGNED NOT NULL,
  `code` VARCHAR(45) NOT NULL,
  `total_amount` DECIMAL(20,0) UNSIGNED NOT NULL,
  `created_by` BIGINT(20) UNSIGNED NOT NULL,
  `created_at` TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_recap_sales_1_idx` (`created_by` ASC),
  INDEX `fk_recap_sales_2_idx` (`partnership_id` ASC),
  CONSTRAINT `fk_recap_sales_1`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_recap_sales_2`
    FOREIGN KEY (`partnership_id`)
    REFERENCES `partnership` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;