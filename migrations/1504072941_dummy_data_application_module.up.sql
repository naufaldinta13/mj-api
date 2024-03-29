SET FOREIGN_KEY_CHECKS = 0;
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (1,NULL,'Dashboard','dashboard','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (2,NULL,'Sales','sales','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (3,NULL,'Purchase','purchase','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (4,NULL,'Inventory','inventory','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (5,NULL,'Workorder','workorder','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (6,NULL,'Finance','finance','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (7,NULL,'Report','report','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (8,NULL,'Partnership','partnership','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (9,NULL,'Setting','setting','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (10,1,'Display Dashboard','dashboard_display','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (11,2,'Sales Order','sales_order','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (12,2,'Sales Invoice','sales_invoice','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (13,2,'Sales Return','sales_return','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (14,3,'Purchase Order','purchase_order','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (15,3,'Purchase Invoice','purchase_invoice','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (16,3,'Purchase Return','purchase_return','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (17,4,'Item','inventory_item','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (18,4,'Category','inventory_category','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (19,4,'Unit of Measurement','inventory_uom','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (20,4,'Direct Placement','inventory_direct_placement','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (21,4,'Stock Opname','inventory_stock_opname','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (22,5,'Fulfillment','workorder_fulfillment','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (23,5,'Shipment','workorder_shipment','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (24,5,'Receiving','workorder_receiving','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (25,6,'Revenue','finance_revenue','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (26,6,'Expense','finance_expense','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (27,6,'Invoice Receipt','finance_invoice_receipt','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (28,7,'Sales Report','report_sales','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (29,7,'Purchase Report','report_purchase','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (30,7,'Sales Item Report','report_sales_item','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (31,7,'Purchase Item Report','report_purchase_item','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (32,7,'Debt Report','report_debt','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (33,7,'Credit Report','report_credit','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (34,7,'Customer Sales Report','report_customer_sales','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (35,9,'Application Setting','setting_application','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (36,9,'Pricing Type Setting','setting_pricing_type','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (37,9,'User Setting','setting_user','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (38,11,'Create Sales Order','sales_order_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (39,11,'Read Sales Order','sales_order_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (40,11,'Update Sales Order','sales_order_update','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (41,11,'Request Cancel Sales Order','sales_order_request_cancel','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (42,11,'Cancel Sales Order','sales_order_cancel','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (43,11,'Show Sales Order','sales_order_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (44,12,'Create Sales Invoice','sales_invoice_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (45,12,'Read Sales Invoice','sales_invoice_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (46,12,'Update Sales Invoice','sales_invoice_update','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (47,12,'Show Sales Invoice','sales_invoice_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (48,13,'Create Sales Return','sales_return_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (49,13,'Read Sales Return','sales_return_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (50,13,'Update Sales Return','sales_return_update','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (51,13,'Cancel Sales Return','sales_return_cancel','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (52,13,'Show Sales Return','sales_return_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (53,14,'Create Purchase Order','purchase_order_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (54,14,'Read Purchase Order','purchase_order_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (55,14,'Update Purchase Order','purchase_order_update','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (56,14,'Cancel Purchase Order','purchase_order_cancel','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (57,14,'Show Purchase Order','purchase_order_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (58,15,'Create Purchase Invoice','purchase_invoice_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (59,15,'Read Purchase Invoice','purchase_invoice_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (60,15,'Update Purchase Invoice','purchase_invoice_update','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (61,15,'Show Purchase Invoice','purchase_invoice_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (62,16,'Create Purchase Return','purchase_return_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (63,16,'Read Purchase Return','purchase_return_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (64,16,'Update Purchase Return','purchase_return_update','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (65,16,'Cancel Purchase Return','purchase_return_cancel','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (66,16,'Show Purchase Return','purchase_return_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (67,17,'Create Item','item_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (68,17,'Read Item','item_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (69,17,'Update Item','item_update','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (70,17,'Show Item','item_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (71,17,'Archive Item','item_archive','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (72,17,'Delete Item','item_delete','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (73,18,'Create Category','category_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (74,18,'Read Category','category_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (75,18,'Update Category','category_update','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (76,18,'Show Category','category_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (77,18,'Delete Category','category_delete','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (78,19,'Create Unit of Measurement','uom_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (79,19,'Read Unit of Measurement','uom_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (80,19,'Update Unit of Measurement','uom_update','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (81,19,'Show Unit of Measurement','uom_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (82,19,'Delete Unit of Measurement','uom_delete','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (83,20,'Create Direct Placement','direct_placement_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (84,20,'Read Direct Placement','direct_placement_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (85,20,'Show Direct Placement','direct_placement_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (86,21,'Create Stock Opname','stock_opname_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (87,21,'Read Stock Opname','stock_opname_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (88,21,'Show Stock Opname','stock_opname_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (89,22,'Create Fulfillment','fulfillment_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (90,22,'Read Fulfillment','fulfillment_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (91,22,'Update Fulfillment','fulfillment_update','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (92,22,'Show Fulfillment','fulfillment_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (93,22,'Approve Fulfillment','fulfillment_approve','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (94,23,'Create Shipment','shipment_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (95,23,'Read Shipment','shipment_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (96,23,'Update Shipment','shipment_update','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (97,23,'Show Shipment','shipment_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (98,23,'Approve Shipment','shipment_approve','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (99,23,'Cancel Shipment','shipment_cancel','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (100,24,'Create Receiving','receiving_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (101,24,'Read Receiving','receiving_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (102,24,'Show Receiving','receiving_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (103,25,'Create Revenue','revenue_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (104,25,'Read Revenue','revenue_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (105,25,'Update Revenue','revenue_update','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (106,25,'Show Revenue','revenue_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (107,25,'Approve Revenue','revenue_approve','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (108,26,'Create Expense','expense_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (109,26,'Read Expense','expense_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (110,26,'Update Expense','expense_update','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (111,26,'Show Expense','expense_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (112,26,'Approve Expense','expense_approve','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (113,27,'Create Invoice Receipt','invoice_receipt_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (114,27,'Read Invoice Receipt','invoice_receipt_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (115,27,'Show Invoice Receipt','invoice_receipt_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (116,27,'Pay Invoice Receipt','invoice_receipt_pay','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (117,28,'Display Sales Report','sales_report_display','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (118,29,'Display Purchase Report','purchase_report_display','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (119,30,'Display Sales Item Report','sales_item_report_display','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (120,31,'Display Purchase Item Report','purchase_item_report_display','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (121,32,'Display Debt Report','debt_report_display','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (122,33,'Display Credit Report','credit_report_display','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (123,34,'Create Customer Sales Report','customer_sales_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (124,34,'Read Customer Sales Report','customer_sales_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (125,34,'Show Customer Sales Report','customer_sales_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (126,34,'Delete Customer Sales Report','customer_sales_delete','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (127,35,'Read Application Setting','application_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (128,35,'Update Application Setting','application_update','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (129,35,'Show Application Setting','application_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (130,36,'Create Pricing Type','pricing_type_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (131,36,'Read Pricing Type','pricing_type_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (132,36,'Update Pricing Type','pricing_type_update','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (133,36,'Show Pricing Type','pricing_type_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (134,37,'Create User','user_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (135,37,'Read User','user_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (136,37,'Change Password User','user_change_password','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (137,37,'Show User','user_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (138,8,'Create Partnership','partnership_create','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (139,8,'Read Partnership','partnership_read','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (140,8,'Update Partnership','partnership_update','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (141,8,'Archive Partnership','partnership_archive','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (142,8,'Show Partnership','partnership_show','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (143,8,'Delete Partnership','partnership_delete','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (144,14,'Print Purchase Order','purchase_order_print','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (145,15,'Print Purchase Invoice','purchase_invoice_print','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (146,16,'Print Purchase Return','purchase_return_print','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (147,27,'Print Invoice Receipt','invoice_receipt_print','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (148,11,'Print Sales Order','sales_order_print','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (149,12,'Print Sales Invoice','sales_invoice_print','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (150,23,'Print Shipment','shipment_print','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (151,37,'Inactive User','user_inactive','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (152,28,'Print Sales Report','sales_report_print','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (153,29,'Print Purchase Report','purchase_report_print','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (154,30,'Print Sales Item Report','sales_item_report_print','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (155,31,'Print Purchase Item Report','purchase_item_report_print','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (156,32,'Print Debt Report','debt_report_print','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (157,33,'Print Credit Report','credit_report_print','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (158,34,'Print Customer Sales Report','customer_sales_print','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (159,22,'Print Fulfillment','fulfillment_print','',1);
INSERT INTO `application_module` (`id`,`parent_module_id`,`module_name`,`alias`,`note`,`is_active`) VALUES (160,11,'Reject Cancel Sales Order','sales_order_reject_cancel','',1);