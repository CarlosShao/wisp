/* ============================================================================
   Result stream (ticket 77; demo look pass 2026-09-25)
   ----------------------------------------------------------------------------
   Assistant text as the Go pump pushes it, restyled onto the demo chat
   rhythm: each chunk is a plain body-size block (demo chat.js draws answers as
   text-sm/leading-relaxed on the page foreground, not as inset cards), and a
   newly mounted block fades up with the demo's animated-list-item so the
   stream reads as arrival, not as appends to a spreadsheet.

   The reveal itself is unchanged: PLAN.md:3476's 逐段 caret lives in
   src/components/reveal-text.tsx, and scripts/render-stream.tsx pins the whole
   contract (no vendored split-on-space atom, done text verbatim, no caret on
   finished blocks). Ticket 35 owns the pump that feeds `chunks`.
   ============================================================================ */

import { RevealText } from "@/components/reveal-text";
import type { ResultChunkView } from "@/lib/panel";

export function ResultStream({ chunks }: { chunks: ResultChunkView[] }) {
  if (chunks.length === 0) {
    return <p className="text-[12px] text-muted-foreground">尚无结果。</p>;
  }
  return (
    <ul className="m-0 flex min-w-0 list-none flex-col gap-4 p-0">
      {chunks.map((chunk) => (
        <li
          className="animated-list-item min-w-0 text-[13px] leading-relaxed text-foreground"
          key={chunk.correlationId}
        >
          {chunk.done ? chunk.text : <RevealText text={chunk.text} />}
        </li>
      ))}
    </ul>
  );
}
