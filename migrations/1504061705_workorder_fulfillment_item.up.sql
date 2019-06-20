SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `workorder_fulfillment_item` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `workorder_fulfillment_id` BIGINT(20) UNSIGNED NOT NULL,
  `sales_order_item_id` BIGINT(20) UNSIGNED NOT NULL,
  `quantity` FLOAT UNSIGNED NOT NULL,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_workorder_fulfillment_item_1_idx` (`workorder_fulfillment_id` ASC),
  INDEX `fk_workorder_fulfillment_item_2_idx` (`sales_order_item_id` ASC),
  CONSTRAINT `fk_workorder_fulfillment_item_1`
    FOREIGN KEY (`workorder_fulfillment_id`)
    REFERENCES `workorder_fulfillment` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_workorder_fulfillment_item_2`
    FOREIGN KEY (`sales_order_item_id`)
    REFERENCES `sales_order_item` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
