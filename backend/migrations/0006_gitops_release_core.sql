-- 0006_gitops_release_core.sql
-- Core schema for 004 gitops and release foundation.

CREATE TABLE IF NOT EXISTS gitops_delivery_sources (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(128) NOT NULL,
  source_type VARCHAR(32) NOT NULL,
  endpoint VARCHAR(1024) NOT NULL,
  default_ref VARCHAR(256) NULL,
  credential_ref VARCHAR(256) NULL,
  workspace_id BIGINT UNSIGNED NULL,
  project_id BIGINT UNSIGNED NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'pending',
  last_verified_at TIMESTAMP NULL,
  last_error_message TEXT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_gitops_sources_workspace (workspace_id),
  KEY idx_gitops_sources_project (project_id),
  KEY idx_gitops_sources_status (status),
  UNIQUE KEY uk_gitops_sources_scope_name (workspace_id, project_id, name)
);

CREATE TABLE IF NOT EXISTS gitops_cluster_target_groups (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(128) NOT NULL,
  workspace_id BIGINT UNSIGNED NOT NULL,
  project_id BIGINT UNSIGNED NULL,
  cluster_refs_json LONGTEXT NULL,
  cluster_selector_snapshot LONGTEXT NULL,
  description VARCHAR(1024) NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'active',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_gitops_target_groups_workspace (workspace_id),
  KEY idx_gitops_target_groups_project (project_id),
  KEY idx_gitops_target_groups_status (status),
  UNIQUE KEY uk_gitops_target_groups_scope_name (workspace_id, project_id, name)
);

CREATE TABLE IF NOT EXISTS gitops_delivery_units (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(128) NOT NULL,
  workspace_id BIGINT UNSIGNED NOT NULL,
  project_id BIGINT UNSIGNED NULL,
  source_id BIGINT UNSIGNED NOT NULL,
  source_path VARCHAR(1024) NULL,
  default_namespace VARCHAR(255) NULL,
  sync_mode VARCHAR(32) NOT NULL DEFAULT 'manual',
  release_policy_json LONGTEXT NULL,
  desired_revision VARCHAR(256) NULL,
  desired_app_version VARCHAR(128) NULL,
  desired_config_version VARCHAR(128) NULL,
  paused TINYINT(1) NOT NULL DEFAULT 0,
  delivery_status VARCHAR(32) NOT NULL DEFAULT 'unknown',
  last_synced_at TIMESTAMP NULL,
  last_release_id BIGINT UNSIGNED NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_gitops_units_workspace (workspace_id),
  KEY idx_gitops_units_project (project_id),
  KEY idx_gitops_units_source (source_id),
  KEY idx_gitops_units_status (delivery_status),
  UNIQUE KEY uk_gitops_units_scope_name (workspace_id, project_id, name)
);

CREATE TABLE IF NOT EXISTS gitops_environment_stages (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  delivery_unit_id BIGINT UNSIGNED NOT NULL,
  name VARCHAR(128) NOT NULL,
  order_index INT NOT NULL,
  target_group_id BIGINT UNSIGNED NOT NULL,
  promotion_mode VARCHAR(32) NOT NULL DEFAULT 'manual',
  paused TINYINT(1) NOT NULL DEFAULT 0,
  status VARCHAR(32) NOT NULL DEFAULT 'idle',
  last_entered_at TIMESTAMP NULL,
  last_completed_at TIMESTAMP NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_gitops_stages_unit (delivery_unit_id),
  KEY idx_gitops_stages_target_group (target_group_id),
  KEY idx_gitops_stages_status (status),
  UNIQUE KEY uk_gitops_stages_order (delivery_unit_id, order_index),
  UNIQUE KEY uk_gitops_stages_name (delivery_unit_id, name)
);

CREATE TABLE IF NOT EXISTS gitops_configuration_overlays (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  delivery_unit_id BIGINT UNSIGNED NOT NULL,
  environment_stage_id BIGINT UNSIGNED NULL,
  overlay_type VARCHAR(32) NOT NULL,
  overlay_ref VARCHAR(1024) NOT NULL,
  precedence INT NOT NULL DEFAULT 0,
  effective_scope_json LONGTEXT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_gitops_overlays_unit (delivery_unit_id),
  KEY idx_gitops_overlays_stage (environment_stage_id),
  KEY idx_gitops_overlays_type (overlay_type)
);

CREATE TABLE IF NOT EXISTS gitops_release_revisions (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  delivery_unit_id BIGINT UNSIGNED NOT NULL,
  source_revision VARCHAR(256) NOT NULL,
  app_version VARCHAR(128) NULL,
  config_version VARCHAR(128) NULL,
  effective_scope_json LONGTEXT NULL,
  release_notes_summary TEXT NULL,
  created_by BIGINT UNSIGNED NOT NULL,
  rollback_available TINYINT(1) NOT NULL DEFAULT 0,
  status VARCHAR(32) NOT NULL DEFAULT 'historical',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_gitops_revisions_unit (delivery_unit_id),
  KEY idx_gitops_revisions_status (status),
  KEY idx_gitops_revisions_created_by (created_by)
);

CREATE TABLE IF NOT EXISTS gitops_delivery_operations (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  request_id VARCHAR(64) NOT NULL,
  operator_id BIGINT UNSIGNED NOT NULL,
  delivery_unit_id BIGINT UNSIGNED NOT NULL,
  environment_stage_id BIGINT UNSIGNED NULL,
  action_type VARCHAR(32) NOT NULL,
  target_release_id BIGINT UNSIGNED NULL,
  payload_json LONGTEXT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'pending',
  progress_percent INT NOT NULL DEFAULT 0,
  result_summary TEXT NULL,
  failure_reason TEXT NULL,
  started_at TIMESTAMP NULL,
  completed_at TIMESTAMP NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_gitops_operations_request_id (request_id),
  KEY idx_gitops_operations_unit (delivery_unit_id),
  KEY idx_gitops_operations_stage (environment_stage_id),
  KEY idx_gitops_operations_action_type (action_type),
  KEY idx_gitops_operations_status (status),
  KEY idx_gitops_operations_operator (operator_id),
  KEY idx_gitops_operations_target_release (target_release_id)
);
