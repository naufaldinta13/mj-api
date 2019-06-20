SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS`application_menu` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_menu_id` BIGINT(20) UNSIGNED NULL COMMENT 'parent of menu',
  `application_module_id` BIGINT(20) UNSIGNED NOT NULL COMMENT 'FK to application_module table',
  `menu_name` VARCHAR(255) NOT NULL,
  `icon` VARCHAR(255) NULL DEFAULT NULL,
  `route` VARCHAR(255) NOT NULL COMMENT 'Route alias',
  `order` INT(10) UNSIGNED NOT NULL COMMENT 'sequence of menu',
  `is_active` TINYINT(1) NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  INDEX `fk_application_menu_1_idx` (`parent_menu_id` ASC),
  INDEX `fk_application_menu_2_idx` (`application_module_id` ASC),
  CONSTRAINT `fk_application_menu_1`
    FOREIGN KEY (`parent_menu_id`)
    REFERENCES `application_menu` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_application_menu_2`
    FOREIGN KEY (`application_module_id`)
    REFERENCES `application_module` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB
COMMENT = 'List of application menu';
