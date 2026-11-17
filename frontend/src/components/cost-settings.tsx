/* ============================================================================
   成本 settings category - ticket 194's cost door gets its page reader
   (ticket 40, the census re-cut).
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
   currency), tasks. The IN/OUT labels are the wire's own dimension names
   (§17.5's micro-label rule), spelled as IN/OUT text - no chart library, no
   invented currency symbol: costMicros is shown as the integer the row
   carries, with the unit sentence, because the currency column does not ride
   THIS row.

   WHAT THE WIRE DOES NOT CARRY - each named on the page, never drawn as a
   zero (the handler's own comment bans derived totals: "a summary that
   derives a number is a second source of truth", cost_handlers.go:26-29, and
   this page obeys the same rule):
     CACHED          no cached column on CostDayRow.
     月/累计聚合       no aggregate fields on costPayload; the page does not
                     add them up itself.
     预算条           the Go-side judgment IS LIVE ([cost] thresholds at
                      internal/config/schema.go:608-613, the 80% warn / 100%
                      over verdict at internal/observe/costbudget.go:103, the
                      EvWarn exit at internal/agent/costbook_44.go:147), but
                      the cost.summary payload does not carry a budget field
                      yet - the bar waits for that field, never for a
                      client-side guess.
     计价版本脚注      no pricing-version field.
   Landing each one means a Go-side field on costPayload (cost_handlers.go:124)
   plus a page read - filed in ticket 40's report, not faked here.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { Button } from "@/components/ai-native/button";
import { SettingsCard, StateTag } from "@/components/settings-controls";
import { useCostSummary, type CostDayRowView } from "@/lib/cost-view";

/** The pure per-day table. Fed rows come back out verbatim - the behavioural
    tooth the render instrument feeds (scripts/render-security-cost-40.tsx). */
export function CostDayTable({ days, count, skipped }: { days: CostDayRowView[]; count: number; skipped: number }) {
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

/** The wire's named absences. Each row says which Go-side field it waits for,
    with the file that owns it - landing one means adding the field there and
    a page read here, in the same change, never a client-side guess. */
const ABSENCES: readonly { label: string; reason: string }[] = [
  {
    label: "CACHED（缓存命中 token）",
    reason: "CostDayRow（internal/panel/cost_handlers.go:30-36）没有 cached 列——IN/OUT 之外的字段先加在 Go 侧。",
  },
  {
    label: "月 / 累计聚合",
    reason:
      "costPayload（internal/panel/cost_handlers.go:124-127）只有逐日行；处理器自己写明不派生总数（:26-29），页面也不自己加总——聚合要 Go 侧加字段。",
  },
  {
    label: "预算条（80% 警告 / 100% 暂停）",
    reason:
      "预算判定 Go 侧已活（internal/config/schema.go:608-613 的 [cost] 阈值、internal/observe/costbudget.go:103 Judge 出 80% 警告 / 100% 超支、internal/agent/costbook_44.go:147 EvWarn），但 cost.summary payload（internal/panel/cost_handlers.go:124）尚未携带预算字段——字段落 wire 之前一格预算数字都不画，不造读数。",
  },
  {
    label: "计价版本脚注",
    reason: "pricing-version 不在 wire 上——脚注没得写，写一个猜的版本号等于假话。",
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
          <SettingsCard title="按天成本" description="cost_daily 表逐行上屏，新的一天在前；页面不派生任何表里没有的数字。">
            <CostDayTable days={state.days} count={state.count} skipped={state.skipped} />
          </SettingsCard>
        </>
      ) : null}

      <SettingsCard title="这一页还没接的面" description="每行写明它缺哪个 Go 侧字段；字段落到 wire 之前，一格数字都不画。">
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
              <StateTag tone="neutral" text="未接线" />
            </div>
          ))}
        </div>
      </SettingsCard>
    </section>
  );
}
