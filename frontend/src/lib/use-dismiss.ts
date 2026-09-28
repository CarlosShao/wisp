/* ============================================================================
   useDismiss - close a floating layer when the user acts outside it (2026-09-28).
   ----------------------------------------------------------------------------
   Owner's complaint, verbatim: 「必须要点那个图标，这个悬浮卡片才会消失，不是应该
   点别的地方也消失吗？」 The monitor popover shipped closing ONLY on its own
   trigger, which means a popover you opened covers the page and there is no
   way out of it except the one button that made it. Every harness in the wild
   dismisses on outside pointerdown, on Escape, and on its own trigger.

   Three ways out, because each covers a case the others miss:
     pointerdown outside  - the mouse case, and the one that must win over the
                            trigger's own click (opening then immediately
                            closing would be a dead button).
     keydown Escape       - the keyboard case; a layer you cannot Tab out of is
                            a trap.
     the trigger itself   - handled by the caller toggling, kept out of here so
                            this hook never has to know where the trigger is.

   capture phase on purpose. A pointerdown that lands on something which then
   unmounts (a menu row, a dialog) never reaches a bubble-phase listener, and
   the layer would be left open over a page that moved on.
   ============================================================================ */

import { useEffect, type RefObject } from "react";

export function useDismiss(
  ref: RefObject<HTMLElement | null>,
  open: boolean,
  onDismiss: () => void,
): void {
  useEffect(() => {
    if (!open) return;

    function onPointerDown(event: PointerEvent) {
      const node = ref.current;
      if (!node) return;
      const target = event.target;
      if (target instanceof Node && node.contains(target)) return;
      onDismiss();
    }

    function onKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") onDismiss();
    }

    document.addEventListener("pointerdown", onPointerDown, true);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("pointerdown", onPointerDown, true);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [ref, open, onDismiss]);
}
