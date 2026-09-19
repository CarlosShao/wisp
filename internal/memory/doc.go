// Package memory owns persistence beyond config (SPEC-01 §3): L1/L2/L3 memory
// stores (C13), the SQLite lifecycle, and the RetentionJob.
//
// Responsibilities:
//   - SQLite (WAL, single db-writer goroutine per D38b) with the schema of
//     SPEC-02 §3
//   - task_log.error_class: validated at the DAO boundary with observe's
//     ValidateErrorClass — exactly the 17 D37 values, no SQL CHECK (the DDL
//     of SPEC-02 §3 is contract-level and stays byte-identical; see
//     models.go validateTaskLog).
//   - retention windows (task_log/tool_call >30d, cost_daily >400d, ...)
//
// Non-responsibilities:
//   - no config storage semantics (config), no secret blobs (secret), no
//     cost metering policy (observe)
//
// DEFERRED(sqlite/schema): the storage core landed with ticket 04 (schema v1,
// WAL, single db-writer, retention, privacy ops); ticket 29 builds the L1/L2
// memory extraction flows on top of the DAOs here.
package memory
