/* ============================================================================
   Turn rail - the line sidebar for long conversations (2026-09-28).
   ----------------------------------------------------------------------------
   Owner's requirement: 「当一个会话出现多轮对话，然后想查看之前的对话，需要有
   reactbits上的Line Sidebar这种动画效果，这个效果主流harness都有」

   What it is: a hairline down the left edge of the conversation column with one
   mark per turn. Hover a mark and the line brightens to that span; click it and
   the column scrolls to that turn, with the mark growing into a label.

   Why it is a rail and not a list
   -------------------------------
   A turn list duplicates the conversation: you would be reading the same turns
   twice, in two places, and keeping them in sync is a whole feature. A rail
   costs one hairline and never repeats a word. It is also the only shape that
   survives a 60-turn session without pushing the conversation off screen.

   A note on react-bits, because the header of the vendored file claims
   otherwise: the piece vendored as react-bits/animated-list.tsx is a
   scale-and-fade-on-scroll list, not this. There is no Line Sidebar in that
   repo to vendor, so this is built here over the same vocabulary the library
   uses (a measured marker, transitions on transform/opacity only). Claiming it
   were vendored would be a false provenance note in the tree.

   It stays hidden until there is more than one turn: a rail with a single mark
   on it is furniture pretending to be a control.
   ============================================================================ */

import { useEffect, useRef, useState } from "react";
import { cn } from "@/lib/cn";

export interface TurnMark {
  id: string;
  /** One line: what this turn was about, verbatim from the turn. */
  label: string;
}

export function TurnRail({
  turns,
  activeId,
  onJump,
}: {
  turns: readonly TurnMark[];
  activeId: string | null;
  onJump: (id: string) => void;
}) {
  const [hovered, setHovered] = useState<string | null>(null);
  const [pinned, setPinned] = useState(false);
  const listRef = useRef<HTMLDivElement>(null);
  const hideTimer = useRef<number | undefined>(undefined);

  // Unpin once the pointer leaves the rail for good, so the column goes back to
  // being conversation-only. A rail that stays open after you left it is a
  // permanent tax on the reading column.
  useEffect(() => {
    return () => {
      if (hideTimer.current !== undefined) window.clearTimeout(hideTimer.current);
    };
  }, []);

  function show() {
    if (hideTimer.current !== undefined) window.clearTimeout(hideTimer.current);
    setPinned(true);
  }

  function scheduleHide() {
    hideTimer.current = window.setTimeout(() => {
      setPinned(false);
      setHovered(null);
    }, 900);
  }

  if (turns.length < 2) return null;

  const lit = hovered ?? activeId;

  return (
    <div
      className="relative flex w-11 shrink-0 justify-end"
      onMouseEnter={show}
      onMouseLeave={scheduleHide}
    >
      {/* the hairline itself, always there - it is the hint that the rail
          exists even when the rail is collapsed to nothing */}
      <div
        aria-hidden="true"
        className="absolute bottom-3 right-[18px] top-3 w-px bg-line"
      />
      <div
        ref={listRef}
        aria-label="对话轮次"
        className={cn(
          "flex w-full flex-col justify-center gap-1 py-2 pr-1 transition-opacity duration-200",
          pinned ? "opacity-100" : "opacity-35 hover:opacity-100",
        )}
        role="tablist"
      >
        {turns.map((turn) => {
          const on = turn.id === lit;
          return (
            <button
              key={turn.id}
              type="button"
              role="tab"
              aria-selected={turn.id === activeId}
              aria-label={turn.label}
              title={turn.label}
              onClick={() => onJump(turn.id)}
              onMouseEnter={() => setHovered(turn.id)}
              onFocus={() => setHovered(turn.id)}
              className="group relative flex h-6 items-center justify-end pr-[14px]"
            >
              {/* the mark grows along the line rather than appearing beside it,
                  so the rail's width never changes and the column never
                  reflows under the pointer. */}
              <span
                aria-hidden="true"
                className={cn(
                  "absolute right-[15px] rounded-full transition-[width,height,background-color] duration-200",
                  on ? "h-2.5 w-2.5 bg-accent" : "h-1.5 w-1.5 bg-line-strong",
                )}
              />
              {on ? (
                <span
                  className="max-w-24 truncate rounded-chip bg-surface px-1.5 py-0.5 text-[10.5px] text-ink-2 shadow-hairline"
                  style={{ animation: "fade-in 140ms ease both" }}
                >
                  {turn.label}
                </span>
              ) : null}
            </button>
          );
        })}
      </div>
    </div>
  );
}
