/* ============================================================================
   Vendored file - ticket 77 AC#7 (ledger: frontend/VENDORED.md).
   ----------------------------------------------------------------------------
   Source repo    : github.com/TurboKach/ai-native-react-components  (registry name "ai-native", beautifului.dev)
   Source file    : components/atoms/Shimmer.tsx
   Component      : Shimmer
   License        : MIT (upstream LICENSE: Copyright (c) 2026 Turbo)
   Upstream commit: 05dab2d2b5f1f3e40029776e339a486d70491079
   Local changes  : provenance header added; the Next.js-only "use client"
                    directive removed; line endings normalized to LF;
                    0 dingbat glyph(s) ASCII-ized for D23/ban #8. (none found in this file)
                    Nothing else - see scripts/vendor.mjs.
   Panel status   : mounted: the waiting label in src/components/PanelSkeleton.tsx
   ============================================================================ */

/* ─────────────────────────────────────────────────────────
 * SHIMMER — a light sweep travelling across a text label.
 *
 * A 200%-wide gradient is clipped to the glyphs and slid
 * from right to left, so the highlight reads as motion
 * through the word rather than a flashing opacity.
 * ───────────────────────────────────────────────────────── */

export function Shimmer({
  children,
  className = "",
}: {
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <span
      className={`bg-clip-text text-transparent ${className}`}
      style={{
        backgroundImage:
          "linear-gradient(90deg, var(--ink-3) 35%, var(--ink) 50%, var(--ink-3) 65%)",
        backgroundSize: "200% 100%",
        animation: "shimmer-text 1.4s linear infinite",
      }}
    >
      {children}
    </span>
  );
}
