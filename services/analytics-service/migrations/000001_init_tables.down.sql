-- Удаляем индексы для таблицы clicks
DROP INDEX IF EXISTS idx_clicks_short_code;
DROP INDEX IF EXISTS idx_clicks_clicked_at;
DROP INDEX IF EXISTS idx_clicks_ip_address;

-- Удаляем индекс для таблицы url_stats
DROP INDEX IF EXISTS idx_url_stats_updated_at;

-- Удаляем таблицы (важно: порядок важен из-за внешних ключей, если они будут)
DROP TABLE IF EXISTS clicks;
DROP TABLE IF EXISTS url_stats;
