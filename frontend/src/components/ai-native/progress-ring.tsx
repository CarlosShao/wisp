/* ============================================================================
   Vendored file - frontend session, 2026-09-25 (owner: "组件都按照真组件库的来")
   ----------------------------------------------------------------------------
   Source repo    : github.com/slev12397/beautiful-ui  (HEAD 44a274e598395ab61e7c96c26fda2758780253b7)
   Source file    : components/atoms/ProgressRing.tsx
   Component      : ProgressRing - small progress ring with content in the
                    center (the task-badge from the library's refs)
   License        : MIT (upstream LICENSE: Copyright (c) 2026 Shane Levine)
   Local changes  : provenance header added. Nothing else - the implementation
                    below is upstream verbatim.
   Token note     : the tone table reads var(--orange) / var(--green) /
                    var(--red) / var(--accent) for the arc and var(--line) for
                    the track; all five resolve in the generated table
                    (src/styles/tokens.generated.css, fourth generation adopts
                    this library's theme verbatim). No colour literal entered
                    with the file.
   Panel status   : mounted: the live-precision ring on the one live control in
                    src/components/config-screen.tsx, plus the ball-state demos
                    in src/components/showcase.tsx (section 05)
   ============================================================================ */

/** Small progress ring with content in the center — the task-badge from the refs. */
export function ProgressRing({
  progress,
  tone = "orange",
  children,
  size = 28,
}: {
  progress: number; // 0..1
  tone?: "orange" | "green" | "red" | "accent";
  children?: React.ReactNode;
  size?: number;
}) {
  const stroke = 2;
  const r = (size - stroke) / 2;
  const c = 2 * Math.PI * r;
  const tones = {
    orange: "var(--orange)",
    green: "var(--green)",
    red: "var(--red)",
    accent: "var(--accent)",
  };
  return (
    <span
      className="relative inline-flex items-center justify-center"
      style={{ width: size, height: size }}
    >
      <svg width={size} height={size} className="-rotate-90 absolute inset-0">
        <circle
          cx={size / 2}
          cy={size / 2}
          r={r}
          fill="none"
          stroke="var(--line)"
          strokeWidth={stroke}
        />
        <circle
          cx={size / 2}
          cy={size / 2}
          r={r}
          fill="none"
          stroke={tones[tone]}
          strokeWidth={stroke}
          strokeLinecap="round"
          strokeDasharray={c}
          strokeDashoffset={c * (1 - Math.min(1, Math.max(0, progress)))}
          style={{ transition: "stroke-dashoffset 400ms cubic-bezier(0.23,1,0.32,1)" }}
        />
      </svg>
      <span className="relative text-[12px] font-semibold tabular-nums">
        {children}
      </span>
    </span>
  );
}
