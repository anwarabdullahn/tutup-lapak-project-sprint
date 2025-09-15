-- Create purchases table
CREATE TABLE IF NOT EXISTS purchases (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    payment_proof_ids TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create purchase_items table
CREATE TABLE IF NOT EXISTS purchase_items (
    id UUID PRIMARY KEY,
    purchase_id UUID NOT NULL REFERENCES purchases(id) ON DELETE CASCADE,
    product_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    qty INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    sku VARCHAR(255),
    file_id VARCHAR(255),
    file_uri TEXT,
    file_thumbnail_uri TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create purchase_senders table
CREATE TABLE IF NOT EXISTS purchase_senders (
    id UUID PRIMARY KEY,
    purchase_id UUID NOT NULL REFERENCES purchases(id) ON DELETE CASCADE,
    sender_name VARCHAR(255) NOT NULL,
    sender_contact_type VARCHAR(50) NOT NULL,
    sender_contact_detail VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_purchases_user_id ON purchases(user_id);
CREATE INDEX IF NOT EXISTS idx_purchases_created_at ON purchases(created_at);
CREATE INDEX IF NOT EXISTS idx_purchase_items_purchase_id ON purchase_items(purchase_id);
CREATE INDEX IF NOT EXISTS idx_purchase_items_product_id ON purchase_items(product_id);
CREATE INDEX IF NOT EXISTS idx_purchase_senders_purchase_id ON purchase_senders(purchase_id);

-- Add comments for documentation
COMMENT ON TABLE purchases IS 'Stores purchase orders with user information and payment proof';
COMMENT ON TABLE purchase_items IS 'Stores individual items within a purchase order';
COMMENT ON TABLE purchase_senders IS 'Stores sender contact information for purchase orders';

COMMENT ON COLUMN purchases.user_id IS 'ID of the user who made the purchase';
COMMENT ON COLUMN purchases.payment_proof_ids IS 'JSON array of file IDs for payment proof images';
COMMENT ON COLUMN purchase_items.product_id IS 'ID of the product being purchased';
COMMENT ON COLUMN purchase_items.qty IS 'Quantity of the product purchased';
COMMENT ON COLUMN purchase_items.price IS 'Price per unit of the product';
COMMENT ON COLUMN purchase_senders.sender_contact_type IS 'Type of contact: email or phone';
