CREATE TABLE IF NOT EXISTS hourly_metrics (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    date_time TIMESTAMP WITH TIME ZONE NOT NULL,
    value DOUBLE PRECISION NOT NULL DEFAULT 0,
    count BIGINT NOT NULL DEFAULT 0,
    metadata JSONB,
    version BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(name, date_time)
);

CREATE INDEX IF NOT EXISTS idx_hourly_metrics_name ON hourly_metrics(name);
CREATE INDEX IF NOT EXISTS idx_hourly_metrics_date_time ON hourly_metrics(date_time);
CREATE INDEX IF NOT EXISTS idx_hourly_metrics_name_date_time ON hourly_metrics(name, date_time);
