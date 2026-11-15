/* ============================================================================
   Vendored file - frontend session, 2026-09-25 (owner: "组件都按照真组件库的来")
   ----------------------------------------------------------------------------
   Source repo    : github.com/slev12397/beautiful-ui  (HEAD 44a274e598395ab61e7c96c26fda2758780253b7)
   Source file    : components/atoms/ValuePill.tsx
   Component      : ValuePill - inline value badge, softer than a StatusPill
                    (no dot) and not a mono token (see Chip)
   License        : MIT (upstream LICENSE: Copyright (c) 2026 Shane Levine)
   Local changes  : provenance header added. Nothing else - the implementation
                    below is upstream verbatim.
   Token note     : the tone table's classes (bg-field / text-ink-2 /
                    bg-green-tint / text-green / bg-orange-tint / text-orange /
                    bg-red-tint / text-red / bg-accent-tint / text-accent-ink)
                    resolve in src/styles/theme.css's @theme inline namespace
                    block. The ring values are color-mix() over var(--green) /
                    var(--orange) / var(--red) / var(--accent) - a colour
                    formula, not a literal, and the hue vars resolve in the
                    fourth-generation table.
   Panel status   : mounted: the privacy cards (section 07) and the three cost
                    stat cards (section 08) in src/components/showcase.tsx
   ============================================================================ */

type Tone = "neutral" | "green" | "orange" | "red" | "accent";

const TONES: Record<Tone, { cls: string; ring: string }> = {
  neutral: { cls: "bg-field text-ink-2", ring: "var(--shadow-hairline)" },
  green: { cls: "bg-green-tint text-green", ring: "0 0 0 1px color-mix(in oklch, var(--green) 28%, transparent)" },
  orange: { cls: "bg-orange-tint text-orange", ring: "0 0 0 1px color-mix(in oklch, var(--orange) 28%, transparent)" },
  red: { cls: "bg-red-tint text-red", ring: "0 0 0 1px color-mix(in oklch, var(--red) 28%, transparent)" },
  accent: { cls: "bg-accent-tint text-accent-ink", ring: "0 0 0 1px color-mix(in oklch, var(--accent) 28%, transparent)" },
};

/** Inline value badge — a plain value (a date, a name, a count) set off in
 *  prose. Softer than a StatusPill (no dot) and not a mono token (see Chip). */
export function ValuePill({
  children,
  tone = "neutral",
  className = "",
}: {
  children: React.ReactNode;
  tone?: Tone;
  className?: string;
}) {
  const t = TONES[tone];
  return (
    <span
      className={`mx-0.5 inline-flex items-center rounded-full px-1.5 py-0
        align-middle text-[12px] font-medium ${t.cls} ${className}`}
      style={{ boxShadow: t.ring }}
    >
      {children}
    </span>
  );
}
