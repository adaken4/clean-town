-- Users table
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    
    -- Basic Information
    name TEXT NOT NULL,
    email CITEXT UNIQUE NOT NULL,
    phone TEXT,

    -- Authentication
    password_hash BYTEA NOT NULL,

    -- Profile
    town TEXT,
    bio TEXT,
    age SMALLINT,
    avatar_url TEXT,
    skills TEXT[], -- Array of skills (e.g., '{recycling, community_org, photography}')

    -- Authorization
    role user_role NOT NULL DEFAULT 'volunteer',
    status user_status NOT NULL DEFAULT 'pending',

    -- Verification
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    phone_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verification_token TEXT,
    verification_token_expiry TIMESTAMP WITH TIME ZONE,

    -- Password Reset
    password_reset_token TEXT,
    password_reset_expiry TIMESTAMP WITH TIME ZONE,

    -- Metadata
    last_login_at TIMESTAMP WITH TIME ZONE,
    login_count INTEGER NOT NULL DEFAULT 0,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Additional metadata (for future extensibility)
    metadata JSONB DEFAULT '{}'::JSONB,
    
    -- Constraints
    CONSTRAINT email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT phone_format CHECK (phone IS NULL OR phone ~* '^\+?[0-9]{10,15}$')
);

-- Create indexes for common queries
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_phone ON users(phone) WHERE deleted_at IS NULL AND phone IS NOT NULL;
CREATE INDEX idx_users_town ON users(town) WHERE deleted_at IS NULL AND town IS NOT NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at DESC);
CREATE INDEX idx_users_verification_token ON users(verification_token) WHERE verification_token IS NOT NULL;
CREATE INDEX idx_users_password_reset_token ON users(password_reset_token) WHERE password_reset_token IS NOT NULL;

-- Create index for metadata JSONB queries (if you plan to query specific JSON fields)
CREATE INDEX idx_users_metadata ON users USING GIN(metadata);

-- Trigger to auto-update updated_at
CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments for documentation
COMMENT ON TABLE users IS 'Core user table for volunteers, organizers, and admins';
COMMENT ON COLUMN users.id IS 'Primary key - internal user ID';
COMMENT ON COLUMN users.email IS 'User email - case insensitive, unique, required for login';
COMMENT ON COLUMN users.phone IS 'User phone number in E.164 format (e.g., +254700123456)';
COMMENT ON COLUMN users.password_hash IS 'Bcrypt hashed password';
COMMENT ON COLUMN users.role IS 'User role: volunteer (default), organizer, or admin';
COMMENT ON COLUMN users.status IS 'Account status: pending activation, active, suspended, or deleted';
COMMENT ON COLUMN users.skills IS 'Array of user skills/interests for better event matching';
COMMENT ON COLUMN users.metadata IS 'Flexible JSONB field for additional user data';
COMMENT ON COLUMN users.deleted_at IS 'Soft delete timestamp - NULL means not deleted';
