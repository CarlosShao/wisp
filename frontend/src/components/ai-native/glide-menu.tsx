/* ============================================================================
   Vendored file - frontend session, 2026-09-27
   ----------------------------------------------------------------------------
   Source repo    : github.com/slev12397/beautiful-ui  (HEAD 44a274e598395ab61e7c96c26fda2758780253b7)
   Source file    : components/primitives/GlideMenu.tsx
   Component      : GlideMenu - one hover layer that glides between menu rows
   License        : MIT (upstream LICENSE: Copyright (c) 2026 Shane Levine)
   Local changes  : provenance header added; the Next.js-only "use client"
                    directive removed; default export changed to a named one to
                    match this tree's import convention (every other vendored
                    atom is imported by name). Behaviour is otherwise upstream
                    verbatim - same 220ms cubic-bezier(0.23,1,0.32,1), same
                    rowSelector hook, same focus capture.

   WHY THIS EXISTS ALONGSIDE glide-highlight.tsx - the two are not duplicates,
   and collapsing them would break one of the two jobs:

     GlideMenu      HOVER layer. It follows the pointer, and it disappears when
                    the pointer leaves. Right for a menu or a list where
                    hovering is the whole interaction.
     GlideHighlight  CONTROLLED selection. It sits on an activeKey the owner
                    decided, so it stays put when the pointer leaves - which is
                    what a selected tab or a keyboard cursor needs. The library
                    has no equivalent (react-bits' GlideSelect is the nearest
                    thing and is a whole select widget, not a layer).

   Using this one where hover is the interaction is the point: the sidebar and
   the tab dropdown were each carrying a hand-written copy of exactly this
   measurement, and the first version of the selected-state work went to
   react-bits for it when this file was sitting unread in the library.
   ============================================================================ */

import { useRef, useState, type ReactNode } from "react";

export interface GlideMenuProps {
  children: ReactNode;
  className?: string;
  highlightClassName?: string;
  rowSelector?: string;
}

export function GlideMenu({
  children,
  className = "",
  highlightClassName = "inset-x-0 rounded-[8px] bg-hover",
  rowSelector = "[data-menu-row]",
}: GlideMenuProps) {
  const ref = useRef<HTMLDivElement>(null);
  const [box, setBox] = useState<{ top: number; height: number } | null>(null);
  const [visible, setVisible] = useState(false);

  const moveTo = (target: EventTarget | null) => {
    const container = ref.current;
    if (!(target instanceof Element) || !container) return;
    const row = target.closest(rowSelector);
    if (!(row instanceof HTMLElement) || !container.contains(row)) return;
    const containerRect = container.getBoundingClientRect();
    const rowRect = row.getBoundingClientRect();
    setBox({ top: rowRect.top - containerRect.top, height: rowRect.height });
    setVisible(true);
  };

  return (
    <div
      className={`group/glide-menu relative ${className}`}
      onBlurCapture={(event) => {
        if (!ref.current?.contains(event.relatedTarget as Node | null)) setVisible(false);
      }}
      onFocusCapture={(event) => moveTo(event.target)}
      onMouseLeave={() => setVisible(false)}
      onMouseOver={(event) => moveTo(event.target)}
      ref={ref}
    >
      <span
        aria-hidden="true"
        className={`m-glide pointer-events-none absolute ${highlightClassName}`}
        style={{
          top: box?.top ?? 0,
          height: box?.height ?? 0,
          opacity: box && visible ? 1 : 0,
          transition:
            "top 220ms cubic-bezier(0.23,1,0.32,1), height 220ms cubic-bezier(0.23,1,0.32,1), opacity 150ms ease",
        }}
      />
      {children}
    </div>
  );
}
