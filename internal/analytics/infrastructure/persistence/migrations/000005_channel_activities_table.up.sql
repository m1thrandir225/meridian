CREATE TABLE IF NOT EXISTS channel_activities (
    id UUID PRIMARY KEY,
    channel_id UUID NOT NULL UNIQUE,
    messages_count BIGINT NOT NULL DEFAULT 0,
    members_count BIGINT NOT NULL DEFAULT 0,
    last_message_at TIMESTAMP WITH TIME ZONE NOT NULL,
    activity_score DOUBLE PRECISION NOT NULL DEFAULT 0,
    version BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_channel_activities_channel_id ON channel_activities(channel_id);
CREATE INDEX IF NOT EXISTS idx_channel_activities_last_message_at ON channel_activities(last_message_at);
CREATE INDEX IF NOT EXISTS idx_channel_activities_activity_score ON channel_activities(activity_score);
