/* ============================================================================
   Ball-state screen (showcase component, 2026-09-25)
   ----------------------------------------------------------------------------
   The ball's state machine as one readable table, plus the ProgressRing atom
   wearing two real values. The rows are the state machine's own words -
   PLAN.md section 2's state/visual table, the same list D43 freezes (20
   states, 40 transitions; C12 freezes state + transition table + timeouts
   together, so this table shows the STATE side and names where the
   transition table lives - it does not paraphrase it).

   Built as an independent props-driven component, not inline showcase JSX,
   so the future Go-fed 球状态 screen mounts the same shape with snapshot
   data instead of fixtures.
   ============================================================================ */

import { ProgressRing } from "@/components/ai-native/progress-ring";

export interface BallStateRow {
  /** The state's identifier, as the state machine spells it. */
  state: string;
  /** The visual/motion the frozen row gives this state. */
  visual: string;
  /** What brings the ball into this state. */
  enter: string;
  /** What takes it out, including the timeout. */
  exit: string;
}

export interface BallRingDemo {
  label: string;
  progress: number;
  tone: "orange" | "green" | "red" | "accent";
}

export function BallScreen({
  states,
  rings = [],
}: {
  states: readonly BallStateRow[];
  rings?: readonly BallRingDemo[];
}) {
  return (
    <div className="flex w-full flex-col gap-3">
      <div className="overflow-hidden rounded-card border border-line bg-surface shadow-card">
        <div className="grid grid-cols-[minmax(0,1.1fr)_minmax(0,1.5fr)_minmax(0,1.2fr)_minmax(0,1.4fr)] gap-2 border-b border-line px-3 py-2 text-[11.5px] font-medium text-ink-3">
          <span>态</span>
          <span>视觉动效</span>
          <span>进入</span>
          <span>退出/超时</span>
        </div>
        {states.map((row, i) => (
          <div
            className="grid grid-cols-[minmax(0,1.1fr)_minmax(0,1.5fr)_minmax(0,1.2fr)_minmax(0,1.4fr)] items-start gap-2 border-b border-line px-3 py-2 text-[12px] leading-[1.6] transition-colors duration-100 last:border-b-0 hover:bg-hover"
            key={row.state}
            style={{ animation: `fade-up 450ms cubic-bezier(0.23,1,0.32,1) ${Math.min(i, 10) * 40}ms both` }}
          >
            <span className="font-mono text-[11.5px] text-ink">{row.state}</span>
            <span className="text-ink-2">{row.visual}</span>
            <span className="text-ink-2">{row.enter}</span>
            <span className="text-ink-3">{row.exit}</span>
          </div>
        ))}
      </div>

      {rings.length > 0 && (
        <div className="flex items-center gap-5 rounded-card border border-line bg-surface px-3 py-2.5 shadow-card">
          {rings.map((ring) => (
            <span className="flex items-center gap-2" key={ring.label}>
              <ProgressRing progress={ring.progress} tone={ring.tone}>
                {Math.round(ring.progress * 100)}
              </ProgressRing>
              <span className="text-[12px] text-ink-2">{ring.label}</span>
            </span>
          ))}
        </div>
      )}

      <p className="text-[11px] leading-[1.6] text-ink-3">
        C12 冻结的是「态 + 转移表 + 超时」三者：转移表见 PLAN.md D43（20 态 · 40 条转移，未列出的转移一律非法），本表只照抄态与视觉两列。
      </p>
    </div>
  );
}
