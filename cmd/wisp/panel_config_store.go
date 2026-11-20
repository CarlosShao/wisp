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

	"github.com/CarlosShao/wisp/internal/audio/aec"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/secret"
)

// credentialValueKey is the ONE JSON key in this repository that carries a
// credential value. It is spelled here once, on purpose, so the test that asks
// "does the shared envelope declare a key this type binds?" has both sides to
// compare, and so a later reader cannot mistake the field below for a normal
// settings payload.
const credentialValueKey = "credentialValue"

// keyPathPanelAnimation193 is panel.animation spelled the way internal/config's
// guarded write path names it (settings_animation_193.go's own
// keyPathPanelAnimation - unexported there, so this leg spells the same
// literal once, HERE, and the wiring test pins the receipt against it: a
// drift between the two spellings reddens the end-to-end tooth, never a
// silent wrong receipt). The tier lookup and the receipt both use it.
const keyPathPanelAnimation193 = "panel.animation"

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
	// 本机音频 half (the 2026-11-11 S-leg alignment): the three keys the page's
	// waiting card graduates on. The locals exist because the view outlives
	// this call - they snapshot the values off the per-call Config() read
	// instead of handing out addresses into the config object. This read is
	// also what cmd/wisp/config_readers_255.go's [audio] row now cites as the
	// section's live reader (the same per-call class the [llm] row rests on).
	aecOn := cfg.Voice.AEC.Enabled
	inputDevice := cfg.Audio.InputDevice
	sampleRate := cfg.Audio.SampleRate
	// 票 452 AC#2 的"界面可判"那一半. The AEC CAPABILITY bit, not a config value:
	// aec.BindingBuilt() is answered by the build shape itself (`-tags
	// webrtcap` present or not), so this leg can state "this binary cannot
	// cancel echo" no matter what config.toml says - and the shipped build has
	// no binding at all (scripts/build.ps1 passes no -tags). It is deliberately
	// its OWN field beside AecEnabled, never folded into it: one is "你想不想开"
	// (AecEnabled, writable), the other is "这枚二进制能不能" (this), and a page
	// that only had the first would print "AEC 开" on a build that cannot cancel
	// anything. The local exists for the same reason as the three above: the
	// view outlives this call, so it snapshots the answer instead of handing out
	// an address into anything per-call.
	aecBuildAvailable := aec.BindingBuilt()
	view := panel.SettingsView{
		Readable:       true,
		ChatModel:      cfg.LLM.Roles.Chat.Model,
		ChatModelKnown: strings.TrimSpace(cfg.LLM.Roles.Chat.Model) != "",
		// 票 313：选择器的数据源 ＋ "用户配过没有" 那一枚旗。
		Presets: config.PresetNames(), PresetsKnown: true,
		ChatProvider:    cfg.LLM.Roles.Chat.Provider,
		Configured:      len(cfg.LLM.Providers) > 0 || strings.TrimSpace(cfg.LLM.Roles.Chat.Provider) != "",
		ConfiguredKnown: true,
		Voice: panel.VoiceSettingsView{
			CloudASRChain:    cfg.Voice.CloudASRChain,
			CloudTTSChain:    cfg.Voice.CloudTTSChain,
			RealtimeEnabled:  cfg.Voice.Realtime.Enabled,
			RealtimeProvider: cfg.Voice.Realtime.Provider,
			RealtimeModel:    cfg.Voice.Realtime.Model,
			AecEnabled:       &aecOn,
			InputDevice:      &inputDevice,
			SampleRate:       &sampleRate,
			// 票 452 AC#2: the capability half, from the build shape.
			AecBuildAvailable: &aecBuildAvailable,
		},
		// Readface wave 1 (2026-11-17), ticket 39 AC#1/AC#4: the FULL D36
		// section tree (config.SectionRows - the registry's own tier words)
		// and the disk file's structured load errors (the Manager's own
		// read face, unknown keys / wrong types with their 1-based line).
		// The tree comes off the config object this read already made; the
		// errors re-read the FILE, because the Manager keeps running on its
		// last good config while the disk file is broken - exactly the state
		// the editor page must see.
		Sections:   configSectionRows456(cfg),
		LoadErrors: configLoadErrors456(s.mgr.ParseErrors()),
	}
	// Ticket 39 AC#4's drift handshake, read half: the page needs the disk
	// file's fingerprint to echo back on config.set. An unreadable file
	// leaves it "" - the key stays off the wire (absent = this leg did not
	// answer), and a page that never got a fingerprint cannot be accused of
	// holding a stale one.
	if sha, shaErr := s.mgr.FileSHA256(); shaErr == nil {
		view.BaseSHA = sha
	}

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
			// The row's identity, read back for the page (see the field notes in
			// panel.ProviderSettings): with two rows of one vendor on screen, the
			// base URL and the billing word are what tell them apart.
			BaseURL:  p.BaseURL,
			Protocol: p.Protocol,
			Billing:  p.Billing,
		}
		modelNames := make([]string, 0, len(p.Models))
		for id := range p.Models {
			modelNames = append(modelNames, id)
		}
		sort.Strings(modelNames)
		row.ModelCount = len(modelNames)
		row.Models = modelNames
		// Audit R1/P3's read face: each catalog model's two audio bits, read
		// off the same specs the names above come from. A model the voice
		// assembly refuses (cloudcatalog.go's lacks capabilities.audio_in/out)
		// is now visible to the page that let the chain be saved, and the
		// model page's switches have a real "当前" to show instead of a guess.
		if len(p.Models) > 0 {
			caps := make(map[string]panel.ModelCapabilityBits, len(p.Models))
			enabled := make(map[string]bool, len(p.Models))
			quotas := make(map[string]panel.ModelQuotaLimits, len(p.Models))
			for id, spec := range p.Models {
				caps[id] = panel.ModelCapabilityBits{
					AudioIn:  spec.Capabilities.AudioIn,
					AudioOut: spec.Capabilities.AudioOut,
					// 票 40 AC#3：realtime 角色下拉的过滤位。实时语音头是语音角色，
					// 它的过滤位此前没有读面（两个 audio 位自 R1/P3 起就在这张
					// 图上）；位来源＝模型行的真字段 spec.Capabilities.Realtime，
					// 不是页面造的名单。
					Realtime: spec.Capabilities.Realtime,
				}
				// R1-recheck residual 2's third gate: the model's enabled flag,
				// the same specs the names and bits above come from. The boot
				// refuses a chain element naming an enabled=false model
				// (internal/speech/cloudcatalog.go:87-90); without this read
				// face the page could only judge the gate by guessing.
				enabled[id] = spec.Enabled
				// 票 40 AC#3 的 quota 编辑半边读口：当前额度值。页面要改一行额度
				// 就得先知道它现在是多少——没有这个读面，编辑框就是蒙着眼覆盖。
				// 0 = 不封顶是真答案（schema 的语义），不折成"没答"。
				quotas[id] = panel.ModelQuotaLimits{
					DailyMicro:   spec.QuotaDailyMicro,
					MonthlyMicro: spec.QuotaMonthlyMicro,
				}
			}
			row.ModelCapabilities = caps
			row.ModelsEnabled = enabled
			row.ModelQuotas = quotas
		}
		row.PlanCreditTotalMicro = p.PlanCreditTotalMicro
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
	// Ticket 39 AC#4's drift handshake, write half: a page that holds a
	// fingerprint from its last config.get must not silently overwrite what
	// changed behind it. The echo is compared against the file's OWN hash
	// right now; a mismatch refuses with refused-drift before anything is
	// written, and the page is told the disk moved (it re-reads and decides
	// again). Empty echo = the page has no fingerprint to answer for (an
	// older page, a leg that never answered the read half) and the write
	// runs without this check - the pre-39 shape, unchanged.
	if base := strings.TrimSpace(w.BaseSHA); base != "" {
		now, shaErr := s.mgr.FileSHA256()
		switch {
		case shaErr != nil:
			return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
				observe.New(observe.ClassConfig,
					"config: refused-drift：config.toml 现在读不出来，无法核对盘上是否已变；这次保存一个键都没写。请先确认文件可读。")
		case now != base:
			return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
				observe.New(observe.ClassConfig,
					"config: refused-drift：盘上已变——config.toml 在这次读取之后又被改过（外部编辑或另一处写入），"+
						"这次保存拿的是旧视图，拒绝静默覆盖：一个键都没写。请重读配置确认盘上现值后再保存。")
		}
	}
	var (
		written []string
		err     error
		// The appearance wave's measured-tier inputs, captured by their cases
		// and read after the switch (the same shape FieldBallStyle uses for
		// its applier answer below).
		waveOpacity       float64
		waveClickOn       bool
		waveOpacityParsed bool
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
	case panel.FieldModelQuotaDaily, panel.FieldModelQuotaMonthly, panel.FieldProviderPlanCredit:
		// 票 40 AC#3 的 quota 编辑半边（后三半派单的第三半）：值＝micro-USD 的
		// 整数字符串。带的上限＝int64 自身（ParseInt 是它的执行者，越界解析
		// 即拒），下限＝0，负数由 internal/config 的 setter 具名拒绝（0 = 不
		// 封顶/未录入是合法值，schema 与 quota_44 的现量语义，这里不另造上限
		// 数字）；这里只替它答一句形状话（上面数值字段的同款）。档位沿用
		// [llm] 各写入腿的 restart 盖章（下面默认盖）。
		n, perr := strconv.ParseInt(strings.TrimSpace(w.Value), 10, 64)
		if perr != nil {
			return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
				fmt.Errorf("字段 %q 要的是 micro-USD 的整数（0 = 不封顶/未录入），收到的值解析不了", w.Field)
		}
		switch w.Field {
		case panel.FieldModelQuotaDaily:
			written, err = s.mgr.SetModelQuotaDaily(w.Provider, w.Model, n)
		case panel.FieldModelQuotaMonthly:
			written, err = s.mgr.SetModelQuotaMonthly(w.Provider, w.Model, n)
		default:
			written, err = s.mgr.SetProviderPlanCreditTotal(w.Provider, n)
		}
	case panel.FieldRoleChatModel:
		written, err = s.mgr.SetRoleChatModel(w.Value)
	case panel.FieldRoleChatThinking:
		// 票 187：思考档。值＝枚举词本身，词表（含当前模型自报的 levels）由
		// config 侧 setter 在落盘前把关（320 裁定的同一条形状）；档位是
		// next_task，面板门在 187 名册里按 fieldTier187 改判，这里照 [llm]
		// 各字段的同款先盖 restart 戳。
		written, err = s.mgr.SetRoleChatThinking(w.Value)
	case panel.FieldVoiceCloudASRChain:
		written, err = s.mgr.SetVoiceCloudASRChain(splitChainValue(w.Value))
	case panel.FieldVoiceCloudTTSChain:
		written, err = s.mgr.SetVoiceCloudTTSChain(splitChainValue(w.Value))
	case panel.FieldVoiceRealtimeProvider:
		written, err = s.mgr.SetVoiceRealtimeProvider(w.Value)
	case panel.FieldVoiceRealtimeModel:
		written, err = s.mgr.SetVoiceRealtimeModel(w.Value)
	case panel.FieldVoiceAecEnabled:
		// 本机音频（S 腿对齐，2026-11-11）：插话监测的配置开关。门已拒空值；
		// 非布尔词在这里替它答一句形状话（BallClickThrough 的同款），让页面
		// 知道是哪一枚字段给错了。更深的把关只有一条单源：writeOneKey 落盘前
		// 对"写完会是什么"跑 validate()，这枚键不另立第二套词表。档位沿用本
		// 店 [voice] 各写入腿的 restart 盖章（与两条云链写一致）。
		on, perr := strconv.ParseBool(strings.TrimSpace(w.Value))
		if perr != nil {
			return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
				fmt.Errorf("字段 %q 要的是 true/false，收到的值解析不了", w.Field)
		}
		written, err = s.mgr.SetVoiceAecEnabled(on)
	case panel.FieldAudioInputDevice:
		// 本机音频：采集设备名。空值已被门拒；"default"（和空串）是 schema 与
		// internal/audio 约定的「跟随系统默认」词。除此之外任何名字照写，这里
		// 不预判枚举——设备在不在场是开机后的事实，不是配置层的词表；打开时
		// 解析不到会按设备名响亮拒绝并列出在场端点（internal/audio 的
		// selectCaptureDevice），那才是这条链的拒绝单源。
		written, err = s.mgr.SetAudioInputDevice(strings.TrimSpace(w.Value))
	case panel.FieldProviderRowCreate:
		written, err = s.mgr.CreateProviderRow(w.Provider, w.Model)
	case panel.FieldProviderBilling:
		// 票 320 AC#1: the value is the enum word itself; the config setter owns
		// the vocabulary check so the page is told which field it got wrong.
		written, err = s.mgr.SetProviderBilling(w.Provider, w.Value)
	case panel.FieldProviderRowCreateCustom:
		// 票 320 AC#2: the wire value is one "协议,base_url" pair; the split is this
		// leg's (the panel's own gate only knows the field takes a value).
		protocol, baseURL := splitCustomRowValue(w.Value)
		written, err = s.mgr.CreateCustomProviderRow(w.Provider, protocol, baseURL, w.Model)
	case panel.FieldModelRowCreate:
		// 票 320 AC#3: provider and model are the two selectors, no value - the
		// entry's window and price leaves keep their own fields.
		written, err = s.mgr.CreateModelRow(w.Provider, w.Model)
	case panel.FieldModelAudioIn, panel.FieldModelAudioOut:
		// Audit R1/P3: the capability bits the voice assembly refuses on, now
		// writable from the page. Value = true/false; the leg's only job is
		// the word shape (the gate already refused an empty one) - the row
		// existence ruler and the guarded leaf write live in internal/config's
		// SetModelAudioIn/SetModelAudioOut, exactly like the price leaves.
		on, perr := strconv.ParseBool(strings.TrimSpace(w.Value))
		if perr != nil {
			return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
				fmt.Errorf("字段 %q 要的是 true/false，收到的值解析不了", w.Field)
		}
		if w.Field == panel.FieldModelAudioIn {
			written, err = s.mgr.SetModelAudioIn(w.Provider, w.Model, on)
		} else {
			written, err = s.mgr.SetModelAudioOut(w.Provider, w.Model, on)
		}
	case panel.FieldModelRowDelete:
		// The taking-back twin of the line above (owner follow-up): provider+model
		// deletes one catalog entry. The setter refuses by name when the entry is
		// not there, so a double click cannot pretend it removed something.
		written, err = s.mgr.RemoveModelRow(w.Provider, w.Model)
	case panel.FieldProviderRowDelete:
		// Deletes a whole row; the setter refuses by name while the chat role or a
		// voice chain still points at it.
		written, err = s.mgr.RemoveProviderRow(w.Provider)
	case panel.FieldRoleChatPair:
		written, err = s.mgr.SetRoleChatPair(w.Provider, w.Model)
	case panel.FieldBallStyle:
		// 票 65: 外观风格切换。值＝风格词本身（glass|jarvis），枚举由 config 侧
		// setter 在落盘前把关（320 裁定的同一条形状）；盘上的 [ball].style 键、
		// 校验与 glass 默认是票 65 实现腿（30d8f2e38）早已落的，这里只接写入腿。
		written, err = s.mgr.SetBallStyle(w.Value)
	case panel.FieldBallSize:
		// 票 413 契约 §2.2、由 422 接线：值＝整数 px 的字符串（"90"）。44..160
		// 的带由 config 侧 SetBallSize 先拒后写，拒绝句＝validate.go 的
		// ballSizeRangeError 单源句，与加载闸逐字同一串，永不分叉；这里只在值
		// 根本不是整数时替它答一句形状话（上面 222 行 [llm] 数值字段的同款）。
		// 档位静态答 restart（下面盖的默认档，无需改判）：size 没有 65 那种进程
		// 内 applier，球窗几何在建窗时从开机种子读定，写完内存热生效但可见直径
		// 要等下次建窗——像 187 的静态标法，不像 ball_style 要实测。
		v, perr := strconv.Atoi(strings.TrimSpace(w.Value))
		if perr != nil {
			return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
				fmt.Errorf("字段 %q 要的是 44..160 的整数（px），收到的值解析不了", w.Field)
		}
		written, err = s.mgr.SetBallSize(v)
	case panel.FieldBallOpacityIdle:
		// 外观波：空闲不透明度。值＝0..1 的浮点字符串（"0.9"）。带由 config 侧
		// SetBallOpacityIdle 先拒后写，拒绝句＝validate.go 的
		// ballOpacityIdleRangeError 单源句，与加载闸逐字同一串，永不分叉；
		// 这里只在值根本不是浮点数时替它答一句形状话（上面数值字段的同款）。
		// 档位是实测的（下面）：本进程有球＝applyBallIdleOpacity 把空闲观感换
		// 过去，答 now；无球＝下次启动的种子读取会落到它，答 restart。
		v, perr := strconv.ParseFloat(strings.TrimSpace(w.Value), 64)
		if perr != nil {
			return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
				fmt.Errorf("字段 %q 要的是 0..1 的数字，收到的值解析不了", w.Field)
		}
		waveOpacity, waveOpacityParsed = v, true
		written, err = s.mgr.SetBallOpacityIdle(v)
	case panel.FieldBallClickThrough:
		// 外观波：空闲鼠标穿透开关。值＝布尔字面量（"true"/"false"）；布尔没有
		// 带可单源，config 侧 setter 直接走受守写路径。档位同上实测：有球＝
		// WS_EX_TRANSPARENT 位当场翻，答 now；无球＝restart。
		on, perr := strconv.ParseBool(strings.TrimSpace(w.Value))
		if perr != nil {
			return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
				fmt.Errorf("字段 %q 要的是 true/false，收到的值解析不了", w.Field)
		}
		waveClickOn = on
		written, err = s.mgr.SetBallClickThrough(on)
	case panel.FieldPanelWidth:
		// 外观波：面板宽度。值＝整数 px 的字符串。config 侧 SetPanelWidth 拒
		// 非正值；无 live applier（面板宿主每次建窗现读，写完内存热生效但可见
		// 尺寸要等下一次显示请求——像 ball_size 的静态标法，但读者不是重启才
		// 到：收据的 Note 把「关窗再开即跟上」写明，档位仍答 restart，这是词表
		// 里最接近而不说谎的一枚）。
		v, perr := strconv.Atoi(strings.TrimSpace(w.Value))
		if perr != nil {
			return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
				fmt.Errorf("字段 %q 要的是正整数（px），收到的值解析不了", w.Field)
		}
		written, err = s.mgr.SetPanelWidth(v)
	case panel.FieldPanelFontSize:
		// 外观波：面板正文字号（panel.font_size，本波 schema 新键，默认 12＝
		// 面板文档正文今天的事实字号，默认下接线不动观感）。同 width 的形状：
		// 页面侧消费随前端接线腿落地；写路由先通，档位与 Note 同款。
		v, perr := strconv.Atoi(strings.TrimSpace(w.Value))
		if perr != nil {
			return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
				fmt.Errorf("字段 %q 要的是正整数（px 字号），收到的值解析不了", w.Field)
		}
		written, err = s.mgr.SetPanelFontSize(v)
	case panel.FieldPanelAnimation:
		// 票 193 AC#3 接线第三步：动画档（panel.animation，本票 schema 新键，
		// 默认 reduce＝owner 待拍占位）。值＝档位词本身（full|reduce|off），
		// 词表只有一条单源：config 侧 SetPanelAnimation 落盘前把关，拒绝句＝
		// validate.go 的 panelAnimationEnumError——与加载闸逐字同一串，这里
		// 不另立第二套词表（320/65/413 的同一条形状）。档位不在这里盖：
		// 下面由 registryEffectiveTier39 问 TierRegistry 答（[panel] 整节行，
		// hot——D36 的 [panel] 行钉的就是这一档）。
		written, err = s.mgr.SetPanelAnimation(strings.TrimSpace(w.Value))
	case panel.FieldPluginEnabled:
		// 票 215a：插件开关。值＝"true"/"false"，插件 id 在 w.Plugin（门已
		// 拒空值）；布尔形状在这里替它答一句形状话（上面布尔字段的同款）。
		// 写入走 internal/config 的受保护分节写（SetPluginEnabled over
		// WriteSection，票 248/195 的现成写路径，⛔ 不另造新路）：关＝收紧，
		// 立即落盘；开＝锁定段放宽，面板来源被结构性拒绝（R20：面板是权限
		// 输入口，不是权限当局）——拒绝句原样上屏。审计两处都落：这一层的
		// CONFIG-ROUTE 行＋internal/config 的 SECTION-WRITE／PLUGIN-TOGGLE 行。
		on, perr := strconv.ParseBool(strings.TrimSpace(w.Value))
		if perr != nil {
			return panel.SettingWriteResult{Field: w.Field, Tier: panel.EffectiveNotApplied},
				fmt.Errorf("字段 %q 要的是 true/false，收到的值解析不了", w.Field)
		}
		written, err = s.mgr.SetPluginEnabled(w.Plugin, on, config.SourcePanel, "panel.config.set")
	case panel.FieldAppLanguage:
		// 票 39 AC#1：重启档代表键。值＝界面语言标签的字符串（schema 无枚举，
		// 门已拒空值）；写入走受守路径 SetAppLanguage。档位不在这里盖——下面
		// 由 registryEffectiveTier39 问 TierRegistry 答（app.language 的
		// per-key restart 行），行为半（只落盘、内存保开机值）在
		// internal/config 的 tier-aware 写路径里。
		written, err = s.mgr.SetAppLanguage(strings.TrimSpace(w.Value))
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
	// Ticket 39 AC#4: the reload event advances the page's fingerprint to
	// the file as THIS write left it, so back-to-back saves never trip the
	// drift check against their own predecessor. A hash failure here stays
	// silent at the wire level ("" = absent key) - the page's next read
	// refreshes the cache instead.
	if sha, shaErr := s.mgr.FileSHA256(); shaErr == nil {
		res.BaseSHAAfter = sha
	}
	res.Tier = panel.EffectiveRestart
	if w.Field == panel.FieldAppLanguage {
		// 票 39 AC#1：档位与注册表同源。app.language 的档是 tiers.go 的
		// per-key 行答的（TierOf），不是这里硬写；本字段是这张映射的第一枚
		// 生产消费者，旧字段的实测档位戳（187/65/422 的名册）刻意不动。
		if len(written) > 0 {
			if t, ok := registryEffectiveTier39(written[0]); ok {
				res.Tier = t
			}
		}
	}
	if w.Field == panel.FieldPluginEnabled {
		// 票 215a 的档位与真话：配置层立即生效（[plugins] 锁定段的收紧在
		// manager.go 走同一条「hot-applied」路），但 Note 必须把「谁消费它」
		// 写明——插件引擎未在装配根接线（票 50/51 组件已在库、装配根未接），
		// 没有任何运行中的东西随这枚开关变化；清单上该行的生效位由装配探针
		// 判（cmd/wisp plugins_inventory_215a.go），探针没有这枚＝未生效。
		// 静默假装生效是票 215 点名禁止的那一种；面板侧也开不了（开＝放宽，
		// 被拒）——这句话是收据上的，不是装饰。
		res.Tier = panel.EffectiveNow
		res.Note = "开关已写入配置层并立即生效（locked 段收紧走同一条热路）。" +
			"插件引擎未在装配根接线：没有引擎消费这枚开关，清单上该行的生效位由装配探针判，探针上没有这枚＝未生效。" +
			"面板侧只能关不能开：开启是 [plugins] 锁定段的放宽，要走原生二次确认那条路。"
	}
	if w.Field == panel.FieldBallStyle {
		// 票 65：档位是这一写实测出来的，不是盖出来的（187 的 fieldTier187 能
		// 静态标是因为它的读者行为固定；ball_style 的读者＝本进程有没有球）。
		// applier 在（球活着）＝ApplyStyle 已把进程级家族与 orb 页面换过去，
		// 如实答 now；applier 不在（无窗启动）＝没有东西可即时换，下次启动的
		// 种子读取会落到它，答 restart（上面盖的默认档）。
		if applyBallStyle65(strings.TrimSpace(w.Value)) {
			res.Tier = panel.EffectiveNow
		}
	}
	// 外观波的两枚 [ball] 写：档位同样实测（65 的同款），不是盖出来的。
	if w.Field == panel.FieldBallOpacityIdle && waveOpacityParsed {
		if applyBallIdleOpacity(waveOpacity) {
			res.Tier = panel.EffectiveNow
		}
	}
	if w.Field == panel.FieldBallClickThrough {
		if applyBallClickThrough(waveClickOn) {
			res.Tier = panel.EffectiveNow
		}
	}
	// 外观波的两枚 [panel] 写：档位盖 restart（词表里最接近的静态档），但
	// Note 把真实的读者节奏写明——面板宿主每次建窗现读 [panel]，已建好的窗
	// 在下一次重新显示时也会把此刻解析出的数发一次（票 255-r1），所以「关窗
	// 再开」就跟上，不需要重启进程。重启档在这里是保守档，不是谎报的档。
	if w.Field == panel.FieldPanelWidth || w.Field == panel.FieldPanelFontSize {
		res.Note = "面板宿主每次建窗现读 [panel]：关窗再开面板（下一次显示请求）即按新值重算，无需重启进程。"
	}
	// 票 193 AC#3 接线：动画档的档位与注册表同源（票 39 立的规矩「档位判断
	// 是注册表的，不是手盖的」，config_tier_39.go）。[panel] 在 TierRegistry
	// 里是整节行（tiers.go: "panel" = hot；D36 的 [panel] 行），没有 per-key
	// 行，所以问注册表的是键路径的节名——整节 hot 即本键 hot。hot ⇒ now
	// （已立即生效列）——这是实测为真的档，不是保守档：写入即落到进程内存
	// （hot 档语义），快照泵的 PanelAnimation 读面每次发包现读同一个 Manager
	// （run.go 的 PumpSources 与常驻泵两处接线，本 commit 同落），所以面板
	// 收到的下一个快照就带着新档位词，页面随包消费，无需重启进程。这正是
	// 193 证据件 §③ 定 hot 的四条理由里最重的一条：重启档会把开关变成装饰
	// 品（页面的消费靠泵送，不在建窗时）。
	if w.Field == panel.FieldPanelAnimation {
		animationSection := strings.Split(keyPathPanelAnimation193, ".")[0]
		if t, ok := registryEffectiveTier39(animationSection); ok {
			res.Tier = t
		}
		res.Note = "动画档随快照走：快照每次发送都现读 panel.animation，面板下一次收到快照" +
			"（每次 show／状态变化都会推送）就按新档位运行，无需重启进程。"
	}
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

// configSectionRows456 maps internal/config's section-tree walk onto the
// panel's row shape. A copy, not a re-walk: the tree is computed in one place
// (config.SectionRows) and this leg renders it, so a second walker cannot
// disagree about what is in the file.
func configSectionRows456(cfg *config.Config) []panel.ConfigSectionRow {
	rows := config.SectionRows(cfg)
	if rows == nil {
		return nil
	}
	out := make([]panel.ConfigSectionRow, 0, len(rows))
	for _, r := range rows {
		out = append(out, panel.ConfigSectionRow{Key: r.Key, Value: r.Value, Tier: r.Tier})
	}
	return out
}

// configLoadErrors456 maps internal/config's structured load errors onto the
// panel's row shape - the same copy-not-rewalk rule as above: the pipeline
// and its rows are computed in one place (config.ParseErrors*), this leg
// renders them.
func configLoadErrors456(rows []config.LoadError) []panel.ConfigLoadError {
	if rows == nil {
		return nil
	}
	out := make([]panel.ConfigLoadError, 0, len(rows))
	for _, r := range rows {
		out = append(out, panel.ConfigLoadError{Line: r.Line, Message: r.Message})
	}
	return out
}

// splitChainValue turns the panel's comma-separated chain value into the element
// slice config's setter takes. The panel's own gate refuses an empty value for
// these fields, so this never has to decide what "" means.
func splitChainValue(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.TrimSpace(p))
	}
	return out
}

// splitCustomRowValue splits ticket 320 AC#2's row value ("<protocol>,<base_url>")
// at the FIRST comma, and trims both halves - the same courtesy splitChainValue
// shows an operator who types a space after the comma. A value with no comma
// yields an empty second half, which the config setter refuses by name; this
// helper never guesses that a missing half means "take the preset default",
// because a non-preset row has no preset to fall back to.
func splitCustomRowValue(v string) (protocol, baseURL string) {
	proto, url, _ := strings.Cut(v, ",")
	return strings.TrimSpace(proto), strings.TrimSpace(url)
}
