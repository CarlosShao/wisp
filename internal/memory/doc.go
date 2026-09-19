// Package memory owns persistence beyond config (SPEC-01 §3): L1/L2/L3 memory
// stores (C13), the SQLite lifecycle, and the RetentionJob.
//
// Responsibilities:
//   - SQLite (WAL, single db-writer goroutine per D38b) with the schema of
//     SPEC-02 §3
//   - task_log.error_class column: a CHECK-constrained enum whose 17 values
//     are exactly observe's ErrorClass set (D37). The DAO-level constraint is
//     RESERVED here and implemented with the schema in ticket 04; observe
//     ValidateErrorClass is the shared validator until then.
//   - retention windows (task_log/tool_call >30d, cost_daily >400d, ...)
//
// Non-responsibilities:
//   - no config storage semantics (config), no secret blobs (secret), no
//     cost metering policy (observe)
//
// DEFERRED(sqlite/schema): implemented by ticket 04 (core) and ticket 29
// (L1/L2 memory). This ticket only freezes the package boundary and the
// error_class constraint contract above.
package memory
