-- Расширение для uuid
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Таблица для хранения статистики по переходам за день
CREATE TABLE url_day_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    short_code VARCHAR(10) NOT NULL,
    date DATE NOT NULL,
    total_clicks BIGINT NOT NULL DEFAULT 0,
    unique_visitors BIGINT NOT NULL DEFAULT 0,
    referers JSONB DEFAULT '[]'::jsonb,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индексы для url_day_stats
CREATE INDEX idx_url_day_stats_updated_at ON url_day_stats(updated_at);
CREATE INDEX idx_url_day_stats_short_code ON url_day_stats(short_code);

-- Триггеры
CREATE TRIGGER trigger_update_url_day_stats_updated_at
    BEFORE UPDATE ON url_day_stats
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
