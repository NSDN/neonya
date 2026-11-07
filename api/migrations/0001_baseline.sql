-- +goose Up
-- pgcrypto：gen_random_uuid()；pg_trgm：标题与正文的模糊搜索（见数据模型 §6）。
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- +goose Down
-- 扩展不回收：可能已被其他对象依赖。
SELECT 1;
