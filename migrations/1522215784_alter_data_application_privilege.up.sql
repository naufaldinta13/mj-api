SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO `application_privilege` (`application_module_id`, `usergroup_id`)
SELECT `application_module_id`, '2'
FROM `application_privilege`
WHERE `application_module_id` IN (135,136,137,37) AND ((`id` = '135') OR (`id` = '136') OR (`id` = '137'));

INSERT INTO `application_privilege` (`application_module_id`, `usergroup_id`)
SELECT `application_module_id`, '3'
FROM `application_privilege`
WHERE `application_module_id` IN (135,136,137,37) AND ((`id` = '485') OR (`id` = '486') OR (`id` = '487'));

INSERT INTO `application_privilege` (`application_module_id`, `usergroup_id`)
SELECT `application_module_id`, '4'
FROM `application_privilege`
WHERE `application_module_id` IN (135,136,137,37) AND ((`id` = '488') OR (`id` = '489') OR (`id` = '490'));


INSERT INTO `application_privilege` (`application_module_id`, `usergroup_id`)
SELECT '37', '2'
FROM `application_privilege`
WHERE `application_module_id` IN (37) AND ((`id` = '37'));

INSERT INTO `application_privilege` (`application_module_id`, `usergroup_id`)
SELECT '37', '3'
FROM `application_privilege`
WHERE `application_module_id` IN (37) AND ((`id` = '494'));

INSERT INTO `application_privilege` (`application_module_id`, `usergroup_id`)
SELECT '37', '4'
FROM `application_privilege`
WHERE `application_module_id` IN (37) AND ((`id` = '495'));