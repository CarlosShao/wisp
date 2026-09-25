/* ============================================================================
   Result stream (ticket 77) - assistant text as the Go pump pushes it.

   The reveal used to be the vendored Beautiful UI StreamText atom. It is now
   src/components/reveal-text.tsx: PLAN.md:3476 wants 逐段, and the atom built
   its segments with split(" "), which has nothing to split on in a Chinese
   answer. Ticket 35 owns the pump that feeds `chunks`.
   ============================================================================ */

import { RevealText } from "@/components/reveal-text";
import type { ResultChunkView } from "@/lib/panel";

export function ResultStream({ chunks }: { chunks: ResultChunkView[] }) {
  if (chunks.length === 0) {
    return <p className="text-[12px] text-ink-3">尚无结果。</p>;
  }
  return (
    <ul className="flex flex-col gap-2">
      {chunks.map((chunk) => (
        <li
          className="glass-inset rounded-control px-3 py-2 text-[12.5px] leading-[var(--lh-body)] text-ink"
          key={chunk.correlationId}
        >
          {chunk.done ? chunk.text : <RevealText text={chunk.text} />}
        </li>
      ))}
    </ul>
  );
}
