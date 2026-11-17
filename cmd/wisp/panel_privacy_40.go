package main

// Readface wave 1 (2026-11-17), ticket 40 AC#1's Go half: the three privacy
// doors get their real leg here. internal/panel names the three roster doors
// (bridge.go's privacy.delete-one / privacy.purge / privacy.export, routed by
// composer_dispatch.go) and this file attaches the handler that owns the
// objects: internal/memory's privacy layer (privacy.go's DeletePrivacyItem /
// PurgePrivacy / ExportPrivacy) - the layer SPEC-02 §5 names as the ONLY route
// privacy-relevant data may travel, the same layer ticket 40's census recorded
// as having zero production callers until now (ledger A365).
//
// SHAPE. The settings-route shape, for the same reason the config doors have
// it: a privacy receipt has no second voice to defer to, and the page-visible
// fact AC#1 measures is the count - "delete-one/purge/export e2e against live
// DB (counts match)". So the delete/purge receipts name the count the store
// returned, and the export carries the store's own indented JSON as one
// payload line behind its marker.
//
// THE DOMAIN IS THE SCOPE. There is no "purge everything" argument anywhere on
// this door: the domain selector is required, validated against the five D35
// domains internal/memory owns (privacy.go's PrivacyDomains - asked, never
// copied), and the destructive verb always carries it. An unknown domain is a
// named refusal that lists the five, never a guess.
//
// WHY OPEN PER READ. The same shape every read/write outlet in this package
// already uses (sessions.go's outlets, panel_cost_194.go): memory.Open, one
// operation, defer Close. A privacy request is a user-triggered, rare
// operation - not a hot path - and a handle held for the process's whole life
// would pin wisp.db open for a machine that may never ask. Per-read keeps
// every failure on the request that actually hit it: a database that will not
// open refuses THAT request loudly, audited, and the next one starts clean.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/panel"
)

// PrivacyExportPayloadMarker starts the export door's payload line: the JSON
// the page parses (memory.Store.ExportPrivacy's own envelope: domain,
// exported_at, count, records). Same naming properties as every marker in the
// family: no `panel.` prefix and no dot - it is neither a method name nor a
// route.
const PrivacyExportPayloadMarker = "wisp-privacy-export:"

// privacyHandler40 is the privacy leg: the data root it opens wisp.db over,
// the audit sink its refusals go to, and the actor for the audit line.
type privacyHandler40 struct {
	dir       string
	auditf    panel.AuditFunc
	actorName string
}

// newPrivacyHandler40 builds the leg over the data root the chain was given.
// No database is opened here - see the file comment.
func newPrivacyHandler40(dir string, auditf panel.AuditFunc, actor string) *privacyHandler40 {
	return &privacyHandler40{dir: dir, auditf: auditf, actorName: actor}
}

// validDomain parses the envelope's domain selector against internal/memory's
// own vocabulary. The word list is ASKED of the layer that owns it
// (memory.PrivacyDomains), never copied here: a sixth domain the schema gains
// tomorrow is valid the moment the layer says so.
func validPrivacyDomain(word string) (memory.PrivacyDomain, bool) {
	word = strings.TrimSpace(word)
	for _, d := range memory.PrivacyDomains() {
		if string(d) == word {
			return d, true
		}
	}
	return "", false
}

// openStore opens the per-read handle. The error is the caller's refusal
// material, with the root named.
func (h *privacyHandler40) openStore() (*memory.Store, error) {
	return memory.Open(h.dir)
}

// HandlePrivacyDeleteOne answers privacy.delete-one: one item out of one
// domain, the receipt naming the id. The handler refuses the two missing
// selectors by name before the leg is reached; the store refuses an id that
// does not exist ("the setter refuses by name when the entry is not there"),
// and that refusal is this receipt's raw material.
func (h *privacyHandler40) HandlePrivacyDeleteOne(ctx context.Context, req panel.ComposerRequest) (string, error) {
	domain, ok := validPrivacyDomain(req.PrivacyDomain)
	if !ok {
		err := h.refuse(req, "delete-one", fmt.Errorf(
			"privacyDomain %q 不是 D35 隐私域；可用的只有 %v", req.PrivacyDomain, memory.PrivacyDomains()))
		return "", err
	}
	if strings.TrimSpace(req.PrivacyItemID) == "" {
		err := h.refuse(req, "delete-one", fmt.Errorf(
			"privacyItemId 是空的：delete-one 只删一条，不给它 id 就什么都不删"))
		return "", err
	}
	st, err := h.openStore()
	if err != nil {
		err := h.refuse(req, "delete-one", fmt.Errorf("隐私存储不可用（%s）：%w", h.dir, err))
		return "", err
	}
	defer func() { _ = st.Close() }()
	if err := st.DeletePrivacyItem(ctx, domain, strings.TrimSpace(req.PrivacyItemID)); err != nil {
		err := h.refuse(req, "delete-one", err)
		return "", err
	}
	h.record(req, "delete-one", "deleted")
	return fmt.Sprintf("已删除 %s 域的 1 条记录（id=%s）。该域其余记录未动。", domain, req.PrivacyItemID), nil
}

// HandlePrivacyPurge answers privacy.purge: ONE whole domain out, the receipt
// naming the count the store returned. The domain selector is required and
// validated before the leg is reached - the destructive verb is always
// scoped, never "all of them".
func (h *privacyHandler40) HandlePrivacyPurge(ctx context.Context, req panel.ComposerRequest) (string, error) {
	domain, ok := validPrivacyDomain(req.PrivacyDomain)
	if !ok {
		err := h.refuse(req, "purge", fmt.Errorf(
			"privacyDomain %q 不是 D35 隐私域；purge 只按域执行，可用的只有 %v", req.PrivacyDomain, memory.PrivacyDomains()))
		return "", err
	}
	st, err := h.openStore()
	if err != nil {
		err := h.refuse(req, "purge", fmt.Errorf("隐私存储不可用（%s）：%w", h.dir, err))
		return "", err
	}
	defer func() { _ = st.Close() }()
	n, err := st.PurgePrivacy(ctx, domain)
	if err != nil {
		err := h.refuse(req, "purge", err)
		return "", err
	}
	h.record(req, "purge", "purged")
	return fmt.Sprintf("已清空 %s 域，共删除 %d 条记录。其他域未动。", domain, n), nil
}

// HandlePrivacyExport answers privacy.export: one domain as indented JSON,
// carried as the settings-route shape's payload line behind
// PrivacyExportPayloadMarker. The records are the store's own rows (the
// export is the user's data, not a summary) and the count in the receipt is
// the envelope's own count, read back off the JSON the leg hands out.
func (h *privacyHandler40) HandlePrivacyExport(ctx context.Context, req panel.ComposerRequest) (string, error) {
	domain, ok := validPrivacyDomain(req.PrivacyDomain)
	if !ok {
		err := h.refuse(req, "export", fmt.Errorf(
			"privacyDomain %q 不是 D35 隐私域；可用的只有 %v", req.PrivacyDomain, memory.PrivacyDomains()))
		return "", err
	}
	st, err := h.openStore()
	if err != nil {
		err := h.refuse(req, "export", fmt.Errorf("隐私存储不可用（%s）：%w", h.dir, err))
		return "", err
	}
	defer func() { _ = st.Close() }()
	data, err := st.ExportPrivacy(ctx, domain)
	if err != nil {
		err := h.refuse(req, "export", err)
		return "", err
	}
	// The count the receipt states is read off the SAME bytes the payload
	// carries: the store's envelope records the row list (and internal/memory
	// leaves its Count field unfilled today), so the receipt states
	// len(records) - the export's own row count, not a second count from a
	// second query.
	var envelope struct {
		Count   int             `json:"count"`
		Records json.RawMessage `json:"records"`
	}
	count := -1
	if jerr := json.Unmarshal(data, &envelope); jerr == nil {
		var rows []json.RawMessage
		if rerr := json.Unmarshal(envelope.Records, &rows); rerr == nil {
			count = len(rows)
		}
	}
	// The payload line is ONE line in this family (the page finds the marker
	// and parses the rest of the line as JSON), while ExportPrivacy renders
	// indented JSON for its file/bundle shape. The wire half therefore
	// re-renders the SAME envelope compact: a re-marshal of the decoded JSON,
	// never a re-count or a re-shape - every field the store exported is on
	// the wire verbatim, only whitespace is dropped.
	var compact map[string]any
	if jerr := json.Unmarshal(data, &compact); jerr == nil {
		if b, merr := json.Marshal(compact); merr == nil {
			data = b
		}
	}
	h.record(req, "export", "exported")
	return fmt.Sprintf("已导出 %s 域：%d 条记录，完整行内数据在下面这一行 JSON（导出即用户自己的数据，不是摘要）。",
		domain, count) + "\n" + PrivacyExportPayloadMarker + string(data), nil
}

// refuse records and returns one refusal; a nil audit sink silences the line,
// never the refusal.
func (h *privacyHandler40) refuse(req panel.ComposerRequest, verb string, cause error) error {
	err := fmt.Errorf("panel: privacy.%s 被拒绝（requestId=%q）：%w", verb, req.RequestID, cause)
	h.record(req, verb, "refused")
	return err
}

func (h *privacyHandler40) actor() string {
	if h == nil || h.actorName == "" {
		return "unknown"
	}
	return h.actorName
}

// record writes one audit line per outcome; it never carries item content -
// the domain word, the verb and the outcome are what an audit line needs,
// and a privacy request's payload is the user's data.
func (h *privacyHandler40) record(req panel.ComposerRequest, verb, outcome string) {
	if h == nil || h.auditf == nil {
		return
	}
	h.auditf("panel: PRIVACY-ROUTE request=%q method=%q verb=%q domain=%q actor=%q outcome=%q",
		req.RequestID, req.Method, verb, req.PrivacyDomain, h.actor(), outcome)
}
