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

   The keys are the Go side's CostDayRow (internal/panel/cost_handlers.go:30-46,
   the JSON key set reconciled by TestComposerContractTypesMatchFrontend's
   family) and the wrapper costPayload (cost_handlers.go:214-225):
     days[]     one row per cost_daily row, NEWEST DAY FIRST, carried verbatim:
                day (YYYY-MM-DD), tokensIn, tokensOut, costMicros (integer
                micro units - 1e-6 currency, no floats cross this seam),
                tasks, cachedTokens (readface wave 1: always 0 on this
                build - cost_daily has no cached column and the leg may not
                widen the D35 table; the payload's cachedKnown=false gate is
                what makes the page show 未知 instead of the rows' zero).
     count      the row count the handler read (and the days length; the two
                agree by construction).
     cachedKnown     false = this build's storage aggregates no cached tokens;
                the cached column renders 未知, never a number (never the
                rows' zero dressed up as a reading).
     aggregateKnown  true = the leg answered the OPTIONAL aggregate face
                (CostAggregateReader, cost_handlers.go:92) and `aggregate`
                carries the month totals, the [cost] budgets, the three-state
                verdict word and the pricing version. false = the leg has no
                aggregate face or its read failed (the prose says so too);
                the page draws NO aggregate/budget/pricing number, it keeps
                the named absence.
     aggregate  CostAggregateRow verbatim (cost_handlers.go:54-83): ten
                fields, integers only, budgetState is the adapter's own word
                ("normal" | "warn" | "paused") computed in internal/observe
                on integer basis points - the page never re-judges it.

   WHAT STILL NEVER ARRIVES (and what this file therefore refuses to invent):
     cached numbers while cachedKnown=false - the gate is the fact.
     aggregate/budget/pricing numbers while aggregateKnown=false - a partial
                answer that SAYS it is partial is not a silence to paper
                over, and no client-side total is derived anywhere (the
                handler's own ban stands, cost_handlers.go:26-29).

   THE HOOK. useCostSummary() sends one cost.summary on mount and exposes
   reload() for a re-read. "no-host" (plain browser) and "unreadable" (the
   host answered with a refusal sentence, no payload line) stay separate
   states, exactly like the config read's pair: a named refusal and an empty
   table are different facts, and rendering them alike is the lie.
   ============================================================================ */

import { useCallback, useEffect, useRef, useState } from "react";
import { requestCostSummary } from "@/lib/panel";

/** One cost_daily row, as the wire carries it (CostDayRow). cachedTokens is
    always 0 on this build (no cached column on cost_daily); the number is
    displayed ONLY while cachedKnown=true, never as a reading on its own. */
export interface CostDayRowView {
  day: string;
  tokensIn: number;
  tokensOut: number;
  /** Integer micro units (1e-6 currency). Never turned into a float here. */
  costMicros: number;
  tasks: number;
  cachedTokens: number;
}

/** The aggregate face of the cost read (ticket 40 AC#3/AC#4), verbatim off
    CostAggregateRow (internal/panel/cost_handlers.go:54-83). Every number is
    the leg's own integer; budgetState is the adapter's own three-state word
    (internal/observe's basis-point judge) - the page re-judges nothing. */
export interface CostAggregateView {
  dayCostMicros: number;
  monthTokensIn: number;
  monthTokensOut: number;
  monthCostMicros: number;
  dailyBudgetMicros: number;
  monthlyBudgetMicros: number;
  budgetState: string;
  budgetAllowNewTask: boolean;
  pricingVersion: number;
  pricingUpdated: string;
}

/** The wire marker. Spelled once here; internal/panel's
    CostSummaryPayloadMarker is the other side of the same string. */
export const COST_SUMMARY_MARKER = "wisp-cost-summary:";

/** What one cost.summary reply turned out to be. A failure carries only prose
    (the handler's refusal sentence), which the page shows verbatim in place.
    A parsed reply carries the two gates verbatim: cachedKnown (false = this
    build aggregates no cached tokens - the cached column is 未知, not a
    number) and aggregateKnown (false = no aggregate face or its read failed -
    no aggregate/budget/pricing number is drawn at all). aggregateKnown=true
    with a malformed aggregate block is a payload-error: the payload claimed a
    read it cannot hand over, and guessing fields would be worse. */
export type CostSummaryReply =
  | {
      kind: "parsed";
      days: CostDayRowView[];
      count: number;
      skipped: number;
      prose: string;
      cachedKnown: boolean;
      aggregateKnown: boolean;
      aggregate: CostAggregateView | null;
    }
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
  const raw = data as { days?: unknown; count?: unknown; cachedKnown?: unknown; aggregateKnown?: unknown; aggregate?: unknown };
  const days: CostDayRowView[] = [];
  // A row whose day is not a non-empty string, or whose numeric fields are not
  // finite numbers, is SKIPPED AND COUNTED - never folded into zeros (a zero
  // row would read as "that day cost nothing", a fact the wire did not state).
  // The view states the skipped count instead of hiding the damage.
  let skipped = 0;
  const cachedKnown = raw.cachedKnown === true;
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
      // cachedTokens rides with the rows but is a READING only while the gate
      // is open: while cachedKnown is false the field is the adapter's
      // structural zero and must not fail a row - while it is true a missing
      // or non-finite cached number is real damage (the page would show a
      // made-up column).
      const cachedTokens = typeof row.cachedTokens === "number" && Number.isFinite(row.cachedTokens) ? row.cachedTokens : cachedKnown ? NaN : 0;
      if (day === "" || Number.isNaN(tokensIn) || Number.isNaN(tokensOut) || Number.isNaN(costMicros) || Number.isNaN(tasks) || Number.isNaN(cachedTokens)) {
        skipped += 1;
        continue;
      }
      days.push({ day, tokensIn, tokensOut, costMicros, tasks, cachedTokens });
    }
  }
  const count = typeof raw.count === "number" ? raw.count : days.length;
  // The aggregate block: gated by aggregateKnown, validated as a whole. A
  // block that claims to be known but misses a field or carries a non-finite
  // number is refused as a payload-error, never partially drawn.
  const aggregateKnown = raw.aggregateKnown === true;
  let aggregate: CostAggregateView | null = null;
  if (aggregateKnown) {
    if (typeof raw.aggregate !== "object" || raw.aggregate === null) {
      return { kind: "payload-error", message: "aggregateKnown=true 但 aggregate 块不在载荷上（页面拒绝猜字段）", raw: reply };
    }
    const a = raw.aggregate as Record<string, unknown>;
    const num = (v: unknown): number => (typeof v === "number" && Number.isFinite(v) ? v : NaN);
    const dayCostMicros = num(a.dayCostMicros);
    const monthTokensIn = num(a.monthTokensIn);
    const monthTokensOut = num(a.monthTokensOut);
    const monthCostMicros = num(a.monthCostMicros);
    const dailyBudgetMicros = num(a.dailyBudgetMicros);
    const monthlyBudgetMicros = num(a.monthlyBudgetMicros);
    const budgetState = typeof a.budgetState === "string" ? a.budgetState : "";
    const budgetAllowNewTask = typeof a.budgetAllowNewTask === "boolean" ? a.budgetAllowNewTask : null;
    const pricingVersion = num(a.pricingVersion);
    const pricingUpdated = typeof a.pricingUpdated === "string" ? a.pricingUpdated : "";
    if (
      Number.isNaN(dayCostMicros) || Number.isNaN(monthTokensIn) || Number.isNaN(monthTokensOut) ||
      Number.isNaN(monthCostMicros) || Number.isNaN(dailyBudgetMicros) || Number.isNaN(monthlyBudgetMicros) ||
      budgetState === "" || budgetAllowNewTask === null || Number.isNaN(pricingVersion)
    ) {
      return { kind: "payload-error", message: "aggregateKnown=true 但 aggregate 块字段残缺/类型不对（页面拒绝猜字段）", raw: reply };
    }
    aggregate = {
      dayCostMicros,
      monthTokensIn,
      monthTokensOut,
      monthCostMicros,
      dailyBudgetMicros,
      monthlyBudgetMicros,
      budgetState,
      budgetAllowNewTask,
      pricingVersion,
      pricingUpdated,
    };
  }
  return { kind: "parsed", days, count, skipped, prose: reply.slice(0, at).trimEnd(), cachedKnown, aggregateKnown, aggregate };
}

/** What the 成本 category renders from. ready carries the two gates and the
    aggregate block (null unless the leg answered the aggregate face). */
export type CostReadState =
  | { phase: "loading" }
  | { phase: "no-host" }
  | {
      phase: "ready";
      days: CostDayRowView[];
      count: number;
      skipped: number;
      prose: string;
      raw: string;
      cachedKnown: boolean;
      aggregateKnown: boolean;
      aggregate: CostAggregateView | null;
    }
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
          cachedKnown: parsed.cachedKnown,
          aggregateKnown: parsed.aggregateKnown,
          aggregate: parsed.aggregate,
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
