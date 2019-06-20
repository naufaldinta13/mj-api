SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `sales_return_item` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `sales_return_id` BIGINT(20) UNSIGNED NOT NULL,
  `sales_order_item_id` BIGINT(20) UNSIGNED NOT NULL,
  `quantity` FLOAT UNSIGNED NOT NULL,
  `note` TINYTEXT NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_sales_return_item_1_idx` (`sales_return_id` ASC),
  INDEX `fk_sales_return_item_2_idx` (`sales_order_item_id` ASC),
  CONSTRAINT `fk_sales_return_item_1`
    FOREIGN KEY (`sales_return_id`)
    REFERENCES `sales_return` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_return_item_2`
    FOREIGN KEY (`sales_order_item_id`)
    REFERENCES `sales_order_item` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
