SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `workorder_shipment_item` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `workorder_shipment_id` BIGINT(20) UNSIGNED NOT NULL,
  `workorder_fulfillment_id` BIGINT(20) UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_workorder_shipment_item_1_idx` (`workorder_shipment_id` ASC),
  INDEX `fk_workorder_shipment_item_2_idx` (`workorder_fulfillment_id` ASC),
  CONSTRAINT `fk_workorder_shipment_item_1`
    FOREIGN KEY (`workorder_shipment_id`)
    REFERENCES `workorder_shipment` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_workorder_shipment_item_2`
    FOREIGN KEY (`workorder_fulfillment_id`)
    REFERENCES `workorder_fulfillment` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
