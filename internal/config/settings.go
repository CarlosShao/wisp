package config

// The settings-page write path (ticket 248 AC#3/AC#7/AC#11).
//
// WHY THESE METHODS EXIST AND SaveFile DOES NOT. Ticket 226 measured that
// SaveFile is a serializer, not a merge: it re-lays config.toml from whatever this
// process loaded at start-up, so any key the operator hand-edited while the process
// ran is silently reverted by the next programmatic write, and the write then
// claims the file's new mtime+size as its own. Every write this file makes goes
// through the same guarded key-by-key primitive the three existing setters use
// ((*Manager).mergeWrite), so a settings save cannot reopen the hole ticket 226
// closed. AC#11 forbids anything else, and mergeWrite's own body is ticket 226's
// to change, not this ticket's - which is why the report below is computed by
// diffing the file around the write instead of by editing that primitive.
//
// THE ORDER IS THE POINT (AC#7). validate() runs BEFORE mergeWrite, on the base the
// write would actually produce - the file's own content when the file can be read.
// That is not paranoia about a bad page: validate() is called at exactly two places
// in this package today (loader.go's readConfigFile and migrate.go), and neither is
// on the write path, so a settings route that skipped it would give one mistyped
// field the power to turn the NEXT start into an Unconfigured process. A refused
// write touches nothing: the file, and the in-memory config the runtime acts on,
// both stay byte-for-byte what they were.
//
// ONE KEY PER CALL (AC#11's second clause). Each method owns exactly one dotted key
// path and says so; a page that saves five fields makes five guarded writes and the
// receipt lists five paths. There is no multi-key atomic write here, because
// inventing one is ticket 226's territory and needs its own ruling.
//
// THREE REFUSALS, THREE SENTENCES (ticket 257 AC#2). A settings write can be
// refused for three different reasons and the page must never fold them into one
// "配置未生效" line (票 223 AC#4's discipline, same shape):
//
//	1. the file is not there at all  -> loader.go's readConfigFile says so by name;
//	   a settings write never creates config.toml.
//	2. the row the write addresses is not in the file -> the reason tags below,
//	   which is also where ticket 257's chosen shape (c) lives: the page does not
//	   mint a provider or a catalog row, and every refusal that is about a missing
//	   row names the one place a row does get created - the file itself.
//	3. the value would not survive validate() -> the pre-write gate below.
//
// Why the guidance spells [llm.providers.<名>] and never [providers.x]: a
// top-level [providers] section is an unknown key to parse.go's
// DisallowUnknownFields, so following that spelling would not unlock a row, it
// would stop the file from loading at all (measured by 257-a1, ledger A543).
//
// WHAT THIS FILE CANNOT DO: it offers no method for [risk] or [fs]. Permission mode
// changes go through the door that owes the L2 card (internal/perm plus the panel's
// panel.mode.request), and allowed-directory widening carries the same debt - see
// allowdirs.go's own header. The absence is the API: a caller cannot reach those
// keys through these setters, and no name lookup can be coerced into it.

import (
	"fmt"
	"os"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
)

// keyPathLLMProviderBaseURL / ... build the dotted key paths these setters own,
// spelled the way diffKeyPaths spells them so one key has one name in a log line,
// in a receipt and in a D36 verdict.
func keyPathProviderField(provider, leaf string) string {
	return "llm.providers." + provider + "." + leaf
}

func keyPathProviderModelField(provider, model, leaf string) string {
	return "llm.providers." + provider + ".models." + model + "." + leaf
}

// The three refusal reasons, each with its own tag (ticket 257 AC#2: 文件没建／行
// 不存在／校验不过 are three different things an operator fixes three different
// ways, and one shared line would be a lie by omission - the same rule 票 223 AC#4
// set for reload failures). Every refusal carries exactly one tag, so a caller -
// or a test - can tell which of the three it is looking at without parsing English.
const (
	refusalFileMissing = "第 1 种拒因：文件没建"
	refusalRowMissing  = "第 2 种拒因：行不存在"
	refusalInvalid     = "第 3 种拒因：校验不过"
)

// handAddGuidance is ticket 257 AC#0's chosen shape (c) in one sentence per row
// kind: the settings page does not create rows, and it says out loud where a row
// is created instead. The path is the one the schema actually reads.
const (
	guidanceProviderRow = "请在配置文件里添加一节 [llm.providers.<名>]（对上内置预设名的，protocol 与 base_url 可以留空；" +
		"非预设名必须自己写 protocol，否则这份文件加载不过）。"
	guidanceModelRow = "请在配置文件里添加这一节 [llm.providers.<名>.models.<模型 id>]，并写 enabled = true" +
		"（缺这枚键的条目视为关闭，点名它会在起动时被拒）。"
	guidanceRolePair = "请先在配置文件里添加一节 [llm.providers.<名>]，再在 [llm.roles.chat] 下把 provider 与 model " +
		"两枚一起点名（这一页写得了 model 那一半，写不了 provider 那一半，所以它替不成这一对）。"
)

// SetProviderBaseURL writes llm.providers.<name>.base_url for a provider this
// config already declares. An unknown provider is refused rather than created: a
// section that appears out of a settings page would carry no protocol, no billing
// and no catalog, and validate() would (correctly) reject the file on the next
// read - which is precisely the "one click bricks the next start" shape AC#7 names.
func (m *Manager) SetProviderBaseURL(provider, url string) ([]string, error) {
	u := strings.TrimSpace(url)
	return m.writeOneKey(keyPathProviderField(provider, "base_url"), func(base *Config) bool {
		p, ok := base.LLM.Providers[provider]
		if !ok {
			return false
		}
		p.BaseURL = u
		base.LLM.Providers[provider] = p
		return true
	}, func(base *Config) error {
		if _, ok := base.LLM.Providers[provider]; !ok {
			return unknownProviderErr(provider)
		}
		return nil
	})
}

// SetProviderAPIKeyRef writes llm.providers.<name>.api_key_ref. The value is a
// REFERENCE (dpapi:<blob-id> / env:NAME), never a key: validateAPIKeyRefs rejects
// anything else, so a plaintext value cannot land here even if a caller tried.
func (m *Manager) SetProviderAPIKeyRef(provider, ref string) ([]string, error) {
	r := strings.TrimSpace(ref)
	return m.writeOneKey(keyPathProviderField(provider, "api_key_ref"), func(base *Config) bool {
		p, ok := base.LLM.Providers[provider]
		if !ok {
			return false
		}
		p.APIKeyRef = r
		base.LLM.Providers[provider] = p
		return true
	}, func(base *Config) error {
		if _, ok := base.LLM.Providers[provider]; !ok {
			return unknownProviderErr(provider)
		}
		return nil
	})
}

// SetModelContextWindow writes llm.providers.<p>.models.<m>.context_window. The
// value has already been parsed to an int by the caller; a negative token count is
// refused here before validate() ever sees it, so the reason the page gets back
// names the field rather than a TOML parse error.
func (m *Manager) SetModelContextWindow(provider, model string, tokens int) ([]string, error) {
	if tokens < 0 {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config: llm.providers.%s.models.%s.context_window: %d is not a token count (0 means unknown)",
			provider, model, tokens))
	}
	return m.writeOneKey(keyPathProviderModelField(provider, model, "context_window"),
		func(base *Config) bool {
			p, ok := base.LLM.Providers[provider]
			if !ok {
				return false
			}
			spec, ok := p.Models[model]
			if !ok {
				return false
			}
			spec.ContextWindow = tokens
			p.Models[model] = spec
			base.LLM.Providers[provider] = p
			return true
		}, func(base *Config) error {
			return requireCatalogEntry(base, provider, model)
		})
}

// SetModelPriceIn / SetModelPriceOut write one leaf of a model's price card each
// (micro-USD per 1M tokens, schema.go's unit). They are separate methods on
// purpose: one key per guarded write is AC#11's second clause, and the five leaves
// of Price are not one key.
func (m *Manager) SetModelPriceIn(provider, model string, per1M float64) ([]string, error) {
	return m.setModelPriceLeaf(provider, model, "in", per1M)
}

// SetModelPriceOut writes llm.providers.<p>.models.<m>.price.out.
func (m *Manager) SetModelPriceOut(provider, model string, per1M float64) ([]string, error) {
	return m.setModelPriceLeaf(provider, model, "out", per1M)
}

func (m *Manager) setModelPriceLeaf(provider, model, leaf string, per1M float64) ([]string, error) {
	if per1M < 0 {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"config: llm.providers.%s.models.%s.price.%s: a negative price is not a price",
			provider, model, leaf))
	}
	return m.writeOneKey(keyPathProviderModelField(provider, model, "price."+leaf),
		func(base *Config) bool {
			p, ok := base.LLM.Providers[provider]
			if !ok {
				return false
			}
			spec, ok := p.Models[model]
			if !ok {
				return false
			}
			switch leaf {
			case "in":
				spec.Price.In = per1M
			case "out":
				spec.Price.Out = per1M
			default:
				return false
			}
			p.Models[model] = spec
			base.LLM.Providers[provider] = p
			return true
		}, func(base *Config) error {
			return requireCatalogEntry(base, provider, model)
		})
}

// SetRoleChatModel writes llm.roles.chat.model. The chat role's provider is not
// rewritten by this call, and validateCatalog is what says whether the resulting
// provider/model pair exists - which is exactly why this runs through the same
// pre-write validation as every other field here.
//
// The one exception, and it is ticket 257's: on a config that names no chat
// provider at all, the failure is reason 2 (the row is not there), not reason 3
// (a value that would not validate). Saying which one it is is the whole
// difference between "去文件里加一节" and a message that reads like the value was
// rejected. A config that DOES name a provider still goes through validate(), so
// an unknown pair keeps being reported as an unknown pair.
func (m *Manager) SetRoleChatModel(model string) ([]string, error) {
	v := strings.TrimSpace(model)
	return m.writeOneKey("llm.roles.chat.model", func(base *Config) bool {
		base.LLM.Roles.Chat.Model = v
		return true
	}, func(base *Config) error { return requireChatProvider(base, v) })
}

// requireChatProvider refuses the chat-model write when this config has no chat
// provider for the model to pair with (and none is named by the write - it only
// owns llm.roles.chat.model). Empty model = nothing to pair, so validate() decides.
func requireChatProvider(base *Config, model string) error {
	if base == nil || model == "" || base.LLM.Roles.Chat.Provider != "" {
		return nil
	}
	return observe.New(observe.ClassConfig,
		"config: llm.roles.chat names no provider, and this write sets llm.roles.chat.model only: "+
			"a settings write pairs a model with a provider this config already names, it does not name one. "+
			refusalRowMissing+"："+guidanceRolePair)
}

// writeOneKey is the shared body: pre-read the file (that snapshot is both the
// validation base and the report's baseline), apply, validate the base the write
// would actually produce, only then hand the change to the guarded merge writer,
// and finally report the key paths the FILE changed by.
//
// The three shapes mergeWrite refuses to be silent about are unchanged and this
// body adds no fourth: an unreadable file is still a refused write, a file that
// already carries the value is still not rewritten, and a write that preserved
// foreign hand edits is still not claimed as this process's own.
//
// apply reports false when the target key turned out not to exist on the base it
// was given (the file changed underneath between the pre-read and mergeWrite's own
// re-read). That is not an error and it is not a success to pretend about: nothing
// changes, the diff comes back empty, and the receipt says so.
func (m *Manager) writeOneKey(ownedKey string, apply func(base *Config) bool,
	check func(base *Config) error,
) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pre, preExists, err := m.readCurrentFile()
	if err != nil {
		// The file is here and this process cannot read it. Overwriting it is the
		// exact destruction ticket 226 was filed for, so the settings page gets the
		// refusal, not a write.
		return nil, err
	}
	candidate := deepCopyConfig(pre)
	if !preExists {
		candidate = deepCopyConfig(m.cur)
		if candidate == nil {
			candidate = NewDefaults()
		}
	}
	if check != nil {
		if err := check(candidate); err != nil {
			return nil, err
		}
	}
	if !apply(candidate) {
		return nil, observe.New(observe.ClassConfig,
			"config: "+ownedKey+" does not exist in the file this write is based on; nothing was written")
	}
	// AC#7's gate: the SAME validate() the start-up path runs, over the base the
	// write would produce, before a single byte is written.
	if err := validate(candidate); err != nil {
		return nil, observe.Wrap(observe.ClassConfig, err,
			"config: refused to write "+ownedKey+"; the value would not survive validation, so config.toml is unchanged。"+
				refusalInvalid+"：值本身过不了这份 schema 的校验，这一行要改的是值；文件一个字节都没动。")
	}

	// Memory first, then the file, then memory back if the file refused - the same
	// order and rollback the three older setters use.
	memOld := deepCopyConfig(m.cur)
	if memOld == nil {
		return nil, observe.New(observe.ClassConfig, "config: no loaded config in this Manager, refusing to write one")
	}
	apply(m.cur)
	if err := m.mergeWrite(ownedKey, func(base *Config) { apply(base) }); err != nil {
		m.cur = memOld
		return nil, observe.Wrap(observe.ClassConfig, err,
			"config: settings write for "+ownedKey+" failed; the values in memory are the ones this process started with")
	}

	return m.writtenKeyPaths(pre, preExists, ownedKey), nil
}

// writtenKeyPaths is the receipt's factual half: what the file actually changed by.
// It diffs the pre-write snapshot against the file as it now stands, so the answer
// is about the bytes on disk and not about what this process believes it wrote.
func (m *Manager) writtenKeyPaths(pre *Config, preExists bool, ownedKey string) []string {
	if !preExists {
		// There was no file: the write created one, and the key this setter owns is
		// the honest single answer rather than a diff against a snapshot that never
		// existed.
		return []string{ownedKey}
	}
	post, exists, err := m.readCurrentFile()
	if err != nil || !exists {
		// The report is a courtesy on top of a write that already happened; a file
		// that cannot be re-read now is said out loud, not papered over with the
		// key this setter expected.
		return []string{ownedKey + "（写后无法复读文件，键路径以本 setter 所持的键为准）"}
	}
	return diffKeyPaths(pre, post)
}

// readCurrentFile returns the watched file's content as a Config, plus whether a
// file is there at all. A missing file is not an error (first-run writes one);
// an unreadable one is, and the caller refuses the write on it.
func (m *Manager) readCurrentFile() (*Config, bool, error) {
	_, statErr := os.Stat(m.path)
	if statErr != nil {
		if fileMissing(statErr) {
			return nil, false, nil
		}
		return nil, false, observe.Wrap(observe.ClassConfig, statErr,
			"config: cannot read config.toml before writing it; nothing was written")
	}
	cfg, err := readConfigFile(m.path)
	if err != nil {
		return nil, true, observe.Wrap(observe.ClassConfig, err,
			"config: cannot read config.toml before writing it; nothing was written")
	}
	return cfg, true, nil
}

// requireCatalogEntry refuses a model-level write for a provider or a model this
// config does not already list, so the settings page addresses an existing entry
// instead of minting a catalog row. Both branches are reason 2, and both say where
// the row is created instead (ticket 257's shape (c)).
func requireCatalogEntry(base *Config, provider, model string) error {
	if base == nil {
		return observe.New(observe.ClassConfig, "config: no config to check a catalog entry against")
	}
	p, ok := base.LLM.Providers[provider]
	if !ok {
		return unknownProviderErr(provider)
	}
	if _, ok := p.Models[model]; !ok {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"config: llm.providers.%s has no model %q in its catalog; refusing to invent one from a settings write. "+
				"%s：%s", provider, model, refusalRowMissing, guidanceModelRow))
	}
	return nil
}

// unknownProviderErr is reason 2 for a provider-level key: the row this write
// would edit is not in the file. Ticket 257 AC#0 chose shape (c), so the refusal
// also names where a provider row does come from - the file - instead of leaving
// the operator to guess that the page cannot make one.
func unknownProviderErr(provider string) error {
	return observe.New(observe.ClassConfig, fmt.Sprintf(
		"config: no provider %q in this config; a settings write addresses an existing entry, it does not create one. "+
			"%s：这一页改不了服务商的存在性。%s", provider, refusalRowMissing, guidanceProviderRow))
}
