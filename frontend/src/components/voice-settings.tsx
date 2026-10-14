/* ============================================================================
   语音 settings page - ticket 313's page half for the voice category.
   ----------------------------------------------------------------------------
   WHAT THIS IS. The 语音 category wired to the same settings door as the 模型
   page: one config.get read (shared gate in components/config-read-state.tsx)
   and config.set writes against the two cloud chains and the realtime head
   (313-panel-ask §3):

     voice_cloud_asr_chain / voice_cloud_tts_chain  value = the chain, elements
       "provider/model", comma-separated. Each element is checked on the Go
       side: the provider must be in the registry, the model must be in that
       provider's catalog, and local-sherpa is refused outright (it is the
       local fallback and never enters a cloud chain).
     voice_realtime_provider / voice_realtime_model  the realtime head as two
       single-key writes (A877 ruling 1's beta shape - one key per call).

   THE TWO NAMED GAPS, printed on the page instead of papered over (§5):
     1. An EMPTY chain value is refused by the gate ("空值不写成默认", no clear
        sentinel exists yet) - so this page does not offer a clear button that
        could only ever fail, and says why next to the inputs.
     2. voice.realtime.enabled is NOT writable (the ruled roster has
        provider/model only) - the page can point the head but cannot switch
        it on, and says so where the realtime rows are.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useRef, useState } from "react";
import { RefreshCw } from "lucide-react";
import { Button } from "@/components/ai-native/button";
import { ConfigReadGate } from "@/components/config-read-state";
import { cn } from "@/lib/cn";
import type { ConfigView } from "@/lib/config-view";
import { requestConfigWrite } from "@/lib/panel";

export function VoiceSettings({
  heading = "语音",
  lede = "云回退链与实时头。链的每个元素写成「服务商/模型」，多个用逗号隔开；写入要重启进程才用得上，回执里会说。",
}: {
  /** The settings page passes its own words here - same convention as
      config-screen.tsx (see the note there). */
  heading?: string;
  lede?: string;
}) {
  return (
    <section aria-label="语音" className="flex min-w-0 flex-col gap-1">
      <h2 className="text-[15px] font-semibold text-ink">{heading}</h2>
      <p className="text-[12.5px] text-ink-2">{lede}</p>
      <ConfigReadGate>
        {(view, prose, reload) => <VoiceBody onChanged={reload} prose={prose} view={view} />}
      </ConfigReadGate>
    </section>
  );
}

interface Receipt {
  step: string;
  text: string;
  refused: boolean;
}

function chainText(chain: string[] | null): string {
  if (chain === null) return "";
  return chain.join(", ");
}

function VoiceBody({
  view,
  prose,
  onChanged,
}: {
  view: ConfigView;
  prose: string;
  onChanged: () => void;
}) {
  const [asr, setAsr] = useState(() => chainText(view.voice.cloudASRChain));
  const [tts, setTts] = useState(() => chainText(view.voice.cloudTTSChain));
  // The realtime inputs are seeded from the read once; they are field VALUES
  // from the config, not credentials, so keeping them in state is the same
  // thing the model page does with the provider name.
  const [rtProvider, setRtProvider] = useState(view.voice.realtimeProvider);
  const [rtModel, setRtModel] = useState(view.voice.realtimeModel);
  const [receipts, setReceipts] = useState<Receipt[]>([]);
  const [busy, setBusy] = useState<string | null>(null);
  const seq = useRef(0);

  async function run(step: string, send: () => Promise<string | null>) {
    const mine = ++seq.current;
    setBusy(step);
    const reply = await send();
    if (mine !== seq.current) return;
    setBusy(null);
    if (reply === null) {
      setReceipts((all) => [{ step, text: "没接到宿主，请求没有发出去。", refused: true }, ...all].slice(0, 6));
      return;
    }
    setReceipts((all) => [{ step, text: reply, refused: reply.includes("被拒绝") }, ...all].slice(0, 6));
    onChanged();
  }

  return (
    <div className="mt-2 flex min-w-0 flex-col gap-3">
      {/* 读到的现状。 */}
      <div className="flex flex-col gap-1.5 rounded-card border border-line bg-surface px-3.5 py-3">
        <div className="flex items-center justify-between gap-2">
          <p className="text-[11px] font-medium uppercase tracking-[0.06em] text-ink-3">配置读回</p>
          <Button aria-label="重读配置" onClick={onChanged} size="xs" variant="quiet">
            <RefreshCw aria-hidden="true" size={11} strokeWidth={2} />
            重读
          </Button>
        </div>
        <pre className="whitespace-pre-wrap text-[12px] leading-[1.7] text-ink-2">{prose}</pre>
        <div className="flex flex-col gap-1 pt-1">
          <ReadLine label="云 ASR 链" value={view.voice.cloudASRChain} />
          <ReadLine label="云 TTS 链" value={view.voice.cloudTTSChain} />
          <p className="text-[11.5px] text-ink-3">
            实时头：
            {view.voice.realtimeEnabled ? "已开" : "未开（enabled 不可写，见下）"}
            {view.voice.realtimeProvider !== "" ? ` · ${view.voice.realtimeProvider}` : " · provider 未设"}
            {view.voice.realtimeModel !== "" ? ` / ${view.voice.realtimeModel}` : " / model 未设"}
          </p>
        </div>
      </div>

      {!view.readable ? (
        <div className="rounded-card border border-line bg-surface px-3.5 py-3">
          <p className="text-[12px] leading-[1.7] text-ink-2">
            配置当前不可读（这不是未配置）：先把它修好，下面的输入框才谈得上有意义。
          </p>
        </div>
      ) : (
        <>
          {/* 两条云链。 */}
          <ChainRow
            busy={busy !== null}
            caption="ASR（听写）回退链"
            field="voice_cloud_asr_chain"
            onChange={setAsr}
            onWrite={() =>
              run("写 ASR 链", () =>
                requestConfigWrite({ field: "voice_cloud_asr_chain", value: asr.trim() }),
              )
            }
            placeholder="stepfun/模型id, 另一家/模型id"
            value={asr}
          />
          <ChainRow
            busy={busy !== null}
            caption="TTS（朗读）回退链"
            field="voice_cloud_tts_chain"
            onChange={setTts}
            onWrite={() =>
              run("写 TTS 链", () =>
                requestConfigWrite({ field: "voice_cloud_tts_chain", value: tts.trim() }),
              )
            }
            placeholder="stepfun/模型id"
            value={tts}
          />
          <p className="text-[11px] leading-[1.6] text-ink-3">
            元素规则：provider 必须在注册表里、model 必须在该服务商的目录里；local-sherpa 是本地兜底，永远不会进云链。
            空值不写成默认，所以这里没有「清空」按钮：清空链暂时没有界面口。
          </p>

          {/* 实时头：两个单键写入（β 形）。 */}
          <div className="flex flex-col gap-2 rounded-card border border-line bg-surface px-3.5 py-3">
            <p className="text-[12.5px] font-medium text-ink">实时头（陪聊路径）</p>
            <div className="flex flex-wrap items-center gap-2">
              <span className="text-[11.5px] text-ink-3">服务商</span>
              <input
                aria-label="实时头服务商"
                className="h-[27px] w-44 rounded-control bg-inset px-2 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3"
                onChange={(e) => setRtProvider(e.currentTarget.value)}
                placeholder="服务商名"
                value={rtProvider}
              />
              <Button
                disabled={busy !== null || rtProvider.trim() === ""}
                onClick={() =>
                  run("写实时服务商", () =>
                    requestConfigWrite({ field: "voice_realtime_provider", value: rtProvider.trim() }),
                  )
                }
                size="sm"
                variant="secondary"
              >
                写入 provider
              </Button>
            </div>
            <div className="flex flex-wrap items-center gap-2">
              <span className="text-[11.5px] text-ink-3">模型</span>
              <input
                aria-label="实时头模型"
                className="h-[27px] w-44 rounded-control bg-inset px-2 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3"
                onChange={(e) => setRtModel(e.currentTarget.value)}
                placeholder="模型 id"
                value={rtModel}
              />
              <Button
                disabled={busy !== null || rtModel.trim() === ""}
                onClick={() =>
                  run("写实时模型", () =>
                    requestConfigWrite({ field: "voice_realtime_model", value: rtModel.trim() }),
                  )
                }
                size="sm"
                variant="secondary"
              >
                写入 model
              </Button>
            </div>
            <p className="text-[11px] leading-[1.6] text-ink-3">
              两枚是两次单键写入（一字段一写入，β 形）。enabled 不可写：这里配得成对，但今天点不亮实时头。
            </p>
          </div>

          {/* 回执：每一句都是宿主发的原文。 */}
          {receipts.length > 0 ? (
            <div className="flex flex-col gap-2">
              {receipts.map((r, i) => (
                <div
                  className={cn(
                    "rounded-card border px-3.5 py-2.5",
                    r.refused ? "border-red/40 bg-red-tint" : "border-line bg-surface",
                  )}
                  key={`${r.step}-${i}`}
                >
                  <p className="text-[11px] font-medium text-ink-3">{r.step}</p>
                  <pre className="whitespace-pre-wrap text-[11.5px] leading-[1.7] text-ink-2">{r.text}</pre>
                </div>
              ))}
            </div>
          ) : null}

          <p className="text-[11px] leading-[1.6] text-ink-3">
            写入的生效级别由回执里那句话说，页面不另说。
          </p>
        </>
      )}
    </div>
  );
}

/** One line of the read-back: the chain as the config carries it, or the
    honest "not answered" when the array is null (which is a different fact
    from an empty chain). */
function ReadLine({ label, value }: { label: string; value: string[] | null }) {
  return (
    <p className="text-[11.5px] leading-[1.6] text-ink-2">
      {label}：
      {value === null ? (
        <span className="text-ink-3">没读到（腿还没答）</span>
      ) : value.length === 0 ? (
        <span className="text-ink-3">空（未配置）</span>
      ) : (
        <span className="font-mono">{value.join(", ")}</span>
      )}
    </p>
  );
}

function ChainRow({
  caption,
  field,
  value,
  placeholder,
  busy,
  onChange,
  onWrite,
}: {
  caption: string;
  field: string;
  value: string;
  placeholder: string;
  busy: boolean;
  onChange: (next: string) => void;
  onWrite: () => void;
}) {
  return (
    <div className="flex flex-col gap-2 rounded-card border border-line bg-surface px-3.5 py-3">
      <div className="flex flex-wrap items-center gap-2">
        <span className="text-[12.5px] font-medium text-ink">{caption}</span>
        <span className="font-mono text-[10.5px] text-ink-3">{field}</span>
      </div>
      <div className="flex flex-wrap items-center gap-2">
        <input
          aria-label={caption}
          className="h-[27px] min-w-[280px] flex-1 rounded-control bg-inset px-2 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3"
          onChange={(e) => onChange(e.currentTarget.value)}
          placeholder={placeholder}
          value={value}
        />
        <Button disabled={busy || value.trim() === ""} onClick={onWrite} size="sm" variant="primary">
          写入链
        </Button>
      </div>
    </div>
  );
}
