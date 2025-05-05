CREATE TABLE channels (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    topic TEXT NOT NULL,
    creator_user_id UUID NOT NULL,
    creation_time TIMESTAMPTZ NOT NULL,
    last_message_time TIMESTAMPTZ NOT NULL,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    version BIGINT NOT NULL DEFAULT 1
);

CREATE INDEX idx_channels_creator_user_id ON channels (creator_user_id);
CREATE INDEX idx_channels_is_archived ON channels (is_archived);

