CREATE TABLE IF NOT EXISTS daily_metrics (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    value DOUBLE PRECISION NOT NULL DEFAULT 0,
    count BIGINT NOT NULL DEFAULT 0,
    metadata JSONB,
    version BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(name, date)
);

CREATE INDEX IF NOT EXISTS idx_daily_metrics_name ON daily_metrics(name);
CREATE INDEX IF NOT EXISTS idx_daily_metrics_date ON daily_metrics(date);
CREATE INDEX IF NOT EXISTS idx_daily_metrics_name_date ON daily_metrics(name, date);
