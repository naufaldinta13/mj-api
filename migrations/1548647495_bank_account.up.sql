SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `bank_account` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `bank_name` VARCHAR(20) NOT NULL,
  `bank_number` VARCHAR(20) NOT NULL,
  `is_default` TINYINT(1) NULL DEFAULT '0',
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


INSERT INTO `application_module` (`id`,`parent_module_id`, `module_name`, `alias`, `is_active`) VALUES ('161','7', 'Bank Report', 'report_bank', '1');

INSERT INTO `application_privilege` (`application_module_id`, `usergroup_id`) VALUES ('161', '1');
INSERT INTO `application_privilege` (`application_module_id`, `usergroup_id`) VALUES ('161', '2');
