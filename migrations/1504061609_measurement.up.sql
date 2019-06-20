SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `measurement` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `measurement_name` VARCHAR(45) NOT NULL COMMENT 'name of UOM',
  `note` TINYTEXT NULL DEFAULT NULL,
  `is_deleted` TINYINT(1) NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
COMMENT = 'Unit of measurement';