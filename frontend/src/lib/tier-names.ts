/* ============================================================================
   Tier names - the native mode names and the Chinese the panel shows for them.
   ----------------------------------------------------------------------------
   Spelled ONCE, in one file, so a fourth tier cannot be invented twice
   differently (App.tsx's header comment records the origin of this map). Two
   readers today: App.tsx's composer permission picker and the settings page's
   安全 category, which shows the mode the snapshot carries RIGHT NOW.

   The values are display words only. Nothing here judges a mode or applies
   one: switching costs the native L2 confirmation on the Go side
   (panel.mode.request), and this file has no way to record an answer.
   An unknown name falls back to the raw wire string at the reader - never to
   a safer-sounding word.
   ============================================================================ */

export const TIER_NAMES: Record<string, string> = {
  ask_every_step: "每步都问",
  ask_high_risk: "高风险才问",
  auto_approve: "全自动",
  unknown: "档位未知",
};

/** One tier's display word: the map's entry, or the raw wire name when the
    vocabulary grew a word this file has not learned. */
export function tierName(mode: string): string {
  return TIER_NAMES[mode] ?? mode;
}
