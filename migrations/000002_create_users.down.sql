-- Drop trigger before dropping the table
DROP TRIGGER IF EXISTS set_updated_at ON users;

-- Drop indexes explicitly
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_phone;
DROP INDEX IF EXISTS idx_users_town;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_verification_token;
DROP INDEX IF EXISTS idx_users_password_reset_token;
DROP INDEX IF EXISTS idx_users_metadata;

-- Drop users table
DROP TABLE IF EXISTS users CASCADE;