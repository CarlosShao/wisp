// Package config owns config.toml as the single source of truth (SPEC-01 sec 3;
// D6/D36, C15): schema, validation with line numbers, hot reload, migration.
//
// Responsibilities:
//   - full config schema (C15/D36) with sections, unknown-key errors, and the
//     hot/cold reload split
//   - migration of older config versions (backup + fail-closed to
//     Unconfigured, never silent fallback)
//
// Non-responsibilities:
//   - no secret storage or DPAPI blobs (secret), no env var fork policy
//     (proc env layout), no GUI editor (panel, ticket 39)
//
// DEFERRED(schema/hot-reload/migration): implemented by ticket 05. This
// ticket only freezes the package boundary.
package config
