SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `item_variant_price` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `item_variant_id` BIGINT(20) UNSIGNED NOT NULL,
  `pricing_type_id` BIGINT(20) UNSIGNED NOT NULL,
  `unit_price` DECIMAL(20,0) NOT NULL,
  `note` TINYTEXT NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_item_variant_price_1_idx` (`item_variant_id` ASC),
  INDEX `fk_item_variant_price_2_idx` (`pricing_type_id` ASC),
  CONSTRAINT `fk_item_variant_price_1`
    FOREIGN KEY (`item_variant_id`)
    REFERENCES `item_variant` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_item_variant_price_2`
    FOREIGN KEY (`pricing_type_id`)
    REFERENCES `pricing_type` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB
COMMENT = 'Relation between item variant with pricing type';