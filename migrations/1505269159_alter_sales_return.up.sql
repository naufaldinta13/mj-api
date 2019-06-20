SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `sales_return`
CHANGE  `document_status`  `document_status` ENUM('new','active','finished','cancelled') NULL DEFAULT 'new';