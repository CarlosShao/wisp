/* ============================================================================
   Ticket 40 AC#1's behavioural tooth (render evidence), 2026-11-17 - the
   privacy page leg.
   ----------------------------------------------------------------------------
   WHAT THIS FILE PROVES, in the shape render-security-cost-40 established:
   every check below feeds REAL-SHAPED data through the PRODUCTION component
   and asserts the SAME VALUE comes back out in the markup - "喂真 payload ⇒
   三支逐字上屏".

   The fixtures are the wire's own documented shapes, copied from the Go leg
   that answers them (cmd/wisp/panel_privacy_40.go, the door
   bridge.go:141-143 names / composer_dispatch.go:174 routes /
   panel_inbound.go:411 wires):
     - delete-one receipt: HandlePrivacyDeleteOne's fmt.Sprintf sentence,
       naming the deleted id;
     - purge receipt: HandlePrivacyPurge's sentence, carrying the store's own
       PurgePrivacy count;
     - export receipt: HandlePrivacyExport's prose plus the ONE payload line
       behind PrivacyExportPayloadMarker ("wisp-privacy-export:"), whose JSON
       is the store's own compacted envelope {domain, exported_at, count,
       records};
     - refusals: RefusedEnvelopeForUser's wrapper sentence
       ("面板请求被拒绝 [method id]：panel: privacy.<verb> 被拒绝（...）：<cause>"),
       including the unknown-domain refusal that NAMES the D35 vocabulary -
       the sentence Test40PrivacyRefusalsNameTheDomainVocabulary pins on the
       Go side. The page's job is to show it verbatim, which is how the user
       "asks PrivacyDomains()": the vocabulary lives on the door, the refusal
       is the answer, and this tooth pins the words on screen.

   Where a reading cannot be had, the check pins the 未判 rule: no payload
   line ⇒ the page says 未判 and NEVER draws records=0; a payload without
   records ⇒ payload-error, named, never partially drawn.

   Source-level teeth:
     - the three route literals live in panel.ts (the ONE door file), with
       the envelope's privacyDomain/privacyItemId keys spelled as bridge.go's
       JSON tags carry them;
     - the page-gap ledger in internal/panel/bridge_test.go no longer
       registers the three names - the page started carrying them, so the
       stale entries had to leave in the same commit (the ledger's own
       both-ways rule);
     - the settings page's 隐私 cell draws the real control;
     - purge CANNOT fire without the confirm card (destruction asks first,
       deletion stays one click, and the page invents no "purge everything");
     - the page never copies the domain vocabulary into a picker.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import { PrivacyReplyBlock, PrivacySettings } from "@/components/privacy-settings";
import { parsePrivacyReply } from "@/lib/privacy-view";

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
   1. The export reply parses: prose unparsed, payload's own records counted
   --------------------------------------------------------------------------- */

// HandlePrivacyExport's exact prose shape, then the marker line whose JSON is
// the store's own envelope (domain / exported_at / count / records - count is
// the envelope's own key, deliberately NOT the page's reading: the page reads
// len(records), the same reading the prose states).
const EXPORT_PROSE =
  "已导出 memory 域：2 条记录，完整行内数据在下面这一行 JSON（导出即用户自己的数据，不是摘要）。";
const EXPORT_PAYLOAD =
  `{"domain":"memory","exported_at":"2026-11-17T03:12:00Z","count":2,` +
  `"records":[{"id":11,"content":"第一条"},{"id":12,"content":"第二条"}]}`;
const EXPORT_REPLY = `${EXPORT_PROSE}\nwisp-privacy-export:${EXPORT_PAYLOAD}`;

const exportParsed = parsePrivacyReply(EXPORT_REPLY, "export");
check("the export reply must parse as the export state", exportParsed.phase === "export");
if (exportParsed.phase === "export") {
  check("the export prose must arrive verbatim, unparsed", exportParsed.prose === EXPORT_PROSE, exportParsed.prose);
  check("records=N must be len(records) read off the payload", exportParsed.envelope.recordsCount === 2);
  check("the payload's domain must ride verbatim", exportParsed.envelope.domain === "memory");
  check("the payload's exported_at must ride verbatim", exportParsed.envelope.exportedAt === "2026-11-17T03:12:00Z");
  check("the rows must be carried verbatim (the user's own data)", exportParsed.envelope.records.length === 2);
}

/* ---------------------------------------------------------------------------
   2. The three receipts render verbatim (delete names the id, purge names
      the store's count, export states records=N off the payload)
   --------------------------------------------------------------------------- */

const DELETE_RECEIPT = "已删除 memory 域的 1 条记录（id=42）。该域其余记录未动。";
const deleteHtml = renderToStaticMarkup(
  <PrivacyReplyBlock state={{ phase: "receipt", text: DELETE_RECEIPT }} />,
);
check("the delete receipt must ride verbatim", deleteHtml.includes(esc(DELETE_RECEIPT)));
check("the delete receipt must name the id on screen", deleteHtml.includes(esc("id=42")));
check("the delete reply block must hold zero emoji", !EMOJI_RE.test(deleteHtml));

const PURGE_RECEIPT = "已清空 memory 域，共删除 3 条记录。其他域未动。";
const purgeHtml = renderToStaticMarkup(
  <PrivacyReplyBlock state={{ phase: "receipt", text: PURGE_RECEIPT }} />,
);
check("the purge receipt must ride verbatim", purgeHtml.includes(esc(PURGE_RECEIPT)));
check("the purge receipt must carry the store's own count on screen", purgeHtml.includes(esc("共删除 3 条记录")));
check("the purge reply block must hold zero emoji", !EMOJI_RE.test(purgeHtml));

const exportHtml = renderToStaticMarkup(<PrivacyReplyBlock state={exportParsed} />);
check(
  "the export block must state records=N off the payload, verbatim",
  exportHtml.includes(`records=${exportParsed.phase === "export" ? exportParsed.envelope.recordsCount : "?"}`),
);
check("the export prose must render verbatim", exportHtml.includes(esc(EXPORT_PROSE)));
check("the exported rows must stay on screen verbatim (the user's data, not a summary)", exportHtml.includes(esc("第一条")) && exportHtml.includes(esc("第二条")));
check("the export block must hold zero emoji", !EMOJI_RE.test(exportHtml));

/* ---------------------------------------------------------------------------
   3. Refusals ride verbatim - including the refusal that NAMES the D35
      vocabulary (the page's way of asking PrivacyDomains())
   --------------------------------------------------------------------------- */

// RefusedEnvelopeForUser's wrapper around panel_privacy_40.go's
// unknown-domain refusal: the cause sentence lists the five domains the
// memory layer owns. Fed to the page, the words must come back out.
const PURGE_REFUSAL =
  `面板请求被拒绝 [privacy.purge pc-9-abc]：panel: privacy.purge 被拒绝（requestId="pc-9-abc"）：` +
  `privacyDomain "sessions" 不是 D35 隐私域；purge 只按域执行，可用的只有 [profile memory task_log tool_call artifacts]`;
const refusalParsed = parsePrivacyReply(PURGE_REFUSAL, "plain");
check("a refusal must classify as the refused state", refusalParsed.phase === "refused");
const refusalHtml = renderToStaticMarkup(<PrivacyReplyBlock state={refusalParsed} />);
check("the refusal sentence must ride verbatim", refusalHtml.includes(esc(PURGE_REFUSAL)));
check("the refusal must name the D35 vocabulary on screen (task_log, artifacts)", refusalHtml.includes(esc("task_log")) && refusalHtml.includes(esc("artifacts")));
check("the refusal must render in the danger tint", refusalHtml.includes("bg-red-tint"));
check("the refusal block must hold zero emoji", !EMOJI_RE.test(refusalHtml));

// The empty-id refusal the door names (Test40PrivacyRefusalsNameTheDomainVocabulary
// pins the cause on the Go side; the page shows it whole).
const NOID_REFUSAL =
  `面板请求被拒绝 [privacy.delete-one pc-2-noid]：panel: privacy.delete-one 被拒绝（requestId="pc-2-noid"）：` +
  `privacyItemId 是空的：delete-one 只删一条，不给它 id 就什么都不删`;
const noidParsed = parsePrivacyReply(NOID_REFUSAL, "plain");
check("the empty-id refusal must classify as refused", noidParsed.phase === "refused");
const noidHtml = renderToStaticMarkup(<PrivacyReplyBlock state={noidParsed} />);
check("the empty-id refusal must ride verbatim", noidHtml.includes(esc(NOID_REFUSAL)));

// An export reply WITHOUT its payload line: the reading stays 未判 and no
// records=0 may be drawn (喂不到 ⇒ 显式未判, 禁止假 0).
const noPayloadParsed = parsePrivacyReply("已导出 memory 域：2 条记录。", "export");
check("a marker-less export reply must parse as no-payload", noPayloadParsed.phase === "no-payload");
const noPayloadHtml = renderToStaticMarkup(<PrivacyReplyBlock state={noPayloadParsed} />);
check("the no-payload block must say 未判 explicitly", noPayloadHtml.includes(esc("未判")));
check("the no-payload block must NEVER draw a fabricated records=0", !noPayloadHtml.includes("records=0"));
check("the no-payload block must carry the raw reply verbatim", noPayloadHtml.includes(esc("已导出 memory 域：2 条记录。")));
check("the no-payload block must hold zero emoji", !EMOJI_RE.test(noPayloadHtml));

// A payload line that claims records it cannot hand over: payload-error,
// named, never partially drawn.
const payloadErrParsed = parsePrivacyReply(`x\nwisp-privacy-export:{"domain":"memory","count":2}`, "export");
check("a payload without records must refuse as payload-error", payloadErrParsed.phase === "payload-error");
if (payloadErrParsed.phase === "payload-error") {
  const payloadErrHtml = renderToStaticMarkup(<PrivacyReplyBlock state={payloadErrParsed} />);
  check("the payload-error message must ride verbatim", payloadErrHtml.includes(esc(payloadErrParsed.message)));
  check("the payload-error block must not draw a fabricated records count (the message may name the rule, never a number)", !payloadErrHtml.includes("records=0") && !/>records=\d+</.test(payloadErrHtml));
  check("the payload-error block must hold zero emoji", !EMOJI_RE.test(payloadErrHtml));
}

// The no-host state (plain browser / hand-assembled harness): explicit, and
// nothing that reads like a result.
const noHostHtml = renderToStaticMarkup(<PrivacyReplyBlock state={{ phase: "no-host" }} />);
check("the no-host block must say the reading is 未判", noHostHtml.includes(esc("未判")));
check("the no-host block must hold zero emoji", !EMOJI_RE.test(noHostHtml));

/* ---------------------------------------------------------------------------
   4. Source-level wiring: the literals, the ledger, the confirm card
   --------------------------------------------------------------------------- */

const panelSrc = readFileSync("src/lib/panel.ts", "utf8");
check("panel.ts must spell privacy.delete-one as a literal (the ledger's closed gap)", panelSrc.includes('"privacy.delete-one"'));
check("panel.ts must spell privacy.purge as a literal", panelSrc.includes('"privacy.purge"'));
check("panel.ts must spell privacy.export as a literal", panelSrc.includes('"privacy.export"'));
check(
  "the envelope keys must be spelled as bridge.go's JSON tags carry them",
  panelSrc.includes("privacyDomain") && panelSrc.includes("privacyItemId"),
);

// The page-gap ledger's own both-ways rule: the page started carrying the
// three literals, so the stale registrations had to leave in the same commit.
const ledgerSrc = readFileSync("../internal/panel/bridge_test.go", "utf8");
check(
  "the page-gap ledger must no longer register the three privacy names (the page literal landed)",
  !ledgerSrc.includes('"privacy.delete-one":') &&
    !ledgerSrc.includes('"privacy.purge":') &&
    !ledgerSrc.includes('"privacy.export":'),
);

const settingsSrc = readFileSync("src/components/harness/settings-page.tsx", "utf8");
check("the settings page must draw 隐私 with the privacy control", settingsSrc.includes('control: "privacy"'));

const privacySrc = readFileSync("src/components/privacy-settings.tsx", "utf8");
const privacyLibSrc = readFileSync("src/lib/privacy-view.ts", "utf8");
check("the marker must be spelled in the lib and match the Go constant", privacyLibSrc.includes('"wisp-privacy-export:"'));
check("the page's requests must come through the ONE door file (lib/panel.ts)", privacyLibSrc.includes('from "@/lib/panel"'));
check("each branch must use its own request function (no shared runner that could cross the verbs)", privacyLibSrc.includes("requestPrivacyDeleteOne") && privacyLibSrc.includes("requestPrivacyPurge") && privacyLibSrc.includes("requestPrivacyExport"));

// DESTRUCTION ASKS FIRST. The purge branch's first button only opens the
// confirm card; the request is sent from exactly one place - the card's
// confirm button - so a purge can never fire on the first click.
check("the purge first button must only open the confirm card", privacySrc.includes("setPurgeConfirmOpen(true)"));
check(
  "privacy.purge must be sent from exactly one place (the confirm card's confirm button)",
  (privacySrc.match(/purge\.submit\(/g) ?? []).length === 1,
);
check("the confirm card must be a real dialog surface", privacySrc.includes('role="alertdialog"'));
check("the confirm card must name the irreversibility", privacySrc.includes("不可逆"));
check("the confirm card must carry the 取消 exit", privacySrc.includes("取消"));

// DELETION STAYS ONE CLICK BUT REQUIRES THE ID: the dispatch's 需选 id rule,
// enforced where the button is drawn.
check("the delete button must be disabled until an id is filled", privacySrc.includes('disabled={deleteItemId.trim() === ""}'));
check(
  "the page must say why the id is typed (no privacy list read face - no invented row picker)",
  privacySrc.includes("没有隐私条目的清单读面"),
);

// THE VOCABULARY IS ASKED, NEVER COPIED: the page spells no domain word as a
// picker option - the door validates against PrivacyDomains() and its
// refusal (checked verbatim above) is the vocabulary the user sees.
check(
  "the page must not copy the D35 vocabulary into a local word list",
  !privacySrc.includes('"profile"') && !privacySrc.includes('"task_log"') && !privacySrc.includes('"artifacts"'),
);
check("the page must name the vocabulary's owner (PrivacyDomains)", privacySrc.includes("PrivacyDomains()"));

// The whole page in its SSR state: no emoji anywhere, and the three branch
// captions present.
const pageHtml = renderToStaticMarkup(<PrivacySettings />);
check("the privacy page must hold zero emoji", !EMOJI_RE.test(pageHtml) && !EMOJI_RE.test(privacySrc));
check("the page must draw the three branches", pageHtml.includes(esc("删除一条")) && pageHtml.includes(esc("清空某域")) && pageHtml.includes(esc("导出某域")));
check("the page must carry the ask-don't-copy caption", pageHtml.includes(esc("PrivacyDomains()")));
check(
  "the SSR page must not carry any receipt before a request is sent (idle draws nothing)",
  !pageHtml.includes(esc("回执逐字")),
);

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-privacy40: FAIL " + f);
  console.error(`render-privacy40: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-privacy40: OK, ticket 40 AC#1's privacy page consumes the real doors " +
    "(delete-one receipt names the id verbatim, purge receipt carries the store's own count verbatim, " +
    "export prose verbatim with records=N read off the payload and the user's own rows on screen, " +
    "refusals ride verbatim including the sentence that names the D35 vocabulary, " +
    "a marker-less export stays 未判 and never draws records=0, " +
    "payload without records is a named payload-error, " +
    "the three literals live in panel.ts with the ledger entries gone in the same commit, " +
    "purge fires only from the confirm card, delete requires the id, the domain vocabulary is asked never copied, zero emoji)",
);
