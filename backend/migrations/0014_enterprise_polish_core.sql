CREATE TABLE IF NOT EXISTS enterprise_permission_change_trails (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  workspace_id BIGINT UNSIGNED NOT NULL,
  project_id BIGINT UNSIGNED NULL,
  subject_type VARCHAR(64) NOT NULL,
  subject_ref VARCHAR(128) NOT NULL,
  source_identity VARCHAR(128) NULL,
  change_type VARCHAR(64) NOT NULL,
  before_state TEXT NULL,
  after_state TEXT NULL,
  authorization_basis TEXT NULL,
  approval_reference TEXT NULL,
  scope_type VARCHAR(64) NOT NULL,
  scope_ref VARCHAR(128) NOT NULL,
  evidence_completeness VARCHAR(32) NOT NULL,
  changed_at DATETIME NOT NULL,
  changed_by BIGINT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS enterprise_key_operation_traces LIKE enterprise_permission_change_trails;
CREATE TABLE IF NOT EXISTS enterprise_cross_team_authorization_snapshots LIKE enterprise_permission_change_trails;
CREATE TABLE IF NOT EXISTS enterprise_governance_risk_events LIKE enterprise_permission_change_trails;
CREATE TABLE IF NOT EXISTS enterprise_governance_coverage_snapshots LIKE enterprise_permission_change_trails;
CREATE TABLE IF NOT EXISTS enterprise_governance_report_packages LIKE enterprise_permission_change_trails;
CREATE TABLE IF NOT EXISTS enterprise_export_records LIKE enterprise_permission_change_trails;
CREATE TABLE IF NOT EXISTS enterprise_delivery_artifacts LIKE enterprise_permission_change_trails;
CREATE TABLE IF NOT EXISTS enterprise_delivery_readiness_bundles LIKE enterprise_permission_change_trails;
CREATE TABLE IF NOT EXISTS enterprise_delivery_checklist_items LIKE enterprise_permission_change_trails;
CREATE TABLE IF NOT EXISTS enterprise_governance_action_items LIKE enterprise_permission_change_trails;
