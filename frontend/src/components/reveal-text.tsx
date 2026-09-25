/* ============================================================================
   RevealText (frontend session, 2026-09-25) - the streaming answer's segment
   reveal, plus the caret PLAN.md:3476 specifies.

   This is OUR component, not the vendored atom. src/components/ai-native/
   stream-text.tsx stays byte-for-byte upstream and is now NOT mounted; forking
   it by hand would be reverted silently by the next `npm run vendor:beautifului`
   (the same trap scripts/vendor.mjs records at its lines 51-52 for the U+2212
   hand-fix in 3b59512). The house pattern for "keep the reference, ship our
   adaptation" is approval-card.tsx vs l2-approval-card.tsx.

   What is deliberately unchanged from the atom, so Latin text looks exactly as
   it did: the 46 ms per-segment cadence, the 420 ms stream-in blur per segment,
   the caret markup, and the same word rhythm for space-separated text.
   ============================================================================ */

import { useEffect, useState } from "react";
import { revealSegments } from "@/lib/reveal";

const SEGMENT_MS = 46;

export function RevealText({ text }: { text: string }) {
  const segments = revealSegments(text);
  const [count, setCount] = useState(0);
  const done = count >= segments.length;

  // No reset when `text` grows: an arriving SSE chunk extends the answer, and
  // rewinding would blank out what the reader has already read. The atom did
  // exactly that (its setCount(0) fires on every text change); whether it ever
  // bites depends on whether ticket 35's pump pushes cumulative text or fresh
  // chunks, which is not decided on disk yet.
  useEffect(() => {
    if (done) return;
    const timer = setTimeout(() => setCount((c) => c + 1), SEGMENT_MS);
    return () => clearTimeout(timer);
  }, [count, done]);

  return (
    <>
      {segments.slice(0, count).map((segment, i) => (
        <span
          className="inline [will-change:filter,opacity]"
          key={i}
          style={{
            animation: "stream-in 420ms cubic-bezier(0.22,0.61,0.25,1) both",
          }}
        >
          {segment}
        </span>
      ))}
      {!done && <span aria-hidden className="stream-caret is-streaming" />}
    </>
  );
}
