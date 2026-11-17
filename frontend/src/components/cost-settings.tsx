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
        <span className="text-[11.5px] text-ink-3">IN = 输入 token，OUT = 输出 token，金额 = costMicros（1e-6 货币单位，整数逐字；币种列不在这一行上）</span>
      </div>
      {days.length === 0 ? (
        <p className="text-[12px] leading-[1.7] text-ink-2">
          成本表是空的——这是读到的表事实（这一天没有任务记账），不是读不到。
        </p>
      ) : (
        <table className="w-full border-collapse text-left">
          <thead>
            <tr className="border-b border-line">
              <th className="py-1.5 pr-3 text-[11px] font-medium text-ink-3">日期</th>
              <th className="py-1.5 pr-3 text-[11px] font-medium text-ink-3">IN</th>
              <th className="py-1.5 pr-3 text-[11px] font-medium text-ink-3">OUT</th>
              <th className="py-1.5 pr-3 text-[11px] font-medium text-ink-3">CACHED</th>
              <th className="py-1.5 pr-3 text-[11px] font-medium text-ink-3">金额（costMicros）</th>
              <th className="py-1.5 text-[11px] font-medium text-ink-3">任务</th>
            </tr>
          </thead>
          <tbody>
            {days.map((row) => (
              <tr className="border-b border-line last:border-b-0" key={row.day}>
                <td className="py-1.5 pr-3 font-mono text-[12px] text-ink">{row.day}</td>
                <td className="py-1.5 pr-3 font-mono text-[12px] tabular-nums text-ink-2">{row.tokensIn}</td>
                <td className="py-1.5 pr-3 font-mono text-[12px] tabular-nums text-ink-2">{row.tokensOut}</td>
                <td className="py-1.5 pr-3 font-mono text-[12px] tabular-nums text-ink-2">
                  {cachedKnown ? row.cachedTokens : "未知"}
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
          有 {skipped} 行字段残缺，没有上屏——宁可少一行，不补零。
        </p>
      ) : null}
      <p className="text-[11px] leading-[1.6] text-ink-3">共 {count} 行（新的一天在前，行序是读口自己的）。</p>
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
    { key: "normal", label: "未满 80%（normal）" },
    { key: "warn", label: "80% 及以上（warn）" },
    { key: "paused", label: "100% 及以上（paused）" },
  ];
  return (
    <div className="flex flex-col gap-1.5">
      <div aria-label={`预算状态条：${state}`} className="flex flex-wrap items-center gap-1.5" role="img">
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
        状态词读自 aggregate.budgetState（逐字：{state}），三态判定在原生侧（internal/observe 的整数基点判定），页面不复算。
        {active === null ? " 这个词不在三态词表里：条不定格，词照原样显示。" : ""}
        新任务准入读 aggregate.budgetAllowNewTask（逐字）：
        {allowNewTask ? "允许" : "停止"}。
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
    micros === 0 ? "不封顶（0 = 未设预算）" : `${micros} micro`;
  return (
    <div className="flex flex-col gap-2">
      <table className="w-full border-collapse text-left">
        <tbody>
          <tr className="border-b border-line">
            <td className="py-1.5 pr-3 text-[11px] text-ink-3">月累计 IN</td>
            <td className="py-1.5 font-mono text-[12px] tabular-nums text-ink">{aggregate.monthTokensIn}</td>
          </tr>
          <tr className="border-b border-line">
            <td className="py-1.5 pr-3 text-[11px] text-ink-3">月累计 OUT</td>
            <td className="py-1.5 font-mono text-[12px] tabular-nums text-ink">{aggregate.monthTokensOut}</td>
          </tr>
          <tr className="border-b border-line">
            <td className="py-1.5 pr-3 text-[11px] text-ink-3">月累计金额（micro，整数逐字）</td>
            <td className="py-1.5 font-mono text-[12px] tabular-nums text-ink">{aggregate.monthCostMicros}</td>
          </tr>
          <tr className="border-b border-line">
            <td className="py-1.5 pr-3 text-[11px] text-ink-3">今日金额（micro，整数逐字）</td>
            <td className="py-1.5 font-mono text-[12px] tabular-nums text-ink">{aggregate.dayCostMicros}</td>
          </tr>
          <tr className="border-b border-line">
            <td className="py-1.5 pr-3 text-[11px] text-ink-3">日预算</td>
            <td className="py-1.5 font-mono text-[12px] tabular-nums text-ink">{budgetWord(aggregate.dailyBudgetMicros)}</td>
          </tr>
          <tr>
            <td className="py-1.5 pr-3 text-[11px] text-ink-3">月预算</td>
            <td className="py-1.5 font-mono text-[12px] tabular-nums text-ink">{budgetWord(aggregate.monthlyBudgetMicros)}</td>
          </tr>
        </tbody>
      </table>
      <BudgetStateBar allowNewTask={aggregate.budgetAllowNewTask} state={aggregate.budgetState} />
      <p className="text-[11px] leading-[1.6] text-ink-3">
        {aggregate.pricingVersion === 0 || aggregate.pricingUpdated === ""
          ? "计价版本未知——读口没能读到价表（0/空 = 未知，页面不编一个版本号）。"
          : `计价版本 ${aggregate.pricingVersion}（${aggregate.pricingUpdated} 更新），读自 aggregate.pricingVersion/pricingUpdated 逐字。`}
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
const ABSENCES: readonly { label: string; reason: string }[] = [
  {
    label: "CACHED（缓存命中 token）",
    reason:
      "本 build 的 cost_daily 没有 cached 列（D35 表，读面腿不许加宽，cost_handlers.go:36-45）⇒ payload 恒带 cachedKnown=false，缓存列显示「未知」，绝不把结构零画成读数。cachedKnown 翻 true（存储开始聚合）时数字才上屏。",
  },
  {
    label: "月累计 / 预算 / 计价版本（聚合面未取到时）",
    reason:
      "聚合面是可选的第二读面（CostAggregateReader，cost_handlers.go:92；CostAggregateRow 在 :54-83）：聚合面没接或读取失败时 payload 带 aggregateKnown=false，prose 也写明「未取到」——这一格一格数字都不画，页面也不自己加总（:26-29 禁派生总数）。聚合面接上时逐字段上屏，预算三态读 state 词、不复算。",
  },
];

export function CostSettings({
  heading = "成本",
  lede = "每轮花费、本月累计、预算与超支提醒。",
}: {
  heading?: string;
  lede?: string;
}) {
  const { state, reload } = useCostSummary();

  return (
    <section aria-label="成本" className="flex min-w-0 flex-col gap-1">
      <h2 className="text-[15px] font-semibold text-ink">{heading}</h2>
      <p className="text-[12.5px] text-ink-2">{lede}</p>

      <div className="mt-2 flex items-center gap-2">
        <Button aria-label="刷新成本" onClick={reload} size="xs" variant="quiet">
          刷新
        </Button>
      </div>

      {state.phase === "loading" ? <p className="mt-2 text-[12px] text-ink-3">正在读取成本...</p> : null}
      {state.phase === "no-host" ? (
        <div className="mt-2 rounded-card border border-line bg-surface px-3.5 py-3">
          <p className="text-[12px] leading-[1.7] text-ink-2">
            暂未连接到本机服务，读不到成本（cost.summary 门在本机侧，浏览器页没有宿主）。
          </p>
        </div>
      ) : null}
      {state.phase === "payload-error" ? (
        <div className="mt-2 flex flex-col gap-2 rounded-card border border-line bg-surface px-3.5 py-3">
          <p className="text-[12px] leading-[1.7] text-ink-2">
            成本读口的载荷没排出来（payloadError），原话如下：
          </p>
          <p className="font-mono text-[11px] leading-[1.6] text-ink-3">{state.message}</p>
          <Button className="self-start" onClick={reload} size="xs" variant="secondary">
            重试
          </Button>
        </div>
      ) : null}
      {state.phase === "unreadable" ? (
        <div className="mt-2 flex flex-col gap-2 rounded-card border border-line bg-surface px-3.5 py-3">
          <p className="text-[12px] leading-[1.7] text-ink-2">读口拒绝了这个请求，原话逐字如下（没有载荷行）：</p>
          <p className="whitespace-pre-wrap font-mono text-[11px] leading-[1.6] text-ink-3">{state.raw}</p>
          <Button className="self-start" onClick={reload} size="xs" variant="secondary">
            重试
          </Button>
        </div>
      ) : null}
      {state.phase === "ready" ? (
        <>
          <div className="mt-2 flex items-center gap-2">
            <StateTag tone="ok" text="读自 cost.summary（票 194 的只读门）" />
          </div>
          {state.prose !== "" ? (
            <p className="mt-2 text-[12px] leading-[1.7] text-ink-2">{state.prose}</p>
          ) : null}
          <SettingsCard title="按天成本" description="cost_daily 表逐行上屏，新的一天在前；页面不派生任何表里没有的数字。CACHED 列听 cachedKnown 门的：false 显示未知，不画结构零。">
            <CostDayTable cachedKnown={state.cachedKnown} count={state.count} days={state.days} skipped={state.skipped} />
          </SettingsCard>
          {state.aggregateKnown && state.aggregate ? (
            <SettingsCard title="本月累计与预算" description="CostAggregateRow 逐字段上屏（月累计/预算/三态/计价版本）；三态条只听 state 词，页面不复算百分比。">
              <CostAggregateCard aggregate={state.aggregate} />
            </SettingsCard>
          ) : null}
        </>
      ) : null}

      <SettingsCard title="这扇门的门与缺" description="每行写明 wire 上的门（cachedKnown）与可选面（aggregateKnown）的语义；门没开就不画数字。">
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
              <StateTag tone="neutral" text="门控" />
            </div>
          ))}
        </div>
      </SettingsCard>
    </section>
  );
}
