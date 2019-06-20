SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `pricing_type`
CHANGE `rule_type` `rule_type` enum('increment','decrement','none') NULL DEFAULT 'none' AFTER `note`;
ALTER TABLE `pricing_type`
CHANGE `parent_type` `parent_type_id` bigint(20) unsigned NULL COMMENT 'Base pricing type' AFTER `id`;
