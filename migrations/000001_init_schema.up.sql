-- Enable useful extensions
CREATE EXTENSION IF NOT EXISTS citext;

-- Create enum types for shared use
CREATE TYPE user_role AS ENUM ('volunteer', 'organizer', 'admin');
CREATE TYPE user_status AS ENUM ('pending', 'active', 'suspended', 'deleted');

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;