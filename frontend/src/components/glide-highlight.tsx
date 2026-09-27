/* ============================================================================
   GlideHighlight - the sliding selection block, as a reusable primitive.
   ----------------------------------------------------------------------------
   Taken from what the harness sidebar was already doing correctly by hand (a
   measured absolutely-positioned block that slides to the hovered or active
   row) and from react-bits' Micro/GlideSelect, which is the same idea done
   properly. Four things make it feel right:

     1. The block is a SIBLING of the rows, absolutely positioned, and it
        MOVES (translate-free top/height transitions) rather than each row
        painting its own highlight. Ten rows each swapping a background repaints
        ten times; one block sliding composites once. That is the whole
        difference between a highlight and a cheap one.

     2. The first placement does not animate. Sliding in from nowhere is a lie
        about where the selection came from; only later moves glide.
        GlideSelect calls this its `instant` flag and it is the single biggest
        thing that makes the motion read as deliberate rather than floaty.

     3. Nothing is measured until the rows exist, and the measure runs in
        useLayoutEffect so the block is placed before the browser paints. An
        effect that ran after paint would show a frame of unhighlighted rows on
        every change.

     4. When the highlighted row goes away (a filter dropped it, a group folded)
        the block hides rather than staying parked over whatever row now
        happens to sit at that offset.

   THE API IS THE ROW ATTRIBUTE, NOT A REF MAP. The first version took
   containerRef + a rowRefs record + a callback ref per row, and it silently
   did nothing because one of the three was wired a moment late. A block that
   finds its own target cannot be half-wired: the only thing a caller owes is
   `data-glide-key` on the row. Anything more would be plumbing for a bug.

   The motion tiers reach it through CSS: the block transitions, so
   data-motion="off" pins it, and "lite" lets it glide (it only moves when the
   selection moves - not a loop, and nothing at all while it is idle).
   ============================================================================ */

import { useLayoutEffect, useRef, useState } from "react";
import { cn } from "@/lib/cn";

/** The block's box inside its container. All four edges are measured, not just
    the two a vertical list needs: a horizontal tab strip has to slide in X, and
    a primitive that only knows about Y silently does nothing there - which is
    exactly what the first version did. */
interface Rect {
  top: number;
  height: number;
  left: number;
  width: number;
}

/** Sub-pixel tolerance. Layout rects land on fractions, so an exact compare
    would re-render on every frame of a scroll for no visible change. */
const EPS = 0.5;

function sameRect(a: Rect, b: Rect): boolean {
  return (
    Math.abs(a.top - b.top) < EPS &&
    Math.abs(a.height - b.height) < EPS &&
    Math.abs(a.left - b.left) < EPS &&
    Math.abs(a.width - b.width) < EPS
  );
}

export interface GlideHighlightProps {
  /** The data-glide-key of the row that should be highlighted, or null. */
  activeKey: string | null;
  /**
   * Change this whenever the rows changed shape (a filter ran, a group folded,
   * a tab opened) so the block re-measures. The open-tab list joined into a
   * string is the usual value.
   */
  layoutKey?: unknown;
  /** Slide time. 220ms is where this tree's sidebar already sits: long enough
      to read as movement, short enough to still feel like a click. */
  duration?: number;
  /** Corner radius, so the block matches whatever it is highlighting. */
  radius?: string;
  /** Background class for the block; defaults to the tree's hover token. */
  blockClassName?: string;
  className?: string;
}

export function GlideHighlight({
  activeKey,
  layoutKey,
  duration = 220,
  radius = "7px",
  blockClassName,
  className,
}: GlideHighlightProps) {
  const blockRef = useRef<HTMLSpanElement>(null);
  const first = useRef(true);
  const [box, setBox] = useState<Rect | null>(null);

  useLayoutEffect(() => {
    const block = blockRef.current;
    const parent = block?.parentElement;
    if (!parent || activeKey === null) {
      setBox((current) => (current === null ? current : null));
      return;
    }
    // Own parent, own lookup. CSS.escape because a key is caller-supplied and
    // a stray quote in it would otherwise throw inside a layout effect.
    const target = parent.querySelector<HTMLElement>(
      `[data-glide-key="${CSS.escape(activeKey)}"]`,
    );
    if (!target) {
      setBox((current) => (current === null ? current : null));
      return;
    }
    const p = parent.getBoundingClientRect();
    const t = target.getBoundingClientRect();
    const next: Rect = {
      top: t.top - p.top,
      height: t.height,
      left: t.left - p.left,
      width: t.width,
    };
    setBox((current) =>
      current !== null && sameRect(current, next) ? current : next,
    );
  }, [activeKey, layoutKey]);

  // One flag, cleared after the first real placement. It is read during render,
  // so the render that shows the moved block is the one that turns the
  // transition back on - the first placement stays transitionless, every later
  // one glides.
  const instant = first.current && box === null;

  useLayoutEffect(() => {
    if (box !== null) first.current = false;
  }, [box]);

  return (
    <span
      aria-hidden="true"
      className={cn(
        "m-glide pointer-events-none absolute top-0 left-0 z-0 bg-hover",
        className,
        blockClassName,
      )}
      ref={blockRef}
      style={{
        top: box?.top ?? 0,
        height: box?.height ?? 0,
        left: box?.left ?? 0,
        width: box?.width ?? 0,
        opacity: box ? 1 : 0,
        borderRadius: radius,
        // No transition on the very first placement (see `instant`).
        transition: instant
          ? "none"
          : `top ${duration}ms var(--ease-out-strong), left ${duration}ms var(--ease-out-strong),` +
            ` height ${duration}ms var(--ease-out-strong), width ${duration}ms var(--ease-out-strong),` +
            ` opacity 150ms ease`,
      }}
    />
  );
}
