CREATE TABLE channel_invites (
    id UUID PRIMARY KEY,
    channel_id UUID NOT NULL REFERENCES channels (id) ON DELETE CASCADE,
    created_by_user_id UUID NOT NULL,
    invite_code VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    max_uses INTEGER DEFAULT NULL,
    current_uses INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT 'now()',
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_channel_invites_channel_id ON channel_invites (channel_id);
CREATE INDEX idx_channel_invites_invite_code ON channel_invites (invite_code);
CREATE INDEX idx_channel_invites_expires_at ON channel_invites (expires_at);
CREATE INDEX idx_channel_invites_is_active ON channel_invites (is_active);
