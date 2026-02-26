CREATE TABLE urls (
  id BIGSERIAL PRIMARY KEY,
  short_code VARCHAR(10) UNIQUE NOT NULL,
  original_url TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMP NULL,
  user_id VARCHAR(50) NULL,
  is_active BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_urls_short_code ON urls(short_code);
CREATE INDEX idx_urls_created_at ON urls(created_at);
