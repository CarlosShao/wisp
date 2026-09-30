package session

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

// ---------------------------------------------------------------------------
// The host-minted session identity (ticket 224, approval record A435)
//
// A435 states the definition in two hard clauses, and this file is only the
// first of them:
//
//	1. 不许派生 - the id is NOT derived from the process number, the clock, the
//	   working directory or anything else recomputable. It is minted from
//	   crypto/rand and nothing else.
//	2. 结束点今天＝进程退出 - the identity this process minted dies with it, and
//	   the next process mints a different one.
//
// Why clause 1 is a security property and not a style preference: every D45
// session grant row on disk is keyed by this string (approval_grant.session_id,
// SPEC-02 §3). If the string were recomputable, the next boot would recompute
// the previous boot's key, ListGrantsBySession would hand those rows back, and
// "本次会话内允许" would silently BECOME the permanent免审通行证 that
// internal/perm/store.go:19-27 and PLAN.md:1642 exist to keep out. The mint is
// therefore the only thing standing between "session-scoped" and "forever",
// which is why ticket 224's AC#3 (重启必失效) is satisfied by construction here
// rather than by a cleanup step somewhere else.
//
// The construction mirrors approval's native nonce (internal/agent/approval/
// approval.go:242 mintGrant): crypto/rand + hex + a readable prefix, so an audit
// line can point at a session without carrying anything that looks like a path
// or a tool name.
// ---------------------------------------------------------------------------

// idBytes is the entropy behind one minted identity: 128 bits. It is not a
// tunable - the collision resistance is what makes "guess the live session key"
// unattractive, and approval's own grant nonce uses the same shape.
const idBytes = 16

// idPrefix exists so a session id is recognisable in a log line and so Valid
// has something structural to check. It is also the reason a hand-written
// literal like "session-before-restart" can never be mistaken for one: that
// string does not carry the prefix and is not hex.
const idPrefix = "sess_"

// ID is one host-minted session identity. The zero value is deliberately NOT a
// usable identity: Valid reports false for it, so a session that was never
// minted cannot authorize anything.
type ID string

// String renders the identity verbatim. It is the value stored in
// approval_grant.session_id and read back by ListGrantsBySession.
func (id ID) String() string { return string(id) }

// Valid reports whether s has the shape Mint produces: the prefix, then exactly
// 2*idBytes lowercase hex characters.
//
// It is exported for one reason - it is the lock ticket 224's nail precheck
// demands. The guard in cmd/wisp/run_mode101_test.go writes grants under
// literals like "session-before-restart"; those literals must be provably
// OUTSIDE the set this function accepts, forever, or a production mint could one
// day collide with a test fixture and the anti-cross-session control (AC#4)
// would go green while granting nothing.
func (id ID) Valid() bool {
	s := string(id)
	if !strings.HasPrefix(s, idPrefix) {
		return false
	}
	body := s[len(idPrefix):]
	if len(body) != 2*idBytes {
		return false
	}
	for _, c := range body {
		switch {
		case c >= '0' && c <= '9':
		case c >= 'a' && c <= 'f':
		default:
			return false
		}
	}
	return true
}

// Mint mints one session identity for this process. Call it once, at the top of
// a composition root.
//
// A failure is a broken RNG and nothing else, and it is returned rather than
// papered over: a host that cannot mint an identity has no session, and with no
// session every D45 grant row it writes would be unkeyable. Fail closed at the
// one call site that knows what to do about it.
func Mint() (ID, error) {
	buf := make([]byte, idBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("session: 无法铸造会话身份（随机源故障），已 fail-closed: %w", err)
	}
	return ID(idPrefix + hex.EncodeToString(buf)), nil
}
