package config

// Ticket 39 AC#4's third half: the disk file's fingerprint - the fact the
// config editor's optimistic-concurrency handshake is made of.
//
// WHY A HASH. The page reads config.get (a snapshot of the tree plus, since
// readface wave 1, the load errors), the operator hand-edits config.toml
// behind the process, and the page then saves a field with the value its
// stale form still holds. mergeWrite (ticket 226) already guarantees the
// write never reverts FOREIGN hand edits - it re-reads and merges - but a
// stale save of the key the page itself addresses would still silently
// replace the external edit of that same key, and nothing on the receipt
// would say the disk had moved. The fix is the standard one (an ETag, an
// If-Match header, an expected-version field): the read face hands the page
// the file's sha256 as it stands at read time, the page echoes it back on
// config.set, and the write leg refuses with refused-drift when the file's
// bytes are no longer what the page saw. No silent overwrite: the page is
// told the disk changed and re-reads.
//
// WHY SHA256 AND NOT MTIME-SIZE: a wall-clock stamp cannot tell "rewritten
// with the same bytes" from "rewritten with different bytes" and lies about
// order after any clock change; the hash judges content, and content is the
// thing the no-op and drift contracts are about. The hash is of the RAW disk
// bytes - never of the migrated or canonicalized config - because the bytes
// are what a hand edit changes and what mergeWrite's no-op promise (write
// nothing when the file already carries the value) is about.

import (
	"crypto/sha256"
	"encoding/hex"
	"os"

	"github.com/CarlosShao/wisp/internal/observe"
)

// FileSHA256 returns the lowercase hex sha256 of the watched config file's
// bytes as they are on disk right now. It reads nothing into the Manager and
// changes nothing: no watch state is touched, no migration runs, no handle is
// kept.
//
// A MISSING file is not an error here: it answers ("", nil) - "there is
// nothing to fingerprint", which the wire renders as the baseSHA key being
// absent, and a page that never received a fingerprint cannot be accused of
// holding a stale one. An unreadable file IS an error: the caller must not
// guess whether the page's view is current when the disk says nothing.
func (m *Manager) FileSHA256() (string, error) {
	if m == nil {
		return "", observe.New(observe.ClassConfig,
			"config: no loaded Manager, the config file's fingerprint has no answer")
	}
	raw, err := os.ReadFile(m.path)
	if err != nil {
		if fileMissing(err) {
			// No file on disk: first-run has not written one yet. There is no
			// base to compare against, so the honest answer is "no fingerprint",
			// never a hash of an empty read.
			return "", nil
		}
		return "", observe.Wrap(observe.ClassConfig, err,
			"config: cannot read config.toml to fingerprint it; whether the page's view is current has no answer")
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}
