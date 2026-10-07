package panel

// The composer's request envelope (ticket 92 AC#1/AC#2/AC#3).
//
// There is ONE way for the page to ask for anything: WebView2's postMessage text
// pipe (C17), the same envelope the approval card already uses. This file does
// not open a second one - it gives the existing pipe a typed reader so the
// native side can (a) route a composer request to the one handler that owns it
// and (b) check who says they are.
//
// Why a source field is worth its bytes: the approval route and the composer
// route share a pipe, and "panel-composer" requests are the ones whose payload
// the panel controls completely (file bytes, a folder spelling, a mode name).
// A handler that accepted them from any sender could be fed by anything that can
// postMessage. The check is fail-closed and it is NOT an authority grant: after
// ParseComposerRequest the payload still goes through the same guards
// (ModeRequest.Parse / DecodeAttachmentPayload / RequestWorkspaceSwitch), so
// what the source field buys is a routing claim, nothing more.

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ComposerRequestSource is the only sender identity the composer route accepts.
const ComposerRequestSource = "panel-composer"

// Methods the composer route answers. Renaming one on either side goes red in
// TestTheRendererHoldsExactlyOneDoorToTheHost (composer_test.go:502), whose
// composerRouteLiterals() is built from these four constants and refuses any
// "panel.*" literal the Go side does not answer; that nail's own positive
// control is TestPlantedRendererDoorShapesGoRed (composer_test.go:533), which
// points the same scan at a knowingly wrong tree. (This pointer used to name a
// test that does not exist in this repository; ticket 35's snapshot pump was the
// step that checked it, and the behaviour was already covered - only the name
// was wrong.)
const (
	MethodModeRequest      = "panel.mode.request"
	MethodWorkspaceRequest = "panel.workspace.request"
	MethodAttachmentAdd    = "panel.attachment.add"
	MethodMessageSend      = "panel.message.send"
	// MethodConfigGet and MethodConfigSet are ticket 248 AC#1's two settings
	// routes. Three facts about the naming, because each one is a decision
	// somebody could later "fix" by accident:
	//
	//   - The names are NOT invented here. They are the two spellings the
	//     product spec already carries (docs/specs/SPEC-08-ui-ball-panel.md:168,
	//     read-only for this leg), so this file follows a document rather than
	//     minting a fourth vocabulary.
	//   - They carry no `panel.` prefix on purpose, and that is load-bearing
	//     twice over: ticket 181's anti-drift nail in this package counts the
	//     `panel.*` literals of this file (git_test.go's whitelistMethodsFromSource
	//     extractor), and the frozen grant-boundary instrument assembles its
	//     approval-shaped-name grid from `panel*` namespaces. Both readings were
	//     re-measured when this route landed; the census for it is in
	//     docs/evidence/s1/248-settings-write-path-r1.md.
	//   - Neither name contains any word of that instrument's verdict vocabulary,
	//     which is the whole reason "settings" is spelled `config.set` and not
	//     `config.approve`/`config.grant`: a route that lets the page change
	//     configuration is not a route that lets the page *approve* anything
	//     (AGENTS.md §1.2 keeps panel-side L2 "allow" off the table entirely).
	MethodConfigGet = "config.get"
	MethodConfigSet = "config.set"
)

// ErrComposerRequest is every way a composer request can fail to be usable.
var ErrComposerRequest = errors.New("panel: composer request refused")

// requestIDBytes is how much randomness a request id carries. It is a log
// correlation key, not a nonce and not authority, which is why it is not
// trusted for anything below.
const requestIDBytes = 8

// NewRequestID mints a correlation id for one request. A collision is a log
// annoyance, not a security event, so 64 random bits is generous here.
func NewRequestID() string {
	var b [requestIDBytes]byte
	if _, err := rand.Read(b[:]); err != nil {
		// No clock fallback: a request id that silently becomes a timestamp
		// would look like a sequencing guarantee nobody implemented.
		return "rid-unavailable"
	}
	return hex.EncodeToString(b[:])
}

// ComposerRequest is one parsed envelope from the composer.
type ComposerRequest struct {
	Method    string `json:"method"`
	RequestID string `json:"requestId"`
	Source    string `json:"source"`

	To   string `json:"to,omitempty"`
	Path string `json:"path,omitempty"`
	Text string `json:"text,omitempty"`
	AttachmentPayload
	Attachments []AttachmentRef `json:"attachments,omitempty"`
	// The four fields below are ticket 248's settings selectors, and they are
	// deliberately selectors plus one plain value - never a key path the page
	// names. `config.set` therefore cannot address a key this file has not
	// enumerated: the accepted spellings of configField live in
	// config_handlers.go's field table, and anything else is refused before the
	// write leg is ever reached. That narrowness is the ruling, not a
	// convenience - an arbitrary key/value pair on this envelope would put the
	// whole of config.toml, risk.permission_mode and fs.allowed_dirs included,
	// inside the panel's reach.
	//
	// What is NOT here is a credential VALUE. This envelope is shared by every
	// answered method, so a field added to it is readable by all of them; the
	// key a user types into the settings page must not live anywhere the message
	// route, the snapshot or the audit line can pick up. The credential write
	// hands its raw bytes to the leg that owns it and that leg decodes them into
	// a write-only shape of its own (cmd/wisp), uses them once, and keeps nothing.
	ConfigField    string `json:"configField,omitempty"`
	ConfigProvider string `json:"configProvider,omitempty"`
	ConfigModel    string `json:"configModel,omitempty"`
	ConfigValue    string `json:"configValue,omitempty"`
}

// ParseComposerRequest reads one postMessage string. Every refusal is an error
// wrapping ErrComposerRequest with the reason, so a host can tell the user their
// own request was dropped instead of looking like it worked.
func ParseComposerRequest(raw string) (ComposerRequest, error) {
	var r ComposerRequest
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		return r, fmt.Errorf("%w: 不是可解析的封套：%v", ErrComposerRequest, err)
	}
	r.Method = strings.TrimSpace(r.Method)
	if !knownComposerMethod(r.Method) {
		return r, fmt.Errorf("%w: 方法 %q 不是面板 composer 通路的能力入口", ErrComposerRequest, r.Method)
	}
	if strings.TrimSpace(r.Source) != ComposerRequestSource {
		return r, fmt.Errorf("%w: 来源 %q 不是 %q，按伪造/串台拒绝（requestId=%q）",
			ErrComposerRequest, r.Source, ComposerRequestSource, r.RequestID)
	}
	if strings.TrimSpace(r.RequestID) == "" {
		return r, fmt.Errorf("%w: 缺少 requestId，无法与审计/卡片对齐（method=%q）",
			ErrComposerRequest, r.Method)
	}
	return r, nil
}

func knownComposerMethod(m string) bool {
	switch m {
	case MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend, MethodConfigGet, MethodConfigSet:
		return true
	}
	return false
}

// RefusedEnvelopeForUser is the answer the host renders when a request never
// made it past ParseComposerRequest: it names the id so the user's own action is
// traceable, and it never degrades into a silent drop.
func RefusedEnvelopeForUser(r ComposerRequest, err error) string {
	id := r.RequestID
	if id == "" {
		id = "(无 requestId)"
	}
	method := r.Method
	if method == "" {
		method = "(无 method)"
	}
	return fmt.Sprintf("面板请求被拒绝 [%s %s]：%v", method, id, err)
}
