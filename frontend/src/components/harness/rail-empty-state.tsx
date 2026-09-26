/* ============================================================================
   RailEmptyState - the right rail's empty state, Qoder shape (owner
   2026-09-26 screenshot: centered big action cards over a faint radial glow,
   each card icon + label + kbd hint).
   ----------------------------------------------------------------------------
   Props-driven: the assembler passes the action list (e.g. 任务运行 / 审批 /
   画布 tabs presented as cards when nothing is pinned yet) and the click
   callback. Zero data, zero literals.
   ============================================================================ */

import type { LucideIcon } from "lucide-react";

export interface RailEmptyAction {
  key: string;
  label: string;
  kbd?: string;
  icon: LucideIcon;
}

export interface RailEmptyStateProps {
  actions: readonly RailEmptyAction[];
  onAction: (key: string) => void;
}

export function RailEmptyState({ actions, onAction }: RailEmptyStateProps) {
  return (
    <div
      className="flex min-h-0 flex-1 flex-col items-center justify-center gap-3 px-4 py-6"
      style={{
        backgroundImage:
          "radial-gradient(circle at 50% 60%, var(--accent-tint) 0%, transparent 62%)",
      }}
    >
      <p className="text-[13px] font-semibold text-ink">打开面板</p>
      <p className="mb-1 text-[11.5px] text-ink-3">选择要在侧边面板中查看的内容。</p>
      <div className="flex w-full flex-col gap-2.5">
        {actions.map((action) => (
          <button
            className="flex w-full items-center gap-3 rounded-card border border-line bg-surface p-4 text-left shadow-card transition-colors duration-150 hover:border-accent/40"
            key={action.key}
            onClick={() => onAction(action.key)}
            type="button"
          >
            <span className="flex size-8 shrink-0 items-center justify-center rounded-control bg-field text-ink-2">
              <action.icon aria-hidden="true" size={15} strokeWidth={1.8} />
            </span>
            <span className="min-w-0 flex-1">
              <span className="block truncate text-[12.5px] font-medium text-ink">{action.label}</span>
              {action.kbd ? (
                <span className="block text-[10.5px] text-ink-3">{action.kbd}</span>
              ) : null}
            </span>
          </button>
        ))}
      </div>
    </div>
  );
}
