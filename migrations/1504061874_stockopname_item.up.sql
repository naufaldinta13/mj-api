SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `stockopname_item` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `stockopname_id` BIGINT(20) UNSIGNED NOT NULL,
  `item_variant_stock_id` BIGINT(20) UNSIGNED NOT NULL,
  `quantity` FLOAT NOT NULL,
  `note` TINYTEXT NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_stockopname_item_q_idx` (`stockopname_id` ASC),
  INDEX `fk_stockopname_item_2_idx` (`item_variant_stock_id` ASC),
  CONSTRAINT `fk_stockopname_item_1`
    FOREIGN KEY (`stockopname_id`)
    REFERENCES `stockopname` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_stockopname_item_2`
    FOREIGN KEY (`item_variant_stock_id`)
    REFERENCES `item_variant_stock` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
