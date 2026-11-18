/* ============================================================================
   Session-rules view - the security page's grants doors (ticket 40 AC#2,
   readface wave 4), their parser, and the hooks the 安全 category reads
   them through.
   ----------------------------------------------------------------------------
   WHAT ARRIVES. Both doors answer through the settings door
   (window.wispDispatch) with ONE string, exactly like cost.summary:

     - list reply (cmd/wisp/panel_session_rules_40.go
       HandleSessionRulesList): one prose sentence followed by ONE line
       starting with SessionRulesPayloadMarker ("wisp-session-rules:") whose
       remainder is the door's own compact JSON: {session, rows, zero}.
       rows[] are the LIVE rows of THIS session's D45-2 ledger, each the
       store's own approval_grant fields verbatim - id, scope, tool,
       pattern, session, createdAt, expiresAt (unix seconds, carried as
       numbers, never re-shaped into dates here). The prose is for humans
       and is NEVER parsed (the same rule lib/config-view.ts and
       lib/cost-view.ts state for their markers).
       zero is the EXPLICIT empty state: "" on a non-empty ledger, the
       door's own sentence on an empty one. The page renders it verbatim -
       an empty ledger is said OUT LOUD ("还没有生效中的授权"), never drawn
       as an empty table that reads like a normal reading.
       A revoked (spent) row LEAVES the live list - what is listed is what
       is in force; the audited row stays on disk (the store's business,
       not this page's).

     - revoke receipt (HandleSessionRulesRevoke): ONE sentence naming the
       row ("已撤销会话授权 #<id>：...下一次调用会重新弹卡确认。"), verbatim.

     - a refusal: RefusedEnvelopeForUser's one sentence ("面板请求被拒绝
       [method id]：..."). The no-ledger refusal NAMES the missing engine
       (会话授权读面未接入 / 会话授权撤销通路未接入) - the page shows the
       whole sentence verbatim, never an empty list that reads as "no
       grants".

   WHAT NEVER ARRIVES gets a named state, never a drawn zero: a list reply
   without its payload line stays 未判 (no-payload), and a payload whose
   rows key is missing or whose row fields are damaged is a payload-error
   (damaged rows are counted and dropped, never folded into zeros - the
   cost-view rule). Nothing here fabricates a row or a count.

   THE HOOKS. useSessionRules sends one list on mount and exposes reload()
   (the useCostSummary shape). useSessionRulesRevoke is send-on-demand and
   takes an after() callback the caller uses to re-read the list after a
   receipt - the page keeps no local picture of the ledger to patch (the
   settings-door rule: every render comes from a reply to a request this
   page just sent).
   ============================================================================ */

import { useCallback, useEffect, useRef, useState } from "react";
import { requestSessionRules, requestSessionRulesRevoke } from "@/lib/panel";
import { REFUSED_ENVELOPE_PREFIX } from "@/lib/privacy-view";

/** The wire marker, spelled once here; cmd/wisp's SessionRulesPayloadMarker
    ("wisp-session-rules:") is the other side of the same string. */
export const SESSION_RULES_MARKER = "wisp-session-rules:";

/** One live grant row, verbatim off the payload (sessionRulesRow40's JSON
    keys). createdAt/expiresAt are the store's own unix seconds - displayed
    as the numbers they are, never re-derived into dates on this page. */
export interface SessionRuleRowView {
  id: number;
  scope: string;
  tool: string;
  pattern: string;
  session: string;
  createdAt: number;
  expiresAt: number;
}

/** The list door's wire shape, as the page reads it. zero carries the door's
    own empty-ledger sentence ("" on a non-empty ledger) - the page renders
    it verbatim and never invents an empty-table reading of its own. */
export interface SessionRulesEnvelope {
  session: string;
  rows: SessionRuleRowView[];
  zero: string;
}

/** What one session.rules.list reply turned out to be. no-host (plain
    browser), unreadable (refusal without a payload line), payload-error and
    the parsed state stay separate - a named refusal and an empty ledger are
    different facts, and rendering them alike is the lie. */
export type SessionRulesReply =
  | {
      kind: "parsed";
      prose: string;
      envelope: SessionRulesEnvelope;
      skipped: number;
    }
  | { kind: "payload-error"; message: string; raw: string }
  | { kind: "no-payload"; raw: string };

/** Parse one session.rules.list reply. Rows are validated field by field:
    a row missing a real field, or carrying a non-finite number where the
    wire said number, is SKIPPED AND COUNTED - never guessed into shape (the
    cost-view rule). zero is carried through verbatim whatever it holds. */
export function parseSessionRulesReply(reply: string): SessionRulesReply {
  const at = reply.indexOf(SESSION_RULES_MARKER);
  if (at < 0) {
    return { kind: "no-payload", raw: reply };
  }
  const afterMarker = reply.slice(at + SESSION_RULES_MARKER.length);
  const line = afterMarker.split("\n", 1)[0] ?? "";
  let data: unknown;
  try {
    data = JSON.parse(line);
  } catch {
    return { kind: "no-payload", raw: reply };
  }
  if (typeof data !== "object" || data === null) {
    return { kind: "no-payload", raw: reply };
  }
  const raw = data as { session?: unknown; rows?: unknown; zero?: unknown };
  if (typeof raw.session !== "string") {
    return {
      kind: "payload-error",
      message: "载荷上没有 session（页面拒绝猜会话）",
      raw: reply,
    };
  }
  if (!Array.isArray(raw.rows)) {
    return {
      kind: "payload-error",
      message: "载荷上 rows 不在或不是数组——行清单未判，不画 0 行",
      raw: reply,
    };
  }
  const rows: SessionRuleRowView[] = [];
  let skipped = 0;
  for (const entry of raw.rows) {
    if (typeof entry !== "object" || entry === null) {
      skipped += 1;
      continue;
    }
    const row = entry as Record<string, unknown>;
    const num = (v: unknown): number => (typeof v === "number" && Number.isFinite(v) ? v : NaN);
    const id = num(row.id);
    const createdAt = num(row.createdAt);
    const expiresAt = num(row.expiresAt);
    const tool = typeof row.tool === "string" ? row.tool : "";
    const pattern = typeof row.pattern === "string" ? row.pattern : "";
    const scope = typeof row.scope === "string" ? row.scope : "";
    const session = typeof row.session === "string" ? row.session : "";
    if (
      Number.isNaN(id) || id <= 0 || Number.isNaN(createdAt) || Number.isNaN(expiresAt) ||
      tool === "" || pattern === "" || scope === "" || session === ""
    ) {
      skipped += 1;
      continue;
    }
    rows.push({ id, scope, tool, pattern, session, createdAt, expiresAt });
  }
  return {
    kind: "parsed",
    prose: reply.slice(0, at).trimEnd(),
    envelope: {
      session: raw.session,
      rows,
      zero: typeof raw.zero === "string" ? raw.zero : "",
    },
    skipped,
  };
}

/** What the 安全 category renders from for the list. */
export type SessionRulesState =
  | { phase: "loading" }
  | { phase: "no-host" }
  | { phase: "ready"; prose: string; envelope: SessionRulesEnvelope; skipped: number; raw: string }
  | { phase: "unreadable"; raw: string }
  | { phase: "payload-error"; message: string; raw: string };

/** Read the session's live grants once, and expose reload() for a re-read. */
export function useSessionRules(): {
  state: SessionRulesState;
  reload: () => void;
} {
  const [state, setState] = useState<SessionRulesState>({ phase: "loading" });
  // A per-call token: a reload that lands after a newer one must not
  // overwrite the newer read (lib/config-view.ts's useConfigView shape).
  const seq = useRef(0);

  const reload = useCallback(() => {
    const mine = ++seq.current;
    void requestSessionRules().then((reply) => {
      if (mine !== seq.current) return;
      if (reply === null) {
        setState({ phase: "no-host" });
        return;
      }
      if (reply.startsWith(REFUSED_ENVELOPE_PREFIX)) {
        setState({ phase: "unreadable", raw: reply });
        return;
      }
      const parsed = parseSessionRulesReply(reply);
      if (parsed.kind === "parsed") {
        setState({ phase: "ready", prose: parsed.prose, envelope: parsed.envelope, skipped: parsed.skipped, raw: reply });
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

/** What one revoke reply turned out to be. idle/sending are the caller's
    phases; the rest classify the reply itself. */
export type SessionRulesRevokeState =
  | { phase: "idle" }
  | { phase: "sending" }
  | { phase: "no-host" }
  /** The receipt, verbatim - it names the row it spent. */
  | { phase: "receipt"; text: string }
  /** RefusedEnvelopeForUser's sentence, verbatim - including the refusal
      that names the missing 会话授权撤销通路. */
  | { phase: "refused"; raw: string };

/** Spend one grant on demand and classify the reply. `after` runs only when
    a receipt really landed: the caller re-reads the list there, so the page
    never patches its own picture of the ledger. A stale reply can never
    overwrite a newer one (the seq guard). */
export function useSessionRulesRevoke(after?: () => void): {
  state: SessionRulesRevokeState;
  submit: (ruleId: string) => void;
} {
  const [state, setState] = useState<SessionRulesRevokeState>({ phase: "idle" });
  const seq = useRef(0);

  const submit = useCallback(
    (ruleId: string) => {
      const mine = ++seq.current;
      setState({ phase: "sending" });
      void requestSessionRulesRevoke(ruleId).then((reply) => {
        if (mine !== seq.current) return;
        if (reply === null) {
          setState({ phase: "no-host" });
          return;
        }
        if (reply.startsWith(REFUSED_ENVELOPE_PREFIX)) {
          setState({ phase: "refused", raw: reply });
          return;
        }
        setState({ phase: "receipt", text: reply });
        after?.();
      });
    },
    [after],
  );

  return { state, submit };
}
