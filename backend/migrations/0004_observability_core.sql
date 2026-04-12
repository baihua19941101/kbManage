-- 0004_observability_core.sql
-- Core schema for observability data sources, alert governance and incident snapshots.

CREATE TABLE IF NOT EXISTS observability_data_sources (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  cluster_id BIGINT UNSIGNED NULL,
  type VARCHAR(32) NOT NULL,
  provider_kind VARCHAR(64) NOT NULL,
  name VARCHAR(128) NOT NULL,
  base_url VARCHAR(1024) NOT NULL,
  auth_secret_ref VARCHAR(256) NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'pending',
  last_verified_at TIMESTAMP NULL,
  last_error VARCHAR(1024) NULL,
  created_by BIGINT UNSIGNED NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_obs_ds_cluster_type_name (cluster_id, type, name),
  KEY idx_obs_ds_cluster (cluster_id),
  KEY idx_obs_ds_type (type),
  KEY idx_obs_ds_status (status),
  CONSTRAINT fk_obs_ds_cluster FOREIGN KEY (cluster_id) REFERENCES clusters(id)
);

CREATE TABLE IF NOT EXISTS observability_alert_rules (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(128) NOT NULL,
  description VARCHAR(1024) NULL,
  severity VARCHAR(16) NOT NULL,
  scope_snapshot_json LONGTEXT NULL,
  condition_expression TEXT NOT NULL,
  evaluation_window VARCHAR(64) NULL,
  notification_strategy LONGTEXT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'enabled',
  created_by BIGINT UNSIGNED NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_obs_rule_name (name),
  KEY idx_obs_rule_status (status),
  KEY idx_obs_rule_severity (severity)
);

CREATE TABLE IF NOT EXISTS observability_notification_targets (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(128) NOT NULL,
  target_type VARCHAR(32) NOT NULL,
  config_ref VARCHAR(256) NULL,
  scope_snapshot LONGTEXT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'active',
  created_by BIGINT UNSIGNED NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_obs_target_name (name),
  KEY idx_obs_target_type (target_type),
  KEY idx_obs_target_status (status)
);

CREATE TABLE IF NOT EXISTS observability_silence_windows (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(128) NOT NULL,
  scope_snapshot LONGTEXT NULL,
  reason VARCHAR(1024) NULL,
  starts_at TIMESTAMP NOT NULL,
  ends_at TIMESTAMP NOT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'scheduled',
  created_by BIGINT UNSIGNED NULL,
  canceled_by BIGINT UNSIGNED NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_obs_silence_window (starts_at, ends_at),
  KEY idx_obs_silence_status (status)
);

CREATE TABLE IF NOT EXISTS observability_alert_incidents (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  source_incident_key VARCHAR(256) NOT NULL,
  rule_id BIGINT UNSIGNED NULL,
  cluster_id BIGINT UNSIGNED NULL,
  workspace_id BIGINT UNSIGNED NULL,
  project_id BIGINT UNSIGNED NULL,
  resource_kind VARCHAR(64) NULL,
  resource_name VARCHAR(255) NULL,
  namespace VARCHAR(255) NULL,
  severity VARCHAR(16) NOT NULL,
  status VARCHAR(32) NOT NULL,
  summary VARCHAR(1024) NULL,
  starts_at TIMESTAMP NULL,
  acknowledged_at TIMESTAMP NULL,
  resolved_at TIMESTAMP NULL,
  last_synced_at TIMESTAMP NULL,
  timeline_json LONGTEXT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_obs_incident_source_key (source_incident_key),
  KEY idx_obs_incident_rule (rule_id),
  KEY idx_obs_incident_cluster (cluster_id),
  KEY idx_obs_incident_workspace (workspace_id),
  KEY idx_obs_incident_project (project_id),
  KEY idx_obs_incident_status (status),
  CONSTRAINT fk_obs_incident_rule FOREIGN KEY (rule_id) REFERENCES observability_alert_rules(id),
  CONSTRAINT fk_obs_incident_cluster FOREIGN KEY (cluster_id) REFERENCES clusters(id),
  CONSTRAINT fk_obs_incident_workspace FOREIGN KEY (workspace_id) REFERENCES workspaces(id),
  CONSTRAINT fk_obs_incident_project FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE TABLE IF NOT EXISTS observability_alert_handling_records (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  incident_id BIGINT UNSIGNED NOT NULL,
  action_type VARCHAR(32) NOT NULL,
  content TEXT NULL,
  acted_by BIGINT UNSIGNED NULL,
  acted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_obs_record_incident (incident_id),
  KEY idx_obs_record_action_type (action_type),
  CONSTRAINT fk_obs_record_incident FOREIGN KEY (incident_id) REFERENCES observability_alert_incidents(id)
);
