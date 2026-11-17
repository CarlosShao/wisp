/* ============================================================================
   安全 settings category - the security page's wired half plus its named
   absences (ticket 40, the census re-cut).
   ----------------------------------------------------------------------------
   WHAT IS REAL HERE. The permission tier is pushed EVERY packet
   (cmd/wisp/run.go:1187 -> internal/panel/pump.go:372-373 -> the composer's
   mode section): the census called the old "还没收到策略快照" sentence a
   假免接线 for exactly this reason. So this page reads composer.mode and
   nothing else:

     current         the mode the live process runs RIGHT NOW; "unknown" is
                     the wire's own marker for "this pump did not read the
                     mode" and renders as 档位未知, never as the
                     safest-sounding default (pump.go AC#4's rule).
     names           risk.ModeNames() verbatim - the vocabulary, displayed as
                     facts, not a choice set this page offers.
     l2ConfirmNames  the transitions that cost an L2 strong confirmation
                     (R20/M4), named by the native side itself.

   Changing the tier happens on the composer row (panel.mode.request, the
   native L2 confirmation); this page deliberately mounts NO switch - a
   settings page that could change the security posture outside the composer's
   confirmation flow would be a second decision surface, which ticket 92
   forbids. The page says so instead of drawing a control.

   WHAT IS NOT REAL HERE (each an absence row naming its missing engine, per
   the census rule - a dormant row names what it waits for, it never draws a
   value):

     门控阈值        thresholds are the native frozen face
                     (internal/risk/thresholds.go); no snapshot section, no
                     C17 route (internal/panel/bridge.go's roster).
     会话授权/撤销    D45 grants production writers are live - the census's
                     journal.go:101 was the census's claim; ticket 49 proved
                     production writes ARE live: gate.go:859
                     g.grants.Record / gate.go:844 rt.RecordTarget, landing
                     through internal/session/grants.go:201/274 InsertGrant,
                     with the id carried back to the call face at
                     tools/bridge.go:1445 (tc.GrantID = &g); session close
                     stamps revoked_at and keeps the rows for audit
                     (grants.go:346 RevokeAll). The absence is PANEL-side
                     only: the snapshot carries no grants section and the
                     architecture bans Go-inbound grant routes
                     (grants.list/grants.revoke are not in C17's roster), so
                     there is no grants list for this page to read. The
                     operating exit is the console verbs: grants (list this
                     session's live grants) and revoke-grant <编号>
                     (cmd/wisp/approval_reply.go:415/:474).
     B-list 覆盖行    no snapshot section.
     taint-hit 历史  the tool_call table exists (internal/memory), but no
                     panel read face names it.
     审批审计浏览     the [audit] lines exist, but there is no read route that
                     joins them by correlationId.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { SettingsCard, StateTag } from "@/components/settings-controls";
import { tierName } from "@/lib/tier-names";
import type { ComposerMode } from "@/lib/panel";

/** One absence row: the surface ticket 40 names, and the engine it waits
    for. `surface` is the page's own vocabulary (未接线), never a tier word. */
const ABSENCES: readonly { label: string; reason: string }[] = [
  {
    label: "门控阈值",
    reason:
      "阈值表是原生冻结面（internal/risk/thresholds.go）：快照没有对应节，C17 名册（internal/panel/bridge.go:41-107）也没有读取路由。",
  },
  {
    label: "会话授权面板读面",
    reason:
      "会话授权生产在写（gate.go:859 g.grants.Record / gate.go:844 RecordTarget，落库走 internal/session/grants.go:201/274 InsertGrant，tools/bridge.go:1445 把 grant_id 带回调用面；会话收口盖 revoked_at、行保留审计，grants.go:346）。面板侧没有 grants 列表读面：快照无 grants 节，架构也禁 Go 入站 grant 路由（grants.list/grants.revoke 不在 C17 名册）。操作出口＝console 动词 grants（列本会话生效授权）/ revoke-grant <编号>（cmd/wisp/approval_reply.go:415/:474）。",
  },
  {
    label: "B-list 覆盖与审计行",
    reason: "黑名单覆盖清单没有快照节，审计行没有读取路由——两面都缺 Go 侧读面。",
  },
  {
    label: "taint-hit（R4）历史",
    reason:
      "tool_call 表在（internal/memory），但没有面板读面把按来源命名的污染命中报上来。",
  },
  {
    label: "审批决定审计浏览",
    reason: "审计行在，但没有任何路由按 correlationId 关联后摊给人看。",
  },
];

export function SecuritySettings({
  mode,
  heading = "安全",
  lede = "权限档位、门控阈值、批量授权的边界。",
}: {
  /** The snapshot's composer.mode section - the ONLY security fact this page
      holds. The tier rides every packet, so this is a read, not a remember. */
  mode: ComposerMode;
  heading?: string;
  lede?: string;
}) {
  // The vocabulary is risk.ModeNames() on the wire; an empty list means this
  // pump's mode section arrived without one (the hand-assembled snapshot
  // shape) - a named absence, not a one-word vocabulary.
  const vocabulary =
    mode.names.length > 0
      ? mode.names.map((name) => `${tierName(name)}（${name}）`).join("、")
      : "";
  const l2Transitions =
    mode.l2ConfirmNames.length > 0
      ? mode.l2ConfirmNames.map((name) => tierName(name)).join("、")
      : "";

  return (
    <section aria-label="安全" className="flex min-w-0 flex-col gap-1">
      <h2 className="text-[15px] font-semibold text-ink">{heading}</h2>
      <p className="text-[12.5px] text-ink-2">{lede}</p>

      <div className="mt-2 flex flex-col gap-3">
        <SettingsCard title="当前权限档位" description="这一泵正在跑的档位，逐包随快照推上来。切换去底部输入行的档位控件：面板这边不画第二个切换面。">
          <div className="flex flex-wrap items-center gap-2">
            <span className="text-[13px] font-medium text-ink">{tierName(mode.current)}</span>
            {mode.current === "unknown" ? (
              <StateTag tone="neutral" text="本泵没读到档位读数" />
            ) : (
              <StateTag tone="ok" text="读自快照 composer.mode" />
            )}
          </div>
          <p className="mt-2 text-[11.5px] leading-[1.6] text-ink-3">
            切档的确认走原生侧：切入需要 L2 强确认的档位时，原生先弹卡再落盘（R20/M4），页面只发起请求、永远不记录结果。
          </p>
        </SettingsCard>

        <SettingsCard title="档位词表" description="原生侧 risk.ModeNames() 的全部档位，逐字上屏；词表是显示事实，不是本页可选项。">
          {vocabulary === "" ? (
            <p className="text-[12px] leading-[1.7] text-ink-2">
              这一泵的 mode 节没带词表（names 为空）——等下一包快照，不编一个词表出来。
            </p>
          ) : (
            <p className="text-[12px] leading-[1.7] text-ink-2">{vocabulary}</p>
          )}
        </SettingsCard>

        <SettingsCard title="切入需 L2 强确认的档位" description="原生侧自己点名的强确认名单（l2ConfirmNames），逐字上屏。">
          {l2Transitions === "" ? (
            <p className="text-[12px] leading-[1.7] text-ink-2">
              这一泵没有点名任何一档——这是读到的事实（空名单），不是「不知道」。
            </p>
          ) : (
            <p className="text-[12px] leading-[1.7] text-ink-2">{l2Transitions}</p>
          )}
        </SettingsCard>

        <SettingsCard title="这一页还没接的面" description="每行写明它缺哪台引擎；画一个没有读口的控件，等于让改不动的旋钮看起来像能改。">
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
      </div>
    </section>
  );
}
