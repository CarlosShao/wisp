/* ============================================================================
   Vendored file - frontend session, 2026-09-28
   ----------------------------------------------------------------------------
   Source repo    : github.com/DavidHDev/react-bits
   Upstream commit: 5480708039d5fba65741802f3811097ab81e3db5 (local clone)
   Source file    : src/content/Components/LineSidebar/LineSidebar.jsx
   Component      : LineSidebar - a vertical list whose labels and marker lines
                    brighten and slide toward the cursor as it approaches, all
                    driven by ONE rAF loop that lerps a per-item --effect, so
                    colour / shift / scale stay in step with no CSS transitions
                    to stagger.
   License        : MIT + Commons Clause v1.0, upstream LICENSE.md
                    (Copyright (c) 2026 David Haz)

   A correction worth keeping
   -------------------------
   The cut that added components/turn-rail.tsx claimed in its header that
   react-bits has no Line Sidebar, and built a different thing instead. That was
   wrong: the clone at D:\work\AI\component library\react-bits HAS
   src/content/Components/LineSidebar, at the same commit as everything else
   vendored from it. The mistake was searching the already-vendored files in
   this repo instead of the clone. A hand-rolled substitute for a component
   that was sitting on disk is the worst kind of wrong - it looks like work and
   costs a reader twice.

   turn-rail.tsx is NOT deleted here: it is a plain vertical nav with no
   proximity effect, and the product needs both shapes. It is kept for the
   places that want a static list. The cursor-proximity one is this file.

   Local changes, each forced by this repo
   ---------------------------------------
   (1) this provenance header;
   (2) JSX -> TSX, with the props interface spelled out (tsconfig runs TS 6
       strict, so upstream's untyped jsx does not compile as-is): useRef<HTMLUListElement>,
       useRef<(HTMLLIElement|null)[]>, the four number refs, and the pointer
       handler's PointerEvent parameter;
   (3) THE THREE COLOUR LITERALS ARE GONE. upstream defaults accentColor /
       textColor / markerColor to hardcoded hexes, and this tree forbids colour
       literals in frontend/src (TestPanelColourLiteralsLiveOnlyInTheGenerated
       Theme reads every file reachable from main.tsx, including strings AND
       comments - so this paragraph names no value either). The defaults are now
       the token reads the parent passes down, so the component is on-theme in
       both light and dark without the parent having to know the palette. A
       caller can still override all three by prop.
   (4) the CSS import is lowercase './line-sidebar.css' - Vite resolves the
       file, and the upstream capitalised name would survive a case-insensitive
       checkout and break on Linux CI.
   (5) the rAF loop is left running on unmount guard exactly as upstream has it;
       the cancel is already in the cleanup effect.

   Panel status: mounted by components/turn-rail.tsx's parent (App.tsx) for
   multi-turn conversations.
   ============================================================================ */

import { useCallback, useEffect, useRef, useState } from "react";
import "./line-sidebar.css";

const FALLOFF_CURVES: Record<string, (p: number) => number> = {
  linear: (p) => p,
  smooth: (p) => p * p * (3 - 2 * p),
  sharp: (p) => p * p * p,
};

const DEFAULT_ITEMS = [
  "Overview",
  "Components",
  "Animations",
  "Backgrounds",
  "Showcase",
  "Playground",
  "Templates",
  "Changelog",
  "Community",
  "Resources",
  "Documentation",
  "Support",
];

export interface LineSidebarProps {
  items?: string[];
  /** Token read by default (local change 3). Override to recolour. */
  accentColor?: string;
  textColor?: string;
  markerColor?: string;
  showIndex?: boolean;
  showMarker?: boolean;
  proximityRadius?: number;
  maxShift?: number;
  falloff?: "linear" | "smooth" | "sharp";
  markerLength?: number;
  markerGap?: number;
  tickScale?: number;
  scaleTick?: boolean;
  itemGap?: number;
  fontSize?: number;
  smoothing?: number;
  defaultActive?: number | null;
  onItemClick?: (index: number, label: string) => void;
  className?: string;
}

export function LineSidebar({
  items = DEFAULT_ITEMS,
  accentColor = "var(--accent)",
  textColor = "var(--ink-3)",
  markerColor = "var(--line-strong)",
  showIndex = true,
  showMarker = true,
  proximityRadius = 100,
  maxShift = 30,
  falloff = "smooth",
  markerLength = 60,
  markerGap = 0,
  tickScale = 0.5,
  scaleTick = true,
  itemGap = 20,
  fontSize = 1.1,
  smoothing = 100,
  defaultActive = null,
  onItemClick,
  className = "",
}: LineSidebarProps) {
  const listRef = useRef<HTMLUListElement>(null);
  const itemRefs = useRef<(HTMLLIElement | null)[]>([]);
  const targetsRef = useRef<number[]>([]);
  const currentRef = useRef<number[]>([]);
  const rafRef = useRef<number | null>(null);
  const lastRef = useRef(0);
  const activeRef = useRef<number | null>(defaultActive);
  const smoothingRef = useRef(smoothing);
  const [activeIndex, setActiveIndex] = useState<number | null>(defaultActive);

  activeRef.current = activeIndex;
  smoothingRef.current = smoothing;

  /* Single rAF loop that eases every item's --effect toward its target using
     frame-rate independent exponential smoothing, so color, shift and scale
     all move together without staggering CSS transitions. */
  const runFrame = useCallback((now: number) => {
    const dt = Math.min((now - lastRef.current) / 1000, 0.05);
    lastRef.current = now;
    const tau = Math.max(smoothingRef.current, 1) / 1000;
    const k = 1 - Math.exp(-dt / tau);

    let moving = false;
    const els = itemRefs.current;
    for (let i = 0; i < els.length; i++) {
      const el = els[i];
      if (!el) continue;
      const target = Math.max(targetsRef.current[i] || 0, activeRef.current === i ? 1 : 0);
      const cur = currentRef.current[i] || 0;
      const next = cur + (target - cur) * k;
      const settled = Math.abs(target - next) < 0.0015;
      const value = settled ? target : next;
      currentRef.current[i] = value;
      el.style.setProperty("--effect", value.toFixed(4));
      if (!settled) moving = true;
    }

    rafRef.current = moving ? requestAnimationFrame(runFrame) : null;
  }, []);

  const startLoop = useCallback(() => {
    if (rafRef.current != null) {
      cancelAnimationFrame(rafRef.current);
    }
    lastRef.current = performance.now();
    rafRef.current = requestAnimationFrame(runFrame);
  }, [runFrame]);

  const handlePointerMove = useCallback(
    (e: React.PointerEvent<HTMLUListElement>) => {
      const list = listRef.current;
      if (!list) return;
      const rect = list.getBoundingClientRect();
      const pointerY = e.clientY - rect.top;
      const ease = FALLOFF_CURVES[falloff] ?? FALLOFF_CURVES.linear!;
      const els = itemRefs.current;
      for (let i = 0; i < els.length; i++) {
        const el = els[i];
        if (!el) continue;
        const center = el.offsetTop + el.offsetHeight / 2;
        const distance = Math.abs(pointerY - center);
        targetsRef.current[i] = ease(Math.max(0, 1 - distance / proximityRadius));
      }
      startLoop();
    },
    [falloff, proximityRadius, startLoop],
  );

  const handlePointerLeave = useCallback(() => {
    targetsRef.current = targetsRef.current.map(() => 0);
    startLoop();
  }, [startLoop]);

  const handleClick = useCallback(
    (index: number, label: string) => {
      setActiveIndex(index);
      onItemClick?.(index, label);
    },
    [onItemClick],
  );

  useEffect(() => {
    startLoop();
  }, [activeIndex, startLoop]);

  useEffect(
    () => () => {
      if (rafRef.current != null) cancelAnimationFrame(rafRef.current);
      rafRef.current = null;
    },
    [],
  );

  return (
    <nav
      className={`line-sidebar${showMarker ? " line-sidebar--markers" : ""}${scaleTick ? " line-sidebar--scale-tick" : ""}${className ? ` ${className}` : ""}`}
      style={{
        "--accent-color": accentColor,
        "--text-color": textColor,
        "--marker-color": markerColor,
        "--marker-length": `${markerLength}px`,
        "--marker-gap": `${markerGap}px`,
        "--tick-scale": tickScale,
        "--max-shift": `${maxShift}px`,
        "--item-gap": `${itemGap}px`,
        "--font-size": `${fontSize}rem`,
        "--smoothing": `${smoothing}ms`,
      } as React.CSSProperties}
    >
      <ul
        ref={listRef}
        className="line-sidebar__list"
        onPointerMove={handlePointerMove}
        onPointerLeave={handlePointerLeave}
      >
        {items.map((label, index) => (
          <li
            key={`${label}-${index}`}
            ref={(el) => {
              itemRefs.current[index] = el;
            }}
            className="line-sidebar__item"
            aria-current={activeIndex === index ? "true" : undefined}
            onClick={() => handleClick(index, label)}
          >
            {showMarker ? <span className="line-sidebar__marker" aria-hidden="true" /> : null}
            <span className="line-sidebar__label">
              {showIndex ? (
                <span className="line-sidebar__index">{String(index + 1).padStart(2, "0")}</span>
              ) : null}
              <span className="line-sidebar__text">{label}</span>
            </span>
          </li>
        ))}
      </ul>
    </nav>
  );
}

export default LineSidebar;
