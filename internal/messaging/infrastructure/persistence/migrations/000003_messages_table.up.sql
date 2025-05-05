CREATE TABLE messages (
  id UUID PRIMARY KEY,
  channel_id UUID NOT NULL REFERENCES channels (id) ON DELETE CASCADE,
  sender_user_id UUID,
  integration_id UUID,
  content_text TEXT NOT NULL,
  content_mentions UUID[],
  content_link TEXT[],
  content_formatted BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL,
  parent_message_id UUID REFERENCES messages(id) ON DELETE SET NULL
);

CREATE INDEX idx_messages_channel_id_timestamp ON messages (channel_id, created_at DESC);
CREATE INDEX idx_messages_sender_user_id ON messages (sender_user_id);
CREATE INDEX idx_messages_integration_id ON messages (integration_id);
CREATE INDEX idx_messages_parent_message_id ON messages (parent_message_id);
