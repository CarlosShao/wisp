/* ============================================================================
   安全 settings category - the LIVE grants card (ticket 40 AC#2, readface
   wave 4's page half).
   ----------------------------------------------------------------------------
   WHAT IS REAL HERE. The Go half landed readface wave 4 (2026-11-18): the two
   roster doors (bridge.go's session.rules.list / session.rules.revoke) are
   routed by composer_dispatch.go, wired in panel_inbound.go, and answered by
   cmd/wisp/panel_session_rules_40.go over the runtime's OWN D45-2 session
   ledger (internal/session's Ledger - the object the answer route writes
   through, so what this card lists and what the gate honours cannot
   diverge). This card sends every read and every click through the real
   door (lib/panel.ts's two requests on the settings envelope) and shows the
   reply verbatim. There is no local ledger copy anywhere: after a revoke
   receipt the card RE-READS the list (the reload-after-receipt wiring), it
   never patches its own picture of the rows.

   THE ROWS ARE THE STORE'S OWN FIELDS. Each rendered row shows the payload's
   verbatim {id, scope, tool, pattern, session, createdAt, expiresAt} - unix
   seconds displayed as the numbers they are, no date re-derivation on this
   page, no count beyond rows.length. A damaged row is skipped and counted
   (the skipped number is stated, never folded into zeros).

   THE EMPTY LEDGER IS SAID OUT LOUD. When the payload's rows are empty, the
   door's own zero sentence is rendered verbatim and NO table chrome is
   drawn - an empty ledger reads as "还没有生效中的授权", never as an
   empty table that looks like a normal reading (tooth ②'s shape).

   THE REVOKE IS ONE CLICK, PER ROW. SPEC-06 §8.2's 一键撤销: revoking
   REMOVES an authorization - the worst outcome is the gate asking again,
   the safe direction - so there is no confirm card, and the wire call
   carries no verdict vocabulary anywhere. The button sends sessionRuleId =
   the row's own id echoed verbatim; a refused revoke shows the refusal's
   own sentence (including the one that names a missing 会话授权撤销通路).

   NAMING, stated once where the user can read it: the doors ride the
   synonym spellings session.rules.* because the frozen grant-boundary
   instrument bans the spec's grants.list / grants.revoke on the wire; the
   semantic target is unchanged (one approval_grant row) and the console
   verbs grants / revoke-grant remain as the operator's second exit.

   Colours: token utilities only, zero literals; no emoji.
   ============================================================================ */

import { Button } from "@/components/ai-native/button";
import { SettingsCard } from "@/components/settings-controls";
import {
  useSessionRules,
  useSessionRulesRevoke,
  type SessionRuleRowView,
  type SessionRulesRevokeState,
  type SessionRulesState,
} from "@/lib/session-rules-view";

/** One grant row as the payload carried it: every field verbatim, plus the
    row's own revoke button (one click, the id echoed, no second judge). */
function SessionRuleRow({ row, onRevoke, revokeBusy }: {
  row: SessionRuleRowView;
  onRevoke: (id: string) => void;
  revokeBusy: boolean;
}) {
  return (
    <div className="flex items-start justify-between gap-4 border-b border-line py-2.5 last:border-b-0">
      <div className="min-w-0">
        <p className="font-mono text-[11.5px] leading-[1.6] text-ink">
          #{row.id}　{row.tool}「{row.pattern}」
        </p>
        <p className="mt-0.5 font-mono text-[11px] leading-[1.6] text-ink-3">
          scope={row.scope}　session={row.session}　created={row.createdAt}　expires={row.expiresAt}
        </p>
      </div>
      <Button
        aria-label={`撤销会话授权 #${row.id}`}
        disabled={revokeBusy}
        onClick={() => onRevoke(String(row.id))}
        size="xs"
        variant="secondary"
      >
        撤销
      </Button>
    </div>
  );
}

/** The list reply area: whatever the door answered, rendered the way the
    wire classified it. This is the behavioural feed point the render
    instrument drives (scripts/render-sessionrules-40.tsx): fed real rows it
    shows them field by field; fed an empty ledger it shows the door's own
    zero sentence verbatim and NO table; fed a refusal or a damaged payload
    it says so instead of drawing a fabricated reading. onRevoke/revokeBusy
    are the card's live revoke wiring; the defaults keep the block renderable
    standalone (the instrument's feed shapes) while the CARD's rows always
    carry the real handler - a dead revoke button would be a dormant shape,
    the thing the census called 假免接线. */
export function SessionRulesReplyBlock({ state, onRevoke, revokeBusy = false }: {
  state: SessionRulesState;
  /** The card's own revoke submit: the row's id echoed verbatim onto the
      session.rules.revoke wire. */
  onRevoke?: (id: string) => void;
  revokeBusy?: boolean;
}) {
  if (state.phase === "loading") {
    return <p className="text-[12px] leading-[1.7] text-ink-3">正在读取本会话的授权...</p>;
  }
  if (state.phase === "no-host") {
    return (
      <div className="rounded-card border border-line bg-surface px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-ink-2">
          没接到本机宿主（这一页没有 wispDispatch 门）：请求没有发出去，读数未判——这里不画任何假清单。
        </p>
      </div>
    );
  }
  if (state.phase === "unreadable") {
    return (
      <div className="rounded-card border border-red/40 bg-red-tint px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-red">门拒绝了这个读取，原话逐字如下：</p>
        <p className="mt-1.5 whitespace-pre-wrap break-all font-mono text-[11px] leading-[1.6] text-ink">
          {state.raw}
        </p>
      </div>
    );
  }
  if (state.phase === "payload-error") {
    return (
      <div className="rounded-card border border-orange/40 bg-orange-tint px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-orange">{state.message}</p>
        <p className="mt-1.5 max-h-32 overflow-y-auto whitespace-pre-wrap break-all font-mono text-[11px] leading-[1.6] text-ink">
          {state.raw}
        </p>
      </div>
    );
  }
  // phase === "ready". THE ZERO BRANCH IS THE DOOR'S OWN SENTENCE: rendered
  // verbatim, no row chrome, no fabricated "0 行" table. The rows branch is
  // the payload's own rows, verbatim, with the skipped count stated when the
  // wire delivered damage.
  if (state.envelope.rows.length === 0) {
    return (
      <div className="rounded-card border border-line bg-surface px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-ink">{state.envelope.zero}</p>
        <p className="mt-1 text-[11px] leading-[1.6] text-ink-3">
          这是门自报的空账读数（zero 字段逐字），不是本页自己拼的空表。
        </p>
      </div>
    );
  }
  return (
    <div className="rounded-card border border-line bg-surface px-3.5 py-3">
      <p className="text-[12px] leading-[1.7] text-ink">
        生效中的授权 {state.envelope.rows.length} 行（读自载荷 rows 数组的长度）。
      </p>
      {state.skipped > 0 ? (
        <p className="mt-1 text-[11.5px] leading-[1.6] text-orange">
          另有 {state.skipped} 行字段残缺被丢弃——计数丢弃，不补零、不猜字段。
        </p>
      ) : null}
      <div className="mt-2">
        {state.envelope.rows.map((row) => (
          <SessionRuleRow
            key={row.id}
            revokeBusy={revokeBusy}
            row={row}
            onRevoke={onRevoke ?? (() => undefined)}
          />
        ))}
      </div>
    </div>
  );
}

/** The revoke reply area, verbatim per classification. */
export function SessionRulesRevokeReplyBlock({ state }: { state: SessionRulesRevokeState }) {
  if (state.phase === "idle") return null;
  if (state.phase === "sending") {
    return <p className="text-[12px] leading-[1.7] text-ink-3">撤销请求已在桥上，等门回话...</p>;
  }
  if (state.phase === "no-host") {
    return (
      <div className="rounded-card border border-line bg-surface px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-ink-2">
          没接到本机宿主：撤销请求没有发出去，这行授权未被触碰——回执未判，不画假成功。
        </p>
      </div>
    );
  }
  if (state.phase === "refused") {
    return (
      <div className="rounded-card border border-red/40 bg-red-tint px-3.5 py-3">
        <p className="text-[12px] leading-[1.7] text-red">门拒绝了这个撤销，原话逐字如下：</p>
        <p className="mt-1.5 whitespace-pre-wrap break-all font-mono text-[11px] leading-[1.6] text-ink">
          {state.raw}
        </p>
      </div>
    );
  }
  return (
    <div className="rounded-card border border-green/40 bg-green-tint px-3.5 py-3">
      <p className="text-[12px] leading-[1.7] text-ink">撤销回执逐字：</p>
      <p className="mt-1 whitespace-pre-wrap font-mono text-[11.5px] leading-[1.6] text-ink">{state.text}</p>
    </div>
  );
}

/** The LIVE grants card: the read face and the per-row one-click revoke,
    each through its real door. after() re-reads the list when a receipt
    lands - the card never keeps a patched copy of the ledger. */
export function SessionRulesCard() {
  const list = useSessionRules();
  const revoke = useSessionRulesRevoke(list.reload);

  const onRevoke = (id: string) => {
    revoke.submit(id);
  };

  return (
    <SettingsCard
      title="本会话生效中的授权"
      description="session.rules.list：实时列出本会话（D45-2）生效中的授权，行内字段逐字来自账本（生产写入面 gate.go:859 g.grants.Record / gate.go:844 RecordTarget，落库 grants.go:201/274 InsertGrant，id 回带 bridge.go:1445）；空账时门自报的 zero 句逐字上屏，本页不画空表。撤销键 = session.rules.revoke（一键，撤的是授权本身：下一次同工具调用会重新弹卡；行保留在审计里，生效清单不再列出它）。路由名走 session.rules.* 而非规范里的 grants.list/grants.revoke：冻结仪器禁「grant」上桥，语义目标不变（一行 approval_grant）；console 动词 grants / revoke-grant 仍是操作员的第二个出口。"
    >
      <SessionRulesReplyBlock state={list.state} onRevoke={onRevoke} revokeBusy={revoke.state.phase === "sending"} />
      <div className="mt-2">
        <SessionRulesRevokeReplyBlock state={revoke.state} />
      </div>
      <p className="mt-2 text-[11px] leading-[1.6] text-ink-3">
        读数与撤销走同一本账（运行时的会话账本）：这里列出的，就是门控在问「已授权过吗」时读的那几行；撤销落库生效，下一次调用重新弹卡。
      </p>
    </SettingsCard>
  );
}
