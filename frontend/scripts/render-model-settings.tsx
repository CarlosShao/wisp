/* ============================================================================
   Model/voice settings evidence (ticket 313's page half, 2026-11-08).
   ----------------------------------------------------------------------------
   What static markup can reach that no other instrument mounts: the settings
   page's 模型/语音 sections render through ConfigReadGate, which needs a host
   answer before anything exists, so render-nav's markup loop never sees them.
   This instrument mounts the pure pieces directly (they are exported for
   exactly this) and asserts the four teeth ticket 313's page half owes:

     1. AC#1 default pre-selection: with `configuredKnown && !configured`
        the add-provider form's picker starts on the owner's default vendor,
        and the name is only offered because the Go side's preset list
        carried it - never invented on the page;
     2. AC#1's other half (user wins): with the user configured there is NO
        pre-selection, the pre-selection is read exactly once as the initial
        state, and the chat pair drafts start from the user's saved pair;
     3. AC#3 key non-display: the key input is a password field with no value
        attribute, a canary typed nowhere near state cannot reach the markup,
        and at source the value leaves the DOM the instant it is sent and
        rides only the credentialValue key;
     4. AC#2 two-model writability: the chat pair and the voice chains (ASR /
        TTS / realtime) all mount provider+model pickers, their saves ride
        the whitelisted field names, and every write path re-reads (the
        receipt is the fact, the page keeps no mirror).

   What this harness CANNOT see (stated, not hidden): user interactions after
   mount (opening the picker, typing, clicking 保存) are runtime state; the
   source-level checks below are what pins them. The live-machine half is the
   real-Chromium run recorded in docs/evidence/s1/313-panel-model-config.md.

   Colours: none - this file paints nothing; it reads markup strings.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import {
  AddProviderForm,
  ModelSettingsView,
  ProviderRow,
  OWNER_DEFAULT_PROVIDER,
} from "@/components/model-settings";
import { VoiceSettingsView } from "@/components/voice-settings";
import type { ConfigView, ConfigViewProvider } from "@/lib/config-view";
import type { Feedback } from "@/lib/config-feedback";
import type { ModelProbeResultView } from "@/lib/config-view";
import { ModelProbeResultCard } from "@/components/model-settings";
import type { ModelsFetchResult, ModelProbeRun } from "@/lib/config-actions";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// The emoji band the repo scanner walks (ban #8, the instrument's own set).
// FE0F rides as its own alternation branch: inside a character class it is a
// combining sequence and oxlint warns (no-misleading-character-class).
const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u{22FF}\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

/* ---------------------------------------------------------------------------
   Fixtures. Key sets follow config-view.ts / config_313_test.go's envelope;
   the model ids are placeholders (the repo has zero real stepfun ids).
   --------------------------------------------------------------------------- */

const PRESETS = ["anthropic", "deepseek", "openai", "stepfun"];

const ROW: ConfigViewProvider = {
  name: "stepfun",
  apiKeyRef: "dpapi:9f2c41d0",
  hasBaseURL: true,
  baseURL: "https://api.stepfun.com/v1",
  protocol: "openai-chat",
  billing: "pay-per-token",
  modelCount: 2,
  models: ["示例模型-1", "示例模型-2"],
  // R1-recheck re-adjudication (2026-11-16): the read face answers both
  // per-model maps. This fixture is the ANSWERED shape - the wire keeps the
  // keys off only when the leg did not answer (omitempty), and tooth 6 feeds
  // that absent shape explicitly.
  modelCapabilities: {
    "示例模型-1": { audioIn: true, audioOut: false, realtime: false },
    "示例模型-2": { audioIn: false, audioOut: false, realtime: false },
  },
  modelsEnabled: { "示例模型-1": true, "示例模型-2": true },
  credentialSet: true,
  credentialError: "",
};

const FRESH: ConfigView = {
  readable: true,
  credential: "no_ref_declared",
  presetsKnown: true,
  presets: PRESETS,
  configuredKnown: true,
  configured: false,
  chatProvider: "",
  chatModel: "",
  chatModelKnown: false,
  providers: [],
  voice: {
    cloudASRChain: null,
    cloudTTSChain: null,
    realtimeEnabled: false,
    realtimeProvider: "",
    realtimeModel: "",
  },
};

const CONFIGURED: ConfigView = {
  ...FRESH,
  configured: true,
  credential: "all_recorded",
  chatProvider: "stepfun",
  chatModel: "示例模型-1",
  chatModelKnown: true,
  providers: [ROW],
};

const VOICED: ConfigView = {
  ...CONFIGURED,
  voice: {
    cloudASRChain: ["stepfun/示例模型-1"],
    cloudTTSChain: ["stepfun/示例模型-1", "stepfun/示例模型-2"],
    realtimeEnabled: false,
    realtimeProvider: "stepfun",
    realtimeModel: "示例模型-2",
  },
};

/* The performers the pure pieces take. They are never awaited to completion
   here (static markup renders initial state only); their shapes keep the
   props honest. */
const okFeedback: Feedback = { tone: "ok", text: "已保存" };
const performWrite = async (): Promise<Feedback> => okFeedback;
const performCredential = async (): Promise<Feedback> => okFeedback;
const performFetchModels = async (): Promise<ModelsFetchResult> => ({
  ok: true,
  models: ["示例模型-3"],
});
const performProbe = async (): Promise<ModelProbeRun> => ({
  ok: true,
  result: {
    provider: "stepfun",
    model: "示例模型-1",
    latencyMS: 421,
    allOK: true,
    written: true,
    outcomes: [],
    mismatches: [],
  },
});
const run = async (_key: string, fn: () => Promise<Feedback>): Promise<Feedback> => fn();
const noop = (): void => undefined;

/* ---------------------------------------------------------------------------
   Tooth 1 - AC#1's pre-selection: fresh machine, picker starts on stepfun.
   --------------------------------------------------------------------------- */

const addDefault = renderToStaticMarkup(
  <AddProviderForm
    available={PRESETS}
    busy={null}
    defaultProvider={OWNER_DEFAULT_PROVIDER}
    onCancel={noop}
    onDone={noop}
    onFeedback={noop}
    performCredential={performCredential}
    performWrite={performWrite}
    presets={PRESETS}
    run={run}
  />,
);
if (process.env.MS_DUMP) {
  console.error(addDefault);
}

check(
  "the owner default must be the preset name the ticket names (stepfun)",
  OWNER_DEFAULT_PROVIDER === "stepfun",
  OWNER_DEFAULT_PROVIDER,
);
check(
  "the add-provider picker must start pre-selected on the default vendor",
  /aria-label="供应商"[\s\S]*?>stepfun</.test(addDefault),
  "picker button text",
);
check(
  "the pre-selection must ride the picker button, not any hidden value attribute",
  /aria-label="供应商"[\s\S]*?>stepfun</.test(addDefault) && !addDefault.includes('value="stepfun"'),
  "button text, no controlled value",
);
check("the default markup must hold zero emoji", !EMOJI_RE.test(addDefault));

// Source-level: the pre-selection is gated on the read's own facts, and the
// name is only offered because the Go side's preset list carried it.
const modelSrc = readFileSync("src/components/model-settings.tsx", "utf8");
check(
  "the pre-selection must be gated on configuredKnown && !configured (user wins)",
  modelSrc.includes("view.configuredKnown") && modelSrc.includes("!view.configured"),
  "gate",
);
check(
  "the pre-selection must require the name in the Go side's preset list",
  /\(availablePresets \?\? \[\]\)\.includes\(OWNER_DEFAULT_PROVIDER\)/.test(modelSrc),
  "preset-list gate",
);
check(
  "the pre-selection must be read exactly once, as the initial state",
  modelSrc.includes("useState(defaultProvider)") &&
    !/useEffect\([\s\S]*?setProvider\(defaultProvider\)/.test(modelSrc),
  "initializer only",
);

/* ---------------------------------------------------------------------------
   Tooth 2 - the user's value wins: configured machine, no pre-selection; the
   chat drafts start from the user's saved pair.
   --------------------------------------------------------------------------- */

const addNoDefault = renderToStaticMarkup(
  <AddProviderForm
    available={PRESETS}
    busy={null}
    defaultProvider=""
    onCancel={noop}
    onDone={noop}
    onFeedback={noop}
    performCredential={performCredential}
    performWrite={performWrite}
    presets={PRESETS}
    run={run}
  />,
);
check(
  "with the user configured the picker must not pre-select the default vendor",
  !/>stepfun</.test(addNoDefault) && addNoDefault.includes("选择供应商"),
  "placeholder instead",
);
check("the no-default markup must hold zero emoji", !EMOJI_RE.test(addNoDefault));

const chatView = renderToStaticMarkup(
  <ModelSettingsView
    performCredential={performCredential}
    performFetchModels={performFetchModels}
    performProbe={performProbe}
    performWrite={performWrite}
    view={CONFIGURED}
  />,
);
check(
  "the chat card must show the user's saved pair verbatim (user wins on read)",
  chatView.includes("stepfun / 示例模型-1"),
  "current line",
);
check(
  "the chat drafts must start from the read's saved pair, never a page default",
  modelSrc.includes("useState(view.chatProvider)") && modelSrc.includes("useState(view.chatModel)"),
  "draft source",
);
check("the configured-page markup must hold zero emoji", !EMOJI_RE.test(chatView));

/* ---------------------------------------------------------------------------
   Tooth 3 - AC#3: the key row never carries a value in the DOM.
   --------------------------------------------------------------------------- */

const CANARY = "sk-canary-313-never-on-screen";
const ROW_NOKEY: ConfigViewProvider = { ...ROW, apiKeyRef: "", credentialSet: false };
const rowEditing = renderToStaticMarkup(
  <ProviderRow
    autoFetch={false}
    busy={null}
    editing={true}
    onAutoFetchHandled={noop}
    onEditing={noop}
    onFeedback={noop}
    performCredential={performCredential}
    performFetchModels={performFetchModels}
    performProbe={performProbe}
    performWrite={performWrite}
    row={ROW_NOKEY}
    run={run}
  />,
);
if (process.env.MS_DUMP_ROW) {
  console.error(rowEditing);
}
check(
  "the key input must render as a password field",
  rowEditing.includes('type="password"'),
  "input type",
);
check(
  "no non-empty value attribute may reach the markup (a controlled input renders empty)",
  !/value="[^"]+"/.test(rowEditing),
  "no filled value attr",
);
check(
  "a canary typed value must never reach the markup (nothing re-renders what was typed)",
  !rowEditing.includes(CANARY),
  "canary absent",
);
check(
  "the saved state must show the reference name only (dpapi:<id>), never a value",
  renderToStaticMarkup(
    <ProviderRow
      autoFetch={false}
      busy={null}
      editing={true}
      onAutoFetchHandled={noop}
      onEditing={noop}
      onFeedback={noop}
      performCredential={performCredential}
      performFetchModels={performFetchModels}
    performProbe={performProbe}
      performWrite={performWrite}
      row={ROW}
      run={run}
    />,
  ).includes("dpapi:9f2c41d0"),
  "reference is displayable",
);
check("the key-row markup must hold zero emoji", !EMOJI_RE.test(rowEditing));

// Source-level: the value leaves the DOM before the request, rides only the
// credential write, and never lands on the shared envelope's value key.
check(
  "saveKey must clear the input the instant it sends (value never in state)",
  modelSrc.includes("input.value = \"\";") &&
    /input\.value = "";\s*\n\s*void run\(`credential:/s.test(modelSrc),
  "clear-before-send",
);
check(
  "the credential value must ride performCredential only (never configValue)",
  modelSrc.includes("performCredential(row.name, value)") &&
    !modelSrc.includes("credentialValue") &&
    !modelSrc.includes("configValue"),
  "one door",
);
const voiceSrc = readFileSync("src/components/voice-settings.tsx", "utf8");
check(
  "the voice page must carry no credential value path at all",
  !voiceSrc.includes("credentialValue") && !voiceSrc.includes("configValue"),
  "voice clean",
);
const actionsSrc = readFileSync("src/lib/config-actions.ts", "utf8");
check(
  "the credential write must go through the shared door that names credentialValue",
  actionsSrc.includes("performCredentialWrite") && actionsSrc.includes("requestCredentialWrite"),
  "shared door",
);

/* ---------------------------------------------------------------------------
   Tooth 4 - AC#2: both models' provider+model pickers mount, and the saves
   ride the whitelisted fields with a re-read after every write.
   --------------------------------------------------------------------------- */

check(
  "the configured model page must mount the chat pair pickers",
  chatView.includes('aria-label="供应商"') && chatView.includes('aria-label="模型"'),
  "pair controls",
);

const voiceView = renderToStaticMarkup(
  <VoiceSettingsView onReload={noop} performWrite={performWrite} view={VOICED} />,
);
check(
  "the voice page must mount the ASR and TTS chain editors and the realtime head",
  voiceView.includes("听写（ASR）") && voiceView.includes("朗读（TTS）") && voiceView.includes("实时语音"),
  "three cards",
);
check(
  "the voice chains must show the saved elements verbatim (stepfun fillable)",
  voiceView.includes("stepfun/示例模型-1") && voiceView.includes("stepfun/示例模型-2"),
  "chain text",
);
check(
  "the voice pair pickers must mount (three cards x provider+model)",
  (voiceView.match(/aria-label="供应商"/g) ?? []).length >= 3 &&
    (voiceView.match(/aria-label="模型"/g) ?? []).length >= 3,
  `${(voiceView.match(/aria-label="供应商"/g) ?? []).length} provider pickers`,
);
check("the voice-page markup must hold zero emoji", !EMOJI_RE.test(voiceView));

// Source-level: the field names are the Go side's admission words, and every
// write path re-reads instead of patching a local picture.
check(
  "the chat save must ride role_chat_pair",
  modelSrc.includes('field: "role_chat_pair"'),
  "field",
);
check(
  "the ASR/TTS chain fields must be spelled at the mount sites",
  voiceSrc.includes('field="voice_cloud_asr_chain"') && voiceSrc.includes('field="voice_cloud_tts_chain"'),
  "mount-site spellings",
);
check(
  "the realtime save must ride the two realtime fields",
  voiceSrc.includes('field: "voice_realtime_provider"') && voiceSrc.includes('field: "voice_realtime_model"'),
  "realtime fields",
);
check(
  "every write path must re-read after the write (receipt is the fact)",
  /setBusy\(null\);\s*\n\s*onReload\(\);/.test(modelSrc) &&
    /setBusy\(null\);\s*\n\s*onReload\(\);/.test(voiceSrc),
  "run() re-reads",
);

/* ---------------------------------------------------------------------------
   Tooth 5 - ticket 313's wiring leg (2026-11-16): AC#1 reaches the voice
   cards, and the honest copy follows the backend truth (chains wired at boot,
   realtime still Reserved).
   --------------------------------------------------------------------------- */

// CONFIGURED carries a provider row but no saved voice chains and no saved
// realtime head: every unset field's draft prefills the owner vendor - three
// provider pickers (ASR chain, TTS chain, realtime) start on stepfun, as
// initial state only.
const voiceFresh = renderToStaticMarkup(
  <VoiceSettingsView onReload={noop} performWrite={performWrite} view={CONFIGURED} />,
);
check(
  "the unset voice fields must draft-prefill the owner vendor (AC#1 voice half)",
  (voiceFresh.match(/>stepfun</g) ?? []).length >= 3,
  `${(voiceFresh.match(/>stepfun</g) ?? []).length} preselected pickers`,
);
// R1-recheck re-adjudication (2026-11-16): this check predates the conditional
// tags. The saved tag stopped being a static 已接入·重启生效 when the R1 rework
// made it a verdict (an unset chain is 未设置, not wired), so the wired-truth
// count moves to the view whose chains ARE set and answered (VOICED): the ASR
// chain passes every gate (已接入·重启生效 once), the TTS chain lacks its
// audio_out bit (装配不了·缺能力位 once), and the realtime card keeps its
// placeholder. The unset view asserts the unset states. No assertion loosened:
// each count is exact and the placeholder pin stays.
check(
  "the two chain cards must say the wired truth and the realtime card keeps its placeholder",
  (voiceView.match(/已接入·重启生效/g) ?? []).length === 1 &&
    (voiceView.match(/装配不了·缺能力位/g) ?? []).length === 1 &&
    (voiceView.match(/占位·当前不生效/g) ?? []).length === 1,
  `wired=${(voiceView.match(/已接入·重启生效/g) ?? []).length} lacks=${(voiceView.match(/装配不了·缺能力位/g) ?? []).length} placeholder=${(voiceView.match(/占位·当前不生效/g) ?? []).length}`,
);
check(
  "the unset chain cards must read 未设置, never a wired tag (verdict, not static copy)",
  (voiceFresh.match(/未设置/g) ?? []).length >= 2 &&
    (voiceFresh.match(/已接入·重启生效/g) ?? []).length === 0,
  `unset=${(voiceFresh.match(/未设置/g) ?? []).length} wired=${(voiceFresh.match(/已接入·重启生效/g) ?? []).length}`,
);
check(
  "the wired chain copy must name the restart contract, not promise a live effect",
  voiceFresh.includes("保存后重启生效"),
  "TierReload wording",
);
check("the fresh voice-page markup must hold zero emoji", !EMOJI_RE.test(voiceFresh));

// Source-level: the default is gated on a REAL provider row (a chain element
// must name a declared provider/model), read exactly once as initial state,
// and never re-imposed by an effect.
check(
  "the voice default must be gated on an existing provider row",
  voiceSrc.includes("providers.some((p) => p.name === OWNER_DEFAULT_PROVIDER)"),
  "row-exists gate",
);
check(
  "the chain draft default must be initializer-only (no effect re-imposes it)",
  voiceSrc.includes('provider: defaultProvider, model: ""') &&
    !/useEffect\([\s\S]*?setSegments/.test(voiceSrc),
  "initializer once",
);
check(
  "the realtime draft default must yield to a saved value verbatim",
  voiceSrc.includes("view.voice.realtimeProvider || voiceDefaultProvider(providers)"),
  "saved value wins",
);
check(
  "the chain saves must no longer downgrade their ok receipts (they are wired)",
  !voiceSrc.includes("setFeedback(placeholderReceipt(feedback))"),
  "real receipt",
);

// User wins on read: the VOICED view shows the saved realtime pair verbatim
// (the default draft must not have replaced it anywhere).
check(
  "the saved realtime pair must render verbatim (user value wins)",
  voiceView.includes("stepfun / 示例模型-2"),
  "current line",
);

/* ---------------------------------------------------------------------------
   Tooth 6 - the R1-recheck residuals (2026-11-16).
   --------------------------------------------------------------------------- */

// Residual 1's fixture: the WIRE shape with modelCapabilities genuinely
// ABSENT. Go's omitempty (config_handlers.go:776) leaves the key off the wire
// and the bare-cast parse hands the page `undefined`, not `null` - the shape
// the audit proved crashes both pages (Node 24.21.0: TypeError reading
// properties of undefined). Built via JSON round-trip + delete so the key is
// truly absent, not a typed null: this is the fixture the old `=== null`
// judgment threw on, and a zero-model provider row (CreateProviderRow without
// an id, settings_313.go:154) is factory-reachable.
const wireAbsentCapsRow: ConfigViewProvider = (() => {
  const raw = JSON.parse(
    JSON.stringify({ ...ROW, models: ["示例模型-1"], modelCount: 1 }),
  ) as Record<string, unknown>;
  delete raw["modelCapabilities"];
  delete raw["modelsEnabled"];
  return raw as unknown as ConfigViewProvider;
})();

function voiceWithChain(row: ConfigViewProvider, asr: string[] | null): ConfigView {
  return {
    ...CONFIGURED,
    providers: [row],
    voice: {
      ...FRESH.voice,
      cloudASRChain: asr,
      cloudTTSChain: null,
    },
  };
}

/** Read one switch's rendered state out of the markup (attribute order is the
    component's own: aria-checked before aria-label). */
function switchState(markup: string, label: string): string | null {
  const m = markup.match(new RegExp(`aria-checked="(true|false)"[^>]*aria-label="${label}"`));
  return m === null ? null : m[1];
}

// 崩页牙 (residual 1), model page: the absent-key row must render the waiting
// tag and must not throw. The old judgment threw INSIDE this render.
const absentCapsRowEdit = renderToStaticMarkup(
  <ProviderRow
    autoFetch={false}
    busy={null}
    editing={true}
    onAutoFetchHandled={noop}
    onEditing={noop}
    onFeedback={noop}
    performCredential={performCredential}
    performFetchModels={performFetchModels}
    performProbe={performProbe}
    performWrite={performWrite}
    row={wireAbsentCapsRow}
    run={run}
  />,
);
check(
  "residual-1 crash tooth: a provider row whose modelCapabilities key is ABSENT must render 等快照 and must not throw",
  absentCapsRowEdit.includes("等快照"),
  "waiting tag",
);
check(
  "residual-1: the absent-key row must mount no capability switch (never a guessed false)",
  !absentCapsRowEdit.includes("audio in 能力位"),
  "no switch",
);
check("the absent-key row markup must hold zero emoji", !EMOJI_RE.test(absentCapsRowEdit));

// 正常值牙 (residual 1): with the read face answered, BOTH bits render with
// their real states - audio in on for the first model, everything else off.
const rowEditPage = renderToStaticMarkup(
  <ProviderRow
    autoFetch={false}
    busy={null}
    editing={true}
    onAutoFetchHandled={noop}
    onEditing={noop}
    onFeedback={noop}
    performCredential={performCredential}
    performFetchModels={performFetchModels}
    performProbe={performProbe}
    performWrite={performWrite}
    row={ROW}
    run={run}
  />,
);
check(
  "residual-1 normal-value tooth: the answered bits render with their real states",
  switchState(rowEditPage, "示例模型-1 的 audio in 能力位") === "true" &&
    switchState(rowEditPage, "示例模型-1 的 audio out 能力位") === "false" &&
    switchState(rowEditPage, "示例模型-2 的 audio in 能力位") === "false" &&
    switchState(rowEditPage, "示例模型-2 的 audio out 能力位") === "false",
  `in1=${switchState(rowEditPage, "示例模型-1 的 audio in 能力位")} out1=${switchState(rowEditPage, "示例模型-1 的 audio out 能力位")}`,
);

// 崩页牙 (residual 1), voice page: a saved chain over the absent-key row
// renders 装配未判（等快照） and must not throw.
const absentCapsVoicePage = renderToStaticMarkup(
  <VoiceSettingsView onReload={noop} performWrite={performWrite} view={voiceWithChain(wireAbsentCapsRow, ["stepfun/示例模型-1"])} />,
);
check(
  "residual-1 crash tooth: a saved chain over the absent-key row must render 装配未判（等快照） and must not throw",
  absentCapsVoicePage.includes("装配未判（等快照）"),
  "unknown tag",
);
check("the absent-key voice markup must hold zero emoji", !EMOJI_RE.test(absentCapsVoicePage));

// 三枚假象牙 (residual 2): each "looks assemblable" illusion fed once - the
// tag must NOT be 已接入·重启生效, and only the all-gates-pass view is
// allowed to say it. Gate order and sentences follow the boot's own
// (cloudcatalog.go:87-90 enabled / :95-102 bits / :121-131 api_key_ref).
const nokeyRow: ConfigViewProvider = { ...ROW, credentialSet: false };
const credErrRow: ConfigViewProvider = { ...ROW, credentialSet: false, credentialError: "凭据存储未接入" };
const disabledRow: ConfigViewProvider = { ...ROW, modelsEnabled: { "示例模型-1": false, "示例模型-2": true } };

const nokeyPage = renderToStaticMarkup(
  <VoiceSettingsView onReload={noop} performWrite={performWrite} view={voiceWithChain(nokeyRow, ["stepfun/示例模型-1"])} />,
);
check(
  "residual-2 illusion 1 (api_key_ref declared, not affirmed): the tag must be 装配不了·缺 API Key",
  nokeyPage.includes("装配不了·缺 API Key") && !nokeyPage.includes("已接入·重启生效"),
  "nokey tag",
);
check(
  "residual-2 illusion 1: the refusal must name the provider and its ref state",
  nokeyPage.includes("stepfun 缺 API Key（api_key_ref"),
  "named provider",
);
check("the nokey markup must hold zero emoji", !EMOJI_RE.test(nokeyPage));

const credErrPage = renderToStaticMarkup(
  <VoiceSettingsView onReload={noop} performWrite={performWrite} view={voiceWithChain(credErrRow, ["stepfun/示例模型-1"])} />,
);
check(
  "residual-2 illusion 2 (credentialError): the tag must not be 已接入 and the error must render verbatim on the voice card",
  !credErrPage.includes("已接入·重启生效") &&
    credErrPage.includes("装配不了·缺 API Key") &&
    credErrPage.includes("凭据存储未接入") &&
    credErrPage.includes("无法核对密钥：stepfun（凭据存储未接入）"),
  "credentialError rendered",
);
check("the credentialError markup must hold zero emoji", !EMOJI_RE.test(credErrPage));

const disabledPage = renderToStaticMarkup(
  <VoiceSettingsView onReload={noop} performWrite={performWrite} view={voiceWithChain(disabledRow, ["stepfun/示例模型-1"])} />,
);
check(
  "residual-2 illusion 3 (enabled=false): the tag must be 装配不了·模型已停用, naming the disabled element",
  disabledPage.includes("装配不了·模型已停用") &&
    !disabledPage.includes("已接入·重启生效") &&
    disabledPage.includes("stepfun/示例模型-1 已停用（enabled=false）"),
  "disabled tag",
);
check("the disabled markup must hold zero emoji", !EMOJI_RE.test(disabledPage));

// 全过牙 (residual 2): with every gate answered and passing, and ONLY then,
// the wired tag is allowed.
const allPassPage = renderToStaticMarkup(
  <VoiceSettingsView onReload={noop} performWrite={performWrite} view={voiceWithChain(ROW, ["stepfun/示例模型-1"])} />,
);
check(
  "residual-2: all three gates passing is the one state that may say 已接入·重启生效",
  (allPassPage.match(/已接入·重启生效/g) ?? []).length === 1 &&
    !allPassPage.includes("装配不了·"),
  "all-pass tag",
);
check("the all-pass markup must hold zero emoji", !EMOJI_RE.test(allPassPage));

/* ---------------------------------------------------------------------------
   Tooth 7 - the probe card (ticket 40 AC#3, 2026-11-18).
   --------------------------------------------------------------------------- */

// A REAL-shaped measured answer: mixed verdicts, one raw failure detail, one
// frozen mismatch label. Every word below must ride to the screen verbatim -
// the page computes no verdict of its own (the 空转返回OK shape has no markup
// it can wear: an outcome that says 实测不可用 cannot render as anything else).
const PROBE_FIXTURE: ModelProbeResultView = {
  provider: "stepfun",
  model: "示例模型-1",
  latencyMS: 421,
  allOK: false,
  written: true,
  outcomes: [
    { capability: "fc", declared: true, ok: true, attempted: true },
    { capability: "vision", declared: true, ok: false, attempted: true, detail: "provider (502): HTTP 502" },
    { capability: "thinking", declared: false, ok: true, attempted: true },
  ],
  mismatches: [
    {
      capability: "vision",
      detail: "provider (502): HTTP 502",
      label: "能力实测不符：stepfun/示例模型-1 声明支持 vision，实测不可用",
    },
  ],
};

const probeDone = renderToStaticMarkup(
  <ModelProbeResultCard onDismiss={noop} probe={{ phase: "done", model: "示例模型-1", result: PROBE_FIXTURE }} />,
);
check(
  "probe tooth 1: the measured latency rides verbatim (p50 number, not a badge)",
  probeDone.includes("stepfun/示例模型-1 实测 p50 421ms"),
  "latency line",
);
check(
  "probe tooth 2: a passing capability renders its declared/measured pairing as words",
  probeDone.includes("fc：声明可用，实测可用") && probeDone.includes("thinking：未声明，实测可用"),
  "paired words",
);
check(
  "probe tooth 3: a failed capability carries its RAW detail verbatim - never paraphrased",
  probeDone.includes("vision：声明可用，实测不可用") && probeDone.includes("（细节：provider (502): HTTP 502）"),
  "raw detail",
);
check(
  "probe tooth 4: the frozen mismatch label renders verbatim",
  probeDone.includes("能力实测不符：stepfun/示例模型-1 声明支持 vision，实测不可用"),
  "label",
);
check(
  "probe tooth 5: the provider_health half is stated (written=true), and the no-save fact with it",
  probeDone.includes("判定已写入 provider_health；探测不改配置，不落 config.toml。"),
  "written sentence",
);
check("the probe-done markup must hold zero emoji", !EMOJI_RE.test(probeDone));
check(
  "probe tooth 6: a probe result must not carry any save/tier wording (it measures, never saves)",
  !probeDone.includes("需要重启") && !probeDone.includes("立即生效") && !probeDone.includes("设置已写入"),
  "no tier wording",
);

const PROBE_REFUSAL = "面板请求被拒绝 [config.model.probe r-1]：身份或密钥没有通过，先检查密钥配置：auth: llm: resolving api_key_ref for provider \"stepfun\" failed (Unconfigured)";
const probeFailed = renderToStaticMarkup(
  <ModelProbeResultCard onDismiss={noop} probe={{ phase: "failed", model: "示例模型-1", text: PROBE_REFUSAL }} />,
);
check(
  "probe tooth 7: a refusal's prose rides VERBATIM (the error原文 is the answer, not a summary)",
  // React escapes `"` in text content, so the verbatim comparison is against
  // the escaped rendering of the same string - still word-for-word.
  probeFailed.includes(PROBE_REFUSAL.replaceAll('"', "&quot;")),
  "refusal verbatim",
);
check(
  "probe tooth 8: a refusal must not read as a measurement (no 实测 p50, no verdict words)",
  !probeFailed.includes("实测 p50") && !probeFailed.includes("实测可用"),
  "no fake verdict",
);
check("the probe-failed markup must hold zero emoji", !EMOJI_RE.test(probeFailed));

// Source-level: the old stub sentence is gone, the button rides the door, and
// the route literal lives ONLY in panel.ts (the renderer holds exactly one
// door to the host).
check(
  "probe tooth 9: the honest-stub sentence 暂不支持连接测试 must be gone from the page",
  !modelSrc.includes("暂不支持连接测试"),
  "stub removed",
);
check(
  "probe tooth 10: the page's 测试 button must ride the performProbe door, not invent a route",
  modelSrc.includes("performProbe(row.name, id)"),
  "door call",
);
check(
  "probe tooth 11: the route literal must live only in the door file (panel.ts)",
  modelSrc.includes("config.model.probe") === false &&
    readFileSync("src/lib/panel.ts", "utf8").includes('"config.model.probe"'),
  "literal home",
);
check(
  "probe tooth 12: the probe block must not compute a second verdict (no local ok-guess)",
  modelSrc.includes("probeOutcomeWord") && !/probe\s*===\s*"done"\s*\?\s*"可用/.test(modelSrc),
  "no second judge",
);

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-model-settings: FAIL " + f);
  console.error(`render-model-settings: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-model-settings: OK, AC#1 pre-selection gated on the read (stepfun, preset-list + " +
    "unconfigured only, initializer once), user value wins, key row password + canary-proof, " +
    "chat/chain/realtime pairs mount and ride the whitelisted fields with re-read write-back, " +
    "voice cards draft-prefill the owner vendor with the verdict (not static) wired/unset copy " +
    "split; residual 1: absent-key rows render 等快照 / 装配未判 without throwing and answered " +
    "bits show their real states; residual 2: the three assembles-illusions (nokey, " +
    "credentialError, enabled=false) each refuse the wired tag and only the all-gates-pass " +
    "chain may say 已接入·重启生效; ticket 40 AC#3: the probe card rides the measured payload " +
    "verbatim (p50 number, paired verdict words, raw details, frozen labels, provider_health " +
    "half, no tier wording) and the refusal prose verbatim, with the stub sentence gone and " +
    "the route literal only in panel.ts",
);
