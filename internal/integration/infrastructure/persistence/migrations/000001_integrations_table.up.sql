CREATE TABLE IF NOT EXISTS integrations (
    id UUID PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    creator_user_id VARCHAR(255) NOT NULL,
    api_token_hash VARCHAR(255) NOT NULL,
    token_lookup_hash VARCHAR(255) NOT NULL UNIQUE,
    target_channel_ids TEXT[] NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_integrations_creator ON integrations(creator_user_id);
CREATE INDEX IF NOT EXISTS idx_integrations_lookup_hash ON integrations(token_lookup_hash);
