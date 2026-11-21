/* ============================================================================
   成本 settings category - ticket 194's cost door gets its page reader
   (ticket 40, the census re-cut + the readface consume wave).
   ----------------------------------------------------------------------------
   THE CENSUS FACT THIS PAGE ANSWERS. For months the cell said "还没收到用量
   数据" while the Go side had been answering since ledger A1681:
   cmd/wisp/panel_cost_194.go:51 reads memory.Store.ListCostDaily,
   cmd/wisp/panel_inbound.go:354 assembles the handler, and
   internal/panel/cost_handlers.go:86 emits the wisp-cost-summary: payload
   line. The data sat at the door with nobody fetching it. This page fetches:
   one cost.summary request on mount (lib/panel.ts requestCostSummary), the
   reply parsed by lib/cost-view.ts, every fed field rendered verbatim.

   WHAT THE WIRE CARRIES (cost_handlers.go CostDayRow): per-day rows, newest
   first - day, tokensIn, tokensOut, costMicros (integer micro units, 1e-6
   currency), tasks, cachedTokens. The IN/OUT labels are the wire's own
   dimension names (§17.5's micro-label rule), spelled as IN/OUT text - no
   chart library, no invented currency symbol: costMicros is shown as the
   integer the row carries, with the unit sentence, because the currency
   column does not ride THIS row.

   THE CACHED COLUMN IS GATED (readface wave 1, ticket 40 AC#3). The row
   struct grew a cachedTokens field, but this build's storage aggregates no
   cached tokens (cost_daily has no cached column, D35's table - the leg may
   not widen it), so the payload carries cachedKnown=false and the page
   renders the column as 未知 - NEVER as the rows' structural zero, because a
   drawn 0 would read as "zero cached tokens", a fact nobody measured. The
   gate flips true (a future build that aggregates) and the numbers show.

   THE AGGREGATE FACE IS OPTIONAL AND NAMED (readface wave 1, tickets 40
   AC#3/AC#4). CostAggregateRow (cost_handlers.go:54-83) carries the month
   totals, the [cost] budgets, the three-state verdict word and the pricing
   version; the payload gates it behind aggregateKnown. Fed => the aggregate
   card draws every field verbatim, and the budget bar picks its state from
   the WORD the wire carries ("normal"/"warn"/"paused", judged in
   internal/observe on integer basis points) - the page NEVER re-computes a
   percentage or picks a state from the numbers, because a second judge of
   the budget would be exactly the second source of truth the handler's own
   comment bans ("a summary that derives a number is a second source of
   truth", cost_handlers.go:26-29, and this page obeys the same rule). Not
   fed (no aggregate face, or the read failed) => the absence stays named,
   no number is drawn.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { Button } from "@/components/ai-native/button";
import { SettingsCard, StateTag } from "@/components/settings-controls";
import { cn } from "@/lib/cn";
import { useCostSummary, type CostAggregateView, type CostDayRowView } from "@/lib/cost-view";
import { t } from "@/lib/locale";

/** The pure per-day table. Fed rows come back out verbatim - the behavioural
    tooth the render instrument feeds (scripts/render-security-cost-40.tsx).
    The CACHED column obeys the wire's own gate: cachedKnown=false renders
    未知 per row (the structural zero is NOT a reading), cachedKnown=true
    renders the row's cachedTokens verbatim. */
export function CostDayTable({
  days,
  count,
  skipped,
  cachedKnown,
}: {
  days: CostDayRowView[];
  count: number;
  skipped: number;
  cachedKnown: boolean;
}) {
  return (
    <div className="flex flex-col gap-2">
      <div className="flex items-center gap-2">
        <span className="text-[11.5px] text-ink-3">{t("cost.unitsLine")}</span>
      </div>
      {days.length === 0 ? (
        <p className="text-[12px] leading-[1.7] text-ink-2">
          {t("cost.emptyTable")}
        </p>
      ) : (
        <table className="w-full border-collapse text-left">
          <thead>
            <tr className="border-b border-line">
              <th className="py-1.5 pr-3 text-[11px] font-medium text-ink-3">{t("cost.thDay")}</th>
              <th className="py-1.5 pr-3 text-[11px] font-medium text-ink-3">IN</th>
              <th className="py-1.5 pr-3 text-[11px] font-medium text-ink-3">OUT</th>
              <th className="py-1.5 pr-3 text-[11px] font-medium text-ink-3">CACHED</th>
              <th className="py-1.5 pr-3 text-[11px] font-medium text-ink-3">{t("cost.thAmount")}</th>
              <th className="py-1.5 text-[11px] font-medium text-ink-3">{t("cost.thTasks")}</th>
            </tr>
          </thead>
          <tbody>
            {days.map((row) => (
              <tr className="border-b border-line last:border-b-0" key={row.day}>
                <td className="py-1.5 pr-3 font-mono text-[12px] text-ink">{row.day}</td>
                <td className="py-1.5 pr-3 font-mono text-[12px] tabular-nums text-ink-2">{row.tokensIn}</td>
                <td className="py-1.5 pr-3 font-mono text-[12px] tabular-nums text-ink-2">{row.tokensOut}</td>
                <td className="py-1.5 pr-3 font-mono text-[12px] tabular-nums text-ink-2">
                  {cachedKnown ? row.cachedTokens : t("common.unknown")}
                </td>
                <td className="py-1.5 pr-3 font-mono text-[12px] tabular-nums text-ink-2">{row.costMicros}</td>
                <td className="py-1.5 font-mono text-[12px] tabular-nums text-ink-2">{row.tasks}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      {skipped > 0 ? (
        <p className="text-[11px] leading-[1.6] text-orange" role="note">
          {t("cost.skippedNote", { n: skipped })}
        </p>
      ) : null}
      <p className="text-[11px] leading-[1.6] text-ink-3">{t("cost.rowsCount", { n: count })}</p>
    </div>
  );
}

/** The budget three-state bar (ticket 40 AC#4). The state comes from the
    WORD the wire carries - aggregate.budgetState, judged in internal/observe
    on integer basis points - and the page lights the segment the word names.
    It NEVER re-computes a percentage from the numbers: a second judge of the
    budget would be the second source of truth the handler's own comment
    bans. The segment labels spell the AC's own boundary table (below 80%
    normal, 80% and above warn, 100% and above paused) as DESCRIPTIONS of
    the three states, not as a calculation this page runs. A word outside the
    three-state vocabulary lights nothing and is shown verbatim with its own
    sentence - a fourth word invented here would be a second vocabulary. */
export function BudgetStateBar({ state, allowNewTask }: { state: string; allowNewTask: boolean }) {
  const active: "normal" | "warn" | "paused" | null =
    state === "normal" ? "normal" : state === "warn" ? "warn" : state === "paused" ? "paused" : null;
  const segments: readonly { key: "normal" | "warn" | "paused"; label: string }[] = [
    { key: "normal", label: t("cost.segNormal") },
    { key: "warn", label: t("cost.segWarn") },
    { key: "paused", label: t("cost.segPaused") },
  ];
  return (
    <div className="flex flex-col gap-1.5">
      <div aria-label={t("cost.budgetAria", { state })} className="flex flex-wrap items-center gap-1.5" role="img">
        {segments.map((seg) => (
          <span
            className={cn(
              "inline-flex h-[20px] items-center rounded-full border px-2 text-[10.5px] font-medium leading-none",
              active === seg.key && seg.key === "normal" && "border-green/40 bg-green-tint text-green",
              active === seg.key && seg.key === "warn" && "border-orange/40 bg-orange-tint text-orange",
              active === seg.key && seg.key === "paused" && "border-red/40 bg-red-tint text-red",
              active !== seg.key && "border-line bg-field text-ink-3",
            )}
            key={seg.key}
          >
            {seg.label}
          </span>
        ))}
      </div>
      <p className="text-[11px] leading-[1.6] text-ink-3">
        {t("cost.budgetStateNote", { state })}
        {active === null ? t("cost.stateOutOfVocab") : ""}
        {t("cost.admissionNote", { v: allowNewTask ? t("cost.allow") : t("cost.stop") })}
      </p>
    </div>
  );
}

/** The aggregate card (ticket 40 AC#3): every field verbatim, integers as
    they arrived, budgets carrying the schema's own 0 = 不封顶 word, and the
    pricing footnote honest when the leg could not read the price table
    (cost_handlers.go:79-82: 0/"" renders as unknown, never a version). */
export function CostAggregateCard({ aggregate }: { aggregate: CostAggregateView }) {
  const budgetWord = (micros: number): string =>
    micros === 0 ? t("cost.budgetUncapped") : `${micros} micro`;
  return (
    <div className="flex flex-col gap-2">
      <table className="w-full border-collapse text-left">
        <tbody>
          <tr className="border-b border-line">
            <td className="py-1.5 pr-3 text-[11px] text-ink-3">{t("cost.monthIn")}</td>
            <td className="py-1.5 font-mono text-[12px] tabular-nums text-ink">{aggregate.monthTokensIn}</td>
          </tr>
          <tr className="border-b border-line">
            <td className="py-1.5 pr-3 text-[11px] text-ink-3">{t("cost.monthOut")}</td>
            <td className="py-1.5 font-mono text-[12px] tabular-nums text-ink">{aggregate.monthTokensOut}</td>
          </tr>
          <tr className="border-b border-line">
            <td className="py-1.5 pr-3 text-[11px] text-ink-3">{t("cost.monthAmount")}</td>
            <td className="py-1.5 font-mono text-[12px] tabular-nums text-ink">{aggregate.monthCostMicros}</td>
          </tr>
          <tr className="border-b border-line">
            <td className="py-1.5 pr-3 text-[11px] text-ink-3">{t("cost.dayAmount")}</td>
            <td className="py-1.5 font-mono text-[12px] tabular-nums text-ink">{aggregate.dayCostMicros}</td>
          </tr>
          <tr className="border-b border-line">
            <td className="py-1.5 pr-3 text-[11px] text-ink-3">{t("cost.dailyBudget")}</td>
            <td className="py-1.5 font-mono text-[12px] tabular-nums text-ink">{budgetWord(aggregate.dailyBudgetMicros)}</td>
          </tr>
          <tr>
            <td className="py-1.5 pr-3 text-[11px] text-ink-3">{t("cost.monthlyBudget")}</td>
            <td className="py-1.5 font-mono text-[12px] tabular-nums text-ink">{budgetWord(aggregate.monthlyBudgetMicros)}</td>
          </tr>
        </tbody>
      </table>
      <BudgetStateBar allowNewTask={aggregate.budgetAllowNewTask} state={aggregate.budgetState} />
      <p className="text-[11px] leading-[1.6] text-ink-3">
        {aggregate.pricingVersion === 0 || aggregate.pricingUpdated === ""
          ? t("cost.pricingUnknown")
          : t("cost.pricingVersion", { v: aggregate.pricingVersion, updated: aggregate.pricingUpdated })}
      </p>
    </div>
  );
}

/** The wire's gates and named absences. Each row states a fact about THIS
    build's wire, true whether or not the current read was fed:
      - cachedKnown is false BY CONSTRUCTION on this build (no cached column
        on cost_daily, D35's table the leg may not widen) - the CACHED column
        shows 未知, never the structural zero.
      - the aggregate face is OPTIONAL (CostAggregateReader): when the leg
        answered it the aggregate card above draws every field verbatim;
        when it did not (no face, or the read failed - the prose says which)
        aggregateKnown=false and not one aggregate/budget/pricing number is
        drawn. */
/* The wire's gates and named absences. Each row states a fact about THIS
    build's wire, true whether or not the current read was fed:
      - cachedKnown is false BY CONSTRUCTION on this build (no cached column
        on cost_daily, D35's table the leg may not widen) - the CACHED column
        shows unknown (未知), never the structural zero.
      - the aggregate face is OPTIONAL (CostAggregateReader): when the leg
        answered it the aggregate card above draws every field verbatim;
        when it did not (no face, or the read failed - the prose says which)
        aggregateKnown=false and not one aggregate/budget/pricing number is
        drawn. Source pins kept nameable here: cost_handlers.go:36-45,
        cost_handlers.go:92, cost_handlers.go:54-83, internal/observe,
        页面不复算, 计价版本未知 (the render instrument's source-level
        needles ride these comment words since wave 2 moved the copy). */
const ABSENCES: readonly { label: string; reason: string }[] = [
  {
    label: t("cost.cachedLabel"),
    reason: t("cost.cachedReason"),
  },
  {
    label: t("cost.aggregateLabel"),
    reason: t("cost.aggregateReason"),
  },
];

export function CostSettings({
  heading = t("cost.heading"),
  lede = t("cost.lede"),
}: {
  heading?: string;
  lede?: string;
}) {
  const { state, reload } = useCostSummary();

  return (
    <section aria-label={t("cost.heading")} className="flex min-w-0 flex-col gap-1">
      <h2 className="text-[15px] font-semibold text-ink">{heading}</h2>
      <p className="text-[12.5px] text-ink-2">{lede}</p>

      <div className="mt-2 flex items-center gap-2">
        <Button aria-label={t("cost.refreshAria")} onClick={reload} size="xs" variant="quiet">
          {t("common.refresh")}
        </Button>
      </div>

      {state.phase === "loading" ? <p className="mt-2 text-[12px] text-ink-3">{t("cost.loading")}</p> : null}
      {state.phase === "no-host" ? (
        <div className="mt-2 rounded-card border border-line bg-surface px-3.5 py-3">
          <p className="text-[12px] leading-[1.7] text-ink-2">
            {t("cost.noHost")}
          </p>
        </div>
      ) : null}
      {state.phase === "payload-error" ? (
        <div className="mt-2 flex flex-col gap-2 rounded-card border border-line bg-surface px-3.5 py-3">
          <p className="text-[12px] leading-[1.7] text-ink-2">
            {t("cost.payloadError")}
          </p>
          <p className="font-mono text-[11px] leading-[1.6] text-ink-3">{state.message}</p>
          <Button className="self-start" onClick={reload} size="xs" variant="secondary">
            {t("cost.retry")}
          </Button>
        </div>
      ) : null}
      {state.phase === "unreadable" ? (
        <div className="mt-2 flex flex-col gap-2 rounded-card border border-line bg-surface px-3.5 py-3">
          <p className="text-[12px] leading-[1.7] text-ink-2">{t("cost.refused")}</p>
          <p className="whitespace-pre-wrap font-mono text-[11px] leading-[1.6] text-ink-3">{state.raw}</p>
          <Button className="self-start" onClick={reload} size="xs" variant="secondary">
            {t("cost.retry")}
          </Button>
        </div>
      ) : null}
      {state.phase === "ready" ? (
        <>
          <div className="mt-2 flex items-center gap-2">
            <StateTag tone="ok" text={t("cost.gateTag")} />
          </div>
          {state.prose !== "" ? (
            <p className="mt-2 text-[12px] leading-[1.7] text-ink-2">{state.prose}</p>
          ) : null}
          <SettingsCard title={t("cost.dailyTitle")} description={t("cost.dailyDesc")}>
            <CostDayTable cachedKnown={state.cachedKnown} count={state.count} days={state.days} skipped={state.skipped} />
          </SettingsCard>
          {state.aggregateKnown && state.aggregate ? (
            <SettingsCard title={t("cost.aggregateTitle")} description={t("cost.aggregateDesc")}>
              <CostAggregateCard aggregate={state.aggregate} />
            </SettingsCard>
          ) : null}
        </>
      ) : null}

      <SettingsCard title={t("cost.gatesTitle")} description={t("cost.gatesDesc")}>
        <div className="flex flex-col">
          {ABSENCES.map((row) => (
            <div
              className="flex items-start justify-between gap-4 border-b border-line py-2.5 last:border-b-0"
              key={row.label}
            >
              <div className="min-w-0">
                <span className="text-[12.5px] font-medium text-ink">{row.label}</span>
                <p className="mt-0.5 text-[11px] leading-[1.6] text-ink-3">{row.reason}</p>
              </div>
              <StateTag tone="neutral" text={t("cost.gated")} />
            </div>
          ))}
        </div>
      </SettingsCard>
    </section>
  );
}
