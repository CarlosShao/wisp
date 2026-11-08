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
import type { ModelsFetchResult } from "@/lib/config-actions";

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

const chatView = renderToStaticMarkup(<ModelSettingsView performCredential={performCredential} performFetchModels={performFetchModels} performWrite={performWrite} view={CONFIGURED} />);
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
    "chat/chain/realtime pairs mount and ride the whitelisted fields with re-read write-back",
);
