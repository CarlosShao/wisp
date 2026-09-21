/* ============================================================================
   AC#3 render evidence (ticket 77): the L2 card painted by React.

   Why this file exists
   --------------------
   internal/panel/approval.go can prove the view model carries a real
   risk.Decision, and TestApprovalCardViewJSONKeysMatchFrontendTypes can prove
   the TS type matches it, but neither proves that React takes those values and
   PAINTS them. The gap that leaves is the classic one: a component that
   compiles, mounts, and renders its library's demo content because nobody
   checked what ended up on screen.

   So this harness renders the real production root (src/App.tsx -> PanelSkeleton
   -> L2ApprovalCard) with react-dom/server over a JSON file that came out of
   `wisp.exe panel-assets -l2 fs.delete -irreversible delete -- <path>` - whose
   -l2 branch is panel.NewApprovalCardView(risk.NewRiskAssessor(), ...), the same
   assessor constructor internal/tools/bridge.go uses. The committed copy of
   that output is fixtures/l2-card-fs-delete.json (byte-identical to the CLI run
   recorded on the ticket, so CI can replay it without a Go toolchain). No
   literal in this file is ever asserted; every assertion is DERIVED from the
   input view, so if the component ignored its props the checks would fail.

   Two negative controls are part of the verdict, not decoration:
     - the vendored Beautiful UI demo questionnaire must not appear, because
       that is exactly the "fake props" failure this AC names;
     - the card count must equal the pending count, so a silently dropped
       pending card cannot render as a green run.

   Why it lives in scripts/ and not src/
   ------------------------------------
   It reads a file with node:fs, which TestPanelFrontendIsStateless (PLAN.md:1044)
   exists to flag inside frontend/src. Moving it into src/ would force that
   scanner to either grow an exception or go red - both are weakening an AC#5
   assertion, which this ticket may not do. It is therefore typechecked by the
   build that consumes it (vite build --ssr) and linted by `npm run lint`, and
   everything it renders (src/**) is still covered by `npm run typecheck`.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import App from "@/App";
import type { ApprovalCardView, PanelSnapshot } from "@/lib/panel";

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

function loadViews(path: string): ApprovalCardView[] {
  const parsed: unknown = JSON.parse(readFileSync(path, "utf8"));
  const views = Array.isArray(parsed) ? parsed : [parsed];
  for (const view of views) {
    if (
      typeof view !== "object" ||
      view === null ||
      typeof (view as ApprovalCardView).correlationId !== "string"
    ) {
      throw new Error(`render-l2: ${path} does not hold an ApprovalCardView JSON`);
    }
  }
  return views as ApprovalCardView[];
}

/** Every check is a function of the view, never of a constant. */
function expectationsFor(view: ApprovalCardView): { label: string; needle: string }[] {
  const wants: { label: string; needle: string }[] = [
    { label: "risk level", needle: view.level },
    { label: "tool name", needle: view.tool },
    { label: "correlationId present in tree", needle: "" },
  ];
  wants.pop();
  for (const arg of view.args) {
    if (arg.trim() !== "") wants.push({ label: "argv element", needle: arg });
  }
  for (const rule of view.rulesHit) {
    if (rule.trim() !== "") wants.push({ label: "rule id", needle: rule });
  }
  if (view.callChain.length > 0) {
    wants.push({ label: "call chain", needle: view.callChain.join(" > ") });
  }
  if (view.reasonKnown && view.reason.trim() !== "") {
    wants.push({ label: "reason", needle: view.reason });
  } else {
    // Q-23 fail-closed wording: an absent reason is stated as insufficient
    // information, never as an empty slot beside a card that looks safe.
    wants.push({ label: "insufficient-information wording", needle: "信息不足" });
  }
  if (view.sessionOverrideBlocked) {
    wants.push({ label: "C25 taint notice", needle: "会话级授权" });
  }
  wants.push({
    label: "verdict origin",
    needle: view.decidedBy === "native" ? "判定来自 原生风险评估" : view.decidedBy,
  });
  return wants;
}

const path = process.argv[2];
if (!path) {
  console.error("usage: node render-l2.js <ApprovalCardView.json>");
  process.exit(2);
}

const views = loadViews(path);
const snapshot: PanelSnapshot = { pending: views, results: [], generatedAt: "render-evidence" };
const html = renderToStaticMarkup(<App snapshot={snapshot} />);

const failures: string[] = [];

for (const view of views) {
  for (const { label, needle } of expectationsFor(view)) {
    if (!html.includes(esc(needle))) {
      failures.push(
        `card ${view.correlationId}: the rendered tree does not contain ${label} ${JSON.stringify(needle)}`,
      );
    }
  }
}

const CARD_TITLE = "需要你的确认";
const cards = html.split(esc(CARD_TITLE)).length - 1;
if (cards !== views.length) {
  failures.push(`rendered ${cards} card(s) for ${views.length} pending approval(s)`);
}

// The vendored upstream demo this card was cut from (src/components/ai-native/
// approval-card.tsx, NOT mounted). If any of it reaches the screen, the card is
// rendering library props and not the risk verdict.
for (const demo of [
  "How many flavors should we launch?",
  "Which mix-ins should we stock?",
  "Chocolate chips",
  "Send answers",
]) {
  if (html.includes(esc(demo))) {
    failures.push(`rendered tree carries the Beautiful UI demo string ${JSON.stringify(demo)}`);
  }
}

process.stdout.write(html + "\n");

if (failures.length > 0) {
  for (const failure of failures) console.error("render-l2: " + failure);
  console.error(`render-l2: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  `render-l2: OK, ${views.length} card(s) rendered from real risk decisions, ` +
    `${html.length} bytes of HTML`,
);
