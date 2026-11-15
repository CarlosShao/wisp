/* ============================================================================
   Vendored file - ticket 77 AC#7 (ledger: frontend/VENDORED.md).
   ----------------------------------------------------------------------------
   Source repo    : github.com/TurboKach/ai-native-react-components  (registry name "ai-native", beautifului.dev)
   Source file    : components/atoms/StreamText.tsx
   Component      : StreamText
   License        : MIT (upstream LICENSE: Copyright (c) 2026 Turbo)
   Upstream commit: 05dab2d2b5f1f3e40029776e339a486d70491079
   Local changes  : provenance header added; the Next.js-only "use client"
                    directive removed; line endings normalized to LF;
                    0 dingbat glyph(s) ASCII-ized for D23/ban #8. (none found in this file)
                    Nothing else - see scripts/vendor.mjs.
   Panel status   : NOT mounted: it builds segments with split on a single space, so a Chinese answer reveals as one unit; PLAN.md:3476 逐段 is served by our own src/components/reveal-text.tsx
   ============================================================================ */

import { useEffect, useState } from "react";

/* ─────────────────────────────────────────────────────────
 * STREAM TEXT — word-by-word reveal with a trailing caret.
 *
 * Each word blurs in on arrival rather than popping, which
 * keeps a long line from flickering as it fills. onProgress
 * fires per word so callers can re-anchor anything tracking
 * the text; onDone fires once the last word lands.
 * ───────────────────────────────────────────────────────── */

const WORD_MS = 46;

export function StreamText({
  text,
  onProgress,
  onDone,
}: {
  text: string;
  onProgress?: () => void;
  onDone?: () => void;
}) {
  const words = text.split(" ");
  const [count, setCount] = useState(0);
  const done = count >= words.length;

  useEffect(() => {
    setCount(0);
  }, [text]);

  useEffect(() => {
    if (done) {
      onDone?.();
      return;
    }
    const t = setTimeout(() => setCount((c) => c + 1), WORD_MS);
    return () => clearTimeout(t);
    // onDone is intentionally uncontrolled — callers pass inline closures.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [count, done]);

  useEffect(() => {
    onProgress?.();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [count]);

  return (
    <>
      {words.slice(0, count).map((word, i) => (
        <span
          key={i}
          className="inline [will-change:filter,opacity]"
          style={{
            animation: "stream-in 420ms cubic-bezier(0.22,0.61,0.25,1) both",
          }}
        >
          {word}
          {i < words.length - 1 ? " " : ""}
        </span>
      ))}
      {!done && <span className="stream-caret is-streaming" aria-hidden />}
    </>
  );
}
