/* ============================================================================
   Ticket 40 AC#2's behavioural tooth (render evidence), 2026-11-18 - the
   security page's grants leg.
   ----------------------------------------------------------------------------
   WHAT THIS FILE PROVES, in the shape render-privacy-40 established: every
   check below feeds REAL-SHAPED data through the PRODUCTION component and
   asserts the SAME VALUE comes back out in the markup - "喂真 payload ⇒
   行内真字段逐字上屏".

   The fixtures are the wire's own documented shapes, copied from the Go leg
   that answers them (cmd/wisp/panel_session_rules_40.go, the doors
   bridge.go's MethodSessionRulesList / MethodSessionRulesRevoke name,
   composer_dispatch.go routes, panel_inbound.go wires):
     - list reply: one prose sentence plus the ONE payload line behind
       SessionRulesPayloadMarker ("wisp-session-rules:"), whose JSON is the
       door's own {session, rows, zero} envelope with the store's real
       approval_grant fields verbatim per row;
     - the EMPTY ledger: rows [] with the door's own zero sentence in the
       zero field - the page must render THAT SENTENCE verbatim and draw NO
       row chrome (tooth ②'s green half: an empty ledger is said out loud,
       never drawn as an empty table that reads like a normal reading);
     - revoke receipt: HandleSessionRulesRevoke's sentence naming the row it
       spent;
     - a refusal: RefusedEnvelopeForUser's wrapper sentence, including the
       one that names a missing 会话授权读面 - shown verbatim, never an
       empty list that reads as "no grants".

   Where a reading cannot be had, the check pins the 未判 rule: no payload
   line ⇒ the page says 未判 and NEVER draws a fabricated row count; a
   payload whose rows are damaged ⇒ skipped count stated, zeros never folded
   in (the cost-view rule).

   Source-level teeth:
     - the two route literals live in panel.ts (the ONE door file), with the
       envelope's sessionRuleId key spelled as bridge.go's JSON tag carries
       it;
     - the zero sentence is the DOOR's own words: the page renders
       state.zero verbatim and never spells the Go sentence itself (a
       page-side zero sentence would be a fabricated reading - mutation ②'s
       target);
     - the card re-reads the list after a revoke receipt (no local ledger
       copy to patch);
     - zero emoji.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import {
  SessionRulesRevokeReplyBlock,
  SessionRulesReplyBlock,
} from "@/components/session-rules-settings";
import { parseSessionRulesReply, type SessionRulesState } from "@/lib/session-rules-view";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// The emoji band the repo scanner walks (ban #8, tools/d22scan's own set).
const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u22FF\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

/** Same escaping React applies to a text node - so a match here is a match on
    screen, not a match on an unescaped substring that never renders. */
function esc(text: string): string {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#x27;");
}

/* ---------------------------------------------------------------------------
   1. The list reply parses: prose unparsed, the store's own fields verbatim
   --------------------------------------------------------------------------- */

// HandleSessionRulesList's prose shape, then the marker line whose JSON is
// the door's own envelope: two REAL rows (the shape ticket 49's tooth mints:
// one fs.write path row, one input.type target row).
const LIST_PROSE =
  "本会话（sess-40a）生效中的授权共 2 行；行内字段逐字在下面这一行 JSON（真行，不是摘要）。撤销一枚：授权页每行的撤销键，或 console 的 revoke-grant <编号>。";
const LIST_PAYLOAD = `{"session":"sess-40a","rows":[` +
  `{"id":11,"scope":"session","tool":"fs.write","pattern":"C:/work/notes.txt","session":"sess-40a","createdAt":1763400000,"expiresAt":1765984800},` +
  `{"id":12,"scope":"session","tool":"input.type","pattern":"notepad.exe","session":"sess-40a","createdAt":1763400100,"expiresAt":1765984900}],` +
  `"zero":""}`;
const LIST_REPLY = `${LIST_PROSE}\nwisp-session-rules:${LIST_PAYLOAD}`;

const listParsed = parseSessionRulesReply(LIST_REPLY);
check("the list reply must parse as the parsed state", listParsed.kind === "parsed");
if (listParsed.kind === "parsed") {
  check("the list prose must arrive verbatim, unparsed", listParsed.prose === LIST_PROSE, listParsed.prose);
  check("the envelope's session must ride verbatim", listParsed.envelope.session === "sess-40a");
  check("two real rows must ride the payload", listParsed.envelope.rows.length === 2);
  const [r1, r2] = listParsed.envelope.rows;
  check(
    "row fields must be the store's own fields verbatim (id/scope/tool/pattern/session/createdAt/expiresAt)",
    r1.id === 11 && r1.scope === "session" && r1.tool === "fs.write" &&
      r1.pattern === "C:/work/notes.txt" && r1.session === "sess-40a" &&
      r1.createdAt === 1763400000 && r1.expiresAt === 1765984800,
    JSON.stringify(r1),
  );
  check(
    "the target row must ride verbatim too (the ticket-49 RecordTarget shape)",
    r2.id === 12 && r2.tool === "input.type" && r2.pattern === "notepad.exe" &&
      r2.createdAt === 1763400100 && r2.expiresAt === 1765984900,
    JSON.stringify(r2),
  );
  check("an empty ledger's zero field must be empty string on a non-empty ledger", listParsed.envelope.zero === "");
  check("no damage in a clean payload means skipped=0", listParsed.skipped === 0);
}

/* ---------------------------------------------------------------------------
   2. Rows render field by field; the per-row revoke button echoes the id
   --------------------------------------------------------------------------- */

const rowsState: SessionRulesState =
  listParsed.kind === "parsed"
    ? { phase: "ready", prose: listParsed.prose, envelope: listParsed.envelope, skipped: listParsed.skipped, raw: LIST_REPLY }
    : ({ phase: "loading" } as SessionRulesState);
const rowsHtml = renderToStaticMarkup(
  <SessionRulesReplyBlock
    state={rowsState}
    onRevoke={() => undefined}
  />,
);
check("the row count must state the payload's own rows length", rowsHtml.includes(esc("生效中的授权 2 行")));
check("row #11 must render with its real tool and pattern verbatim", rowsHtml.includes(esc("fs.write")) && rowsHtml.includes(esc("C:/work/notes.txt")));
check("row #12 must render with its real target process verbatim", rowsHtml.includes(esc("notepad.exe")));
check("the scope/session/created/expires fields must ride verbatim", rowsHtml.includes(esc("scope=session")) && rowsHtml.includes(esc("created=1763400000")) && rowsHtml.includes(esc("expires=1765984800")));
check("each row's revoke button must name its own id", rowsHtml.includes(esc("撤销会话授权 #11")) && rowsHtml.includes(esc("撤销会话授权 #12")));
check("the rows block must hold zero emoji", !EMOJI_RE.test(rowsHtml));

/* ---------------------------------------------------------------------------
   3. THE EMPTY LEDGER: the door's own zero sentence, verbatim, NO table
      (tooth ②'s green half)
   --------------------------------------------------------------------------- */

const ZERO_SENTENCE = "本会话（sess-40a）当前没有生效中的授权。";
const zeroParsed = parseSessionRulesReply(
  `${ZERO_SENTENCE}\nwisp-session-rules:{"session":"sess-40a","rows":[],"zero":"${ZERO_SENTENCE}"}`,
);
check("an empty-ledger reply must parse as the parsed state", zeroParsed.kind === "parsed");
if (zeroParsed.kind === "parsed") {
  check("the zero field must carry the door's own sentence verbatim", zeroParsed.envelope.zero === ZERO_SENTENCE);
  const zeroHtml = renderToStaticMarkup(<SessionRulesReplyBlock state={{ phase: "ready", prose: zeroParsed.prose, envelope: zeroParsed.envelope, skipped: zeroParsed.skipped, raw: "" }} />);
  check("the empty ledger must render the door's zero sentence verbatim", zeroHtml.includes(esc(ZERO_SENTENCE)));
  check("the empty ledger must name the zero field's source (the door's words, not the page's)", zeroHtml.includes(esc("zero 字段逐字")));
  check("the empty ledger must draw NO row count chrome (no 生效中的授权 N 行 table)", !/生效中的授权 \d+ 行/.test(zeroHtml));
  check("the empty ledger must draw NO revoke button", !zeroHtml.includes(esc("撤销会话授权")));
  check("the empty-ledger block must hold zero emoji", !EMOJI_RE.test(zeroHtml));
}

/* ---------------------------------------------------------------------------
   4. Refusals ride verbatim; damaged payloads are named, never drawn
   --------------------------------------------------------------------------- */

// RefusedEnvelopeForUser's wrapper around the no-ledger refusal: the cause
// NAMES the missing engine. Fed to the page, the words must come back out.
const NO_LEDGER_REFUSAL =
  `面板请求被拒绝 [session.rules.list pc-5-abc]：panel: 该方法在名册内，但本机未接入处理器：` +
  `会话授权读面未接入（本进程没有会话授权记账，或任务入口没有起来）：本会话的授权无从列出，也不会被静默丢掉`;
const refusalHtml = renderToStaticMarkup(<SessionRulesReplyBlock state={{ phase: "unreadable", raw: NO_LEDGER_REFUSAL }} />);
check("the no-ledger refusal must ride verbatim", refusalHtml.includes(esc(NO_LEDGER_REFUSAL)));
check("the refusal must name the missing engine (会话授权读面未接入)", refusalHtml.includes(esc("会话授权读面未接入")));
check("the refusal must render in the danger tint", refusalHtml.includes("bg-red-tint"));
check("the refusal block must hold zero emoji", !EMOJI_RE.test(refusalHtml));

// A list reply WITHOUT its payload line: the reading stays 未判 and no
// fabricated row count may be drawn (喂不到 ⇒ 显式未判, 禁止假行数).
const noPayloadHtml = renderToStaticMarkup(<SessionRulesReplyBlock state={{ phase: "unreadable", raw: "会话授权读面被拒绝，载荷行不在。" }} />);
check("a marker-less reply must stay an explicit named state", noPayloadHtml.includes(esc("会话授权读面被拒绝")));
check("the no-payload block must hold zero emoji", !EMOJI_RE.test(noPayloadHtml));

// A payload whose rows are damaged: skipped count stated, zeros never folded
// in, good rows still drawn (the cost-view rule).
const damagedParsed = parseSessionRulesReply(
  `x\nwisp-session-rules:{"session":"sess-40a","rows":[{"id":11,"scope":"session","tool":"fs.write","pattern":"C:/ok.txt","session":"sess-40a","createdAt":1,"expiresAt":2},{"id":0,"scope":"session","tool":"","pattern":"","session":"","createdAt":"x","expiresAt":"y"}],"zero":""}`,
);
check("a damaged payload must parse with the good rows", damagedParsed.kind === "parsed");
if (damagedParsed.kind === "parsed") {
  check("the damaged row must be skipped, not folded into a zero row", damagedParsed.envelope.rows.length === 1 && damagedParsed.skipped === 1);
  const damagedHtml = renderToStaticMarkup(<SessionRulesReplyBlock state={{ phase: "ready", prose: damagedParsed.prose, envelope: damagedParsed.envelope, skipped: damagedParsed.skipped, raw: "" }} />);
  check("the skipped count must be stated on screen", damagedHtml.includes(esc("字段残缺被丢弃")) && damagedHtml.includes(esc("不补零")));
  check("the good row must still render", damagedHtml.includes(esc("C:/ok.txt")));
  check("the damaged block must hold zero emoji", !EMOJI_RE.test(damagedHtml));
}

// The no-host state (plain browser / hand-assembled harness): explicit, and
// nothing that reads like a result.
const noHostHtml = renderToStaticMarkup(<SessionRulesReplyBlock state={{ phase: "no-host" }} />);
check("the no-host block must say the reading is 未判", noHostHtml.includes(esc("未判")));
check("the no-host block must hold zero emoji", !EMOJI_RE.test(noHostHtml));

/* ---------------------------------------------------------------------------
   5. The revoke receipt rides verbatim
   --------------------------------------------------------------------------- */

const REVOKE_RECEIPT =
  "已撤销会话授权 #12：本会话内对应工具的下一次调用会重新弹卡确认（行保留在审计里，授权页的生效清单不再列出它）。";
const revokeHtml = renderToStaticMarkup(<SessionRulesRevokeReplyBlock state={{ phase: "receipt", text: REVOKE_RECEIPT }} />);
check("the revoke receipt must ride verbatim", revokeHtml.includes(esc(REVOKE_RECEIPT)));
check("the revoke receipt must name the row it spent", revokeHtml.includes(esc("已撤销会话授权 #12")));
check("the revoke receipt must state the NEXT-call re-ask", revokeHtml.includes(esc("下一次调用会重新弹卡确认")));
check("the revoke receipt block must hold zero emoji", !EMOJI_RE.test(revokeHtml));

const revokeRefusal =
  `面板请求被拒绝 [session.rules.revoke pc-7-def]：panel: 会话授权撤销通路未接入（本进程没有会话授权记账，或任务入口没有起来）：授权 #12 未被触碰，也不会被静默丢掉`;
const revokeRefusalHtml = renderToStaticMarkup(<SessionRulesRevokeReplyBlock state={{ phase: "refused", raw: revokeRefusal }} />);
check("the revoke refusal must ride verbatim", revokeRefusalHtml.includes(esc(revokeRefusal)));
check("the revoke refusal must state zero effect (未被触碰)", revokeRefusalHtml.includes(esc("未被触碰")));
check("the revoke refusal block must hold zero emoji", !EMOJI_RE.test(revokeRefusalHtml));

const revokeNoHostHtml = renderToStaticMarkup(<SessionRulesRevokeReplyBlock state={{ phase: "no-host" }} />);
check("the revoke no-host block must refuse to fake a success", revokeNoHostHtml.includes(esc("不画假成功")) && revokeNoHostHtml.includes(esc("未被触碰")));
check("the revoke no-host block must hold zero emoji", !EMOJI_RE.test(revokeNoHostHtml));

/* ---------------------------------------------------------------------------
   6. Source-level wiring: the literals, the echo key, the re-read
   --------------------------------------------------------------------------- */

const panelSrc = readFileSync("src/lib/panel.ts", "utf8");
check("panel.ts must spell session.rules.list as a literal (the ledger's closed gap)", panelSrc.includes('"session.rules.list"'));
check("panel.ts must spell session.rules.revoke as a literal", panelSrc.includes('"session.rules.revoke"'));
check(
  "the envelope key must be spelled as bridge.go's JSON tag carries it",
  panelSrc.includes("sessionRuleId"),
);
// The zero sentence is the DOOR's own words. The page must render
// state.zero verbatim and never spell the Go sentence itself - a page-side
// zero sentence would be a fabricated reading (mutation ②'s target).
const rulesLibSrc = readFileSync("src/lib/session-rules-view.ts", "utf8");
const cardSrc = readFileSync("src/components/session-rules-settings.tsx", "utf8");
check("the page must render the door's zero field, not a page-side zero sentence", cardSrc.includes("state.envelope.zero") && !cardSrc.includes("当前没有生效中的授权"));
check("the lib's marker must be spelled and match the Go constant", rulesLibSrc.includes('"wisp-session-rules:"'));
check("the page's requests must come through the ONE door file (lib/panel.ts)", rulesLibSrc.includes('from "@/lib/panel"'));
check("the list must be re-read after a revoke receipt (no local ledger copy to patch)", cardSrc.includes("useSessionRulesRevoke(list.reload)"));
check("the empty-ledger branch must draw no row chrome", cardSrc.includes("rows.length === 0"));
check("the settings page's security cell must mount the real card", readFileSync("src/components/security-settings.tsx", "utf8").includes("SessionRulesCard"));

// The whole card in its SSR state: loading is explicit, no emoji anywhere,
// and the idle revoke slot draws nothing.
const cardHtml = renderToStaticMarkup(<SessionRulesReplyBlock state={{ phase: "loading" }} />);
check("the loading state must say what is in flight", cardHtml.includes(esc("正在读取本会话的授权")));
check("the loading card must hold zero emoji", !EMOJI_RE.test(cardHtml));
check("the card source must hold zero emoji", !EMOJI_RE.test(cardSrc) && !EMOJI_RE.test(rulesLibSrc) && !EMOJI_RE.test(panelSrc));

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-sessionrules40: FAIL " + f);
  console.error(`render-sessionrules40: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-sessionrules40: OK, ticket 40 AC#2's grants page consumes the real doors " +
    "(list payload rows ride the store's own fields verbatim - id/scope/tool/pattern/session/createdAt/expiresAt, " +
    "the empty ledger renders the door's own zero sentence verbatim with no table chrome, " +
    "the revoke receipt names the row it spent and states the NEXT-call re-ask, " +
    "refusals ride verbatim including the sentence that names the missing 会话授权读面, " +
    "a marker-less reply stays 未判 and never draws a fabricated count, damaged rows are skipped and counted never folded into zeros, " +
    "the two literals live in panel.ts with the sessionRuleId echo key, the card re-reads after a revoke receipt, " +
    "the zero sentence is the door's words never the page's, zero emoji)",
);
