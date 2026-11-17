package config

// Ticket 39 AC#4's structured load-error read face (readface wave 1,
// 2026-11-17): `errors[]{line,message}` for a config.toml that was edited
// externally - unknown keys and wrong types carry their 1-based line, exactly
// what the strict decoder already had and only the prose path surfaced.
//
// WHY A SECOND DECODE. The Manager keeps the last good config when a load
// fails (CheckAndReload's contract), so the runtime keeps running while the
// file on disk is broken - and the page editing that file needs the errors
// NOW, not at the next reload. ParseErrors re-runs the READ pipeline's
// content half against the bytes on disk (parse -> migration -> strict decode
// -> validation) and returns what it finds without touching the Manager's
// state: it writes no bytes, keeps no handle, and cannot change what the
// process is running.
//
// Line numbers are honest or absent: only the strict decoder's two error
// shapes carry a position, and a validation failure (a semantic rule, not a
// line) reports Line 0 with its message verbatim. A caller must render "no
// line" for 0, never pretend the file pointed at line 1.

import (
	"errors"
	"fmt"
	"os"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
)

// LoadError is one problem the read pipeline found in the file: the 1-based
// line the decoder pointed at (0 = no line - a read failure or a semantic
// validation failure, which name no line), and the message.
type LoadError struct {
	Line    int
	Message string
}

// ParseErrorsFile runs the read pipeline's content half over the file at path
// and returns every problem it finds, first failure stage first. A file that
// loads clean answers nil - which the caller renders as "no errors", never as
// "the file could not be checked".
func ParseErrorsFile(path string) []LoadError {
	raw, err := os.ReadFile(path)
	if err != nil {
		return []LoadError{{Line: 0, Message: fmt.Sprintf("config.toml read: %v", err)}}
	}
	return ParseErrorsBytes(raw)
}

// ParseErrors is the Manager's own read face for the same answers, over the
// path this Manager watches. It reads m.path under the lock and decodes the
// bytes off the lock: the decode reads the file, never cur.
func (m *Manager) ParseErrors() []LoadError {
	if m == nil {
		return nil
	}
	m.mu.Lock()
	path := m.path
	m.mu.Unlock()
	return ParseErrorsFile(path)
}

// ParseErrorsBytes is the content half of the read pipeline over raw bytes,
// shaped as structured errors instead of the pipeline's first error. The
// stage order is readConfigFile's own (loader.go): version probe -> parse ->
// migration -> strict decode -> validation. Every stage that fails contributes
// its rows and stops the walk - a file whose syntax is broken has no semantic
// errors worth listing yet.
func ParseErrorsBytes(raw []byte) []LoadError {
	ver, peekErr := peekSchemaVersion(raw)
	if ver == 0 && peekErr != nil {
		return []LoadError{{Line: 0, Message: "config.toml parse: " + peekErr.Error()}}
	}
	if peekErr != nil && ver >= SchemaVersionCurrent {
		return []LoadError{{Line: 0, Message: "config.toml parse: " + peekErr.Error()}}
	}
	if ver > SchemaVersionCurrent {
		return []LoadError{{Line: 0, Message: fmt.Sprintf(
			"config.toml: schema_version %d was written by a newer build (this build understands %d); upgrade Wisp or restore a backup",
			ver, SchemaVersionCurrent)}}
	}
	if ver != SchemaVersionCurrent {
		migrated, err := migrateRaw(raw, ver)
		if err != nil {
			return []LoadError{{Line: 0, Message: err.Error()}}
		}
		raw = migrated
	}
	cfg := NewDefaults()
	if err := decodeStrictInto(raw, cfg); err != nil {
		return decodeErrorRows(err)
	}
	if verr := validate(cfg); verr != nil {
		// Validation names a rule, not a line: Line 0 is the "no line" marker
		// by contract (see LoadError).
		return []LoadError{{Line: 0, Message: verr.Error()}}
	}
	return nil
}

// decodeErrorRows turns one decode error into structured rows. It is the
// structured twin of formatDecodeError (parse.go): the strict shape may carry
// MANY offending keys (one row each, with the key's own line), the single
// DecodeError shape carries one position, and anything else is reported as
// one lineless row. formatDecodeError builds its prose from this, so the two
// faces cannot disagree about what the decoder said.
func decodeErrorRows(err error) []LoadError {
	var sme *toml.StrictMissingError
	if errors.As(err, &sme) {
		rows := make([]LoadError, 0, len(sme.Errors))
		for _, de := range sme.Errors {
			row, _ := de.Position()
			rows = append(rows, LoadError{
				Line:    row,
				Message: fmt.Sprintf("unknown key %q at line %d", strings.Join(de.Key(), "."), row),
			})
		}
		return rows
	}
	var de *toml.DecodeError
	if errors.As(err, &de) {
		row, col := de.Position()
		msg := strings.TrimPrefix(de.Error(), "toml: ")
		return []LoadError{{Line: row, Message: fmt.Sprintf("line %d, col %d: %s", row, col, msg)}}
	}
	return []LoadError{{Line: 0, Message: err.Error()}}
}
