/* ============================================================================
   index.ts - the minimal locale runtime, ticket 58 AC#1.
   ----------------------------------------------------------------------------
   Hand-written, zero npm dependencies, no i18n library (dispatch rule). The
   bundles live in zh.ts / en.ts; t() picks a line from the current bundle and
   interpolates {name} placeholders from its params argument.

   Language selection: the module defaults to zh and setLanguage() flips it.
   Nothing persists the choice here - the panel is stateless by contract
   (PLAN.md:1044, TestPanelFrontendIsStateless), so when the config wire grows
   an app.language push (ticket 58's switch AC) the host snapshot will drive
   setLanguage(). Until that lands every component re-reads t() at render time
   and the module default keeps the current product face unchanged.
   ============================================================================ */

import { zh, type LocaleKey } from "./zh";
import { en } from "./en";

export type Language = "zh" | "en";

const BUNDLES: Record<Language, Record<LocaleKey, string>> = { zh, en };

let current: Language = "zh";

/** The language t() reads from right now. Defaults to zh. */
export function currentLanguage(): Language {
  return current;
}

/** Flip the bundle t() reads from. The caller owns re-rendering (the ticket
    58 switch AC will drive this off the config wire; no persistence here -
    the panel is stateless by contract). */
export function setLanguage(lang: Language): void {
  current = lang;
}

/** The one lookup the render path uses. A key missing from a bundle falls
    through to the zh bundle and, failing that, to the key itself - a missing
    line shows its own key on screen instead of an empty string, so the loss
    is visible rather than silent. */
export function t(key: LocaleKey, params?: Record<string, string | number>): string {
  const raw = BUNDLES[current][key] ?? BUNDLES.zh[key];
  if (params === undefined) return raw;
  return raw.replace(/\{(\w+)\}/g, (whole, name: string) =>
    Object.prototype.hasOwnProperty.call(params, name) ? String(params[name]) : whole,
  );
}
