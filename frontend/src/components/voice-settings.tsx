/* ============================================================================
   语音 settings page - v5 (R1-recheck residuals fix, 2026-11-16).
   ----------------------------------------------------------------------------
   Same interface as v1 (313-panel-ask.md §3): two cloud fallback chains and the
   realtime head. v2 made it a product surface (chain editors, plain language,
   pair rule at save time); v3 made the honest copy follow the backend truth.
   v4 is the audit-R1 fix: the two cards used to print a static 「已接入·重启
   生效」while the boot assembly refuses any chain whose models never declared
   capabilities.audio_in / audio_out - the read face did not even carry the
   bits, so the page had no way to know. Now the read face carries them
   (ConfigViewProvider.modelCapabilities), the card's tag and copy are
   CONDITIONAL on that verdict (缺能力位 = named, 醒目), the copy names the
   threshold (「需该模型在能力位里勾上 audio in/out」), and the save receipt
   branches on "can this chain assemble now" - the same ruler the boot runs,
   computed from the same read face. The capability bits themselves are
   switchable on the model page's model rows (model_audio_in / model_audio_out
   through the ticket-248 door).

   v5 closes the two residuals the R1 recheck caught:

     - residual 1 (the absent-key crash): the wire keeps modelCapabilities /
       modelsEnabled behind omitempty, so "the leg did not answer" arrives as
       an ABSENT key (undefined), not null - and the old `=== null` judgments
       threw on it (a zero-model provider row is factory-reachable,
       settings_313.go:154). Every read-face lookup here is now three-state
       (有值 / 腿没答 / 该模型没答), the waiting states render, and nothing
       throws.
     - residual 2 (the two invisible gates): the boot refuses a chain over
       THREE gates, the page modeled one. chainAssembly now judges all three
       off the read face (enabled=false via the new ConfigViewProvider.
       modelsEnabled read surface; a declared api_key_ref without an affirmed
       credential via apiKeyRef/credentialSet/credentialError), the tag
       says 已接入·重启生效 ONLY when all three pass, and credentialError
       renders on this page for the first time.

     - 听写 / 朗读 are now WIRED: the assembly root resolves the two chains at
       boot (cmd/wisp/voice_cloudchain_313_windows.go) and the resident voice
       loaders build speech.NewCloudASRChain / NewCloudTTSChain with the local
       cascade as the terminal fallback - the chains were "savable but read by
       no runtime path" before this leg (grep: zero NewCloud*Chain call sites).
       The chains are TierReload config, so the honest receipt is "保存后重启
       生效", and the placeholder warn-downgrade below no longer applies to
       these two cards.
     - 实时语音 stays a PLACEHOLDER card, exactly as v2 graded it: the
       RealtimeEngine seam is Reserved with no implementation, so its writes
       are still real-save/no-effect and keep the warn receipt.
     - 本机音频 graduates from the WAITING state (the snapshot keys and write fields landed Go-side 2026-11-11 (leg R); wired by this leg
       the two writable keys use real controls below, and the null reads still
       render the waiting tag per row (the pointer gate survived the wiring).

   A chain element must name a provider AND a model the config already lists
   (checkChainElement at load and on write), so the model control lists the
   provider's catalog - there is no typed-id path here, and a provider with no
   models says so instead of offering a doomed save.

   AC#1 (owner's default vendor): the chain editors and the realtime head
   prefill the owner-named default vendor (OWNER_DEFAULT_PROVIDER, the same
   constant the model page defines) ONLY as the draft state of an unset field -
   the read's saved value always wins, nothing writes until the user saves,
   and the name is only offered because a provider row actually carries it.

   STATE ALIGNMENT history (kept so the next leg does not re-litigate): the
   本机音频 keys - voice.aec.enabled feeds the barge-in monitor
   (cmd/wisp/resident_voice_windows.go:1992 -> :2010), audio.input_device
   resolves per capture open and hotplug reopen
   (cmd/wisp/resident_audio_windows.go:307), and audio.sample_rate is a
   load-time guard pinned to 16000 (internal/config/validate.go) - all have
   real backend readers. The v2 WAITING state was honest because the snapshot
   carried no values and the write face had no fields; leg R (2026-11-11)
   landed both halves on the Go side (snapshot pointers + the two flat write
   fields), and this leg wires the page to them. A null in the read payload
   still renders the waiting tag per row - never a guessed value.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useState, useRef } from "react";
import { ChevronDown, ChevronUp, Plus, RefreshCw, X } from "lucide-react";
import { Button } from "@/components/ai-native/button";
import { Switch } from "@/components/ai-native/switch";
import { ConfigReadGate } from "@/components/config-read-state";
import { COPY, configStatus, type Feedback } from "@/lib/config-feedback";
import type { ConfigView, ConfigViewProvider } from "@/lib/config-view";
import { performConfigWrite, type ConfigWireWrite } from "@/lib/config-actions";
import { OWNER_DEFAULT_PROVIDER } from "@/components/model-settings";
import {
  EmptyState,
  FeedbackLine,
  PairFields,
  SettingsCard,
  StateTag,
} from "@/components/settings-controls";

/** The chain directions and the capability bit each one consumes. The boot
    assembly refuses an ASR element whose model lacks capabilities.audio_in and
    a TTS element whose model lacks capabilities.audio_out
    (cloudcatalog.go:95-102); the read face's modelCapabilities carries the
    same two bits (audit R1/P3), so the page can judge a chain against the same
    ruler the boot will run instead of printing a static 已接入. */
type Direction = "asr" | "tts";

const DIRECTION_BIT: Record<Direction, "audioIn" | "audioOut"> = {
  asr: "audioIn",
  tts: "audioOut",
};
const DIRECTION_CONFIG_KEY: Record<Direction, "audio_in" | "audio_out"> = {
  asr: "audio_in",
  tts: "audio_out",
};

/* ---------------------------------------------------------------------------
   The read face's per-model answers, three-state (R1-recheck residual 1).
   Go's wire keeps modelCapabilities / modelsEnabled behind omitempty
   (config_handlers.go:776/:785), so an unanswering leg leaves BOTH keys off the
   wire and the bare-cast parse (config-view.ts) hands us `undefined`, not
   `null`. Absence and null are the same fact - "the leg did not answer" -
   and must never read as "every bit false" or "every model enabled": that
   folding is exactly the 能装配假象 both residuals exist to kill. A record
   that answered but skips this model id stays unanswered too.
   --------------------------------------------------------------------------- */

type CapsRecord = NonNullable<ConfigViewProvider["modelCapabilities"]>;

function capabilitiesOf(row: ConfigViewProvider): CapsRecord | null {
  return row.modelCapabilities ?? null;
}

function capsOf(caps: CapsRecord | null, model: string): CapsRecord[string] | null {
  return caps === null ? null : (caps[model] ?? null);
}

function enabledOf(row: ConfigViewProvider): Record<string, boolean> | null {
  return row.modelsEnabled ?? null;
}

function enabledOfModel(enabled: Record<string, boolean> | null, model: string): boolean | null {
  return enabled === null ? null : (enabled[model] ?? null);
}

/* ---------------------------------------------------------------------------
   票 40 AC#3 半一：capability 过滤的语音角色下拉（voice roles hide non-voice
   models）。位来源＝模型行的真字段（读面 modelCapabilities，Go 侧从
   spec.Capabilities.* 原样搬来），前端硬编码名单是禁区——本文件里没有任何
   模型名清单，候选只由这一位的真假决定。

   三态纪律沿用本文件 residual-1 的裁定：
     - 位已答 true  ⇒ 候选；
     - 位已答 false ⇒ 剔除（塞进语音角色的那个非语音模型就是它，剔除才不红）；
     - 位没答（模型条目缺） ⇒ 留在候选里——没答不是 false，把没答折成
       "非语音"正是 residual-1 禁的那一类猜；
     - 整张能力图没答（腿没答/等快照） ⇒ 返回 null，调用方不过滤：
       全隐藏＝把"没答"读成"全员 false"，同样禁。
   --------------------------------------------------------------------------- */

export type VoiceBit = "audioIn" | "audioOut" | "realtime";

export function voiceModelCandidates(
  row: ConfigViewProvider | null,
  bit: VoiceBit,
): string[] | null {
  if (row === null || row.models === null) return null;
  const caps = row.modelCapabilities ?? null;
  if (caps === null) return null;
  return row.models.filter((id) => {
    const bits = caps[id];
    if (!bits) return true;
    return bits[bit];
  });
}

/** One reorder of the chain draft: segment `from` moved to index `to`, both
    clamped, every other segment keeping its relative order. Exported pure for
    the render instrument; the drag handles and the two move buttons both land
    here, so there is exactly one reorder rule and no second one to drift. */
export function moveChainSegment(all: ChainSegment[], from: number, to: number): ChainSegment[] {
  if (from < 0 || from >= all.length || to < 0 || to >= all.length || from === to) return all;
  const next = [...all];
  const [moved] = next.splice(from, 1);
  next.splice(to, 0, moved);
  return next;
}

/** The assembly verdict for one SAVED chain, judged off the read face against
    the boot's gates in the boot's own order (internal/speech/cloudcatalog.go
    resolveChainElement): a chain element names an existing catalog pair, the
    model must be enabled, the model must DECLARE the direction's capability
    bit, and a declared api_key_ref must resolve to a credential the panel can
    affirm. The verdict kinds follow those gates:

      - "disabled"  an element's model is enabled=false (cloudcatalog.go:87-90),
                    named element-first with the boot's own vocabulary;
      - "nokey"     a provider declares api_key_ref but the read face cannot
                    affirm a resolvable credential behind it (CredentialSet
                    false, or CredentialError: the panel could not even check -
                    boot side both refuse, cloudcatalog.go:121-131);
      - "lacks"     the direction's capability bit is off (cloudcatalog.go:95-102),
                    or the element names a pair the catalog does not list;
      - "unknown"   a snapshot half was never answered, so no judgement is
                    invented (装配未判);
      - "assembles" every answered gate passes and none is missing - the ONLY
                    state allowed to say 已接入·重启生效.

    Precedence when several gates fail across the chain: the hard refusals
    win over "unknown" (the boot refuses regardless of what we could not
    read), and among them disabled > nokey > lacks, matching the boot's own
    check order within an element; the detail lines name every hit. */
type ChainAssembly =
  | { kind: "assembles" }
  | { kind: "unknown" }
  | { kind: "lacks"; missing: string[] }
  | { kind: "nokey"; providers: string[] }
  | { kind: "disabled"; models: string[] };

function chainAssembly(
  chain: string[] | null,
  providers: readonly ConfigViewProvider[],
  direction: Direction,
): ChainAssembly {
  const missing: string[] = [];
  const nokey: string[] = [];
  const disabled: string[] = [];
  let allAnswered = true;
  for (const el of chain ?? []) {
    const cut = el.indexOf("/");
    const provider = cut < 0 ? el : el.slice(0, cut);
    const model = cut < 0 ? "" : el.slice(cut + 1);
    const row = providers.find((p) => p.name === provider);
    if (provider === "" || model === "" || !row || row.models === null || !row.models.includes(model)) {
      missing.push(`${el}：目录里没有这一行`);
      continue;
    }
    // Gate 2 (cloudcatalog.go:87-90): enabled=false refuses by name.
    const enabled = enabledOfModel(enabledOf(row), model);
    if (enabled === false) {
      disabled.push(`${el} 已停用（enabled=false）`);
      continue;
    }
    if (enabled === null) allAnswered = false;
    // Gate 3 (cloudcatalog.go:95-102): the direction's capability bit.
    const bits = capsOf(capabilitiesOf(row), model);
    if (bits === null) {
      allAnswered = false;
    } else if (!bits[DIRECTION_BIT[direction]]) {
      missing.push(`${el} 缺 capabilities.${DIRECTION_CONFIG_KEY[direction]}`);
    }
    // Gate 4 (cloudcatalog.go:121-131): a declared api_key_ref the read face
    // cannot affirm a resolvable credential behind. An EMPTY ref passes - the
    // boot skips resolution for one too (the if ref != "" gate).
    if (row.apiKeyRef !== "" && (row.credentialError !== "" || !row.credentialSet)) {
      nokey.push(
        row.credentialError !== ""
          ? `${provider} 无法核对密钥（${row.credentialError}）`
          : `${provider} 缺 API Key（api_key_ref ${row.apiKeyRef} 未就绪）`,
      );
    }
  }
  if (disabled.length > 0) return { kind: "disabled", models: disabled };
  if (nokey.length > 0) return { kind: "nokey", providers: nokey };
  if (missing.length > 0) return { kind: "lacks", missing };
  return allAnswered ? { kind: "assembles" } : { kind: "unknown" };
}

/** 票 313 AC#1 的语音半边：机主点名的默认供应商，与 model-settings.tsx 的
    OWNER_DEFAULT_PROVIDER 是同一个常量（一处定义，两页引用，不许漂）。
    语音卡上的用法比模型页更窄：它只是「未设置字段的草稿初值」——链条为空时
    第一段的供应商预选、realtime 未配置时的草稿。门槛是这行供应商真的存在
    （链元素必须点名已建行的 provider/model，给一个不存在行的预选就是必然
    保存失败的坑），用户存过的值永远原样显示，从不被默认值盖写。 */
function voiceDefaultProvider(providers: readonly ConfigViewProvider[]): string {
  return providers.some((p) => p.name === OWNER_DEFAULT_PROVIDER)
    ? OWNER_DEFAULT_PROVIDER
    : "";
}

/** The realtime card's receipt: the save is real, the effect is not (the
    RealtimeEngine seam is still Reserved). The generic translated receipt
    ("已保存，重启后生效") would promise an effect no runtime path can deliver,
    so the ok answer is restated as a warn - the v2 tone downgrade, kept
    exactly for the card that still needs it. */
function placeholderReceipt(feedback: Feedback): Feedback {
  return feedback.tone === "ok"
    ? { tone: "warn", text: "已写入配置文件；这一项还没有接到运行路径，当前不生效" }
    : feedback;
}


export function VoiceSettings({
  heading = "语音",
  lede = "云端听写链与云端朗读链已接入常驻语音的装配——门槛：链里每个模型都要在能力位里勾上对应的 audio in / audio out（「模型」页的模型行），保存后重启生效，链尽回退本地引擎；实时语音的引擎尚未实现仍是占位；本机音频三项显示真当前值，回声消除与输入设备可改可存。",
}: {
  heading?: string;
  lede?: string;
}) {
  return (
    <section aria-label="语音" className="flex min-w-0 flex-col gap-1">
      <h2 className="text-[15px] font-semibold text-ink">{heading}</h2>
      <p className="text-[12.5px] text-ink-2">{lede}</p>
      <ConfigReadGate>
        {(view, _prose, reload) => (
          <VoiceSettingsView
            onReload={reload}
            performWrite={performConfigWrite}
            view={view}
          />
        )}
      </ConfigReadGate>
    </section>
  );
}

export interface VoiceSettingsViewProps {
  view: ConfigView;
  performWrite: (write: ConfigWireWrite) => Promise<Feedback>;
  onReload: () => void;
}

export function VoiceSettingsView({ view, performWrite, onReload }: VoiceSettingsViewProps) {
  const status = configStatus(view);
  const [busy, setBusy] = useState<string | null>(null);

  async function run(key: string, fn: () => Promise<Feedback>): Promise<Feedback> {
    setBusy(key);
    const feedback = await fn();
    setBusy(null);
    onReload();
    return feedback;
  }

  return (
    <div className="mt-2 flex min-w-0 flex-col gap-3">
      <div className="flex items-center gap-2">
        <StateTag text={status.label} tone={status.tone} />
        <span className="min-w-0 flex-1" />
        <Button aria-label="刷新" onClick={onReload} size="xs" variant="quiet">
          <RefreshCw aria-hidden="true" size={11} strokeWidth={2} />
          刷新
        </Button>
      </div>

      {!view.readable ? (
        <SettingsCard title="配置">
          <p className="text-[12px] leading-[1.7] text-ink-2">暂时读不到配置，请稍后刷新。</p>
        </SettingsCard>
      ) : view.providers.length === 0 ? (
        <SettingsCard title="语音">
          <EmptyState text="先在「模型」页添加一个供应商，再回来配置语音。" />
        </SettingsCard>
      ) : (
        <>
          {/* R1-recheck residual 2: the credential read face's cannot-check
              state (credentialError) rode every provider row but no voice
              card rendered it - the one page that assembles chains over that
              key never said why the key is not affirmed. Rendered verbatim:
              it is a state the leg answered, not a guess. */}
          {view.providers
            .filter((p) => p.credentialError !== "")
            .map((p) => (
              <p className="text-[11.5px] text-orange" key={p.name}>
                无法核对密钥：{p.name}（{p.credentialError}）
              </p>
            ))}
          <ChainEditor
            busy={busy}
            defaultProvider={voiceDefaultProvider(view.providers)}
            description="已接入：常驻语音启动时按这条链的顺序装配云端听写，链尽回退本地引擎。链段可拖动（或上移/下移）重排，保存按现序写入配置文件，重启后按同序读回。门槛：链里每个模型都要在能力位里勾上 audio in（「模型」页的模型行），保存后重启生效。"
            direction="asr"
            field="voice_cloud_asr_chain"
            initial={view.voice.cloudASRChain}
            performWrite={performWrite}
            providers={view.providers}
            run={run}
            title="听写（ASR）"
          />
          <ChainEditor
            busy={busy}
            defaultProvider={voiceDefaultProvider(view.providers)}
            description="已接入：常驻语音启动时按这条链的顺序装配云端朗读，链尽回退本地引擎。链段可拖动（或上移/下移）重排，保存按现序写入配置文件，重启后按同序读回。门槛：链里每个模型都要在能力位里勾上 audio out（「模型」页的模型行），保存后重启生效。"
            direction="tts"
            field="voice_cloud_tts_chain"
            initial={view.voice.cloudTTSChain}
            performWrite={performWrite}
            providers={view.providers}
            run={run}
            title="朗读（TTS）"
          />
          <RealtimeCard
            busy={busy}
            performWrite={performWrite}
            providers={view.providers}
            run={run}
            view={view}
          />
          <LocalAudioCard busy={busy} performWrite={performWrite} run={run} view={view} />
        </>
      )}
    </div>
  );
}

/* ---------------------------------------------------------------------------
   One chain.
   --------------------------------------------------------------------------- */

interface ChainSegment {
  provider: string;
  model: string;
}

/** "provider/model" as the config spells a chain element; a hand-written
    element without a slash keeps its whole text in the provider half so the
    editor never silently drops it. */
function parseChain(chain: string[] | null): ChainSegment[] {
  if (chain === null) return [];
  return chain.map((el) => {
    const cut = el.indexOf("/");
    return cut < 0 ? { provider: el, model: "" } : { provider: el.slice(0, cut), model: el.slice(cut + 1) };
  });
}

function chainText(chain: string[] | null): string {
  if (chain === null || chain.length === 0) return "未设置";
  return chain.join(", ");
}

function ChainEditor({
  title,
  description,
  direction,
  field,
  initial,
  providers,
  busy,
  run,
  performWrite,
  defaultProvider,
}: {
  title: string;
  description: string;
  /** Which direction this card assembles; picks the capability bit the boot
      refuses on (the audit-R1 conditional copy's 点名门槛). */
  direction: Direction;
  field: string;
  initial: string[] | null;
  providers: readonly ConfigViewProvider[];
  busy: string | null;
  run: (key: string, fn: () => Promise<Feedback>) => Promise<Feedback>;
  performWrite: (write: ConfigWireWrite) => Promise<Feedback>;
  /** AC#1's draft default for an unset chain: the owner vendor's name as the
      first empty segment's provider, only when a provider row carries it.
      Initial state only - no effect re-imposes it, and the read's saved
      chain always parses in verbatim below. */
  defaultProvider: string;
}) {
  const [segments, setSegments] = useState<ChainSegment[]>(() => {
    const parsed = parseChain(initial);
    if (parsed.length > 0) return parsed;
    return [{ provider: defaultProvider, model: "" }];
  });
  const [feedback, setFeedback] = useState<Feedback | null>(null);

  // The saved chain's assembly verdict, judged off the read face (audit R1 +
  // R1-recheck residual 2): the tag stops being a static 已接入 and states
  // what is actually true about THIS chain against the ruler the boot runs -
  // all three of its gates, not only the capability bit. 已接入·重启生效 is
  // reachable ONLY from the assembles verdict, where every answered gate
  // passed and none is missing.
  const savedVerdict = chainAssembly(initial, providers, direction);
  const savedTag = ((): { text: string; tone: "ok" | "warn" | "neutral" | "error" } => {
    if (initial === null || initial.length === 0) return { text: "未设置", tone: "neutral" };
    switch (savedVerdict.kind) {
      case "disabled":
        return { text: "装配不了·模型已停用", tone: "warn" };
      case "nokey":
        return { text: "装配不了·缺 API Key", tone: "warn" };
      case "lacks":
        return { text: "装配不了·缺能力位", tone: "warn" };
      case "unknown":
        return { text: "已接入·装配未判（等快照）", tone: "neutral" };
      default:
        return { text: "已接入·重启生效", tone: "ok" };
    }
  })();
  // Every gate's hit list, element-first in the boot's own vocabulary: the
  // tag names the first gate the boot would refuse on, this line names them
  // all, so a chain failing two gates hides neither.
  const refusalLines = ((): string[] => {
    switch (savedVerdict.kind) {
      case "disabled":
        return savedVerdict.models;
      case "nokey":
        return savedVerdict.providers;
      case "lacks":
        return savedVerdict.missing;
      default:
        return [];
    }
  })();

  function update(index: number, next: Partial<ChainSegment>) {
    setSegments((all) => all.map((seg, i) => (i === index ? { ...seg, ...next } : seg)));
  }

  // 票 40 AC#3 半二：拖拽排序。dragging ref 存起拖的那一段的下标，drop 时经
  // moveChainSegment 重排；数组顺序就是保存时写进 config.toml 的链顺序（save
  // 照 segments 的现序逐段拼值走既有写门），拖完不保存不改盘上。上移/下移
  // 按钮是同一条重排规则的第二个入口（键盘可达，仪器钉得住）。
  const dragFrom = useRef<number | null>(null);
  const [dragOver, setDragOver] = useState<number | null>(null);

  function moveTo(from: number, to: number) {
    setSegments((all) => moveChainSegment(all, from, to));
  }

  function save() {
    if (segments.length === 0) {
      setFeedback({ tone: "error", text: COPY.chainEmpty });
      return;
    }
    for (const seg of segments) {
      if (seg.provider === "" || seg.model === "") {
        const row = providers.find((p) => p.name === seg.provider);
        const noModels = seg.provider !== "" && row !== undefined && (row.models === null || row.models.length === 0);
        setFeedback({
          tone: "error",
          text: noModels ? "这家供应商还没有模型" : COPY.chainSegmentMissing,
        });
        return;
      }
    }
    const savedChain = segments.map((seg) => `${seg.provider}/${seg.model}`);
    void run(`${field}`, () =>
      performWrite({
        field,
        value: savedChain.join(","),
      }).then((feedback) => {
        if (feedback.tone !== "ok") return feedback;
        // The receipt branches on "can this chain assemble NOW" (audit R1 +
        // R1-recheck residual 2): the same verdict the boot will compute -
        // all three gates - from the read face. Each branch names the gate it
        // refused on and the fix the user can make.
        const verdict = chainAssembly(savedChain, providers, direction);
        if (verdict.kind === "disabled") {
          return {
            tone: "error",
            text: `已写入配置文件，但这条链现在装配不了：${verdict.models.join("；")}。把这个模型启用或换一个模型，再回来存这条链，重启生效。`,
          };
        }
        if (verdict.kind === "nokey") {
          return {
            tone: "error",
            text: `已写入配置文件，但这条链现在装配不了：${verdict.providers.join("；")}。在「模型」页把这一家的密钥录好，再回来存这条链，重启生效。`,
          };
        }
        if (verdict.kind === "lacks") {
          return {
            tone: "error",
            text: `已写入配置文件，但这条链现在装配不了：${verdict.missing.join("；")}。在「模型」页把缺的那一位打开并保存，再回来存这条链，重启生效。`,
          };
        }
        if (verdict.kind === "unknown") {
          return {
            tone: "warn",
            text: "已保存，重启后生效；这条链现在能不能装配判不了（快照还没读到能力位或启用位）。",
          };
        }
        return { tone: "ok", text: "已接入：保存后重启生效" };
      }),
    ).then((feedback) => setFeedback(feedback));
  }

  return (
    <SettingsCard title={title} description={description}>
      <div className="flex flex-col gap-2">
        <div className="flex flex-wrap items-center gap-2">
          <StateTag text={savedTag.text} tone={savedTag.tone} />
          {refusalLines.length > 0 ? (
            <span className="text-[11.5px] text-orange">现在装配不了：{refusalLines.join("；")}</span>
          ) : null}
        </div>
        <div className="flex flex-col gap-1.5">
          {segments.map((seg, index) => (
            <div
              className={
                "flex flex-wrap items-center gap-2 rounded-control px-1 py-0.5" +
                (dragOver === index && dragFrom.current !== null && dragFrom.current !== index
                  ? " bg-inset"
                  : "")
              }
              data-chain-segment={index + 1}
              draggable
              key={`${index}-${seg.provider}-${seg.model}`}
              onDragEnd={() => {
                dragFrom.current = null;
                setDragOver(null);
              }}
              onDragOver={(e) => {
                e.preventDefault();
                setDragOver(index);
              }}
              onDragStart={() => {
                dragFrom.current = index;
              }}
              onDrop={(e) => {
                e.preventDefault();
                const from = dragFrom.current;
                dragFrom.current = null;
                setDragOver(null);
                if (from !== null) moveTo(from, index);
              }}
            >
              <span className="w-8 shrink-0 cursor-grab font-mono text-[11px] text-ink-3" title="拖动调整顺序">
                {index + 1}
              </span>
              <PairFields
                manualFallback={false}
                model={seg.model}
                modelCandidates={(row) => voiceModelCandidates(row, DIRECTION_BIT[direction])}
                onModel={(next) => update(index, { model: next })}
                onProvider={(next) => update(index, { provider: next, model: "" })}
                provider={seg.provider}
                providers={providers}
              />
              <Button
                aria-label={`第 ${index + 1} 段上移`}
                disabled={index === 0 || busy !== null}
                onClick={() => moveTo(index, index - 1)}
                size="xs"
                variant="quiet"
              >
                <ChevronUp aria-hidden="true" size={12} strokeWidth={2} />
              </Button>
              <Button
                aria-label={`第 ${index + 1} 段下移`}
                disabled={index === segments.length - 1 || busy !== null}
                onClick={() => moveTo(index, index + 1)}
                size="xs"
                variant="quiet"
              >
                <ChevronDown aria-hidden="true" size={12} strokeWidth={2} />
              </Button>
              <Button
                aria-label={`移除第 ${index + 1} 段`}
                onClick={() => setSegments((all) => all.filter((_, i) => i !== index))}
                size="xs"
                variant="quiet"
              >
                <X aria-hidden="true" size={12} strokeWidth={2} />
              </Button>
            </div>
          ))}
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <Button
            onClick={() => setSegments((all) => [...all, { provider: "", model: "" }])}
            size="xs"
            variant="secondary"
          >
            <Plus aria-hidden="true" size={12} strokeWidth={2} />
            添加一段
          </Button>
          <Button disabled={busy !== null} onClick={save} size="sm" variant="primary">
            {busy === field ? COPY.saving : COPY.save}
          </Button>
        </div>
        <p className="text-[11.5px] text-ink-3">当前：{chainText(initial)}</p>
        <FeedbackLine feedback={feedback} />
      </div>
    </SettingsCard>
  );
}

/* ---------------------------------------------------------------------------
   本机音频 (voice.aec.enabled / audio.input_device / audio.sample_rate).
   --------------------------------------------------------------------------- */

/** The 本机音频 card, wired (this leg, 2026-11-16): the read payload's three
    current values render verbatim, and the two keys the write face owns
    (voice_aec_enabled / audio_input_device, config_write_313.go) become real
    controls riding the same performConfigWrite path as every other settings
    field - the ticket-248 protected config.set setter, no second route.
    The third key stays a read-only fact on purpose: audio.sample_rate is a
    load-time 16000 guard (internal/config/validate.go refuses every other
    value before a byte is written), so a control here would be ticket 83's
    lie wearing a slider - the current number shows, the guard reason is
    said, and nothing saves. AEC adds the capability dimension (ticket 452
    AC#2): aecBuildAvailable=false greys the switch out and says so, because
    "你想不想开" and "这枚二进制能不能" are two facts, not one. */
function LocalAudioCard({
  view,
  busy,
  run,
  performWrite,
}: {
  view: ConfigView;
  busy: string | null;
  run: (key: string, fn: () => Promise<Feedback>) => Promise<Feedback>;
  performWrite: (write: ConfigWireWrite) => Promise<Feedback>;
}) {
  // Draft state, seeded from the read's answered values. A null read renders
  // the waiting rows below and keeps the controls out of the tree entirely.
  const [aecDraft, setAecDraft] = useState<boolean>(view.voice.aecEnabled ?? false);
  const [deviceDraft, setDeviceDraft] = useState<string>(view.voice.inputDevice ?? "");
  const [aecFeedback, setAecFeedback] = useState<Feedback | null>(null);
  const [deviceFeedback, setDeviceFeedback] = useState<Feedback | null>(null);

  // The 本机音频 half's answered flags, three-state (R1-recheck residual 1's
  // own class): the four keys ride reader-gated pointers behind omitempty
  // (config_handlers.go:787-794), so an unanswering leg leaves them OFF the
  // wire and the bare-cast parse hands us `undefined`, not `null`. Judging
  // `=== null` let absence read as an answered "关"/default - a fabricated
  // fact. `?? null` folds absence and null into "the leg did not answer".
  const aecRead = view.voice.aecEnabled ?? null;
  const deviceRead = view.voice.inputDevice ?? null;
  const rateRead = view.voice.sampleRate ?? null;
  const aecKnown = aecRead !== null;
  const deviceKnown = deviceRead !== null;
  const rateKnown = rateRead !== null;
  // null = the leg did not answer the capability bit: keep the control rather
  // than pretend an absence that was never measured. false = answered absence.
  const aecBuildable = (view.voice.aecBuildAvailable ?? null) !== false;
  const deviceShown =
    view.voice.inputDevice === "" || view.voice.inputDevice === "default"
      ? "default（跟随系统默认）"
      : (view.voice.inputDevice ?? "");

  function saveAec() {
    void run("voice_aec_enabled", () =>
      performWrite({ field: "voice_aec_enabled", value: aecDraft ? "true" : "false" }),
    ).then(setAecFeedback);
  }

  function saveDevice() {
    const v = deviceDraft.trim();
    if (v === "") {
      setDeviceFeedback({ tone: "error", text: "输入设备名不能为空；跟随系统默认请填 default" });
      return;
    }
    void run("audio_input_device", () =>
      performWrite({ field: "audio_input_device", value: v }),
    ).then(setDeviceFeedback);
  }

  return (
    <SettingsCard
      title="本机音频"
      description="三项显示配置里的真当前值；回声消除与输入设备可改可存（保存后重启生效），采样率是采集链的加载期守卫（只跑 16000），只显示不提供修改。"
    >
      <div className="flex flex-col gap-2.5">
        <div className="flex flex-wrap items-center gap-x-2 gap-y-1">
          <span className="text-[12px] font-medium text-ink">回声消除</span>
          <span className="font-mono text-[10.5px] text-ink-3">voice.aec.enabled</span>
          <span className="min-w-0 flex-1" />
          {aecKnown ? (
            <>
              <span className="text-[11.5px] text-ink-2">
                当前：{view.voice.aecEnabled ? "开" : "关"}
              </span>
              {aecBuildable ? (
                <Switch checked={aecDraft} label="回声消除开关" onChange={setAecDraft} />
              ) : (
                <StateTag text="本包无 AEC" tone="warn" />
              )}
              <Button
                disabled={busy !== null || !aecBuildable}
                onClick={saveAec}
                size="sm"
                variant="primary"
              >
                {busy === "voice_aec_enabled" ? COPY.saving : COPY.save}
              </Button>
            </>
          ) : (
            <StateTag text="等快照" tone="neutral" />
          )}
        </div>
        {!aecKnown ? (
          <p className="text-[11.5px] leading-[1.6] text-ink-2">
            后端已接：喂给打断监测的回声消除开关（重启后生效）。
          </p>
        ) : aecBuildable ? (
          <p className="text-[11.5px] leading-[1.6] text-ink-3">
            喂给打断监测的回声消除开关；写入走 config.set 的受保护写路径，重启后生效。
          </p>
        ) : (
          <p className="text-[11.5px] leading-[1.6] text-orange">
            这枚二进制没有编译回声消除（缺 webrtcap 构建标签），开关只能保存配置意图，当前不会生效。
          </p>
        )}
        <FeedbackLine feedback={aecFeedback} />

        <div className="flex flex-wrap items-center gap-x-2 gap-y-1">
          <span className="text-[12px] font-medium text-ink">输入设备</span>
          <span className="font-mono text-[10.5px] text-ink-3">audio.input_device</span>
          <span className="min-w-0 flex-1" />
          {deviceKnown ? (
            <span className="text-[11.5px] text-ink-2">当前：{deviceShown}</span>
          ) : (
            <StateTag text="等快照" tone="neutral" />
          )}
        </div>
        {deviceKnown ? (
          <div className="flex flex-wrap items-center gap-2">
            <input
              aria-label="输入设备名"
              className="h-[30px] w-[240px] rounded-control bg-inset px-2.5 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3"
              onChange={(e) => setDeviceDraft(e.currentTarget.value)}
              placeholder="设备名（default = 跟随系统默认）"
              value={deviceDraft}
            />
            <Button disabled={busy !== null} onClick={saveDevice} size="sm" variant="primary">
              {busy === "audio_input_device" ? COPY.saving : COPY.save}
            </Button>
            <span className="min-w-0 flex-1 basis-40 text-[11px] leading-[1.6] text-ink-3">
              每次开麦与热插拔重开都按这个名字解析设备；解析不到会在开麦时按名拒绝并列出在场端点。
            </span>
          </div>
        ) : (
          <p className="text-[11.5px] leading-[1.6] text-ink-2">
            后端已接：每次开麦与热插拔重开都按这个名字解析设备。
          </p>
        )}
        <FeedbackLine feedback={deviceFeedback} />

        <div className="flex flex-wrap items-center gap-x-2 gap-y-1">
          <span className="text-[12px] font-medium text-ink">采样率</span>
          <span className="font-mono text-[10.5px] text-ink-3">audio.sample_rate</span>
          <span className="min-w-0 flex-1" />
          {rateKnown ? (
            <span className="text-[11.5px] text-ink-2">当前：{view.voice.sampleRate} Hz</span>
          ) : (
            <StateTag text="等快照" tone="neutral" />
          )}
          <StateTag text="守卫·只读" tone="neutral" />
        </div>
        <p className="text-[11.5px] leading-[1.6] text-ink-3">
          后端守卫：采集链只跑 16000，配置里写别的值会被加载闸按名拒绝；本页不提供修改，也不写坏值。
        </p>
      </div>
    </SettingsCard>
  );
}

/* ---------------------------------------------------------------------------
   实时语音.
   --------------------------------------------------------------------------- */

function RealtimeCard({
  view,
  providers,
  busy,
  run,
  performWrite,
}: {
  view: ConfigView;
  providers: readonly ConfigViewProvider[];
  busy: string | null;
  run: (key: string, fn: () => Promise<Feedback>) => Promise<Feedback>;
  performWrite: (write: ConfigWireWrite) => Promise<Feedback>;
}) {
  // AC#1's draft default for an unset head: the owner vendor when no realtime
  // provider is saved. The read's saved value wins verbatim; the default is
  // read once here as the initial state and never re-imposed by an effect.
  const [provider, setProvider] = useState(
    view.voice.realtimeProvider || voiceDefaultProvider(providers),
  );
  const [model, setModel] = useState(view.voice.realtimeModel);
  const [feedback, setFeedback] = useState<Feedback | null>(null);

  const current =
    view.voice.realtimeProvider !== "" || view.voice.realtimeModel !== ""
      ? `${view.voice.realtimeProvider} / ${view.voice.realtimeModel}`
      : "未设置";

  function save() {
    if (provider === "" || model === "") {
      setFeedback({ tone: "error", text: COPY.pairMissing });
      return;
    }
    void run("realtime", async () => {
      const head = await performWrite({ field: "voice_realtime_provider", value: provider });
      if (head.tone !== "ok") return head;
      const tail = await performWrite({ field: "voice_realtime_model", value: model });
      return tail.tone === "ok" ? placeholderReceipt(tail) : tail;
    }).then(setFeedback);
  }

  return (
    <SettingsCard title="实时语音" description="占位卡：实时通话引擎尚未实现，这里配置的供应商与模型暂时没有通话链使用；写入会保存，接入前不生效。">
      <div className="flex flex-col gap-2">
        <StateTag text="占位·当前不生效" tone="warn" />
        <div className="flex flex-wrap items-center gap-2">
          <PairFields
            manualFallback
            model={model}
            // 票 40 AC#3 半一：realtime 也是语音角色，候选只列声明 realtime 位
            // 的模型（位来源＝读面 modelCapabilities，非硬编码名单）。
            modelCandidates={(row) => voiceModelCandidates(row, "realtime")}
            onModel={setModel}
            onProvider={(next) => {
              setProvider(next);
              setModel("");
            }}
            provider={provider}
            providers={providers}
          />
          <Button disabled={busy !== null} onClick={save} size="sm" variant="primary">
            {busy === "realtime" ? COPY.saving : COPY.save}
          </Button>
        </div>
        <p className="text-[11.5px] text-ink-3">当前：{current}</p>
        <FeedbackLine feedback={feedback} />
      </div>
    </SettingsCard>
  );
}
