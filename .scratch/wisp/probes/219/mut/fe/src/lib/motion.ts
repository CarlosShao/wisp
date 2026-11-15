/* ============================================================================
   Motion tiers - one vocabulary, one attribute, three levels.
   ----------------------------------------------------------------------------
   The panel is allowed to be pretty, and it is also allowed to be cheap. Those
   used to be a binary switch (animate / stop), which is why every animation in
   this tree grew its own inline transition with no way to switch the expensive
   ones off without giving up all of them. Three tiers fix that:

     full   everything runs. The default, and the reason the panel looks like
            something rather than a screenshot of something.
     lite   only the CPU hogs stop. Those are, by definition: animation loops
            (infinite means the compositor never gets to idle) and anything
            animating a filter (blur repaints). One-shot transform/opacity
            entrances and hover transitions stay - they cost a frame or two
            when they fire and nothing at all when they do not.
     off    nothing runs. The floor, and what the OS-level
            prefers-reduced-motion path does regardless of the user's choice.

   The attribute is `data-motion` on <html>, and theme.css is where the rules
   live. Keeping the tier in an attribute rather than in React state is
   deliberate: the CSS has to be able to answer "should this loop run" without
   asking a component tree that may not be mounted yet, and the value survives
   anything that re-renders the root.

   Memory-only, like every other preference on the settings screen: nothing
   leaves the WebView and nothing is persisted, so a restart comes back at the
   default. That is the panel's no-state rule (PLAN.md:1044), not an oversight.
   ============================================================================ */

/** What the settings row offers, in the order it offers them. */
export const MOTION_TIERS = ["全开", "适中", "关闭"] as const;

export type MotionTier = (typeof MOTION_TIERS)[number];

/** The value each tier writes onto <html data-motion>. "full" is written as
    the absence of nothing - it IS the value, because "full" has to be the
    default and a document that forgot to set an attribute should still be
    allowed to animate. */
const TIER_ATTR: Record<MotionTier, "full" | "lite" | "off"> = {
  全开: "full",
  适中: "lite",
  关闭: "off",
};

export const MOTION_ATTR = "data-motion";

/**
 * Apply a tier to the current document. Safe to call before there is a document
 * (the SSR render harnesses), where it does nothing.
 */
export function applyMotionTier(tier: MotionTier): void {
  if (typeof document === "undefined") return;
  document.documentElement.setAttribute(MOTION_ATTR, TIER_ATTR[tier]);
}

/** The tier this document is currently in, read back from the attribute rather
    than from a mirror of it - so a tier set anywhere (including by a test) is
    the tier the UI reports. */
export function readMotionTier(): MotionTier {
  if (typeof document === "undefined") return "全开";
  const v = document.documentElement.getAttribute(MOTION_ATTR);
  if (v === "off") return "关闭";
  if (v === "lite") return "适中";
  return "全开";
}

/** True when JS-driven animation should not run: the reveal timer, the
    counting numbers, the magnet pull. CSS can answer this for the loops it
    owns; these two live in React and have to ask. */
export function motionHeavyAllowed(): boolean {
  return readMotionTier() === "全开";
}
