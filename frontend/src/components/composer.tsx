/* ============================================================================
   Composer (ticket 92; demo look pass 2026-09-25)
   ----------------------------------------------------------------------------
   The main panel's input row: attachment strip, workspace row, prompt bar and
   the permission-mode row, laid out like the demo chat screen's composer
   (design/doubao/demo/screens/chat.js:296-378). The mode and the workspace
   stay PERMISSION INPUTS, exactly as owner ruled it in R20/M5: everything
   shown is read from the snapshot Go pushes, and every interaction is a
   REQUEST through lib/panel's one envelope - requestModeSwitch /
   requestWorkspaceChange / submitAttachment / sendMessage - with the
   confirmation, the C26 resolution and the audit line happening natively.

   Demo affordances with no route behind them are absent rather than wired to
   a toast: no @ source button, no slash command button, no model chip, no mic
   button, no stop button, and no attachment remove button (there is no
   removal route; the strip renders the native verdicts only). A button that
   fakes a request would be worse than a missing button.

   Attachment BYTES do leave this file: submitAttachment reads the File and the
   payload travels as base64 in the same postMessage envelope, where the native
   side decodes, sniffs and stores it. The rows below are the native verdicts,
   so an unsupported type is reported to the user with its reason rather than
   disappearing.
   ============================================================================ */

import { useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import {
  ChevronRight,
  Film,
  FolderCog,
  Image,
  Paperclip,
  Send,
  type LucideIcon,
} from "lucide-react";
import { cn } from "@/lib/cn";
import {
  requestModeSwitch,
  requestWorkspaceChange,
  sendMessage,
  submitAttachment,
  type ComposerAttachment,
  type ComposerState,
} from "@/lib/panel";

/** Native mode names, spelled once so a typo cannot invent a fourth档. */
const MODE_LABELS: Record<string, string> = {
  ask_every_step: "每步都问",
  ask_high_risk: "高风险才问",
  auto_approve: "全自动",
  unknown: "档位未知",
};

function modeLabel(name: string): string {
  return MODE_LABELS[name] ?? name;
}

function bytes(n: number): string {
  if (n < 1024) return `${n} B`;
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`;
  return `${(n / (1024 * 1024)).toFixed(1)} MB`;
}

function attachmentIcon(kind: string): LucideIcon {
  if (kind === "image") return Image;
  if (kind === "video") return Film;
  return Paperclip;
}

export function Composer({
  state,
  onUserError,
}: {
  state: ComposerState;
  onUserError?: (message: string) => void;
}) {
  const [draft, setDraft] = useState("");
  const [busy, setBusy] = useState(false);
  const [workspaceDraft, setWorkspaceDraft] = useState("");
  const fileInput = useRef<HTMLInputElement>(null);

  // The native verdicts, in the order Go reported them: stored rows and refused
  // rows both come back through the snapshot, and each renders as itself.
  const accepted = state.attachments.filter((a) => a.stored);

  // The vocabulary plus the current value when the host has not listed it, so
  // an unreadable mode still shows as its own pill instead of vanishing.
  const modeNames = state.mode.names.includes(state.mode.current)
    ? state.mode.names
    : [state.mode.current, ...state.mode.names];

  function askMode(to: string) {
    // Request only. The native side raises the L2 card when `to` is the档 that
    // stops asking, and nothing here learns the outcome except through the
    // next snapshot.
    if (to === state.mode.current) return;
    requestModeSwitch(to);
  }

  async function attach(files: FileList | File[] | null) {
    if (!files) return;
    setBusy(true);
    try {
      for (const f of Array.from(files)) {
        await submitAttachment(f);
      }
    } catch (err) {
      // A failed read is told to the user, not swallowed: the alternative is
      // the product eating the intent it was given.
      onUserError?.(String(err instanceof Error ? err.message : err));
    } finally {
      setBusy(false);
      if (fileInput.current) fileInput.current.value = "";
    }
  }

  function send() {
    const text = draft.trim();
    if (!text && accepted.length === 0) return;
    sendMessage(text, accepted);
    setDraft("");
  }

  return (
    <div className="mt-2 flex w-full flex-col gap-2 border-t border-border px-1 pb-1 pt-3">
      {state.attachments.length > 0 ? (
        <ul className="m-0 flex list-none flex-wrap items-center gap-2 p-0">
          {state.attachments.map((a) => (
            <AttachmentRow key={a.id} a={a} />
          ))}
        </ul>
      ) : null}
      {state.attachmentError ? (
        <p className="text-[11px] text-destructive">{state.attachmentError}</p>
      ) : null}

      {/* 工作区行: the C26-authorized scope, plus the one real route that can
          move it. The new path is typed into the small mono field; the request
          goes out verbatim and the verdict comes back in the next snapshot. */}
      <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
        <FolderCog aria-hidden="true" className="size-3.5 shrink-0" />
        <span
          className="min-w-0 flex-1 truncate font-mono"
          title={state.workspace.set ? state.workspace.canonical : state.workspace.reason}
        >
          {state.workspace.set ? state.workspace.canonical : state.workspace.reason}
        </span>
        {state.workspace.reparse ? (
          <span className="shrink-0 text-warn">经 reparse 点（例外名单）</span>
        ) : null}
        <input
          aria-label="工作区路径（由原生侧解析与校验）"
          className="w-32 shrink-0 rounded-md border border-border bg-card px-1.5 py-0.5 font-mono text-[11px] text-foreground outline-none focus:border-ring"
          onChange={(e) => setWorkspaceDraft(e.target.value)}
          placeholder={state.workspace.spelling || "D:\\work\\project"}
          value={workspaceDraft}
        />
        <button
          className="flex shrink-0 items-center gap-0.5 transition-colors hover:text-foreground disabled:pointer-events-none disabled:opacity-50"
          disabled={workspaceDraft.trim() === ""}
          onClick={() => {
            requestWorkspaceChange(workspaceDraft.trim());
            setWorkspaceDraft("");
          }}
          type="button"
        >
          <ChevronRight aria-hidden="true" className="size-3" />
          切换
        </button>
      </div>

      {/* Prompt bar (demo .prompt-bar): attachment pick and send are the two
          affordances with real routes; the textarea is the message itself. */}
      <div className="prompt-bar">
        <button
          aria-label="添加附件"
          className="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
          disabled={busy}
          onClick={() => fileInput.current?.click()}
          title={`支持 ${state.acceptedAttachmentMimes.join(" / ")}，单个 <= ${bytes(state.maxAttachmentBytes)}`}
          type="button"
        >
          <Paperclip aria-hidden="true" className="size-4" />
        </button>
        <input
          accept={state.acceptedAttachmentMimes.join(",")}
          className="hidden"
          multiple
          onChange={(e) => void attach(e.target.files)}
          ref={fileInput}
          type="file"
        />
        <textarea
          aria-label="给 Wisp 的指令"
          onChange={(e) => setDraft(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              send();
            }
          }}
          onPaste={(e) => {
            const items = e.clipboardData?.files;
            if (items && items.length > 0) {
              e.preventDefault();
              void attach(items);
            }
          }}
          placeholder="输入消息或粘贴图片…"
          rows={1}
          value={draft}
        />
        <div className="prompt-actions">
          <Button aria-label="发送" onClick={send} size="icon-sm" variant="default">
            <Send aria-hidden="true" className="size-4" />
          </Button>
        </div>
      </div>

      {/* 底部行: the permission档 read out of the snapshot, then the input
          hint. The group label is the accessibility contract: pills only ever
          ASK, the native side answers. */}
      <div className="flex flex-wrap items-center justify-between gap-2">
        <div
          aria-label="权限档位（只读显示，点击仅发起切换请求）"
          className="flex flex-wrap items-center gap-1"
          role="group"
        >
          {modeNames.map((m) => (
            <button
              className={cn("sub-tab", m === state.mode.current && "active")}
              key={m}
              onClick={() => askMode(m)}
              type="button"
            >
              {modeLabel(m)}
              {state.mode.l2ConfirmNames.includes(m) ? (
                <span className="ml-1 text-[10px] text-warn">需原生 L2 强确认</span>
              ) : null}
            </button>
          ))}
        </div>
        <span className="shrink-0 text-[11px] text-muted-foreground">
          Enter 发送 · Shift+Enter 换行
        </span>
      </div>
    </div>
  );
}

/** One attachment: the native verdict, including the loud "no", in the demo
    strip's thumbnail shape. There is no remove button on purpose - the panel
    has no attachment removal route, and faking one would eat the user's file. */
function AttachmentRow({ a }: { a: ComposerAttachment }) {
  const Icon = attachmentIcon(a.kind);
  return (
    <li className="flex items-center gap-2 rounded-md border border-border bg-muted/40 px-2 py-1">
      <div className="flex size-8 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary">
        <Icon aria-hidden="true" className="size-4" />
      </div>
      <div className="flex min-w-0 flex-col leading-tight">
        <span className="truncate text-xs font-medium text-foreground">{a.name}</span>
        {a.stored ? (
          <span className="font-mono text-[10px] text-muted-foreground">
            {a.mime} · {bytes(a.sizeBytes)}
          </span>
        ) : null}
        <span
          className={a.stored ? "text-[10px] text-muted-foreground" : "text-[10px] text-destructive"}
        >
          {a.stored
            ? a.deduplicated
              ? "已存在相同内容的附件"
              : "已存入附件目录"
            : `未发送：${a.reason}`}
        </span>
      </div>
    </li>
  );
}
