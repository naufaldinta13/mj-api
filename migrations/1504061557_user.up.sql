SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `user` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `usergroup_id` BIGINT(20) UNSIGNED NOT NULL,
  `full_name` VARCHAR(100) NOT NULL,
  `username` VARCHAR(100) NOT NULL,
  `password` VARCHAR(145) NOT NULL,
  `remember_token` VARCHAR(128) NULL DEFAULT NULL,
  `is_active` TINYINT(1) NULL DEFAULT 1,
  `last_login` TIMESTAMP NULL DEFAULT NULL,
  `created_by` BIGINT(20) UNSIGNED NULL,
  `updated_by` BIGINT(20) UNSIGNED NULL DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_user_1_idx` (`usergroup_id` ASC),
  INDEX `fk_user_2_idx` (`created_by` ASC),
  INDEX `fk_user_3_idx` (`updated_by` ASC),
  CONSTRAINT `fk_user_1`
    FOREIGN KEY (`usergroup_id`)
    REFERENCES `usergroup` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_user_2`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_user_3`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
