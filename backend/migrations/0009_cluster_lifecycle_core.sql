CREATE TABLE IF NOT EXISTS cluster_lifecycle_records (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(128) NOT NULL,
  display_name VARCHAR(128) NOT NULL,
  lifecycle_mode VARCHAR(32) NOT NULL,
  infrastructure_type VARCHAR(64) NOT NULL,
  driver_ref VARCHAR(128) NOT NULL,
  driver_version VARCHAR(64) NOT NULL,
  workspace_id BIGINT UNSIGNED NOT NULL,
  project_id BIGINT UNSIGNED NULL,
  status VARCHAR(32) NOT NULL,
  registration_status VARCHAR(32) NOT NULL,
  health_status VARCHAR(32) NOT NULL,
  kubernetes_version VARCHAR(64) NOT NULL,
  target_version VARCHAR(64) NULL,
  node_pool_summary TEXT NULL,
  last_validation_status VARCHAR(32) NOT NULL,
  last_validation_at DATETIME NULL,
  last_operation_id BIGINT UNSIGNED NULL,
  template_id BIGINT UNSIGNED NULL,
  retirement_reason TEXT NULL,
  created_by BIGINT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cluster_driver_versions (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  driver_key VARCHAR(64) NOT NULL,
  version VARCHAR(64) NOT NULL,
  display_name VARCHAR(128) NOT NULL,
  provider_type VARCHAR(64) NOT NULL,
  status VARCHAR(32) NOT NULL,
  capability_profile_version VARCHAR(64) NOT NULL,
  schema_version VARCHAR(64) NOT NULL,
  release_notes TEXT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uniq_cluster_driver_versions (driver_key, version)
);

CREATE TABLE IF NOT EXISTS cluster_capability_matrix_entries (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  owner_type VARCHAR(32) NOT NULL,
  owner_ref VARCHAR(128) NOT NULL,
  capability_domain VARCHAR(64) NOT NULL,
  support_level VARCHAR(32) NOT NULL,
  compatibility_status VARCHAR(32) NOT NULL,
  constraints_summary TEXT NULL,
  recommended_for TEXT NULL,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uniq_cluster_capability_matrix_entries (owner_type, owner_ref, capability_domain)
);

CREATE TABLE IF NOT EXISTS cluster_templates (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(128) NOT NULL UNIQUE,
  description TEXT NULL,
  infrastructure_type VARCHAR(64) NOT NULL,
  driver_key VARCHAR(64) NOT NULL,
  driver_version_range VARCHAR(128) NOT NULL,
  required_capabilities TEXT NULL,
  parameter_schema TEXT NULL,
  default_values TEXT NULL,
  status VARCHAR(32) NOT NULL,
  created_by BIGINT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cluster_lifecycle_operations (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  cluster_id BIGINT UNSIGNED NULL,
  operation_type VARCHAR(32) NOT NULL,
  trigger_source VARCHAR(32) NOT NULL,
  status VARCHAR(32) NOT NULL,
  risk_level VARCHAR(32) NOT NULL,
  requested_by BIGINT UNSIGNED NOT NULL,
  request_snapshot TEXT NULL,
  result_summary TEXT NULL,
  failure_reason TEXT NULL,
  started_at DATETIME NULL,
  completed_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cluster_upgrade_plans (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  cluster_id BIGINT UNSIGNED NOT NULL,
  from_version VARCHAR(64) NOT NULL,
  to_version VARCHAR(64) NOT NULL,
  window_start DATETIME NULL,
  window_end DATETIME NULL,
  precheck_status VARCHAR(32) NOT NULL,
  impact_summary TEXT NULL,
  status VARCHAR(32) NOT NULL,
  last_operation_id BIGINT UNSIGNED NULL,
  created_by BIGINT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cluster_node_pool_profiles (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  cluster_id BIGINT UNSIGNED NOT NULL,
  name VARCHAR(128) NOT NULL,
  role VARCHAR(32) NOT NULL,
  desired_count INT NOT NULL,
  current_count INT NOT NULL,
  min_count INT NOT NULL,
  max_count INT NOT NULL,
  version VARCHAR(64) NOT NULL,
  zone_refs TEXT NULL,
  status VARCHAR(32) NOT NULL,
  last_operation_id BIGINT UNSIGNED NULL,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cluster_lifecycle_audit_events (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  action VARCHAR(128) NOT NULL,
  actor_user_id BIGINT UNSIGNED NOT NULL,
  cluster_id BIGINT UNSIGNED NULL,
  target_type VARCHAR(32) NOT NULL,
  target_ref VARCHAR(128) NOT NULL,
  outcome VARCHAR(32) NOT NULL,
  detail_snapshot TEXT NULL,
  occurred_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
