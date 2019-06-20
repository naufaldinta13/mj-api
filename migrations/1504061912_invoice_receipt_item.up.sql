SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `invoice_receipt_item` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `invoice_receipt_id` BIGINT(20) UNSIGNED NOT NULL,
  `sales_invoice_id` BIGINT(20) UNSIGNED NOT NULL,
  `subtotal` DECIMAL(20,0) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `fk_invoice_receipt_item_idx` (`invoice_receipt_id` ASC),
  INDEX `fk_invoice_receipt_item_2_idx` (`sales_invoice_id` ASC),
  CONSTRAINT `fk_invoice_receipt_item_1`
    FOREIGN KEY (`invoice_receipt_id`)
    REFERENCES `invoice_receipt` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_invoice_receipt_item_2`
    FOREIGN KEY (`sales_invoice_id`)
    REFERENCES `sales_invoice` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
