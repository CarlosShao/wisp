/* ============================================================================
   Vendored file - frontend session, 2026-09-25 (owner: "组件都按照真组件库的来")
   ----------------------------------------------------------------------------
   Source repo    : github.com/slev12397/beautiful-ui  (HEAD 44a274e598395ab61e7c96c26fda2758780253b7)
   Source file    : components/atoms/TextRow.tsx
   Component      : TextRow - label-left / value-right row ("the workhorse of
                    every card in the refs")
   License        : MIT (upstream LICENSE: Copyright (c) 2026 Shane Levine)
   Local changes  : provenance header added. Nothing else - the implementation
                    below is upstream verbatim (the file carries no "use
                    client" upstream, so none had to be removed).
   Token note     : text-ink / text-ink-2 / text-ink-3 resolve in
                    src/styles/theme.css's @theme inline namespace block over
                    the fourth-generation table. No colour literal entered
                    with the file.
   Panel status   : mounted: the L1 profile / L2 memory privacy cards in
                    src/components/showcase.tsx (section 07)
   ============================================================================ */

/** Label-left / value-right row — the workhorse of every card in the refs. */
export function TextRow({
  label,
  value,
  meta,
  className = "",
}: {
  label: React.ReactNode;
  value: React.ReactNode;
  meta?: React.ReactNode;
  className?: string;
}) {
  return (
    <div
      className={`flex min-h-11 items-center justify-between gap-4 py-2 ${className}`}
    >
      <span className="text-sm text-ink-2">{label}</span>
      <span className="flex items-baseline gap-2 text-right">
        <span className="text-sm font-medium text-ink tabular-nums">
          {value}
        </span>
        {meta && <span className="text-[13px] text-ink-3">{meta}</span>}
      </span>
    </div>
  );
}
