SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `recap_sales_item` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `recap_sales_id` BIGINT(20) UNSIGNED NOT NULL,
  `sales_order_id` BIGINT(20) UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_recap_sales_item_1_idx` (`recap_sales_id` ASC),
  INDEX `fk_recap_sales_item_2_idx` (`sales_order_id` ASC),
  CONSTRAINT `fk_recap_sales_item_1`
    FOREIGN KEY (`recap_sales_id`)
    REFERENCES `recap_sales` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_recap_sales_item_2`
    FOREIGN KEY (`sales_order_id`)
    REFERENCES `sales_order` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;