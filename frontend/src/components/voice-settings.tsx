/* ============================================================================
   语音 settings page - v3 (ticket 313's wiring leg, 2026-11-16).
   ----------------------------------------------------------------------------
   Same interface as v1 (313-panel-ask.md §3): two cloud fallback chains and the
   realtime head. v2 made it a product surface (chain editors, plain language,
   pair rule at save time); v3 makes the honest copy follow the backend truth:

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
     - 本机音频 keeps the WAITING state (snapshot keys and write fields still
       absent backend-side; the v2 STATE ALIGNMENT note below names the three
       keys); unchanged by this leg.

   A chain element must name a provider AND a model the config already lists
   (checkChainElement at load and on write), so the model control lists the
   provider's catalog - there is no typed-id path here, and a provider with no
   models says so instead of offering a doomed save.

   AC#1 (owner's default vendor): the chain editors and the realtime head
   prefill the owner-named default vendor (OWNER_DEFAULT_PROVIDER, the same
   constant the model page defines) ONLY as the draft state of an unset field -
   the read's saved value always wins, nothing writes until the user saves,
   and the name is only offered because a provider row actually carries it.

   STATE ALIGNMENT kept from v2 (2026-11-11, still true): the 本机音频 keys -
   voice.aec.enabled feeds the barge-in monitor
   (cmd/wisp/resident_voice_windows.go:1992 -> :2010), audio.input_device
   resolves per capture open and hotplug reopen
   (cmd/wisp/resident_audio_windows.go:307), and audio.sample_rate is a
   load-time guard pinned to 16000 (internal/config/validate.go) - all have
   real backend readers, but the config.get snapshot's voice payload carries
   five keys only (pinned in config_313_test.go) and the config.set whitelist
   has no write field for any of the three. So that card renders the WAITING
   state; drawing switchable controls now would re-open the C-class lie.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useState } from "react";
import { Plus, RefreshCw, X } from "lucide-react";
import { Button } from "@/components/ai-native/button";
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
  lede = "云端听写链与云端朗读链已接入常驻语音的装配（保存后重启生效，链尽回退本地引擎）；实时语音的引擎尚未实现仍是占位；本机音频的后端读者已接真，面板先等快照。",
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
          <ChainEditor
            busy={busy}
            defaultProvider={voiceDefaultProvider(view.providers)}
            description="已接入：常驻语音启动时按这条链的顺序装配云端听写，链尽回退本地引擎；保存后重启生效。"
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
            description="已接入：常驻语音启动时按这条链的顺序装配云端朗读，链尽回退本地引擎；保存后重启生效。"
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
          <LocalAudioCard />
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

  function update(index: number, next: Partial<ChainSegment>) {
    setSegments((all) => all.map((seg, i) => (i === index ? { ...seg, ...next } : seg)));
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
    void run(`${field}`, () =>
      performWrite({
        field,
        value: segments.map((seg) => `${seg.provider}/${seg.model}`).join(","),
      }),
    ).then((feedback) => setFeedback(feedback));
  }

  return (
    <SettingsCard title={title} description={description}>
      <div className="flex flex-col gap-2">
        <StateTag text="已接入·重启生效" tone="ok" />
        <div className="flex flex-col gap-1.5">
          {segments.map((seg, index) => (
            <div className="flex flex-wrap items-center gap-2" key={`${index}-${seg.provider}`}>
              <span className="w-8 shrink-0 font-mono text-[11px] text-ink-3">{index + 1}</span>
              <PairFields
                manualFallback={false}
                model={seg.model}
                onModel={(next) => update(index, { model: next })}
                onProvider={(next) => update(index, { provider: next, model: "" })}
                provider={seg.provider}
                providers={providers}
              />
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

/** One row of the local-audio card: the key as the config spells it, the
    backend fact (why it is real today), and the waiting state. */
const LOCAL_AUDIO_ROWS: readonly { key: string; name: string; fact: string }[] = [
  {
    key: "voice.aec.enabled",
    name: "回声消除",
    fact: "后端已接：喂给打断监测的回声消除开关（重启后生效）",
  },
  {
    key: "audio.input_device",
    name: "输入设备",
    fact: "后端已接：每次开麦与热插拔重开都按这个名字解析设备",
  },
  {
    key: "audio.sample_rate",
    name: "采样率",
    fact: "后端守卫：采集链只跑 16000，配置里写别的值会被拒",
  },
];

/** The three local-audio keys the backend already reads for real, but the
    panel cannot show or change yet: the config.get snapshot's voice payload
    carries five keys only (pinned in config_313_test.go), none of them these
    three, and the config.set whitelist has no write field for them either.
    So this card is the honest WAITING state - the backend fact as words plus
    "等快照" - never a guessed value, never a control that looks switchable.
    When Go adds the snapshot keys and the write fields, this card graduates
    to real controls; drawing them now would re-open the C-class lie. */
function LocalAudioCard() {
  return (
    <SettingsCard
      title="本机音频"
      description="这三项的后端读者/守卫已接（真），但配置读口的快照还没有它们的当前值，写入路由也未接。快照补上之前这里显示等待态，不画能点但改不动的控件。"
    >
      <div className="flex flex-col gap-1.5">
        {LOCAL_AUDIO_ROWS.map((row) => (
          <div className="flex flex-wrap items-baseline gap-x-2 gap-y-0.5" key={row.key}>
            <span className="text-[12px] font-medium text-ink">{row.name}</span>
            <span className="font-mono text-[10.5px] text-ink-3">{row.key}</span>
            <span className="min-w-0 flex-1 basis-40 text-[11.5px] leading-[1.6] text-ink-2">
              {row.fact}
            </span>
            <StateTag text="等快照" tone="neutral" />
          </div>
        ))}
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
