package main

// Ticket 248 AC#3/AC#8/AC#11 - the settings leg at the assembly root.
//
// WHAT THIS IS. internal/panel declares the socket (panel.ConfigStore) and refuses
// everything its field table does not name; this file is the half that owns the
// objects, built where both halves are reachable: the one config.Manager this
// process has, and the one secret.Store DPAPI owns. Same division ticket 114 set
// up for the permission mode (the panel requests, internal/perm writes), and the
// same reason: a page that could write configuration directly would be a page
// granting itself something.
//
// THE CREDENTIAL PATH, and the one rule it is made of. The value a user types into
// the settings page is bound by credentialWriteRequest below - a type that exists
// ONLY in this file, on a JSON key that panel.ComposerRequest does not declare. The
// shared envelope is decoded by every answered method, so anything readable from it
// is readable from the message route too; keeping this key off that struct is the
// whole of "the value never enters the shared envelope" (ticket 248's ruling J4).
// From the moment it is unmarshalled here the value is handed to the store and
// never carried again: not in a receipt, not in an error string, not in the audit
// line, not in a field of any type this package returns. What the page gets back is
// the reference the value now sits behind, which the schema itself writes into
// config.toml and which names nothing about the vendor.
//
// THE TIER, stated because AC#8 is a wording rule and not a comment. Every field
// this leg writes lives in [llm], and the model path a running process uses is built
// once at assembly (run.go) and never rebuilt by a reload - OnReload has no
// production assignment in this package. So every accepted write answers
// "restart required", even though the spec's table calls [llm] hot. The gap between
// those two facts is the reason this ticket forbids any "saved, already in force"
// wording: it would be the one sentence a user could act on and be misled by.
// Rebuilding that hop is not this ticket's to add.

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/secret"
)

// credentialValueKey is the ONE JSON key in this repository that carries a
// credential value. It is spelled here once, on purpose, so the test that asks
// "does the shared envelope declare a key this type binds?" has both sides to
// compare, and so a later reader cannot mistake the field below for a normal
// settings payload.
const credentialValueKey = "credentialValue"

// credentialWriteRequest is the write-only envelope of the credential field. Its
// two value-bearing keys are NOT on panel.ComposerRequest: the request that arrives
// on the message route decodes into that one shared type, and this type is only
// ever unmarshalled here, for this one field, once.
type credentialWriteRequest struct {
	Method         string `json:"method"`
	ConfigField    string `json:"configField"`
	ConfigProvider string `json:"configProvider"`
	// CredentialValue holds the secret between "read off the wire" and "handed to
	// DPAPI". Nothing below returns it, logs it or stores it in a field.
	CredentialValue string `json:"credentialValue"`
}

// configStore is the settings leg: the two objects it writes through, and the audit
// sink its refusals go to.
type configStore struct {
	mgr     *config.Manager
	secrets *secret.Store
	auditf  panel.AuditFunc
	actor   string
}

// newConfigStore wires the leg. Both halves come from the same assembly that built
// the permission-mode chain, so one process has one Manager and one Store - the
// "loading the file twice over is the split state" rule run.go states for the mode
// applies unchanged here.
func newConfigStore(mgr *config.Manager, secrets *secret.Store, auditf panel.AuditFunc, actor string) *configStore {
	return &configStore{mgr: mgr, secrets: secrets, auditf: auditf, actor: actor}
}

// ReadSettings answers "配了没配" with counts, references and states. The values a
// provider would need to answer "is the key actually there" are never returned:
// Existence is asked, content is not.
func (s *configStore) ReadSettings(_ context.Context) (panel.SettingsView, error) {
	if s == nil || s.mgr == nil {
		return panel.SettingsView{Credential: panel.CredentialUnknown}, fmt.Errorf("设置读写腿没有配置对象")
	}
	cfg := s.mgr.Config()
	if cfg == nil {
		return panel.SettingsView{Readable: false, Credential: panel.CredentialConfigUnreadable}, nil
	}
	view := panel.SettingsView{Readable: true, ChatModel: cfg.LLM.Roles.Chat.Model, ChatModelKnown: strings.TrimSpace(cfg.LLM.Roles.Chat.Model) != ""}

	names := make([]string, 0, len(cfg.LLM.Providers))
	for name := range cfg.LLM.Providers {
		names = append(names, name)
	}
	sort.Strings(names)

	var declared, missing int
	for _, name := range names {
		p := cfg.LLM.Providers[name]
		row := panel.ProviderSettings{
			Name:       name,
			APIKeyRef:  p.APIKeyRef,
			HasBaseURL: strings.TrimSpace(p.BaseURL) != "",
		}
		for range p.Models {
			row.ModelCount++
		}
		ref := strings.TrimSpace(p.APIKeyRef)
		if ref != "" {
			declared++
			present, err := s.refRecorded(ref)
			switch {
			case err != nil:
				row.CredentialError = err.Error()
				missing++
			case !present:
				missing++
			default:
				row.CredentialSet = true
			}
		}
		view.Providers = append(view.Providers, row)
	}
	view.Credential = aggregateCredential(declared, missing)
	return view, nil
}

// aggregateCredential is the read side's whole answer about keys, and it is a count
// comparison rather than a probe of any value.
func aggregateCredential(declared, missing int) panel.CredentialState {
	switch {
	case declared == 0:
		return panel.CredentialNoneDeclared
	case missing == 0:
		return panel.CredentialAllRecorded
	default:
		return panel.CredentialPartlyMissing
	}
}

// refRecorded asks the two reference families separately, because the store answers
// only one of them: Exists rejects env: refs on purpose (they name process state,
// not a blob file), and a single ruler that asked both would report a correct "not
// recorded" for the wrong reason.
func (s *configStore) refRecorded(ref string) (bool, error) {
	kind, value, err := secret.ParseRef(ref)
	if err != nil {
		return false, fmt.Errorf("引用形状不合法")
	}
	if kind == secret.RefKindEnv {
		// An env: ref names process state, and this process cannot fill the page's
		// operator's environment. Absent here is reported as a state, not as "the
		// key is missing everywhere".
		return strings.TrimSpace(os.Getenv(value)) != "", nil
	}
	if s.secrets == nil {
		return false, fmt.Errorf("凭据存储未接入")
	}
	return s.secrets.Exists(ref)
}

// ApplySetting maps one whitelisted field onto the config key it owns and writes
// that ONE key through the guarded path. The panel already refused any name this
// switch does not carry; the default branch stays because a field added to the
// panel's table without a leg behind it must not become a silent success.
func (s *configStore) ApplySetting(_ context.Context, w panel.SettingWrite) (panel.SettingWriteResult, error) {
	if s == nil || s.mgr == nil {
		return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied}, fmt.Errorf("设置写入腿没有配置对象，%q 未落盘", w.Field)
	}
	res := panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied}
	var (
		written []string
		err     error
	)
	switch w.Field {
	case panel.FieldProviderBaseURL:
		written, err = s.mgr.SetProviderBaseURL(w.Provider, w.Value)
	case panel.FieldProviderAPIKeyRef:
		ref := strings.TrimSpace(w.Value)
		if _, _, perr := secret.ParseRef(ref); perr != nil {
			// The schema keeps plaintext out structurally (D36 rule 5); validate()
			// is the second gate. This one fires first so the page is told which
			// field it got wrong instead of a TOML-shaped complaint.
			return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
				fmt.Errorf("引用形状不合法（只接受 dpapi:<id> 或 env:NAME）：字段 %q 未落盘", w.Field)
		}
		written, err = s.mgr.SetProviderAPIKeyRef(w.Provider, ref)
	case panel.FieldModelContextWindow:
		n, cerr := strconv.Atoi(strings.TrimSpace(w.Value))
		if cerr != nil {
			return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
				fmt.Errorf("字段 %q 要的是整数 token 数，收到的值解析不了（0 表示未知）", w.Field)
		}
		written, err = s.mgr.SetModelContextWindow(w.Provider, w.Model, n)
	case panel.FieldModelPriceIn:
		f, cerr := strconv.ParseFloat(strings.TrimSpace(w.Value), 64)
		if cerr != nil {
			return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
				fmt.Errorf("字段 %q 要的是每 1M token 的 micro-USD 数字，收到的值解析不了", w.Field)
		}
		written, err = s.mgr.SetModelPriceIn(w.Provider, w.Model, f)
	case panel.FieldModelPriceOut:
		f, cerr := strconv.ParseFloat(strings.TrimSpace(w.Value), 64)
		if cerr != nil {
			return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
				fmt.Errorf("字段 %q 要的是每 1M token 的 micro-USD 数字，收到的值解析不了", w.Field)
		}
		written, err = s.mgr.SetModelPriceOut(w.Provider, w.Model, f)
	case panel.FieldRoleChatModel:
		written, err = s.mgr.SetRoleChatModel(w.Value)
	case panel.FieldProviderCredential:
		return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
			fmt.Errorf("凭据值不走这条写入路径（它只写不回显，必须走 StoreCredential），字段 %q 未落盘", w.Field)
	default:
		return res, fmt.Errorf("字段 %q 没有对应的写入腿，未落盘", w.Field)
	}
	if err != nil {
		return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied}, err
	}
	res.Written = written
	res.Tier = panel.EffectiveRestart
	return res, nil
}

// StoreCredential is the write-only half of the settings path: bind the value off
// the key the shared envelope does not declare, hand it to the one store, record the
// reference it landed behind, and answer with status.
//
// Nothing returned from here carries the value, and no branch of it echoes what it
// was given - including the failures, which name the field and the reason.
func (s *configStore) StoreCredential(_ context.Context, raw string) (panel.SettingWriteResult, error) {
	res := panel.SettingWriteResult{Field: panel.FieldProviderCredential, Tier: panel.EffectiveNotApplied}
	if s == nil || s.mgr == nil || s.secrets == nil {
		return res, fmt.Errorf("凭据存储或配置对象未接入，密钥未录入")
	}
	var req credentialWriteRequest
	if err := json.Unmarshal([]byte(raw), &req); err != nil {
		return res, fmt.Errorf("凭据写入的封套解析不了（未录入任何密钥）：%v", err)
	}
	if strings.TrimSpace(req.ConfigField) != panel.FieldProviderCredential {
		return res, fmt.Errorf("封套里的字段不是 %q，凭据写入未被执行", panel.FieldProviderCredential)
	}
	provider := strings.TrimSpace(req.ConfigProvider)
	if provider == "" {
		return res, fmt.Errorf("凭据写入必须指名服务商，未录入任何密钥")
	}
	value := req.CredentialValue
	if value == "" {
		// The store refuses an empty secret on its own; saying it here keeps the
		// sentence about the field rather than about a file write.
		return res, fmt.Errorf("凭据字段是空的，未录入任何密钥")
	}
	// The reference is minted before the write so the blob file name and the ref
	// config will record are the same string; the id is random on purpose, which
	// is what makes a ref safe to show and a value not.
	ref := secret.NewRef()
	if err := s.secrets.Store(ref, value); err != nil {
		return res, fmt.Errorf("密钥没有写进凭据存储（配置未改动，未录入任何密钥）：%v", err)
	}
	written, err := s.mgr.SetProviderAPIKeyRef(provider, ref)
	if err != nil {
		// The blob exists now, so the honest sentence says both halves: the secret
		// is stored, the reference did not land. Cleaning the blob up silently
		// would be worse than reporting a state the user can finish by hand.
		return res, fmt.Errorf("密钥已进凭据存储，但 %q 的引用没写进配置：%w", provider, err)
	}
	res.Written = written
	res.Tier = panel.EffectiveRestart
	res.Credential = panel.CredentialSummary{Ref: ref, Recorded: true, Provider: provider}
	res.Note = "值本身不在配置里，也不在任何回话里。"
	return res, nil
}

// credentialStatus is the pump's reader for the snapshot's credential dimension
// (ticket 248 AC#2). It is a reader, not a value, for the same reason every other
// reader on that pump is one: a key the packet carries must be answered by something
// that actually looked.
func (s *configStore) credentialStatus() panel.CredentialState {
	view, err := s.ReadSettings(context.Background())
	if err != nil {
		return panel.CredentialConfigUnreadable
	}
	return view.Credential
}
