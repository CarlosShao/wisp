/* ============================================================================
   Composer v2 (ticket 92) — the main panel's input row: 档位显示 + 附件 +
   工作区选择, exactly as owner ruled it in R20/M5, with git 分支/仓库切换 left
   out because he cut it ("我们产品不要求 coding 能力").

   The rule this component is written against, and the reason it looks dull:
   the mode and the workspace are PERMISSION INPUTS. PLAN.md:1588 keeps an allow
   decision on the native side, and the same logic says a renderer must not be
   able to move the session onto a looser档 or a different tree by itself - a
   compromised page that could do either would hold a permanent approval pass.
   So:

     - everything shown here is read out of the snapshot the Go side pushes
       (AC#4: a closed panel reopens with the same mode and the same workspace,
       because neither was ever stored here);
     - every interaction sends a REQUEST through lib/panel's one envelope -
       requestModeSwitch / requestWorkspaceChange / submitAttachment - and the
       confirmation, the C26 resolution and the audit line happen natively;
     - there is no page-side mode setter in this file, and no second channel to
       hang one on; internal/panel/composer_test.go plants that mistake in a
       fixture and asserts two separate instruments go red when it appears.

   Attachment BYTES do leave this file: submitAttachment reads the File and the
   payload travels as base64 in the same postMessage envelope, where the native
   side decodes, sniffs and stores it (internal/panel/attachments.go). The rows
   rendered below are the native verdicts, so an unsupported type is reported to
   the user with its reason rather than disappearing.
   ============================================================================ */

import { useRef, useState } from "react";
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

  function askMode(to: string) {
    // Request only. The native side raises the L2 card when `to` is the档 that
    // stops asking, and nothing here learns the outcome except through the next
    // snapshot.
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
    <div className="glass-inset flex w-full flex-col gap-2 rounded-control border border-line p-3">
      <div className="flex items-center gap-2 text-[11.5px]">
        <label className="flex items-center gap-1.5" title="档位由原生侧决定并记录；这里只能发起请求">
          <span className="text-ink-3">档位</span>
          <select
            aria-label="权限档位（只读显示，点击仅发起切换请求）"
            className="rounded-control border border-line bg-bg-inset px-1.5 py-0.5 text-ink"
            value={state.mode.current}
            onChange={(e) => askMode(e.target.value)}
          >
            {(state.mode.names.includes(state.mode.current)
              ? state.mode.names
              : [state.mode.current, ...state.mode.names]
            ).map((m) => (
              <option key={m} value={m}>
                {modeLabel(m)}
                {state.mode.l2ConfirmNames.includes(m) ? "（需原生 L2 强确认）" : ""}
              </option>
            ))}
          </select>
        </label>
        <span className="text-ink-3">工作区</span>
        <span className="truncate text-ink" title={state.workspace.canonical || state.workspace.reason}>
          {state.workspace.set ? state.workspace.canonical : state.workspace.reason}
        </span>
        {state.workspace.reparse ? (
          <span className="text-caution">经 reparse 点（例外名单）</span>
        ) : null}
      </div>

      <textarea
        aria-label="给 Wisp 的指令"
        className="min-h-[52px] w-full resize-y rounded-control border border-line bg-bg-inset px-2 py-1.5 text-[12.5px] leading-[var(--lh-body)] text-ink"
        placeholder="要做什么？可加附件（图片/视频），或改工作区"
        value={draft}
        onChange={(e) => setDraft(e.target.value)}
        onPaste={(e) => {
          const items = e.clipboardData?.files;
          if (items && items.length > 0) {
            e.preventDefault();
            void attach(items);
          }
        }}
      />

      <div className="flex flex-wrap items-center gap-2 text-[11.5px]">
        <button
          type="button"
          className="rounded-control border border-line px-2 py-1 text-ink"
          onClick={() => fileInput.current?.click()}
          disabled={busy}
        >
          添加附件
        </button>
        <input
          ref={fileInput}
          type="file"
          multiple
          className="hidden"
          accept={state.acceptedAttachmentMimes.join(",")}
          onChange={(e) => void attach(e.target.files)}
        />
        <span className="text-ink-3">
          支持 {state.acceptedAttachmentMimes.join(" / ")}，单个 ≤ {bytes(state.maxAttachmentBytes)}
        </span>
        <button
          type="button"
          className="ml-auto rounded-control border border-line px-2 py-1 text-ink"
          onClick={send}
        >
          发送
        </button>
      </div>

      {state.attachments.length > 0 ? (
        <ul className="flex flex-col gap-1 text-[11.5px]">
          {state.attachments.map((a) => (
            <AttachmentRow key={a.id} a={a} />
          ))}
        </ul>
      ) : null}
      {state.attachmentError ? <p className="text-caution">{state.attachmentError}</p> : null}

      <div className="flex items-center gap-2 text-[11.5px]">
        <input
          aria-label="工作区路径（由原生侧解析与校验）"
          className="min-w-0 flex-1 rounded-control border border-line bg-bg-inset px-2 py-1 text-ink"
          placeholder={state.workspace.spelling || "D:\\work\\project"}
          value={workspaceDraft}
          onChange={(e) => setWorkspaceDraft(e.target.value)}
        />
        <button
          type="button"
          className="rounded-control border border-line px-2 py-1 text-ink"
          disabled={workspaceDraft.trim() === ""}
          onClick={() => {
            requestWorkspaceChange(workspaceDraft.trim());
            setWorkspaceDraft("");
          }}
        >
          换工作区
        </button>
      </div>
    </div>
  );
}

/** One attachment row: the native verdict, including the loud "no". */
function AttachmentRow({ a }: { a: ComposerAttachment }) {
  return (
    <li className={a.stored ? "flex gap-2 text-ink-3" : "flex gap-2 text-caution"}>
      <span className="truncate">
        {a.name}
        {a.stored ? ` · ${a.mime} · ${bytes(a.sizeBytes)}` : ""}
      </span>
      <span className="text-ink-3">
        {a.stored
          ? a.deduplicated
            ? "已存在相同内容的附件"
            : "已存入附件目录"
          : `未发送：${a.reason}`}
      </span>
    </li>
  );
}
