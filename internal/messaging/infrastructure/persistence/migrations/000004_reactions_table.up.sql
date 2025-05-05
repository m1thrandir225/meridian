CREATE TABLE reactions (
    id UUID PRIMARY KEY,
    message_id UUID NOT NULL REFERENCES messages (id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    reaction_type VARCHAR(100) NOT NULL, -- Emoji code or name
    created_at TIMESTAMPTZ NOT NULL,
    UNIQUE (message_id, user_id, reaction_type)
);

CREATE INDEX idx_reactions_message_id ON reactions (message_id);
CREATE INDEX idx_reactions_user_id ON reactions (user_id);
