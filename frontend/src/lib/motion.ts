/* ============================================================================
   Motion tiers - one vocabulary, one attribute, three levels.
   ----------------------------------------------------------------------------
   Ticket 193 (owner 2026-09-28, ledger A373): the animation tier is a USER
   OPTION backed by config, not a page-local preference - "动画会有一个选项
   让用户选，切换最低就是关闭动画了，那就不会卡了". The vocabulary is the Go
   side's own three words (internal/config/validate.go's
   panelAnimationVocabulary, copied verbatim; the render:animation193
   instrument re-reads that Go line so the two cannot drift):

     full   everything runs.
     reduce only the CPU hogs stop - the long-standing boundary kept as is:
            animation loops are compressed to a single 1ms iteration
            (nothing loops), while one-shot entrances and hover transitions
            stay. The registered rAF loop factory declines to schedule here
            too: a loop is the one shape that never ends on its own.
     off    nothing runs - and "runs unseen" does NOT count as stopped: the
            theme.css off rules attach NO animation and NO transition (so
            getAnimations() stays empty by construction), the flip to off
            cancels what is already in flight, and the loop factory declines
            every start.

   The attribute is `data-motion` on <html>; theme.css is where the rules
   live. App applies the tier straight off the snapshot's panelAnimation
   section on every push (the same seam applyPanelFontSize rides); an absent
   section, an empty tier or an unknown word withdraws the attribute -
   unknown is UNKNOWN, never painted as off or full ("不发键不等于 off",
   the Go carrier's own rule, mirrored here).

   The page-local demo row this file used to serve (the memory-only
   data-motion switch on the settings screen) was retired by the same
   ticket: two switches for one thing was a second judge, and the
   config-backed row (panel.animation, riding config.set) is the real one.
   ============================================================================ */

/** The tier words the wire carries (validate.go's panelAnimationVocabulary,
    verbatim - the instrument pins the equality against the Go source). */
export const ANIMATION_TIERS = ["full", "reduce", "off"] as const;

export type AnimationTier = (typeof ANIMATION_TIERS)[number];

/** The attribute the CSS rules key on (theme.css's motion-tier block). */
export const MOTION_ATTR = "data-motion";

/** The labels the settings row offers, one per wire word. The reduce option
    is the owner-pending default placeholder (schema default "reduce"; ticket
    193: 默认值由 owner 定) - the row marks it （默认，待定） instead of
    treating the placeholder as final. */
export const ANIMATION_TIER_LABELS: Record<AnimationTier, string> = {
  full: "全开",
  reduce: "减少",
  off: "关闭",
};

function isAnimationTier(v: unknown): v is AnimationTier {
  return v === "full" || v === "reduce" || v === "off";
}

/** Apply one tier to the current document. Safe to call before there is a
    document (the SSR render harnesses), where it does nothing. */
export function applyMotionTier(tier: AnimationTier): void {
  if (typeof document === "undefined") return;
  document.documentElement.setAttribute(MOTION_ATTR, tier);
}

/** The tier this document is currently in, or null when none is applied. An
    absent attribute is UNKNOWN - the wire's no-key shape - and comes back as
    null, never as one of the three words. */
export function readMotionTier(): AnimationTier | null {
  if (typeof document === "undefined") return null;
  const v = document.documentElement.getAttribute(MOTION_ATTR);
  return isAnimationTier(v) ? v : null;
}

/** Whether heavy JS-driven animation work (rAF loops) may run on the given
    document: only under full. reduce keeps the no-loop boundary (loops stop;
    entrances and hover transitions stay); off stops everything. A document
    without the attribute (unknown) keeps the factory default - heavy allowed
    - because unknown is not a tier and must not silently stop the page; the
    settings row says 档位未知 instead. The document is injectable so the
    loop factory gates on the document it actually runs against (the
    visibility instrument feeds a bare fake without a root - factory
    default). */
export function motionHeavyAllowedOn(
  doc: { documentElement?: { getAttribute(name: string): string | null } } | null | undefined,
): boolean {
  const v = doc?.documentElement?.getAttribute(MOTION_ATTR) ?? null;
  return v !== "reduce" && v !== "off";
}

/** The same gate against the page's own document. No document (SSR, the
    render harnesses): the honest answer for "nowhere to read the fact from"
    is the paused one, matching panelWorkAllowed's shape. */
export function motionHeavyAllowed(): boolean {
  if (typeof document === "undefined") return false;
  return motionHeavyAllowedOn(document);
}

/**
 * The snapshot's panelAnimation section, applied to the document - the page
 * half of ticket 193's switch (App calls this on every push, the same seam
 * applyPanelFontSize rides). One function, one mechanism:
 *
 *   - a wire tier word is applied verbatim (the page derives nothing);
 *   - "off" additionally CANCELS every animation the CSS engine is currently
 *     running (the Web Animations API) - the flip to off must stop what is
 *     already in flight, not let it finish unseen;
 *   - an absent section, an empty tier or an unknown word withdraws the
 *     attribute: unknown is unknown - the stylesheet's factory default
 *     stands, and neither off nor full is painted in its place.
 *
 * The registered rAF loops are the visibility module's half: App pairs this
 * call with syncPanelAnimationWorkToMotionTier (lib/panel-visibility), the
 * same one-two the render:animation193 instrument drives.
 */
export function applyPanelAnimationTier(tier: string | null | undefined): void {
  if (typeof document === "undefined") return;
  const root = document.documentElement;
  if (!isAnimationTier(tier)) {
    root.removeAttribute(MOTION_ATTR);
    return;
  }
  applyMotionTier(tier);
  if (tier !== "off") return;
  // 真停: what is already in flight stops NOW, mid-flight - not "runs its
  // remaining instant unseen". getAnimations exists on every WebView2
  // (Chromium 84+); a missing implementation leaves the CSS none-rules in
  // charge and cancels nothing.
  const live = document.getAnimations?.() ?? [];
  for (const a of live) a.cancel();
}
