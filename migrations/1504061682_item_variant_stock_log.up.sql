SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `item_variant_stock_log` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `item_variant_stock_id` BIGINT(20) UNSIGNED NOT NULL,
  `ref_id` BIGINT(20) UNSIGNED NOT NULL COMMENT 'refer to id of referred document',
  `ref_type` ENUM('workorder_fulfillment', 'workorder_receiving', 'stockopname', 'direct_placement') NULL DEFAULT 'workorder_fulfillment' COMMENT 'referred document type',
  `log_type` ENUM('in', 'out') NULL DEFAULT 'in',
  `quantity` FLOAT UNSIGNED NOT NULL,
  `final_stock` FLOAT UNSIGNED NOT NULL COMMENT 'quantity afterwards',
  PRIMARY KEY (`id`),
  INDEX `fk_item_stock_log_1_idx` (`item_variant_stock_id` ASC),
  CONSTRAINT `fk_item_stock_log_1`
    FOREIGN KEY (`item_variant_stock_id`)
    REFERENCES `item_variant_stock` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB
COMMENT = 'history of item variant stock';
