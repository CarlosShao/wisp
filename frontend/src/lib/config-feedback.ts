/* ============================================================================
   Config copy - every user-visible sentence the settings pages say (v2).
   ----------------------------------------------------------------------------
   WHY THIS FILE EXISTS. The v1 pages were built straight off the interface
   spec, and the spec's internal vocabulary went onto the screen ("腿还没答",
   "两半", raw enum values like no_ref_declared). Owner's ruling: everything a
   user reads is plain language, and the technical words stay between the page
   and the backend (313-panel-ask.md's 措辞修正 section is the replacement
   table this file implements).

   Two rules shape what is here:
     1. A status is one of three words - 未读取 (nobody answered yet), 未配置
        (answered: nothing set), 已配置 (answered: something set) - and an
        error says what the user can do, never which subsystem failed.
     2. A write's answer is TRANSLATED, not echoed: the receipt's key paths and
        tier wording become one plain sentence ("已保存，需要重启后生效"), and a
        refusal becomes the one line a person can act on. The raw reply is
        never rendered.

   Matched substrings on refusals are sealed sentences on the Go side
   (internal/config/settings.go's refusal tags, config_handlers.go's gate
   messages, bridge.go's RefusedEnvelopeForUser prefix). They are matched, not
   parsed: an unrecognised refusal renders the generic line, never a guess.
   ============================================================================ */

import type { ConfigView } from "@/lib/config-view";
import { CONFIG_RELOAD_MARKER } from "@/lib/config-view";
import { t } from "@/lib/locale";
import { GO_KEY_STORED, GO_NEXT_TASK, GO_NO_CHANGE, GO_REFUSED, GO_RESTART, GO_WRITTEN, REFUSAL_RULES } from "@/lib/config-protocol";

/** How a line reads: gray = plain information, green = good, orange =
    attention, red = did not work. */
export type Tone = "neutral" | "ok" | "warn" | "error";

export interface Feedback {
  tone: Tone;
  text: string;
}

/* ---------------------------------------------------------------------------
   Status vocabulary.
   --------------------------------------------------------------------------- */

export interface ConfigStatus {
  label: string;
  tone: Tone;
}

/** The page-level three states. "未读取" is a fact about the READ, not about
    the config: nothing has been set AND nothing has been answered must not
    look alike. */
export function configStatus(view: ConfigView): ConfigStatus {
  if (!view.readable) return { label: t("status.unread"), tone: "neutral" };
  if (!view.configuredKnown) return { label: t("status.unread"), tone: "neutral" };
  return view.configured
    ? { label: t("status.configured"), tone: "ok" }
    : { label: t("status.unconfigured"), tone: "warn" };
}

/** The five credential states as words. `unknown` is "nobody answered", which
    is not the same sentence as "no key yet". */
export function credentialLabel(state: string): string {
  switch (state) {
    case "unknown":
      return t("credential.unknown");
    case "config_unreadable":
      return t("credential.unreadable");
    case "no_ref_declared":
      return t("credential.noRef");
    case "all_recorded":
      return t("credential.allReady");
    case "partly_missing":
      return t("credential.someMissing");
    default:
      return t("credential.unknown");
  }
}

/** One provider's key state as a dot's words. */
export function keyLabel(row: { credentialSet: boolean; credentialError: string }): {
  label: string;
  tone: Tone;
} {
  if (row.credentialError !== "") return { label: t("key.unread"), tone: "neutral" };
  return row.credentialSet
    ? { label: t("key.saved"), tone: "ok" }
    : { label: t("key.unsaved"), tone: "warn" };
}

/* ---------------------------------------------------------------------------
   Write answers, translated.
   --------------------------------------------------------------------------- */

/** The refusal matcher: sealed sentences, most specific first. Each maps to
    the one line a person can act on. Order matters where two sentences share
    a word (a model row's "already exists" must not be read as a provider's;
    the drift refusal's 盘上已变 is matched before anything generic could eat
    it). */
/** The refusal half both write translators share: a matched sealed sentence
    becomes its human line; anything else is the generic "保存失败", never a
    guess dressed up as a diagnosis. */
function refusalFeedback(reply: string): Feedback {
  for (const rule of REFUSAL_RULES) {
    if (reply.includes(rule.match)) return { tone: "error", text: t(rule.key) };
  }
  return { tone: "error", text: t("save.failed") };
}

/** Translate one write reply into the single sentence shown to the user.
    `null` means no host answered at all. A no-op save (the file already
    carries the value, ticket 39 AC#4) is its own answer - "没有变化", not a
    green 已保存 that reads as if something was written. */
export function translateWriteReply(reply: string | null): Feedback {
  if (reply === null) return { tone: "error", text: t("save.noHost") };
  if (reply.includes(GO_REFUSED)) {
    return refusalFeedback(reply);
  }
  if (reply.includes(GO_NO_CHANGE)) {
    return { tone: "neutral", text: t("save.noChange") };
  }
  if (reply.includes(GO_KEY_STORED)) {
    return { tone: "ok", text: restartTier(reply) ? t("save.keySavedRestart") : t("key.saved") };
  }
  if (reply.includes(GO_WRITTEN)) {
    // Ticket 39 AC#1: the restart-required line on screen carries 需要重启
    // verbatim - the same word the Go receipt's restart sentence is pinned by
    // (tierSentence), so the page and the receipt say one thing word for word.
    return { tone: "ok", text: restartTier(reply) ? t("save.savedRestart") : t("save.saved") };
  }
  // Anything that is neither a known receipt nor a refusal: say so plainly
  // rather than claim a save.
  return { tone: "error", text: t("save.failed") };
}

/** The composer row's model/thinking switches (ticket 187) translate through
    this one instead of translateWriteReply. The Go receipt states the
    effective tier itself ("这一项在下一次任务开始时生效" for a thinking
    write, "这一项要重启进程…" for a model write - config_handlers.go's
    tierSentence), and the two halves of the row have DIFFERENT tiers: a
    translator that collapsed both into one generic "已保存" would let the
    model half look in force before a restart and the tier half wait for one.
    So the receipt's own tier wording is kept. Refusals match exactly like
    translateWriteReply. */
export function translateTieredWriteReply(reply: string | null): Feedback {
  if (reply === null) return { tone: "error", text: t("save.noHost") };
  if (reply.includes(GO_REFUSED)) {
    return refusalFeedback(reply);
  }
  if (reply.includes(GO_NO_CHANGE)) {
    return { tone: "neutral", text: t("save.noChange") };
  }
  if (reply.includes(GO_WRITTEN)) {
    if (reply.includes(GO_NEXT_TASK)) return { tone: "ok", text: t("save.savedNextTask") };
    if (reply.includes(GO_RESTART)) return { tone: "ok", text: t("save.savedRestartProcess") };
    return { tone: "ok", text: t("save.saved") };
  }
  return { tone: "error", text: t("save.failed") };
}

/** All settings writes are restart-tier today; the receipt carries the tier
    sentence and this reads it rather than assuming. */
function restartTier(reply: string): boolean {
  return reply.includes(GO_RESTART);
}

/* ---------------------------------------------------------------------------
   The plugin toggle's own receipt translation (ticket 407).
   ---------------------------------------------------------------------------
   WHY THIS IS NOT translateTieredWriteReply. The plugin toggle's receipt is
   not a generic save line: the accepted receipt carries the leg's honest
   note (配置层已生效；引擎未接线；面板侧只能关不能开 - cmd/wisp
   panel_config_store.go:567-579) whose whole point is that enabled=true is
   NOT "running", and a refused plugin write carries sentences the user must
   see word for word - the missing-configPlugin gate refusal
   (internal/panel/config_handlers.go:575-579, "必须指名插件 id（configPlugin
   ＝[plugins.<id>] 的 <id>…）") and the locked-loosening refusal (R20: the
   panel is a permission input, not an authority). A translator that
   collapsed those into 已保存/保存失败 would swallow exactly the sentences
   ticket 407's teeth exist to keep on screen. Precedent for verbatim prose:
   translateFetchReply and the probe card render the host's 原文 when the
   原文 IS the answer.

   The machine half of an accepted reply (the CONFIG_RELOAD_MARKER event
   line, config_handlers.go:623) is dropped here - notifyConfigReloadListeners
   already parsed it; it never belongs on screen.
   -------------------------------------------------------------------------- */

/** The reply's prose minus the reload-event machine line (and surrounding
    blank lines), ready to be shown verbatim. */
function pluginReceiptProse(reply: string): string {
  return reply
    .split("\n")
    .filter((line) => !line.startsWith(CONFIG_RELOAD_MARKER))
    .join("\n")
    .trim();
}

/** The plugin toggle's receipt: the host's own words, verbatim, never a
    summary. tone follows the reply's own verdict word (被拒绝 = error,
    没有任何键路径发生变化 = neutral no-op, otherwise an accepted write). */
export function translatePluginToggleReply(reply: string | null): Feedback {
  if (reply === null) return { tone: "error", text: t("save.noHost") };
  const prose = pluginReceiptProse(reply);
  if (reply.includes(GO_REFUSED)) return { tone: "error", text: prose };
  if (reply.includes(GO_NO_CHANGE)) return { tone: "neutral", text: prose };
  return { tone: "ok", text: prose };
}


/** The model-list fetch's answer when it did not carry a payload: a refusal
    gets its matched human line; anything else (the host's D37 prose) is shown
    verbatim - the fetch has no receipt to translate, and inventing a summary
    for prose this page never parsed would be a guess. */
export function translateFetchReply(reply: string | null): Feedback {
  if (reply === null) return { tone: "error", text: t("fetch.noHost") };
  if (reply.includes(GO_REFUSED)) {
    for (const rule of REFUSAL_RULES) {
      if (reply.includes(rule.match)) return { tone: "error", text: t(rule.key) };
    }
  }
  const text = reply.trim();
  return { tone: "error", text: text === "" ? t("fetch.failed") : text };
}

/** Local validation lines (checked before anything is sent - the pair rule is
    stated at save time, not by disabling the button, which reads as a black
    box to the person looking at it). */
export const COPY = {
  pairMissing: t("copy.pairMissing"),
  chainSegmentMissing: t("copy.chainSegmentMissing"),
  chainEmpty: t("copy.chainEmpty"),
  chatPairMissing: t("copy.pairMissing"),
  testKeyUnavailable: t("copy.testKeyUnavailable"),
  deleteUnavailable: t("copy.deleteUnavailable"),
  editUnavailable: t("copy.editUnavailable"),
  baseURLSaved: t("save.savedRestart"),
  noProvidersYet: t("copy.noProvidersYet"),
  addProvider: t("copy.addProvider"),
  save: t("copy.save"),
  saving: t("copy.saving"),
} as const;
