CREATE TABLE IF NOT EXISTS backup_policies (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  description TEXT NULL,
  scope_type VARCHAR(64) NOT NULL,
  scope_ref VARCHAR(128) NOT NULL,
  workspace_id BIGINT UNSIGNED NOT NULL,
  project_id BIGINT UNSIGNED NULL,
  execution_mode VARCHAR(32) NOT NULL,
  schedule_expression VARCHAR(256) NULL,
  retention_rule VARCHAR(256) NOT NULL,
  consistency_level VARCHAR(64) NOT NULL,
  status VARCHAR(32) NOT NULL,
  owner_user_id BIGINT UNSIGNED NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  UNIQUE KEY uk_backup_policy_scope_name (scope_type, scope_ref, name),
  KEY idx_backup_policy_workspace (workspace_id),
  KEY idx_backup_policy_project (project_id)
);

CREATE TABLE IF NOT EXISTS restore_points (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  policy_id BIGINT UNSIGNED NOT NULL,
  workspace_id BIGINT UNSIGNED NOT NULL,
  project_id BIGINT UNSIGNED NULL,
  scope_snapshot JSON NOT NULL,
  backup_started_at DATETIME(3) NOT NULL,
  backup_completed_at DATETIME(3) NULL,
  duration_seconds INT NOT NULL DEFAULT 0,
  result VARCHAR(32) NOT NULL,
  consistency_summary TEXT NULL,
  failure_reason TEXT NULL,
  storage_ref VARCHAR(256) NULL,
  expires_at DATETIME(3) NULL,
  created_by BIGINT UNSIGNED NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  KEY idx_restore_points_policy (policy_id),
  KEY idx_restore_points_workspace (workspace_id)
);

CREATE TABLE IF NOT EXISTS restore_jobs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  restore_point_id BIGINT UNSIGNED NOT NULL,
  workspace_id BIGINT UNSIGNED NOT NULL,
  project_id BIGINT UNSIGNED NULL,
  job_type VARCHAR(64) NOT NULL,
  source_environment VARCHAR(128) NULL,
  target_environment VARCHAR(128) NOT NULL,
  scope_selection JSON NOT NULL,
  conflict_summary TEXT NULL,
  consistency_notice TEXT NULL,
  status VARCHAR(32) NOT NULL,
  result_summary TEXT NULL,
  failure_reason TEXT NULL,
  requested_by BIGINT UNSIGNED NOT NULL,
  started_at DATETIME(3) NULL,
  completed_at DATETIME(3) NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  KEY idx_restore_jobs_restore_point (restore_point_id),
  KEY idx_restore_jobs_workspace (workspace_id)
);

CREATE TABLE IF NOT EXISTS migration_plans (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  workspace_id BIGINT UNSIGNED NOT NULL,
  project_id BIGINT UNSIGNED NULL,
  source_cluster_id BIGINT UNSIGNED NOT NULL,
  target_cluster_id BIGINT UNSIGNED NOT NULL,
  scope_selection JSON NOT NULL,
  mapping_rules JSON NULL,
  cutover_steps JSON NULL,
  status VARCHAR(32) NOT NULL,
  created_by BIGINT UNSIGNED NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  UNIQUE KEY uk_migration_plan_name (name)
);

CREATE TABLE IF NOT EXISTS dr_drill_plans (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  description TEXT NULL,
  workspace_id BIGINT UNSIGNED NOT NULL,
  project_id BIGINT UNSIGNED NULL,
  scope_selection JSON NOT NULL,
  rpo_target_minutes INT NOT NULL,
  rto_target_minutes INT NOT NULL,
  role_assignments JSON NULL,
  cutover_procedure JSON NOT NULL,
  validation_checklist JSON NOT NULL,
  status VARCHAR(32) NOT NULL,
  created_by BIGINT UNSIGNED NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  UNIQUE KEY uk_dr_drill_plan_name (name)
);

CREATE TABLE IF NOT EXISTS dr_drill_records (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  plan_id BIGINT UNSIGNED NOT NULL,
  workspace_id BIGINT UNSIGNED NOT NULL,
  project_id BIGINT UNSIGNED NULL,
  started_at DATETIME(3) NOT NULL,
  completed_at DATETIME(3) NULL,
  actual_rpo_minutes INT NOT NULL DEFAULT 0,
  actual_rto_minutes INT NOT NULL DEFAULT 0,
  status VARCHAR(32) NOT NULL,
  step_results JSON NULL,
  validation_results JSON NULL,
  incident_notes TEXT NULL,
  executed_by BIGINT UNSIGNED NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  KEY idx_dr_drill_record_plan (plan_id)
);

CREATE TABLE IF NOT EXISTS dr_drill_reports (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  drill_record_id BIGINT UNSIGNED NOT NULL,
  goal_assessment TEXT NOT NULL,
  gap_summary TEXT NULL,
  issues_found JSON NULL,
  improvement_actions JSON NOT NULL,
  published_at DATETIME(3) NOT NULL,
  published_by BIGINT UNSIGNED NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  UNIQUE KEY uk_dr_drill_report_record (drill_record_id)
);

CREATE TABLE IF NOT EXISTS backup_audit_events (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  action VARCHAR(128) NOT NULL,
  actor_user_id BIGINT UNSIGNED NOT NULL,
  target_type VARCHAR(64) NOT NULL,
  target_ref VARCHAR(128) NOT NULL,
  workspace_id BIGINT UNSIGNED NOT NULL,
  project_id BIGINT UNSIGNED NULL,
  scope_snapshot JSON NULL,
  outcome VARCHAR(32) NOT NULL,
  detail_snapshot JSON NULL,
  occurred_at DATETIME(3) NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  KEY idx_backup_audit_workspace (workspace_id),
  KEY idx_backup_audit_target (target_type, target_ref),
  KEY idx_backup_audit_action (action),
  KEY idx_backup_audit_occurred (occurred_at)
);
