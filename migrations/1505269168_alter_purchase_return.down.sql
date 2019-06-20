ALTER TABLE `purchase_return`
CHANGE  `document_status`   `document_status` ENUM('new', 'active', 'finished') NULL DEFAULT 'new';