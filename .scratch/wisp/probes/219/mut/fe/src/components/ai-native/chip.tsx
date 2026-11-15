/* ============================================================================
   Vendored file - frontend session, 2026-09-25 (owner: "组件都按照真组件库的来")
   ----------------------------------------------------------------------------
   Source repo    : github.com/slev12397/beautiful-ui  (HEAD 44a274e598395ab61e7c96c26fda2758780253b7)
   Source file    : components/atoms/Chip.tsx
   Component      : Chip - "Monospace token chip, for code values like updated_at"
   License        : MIT (upstream LICENSE: Copyright (c) 2026 Shane Levine)
   Local changes  : provenance header added; upstream doc comment kept verbatim
                    on the function below; nothing else.
   Token note     : the three tone classes (bg-inset / bg-accent-tint /
                    bg-orange-tint, text-ink-2 / text-accent-ink / text-orange)
                    resolve through src/styles/theme.css's library-vocabulary
                    alias block, which maps each one onto a token the generated
                    table derives from design/doubao/demo/styles.css. That alias
                    block is why this file can be mounted at all: before it,
                    16 of the library's utility names existed in no theme of
                    ours, so every component using them rendered unstyled.
   Panel status   : mounted: the config key labels in
                    src/components/config-screen.tsx
   ============================================================================ */

/** Monospace token chip — for code values like `updated_at`. */
export function Chip({
  children,
  tone = "neutral",
  className = "",
}: {
  children: React.ReactNode;
  tone?: "neutral" | "accent" | "orange";
  className?: string;
}) {
  const tones = {
    neutral: "bg-inset text-ink-2",
    accent: "bg-accent-tint text-accent-ink",
    orange: "bg-orange-tint text-orange",
  };
  return (
    <code
      className={`inline rounded-md px-1.5 py-0.5 font-mono text-[12px]
        leading-none align-[-1px] ${tones[tone]} ${className}`}
    >
      {children}
    </code>
  );
}
