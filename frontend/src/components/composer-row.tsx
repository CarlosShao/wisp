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

   THE FOUR COMPOSER CALLBACKS (wave128's finding: App.tsx was passing
   () => undefined stubs into the PromptBar, so every control LOOKED alive and
   none of them did anything). All four ride routes that already exist on the
   one C17 envelope in lib/panel - none is invented here, none is a fifth
   route (Q-50 = 甲):

     - 档位    onPermissionChange -> requestModeSwitch (panel.mode.request).
               Request only: the L2 strong confirmation and the audit line are
               the native side's, and the receipt is the NEXT snapshot's mode -
               nothing here is optimistic.
     - 工作区  switchWorkspace -> requestWorkspaceChange
               (panel.workspace.request). C26 resolves and answers through the
               snapshot's workspace section; the row reads it verbatim.
     - 附件    attach -> submitAttachment (panel.attachment.add), one request
               per picked file, bytes included. The stored/refused verdicts
               come back as snapshot rows; refusals render with their reason,
               and the chips (PromptBar) carry only the stored ones - a
               refused file drawn as a chip would be a lie about the user's
               pick. There is no attachment-removal route, so
               onRemoveAttachment is NOT passed: PromptBar hides its remove
               control when the handler is absent (faking removal would eat
               the user's file).
     - 发送    send -> sendMessage (panel.message.send) with the stored
               attachments; results arrive as pushed chunks in the column
               above, which is the receipt. A failure lands on the feedback
               line, never swallowed.
   ============================================================================ */

import { useRef, useState } from "react";
import { ChevronRight, FolderCog } from "lucide-react";
import { PromptBar, type PromptEffort, type PromptModel, type PromptPermission } from "@/components/ai-native/prompt-bar";
import { ContextMeter } from "@/components/context-meter";
import { FeedbackLine } from "@/components/settings-controls";
import {
  performComposerModelWrite,
  performComposerThinkingWrite,
} from "@/lib/config-actions";
import type { Feedback } from "@/lib/config-feedback";
import {
  requestModeSwitch,
  requestWorkspaceChange,
  sendMessage,
  submitAttachment,
  type ComposerState,
  type ModelSwitchSection,
} from "@/lib/panel";

/** The one error sentence shape for a request the envelope refused to carry
    (no host attached is the usual case - lib/panel's sendRequest throws).
    Told to the user on the feedback line, never swallowed. */
function requestError(err: unknown): string {
  return String(err instanceof Error ? err.message : err);
}

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
  composer,
  permission,
  permissions,
}: {
  /** The snapshot's modelSwitch section; undefined (key absent) draws no
      picker - a read the pump never assembled is not a picker with nothing
      in it. */
  modelSwitch?: ModelSwitchSection;
  /** The snapshot's composer section - the write targets the four callbacks
      need (the attachments a send carries, the workspace in force, the
      accepted mime list). Undefined (a pump without the composer reader)
      draws no workspace row and no attachment strip: honest absence. */
  composer?: ComposerState;
  permission: string;
  permissions: readonly PromptPermission[];
}) {
  const [feedback, setFeedback] = useState<Feedback | null>(null);
  const [workspaceDraft, setWorkspaceDraft] = useState("");
  const [attaching, setAttaching] = useState(false);
  const fileInput = useRef<HTMLInputElement>(null);

  // The native verdicts, split the way the retired composer split them: stored
  // rows become chips in the bar, refused rows render below with their full
  // reason (a truncated refusal reads as if the file had been accepted).
  const stored = (composer?.attachments ?? []).filter((a) => a.stored);
  const refused = (composer?.attachments ?? []).filter((a) => !a.stored);

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

  /* ── The four composer callbacks, each one request on the C17 envelope ── */

  // 档位: a REQUEST, never a record. The native side raises the L2 card when
  // `to` is the档 that stops asking, and the outcome arrives as the next
  // snapshot's mode; clicking the tier already in force asks for nothing.
  function askMode(to: string) {
    if (composer !== undefined && to === composer.mode.current) return;
    try {
      requestModeSwitch(to);
    } catch (err) {
      setFeedback({ tone: "error", text: `切档请求没有发出去：${requestError(err)}` });
    }
  }

  // 工作区: the path goes out verbatim; C26 resolves it natively and the
  // verdict lands in the snapshot's workspace section (set/spelling/canonical).
  function switchWorkspace() {
    const path = workspaceDraft.trim();
    if (path === "") return;
    try {
      requestWorkspaceChange(path);
      setWorkspaceDraft("");
    } catch (err) {
      setFeedback({ tone: "error", text: `工作区切换请求没有发出去：${requestError(err)}` });
    }
  }

  // 附件: the picked files' BYTES, one request each. The verdict (stored /
  // refused + reason) comes back as snapshot rows; a read failure is told to
  // the user, not swallowed - the alternative is the product eating the intent
  // it was given.
  async function attach(files: FileList | File[] | null) {
    if (!files || files.length === 0) return;
    setAttaching(true);
    try {
      for (const f of Array.from(files)) {
        await submitAttachment(f);
      }
    } catch (err) {
      setFeedback({ tone: "error", text: `附件没有送出去：${requestError(err)}` });
    } finally {
      setAttaching(false);
      if (fileInput.current) fileInput.current.value = "";
    }
  }

  // 发送: the text plus the attachments native ACCEPTED (a refused row is not
  // sent - it never became an attachment). Results arrive as pushed chunks in
  // the conversation column; only a failure of the request itself shows here.
  function send(text: string) {
    const body = text.trim();
    if (body === "" && stored.length === 0) return;
    try {
      sendMessage(body, stored);
    } catch (err) {
      setFeedback({ tone: "error", text: `消息没有发出去：${requestError(err)}` });
    }
  }

  return (
    <>
      {modelSwitch !== undefined && models.length === 0 && modelSwitch.reason !== "" ? (
        <p className="mb-1.5 text-[11px] leading-[1.6] text-ink-3" role="note">
          {modelSwitch.reason}（模型选择器不画）
        </p>
      ) : null}
      {/* 工作区行: the C26-authorized scope, plus the one real route that can
          move it. Drawn only when the pump actually carries a composer section
          - a read that was never assembled is not a row with nothing in it. */}
      {composer ? (
        <div className="mb-1.5 flex items-center gap-1.5 text-[11.5px] text-ink-3">
          <FolderCog aria-hidden="true" className="size-3.5 shrink-0" />
          <span
            className="min-w-0 flex-1 truncate font-mono"
            title={composer.workspace.set ? composer.workspace.canonical : composer.workspace.reason}
          >
            {composer.workspace.set ? composer.workspace.canonical : composer.workspace.reason}
          </span>
          {composer.workspace.reparse ? (
            <span className="shrink-0 text-warn">经 reparse 点（例外名单）</span>
          ) : null}
          <input
            aria-label="工作区路径（由原生侧解析与校验）"
            className="w-32 shrink-0 rounded-chip border border-line bg-surface px-1.5 py-0.5
              font-mono text-[11px] text-ink outline-none transition-colors duration-100
              focus:border-line-strong"
            onChange={(e) => setWorkspaceDraft(e.target.value)}
            placeholder={composer.workspace.spelling || "D:\\work\\project"}
            value={workspaceDraft}
          />
          <button
            aria-label="请求切换工作区"
            className="flex shrink-0 items-center gap-0.5 rounded-chip px-1 py-0.5 text-ink-3
              transition-colors duration-100 hover:bg-hover hover:text-ink
              disabled:pointer-events-none disabled:opacity-50"
            disabled={workspaceDraft.trim() === ""}
            onClick={switchWorkspace}
            type="button"
          >
            <ChevronRight aria-hidden="true" className="size-3" />
            切换
          </button>
        </div>
      ) : null}
      <ContextMeter usage={{}} className="mb-1.5" />
      <PromptBar
        attachments={stored.map((a) => ({ key: a.id, name: a.name }))}
        effort={modelSwitch !== undefined ? currentThinking : undefined}
        efforts={efforts}
        model={modelSwitch !== undefined ? currentPair : undefined}
        models={models}
        onAttach={composer ? () => fileInput.current?.click() : undefined}
        onEffortChange={(key) => void switchThinking(key)}
        onModelChange={(key) => void switchModel(key)}
        onPermissionChange={askMode}
        onSend={send}
        permission={permission}
        permissions={permissions}
      />
      <input
        accept={composer && composer.acceptedAttachmentMimes.length > 0 ? composer.acceptedAttachmentMimes.join(",") : undefined}
        className="hidden"
        multiple
        onChange={(e) => void attach(e.target.files)}
        ref={fileInput}
        type="file"
      />
      {/* The loud half of the attachment verdicts, verbatim from the snapshot:
          the bar's chips carry only stored rows, so a refusal lives here or it
          does not exist. */}
      {refused.length > 0 || (composer?.attachmentError ?? "") !== "" ? (
        <div className="mt-1">
          {composer?.attachmentError ? (
            <p className="text-[11px] leading-[1.6] text-red">{composer.attachmentError}</p>
          ) : null}
          {refused.map((a) => (
            <p className="text-[11px] leading-[1.6] text-red" key={a.id}>
              未发送：{a.name}——{a.reason}
            </p>
          ))}
        </div>
      ) : null}
      {attaching ? (
        <div className="mt-1">
          <FeedbackLine feedback={{ tone: "neutral", text: "附件读取中…" }} />
        </div>
      ) : null}
      {feedback !== null ? (
        <div className="mt-1">
          <FeedbackLine feedback={feedback} />
        </div>
      ) : null}
    </>
  );
}
