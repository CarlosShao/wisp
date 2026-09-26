/* ============================================================================
   HarnessMain - the harness APP page's main area (showcase demo data,
   ?harness=1 only; the product path never mounts this).
   ----------------------------------------------------------------------------
   Owner's shape ruling (2026-09-26): the main area has three states, selected
   by the `view` prop - the new-task entry page, the conversation view, and
   the no-history empty state. Pairs with HarnessSidebar (sidebar.tsx) to form
   the application skeleton; the right context column and the overlays are
   another agent's files.

   PROPS CONTRACT - full interface in src/fixtures/harness-app.ts
   (HarnessMainProps); the component is props-driven, no data lives here:

     view: "new" | "session" | "empty"
         "new"    - greeting, the big input card (textarea + decorative send
                    spot), the workspace bar (mono path + 更换), quick-command
                    chips, and the recent-task rows.
         "session" - session header (title + workspace EntityChip + time) and
                    the message flow (user bubbles right-aligned on bg-field;
                    assistant blocks render Thinking / ToolChips / StreamingText
                    / cost line exactly when the message carries those fields),
                    with the composer pinned to the bottom (mt-auto/flex).
                    `session: null` falls back to the empty state.
         "empty"  - centered inbox glyph + hint + optional 新建任务 button.
     greeting / workspace / quickCommands / recentTasks / session / emptyHint
         Data props as listed in HarnessMainProps.
     newPlaceholder? / composerPlaceholder?
         Input placeholders; defaults "描述你要做的事…" and "输入消息…".
     onPickWorkspace? / onPickCommand? / onPickTask? / onNewTask?
         Callbacks; the composer does NOT send requests - it keeps its draft
         in local demo state and the send button just clears it (the hint
         line under the shell says so).

   Motion policy: html[data-motion="off"] (theme.css) already silences every
   CSS keyframe in the tree. The one JS-driven mount reachable here is
   RevealText (inside StreamingText): its timer does not read the attribute,
   so the streaming flag is forced false under motion-off and the static end
   value renders instead - the one-line condition motionOff() below.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useState } from "react";
import { ChevronRight, FolderCog, Inbox, Paperclip, Plus, Send } from "lucide-react";
import { cn } from "@/lib/cn";
import { Button } from "@/components/ai-native/button";
import { EntityChip } from "@/components/ai-native/entity-chip";
import { StreamingText } from "@/components/ai-native/streaming-text";
import { Thinking } from "@/components/ai-native/thinking";
import { ToolChips } from "@/components/ai-native/tool-chips";
import type {
  HarnessAssistantMessage,
  HarnessMainProps,
  HarnessRecentTask,
  HarnessSession,
} from "@/fixtures/harness-app";

/* RevealText 的逐段计时器不读 html[data-motion="off"]，挂载点在这里代读：
   off 时 streaming=false，直接渲染静态终值（一行条件）。 */
function motionOff(): boolean {
  return document.documentElement.dataset.motion === "off";
}

/* 最近任务行的语义色点（与侧栏徽标同一套语义）。 */
const STATUS_DOTS: Record<HarnessRecentTask["status"], string> = {
  done: "bg-green",
  streaming: "bg-accent",
  denied: "bg-red",
};

export function HarnessMain({
  view,
  greeting,
  workspace,
  quickCommands,
  recentTasks,
  session,
  emptyHint,
  newPlaceholder = "描述你要做的事…",
  composerPlaceholder = "输入消息…",
  onPickWorkspace,
  onPickCommand,
  onPickTask,
  onNewTask,
  className,
}: HarnessMainProps) {
  if (view === "session") {
    return session ? (
      <SessionView composerPlaceholder={composerPlaceholder} className={className} session={session} />
    ) : (
      <EmptyView className={className} hint={emptyHint} onNewTask={onNewTask} />
    );
  }
  if (view === "empty") {
    return <EmptyView className={className} hint={emptyHint} onNewTask={onNewTask} />;
  }
  return (
    <NewTaskView
      className={className}
      greeting={greeting}
      newPlaceholder={newPlaceholder}
      onPickCommand={onPickCommand}
      onPickTask={onPickTask}
      onPickWorkspace={onPickWorkspace}
      quickCommands={quickCommands}
      recentTasks={recentTasks}
      workspace={workspace}
    />
  );
}

/* ---------------------------------------------------------------------------
   view="new" - 新任务入口页
   --------------------------------------------------------------------------- */

function NewTaskView({
  greeting,
  workspace,
  quickCommands,
  recentTasks,
  newPlaceholder,
  onPickWorkspace,
  onPickCommand,
  onPickTask,
  className,
}: {
  greeting: string;
  workspace: string;
  quickCommands: HarnessMainProps["quickCommands"];
  recentTasks: HarnessRecentTask[];
  newPlaceholder: string;
  onPickWorkspace?: () => void;
  onPickCommand?: (command: HarnessMainProps["quickCommands"][number]) => void;
  onPickTask?: (task: HarnessRecentTask) => void;
  className?: string;
}) {
  // 演示页允许本地 state：大输入框只是摆出可输入的样子，不发起任何请求。
  const [draft, setDraft] = useState("");
  const filled = draft.trim().length > 0;

  return (
    <main
      className={cn(
        "screen-enter flex h-full min-w-0 flex-1 flex-col overflow-y-auto bg-page",
        className,
      )}
    >
      <div className="mx-auto flex w-full max-w-[680px] flex-1 flex-col justify-center gap-6 px-8 py-10">
        <h1 className="text-[22px] font-semibold tracking-wide text-ink">{greeting}</h1>

        {/* 大输入框卡：rounded-card + 右下发送钮位装饰（演示页不发请求） */}
        <div className="rounded-card border border-line bg-surface p-4 shadow-card">
          <textarea
            aria-label="新任务描述"
            className="min-h-[72px] w-full resize-none bg-transparent text-[13.5px] leading-relaxed text-ink outline-none [overflow-wrap:anywhere] placeholder:text-ink-3"
            onChange={(e) => setDraft(e.target.value)}
            placeholder={newPlaceholder}
            rows={3}
            value={draft}
          />
          <div className="mt-2 flex items-center justify-end">
            <span
              aria-hidden="true"
              className="flex size-8 items-center justify-center rounded-[8px] transition-colors duration-200"
              style={{
                background: filled ? "var(--ink)" : "var(--line-strong)",
                color: filled ? "var(--surface)" : "var(--ink-2)",
              }}
            >
              <Send className="size-4" strokeWidth={2.4} />
            </span>
          </div>
        </div>

        {/* 工作区选择条 */}
        <div className="flex items-center gap-1.5 text-[11.5px] text-ink-3">
          <FolderCog aria-hidden="true" className="size-3.5 shrink-0" strokeWidth={2} />
          <span className="min-w-0 flex-1 truncate font-mono" title={workspace}>
            {workspace}
          </span>
          <button
            className="flex shrink-0 items-center gap-0.5 rounded-chip px-1.5 py-0.5 text-ink-2
              transition-colors duration-100 hover:bg-hover hover:text-ink"
            onClick={onPickWorkspace}
            type="button"
          >
            <ChevronRight aria-hidden="true" className="size-3" />
            更换
          </button>
        </div>

        {/* 快捷指令 chips */}
        <div className="flex flex-wrap items-center gap-1.5">
          {quickCommands.map((command, i) => (
            <button
              className="rounded-full bg-accent-tint px-2.5 py-1 text-[11.5px] font-medium
                text-accent-ink transition-[opacity,transform] duration-150
                hover:opacity-90 active:scale-[0.96]"
              key={command.id}
              onClick={() => onPickCommand?.(command)}
              style={{ animation: `fade-up 300ms var(--ease-out-strong) ${i * 80}ms both` }}
              type="button"
            >
              {command.label}
            </button>
          ))}
        </div>

        {/* 最近任务三行 */}
        <div>
          <p className="px-1 pb-1.5 text-[11px] uppercase tracking-wide text-ink-3">最近任务</p>
          <div className="flex flex-col">
            {recentTasks.map((task, i) => (
              <button
                className="flex w-full items-center gap-2 rounded-control px-2 py-1.5 text-left
                  transition-colors duration-100 hover:bg-hover"
                key={task.id}
                onClick={() => onPickTask?.(task)}
                style={{ animation: `fade-up 320ms var(--ease-out-strong) ${i * 90}ms both` }}
                type="button"
              >
                <span
                  aria-hidden="true"
                  className={cn("size-1.5 shrink-0 rounded-full", STATUS_DOTS[task.status])}
                />
                <span className="min-w-0 flex-1 truncate text-[12.5px] text-ink">{task.title}</span>
                <span className="shrink-0 text-[11px] text-ink-3">{task.meta}</span>
              </button>
            ))}
          </div>
        </div>
      </div>
    </main>
  );
}

/* ---------------------------------------------------------------------------
   view="session" - 会话视图（头 + 消息流 + 钉底 composer）
   --------------------------------------------------------------------------- */

function SessionView({
  session,
  composerPlaceholder,
  className,
}: {
  session: HarnessSession;
  composerPlaceholder: string;
  className?: string;
}) {
  // 演示页允许本地 state：composer 不发真请求，发送钮只清空草稿。
  const [draft, setDraft] = useState("");
  const canSend = draft.trim().length > 0;

  return (
    <main
      className={cn("screen-enter flex h-full min-w-0 flex-1 flex-col bg-page", className)}
    >
      {/* 会话头：标题 + 工作区 EntityChip + 日期 */}
      <header className="flex items-center gap-2.5 border-b border-line px-6 py-3">
        <h1 className="min-w-0 truncate text-[14px] font-semibold text-ink">{session.title}</h1>
        <EntityChip color={session.workspaceColor} name={session.workspace} />
        <span className="ml-auto shrink-0 text-[11.5px] tabular-nums text-ink-3">
          {session.time}
        </span>
      </header>

      {/* 消息流：纵向滚动容器，内容限宽居中 */}
      <div className="flex min-h-0 flex-1 flex-col overflow-y-auto px-6 py-5">
        <div className="mx-auto flex w-full max-w-[720px] flex-col gap-5">
          {session.messages.map((msg) =>
            msg.role === "user" ? (
              <div className="flex justify-end" key={msg.id}>
                <p
                  className="max-w-[80%] rounded-card bg-field px-3.5 py-2.5 text-[13px]
                    leading-relaxed text-ink shadow-hairline"
                  style={{ animation: "fade-up 300ms var(--ease-out-strong) both" }}
                >
                  {msg.text}
                </p>
              </div>
            ) : (
              <AssistantBlock key={msg.id} msg={msg} />
            ),
          )}
        </div>
      </div>

      {/* composer 钉底（chat.tsx:152 外壳词汇：rounded-control / bg-field /
          focus-within:border-line-strong / shadow-btn）。演示不发真请求。 */}
      <div className="border-t border-line px-6 pb-4 pt-3">
        <div className="mx-auto w-full max-w-[720px]">
          <div
            className="flex cursor-text flex-col gap-1.5 rounded-control border border-line
              bg-field p-2.5 shadow-btn transition-[border-color,box-shadow] duration-150
              focus-within:border-line-strong"
          >
            <div className="flex items-end gap-1">
              {/* 附件位装饰（aria-hidden，无路由不假装可点） */}
              <span
                aria-hidden="true"
                className="flex size-7 shrink-0 items-center justify-center rounded-[8px] text-ink-3"
              >
                <Paperclip className="size-4" strokeWidth={2} />
              </span>
              <textarea
                aria-label="给一缕的指令"
                className="min-h-7 min-w-0 w-full resize-none bg-transparent px-1 py-[5px]
                  text-[13px] leading-[18px] text-ink outline-none [overflow-wrap:anywhere]
                  placeholder:text-ink-3"
                onChange={(e) => setDraft(e.target.value)}
                placeholder={composerPlaceholder}
                rows={1}
                value={draft}
              />
              <button
                aria-label="发送（演示位，不发起请求）"
                className="flex size-7 shrink-0 items-center justify-center rounded-[8px]
                  transition-[background-color,color,transform] duration-200
                  enabled:active:scale-[0.94] disabled:pointer-events-none disabled:opacity-50"
                disabled={!canSend}
                onClick={() => setDraft("")}
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
          <p className="mt-1.5 text-[10.5px] text-ink-3">
            Enter 发送 · Shift+Enter 换行 · 演示页不发起真实请求
          </p>
        </div>
      </div>
    </main>
  );
}

/* 助手块：thinking / 工具条 / 正文流式 / 成本行，字段缺省即不渲染。 */
function AssistantBlock({ msg }: { msg: HarnessAssistantMessage }) {
  return (
    <div className="flex flex-col gap-2.5">
      {msg.thinking ? (
        <Thinking seconds={msg.thinking.seconds} steps={msg.thinking.steps} />
      ) : null}
      {msg.tools && msg.tools.length > 0 ? <ToolChips tools={msg.tools} /> : null}
      {msg.text ? (
        <StreamingText
          sources={msg.sources}
          streaming={motionOff() ? false : Boolean(msg.streaming)}
          text={msg.text}
        />
      ) : null}
      {msg.costLine ? (
        <p className="font-mono text-[11px] tabular-nums text-ink-3">{msg.costLine}</p>
      ) : null}
    </div>
  );
}

/* ---------------------------------------------------------------------------
   view="empty" - 无历史空态
   --------------------------------------------------------------------------- */

function EmptyView({
  hint,
  onNewTask,
  className,
}: {
  hint: string;
  onNewTask?: () => void;
  className?: string;
}) {
  return (
    <main
      className={cn(
        "screen-enter flex h-full min-w-0 flex-1 flex-col items-center justify-center gap-3",
        "bg-page px-8 text-center",
        className,
      )}
    >
      <span className="flex size-11 items-center justify-center rounded-full bg-field shadow-hairline">
        <Inbox aria-hidden="true" className="size-5 text-ink-3" strokeWidth={1.8} />
      </span>
      <p className="text-[13.5px] font-medium text-ink">还没有会话</p>
      <p className="max-w-[320px] text-[12px] leading-relaxed text-ink-3">{hint}</p>
      {onNewTask ? (
        <Button className="mt-1" onClick={onNewTask} size="sm" variant="accent">
          <Plus aria-hidden="true" className="size-3.5" strokeWidth={2.2} />
          新建任务
        </Button>
      ) : null}
    </main>
  );
}
