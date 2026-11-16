/* ============================================================================
   PluginsScreen - ticket 215a's plus-menu third category: the plugin
   inventory as a page.
   ----------------------------------------------------------------------------
   WHAT IS REAL HERE. Everything comes from the snapshot's plugins section
   (internal/panel/plugins_inventory_215.go): rows are the config's
   [plugins.<id>] entries, name/version come from the installed manifest
   (empty = no install on disk - measured absence, never invented), and each
   row's 生效与否 label is the wire's own three-state spelling
   (本轮生效/需重启/未生效) rendered verbatim - this file maps the known
   words to tones and shows an unknown state as its own text, never as one
   of the three.

   THE HONEST-STATE SENTENCES travel verbatim (engineState: the engine is
   not wired at the assembly root, so enabled=true is configured, not
   running; skillState: the skills dimension has no config keys - 215b,
   awaiting the owner's word on the config tree). A page that paraphrased
   them could outlive the facts they describe.

   THE SWITCH rides config.set's plugin_enabled field with the row's plugin
   id (lib/config-actions' performPluginToggleWrite) - the settings door,
   not a new route. Nothing is optimistic: the receipt goes on the feedback
   line, and the row's state is the NEXT snapshot's. The receipt itself
   says what the native side decided (a panel-sourced enable is a locked
   loosening and is refused there) - this page renders it verbatim and does
   not soften it.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useEffect, useState } from "react";
import { Puzzle, X } from "lucide-react";
import { FeedbackLine } from "@/components/settings-controls";
import { performPluginToggleWrite } from "@/lib/config-actions";
import type { Feedback } from "@/lib/config-feedback";
import type { PluginRowView, PluginsSection } from "@/lib/panel";

/** The three wire states' tones. The WORDS stay the wire's verbatim (the
    state text rendered on the row is row.state itself); only the tone is
    mapped here, and an unknown word falls to neutral with its own text. */
function stateTone(state: string): string {
  if (state === "本轮生效") return "bg-accent-tint text-accent-ink";
  if (state === "需重启") return "bg-warn-soft text-warn";
  return "bg-inset text-ink-3";
}

export function PluginsScreen({
  section,
  onClose,
}: {
  /** The snapshot's plugins section; undefined (key absent - a pump without
      the reader) draws the waiting sentence, never an empty list that would
      read as "no plugins configured". */
  section?: PluginsSection;
  onClose: () => void;
}) {
  const [feedback, setFeedback] = useState<Feedback | null>(null);
  const [busyId, setBusyId] = useState<string | null>(null);

  // Escape closes, the same way the command palette closes (App.tsx's own
  // handler stops at the document level; this one covers the overlay).
  useEffect(() => {
    function onKey(event: KeyboardEvent) {
      if (event.key === "Escape") onClose();
    }
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [onClose]);

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

  if (section === undefined) {
    return (
      <div className="flex flex-col px-4 py-4">
        <ScreenHeader onClose={onClose} />
        <p className="mt-3 text-[12px] leading-[1.7] text-ink-3">
          还没收到原生侧的插件清单（快照里没有 plugins 键）。这不是「没有插件」：
          清单键只在读口装配时上线上，缺键是读口不在，不是清单为空。
        </p>
      </div>
    );
  }

  return (
    <div className="flex max-h-[70vh] flex-col px-4 py-4">
      <ScreenHeader tier2Enabled={section.tier2Enabled} onClose={onClose} />
      <p className="mt-2 text-[11px] leading-[1.7] text-ink-3" role="note">
        {section.engineState}
      </p>
      <p className="mt-1 text-[11px] leading-[1.7] text-ink-3" role="note">
        {section.skillState}
      </p>

      {section.rows.length > 0 ? (
        <ul className="mt-3 flex min-h-0 flex-1 list-none flex-col gap-1.5 overflow-y-auto p-0">
          {section.rows.map((row) => (
            <li
              className="flex flex-col gap-1 rounded-card border border-line bg-surface px-3 py-2.5"
              key={row.id}
            >
              <div className="flex min-w-0 items-center gap-2">
                <Puzzle aria-hidden="true" className="size-3.5 shrink-0 text-ink-3" />
                <span className="shrink-0 text-[12.5px] font-medium text-ink">
                  {row.name !== "" ? row.name : row.id}
                </span>
                {row.version !== "" ? (
                  <span className="shrink-0 font-mono text-[10.5px] text-ink-3">v{row.version}</span>
                ) : null}
                <span className="min-w-0 flex-1 truncate font-mono text-[10.5px] text-ink-3">
                  {row.name !== "" ? row.id : ""}
                </span>
                <span
                  className={`shrink-0 rounded-full px-2 py-0.5 text-[10.5px] font-medium ${stateTone(row.state)}`}
                  title="生效与否由装配探针判，不是配置里的开关字面量"
                >
                  {row.state}
                </span>
                <button
                  aria-label={row.enabled ? `关闭插件 ${row.id}` : `开启插件 ${row.id}`}
                  className="shrink-0 rounded-chip bg-inset px-2 py-0.5 text-[11px] font-medium text-ink-2
                    transition-colors duration-100 hover:bg-hover hover:text-ink
                    disabled:pointer-events-none disabled:opacity-50"
                  disabled={busyId !== null}
                  onClick={() => void toggle(row)}
                  type="button"
                >
                  {row.enabled ? "关闭" : "开启"}
                </button>
              </div>
              {row.note !== undefined && row.note !== "" ? (
                <p className="text-[10.5px] leading-[1.6] text-warn" role="note">
                  {row.note}
                </p>
              ) : null}
              <p className="truncate text-[10.5px] leading-[1.6] text-ink-3">
                能力 {row.capabilities.length > 0 ? row.capabilities.join(" / ") : "（无）"}
                {row.netAllowlist.length > 0 ? ` · 网络白名单 ${row.netAllowlist.join(" / ")}` : ""}
                {row.hostApi.length > 0 ? ` · 宿主面 ${row.hostApi.join(" / ")}` : ""}
              </p>
            </li>
          ))}
        </ul>
      ) : (
        <p className="mt-3 text-[12px] leading-[1.7] text-ink-3">
          配置里没有声明任何 [plugins.&lt;id&gt;] 条目——这是读到的清单事实，不是读不到。
        </p>
      )}

      {feedback !== null ? (
        <div className="mt-3 shrink-0">
          <FeedbackLine feedback={feedback} />
        </div>
      ) : null}
    </div>
  );
}

function ScreenHeader({
  tier2Enabled,
  onClose,
}: {
  /** undefined = the section never arrived; the chip is not drawn, because a
      "关" here would be an invented fact about a section nobody read. */
  tier2Enabled?: boolean;
  onClose: () => void;
}) {
  return (
    <div className="flex shrink-0 items-center gap-2">
      <span className="text-[13px] font-medium text-ink">插件清单</span>
      {tier2Enabled !== undefined ? (
        <span className="shrink-0 rounded-full bg-inset px-2 py-0.5 text-[10.5px] font-medium text-ink-3">
          tier2_enabled {tier2Enabled ? "开" : "关"}
        </span>
      ) : null}
      <span className="min-w-0 flex-1" />
      <button
        aria-label="关闭插件清单"
        className="flex size-6 items-center justify-center rounded-chip text-ink-3 transition-colors duration-100 hover:bg-hover hover:text-ink"
        onClick={onClose}
        type="button"
      >
        <X aria-hidden="true" size={13} strokeWidth={2.4} />
      </button>
    </div>
  );
}
