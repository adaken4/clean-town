CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,

    -- User Context
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    user_ip INET, -- Capture IP address for security
    user_agent TEXT, -- Browser/device info

    -- Action Details
    action TEXT NOT NULL, -- TODO: change to enum
    resource_type TEXT, -- e.g., 'user', 'organizer', 'reported_site', 'event' TODO: change to enum
    resource_id BIGINT, -- ID of the affected resource
    old_values JSONB DEFAULT '{}'::JSONB, -- State before change (for updates)
    new_values JSONB DEFAULT '{}'::JSONB, -- State after change (for updates)
    
    -- Request Context
    request_id TEXT, -- Unique ID for tracing entire request chain
    endpoint TEXT, -- API endpoint or route
    http_method TEXT, -- GET, POST, PUT, DELETE
    
    -- Status & Outcome
    status_code INTEGER, -- HTTP status code (200, 400, 500, etc.)
    error_message TEXT, -- If action failed
    duration_ms INTEGER, -- How long the action took

    -- Metadata    
    metadata JSONB DEFAULT '{}'::JSONB,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Critical Indexes for Performance
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_user_action ON audit_logs(user_id, action);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_request_id ON audit_logs(request_id);
CREATE INDEX idx_audit_logs_status_code ON audit_logs(status_code) WHERE status_code >= 400;

COMMENT ON TABLE audit_logs IS 'Comprehensive audit trail of user actions, system events, and API requests';
COMMENT ON COLUMN audit_logs.user_ip IS 'IP address of the user for security and geo-tracking';
COMMENT ON COLUMN audit_logs.resource_type IS 'Type of resource being modified (user, organizer, site, etc.)';
COMMENT ON COLUMN audit_logs.old_values IS 'JSON snapshot of resource state before modification';
COMMENT ON COLUMN audit_logs.new_values IS 'JSON snapshot of resource state after modification';
COMMENT ON COLUMN audit_logs.request_id IS 'Unique identifier for tracing requests across microservices';
COMMENT ON COLUMN audit_logs.duration_ms IS 'Request duration in milliseconds for performance monitoring';