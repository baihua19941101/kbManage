-- 0003_scope_bindings.sql
-- Scope-level bindings and role metadata.

CREATE TABLE IF NOT EXISTS workspace_cluster_bindings (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  workspace_id BIGINT UNSIGNED NOT NULL,
  cluster_id BIGINT UNSIGNED NOT NULL,
  default_namespaces VARCHAR(1024) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_workspace_cluster (workspace_id, cluster_id),
  CONSTRAINT fk_wcb_workspace FOREIGN KEY (workspace_id) REFERENCES workspaces(id),
  CONSTRAINT fk_wcb_cluster FOREIGN KEY (cluster_id) REFERENCES clusters(id)
);

CREATE TABLE IF NOT EXISTS project_cluster_bindings (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  project_id BIGINT UNSIGNED NOT NULL,
  cluster_id BIGINT UNSIGNED NOT NULL,
  default_namespaces VARCHAR(1024) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_project_cluster (project_id, cluster_id),
  CONSTRAINT fk_pcb_project FOREIGN KEY (project_id) REFERENCES projects(id),
  CONSTRAINT fk_pcb_cluster FOREIGN KEY (cluster_id) REFERENCES clusters(id)
);

CREATE TABLE IF NOT EXISTS scope_roles (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  scope_type VARCHAR(32) NOT NULL,
  role_key VARCHAR(128) NOT NULL,
  name VARCHAR(128) NOT NULL,
  description VARCHAR(512) NULL,
  metadata_json JSON NULL,
  is_system TINYINT(1) NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_scope_role (scope_type, role_key)
);

ALTER TABLE scope_roles
  ADD COLUMN IF NOT EXISTS metadata_json JSON NULL AFTER description;

CREATE TABLE IF NOT EXISTS scope_role_bindings (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  subject_type VARCHAR(32) NOT NULL,
  subject_id BIGINT UNSIGNED NOT NULL,
  scope_type VARCHAR(32) NOT NULL,
  scope_id BIGINT UNSIGNED NOT NULL,
  scope_role_id BIGINT UNSIGNED NOT NULL,
  granted_by BIGINT UNSIGNED NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_scope_binding (subject_type, subject_id, scope_type, scope_id, scope_role_id),
  KEY idx_scope_binding_subject (subject_type, subject_id),
  KEY idx_scope_binding_scope (scope_type, scope_id),
  CONSTRAINT fk_srb_scope_role FOREIGN KEY (scope_role_id) REFERENCES scope_roles(id)
);

INSERT INTO scope_roles (scope_type, role_key, name, description, metadata_json, is_system)
VALUES
  ('workspace', 'platform-admin', 'Platform Admin', 'Workspace scoped super access', JSON_OBJECT('matrix', 'v1', 'tier', 'core'), 1),
  ('workspace', 'ops-operator', 'Ops Operator', 'Workspace scoped operations access', JSON_OBJECT('matrix', 'v1', 'tier', 'core'), 1),
  ('workspace', 'audit-reader', 'Audit Reader', 'Workspace scoped audit read access', JSON_OBJECT('matrix', 'v1', 'tier', 'core'), 1),
  ('workspace', 'readonly', 'Read Only', 'Workspace scoped read-only access', JSON_OBJECT('matrix', 'v1', 'tier', 'core'), 1),
  ('project', 'platform-admin', 'Platform Admin', 'Project scoped super access', JSON_OBJECT('matrix', 'v1', 'tier', 'core'), 1),
  ('project', 'ops-operator', 'Ops Operator', 'Project scoped operations access', JSON_OBJECT('matrix', 'v1', 'tier', 'core'), 1),
  ('project', 'audit-reader', 'Audit Reader', 'Project scoped audit read access', JSON_OBJECT('matrix', 'v1', 'tier', 'core'), 1),
  ('project', 'readonly', 'Read Only', 'Project scoped read-only access', JSON_OBJECT('matrix', 'v1', 'tier', 'core'), 1),
  ('workspace', 'workspace-owner', 'Workspace Owner', 'Workspace full access', JSON_OBJECT('compat', true), 1),
  ('workspace', 'workspace-viewer', 'Workspace Viewer', 'Workspace read-only', JSON_OBJECT('compat', true), 1),
  ('project', 'project-owner', 'Project Owner', 'Project full access', JSON_OBJECT('compat', true), 1),
  ('project', 'project-viewer', 'Project Viewer', 'Project read-only', JSON_OBJECT('compat', true), 1)
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  description = VALUES(description),
  metadata_json = COALESCE(scope_roles.metadata_json, VALUES(metadata_json)),
  is_system = VALUES(is_system);
