/* ============================================================================
   Context + usage meter (product panel, 2026-09-28).
   ----------------------------------------------------------------------------
   The line that sits just above the composer, in the shape the owner asked
   for: tokens, cache hit rate, context occupancy - the three numbers you
   actually watch while an agent works.

   Why it renders an ABSENCE instead of zeros
   ------------------------------------------
   PanelSnapshot has no field for any of the three yet, and the owner's P9 red
   line ("不得造假数据当真实字段") is exactly about this. A row reading
   `0 tok . 0% . 0%` is a lie in a worse way than no row: it says the agent has
   used nothing, which is a claim about the world, not about the panel's
   knowledge. So with no numbers the row says it has not been told, and says
   which ticket would fill it.

   The three-part readout is written here rather than assembled by the caller
   so the SHAPE is pinned in one place: when Go grows the fields this file
   fills in and nothing else has to move. A dash is used for a number Go has
   not sent, never a zero.
   ============================================================================ */

import { cn } from "@/lib/cn";

export interface ContextUsage {
  /** Cumulative tokens on this session, as the host counts them. */
  tokens?: number;
  /** Fraction of prompt tokens served from cache, 0..1. */
  cacheHitRate?: number;
  /** Fraction of the context window occupied, 0..1. */
  occupancy?: number;
}

/** `75.9M tok` - the dsh bottom-bar shape: a short count, not a digit
    avalanche. Below 1000 stays exact because below that nobody wants a unit. */
function shortTokens(n: number): string {
  if (n < 1000) return String(n);
  if (n < 1_000_000) return `${(n / 1000).toFixed(n < 10_000 ? 1 : 0)}k`;
  return `${(n / 1_000_000).toFixed(1)}M`;
}

function percent(fraction: number): string {
  return `${Math.round(fraction * 100)}%`;
}

/** The dot between parts. Muted so it reads as punctuation, not content. */
function Dot() {
  return (
    <span aria-hidden="true" className="shrink-0 text-ink-3">
      ·
    </span>
  );
}

export function ContextMeter({
  usage,
  className,
}: {
  usage: ContextUsage;
  className?: string;
}) {
  const any = usage.tokens !== undefined || usage.cacheHitRate !== undefined || usage.occupancy !== undefined;

  if (!any) {
    return (
      <p className={cn("truncate text-[10.5px] text-ink-3", className)}>
        上下文与用量还没有收到数据，面板不自己造这个数字。
      </p>
    );
  }

  return (
    <div
      aria-label="上下文与用量"
      className={cn("flex items-center gap-1.5 text-[10.5px] tabular-nums text-ink-3", className)}
    >
      {usage.tokens !== undefined ? (
        <span className="shrink-0">
          <span className="text-ink-2">{shortTokens(usage.tokens)}</span> tok
        </span>
      ) : null}
      {usage.tokens !== undefined && usage.cacheHitRate !== undefined ? <Dot /> : null}
      {usage.cacheHitRate !== undefined ? (
        <span className="shrink-0">缓存命中 {percent(usage.cacheHitRate)}</span>
      ) : null}
      {usage.occupancy !== undefined ? (
        <>
          {usage.tokens !== undefined || usage.cacheHitRate !== undefined ? <Dot /> : null}
          <span className="shrink-0">占用 {percent(usage.occupancy)}</span>
        </>
      ) : null}
    </div>
  );
}
