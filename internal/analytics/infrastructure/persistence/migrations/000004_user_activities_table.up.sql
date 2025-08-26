CREATE TABLE IF NOT EXISTS user_activities (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    last_active_at TIMESTAMP WITH TIME ZONE NOT NULL,
    messages_sent BIGINT NOT NULL DEFAULT 0,
    channels_joined BIGINT NOT NULL DEFAULT 0,
    reactions_given BIGINT NOT NULL DEFAULT 0,
    session_duration INTERVAL,
    version BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_activities_user_id ON user_activities(user_id);
CREATE INDEX IF NOT EXISTS idx_user_activities_last_active_at ON user_activities(last_active_at);
