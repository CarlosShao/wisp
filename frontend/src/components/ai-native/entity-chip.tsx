/* ============================================================================
   Vendored file - frontend session, 2026-09-25 (owner: "组件都按照真组件库的来")
   ----------------------------------------------------------------------------
   Source repo    : github.com/slev12397/beautiful-ui  (HEAD 44a274e598395ab61e7c96c26fda2758780253b7)
   Source file    : components/atoms/EntityChip.tsx
   Component      : Monogram + EntityChip - a colored monogram disc and the
                    inline entity reference built on it
   License        : MIT (upstream LICENSE: Copyright (c) 2026 Shane Levine)
   Local changes  : provenance header added. ONE further change, the only one
                    sanctioned by the tree's no-colour-literal rule: upstream
                    gave Monogram's `color` prop the hex default and left
                    EntityChip's `color` optional, so rendering EntityChip with
                    no color silently painted the upstream default. The default
                    is deleted and BOTH color props are now required - a caller
                    must name a colour (our callers pass var(...) tokens), so
                    no upstream hue can enter the tree through a forgotten
                    prop. Nothing else - the implementation below is otherwise
                    upstream verbatim.
   Token note     : bg-field / text-ink / shadow-hairline resolve in
                    src/styles/theme.css's @theme inline namespace block. The
                    disc colour is whatever the caller passes; this repo's
                    callers pass var(--accent) / var(--orange) / var(--red)
                    style tokens, never a literal.
   Panel status   : mounted: the actor marks in the authorization audit table
                    (showcase section 06, src/components/showcase.tsx)
   ============================================================================ */

/** Monogram mark — a colored disc with an initial or short glyph.
 *  The shared building block for entity chips and monogram headings. */
export function Monogram({
  children,
  color,
  className = "",
}: {
  children: React.ReactNode;
  color: string;
  className?: string;
}) {
  return (
    <span
      className={`flex size-4 shrink-0 items-center justify-center rounded-full
        text-[9px] font-semibold leading-none text-white ${className}`}
      style={{ background: color }}
    >
      {children}
    </span>
  );
}

/** Inline entity reference — a monogram + name in a soft field pill.
 *  Names a supplier, person, or record inside running text. Softer than a
 *  StatusPill (no dot, no state) and not a mono token (see Chip). */
export function EntityChip({
  name,
  color,
  monogram,
  className = "",
}: {
  name: string;
  color: string;
  monogram?: React.ReactNode;
  className?: string;
}) {
  return (
    <span
      className={`mx-0.5 inline-flex items-center gap-1 rounded-full bg-field
        py-px pl-[3px] pr-1.5 align-middle shadow-hairline ${className}`}
    >
      <Monogram color={color}>{monogram ?? name.charAt(0)}</Monogram>
      <span className="text-[12px] font-medium text-ink">{name}</span>
    </span>
  );
}
