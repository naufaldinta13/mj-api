SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `purchase_order_item` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `purchase_order_id` BIGINT(20) UNSIGNED NOT NULL,
  `item_variant_id` BIGINT(20) UNSIGNED NOT NULL,
  `quantity` FLOAT UNSIGNED NOT NULL,
  `unit_price` DECIMAL(20,0) UNSIGNED NULL DEFAULT 0,
  `discount` FLOAT NULL DEFAULT 0,
  `subtotal` DECIMAL(20,0) UNSIGNED NOT NULL,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_purchase_order_item_1_idx` (`purchase_order_id` ASC),
  INDEX `fk_purchase_order_item_2_idx` (`item_variant_id` ASC),
  CONSTRAINT `fk_purchase_order_item_1`
    FOREIGN KEY (`purchase_order_id`)
    REFERENCES `purchase_order` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_purchase_order_item_2`
    FOREIGN KEY (`item_variant_id`)
    REFERENCES `item_variant` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = big5;
