/* ============================================================================
   Terminal region (ticket 191 wave 2's UI half, 2026-11-20).
   ----------------------------------------------------------------------------
   The panel-embedded terminal's read face: one collapsible region showing the
   INPUT PORT (the command line the user typed) on top and the OUTPUT STREAM's
   rolling tail below, both verbatim off the snapshot's terminal section
   (internal/panel/terminal_carrier_191.go). This page derives nothing: the
   text on screen is the wire's text, the truncation line's N is the wire's
   droppedRunes.

   The carrier's three facts each render as their own shape, never collapsed
   into one another:
     - key absent  (no Terminal reader on this pump)  -> 终端未连接
     - key present, both rows absent ("terminal": {}) -> 尚无输入 / 尚无输出
       (wired, nothing arrived yet - a different fact from disconnected)
     - rows present                                   -> the tails verbatim,
       each with its own 已截断，另有 N 字符未显示 line when its bound
       evicted anything (a present row with truncated=false carries no
       notice - the window holding everything is information, not absence).

   READ ONLY, AND THE SHAPE PROOF IS THE BUTTON COUNT: the region mounts
   exactly ONE control, the disclosure toggle (a view preference that dies
   with the document, the same class as the theme switch). No send, no run,
   no clear, no bridge verb anywhere - the typed command travels the same
   gated tool leg as every other command (ticket 191 AC#0's 甲 ruling), and a
   region that drew an executable control would be a second execution path
   the gate never judged. render:terminal191 pins the one-button shape and
   the no-bridge-verb shape.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { useEffect, useRef, useState, type RefObject } from "react";
import { ChevronDown, SquareTerminal } from "lucide-react";
import type { TerminalSection, TerminalWindow } from "@/lib/panel";

/** The absent-section sentence. The wording is load-bearing: render:terminal191's
    tooth 1 reads it off the markup - a region that renders blank on a missing
    key (or that paints the empty state and the disconnected state the same)
    goes red there. */
export const TERMINAL_NOT_CONNECTED = "终端未连接";

const NOT_CONNECTED_REASON =
  "快照没有 terminal 键：这台进程的面板装配没有接终端读取器（PumpSources.Terminal）。";

const NO_INPUT = "尚无输入";
const NO_OUTPUT = "尚无输出";

/** The truncation line, verbatim shape 「已截断，另有 N 字符未显示」: N is the
    row's own droppedRunes - the wire's count, never a page-side derivation
    (tooth 3 renders two windows with different counts and reads both back). */
function truncationNotice(row: TerminalWindow): string | null {
  return row.truncated ? `已截断，另有 ${row.droppedRunes} 字符未显示` : null;
}

/** One read window of the region: label, optional truncation line, and either
    the row's tail verbatim or the row-level honest-empty word. The tail never
    re-wraps the producer's own line breaks (whitespace-pre-wrap keeps them)
    and never re-derives the text. */
function TerminalWindowBlock({
  label,
  row,
  emptyWord,
  scroll,
  tailRef,
}: {
  label: string;
  row?: TerminalWindow;
  emptyWord: string;
  /** The output stream is the tail window: bounded height, scrollable, and
      pinned to its newest line by the region's effect. The input port shows
      the line the user is on - most recent wins, no history to scroll. */
  scroll?: boolean;
  tailRef?: RefObject<HTMLDivElement | null>;
}) {
  const notice = row === undefined ? null : truncationNotice(row);
  return (
    <div aria-label={label} className="rounded-card border border-line bg-inset px-3 py-2">
      <p className="text-[10.5px] font-medium text-ink-3">{label}</p>
      {notice !== null ? <p className="mt-0.5 text-[10.5px] text-warn">{notice}</p> : null}
      {row === undefined ? (
        <p className="mt-0.5 text-[11.5px] text-ink-3">{emptyWord}</p>
      ) : (
        <div
          ref={tailRef}
          className={`mt-0.5 whitespace-pre-wrap break-words font-mono text-[11.5px] leading-[1.8] text-ink${
            scroll ? " max-h-[220px] overflow-y-auto" : ""
          }`}
        >
          {row.text}
        </div>
      )}
    </div>
  );
}

/** The panel's terminal region. Key absent (no reader on the pump) renders
    终端未连接 - a blank region would read as "everything is fine and quiet",
    which is exactly the fact nobody sent. The section's own presence with
    both rows absent renders the two honest 尚无 words instead. */
export function TerminalRegion({ terminal }: { terminal?: TerminalSection }) {
  const [open, setOpen] = useState(true);

  // The tail window follows the newest line: without this the push that
  // appends output leaves the view parked on the head - a tail window that
  // shows the head is a lie about its own shape. Re-runs per push and on
  // unfold, never scrolls anything the user is reading upward.
  const tailRef = useRef<HTMLDivElement | null>(null);
  const outputTail = terminal?.output?.text ?? "";
  useEffect(() => {
    const el = tailRef.current;
    if (el) el.scrollTop = el.scrollHeight;
  }, [outputTail, open]);

  return (
    <section aria-label="终端" className="shrink-0 border-t border-line bg-surface">
      <div className="flex items-center gap-2 px-6 pt-2.5">
        {/* The region's ONE control: the disclosure toggle. A view preference
            that dies with the document - not an execution path. */}
        <button
          type="button"
          aria-expanded={open}
          aria-controls="wisp-terminal-body"
          onClick={() => setOpen((v) => !v)}
          className="flex items-center gap-1.5 rounded-control px-1.5 py-1 text-[12px] font-medium text-ink-2 transition-colors duration-150 hover:bg-hover hover:text-ink"
        >
          <SquareTerminal aria-hidden="true" size={13} strokeWidth={1.9} />
          终端
          <ChevronDown
            aria-hidden="true"
            className={`shrink-0 text-ink-3 transition-transform duration-150 ${open ? "rotate-180" : ""}`}
            size={11}
            strokeWidth={2.2}
          />
        </button>
        <span className="text-[10.5px] text-ink-3">只读</span>
      </div>
      {open ? (
        <div id="wisp-terminal-body" className="flex flex-col gap-2 px-6 pb-3 pt-1.5">
          {terminal === undefined ? (
            <div aria-label={TERMINAL_NOT_CONNECTED} className="rounded-card border border-line bg-inset px-3.5 py-3">
              <p className="text-[12.5px] font-medium text-ink">{TERMINAL_NOT_CONNECTED}</p>
              <p className="mt-1 text-[11.5px] leading-[1.7] text-ink-3">{NOT_CONNECTED_REASON}</p>
            </div>
          ) : (
            <>
              <TerminalWindowBlock label="命令行" row={terminal.input} emptyWord={NO_INPUT} />
              <TerminalWindowBlock label="输出" row={terminal.output} emptyWord={NO_OUTPUT} scroll tailRef={tailRef} />
            </>
          )}
        </div>
      ) : null}
    </section>
  );
}
