SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `application_privilege` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `application_module_id` BIGINT(20) UNSIGNED NOT NULL COMMENT 'FK to application_module table',
  `usergroup_id` BIGINT(20) UNSIGNED NOT NULL COMMENT 'FK to usergroup table\n',
  PRIMARY KEY (`id`),
  INDEX `fk_application_privilege_1_idx` (`application_module_id` ASC),
  INDEX `fk_application_privilege_2_idx` (`usergroup_id` ASC),
  CONSTRAINT `fk_application_privilege_1`
    FOREIGN KEY (`application_module_id`)
    REFERENCES `application_module` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_application_privilege_2`
    FOREIGN KEY (`usergroup_id`)
    REFERENCES `usergroup` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB
COMMENT = 'Relation between usergroup and module';
