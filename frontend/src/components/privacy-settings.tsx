/* ============================================================================
   隐私 settings category - the three privacy doors get their page buttons
   (ticket 40 AC#1, the readface privacy leg).
   ----------------------------------------------------------------------------
   WHAT IS REAL HERE. The Go half landed readface wave 1 (2026-11-17): the
   three roster doors (bridge.go:141-143 privacy.delete-one / privacy.purge /
   privacy.export) are routed by composer_dispatch.go:174's table, wired at
   panel_inbound.go:411, and answered
   by cmd/wisp/panel_privacy_40.go over internal/memory's privacy layer - the
   layer with zero production callers until then (census A365). This page
   sends every click through the real door (lib/panel.ts's three requests on
   the settings envelope) and shows the receipt or refusal verbatim. There is
   no local result anywhere: a button that drew its own success line would be
   exactly the 造假 the dispatch bans.

   THE DOMAIN IS THE DOOR'S, NOT THE PAGE'S. The valid domain words are owned
   by internal/memory's PrivacyDomains() (the five D35 domains) and the page
   NEVER copies that vocabulary into a picker: the word the user types rides
   the envelope's privacyDomain key verbatim, the door validates it against
   PrivacyDomains(), and an unknown word comes back as the refusal sentence
   that NAMES the vocabulary - shown verbatim below the input. A page-side
   word list would be a second vocabulary that silently rots when the schema
   grows a sixth domain; the refusal is the ask, kept alive.

   THE ID IS TYPED, AND THE PAGE SAYS WHY. There is no privacy LIST route on
   the wire (bridge.go's roster carries the three write/read doors, no
   privacy.list), so the page cannot offer a row picker without inventing
   rows. The delete-one branch therefore asks for the id, and the caption
   names the missing read face instead of pretending a list exists.

   THE CONFIRM CARD. privacy.purge is irreversible (SPEC-08:170's own "purge
   需确认"), so the purge branch is a two-step control: the first button only
   opens a danger-styled confirmation card naming the domain and the
   irreversibility; the request is sent ONLY from the card's confirm button.
   delete-one is one item out of one domain and export writes nothing, so
   they stay one-click (deletion must be one click, destruction must ask
   first). There is no "purge everything" control on this page: the door
   itself has no such argument (panel_privacy_40.go's file comment), and the
   page invents none.

   THE GO-SIDE CONFIRM BIT (2026-11-17, the audit R1 fix). The door now also
   re-judges the confirmation itself: a destructive request without
   privacyConfirmed:true is refused by name and deletes ZERO rows. The bit is
   placed at the wire layer (lib/panel.ts's two destructive request functions
   - the only calls this card's 确认清空 key and the delete button's 发起删除
   click can reach) and nowhere else: the first card button sends nothing, the
   取消 exit sends nothing, and export never carries the bit. Confirm here,
   the bit rides; skip this card, the door refuses.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useState } from "react";
import { Button } from "@/components/ai-native/button";
import { SettingsCard } from "@/components/settings-controls";
import { cn } from "@/lib/cn";
import { privacyRunners, usePrivacyRequest, type PrivacyActionState } from "@/lib/privacy-view";

const INPUT_CLASS =
  "h-[30px] rounded-control bg-inset px-2.5 font-mono text-[12px] text-ink shadow-hairline outline-none placeholder:font-sans placeholder:text-ink-3";

/** The reply area: whatever the door answered, verbatim, in the branch's own
    slot. This is the behavioural feed point the render instrument drives
    (scripts/render-privacy-40.tsx): fed a receipt it shows the receipt's own
    words; fed a refusal it shows the refusal's own words (including the
    domain-vocabulary sentence); fed nothing parseable it says 未判 and never
    draws a fabricated count. */
export function PrivacyReplyBlock({ state }: { state: PrivacyActionState }) {
  if (state.phase === "idle") return null;
  if (state.phase === "sending") {
    return <p className="text-[12px] leading-[1.7] text-ink-3">请求已在桥上，等门回话...</p>;
  }
  if (state.phase === "no-host") {
    return (
      <div className="rounded-card border border-line bg-surface px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-ink-2">
          没接到本机宿主（这一页没有 wispDispatch 门）：请求没有发出去，回执未判——这里不画任何假结果。
        </p>
      </div>
    );
  }
  if (state.phase === "refused") {
    return (
      <div className="rounded-card border border-red/40 bg-red-tint px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-red">门拒绝了这个请求，原话逐字如下：</p>
        <p className="mt-1.5 whitespace-pre-wrap break-all font-mono text-[11px] leading-[1.6] text-ink">
          {state.raw}
        </p>
      </div>
    );
  }
  if (state.phase === "receipt") {
    return (
      <div className="rounded-card border border-green/40 bg-green-tint px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-ink">回执逐字：</p>
        <p className="mt-1 whitespace-pre-wrap font-mono text-[11.5px] leading-[1.6] text-ink">{state.text}</p>
      </div>
    );
  }
  if (state.phase === "no-payload") {
    return (
      <div className="rounded-card border border-orange/40 bg-orange-tint px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-orange">
          回执上没有解析到导出载荷行（wisp-privacy-export: 行不在或不是 JSON）：records 未判，不画 0。原话逐字如下：
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
  // phase === "export": the prose verbatim, the record count READ OFF THE
  // PAYLOAD (the same len(records) the receipt sentence states), and the
  // store's own JSON carried through unreshaped - the export is the user's
  // data, not a summary, so the rows stay on screen verbatim.
  return (
    <div className="rounded-card border border-green/40 bg-green-tint px-3.5 py-3">
      <p className="text-[12px] leading-[1.7] text-ink">回执逐字：</p>
      <p className="mt-1 whitespace-pre-wrap font-mono text-[11.5px] leading-[1.6] text-ink">{state.prose}</p>
      <p className="mt-2 text-[12px] leading-[1.7] text-ink">
        records={state.envelope.recordsCount}
        <span className="text-ink-3">（读自载荷 records 数组的长度，与回执句子同一读数）</span>
      </p>
      <p className="mt-1 text-[11px] leading-[1.6] text-ink-3">
        域（载荷自报）：{state.envelope.domain}
        {state.envelope.exportedAt === "" ? "" : `　导出时刻：${state.envelope.exportedAt}`}
      </p>
      <pre className="mt-2 max-h-48 overflow-auto whitespace-pre-wrap break-all rounded-control border border-line bg-inset px-3 py-2 font-mono text-[11px] leading-[1.5] text-ink">
        {JSON.stringify(state.envelope.records)}
      </pre>
    </div>
  );
}

/** One labelled text field shared by the three branches. The value is the
    parent's (a submitted word must stay visible while its receipt shows). */
function PrivacyField({
  label,
  value,
  onChange,
  placeholder,
}: {
  label: string;
  value: string;
  onChange: (next: string) => void;
  placeholder: string;
}) {
  return (
    <label className="flex min-w-0 flex-1 flex-col gap-1">
      <span className="text-[11px] text-ink-3">{label}</span>
      <input
        className={INPUT_CLASS}
        onChange={(e) => onChange(e.currentTarget.value)}
        placeholder={placeholder}
        value={value}
      />
    </label>
  );
}

/** The caption under the domain inputs - the ask-don't-copy rule, stated
    where the user can read it. */
const DOMAIN_CAPTION =
  "域词由 internal/memory 的 PrivacyDomains() 拥有（D35 五域），本页不复制词表：你填的词原样上桥，门拒绝时拒绝句会点名可用的域，逐字显示在下面。";

export function PrivacySettings({
  heading = "隐私",
  lede = "本机数据的删除、清空与导出（D35 五域，走真桥）。",
}: {
  heading?: string;
  lede?: string;
}) {
  // One state slot per branch, so no receipt can read under the wrong button.
  const [deleteDomain, setDeleteDomain] = useState("");
  const [deleteItemId, setDeleteItemId] = useState("");
  const [purgeDomain, setPurgeDomain] = useState("");
  const [purgeConfirmOpen, setPurgeConfirmOpen] = useState(false);
  const [exportDomain, setExportDomain] = useState("");

  const del = usePrivacyRequest("plain", privacyRunners.deleteOne);
  const purge = usePrivacyRequest("plain", privacyRunners.purge);
  const exp = usePrivacyRequest("export", privacyRunners.export);

  return (
    <section aria-label="隐私" className="flex min-w-0 flex-col gap-1">
      <h2 className="text-[15px] font-semibold text-ink">{heading}</h2>
      <p className="text-[12.5px] text-ink-2">{lede}</p>

      <SettingsCard
        title="删除一条"
        description="privacy.delete-one：一域一条，回执点名 id。面板没有隐私条目的清单读面（名册只有这三扇门），id 需要你手头有；不填 id 这一格不发起。这一次点击就是确认：请求带 privacyConfirmed:true，门上没有确认位的同源伪造请求会被具名拒绝、零删除。"
      >
        <div className="flex flex-wrap items-end gap-3">
          <PrivacyField label="域（privacyDomain）" onChange={setDeleteDomain} placeholder="域词" value={deleteDomain} />
          <PrivacyField label="对象 id（privacyItemId）" onChange={setDeleteItemId} placeholder="记录 id" value={deleteItemId} />
          <Button
            aria-label="删除一条隐私记录"
            disabled={deleteItemId.trim() === ""}
            onClick={() => del.submit(deleteDomain, deleteItemId)}
            size="sm"
            variant="primary"
          >
            发起删除
          </Button>
        </div>
        <p className="mt-2 text-[11px] leading-[1.6] text-ink-3">{DOMAIN_CAPTION}</p>
        <div className="mt-2">
          <PrivacyReplyBlock state={del.state} />
        </div>
      </SettingsCard>

      <SettingsCard
        title="清空某域"
        description="privacy.purge：整域清空，不可逆（SPEC-08:170 自己写着 purge 需确认）。这里没有一键清空：第一下只开确认卡，请求只从卡上的确认键发出；门上也没有「全清」参数，域词必填。"
      >
        <div className="flex flex-wrap items-end gap-3">
          <PrivacyField label="域（privacyDomain）" onChange={setPurgeDomain} placeholder="域词" value={purgeDomain} />
          <Button
            aria-label="清空域（先确认）"
            disabled={purgeDomain.trim() === ""}
            onClick={() => setPurgeConfirmOpen(true)}
            size="sm"
            variant="secondary"
          >
            清空这一域
          </Button>
        </div>
        {purgeConfirmOpen ? (
          <div className="mt-3 rounded-card border border-red/40 bg-red-tint px-3.5 py-3" role="alertdialog" aria-label="确认清空">
            <p className="text-[12.5px] font-semibold text-red">确认清空</p>
            <p className="mt-1 text-[12px] leading-[1.7] text-ink">
              privacy.purge 不可逆：整域清空后无法恢复。即将清空的域：{purgeDomain.trim()}
              （词原样上桥，门按 PrivacyDomains() 判）。其他域不会动；清空的条数以门回执自报的计数为准。
            </p>
            <div className="mt-2.5 flex items-center gap-2">
              <Button onClick={() => setPurgeConfirmOpen(false)} size="xs" variant="secondary">
                取消
              </Button>
              <Button
                aria-label="确认清空"
                className={cn("bg-red text-white hover:bg-red/90")}
                onClick={() => {
                  setPurgeConfirmOpen(false);
                  purge.submit(purgeDomain);
                }}
                size="xs"
                variant="primary"
              >
                确认清空
              </Button>
            </div>
          </div>
        ) : null}
        <p className="mt-2 text-[11px] leading-[1.6] text-ink-3">{DOMAIN_CAPTION}</p>
        <div className="mt-2">
          <PrivacyReplyBlock state={purge.state} />
        </div>
      </SettingsCard>

      <SettingsCard
        title="导出某域"
        description="privacy.export：整域导出为存储自报的 JSON（导出即你自己的数据，不是摘要）；回执 prose 逐字上屏，records=N 读自载荷 records 数组，载荷行不在就明说未判，不画 0。"
      >
        <div className="flex flex-wrap items-end gap-3">
          <PrivacyField label="域（privacyDomain）" onChange={setExportDomain} placeholder="域词" value={exportDomain} />
          <Button
            aria-label="导出域 JSON"
            onClick={() => exp.submit(exportDomain)}
            size="sm"
            variant="secondary"
          >
            导出 JSON
          </Button>
        </div>
        <p className="mt-2 text-[11px] leading-[1.6] text-ink-3">{DOMAIN_CAPTION}</p>
        <div className="mt-2">
          <PrivacyReplyBlock state={exp.state} />
        </div>
      </SettingsCard>
    </section>
  );
}
