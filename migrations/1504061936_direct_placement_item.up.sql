SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE IF NOT EXISTS `direct_placement_item` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `direct_placment_id` BIGINT(20) UNSIGNED NOT NULL,
  `item_variant_id` BIGINT(20) UNSIGNED NOT NULL,
  `quantity` FLOAT NULL DEFAULT 0,
  `unti_price` DECIMAL(20,0) UNSIGNED NOT NULL DEFAULT 0,
  `total_price` DECIMAL(20,0) UNSIGNED NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `fk_direct_placement_item_2_idx` (`item_variant_id` ASC),
  CONSTRAINT `fk_direct_placement_item_1`
    FOREIGN KEY (`direct_placment_id`)
    REFERENCES `direct_placement` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_direct_placement_item_2`
    FOREIGN KEY (`item_variant_id`)
    REFERENCES `item_variant` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
