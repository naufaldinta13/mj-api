SET FOREIGN_KEY_CHECKS = 0;
UPDATE sales_order set is_reported = 1 where id in (SELECT sales_order_id FROM recap_sales_item) AND id > 0;