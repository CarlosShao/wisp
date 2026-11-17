/* ============================================================================
   MonitorPopover - the task-monitor popover, ZCode shape (owner 2026-09-26
   screenshot: a small card anchored under the header button - title + pin
   slot, an 环境信息 group of icon/label/value rows, a 子智能体 group with
   gradient-avatar agent rows and status dots).
   ----------------------------------------------------------------------------
   Props-driven; the ANCHOR (absolute positioning shell) belongs to the
   parent - this component only paints the card. Zero data of its own.
   ============================================================================ */

import { useRef, useState } from "react";
import { Pin } from "lucide-react";
import { GlideMenu } from "@/components/ai-native/glide-menu";
import { cn } from "@/lib/cn";
import { formatHistoryTs, scanTaskForensics, appliedStepsReading, type AppliedStepsReading, type HistoryLoader } from "@/lib/history-view";
import { REFUSED_ENVELOPE_PREFIX } from "@/lib/privacy-view";
import type { RunUsageRow } from "@/lib/panel";

export interface MonitorEnvRow {
  label: string;
  value: string;
}

/** One row of the task monitor as the WIRE carries it (App derives these from
    lib/panel.ts TaskRowView, ticket 197). Two rules ride over from the wire:
    blockedOnApproval is the host's own join - a card for this very task is on
    the queue right now (tickets 197/220) - and status is a D43 name that is
    empty exactly when statusKnown is false. This layer displays those words;
    it never re-judges a status and never invents one. */
export interface MonitorAgentRow {
  id: string;
  name: string;
  blockedOnApproval: boolean;
  status: string;
  statusKnown: boolean;
  /** The host's own sentence for a status it did not file (wire statusReason). */
  statusReason?: string;
  /** This row's own work page cost, verbatim off the wire (TaskRowView's
      stream triple, internal/panel/subagent_roster_197.go:140-142). A DISPLAY
      suffix, never a status input: "少了一段" rides beside the status word,
      it does not replace or recolour it. Undefined = the wire did not carry
      the row-level numbers (older hosts), which renders nothing - a missing
      reading is never drawn as a zero. */
  streamTruncated?: boolean;
  streamElidedRunes?: number;
  streamDropped?: boolean;
  /** ---- readface wave 2 (2026-11-17), ticket 38 AC#4's three task-card
      keys, copied verbatim off the wire (App's monitorAgents maps them
      1:1 from TaskRowView). All three are omitempty on the wire, so ABSENT
      is the host's own "stamped nothing" - the row renders 未判 for that
      dimension and NEVER an epoch, a zero cost, or a fake step name.
      EXCEPTION: currentStep's absence is the wire's own "no call in
      flight" (the loop retires the row at EvToolEnd), which reads 无在飞 -
      that is a reading, not an unknown. */
  startedAt?: number;
  costToDate?: RunUsageRow;
  currentStep?: string;
}

/** The doors a hosted panel hands the cancel flow: roster.cancel (the ONE
    stop path the C17 roster carries) and the history.page loader the
    applied-steps reading counts its forensics rows from. Both undefined in
    the showcase/demo shapes: rows then draw no 停止 control at all - a
    button wired to a door that does not exist is the lie the render-nav
    gate names. */
export interface TaskCancelDoors {
  cancelTask: (taskId: string) => Promise<string | null>;
  historyLoad: HistoryLoader;
}

/** What one roster.cancel exchange turned out to be, per task. The receipt
    and the refusal ride VERBATIM (the door's own words); the applied-steps
    reading arrives from the forensics scan and stays undefined until it
    completes - an incomplete scan is 未判, never a partial count. */
export type TaskCancelState =
  | { phase: "sending" }
  | { phase: "receipt"; text: string; scanning: boolean; reading?: AppliedStepsReading; scanNote?: string }
  | { phase: "refused"; raw: string }
  | { phase: "no-host" };

/** The dot's vocabulary - a DISPLAY tier, not a status judgement. blocked is
    the 等批准 tier: the row's task has its own card waiting. idle is "the host
    filed no status word", a different fact from any of the others. */
type AgentDot = "running" | "done" | "failed" | "blocked" | "idle";

/** D43 状态名 -> the Chinese the row shows (internal/statemachine/states.go is
    the source of the names; App's TIER_NAMES is the precedent for spelling the
    display map in the component that draws it). A name missing here falls to
    the honest line below, never to a guessed word. */
const D43_LABELS: Record<string, string> = {
  FirstRun: "首次运行",
  Sleeping: "休眠",
  Armed: "已就绪",
  Muted: "已静音",
  Listening: "在听",
  Thinking: "思考中",
  Acting: "执行中",
  Speaking: "回话中",
  Warm: "预热",
  Conversation: "对话中",
  Confirming: "待确认（L1）",
  AwaitingApproval: "待审批（L2）",
  Settling: "收尾中",
  Downloading: "下载中",
  Unconfigured: "未配置",
  NoNetwork: "无网络",
  Error: "出错",
  Queued: "排队中",
  Stuck: "卡住",
  WatchdogAlert: "看门狗告警",
};

/** Only the failure family goes red; Settling is the one name that reads as
    wrapping up. Every other known status is live work (accent); a status the
    host never filed is idle grey. */
const D43_DOT: Record<string, AgentDot> = {
  Error: "failed",
  Stuck: "failed",
  WatchdogAlert: "failed",
  Settling: "done",
};

/** One row's displayed word and dot, derived from the wire fields and nothing
    else. 等批准 outranks the status word: a task whose own card is on the
    queue is waiting on the user, whatever its D43 name currently says. */
function agentDisplay(row: MonitorAgentRow): { label: string; dot: AgentDot } {
  if (row.blockedOnApproval) return { label: "等批准", dot: "blocked" };
  if (row.statusKnown && D43_LABELS[row.status] !== undefined) {
    return { label: D43_LABELS[row.status]!, dot: D43_DOT[row.status] ?? "running" };
  }
  return { label: row.statusReason?.trim() || "宿主没有给出状态", dot: "idle" };
}

const DOT_CLASS: Record<AgentDot, string> = {
  running: "bg-accent",
  done: "bg-green",
  failed: "bg-red",
  blocked: "bg-warn",
  idle: "bg-ink-3",
};

/** The row's own truncation suffix, or null when this row's stream is whole.
    Dropped outranks truncated: a dropped stream is not "partly shown", the
    page for that row is GONE. Only an actual loss renders anything - a clean
    row stays clean, because a permanent "0 字省略" chip on every row would be
    alarm fatigue (the same argument App's header uses for the approval chip:
    a badge that is always on trains you to ignore it). */
function rowBoundLabel(row: MonitorAgentRow): string | null {
  if (row.streamDropped) return "流已丢";
  if (row.streamTruncated) return `省 ${row.streamElidedRunes ?? 0} 字`;
  return null;
}

/** The whole log's truncation sentence off the SECTION numbers (wire
    TaskRosterSection's streamTruncated/streamElidedRunes/droppedCount, Go
    :166-183), or null on a clean log. Same display rule as the row suffix:
    abnormal states only. Both facts when both hold, in one line. */
function sectionBoundLine(
  truncated: boolean,
  elidedRunes: number | undefined,
  droppedCount: number | undefined,
): string | null {
  const dropped = droppedCount ?? 0;
  if (truncated && dropped > 0) {
    return `流已截断：省略 ${elidedRunes ?? 0} 字，已丢 ${dropped} 条`;
  }
  if (truncated) return `流已截断：省略 ${elidedRunes ?? 0} 字`;
  if (dropped > 0) return `已丢 ${dropped} 条流`;
  return null;
}

export interface MonitorPopoverProps {
  open: boolean;
  env: readonly MonitorEnvRow[];
  agents: readonly MonitorAgentRow[];
  /** Whether the snapshot carried a tasks section AT ALL. Absent = the pump
      was assembled without a roster reader, which lib/panel.ts's own rule says
      is a different fact from an empty roster - the empty line says which one
      it is instead of leaving both looking the same. */
  rosterPresent?: boolean;
  /** The roster pool numbers (wire inFlightSlots/poolCap, Go :162-163), read
      off the roster's own counter - NOT counted from the rows, because a slot
      is reserved before its row exists (Go TaskRosterState's comment). Either
      one undefined = the snapshot carried no tasks key: nothing pool-shaped
      renders at all, and NEVER a "在跑 0 / 上限 0" that would dress an absent
      reading up as data. */
  poolInFlight?: number;
  poolCap?: number;
  /** The whole-log truncation triple (wire streamTruncated/streamElidedRunes/
      droppedCount). All three undefined = no tasks key on the snapshot; a
      clean log (all false/zero) is the normal case and renders no line - the
      warning exists for the day the bound actually bites. */
  streamTruncated?: boolean;
  streamElidedRunes?: number;
  streamDroppedCount?: number;
  /** The two doors the cancel flow needs (readface wave 2, ticket 38 AC#4):
      roster.cancel and the history.page loader. Absent (demo/showcase
      shapes) = rows draw no 停止 control at all, never a dead button. */
  cancelDoors?: TaskCancelDoors;
  onPin?: () => void;
}

/** 子智能体头像：conic 渐变圆片（token 色相，零字面量）。 */
function AgentAvatar() {
  return (
    <span
      aria-hidden="true"
      className="size-4 shrink-0 rounded-full"
      style={{
        background:
          "conic-gradient(var(--accent), var(--green), var(--orange), var(--accent))",
      }}
    />
  );
}

function StateDot({ state }: { state: AgentDot }) {
  return <span aria-hidden="true" className={cn("size-1.5 shrink-0 rounded-full", DOT_CLASS[state])} />;
}

/* ---------------------------------------------------------------------------
   Readface wave 2 (2026-11-17), ticket 38 AC#4: the task-card line and the
   cancel block. Both are exported PURE display components so the render
   instrument feeds real wire shapes through the same functions the popover
   uses - no second rendering path.
   --------------------------------------------------------------------------- */

/** The three task-card keys, one line per row (startedAt / costToDate /
    currentStep). startedAt renders through the SAME formatter the history
    rows use (0/absent stays 未判, never epoch); the cost numbers ride
    verbatim with their currency pair (micros and currency travel together
    per the wire's own pairing rule); currentStep's absence is the wire's
    own "no call in flight" and reads 无在飞, not 未判. */
export function TaskCardLine({ row }: { row: MonitorAgentRow }) {
  const start =
    row.startedAt === undefined ? "未判" : formatHistoryTs(row.startedAt);
  const cost =
    row.costToDate === undefined
      ? "未判"
      : `${row.costToDate.micros} micros ${row.costToDate.currency}（tokens ${row.costToDate.tokensIn}/${row.costToDate.tokensOut}${row.costToDate.cachedTokens === undefined ? "" : `，缓存 ${row.costToDate.cachedTokens}`}）`;
  const step = row.currentStep === undefined || row.currentStep === "" ? "无在飞" : row.currentStep;
  return (
    <p
      className="flex flex-wrap items-center gap-x-2 gap-y-0 px-1.5 pb-1 text-[10.5px] leading-[1.5] text-ink-3"
      aria-label="任务卡三键"
    >
      <span>
        开始 <span className="font-mono tabular-nums">{start}</span>
      </span>
      <span>
        成本 <span className="font-mono tabular-nums">{cost}</span>
      </span>
      <span>
        当前步 <span className="font-mono">{step}</span>
      </span>
    </p>
  );
}

/** The applied-steps reading line (ticket 38 AC#4's three readings). Each
    slot shows ITS value when the scan carries it and 未判 when it does not -
    a count read from an incomplete store would be a fake number, so an
    incomplete scan never renders a number. */
export function AppliedReadingLine({
  reading,
  scanning,
}: {
  reading?: AppliedStepsReading;
  scanning?: boolean;
}) {
  const terminal = reading?.terminalKnown === true ? reading.terminal : "未判";
  const landed =
    reading !== undefined && reading.landed !== null ? `已落地 ${reading.landed} 步` : "已落地步数：未判";
  const aborted =
    reading !== undefined && reading.aborted !== null ? `被中止 ${reading.aborted} 步` : "被中止步数：未判";
  return (
    <p className="px-1.5 pb-1 text-[10.5px] leading-[1.5] text-ink-3" aria-label="取消后读数">
      终态 <span className="font-mono">{terminal}</span> · {landed} · {aborted}
      {scanning === true ? "　（取证行读取中…）" : ""}
    </p>
  );
}

/** One roster.cancel exchange, rendered verbatim: the receipt is the door's
    own sentence (已对任务 <id> 发出停止…), the refusal is the host's own
    words, and a receipt carries the applied-steps reading underneath. */
export function TaskCancelReply({ state }: { state: TaskCancelState }) {
  if (state.phase === "sending") {
    return <p className="px-1.5 pb-1 text-[10.5px] leading-[1.5] text-ink-3">停止请求已在桥上，等门回话…</p>;
  }
  if (state.phase === "no-host") {
    return (
      <p className="px-1.5 pb-1 text-[10.5px] leading-[1.5] text-ink-3">
        没接到本机宿主（这一页没有 wispDispatch 门）：停止请求没有发出去，任务未被触碰。
      </p>
    );
  }
  if (state.phase === "refused") {
    return (
      <div className="mx-1.5 mb-1 rounded-card border border-red/40 bg-red-tint px-2 py-1.5">
        <p className="text-[10.5px] leading-[1.5] text-red">门拒绝了停止，原话逐字：</p>
        <p className="mt-0.5 whitespace-pre-wrap break-all font-mono text-[10px] leading-[1.5] text-ink">{state.raw}</p>
      </div>
    );
  }
  return (
    <div className="mx-1.5 mb-1 rounded-card border border-green/40 bg-green-tint px-2 py-1.5">
      <p className="text-[10.5px] leading-[1.5] text-ink">回执逐字：{state.text}</p>
      <AppliedReadingLine reading={state.reading} scanning={state.scanning} />
      {state.scanNote !== undefined ? (
        <p className="whitespace-pre-wrap break-all text-[10px] leading-[1.5] text-orange">{state.scanNote}</p>
      ) : null}
    </div>
  );
}

export function MonitorPopover({
  open,
  env,
  agents,
  rosterPresent,
  poolInFlight,
  poolCap,
  streamTruncated,
  streamElidedRunes,
  streamDroppedCount,
  cancelDoors,
  onPin,
}: MonitorPopoverProps) {
  // The cancel flow's ephemeral view state: which row's confirm card is up,
  // and each task's own roster.cancel exchange. It dies with the document
  // (the same class as the palette's query); the task's settled truth comes
  // back in the snapshot pushes, never from this state.
  const [confirmId, setConfirmId] = useState<string | null>(null);
  const [cancelStates, setCancelStates] = useState<Record<string, TaskCancelState>>({});
  const agentsRef = useRef(agents);
  agentsRef.current = agents;

  /** The confirm card's confirm button is the ONE place a stop request is
      sent from: the 停止 button only opens the card (先确认再取消，禁止一键杀),
      exactly the two-step shape the privacy purge card established. */
  function confirmStop(taskId: string) {
    if (cancelDoors === undefined) return;
    setConfirmId(null);
    setCancelStates((c) => ({ ...c, [taskId]: { phase: "sending" } }));
    void cancelDoors.cancelTask(taskId).then((reply) => {
      if (reply === null) {
        setCancelStates((c) => ({ ...c, [taskId]: { phase: "no-host" } }));
        return;
      }
      if (reply.startsWith(REFUSED_ENVELOPE_PREFIX)) {
        setCancelStates((c) => ({ ...c, [taskId]: { phase: "refused", raw: reply } }));
        return;
      }
      setCancelStates((c) => ({ ...c, [taskId]: { phase: "receipt", text: reply, scanning: true } }));
      // The task settles on its own goroutine; the receipt does not wait.
      // The applied-steps reading is counted off the history read face once
      // the WHOLE history has been scanned - an incomplete scan stays 未判.
      void scanTaskForensics(cancelDoors.historyLoad, taskId).then((outcome) => {
        setCancelStates((c) => {
          const current = c[taskId];
          if (current === undefined || current.phase !== "receipt") return c;
          if (!outcome.complete) {
            const note = outcome.noHost
              ? "取证读取没接到宿主：读数未判"
              : outcome.refused !== undefined
                ? `取证读取被门拒绝，原话逐字：${outcome.refused}`
                : outcome.payloadError !== undefined
                  ? `取证读取未判：${outcome.payloadError}`
                  : "取证读取未判（history.page 载荷行不在）";
            return { ...c, [taskId]: { ...current, scanning: false, scanNote: note } };
          }
          const live = agentsRef.current.find((a) => a.id === taskId);
          const reading = appliedStepsReading(
            outcome.rows,
            taskId,
            live !== undefined ? { status: live.status, statusKnown: live.statusKnown } : undefined,
            true,
          );
          return { ...c, [taskId]: { ...current, scanning: false, reading } };
        });
      });
    });
  }

  if (!open) return null;
  const boundLine = sectionBoundLine(streamTruncated ?? false, streamElidedRunes, streamDroppedCount);
  return (
    <div
      className="w-[300px] rounded-card border border-line bg-surface p-3 shadow-overlay"
      role="dialog"
      aria-label="任务监控"
      style={{ animation: "pop-in 180ms var(--ease-out-strong) both" }}
    >
      <GlideMenu className="flex flex-col" highlightClassName="inset-x-0 rounded-chip bg-hover">
      <div className="flex items-center justify-between gap-2 pb-1">
        <p className="text-[12.5px] font-semibold text-ink">任务监控</p>
        <button
          aria-label="固定任务监控"
          className="flex size-6 items-center justify-center rounded-control text-ink-3 transition-colors duration-150 hover:bg-hover hover:text-ink"
          onClick={onPin}
          title="固定"
          type="button"
        >
          <Pin aria-hidden="true" size={13} strokeWidth={1.8} />
        </button>
      </div>

      {/* 环境信息：图标/标签 + 右对齐值（ZCode 的 未提交 +10,029 -5,916 行式） */}
      <p className="px-1 pb-1 pt-1.5 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3">
        环境信息
      </p>
      <div className="flex flex-col gap-0.5">
        {env.map((row) => (
          <div className="relative flex items-center gap-2 rounded-chip px-1.5 py-1.5 text-[12px]" data-menu-row key={row.label}>
            <span className="min-w-0 flex-1 truncate text-ink-2">{row.label}</span>
            <span className="shrink-0 font-mono text-[11.5px] tabular-nums text-ink" title={row.value}>
              {row.value}
            </span>
          </div>
        ))}
      </div>

      {/* 子智能体：渐变圆点头像 + 名 + 截断标 + 状态字 + 状态点。零行时把「没有带名册」
          和「名册在、还没有子任务」分开说，绝不共用一句含糊话。池号在组标题右侧：
          读的是名册自己的计数（wire inFlightSlots/poolCap），不从行数现数——
          槽在行之前就被占住，两数偶尔不同是 Go 侧写明的真相，页面不替它挑一个。 */}
      <div className="flex items-center justify-between gap-2 px-1 pb-1 pt-2.5">
        <p className="text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3">
          子智能体
        </p>
        {poolInFlight !== undefined && poolCap !== undefined ? (
          <span
            className="shrink-0 font-mono text-[11px] tabular-nums text-ink-3"
            title="名册池：宿主正占着的槽位与上限"
          >
            在跑 {poolInFlight} / 上限 {poolCap}
          </span>
        ) : null}
      </div>
      <div className="flex flex-col gap-0.5">
        {agents.length === 0 ? (
          <p className="px-1.5 py-1.5 text-[11.5px] leading-relaxed text-ink-3">
            {rosterPresent ? "还没有子任务在跑" : "还没收到名册（宿主没带任务这一节）"}
          </p>
        ) : (
          agents.map((agent) => {
            const display = agentDisplay(agent);
            const bound = rowBoundLabel(agent);
            const cancelState = cancelStates[agent.id];
            return (
              <div key={agent.id} className="flex flex-col">
                <div className="relative flex items-center gap-2 rounded-chip px-1.5 py-1.5 text-[12px]" data-menu-row>
                  <AgentAvatar />
                  <span className="min-w-0 flex-1 truncate text-ink-2" title={agent.name}>
                    {agent.name}
                  </span>
                  {bound ? (
                    <span
                      className="shrink-0 font-mono text-[10.5px] tabular-nums text-warn"
                      title={bound === "流已丢" ? "这一行的流页已被日志丢弃" : "这一行的流被日志折界，页面上少这么多字"}
                    >
                      {bound}
                    </span>
                  ) : null}
                  <span
                    className={cn(
                      "shrink-0 text-[11px] tabular-nums",
                      display.dot === "blocked" ? "text-warn" : "text-ink-3",
                    )}
                  >
                    {display.label}
                  </span>
                  <StateDot state={display.dot} />
                  {cancelDoors !== undefined ? (
                    <button
                      type="button"
                      aria-label={`停止任务 ${agent.id}`}
                      onClick={() => setConfirmId(agent.id)}
                      className="shrink-0 rounded-chip px-1.5 py-0.5 text-[10.5px] text-ink-3 transition-colors duration-150 hover:bg-hover hover:text-red"
                    >
                      停止
                    </button>
                  ) : null}
                </div>
                <TaskCardLine row={agent} />
                {cancelState !== undefined ? <TaskCancelReply state={cancelState} /> : null}
              </div>
            );
          })
        )}
        {boundLine ? (
          <p className="px-1.5 py-1.5 text-[11.5px] leading-relaxed text-warn">{boundLine}</p>
        ) : null}
      </div>

      {/* The stop confirm card (ticket 38 AC#4's 真门): the first click only
          opens it, the request leaves ONLY from its confirm button, and it
          names the row and the irreversibility of the request (the task
          itself settles and lands its terminal rows). */}
      {confirmId !== null ? (
        <div className="mt-2 rounded-card border border-red/40 bg-red-tint px-2.5 py-2" role="alertdialog" aria-label="确认停止任务">
          <p className="text-[11.5px] font-semibold text-red">确认停止任务</p>
          <p className="mt-1 text-[11px] leading-[1.7] text-ink">
            将对任务 <span className="font-mono">{confirmId}</span> 发出停止（roster.cancel，经既有取消通路）。停止不可从面板撤回；
            任务会自行结算并落下终态行，已落地的步骤以 tool_call 取证行与任务终态为准。只停这一行，不级联。
          </p>
          <div className="mt-2 flex items-center gap-2">
            <button
              type="button"
              onClick={() => setConfirmId(null)}
              className="rounded-chip bg-surface px-2 py-0.5 text-[11px] text-ink-2 transition-colors duration-150 hover:bg-hover"
            >
              取消
            </button>
            <button
              type="button"
              aria-label="确认停止"
              onClick={() => confirmStop(confirmId)}
              className="rounded-chip bg-red px-2 py-0.5 text-[11px] text-white transition-colors duration-150 hover:bg-red/90"
            >
              确认停止
            </button>
          </div>
        </div>
      ) : null}
    </GlideMenu>
    </div>
  );
}

/* 演示数据：环境信息取 Wisp 仓现状口径（demo 假数据）；行形状＝wire 的
   TaskRowView 派生形（等批准、D43 名、宿主的未知状态句各占一枚）。 */
export const RB_MONITOR: Pick<MonitorPopoverProps, "env" | "agents"> = {
  env: [
    { label: "工作区", value: "D:\work\workspace\projects plans\Wisp" },
    { label: "分支", value: "dev" },
    { label: "未提交", value: "+2,244 -1,233" },
  ],
  agents: [
    { id: "a1", name: "前端会话看门狗", blockedOnApproval: false, status: "Acting", statusKnown: true },
    { id: "a2", name: "票 145 快照字段", blockedOnApproval: true, status: "Acting", statusKnown: true },
    { id: "a3", name: "票 153 门禁修复", blockedOnApproval: false, status: "Settling", statusKnown: true },
    {
      id: "a4",
      name: "票 146 LiveApprovals",
      blockedOnApproval: false,
      status: "",
      statusKnown: false,
      statusReason: "宿主没有给出这一维的读数（fail-closed：不替任务编一个状态）",
    },
  ],
};
