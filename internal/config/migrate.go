package config

import (
	"fmt"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Schema migration (SPEC-03 sec 4.4). The full migration chain (registry +
// backup + unmigratable handling) is filled in by the migration unit; this
// stub only rejects pre-v2 files until then.

// applyMigrations rewrites raw from version ver to SchemaVersionCurrent,
// backing up the original file. Stub: no migrations registered yet.
func applyMigrations(path string, raw []byte, ver int) ([]byte, error) {
	return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
		"config.toml: schema_version %d requires migration (not yet available in this build)", ver))
}
