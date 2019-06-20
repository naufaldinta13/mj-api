SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `invoice_receipt_return` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `invoice_receipt_id` BIGINT(20) UNSIGNED NOT NULL,
  `sales_return_id` BIGINT(20) UNSIGNED NOT NULL,
  `subtotal` DECIMAL(20,0) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `fk_invoice_receipt_item_idx` (`invoice_receipt_id` ASC),
  INDEX `fk_invoice_receipt_item_20_idx` (`sales_return_id` ASC),
  CONSTRAINT `fk_invoice_receipt_item_10`
    FOREIGN KEY (`invoice_receipt_id`)
    REFERENCES `invoice_receipt` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_invoice_receipt_item_20`
    FOREIGN KEY (`sales_return_id`)
    REFERENCES `sales_return` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
