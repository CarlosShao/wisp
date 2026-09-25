/* ============================================================================
   Vendored file - frontend session, 2026-09-25 (owner: "组件都按照真组件库的来")
   ----------------------------------------------------------------------------
   Source repo    : github.com/slev12397/beautiful-ui  (HEAD 44a274e598395ab61e7c96c26fda2758780253b7)
   Source file    : components/atoms/Switch.tsx
   Component      : Switch - role="switch" pill with a sliding thumb
   License        : MIT (upstream LICENSE: Copyright (c) 2026 Shane Levine)
   Local changes  : provenance header added; the Next.js-only "use client"
                    directive removed. ONE further change, the only one
                    sanctioned by the tree's no-colour-literal rule: the
                    thumb's drop-shadow was written as a colour literal and is
                    replaced with the shadow-btn token utility (the same word
                    the pressed pills use). Nothing else - the implementation
                    below is upstream verbatim.
   Token note     : bg-ink / bg-line-strong / bg-white resolve via the @theme
                    inline namespace block over the fourth-generation table
                    plus Tailwind's built-in white.
   Panel status   : mounted: the local-retention toggle on the L2 memory card
                    (showcase section 07, src/components/showcase.tsx)
   ============================================================================ */

export function Switch({
  checked,
  onChange,
  label,
}: {
  checked: boolean;
  onChange: (v: boolean) => void;
  label?: string;
}) {
  return (
    <button
      role="switch"
      aria-checked={checked}
      aria-label={label}
      onClick={() => onChange(!checked)}
      className={`relative h-6 w-10 shrink-0 rounded-full transition-colors duration-200
        ${checked ? "bg-ink" : "bg-line-strong"}`}
    >
      <span
        className="absolute top-0.5 left-0.5 size-5 rounded-full bg-white
          shadow-btn transition-transform duration-200"
        style={{
          transform: checked ? "translateX(16px)" : "translateX(0)",
          transitionTimingFunction: "cubic-bezier(0.23, 1, 0.32, 1)",
        }}
      />
    </button>
  );
}
