/* ============================================================================
   Ticket 38 AC#3's page-half tooth (render evidence), 2026-11-17 - the tool
   discovery leg (consume wave 3).
   ----------------------------------------------------------------------------
   WHAT THIS FILE PROVES, in the shape render-history-page-36 established:
   every check below feeds REAL-SHAPED data through the PRODUCTION code and
   asserts the SAME VALUE comes back out - no browser needed, because every
   piece under test is a pure function or a pure render over the wire's own
   shapes.

   THE FIXTURES are the Go leg's own documented shapes, copied from the Go
   half that answers them (cmd/wisp/panel_tools_discovery.go
   HandleToolsList: the prose sentence plus one "wisp-tools-list:" marker
   line carrying {"rows":[ToolRosterRow]}, ToolRosterRow being
   internal/panel/discovery_force_w3.go's shape with riskPill rendered by
   internal/risk's Level.String()).

   THE TEETH (the dispatch's three, plus the door facts that make the page
   half honest):
     1. 喂真 payload => 每行风险丸逐字: the pill word on screen is EXACTLY
        the wire's riskPill - no "轻确认"/"强确认" decoration, no page-side
        level mapping;
     2. 3 枚不同级 => 三色三字: L0/L1/L2 rows render three DISTINCT tones
        and the three words verbatim;
     3. 缺字段 => 未判: a payload without riskPill renders 未判 with the
        unjudged tone, never a level the page guessed;
     4. refusals ride verbatim, a receipt without its payload line stays
        未判-shaped, and no-host says so by name;
     5. the launch key is the REAL door: one button per row, the pill never
        gates it, the page's only launch wire call is lib/panel's
        requestToolLaunch, and no allow/decide vocabulary exists anywhere on
        the page half;
     6. source level: the three wave-3 literals live in panel.ts, the
        page-gap ledger's three registrations are GONE (the both-ways rule),
        and zero emoji everywhere.
   ============================================================================ */

import { readFileSync } from "node:fs";
import { renderToStaticMarkup } from "react-dom/server";
import { PaletteScreen, type PaletteToolRow } from "@/components/palette-screen";
import { parseToolsListReply, riskPillTone } from "@/lib/tools-discovery-view";

const failures: string[] = [];
function check(label: string, ok: boolean, detail = ""): void {
  if (!ok) failures.push(`${label}${detail ? ` [${detail}]` : ""}`);
}

// The emoji band the repo scanner walks (ban #8, tools/d22scan's own set).
const EMOJI_RE =
  /[\u{1F000}-\u{1FAFF}\u{2200}-\u22FF\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{1F1E6}-\u{1F1FF}]|\uFE0F/u;

/* ---------------------------------------------------------------------------
   Fixtures: the Go reply, byte-shaped. The rows are the registry family the
   D34 table declares at three levels (L0 read, L1 write, L2 delete) plus the
   Deny word - each pill is the exact string internal/risk's Level.String()
   renders, which is what the Go handler puts on the wire.
   --------------------------------------------------------------------------- */

const TOOLS_PROSE = "已安装工具共 4 只（riskPill 为 C19 风险级真值，只读名册；启动走 tool.launch，照样过审批门）。";

const TOOLS_ROWS = [
  { name: "fs.read", title: "读取工作区内的一个文件", riskPill: "L0", argSummary: "path" },
  { name: "fs.write", title: "写入或创建工作区内的一个文件", riskPill: "L1", argSummary: "path,content" },
  { name: "fs.delete", title: "删除工作区内的一个文件", riskPill: "L2", argSummary: "path" },
  { name: "net.fetch", title: "按网络白名单取一个 URL", riskPill: "Deny", argSummary: "url" },
] satisfies PaletteToolRow[];

function toolsReply(rows: unknown[]): string {
  const payload = JSON.stringify({ rows });
  return `${TOOLS_PROSE}\nwisp-tools-list:${payload}`;
}

/* ---------------------------------------------------------------------------
   1-2. 喂真 payload => 逐字风险丸; 3 枚不同级 => 三色三字
   --------------------------------------------------------------------------- */

const fed = parseToolsListReply(toolsReply(TOOLS_ROWS));
check("the real reply must parse as a roster", fed.phase === "roster");
if (fed.phase === "roster") {
  check("the receipt prose must ride verbatim, unparsed", fed.prose === TOOLS_PROSE);
  check("the roster must keep all four rows", fed.rows.length === 4);
  for (const row of fed.rows) {
    const wire = TOOLS_ROWS.find((r) => r.name === row.name);
    check(`row ${row.name}: the pill must be the wire's own word, verbatim`, row.riskPill === wire?.riskPill);
  }
}

const fedHtml = renderToStaticMarkup(
  <PaletteScreen
    groups={[]}
    notes={[]}
    onQueryChange={() => undefined}
    query=""
    toolsPhase="roster"
    toolsProse={TOOLS_PROSE}
    toolsRows={TOOLS_ROWS}
  />,
);
check("the L0 row's pill word must reach the screen verbatim", fedHtml.includes(">L0</span>"));
check("the L1 row's pill word must reach the screen verbatim", fedHtml.includes(">L1</span>"));
check("the L2 row's pill word must reach the screen verbatim", fedHtml.includes(">L2</span>"));
check("the Deny row's pill word must reach the screen verbatim", fedHtml.includes(">Deny</span>"));
// 三色三字: three distinct tones for the three levels (plus danger for Deny),
// and never the approval card's level labels - the page maps WORDS to hues,
// never a level to new words.
check("L0 must take the ok tone", fedHtml.includes('data-pill-tone="ok"'));
check("L1 must take the warn tone", fedHtml.includes('data-pill-tone="warn"'));
check("L2 must take the danger tone", fedHtml.includes('data-pill-tone="danger"'));
check(
  "the page must not decorate the pill with a level label of its own",
  !fedHtml.includes("轻确认") && !fedHtml.includes("强确认"),
);
check("the tone derivation must be a pure word map", riskPillTone("L0") === "ok" && riskPillTone("L1") === "warn" && riskPillTone("L2") === "danger" && riskPillTone("Deny") === "danger");
check("an unknown pill word earns no colour but stays verbatim", riskPillTone("Level(9)") === "neutral");

// The row facts round-trip: name and title verbatim, argSummary carried.
check("the row name must reach the screen", fedHtml.includes("fs.read"));
check("the row title must reach the screen", fedHtml.includes("读取工作区内的一个文件"));
check("the arg summary must reach the screen", fedHtml.includes("path,content"));

/* ---------------------------------------------------------------------------
   3. 缺字段 => 未判, never a fake level
   --------------------------------------------------------------------------- */

const toothless = parseToolsListReply(toolsReply([{ name: "fs.read", title: "读取工作区内的一个文件" }]));
check("a pill-less row must still parse as a roster", toothless.phase === "roster");
if (toothless.phase === "roster") {
  check("the missing pill must stay absent, not defaulted", toothless.rows[0]?.riskPill === undefined);
}
const toothlessHtml = renderToStaticMarkup(
  <PaletteScreen
    groups={[]}
    notes={[]}
    onQueryChange={() => undefined}
    query=""
    toolsPhase="roster"
    toolsRows={[{ name: "fs.read", title: "读取工作区内的一个文件" }]}
  />,
);
check("a missing pill must render 未判", toothlessHtml.includes("未判"));
check("the missing pill must carry the unjudged tone", toothlessHtml.includes('data-pill-tone="unjudged"'));
check(
  "a missing pill must never render as a level word",
  !/>L0<|>L1<|>L2<|>Deny</.test(toothlessHtml),
);

// A payload defect is a named refusal, never a partially drawn list.
const broken = parseToolsListReply(toolsReply(["not-an-object"]));
check("a non-object row must refuse the payload whole", broken.phase === "payload-error");
const noRows = parseToolsListReply(`${TOOLS_PROSE}\nwisp-tools-list:{"total":4}`);
check("a payload without rows must refuse by name", noRows.phase === "payload-error");
const noMarker = parseToolsListReply(TOOLS_PROSE);
check("a receipt without its payload line must stay 未判-shaped", noMarker.phase === "no-payload");

/* ---------------------------------------------------------------------------
   4. Refusals ride verbatim; no-host says so by name
   --------------------------------------------------------------------------- */

const refusedHtml = renderToStaticMarkup(
  <PaletteScreen
    groups={[]}
    notes={[]}
    onQueryChange={() => undefined}
    query=""
    toolsNote="面板请求被拒绝 [tools.list：名册未接入（本进程没有任务管线）]（requestId=pc-1-x）"
    toolsPhase="refused"
  />,
);
check(
  "the door's refusal must ride verbatim",
  refusedHtml.includes("面板请求被拒绝 [tools.list：名册未接入（本进程没有任务管线）]") && refusedHtml.includes("data-tool-refusal"),
);
const nohostHtml = renderToStaticMarkup(
  <PaletteScreen groups={[]} notes={[]} onQueryChange={() => undefined} query="" toolsPhase="nohost" />,
);
check("no host must say so by name, never as an empty roster", nohostHtml.includes("没接到宿主"));

/* ---------------------------------------------------------------------------
   5. The launch key is the REAL door
   --------------------------------------------------------------------------- */

// One launch button per row, rendered for EVERY row - the pill never gates it
// (an L2 row and an L0 row carry exactly the same key).
const launchCount = (fedHtml.match(/data-launch-tool/g) ?? []).length;
check("every roster row must carry exactly one launch key", launchCount === TOOLS_ROWS.length);
check(
  "the L2 row's launch key must not be gated by its pill",
  fedHtml.includes("fs.delete") && fedHtml.includes("启动工具 fs.delete"),
);

const paletteSrc = readFileSync("src/components/palette-screen.tsx", "utf8");
const appSrc = readFileSync("src/App.tsx", "utf8");
const panelSrc = readFileSync("src/lib/panel.ts", "utf8");
const viewSrc = readFileSync("src/lib/tools-discovery-view.ts", "utf8");

check(
  "App must bind the launch key to lib/panel's ONE wire function",
  appSrc.includes("onLaunchTool={launchTool}") && appSrc.includes("requestToolLaunch(row.name)"),
);
check("the launch wire function must spell the tool.launch literal", panelSrc.includes(`method: "tool.launch"`));
check(
  "the launch receipt/refusal must land on the palette's own line, verbatim",
  appSrc.includes("paletteLaunchNote"),
);
check(
  "the page half must carry no allow/decide vocabulary",
  !/approval\.decide|allowAll|"allow"/.test(paletteSrc) &&
    !/approval\.decide|allowAll|"allow"/.test(appSrc),
);
check(
  "the page must not predict a verdict: no pill-gated launch condition",
  !paletteSrc.includes("riskPill ===") && !paletteSrc.includes("riskPill !=="),
);
check(
  "the pill's only page-side derivation is the word->hue map",
  viewSrc.includes('if (pill === "L0")') &&
    viewSrc.includes('if (pill === "L1")') &&
    viewSrc.includes('pill === "L2"'),
);

/* ---------------------------------------------------------------------------
   6. Source level: the three literals, the ledger's departure, zero emoji
   --------------------------------------------------------------------------- */

for (const literal of ['"tools.list"', '"tool.launch"', '"rules.force"']) {
  check(`panel.ts must spell the ${literal} literal (the page-gap half)`, panelSrc.includes(literal));
}
const bridgeTestSrc = readFileSync("../internal/panel/bridge_test.go", "utf8");
check(
  "the page-gap ledger's three wave-3 entries must be GONE (the both-ways rule)",
  bridgeTestSrc.includes("var rosterPageGapLedger194 = map[string]string{}") &&
    !bridgeTestSrc.includes('"tools.list":') &&
    !bridgeTestSrc.includes('"tool.launch":') &&
    !bridgeTestSrc.includes('"rules.force":'),
);

for (const [name, src] of [
  ["palette-screen", paletteSrc],
  ["App", appSrc],
  ["src/lib/panel.ts", panelSrc],
  ["tools-discovery-view", viewSrc],
  ["this instrument", readFileSync("scripts/render-discovery38.tsx", "utf8")],
] as const) {
  check(`${name} must hold zero emoji`, !EMOJI_RE.test(src));
}
check("the rendered palette markup must hold zero emoji", !EMOJI_RE.test(fedHtml));

/* ---------------------------------------------------------------------------
   The verdict.
   --------------------------------------------------------------------------- */

if (failures.length > 0) {
  for (const f of failures) console.error("render-discovery38: FAIL " + f);
  console.error(`render-discovery38: FAIL, ${failures.length} check(s) failed`);
  process.exit(1);
}
console.error(
  "render-discovery38: OK, ticket 38 AC#3's discovery face consumes tools.list " +
    "(pills verbatim off the wire, L0/L1/L2 three tones three words, missing pill 未判, " +
    "refusals verbatim, one ungated launch key per row through tool.launch, " +
    "the three page literals landed and the ledger entries left)",
);
