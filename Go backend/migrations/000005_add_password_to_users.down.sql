-- Remove password column from users table
ALTER TABLE users DROP COLUMN IF EXISTS password;

