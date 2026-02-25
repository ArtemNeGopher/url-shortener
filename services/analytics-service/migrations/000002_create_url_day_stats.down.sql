-- Удаляем триггер
DROP TRIGGER IF EXISTS trigger_update_url_day_stats_updated_at ON url_day_stats;

-- Удаляем индексы
DROP INDEX IF EXISTS idx_url_day_stats_updated_at;
DROP INDEX IF EXISTS idx_url_day_stats_short_code;

-- Удаляем таблицу
DROP TABLE IF EXISTS url_day_stats;
