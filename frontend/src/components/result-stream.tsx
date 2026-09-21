/* ============================================================================
   Result stream (ticket 77) - the one place the vendored Beautiful UI
   StreamText atom is mounted: assistant text reveals word by word with the
   upstream caret, driven by the Go push (ticket 35 owns that pump).
   ============================================================================ */

import { StreamText } from "@/components/ai-native/stream-text";
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
          {chunk.done ? chunk.text : <StreamText text={chunk.text} />}
        </li>
      ))}
    </ul>
  );
}
