-- CREATE EXTENSION IF NOT EXISTS geoip;

-- Таблица для хранения кликов по коротким ссылкам
CREATE TABLE clicks (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(10) NOT NULL,
    clicked_at TIMESTAMP NOT NULL DEFAULT NOW(),
    ip_address INET NULL,
    user_agent TEXT NULL,
    referer TEXT NULL,
    country VARCHAR(2) NULL
);

-- Индексы для таблицы clicks
CREATE INDEX idx_clicks_short_code ON clicks(short_code);
CREATE INDEX idx_clicks_clicked_at ON clicks(clicked_at);
CREATE INDEX idx_clicks_ip_address ON clicks(ip_address);

-- Таблица для агрегированной статистики по URL
CREATE TABLE url_stats (
    short_code VARCHAR(10) PRIMARY KEY,
    total_clicks BIGINT NOT NULL DEFAULT 0,
    unique_visitors BIGINT NOT NULL DEFAULT 0,
    last_clicked_at TIMESTAMP NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индекс для таблицы url_stats
CREATE INDEX idx_url_stats_updated_at ON url_stats(updated_at);
