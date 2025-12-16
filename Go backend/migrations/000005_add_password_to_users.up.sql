-- Add password column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS password VARCHAR(255) NOT NULL DEFAULT '';

-- Populate existing users with default password "password"
UPDATE users SET password = 'password' WHERE password = '';

