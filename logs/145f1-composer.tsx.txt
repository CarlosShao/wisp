/* ============================================================================
   Composer (ticket 92; ai-native look pass 2026-09-25)
   ----------------------------------------------------------------------------
   The main panel's input row. Visuals follow the library's prompt-bar and the
   chat.tsx composer shell: a rounded-control bg-field input card with
   focus-within:border-line-strong, the attachment strip above it, and the
   permission row below it. The mode and the workspace stay PERMISSION INPUTS,
   exactly as owner ruled in R20/M5: everything shown is read from the snapshot
   Go pushes, and every interaction is a REQUEST through lib/panel's one
   envelope - requestModeSwitch / requestWorkspaceChange / submitAttachment /
   sendMessage - with the confirmation, the C26 resolution and the audit line
   happening natively.

   来源 = derived from TurboKach/ai-native-react-components
          components/prompt-bar.tsx (MIT, Copyright (c) 2026 Turbo,
          upstream commit 05dab2d2); the input-card shell is chat.tsx:152's
          Tailwind string with its rgba shadows swapped for shadow-btn.
   本地改动：
   - props 化: 渲染逻辑照旧（ticket 92 的每条断言原样保留），只有观感换到库
     词汇；prompt-bar 的演示回路（AUTO_STEPS 自跑、@ / 斜杠菜单、模型选择、
     听写、glimm 扫光）没有对应路由，全部不渲染 - 一个假按钮比缺一个按钮糟。
   - 无路由不渲染: @、/、模型下拉、麦克风、停止、附件删除一概不画；附件条
     只渲染原生侧的裁决（存入/去重/拒绝原话）。
   - 字面量换 token: 附件 chip 用 bg-field / shadow-hairline，拒绝 chip 用
     bg-red-tint / text-red，档位 pill 用 rounded-full text-[11px] px-2.5 h-6
     （active 为 bg-accent-tint text-accent-ink），工作区输入框 border-line
     focus:border-line-strong；全文件无颜色字面量。
   ============================================================================ */

import { useRef, useState } from "react";
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
  const promptInput = useRef<HTMLTextAreaElement>(null);

  // The native verdicts, in the order Go reported them: stored rows and refused
  // rows both come back through the snapshot, and each renders as itself.
  const accepted = state.attachments.filter((a) => a.stored);
  const canSend = draft.trim().length > 0 || accepted.length > 0;

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
    <div className="mt-2 flex w-full flex-col gap-2.5 border-t border-line px-1 pb-1 pt-3">
      {state.attachments.length > 0 ? (
        <ul className="m-0 flex list-none flex-wrap items-center gap-1.5 p-0">
          {state.attachments.map((a) => (
            <AttachmentRow a={a} key={a.id} />
          ))}
        </ul>
      ) : null}
      {state.attachmentError ? (
        <p className="text-[11px] text-red">{state.attachmentError}</p>
      ) : null}

      {/* 工作区行: the C26-authorized scope, plus the one real route that can
          move it. The new path is typed into the small mono field; the request
          goes out verbatim and the verdict comes back in the next snapshot. */}
      <div className="flex items-center gap-1.5 text-[11.5px] text-ink-3">
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
          className="w-32 shrink-0 rounded-chip border border-line bg-surface px-1.5 py-0.5
            font-mono text-[11px] text-ink outline-none transition-colors duration-100
            focus:border-line-strong"
          onChange={(e) => setWorkspaceDraft(e.target.value)}
          placeholder={state.workspace.spelling || "D:\\work\\project"}
          value={workspaceDraft}
        />
        <button
          className="flex shrink-0 items-center gap-0.5 rounded-chip px-1 py-0.5 text-ink-3
            transition-colors duration-100 hover:bg-hover hover:text-ink
            disabled:pointer-events-none disabled:opacity-50"
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

      {/* Prompt bar: chat.tsx:152 的外壳 Tailwind 原样内联，rgba 阴影换成
          shadow-btn。附件挑选与发送是唯二有真实路由的动作。 */}
      <div
        className="flex cursor-text flex-col gap-1.5 rounded-control border border-line bg-field
          p-2.5 shadow-btn transition-[border-color,box-shadow] duration-150
          focus-within:border-line-strong"
        onClick={() => promptInput.current?.focus()}
        role="presentation"
      >
        <div className="flex items-end gap-1">
          <button
            aria-label="添加附件"
            className="flex size-7 shrink-0 items-center justify-center rounded-[8px] text-ink-3
              transition-[background-color,color,transform] duration-150 hover:bg-hover
              hover:text-ink active:scale-[0.94] disabled:pointer-events-none disabled:opacity-50"
            disabled={busy}
            onClick={() => fileInput.current?.click()}
            title={`支持 ${state.acceptedAttachmentMimes.join(" / ")}，单个 <= ${bytes(state.maxAttachmentBytes)}`}
            type="button"
          >
            <Paperclip aria-hidden="true" className="size-4" strokeWidth={2} />
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
            className="min-h-7 min-w-0 w-full resize-none bg-transparent px-1 py-[5px]
              text-[13px] leading-[18px] text-ink outline-none [overflow-wrap:anywhere]
              placeholder:text-ink-3"
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
            ref={promptInput}
            rows={1}
            value={draft}
          />
          <button
            aria-label="发送"
            className="flex size-7 shrink-0 items-center justify-center rounded-[8px]
              transition-[background-color,color,transform] duration-200
              enabled:active:scale-[0.94] disabled:pointer-events-none disabled:opacity-50"
            disabled={!canSend}
            onClick={send}
            style={{
              background: canSend ? "var(--ink)" : "var(--line-strong)",
              color: canSend ? "var(--surface)" : "var(--ink-2)",
            }}
            type="button"
          >
            <Send aria-hidden="true" className="size-4" strokeWidth={2.4} />
          </button>
        </div>
      </div>

      {/* 底部行: the permission档 read out of the snapshot, then the input
          hint. The group label is the accessibility contract: pills only ever
          ASK, the native side answers. */}
      <div className="flex flex-wrap items-center justify-between gap-2">
        <div
          aria-label="权限档位（只读显示，点击仅发起切换请求）"
          className="flex flex-wrap items-center gap-1.5"
          role="group"
        >
          {modeNames.map((m) => (
            <button
              className={cn(
                "inline-flex h-6 items-center rounded-full px-2.5 text-[11px] font-medium transition-colors duration-100",
                m === state.mode.current
                  ? "bg-accent-tint text-accent-ink"
                  : "bg-inset text-ink-2 hover:bg-hover",
              )}
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
        <span className="shrink-0 text-[11px] text-ink-3">Enter 发送 · Shift+Enter 换行</span>
      </div>
    </div>
  );
}

/** One attachment: the native verdict, including the loud "no", in the prompt
    bar's chip shape. There is no remove button on purpose - the panel has no
    attachment removal route, and faking one would eat the user's file. A
    refusal renders with its full reason: a truncated refusal reads as if the
    file had been accepted. */
function AttachmentRow({ a }: { a: ComposerAttachment }) {
  const Icon = attachmentIcon(a.kind);
  if (!a.stored) {
    return (
      <li className="flex max-w-full flex-col gap-0.5 rounded-chip bg-red-tint px-2 py-1.5 shadow-hairline">
        <span className="flex items-center gap-1.5 text-[11.5px] font-medium text-red">
          <Icon aria-hidden="true" className="size-3.5 shrink-0" />
          {a.name}
        </span>
        <span className="text-[10.5px] leading-[1.6] text-red">未发送：{a.reason}</span>
      </li>
    );
  }
  return (
    <li
      className="flex h-6.5 max-w-full items-center gap-1.5 rounded-chip bg-field py-1 pr-2.5
        pl-1.5 text-[11.5px] text-ink-2 shadow-hairline"
      style={{ animation: "pop-in 200ms var(--ease-out-strong) both" }}
    >
      <Icon aria-hidden="true" className="size-3.5 shrink-0" />
      <span className="max-w-36 truncate">{a.name}</span>
      <span className="shrink-0 font-mono text-[10.5px] text-ink-3">
        {a.mime} · {bytes(a.sizeBytes)}
      </span>
      <span className="shrink-0 text-[10.5px] text-green">
        {a.deduplicated ? "已存在相同内容的附件" : "已存入附件目录"}
      </span>
    </li>
  );
}
