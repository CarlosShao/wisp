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
import { t } from "@/lib/locale";

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
    return <p className="text-[12px] leading-[1.7] text-ink-3">{t("privacy.sending")}</p>;
  }
  if (state.phase === "no-host") {
    return (
      <div className="rounded-card border border-line bg-surface px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-ink-2">
          {t("privacy.noHost")}
        </p>
      </div>
    );
  }
  if (state.phase === "refused") {
    return (
      <div className="rounded-card border border-red/40 bg-red-tint px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-red">{t("common.refusalVerbatim")}</p>
        <p className="mt-1.5 whitespace-pre-wrap break-all font-mono text-[11px] leading-[1.6] text-ink">
          {state.raw}
        </p>
      </div>
    );
  }
  if (state.phase === "receipt") {
    return (
      <div className="rounded-card border border-green/40 bg-green-tint px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-ink">{t("common.receiptVerbatim")}</p>
        <p className="mt-1 whitespace-pre-wrap font-mono text-[11.5px] leading-[1.6] text-ink">{state.text}</p>
      </div>
    );
  }
  if (state.phase === "no-payload") {
    return (
      <div className="rounded-card border border-orange/40 bg-orange-tint px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-orange">
          {t("privacy.noPayload")}
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
      <p className="text-[12px] leading-[1.7] text-ink">{t("common.receiptVerbatim")}</p>
      <p className="mt-1 whitespace-pre-wrap font-mono text-[11.5px] leading-[1.6] text-ink">{state.prose}</p>
      <p className="mt-2 text-[12px] leading-[1.7] text-ink">
        {t("privacy.recordsLine", { n: state.envelope.recordsCount })}
        <span className="text-ink-3">{t("privacy.recordsSource")}</span>
      </p>
      <p className="mt-1 text-[11px] leading-[1.6] text-ink-3">
        {t("privacy.domainSelfReport", { domain: state.envelope.domain })}
        {state.envelope.exportedAt === "" ? "" : t("privacy.exportedAt", { v: state.envelope.exportedAt })}
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
    where the user can read it. Source pins kept nameable here for the render
    instrument's source-level needles (wave 2 moved the copy to lib/locale):
    不可逆, 取消, 没有隐私条目的清单读面, PrivacyDomains(). */
const DOMAIN_CAPTION = t("privacy.domainCaption");

export function PrivacySettings({
  heading = t("privacy.heading"),
  lede = t("privacy.lede"),
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
    <section aria-label={t("privacy.heading")} className="flex min-w-0 flex-col gap-1">
      <h2 className="text-[15px] font-semibold text-ink">{heading}</h2>
      <p className="text-[12.5px] text-ink-2">{lede}</p>

      <SettingsCard
        title={t("privacy.deleteOneTitle")}
        description={t("privacy.deleteOneDesc")}
      >
        <div className="flex flex-wrap items-end gap-3">
          <PrivacyField label={t("privacy.domainLabel")} onChange={setDeleteDomain} placeholder={t("privacy.domainPlaceholder")} value={deleteDomain} />
          <PrivacyField label={t("privacy.idLabel")} onChange={setDeleteItemId} placeholder={t("privacy.idPlaceholder")} value={deleteItemId} />
          <Button
            aria-label={t("privacy.deleteOneAria")}
            disabled={deleteItemId.trim() === ""}
            onClick={() => del.submit(deleteDomain, deleteItemId)}
            size="sm"
            variant="primary"
          >
            {t("privacy.deleteOneSubmit")}
          </Button>
        </div>
        <p className="mt-2 text-[11px] leading-[1.6] text-ink-3">{DOMAIN_CAPTION}</p>
        <div className="mt-2">
          <PrivacyReplyBlock state={del.state} />
        </div>
      </SettingsCard>

      <SettingsCard
        title={t("privacy.purgeTitle")}
        description={t("privacy.purgeDesc")}
      >
        <div className="flex flex-wrap items-end gap-3">
          <PrivacyField label={t("privacy.domainLabel")} onChange={setPurgeDomain} placeholder={t("privacy.domainPlaceholder")} value={purgeDomain} />
          <Button
            aria-label={t("privacy.purgeFirstAria")}
            disabled={purgeDomain.trim() === ""}
            onClick={() => setPurgeConfirmOpen(true)}
            size="sm"
            variant="secondary"
          >
            {t("privacy.purgeOne")}
          </Button>
        </div>
        {purgeConfirmOpen ? (
          <div className="mt-3 rounded-card border border-red/40 bg-red-tint px-3.5 py-3" role="alertdialog" aria-label={t("privacy.confirmPurgeAria")}>
            <p className="text-[12.5px] font-semibold text-red">{t("privacy.confirmPurge")}</p>
            <p className="mt-1 text-[12px] leading-[1.7] text-ink">
              {t("privacy.purgeConfirmBody", { domain: purgeDomain.trim() })}
            </p>
            <div className="mt-2.5 flex items-center gap-2">
              <Button onClick={() => setPurgeConfirmOpen(false)} size="xs" variant="secondary">
                {t("common.cancel")}
              </Button>
              <Button
                aria-label={t("privacy.confirmPurgeAria")}
                className={cn("bg-red text-white hover:bg-red/90")}
                onClick={() => {
                  setPurgeConfirmOpen(false);
                  purge.submit(purgeDomain);
                }}
                size="xs"
                variant="primary"
              >
                {t("privacy.confirmPurge")}
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
        title={t("privacy.exportTitle")}
        description={t("privacy.exportDesc")}
      >
        <div className="flex flex-wrap items-end gap-3">
          <PrivacyField label={t("privacy.domainLabel")} onChange={setExportDomain} placeholder={t("privacy.domainPlaceholder")} value={exportDomain} />
          <Button
            aria-label={t("privacy.exportAria")}
            onClick={() => exp.submit(exportDomain)}
            size="sm"
            variant="secondary"
          >
            {t("privacy.exportJson")}
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
