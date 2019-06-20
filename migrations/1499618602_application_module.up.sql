SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `application_module` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_module_id` BIGINT(20) UNSIGNED NULL DEFAULT NULL COMMENT 'parent of module',
  `module_name` VARCHAR(255) NULL DEFAULT NULL,
  `alias` VARCHAR(255) NOT NULL COMMENT 'unique, lower case, no space',
  `note` TINYTEXT NULL DEFAULT NULL,
  `is_active` TINYINT(1) NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  INDEX `fk_application_module_1_idx` (`parent_module_id` ASC),
  CONSTRAINT `fk_application_module_1`
    FOREIGN KEY (`parent_module_id`)
    REFERENCES `application_module` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB
COMMENT = 'List of application module';
