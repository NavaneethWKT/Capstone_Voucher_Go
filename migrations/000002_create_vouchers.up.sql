-- Create vouchers table
CREATE TABLE IF NOT EXISTS vouchers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    valid_from TIMESTAMP NOT NULL,
    valid_to TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_valid_dates CHECK (valid_to > valid_from)
);

-- Create indexes for faster searches
CREATE INDEX IF NOT EXISTS idx_vouchers_category ON vouchers(category);
CREATE INDEX IF NOT EXISTS idx_vouchers_price ON vouchers(price);
CREATE INDEX IF NOT EXISTS idx_vouchers_validity ON vouchers(valid_from, valid_to);

