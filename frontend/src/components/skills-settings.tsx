/* ============================================================================
   技能 settings category - the plugins inventory, mounted as a settings
   category (ticket 40, the census re-cut).
   ----------------------------------------------------------------------------
   THE CENSUS FACT THIS PAGE ANSWERS. The old waiting sentence said "还没收到
   技能清单", but the plugins inventory has had a source all along
   (internal/panel/plugins_inventory_215.go:176-183; the reader is attached at
   cmd/wisp/run.go's PumpSources.Plugins) - the census called that a 半假:
   the PLUGIN half had a feed nobody rendered here, while the SKILL half
   genuinely has no config keys to enumerate (SkillsNotInConfigModel215, the
   section's own skillState sentence).

   So this page mounts the inventory's body - the SAME PluginInventoryList the
   plus-menu overlay mounts, not a second copy - and the toggle rides
   config.set's plugin_enabled field through the settings door, exactly like
   the overlay's switch does. The receipt goes on the feedback line; nothing
   is optimistic, the row's state is the NEXT snapshot's.

   The two honest sentences travel verbatim (engineState: enabled=true is
   configured, not running; skillState: the skills dimension has no config
   keys). A page that paraphrased them could outlive the facts they describe.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useState } from "react";
import { PluginInventoryList } from "@/components/plugins-screen";
import { FeedbackLine } from "@/components/settings-controls";
import { performPluginToggleWrite } from "@/lib/config-actions";
import type { Feedback } from "@/lib/config-feedback";
import type { PluginRowView, PluginsSection } from "@/lib/panel";

export function SkillsSettings({
  section,
  heading = "技能",
  lede = "装了哪些技能、每个技能的权限范围。",
}: {
  /** The snapshot's plugins section; null (a pump without the reader - the
      key was ABSENT on the wire) draws the honest waiting sentence, never an
      empty list that would read as "nothing configured". */
  section: PluginsSection | null;
  heading?: string;
  lede?: string;
}) {
  const [feedback, setFeedback] = useState<Feedback | null>(null);
  const [busyId, setBusyId] = useState<string | null>(null);

  async function toggle(row: PluginRowView) {
    if (busyId !== null) return;
    setBusyId(row.id);
    setFeedback({ tone: "neutral", text: "写入中…" });
    try {
      setFeedback(await performPluginToggleWrite(row.id, !row.enabled));
    } finally {
      setBusyId(null);
    }
  }

  return (
    <section aria-label="技能" className="flex min-w-0 flex-col gap-1">
      <h2 className="text-[15px] font-semibold text-ink">{heading}</h2>
      <p className="text-[12.5px] text-ink-2">{lede}</p>

      {section === null ? (
        <div className="mt-2 rounded-card border border-line bg-surface px-3.5 py-3">
          <p className="text-[12px] leading-[1.7] text-ink-2">
            快照里没有 plugins 键：这一泵没装配插件清单读取器（PumpSources.Plugins，internal/panel/pump.go）。
            缺键是读口不在，不是清单为空——这一格等快照，不画猜出来的清单。
          </p>
        </div>
      ) : (
        <div className="mt-2 flex min-w-0 flex-col">
          <PluginInventoryList busyId={busyId} onToggle={(row) => void toggle(row)} section={section} />
        </div>
      )}

      {feedback !== null ? (
        <div className="mt-2">
          <FeedbackLine feedback={feedback} />
        </div>
      ) : null}
    </section>
  );
}
