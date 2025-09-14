-- Insert sample data for testing (optional)
-- This migration can be skipped in production

-- Sample purchase data
INSERT INTO purchases (id, user_id, payment_proof_ids, created_at, updated_at) VALUES
    ('01234567-89ab-cdef-0123-456789abcdef', '01234567-89ab-cdef-0123-456789abcde0', '["proof1.jpg", "proof2.jpg"]', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),
    ('01234567-89ab-cdef-0123-456789abcde1', '01234567-89ab-cdef-0123-456789abcde0', NULL, NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours')
ON CONFLICT (id) DO NOTHING;

-- Sample purchase items
INSERT INTO purchase_items (id, purchase_id, product_id, name, category, qty, price, sku, file_id, file_uri, file_thumbnail_uri, created_at, updated_at) VALUES
    ('01234567-89ab-cdef-0123-456789abcde2', '01234567-89ab-cdef-0123-456789abcdef', 'prod-001', 'Sample Product 1', 'Electronics', 2, 99.99, 'SKU-001', 'file-001', 'https://example.com/file1.jpg', 'https://example.com/thumb1.jpg', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),
    ('01234567-89ab-cdef-0123-456789abcde3', '01234567-89ab-cdef-0123-456789abcdef', 'prod-002', 'Sample Product 2', 'Clothing', 1, 49.99, 'SKU-002', 'file-002', 'https://example.com/file2.jpg', 'https://example.com/thumb2.jpg', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),
    ('01234567-89ab-cdef-0123-456789abcde4', '01234567-89ab-cdef-0123-456789abcde1', 'prod-003', 'Sample Product 3', 'Books', 3, 19.99, 'SKU-003', 'file-003', 'https://example.com/file3.jpg', 'https://example.com/thumb3.jpg', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours')
ON CONFLICT (id) DO NOTHING;

-- Sample purchase senders
INSERT INTO purchase_senders (id, purchase_id, sender_name, sender_contact_type, sender_contact_detail, created_at, updated_at) VALUES
    ('01234567-89ab-cdef-0123-456789abcde5', '01234567-89ab-cdef-0123-456789abcdef', 'John Doe', 'email', 'john.doe@example.com', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),
    ('01234567-89ab-cdef-0123-456789abcde6', '01234567-89ab-cdef-0123-456789abcde1', 'Jane Smith', 'phone', '+1234567890', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours')
ON CONFLICT (id) DO NOTHING;
