CREATE TABLE IF NOT EXISTS analytics_metrics (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    channel_id UUID,
    user_id UUID,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    metadata JSONB,
    version BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_analytics_metrics_name ON analytics_metrics(name);
CREATE INDEX IF NOT EXISTS idx_analytics_metrics_timestamp ON analytics_metrics(timestamp);
CREATE INDEX IF NOT EXISTS idx_analytics_metrics_channel_id ON analytics_metrics(channel_id);
CREATE INDEX IF NOT EXISTS idx_analytics_metrics_user_id ON analytics_metrics(user_id);
CREATE INDEX IF NOT EXISTS idx_analytics_metrics_name_timestamp ON analytics_metrics(name, timestamp);
