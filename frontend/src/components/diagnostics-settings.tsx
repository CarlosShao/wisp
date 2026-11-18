/* ============================================================================
   诊断 settings category - the diagnostics doors get their page (ticket 40
   AC#5, the readface diagnostics leg).
   ----------------------------------------------------------------------------
   WHAT IS REAL HERE. The Go half landed readface wave 5 (2026-11-18): the
   two roster doors (bridge.go's diagnostics.preview / diagnostics.export -
   the export spelling is SPEC-08 §5.2:173's own) are routed by
   composer_dispatch.go, wired at panel_inbound.go, and answered by
   cmd/wisp/panel_diagnostics_40.go over the ticket-45 engine - the SAME
   diagnosticsBundleOptions180P assembler and the SAME Preview/Build pair the
   `wisp diagnostics` CLI uses, under the SAME [privacy] diagnostics_opt_in
   total switch. This page sends every click through the real door
   (lib/panel.ts's two requests on the settings envelope) and shows the
   engine's own words verbatim. There is no local result anywhere: a page
   that derived its own entry list would be exactly the "preview writes one
   list, the zip carries another" shape the AC bans.

   THE PREVIEW IS THE REVIEW. The preview card renders the engine's own
   BundlePreview: every category with its policy detail verbatim, every
   included item with name + bytes + the redacted head (head passed through
   the same redactor the builder uses, so the page labels it 脱敏), every
   excluded item with the engine's own reason verbatim. The user reviews
   THIS list before clicking export (SPEC-14 §14.4's review-before-send).

   AN EMPTY PREVIEW IS NOT A NORMAL READING. A machine whose bundle carries
   zero entries renders the excluded categories' own statements out loud
   (each category's payload detail, verbatim - the engine owns those words)
   and an explicit 条目 0 个 line - never a bare empty list that reads as
   "all normal". The categories' words come from the payload; the page never
   writes the exclusion sentences itself.

   THE EXPORT BUTTON WALKS THE REAL EXPORT. The user has seen the preview
   card; export runs the SAME options with a real path, the preview FIRST (a
   failure aborts before anything is written), then the builder. The receipt
   names the written path verbatim; the zip stays on disk - the bundle is a
   file the user opens, not bytes over this bridge. A refusal (including the
   switch-off one that names diagnostics_opt_in) renders verbatim in the
   danger tint; the page never fakes a success.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { Button } from "@/components/ai-native/button";
import { SettingsCard } from "@/components/settings-controls";
import {
  useDiagnosticsExport,
  useDiagnosticsPreview,
  type DiagnosticsCategoryView,
  type DiagnosticsExportState,
  type DiagnosticsPreviewState,
} from "@/lib/diagnostics-view";

/** One category's rows: included items carry name + bytes + the redacted
    head under the 脱敏 label; excluded items carry the engine's own reason
    verbatim. No emoji, no fabricated rows. */
export function DiagnosticsCategoryBlock({ category }: { category: DiagnosticsCategoryView }) {
  return (
    <div className="rounded-card border border-line bg-inset px-3 py-2.5">
      <p className="text-[12px] font-medium text-ink">
        {category.category}
        <span className="ml-2 font-normal text-ink-3">
          {category.defaultIncluded ? "默认包含" : "默认排除"}
          {category.collected ? "" : "　未收集"}
        </span>
      </p>
      <p className="mt-1 text-[11px] leading-[1.6] text-ink-2">{category.detail}</p>
      {category.items.length > 0 ? (
        <ul className="mt-1.5 flex flex-col gap-1">
          {category.items.map((item) => (
            <li key={item.name} className="text-[11.5px] leading-[1.6] text-ink">
              <span className="font-mono">{item.name}</span>
              <span className="ml-2 text-ink-3">{item.bytes} 字节</span>
              {item.included ? null : <span className="ml-2 text-ink-3">{item.reason}</span>}
              {item.head !== "" ? (
                <span className="mt-0.5 block break-all font-mono text-[10.5px] leading-[1.5] text-ink-3">
                  <span className="mr-1.5 rounded-control bg-hover px-1 py-px text-[10px] text-ink-2">脱敏</span>
                  {item.head}
                </span>
              ) : null}
            </li>
          ))}
        </ul>
      ) : null}
    </div>
  );
}

/** The preview reply area: the engine's own words verbatim. This is the
    behavioural feed point the render instrument drives
    (scripts/render-diagnostics-40.tsx): fed a real preview it shows the
    entry list verbatim; fed a refusal it shows the refusal's own words
    (including the one naming diagnostics_opt_in); fed nothing parseable it
    says 未判 and never draws a fabricated entry. */
export function DiagnosticsPreviewBlock({ state }: { state: DiagnosticsPreviewState }) {
  if (state.phase === "loading") {
    return <p className="text-[12px] leading-[1.7] text-ink-3">正在读取诊断包预览...</p>;
  }
  if (state.phase === "no-host") {
    return (
      <div className="rounded-card border border-line bg-surface px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-ink-2">
          没接到本机宿主（这一页没有 wispDispatch 门）：请求没有发出去，预览未判——这里不画任何假清单。
        </p>
      </div>
    );
  }
  if (state.phase === "unreadable") {
    return (
      <div className="rounded-card border border-red/40 bg-red-tint px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-red">门拒绝了这个请求，原话逐字如下：</p>
        <p className="mt-1.5 whitespace-pre-wrap break-all font-mono text-[11px] leading-[1.6] text-ink">
          {state.raw}
        </p>
      </div>
    );
  }
  if (state.phase === "payload-error") {
    return (
      <div className="rounded-card border border-orange/40 bg-orange-tint px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-orange">{state.message}</p>
        <p className="mt-1.5 max-h-32 overflow-y-auto whitespace-pre-wrap break-all font-mono text-[11px] leading-[1.6] text-ink">
          {state.raw}
        </p>
      </div>
    );
  }
  const { envelope } = state;
  return (
    <div className="flex flex-col gap-2">
      <p className="whitespace-pre-wrap text-[12px] leading-[1.7] text-ink">{state.prose}</p>
      {envelope.entries.length === 0 ? (
        /* THE EMPTY-PREVIEW READING. Zero entries is said OUT LOUD with the
           payload's own exclusion statements carried verbatim - never a bare
           empty list that reads as a normal, nothing-to-see reading. */
        <div className="rounded-card border border-orange/40 bg-orange-tint px-3.5 py-3">
          <p className="text-[12.5px] font-semibold text-orange">条目 0 个：这个包当前不携带任何文件</p>
          <p className="mt-1 text-[11.5px] leading-[1.7] text-ink-2">
            不是正常的空读数：每一类都给出了引擎自己的未收集原因，逐类逐字如下——
          </p>
          <ul className="mt-1.5 flex flex-col gap-1">
            {envelope.categories.map((c) => (
              <li key={c.category} className="text-[11.5px] leading-[1.6] text-ink">
                <span className="font-mono">{c.category}</span>
                <span className="ml-2 text-ink-2">{c.detail}</span>
              </li>
            ))}
          </ul>
        </div>
      ) : (
        <div className="rounded-card border border-line bg-surface px-3.5 py-3">
          <p className="text-[12px] text-ink">
            zip 条目清单（{envelope.entries.length} 条，与导出同一个来源，逐字）：
          </p>
          <ul className="mt-1.5 flex flex-col gap-0.5">
            {envelope.entries.map((entry) => (
              <li key={entry} className="break-all font-mono text-[11px] leading-[1.5] text-ink">
                {entry}
              </li>
            ))}
          </ul>
          <p className="mt-1.5 text-[11px] leading-[1.6] text-ink-3">合计 {envelope.totalBytes} 字节</p>
        </div>
      )}
      {envelope.categories.map((c) => (
        <DiagnosticsCategoryBlock key={c.category} category={c} />
      ))}
      {state.skipped > 0 ? (
        <p className="text-[11px] leading-[1.6] text-orange">
          载荷上有 {state.skipped} 条残缺记录被丢弃（字段残缺被丢弃，不补零）。
        </p>
      ) : null}
    </div>
  );
}

/** The export reply area: the receipt names the written path verbatim; the
    payload is the same preview JSON whose out_path is now the real target. */
export function DiagnosticsExportBlock({ state }: { state: DiagnosticsExportState }) {
  if (state.phase === "idle") return null;
  if (state.phase === "sending") {
    return <p className="text-[12px] leading-[1.7] text-ink-3">导出已在桥上，等门回话...</p>;
  }
  if (state.phase === "no-host") {
    return (
      <div className="rounded-card border border-line bg-surface px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-ink-2">
          没接到本机宿主（这一页没有 wispDispatch 门）：请求没有发出去，导出未判——不画假成功。
        </p>
      </div>
    );
  }
  if (state.phase === "refused") {
    return (
      <div className="rounded-card border border-red/40 bg-red-tint px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-red">门拒绝了这个导出，原话逐字如下：</p>
        <p className="mt-1.5 whitespace-pre-wrap break-all font-mono text-[11px] leading-[1.6] text-ink">
          {state.raw}
        </p>
      </div>
    );
  }
  if (state.phase === "no-payload") {
    return (
      <div className="rounded-card border border-orange/40 bg-orange-tint px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-orange">
          回执上没有解析到导出载荷行（wisp-diagnostics-export: 行不在或不是 JSON）：写盘结果未判，不画假成功。原话逐字如下：
        </p>
        <p className="mt-1.5 whitespace-pre-wrap break-all font-mono text-[11px] leading-[1.6] text-ink">
          {state.raw}
        </p>
      </div>
    );
  }
  if (state.phase === "payload-error") {
    return (
      <div className="rounded-card border border-orange/40 bg-orange-tint px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-orange">{state.message}</p>
        <p className="mt-1.5 max-h-32 overflow-y-auto whitespace-pre-wrap break-all font-mono text-[11px] leading-[1.6] text-ink">
          {state.raw}
        </p>
      </div>
    );
  }
  return (
    <div className="rounded-card border border-green/40 bg-green-tint px-3.5 py-3">
      <p className="text-[12px] leading-[1.7] text-ink">回执逐字：</p>
      <p className="mt-1 whitespace-pre-wrap break-all font-mono text-[11.5px] leading-[1.6] text-ink">{state.prose}</p>
      <p className="mt-2 text-[12px] leading-[1.7] text-ink">
        写盘路径（载荷 out_path 逐字）：<span className="font-mono">{state.envelope.outPath}</span>
      </p>
      <p className="mt-1 text-[11.5px] leading-[1.6] text-ink-2">
        zip 条目 {state.envelope.entries.length} 条，与预览同一来源；音频、密钥、完整转写按引擎默认排除（各类说明见预览卡，逐字来自载荷）。
      </p>
    </div>
  );
}

const EXPORT_CAPTION =
  "导出走真导出通路（diagnostics.export）：与预览同一份文件清单来源（cmd/wisp 的 diagnosticsBundleOptions180P 组装、internal/observe 同一对 Preview/Build 函数），先跑预览、预览失败就中止且不写任何文件；zip 落在数据根目录，回执点名路径。总闸与 wisp diagnostics CLI 是同一把 [privacy] diagnostics_opt_in。";

export function DiagnosticsSettings({
  heading = "诊断",
  lede = "诊断包预览与导出：先看清单，再导出（票 45 引擎，票 40 AC#5）。",
}: {
  heading?: string;
  lede?: string;
}) {
  const preview = useDiagnosticsPreview();
  // The export receipt re-reads the preview, so the page never patches its
  // own picture of the bundle (the settings-door rule).
  const exp = useDiagnosticsExport(preview.reload);

  return (
    <section aria-label="诊断" className="flex min-w-0 flex-col gap-1">
      <h2 className="text-[15px] font-semibold text-ink">{heading}</h2>
      <p className="text-[12.5px] text-ink-2">{lede}</p>

      <SettingsCard
        title="预览清单（导出前先看）"
        description="diagnostics.preview：引擎自报的 zip 条目清单与逐类说明（含文件名+大小+脱敏后的内容首行；排除类给出引擎自己的原因，逐字）。这个清单与导出走同一个来源，导出前先审它。"
      >
        <div className="flex flex-wrap items-center gap-2">
          <Button onClick={preview.reload} size="sm" variant="secondary">
            重新读预览
          </Button>
        </div>
        <div className="mt-2">
          <DiagnosticsPreviewBlock state={preview.state} />
        </div>
      </SettingsCard>

      <SettingsCard title="导出诊断包" description="diagnostics.export：走真导出，回执点名写盘路径；总闸关着时门拒绝并点名 diagnostics_opt_in。">
        <div className="flex flex-wrap items-center gap-2">
          <Button onClick={exp.submit} size="sm" variant="primary">
            导出诊断包
          </Button>
        </div>
        <p className="mt-2 text-[11px] leading-[1.6] text-ink-3">{EXPORT_CAPTION}</p>
        <div className="mt-2">
          <DiagnosticsExportBlock state={exp.state} />
        </div>
      </SettingsCard>
    </section>
  );
}
