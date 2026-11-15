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

// ---------------------------------------------------------------------------
// The four shapes PLAN.md:3481 names for this card, added 2026-09-25.
//
// Why they live in THIS harness and not a new one: `render:l2` is already CI
// step 6, so a check put here gates on the next push. A new `npm run` script
// would need a new step in .github/workflows/ci.yml, which is not frontend's to
// edit. The step is named "renders the real risk fields"; these checks are about
// the same card's shape, so the name now reads narrower than the content - that
// drift is reported rather than hidden.
//
// Where a number is expected it is parsed out of PLAN.md at run time (same rule
// render-stream.tsx follows), so the spec can be reworded and this file cannot
// quietly keep agreeing with an old reading of it.
// ---------------------------------------------------------------------------

function planCell(rowHead: string, cell: number, label: string): string {
  let plan: string;
  try {
    plan = readFileSync("../docs/PLAN.md", "utf8");
  } catch {
    throw new Error("render-l2: cannot read ../docs/PLAN.md, which these checks derive from");
  }
  const row = plan.split("\n").find((l) => l.startsWith("|") && l.includes(rowHead));
  if (!row) throw new Error(`render-l2: no '${rowHead}' row in PLAN.md - ${label} has no source`);
  const cells = row.split("|").map((c) => c.trim());
  const out = cells[cell];
  if (!out) throw new Error(`render-l2: '${rowHead}' row has no cell ${cell} (${label})`);
  return out;
}

const themeCss = readFileSync("src/styles/theme.css", "utf8");
function cssRule(selector: string): string {
  const m = themeCss.match(new RegExp(`^${selector.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}\\s*\\{([^}]*)\\}`, "m"));
  return m?.[1] ?? "";
}

// 顶部 2px --danger 横条
const barPx = planCell("L2 确认卡（面板内）", 2, "the danger bar width").match(/顶部\s*(\d+)px/)?.[1];
if (!barPx) throw new Error("render-l2: PLAN.md's L2 row no longer states the top bar in Npx");
const barClass = `h-[${barPx}px]`;
const barCount = html.split(barClass).length - 1;
if (barCount !== views.length) {
  failures.push(`${barCount} top bar(s) of ${barClass} for ${views.length} card(s)`);
}
if (!/[,\s]bg-stop(?![\w-])/.test(html)) {
  failures.push("the top bar is not painted with bg-stop");
}
if (!/--stop:\s*var\(--danger\)/.test(cssRule(":root") + themeCss)) {
  failures.push("bg-stop does not resolve to the C21 danger token --stop: var(--danger)");
}

// 悬浮球图示的环 2s 呼吸
const ringSec = planCell("L2 确认卡（面板内）", 3, "the ring period").match(/(\d+(?:\.\d+)?)s/)?.[1];
if (!ringSec) throw new Error("render-l2: PLAN.md's L2 row no longer states the ring period in Ns");
const ringRule = cssRule(".l2-ball-ring");
if (!ringRule.includes(` ${ringSec}s`)) {
  failures.push(`.l2-ball-ring animation is not ${ringSec}s: '${ringRule.trim() || "no rule"}'`);
}
if (!ringRule.includes("infinite")) {
  failures.push(".l2-ball-ring must breathe on a loop while the card is up");
}
if (html.split("l2-ball-ring").length - 1 !== views.length) {
  failures.push(`the ball glyph is not drawn once per card (${views.length} expected)`);
}

// 底部常驻一行「点击悬浮球以批准」: PLAN names the words, so the words are the needle.
// Anchored on 常驻一行提示 because the same cell quotes a SECOND 「...」 first (the C19
// rule example) - a bare 「([^」]+)」 grabbed that one and the check then demanded the
// rule text on the card, which is how this line got its anchor.
const hintText = planCell("L2 确认卡（面板内）", 2, "the standing hint")
  .match(/常驻一行提示[「]([^」]+)[」]/)?.[1];
if (!hintText) throw new Error("render-l2: PLAN.md's L2 row no longer quotes the standing hint");
if (html.split(esc(hintText)).length - 1 !== views.length) {
  failures.push(`${JSON.stringify(hintText)} is not present once per card`);
}

// 按钮区只有 拒绝 与 查看完整参数 —— and, in the row's own 明确不用 column, no allow.
for (const label of ["拒绝", "查看完整参数"]) {
  if (html.split(esc(label)).length - 1 !== views.length) {
    failures.push(`button ${JSON.stringify(label)} is not present once per card`);
  }
}

/** A card's whole job is to withhold this one affordance, so the control is
    spelled out rather than inferred from "nobody added one". Every labelled
    button is collected - visible text plus aria-label, because a name can ride
    either - and none may read as an approval. */
const buttonNames: string[] = [];
for (const m of html.matchAll(/<button\b([^>]*)>([\s\S]*?)<\/button>/g)) {
  const aria = m[1]?.match(/aria-label="([^"]*)"/)?.[1] ?? "";
  const text = (m[2] ?? "").replace(/<[^>]*>/g, "").replace(/&quot;/g, '"').trim();
  buttonNames.push(`${text} ${aria}`.trim());
}
for (const name of buttonNames) {
  if (/允许|同意|批准|通过|approve|allow|grant/i.test(name)) {
    failures.push(`the panel card offers an approval-shaped button ${JSON.stringify(name)}`);
  }
}
if (readFileSync("src/components/l2-approval-card.tsx", "utf8").includes('send("grant"')) {
  failures.push("l2-approval-card.tsx still sends a grant");
}

// The badge is an outline capsule, never a filled chip (PLAN.md:3478's 风险级徽标
// rule, applied to the same badge this card draws in its title).
for (const fill of ["bg-caution/15", "bg-stop/15"]) {
  if (html.includes(fill)) {
    failures.push(`the risk badge is filled with ${fill}; 1px 描边胶囊 forbids a fill`);
  }
}
const titles = [...html.matchAll(/data-slot="card-title"[^>]*>([\s\S]*?)<\/div>/g)].map((m) => m[1] ?? "");
if (titles.length !== views.length) {
  failures.push(`${titles.length} card title(s) for ${views.length} card(s) - the title element is not one per card`);
}
for (const [i, view] of views.entries()) {
  const title = titles[i] ?? "";
  // 标题 = 工具名 + L2 徽标: both must sit INSIDE the title element, which is what
  // distinguishes this check from the older ones that only proved the two
  // strings appear somewhere on the card.
  if (!title.includes(esc(view.tool))) {
    failures.push(`card ${view.correlationId}: title does not carry the tool name`);
  }
  if (!title.includes("data-slot=\"badge\"") || !title.includes(esc(view.level))) {
    failures.push(`card ${view.correlationId}: title does not carry the ${view.level} badge`);
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
