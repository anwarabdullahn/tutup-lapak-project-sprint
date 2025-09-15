-- Drop indexes first
DROP INDEX IF EXISTS idx_purchase_senders_purchase_id;
DROP INDEX IF EXISTS idx_purchase_items_product_id;
DROP INDEX IF EXISTS idx_purchase_items_purchase_id;
DROP INDEX IF EXISTS idx_purchases_created_at;
DROP INDEX IF EXISTS idx_purchases_user_id;

-- Drop tables in reverse order (due to foreign key constraints)
DROP TABLE IF EXISTS purchase_senders;
DROP TABLE IF EXISTS purchase_items;
DROP TABLE IF EXISTS purchases;
