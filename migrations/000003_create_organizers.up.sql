CREATE TYPE organization_type_enum AS ENUM ('nonprofit', 'corporate', 'community_group', 'government', 'school');

CREATE TABLE IF NOT EXISTS organizers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,

    -- Organization Details
    organization_name TEXT NOT NULL,
    description TEXT,
    website_url TEXT,

    -- Contact Information (organization-specific)
    contact_email CITEXT,
    contact_phone TEXT,

    -- Location Details
    operating_towns TEXT[] NOT NULL DEFAULT '{}', -- Multiple towns they operate in
    headquarters_town TEXT, -- Main base of operations

    -- Organization Type/Category
    organization_type organization_type_enum NOT NULL DEFAULT 'community_group', -- e.g., 'nonprofit', 'corporate', community_group', 'government'

    -- Verification/Trust Signals
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    verification_documents TEXT[] NOT NULL DEFAULT '{}', -- Array of document URLs

    -- Capacity & Impact Tracking
    years_operating SMALLINT,
    previous_cleanup_events INTEGER DEFAULT 0,
    avg_volunteer_turnout INTEGER, -- Average volunteers per event

    -- Social Proof
    social_media_links JSONB DEFAULT '{}', -- {facebook: url, twitter: url, etc.}

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for better query performance
CREATE INDEX idx_organizers_operating_towns ON organizers USING GIN(operating_towns);
CREATE INDEX idx_organizers_organization_type ON organizers(organization_type);
CREATE INDEX idx_organizers_verified ON organizers(verified) WHERE verified = true;
CREATE INDEX idx_organizers_headquarters_town ON organizers(headquarters_town);

CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON organizers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE organizers IS 'Organizer profiles for waste management organizations, community groups, and corporate sponsors';
COMMENT ON COLUMN organizers.operating_towns IS 'Array of towns where this organization operates cleanup events';
COMMENT ON COLUMN organizers.organization_type IS 'Type of organization: nonprofit, corporate, community_group, government, school';
COMMENT ON COLUMN organizers.verified IS 'Whether the organization has been verified by CleanTown admins';