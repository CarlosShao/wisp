/* ============================================================================
   MonitorPopover - the task-monitor popover, ZCode shape (owner 2026-09-26
   screenshot: a small card anchored under the header button - title + pin
   slot, an 环境信息 group of icon/label/value rows, a 子智能体 group with
   gradient-avatar agent rows and status dots).
   ----------------------------------------------------------------------------
   Props-driven; the ANCHOR (absolute positioning shell) belongs to the
   parent - this component only paints the card. Zero data of its own.
   ============================================================================ */

import { Pin } from "lucide-react";
import { cn } from "@/lib/cn";

export interface MonitorEnvRow {
  label: string;
  value: string;
}

export interface MonitorAgentRow {
  id: string;
  name: string;
  state: "running" | "done" | "failed";
}

export interface MonitorPopoverProps {
  open: boolean;
  env: readonly MonitorEnvRow[];
  agents: readonly MonitorAgentRow[];
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

function StateDot({ state }: { state: MonitorAgentRow["state"] }) {
  return (
    <span
      aria-hidden="true"
      className={cn(
        "size-1.5 shrink-0 rounded-full",
        state === "running" ? "bg-accent" : state === "done" ? "bg-green" : "bg-red",
      )}
    />
  );
}

export function MonitorPopover({ open, env, agents, onPin }: MonitorPopoverProps) {
  if (!open) return null;
  return (
    <div
      className="w-[300px] rounded-card border border-line bg-surface p-3 shadow-overlay"
      role="dialog"
      aria-label="任务监控"
      style={{ animation: "pop-in 180ms var(--ease-out-strong) both" }}
    >
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
          <div className="flex items-center gap-2 rounded-chip px-1.5 py-1.5 text-[12px] hover:bg-hover" key={row.label}>
            <span className="min-w-0 flex-1 truncate text-ink-2">{row.label}</span>
            <span className="shrink-0 font-mono text-[11.5px] tabular-nums text-ink" title={row.value}>
              {row.value}
            </span>
          </div>
        ))}
      </div>

      {/* 子智能体：渐变圆点头像 + 名 + 状态点 */}
      <p className="px-1 pb-1 pt-2.5 text-[10.5px] font-medium uppercase tracking-[0.08em] text-ink-3">
        子智能体
      </p>
      <div className="flex flex-col gap-0.5">
        {agents.map((agent) => (
          <div className="flex items-center gap-2 rounded-chip px-1.5 py-1.5 text-[12px] hover:bg-hover" key={agent.id}>
            <AgentAvatar />
            <span className="min-w-0 flex-1 truncate text-ink-2" title={agent.name}>
              {agent.name}
            </span>
            <StateDot state={agent.state} />
          </div>
        ))}
      </div>
    </div>
  );
}

/* 演示数据：环境信息取 Wisp 仓现状口径（demo 假数据）。 */
export const RB_MONITOR: Pick<MonitorPopoverProps, "env" | "agents"> = {
  env: [
    { label: "工作区", value: "D:\work\workspace\projects plans\Wisp" },
    { label: "分支", value: "dev" },
    { label: "未提交", value: "+2,244 -1,233" },
  ],
  agents: [
    { id: "a1", name: "前端会话看门狗", state: "running" },
    { id: "a2", name: "票 145 快照字段", state: "running" },
    { id: "a3", name: "票 153 门禁修复", state: "done" },
    { id: "a4", name: "票 146 LiveApprovals", state: "done" },
  ],
};
