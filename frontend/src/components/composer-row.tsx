/* ============================================================================
   ComposerRow (ticket 187's UI half; the bottom block of the conversation
   column, extracted from App.tsx so the switch feedback state can live here -
   render-nav pins App.tsx to exactly one useState, and the state object there
   is view preference, not a write receipt).
   ----------------------------------------------------------------------------
   WHAT IS REAL HERE. Everything this row offers comes from the snapshot's
   modelSwitch section (internal/panel/model_switch_187.go, read by
   internal/config's ChatSwitchOptions at the assembly root):

     - models[]      the offerable catalog rows (enabled entries only) - this
                     file invents no model id and no provider name (AC#4's
                     宁缺毋造 on the page side);
     - thinkingLevels[] the tiers the CURRENT model declares and its adapter
                     maps - possibly fewer than four, possibly none; an empty
                     list renders no tier control at all, never a guessed
                     four-value menu;
     - reason        the host's own "no rows can be offered" sentence, shown
                     verbatim when the list is empty;
     - current / currentThinking what the live config holds RIGHT NOW - the
                     button labels read these, and nothing here is optimistic:
                     after a write the next snapshot's section is the truth
                     (panel-ask §3), so a restart-pending model switch keeps
                     showing the old pair until Go says otherwise.

   WRITES. Both ride the settings door already in the C17 roster, via
   lib/config-actions: role_chat_pair for a model pick, role_chat_thinking for
   a tier pick. The receipt is translated with translateTieredWriteReply so
   the two halves keep their different effective tiers (model = restart,
   thinking = next task) instead of one generic "已保存".

   ABSENT KEY. A pump assembled without the ModelSwitch reader sends no
   modelSwitch key at all; this row then draws no picker - the same shape
   today's row has, which is honest absence, not a hidden failure.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useState } from "react";
import { PromptBar, type PromptEffort, type PromptModel, type PromptPermission } from "@/components/ai-native/prompt-bar";
import { ContextMeter } from "@/components/context-meter";
import { FeedbackLine } from "@/components/settings-controls";
import {
  performComposerModelWrite,
  performComposerThinkingWrite,
} from "@/lib/config-actions";
import type { Feedback } from "@/lib/config-feedback";
import type { ModelSwitchSection } from "@/lib/panel";

/** Thinking-tier enum words -> the Chinese the effort button shows. The write
    value stays the enum word verbatim (PromptEffort.key carries it); a word
    this table does not know shows its own text instead of being dropped. */
const TIER_LABELS: Record<string, string> = {
  off: "关",
  low: "低",
  medium: "中",
  high: "高",
};

export function ComposerRow({
  modelSwitch,
  permission,
  permissions,
  onPermissionChange,
  onSend,
}: {
  /** The snapshot's modelSwitch section; undefined (key absent) draws no
      picker - a read the pump never assembled is not a picker with nothing
      in it. */
  modelSwitch?: ModelSwitchSection;
  permission: string;
  permissions: readonly PromptPermission[];
  onPermissionChange?: (key: string) => void;
  onSend?: (text: string) => void;
}) {
  const [feedback, setFeedback] = useState<Feedback | null>(null);

  // The offerable rows, verbatim from the snapshot: the key is the pair the
  // role_chat_pair write addresses, the name is the model id, the tag names
  // the provider. No display name is invented here for rows the catalog does
  // not describe.
  const models: PromptModel[] = (modelSwitch?.models ?? []).map((m) => ({
    key: `${m.provider}/${m.model}`,
    name: m.model,
    tag: m.provider,
  }));
  const efforts: PromptEffort[] = (modelSwitch?.thinkingLevels ?? []).map((level) => ({
    key: level,
    name: TIER_LABELS[level] ?? level,
  }));

  // The pair in force RIGHT NOW. An unset role stays "" so the button says
  // 未设置 instead of standing in the first row for it (PromptBar's own
  // honesty rule for a prop that matches nothing shows the raw text, which
  // is how a restart-pending switch keeps telling the truth).
  const currentPair =
    modelSwitch !== undefined &&
    (modelSwitch.current.provider !== "" || modelSwitch.current.model !== "")
      ? `${modelSwitch.current.provider}/${modelSwitch.current.model}`
      : "";
  // "" on the wire means the file carries no tier key; the schema default
  // spells "off" on disk (model_switch_187.go's own comment), so "off" is the
  // value in force, not an invention.
  const currentThinking =
    modelSwitch !== undefined && modelSwitch.currentThinking !== ""
      ? modelSwitch.currentThinking
      : "off";

  async function switchModel(pairKey: string) {
    setFeedback({ tone: "neutral", text: "保存中…" });
    // First slash splits provider/model, the same cut the voice chain editor
    // makes - a model id with a slash keeps its tail intact.
    const cut = pairKey.indexOf("/");
    setFeedback(
      await performComposerModelWrite(
        cut < 0 ? pairKey : pairKey.slice(0, cut),
        cut < 0 ? "" : pairKey.slice(cut + 1),
      ),
    );
  }

  async function switchThinking(level: string) {
    setFeedback({ tone: "neutral", text: "保存中…" });
    setFeedback(await performComposerThinkingWrite(level));
  }

  return (
    <>
      {modelSwitch !== undefined && models.length === 0 && modelSwitch.reason !== "" ? (
        <p className="mb-1.5 text-[11px] leading-[1.6] text-ink-3" role="note">
          {modelSwitch.reason}（模型选择器不画）
        </p>
      ) : null}
      <ContextMeter usage={{}} className="mb-1.5" />
      <PromptBar
        effort={modelSwitch !== undefined ? currentThinking : undefined}
        efforts={efforts}
        model={modelSwitch !== undefined ? currentPair : undefined}
        models={models}
        onEffortChange={(key) => void switchThinking(key)}
        onModelChange={(key) => void switchModel(key)}
        onPermissionChange={onPermissionChange}
        onSend={onSend}
        permission={permission}
        permissions={permissions}
      />
      {feedback !== null ? (
        <div className="mt-1">
          <FeedbackLine feedback={feedback} />
        </div>
      ) : null}
    </>
  );
}
