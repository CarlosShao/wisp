/* ============================================================================
   语音 settings page - v2 (owner review 2026-10-14; see model-settings.tsx's
   header for what changed and why).
   ----------------------------------------------------------------------------
   Same interface as v1 (313-panel-ask.md §3): two cloud fallback chains and the
   realtime head. v2 makes it a product surface:

     - 听写 / 朗读 are CHAIN EDITORS (rows of provider+model, add/remove), which
       is what the field actually is - a comma-separated fallback list. v1 wrote
       a whole chain from one typed string, which silently overwrote a multi-hop
       chain the operator may have written by hand.
     - everything reads in plain language; the enabled-switch gap is stated as a
       scope line, not as internal jargon.
     - the pair rule is checked at SAVE time, per segment.

   A chain element must name a provider AND a model the config already lists
   (checkChainElement at load and on write), so the model control lists the
   provider's catalog - there is no typed-id path here, and a provider with no
   models says so instead of offering a doomed save.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useState } from "react";
import { Plus, RefreshCw, X } from "lucide-react";
import { Button } from "@/components/ai-native/button";
import { ConfigReadGate } from "@/components/config-read-state";
import { COPY, configStatus, type Feedback } from "@/lib/config-feedback";
import type { ConfigView, ConfigViewProvider } from "@/lib/config-view";
import { performConfigWrite, type ConfigWireWrite } from "@/lib/config-actions";
import {
  EmptyState,
  FeedbackLine,
  PairFields,
  SettingsCard,
  StateTag,
} from "@/components/settings-controls";


export function VoiceSettings({
  heading = "语音",
  lede = "配置听写、朗读与实时语音使用的供应商和模型。",
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
            description="听写优先走这里列出的顺序，都不可用时回退到本地引擎。"
            field="voice_cloud_asr_chain"
            initial={view.voice.cloudASRChain}
            performWrite={performWrite}
            providers={view.providers}
            run={run}
            title="听写（ASR）"
          />
          <ChainEditor
            busy={busy}
            description="朗读优先走这里列出的顺序，都不可用时回退到本地引擎。"
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
}: {
  title: string;
  description: string;
  field: string;
  initial: string[] | null;
  providers: readonly ConfigViewProvider[];
  busy: string | null;
  run: (key: string, fn: () => Promise<Feedback>) => Promise<Feedback>;
  performWrite: (write: ConfigWireWrite) => Promise<Feedback>;
}) {
  const [segments, setSegments] = useState<ChainSegment[]>(() => {
    const parsed = parseChain(initial);
    return parsed.length > 0 ? parsed : [{ provider: "", model: "" }];
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
    ).then(setFeedback);
  }

  return (
    <SettingsCard title={title} description={description}>
      <div className="flex flex-col gap-2">
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
  const [provider, setProvider] = useState(view.voice.realtimeProvider);
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
      return tail.tone === "ok" ? { tone: "ok" as const, text: "已保存，重启后生效" } : tail;
    }).then(setFeedback);
  }

  return (
    <SettingsCard title="实时语音" description="陪聊模式的实时通话；通话开关由应用控制，这里配置供应商与模型。">
      <div className="flex flex-col gap-2">
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
