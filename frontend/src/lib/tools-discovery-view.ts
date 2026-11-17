/* ============================================================================
   Tool discovery + force rules view - the tools.list / tool.launch /
   rules.force read faces' replies and the page-side derivations that are
   display, not judgement (readface wave 3, ticket 38 AC#3 / AC#5 page half,
   2026-11-17).
   ----------------------------------------------------------------------------
   WHAT ARRIVES. All three doors answer over the same settings door
   (window.wispDispatch) with ONE string, exactly like history.page and the
   privacy doors (cmd/wisp/panel_tools_discovery.go):

     - tools.list: prose ("已安装工具共 <n> 只...") followed by ONE line
       starting with the door's own marker ("wisp-tools-list:") whose
       remainder is the roster JSON {rows:[{name,title,riskPill,argSummary}]}.
       The prose is for humans and is NEVER parsed; the page's facts are read
       off the payload only.

     - rules.force: prose (the 下一 prompt sentence, or the default sentence)
       followed by ONE line starting with "wisp-rules-force:" whose remainder
       is the ForceRulesView JSON. Same rule: prose rides verbatim, facts come
       off the payload.

     - tool.launch: the receipt sentence, no payload line - the launch does
       not wait for a card or a result, so there is nothing to parse and the
       sentence is shown verbatim.

     - a refusal: RefusedEnvelopeForUser's one sentence, shown verbatim.

   THE PILL IS THE TRUTH. riskPill is internal/risk's own Level.String()
   vocabulary carried on the wire. The page renders the WORD it receives,
   tinted by a fixed word->hue table below - the table maps WORDS to token
   colours and never renames, rescales or reorders a level. A payload that
   arrives without the pill renders 未判: a missing field is a different fact
   from any level, and inventing one (or defaulting to "safe") is exactly the
   fake value the dispatch forbids.

   THE CHIP IS DISPLAY + ONE PAYLOAD KEY. The per-task override chip's two
   words are D11③'s own force_tool / force_chat; they leave through
   panel.message.send's forceMode payload (lib/panel.ts sendMessage), and the
   Go side judges the word and consumes it at admission. Nothing here stages,
   clears or reads back an override: rules.force is the read face, the chip is
   the send-time payload, and the two never touch each other on this page.
   ============================================================================ */

import { REFUSED_ENVELOPE_PREFIX } from "@/lib/privacy-view";
import type { ToolRosterRowView } from "@/lib/panel";

/** The two read faces' wire markers, spelled once here; cmd/wisp's
    toolsListPayloadMarker / rulesForcePayloadMarker are the other sides of
    the same strings. The launch receipt carries no marker line. */
export const TOOLS_LIST_MARKER = "wisp-tools-list:";
export const RULES_FORCE_MARKER = "wisp-rules-force:";

/* ---------------------------------------------------------------------------
   tools.list
   --------------------------------------------------------------------------- */

/** What one tools.list reply turned out to be, classified on the wire's own
    words (the refusal wrapper prefix, then the marker line) - the same shape
    lib/history-view.ts established for history.page. */
export type ToolsListState =
  | { phase: "roster"; prose: string; rows: ToolRosterRowView[] }
  | { phase: "refused"; raw: string }
  /** A reply without its payload line: 未判, nothing drawn. */
  | { phase: "no-payload"; raw: string }
  /** A payload that cannot hand over its rows: named, never partially drawn. */
  | { phase: "payload-error"; message: string; raw: string };

/** Classify one tools.list reply. Rows keep the wire's own fields verbatim;
    riskPill/argSummary are kept OPTIONAL exactly as the tooth demands - a
    row without the pill is a row the page renders as 未判, not a row it
    guesses a level for. A row without name or title is a payload defect:
    the whole list refuses rather than rendering an invented row. */
export function parseToolsListReply(reply: string): ToolsListState {
  if (reply.startsWith(REFUSED_ENVELOPE_PREFIX)) {
    return { phase: "refused", raw: reply };
  }
  const at = reply.indexOf(TOOLS_LIST_MARKER);
  if (at < 0) {
    return { phase: "no-payload", raw: reply };
  }
  const afterMarker = reply.slice(at + TOOLS_LIST_MARKER.length);
  const line = afterMarker.split("\n", 1)[0] ?? "";
  let data: unknown;
  try {
    data = JSON.parse(line);
  } catch {
    return { phase: "no-payload", raw: reply };
  }
  if (typeof data !== "object" || data === null) {
    return { phase: "no-payload", raw: reply };
  }
  const raw = data as { rows?: unknown };
  if (!Array.isArray(raw.rows)) {
    return {
      phase: "payload-error",
      message: "名册载荷上 rows 不在或不是数组——名册未判，不画假行",
      raw: reply,
    };
  }
  const rows: ToolRosterRowView[] = [];
  for (const item of raw.rows) {
    if (typeof item !== "object" || item === null) {
      return { phase: "payload-error", message: "名册载荷有一行不是对象——名册未判", raw: reply };
    }
    const row = item as { name?: unknown; title?: unknown; riskPill?: unknown; argSummary?: unknown };
    if (typeof row.name !== "string" || typeof row.title !== "string") {
      return {
        phase: "payload-error",
        message: "名册载荷有一行缺 name/title——这一行不画，也不替它编名字",
        raw: reply,
      };
    }
    rows.push({
      name: row.name,
      title: row.title,
      ...(typeof row.riskPill === "string" ? { riskPill: row.riskPill } : {}),
      ...(typeof row.argSummary === "string" ? { argSummary: row.argSummary } : {}),
    });
  }
  const prose = reply.slice(0, at).trim();
  return { phase: "roster", prose, rows };
}

/** The pill's word->hue table, the ONLY derivation this file performs on a
    level. "L0"/"L1"/"L2" and "Deny" are internal/risk's own words; any other
    word (and an absent pill) is neutral - the word still renders verbatim,
    it just earns no colour this page cannot justify. */
export type RiskPillTone = "ok" | "warn" | "danger" | "neutral";

export function riskPillTone(pill: string | undefined): RiskPillTone {
  if (pill === "L0") return "ok";
  if (pill === "L1") return "warn";
  if (pill === "L2" || pill === "Deny") return "danger";
  return "neutral";
}

/* ---------------------------------------------------------------------------
   rules.force
   --------------------------------------------------------------------------- */

/** One rules row the payload carries, with the wire's own words. source and
    mode stay optional at the parse boundary: an absent word renders 未判 in
    its slot, never a source or mode this page filled in. (The wire's own
    mirror - Go's ForceRuleRow with every key present - is lib/panel.ts's
    ForceRuleRowView; this is the honest parsed shape a partial payload
    produces.) */
export interface ForceRuleRowParsed {
  source?: string;
  pattern?: string;
  mode?: string;
}

/** The rules.force payload as parsed: the wire's mirror is lib/panel.ts's
    ForceRulesViewView; here every key stays optional exactly as the
    missing-field tooth demands. mode "" is the default reading. */
export interface ForceRulesViewParsed {
  mode?: string;
  source?: string;
  pattern?: string;
  instruction?: string;
  rules: ForceRuleRowParsed[];
}

/** What one rules.force reply turned out to be. */
export type ForceRulesState =
  | { phase: "view"; prose: string; view: ForceRulesViewParsed }
  | { phase: "refused"; raw: string }
  | { phase: "no-payload"; raw: string }
  | { phase: "payload-error"; message: string; raw: string };

/** Classify one rules.force reply. The rules rows are read leniently per
    field (an absent word stays absent and renders 未判); a rules value that
    is not an array refuses the payload whole - a face that reports steering
    must not render a silently shorter rule table. */
export function parseRulesForceReply(reply: string): ForceRulesState {
  if (reply.startsWith(REFUSED_ENVELOPE_PREFIX)) {
    return { phase: "refused", raw: reply };
  }
  const at = reply.indexOf(RULES_FORCE_MARKER);
  if (at < 0) {
    return { phase: "no-payload", raw: reply };
  }
  const afterMarker = reply.slice(at + RULES_FORCE_MARKER.length);
  const line = afterMarker.split("\n", 1)[0] ?? "";
  let data: unknown;
  try {
    data = JSON.parse(line);
  } catch {
    return { phase: "no-payload", raw: reply };
  }
  if (typeof data !== "object" || data === null) {
    return { phase: "no-payload", raw: reply };
  }
  const raw = data as {
    mode?: unknown;
    source?: unknown;
    pattern?: unknown;
    instruction?: unknown;
    rules?: unknown;
  };
  if (!Array.isArray(raw.rules)) {
    return {
      phase: "payload-error",
      message: "force 规则载荷上 rules 不在或不是数组——规则表未判，不画假表",
      raw: reply,
    };
  }
  const rules: ForceRuleRowParsed[] = [];
  for (const item of raw.rules) {
    if (typeof item !== "object" || item === null) {
      return { phase: "payload-error", message: "force 规则载荷有一行不是对象——规则表未判", raw: reply };
    }
    const row = item as { source?: unknown; pattern?: unknown; mode?: unknown };
    rules.push({
      ...(typeof row.source === "string" ? { source: row.source } : {}),
      ...(typeof row.pattern === "string" ? { pattern: row.pattern } : {}),
      ...(typeof row.mode === "string" ? { mode: row.mode } : {}),
    });
  }
  const view: ForceRulesViewParsed = {
    ...(typeof raw.mode === "string" ? { mode: raw.mode } : {}),
    ...(typeof raw.source === "string" ? { source: raw.source } : {}),
    ...(typeof raw.pattern === "string" ? { pattern: raw.pattern } : {}),
    ...(typeof raw.instruction === "string" ? { instruction: raw.instruction } : {}),
    rules,
  };
  const prose = reply.slice(0, at).trim();
  return { phase: "view", prose, view };
}
