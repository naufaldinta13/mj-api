SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `pricing_type` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_type` BIGINT(20) UNSIGNED NULL DEFAULT NULL COMMENT 'Base pricing type',
  `type_name` VARCHAR(45) NOT NULL,
  `note` TINYTEXT NULL DEFAULT NULL,
  `rule_type` ENUM('increment', 'decrement') NOT NULL DEFAULT 'increment',
  `nominal` DECIMAL(20,0) NOT NULL DEFAULT 0,
  `is_percentage` TINYINT(1) NULL DEFAULT 0,
  `is_default` TINYINT(1) NULL DEFAULT 0 COMMENT 'shown by default',
  PRIMARY KEY (`id`),
  INDEX `fk_pricing_type_1_idx` (`parent_type` ASC),
  CONSTRAINT `fk_pricing_type_1`
    FOREIGN KEY (`parent_type`)
    REFERENCES `pricing_type` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB
COMMENT = 'Type of selling price, e.g. wholesale price, retail price, etc.';
