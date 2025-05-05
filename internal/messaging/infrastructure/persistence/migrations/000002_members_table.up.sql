CREATE TABLE members (
    channel_id UUID NOT NULL REFERENCES channels (id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT 'now()',
    last_read TIMESTAMPTZ
);

CREATE INDEX idx_members_user_id ON members (user_id);
