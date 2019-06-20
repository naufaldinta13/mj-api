SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `workorder_receiving_item` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `workorder_receiving_id` BIGINT(20) UNSIGNED NOT NULL,
  `purchase_order_item_id` BIGINT(20) UNSIGNED NOT NULL,
  `quantity` FLOAT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_workorder_receiving_item_1_idx` (`workorder_receiving_id` ASC),
  INDEX `fk_workorder_receiving_item_2_idx` (`purchase_order_item_id` ASC),
  CONSTRAINT `fk_workorder_receiving_item_1`
    FOREIGN KEY (`workorder_receiving_id`)
    REFERENCES `workorder_receiving` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_workorder_receiving_item_2`
    FOREIGN KEY (`purchase_order_item_id`)
    REFERENCES `purchase_order_item` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;