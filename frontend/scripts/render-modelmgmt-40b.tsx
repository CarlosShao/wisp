/* ============================================================================
   Model-management page half B (ticket 40 AC#3 后三半, 2026-11-18).
   ----------------------------------------------------------------------------
   What this instrument pins that no Go test can: the PAGE half of the last
   three AC#3 halves.

     半一 capability 过滤 (voice roles hide non-voice models):
       the candidates a voice role's model dropdown is fed come from the read
       face's real capability bits (row.modelCapabilities, straight from
       spec.Capabilities.*) through voiceModelCandidates - NO hardcoded model
       list anywhere. Teeth: a non-voice model in a voice role's candidates is
       red (asserted absent), a voice model missing from its role's candidates
       is red (asserted present), and the three-state rule holds - an
       unanswered read face answers null (no filtering, never "hide
       everything"), an unanswered per-model entry stays in the candidates
       (没答不等于 false).

     半二 chain 拖拽序持久:
       the chain editor's rows are draggable (plus 上移/下移, one pure reorder
       rule both entries share), and the save writes the segments in their
       CURRENT array order - so a drag that is saved lands on config.toml in
       the new order (the disk half is the Go tooth, fresh LoadFile), and a
       drag that is never saved writes nothing.

     半三 quota 范围:
       the quota editor shows the CURRENT values off the read face verbatim
       (0 = 未设/不封顶, never folded into a guessed number), waits (等快照,
       no inputs) when the leg did not answer modelQuotas, sends edits through
       the whitelisted fields with NO page-side range judge (the Go setter is
       the single refusal source), and surfaces the host's refusal sentence
       verbatim through the same feedback line every card uses.

   What static markup cannot see (stated, not hidden): opening a picker,
   dragging a row and clicking 保存 are runtime interactions; the source-level
   checks below pin the wiring those interactions ride.

   Colours: none - this file paints nothing; it reads markup strings.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import { ModelQuotaEditor } from "@/components/model-settings";
import {
  moveChainSegment,
  VoiceSettingsView,
  voiceModelCandidates,
  type VoiceBit,
} from "@/components/voice-settings";
import { FeedbackLine } from "@/components/settings-controls";
import type { ConfigView, ConfigViewProvider } from "@/lib/config-view";
import type { Feedback } from "@/lib/config-feedback";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// The emoji band the repo scanner walks (ban #8, the instrument's own set).
const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u{22FF}\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

/* ---------------------------------------------------------------------------
   Fixtures. Model ids are made-up placeholders (the 313 rule); every bit is a
   real read-face shape.
   --------------------------------------------------------------------------- */

const MODELS = ["voice-in-only", "voice-out-only", "voice-both", "text-only", "realtime-only"];

function rowWith(caps: Record<string, { audioIn: boolean; audioOut: boolean; realtime: boolean }>): ConfigViewProvider {
  return {
    name: "acme",
    apiKeyRef: "",
    hasBaseURL: true,
    baseURL: "https://api.acme.invalid/v1",
    protocol: "openai-chat",
    billing: "pay-per-token",
    modelCount: MODELS.length,
    models: MODELS,
    modelCapabilities: caps,
    modelsEnabled: Object.fromEntries(MODELS.map((id) => [id, true])),
    credentialSet: false,
    credentialError: "",
  };
}

const ANSWERED_ROW = rowWith({
  "voice-in-only": { audioIn: true, audioOut: false, realtime: false },
  "voice-out-only": { audioIn: false, audioOut: true, realtime: false },
  "voice-both": { audioIn: true, audioOut: true, realtime: false },
  "text-only": { audioIn: false, audioOut: false, realtime: false },
  "realtime-only": { audioIn: false, audioOut: false, realtime: true },
});

/* ---------------------------------------------------------------------------
   半一: capability-filtered voice role candidates.
   --------------------------------------------------------------------------- */

function expectCandidates(
  label: string,
  got: string[] | null,
  wantInclude: string[],
  wantExclude: string[],
): void {
  if (got === null) {
    check(label, false, "candidates = null");
    return;
  }
  for (const id of wantInclude) {
    check(`${label}: ${id} 必须在候选里（语音模型漏掉⇒红）`, got.includes(id), got.join(","));
  }
  for (const id of wantExclude) {
    // 牙一: a non-voice model among the candidates is exactly the shape this
    // tooth exists to redden.
    check(`${label}: ${id} 不得出现在候选里（非语音模型塞进语音角色⇒红）`, !got.includes(id), got.join(","));
  }
}

expectCandidates("ASR(dictation) 候选按 audioIn 位过滤", voiceModelCandidates(ANSWERED_ROW, "audioIn"),
  ["voice-in-only", "voice-both"], ["text-only", "voice-out-only", "realtime-only"]);
expectCandidates("TTS(朗读) 候选按 audioOut 位过滤", voiceModelCandidates(ANSWERED_ROW, "audioOut"),
  ["voice-out-only", "voice-both"], ["text-only", "voice-in-only", "realtime-only"]);
expectCandidates("realtime 候选按 realtime 位过滤", voiceModelCandidates(ANSWERED_ROW, "realtime"),
  ["realtime-only"], ["text-only", "voice-in-only", "voice-out-only", "voice-both"]);

// The three-state rule: an unanswered read face answers null - the caller
// keeps the full catalog instead of folding 没答 into "nobody is a candidate".
check(
  "半一: 腿没答（modelCapabilities 缺键）⇒ null，不过滤（全隐藏＝把没答折成全员 false）",
  voiceModelCandidates({ ...ANSWERED_ROW, modelCapabilities: null }, "audioIn") === null,
  "unanswered read face",
);
// An answered record that skips one model's entry: that model is UNJUDGED, so
// it stays in the candidates (never a guessed false) while the answered bits
// still rule the rest.
const PARTIAL_ROW = rowWith({
  "voice-in-only": { audioIn: true, audioOut: false, realtime: false },
});
const partialGot = voiceModelCandidates(PARTIAL_ROW, "audioIn");
check(
  "半一: 记录答了但缺某模型条目 ⇒ 未判的那枚全部留在候选（没答不等于 false），已答的照真过滤",
  partialGot !== null && partialGot.includes("voice-in-only") && partialGot.includes("text-only") &&
    partialGot.includes("voice-both") && partialGot.includes("voice-out-only") &&
    partialGot.includes("realtime-only"),
  partialGot?.join(",") ?? "null",
);

// Source pins: the bit source is the read face (never a hardcoded list), the
// two chain cards and the realtime head all hand the filter to the pair
// control, and the chat card stays unfiltered (the dispatch names voice
// roles only).
const voiceSrc = readFileSync("src/components/voice-settings.tsx", "utf8");
const controlsSrc = readFileSync("src/components/settings-controls.tsx", "utf8");
check(
  "半一: voiceModelCandidates 的位来源是 row.modelCapabilities（前端不硬编码名单）",
  voiceSrc.includes("row.modelCapabilities ?? null") && voiceSrc.includes("bits[bit]"),
  "bit source",
);
check(
  "半一: ASR/TTS 链编辑器把位过滤喂给 PairFields",
  voiceSrc.includes("modelCandidates={(row) => voiceModelCandidates(row, DIRECTION_BIT[direction])}"),
  "chain wiring",
);
check(
  "半一: realtime 头把 realtime 位喂给 PairFields",
  voiceSrc.includes('voiceModelCandidates(row, "realtime")'),
  "realtime wiring",
);
check(
  "半一: 过滤只跑在答过的目录上；没答的读面回退全清单（settings-controls 的 null 回退）",
  controlsSrc.includes("modelCandidates(row) ?? all"),
  "null fallback",
);

/* ---------------------------------------------------------------------------
   半二: chain drag-order.
   --------------------------------------------------------------------------- */

const SEG = (p: string, m: string) => ({ provider: p, model: m });
const ABC = [SEG("a", "m1"), SEG("b", "m2"), SEG("c", "m3")];
check(
  "半二: moveChainSegment 0→2 ⇒ b,c,a（拖拽重排＝数组顺序重排）",
  JSON.stringify(moveChainSegment(ABC, 0, 2)) === JSON.stringify([SEG("b", "m2"), SEG("c", "m3"), SEG("a", "m1")]),
  "0->2",
);
check(
  "半二: moveChainSegment 2→0 ⇒ c,a,b（反向拖拽同一条规则）",
  JSON.stringify(moveChainSegment(ABC, 2, 0)) === JSON.stringify([SEG("c", "m3"), SEG("a", "m1"), SEG("b", "m2")]),
  "2->0",
);
check(
  "半二: 越界/同位的移动不改原序（一条重排规则，没有第二套）",
  moveChainSegment(ABC, 1, 1) === ABC && moveChainSegment(ABC, 3, 0) === ABC && moveChainSegment(ABC, -1, 0) === ABC,
  "no-op cases",
);

// Source pins: rows drag, the two move buttons share the pure rule, the drop
// lands on it, and the save writes the segments in their CURRENT order (that
// is what makes a saved drag persist - and an unsaved drag write nothing).
check(
  "半二: 链段行可拖拽（draggable）",
  voiceSrc.includes("draggable") && voiceSrc.includes("onDragStart") && voiceSrc.includes("onDrop"),
  "drag attrs",
);
check(
  "半二: 上移/下移按钮与 drop 走同一条 moveChainSegment",
  voiceSrc.includes("moveTo(index, index - 1)") && voiceSrc.includes("moveTo(index, index + 1)") &&
    voiceSrc.includes("moveChainSegment(all, from, to)"),
  "one reorder rule",
);
check(
  "半二: save 按 segments 现序拼值（拖完保存才落盘；不保存不落盘）",
  /savedChain = segments\.map\(\(seg\) =>/.test(voiceSrc) && voiceSrc.includes("value: savedChain.join(\",\")"),
  "save order",
);

// Markup pin: the mounted chain editor carries the drag affordance and the
// two move buttons (aria labels), for a saved chain.
const CHAIN_VIEW: ConfigView = {
  readable: true,
  credential: "no_ref_declared",
  presetsKnown: true,
  presets: ["acme"],
  configuredKnown: true,
  configured: true,
  chatProvider: "",
  chatModel: "",
  chatModelKnown: false,
  providers: [ANSWERED_ROW],
  voice: {
    cloudASRChain: ["acme/voice-in-only"],
    cloudTTSChain: ["acme/voice-out-only"],
    realtimeEnabled: false,
    realtimeProvider: "",
    realtimeModel: "",
  },
};
const noop = (): void => undefined;
const chainMarkup = renderToStaticMarkup(
  <VoiceSettingsView onReload={noop} performWrite={async () => ({ tone: "ok", text: "已接入：保存后重启生效" })} view={CHAIN_VIEW} />,
);
check("半二: 挂载的链编辑器带 draggable 行", chainMarkup.includes("draggable"), "draggable attr");
check(
  "半二: 挂载的链编辑器带上移/下移按钮（aria 名）",
  chainMarkup.includes("上移") && chainMarkup.includes("下移"),
  "move buttons",
);
check("the chain-editor markup must hold zero emoji", !EMOJI_RE.test(chainMarkup));

/* ---------------------------------------------------------------------------
   半三: quota edits, ranges enforced on the Go side, values read verbatim.
   --------------------------------------------------------------------------- */

const QUOTA_ROW: ConfigViewProvider = {
  ...ANSWERED_ROW,
  modelQuotas: {
    "voice-in-only": { dailyMicro: 1000000, monthlyMicro: 20000000 },
    "voice-out-only": { dailyMicro: 0, monthlyMicro: 0 },
    "voice-both": { dailyMicro: 0, monthlyMicro: 0 },
    "text-only": { dailyMicro: 0, monthlyMicro: 0 },
    "realtime-only": { dailyMicro: 0, monthlyMicro: 0 },
  },
  planCreditTotalMicro: 5000000,
};
const quotaMarkup = renderToStaticMarkup(
  <ModelQuotaEditor
    busy={null}
    onFeedback={noop}
    performWrite={async () => ({ tone: "ok", text: "已保存，重启后生效" })}
    row={QUOTA_ROW}
    run={async (_key, fn) => fn()}
  />,
);
check(
  "半三: 当前额度值逐字上屏（1000000/20000000）",
  quotaMarkup.includes("1000000") && quotaMarkup.includes("20000000"),
  "current values",
);
check(
  "半三: 当前值 0 逐字显示为 未设（0＝真答案，不折成猜的数）",
  quotaMarkup.includes("未设（0 = 不封顶）"),
  "zero reading",
);
check(
  "半三: 套餐 credit 总额读口上屏（5000000）",
  quotaMarkup.includes("5000000"),
  "plan credit",
);
// The credit total's own zero reading: 0 = 未录入 is an answered value that
// gates nothing (quota_44.go's own reading), rendered verbatim, never 0-as-fact.
const CREDIT_ZERO_ROW: ConfigViewProvider = { ...ANSWERED_ROW, modelQuotas: undefined, planCreditTotalMicro: 0 };
const creditZeroMarkup = renderToStaticMarkup(
  <ModelQuotaEditor
    busy={null}
    onFeedback={noop}
    performWrite={async () => ({ tone: "ok", text: "ok" })}
    row={CREDIT_ZERO_ROW}
    run={async (_key, fn) => fn()}
  />,
);
check(
  "半三: credit 总额 0 逐字显示为 未录入（0）",
  creditZeroMarkup.includes("未录入（0）"),
  "credit zero reading",
);
check("the quota-editor markup must hold zero emoji", !EMOJI_RE.test(quotaMarkup));

// The waiting state: a leg that did not answer modelQuotas renders 等快照 and
// draws NO quota input (an editor that cannot see the current value would be
// overwriting blind).
const WAIT_ROW: ConfigViewProvider = { ...ANSWERED_ROW, modelQuotas: undefined, planCreditTotalMicro: 0 };
const waitMarkup = renderToStaticMarkup(
  <ModelQuotaEditor
    busy={null}
    onFeedback={noop}
    performWrite={async () => ({ tone: "ok", text: "ok" })}
    row={WAIT_ROW}
    run={async (_key, fn) => fn()}
  />,
);
check(
  "半三: 读面没答 modelQuotas ⇒ 等快照，不画输入框（没有读口就没有盲写）",
  waitMarkup.includes("等快照") && !waitMarkup.includes("日额度") && !waitMarkup.includes("月额度"),
  "waiting state",
);

// Source pins: the three whitelisted fields ride the one door, and the page
// adds NO range judge of its own (the Go setter's refusal is the single
// source; a page-side clamp would be the invented number the dispatch
// forbids).
const modelSrc = readFileSync("src/components/model-settings.tsx", "utf8");
check(
  "半三: 三枚字段字面量落页（model_quota_daily_micro / model_quota_monthly_micro / provider_plan_credit_total_micro）",
  modelSrc.includes('"model_quota_daily_micro"') && modelSrc.includes('"model_quota_monthly_micro"') &&
    modelSrc.includes('"provider_plan_credit_total_micro"'),
  "field literals",
);
check(
  "半三: quota 编辑走 performWrite（同一扇设置门，无第二通道）",
  modelSrc.includes("field,\n        provider,") && quotaSrcDoor(),
  "door",
);
function quotaSrcDoor(): boolean {
  // QuotaLeafInput hands the write to the caller's performWrite verbatim and
  // surfaces the answer verbatim (no re-wording of a refusal).
  return (
    modelSrc.includes("void run(`${field}:${provider}:${model ?? \"-\"}`, () =>") &&
    modelSrc.includes("onFeedback(fb);") &&
    modelSrc.includes("performWrite({")
  );
}
check(
  "半三: 页面不私设范围裁判（无页端负数拦截/上限数字）",
  !/quota.*-1|quotaMax|MAX_QUOTA|额度不能为负/.test(modelSrc),
  "no second judge",
);

// The host's range refusal surfaces verbatim through the same feedback line
// every settings card uses (the sentence below is the Go setter's own shape;
// the page renders it word for word, it never paraphrases a refusal).
const RANGE_REFUSAL =
  "面板请求被拒绝 [config.set r-40q]：config: llm.providers.acme.models.voice-in-only.quota_daily_micro: -1 超出范围（额度不能为负，0 = 不封顶）；未落盘";
const refusalMarkup = renderToStaticMarkup(
  <FeedbackLine feedback={{ tone: "error", text: RANGE_REFUSAL } as Feedback} />,
);
check(
  "半三: 越界拒绝句原话上屏（具名拒绝＝原话，不是转述）",
  refusalMarkup.includes(RANGE_REFUSAL) && refusalMarkup.includes("超出范围"),
  "refusal verbatim",
);
check("the refusal markup must hold zero emoji", !EMOJI_RE.test(refusalMarkup));

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-modelmgmt-40b: FAIL " + f);
  console.error(`render-modelmgmt-40b: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-modelmgmt-40b: OK, ticket 40 AC#3 后三半: voice role candidates filter on the read " +
    "face's real bits (non-voice absent, voice present, unanswered = null/unjudged never folded), " +
    "chain rows drag and share one pure reorder rule with the move buttons and the save writes the " +
    "current order, quota editor reads the current values verbatim (0 = 未设/不录入), waits without " +
    "inputs when the leg did not answer, rides the three whitelisted fields with no page-side range " +
    "judge, and surfaces the host's range refusal verbatim",
);
