/* ============================================================================
   HarnessMain - the harness APP page's main area (showcase demo data,
   ?harness=1 only; the product path never mounts this).
   ----------------------------------------------------------------------------
   Owner's shape ruling (2026-09-26): the main area has three states, selected
   by the `view` prop - the new-task entry page, the conversation view, and
   the no-history empty state. Pairs with HarnessSidebar (sidebar.tsx) to form
   the application skeleton; the right context column and the overlays are
   another agent's files.

   Owner's second form ruling (2026-09-26, ZCode/Qoder screenshots): the
   conversation view carries the mainstream harness furniture, in order -
     1. changes summary capsule  「已整理 N 个文件 +a -r」, FIRST thing under
        the session header; green/red tabular numbers; click expands the file
        list card (mono filename + action badge 已移动/已新建/已删除, first
        three rows then a 显示更多 toggle, rows with a detail line expand
        again to show it).
     2. elapsed row  「已工作 2 分 41 秒 >」 muted; click expands the
        behaviour demo EMPTY state (this walkthrough has no step log).
     3. composer pill row  above the PromptBar: the permission-tier pill
        (orange label + chevron, light bg + hairline border) and the mode
        label (neutral outlined pill). Both labels arrive via props (the
        session fixture fields permissionLabel / modeLabel) - nothing is
        hardcoded here. The PromptBar itself is NOT touched: it stays the
        library composer, untouched input body.
   All three are optional per session: a session without the fields renders
   without them (the contract lives in src/fixtures/harness-app.ts).

   PROPS CONTRACT - full interface in src/fixtures/harness-app.ts
   (HarnessMainProps); the component is props-driven, no data lives here:

     view: "new" | "session" | "empty"
         "new"    - greeting, the big input card (textarea + decorative send
                    spot), the workspace bar (mono path + 更换), quick-command
                    chips, and the recent-task rows.
         "session" - session header (title + workspace EntityChip + time),
                    the changes capsule + elapsed row (when the session
                    carries them), the message flow (user bubbles right-
                    aligned on bg-field; assistant blocks render Thinking /
                    ToolChips / StreamingText / cost line exactly when the
                    message carries those fields), with the composer pinned
                    to the bottom (mt-auto/flex).
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
import { ChevronDown, ChevronRight, FolderCog, Inbox, Plus } from "lucide-react";
import { cn } from "@/lib/cn";
import { Button } from "@/components/ai-native/button";
import { EntityChip } from "@/components/ai-native/entity-chip";
import { PromptBar } from "@/components/ai-native/prompt-bar";
import { StreamingText } from "@/components/ai-native/streaming-text";
import { Thinking } from "@/components/ai-native/thinking";
import { ToolChips } from "@/components/ai-native/tool-chips";
import {
  PROMPT_COMMANDS,
  PROMPT_MODELS,
  PROMPT_SOURCES,
  type HarnessAssistantMessage,
  type HarnessChangeFile,
  type HarnessChanges,
  type HarnessMainProps,
  type HarnessRecentTask,
  type HarnessSession,
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

/* 动作徽标的语义底色：新建绿、删除红、其余（已移动等）accent；未登记的
   动作文本落到中性 field 底。 */
const ACTION_TONES: Record<string, string> = {
  已新建: "bg-green-tint text-green",
  已删除: "bg-red-tint text-red",
  已移动: "bg-accent-tint text-accent-ink",
};

/** 文件清单卡默认只列前几行，其余进「显示更多」。 */
const CHANGES_PREVIEW_COUNT = 3;

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
  // 演示页允许本地 state：模型选择是 PromptBar 的受控演示值，发送不发请求。
  const [promptModel, setPromptModel] = useState(PROMPT_MODELS[0]?.key ?? "");

  return (
    <main
      className={cn(
        "screen-enter flex h-full min-w-0 flex-1 flex-col overflow-y-auto bg-page",
        className,
      )}
    >
      <div className="mx-auto flex w-full max-w-[680px] flex-1 flex-col justify-center gap-6 px-8 py-10">
        <h1 className="text-[22px] font-semibold tracking-wide text-ink">{greeting}</h1>

        {/* 大输入框 = beautiful-ui 的 Prompt Bar（完整解剖：附加菜单 / @ 来源 /
            / 命令 / 模型下拉 / 听写 / 发送）。演示页 onSend 只清空，不发请求。 */}
        <PromptBar
          autoFocus
          commands={PROMPT_COMMANDS}
          models={PROMPT_MODELS}
          model={promptModel}
          onModelChange={setPromptModel}
          onSend={() => undefined}
          placeholder={newPlaceholder}
          sources={PROMPT_SOURCES}
        />

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
   view="session" - 会话视图（头 + 更改摘要条 / 耗时行 + 消息流 + 钉底 composer）
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
  // 演示页允许本地 state：模型选择为受控演示值。
  const [promptModel, setPromptModel] = useState(PROMPT_MODELS[0]?.key ?? "");

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

      {/* 消息流：纵向滚动容器，内容限宽居中；会话头下第一件是更改摘要条
          （owner 截图形态），其后是耗时行，再往下才是消息。 */}
      <div className="flex min-h-0 flex-1 flex-col overflow-y-auto px-6 py-5">
        <div className="mx-auto flex w-full max-w-[720px] flex-col gap-5">
          {(session.changes || session.elapsed) && (
            <div className="flex flex-col gap-1.5">
              {session.changes ? <ChangesBar changes={session.changes} /> : null}
              {session.elapsed ? <ElapsedRow elapsed={session.elapsed} /> : null}
            </div>
          )}
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

      {/* composer 钉底：权限档/模式 pill 行 + beautiful-ui 的 Prompt Bar
          （组件原样保留，不拆）。演示不发真请求，onSend 只清空。 */}
      <div className="border-t border-line px-6 pb-4 pt-3">
        <div className="mx-auto w-full max-w-[720px]">
          {(session.permissionLabel || session.modeLabel) && (
            <div className="mb-1.5 flex items-center gap-1.5">
              {session.permissionLabel ? (
                <button
                  type="button"
                  title="权限档"
                  className="flex items-center gap-1 rounded-chip border border-warn-line bg-warn-soft
                    px-2 py-[3px] text-[11.5px] font-medium text-warn
                    transition-transform duration-100 active:scale-[0.97]"
                >
                  {session.permissionLabel}
                  <ChevronDown aria-hidden="true" size={11} strokeWidth={2.2} />
                </button>
              ) : null}
              {session.modeLabel ? (
                <span
                  className="flex items-center rounded-chip border border-line bg-field
                    px-2 py-[3px] text-[11.5px] text-ink-2"
                >
                  {session.modeLabel}
                </span>
              ) : null}
            </div>
          )}
          <PromptBar
            commands={PROMPT_COMMANDS}
            models={PROMPT_MODELS}
            model={promptModel}
            onModelChange={setPromptModel}
            onSend={() => undefined}
            placeholder={composerPlaceholder}
            sources={PROMPT_SOURCES}
          />
          <p className="mt-1.5 text-[10.5px] text-ink-3">
            Enter 发送 · Shift+Enter 换行 · 演示页不发起真实请求
          </p>
        </div>
      </div>
    </main>
  );
}

/* ---------------------------------------------------------------------------
   更改摘要条：圆角胶囊（已整理 N 个文件 +绿 -红）→ 展开文件清单卡
   （mono 文件名 + 动作徽标；默认前三行 + 显示更多；带 detail 的行可再展开）
   --------------------------------------------------------------------------- */

function ChangesBar({ changes }: { changes: HarnessChanges }) {
  const [open, setOpen] = useState(false);
  const [showAll, setShowAll] = useState(false);
  const [openDetails, setOpenDetails] = useState<ReadonlySet<string>>(new Set());

  const visible =
    showAll || changes.files.length <= CHANGES_PREVIEW_COUNT
      ? changes.files
      : changes.files.slice(0, CHANGES_PREVIEW_COUNT);

  function toggleDetail(name: string) {
    setOpenDetails((current) => {
      const next = new Set(current);
      if (next.has(name)) next.delete(name);
      else next.add(name);
      return next;
    });
  }

  return (
    <div style={{ animation: "fade-up 300ms var(--ease-out-strong) both" }}>
      {/* 胶囊：汇总行 + 审阅钮（ZCode 的 已编辑 5 个文件…审阅 形态） */}
      <div className="flex items-center gap-2">
      <button
        type="button"
        aria-expanded={open}
        onClick={() => setOpen((current) => !current)}
        className="flex items-center gap-2 rounded-full border border-line bg-surface px-3
          py-1.5 text-[12px] shadow-hairline transition-colors duration-150 hover:bg-hover"
      >
        <span className="text-ink-2">已整理 {changes.count} 个文件</span>
        <span className="font-medium tabular-nums text-green">+{changes.additions}</span>
        <span className="font-medium tabular-nums text-red">-{changes.removals}</span>
        <ChevronDown
          aria-hidden="true"
          size={11}
          strokeWidth={2.2}
          className={cn("text-ink-3 transition-transform duration-150", open && "rotate-180")}
        />
      </button>
      <button
        type="button"
        className="shrink-0 rounded-control border border-line bg-surface px-2.5 py-1.5 text-[11.5px]
          font-medium text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
      >
        审阅
      </button>
      </div>

      {/* 文件清单卡 */}
      {open && (
        <div
          className="mt-1.5 rounded-card border border-line bg-surface p-1 shadow-hairline"
          style={{ animation: "pop-in 180ms var(--ease-out-strong) both" }}
        >
          {visible.map((file) => (
            <div key={file.name}>
              {file.detail ? (
                <button
                  type="button"
                  aria-expanded={openDetails.has(file.name)}
                  onClick={() => toggleDetail(file.name)}
                  className="flex w-full items-center gap-2 rounded-chip px-2 py-1.5 text-left
                    transition-colors duration-100 hover:bg-hover"
                >
                  <FileRowCore expanded={openDetails.has(file.name)} file={file} />
                </button>
              ) : (
                <div className="flex w-full items-center gap-2 px-2 py-1.5">
                  <FileRowCore expanded={false} file={file} />
                </div>
              )}
              {file.detail && openDetails.has(file.name) ? (
                <p
                  className="pb-1.5 pl-8 pr-2 font-mono text-[10.5px] leading-relaxed text-ink-3"
                  style={{ animation: "fade-in 150ms ease both" }}
                >
                  {file.detail}
                </p>
              ) : null}
            </div>
          ))}
          {changes.files.length > CHANGES_PREVIEW_COUNT && (
            <button
              type="button"
              onClick={() => setShowAll((current) => !current)}
              className="flex w-full items-center gap-1 rounded-chip px-2 py-1.5 text-[11px]
                text-ink-3 transition-colors duration-100 hover:bg-hover hover:text-ink-2"
            >
              {showAll ? "收起" : `显示更多 ${changes.files.length - CHANGES_PREVIEW_COUNT} 个`}
              <ChevronDown
                aria-hidden="true"
                size={10}
                strokeWidth={2.2}
                className={cn("transition-transform duration-150", showAll && "rotate-180")}
              />
            </button>
          )}
        </div>
      )}
    </div>
  );
}

/** 清单卡一行的共同解剖：展开箭头位（有 detail 才有）+ mono 文件名 + 动作徽标。 */
function FileRowCore({ file, expanded }: { file: HarnessChangeFile; expanded: boolean }) {
  return (
    <>
      <span className="flex w-3.5 shrink-0 justify-center text-ink-3">
        {file.detail ? (
          <ChevronRight
            aria-hidden="true"
            size={11}
            strokeWidth={2.2}
            className={cn("transition-transform duration-150", expanded && "rotate-90")}
          />
        ) : null}
      </span>
      <span className="min-w-0 flex-1 truncate font-mono text-[11.5px] text-ink">{file.name}</span>
      <span
        className={cn(
          "shrink-0 rounded-chip px-1.5 py-0.5 text-[10.5px] font-medium",
          ACTION_TONES[file.action] ?? "bg-field text-ink-2",
        )}
      >
        {file.action}
      </span>
    </>
  );
}

/* ---------------------------------------------------------------------------
   耗时行：muted 文本 + 展开箭头；展开是行为演示空态（demo 无逐步行为记录）
   --------------------------------------------------------------------------- */

function ElapsedRow({ elapsed }: { elapsed: string }) {
  const [open, setOpen] = useState(false);

  return (
    <div>
      <button
        type="button"
        aria-expanded={open}
        onClick={() => setOpen((current) => !current)}
        className="flex items-center gap-1 text-[11.5px] text-ink-3
          transition-colors duration-100 hover:text-ink-2"
      >
        {elapsed}
        <ChevronRight
          aria-hidden="true"
          size={11}
          strokeWidth={2.2}
          className={cn("transition-transform duration-150", open && "rotate-90")}
        />
      </button>
      {open && (
        <div
          className="mt-1.5 rounded-card border border-line bg-inset px-3 py-2.5"
          style={{ animation: "fade-in 150ms ease both" }}
        >
          <p className="text-[11.5px] font-medium text-ink-2">行为时间线</p>
          <p className="mt-0.5 text-[11px] text-ink-3">演示空态：本会话没有可展示的逐步行为。</p>
        </div>
      )}
    </div>
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
      {msg.backgroundTools && msg.backgroundTools.length > 0 ? (
        <div className="flex flex-col gap-1">
          {msg.backgroundTools.map((n, i) => (
            <div className="flex items-center gap-2 text-[11.5px] text-ink-3" key={i}>
              <span className="flex size-4 items-center justify-center rounded-[4px] border border-line text-[9px] tabular-nums">{i + 1}</span>
              执行工具 {n} 次
            </div>
          ))}
        </div>
      ) : null}
      {msg.subagents && msg.subagents.length > 0 ? (
        <div className="flex flex-col gap-1">
          {msg.subagents.map((s) => (
            <div className="flex items-center gap-2 text-[11.5px]" key={s.id}>
              <span
                aria-hidden="true"
                className="size-4 shrink-0 rounded-full"
                style={{ background: "conic-gradient(var(--accent), var(--green), var(--orange), var(--accent))" }}
              />
              <span className="min-w-0 flex-1 truncate text-ink-2">{s.name}</span>
              <span className="shrink-0 tabular-nums text-ink-3">{s.duration}</span>
            </div>
          ))}
        </div>
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
