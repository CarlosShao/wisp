/* ============================================================================
   Result stream (ticket 77; ai-native look pass 2026-09-25)
   ----------------------------------------------------------------------------
   Assistant text as the Go pump pushes it, restyled onto the library chat
   rhythm: each chunk is a plain body block (text-[13px] leading-relaxed on the
   ink foreground, the ai-native chat.tsx answer shape), and a newly mounted
   block fades up with the library's fade-up keyframe so the stream reads as
   arrival, not as appends to a spreadsheet.

   The ruler this file must never break (scripts/render-stream.tsx pins it,
   and it greps this very file, so the atom's import path is not even written
   in a comment here):
     - no import of the vendored split-on-space atom (the stream-text.tsx
       upstream copy);
     - a finished chunk's text renders VERBATIM as one text node - no reveal
       spans, no caret (the "is-streaming" string must not appear);
     - an unfinished chunk reveals through src/components/reveal-text.tsx,
       whose caret is the frozen PLAN.md:3476 breathing bar.
   Ticket 35 owns the pump that feeds `chunks`.
   ============================================================================ */

import { RevealText } from "@/components/reveal-text";
import type { ResultChunkView } from "@/lib/panel";

export function ResultStream({ chunks }: { chunks: ResultChunkView[] }) {
  if (chunks.length === 0) {
    return <p className="text-[12px] text-ink-3">尚无结果。</p>;
  }
  return (
    <ul className="m-0 flex min-w-0 list-none flex-col gap-5 p-0">
      {chunks.map((chunk) => (
        <li
          className="min-w-0 text-[13px] leading-relaxed text-ink"
          key={chunk.correlationId}
          style={{ animation: "fade-up 400ms var(--ease-out-strong) both" }}
        >
          {chunk.done ? chunk.text : <RevealText text={chunk.text} />}
        </li>
      ))}
    </ul>
  );
}
