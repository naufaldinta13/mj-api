SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `purchase_return_item` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `purchase_return_id` BIGINT(20) UNSIGNED NOT NULL,
  `purchase_order_item_id` BIGINT(20) UNSIGNED NOT NULL,
  `quantity` FLOAT UNSIGNED NOT NULL,
  `note` TINYTEXT NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_purchase_return_item_1_idx` (`purchase_return_id` ASC),
  INDEX `fk_purchase_return_item_2_idx` (`purchase_order_item_id` ASC),
  CONSTRAINT `fk_purchase_return_item_1`
    FOREIGN KEY (`purchase_return_id`)
    REFERENCES `purchase_return` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_purchase_return_item_2`
    FOREIGN KEY (`purchase_order_item_id`)
    REFERENCES `purchase_order_item` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
