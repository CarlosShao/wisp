/* ============================================================================
   Cost view - ticket 194's cost.summary payload, its parser, and the hook the
   settings page's 成本 category reads it through.
   ----------------------------------------------------------------------------
   WHAT ARRIVES. cost.summary answers a multi-line string: one prose sentence
   ("成本按天共 N 行（新的一天在前）。") and AFTER it one line beginning with
   CostSummaryPayloadMarker ("wisp-cost-summary:"). What follows that marker on
   the same line is the JSON this file parses. The prose is for humans and is
   NEVER parsed - the same rule lib/config-view.ts states for its marker, and
   for the same reason.

   The keys are the Go side's CostDayRow (internal/panel/cost_handlers.go:30-36,
   the JSON key set reconciled by TestComposerContractTypesMatchFrontend's
   family) and the wrapper costPayload (cost_handlers.go:124-127):
     days[]     one row per cost_daily row, NEWEST DAY FIRST, carried verbatim:
                day (YYYY-MM-DD), tokensIn, tokensOut, costMicros (integer
                micro units - 1e-6 currency, no floats cross this seam),
                tasks.
     count      the row count the handler read (and the days length; the two
                agree by construction).

   WHAT NEVER ARRIVES (and what this file therefore refuses to invent):
     cached tokens   CostDayRow carries no cached column - IN/OUT only.
     aggregates      the handler's own comment bans derived totals on this wire
                     ("a summary that derives a number is a second source of
                     truth", cost_handlers.go:26-29) - the page does the same
                     arithmetic nowhere.
     budget state    the 80% warn / 100% paused bars have no wire fields.
     pricing version the footnote has no wire field.
   Each absence is NAMED on the page instead of drawn as a zero.

   THE HOOK. useCostSummary() sends one cost.summary on mount and exposes
   reload() for a re-read. "no-host" (plain browser) and "unreadable" (the
   host answered with a refusal sentence, no payload line) stay separate
   states, exactly like the config read's pair: a named refusal and an empty
   table are different facts, and rendering them alike is the lie.
   ============================================================================ */

import { useCallback, useEffect, useRef, useState } from "react";
import { requestCostSummary } from "@/lib/panel";

/** One cost_daily row, as the wire carries it (CostDayRow). */
export interface CostDayRowView {
  day: string;
  tokensIn: number;
  tokensOut: number;
  /** Integer micro units (1e-6 currency). Never turned into a float here. */
  costMicros: number;
  tasks: number;
}

/** The wire marker. Spelled once here; internal/panel's
    CostSummaryPayloadMarker is the other side of the same string. */
export const COST_SUMMARY_MARKER = "wisp-cost-summary:";

/** What one cost.summary reply turned out to be. A failure carries only prose
    (the handler's refusal sentence), which the page shows verbatim in place. */
export type CostSummaryReply =
  | { kind: "parsed"; days: CostDayRowView[]; count: number; skipped: number; prose: string }
  | { kind: "payload-error"; message: string; raw: string }
  | { kind: "no-payload"; raw: string };

/** Parse one cost.summary reply. Unknown fields are ignored on purpose (the
    wire shape is the Go side's to grow); day rows whose day is not a non-empty
    string are skipped rather than guessed into shape. */
export function parseCostSummaryReply(reply: string): CostSummaryReply {
  const at = reply.indexOf(COST_SUMMARY_MARKER);
  if (at < 0) return { kind: "no-payload", raw: reply };
  const afterMarker = reply.slice(at + COST_SUMMARY_MARKER.length);
  const line = afterMarker.split("\n", 1)[0] ?? "";
  let data: unknown;
  try {
    data = JSON.parse(line);
  } catch {
    return { kind: "no-payload", raw: reply };
  }
  if (typeof data !== "object" || data === null) return { kind: "no-payload", raw: reply };
  // The Go side's own marshalling failure arrives as {"payloadError": "..."}.
  // A named state, never an empty table.
  const err = (data as { payloadError?: unknown }).payloadError;
  if (typeof err === "string" && err !== "") {
    return { kind: "payload-error", message: err, raw: reply };
  }
  const raw = data as { days?: unknown; count?: unknown };
  const days: CostDayRowView[] = [];
  // A row whose day is not a non-empty string, or whose numeric fields are not
  // finite numbers, is SKIPPED AND COUNTED - never folded into zeros (a zero
  // row would read as "that day cost nothing", a fact the wire did not state).
  // The view states the skipped count instead of hiding the damage.
  let skipped = 0;
  if (Array.isArray(raw.days)) {
    for (const entry of raw.days) {
      if (typeof entry !== "object" || entry === null) {
        skipped += 1;
        continue;
      }
      const row = entry as Record<string, unknown>;
      const day = typeof row.day === "string" ? row.day : "";
      const tokensIn = typeof row.tokensIn === "number" && Number.isFinite(row.tokensIn) ? row.tokensIn : NaN;
      const tokensOut = typeof row.tokensOut === "number" && Number.isFinite(row.tokensOut) ? row.tokensOut : NaN;
      const costMicros = typeof row.costMicros === "number" && Number.isFinite(row.costMicros) ? row.costMicros : NaN;
      const tasks = typeof row.tasks === "number" && Number.isFinite(row.tasks) ? row.tasks : NaN;
      if (day === "" || Number.isNaN(tokensIn) || Number.isNaN(tokensOut) || Number.isNaN(costMicros) || Number.isNaN(tasks)) {
        skipped += 1;
        continue;
      }
      days.push({ day, tokensIn, tokensOut, costMicros, tasks });
    }
  }
  const count = typeof raw.count === "number" ? raw.count : days.length;
  return { kind: "parsed", days, count, skipped, prose: reply.slice(0, at).trimEnd() };
}

/** What the 成本 category renders from. */
export type CostReadState =
  | { phase: "loading" }
  | { phase: "no-host" }
  | { phase: "ready"; days: CostDayRowView[]; count: number; skipped: number; prose: string; raw: string }
  | { phase: "unreadable"; raw: string }
  | { phase: "payload-error"; message: string; raw: string };

/** Read the cost summary once, and expose reload() for a re-read. */
export function useCostSummary(): {
  state: CostReadState;
  reload: () => void;
} {
  const [state, setState] = useState<CostReadState>({ phase: "loading" });
  // A per-call token: a reload that lands after a newer one must not overwrite
  // the newer read (lib/config-view.ts's useConfigView shape).
  const seq = useRef(0);

  const reload = useCallback(() => {
    const mine = ++seq.current;
    void requestCostSummary().then((reply) => {
      if (mine !== seq.current) return;
      if (reply === null) {
        setState({ phase: "no-host" });
        return;
      }
      const parsed = parseCostSummaryReply(reply);
      if (parsed.kind === "parsed") {
        setState({
          phase: "ready",
          days: parsed.days,
          count: parsed.count,
          skipped: parsed.skipped,
          prose: parsed.prose,
          raw: reply,
        });
      } else if (parsed.kind === "payload-error") {
        setState({ phase: "payload-error", message: parsed.message, raw: reply });
      } else {
        setState({ phase: "unreadable", raw: reply });
      }
    });
  }, []);

  useEffect(() => {
    reload();
  }, [reload]);

  return { state, reload };
}
