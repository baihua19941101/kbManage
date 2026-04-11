-- 0002_cluster_sync_status.sql
-- Persist cluster synchronization lifecycle status and timestamps.

ALTER TABLE clusters
  ADD COLUMN IF NOT EXISTS sync_status VARCHAR(32) NOT NULL DEFAULT 'idle' AFTER status,
  ADD COLUMN IF NOT EXISTS last_sync_at TIMESTAMP NULL DEFAULT NULL AFTER sync_status,
  ADD COLUMN IF NOT EXISTS last_success_at TIMESTAMP NULL DEFAULT NULL AFTER last_sync_at,
  ADD COLUMN IF NOT EXISTS sync_failure_reason VARCHAR(1024) NULL AFTER last_success_at;
