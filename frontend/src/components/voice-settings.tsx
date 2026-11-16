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

import { useState } from "react";
import { Plus, RefreshCw, X } from "lucide-react";
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
  lede = "云端听写链与云端朗读链已接入常驻语音的装配（保存后重启生效，链尽回退本地引擎）；实时语音的引擎尚未实现仍是占位；本机音频三项显示真当前值，回声消除与输入设备可改可存。",
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

  const aecKnown = view.voice.aecEnabled !== null;
  const deviceKnown = view.voice.inputDevice !== null;
  const rateKnown = view.voice.sampleRate !== null;
  // null = the leg did not answer the capability bit: keep the control rather
  // than pretend an absence that was never measured. false = answered absence.
  const aecBuildable = view.voice.aecBuildAvailable !== false;
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
