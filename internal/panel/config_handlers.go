package panel

// Ticket 248 AC#1/AC#3/AC#8 - the settings door's own handler.
//
// WHAT THIS IS. bridge.go answers two new names (config.get / config.set) and
// composer_dispatch.go routes them here. This file is the caller-side half of the
// settings path: it decides what a settings request may name, refuses everything
// it may not, and hands the rest to the one leg that owns configuration.
//
// THE FIELD TABLE IS A WHITELIST, NOT A PARSER. A settings route that accepted a
// key path would hand the page the whole of config.toml, and the two families
// where that is not a convenience but the exact thing D33 forbids
// (permission_mode, the allowed-directory widening) live in that file too. So the
// accepted spellings are enumerated below, one per row, and anything else - a
// typo, a dotted path, a name from a locked family - is refused by name before
// the leg is reached. Growing this table is a contract change, not a bug fix.
//
// NO VALUE OF A CREDENTIAL EVER PASSES THROUGH HERE. The shared ComposerRequest
// (bridge.go) is decoded by every answered method, so a field added to it is
// readable by the message route as well, and this package's frozen instrument
// only screens names against its approval vocabulary - it has no word that stops
// a value-shaped key. The credential write therefore travels as the raw envelope
// and the leg that stores it decodes it into a write-only shape of its own
// (cmd/wisp), which is the only reason ConfigRequestHandler takes `raw` at all.
// What this file returns for that request is a status sentence, never an echo.
//
// DECLARED LOCALLY, NO NEW IMPORT EDGE. Same shape as ModeWriter above the
// handler in composer_handlers.go: the leg is an interface this package names and
// the assembly root fills, so internal/panel does not start importing
// internal/config or internal/secret for one ticket's route.
//
// WHAT THIS DOES NOT DO: it opens no second channel, it grants nothing, it cannot
// change the permission mode (J5 refuses that family outright - the mode goes
// through panel.mode.request, which costs the L2 card), and it invents no answer
// when the leg is absent: a nil leg is a refusal with a name.

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ErrNoConfigStore is the refusal for a settings request on a machine with no
// settings leg attached. It is the settings twin of ErrNoModeWriter and it exists
// so "the page asked and nothing happened" is never a silence.
var ErrNoConfigStore = errors.New("panel: 配置读写腿未接入")

// ErrConfigFieldRefused is every way one settings write can be refused before it
// reaches the leg: an unknown field, a missing selector, a value this route does
// not accept, or a name from a family the panel may not write.
var ErrConfigFieldRefused = errors.New("panel: 配置字段不在可写名册内")

// The writable fields. Each one names exactly one existing config.toml leaf, and
// the selector columns say which of the envelope's selector fields it needs.
const (
	FieldProviderBaseURL    = "provider_base_url"
	FieldProviderAPIKeyRef  = "provider_api_key_ref"
	FieldModelContextWindow = "model_context_window"
	FieldModelPriceIn       = "model_price_in"
	FieldModelPriceOut      = "model_price_out"
	FieldRoleChatModel      = "role_chat_model"
	// FieldProviderCredential is the write-only one: its value is not on the shared
	// envelope at all, so the leg reads it off the raw bytes and this route can only
	// ever answer with a status. `provider_api_key_ref` is the field that records
	// which reference the value ended up behind.
	FieldProviderCredential = "provider_credential"
)

// fieldNeedsProvider / fieldNeedsModel / fieldTakesValue spell the field table's
// shape as data, so the checks below are one loop over a table rather than six
// hand-written branches that can drift from it.
var (
	fieldNeedsProvider = map[string]bool{
		FieldProviderBaseURL: true, FieldProviderAPIKeyRef: true,
		FieldModelContextWindow: true, FieldModelPriceIn: true,
		FieldModelPriceOut: true, FieldProviderCredential: true,
	}
	fieldNeedsModel = map[string]bool{
		FieldModelContextWindow: true, FieldModelPriceIn: true, FieldModelPriceOut: true,
	}
	fieldTakesValue = map[string]bool{
		FieldProviderBaseURL: true, FieldProviderAPIKeyRef: true,
		FieldModelContextWindow: true, FieldModelPriceIn: true,
		FieldModelPriceOut: true, FieldRoleChatModel: true,
	}
	// allWritableFields is the whole name册 in one place, for the refusal text and
	// for the test that asks whether a locked-family name was quietly added.
	allWritableFields = []string{
		FieldProviderBaseURL, FieldProviderAPIKeyRef, FieldModelContextWindow,
		FieldModelPriceIn, FieldModelPriceOut, FieldRoleChatModel, FieldProviderCredential,
	}
)

// WritableFields returns the settings fields this route accepts, sorted. A test
// reads the table off here rather than copying it, so a field added to the switch
// and not to the table cannot stay invisible.
func WritableFields() []string {
	out := append([]string(nil), allWritableFields...)
	sort.Strings(out)
	return out
}

// lockedFieldFamilies are the config.toml families a settings write may never
// name, whatever else changes. The refusal says which door each one owns instead
// of leaving the page to guess.
var lockedFieldFamilies = []string{"risk.", "fs.", "net.", "plugins.", "privacy.", "models.", "audio."}

// refusedLockedFamily reports whether a field name reaches into a locked family,
// either by carrying the dotted path itself or by spelling one of the two leaves
// this ticket's ruling names (permission mode, the allowed-directory widening).
// Every reason it returns names the door that key does live behind, because a flat
// "not writable" is a sentence the page cannot act on.
func refusedLockedFamily(field string) (family string, why string) {
	lower := strings.ToLower(strings.TrimSpace(field))
	for _, f := range lockedFieldFamilies {
		if strings.HasPrefix(lower, f) {
			switch f {
			case "risk.":
				return f, "档位与放宽类键只由带原生二次确认的那条路写：权限档位走 panel.mode.request（变宽要那张 L2 卡）"
			case "fs.":
				return f, "可读写根的放宽欠同一张 L2 卡，由那条确认腿落盘，不走设置页的写入腿"
			default:
				return f, "该段属于安全锁定族，放宽必须走带原生二次确认的那条路，不走设置页"
			}
		}
	}
	switch lower {
	case "permission_mode", "permissionmode", "mode":
		return "risk.", "权限档位的变更只由 panel.mode.request 发起（变宽要那张 L2 卡），设置页不写它"
	case "allowed_dirs", "alloweddirs":
		return "fs.", "放宽可读写目录的根欠同一张 L2 卡，不由设置页的写入腿落盘"
	}
	return "", ""
}

// SettingsView is what the read side answers to "is it configured": counts,
// references and enum states. It carries no credential value and no
// partially-redacted value either - not the last four characters, which is the
// CLI terminal's tool (internal/secret's RedactSecret), not a page field.
type SettingsView struct {
	// Readable is the "配置是否可读" half: the leg read config.toml and got a
	// config back. False is a state of its own, never folded into "empty".
	Readable bool
	// Credential is the aggregate "凭据是否已录入" answer over the providers this
	// config declares. See CredentialState for what each value means.
	Credential CredentialState
	// Providers is one row per declared provider, name first so the page can
	// address a field without guessing an index.
	Providers []ProviderSettings
	// ChatModel is the model id the chat role resolves to right now; "" plus
	// ChatModelKnown=false means the leg never answered it.
	ChatModel      string
	ChatModelKnown bool
}

// ProviderSettings is one provider's row. APIKeyRef is a reference (dpapi:<id> /
// env:NAME), which the schema itself puts on disk and which names nothing about
// the vendor; it is the only credential-shaped string this type can hold, and
// internal/secret's own comment says a random ref does not leak which table it
// belongs to.
type ProviderSettings struct {
	Name          string
	APIKeyRef     string
	HasBaseURL    bool
	ModelCount    int
	CredentialSet bool
	// CredentialError says the blob could not be checked (an env: ref, or an
	// unreadable store). It is a state, not a swallowed failure.
	CredentialError string
}

// CredentialState is the read side's whole vocabulary about credentials.
type CredentialState string

const (
	// CredentialUnknown: no reader was attached, so nobody asked. Never rendered
	// as the safest-sounding option (ticket 92 AC#4's rule, inherited).
	CredentialUnknown CredentialState = "unknown"
	// CredentialConfigUnreadable: the config itself could not be read, so "is a
	// key recorded" has no answer at all.
	CredentialConfigUnreadable CredentialState = "config_unreadable"
	// CredentialNoneDeclared: the config reads fine and no provider names a ref.
	CredentialNoneDeclared CredentialState = "no_ref_declared"
	// CredentialAllRecorded: every declared ref resolves to something present.
	CredentialAllRecorded CredentialState = "all_recorded"
	// CredentialPartlyMissing: at least one declared ref is absent or could not
	// be checked. Stated as a count, never as which vendor.
	CredentialPartlyMissing CredentialState = "partly_missing"
)

// EffectiveTier is the "when does this take effect" answer every settings receipt
// carries (ticket 248 AC#8). The three names are the spec's three tiers; what
// this ticket adds is that a tier must be *said*, on the surface the user reads.
type EffectiveTier string

const (
	EffectiveNow        EffectiveTier = "now"
	EffectiveNextTask   EffectiveTier = "next_task"
	EffectiveRestart    EffectiveTier = "restart"
	EffectiveNotApplied EffectiveTier = "not_applied"
)

// SettingWriteResult is the leg's answer to one accepted write: which key paths
// actually changed on disk, and when they take effect. Written comes from
// config's own guarded merge writer, which can name the keys it replaced
// (internal/config's diffKeyPaths) - flattening this into "saved" is the failure
// AC#11's second clause names.
type SettingWriteResult struct {
	Field      string
	Written    []string
	Tier       EffectiveTier
	Err        error
	Note       string
	Credential CredentialSummary
}

// CredentialSummary is the credential write's status half: the reference the
// value now sits behind, and the blob's existence. Both are safe to show; the
// value itself is not carried by this type and cannot be.
type CredentialSummary struct {
	Ref        string
	Recorded   bool
	Provider   string
	ModelCount int
}

// ConfigStore is the injected settings leg. *cmd/wisp's configStore satisfies it
// the way *perm.Store satisfies ModeWriter, and for the same reason: the panel
// requests, the leg that owns the file owns the write.
type ConfigStore interface {
	// ReadSettings answers the "配了没配" question without returning any secret.
	ReadSettings(ctx context.Context) (SettingsView, error)
	// ApplySetting writes ONE named field through the guarded key-by-key write
	// path. It must refuse anything not in the table above, and it must validate
	// before it writes (AC#7).
	ApplySetting(ctx context.Context, w SettingWrite) (SettingWriteResult, error)
	// StoreCredential takes the raw envelope, binds the credential value off the
	// key the shared envelope does NOT declare, hands it to the one secret store,
	// and records the reference. It returns status only.
	StoreCredential(ctx context.Context, raw string) (SettingWriteResult, error)
}

// SettingWrite is one accepted, already-routed field write. The credential value
// is deliberately absent from this type as well: it never leaves the raw bytes.
type SettingWrite struct {
	Field     string
	Provider  string
	Model     string
	Value     string
	RequestID string
	Actor     string
}

// ConfigWriteHandler is the native handler for the two settings requests.
type ConfigWriteHandler struct {
	// Store is the single leg. nil = every settings request is refused.
	Store ConfigStore
	// Audit is the project's existing "[audit]" sink, as in every other handler.
	Audit AuditFunc
	// Actor names who triggered it, for the audit line. Empty stays "unknown".
	Actor string
}

// HandleConfigRequest answers one parsed settings request and returns the sentence
// the page is shown. Every refusal is returned as well as written to the audit,
// so a dropped request never looks like a handled one.
func (h *ConfigWriteHandler) HandleConfigRequest(ctx context.Context, req ComposerRequest, raw string) (string, error) {
	if h == nil || h.Store == nil {
		err := fmt.Errorf("%w: 设置请求 %q（requestId=%q）没有读写腿可用，未发生任何变更",
			ErrNoConfigStore, req.Method, req.RequestID)
		h.record(req, err, "no-leg")
		return "", err
	}
	switch req.Method {
	case MethodConfigGet:
		return h.read(ctx, req)
	case MethodConfigSet:
		return h.write(ctx, req, raw)
	default:
		err := fmt.Errorf("%w: %q 不是设置请求，本处理器不回答（requestId=%q）",
			ErrConfigFieldRefused, req.Method, req.RequestID)
		h.record(req, err, "wrong-method")
		return "", err
	}
}

// read answers config.get. A leg error is a state the user is told about, not an
// empty view: "nothing configured" and "nobody could read it" must not render the
// same way (the same rule pump.go's unreadable branches carry).
func (h *ConfigWriteHandler) read(ctx context.Context, req ComposerRequest) (string, error) {
	view, err := h.Store.ReadSettings(ctx)
	if err != nil {
		refused := fmt.Errorf("panel: 配置读取失败（requestId=%q）：%w", req.RequestID, err)
		h.record(req, refused, "read-failed")
		return "", refused
	}
	if !view.Readable {
		h.record(req, fmt.Errorf("panel: 配置不可读（requestId=%q），已按状态回报而不是空值", req.RequestID), "unreadable")
	}
	return renderSettingsView(view), nil
}

// write is config.set's gate: name the field, check it against the table, refuse
// locked families by name, and only then hand anything to the leg.
func (h *ConfigWriteHandler) write(ctx context.Context, req ComposerRequest, raw string) (string, error) {
	field := strings.TrimSpace(req.ConfigField)
	if field == "" {
		err := fmt.Errorf("%w: 设置写入必须点名一枚字段；可写的只有 %s（requestId=%q）",
			ErrConfigFieldRefused, strings.Join(WritableFields(), "/"), req.RequestID)
		h.record(req, err, "no-field")
		return "", err
	}
	if family, why := refusedLockedFamily(field); family != "" {
		err := fmt.Errorf("%w: %q 属于 %s 族，本路由不写它：%s（requestId=%q）",
			ErrConfigFieldRefused, field, family, why, req.RequestID)
		h.record(req, err, "locked-family")
		return "", err
	}
	if !knownWritableField(field) {
		err := fmt.Errorf("%w: %q 不在可写名册内；可写的只有 %s（requestId=%q）",
			ErrConfigFieldRefused, field, strings.Join(WritableFields(), "/"), req.RequestID)
		h.record(req, err, "unlisted-field")
		return "", err
	}
	provider := strings.TrimSpace(req.ConfigProvider)
	if fieldNeedsProvider[field] && provider == "" {
		err := fmt.Errorf("%w: 字段 %q 必须指名服务商（requestId=%q）", ErrConfigFieldRefused, field, req.RequestID)
		h.record(req, err, "no-provider")
		return "", err
	}
	model := strings.TrimSpace(req.ConfigModel)
	if fieldNeedsModel[field] && model == "" {
		err := fmt.Errorf("%w: 字段 %q 必须指名目录里的模型（requestId=%q）", ErrConfigFieldRefused, field, req.RequestID)
		h.record(req, err, "no-model")
		return "", err
	}
	if field == FieldProviderCredential {
		res, err := h.Store.StoreCredential(ctx, raw)
		if err != nil {
			h.record(req, err, "credential-refused")
			return "", err
		}
		h.record(req, nil, "credential-accepted")
		return renderCredentialReceipt(res), nil
	}
	value := req.ConfigValue
	if fieldTakesValue[field] && strings.TrimSpace(value) == "" {
		err := fmt.Errorf("%w: 字段 %q 需要一个值，空值不会被写成默认（requestId=%q）",
			ErrConfigFieldRefused, field, req.RequestID)
		h.record(req, err, "empty-value")
		return "", err
	}
	res, err := h.Store.ApplySetting(ctx, SettingWrite{
		Field: field, Provider: provider, Model: model, Value: value,
		RequestID: req.RequestID, Actor: h.actor(),
	})
	if err != nil {
		h.record(req, err, "write-refused")
		return "", err
	}
	h.record(req, nil, "write-accepted")
	return renderSettingReceipt(res), nil
}

// knownWritableField asks the table, not a copy of it.
func knownWritableField(field string) bool {
	for _, f := range allWritableFields {
		if f == field {
			return true
		}
	}
	return false
}

// renderSettingReceipt is the page-visible sentence for one accepted write: which
// key paths actually changed, and when they take effect. Two sentences are
// forbidden by AC#8 and neither is generated here - "saved" on its own (it hides
// the tier) and any wording that says a save takes effect in this run.
func renderSettingReceipt(r SettingWriteResult) string {
	var b strings.Builder
	b.WriteString("设置已写入。")
	if len(r.Written) == 0 {
		b.WriteString(" 本次没有任何键路径发生变化（文件里已经是这个值）。")
	} else {
		b.WriteString(" 实际落盘的键路径：")
		b.WriteString(strings.Join(r.Written, ", "))
		b.WriteString("。")
	}
	b.WriteString(" " + tierSentence(r.Tier))
	if r.Note != "" {
		b.WriteString(" " + r.Note)
	}
	return b.String()
}

// renderCredentialReceipt is the credential write's answer, and it says out loud
// that nothing is echoed: a settings page that shows the value back is the same
// leak as a log line that carries it.
func renderCredentialReceipt(r SettingWriteResult) string {
	var b strings.Builder
	b.WriteString("密钥已录入（只写不回显：这一页任何时候都不会再显示它的值，只会显示已录入/未录入）。")
	if r.Credential.Provider != "" {
		b.WriteString(" 服务商 ")
		b.WriteString(r.Credential.Provider)
		b.WriteString("。")
	}
	if r.Credential.Ref != "" {
		b.WriteString(" 记录的是引用 ")
		b.WriteString(r.Credential.Ref)
		b.WriteString("。")
	}
	if len(r.Written) > 0 {
		b.WriteString(" 实际落盘的键路径：")
		b.WriteString(strings.Join(r.Written, ", "))
		b.WriteString("。")
	}
	b.WriteString(" " + tierSentence(r.Tier))
	if r.Note != "" {
		b.WriteString(" " + r.Note)
	}
	return b.String()
}

// tierSentence is AC#8's wording rule in one place: restart-tier fields must be
// said as restart, and no branch here may claim a save is already in force.
func tierSentence(t EffectiveTier) string {
	switch t {
	case EffectiveNow:
		return "这一项立即生效。"
	case EffectiveNextTask:
		return "这一项在下一次任务开始时生效。"
	case EffectiveRestart:
		return "这一项要重启进程并重新运行才算用上（同一次运行里模型通路不会重建）。"
	case EffectiveNotApplied:
		return "这一项没有被应用。"
	default:
		return "生效级别未知，按需要重启处理。"
	}
}

// renderSettingsView turns the read-side facts into the sentence the settings page
// shows. Every field here is a state or a count; there is no value to leak.
func renderSettingsView(v SettingsView) string {
	var b strings.Builder
	if !v.Readable {
		b.WriteString("配置当前不可读，所以无法回答是否已录入密钥（这不是未配置）。")
		return b.String()
	}
	fmt.Fprintf(&b, "配置可读。服务商 %d 家；凭据状态 %s。", len(v.Providers), v.Credential)
	if v.ChatModelKnown {
		fmt.Fprintf(&b, " 当前聊天模型 %s。", v.ChatModel)
	} else {
		b.WriteString(" 当前聊天模型未读到。")
	}
	for _, p := range v.Providers {
		fmt.Fprintf(&b, "\n- %s：引用 %s，目录模型 %d 个，密钥 %s",
			p.Name, refOrNone(p.APIKeyRef), p.ModelCount, recordedWord(p))
	}
	return b.String()
}

func refOrNone(ref string) string {
	if strings.TrimSpace(ref) == "" {
		return "（无）"
	}
	return ref
}

func recordedWord(p ProviderSettings) string {
	switch {
	case strings.TrimSpace(p.APIKeyRef) == "":
		return "不适用"
	case p.CredentialError != "":
		return "无法核对：" + p.CredentialError
	case p.CredentialSet:
		return "已录入"
	default:
		return "未录入"
	}
}

func (h *ConfigWriteHandler) actor() string {
	if h == nil || h.Actor == "" {
		return "unknown"
	}
	return h.Actor
}

// record writes one audit line per outcome. It never carries a value: the fields
// it names are the method, the request id, the source and the outcome of the
// gate, which is everything an audit line needs and nothing a credential is made
// of. A nil sink silences the line, never the refusal.
func (h *ConfigWriteHandler) record(req ComposerRequest, err error, outcome string) {
	if h == nil || h.Audit == nil {
		return
	}
	line := "ok"
	if err != nil {
		line = err.Error()
	}
	h.Audit("panel: CONFIG-ROUTE request=%q method=%q source=%q actor=%q outcome=%q err=%s",
		req.RequestID, req.Method, req.Source, h.actor(), outcome, line)
}
